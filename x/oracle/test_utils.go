// nolint:deadcode unused DONTCOVER
package oracle

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/peggyjv/sommelier/x/oracle/keeper"
	"github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/stretchr/testify/require"
)

var (
	uSDRAmt    = sdk.NewInt(1005 * types.MicroUnit)
	stakingAmt = sdk.TokensFromConsensusPower(10)

	randomExchangeRate        = sdk.NewDec(1700)
	anotherRandomExchangeRate = sdk.NewDecWithPrec(4882, 2) // swap rate
)

func setup_with_small_voting_power(t *testing.T) (keeper.TestInput, sdk.Handler) {
	input := keeper.CreateTestInput(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 1
	params.SlashWindow = 100
	params.RewardDistributionWindow = 100
	input.OracleKeeper.SetParams(input.Ctx, params)
	h := NewHandler(input.OracleKeeper)

	bd := input.StakingKeeper.GetParams(input.Ctx).BondDenom
	acc := input.AccKeeper.NewAccount(input.Ctx, authtypes.NewBaseAccount(keeper.Addrs[0], keeper.AccPubKeys[0], 0, 0))
	input.BankKeeper.SetBalances(input.Ctx, acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(bd, stakingAmt.Add(sdk.NewInt(100)))))
	input.AccKeeper.SetAccount(input.Ctx, acc)

	sh := staking.NewHandler(input.StakingKeeper)
	_, err := sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[0], keeper.PubKeys[0], sdk.TokensFromConsensusPower(1)))
	require.NoError(t, err)

	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	return input, h
}

func setup(t *testing.T) (keeper.TestInput, sdk.Handler) {
	input := keeper.CreateTestInput(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 1
	params.SlashWindow = 100
	params.RewardDistributionWindow = 100
	input.OracleKeeper.SetParams(input.Ctx, params)
	h := NewHandler(input.OracleKeeper)

	sh := staking.NewHandler(input.StakingKeeper)

	bd := input.StakingKeeper.GetParams(input.Ctx).BondDenom
	for i := range []int{0, 1, 2} {
		acc := input.AccKeeper.NewAccount(input.Ctx, authtypes.NewBaseAccount(keeper.Addrs[i], keeper.AccPubKeys[i], uint64(i), 0))
		input.BankKeeper.SetBalances(input.Ctx, acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(bd, stakingAmt.Add(sdk.NewInt(100)))))
		input.AccKeeper.SetAccount(input.Ctx, acc)
	}

	// Validator created
	_, err := sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[0], keeper.PubKeys[0], stakingAmt))
	require.NoError(t, err)
	_, err = sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[1], keeper.PubKeys[1], stakingAmt))
	require.NoError(t, err)
	_, err = sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[2], keeper.PubKeys[2], stakingAmt))
	require.NoError(t, err)
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	return input, h
}

func setupVal5(t *testing.T) (keeper.TestInput, sdk.Handler) {
	input := keeper.CreateTestInput(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 1
	params.SlashWindow = 100
	params.RewardDistributionWindow = 100
	input.OracleKeeper.SetParams(input.Ctx, params)
	h := NewHandler(input.OracleKeeper)

	sh := staking.NewHandler(input.StakingKeeper)

	bd := input.StakingKeeper.GetParams(input.Ctx).BondDenom
	for i := range []int{0, 1, 2, 3, 4} {
		acc := input.AccKeeper.NewAccount(input.Ctx, authtypes.NewBaseAccount(keeper.Addrs[i], keeper.AccPubKeys[i], uint64(i), 0))
		input.BankKeeper.SetBalances(input.Ctx, acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(bd, stakingAmt.Add(sdk.NewInt(100)))))
		input.AccKeeper.SetAccount(input.Ctx, acc)
	}

	// Validator created
	_, err := sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[0], keeper.PubKeys[0], stakingAmt))
	require.NoError(t, err)
	_, err = sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[1], keeper.PubKeys[1], stakingAmt))
	require.NoError(t, err)
	_, err = sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[2], keeper.PubKeys[2], stakingAmt))
	require.NoError(t, err)
	_, err = sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[3], keeper.PubKeys[3], stakingAmt))
	require.NoError(t, err)
	_, err = sh(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[4], keeper.PubKeys[4], stakingAmt))
	require.NoError(t, err)
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	return input, h
}
