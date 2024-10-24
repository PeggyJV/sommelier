package keeper

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/sommelier/v8/x/axelarcork/types"
)

var _ types.MsgServer = Keeper{}

// ScheduleCork implements types.MsgServer
func (k Keeper) ScheduleCork(c context.Context, msg *types.MsgScheduleAxelarCorkRequest) (*types.MsgScheduleAxelarCorkResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParamSet(ctx)
	if !params.Enabled {
		return nil, types.ErrDisabled
	}

	signer := msg.MustGetSigner()
	validatorAddr := k.gravityKeeper.GetOrchestratorValidatorAddress(ctx, signer)
	if validatorAddr == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "signer %s is not a delegate", signer.String())
	}

	validatorAxelarCorkCount := k.GetValidatorAxelarCorkCount(ctx, validatorAddr)
	if validatorAxelarCorkCount >= types.MaxAxelarCorksPerValidator {
		return nil, types.ErrValidatorAxelarCorkCapacityReached
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

	corkID := k.SetScheduledAxelarCork(ctx, config.Id, msg.BlockHeight, validatorAddr, *msg.Cork)
	k.IncrementValidatorAxelarCorkCount(ctx, validatorAddr)

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
	_, cork, ok := k.GetWinningAxelarCork(ctx, config.Id, common.HexToAddress(msg.TargetContractAddress))
	if !ok {
		return nil, fmt.Errorf("no cork on chain %d found for address %s", config.Id, msg.TargetContractAddress)
	}

	// transfer tokens to the module account
	signer := msg.MustGetSigner()
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, signer, types.ModuleName, sdk.NewCoins(msg.Token))
	if err != nil {
		return nil, err
	}

	nonce := k.IncrementAxelarContractCallNonce(ctx, msg.ChainId, cork.TargetContractAddress)
	payload, err := types.EncodeLogicCallArgs(msg.TargetContractAddress, nonce, cork.Deadline, cork.EncodedContractCall)
	if err != nil {
		return nil, err
	}

	axelarMemo := types.AxelarBody{
		DestinationChain:   config.Name,
		DestinationAddress: config.ProxyAddress,
		Payload:            payload,
		Type:               types.PureMessage,
		Fee: &types.Fee{
			Amount:    strconv.FormatUint(msg.Fee, 10),
			Recipient: params.ExecutorAccount,
		},
	}
	memoBz, err := json.Marshal(axelarMemo)
	if err != nil {
		return nil, err
	}

	transferMsg := transfertypes.NewMsgTransfer(
		params.IbcPort,
		params.IbcChannel,
		msg.Token,
		k.GetSenderAccount(ctx).GetAddress().String(),
		params.GmpAccount,
		clienttypes.ZeroHeight(),
		uint64(ctx.BlockTime().Add(time.Duration(params.TimeoutDuration)).UnixNano()),
		string(memoBz),
	)
	_, err = k.transferKeeper.Transfer(c, transferMsg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(types.EventTypeAxelarCorkRelayCalled,
				sdk.NewAttribute(types.AttributeKeyCork, cork.String()),
				sdk.NewAttribute(types.AttributeKeyDeadline, fmt.Sprintf("%d", cork.Deadline)),
			),
		},
	)

	return &types.MsgRelayAxelarCorkResponse{}, nil
}

// RelayProxyUpgrade prepares the payload for IBC transfer if the current blockheight meets
// the payload's execution threshold, and then deletes the upgrade data.
func (k Keeper) RelayProxyUpgrade(c context.Context, msg *types.MsgRelayAxelarProxyUpgradeRequest) (*types.MsgRelayAxelarProxyUpgradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParamSet(ctx)

	if !params.Enabled {
		return nil, types.ErrDisabled
	}

	config, ok := k.GetChainConfigurationByID(ctx, msg.ChainId)
	if !ok {
		return nil, fmt.Errorf("chain by id %d not found", msg.ChainId)
	}

	upgradeData, found := k.GetAxelarProxyUpgradeData(ctx, msg.ChainId)
	if !found {
		return nil, fmt.Errorf("no proxy upgrade data found for chain %d", msg.ChainId)
	}

	if ctx.BlockHeight() < upgradeData.ExecutableHeightThreshold {
		return nil, fmt.Errorf("proxy upgrade call is not executable until height %d", upgradeData.ExecutableHeightThreshold)
	}

	// transfer tokens to the module account
	signer := msg.MustGetSigner()
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, signer, types.ModuleName, sdk.NewCoins(msg.Token))
	if err != nil {
		return nil, err
	}

	axelarMemo := types.AxelarBody{
		DestinationChain:   config.Name,
		DestinationAddress: config.ProxyAddress,
		Payload:            upgradeData.Payload,
		Type:               types.PureMessage,
		Fee: &types.Fee{
			Amount:    strconv.FormatUint(msg.Fee, 10),
			Recipient: params.ExecutorAccount,
		},
	}
	memoBz, err := json.Marshal(axelarMemo)
	if err != nil {
		return nil, err
	}

	transferMsg := transfertypes.NewMsgTransfer(
		params.IbcPort,
		params.IbcChannel,
		msg.Token,
		k.GetSenderAccount(ctx).GetAddress().String(),
		params.GmpAccount,
		clienttypes.ZeroHeight(),
		uint64(ctx.BlockTime().Add(time.Duration(params.TimeoutDuration)).UnixNano()),
		string(memoBz),
	)
	_, err = k.transferKeeper.Transfer(c, transferMsg)
	if err != nil {
		return nil, err
	}

	k.DeleteAxelarProxyUpgradeData(ctx, msg.ChainId)

	return &types.MsgRelayAxelarProxyUpgradeResponse{}, nil
}

func (k Keeper) BumpCorkGas(c context.Context, msg *types.MsgBumpAxelarCorkGasRequest) (*types.MsgBumpAxelarCorkGasResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParamSet(ctx)

	if !params.Enabled {
		return nil, types.ErrDisabled
	}

	memo := msg.MessageId
	transferMsg := transfertypes.NewMsgTransfer(
		params.IbcPort,
		params.IbcChannel,
		msg.Token,
		msg.Signer,
		params.ExecutorAccount,
		clienttypes.ZeroHeight(),
		uint64(ctx.BlockTime().Add(time.Duration(params.TimeoutDuration)).UnixNano()),
		memo,
	)
	_, err := k.transferKeeper.Transfer(c, transferMsg)
	if err != nil {
		return nil, err
	}

	return &types.MsgBumpAxelarCorkGasResponse{}, nil
}

func (k Keeper) CancelScheduledCork(c context.Context, msg *types.MsgCancelAxelarCorkRequest) (*types.MsgCancelAxelarCorkResponse, error) {

	// todo: implement

	return &types.MsgCancelAxelarCorkResponse{}, nil
}
