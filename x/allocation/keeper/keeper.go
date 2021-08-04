package keeper

import (
	"bytes"

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
	gravityKeeper types.GravityKeeper
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper types.StakingKeeper, gravityKeeper types.GravityKeeper,
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
		gravityKeeper: gravityKeeper,
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
	iter := sdk.KVStorePrefixIterator(store, []byte{types.AllocationDelegateKeyPrefix})
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		del := sdk.AccAddress(bytes.TrimPrefix(iter.Key(), []byte{types.AllocationDelegateKeyPrefix}))
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


/////////////
// Cellars //
/////////////

func (k Keeper) SetCellar(ctx sdk.Context, c types.Cellar) {
	bz := k.cdc.MustMarshalBinaryBare(&c)
	ctx.KVStore(k.storeKey).Set(types.GetCellarKey(c.Id), bz)
}

func (k Keeper) GetCellarByID(ctx sdk.Context, id common.Address) (cellar types.Cellar, found bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetCellarKey(id.))
	if len(bz) == 0 {
		return types.Cellar{}, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &cellar)
	return cellar, true
}

func (k Keeper) IterateCellars(ctx sdk.Context, cb func(cellar types.Cellar) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.CellarKeyPrefix})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var cellar types.Cellar
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &cellar)
		if cb(cellar) {
			break
		}
	}
}

func (k Keeper) GetCellars(ctx sdk.Context) (cellars []types.Cellar) {
	k.IterateCellars(ctx, func(cellar types.Cellar) (stop bool) {
		cellars = append(cellars, cellar)
		return false
	})

	return
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
	iter := sdk.KVStorePrefixIterator(store, []byte{types.AllocationPrecommitKeyPrefix})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var precommit types.AllocationPrecommit
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.AllocationPrecommitKeyPrefix}))
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
	iter := sdk.KVStorePrefixIterator(store, []byte{types.AllocationCommitForCellarKeyPrefix})
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.AllocationCommitForCellarKeyPrefix}))
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
	iter := sdk.KVStorePrefixIterator(store, []byte{types.AllocationCommitForCellarKeyPrefix})
	defer iter.Close()

	seenValidators := mapset.NewThreadUnsafeSet()

	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.AllocationCommitForCellarKeyPrefix}))
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

// IterateValidatorAllocationCommits Iterates all of the commits for a provided validator
func (k Keeper) IterateValidatorAllocationCommits(ctx sdk.Context, val sdk.ValAddress, handler func(cellar common.Address, commit types.Allocation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, append([]byte{types.AllocationCommitForCellarKeyPrefix}, val.Bytes()...))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.AllocationCommitForCellarKeyPrefix}))
		keyPair.Next(20)
		cel := common.BytesToAddress(keyPair.Bytes())

		var commit types.Allocation
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &commit)
		if handler(cel, commit) {
			break
		}
	}
}

//////////////////
// CommitPeriod //
//////////////////

// SetCommitPeriodStart sets the current vote period start height
func (k Keeper) SetCommitPeriodStart(ctx sdk.Context, height int64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.CommitPeriodStartKey}, sdk.Uint64ToBigEndian(uint64(height)))
}

// GetCommitPeriodStart returns the vote period start height
func (k Keeper) GetCommitPeriodStart(ctx sdk.Context) (int64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.CommitPeriodStartKey})
	if len(bz) == 0 {
		return 0, false
	}

	return int64(sdk.BigEndianToUint64(bz)), true
}

// HasCommitPeriodStart returns true if the vote period start has been set
func (k Keeper) HasCommitPeriodStart(ctx sdk.Context) bool {
	return ctx.KVStore(k.storeKey).Has([]byte{types.CommitPeriodStartKey})
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

/////////////////////////
// Invalidation Nonces //
/////////////////////////

func (k Keeper) GetLatestInvalidationNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.LatestInvalidationNonceKey})
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLatestInvalidationNonce(ctx sdk.Context, invalidationNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.LatestInvalidationNonceKey}, sdk.Uint64ToBigEndian(invalidationNonce))
}
