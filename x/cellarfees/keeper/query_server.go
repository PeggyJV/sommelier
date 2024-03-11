package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/peggyjv/sommelier/v7/x/cellarfees/types/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) QueryParams(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &types.QueryParamsResponse{
		Params: k.GetParams(sdk.UnwrapSDKContext(c)),
	}, nil
}

func (k Keeper) QueryModuleAccounts(c context.Context, req *types.QueryModuleAccountsRequest) (*types.QueryModuleAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	return &types.QueryModuleAccountsResponse{
		FeesAddress: k.GetFeesAccount(sdk.UnwrapSDKContext(c)).GetAddress().String(),
	}, nil
}

func (k Keeper) QueryLastRewardSupplyPeak(c context.Context, req *types.QueryLastRewardSupplyPeakRequest) (*types.QueryLastRewardSupplyPeakResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	return &types.QueryLastRewardSupplyPeakResponse{LastRewardSupplyPeak: k.GetLastRewardSupplyPeak(sdk.UnwrapSDKContext(c))}, nil
}

func (k Keeper) QueryAPY(c context.Context, _ *types.QueryAPYRequest) (*types.QueryAPYResponse, error) {
	return &types.QueryAPYResponse{
		Apy: k.GetAPY(sdk.UnwrapSDKContext(c)).String(),
	}, nil
}

func (k Keeper) QueryFeeTokenBalance(c context.Context, req *types.QueryFeeTokenBalanceRequest) (*types.QueryFeeTokenBalanceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Denom == "" {
		return nil, status.Error(codes.InvalidArgument, "denom cannot be empty")
	}

	balance, found := k.GetFeeBalance(ctx, req.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, "fee token balance not found")
	}

	tokenPrice, found := k.auctionKeeper.GetTokenPrice(ctx, req.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, "token price not found")
	}

	totalUsdValue, err := k.GetBalanceUsdValue(ctx, balance, &tokenPrice).Float64()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to convert usd value to float")
	}

	feeTokenBalance := types.FeeTokenBalance{
		Balance:  balance,
		UsdValue: totalUsdValue,
	}

	return &types.QueryFeeTokenBalanceResponse{
		Balance: &feeTokenBalance,
	}, nil
}

func (k Keeper) QueryFeeTokenBalances(c context.Context, _ *types.QueryFeeTokenBalancesRequest) (*types.QueryFeeTokenBalancesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	feeBalances := make([]*types.FeeTokenBalance, 0)

	// because we can't get a USD value without a corresponding TokenPrice set in the auction module,
	// this exclude fee token balances that don't have one yet.
	tokenPrices := k.auctionKeeper.GetTokenPrices(ctx)
	for _, tokenPrice := range tokenPrices {
		balance, found := k.GetFeeBalance(ctx, tokenPrice.Denom)
		if !found {
			continue
		}

		if balance.IsZero() {
			continue
		}

		totalUsdValue, err := k.GetBalanceUsdValue(ctx, balance, tokenPrice).Float64()
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to convert usd value to float")
		}

		feeTokenBalance := types.FeeTokenBalance{
			Balance:  balance,
			UsdValue: totalUsdValue,
		}

		feeBalances = append(feeBalances, &feeTokenBalance)
	}

	return &types.QueryFeeTokenBalancesResponse{
		Balances: feeBalances,
	}, nil
}
