package keeper

import (
	"bytes"
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/x"
	"github.com/peggyjv/sommelier/x/allocation/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the oracle store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryMarshaler
	paramSpace    paramtypes.Subspace
	stakingKeeper x.StakingKeeper

	handlerSet bool
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper x.StakingKeeper,
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

//////////////////////////
// MsgDelegateAllocations //
//////////////////////////

// SetValidatorDelegateAddress sets a new address that will have the power to send Allocations on behalf of the validator
func (k Keeper) SetValidatorDelegateAddress(ctx sdk.Context, del sdk.AccAddress, val sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetFeedDelegateKey(del), val.Bytes())
}

// GetValidatorAddressFromDelegate returns the delegate address for a given validator
func (k Keeper) GetValidatorAddressFromDelegate(ctx sdk.Context, del sdk.AccAddress) sdk.ValAddress {
	return ctx.KVStore(k.storeKey).Get(types.GetFeedDelegateKey(del))
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

// IsDelegateAddress returns true if the validator has delegated their Allocations to an address
func (k Keeper) IsDelegateAddress(ctx sdk.Context, del sdk.AccAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetFeedDelegateKey(del))
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

//////////////////////////
// MsgAllocationPrecommit //
//////////////////////////

func (k Keeper) SetAllocationPrecommit(ctx sdk.Context, validatorAddr sdk.ValAddress, precommit types.AllocationPrecommit) {
	bz := k.cdc.MustMarshalBinaryBare(&precommit)
	ctx.KVStore(k.storeKey).Set(types.GetAllocationPrecommitKey(validatorAddr), bz)
}

func (k Keeper) GetAllocationPrecommit(ctx sdk.Context, val sdk.ValAddress) (types.AllocationPrecommit, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetAllocationPrecommitKey(val))
	if len(bz) == 0 {
		return types.AllocationPrecommit{}, false
	}

	var precommit types.AllocationPrecommit
	k.cdc.MustUnmarshalBinaryBare(bz, &precommit)
	return precommit, true
}

func (k Keeper) DeleteAllPrecommits(ctx sdk.Context) {
	k.IterateAllocationPrecommits(ctx, func(val sdk.ValAddress, _ types.AllocationPrecommit) bool {
		k.DeleteAllocationPrecommit(ctx, val)
		return false
	})
}

func (k Keeper) DeleteAllocationPrecommit(ctx sdk.Context, val sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetAllocationPrecommitKey(val))
}

func (k Keeper) HasAllocationPrecommit(ctx sdk.Context, val sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetAllocationPrecommitKey(val))
}

func (k Keeper) IterateAllocationPrecommits(ctx sdk.Context, cb func(val sdk.ValAddress, precommit types.AllocationPrecommit) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AllocationPrecommitKeyPrefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var precommit types.AllocationPrecommit
		val := sdk.ValAddress(bytes.TrimPrefix(iter.Key(), types.AllocationPrecommitKeyPrefix))

		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &precommit)
		if cb(val, precommit) {
			break
		}
	}
}

///////////////////////
// MsgAllocationCommit //
///////////////////////

func (k Keeper) SetAllocationCommit(ctx sdk.Context, validatorAddr sdk.ValAddress, Allocation types.Allocation) {
	bz := k.cdc.MustMarshalBinaryBare(&Allocation)
	ctx.KVStore(k.storeKey).Set(types.GetAllocationCommitKey(validatorAddr), bz)
}

func (k Keeper) GetAllocationCommit(ctx sdk.Context, val sdk.ValAddress) (types.Allocation, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetAllocationCommitKey(val))
	if len(bz) == 0 {
		return types.Allocation{}, false
	}

	var Allocation types.Allocation
	k.cdc.MustUnmarshalBinaryBare(bz, &Allocation)
	return Allocation, true
}

func (k Keeper) DeleteAllocationCommit(ctx sdk.Context, val sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetAllocationCommitKey(val))
}

func (k Keeper) HasAllocationCommit(ctx sdk.Context, val sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetAllocationCommitKey(val))
}

func (k Keeper) IterateAllocationCommits(ctx sdk.Context, handler func(val sdk.ValAddress, vote types.Allocation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AllocationCommitKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {

		val := sdk.ValAddress(bytes.TrimPrefix(iter.Key(), types.AllocationCommitKeyPrefix))

		var Allocation types.Allocation
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &Allocation)
		if handler(val, Allocation) {
			break
		}
	}
}

////////////////
// VotePeriod //
////////////////

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
