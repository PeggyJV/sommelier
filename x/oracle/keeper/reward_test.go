package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/peggyjv/sommelier/x/oracle/types"
)

// Test a reward giving mechanism
func TestRewardBallotWinners(t *testing.T) {
	// initial setup
	input := CreateTestInput(t)
	amt := sdk.TokensFromConsensusPower(100)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx
	bondDenom := input.StakingKeeper.GetParams(ctx).BondDenom

	// Set the account in state
	input.AccKeeper.SetAccount(ctx, input.AccKeeper.NewAccountWithAddress(ctx, Addrs[0]))
	input.AccKeeper.SetAccount(ctx, input.AccKeeper.NewAccountWithAddress(ctx, Addrs[1]))
	require.NoError(t, input.BankKeeper.SetBalances(ctx, Addrs[0], InitCoins))
	require.NoError(t, input.BankKeeper.SetBalances(ctx, Addrs[1], InitCoins))

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], PubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], PubKeys[1], amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	require.Equal(
		t, input.BankKeeper.GetAllBalances(ctx, Addrs[0]),
		sdk.NewCoins(sdk.NewCoin(bondDenom, InitTokens.Sub(amt))),
	)
	require.Equal(t, amt, input.StakingKeeper.Validator(ctx, ValAddrs[0]).GetBondedTokens())
	require.Equal(

		t, input.BankKeeper.GetAllBalances(ctx, Addrs[1]),
		sdk.NewCoins(sdk.NewCoin(bondDenom, InitTokens.Sub(amt))),
	)
	require.Equal(t, amt, input.StakingKeeper.Validator(ctx, ValAddrs[1]).GetBondedTokens())

	// Add claim pools
	claim := types.NewClaim(10, ValAddrs[0])
	claim2 := types.NewClaim(20, ValAddrs[1])
	claims := map[string]types.Claim{
		ValAddrs[0].String(): claim,
		ValAddrs[1].String(): claim2,
	}

	// Prepare reward pool
	givingAmt := sdk.NewCoins(sdk.NewInt64Coin(types.MicroLunaDenom, 30000000))
	acc := input.AccKeeper.GetModuleAccount(ctx, types.ModuleName)
	require.NoError(t, input.BankKeeper.SetBalances(input.Ctx, acc.GetAddress(), givingAmt))
	input.AccKeeper.SetModuleAccount(ctx, acc)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.RewardDistributionWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	input.OracleKeeper.RewardBallotWinners(ctx, claims)
	outstandingRewardsDec := input.DistrKeeper.GetValidatorOutstandingRewards(ctx, ValAddrs[0])
	outstandingRewards, _ := outstandingRewardsDec.GetRewards().TruncateDecimal()
	require.Equal(t, sdk.NewDecFromInt(givingAmt.AmountOf(types.MicroLunaDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(3).TruncateInt(),
		outstandingRewards.AmountOf(types.MicroLunaDenom))

	outstandingRewardsDec1 := input.DistrKeeper.GetValidatorOutstandingRewards(ctx, ValAddrs[1])
	outstandingRewards1, _ := outstandingRewardsDec1.GetRewards().TruncateDecimal()
	require.Equal(t, sdk.NewDecFromInt(givingAmt.AmountOf(types.MicroLunaDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(3).MulInt64(2).TruncateInt(),
		outstandingRewards1.AmountOf(types.MicroLunaDenom))
}
