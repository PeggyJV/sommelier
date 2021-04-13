package keeper

import (
	"bytes"

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

	handlerSet    bool
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
// MsgDelegateDecisions //
//////////////////////////

// SetValidatorDelegateAddress sets a new address that will have the power to send decisions on behalf of the validator
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

// IsDelegateAddress returns true if the validator has delegated their decisions to an address
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
// MsgDecisionPrecommit //
//////////////////////////

func (k Keeper) SetDecisionPrecommit(ctx sdk.Context, validatorAddr sdk.ValAddress, precommit types.DecisionPrecommit) {
	bz := k.cdc.MustMarshalBinaryBare(&precommit)
	ctx.KVStore(k.storeKey).Set(types.GetDecisionPrecommitKey(validatorAddr), bz)
}

func (k Keeper) GetDecisionPrecommit(ctx sdk.Context, val sdk.ValAddress) (types.DecisionPrecommit, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetDecisionPrecommitKey(val))
	if len(bz) == 0 {
		return types.DecisionPrecommit{}, false
	}

	var precommit types.DecisionPrecommit
	k.cdc.MustUnmarshalBinaryBare(bz, &precommit)
	return precommit, true
}

func (k Keeper) DeleteAllPrecommits(ctx sdk.Context) {
	k.IterateDecisionPrecommits(ctx, func(val sdk.ValAddress, _ types.DecisionPrecommit) bool {
		k.DeleteDecisionPrecommit(ctx, val)
		return false
	})
}

func (k Keeper) DeleteDecisionPrecommit(ctx sdk.Context, val sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetDecisionPrecommitKey(val))
}

func (k Keeper) HasDecisionPrecommit(ctx sdk.Context, val sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetDecisionPrecommitKey(val))
}

func (k Keeper) IterateDecisionPrecommits(ctx sdk.Context, cb func(val sdk.ValAddress, precommit types.DecisionPrecommit) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AllocationDecisionPrecommitKeyPrefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var precommit types.DecisionPrecommit
		val := sdk.ValAddress(bytes.TrimPrefix(iter.Key(), types.AllocationDecisionPrecommitKeyPrefix))

		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &precommit)
		if cb(val, precommit) {
			break
		}
	}
}

///////////////////////
// MsgDecisionCommit //
///////////////////////

func (k Keeper) SetDecisionCommit(ctx sdk.Context, validatorAddr sdk.ValAddress, decision types.Decision) {
	bz := k.cdc.MustMarshalBinaryBare(&decision)
	ctx.KVStore(k.storeKey).Set(types.GetDecisionCommitKey(validatorAddr), bz)
}

func (k Keeper) GetDecisionCommit(ctx sdk.Context, val sdk.ValAddress) (types.Decision, bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetDecisionCommitKey(val))
	if len(bz) == 0 {
		return types.Decision{}, false
	}

	var decision types.Decision
	k.cdc.MustUnmarshalBinaryBare(bz, &decision)
	return decision, true
}

func (k Keeper) DeleteDecisionCommit(ctx sdk.Context, val sdk.ValAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetDecisionCommitKey(val))
}

func (k Keeper) HasDecisionCommit(ctx sdk.Context, val sdk.ValAddress) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetDecisionCommitKey(val))
}

func (k Keeper) IterateDecisionCommits(ctx sdk.Context, handler func(val sdk.ValAddress, vote types.Decision) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AllocationDecisionCommitKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {

		val := sdk.ValAddress(bytes.TrimPrefix(iter.Key(), types.AllocationDecisionCommitKeyPrefix))

		var decision types.Decision
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &decision)
		if handler(val, decision) {
			break
		}
	}
}
