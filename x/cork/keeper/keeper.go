package keeper

import (
	"bytes"
	"encoding/binary"
	"sort"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/common"
	corktypes "github.com/peggyjv/sommelier/v8/x/cork/types"
	types "github.com/peggyjv/sommelier/v8/x/cork/types/v2"
)

const corkVoteThresholdStr = "0.67"

// Keeper of the oracle store
type Keeper struct {
	storeKey      storetypes.StoreKey
	cdc           codec.BinaryCodec
	paramSpace    paramtypes.Subspace
	stakingKeeper corktypes.StakingKeeper
	gravityKeeper corktypes.GravityKeeper
	pubsubKeeper  corktypes.PubsubKeeper
}

// NewKeeper creates a new x/cork Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper corktypes.StakingKeeper, gravityKeeper corktypes.GravityKeeper, pubsubKeeper corktypes.PubsubKeeper,
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
		pubsubKeeper:  pubsubKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+corktypes.ModuleName)
}

/////////////////////
// Scheduled Corks //
/////////////////////

func (k Keeper) SetScheduledCork(ctx sdk.Context, blockHeight uint64, val sdk.ValAddress, cork types.Cork) []byte {
	id := cork.IDHash(blockHeight)
	bz := k.cdc.MustMarshal(&cork)
	ctx.KVStore(k.storeKey).Set(corktypes.GetScheduledCorkKey(blockHeight, id, val, common.HexToAddress(cork.TargetContractAddress)), bz)
	return id
}

// GetScheduledCork gets the scheduled cork for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) GetScheduledCork(ctx sdk.Context, blockHeight uint64, id []byte, val sdk.ValAddress, contract common.Address) (types.Cork, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(corktypes.GetScheduledCorkKey(blockHeight, id, val, contract))
	if len(bz) == 0 {
		return types.Cork{}, false
	}

	var cork types.Cork
	k.cdc.MustUnmarshal(bz, &cork)
	return cork, true
}

// DeleteScheduledCork deletes the scheduled cork for a given validator
// CONTRACT: must provide the validator address here not the delegate address
func (k Keeper) DeleteScheduledCork(ctx sdk.Context, blockHeight uint64, id []byte, val sdk.ValAddress, contract common.Address) {
	ctx.KVStore(k.storeKey).Delete(corktypes.GetScheduledCorkKey(blockHeight, id, val, contract))
}

// IterateScheduledCorks iterates over all scheduled corks in the store
func (k Keeper) IterateScheduledCorks(ctx sdk.Context, cb func(val sdk.ValAddress, blockHeight uint64, id []byte, cel common.Address, cork types.Cork) (stop bool)) {
	k.IterateScheduledCorksByPrefix(ctx, corktypes.GetScheduledCorkKeyPrefix(), cb)
}

