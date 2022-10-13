package keeper

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
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

// NewKeeper creates a new distribution Keeper instance
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

/////////////////////
// Scheduled Corks //
/////////////////////

func (k Keeper) SetScheduledCork(ctx sdk.Context, blockHeight uint64, val sdk.ValAddress, cork types.Cork) uint64 {
	id, found := k.GetCorkId(ctx, blockHeight, cork)

	if !found {
		id = k.IncrementScheduledCorkID(ctx)
		k.SetCorkId(ctx, blockHeight, cork, id)
	}

	bz := k.cdc.MustMarshal(&cork)
	ctx.KVStore(k.storeKey).Set(types.GetScheduledCorkKey(blockHeight, id, val, common.HexToAddress(cork.TargetContractAddress)), bz)
	return id
}

// GetScheduledCork gets the scheduled cork for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetScheduledCork(ctx sdk.Context, blockHeight uint64, id uint64, val sdk.ValAddress, contract common.Address) (types.Cork, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetScheduledCorkKey(blockHeight, id, val, contract))
	if len(bz) == 0 {
		return types.Cork{}, false
	}

	var cork types.Cork
	k.cdc.MustUnmarshal(bz, &cork)
	return cork, true
}

// DeleteScheduledCork deletes the scheduled cork for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteScheduledCork(ctx sdk.Context, blockHeight uint64, id uint64, val sdk.ValAddress, contract common.Address) {
	ctx.KVStore(k.storeKey).Delete(types.GetScheduledCorkKey(blockHeight, id, val, contract))
}

// IterateScheduledCorks iterates over all scheduled corks in the store
func (k Keeper) IterateScheduledCorks(ctx sdk.Context, cb func(val sdk.ValAddress, blockHeight uint64, id uint64, cel common.Address, cork types.Cork) (stop bool)) {
	k.IterateScheduledCorksByPrefix(ctx, types.GetScheduledCorkKeyPrefix(), cb)
}

func (k Keeper) IterateScheduledCorksByBlockHeight(ctx sdk.Context, blockHeight uint64, cb func(val sdk.ValAddress, blockHeight uint64, id uint64, cel common.Address, cork types.Cork) (stop bool)) {
	k.IterateScheduledCorksByPrefix(ctx, types.GetScheduledCorkKeyByBlockHeightPrefix(blockHeight), cb)
}

func (k Keeper) IterateScheduledCorksByPrefix(ctx sdk.Context, prefix []byte, cb func(val sdk.ValAddress, blockHeight uint64, id uint64, cel common.Address, cork types.Cork) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var cork types.Cork
		keyPair := bytes.NewBuffer(iter.Key())
		keyPair.Next(1) // trim prefix byte
		blockHeight := sdk.BigEndianToUint64(keyPair.Next(8))
		id := sdk.BigEndianToUint64(keyPair.Next(8))
		val := sdk.ValAddress(keyPair.Next(20))
		contract := common.BytesToAddress(keyPair.Bytes())

		k.cdc.MustUnmarshal(iter.Value(), &cork)
		if cb(val, blockHeight, id, contract, cork) {
			break
		}
	}
}

