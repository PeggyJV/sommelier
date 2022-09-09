package keeper

import (
	"fmt"
	"math/big"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
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
	auctionKeeper types.AuctionKeeper
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
