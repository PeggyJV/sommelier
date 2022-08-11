package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v4/app/params"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	paramSpace    paramtypes.Subspace
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    paramSpace,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

////////////
// Params //
////////////

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

/////////////////////
// Module Accounts //
/////////////////////

func (k Keeper) GetFeesAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

/////////////////////
// Cellar Fee Pool //
/////////////////////

func (k Keeper) GetCellarFeePool(ctx sdk.Context) (cellarFeePool types.CellarFeePool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.CellarFeePoolKey)
	if b == nil {
		panic("Stored cellar fee pool should not have been nil")
	}
	k.cdc.MustUnmarshal(b, &cellarFeePool)
	return
}

func (k Keeper) SetCellarFeePool(ctx sdk.Context, cellarFeePool types.CellarFeePool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&cellarFeePool)
	store.Set(types.CellarFeePoolKey, b)
}

// Appends the coin to the pool coins if denom isn't already present, otherwise add the amount to the
// existing balance
func (k Keeper) AddCoinToPool(ctx sdk.Context, coin sdk.Coin) {
	if (coin == sdk.Coin{}) {
		return
	}

	if coin.Denom == params.BaseCoinUnit || coin.Denom == params.HumanCoinUnit {
		panic("Cannot add SOMM to cellar fee pool")
	}

	pool := k.GetCellarFeePool(ctx)
	pool.Pool = pool.Pool.Add(coin)
	k.SetCellarFeePool(ctx, pool)
}

func (k Keeper) AddCoinsToPool(ctx sdk.Context, coins sdk.Coins) {
	if len(coins) == 0 {
		return
	}

	if !coins.AmountOfNoDenomValidation(params.BaseCoinUnit).IsZero() {
		panic("Cannot add SOMM to cellar fee pool")
	}

	pool := k.GetCellarFeePool(ctx)
	pool.Pool = pool.Pool.Add(coins.Sort()...)
	k.SetCellarFeePool(ctx, pool)
}

func (k Keeper) SendPoolToAuction(ctx sdk.Context) {
	pool := k.GetCellarFeePool(ctx).Pool
	if pool.Empty() {
		return
	}

	// TO-DO: Update when auction module exists. Test setup creates this mock auction module account.
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, "auction", pool)
	if err != nil {
		panic(err)
	}

	// reset pool
	k.SetCellarFeePool(ctx, types.NewEmptyPool())
}
