package oracle

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/sommelier/x/oracle/keeper"
	"github.com/peggyjv/sommelier/x/oracle/types"
	core "github.com/peggyjv/sommelier/x/oracle/types"
)

func TestOracleDrop(t *testing.T) {
	input, h := setup(t)

	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, core.MicroKRWDenom, randomExchangeRate)

	// Account 1, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 0)

	// Immediately swap halt after an illiquid oracle vote
	EndBlocker(input.Ctx, input.OracleKeeper)

	_, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroKRWDenom)
	require.NotNil(t, err)
}

func TestOracleTallyTiming(t *testing.T) {
	input, h := setup(t)

	// all the keeper.Addrs vote for the block ... not last period block yet, so tally fails
	for i := range keeper.Addrs[:2] {
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), i)
	}

	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 10 // set vote period to 10 for now, for convinience
	input.OracleKeeper.SetParams(input.Ctx, params)
	require.Equal(t, 0, int(input.Ctx.BlockHeight()))

	EndBlocker(input.Ctx, input.OracleKeeper)
	_, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroSDRDenom)
	require.Error(t, err)

	input.Ctx = input.Ctx.WithBlockHeight(params.VotePeriod - 1)

	EndBlocker(input.Ctx, input.OracleKeeper)
	_, err = input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroSDRDenom)
	require.NoError(t, err)
}

func TestOracleRewardDistribution(t *testing.T) {
	input, h := setup(t)

	// Account 1, SDR
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroSDRDenom, randomExchangeRate)), 0)

	// Account 2, SDR
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroSDRDenom, randomExchangeRate)), 1)

	rewardsAmt := sdk.NewInt(100000000)
	moduleAcc := input.AccKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), types.ModuleName)
	err := input.BankKeeper.SetBalances(input.Ctx, moduleAcc.GetAddress(), sdk.NewCoins(sdk.NewCoin(core.MicroLunaDenom, rewardsAmt)))
	require.NoError(t, err)

	input.AccKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.RewardDistributionWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	expectedRewardAmt := sdk.NewDecFromInt(rewardsAmt.QuoRaw(2)).QuoInt64(votePeriodsPerWindow).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
}

func TestOracleRewardBand(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = []*types.Denom{{Name: core.MicroKRWDenom, TobinTax: types.DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, types.DefaultTobinTax)

	rewardSpread := randomExchangeRate.Mul(input.OracleKeeper.RewardBand(input.Ctx).QuoInt64(2))

	// no one will miss the vote
	// Account 1, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate.Sub(rewardSpread))), 0)

	// Account 2, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 1)

	// Account 3, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate.Add(rewardSpread))), 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// Account 1 will miss the vote due to raward band condition
	// Account 1, KRW

	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate.Sub(rewardSpread.Add(sdk.OneDec())))), 0)

	// Account 2, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 1)

	// Account 3, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate.Add(rewardSpread))), 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

}

func TestOracleMultiRewardDistribution(t *testing.T) {
	input, h := setup(t)

	// SDR and KRW have the same voting power, but KRW has been chosen as referenceTerra by alphabetical order.
	// Account 1, SDR
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroSDRDenom, randomExchangeRate)), 0)

	// Account 1, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 0)

	// Account 2, SDR
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroSDRDenom, randomExchangeRate)), 1)

	// Account 3, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 2)

	rewardAmt := sdk.NewInt(100000000)
	moduleAcc := input.AccKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), types.ModuleName)
	require.NoError(t, input.BankKeeper.SetBalance(input.Ctx, moduleAcc.GetAddress(), sdk.NewCoin(core.MicroLunaDenom, rewardAmt)))

	input.AccKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rewardDistributedWindow := input.OracleKeeper.RewardDistributionWindow(input.Ctx)

	expectedRewardAmt := sdk.NewDecFromInt(rewardAmt.QuoRaw(3).MulRaw(2)).QuoInt64(rewardDistributedWindow).TruncateInt()
	expectedRewardAmt2 := sdk.ZeroInt() // even vote power is same KRW with SDR, KRW chosen referenceTerra because alphabetical order
	expectedRewardAmt3 := sdk.NewDecFromInt(rewardAmt.QuoRaw(3)).QuoInt64(rewardDistributedWindow).TruncateInt()

	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt2, rewards.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt3, rewards.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
}

