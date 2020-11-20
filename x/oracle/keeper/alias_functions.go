package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

// GetOracleAccount returns oracle ModuleAccount
func (k Keeper) GetOracleAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetRewardPool retrieves the balance of the oracle module account
func (k Keeper) GetRewardPool(ctx sdk.Context) sdk.Coins {
	acc := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	return acc.GetCoins()
}
