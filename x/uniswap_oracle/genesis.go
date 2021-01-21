package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/uniswap_oracle/keeper"
	"github.com/peggyjv/sommelier/x/uniswap_oracle/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	for _, delegation := range data.FeederDelegations {
		// error checked during genesis validation
		delegator, _ := sdk.ValAddressFromBech32(delegation.DelegatorAddress)
		delegate, _ := sdk.AccAddressFromBech32(delegation.DelegateAddress)
		keeper.SetOracleDelegate(ctx, delegator, delegate)
	}

	for _, rate := range data.ExchangeRates {
		keeper.SetUSDExchangeRate(ctx, rate.Denom, rate.Amount)
	}

	for _, missCounter := range data.MissCounters {
		// error checked during genesis validation
		operator, _ := sdk.ValAddressFromBech32(missCounter.ValAddress)
		keeper.SetMissCounter(ctx, operator, missCounter.MissedCounter)
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
	rates := keeper.GetUSDExchangeRates(ctx)
	feederDelegations := keeper.GetOracleDelegations(ctx)
	missCounters := keeper.GetMissCounters(ctx)
	aggregateExchangeRatePrevotes := keeper.GetAggregateExchangeRatePrevotes(ctx)
	aggregateExchangeRateVotes := keeper.GetAggregateExchangeRateVotes(ctx)
	tobinTaxes := keeper.GetTobinTaxes(ctx)

	return types.NewGenesisState(params, rates, feederDelegations, missCounters, aggregateExchangeRatePrevotes, aggregateExchangeRateVotes, tobinTaxes)
}