func TestOracleExchangeRate(t *testing.T) {
	input, h := setup(t)

	krwRandomExchangeRate := sdk.NewDecWithPrec(1000000000, int64(6)).MulInt64(core.MicroUnit)
	uswRandomExchangeRate := sdk.NewDecWithPrec(1000000, int64(6)).MulInt64(core.MicroUnit)

	// KRW has been chosen as referenceTerra by highest voting power
	// Account 1, USD, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroUSDDenom, uswRandomExchangeRate)), 0)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, krwRandomExchangeRate)), 0)

	// Account 2, USD, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroUSDDenom, uswRandomExchangeRate)), 1)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, krwRandomExchangeRate)), 1)

	// Account 3, KRW, SDR
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, krwRandomExchangeRate)), 2)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroSDRDenom, randomExchangeRate)), 2)

	rewardAmt := sdk.NewInt(100000000)
	moduleAcc := input.AccKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), types.ModuleName)
	require.NoError(t, input.BankKeeper.SetBalance(input.Ctx, moduleAcc.GetAddress(), sdk.NewCoin(core.MicroLunaDenom, rewardAmt)))

	input.AccKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rewardDistributedWindow := input.OracleKeeper.RewardDistributionWindow(input.Ctx)
	expectedRewardAmt := sdk.NewDecFromInt(rewardAmt.QuoRaw(5).MulRaw(2)).QuoInt64(rewardDistributedWindow).TruncateInt()
	expectedRewardAmt2 := sdk.NewDecFromInt(rewardAmt.QuoRaw(5).MulRaw(1)).QuoInt64(rewardDistributedWindow).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt2, rewards.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
}

func TestOracleExchangeRateVal5(t *testing.T) {
	input, h := setupVal5(t)

	krwExchangeRate := sdk.NewDecWithPrec(505000, int64(6)).MulInt64(core.MicroUnit)
	krwExchangeRateWithErr := sdk.NewDecWithPrec(500000, int64(6)).MulInt64(core.MicroUnit)
	usdExchangeRate := sdk.NewDecWithPrec(505, int64(6)).MulInt64(core.MicroUnit)
	usdExchangeRateWithErr := sdk.NewDecWithPrec(500, int64(6)).MulInt64(core.MicroUnit)

	// KRW has been chosen as referenceTerra by highest voting power
	// Account 1, KRW, USD
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, krwExchangeRate)), 0)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroUSDDenom, usdExchangeRate)), 0)

	// Account 2, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, krwExchangeRate)), 1)

	// Account 3, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, krwExchangeRate)), 2)

	// Account 4, KRW, USD
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, krwExchangeRateWithErr)), 3)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroUSDDenom, usdExchangeRateWithErr)), 3)

	// Account 5, KRW, USD
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, krwExchangeRateWithErr)), 4)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroUSDDenom, usdExchangeRateWithErr)), 4)

	rewardAmt := sdk.NewInt(100000000)
	moduleAcc := input.AccKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), types.ModuleName)
	require.NoError(t, input.BankKeeper.SetBalance(input.Ctx, moduleAcc.GetAddress(), sdk.NewCoin(core.MicroLunaDenom, rewardAmt)))

	input.AccKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	krw, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroKRWDenom)
	require.NoError(t, err)
	usd, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroUSDDenom)
	require.NoError(t, err)

	// legacy version case
	require.NotEqual(t, usdExchangeRateWithErr, usd)

	// new version case
	require.Equal(t, krwExchangeRate, krw)
	require.Equal(t, usdExchangeRate, usd)

	rewardDistributedWindow := input.OracleKeeper.RewardDistributionWindow(input.Ctx)
	expectedRewardAmt := sdk.NewDecFromInt(rewardAmt.QuoRaw(8).MulRaw(2)).QuoInt64(rewardDistributedWindow).TruncateInt()
	expectedRewardAmt2 := sdk.NewDecFromInt(rewardAmt.QuoRaw(8).MulRaw(1)).QuoInt64(rewardDistributedWindow).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards1 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt2, rewards1.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards2 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt2, rewards2.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards3 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[3])
	require.Equal(t, expectedRewardAmt, rewards3.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
	rewards4 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[4])
	require.Equal(t, expectedRewardAmt, rewards4.Rewards.AmountOf(core.MicroLunaDenom).TruncateInt())
}

func TestInvalidVotesSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = []*types.Denom{{Name: core.MicroKRWDenom, TobinTax: types.DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, types.DefaultTobinTax)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.SlashWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := int64(0); i < sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64(); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 1, KRW
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 0)

		// Account 2, KRW, miss vote
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate.Add(sdk.NewDec(100000000000000)))), 1)

		// Account 3, KRW
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 2)

		EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, i+1, input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// one more miss vote will inccur ValAddrs[1] slashing
	// Account 1, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 0)

	// Account 2, KRW, miss vote
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate.Add(sdk.NewDec(100000000000000)))), 1)

	// Account 3, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 2)

	input.Ctx = input.Ctx.WithBlockHeight(votePeriodsPerWindow - 1)
	EndBlocker(input.Ctx, input.OracleKeeper)
	validator = input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, sdk.OneDec().Sub(slashFraction).MulInt(stakingAmt).TruncateInt(), validator.GetBondedTokens())
}

