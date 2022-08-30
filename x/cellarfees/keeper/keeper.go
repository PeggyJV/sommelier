package keeper

import (
	"fmt"
	"math/big"

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
	corkKeeper    types.CorkKeeper
	gravityKeeper types.GravityKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	corkKeeper types.CorkKeeper,
	gravityKeeper types.GravityKeeper,
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
		corkKeeper:    corkKeeper,
		gravityKeeper: gravityKeeper,
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

func (k Keeper) HandleAuctions(ctx sdk.Context) {
	pool := k.GetCellarFeePool(ctx).Pool
	if pool.Empty() {
		// Schedule next auction a short delay away until an auction occurs
		k.SetScheduledAuctionHeight(ctx, k.GetScheduledAuctionHeight(ctx).Add(sdk.NewInt(100)))
		return
	}

	// Get any active auctions

	// For each denom in pool, start an auction if there isn't an active one

	// Remove denoms from pool that auctions were started for

	k.ScheduleNextAuction(ctx)
}

////////////////////////////////
// Last highest reward supply //
////////////////////////////////

func (k Keeper) GetLastRewardSupplyPeak(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastHighestRewardSupply)
	if b == nil {
		panic("Last highest reward supply should not have been nil")
	}
	var amount big.Int
	return sdk.NewIntFromBigInt((&amount).SetBytes(b))
}

func (k Keeper) SetLastRewardSupplyPeak(ctx sdk.Context, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := amount.BigInt().Bytes()
	store.Set(types.LastHighestRewardSupply, b)
}

////////////////////////
// Auction scheduling //
////////////////////////

func (k Keeper) GetScheduledAuctionHeight(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.ScheduledAuctionHeight)
	if b == nil {
		panic("Auction height should not have been nil")
	}
	var amount big.Int
	return sdk.NewIntFromBigInt((&amount).SetBytes(b))
}

func (k Keeper) SetScheduledAuctionHeight(ctx sdk.Context, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := amount.BigInt().Bytes()
	store.Set(types.ScheduledAuctionHeight, b)
}

func (k Keeper) ScheduleNextAuction(ctx sdk.Context) {
	lastAuctionHeight := k.GetScheduledAuctionHeight(ctx)
	// next = last + delay param
	nextAuctionHeight := lastAuctionHeight.Add(sdk.NewInt(int64(k.GetParams(ctx).AuctionBlockDelay)))
	k.SetScheduledAuctionHeight(ctx, nextAuctionHeight)
}
