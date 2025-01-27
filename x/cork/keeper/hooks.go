package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
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

func (h Hooks) AfterSendToCosmosEvent(ctx sdk.Context, event gravitytypes.SendToCosmosEvent) {}
