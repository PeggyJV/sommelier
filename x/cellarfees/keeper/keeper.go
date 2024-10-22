package keeper

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v8/app/params"
	cellarfeestypes "github.com/peggyjv/sommelier/v8/x/cellarfees/types"
	types "github.com/peggyjv/sommelier/v8/x/cellarfees/types/v2"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	paramSpace    paramtypes.Subspace
	accountKeeper cellarfeestypes.AccountKeeper
	bankKeeper    cellarfeestypes.BankKeeper
	mintKeeper    cellarfeestypes.MintKeeper
	corkKeeper    cellarfeestypes.CorkKeeper
	auctionKeeper cellarfeestypes.AuctionKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper cellarfeestypes.AccountKeeper,
	bankKeeper cellarfeestypes.BankKeeper,
	mintKeeper cellarfeestypes.MintKeeper,
	corkKeeper cellarfeestypes.CorkKeeper,
	auctionKeeper cellarfeestypes.AuctionKeeper,
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
		mintKeeper:    mintKeeper,
		corkKeeper:    corkKeeper,
		auctionKeeper: auctionKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", cellarfeestypes.ModuleName))
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

////////////////////////////////
// Last highest reward supply //
////////////////////////////////

func (k Keeper) GetLastRewardSupplyPeak(ctx sdk.Context) math.Int {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(cellarfeestypes.GetLastRewardSupplyPeakKey())
	if b == nil {
		panic("Last highest reward supply should not have been nil")
	}
	var amount big.Int
	return sdk.NewIntFromBigInt((&amount).SetBytes(b))
}

func (k Keeper) SetLastRewardSupplyPeak(ctx sdk.Context, amount math.Int) {
	store := ctx.KVStore(k.storeKey)
	b := amount.BigInt().Bytes()
	store.Set(cellarfeestypes.GetLastRewardSupplyPeakKey(), b)
}

////////////
// APY    //
////////////

func (k Keeper) GetAPY(ctx sdk.Context) sdk.Dec {
	remainingRewardsSupply := k.bankKeeper.GetBalance(ctx, k.GetFeesAccount(ctx).GetAddress(), params.BaseCoinUnit).Amount
	if remainingRewardsSupply.IsZero() {
		return sdk.ZeroDec()
	}

	mintParams := k.mintKeeper.GetParams(ctx)
	bondedRatio := k.mintKeeper.BondedRatio(ctx)
	totalCoins := k.mintKeeper.StakingTokenSupply(ctx)
	emission := k.GetEmission(ctx, remainingRewardsSupply)
	annualRewards := emission.AmountOf(params.BaseCoinUnit).Mul(sdk.NewInt(int64(mintParams.BlocksPerYear)))

	return sdk.NewDecFromInt(annualRewards).Quo(sdk.NewDecFromInt(totalCoins)).Quo(bondedRatio)
}
