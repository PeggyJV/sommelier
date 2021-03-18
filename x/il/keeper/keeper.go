package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/peggyjv/sommelier/x/il/types"
)

// Keeper of the impermanent store
type Keeper struct {
	cdc        codec.BinaryMarshaler
	storeKey   sdk.StoreKey
	paramSpace paramtypes.Subspace

	oracleKeeper    types.OracleKeeper
	ethBridgeKeeper types.EthBridgeKeeper
}

// NewKeeper constructs a new keeper for oracle
func NewKeeper(cdc codec.BinaryMarshaler, storeKey sdk.StoreKey,
	paramspace paramtypes.Subspace, oracleKeeper types.OracleKeeper,
	ethBridgeKeeper types.EthBridgeKeeper) Keeper {

	// set KeyTable if it has not already been set
	if !paramspace.HasKeyTable() {
		paramspace = paramspace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		paramSpace:      paramspace,
		oracleKeeper:    oracleKeeper,
		ethBridgeKeeper: ethBridgeKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// HasStoplossPosition returns true if the stoploss exists
func (k Keeper) HasStoplossPosition(ctx sdk.Context, address sdk.AccAddress, uniswapPair string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StoplossKeyPrefix)
	return store.Has(types.StoplossKey(address, uniswapPair))
}

// GetStoplossPosition returns a stoploss for the given address an pair.
func (k Keeper) GetStoplossPosition(ctx sdk.Context, address sdk.AccAddress, uniswapPair string) (types.Stoploss, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StoplossKeyPrefix)
	bz := store.Get(types.StoplossKey(address, uniswapPair))
	if len(bz) == 0 {
		return types.Stoploss{}, false
	}
	var stoploss types.Stoploss
	k.cdc.MustUnmarshalBinaryBare(bz, &stoploss)

	return stoploss, true
}

// SetStoplossPosition returns a stoploss for the given address an pair.
func (k Keeper) SetStoplossPosition(ctx sdk.Context, address sdk.AccAddress, stoploss types.Stoploss) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StoplossKeyPrefix)
	bz := k.cdc.MustMarshalBinaryBare(&stoploss)
	store.Set(types.StoplossKey(address, stoploss.UniswapPairID), bz)
}

// DeleteStoplossPosition removes a stoploss position from the store.
func (k Keeper) DeleteStoplossPosition(ctx sdk.Context, address sdk.AccAddress, uniswapPair string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StoplossKeyPrefix)
	store.Delete(types.StoplossKey(address, uniswapPair))
}

// IterateStoplossPositions iterates over the all the stoploss keys and performs a callback.
func (k Keeper) IterateStoplossPositions(ctx sdk.Context, cb func(sdk.AccAddress, types.Stoploss) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StoplossKeyPrefix)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		lpAddress := types.LPAddressFromStoplossKey(append(types.StoplossKeyPrefix, iterator.Key()...))
		var stoploss types.Stoploss
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &stoploss)

		if cb(lpAddress, stoploss) {
			break
		}
	}
}

// IterateStoplossPositionsByAddress iterates over the all the stoploss owned by an address and performs a callback.
func (k Keeper) IterateStoplossPositionsByAddress(ctx sdk.Context, address sdk.AccAddress, cb func(types.Stoploss) (stop bool)) {
	keyPrefix := append(types.StoplossKeyPrefix, address.Bytes()...)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), keyPrefix)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var stoploss types.Stoploss
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &stoploss)

		if cb(stoploss) {
			break
		}
	}
}

// GetLPStoplossPositions retrieves all the positions owned by a given LPs.
func (k Keeper) GetLPStoplossPositions(ctx sdk.Context, address sdk.AccAddress) []types.Stoploss {
	positions := []types.Stoploss{}

	k.IterateStoplossPositionsByAddress(ctx, address, func(stoploss types.Stoploss) bool {
		positions = append(positions, stoploss)
		return false
	})

	return positions
}

// GetLPsStoplossPositions retrieves all the positions from the LPs.
func (k Keeper) GetLPsStoplossPositions(ctx sdk.Context) types.LPsStoplossPositions {
	lps := types.LPsStoplossPositions{}

	addresses := make([]string, 0)
	positionsMap := make(map[string][]types.Stoploss)

	k.IterateStoplossPositions(ctx, func(lpAddress sdk.AccAddress, stoploss types.Stoploss) bool {
		address := lpAddress.String()
		positions, found := positionsMap[address]
		if !found {
			positionsMap[address] = []types.Stoploss{stoploss}
			addresses = append(addresses, address)
		} else {
			positionsMap[address] = append(positions, stoploss)
		}
		return false
	})

	for _, addr := range addresses {
		lps = append(lps, types.StoplossPositions{
			Address:           addr,
			StoplossPositions: positionsMap[addr],
		})
	}

	return lps.Sort()
}

func (k Keeper) GetInvalidationID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.InvalidationIDPrefix)
	if len(bz) == 0 {
		return 0
	}

	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetInvalidationID(ctx sdk.Context, invalidationID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.InvalidationIDPrefix, sdk.Uint64ToBigEndian(invalidationID))
}

// DeleteSubmittedPosition removes an submitted stoploss position from the store.
func (k Keeper) DeleteSubmittedPosition(ctx sdk.Context, timeoutHeight uint64, address sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SubmittedPositionsQueuePrefix)
	store.Delete(types.SubmittedPositionKey(timeoutHeight, address))
}

// TODO: update, this assumes LP position per address

// SetSubmittedPosition sets a submitted stoplos position to the store.
func (k Keeper) SetSubmittedPosition(ctx sdk.Context, timeoutHeight uint64, address sdk.AccAddress, pairID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SubmittedPositionsQueuePrefix)
	store.Set(types.SubmittedPositionKey(timeoutHeight, address), []byte(pairID))
}

// IterateSubmittedQueue iterates over the all the submitted positions in ascending height performs a callback.
func (k Keeper) IterateSubmittedQueue(ctx sdk.Context, cb func(timeoutHeight uint64, address sdk.AccAddress, pairID string) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SubmittedPositionsQueuePrefix)

	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		timeoutHeight, address := types.SplitSubmittedStoplossKey(append(types.StoplossKeyPrefix, iterator.Key()...))
		pairID := string(iterator.Value())

		if cb(timeoutHeight, address, pairID) {
			break
		}
	}
}

// GetSubmittedQueue queries all the submitted positions.
func (k Keeper) GetSubmittedQueue(ctx sdk.Context) []types.SubmittedPosition {
	var submittedQueue []types.SubmittedPosition

	k.IterateSubmittedQueue(ctx, func(timeoutHeight uint64, address sdk.AccAddress, pairID string) bool {
		position := types.SubmittedPosition{
			Address:       address.String(),
			TimeoutHeight: timeoutHeight,
			PairId:        pairID,
		}

		submittedQueue = append(submittedQueue, position)
		return false
	})

	return submittedQueue
}
