package keeper

import (
	"bytes"
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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
	ctx.KVStore(k.storeKey).Set(types.GetAllocationCommitKey(val, cel), bz)
}

// GetAllocationCommit gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetAllocationCommit(ctx sdk.Context, val sdk.ValAddress, cel common.Address) (types.Allocation, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetAllocationCommitKey(val, cel))
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
	ctx.KVStore(k.storeKey).Delete(types.GetAllocationCommitKey(val, cel))
}

// HasAllocationCommit gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasAllocationCommit(ctx sdk.Context, val sdk.ValAddress, cel common.Address) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetAllocationCommitKey(val, cel))
}

// IterateAllocationCommits iterates over all votes in the store
func (k Keeper) IterateAllocationCommits(ctx sdk.Context, handler func(val sdk.ValAddress, cel common.Address, commit types.Allocation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AllocationCommitKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), types.AllocationCommitKeyPrefix))
		val := sdk.ValAddress(keyPair.Next(20))
		cel := common.BytesToAddress(keyPair.Bytes())

		var commit types.Allocation
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &commit)
		if handler(val, cel, commit) {
			break
		}
	}
}

////////////////
// TickWeight //
////////////////

// GetOracleDataType gets oracle data stored for a given type
func (k Keeper) GetAllocationTickWeight(ctx sdk.Context, val sdk.ValAddress, cel common.Address) (types.TickWeights, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetAllocationTickWeightKey(val, cel))
	if len(bz) == 0 {
		return types.TickWeights{}, false
	}

	var tickWeights types.TickWeights
	k.cdc.MustUnmarshalBinaryBare(bz, &tickWeights)
	return tickWeights, true
}


////////////////
// OracleData //
////////////////

// GetOracleDataType gets oracle data stored for a given type
func (k Keeper) GetOracleDataType(ctx sdk.Context, id string) (string, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOracleDataTypeKey(id))
	if len(bz) == 0 {
		return "", false
	}

	return string(bz), true
}

// HasOracleDataType checks the oracle data type associated with a given data identifier
func (k Keeper) HasOracleDataType(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetOracleDataTypeKey(id))
}

// SetOracleDataType sets oracle data type associated with a given data identifier
func (k Keeper) SetOracleDataType(ctx sdk.Context, dataType, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOracleDataTypeKey(id), []byte(dataType))
}

// SetOracleData sets the oracle data in the store
func (k Keeper) SetOracleData(ctx sdk.Context, oracleData types.OracleData) {
	bz, err := k.cdc.MarshalInterface(oracleData)
	if err != nil {
		panic(err)
	}

	ctx.KVStore(k.storeKey).Set(types.GetOracleDataKey(oracleData.Type(), oracleData.GetID()), bz)
}

// GetAggregatedOracleData gets the aggregated oracle data from the store by height, type and id
func (k Keeper) GetAggregatedOracleData(ctx sdk.Context, height int64, dataType, id string) types.OracleData {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAggregatedOracleDataKey(uint64(height), dataType, id)

	bz := store.Get(key)
	if len(bz) == 0 {
		return nil
	}

	var oracleData types.OracleData
	err := k.cdc.UnmarshalInterface(bz, &oracleData)
	if err != nil {
		panic(err)
	}

	return oracleData
}

// GetLatestAggregatedOracleData returns the latest stored aggregated data
func (k Keeper) GetLatestAggregatedOracleData(ctx sdk.Context, dataType, id string) (types.OracleData, int64) {
	// get the latest stored height for the given id
	height := k.GetOracleDataHeight(ctx, id)
	if height == 0 {
		return nil, 0
	}

	return k.GetAggregatedOracleData(ctx, height, dataType, id), height
}

// IterateAggregatedOracleData iterates over all aggregated data in the store
func (k Keeper) IterateAggregatedOracleData(ctx sdk.Context, cb func(height int64, oracleData types.OracleData) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	// NOTE: we use reverse prefix iterator to iterate from latest to earliest height
	iter := sdk.KVStoreReversePrefixIterator(store, types.AggregatedOracleDataKeyPrefix)

	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var oracleData types.OracleData

		height := sdk.BigEndianToUint64(iter.Key()[1:9])
		if err := k.cdc.UnmarshalInterface(iter.Value(), &oracleData); err != nil {
			panic(err)
		}

		if cb(int64(height), oracleData) {
			break
		}
	}
}

// GetAllAggregatedData returns all the aggregated data
func (k Keeper) GetAllAggregatedData(ctx sdk.Context) []types.AggregatedOracleData {
	var aggregates []types.AggregatedOracleData

	k.IterateAggregatedOracleData(ctx, func(height int64, oracleData types.OracleData) bool {
		pair, ok := oracleData.(*types.UniswapPair)
		if !ok {
			return false
		}

		aggregates = append(aggregates, types.AggregatedOracleData{
			Height: height,
			Data:   pair,
		})

		return false
	})

	return aggregates
}

// HasAggregatedOracleData sets the aggregated oracle data in the store by height, type and id
func (k Keeper) HasAggregatedOracleData(ctx sdk.Context, height int64, dataType, id string) bool {
	key := types.GetAggregatedOracleDataKey(uint64(height), dataType, id)
	return ctx.KVStore(k.storeKey).Has(key)
}

// SetAggregatedOracleData sets the aggregated oracle data in the store by height, type and id
func (k Keeper) SetAggregatedOracleData(ctx sdk.Context, height int64, oracleData types.OracleData) {
	bz, err := k.cdc.MarshalInterface(oracleData)
	if err != nil {
		panic(err)
	}

	key := types.GetAggregatedOracleDataKey(uint64(height), oracleData.Type(), oracleData.GetID())

	ctx.KVStore(k.storeKey).Set(key, bz)
}

func (k Keeper) GetOracleDataHeight(ctx sdk.Context, id string) int64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOracleDataHeightKey(id))
	if len(bz) == 0 {
		return 0
	}

	return int64(sdk.BigEndianToUint64(bz))
}

// SetOracleDataHeight sets latest aggregated oracle data height associated with a given data identifier
func (k Keeper) SetOracleDataHeight(ctx sdk.Context, id string, height int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOracleDataHeightKey(id), sdk.Uint64ToBigEndian(uint64(height)))
}

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
