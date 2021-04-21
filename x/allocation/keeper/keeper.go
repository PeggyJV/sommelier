package keeper

import (
	"bytes"
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/x/allocation/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the oracle store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryMarshaler
	paramSpace    paramtypes.Subspace
	stakingKeeper types.StakingKeeper

	oracleHandler types.OracleHandler
	handlerSet    bool
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper types.StakingKeeper,
) Keeper {

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSpace:    paramSpace,
		stakingKeeper: stakingKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// SetHandler sets the oracle handler into the keeper's fields. It will panic
// if the handler is already set.
func (k *Keeper) SetHandler(handlerFn types.OracleHandler) {
	if k.handlerSet {
		panic("oracle handler is already set")
	}

	k.oracleHandler = handlerFn
	k.handlerSet = true
}

////////////////////////
// MsgDelegateAddress //
////////////////////////

// SetValidatorDelegateAddress sets a new address that will have the power to send data on behalf of the validator
func (k Keeper) SetValidatorDelegateAddress(ctx sdk.Context, del sdk.AccAddress, val sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetAllocationDelegateKey(del), val.Bytes())
}

// GetValidatorAddressFromDelegate returns the delegate address for a given validator
func (k Keeper) GetValidatorAddressFromDelegate(ctx sdk.Context, del sdk.AccAddress) sdk.ValAddress {
	return ctx.KVStore(k.storeKey).Get(types.GetAllocationDelegateKey(del))
}

// GetDelegateAddressFromValidator returns the validator address for a given delegate
func (k Keeper) GetDelegateAddressFromValidator(ctx sdk.Context, val sdk.ValAddress) sdk.AccAddress {
	var address sdk.AccAddress
	// TODO: create secondary index
	k.IterateDelegateAddresses(ctx, func(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) bool {
		if !val.Equals(validatorAddr) {
			return false
		}

		address = delegatorAddr
		return true
	})
	return address
}

// IsDelegateAddress returns true if the validator has delegated their feed to an address
func (k Keeper) IsDelegateAddress(ctx sdk.Context, del sdk.AccAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetAllocationDelegateKey(del))
}

// IterateDelegateAddresses iterates over all delegate address pairs in the store
func (k Keeper) IterateDelegateAddresses(ctx sdk.Context, handler func(del sdk.AccAddress, val sdk.ValAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AllocationDelegateKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		del := sdk.AccAddress(bytes.TrimPrefix(iter.Key(), types.AllocationDelegateKeyPrefix))
		val := sdk.ValAddress(iter.Value())
		if handler(del, val) {
			break
		}
	}
}

// GetAllAllocationDelegations returns all the delegations for allocations
func (k Keeper) GetAllAllocationDelegations(ctx sdk.Context) []types.MsgDelegateAllocations {
	allocationDelegations := make([]types.MsgDelegateAllocations, 0)

	k.IterateDelegateAddresses(ctx, func(del sdk.AccAddress, val sdk.ValAddress) bool {
		allocationDelegations = append(allocationDelegations, types.MsgDelegateAllocations{
			Delegate:  del.String(),
			Validator: val.String(),
		})
		return false
	})

	return allocationDelegations
}

//////////////////////////
// Allocation Precommit //
//////////////////////////

// SetAllocationPrecommit sets the precommit for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) SetAllocationPrecommit(ctx sdk.Context, validatorAddr sdk.ValAddress, cellarAddr common.Address, precommit types.AllocationPrecommit) {
	bz := k.cdc.MustMarshalBinaryBare(&precommit)
	ctx.KVStore(k.storeKey).Set(types.GetAllocationPrecommitKey(validatorAddr, cellarAddr), bz)
}

// GetAllocationPrecommit gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetAllocationPrecommit(ctx sdk.Context, val sdk.ValAddress, cel common.Address) (types.AllocationPrecommit, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetAllocationPrecommitKey(val, cel))
	if len(bz) == 0 {
		return types.AllocationPrecommit{}, false
	}

	var precommit types.AllocationPrecommit
	k.cdc.MustUnmarshalBinaryBare(bz, &precommit)
	return precommit, true
}

// DeleteAllPrecommits removes all the prevotes for the current block iteration
func (k Keeper) DeleteAllPrecommits(ctx sdk.Context) {
	k.IterateAllocationPrecommits(ctx, func(val sdk.ValAddress, cel common.Address, _ types.AllocationPrecommit) bool {
		k.DeleteAllocationPrecommit(ctx, val, cel)
		return false
	})
}

// DeleteAllocationPrecommit deletes the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteAllocationPrecommit(ctx sdk.Context, val sdk.ValAddress, cel common.Address) {
	ctx.KVStore(k.storeKey).Delete(types.GetAllocationPrecommitKey(val, cel))
}

