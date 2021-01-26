package keeper

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the oracle store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryMarshaler
	paramSpace paramtypes.Subspace
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace paramtypes.Subspace,
) Keeper {

	// // set KeyTable if it has not already been set
	// if !paramSpace.HasKeyTable() {
	// 	paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	// }

	return Keeper{
		storeKey:   key,
		cdc:        cdc,
		paramSpace: paramSpace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

////////////////////////
// MsgDelegateAddress //
////////////////////////

// SetDelegateAddress sets a new address that will have the power to send data on behalf of the validator
func (k Keeper) SetDelegateAddress(ctx sdk.Context, val, del sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetFeedDelegateKey(val), del.Bytes())
}

// GetDelegateAddress returns the delegate address for a given validator
func (k Keeper) GetDelegateAddress(ctx sdk.Context, val sdk.AccAddress) sdk.AccAddress {
	return sdk.AccAddress(ctx.KVStore(k.storeKey).Get(types.GetFeedDelegateKey(val)))
}

// HasDelegateAddress returns true if the validator has delegated their feed to an address
func (k Keeper) HasDelegateAddress(ctx sdk.Context, val sdk.AccAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetFeedDelegateKey(val))
}

//////////////////////////
// MsgOracleDataPrevote //
//////////////////////////

// SetOracleDataPrevote sets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) SetOracleDataPrevote(ctx sdk.Context, val sdk.AccAddress, hashes [][]byte) {
	ctx.KVStore(k.storeKey).Set(types.GetOracleDataPrevoteKey(val), bytes.Join(hashes, []byte(",")))
}

// GetOracleDataPrevote gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetOracleDataPrevote(ctx sdk.Context, val sdk.AccAddress) [][]byte {
	return bytes.Split(ctx.KVStore(k.storeKey).Get(types.GetOracleDataPrevoteKey(val)), []byte(","))
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
func (k Keeper) IterateOracleDataPrevotes(ctx sdk.Context, handler func(val sdk.AccAddress, hashes [][]byte) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.OracleDataPrevoteKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		val := sdk.AccAddress(iter.Key())
		hashes := bytes.Split(iter.Value(), []byte(","))
		if handler(val, hashes) {
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
func (k Keeper) GetOracleDataVote(ctx sdk.Context, val sdk.AccAddress) (out *types.MsgOracleDataVote) {
	k.cdc.MustUnmarshalBinaryBare(ctx.KVStore(k.storeKey).Get(types.GetOracleDataVoteKey(val)), out)
	return
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
	iter := sdk.KVStorePrefixIterator(store, types.OracleDataPrevoteKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vote *types.MsgOracleDataVote
		val := sdk.AccAddress(iter.Key())
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), vote)
		if handler(val, vote) {
			break
		}
	}
}
