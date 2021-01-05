package oracle

import (
	"fmt"
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/peggyjv/sommelier/x/oracle/keeper"
	"github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestOracleThreshold(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	// Prevote without exchange rate
	salt := "1"
	exchangeRates := sdk.DecCoins{sdk.NewDecCoinFromDec(types.MicroSDRDenom, randomExchangeRate)}
	hash := types.GetAggregateVoteHash(salt, exchangeRates.String(), keeper.ValAddrs[0])

	prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err := h(input.Ctx.WithBlockHeight(0), prevoteMsg)
	require.NoError(t, err)

	// Vote and new Prevote
	voteMsg := types.NewMsgAggregateExchangeRateVote(salt, exchangeRates.String(), keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	_, err = input.OracleKeeper.GetUSDExchangeRate(input.Ctx.WithBlockHeight(1), types.MicroSDRDenom)
	require.NotNil(t, err)

	// More than the threshold signs, msg succeeds
	salt = "1"
	hash = types.GetAggregateVoteHash(salt, exchangeRates.String(), keeper.ValAddrs[0])

	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRates.String(), keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "2"
	hash = types.GetAggregateVoteHash(salt, exchangeRates.String(), keeper.ValAddrs[1])

	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRates.String(), keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "3"
	hash = types.GetAggregateVoteHash(salt, exchangeRates.String(), keeper.ValAddrs[2])

	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[2], keeper.ValAddrs[2])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRates.String(), keeper.Addrs[2], keeper.ValAddrs[2])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rate, err := input.OracleKeeper.GetUSDExchangeRate(input.Ctx.WithBlockHeight(1), types.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, randomExchangeRate, rate)

	val, _ := input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[2])
	input.StakingKeeper.Delegate(input.Ctx.WithBlockHeight(0), keeper.Addrs[2], stakingAmt.MulRaw(3), stakingtypes.Unbonded, val, false)

	salt = "1"
	hash = types.GetAggregateVoteHash(salt, exchangeRates.String(), keeper.ValAddrs[0])

	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRates.String(), keeper.Addrs[0], keeper.ValAddrs[0])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	salt = "2"
	hash = types.GetAggregateVoteHash(salt, exchangeRates.String(), keeper.ValAddrs[1])

	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(0), prevoteMsg)

	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, exchangeRates.String(), keeper.Addrs[1], keeper.ValAddrs[1])
	h(input.Ctx.WithBlockHeight(1), voteMsg)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rate, err = input.OracleKeeper.GetUSDExchangeRate(input.Ctx.WithBlockHeight(1), types.MicroSDRDenom)
	require.NotNil(t, err)
}

func TestOracleMultiVote(t *testing.T) {
	input, h := setup(t)

	// Less than the threshold signs, msg fails
	salt := "1"
	exchangeRates := sdk.DecCoins{sdk.NewDecCoinFromDec(types.MicroSDRDenom, randomExchangeRate)}
	anotherExchangeRates := sdk.DecCoins{sdk.NewDecCoinFromDec(types.MicroSDRDenom, anotherRandomExchangeRate)}
	hash := types.GetAggregateVoteHash(salt, exchangeRates.String(), keeper.ValAddrs[0])

	prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err := h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	hash = types.GetAggregateVoteHash(salt, exchangeRates.String(), keeper.ValAddrs[1])

	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[1])
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	hash = types.GetAggregateVoteHash(salt, exchangeRates.String(), keeper.ValAddrs[2])

	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[2], keeper.ValAddrs[2])
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	hash = types.GetAggregateVoteHash(salt, anotherExchangeRates.String(), keeper.ValAddrs[0])

	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	hash = types.GetAggregateVoteHash(salt, anotherExchangeRates.String(), keeper.ValAddrs[1])

	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[1])
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	hash = types.GetAggregateVoteHash(salt, anotherExchangeRates.String(), keeper.ValAddrs[2])

	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[2], keeper.ValAddrs[2])
	_, err = h(input.Ctx, prevoteMsg)
	require.NoError(t, err)

	// Reveal ExchangeRate
	input.Ctx = input.Ctx.WithBlockHeight(1)
	voteMsg := types.NewMsgAggregateExchangeRateVote(salt, anotherExchangeRates.String(), keeper.Addrs[0], keeper.ValAddrs[0])
	_, err = h(input.Ctx, voteMsg)
	require.NoError(t, err)

	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, anotherExchangeRates.String(), keeper.Addrs[1], keeper.ValAddrs[1])
	_, err = h(input.Ctx, voteMsg)
	require.NoError(t, err)

	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, anotherExchangeRates.String(), keeper.Addrs[2], keeper.ValAddrs[2])
	_, err = h(input.Ctx, voteMsg)
	require.NoError(t, err)

	EndBlocker(input.Ctx, input.OracleKeeper)

	rate, err := input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroSDRDenom)
	require.Nil(t, err)
	require.Equal(t, rate, anotherRandomExchangeRate)
}

