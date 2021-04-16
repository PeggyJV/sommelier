package keeper

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
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

	prevote, found := k.GetAllocationPrecommit(ctx, validatorAddr)
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

	dataVote, found := k.GetAllocationCommit(ctx, validatorAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "data vote")
	}

	return &types.QueryOracleDataVoteResponse{
		Vote: &dataVote,
	}, nil
}

// QueryAggregateData implements QueryServer
func (k Keeper) QueryAggregateData(c context.Context, req *types.QueryAggregateDataRequest) (*types.QueryAggregateDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if strings.TrimSpace(req.Id) == "" {
		return nil, status.Error(codes.InvalidArgument, "empty oracle data id")
	}

	ctx := sdk.UnwrapSDKContext(c)

	var ok bool

	if req.Type == "" {
		req.Type, ok = k.GetOracleDataType(ctx, req.Id)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "empty oracle data type")
		}
	}

	oracleData, height := k.GetLatestAggregatedOracleData(ctx, req.Type, req.Id)
	if oracleData == nil {
		return nil, status.Errorf(codes.NotFound, "aggregated oracle data with id %s", req.Id)
	}

	pair, ok := oracleData.(*types.UniswapPair)
	if !ok {
		return nil, status.Errorf(codes.Internal, "data type is not %s", types.UniswapDataType)
	}

	return &types.QueryAggregateDataResponse{
		OracleData: pair,
		Height:     height,
	}, nil
}

// QueryLatestPeriodAggregateData implements QueryServer
func (k Keeper) QueryLatestPeriodAggregateData(c context.Context, req *types.QueryLatestPeriodAggregateDataRequest) (*types.QueryLatestPeriodAggregateDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.AggregatedOracleDataKeyPrefix)

	var (
		pairs  = []*types.UniswapPair{}
		height int64
	)

	// TODO: this should use reverse iterator
	// Ref: https://github.com/cosmos/cosmos-sdk/issues/8754
	pageRes, err := query.FilteredPaginate(store, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		// check if
		dataHeight := int64(sdk.BigEndianToUint64(key[1:9]))
		if height == 0 {
			height = dataHeight
		} else if dataHeight != height {
			// do not count the current value in the pagination response
			return false, nil
		}

		var oracleData types.OracleData
		if err := k.cdc.UnmarshalInterface(value, &oracleData); err != nil {
			return false, err
		}

		pair, ok := oracleData.(*types.UniswapPair)
		if !ok {
			return false, fmt.Errorf("data type is not %s", types.UniswapDataType)
		}

		if accumulate {
			pairs = append(pairs, pair)
		}
		// count the current value in the pagination response
		return true, nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryLatestPeriodAggregateDataResponse{
		OracleData: pairs,
		Height:     height,
		Pagination: pageRes,
	}, nil
}

// QueryParams implements QueryServer
func (k Keeper) QueryParams(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{
		Params: k.GetParamSet(sdk.UnwrapSDKContext(c)),
	}, nil
}

// QueryVotePeriod implements QueryServer
func (k Keeper) QueryVotePeriod(c context.Context, _ *types.QueryVotePeriodRequest) (*types.QueryVotePeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	votePeriodStart, found := k.GetCommitPeriodStart(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "vote period start not set")
	}

	return &types.QueryVotePeriodResponse{
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
