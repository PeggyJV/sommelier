package keeper

import (
	"context"
	"encoding/hex"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
	corktypes "github.com/peggyjv/sommelier/v10/x/cork/types"
	types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
)

var _ types.MsgServer = Keeper{}

// ScheduleCork implements types.MsgServer
func (k Keeper) ScheduleCork(c context.Context, msg *types.MsgScheduleCorkRequest) (*types.MsgScheduleCorkResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	if k.inSafeMode(ctx) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "x/poa safe mode active: scheduling corks is frozen until the authority set is restored")
	}

	signer := msg.MustGetSigner()
	params := k.GetParamSet(ctx)
	// Fail-closed: an unset authority means no address may schedule a cork.
	// There is deliberately no fallback to the retired validator-delegate path.
	if params.CorkAuthority == "" || signer.String() != params.CorkAuthority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
			"signer %s is not the cork authority", signer.String())
	}

	if !k.HasCellarID(ctx, common.HexToAddress(msg.Cork.TargetContractAddress)) {
		return nil, corktypes.ErrUnmanagedCellarAddress
	}

	if msg.BlockHeight <= uint64(ctx.BlockHeight()) {
		return nil, corktypes.ErrSchedulingInThePast
	}

	corkID := k.SetAuthorityCork(ctx, msg.BlockHeight, *msg.Cork)

	invalidationScope := msg.Cork.InvalidationScope()
	// The nonce the EndBlocker will consume when it submits this cork.
	invalidationNonce := k.GetLatestInvalidationNonce(ctx) + 1

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				corktypes.EventTypeCork,
				sdk.NewAttribute(sdk.AttributeKeyModule, corktypes.AttributeValueCategory),
				sdk.NewAttribute(corktypes.AttributeKeySigner, signer.String()),
				sdk.NewAttribute(corktypes.AttributeKeyCork, msg.Cork.String()),
				sdk.NewAttribute(corktypes.AttributeKeyBlockHeight, fmt.Sprintf("%d", msg.BlockHeight)),
				sdk.NewAttribute(corktypes.AttributeKeyCorkID, hex.EncodeToString(corkID)),
				sdk.NewAttribute(gravitytypes.AttributeKeyContractCallInvalidationScope, fmt.Sprint(invalidationScope)),
				sdk.NewAttribute(gravitytypes.AttributeKeyContractCallInvalidationNonce, fmt.Sprint(invalidationNonce)),
			),
		},
	)

	return &types.MsgScheduleCorkResponse{Id: hex.EncodeToString(corkID)}, nil
}
