package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v7/x/addresses/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.MsgServer = Keeper{}

func (k Keeper) AddAddressMapping(c context.Context, req *types.MsgAddAddressMapping) (*types.MsgAddAddressMappingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signer, err := sdk.AccAddressFromBech32(req.GetSigner())
	if err != nil {
		return &types.MsgAddAddressMappingResponse{}, status.Errorf(codes.InvalidArgument, "invalid signer address %s", req.GetSigner())
	}

	if !common.IsHexAddress(req.EvmAddress) {
		return &types.MsgAddAddressMappingResponse{}, status.Errorf(codes.InvalidArgument, "invalid EVM address %s", req.EvmAddress)
	}

	evmAddr := common.Hex2Bytes(req.EvmAddress)

	k.SetAddressMapping(ctx, signer.Bytes(), evmAddr)

	return &types.MsgAddAddressMappingResponse{}, nil
}

func (k Keeper) RemoveAddressMapping(c context.Context, req *types.MsgRemoveAddressMapping) (*types.MsgRemoveAddressMappingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signer, err := sdk.AccAddressFromBech32(req.GetSigner())
	if err != nil {
		return &types.MsgRemoveAddressMappingResponse{}, status.Errorf(codes.InvalidArgument, "invalid signer address %s", req.GetSigner())
	}

	k.DeleteAddressMapping(ctx, signer.Bytes())

	return &types.MsgRemoveAddressMappingResponse{}, nil
}
