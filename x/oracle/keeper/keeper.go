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
	StakingKeeper types.StakingKeeper

	storeKey   sdk.StoreKey
	cdc        codec.BinaryMarshaler
	paramSpace paramtypes.Subspace
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
		StakingKeeper: stakingKeeper,
		storeKey:      key,
		cdc:           cdc,
		paramSpace:    paramSpace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
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

// GetDelegateAddressFromValidator returns the valdiator address for a given delegate
func (k Keeper) GetDelegateAddressFromValidator(ctx sdk.Context, val sdk.AccAddress) sdk.AccAddress {
	var out sdk.AccAddress
	k.IterateDelegateAddresses(ctx, func(del, sval sdk.AccAddress) bool {
		if val.Equals(sval) {
			out = del
			return true
		}
		return false
	})
	return out
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

// SetOracleData sets the oracle data in the store
func (k Keeper) SetOracleData(ctx sdk.Context, od types.OracleData) {
	switch data := od.(type) {
	case *types.UniswapData:
		ctx.KVStore(k.storeKey).Set(types.GetOracleDataKey(od.Type()), k.cdc.MustMarshalBinaryBare(data))
	default:
		panic("NOT HERE")
	}
}

// GetOracleData gets oracle data stored for a given type
func (k Keeper) GetOracleData(ctx sdk.Context, typ string) types.OracleData {
	switch typ {
	case types.UniswapDataType:
		var data types.UniswapData
		k.cdc.MustUnmarshalBinaryBare(ctx.KVStore(k.storeKey).Get(types.GetOracleDataKey(typ)), &data)
		return &data
	default:
		k.Logger(ctx).Info("data type not supported")
		return nil
	}
}

// DeleteOracleData deletes the data from the oracle of a given type
func (k Keeper) DeleteOracleData(ctx sdk.Context, typ string) {
	ctx.KVStore(k.storeKey).Delete(types.GetOracleDataKey(typ))
}

// HasOracleData returns true if a given type exists in the store
func (k Keeper) HasOracleData(ctx sdk.Context, typ string) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOracleDataKey(typ))
}

////////////////
// VotePeriod //
////////////////

// SetVotePeriodStart sets the vote period start height
func (k Keeper) SetVotePeriodStart(ctx sdk.Context, h int64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(h))
	ctx.KVStore(k.storeKey).Set(types.VotePeriodStartKey, bz)
}

// GetVotePeriodStart returns the vote period start height
func (k Keeper) GetVotePeriodStart(ctx sdk.Context) int64 {
	return int64(binary.BigEndian.Uint64(ctx.KVStore(k.storeKey).Get(types.VotePeriodStartKey)))
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
