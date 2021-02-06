package keeper

import (
	"bytes"
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
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
func (k Keeper) SetHandler(handlerFn types.OracleHandler) {
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
func (k Keeper) SetValidatorDelegateAddress(ctx sdk.Context, val, del sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetFeedDelegateKey(del), val.Bytes())
}

// GetValidatorAddressFromDelegate returns the delegate address for a given validator
func (k Keeper) GetValidatorAddressFromDelegate(ctx sdk.Context, del sdk.AccAddress) sdk.AccAddress {
	return sdk.AccAddress(ctx.KVStore(k.storeKey).Get(types.GetFeedDelegateKey(del)))
}

// GetDelegateAddressFromValidator returns the validator address for a given delegate
func (k Keeper) GetDelegateAddressFromValidator(ctx sdk.Context, val sdk.AccAddress) sdk.AccAddress {
	var address sdk.AccAddress
	// TODO: create secondary index
	k.IterateDelegateAddresses(ctx, func(delegatorAddr, validatorAddr sdk.AccAddress) bool {
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
	return ctx.KVStore(k.storeKey).Has(types.GetFeedDelegateKey(del))
}

// IterateDelegateAddresses iterates over all delegate address pairs in the store
func (k Keeper) IterateDelegateAddresses(ctx sdk.Context, handler func(del, val sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.FeedDelegateKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		del := sdk.AccAddress(bytes.TrimPrefix(iter.Key(), types.FeedDelegateKeyPrefix))
		val := sdk.AccAddress(iter.Value())
		if handler(del, val) {
			break
		}
	}
}

//////////////////////////
// MsgOracleDataPrevote //
//////////////////////////

// SetOracleDataPrevote sets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) SetOracleDataPrevote(ctx sdk.Context, val sdk.AccAddress, prevote *types.MsgOracleDataPrevote) {
	ctx.KVStore(k.storeKey).Set(types.GetOracleDataPrevoteKey(val), k.cdc.MustMarshalBinaryBare(prevote))
}

// GetOracleDataPrevote gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetOracleDataPrevote(ctx sdk.Context, val sdk.AccAddress) *types.MsgOracleDataPrevote {
	var out types.MsgOracleDataPrevote
	k.cdc.MustUnmarshalBinaryBare(ctx.KVStore(k.storeKey).Get(types.GetOracleDataPrevoteKey(val)), &out)
	return &out
}

func (k Keeper) DeleteAllPrevotes(ctx sdk.Context) {
	k.IterateOracleDataPrevotes(ctx, func(val sdk.AccAddress, _ *types.MsgOracleDataPrevote) bool {
		k.DeleteOracleDataPrevote(ctx, val)
		return false
	})
}

// DeleteOracleDataPrevote deletes the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteOracleDataPrevote(ctx sdk.Context, val sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetOracleDataPrevoteKey(val))
}

// HasOracleDataPrevote gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasOracleDataPrevote(ctx sdk.Context, val sdk.AccAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOracleDataPrevoteKey(val))
}

// IterateOracleDataPrevotes iterates over all prevotes in the store
func (k Keeper) IterateOracleDataPrevotes(ctx sdk.Context, handler func(val sdk.AccAddress, msg *types.MsgOracleDataPrevote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.OracleDataPrevoteKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var out types.MsgOracleDataPrevote
		val := sdk.AccAddress(bytes.TrimPrefix(iter.Key(), types.OracleDataPrevoteKeyPrefix))
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &out)
		if handler(val, &out) {
			break
		}
	}
}

///////////////////////
// MsgOracleDataVote //
///////////////////////

// SetOracleDataVote sets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) SetOracleDataVote(ctx sdk.Context, val sdk.AccAddress, msg *types.MsgOracleDataVote) {
	ctx.KVStore(k.storeKey).Set(types.GetOracleDataVoteKey(val), k.cdc.MustMarshalBinaryBare(msg))
}

// GetOracleDataVote gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetOracleDataVote(ctx sdk.Context, val sdk.AccAddress) *types.MsgOracleDataVote {
	var out types.MsgOracleDataVote
	k.cdc.MustUnmarshalBinaryBare(ctx.KVStore(k.storeKey).Get(types.GetOracleDataVoteKey(val)), &out)
	return &out
}