// HasAllocationPrecommit gets the precommit for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasAllocationPrecommit(ctx sdk.Context, val sdk.ValAddress, cel common.Address) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetAllocationPrecommitKey(val, cel))
}

// IterateAllocationPrecommits iterates over all prevotes in the store
func (k Keeper) IterateAllocationPrecommits(ctx sdk.Context, cb func(val sdk.ValAddress, cel common.Address, precommit types.AllocationPrecommit) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AllocationPrecommitKeyPrefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var precommit types.AllocationPrecommit
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), types.AllocationPrecommitKeyPrefix))
		val := sdk.ValAddress(keyPair.Next(20))
		cel := common.BytesToAddress(keyPair.Bytes())

		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &precommit)
		if cb(val, cel, precommit) {
			break
		}
	}
}

///////////////////////
// MsgAllocationCommit //
///////////////////////

// SetAllocationCommit sets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) SetAllocationCommit(ctx sdk.Context, val sdk.ValAddress, cel common.Address, allocationCommit types.Allocation) {
	bz := k.cdc.MustMarshalBinaryBare(&allocationCommit)
	ctx.KVStore(k.storeKey).Set(types.GetAllocationCommitForCellarKey(val, cel), bz)
}

// GetAllocationCommit gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetAllocationCommit(ctx sdk.Context, val sdk.ValAddress, cel common.Address) (types.Allocation, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetAllocationCommitForCellarKey(val, cel))
	if len(bz) == 0 {
		return types.Allocation{}, false
	}

	var allocationCommit types.Allocation
	k.cdc.MustUnmarshalBinaryBare(bz, &allocationCommit)
	return allocationCommit, true
}

// DeleteAllocationCommit deletes the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteAllocationCommit(ctx sdk.Context, val sdk.ValAddress, cel common.Address) {
	ctx.KVStore(k.storeKey).Delete(types.GetAllocationCommitForCellarKey(val, cel))
}

// HasAllocationCommitForCellar gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasAllocationCommitForCellar(ctx sdk.Context, val sdk.ValAddress, cel common.Address) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetAllocationCommitForCellarKey(val, cel))
}


// HasAllocationCommit gets the existence of any commit for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasAllocationCommit(ctx sdk.Context, val sdk.ValAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetAllocationCommitKeyPrefix(val))
	iter := store.Iterator(nil, nil)
	defer iter.Close()

	return iter.Valid()
}

// IterateAllocationCommits iterates over all votes in the store
func (k Keeper) IterateAllocationCommits(ctx sdk.Context, handler func(val sdk.ValAddress, cel common.Address, commit types.Allocation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AllocationCommitForCellarKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), types.AllocationCommitForCellarKeyPrefix))
		val := sdk.ValAddress(keyPair.Next(20))
		cel := common.BytesToAddress(keyPair.Bytes())

		var commit types.Allocation
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &commit)
		if handler(val, cel, commit) {
			break
		}
	}
}

// IterateAllocationCommitValidators iterates over all validators who have committed allocations
func (k Keeper) IterateAllocationCommitValidators(ctx sdk.Context, handler func(val sdk.ValAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AllocationCommitForCellarKeyPrefix)
	defer iter.Close()

	seenValidators := mapset.NewThreadUnsafeSet()

	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), types.AllocationCommitForCellarKeyPrefix))
		val := sdk.ValAddress(keyPair.Next(20))

		// add seen validator to set. if already in set, don't return to consumer
		if !seenValidators.Add(val) {
			continue
		}

		if handler(val) {
			break
		}
	}

}

// Iterates all of the commits for a provided validator
func (k Keeper) IterateValidatorAllocationCommits(ctx sdk.Context, val sdk.ValAddress, handler func(cellar common.Address, commit types.Allocation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefix := append(types.AllocationCommitForCellarKeyPrefix, val.Bytes()...)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), types.AllocationCommitForCellarKeyPrefix))
		keyPair.Next(20)
		cel := common.BytesToAddress(keyPair.Bytes())

		var commit types.Allocation
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &commit)
		if handler(cel, commit) {
			break
		}
	}
}

////////////////
// TickWeight //
////////////////

// GetAllocationTickWeights gets allocation tick weights for a validator for a given cellar
func (k Keeper) GetAllocationTickWeights(ctx sdk.Context, val sdk.ValAddress, cel common.Address) (types.TickWeights, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAllocationTickWeightKey(val, cel))
	if len(bz) == 0 {
		return types.TickWeights{}, false
	}

	var tickWeights types.TickWeights
	k.cdc.MustUnmarshalBinaryBare(bz, &tickWeights)
	return tickWeights, true
}

