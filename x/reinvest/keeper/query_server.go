package keeper

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/reinvest/types"
)

var _ types.QueryServer = Keeper{}

// QuerySubmittedReinvestments implements QueryServer
func (k Keeper) QuerySubmittedReinvestments(c context.Context, _ *types.QuerySubmittedReinvestmentsRequest) (*types.QuerySubmittedReinvestmentsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	// get reinvests
	var reinvests []*types.Reinvestment
	k.IterateReinvestments(ctx, func(_ sdk.ValAddress, _ common.Address, reinvestment types.Reinvestment) (stop bool) {
		reinvests = append(reinvests, &reinvestment)
		return false
	})

	return &types.QuerySubmittedReinvestmentsResponse{
		Reinvests: reinvests,
	}, nil
}

// QueryParams implements QueryServer
func (k Keeper) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{
		Params: k.GetParamSet(sdk.UnwrapSDKContext(c)),
	}, nil
}

// QueryCommitPeriod implements QueryServer
func (k Keeper) QueryCommitPeriod(c context.Context, _ *types.QueryCommitPeriodRequest) (*types.QueryCommitPeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	votePeriodStart, found := k.GetCommitPeriodStart(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "vote period start not set")
	}

	return &types.QueryCommitPeriodResponse{
		VotePeriodStart: votePeriodStart,
		VotePeriodEnd:   votePeriodStart + k.GetParamSet(ctx).VotePeriod,
		CurrentHeight:   ctx.BlockHeight(),
	}, nil
}