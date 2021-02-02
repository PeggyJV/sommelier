package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

var _ types.QueryServer = Keeper{}

// QueryDelegeateAddress implements QueryServer
func (k Keeper) QueryDelegeateAddress(c context.Context, req *types.QueryDelegeateAddressRequest) (*types.QueryDelegeateAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.AccAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}
	return &types.QueryDelegeateAddressResponse{Delegate: k.GetDelegateAddressFromValidator(ctx, val).String()}, nil
}

// QueryValidatorAddress implements QueryServer
func (k Keeper) QueryValidatorAddress(c context.Context, req *types.QueryValidatorAddressRequest) (*types.QueryValidatorAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	del, err := sdk.AccAddressFromBech32(req.Delegate)
	if err != nil {
		return nil, err
	}
	return &types.QueryValidatorAddressResponse{Validator: k.GetValidatorAddressFromDelegate(ctx, del).String()}, nil
}

// QueryOracleDataPrevote implements QueryServer
func (k Keeper) QueryOracleDataPrevote(c context.Context, req *types.QueryOracleDataPrevoteRequest) (*types.QueryOracleDataPrevoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.AccAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}
	if k.StakingKeeper.Validator(ctx, sdk.ValAddress(val)) == nil {
		val = k.GetValidatorAddressFromDelegate(ctx, val)
		if val == nil {
			return nil, fmt.Errorf("not a validator")
		}
	}
	return &types.QueryOracleDataPrevoteResponse{Hashes: k.GetOracleDataPrevote(ctx, val)}, nil
}

// QueryOracleDataVote implements QueryServer
func (k Keeper) QueryOracleDataVote(c context.Context, req *types.QueryOracleDataVoteRequest) (*types.MsgOracleDataVote, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.AccAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}
	if k.StakingKeeper.Validator(ctx, sdk.ValAddress(val)) == nil {
		val = k.GetValidatorAddressFromDelegate(ctx, val)
		if val == nil {
			return nil, fmt.Errorf("not a validator")
		}
	}
	return k.GetOracleDataVote(ctx, val), nil
}

// OracleData implements QueryServer
func (k Keeper) OracleData(c context.Context, req *types.QueryOracleDataRequest) (*types.QueryOracleDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	oda, err := types.PackOracleData(k.GetOracleData(ctx, req.Type))
	if err != nil {
		return nil, err
	}
	return &types.QueryOracleDataResponse{OracleData: oda}, nil
}

// QueryParams implements QueryServer
func (k Keeper) QueryParams(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{Params: k.GetParamSet(sdk.UnwrapSDKContext(c))}, nil

}

// QueryVotePeriod implements QueryServer
func (k Keeper) QueryVotePeriod(c context.Context, req *types.QueryVotePeriodRequest) (*types.VotePeriod, error) {
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
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.AccAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}
	if k.StakingKeeper.Validator(ctx, sdk.ValAddress(val)) == nil {
		val = k.GetValidatorAddressFromDelegate(ctx, val)
		if val == nil {
			return nil, fmt.Errorf("not a validator")
		}
	}
	return &types.QueryMissCounterResponse{MissCounter: k.GetMissCounter(ctx, val)}, nil
}
