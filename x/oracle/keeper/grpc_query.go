package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

var _ types.QueryServer = Keeper{}

// ExchangeRate returns the current exchange rate for a given denom
func (k Keeper) ExchangeRate(c context.Context, req *types.QueryExchangeRateRequest) (*types.QueryExchangeRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	rate, err := k.GetUSDExchangeRate(sdk.UnwrapSDKContext(c), req.Denom)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnknowDenom, req.Denom)
	}

	return &types.QueryExchangeRateResponse{Rate: rate}, nil
}

// ExchangeRates returns all the exchange rates tracked by the system
func (k Keeper) ExchangeRates(c context.Context, req *types.QueryExchangeRatesRequest) (*types.QueryExchangeRatesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	rates := k.GetUSDExchangeRates(sdk.UnwrapSDKContext(c))
	return &types.QueryExchangeRatesResponse{Rates: rates}, nil
}

// Actives returns the active denoms
func (k Keeper) Actives(c context.Context, req *types.QueryActivesRequest) (*types.QueryActivesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	denoms := []string{}
	k.IterateUSDExchangeRates(sdk.UnwrapSDKContext(c), func(denom string, _ sdk.Dec) (stop bool) {
		denoms = append(denoms, denom)
		return false
	})

	return &types.QueryActivesResponse{Denoms: denoms}, nil
}

// Parameters returns the oracle module parameters
func (k Keeper) Parameters(c context.Context, req *types.QueryParametersRequest) (*types.QueryParametersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &types.QueryParametersResponse{Params: k.GetParams(sdk.UnwrapSDKContext(c))}, nil
}

// FeederDelegation returns the address to which a validator is delegating the feeder responsibility
func (k Keeper) FeederDelegation(c context.Context, req *types.QueryFeederDelegationRequest) (*types.QueryFeederDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &types.QueryFeederDelegationResponse{
		Address: k.GetOracleDelegate(sdk.UnwrapSDKContext(c), addr).String(),
	}, nil
}

// MissCounter returns the number of misses for a given validator
func (k Keeper) MissCounter(c context.Context, req *types.QueryMissCounterRequest) (*types.QueryMissCounterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &types.QueryMissCounterResponse{Counter: k.GetMissCounter(sdk.UnwrapSDKContext(c), addr)}, nil
}

// AggregatePrevote returns the latest aggregate prevote from a given validator
func (k Keeper) AggregatePrevote(c context.Context, req *types.QueryAggregatePrevoteRequest) (*types.QueryAggregatePrevoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	prevote, err := k.GetAggregateExchangeRatePrevote(sdk.UnwrapSDKContext(c), addr)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryAggregatePrevoteResponse{Prevote: prevote}, nil
}

// AggregateVote returns the latest aggregate prevote from a given validator
func (k Keeper) AggregateVote(c context.Context, req *types.QueryAggregateVoteRequest) (*types.QueryAggregateVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	addr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	vote, err := k.GetAggregateExchangeRateVote(sdk.UnwrapSDKContext(c), addr)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryAggregateVoteResponse{Vote: vote}, nil
}

// VoteTargets returns the target denoms for voting?
func (k Keeper) VoteTargets(c context.Context, req *types.QueryVoteTargetsRequest) (*types.QueryVoteTargetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &types.QueryVoteTargetsResponse{
		Targets: k.GetVoteTargets(sdk.UnwrapSDKContext(c)),
	}, nil
}

// TobinTax returns the current tobin tax given a denom
func (k Keeper) TobinTax(c context.Context, req *types.QueryTobinTaxRequest) (*types.QueryTobinTaxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	tobinTax, err := k.GetTobinTax(sdk.UnwrapSDKContext(c), req.Denom)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &types.QueryTobinTaxResponse{Rate: tobinTax}, nil
}

// TobinTaxes returns all the tobin taxes tracked by the system
func (k Keeper) TobinTaxes(c context.Context, req *types.QueryTobinTaxesRequest) (*types.QueryTobinTaxesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	denoms := k.GetTobinTaxes(sdk.UnwrapSDKContext(c))
	return &types.QueryTobinTaxesResponse{Rates: denoms}, nil
}