func TestOracleDrop(t *testing.T) {
	input, h := setup(t)

	input.OracleKeeper.SetUSDExchangeRate(input.Ctx, types.MicroKRWDenom, randomExchangeRate)

	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 0)

	// Immediately swap halt after an illiquid oracle vote
	EndBlocker(input.Ctx, input.OracleKeeper)

	_, err := input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroKRWDenom)
	require.NotNil(t, err)
}

func TestOracleTally(t *testing.T) {
	input, _ := setup(t)

	ballot := types.ExchangeRateBallot{}
	rates, valAddrs, stakingKeeper := types.GenerateRandomTestCase()
	input.OracleKeeper.StakingKeeper = stakingKeeper
	h := NewHandler(input.OracleKeeper)
	for i, rate := range rates {

		decExchangeRate := sdk.NewDecWithPrec(int64(rate*math.Pow10(keeper.OracleDecPrecision)), int64(keeper.OracleDecPrecision))

		salt := fmt.Sprintf("%d", i)
		exchangeRates := sdk.DecCoins{sdk.NewDecCoinFromDec(types.MicroSDRDenom, decExchangeRate)}
		hash := types.GetAggregateVoteHash(salt, exchangeRates.String(), valAddrs[i])

		prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(
			hash,
			sdk.AccAddress(valAddrs[i]),
			valAddrs[i],
		)

		_, err := h(input.Ctx.WithBlockHeight(0), prevoteMsg)
		require.NoError(t, err)

		voteMsg := types.NewMsgAggregateExchangeRateVote(
			salt,
			exchangeRates.String(),
			sdk.AccAddress(valAddrs[i]),
			valAddrs[i],
		)

		_, err = h(input.Ctx.WithBlockHeight(1), voteMsg)
		require.NoError(t, err)

		vote := types.NewVoteForTally(types.NewExchangeRateVote(decExchangeRate, types.MicroSDRDenom, valAddrs[i]), stakingAmt.QuoRaw(types.MicroUnit).Int64())
		ballot = append(ballot, vote)

		// change power of every three validator
		if i%3 == 0 {
			stakingKeeper.Validators()[i].SetPower(int64(i + 1))
		}
	}

	rewardees := []sdk.AccAddress{}
	weightedMedian := ballot.WeightedMedian()
	standardDeviation := ballot.StandardDeviation()
	maxSpread := weightedMedian.Mul(input.OracleKeeper.RewardBand(input.Ctx).QuoInt64(2))

	if standardDeviation.GT(maxSpread) {
		maxSpread = standardDeviation
	}

	for _, vote := range ballot {
		if vote.ExchangeRate.GTE(weightedMedian.Sub(maxSpread)) && vote.ExchangeRate.LTE(weightedMedian.Add(maxSpread)) {
			rewardees = append(rewardees, sdk.AccAddress(vote.Voter))
		}
	}

	tallyMedian, ballotWinner := tally(ballot, input.OracleKeeper.RewardBand(input.Ctx))

	require.Equal(t, len(rewardees), len(ballotWinner))
	require.Equal(t, tallyMedian.MulInt64(100).TruncateInt(), weightedMedian.MulInt64(100).TruncateInt())
}

