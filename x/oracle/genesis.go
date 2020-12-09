package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/keeper"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	for delegatorBechAddr, delegateBechAddr := range data.FeederDelegations {
		delegator, err := sdk.ValAddressFromBech32(delegatorBechAddr)
		if err != nil {
			panic(err)
		}
		delegate, err := sdk.AccAddressFromBech32(delegateBechAddr)
		if err != nil {
			panic(err)
		}
		keeper.SetOracleDelegate(ctx, delegator, delegate)
	}

	for _, rate := range data.ExchangeRates {
		keeper.SetLunaExchangeRate(ctx, rate.Denom, rate.Amount)
	}

	for operatorBechAddr, missCounter := range data.MissCounters {
		operator, err := sdk.ValAddressFromBech32(operatorBechAddr)
		if err != nil {
			panic(err)
		}
		keeper.SetMissCounter(ctx, operator, missCounter)
	}

	for _, aggregatePrevote := range data.AggregateExchangeRatePrevotes {
		keeper.AddAggregateExchangeRatePrevote(ctx, aggregatePrevote)
	}

	for _, aggregateVote := range data.AggregateExchangeRateVotes {
		keeper.AddAggregateExchangeRateVote(ctx, aggregateVote)
	}

	if len(data.TobinTaxes) > 0 {
		for _, tt := range data.TobinTaxes {
			keeper.SetTobinTax(ctx, tt.Denom, tt.Amount)
		}
	} else {
		for _, item := range data.Params.Whitelist {
			keeper.SetTobinTax(ctx, item.Denom, item.Amount)
		}
	}

	keeper.SetParams(ctx, data.Params)

	// check if the module account exists
	moduleAcc := keeper.GetOracleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) (data types.GenesisState) {
	params := keeper.GetParams(ctx)
	feederDelegations := make(map[string]string)
	keeper.IterateOracleDelegates(ctx, func(delegator sdk.ValAddress, delegate sdk.AccAddress) (stop bool) {
		bechAddr := delegator.String()
		feederDelegations[bechAddr] = delegate.String()
		return false
	})

	rates := make(sdk.DecCoins, 0)
	keeper.IterateLunaExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		rates = append(rates, sdk.DecCoin{Amount: rate, Denom: denom})
		return false
	})

	missCounters := make(map[string]int64)
	keeper.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter int64) (stop bool) {
		missCounters[operator.String()] = missCounter
		return false
	})

	var aggregateExchangeRatePrevotes []types.AggregateExchangeRatePrevote
	keeper.IterateAggregateExchangeRatePrevotes(ctx, func(aggregatePrevote types.AggregateExchangeRatePrevote) (stop bool) {
		aggregateExchangeRatePrevotes = append(aggregateExchangeRatePrevotes, aggregatePrevote)
		return false
	})

	var aggregateExchangeRateVotes []types.AggregateExchangeRateVote
	keeper.IterateAggregateExchangeRateVotes(ctx, func(aggregateVote types.AggregateExchangeRateVote) bool {
		aggregateExchangeRateVotes = append(aggregateExchangeRateVotes, aggregateVote)
		return false
	})

	tobinTaxes := make(sdk.DecCoins, 0)
	keeper.IterateTobinTaxes(ctx, func(denom string, tobinTax sdk.Dec) (stop bool) {
		tobinTaxes = append(tobinTaxes, sdk.DecCoin{Amount: tobinTax, Denom: denom})
		return false
	})

	return types.NewGenesisState(params, rates, feederDelegations, missCounters, aggregateExchangeRatePrevotes, aggregateExchangeRateVotes, tobinTaxes)
}
