package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v6/x/incentives/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the incentives store
type Keeper struct {
	storeKey           sdk.StoreKey
	cdc                codec.BinaryCodec
	paramSpace         paramtypes.Subspace
	DistributionKeeper types.DistributionKeeper
	BankKeeper         types.BankKeeper
	MintKeeper         types.MintKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	distributionKeeper types.DistributionKeeper,
	bankKeeper types.BankKeeper,
	mintKeeper types.MintKeeper,
) Keeper {
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:           storeKey,
		cdc:                cdc,
		paramSpace:         paramSpace,
		DistributionKeeper: distributionKeeper,
		BankKeeper:         bankKeeper,
		MintKeeper:         mintKeeper,
	}
}

////////////
// Params //
////////////

// GetParamSet returns the parameters
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// setParams sets the parameters in the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

////////////
// APY    //
////////////

func (k Keeper) GetAPY(ctx sdk.Context) sdk.Dec {
	incentivesParams := k.GetParamSet(ctx)
	// check if incentives are enabled
	if uint64(ctx.BlockHeight()) >= incentivesParams.IncentivesCutoffHeight || incentivesParams.DistributionPerBlock.IsZero() {
		return sdk.ZeroDec()
	}

	mintParams := k.MintKeeper.GetParams(ctx)
	bondedRatio := k.MintKeeper.BondedRatio(ctx)
	totalCoins := k.MintKeeper.StakingTokenSupply(ctx)
	annualRewards := incentivesParams.DistributionPerBlock.Amount.Mul(sdk.NewInt(int64(mintParams.BlocksPerYear)))

	return annualRewards.ToDec().Quo(totalCoins.ToDec()).Quo(bondedRatio)
}
