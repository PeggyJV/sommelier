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

// QueryDelegateAddress implements QueryServer
func (k Keeper) QueryDelegateAddress(c context.Context, req *types.QueryDelegateAddressRequest) (*types.QueryDelegateAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	val, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	delegateAddr := k.GetDelegateAddressFromValidator(ctx, val)
	if delegateAddr == nil {
		return nil, status.Errorf(
			codes.NotFound, "delegator address for validator %s", req.Validator,
		)
	}

	return &types.QueryDelegateAddressResponse{
		Delegate: delegateAddr.String(),
	}, nil
}

// QueryValidatorAddress implements QueryServer
func (k Keeper) QueryValidatorAddress(c context.Context, req *types.QueryValidatorAddressRequest) (*types.QueryValidatorAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	del, err := sdk.AccAddressFromBech32(req.Delegate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	validatorAddr := k.GetValidatorAddressFromDelegate(ctx, del)
	if validatorAddr == nil {
		return nil, status.Errorf(
			codes.NotFound, "delegator address for delegate %s", req.Delegate,
		)
	}

	return &types.QueryValidatorAddressResponse{
		Validator: validatorAddr.String(),
	}, nil
}

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
