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
	store.Set(types.StoplossKey(address, stoploss.UniswapPairId), bz)
}

// DeleteStoplossPosition removes a stoploss position from the store.
func (k Keeper) DeleteStoplossPosition(ctx sdk.Context, address sdk.AccAddress, uniswapPair string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StoplossKeyPrefix)
	store.Delete(types.StoplossKey(address, uniswapPair))
}
