package keeper

import (
	"fmt"
	"math/big"
	"sort"
	"strings"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v9/app/params"
	v1types "github.com/peggyjv/sommelier/v9/x/cellarfees/migrations/v1/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	paramSpace    paramtypes.Subspace
	accountKeeper v1types.AccountKeeper
	bankKeeper    v1types.BankKeeper
	mintKeeper    v1types.MintKeeper
	corkKeeper    v1types.CorkKeeper
	gravityKeeper v1types.GravityKeeper
	auctionKeeper v1types.AuctionKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper v1types.AccountKeeper,
	bankKeeper v1types.BankKeeper,
	mintKeeper v1types.MintKeeper,
	corkKeeper v1types.CorkKeeper,
	gravityKeeper v1types.GravityKeeper,
	auctionKeeper v1types.AuctionKeeper,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(v1types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    paramSpace,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		mintKeeper:    mintKeeper,
		corkKeeper:    corkKeeper,
		gravityKeeper: gravityKeeper,
		auctionKeeper: auctionKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", v1types.ModuleName))
}

////////////
// Params //
////////////

func (k Keeper) GetParams(ctx sdk.Context) v1types.Params {
	var p v1types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

func (k Keeper) SetParams(ctx sdk.Context, params v1types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

////////////////////////////////
// Last highest reward supply //
////////////////////////////////

func (k Keeper) GetLastRewardSupplyPeak(ctx sdk.Context) math.Int {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(v1types.GetLastRewardSupplyPeakKey())
	if b == nil {
		panic("Last highest reward supply should not have been nil")
	}
	var amount big.Int
	return sdk.NewIntFromBigInt((&amount).SetBytes(b))
}

func (k Keeper) SetLastRewardSupplyPeak(ctx sdk.Context, amount math.Int) {
	store := ctx.KVStore(k.storeKey)
	b := amount.BigInt().Bytes()
	store.Set(v1types.GetLastRewardSupplyPeakKey(), b)
}

//////////////////////////
// Fee accrual counters //
//////////////////////////

func (k Keeper) GetFeeAccrualCounters(ctx sdk.Context) (counters v1types.FeeAccrualCounters) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(v1types.GetFeeAccrualCountersKey())
	if b == nil {
		panic("Fee accrual counters is nil, it should have been set by InitGenesis")
	}
	if len(b) == 0 {
		return v1types.DefaultFeeAccrualCounters()
	}

	k.cdc.MustUnmarshal(b, &counters)
	return
}

func (k Keeper) SetFeeAccrualCounters(ctx sdk.Context, counters v1types.FeeAccrualCounters) {
	store := ctx.KVStore(k.storeKey)
	counterSlice := make([]v1types.FeeAccrualCounter, 0, len(counters.Counters))
	counterSlice = append(counterSlice, counters.Counters...)
	sort.Slice(counterSlice, func(i, j int) bool {
		return strings.Compare(counterSlice[i].Denom, counterSlice[j].Denom) == -1
	})
	counters.Counters = counterSlice
	b := k.cdc.MustMarshal(&counters)
	store.Set(v1types.GetFeeAccrualCountersKey(), b)
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