func TestOracleTallyTiming(t *testing.T) {
	input, h := setup(t)

	// all the keeper.Addrs vote for the block ... not last period block yet, so tally fails
	for i := range keeper.Addrs[:2] {
		makePrevoteAndVote(t, input, h, 0, types.MicroSDRDenom, randomExchangeRate, i)
	}

	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 10 // set vote period to 10 for now, for convinience
	input.OracleKeeper.SetParams(input.Ctx, params)
	require.Equal(t, 0, int(input.Ctx.BlockHeight()))

	EndBlocker(input.Ctx, input.OracleKeeper)
	_, err := input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroSDRDenom)
	require.Error(t, err)

	input.Ctx = input.Ctx.WithBlockHeight(params.VotePeriod - 1)

	EndBlocker(input.Ctx, input.OracleKeeper)
	_, err = input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroSDRDenom)
	require.NoError(t, err)
}

func TestOracleRewardDistribution(t *testing.T) {
	input, h := setup(t)

	// Account 1, SDR
	makePrevoteAndVote(t, input, h, 0, types.MicroSDRDenom, randomExchangeRate, 0)

	// Account 2, SDR
	makePrevoteAndVote(t, input, h, 0, types.MicroSDRDenom, randomExchangeRate, 1)

	// Account 3, SDR
	makePrevoteAndVote(t, input, h, 0, types.MicroSDRDenom, randomExchangeRate, 2)

	rewardsAmt := sdk.NewInt(100000000)
	moduleAcc := input.AccKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), types.ModuleName)
	require.NoError(t, input.BankKeeper.SetBalances(input.Ctx, moduleAcc.GetAddress(), sdk.NewCoins(sdk.NewCoin(types.MicroLunaDenom, rewardsAmt))))

	input.AccKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.RewardDistributionWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	expectedRewardAmt := sdk.NewDecFromInt(rewardsAmt.QuoRaw(3)).QuoInt64(votePeriodsPerWindow).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
}

func TestOracleRewardBand(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = sdk.DecCoins{{types.MicroKRWDenom, types.DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, types.MicroKRWDenom, types.DefaultTobinTax)

	rewardSpread := randomExchangeRate.Mul(input.OracleKeeper.RewardBand(input.Ctx).QuoInt64(2))
	// no one will miss the vote
	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate.Sub(rewardSpread), 0)

	// Account 2, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 1)

	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate.Add(rewardSpread), 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// Account 1 will miss the vote due to raward band condition
	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate.Sub(rewardSpread.Add(sdk.OneDec())), 0)

	// Account 2, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 1)

	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate.Add(rewardSpread), 2)

	EndBlocker(input.Ctx, input.OracleKeeper)
	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

}

func TestOracleMultiRewardDistribution(t *testing.T) {
	input, h := setup(t)

	// SDR and KRW have the same voting power, but KRW has been chosen as referenceTerra by alphabetical order.
	// Account 1, SDR
	cnz := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec(types.MicroSDRDenom, randomExchangeRate),
		sdk.NewDecCoinFromDec(types.MicroKRWDenom, randomExchangeRate),
	)

	// Account 1, KRW+SDR
	makeAggregatePrevoteAndVote(t, input, h, 0, cnz, 0)

	// Account 2, SDR
	makePrevoteAndVote(t, input, h, 0, types.MicroSDRDenom, randomExchangeRate, 1)

	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 2)

	rewardAmt := sdk.NewInt(100000000)
	moduleAcc := input.AccKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), types.ModuleName)
	require.NoError(t, input.BankKeeper.SetBalances(input.Ctx, moduleAcc.GetAddress(), sdk.NewCoins(sdk.NewCoin(types.MicroLunaDenom, rewardAmt))))

	input.AccKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rewardDistributedWindow := input.OracleKeeper.RewardDistributionWindow(input.Ctx)

	expectedRewardAmt := sdk.NewDecFromInt(rewardAmt.QuoRaw(3).MulRaw(2)).QuoInt64(rewardDistributedWindow).TruncateInt()
	expectedRewardAmt2 := sdk.ZeroInt() // even vote power is same KRW with SDR, KRW chosen referenceTerra because alphabetical order
	expectedRewardAmt3 := sdk.NewDecFromInt(rewardAmt.QuoRaw(3)).QuoInt64(rewardDistributedWindow).TruncateInt()

	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt2, rewards.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt3, rewards.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
}

