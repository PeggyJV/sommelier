package keeper

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/allocation/types"
)

var _ types.QueryServer = Keeper{}

// QueryAllocationPrecommit implements QueryServer
func (k Keeper) QueryAllocationPrecommit(c context.Context, req *types.QueryAllocationPrecommitRequest) (*types.QueryAllocationPrecommitResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	validatorAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	precommit, found := k.GetAllocationPrecommit(ctx, validatorAddr, common.HexToAddress(req.Cellar))
	if !found {
		return nil, status.Error(codes.NotFound, "data precommit")
	}

	return &types.QueryAllocationPrecommitResponse{
		Precommit: &precommit,
	}, nil
}

// QueryAllocationPrecommits implements QueryServer
func (k Keeper) QueryAllocationPrecommits(c context.Context, _ *types.QueryAllocationPrecommitsRequest) (*types.QueryAllocationPrecommitsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var precommits []*types.AllocationPrecommit
	k.IterateAllocationPrecommits(ctx, func(val sdk.ValAddress, cel common.Address, precommit types.AllocationPrecommit) (stop bool) {
		precommits = append(precommits, &precommit)
		return false
	})

	return &types.QueryAllocationPrecommitsResponse{
		Precommits: precommits,
	}, nil
}

// QueryAllocationCommit implements QueryServer
func (k Keeper) QueryAllocationCommit(c context.Context, req *types.QueryAllocationCommitRequest) (*types.QueryAllocationCommitResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	validatorAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	allocationCommit, found := k.GetAllocationCommit(ctx, validatorAddr, common.HexToAddress(req.Cellar))
	if !found {
		return nil, status.Error(codes.NotFound, "data vote")
	}

	return &types.QueryAllocationCommitResponse{
		Commit: &allocationCommit,
	}, nil
}

// QueryAllocationCommits implements QueryServer
func (k Keeper) QueryAllocationCommits(c context.Context, _ *types.QueryAllocationCommitsRequest) (*types.QueryAllocationCommitsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var commits []*types.Allocation
	k.IterateAllocationCommits(ctx, func(val sdk.ValAddress, cel common.Address, commit types.Allocation) (stop bool) {
		commits = append(commits, &commit)
		return false
	})


	return &types.QueryAllocationCommitsResponse{
		Commits: commits,
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

func (k Keeper) QueryCellars(c context.Context, _ *types.QueryCellarsRequest) (*types.QueryCellarsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var cellars []*types.Cellar
	k.IterateCellars(ctx, func(cellar types.Cellar) (stop bool) {
		cellars = append(cellars, &cellar)
		return false
	})

	return &types.QueryCellarsResponse{Cellars: cellars}, nil
}
