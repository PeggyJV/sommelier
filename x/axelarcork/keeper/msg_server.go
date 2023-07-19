package keeper

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
)

var _ types.MsgServer = Keeper{}

// ScheduleCork implements types.MsgServer
func (k Keeper) ScheduleCork(c context.Context, msg *types.MsgScheduleAxelarCorkRequest) (*types.MsgScheduleAxelarCorkResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	if !k.GetParamSet(ctx).Enabled {
		return nil, types.ErrDisabled
	}

	signer := msg.MustGetSigner()
	validatorAddr := k.gravityKeeper.GetOrchestratorValidatorAddress(ctx, signer)
	if validatorAddr == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "signer %s is not a delegate", signer.String())
	}

	config, ok := k.GetChainConfigurationByID(ctx, msg.ChainId)
	if !ok {
		return nil, fmt.Errorf("chain by id %d not found", msg.ChainId)
	}

	if !k.HasCellarID(ctx, config.Id, common.HexToAddress(msg.Cork.TargetContractAddress)) {
		return nil, types.ErrUnmanagedCellarAddress
	}

	if msg.BlockHeight <= uint64(ctx.BlockHeight()) {
		return nil, types.ErrSchedulingInThePast
	}

	corkID := k.SetScheduledCork(ctx, config.Id, msg.BlockHeight, validatorAddr, *msg.Cork)

	if err := ctx.EventManager().EmitTypedEvent(&types.ScheduleCorkEvent{
		Signer:      signer.String(),
		Validator:   validatorAddr.String(),
		Cork:        msg.Cork.String(),
		BlockHeight: msg.BlockHeight,
		ChainId:     config.Id,
	}); err != nil {
		return nil, err
	}

	return &types.MsgScheduleAxelarCorkResponse{Id: hex.EncodeToString(corkID)}, nil
}

func (k Keeper) RelayCork(c context.Context, msg *types.MsgRelayAxelarCorkRequest) (*types.MsgRelayAxelarCorkResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParamSet(ctx)

	if !params.Enabled {
		return nil, types.ErrDisabled
	}

	config, ok := k.GetChainConfigurationByID(ctx, msg.ChainId)
	if !ok {
		return nil, fmt.Errorf("chain by id %d not found", msg.ChainId)
	}

	// winning cork will be deleted during the middleware pass
	cork, ok := k.GetWinningCork(ctx, config.Id, common.HexToAddress(msg.TargetContractAddress))
	if !ok {
		return nil, fmt.Errorf("no cork on chain %d found for address %s", config.Id, msg.TargetContractAddress)
	}

	proxyWrappedMsg := types.ProxyWrapper{
		Target: msg.TargetContractAddress,
		Body:   cork.EncodedContractCall,
	}
	pwbz, err := json.Marshal(proxyWrappedMsg)
	if err != nil {
		return nil, err
	}

	axelarMemo := types.AxelarBody{
		DestinationChain:   config.Name,
		DestinationAddress: config.ProxyAddress,
		Payload:            pwbz,
		Type:               types.PureMessage,
		Fee: &types.Fee{
			Amount:    strconv.FormatUint(msg.Fee, 10),
			Recipient: params.ExecutorAccount,
		},
	}
	bz, err := json.Marshal(axelarMemo)
	if err != nil {
		return nil, err
	}

	transferMsg := transfertypes.NewMsgTransfer(
		params.IbcPort,
		params.IbcChannel,
		msg.Token,
		msg.Signer,
		params.GmpAccount,
		clienttypes.ZeroHeight(),
		uint64(ctx.BlockTime().Add(time.Duration(params.TimeoutDuration)).UnixNano()),
	)
	transferMsg.Memo = string(bz)
	_, err = k.transferKeeper.Transfer(c, transferMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgRelayAxelarCorkResponse{}, nil
}

func (k Keeper) BumpCorkGas(c context.Context, msg *types.MsgBumpAxelarCorkGasRequest) (*types.MsgBumpAxelarCorkGasResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParamSet(ctx)

	if !params.Enabled {
		return nil, types.ErrDisabled
	}

	transferMsg := transfertypes.NewMsgTransfer(
		params.IbcPort,
		params.IbcChannel,
		msg.Token,
		msg.Signer,
		params.ExecutorAccount,
		clienttypes.ZeroHeight(),
		uint64(ctx.BlockTime().Add(time.Duration(params.TimeoutDuration)).UnixNano()),
	)
	transferMsg.Memo = msg.MessageId
	_, err := k.transferKeeper.Transfer(c, transferMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgBumpAxelarCorkGasResponse{}, nil
}