func TestOracleExchangeRate(t *testing.T) {
	input, h := setup(t)

	krwRandomExchangeRate := sdk.NewDecWithPrec(1000000000, int64(6)).MulInt64(types.MicroUnit)
	uswRandomExchangeRate := sdk.NewDecWithPrec(1000000, int64(6)).MulInt64(types.MicroUnit)

	// KRW has been chosen as referenceTerra by highest voting power
	// Account 1, USD, KRW
	cnz := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec(types.MicroUSDDenom, uswRandomExchangeRate),
		sdk.NewDecCoinFromDec(types.MicroKRWDenom, krwRandomExchangeRate),
	)
	makeAggregatePrevoteAndVote(t, input, h, 0, cnz, 0)

	// Account 2, USD, KRW
	makeAggregatePrevoteAndVote(t, input, h, 0, cnz, 1)

	// Account 3, KRW, SDR
	cnz = sdk.NewDecCoins(
		sdk.NewDecCoinFromDec(types.MicroSDRDenom, randomExchangeRate),
		sdk.NewDecCoinFromDec(types.MicroKRWDenom, krwRandomExchangeRate),
	)
	makeAggregatePrevoteAndVote(t, input, h, 0, cnz, 2)

	rewardAmt := sdk.NewInt(100000000)
	moduleAcc := input.AccKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), types.ModuleName)
	require.NoError(t, input.BankKeeper.SetBalances(input.Ctx, moduleAcc.GetAddress(), sdk.NewCoins(sdk.NewCoin(types.MicroLunaDenom, rewardAmt))))

	input.AccKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rewardDistributedWindow := input.OracleKeeper.RewardDistributionWindow(input.Ctx)
	expectedRewardAmt := sdk.NewDecFromInt(rewardAmt.QuoRaw(5).MulRaw(2)).QuoInt64(rewardDistributedWindow).TruncateInt()
	expectedRewardAmt2 := sdk.NewDecFromInt(rewardAmt.QuoRaw(5).MulRaw(1)).QuoInt64(rewardDistributedWindow).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt2, rewards.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
}