func (k Keeper) IterateScheduledCorksByBlockHeight(ctx sdk.Context, blockHeight uint64, cb func(val sdk.ValAddress, blockHeight uint64, id []byte, cel common.Address, cork types.Cork) (stop bool)) {
	k.IterateScheduledCorksByPrefix(ctx, corktypes.GetScheduledCorkKeyByBlockHeightPrefix(blockHeight), cb)
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

func (k Keeper) GetScheduledCorks(ctx sdk.Context) []*types.ScheduledCork {
	var scheduledCorks []*types.ScheduledCork
	k.IterateScheduledCorks(ctx, func(val sdk.ValAddress, blockHeight uint64, id []byte, _ common.Address, cork types.Cork) (stop bool) {
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
	k.IterateScheduledCorksByBlockHeight(ctx, height, func(val sdk.ValAddress, blockHeight uint64, Id []byte, _ common.Address, cork types.Cork) (stop bool) {
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

func (k Keeper) GetScheduledCorksByID(ctx sdk.Context, queriedID []byte) []*types.ScheduledCork {
	var scheduledCorks []*types.ScheduledCork
	k.IterateScheduledCorks(ctx, func(val sdk.ValAddress, blockHeight uint64, id []byte, _ common.Address, cork types.Cork) (stop bool) {
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

///////////////////////////
// ScheduledBlockHeights //
///////////////////////////

func (k Keeper) GetScheduledBlockHeights(ctx sdk.Context) []uint64 {
	var heights []uint64

	latestHeight := uint64(0)
	k.IterateScheduledCorks(ctx, func(_ sdk.ValAddress, blockHeight uint64, _ []byte, _ common.Address, _ types.Cork) (stop bool) {
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
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

/////////////////////////
// Invalidation Nonces //
/////////////////////////

func (k Keeper) GetLatestInvalidationNonce(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{corktypes.LatestInvalidationNonceKey})
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetLatestInvalidationNonce(ctx sdk.Context, invalidationNonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte{corktypes.LatestInvalidationNonceKey}, sdk.Uint64ToBigEndian(invalidationNonce))
}

func (k Keeper) IncrementInvalidationNonce(ctx sdk.Context) uint64 {
	nextNonce := k.GetLatestInvalidationNonce(ctx) + 1
	k.SetLatestInvalidationNonce(ctx, nextNonce)
	return nextNonce
}

//////////////////
// Cork Results //
//////////////////

func (k Keeper) SetCorkResult(ctx sdk.Context, id []byte, corkResult types.CorkResult) {
	bz := k.cdc.MustMarshal(&corkResult)
	ctx.KVStore(k.storeKey).Set(corktypes.GetCorkResultKey(id), bz)
}

func (k Keeper) GetCorkResult(ctx sdk.Context, id []byte) (types.CorkResult, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(corktypes.GetCorkResultKey(id))
	if len(bz) == 0 {
		return types.CorkResult{}, false
	}

	var corkResult types.CorkResult
	k.cdc.MustUnmarshal(bz, &corkResult)
	return corkResult, true
}

func (k Keeper) DeleteCorkResult(ctx sdk.Context, id []byte) {
	ctx.KVStore(k.storeKey).Delete(corktypes.GetCorkResultKey(id))
}

// IterateCorkResults iterates over all cork results in the store
func (k Keeper) IterateCorkResults(ctx sdk.Context, cb func(id []byte, blockHeight uint64, approved bool, approvalPercentage string, corkResult types.CorkResult) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, corktypes.GetCorkResultPrefix())
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
func (k Keeper) GetCorkResults(ctx sdk.Context) []*types.CorkResult {
	var corkResults []*types.CorkResult
	k.IterateCorkResults(ctx, func(id []byte, blockHeight uint64, approved bool, approvalPercentage string, corkResult types.CorkResult) (stop bool) {
		corkResults = append(corkResults, &corkResult)
		return false
	})

	return corkResults
}

///////////
// Votes //
///////////

func (k Keeper) GetApprovedScheduledCorks(ctx sdk.Context) (approvedCorks []types.Cork) {
	currentBlockHeight := uint64(ctx.BlockHeight())
	totalPower := k.stakingKeeper.GetLastTotalPower(ctx)
	corks := []types.Cork{}
	powers := []uint64{}
	k.IterateScheduledCorksByBlockHeight(ctx, currentBlockHeight, func(val sdk.ValAddress, _ uint64, id []byte, addr common.Address, cork types.Cork) (stop bool) {
		validator := k.stakingKeeper.Validator(ctx, val)
		validatorPower := uint64(validator.GetConsensusPower(k.stakingKeeper.PowerReduction(ctx)))
		found := false
		for i, c := range corks {
			if c.Equals(cork) {
				powers[i] += validatorPower

				found = true
				break
			}
		}

		if !found {
			corks = append(corks, cork)
			powers = append(powers, validatorPower)
		}

		k.DeleteScheduledCork(ctx, currentBlockHeight, id, val, addr)
		k.DecrementValidatorCorkCount(ctx, val)

		return false
	})

	threshold := sdk.MustNewDecFromStr(corkVoteThresholdStr)
	for i, power := range powers {
		cork := corks[i]
		approvalPercentage := sdk.NewDecFromInt(sdk.NewIntFromUint64(power)).Quo(sdk.NewDecFromInt(totalPower))
		quorumReached := approvalPercentage.GT(threshold)
		corkResult := types.CorkResult{
			Cork:               &cork,
			BlockHeight:        currentBlockHeight,
			Approved:           quorumReached,
			ApprovalPercentage: approvalPercentage.Mul(sdk.NewDec(100)).String(),
		}
		corkID := cork.IDHash(currentBlockHeight)

		k.SetCorkResult(ctx, corkID, corkResult)

		if quorumReached {
			approvedCorks = append(approvedCorks, cork)
		}
	}

	return approvedCorks
}

/////////////
// Cellars //
/////////////

func (k Keeper) SetCellarIDs(ctx sdk.Context, c types.CellarIDSet) {
	bz := k.cdc.MustMarshal(&c)
	// always sort before writing to the store
	cellarIDs := make([]string, 0, len(c.Ids))
	cellarIDs = append(cellarIDs, c.Ids...)
	sort.Strings(cellarIDs)
	c.Ids = cellarIDs
	ctx.KVStore(k.storeKey).Set(corktypes.MakeCellarIDsKey(), bz)
}

func (k Keeper) GetCellarIDs(ctx sdk.Context) (cellars []common.Address) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(corktypes.MakeCellarIDsKey())

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

///////////////////////////
// Validator Cork counts //
///////////////////////////

func (k Keeper) GetValidatorCorkCount(ctx sdk.Context, val sdk.ValAddress) (count uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(corktypes.GetValidatorCorkCountKey(val))
	if len(bz) == 0 {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) SetValidatorCorkCount(ctx sdk.Context, val sdk.ValAddress, count uint64) {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	ctx.KVStore(k.storeKey).Set(corktypes.GetValidatorCorkCountKey(val), bz)
}

func (k Keeper) IncrementValidatorCorkCount(ctx sdk.Context, val sdk.ValAddress) {
	count := k.GetValidatorCorkCount(ctx, val)
	k.SetValidatorCorkCount(ctx, val, count+1)
}

func (k Keeper) DecrementValidatorCorkCount(ctx sdk.Context, val sdk.ValAddress) {
	count := k.GetValidatorCorkCount(ctx, val)
	if count > 0 {
		k.SetValidatorCorkCount(ctx, val, count-1)
	}
}
