package keeper

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the oracle store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.BinaryCodec
	paramSpace    paramtypes.Subspace
	stakingKeeper types.StakingKeeper

	ics4Wrapper types.ICS4Wrapper
}

// NewKeeper creates a new x/cork Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper types.StakingKeeper, wrapper types.ICS4Wrapper,
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

		ics4Wrapper: wrapper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

/////////////////////
// Scheduled Corks //
/////////////////////

func (k Keeper) SetScheduledCork(ctx sdk.Context, chainID uint64, blockHeight uint64, val sdk.ValAddress, cork types.Cork) []byte {
	id := cork.IDHash(blockHeight)
	bz := k.cdc.MustMarshal(&cork)
	ctx.KVStore(k.storeKey).Set(types.GetScheduledCorkKey(chainID, blockHeight, id, val, common.HexToAddress(cork.TargetContractAddress)), bz)
	return id
}

// GetScheduledCork gets the scheduled cork for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetScheduledCork(ctx sdk.Context, chainID uint64, blockHeight uint64, id []byte, val sdk.ValAddress, contract common.Address) (types.Cork, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetScheduledCorkKey(chainID, blockHeight, id, val, contract))
	if len(bz) == 0 {
		return types.Cork{}, false
	}

	var cork types.Cork
	k.cdc.MustUnmarshal(bz, &cork)
	return cork, true
}

// DeleteScheduledCork deletes the scheduled cork for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteScheduledCork(ctx sdk.Context, chainID uint64, blockHeight uint64, id []byte, val sdk.ValAddress, contract common.Address) {
	ctx.KVStore(k.storeKey).Delete(types.GetScheduledCorkKey(chainID, blockHeight, id, val, contract))
}

// IterateScheduledCorks iterates over all scheduled corks in the store
func (k Keeper) IterateScheduledCorks(ctx sdk.Context, chainID uint64, cb func(val sdk.ValAddress, blockHeight uint64, id []byte, cel common.Address, cork types.Cork) (stop bool)) {
	k.IterateScheduledCorksByPrefix(ctx, types.GetScheduledCorkKeyPrefix(chainID), cb)
}

func (k Keeper) IterateScheduledCorksByBlockHeight(ctx sdk.Context, chainID uint64, blockHeight uint64, cb func(val sdk.ValAddress, blockHeight uint64, id []byte, cel common.Address, cork types.Cork) (stop bool)) {
	k.IterateScheduledCorksByPrefix(ctx, types.GetScheduledCorkKeyByBlockHeightPrefix(chainID, blockHeight), cb)
}

func (k Keeper) IterateScheduledCorksByPrefix(ctx sdk.Context, prefix []byte, cb func(val sdk.ValAddress, blockHeight uint64, id []byte, cel common.Address, cork types.Cork) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var cork types.Cork
		keyPair := bytes.NewBuffer(iter.Key())
		keyPair.Next(1) // trim prefix byte
		blockHeight := sdk.BigEndianToUint64(keyPair.Next(8))
		id := keyPair.Next(32)
		val := sdk.ValAddress(keyPair.Next(20))
		contract := common.BytesToAddress(keyPair.Bytes())

		k.cdc.MustUnmarshal(iter.Value(), &cork)
		if cb(val, blockHeight, id, contract, cork) {
			break
		}
	}
}

func (k Keeper) GetScheduledCorks(ctx sdk.Context, chainID uint64) []*types.ScheduledCork {
	var scheduledCorks []*types.ScheduledCork
	k.IterateScheduledCorks(ctx, chainID, func(val sdk.ValAddress, blockHeight uint64, id []byte, _ common.Address, cork types.Cork) (stop bool) {
		scheduledCorks = append(scheduledCorks, &types.ScheduledCork{
			Validator:   val.String(),
			Cork:        &cork,
			BlockHeight: blockHeight,
			Id:          id,
		})
		return false
	})

	return scheduledCorks
}