func TestWhitelistSlashing(t *testing.T) {
	input, h := setup(t)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.SlashWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := int64(0); i < sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64(); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 2, KRW
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 1)
		// Account 3, KRW
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 2)

		EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, i+1, input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// one more miss vote will inccur Account 1 slashing

	// Account 2, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 1)
	// Account 3, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 2)

	input.Ctx = input.Ctx.WithBlockHeight(votePeriodsPerWindow - 1)
	EndBlocker(input.Ctx, input.OracleKeeper)
	validator = input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, sdk.OneDec().Sub(slashFraction).MulInt(stakingAmt).TruncateInt(), validator.GetBondedTokens())
}

func TestNotPassedBallotSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = []*types.Denom{{Name: core.MicroKRWDenom, TobinTax: types.DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, types.DefaultTobinTax)

	input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

	// Account 1, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 0)

	EndBlocker(input.Ctx, input.OracleKeeper)
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))
}

func TestAbstainSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = []*types.Denom{{Name: core.MicroKRWDenom, TobinTax: types.DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, types.DefaultTobinTax)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.SlashWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := int64(0); i <= sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64(); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 1, KRW
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 0)

		// Account 2, KRW, abstain vote
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, sdk.ZeroDec())), 1)

		// Account 3, KRW
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 2)

		EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())
}

func TestVoteTargets(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = []*types.Denom{{Name: core.MicroKRWDenom, TobinTax: types.DefaultTobinTax}, {Name: core.MicroSDRDenom, TobinTax: types.DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, types.DefaultTobinTax)

	// KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 0)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 1)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	// no missing current
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// vote targets are {KRW, SDR}
	require.Equal(t, []string{core.MicroKRWDenom, core.MicroSDRDenom}, input.OracleKeeper.GetVoteTargets(input.Ctx))

	// tobin tax must be exists for SDR
	sdrTobinTax, err := input.OracleKeeper.GetTobinTax(input.Ctx, core.MicroSDRDenom)
	require.NoError(t, err)
	require.Equal(t, types.DefaultTobinTax, sdrTobinTax)

	// delete SDR
	params.Whitelist = []*types.Denom{{Name: core.MicroKRWDenom, TobinTax: types.DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// KRW, missing
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 0)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 1)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// SDR must be deleted
	require.Equal(t, []string{core.MicroKRWDenom}, input.OracleKeeper.GetVoteTargets(input.Ctx))

	_, err = input.OracleKeeper.GetTobinTax(input.Ctx, core.MicroSDRDenom)
	require.Error(t, err)

	// change KRW tobin tax
	params.Whitelist = []*types.Denom{{Name: core.MicroKRWDenom, TobinTax: sdk.ZeroDec()}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// KRW, no missing
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 0)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 1)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, randomExchangeRate)), 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// KRW tobin tax must be 0
	tobinTax, err := input.OracleKeeper.GetTobinTax(input.Ctx, core.MicroKRWDenom)
	require.NoError(t, err)
	require.True(t, sdk.ZeroDec().Equal(tobinTax))
}

func TestAbstainWithSmallStakingPower(t *testing.T) {
	input, h := setup_with_small_voting_power(t)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, core.MicroKRWDenom, types.DefaultTobinTax)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.NewDecCoins(sdk.NewDecCoinFromDec(core.MicroKRWDenom, sdk.ZeroDec())), 0)

	EndBlocker(input.Ctx, input.OracleKeeper)
	_, err := input.OracleKeeper.GetLunaExchangeRate(input.Ctx, core.MicroKRWDenom)
	require.Error(t, err)
}

// func makeAggregatePrevoteAndVote(t *testing.T, input keeper.TestInput, h sdk.Handler, height int64, denom string, rate sdk.Dec, idx int) {
// 	// Account 1, SDR
// 	salt := "1"
// 	hash := GetVoteHash(salt, rate, denom, keeper.ValAddrs[idx])

// 	prevoteMsg := NewMsgExchangeRatePrevote(hash, denom, keeper.Addrs[idx], keeper.ValAddrs[idx])
// 	_, err := h(input.Ctx.WithBlockHeight(height), prevoteMsg)
// 	require.NoError(t, err)

// 	voteMsg := NewMsgExchangeRateVote(rate, salt, denom, keeper.Addrs[idx], keeper.ValAddrs[idx])
// 	_, err = h(input.Ctx.WithBlockHeight(height+1), voteMsg)
// 	require.NoError(t, err)
// }

func makeAggregatePrevoteAndVote(t *testing.T, input keeper.TestInput, h sdk.Handler, height int64, rates sdk.DecCoins, idx int) {
	// Account 1, SDR
	salt := "1"
	hash := types.GetAggregateVoteHash(salt, rates.String(), keeper.ValAddrs[idx])

	prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[idx], keeper.ValAddrs[idx])
	_, err := h(input.Ctx.WithBlockHeight(height), &prevoteMsg)
	require.NoError(t, err)

	voteMsg := types.NewMsgAggregateExchangeRateVote(salt, rates.String(), keeper.Addrs[idx], keeper.ValAddrs[idx])
	_, err = h(input.Ctx.WithBlockHeight(height+1), &voteMsg)
	require.NoError(t, err)
}
