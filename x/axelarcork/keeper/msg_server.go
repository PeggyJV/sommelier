package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
)

var _ types.MsgServer = Keeper{}

// ScheduleCork implements types.MsgServer
func (k Keeper) ScheduleCork(c context.Context, msg *types.MsgScheduleCorkRequest) (*types.MsgScheduleCorkResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	config, err := k.GetChainConfigurationByNameAndID(ctx, msg.ChainName, msg.ChainId)
	if err != nil {
		return nil, err
	}

	if !config.HasCellarID(common.HexToAddress(msg.Cork.TargetContractAddress)) {
		return nil, types.ErrUnmanagedCellarAddress
	}

	if msg.BlockHeight <= uint64(ctx.BlockHeight()) {
		return nil, types.ErrSchedulingInThePast
	}

	signer := msg.MustGetSigner()
	validator := k.stakingKeeper.Validator(ctx, sdk.ValAddress(signer))
	if validator != nil {
		return nil, fmt.Errorf("validator not found for acc addr %s", signer)
	}
	validatorAddr := sdk.ValAddress(signer)

	corkID := k.SetScheduledCork(ctx, config.Id, msg.BlockHeight, validatorAddr, *msg.Cork)

	if err := ctx.EventManager().EmitTypedEvent(&types.CorkEvent{
		Signer:      signer.String(),
		Validator:   validatorAddr.String(),
		Cork:        msg.Cork.String(),
		BlockHeight: msg.BlockHeight,
		ChainId:     config.Id,
	}); err != nil {
		return nil, err
	}

	return &types.MsgScheduleCorkResponse{Id: hex.EncodeToString(corkID)}, nil
}
