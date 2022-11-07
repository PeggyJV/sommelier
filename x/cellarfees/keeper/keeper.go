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

////////////////////////////////
// Last highest reward supply //
////////////////////////////////

func (k Keeper) GetLastRewardSupplyPeak(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetLastRewardSupplyPeakKey())
	if b == nil {
		panic("Last highest reward supply should not have been nil")
	}
	var amount big.Int
	return sdk.NewIntFromBigInt((&amount).SetBytes(b))
}

func (k Keeper) SetLastRewardSupplyPeak(ctx sdk.Context, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := amount.BigInt().Bytes()
	store.Set(types.GetLastRewardSupplyPeakKey(), b)
}

//////////////////////////
// Fee accrual counters //
//////////////////////////

func (k Keeper) GetFeeAccrualCounters(ctx sdk.Context) (counters types.FeeAccrualCounters) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetFeeAccrualCountersKey())
	if b == nil {
		panic("Fee accrual counters is nil, it should have been set by InitGenesis")
	}
	if len(b) == 0 {
		return types.DefaultFeeAccrualCounters()
	}

	k.cdc.MustUnmarshal(b, &counters)
	return
}

func (k Keeper) SetFeeAccrualCounters(ctx sdk.Context, counters types.FeeAccrualCounters) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&counters)
	store.Set(types.GetFeeAccrualCountersKey(), b)
}
