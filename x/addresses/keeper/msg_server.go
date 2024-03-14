package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	typesv1 "github.com/peggyjv/sommelier/v7/x/addresses/types/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ typesv1.MsgServer = Keeper{}

func (k Keeper) AddAddressMapping(c context.Context, req *typesv1.MsgAddAddressMapping) (*typesv1.MsgAddAddressMappingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signer, err := sdk.AccAddressFromBech32(req.GetSigner())
	if err != nil {
		return &typesv1.MsgAddAddressMappingResponse{}, status.Errorf(codes.InvalidArgument, "invalid signer address %s", req.GetSigner())
	}

	if !common.IsHexAddress(req.EvmAddress) {
		return &typesv1.MsgAddAddressMappingResponse{}, status.Errorf(codes.InvalidArgument, "invalid EVM address %s", req.EvmAddress)
	}

	evmAddr := common.Hex2Bytes(req.EvmAddress)

	k.SetAddressMapping(ctx, signer.Bytes(), evmAddr)

	return &typesv1.MsgAddAddressMappingResponse{}, nil
}

func (k Keeper) RemoveAddressMapping(c context.Context, req *typesv1.MsgRemoveAddressMapping) (*typesv1.MsgRemoveAddressMappingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signer, err := sdk.AccAddressFromBech32(req.GetSigner())
	if err != nil {
		return &typesv1.MsgRemoveAddressMappingResponse{}, status.Errorf(codes.InvalidArgument, "invalid signer address %s", req.GetSigner())
	}

	k.DeleteAddressMapping(ctx, signer.Bytes())

	return &typesv1.MsgRemoveAddressMappingResponse{}, nil
}
