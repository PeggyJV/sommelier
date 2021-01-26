package keeper

import (
	"context"

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

// QueryOracleDataPrevote implements QueryServer
func (k Keeper) QueryOracleDataPrevote(c context.Context, req *types.QueryOracleDataPrevoteRequest) (*types.QueryOracleDataPrevoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	val, err := sdk.AccAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
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
