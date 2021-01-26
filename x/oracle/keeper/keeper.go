package keeper

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
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

	// // set KeyTable if it has not already been set
	// if !paramSpace.HasKeyTable() {
	// 	paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	// }

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
	iter := sdk.KVStorePrefixIterator(store, types.OracleDataPrevoteKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		del := sdk.AccAddress(iter.Key())
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

////////////////
// OracleData //
////////////////

// SetOracleData sets the oracle data in the store
func (k Keeper) SetOracleData(ctx sdk.Context, od types.OracleData) {
	oda, err := types.PackOracleData(od)
	if err != nil {
		panic(err)
	}
	ctx.KVStore(k.storeKey).Set(types.GetOracleDataKey(od.Type()), k.cdc.MustMarshalBinaryBare(oda))
}

// GetOracleData gets oracle data stored for a given type
func (k Keeper) GetOracleData(ctx sdk.Context, typ string) types.OracleData {
	var any *cdctypes.Any
	k.cdc.MustUnmarshalBinaryBare(ctx.KVStore(k.storeKey).Get(types.GetOracleDataKey(typ)), any)
	od, err := types.UnpackOracleData(any)
	if err != nil {
		panic(err)
	}
	return od
}

// HasOracleData returns true if a given type exists in the store
func (k Keeper) HasOracleData(ctx sdk.Context, typ string) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetOracleDataKey(typ))
}
