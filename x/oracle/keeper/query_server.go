package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

var _ types.QueryServer = Keeper{}

// QueryDelegeateAddress implements QueryServer
func (k Keeper) QueryDelegeateAddress(c context.Context, req *types.QueryDelegeateAddressRequest) (*types.QueryDelegeateAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	val, err := sdk.AccAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	delegateAddr := k.GetDelegateAddressFromValidator(ctx, val)
	if delegateAddr == nil {
		return nil, status.Errorf(
			codes.NotFound, "delegator address for validator %s", req.Validator,
		)
	}

	return &types.QueryDelegeateAddressResponse{
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

	val, err := sdk.AccAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if k.stakingKeeper.Validator(ctx, sdk.ValAddress(val)) == nil {
		val = k.GetValidatorAddressFromDelegate(ctx, val)
		if val == nil {
			return nil, status.Errorf(codes.NotFound, "address %s is not a validator", req.Validator)
		}
	}

	dataPrevote := k.GetOracleDataPrevote(ctx, val)
	if dataPrevote == nil {
		return nil, status.Error(codes.NotFound, "data prevote")
	}

	return &types.QueryOracleDataPrevoteResponse{
		Hashes: dataPrevote.Hashes,
	}, nil
}

// QueryOracleDataVote implements QueryServer
func (k Keeper) QueryOracleDataVote(c context.Context, req *types.QueryOracleDataVoteRequest) (*types.MsgOracleDataVote, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	val, err := sdk.AccAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if k.stakingKeeper.Validator(ctx, sdk.ValAddress(val)) == nil {
		val = k.GetValidatorAddressFromDelegate(ctx, val)
		if val == nil {
			return nil, status.Errorf(codes.NotFound, "address %s is not a validator", req.Validator)
		}
	}

	dataVote := k.GetOracleDataVote(ctx, val)
	if dataVote == nil {
		return nil, status.Error(codes.NotFound, "data vote")
	}

	return dataVote, nil
}

// OracleData implements QueryServer
func (k Keeper) OracleData(c context.Context, req *types.QueryOracleDataRequest) (*types.QueryOracleDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	oracleData := k.GetOracleData(ctx, req.Type)
	if oracleData == nil {
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
	vps := k.GetVotePeriodStart(ctx)

	return &types.VotePeriod{
		VotePeriodStart: vps,
		VotePeriodEnd:   vps + k.GetParamSet(ctx).VotePeriod,
		CurrentHeight:   ctx.BlockHeight(),
	}, nil
}

// QueryMissCounter implements QueryServer
func (k Keeper) QueryMissCounter(c context.Context, req *types.QueryMissCounterRequest) (*types.QueryMissCounterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.AccAddressFromBech32(req.Validator)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if k.stakingKeeper.Validator(ctx, sdk.ValAddress(val)) == nil {
		val = k.GetValidatorAddressFromDelegate(ctx, val)
		if val == nil {
			return nil, status.Errorf(codes.NotFound, "address %s is not a validator", req.Validator)
		}
	}

	return &types.QueryMissCounterResponse{
		MissCounter: k.GetMissCounter(ctx, val),
	}, nil
}
