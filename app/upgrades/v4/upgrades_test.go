package v4

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/x/gravity/keeper"
	"github.com/peggyjv/gravity-bridge/module/x/gravity/types"
	"github.com/stretchr/testify/require"
)

func TestV2UpgradeDenomNormalization(t *testing.T) {
	input := keeper.CreateTestEnv(t)
	ctx := input.Context

	addr, _ := sdk.AccAddressFromBech32("cosmos1ahx7f8wyertuus9r20284ej0asrs085case3kn")
	erc20contract := common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5")
	amount := sdk.NewInt(1000)

	// mint some tokens
	incorrectDenom := strings.ToLower(types.GravityDenom(erc20contract))
	gravityCoins := sdk.NewCoins(sdk.NewCoin(incorrectDenom, amount))
	err := input.BankKeeper.MintCoins(ctx, types.ModuleName, gravityCoins)
	require.NoError(t, err)

	// set account's balance
	input.AccountKeeper.NewAccountWithAddress(ctx, addr)
	require.NoError(t, input.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, gravityCoins))
	oldBalance := input.BankKeeper.GetAllBalances(ctx, addr)
	require.Equal(t, oldBalance, gravityCoins)

	normalizeGravityDenoms(ctx, input.BankKeeper)

	normalizedDenom := types.NormalizeDenom(incorrectDenom)
	newBalance := input.BankKeeper.GetAllBalances(ctx, addr)
	require.Equal(t, newBalance, sdk.NewCoins(sdk.NewCoin(normalizedDenom, amount)))
}
