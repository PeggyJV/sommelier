package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the oracle store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	paramSpace    paramtypes.Subspace
	stakingKeeper types.StakingKeeper
	gravityKeeper types.GravityKeeper
}

// NewKeeper creates a new x/cork Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
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

///////////////////////
// MsgSubmitCork //
///////////////////////

// SetCork sets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) SetCork(ctx sdk.Context, val sdk.ValAddress, cork types.Cork) {
	bz := k.cdc.MustMarshal(&cork)
	ctx.KVStore(k.storeKey).Set(types.GetCorkForValidatorAddressKey(val, common.HexToAddress(cork.TargetContractAddress)), bz)
}

// GetCork gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetCork(ctx sdk.Context, val sdk.ValAddress, contract common.Address) (types.Cork, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetCorkForValidatorAddressKey(val, contract))
	if len(bz) == 0 {
		return types.Cork{}, false
	}

	var cork types.Cork
	k.cdc.MustUnmarshal(bz, &cork)
	return cork, true
}

// DeleteCork deletes the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteCork(ctx sdk.Context, val sdk.ValAddress, contract common.Address) {
	ctx.KVStore(k.storeKey).Delete(types.GetCorkForValidatorAddressKey(val, contract))
}

// HasCorkForContract gets the prevote for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasCorkForContract(ctx sdk.Context, val sdk.ValAddress, contract common.Address) bool {
	return ctx.KVStore(k.storeKey).Has(types.GetCorkForValidatorAddressKey(val, contract))
}

// HasCork gets the existence of any commit for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) HasCork(ctx sdk.Context, val sdk.ValAddress) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.GetCorkValidatorKeyPrefix(val))
	iter := store.Iterator(nil, nil)
	defer iter.Close()

	return iter.Valid()
}

// IterateCorks iterates over all votes in the store
func (k Keeper) IterateCorks(ctx sdk.Context, handler func(val sdk.ValAddress, contract common.Address, cork types.Cork) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.CorkForAddressKeyPrefix})
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.CorkForAddressKeyPrefix}))
		val := sdk.ValAddress(keyPair.Next(20))
		contract := common.BytesToAddress(keyPair.Bytes())

		var cork types.Cork
		k.cdc.MustUnmarshal(iter.Value(), &cork)
		if handler(val, contract, cork) {
			break
		}
	}
}

func (k Keeper) GetValidatorCorks(ctx sdk.Context) []*types.ValidatorCork {
	var validatorCorks []*types.ValidatorCork
	k.IterateCorks(ctx, func(val sdk.ValAddress, _ common.Address, cork types.Cork) (stop bool) {
		validatorCorks = append(validatorCorks, &types.ValidatorCork{
			Validator: val.String(),
			Cork:      &cork,
		})
		return false
	})

	return validatorCorks
}

/////////////////////
// Scheduled Corks //
/////////////////////

func (k Keeper) SetScheduledCork(ctx sdk.Context, blockHeight uint64, val sdk.ValAddress, cork types.Cork) {
	bz := k.cdc.MustMarshal(&cork)
	ctx.KVStore(k.storeKey).Set(types.GetScheduledCorkKey(blockHeight, val, common.HexToAddress(cork.TargetContractAddress)), bz)
}

// GetScheduledCork gets the scheduled cork for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetScheduledCork(ctx sdk.Context, val sdk.ValAddress, blockHeight uint64, contract common.Address) (types.Cork, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetScheduledCorkKey(blockHeight, val, contract))
	if len(bz) == 0 {
		return types.Cork{}, false
	}

	var cork types.Cork
	k.cdc.MustUnmarshal(bz, &cork)
	return cork, true
}

// DeleteScheduledCork deletes the scheduled cork for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteScheduledCork(ctx sdk.Context, val sdk.ValAddress, blockHeight uint64, contract common.Address) {
	ctx.KVStore(k.storeKey).Delete(types.GetScheduledCorkKey(blockHeight, val, contract))
}

// IterateScheduledCorks iterates over all scheduled corks in the store
func (k Keeper) IterateScheduledCorks(ctx sdk.Context, cb func(val sdk.ValAddress, blockHeight uint64, cel common.Address, cork types.Cork) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.ScheduledCorkKeyPrefix})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var cork types.Cork
		keyPair := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.ScheduledCorkKeyPrefix}))
		blockHeight := binary.BigEndian.Uint64(keyPair.Next(8))
		val := sdk.ValAddress(keyPair.Next(20))
		cel := common.BytesToAddress(keyPair.Bytes())

		k.cdc.MustUnmarshal(iter.Value(), &cork)
		if cb(val, blockHeight, cel, cork) {
			break
		}
	}
}

func (k Keeper) GetScheduledCorks(ctx sdk.Context) []*types.ScheduledCork {
	var scheduledCorks []*types.ScheduledCork
	k.IterateScheduledCorks(ctx, func(val sdk.ValAddress, blockHeight uint64, _ common.Address, cork types.Cork) (stop bool) {
		scheduledCorks = append(scheduledCorks, &types.ScheduledCork{
			Validator:   val.String(),
			Cork:        &cork,
			BlockHeight: blockHeight,
		})
		return false
	})

	return scheduledCorks
}

