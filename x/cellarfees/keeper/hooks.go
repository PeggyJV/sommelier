package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v2/x/gravity/types"
	"github.com/peggyjv/sommelier/v4/app/params"
	"github.com/peggyjv/sommelier/v4/x/cellarfees/types"
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

// In order to update the CellarFeePool, we check gravity module SendToCosmos transactions to see if the recipient is
// the cellarfees module account, and if the sender is a Cellar contract approved by governance. If these criteria are
// met, we account for the coins as fees in the pool.
func (h Hooks) AfterSendToCosmosEvent(ctx sdk.Context, event gravitytypes.SendToCosmosEvent) {
	// Check if recipient is the cellarfees module account
	moduleAccountAddress := h.k.accountKeeper.GetModuleAddress(types.ModuleName)
	if event.CosmosReceiver != moduleAccountAddress.String() {
		return
	}

	// Check if the sender is an approved Cellar contract
	if h.k.corkKeeper.HasCellarID(ctx, common.HexToAddress(event.EthereumSender)) {
		_, denom := h.k.gravityKeeper.ERC20ToDenomLookup(ctx, common.HexToAddress(event.TokenContract))

		// don't include SOMM in the pool
		if denom == params.BaseCoinUnit {
			return
		}

		balance := h.k.bankKeeper.GetBalance(ctx, moduleAccountAddress, denom)

		// sanity check
		if balance.Amount.LT(event.Amount) {
			panic("Coin balance in module account cannot be less than was sent from Ethereum!")
		}

		fee := sdk.Coin{
			Amount: event.Amount,
			Denom:  denom,
		}
		h.k.AddCoinToPool(ctx, fee)

		ctx.EventManager().EmitEvents(
			sdk.Events{
				sdk.NewEvent(
					types.EventTypeFeeAccrual,
					sdk.NewAttribute(types.AttributeKeyCellar, event.EthereumSender),
					sdk.NewAttribute(types.AttributeKeyTokenContract, event.TokenContract),
					sdk.NewAttribute(types.AttributeKeyDenom, denom),
					sdk.NewAttribute(types.AttributeKeyAmount, event.Amount.String()),
				),
			},
		)
	}
}
