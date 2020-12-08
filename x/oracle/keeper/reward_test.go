package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/stretchr/testify/require"
)

// Test a reward giving mechanism
func TestRewardBallotWinners(t *testing.T) {
	// initial setup
	input := CreateTestInput(t)
	addr, val := ValAddrs[0], PubKeys[0]
	addr1, val1 := ValAddrs[1], PubKeys[1]
	amt := sdk.TokensFromConsensusPower(100)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// TODO: Do this initialization in the test_utils file?
	for i := range []int{0, 1} {
		acc := input.AccKeeper.NewAccount(ctx, authtypes.NewBaseAccount(Addrs[i], AccPubKeys[i], uint64(i), 0))
		input.BankKeeper.SetBalances(ctx, acc.GetAddress(), InitCoins)
		input.AccKeeper.SetAccount(ctx, acc)
	}
	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(addr, val, amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(addr1, val1, amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	require.Equal(
		t, input.BankKeeper.GetAllBalances(ctx, sdk.AccAddress(addr)),
		sdk.NewCoins(sdk.NewCoin(input.StakingKeeper.GetParams(ctx).BondDenom, InitTokens.Sub(amt))),
	)
	require.Equal(t, amt, input.StakingKeeper.Validator(ctx, addr).GetBondedTokens())
	require.Equal(
		t, input.BankKeeper.GetAllBalances(ctx, sdk.AccAddress(addr1)),
		sdk.NewCoins(sdk.NewCoin(input.StakingKeeper.GetParams(ctx).BondDenom, InitTokens.Sub(amt))),
	)
	require.Equal(t, amt, input.StakingKeeper.Validator(ctx, addr1).GetBondedTokens())

	// Add claim pools
	claim := types.NewClaim(10, addr)
	claim2 := types.NewClaim(20, addr1)
	claims := map[string]types.Claim{
		addr.String():  claim,
		addr1.String(): claim2,
	}

	// Prepare reward pool
	givingAmt := sdk.NewCoins(sdk.NewInt64Coin(types.MicroLunaDenom, 30000000))
	acc := input.AccKeeper.GetModuleAccount(ctx, types.ModuleName)
	err = input.BankKeeper.SetBalances(ctx, acc.GetAddress(), givingAmt)
	require.NoError(t, err)
	input.AccKeeper.SetModuleAccount(ctx, acc)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.RewardDistributionWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	input.OracleKeeper.RewardBallotWinners(ctx, claims)
	outstandingRewardsDec := input.DistrKeeper.GetValidatorOutstandingRewards(ctx, addr)
	outstandingRewards, _ := outstandingRewardsDec.Rewards.TruncateDecimal()
	require.Equal(t, sdk.NewDecFromInt(givingAmt.AmountOf(types.MicroLunaDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(3).TruncateInt(),
		outstandingRewards.AmountOf(types.MicroLunaDenom))

	outstandingRewardsDec1 := input.DistrKeeper.GetValidatorOutstandingRewards(ctx, addr1)
	outstandingRewards1, _ := outstandingRewardsDec1.Rewards.TruncateDecimal()
	require.Equal(t, sdk.NewDecFromInt(givingAmt.AmountOf(types.MicroLunaDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(3).MulInt64(2).TruncateInt(),
		outstandingRewards1.AmountOf(types.MicroLunaDenom))
}