// DeleteOracleDataVote deletes the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteOracleDataVote(ctx sdk.Context, val sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetOracleDataVoteKey(val))
}

// HasOracleDataVote gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasOracleDataVote(ctx sdk.Context, val sdk.AccAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOracleDataVoteKey(val))
}

// IterateOracleDataVotes iterates over all votes in the store
func (k Keeper) IterateOracleDataVotes(ctx sdk.Context, handler func(val sdk.AccAddress, vote *types.MsgOracleDataVote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.OracleDataVoteKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote types.MsgOracleDataVote
		val := sdk.AccAddress(bytes.TrimPrefix(iter.Key(), types.OracleDataVoteKeyPrefix))
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &vote)
		if handler(val, &vote) {
			break
		}
	}
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

// GetOracleData gets oracle data stored for a given type
func (k Keeper) GetOracleData(ctx sdk.Context, dataType, id string) (types.OracleData, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOracleDataKey(dataType, id))
	if len(bz) == 0 {
		return nil, false
	}

	var oracleData types.OracleData
	k.cdc.UnmarshalInterface(bz, &oracleData)
	return oracleData, true
}

// DeleteOracleData deletes the data from the oracle of a given type
func (k Keeper) DeleteOracleData(ctx sdk.Context, dataType, id string) {
	ctx.KVStore(k.storeKey).Delete(types.GetOracleDataKey(dataType, id))
}

// HasOracleData returns true if a given type exists in the store
func (k Keeper) HasOracleData(ctx sdk.Context, dataType, id string) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOracleDataKey(dataType, id))
}

////////////////
// VotePeriod //
////////////////

// SetVotePeriodStart sets the vote period start height
func (k Keeper) SetVotePeriodStart(ctx sdk.Context, h int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.VotePeriodStartKey, sdk.Uint64ToBigEndian(uint64(h)))
}

// GetVotePeriodStart returns the vote period start height
func (k Keeper) GetVotePeriodStart(ctx sdk.Context) (int64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VotePeriodStartKey)
	if len(bz) == 0 {
		return 0, false
	}

	return int64(sdk.BigEndianToUint64(bz)), true
}

// HasVotePeriodStart returns true if the vote period start has been set
func (k Keeper) HasVotePeriodStart(ctx sdk.Context) bool {
	return ctx.KVStore(k.storeKey).Has(types.VotePeriodStartKey)
}

/////////////////
// MissCounter //
/////////////////

// IncrementMissCounter increments the miss counter for a validator
func (k Keeper) IncrementMissCounter(ctx sdk.Context, val sdk.AccAddress) {
	if !k.HasMissCounter(ctx, val) {
		k.SetMissCounter(ctx, val, 1)
		return
	}
	k.SetMissCounter(ctx, val, k.GetMissCounter(ctx, val)+1)
}

// GetMissCounter return the miss counter for a validator
func (k Keeper) GetMissCounter(ctx sdk.Context, val sdk.AccAddress) int64 {
	return int64(binary.BigEndian.Uint64(ctx.KVStore(k.storeKey).Get(types.GetMissCounterKey(val))))
}

// SetMissCounter sets the miss counter for a given validator
func (k Keeper) SetMissCounter(ctx sdk.Context, val sdk.AccAddress, misses int64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(misses))
	ctx.KVStore(k.storeKey).Set(types.GetMissCounterKey(val), bz)
}

// HasMissCounter checks if a validator has an existing miss counter
func (k Keeper) HasMissCounter(ctx sdk.Context, val sdk.AccAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetMissCounterKey(val))
}

// DeleteMissCounter removes a validators miss counter
func (k Keeper) DeleteMissCounter(ctx sdk.Context, val sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetMissCounterKey(val))
}

// IterateMissCounters iterates over the miss counters
func (k Keeper) IterateMissCounters(ctx sdk.Context, handler func(val sdk.AccAddress, counter int64) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.MissCounterKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		count := binary.BigEndian.Uint64(iter.Value())
		val := sdk.AccAddress(bytes.TrimPrefix(iter.Key(), types.MissCounterKeyPrefix))
		if handler(val, int64(count)) {
			break
		}
	}
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