// SetAllocationTickWeights sets allocation tick weights for a validator for a given cellar
func (k Keeper) SetAllocationTickWeights(ctx sdk.Context, val sdk.ValAddress, cel common.Address, tickWeights types.TickWeights) {
	bz, err := k.cdc.MarshalInterface(&tickWeights)
	if err != nil {
		panic(err)
	}

	ctx.KVStore(k.storeKey).Set(types.GetAllocationTickWeightKey(val, cel), bz)
}

// DeleteAllocationTickWeights deletes the tick weights for a given validator and cellar
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteAllocationTickWeights(ctx sdk.Context, val sdk.ValAddress, cel common.Address) {
	ctx.KVStore(k.storeKey).Delete(types.GetAllocationTickWeightKey(val, cel))
}

// HasAllocationTickWeights returns if tick weights for a given validator and cellar exist
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasAllocationTickWeights(ctx sdk.Context, val sdk.ValAddress, cel common.Address) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetAllocationTickWeightKey(val, cel))
}

// IterateAllocationTickWeights iterates over all tick weights in the store
func (k Keeper) IterateAllocationTickWeights(ctx sdk.Context, handler func(val sdk.ValAddress, cel common.Address, tickWeights types.TickWeights) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AllocationTickWeightKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), types.AllocationTickWeightKeyPrefix))
		val := sdk.ValAddress(keyPair.Next(20))
		cel := common.BytesToAddress(keyPair.Bytes())

		var tickWeights types.TickWeights
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &tickWeights)
		if handler(val, cel, tickWeights) {
			break
		}
	}
}

// todo: aggregation functions for tick weight data

//////////////////
// CommitPeriod //
//////////////////

// SetCommitPeriodStart sets the current vote period start height
func (k Keeper) SetCommitPeriodStart(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CommitPeriodStartKey, sdk.Uint64ToBigEndian(uint64(height)))
}

// GetCommitPeriodStart returns the vote period start height
func (k Keeper) GetCommitPeriodStart(ctx sdk.Context) (int64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.CommitPeriodStartKey)
	if len(bz) == 0 {
		return 0, false
	}

	return int64(sdk.BigEndianToUint64(bz)), true
}

// HasCommitPeriodStart returns true if the vote period start has been set
func (k Keeper) HasCommitPeriodStart(ctx sdk.Context) bool {
	return ctx.KVStore(k.storeKey).Has(types.CommitPeriodStartKey)
}

/////////////////
// MissCounter //
/////////////////

// IncrementMissCounter increments the miss counter for a validator
func (k Keeper) IncrementMissCounter(ctx sdk.Context, val sdk.ValAddress) {
	missCounter := k.GetMissCounter(ctx, val)
	k.SetMissCounter(ctx, val, missCounter+1)
}

// GetMissCounter return the miss counter for a validator
func (k Keeper) GetMissCounter(ctx sdk.Context, val sdk.ValAddress) int64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMissCounterKey(val))
	if len(bz) == 0 {
		return 0
	}
	return int64(binary.BigEndian.Uint64(bz))
}

// SetMissCounter sets the miss counter for a given validator
func (k Keeper) SetMissCounter(ctx sdk.Context, val sdk.ValAddress, misses int64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(misses))
	ctx.KVStore(k.storeKey).Set(types.GetMissCounterKey(val), bz)
}

// HasMissCounter checks if a validator has an existing miss counter
func (k Keeper) HasMissCounter(ctx sdk.Context, val sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetMissCounterKey(val))
}

// DeleteMissCounter removes a validators miss counter
func (k Keeper) DeleteMissCounter(ctx sdk.Context, val sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetMissCounterKey(val))
}

// IterateMissCounters iterates over the miss counters
func (k Keeper) IterateMissCounters(ctx sdk.Context, handler func(val sdk.ValAddress, counter int64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.MissCounterKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		count := binary.BigEndian.Uint64(iter.Value())
		val := sdk.ValAddress(bytes.TrimPrefix(iter.Key(), types.MissCounterKeyPrefix))
		if handler(val, int64(count)) {
			break
		}
	}
}

// GetAllMissCounters returns all the miss counter values for all validators.
func (k Keeper) GetAllMissCounters(ctx sdk.Context) []types.MissCounter {
	missCounters := make([]types.MissCounter, 0)

	k.IterateMissCounters(ctx, func(validatorAddr sdk.ValAddress, counter int64) bool {
		missCounters = append(missCounters, types.MissCounter{
			Validator: validatorAddr.String(),
			Misses:    counter,
		})
		return false
	})

	return missCounters
}

////////////
// Params //
////////////

// GetParamSet returns the vote period from the parameters
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// SetParams sets the parameters in the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