func (k Keeper) GetScheduledCorks(ctx sdk.Context) []*types.ScheduledCork {
	var scheduledCorks []*types.ScheduledCork
	k.IterateScheduledCorks(ctx, func(val sdk.ValAddress, blockHeight uint64, id uint64, _ common.Address, cork types.Cork) (stop bool) {
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

func (k Keeper) GetScheduledCorksByBlockHeight(ctx sdk.Context, height uint64) []*types.ScheduledCork {
	var scheduledCorks []*types.ScheduledCork
	k.IterateScheduledCorksByBlockHeight(ctx, height, func(val sdk.ValAddress, blockHeight uint64, Id uint64, _ common.Address, cork types.Cork) (stop bool) {
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

func (k Keeper) GetScheduledCorksByID(ctx sdk.Context, id uint64) []*types.ScheduledCork {
	var scheduledCorks []*types.ScheduledCork
	k.IterateScheduledCorks(ctx, func(val sdk.ValAddress, blockHeight uint64, scheduledCorkId uint64, _ common.Address, cork types.Cork) (stop bool) {
		if scheduledCorkId == id {
			scheduledCorks = append(scheduledCorks, &types.ScheduledCork{
				Validator:   val.String(),
				Cork:        &cork,
				BlockHeight: blockHeight,
				Id:          id,
			})
		}

		if scheduledCorkId > id {
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
	k.IterateScheduledCorks(ctx, func(_ sdk.ValAddress, blockHeight uint64, Id uint64, _ common.Address, _ types.Cork) (stop bool) {
		if blockHeight > latestHeight {
			heights = append(heights, blockHeight)
		}
		latestHeight = blockHeight
		return false
	})

	return heights
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

///////////////////////
// Scheduled Cork ID //
///////////////////////

func (k Keeper) GetCorkIdByHash(ctx sdk.Context, hash []byte) (uint64, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetScheduledCorkIdHashKey(hash))
	if len(bz) == 0 {
		return 0, false
	}
	return sdk.BigEndianToUint64(bz), true
}

func (k Keeper) GetCorkId(ctx sdk.Context, blockHeight uint64, cork types.Cork) (uint64, bool) {
	return k.GetCorkIdByHash(ctx, cork.IdHash(blockHeight))
}

func (k Keeper) SetCorkIdByHash(ctx sdk.Context, hash []byte, id uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetScheduledCorkIdHashKey(hash), sdk.Uint64ToBigEndian(id))
}

func (k Keeper) SetCorkId(ctx sdk.Context, blockHeight uint64, cork types.Cork, id uint64) {
	k.SetCorkIdByHash(ctx, cork.IdHash(blockHeight), id)
}

func (k Keeper) GetLatestScheduledCorkID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.LatestCorkIDKey})
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLatestScheduledCorkID(ctx sdk.Context, ScheduledCorkID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{types.LatestCorkIDKey}, sdk.Uint64ToBigEndian(ScheduledCorkID))
}

func (k Keeper) IncrementScheduledCorkID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	nextID := k.GetLatestScheduledCorkID(ctx) + 1
	store.Set([]byte{types.LatestCorkIDKey}, sdk.Uint64ToBigEndian(nextID))
	return nextID
}

//////////////////
// Cork Results //
//////////////////

func (k Keeper) SetCorkResult(ctx sdk.Context, id uint64, corkResult types.CorkResult) {
	bz := k.cdc.MustMarshal(&corkResult)
	ctx.KVStore(k.storeKey).Set(types.GetCorkResultKey(id), bz)
}

func (k Keeper) GetCorkResult(ctx sdk.Context, id uint64) (types.CorkResult, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetCorkResultKey(id))
	if len(bz) == 0 {
		return types.CorkResult{}, false
	}

	var corkResult types.CorkResult
	k.cdc.MustUnmarshal(bz, &corkResult)
	return corkResult, true
}

func (k Keeper) DeleteCorkResult(ctx sdk.Context, id uint64) {
	ctx.KVStore(k.storeKey).Delete(types.GetCorkResultKey(id))
}

// IterateCorksResult iterates over all cork results in the store
func (k Keeper) IterateCorkResults(ctx sdk.Context, cb func(id uint64, blockHeight uint64, approved bool, approvalPercentage string, corkResult types.CorkResult) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetCorkResultPrefix())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var corkResult types.CorkResult
		keyPair := bytes.NewBuffer(iter.Key())
		keyPair.Next(1) // trim prefix byte
		id := sdk.BigEndianToUint64(keyPair.Bytes())

		k.cdc.MustUnmarshal(iter.Value(), &corkResult)
		if cb(id, corkResult.BlockHeight, corkResult.Approved, corkResult.ApprovalPercentage, corkResult) {
			break
		}
	}
}

// GetCorkResults returns CorkResults
func (k Keeper) GetCorkResults(ctx sdk.Context) []*types.CorkResult {
	var corkResults []*types.CorkResult
	k.IterateCorkResults(ctx, func(id uint64, blockHeight uint64, approved bool, approvalPercentage string, corkResult types.CorkResult) (stop bool) {
		corkResults = append(corkResults, &corkResult)
		return false
	})

	return corkResults
}

///////////
// Votes //
///////////

func (k Keeper) GetApprovedScheduledCorks(ctx sdk.Context, currentBlockHeight uint64, threshold sdk.Dec) (approvedCorks []types.Cork) {
	corksForBlockHeight := make(map[uint64][]types.Cork)
	corkPowersForBlockHeight := make(map[uint64][]uint64)

	totalPower := k.stakingKeeper.GetLastTotalPower(ctx)

	k.IterateScheduledCorks(ctx, func(val sdk.ValAddress, scheduledBlockHeight uint64, id uint64, addr common.Address, cork types.Cork) (stop bool) {
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

		k.DeleteScheduledCork(ctx, scheduledBlockHeight, id, val, addr)

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
			corkId, found := k.GetCorkId(ctx, blockHeight, cork)
			if !found {
				// this should be impossible, corks are assigned an ID when written to the store
				panic(fmt.Sprintf("trying to vote on cork with no cork ID at block height %d: %s", blockHeight, cork))
			}
			k.SetCorkResult(ctx, corkId, corkResult)
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