func TestOracleExchangeRateVal5(t *testing.T) {
	input, h := setupVal5(t)

	krwExchangeRate := sdk.NewDecWithPrec(505000, int64(6)).MulInt64(types.MicroUnit)
	krwExchangeRateWithErr := sdk.NewDecWithPrec(500000, int64(6)).MulInt64(types.MicroUnit)
	lunaExchangeRate := sdk.NewDecWithPrec(505, int64(6)).MulInt64(types.MicroUnit)
	lunaExchangeRateWithErr := sdk.NewDecWithPrec(500, int64(6)).MulInt64(types.MicroUnit)

	// KRW has been chosen as referenceTerra by highest voting power
	// Account 1, KRW, USD
	val1 := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec(types.MicroKRWDenom, krwExchangeRate),
		sdk.NewDecCoinFromDec(types.MicroSDRDenom, lunaExchangeRate),
	)
	makeAggregatePrevoteAndVote(t, input, h, 0, val1, 0)

	// Account 2, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, krwExchangeRate, 1)

	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, krwExchangeRate, 2)

	// Account 4, KRW, USD
	val4 := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec(types.MicroKRWDenom, krwExchangeRateWithErr),
		sdk.NewDecCoinFromDec(types.MicroSDRDenom, lunaExchangeRateWithErr),
	)
	makeAggregatePrevoteAndVote(t, input, h, 0, val4, 3)

	// Account 5, KRW, USD
	makeAggregatePrevoteAndVote(t, input, h, 0, val4, 4)

	rewardAmt := sdk.NewInt(100000000)
	moduleAcc := input.AccKeeper.GetModuleAccount(input.Ctx.WithBlockHeight(1), types.ModuleName)
	require.NoError(t, input.BankKeeper.SetBalances(input.Ctx, moduleAcc.GetAddress(), sdk.NewCoins(sdk.NewCoin(types.MicroSDRDenom, rewardAmt))))

	input.AccKeeper.SetModuleAccount(input.Ctx.WithBlockHeight(1), moduleAcc)

	EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	krw, err := input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroKRWDenom)
	require.NoError(t, err)
	usd, err := input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroSDRDenom)
	require.NoError(t, err)

	// legacy version case
	require.NotEqual(t, lunaExchangeRateWithErr, usd)

	// new version case
	require.Equal(t, krwExchangeRate, krw)
	require.Equal(t, lunaExchangeRate, usd)

	rewardDistributedWindow := input.OracleKeeper.RewardDistributionWindow(input.Ctx)
	expectedRewardAmt := sdk.NewDecFromInt(rewardAmt.QuoRaw(8).MulRaw(2)).QuoInt64(rewardDistributedWindow).TruncateInt()
	expectedRewardAmt2 := sdk.NewDecFromInt(rewardAmt.QuoRaw(8).MulRaw(1)).QuoInt64(rewardDistributedWindow).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
	rewards1 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt2, rewards1.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
	rewards2 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt2, rewards2.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
	rewards3 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[3])
	require.Equal(t, expectedRewardAmt, rewards3.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
	rewards4 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[4])
	require.Equal(t, expectedRewardAmt, rewards4.Rewards.AmountOf(types.MicroLunaDenom).TruncateInt())
}

func TestInvalidVotesSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = sdk.DecCoins{sdk.NewDecCoinFromDec(types.MicroKRWDenom, types.DefaultTobinTax)}
	input.OracleKeeper.SetParams(input.Ctx, params)
	input.OracleKeeper.SetTobinTax(input.Ctx, types.MicroKRWDenom, types.DefaultTobinTax)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.SlashWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := int64(0); i < sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64(); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 1, KRW
		makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 0)

		// Account 2, KRW, miss vote
		makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate.Add(sdk.NewDec(100000000000000)), 1)

		// Account 3, KRW
		makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 2)

		EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, i+1, input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// one more miss vote will inccur ValAddrs[1] slashing
	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 0)

	// Account 2, KRW, miss vote
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate.Add(sdk.NewDec(100000000000000)), 1)

	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 2)

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
		makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 1)
		// Account 3, KRW
		makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 2)

		EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, i+1, input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// one more miss vote will inccur Account 1 slashing

	// Account 2, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 1)
	// Account 3, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 2)

	input.Ctx = input.Ctx.WithBlockHeight(votePeriodsPerWindow - 1)
	EndBlocker(input.Ctx, input.OracleKeeper)
	validator = input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[0])
	require.Equal(t, sdk.OneDec().Sub(slashFraction).MulInt(stakingAmt).TruncateInt(), validator.GetBondedTokens())
}

func TestNotPassedBallotSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = sdk.DecCoins{{types.MicroKRWDenom, types.DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, types.MicroKRWDenom, types.DefaultTobinTax)

	input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

	// Account 1, KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 0)

	EndBlocker(input.Ctx, input.OracleKeeper)
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))
}

func TestAbstainSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = sdk.DecCoins{{types.MicroKRWDenom, types.DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, types.MicroKRWDenom, types.DefaultTobinTax)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.SlashWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := int64(0); i <= sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64(); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 1, KRW
		makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 0)

		// Account 2, KRW, abstain vote
		makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, sdk.ZeroDec(), 1)

		// Account 3, KRW
		makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 2)

		EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())
}

