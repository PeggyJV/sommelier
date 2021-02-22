package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
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

// QueryOracleDataPrevote implements QueryServer
func (k Keeper) QueryOracleDataPrevote(c context.Context, req *types.QueryOracleDataPrevoteRequest) (*types.QueryOracleDataPrevoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	validatorAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	prevote, found := k.GetOracleDataPrevote(ctx, validatorAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "data prevote")
	}

	return &types.QueryOracleDataPrevoteResponse{
		Prevote: &prevote,
	}, nil
}

// QueryOracleDataVote implements QueryServer
func (k Keeper) QueryOracleDataVote(c context.Context, req *types.QueryOracleDataVoteRequest) (*types.QueryOracleDataVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	validatorAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	dataVote, found := k.GetOracleDataVote(ctx, validatorAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "data vote")
	}

	return &types.QueryOracleDataVoteResponse{
		Vote: &dataVote,
	}, nil
}

// OracleData implements QueryServer
func (k Keeper) OracleData(c context.Context, req *types.QueryOracleDataRequest) (*types.QueryOracleDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	oracleData, found := k.GetOracleData(ctx, req.Type, "") // TODO: query by id
	if !found {
		return nil, status.Errorf(codes.NotFound, "data type %s", req.Type)
	}

	oracleDataAny, err := types.PackOracleData(oracleData)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryOracleDataResponse{
		OracleData: oracleDataAny,
	}, nil
}

// QueryParams implements QueryServer
func (k Keeper) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{
		Params: k.GetParamSet(sdk.UnwrapSDKContext(c)),
	}, nil
}

// QueryVotePeriod implements QueryServer
func (k Keeper) QueryVotePeriod(c context.Context, _ *types.QueryVotePeriodRequest) (*types.VotePeriod, error) {
	ctx := sdk.UnwrapSDKContext(c)
	votePeriodStart, found := k.GetVotePeriodStart(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "vote period start not set")
	}

	return &types.VotePeriod{
		VotePeriodStart: votePeriodStart,
		VotePeriodEnd:   votePeriodStart + k.GetParamSet(ctx).VotePeriod,
		CurrentHeight:   ctx.BlockHeight(),
	}, nil
}

// QueryMissCounter implements QueryServer
func (k Keeper) QueryMissCounter(c context.Context, req *types.QueryMissCounterRequest) (*types.QueryMissCounterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	validatorAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &types.QueryMissCounterResponse{
		MissCounter: k.GetMissCounter(ctx, validatorAddr),
	}, nil
}