func (k Keeper) GetScheduledCorksByBlockHeight(ctx sdk.Context, chainID uint64, height uint64) []*types.ScheduledCork {
	var scheduledCorks []*types.ScheduledCork
	k.IterateScheduledCorksByBlockHeight(ctx, chainID, height, func(val sdk.ValAddress, blockHeight uint64, Id []byte, _ common.Address, cork types.Cork) (stop bool) {
		scheduledCorks = append(scheduledCorks, &types.ScheduledCork{
			Validator:   val.String(),
			Cork:        &cork,
			BlockHeight: blockHeight,
			Id:          Id,
		})

		return false
	})

	return scheduledCorks
}

func (k Keeper) GetScheduledCorksByID(ctx sdk.Context, chainID uint64, queriedID []byte) []*types.ScheduledCork {
	var scheduledCorks []*types.ScheduledCork
	k.IterateScheduledCorks(ctx, chainID, func(val sdk.ValAddress, blockHeight uint64, id []byte, _ common.Address, cork types.Cork) (stop bool) {
		if bytes.Equal(id, queriedID) {
			scheduledCorks = append(scheduledCorks, &types.ScheduledCork{
				Validator:   val.String(),
				Cork:        &cork,
				BlockHeight: blockHeight,
				Id:          id,
			})
		}

		return false
	})

	return scheduledCorks
}

/////////////////
// WinningCork //
/////////////////

func (k Keeper) SetWinningCork(ctx sdk.Context, chainID uint64, cork types.Cork) {
	bz := k.cdc.MustMarshal(&cork)
	ctx.KVStore(k.storeKey).Set(types.GetWinningCorkKey(chainID, common.HexToAddress(cork.TargetContractAddress)), bz)
}

func (k Keeper) GetWinningCork(ctx sdk.Context, chainID uint64, contract common.Address) (types.Cork, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetWinningCorkKey(chainID, contract))
	if len(bz) == 0 {
		return types.Cork{}, false
	}

	var cork types.Cork
	k.cdc.MustUnmarshal(bz, &cork)
	return cork, true
}

func (k Keeper) DeleteWinningCork(ctx sdk.Context, chainID uint64, cork types.Cork) {
	ctx.KVStore(k.storeKey).Delete(types.GetWinningCorkKey(chainID, common.HexToAddress(cork.TargetContractAddress)))
}

///////////////////////////
// ScheduledBlockHeights //
///////////////////////////

func (k Keeper) GetScheduledBlockHeights(ctx sdk.Context, chainID uint64) []uint64 {
	var heights []uint64

	latestHeight := uint64(0)
	k.IterateScheduledCorks(ctx, chainID, func(_ sdk.ValAddress, blockHeight uint64, _ []byte, _ common.Address, _ types.Cork) (stop bool) {
		if blockHeight > latestHeight {
			heights = append(heights, blockHeight)
		}
		latestHeight = blockHeight
		return false
	})

	return heights
}

/////////////////////////
// Invalidation Nonces //
/////////////////////////

func (k Keeper) GetLatestInvalidationNonce(ctx sdk.Context, chainID uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LatestInvalidationNonceKey(chainID))
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLatestInvalidationNonce(ctx sdk.Context, chainID uint64, invalidationNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LatestInvalidationNonceKey(chainID), sdk.Uint64ToBigEndian(invalidationNonce))
}

func (k Keeper) IncrementInvalidationNonce(ctx sdk.Context, chainID uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	nextNonce := k.GetLatestInvalidationNonce(ctx, chainID) + 1
	store.Set(types.LatestInvalidationNonceKey(chainID), sdk.Uint64ToBigEndian(nextNonce))
	return nextNonce
}

//////////////////
// Cork Results //
//////////////////

func (k Keeper) SetCorkResult(ctx sdk.Context, chainID uint64, id []byte, corkResult types.CorkResult) {
	bz := k.cdc.MustMarshal(&corkResult)
	ctx.KVStore(k.storeKey).Set(types.GetCorkResultKey(chainID, id), bz)
}

func (k Keeper) GetCorkResult(ctx sdk.Context, chainID uint64, id []byte) (types.CorkResult, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCorkResultKey(chainID, id))
	if len(bz) == 0 {
		return types.CorkResult{}, false
	}

	var corkResult types.CorkResult
	k.cdc.MustUnmarshal(bz, &corkResult)
	return corkResult, true
}

