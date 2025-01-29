package axelarcork

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/peggyjv/sommelier/v9/x/axelarcork/keeper"
	"github.com/peggyjv/sommelier/v9/x/axelarcork/types"
)

// NewHandler returns a handler for "axelarcork" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgScheduleAxelarCorkRequest:
			res, err := k.ScheduleCork(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRelayAxelarCorkRequest:
			res, err := k.RelayCork(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRelayAxelarProxyUpgradeRequest:
			res, err := k.RelayProxyUpgrade(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgBumpAxelarCorkGasRequest:
			res, err := k.BumpCorkGas(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized axelar cork message type: %T", msg)
		}
	}
}

func NewProposalHandler(k keeper.Keeper) govtypesv1beta1.Handler {
	return func(ctx sdk.Context, content govtypesv1beta1.Content) error {
		switch c := content.(type) {
		case *types.AddAxelarManagedCellarIDsProposal:
			return keeper.HandleAddManagedCellarsProposal(ctx, k, *c)
		case *types.RemoveAxelarManagedCellarIDsProposal:
			return keeper.HandleRemoveManagedCellarsProposal(ctx, k, *c)
		case *types.AxelarScheduledCorkProposal:
			return keeper.HandleScheduledCorkProposal(ctx, k, *c)
		case *types.AxelarCommunityPoolSpendProposal:
			return keeper.HandleCommunityPoolSpendProposal(ctx, k, *c)
		case *types.AddChainConfigurationProposal:
			return keeper.HandleAddChainConfigurationProposal(ctx, k, *c)
		case *types.RemoveChainConfigurationProposal:
			return keeper.HandleRemoveChainConfigurationProposal(ctx, k, *c)
		case *types.UpgradeAxelarProxyContractProposal:
			return keeper.HandleUpgradeAxelarProxyContractProposal(ctx, k, *c)
		case *types.CancelAxelarProxyContractUpgradeProposal:
			return keeper.HandleCancelAxelarProxyContractUpgradeProposal(ctx, k, *c)

		default:
			return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized axelar cork proposal content type: %T", c)
		}
	}
}