func TestVoteTargets(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.Whitelist = sdk.DecCoins{{types.MicroKRWDenom, types.DefaultTobinTax}, {types.MicroSDRDenom, types.DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, types.MicroKRWDenom, types.DefaultTobinTax)

	// KRW
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 0)
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 1)
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	// no missing current
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// vote targets are {KRW, SDR}
	require.Equal(t, []string{types.MicroKRWDenom, types.MicroSDRDenom}, input.OracleKeeper.GetVoteTargets(input.Ctx))

	// tobin tax must be exists for SDR
	sdrTobinTax, err := input.OracleKeeper.GetTobinTax(input.Ctx, types.MicroSDRDenom)
	require.NoError(t, err)
	require.Equal(t, types.DefaultTobinTax, sdrTobinTax)

	// delete SDR
	params.Whitelist = sdk.DecCoins{{types.MicroKRWDenom, types.DefaultTobinTax}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// KRW, missing
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 0)
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 1)
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// SDR must be deleted
	require.Equal(t, []string{types.MicroKRWDenom}, input.OracleKeeper.GetVoteTargets(input.Ctx))

	_, err = input.OracleKeeper.GetTobinTax(input.Ctx, types.MicroSDRDenom)
	require.Error(t, err)

	// change KRW tobin tax
	params.Whitelist = sdk.DecCoins{{types.MicroKRWDenom, sdk.ZeroDec()}}
	input.OracleKeeper.SetParams(input.Ctx, params)

	// KRW, no missing
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 0)
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 1)
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, randomExchangeRate, 2)

	EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, int64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// KRW tobin tax must be 0
	tobinTax, err := input.OracleKeeper.GetTobinTax(input.Ctx, types.MicroKRWDenom)
	require.NoError(t, err)
	require.True(t, sdk.ZeroDec().Equal(tobinTax))
}

func TestAbstainWithSmallStakingPower(t *testing.T) {
	input, h := setup_with_small_voting_power(t)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	input.OracleKeeper.SetTobinTax(input.Ctx, types.MicroKRWDenom, types.DefaultTobinTax)
	makePrevoteAndVote(t, input, h, 0, types.MicroKRWDenom, sdk.ZeroDec(), 0)

	EndBlocker(input.Ctx, input.OracleKeeper)
	_, err := input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroKRWDenom)
	require.Error(t, err)
}

func makePrevoteAndVote(t *testing.T, input keeper.TestInput, h sdk.Handler, height int64, denom string, rate sdk.Dec, idx int) {
	// Account 1, SDR
	salt := "1"
	exchangeRates := sdk.DecCoins{sdk.NewDecCoinFromDec(denom, rate)}
	hash := types.GetAggregateVoteHash(salt, exchangeRates.String(), keeper.ValAddrs[idx])

	prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[idx], keeper.ValAddrs[idx])
	_, err := h(input.Ctx.WithBlockHeight(height), prevoteMsg)
	require.NoError(t, err)

	voteMsg := types.NewMsgAggregateExchangeRateVote(salt, exchangeRates.String(), keeper.Addrs[idx], keeper.ValAddrs[idx])
	_, err = h(input.Ctx.WithBlockHeight(height+1), voteMsg)
	require.NoError(t, err)
}

func makeAggregatePrevoteAndVote(t *testing.T, input keeper.TestInput, h sdk.Handler, height int64, rates sdk.DecCoins, idx int) {
	// Account 1, SDR
	salt := "1"
	hash := types.GetAggregateVoteHash(salt, rates.String(), keeper.ValAddrs[idx])

	prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[idx], keeper.ValAddrs[idx])
	_, err := h(input.Ctx.WithBlockHeight(height), prevoteMsg)
	require.NoError(t, err)

	voteMsg := types.NewMsgAggregateExchangeRateVote(salt, rates.String(), keeper.Addrs[idx], keeper.ValAddrs[idx])
	_, err = h(input.Ctx.WithBlockHeight(height+1), voteMsg)
	require.NoError(t, err)
}