func (k Keeper) GetScheduledCorksByBlockHeight(ctx sdk.Context, height uint64) []*types.ScheduledCork {
	var scheduledCorks []*types.ScheduledCork
	k.IterateScheduledCorks(ctx, func(val sdk.ValAddress, blockHeight uint64, _ common.Address, cork types.Cork) (stop bool) {
		if blockHeight == height {
			scheduledCorks = append(scheduledCorks, &types.ScheduledCork{
				Validator:   val.String(),
				Cork:        &cork,
				BlockHeight: blockHeight,
			})
		}
		if blockHeight > height {
			return true
		}
		return false
	})

	return scheduledCorks
}

///////////////////////////
// ScheduledBlockHeights //
///////////////////////////

func (k Keeper) GetScheduledBlockHeights(ctx sdk.Context) []uint64 {
	var heights []uint64

	latestHeight := uint64(0)
	k.IterateScheduledCorks(ctx, func(_ sdk.ValAddress, blockHeight uint64, _ common.Address, _ types.Cork) (stop bool) {
		if blockHeight > latestHeight {
			heights = append(heights, blockHeight)
		}
		latestHeight = blockHeight
		return false
	})

	return heights
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

// setParams sets the parameters in the store
func (k Keeper) setParams(ctx sdk.Context, params types.Params) {
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

func (k Keeper) IncrementInvalidationNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	nextNonce := k.GetLatestInvalidationNonce(ctx) + 1
	store.Set([]byte{types.LatestInvalidationNonceKey}, sdk.Uint64ToBigEndian(nextNonce))
	return nextNonce
}

///////////
// Votes //
///////////

func (k Keeper) GetApprovedCorks(ctx sdk.Context, threshold sdk.Dec) (approvedCorks []types.Cork) {

	var corks []types.Cork
	var corkPowers []int64

	totalPower := k.stakingKeeper.GetLastTotalPower(ctx)

	k.IterateCorks(ctx, func(val sdk.ValAddress, addr common.Address, cork types.Cork) (stop bool) {
		fmt.Printf("DEBUG: val = %v", val.String())
		validator := k.stakingKeeper.Validator(ctx, val)
		validatorPower := validator.GetConsensusPower(k.stakingKeeper.PowerReduction(ctx))

		found := false
		for i, c := range corks {
			if c.Equals(cork) {
				corkPowers[i] += validatorPower

				found = true
				break
			}
		}

		if !found {
			corks = append(corks, cork)
			corkPowers = append(corkPowers, validatorPower)
		}

		k.DeleteCork(ctx, val, addr)

		return false
	})

	var winningCorks []types.Cork

	for i, power := range corkPowers {
		quorumReached := sdk.NewDec(power).Quo(totalPower.ToDec()).GT(threshold)
		if quorumReached {
			winningCorks = append(winningCorks, corks[i])
		}
	}

	return winningCorks
}

func (k Keeper) GetApprovedScheduledCorks(ctx sdk.Context, currentBlockHeight uint64, threshold sdk.Dec) (approvedCorks []types.Cork) {

	corksForBlockHeight := make(map[uint64][]types.Cork)
	corkPowersForBlockHeight := make(map[uint64][]uint64)

	totalPower := k.stakingKeeper.GetLastTotalPower(ctx)

	k.IterateScheduledCorks(ctx, func(val sdk.ValAddress, scheduledBlockHeight uint64, addr common.Address, cork types.Cork) (stop bool) {
		if currentBlockHeight != scheduledBlockHeight {
			// only operate on scheduled corksForBlockHeight that are valid, quit early when a further one is found
			return true
		}

		validator := k.stakingKeeper.Validator(ctx, val)
		validatorPower := uint64(validator.GetConsensusPower(k.stakingKeeper.PowerReduction(ctx)))

		found := false
		for i, c := range corksForBlockHeight[scheduledBlockHeight] {
			if c.Equals(cork) {
				corkPowersForBlockHeight[scheduledBlockHeight][i] += validatorPower

				found = true
				break
			}
		}

		if !found {
			corksForBlockHeight[scheduledBlockHeight] = append(corksForBlockHeight[scheduledBlockHeight], cork)
			corkPowersForBlockHeight[scheduledBlockHeight] = append(corkPowersForBlockHeight[scheduledBlockHeight], validatorPower)
		}

		k.DeleteScheduledCork(ctx, val, scheduledBlockHeight, addr)

		return false
	})

	for blockHeight := range corkPowersForBlockHeight {
		for i, power := range corkPowersForBlockHeight[blockHeight] {
			quorumReached := sdk.NewIntFromUint64(power).ToDec().Quo(totalPower.ToDec()).GT(threshold)
			if quorumReached {
				approvedCorks = append(approvedCorks, corksForBlockHeight[blockHeight][i])
			}
		}
	}

	return approvedCorks
}

/////////////
// Cellars //
/////////////

func (k Keeper) SetCellarIDs(ctx sdk.Context, c types.CellarIDSet) {
	bz := k.cdc.MustMarshal(&c)
	ctx.KVStore(k.storeKey).Set(types.MakeCellarIDsKey(), bz)
}

func (k Keeper) GetCellarIDs(ctx sdk.Context) (cellars []common.Address) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MakeCellarIDsKey())

	var cids types.CellarIDSet
	k.cdc.MustUnmarshal(bz, &cids)

	for _, cid := range cids.Ids {
		cellars = append(cellars, common.HexToAddress(cid))
	}

	return cellars
}

func (k Keeper) HasCellarID(ctx sdk.Context, address common.Address) (found bool) {
	found = false
	for _, id := range k.GetCellarIDs(ctx) {
		if id == address {
			found = true
			break
		}
	}

	return found
}