func (k Keeper) DeleteCorkResult(ctx sdk.Context, chainID uint64, id []byte) {
	ctx.KVStore(k.storeKey).Delete(types.GetCorkResultKey(chainID, id))
}

// IterateCorksResult iterates over all cork results in the store
func (k Keeper) IterateCorkResults(ctx sdk.Context, chainID uint64, cb func(id []byte, blockHeight uint64, approved bool, approvalPercentage string, corkResult types.CorkResult) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetCorkResultPrefix(chainID))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var corkResult types.CorkResult
		keyPair := bytes.NewBuffer(iter.Key())
		keyPair.Next(1) // trim prefix byte
		id := keyPair.Next(32)

		k.cdc.MustUnmarshal(iter.Value(), &corkResult)
		if cb(id, corkResult.BlockHeight, corkResult.Approved, corkResult.ApprovalPercentage, corkResult) {
			break
		}
	}
}

// GetCorkResults returns CorkResults
func (k Keeper) GetCorkResults(ctx sdk.Context, chainID uint64) []*types.CorkResult {
	var corkResults []*types.CorkResult
	k.IterateCorkResults(ctx, chainID, func(id []byte, blockHeight uint64, approved bool, approvalPercentage string, corkResult types.CorkResult) (stop bool) {
		corkResults = append(corkResults, &corkResult)
		return false
	})

	return corkResults
}

///////////
// Votes //
///////////

func (k Keeper) GetApprovedScheduledCorks(ctx sdk.Context, chainID uint64, currentBlockHeight uint64, threshold sdk.Dec) (approvedCorks []types.Cork) {
	corksForBlockHeight := make(map[uint64][]types.Cork)
	corkPowersForBlockHeight := make(map[uint64][]uint64)

	totalPower := k.stakingKeeper.GetLastTotalPower(ctx)

	k.IterateScheduledCorks(ctx, chainID, func(val sdk.ValAddress, scheduledBlockHeight uint64, id []byte, addr common.Address, cork types.Cork) (stop bool) {
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

		k.DeleteScheduledCork(ctx, chainID, scheduledBlockHeight, id, val, addr)

		return false
	})

	for blockHeight := range corkPowersForBlockHeight {
		for i, power := range corkPowersForBlockHeight[blockHeight] {
			cork := corksForBlockHeight[blockHeight][i]
			approvalPercentage := sdk.NewIntFromUint64(power).ToDec().Quo(totalPower.ToDec())
			quorumReached := approvalPercentage.GT(threshold)
			corkResult := types.CorkResult{
				Cork:               &cork,
				BlockHeight:        blockHeight,
				Approved:           quorumReached,
				ApprovalPercentage: approvalPercentage.Mul(sdk.NewDec(100)).String(),
			}
			corkID := cork.IDHash(blockHeight)

			k.SetCorkResult(ctx, chainID, corkID, corkResult)

			if quorumReached {
				approvedCorks = append(approvedCorks, cork)
			}
		}
	}

	return approvedCorks
}

/////////////
// Cellars //
/////////////

func (k Keeper) SetCellarIDs(ctx sdk.Context, chainID uint64, c types.CellarIDSet) {
	bz := k.cdc.MustMarshal(&c)
	ctx.KVStore(k.storeKey).Set(types.MakeCellarIDsKey(chainID), bz)
}

func (k Keeper) GetCellarIDs(ctx sdk.Context, chainID uint64) (cellars []common.Address) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MakeCellarIDsKey(chainID))

	var cids types.CellarIDSet
	k.cdc.MustUnmarshal(bz, &cids)

	for _, cid := range cids.Ids {
		cellars = append(cellars, common.HexToAddress(cid))
	}

	return cellars
}

func (k Keeper) HasCellarID(ctx sdk.Context, chainID uint64, address common.Address) (found bool) {
	found = false
	for _, id := range k.GetCellarIDs(ctx, chainID) {
		if id == address {
			found = true
			break
		}
	}

	return found
}
