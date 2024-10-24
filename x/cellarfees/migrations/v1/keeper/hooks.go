package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v5/x/gravity/types"
	"github.com/peggyjv/sommelier/v8/app/params"
	"github.com/peggyjv/sommelier/v8/x/cellarfees/types"
)

type Hooks struct {
	k Keeper
}

var _ gravitytypes.GravityHooks = Hooks{}

// Hooks Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

func (h Hooks) AfterContractCallExecutedEvent(ctx sdk.Context, event gravitytypes.ContractCallExecutedEvent) {
}

func (h Hooks) AfterERC20DeployedEvent(ctx sdk.Context, event gravitytypes.ERC20DeployedEvent) {}

func (h Hooks) AfterSignerSetExecutedEvent(ctx sdk.Context, event gravitytypes.SignerSetTxExecutedEvent) {
}

func (h Hooks) AfterBatchExecutedEvent(ctx sdk.Context, event gravitytypes.BatchExecutedEvent) {}

// Each time we receive a fee accrual from a cellar, we increment a counter for the respective denom. If a counter
// reaches a threshold defined in the cellarfees params, we attempt to start an auction. If the auction is started
// successfully we reset the count for that denom.
func (h Hooks) AfterSendToCosmosEvent(ctx sdk.Context, event gravitytypes.SendToCosmosEvent) {
	// Check if recipient is the cellarfees module account
	moduleAccountAddress := h.k.GetFeesAccount(ctx).GetAddress()
	if event.CosmosReceiver != moduleAccountAddress.String() {
		return
	}

	if event.Amount.IsZero() {
		return
	}

	// Check if the sender is an approved Cellar contract. We don't want to count coins sent from any address
	// as fee accruals.
	if !h.k.corkKeeper.HasCellarID(ctx, common.HexToAddress(event.EthereumSender)) {
		return
	}

	// Denom cannot be SOMM
	_, denom := h.k.gravityKeeper.ERC20ToDenomLookup(ctx, common.HexToAddress(event.TokenContract))
	if denom == params.BaseCoinUnit {
		return
	}

	counters := h.k.GetFeeAccrualCounters(ctx)
	count := counters.IncrementCounter(denom)
	h.k.SetFeeAccrualCounters(ctx, counters)

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeFeeAccrual,
				sdk.NewAttribute(types.AttributeKeyCellar, event.EthereumSender),
				sdk.NewAttribute(types.AttributeKeyTokenContract, event.TokenContract),
				sdk.NewAttribute(types.AttributeKeyDenom, denom),
				sdk.NewAttribute(types.AttributeKeyAmount, event.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyCount, fmt.Sprint(count)),
			),
		},
	)
}
