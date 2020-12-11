package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

var _ types.QueryServer = Keeper{}

// ExchangeRate returns the current exchange rate for a given denom
func (k Keeper) ExchangeRate(c context.Context, req *types.QueryExchangeRateRequest) (*types.QueryExchangeRateResponse, error) {
	rate, err := k.GetUSDExchangeRate(sdk.UnwrapSDKContext(c), req.Denom)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnknowDenom, req.Denom)
	}

	return &types.QueryExchangeRateResponse{Rate: rate}, nil
}

// ExchangeRates returns all the exchange rates tracked by the system
func (k Keeper) ExchangeRates(c context.Context, req *types.QueryExchangeRatesRequest) (*types.QueryExchangeRatesResponse, error) {
	var rates sdk.DecCoins
	k.IterateUSDExchangeRates(sdk.UnwrapSDKContext(c), func(denom string, rate sdk.Dec) (stop bool) {
		rates = append(rates, sdk.NewDecCoinFromDec(denom, rate))
		return false
	})
	return &types.QueryExchangeRatesResponse{Rates: rates}, nil
}

// Actives returns the active denoms
func (k Keeper) Actives(c context.Context, req *types.QueryActivesRequest) (*types.QueryActivesResponse, error) {
	denoms := []string{}
	k.IterateUSDExchangeRates(sdk.UnwrapSDKContext(c), func(denom string, rate sdk.Dec) (stop bool) {
		denoms = append(denoms, denom)
		return false
	})
	return &types.QueryActivesResponse{Denoms: denoms}, nil
}

// Parameters returns the oracle module parameters
func (k Keeper) Parameters(c context.Context, req *types.QueryParametersRequest) (*types.QueryParametersResponse, error) {
	return &types.QueryParametersResponse{Params: k.GetParams(sdk.UnwrapSDKContext(c))}, nil
}

// FeederDelegation returns the address to which a validator is delegating the feeder responsibility
func (k Keeper) FeederDelegation(c context.Context, req *types.QueryFeederDelegationRequest) (*types.QueryFeederDelegationResponse, error) {
	addr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "must give valid validator address")
	}
	return &types.QueryFeederDelegationResponse{Address: k.GetOracleDelegate(sdk.UnwrapSDKContext(c), addr).String()}, nil
}

// MissCounter returns the number of misses for a given validator
func (k Keeper) MissCounter(c context.Context, req *types.QueryMissCounterRequest) (*types.QueryMissCounterResponse, error) {
	addr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "must give valid validator address")
	}
	return &types.QueryMissCounterResponse{Counter: k.GetMissCounter(sdk.UnwrapSDKContext(c), addr)}, nil
}

// AggregatePrevote returns the latest aggregate prevote from a given validator
func (k Keeper) AggregatePrevote(c context.Context, req *types.QueryAggregatePrevoteRequest) (*types.QueryAggregatePrevoteResponse, error) {
	addr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "must give valid validator address")
	}
	prevote, err := k.GetAggregateExchangeRatePrevote(sdk.UnwrapSDKContext(c), addr)
	if err != nil {
		return nil, err
	}
	return &types.QueryAggregatePrevoteResponse{Prevote: prevote}, nil
}

// AggregateVote returns the latest aggregate prevote from a given validator
func (k Keeper) AggregateVote(c context.Context, req *types.QueryAggregateVoteRequest) (*types.QueryAggregateVoteResponse, error) {
	addr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "must give valid validator address")
	}
	vote, err := k.GetAggregateExchangeRateVote(sdk.UnwrapSDKContext(c), addr)
	if err != nil {
		return nil, err
	}
	return &types.QueryAggregateVoteResponse{Vote: vote}, nil
}

// VoteTargets returns the target denoms for voting?
func (k Keeper) VoteTargets(c context.Context, req *types.QueryVoteTargetsRequest) (*types.QueryVoteTargetsResponse, error) {
	return &types.QueryVoteTargetsResponse{Targets: k.GetVoteTargets(sdk.UnwrapSDKContext(c))}, nil
}

// TobinTax returns the current tobin tax given a denom
func (k Keeper) TobinTax(c context.Context, req *types.QueryTobinTaxRequest) (*types.QueryTobinTaxResponse, error) {
	tobinTax, err := k.GetTobinTax(sdk.UnwrapSDKContext(c), req.Denom)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid denom")
	}
	return &types.QueryTobinTaxResponse{Rate: tobinTax}, nil
}

// TobinTaxes returns all the tobin taxes tracked by the system
func (k Keeper) TobinTaxes(c context.Context, req *types.QueryTobinTaxesRequest) (*types.QueryTobinTaxesResponse, error) {
	var denoms sdk.DecCoins
	k.IterateTobinTaxes(sdk.UnwrapSDKContext(c), func(denom string, tobinTax sdk.Dec) (stop bool) {
		denoms = append(denoms, sdk.NewDecCoinFromDec(denom, tobinTax))
		return false
	})
	return &types.QueryTobinTaxesResponse{Rates: denoms}, nil
}
