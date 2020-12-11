package keeper

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestVoteAddDelete(t *testing.T) {
	input := CreateTestInput(t)

	rate := sdk.NewDec(1700)
	vote := types.NewExchangeRateVote(rate, types.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote)

	KVote, err := input.OracleKeeper.getExchangeRateVote(input.Ctx, types.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	require.NoError(t, err)
	require.Equal(t, vote, KVote)

	input.OracleKeeper.DeleteExchangeRateVote(input.Ctx, vote)
	_, err = input.OracleKeeper.getExchangeRateVote(input.Ctx, types.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	require.Error(t, err)
}

func TestVoteIterate(t *testing.T) {
	input := CreateTestInput(t)

	rate := sdk.NewDec(1700)
	vote1 := types.NewExchangeRateVote(rate, types.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote1)

	vote2 := types.NewExchangeRateVote(rate, types.MicroSDRDenom, sdk.ValAddress(Addrs[1]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote2)

	i := 0
	bigger := bytes.Compare(Addrs[0], Addrs[1])
	input.OracleKeeper.IterateExchangeRateVotes(input.Ctx, func(p types.ExchangeRateVote) (stop bool) {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, vote1, p)
		} else {
			require.Equal(t, vote2, p)
		}

		i++
		return false
	})

	vote3 := types.NewExchangeRateVote(rate, types.MicroLunaDenom, sdk.ValAddress(Addrs[2]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote3)

	input.OracleKeeper.iterateExchangeRateVotesWithPrefix(input.Ctx, types.GetVoteKey(types.MicroLunaDenom, sdk.ValAddress{}), func(p types.ExchangeRateVote) (stop bool) {
		require.Equal(t, vote3, p)
		return false
	})
}

func TestVoteCollect(t *testing.T) {
	input := CreateTestInput(t)

	rate := sdk.NewDec(1700)
	vote1 := types.NewExchangeRateVote(rate, types.MicroSDRDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote1)

	vote2 := types.NewExchangeRateVote(rate, types.MicroSDRDenom, sdk.ValAddress(Addrs[1]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote2)

	vote3 := types.NewExchangeRateVote(rate, types.MicroLunaDenom, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote3)

	vote4 := types.NewExchangeRateVote(rate, types.MicroLunaDenom, sdk.ValAddress(Addrs[1]))
	input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote4)

	collectedVotes := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx)

	pb1 := collectedVotes[types.MicroSDRDenom]
	pb2 := collectedVotes[types.MicroLunaDenom]

	bigger := bytes.Compare(Addrs[0], Addrs[1])
	for i, v := range pb1 {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, vote1, v)
		} else {
			require.Equal(t, vote2, v)
		}
	}

	for i, v := range pb2 {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, vote3, v)
		} else {
			require.Equal(t, vote4, v)
		}
	}
}

func TestExchangeRate(t *testing.T) {
	input := CreateTestInput(t)

	cnyExchangeRate := sdk.NewDecWithPrec(839, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	gbpExchangeRate := sdk.NewDecWithPrec(4995, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	krwExchangeRate := sdk.NewDecWithPrec(2838, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	lunaExchangeRate := sdk.NewDecWithPrec(3282384, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)

	// Set & get rates
	input.OracleKeeper.SetUSDExchangeRate(input.Ctx, types.MicroCNYDenom, cnyExchangeRate)
	rate, err := input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroCNYDenom)
	require.NoError(t, err)
	require.Equal(t, cnyExchangeRate, rate)

	input.OracleKeeper.SetUSDExchangeRate(input.Ctx, types.MicroGBPDenom, gbpExchangeRate)
	rate, err = input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroGBPDenom)
	require.NoError(t, err)
	require.Equal(t, gbpExchangeRate, rate)

	input.OracleKeeper.SetUSDExchangeRate(input.Ctx, types.MicroKRWDenom, krwExchangeRate)
	rate, err = input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroKRWDenom)
	require.NoError(t, err)
	require.Equal(t, krwExchangeRate, rate)

	input.OracleKeeper.SetUSDExchangeRate(input.Ctx, types.MicroUSDDenom, lunaExchangeRate)
	rate, _ = input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroUSDDenom)
	require.Equal(t, sdk.OneDec(), rate)

	input.OracleKeeper.DeleteUSDExchangeRate(input.Ctx, types.MicroKRWDenom)
	_, err = input.OracleKeeper.GetUSDExchangeRate(input.Ctx, types.MicroKRWDenom)
	require.Error(t, err)

	numExchangeRates := 0
	handler := func(denom string, exchangeRate sdk.Dec) (stop bool) {
		numExchangeRates = numExchangeRates + 1
		return false
	}
	input.OracleKeeper.IterateUSDExchangeRates(input.Ctx, handler)

	require.True(t, numExchangeRates == 3)
}

func TestIterateUSDExchangeRates(t *testing.T) {
	input := CreateTestInput(t)

	cnyExchangeRate := sdk.NewDecWithPrec(839, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	gbpExchangeRate := sdk.NewDecWithPrec(4995, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	krwExchangeRate := sdk.NewDecWithPrec(2838, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)
	lunaExchangeRate := sdk.NewDecWithPrec(3282384, int64(OracleDecPrecision)).MulInt64(types.MicroUnit)

	// Set & get rates
	input.OracleKeeper.SetUSDExchangeRate(input.Ctx, types.MicroCNYDenom, cnyExchangeRate)
	input.OracleKeeper.SetUSDExchangeRate(input.Ctx, types.MicroGBPDenom, gbpExchangeRate)
	input.OracleKeeper.SetUSDExchangeRate(input.Ctx, types.MicroKRWDenom, krwExchangeRate)
	input.OracleKeeper.SetUSDExchangeRate(input.Ctx, types.MicroLunaDenom, lunaExchangeRate)

	input.OracleKeeper.IterateUSDExchangeRates(input.Ctx, func(denom string, rate sdk.Dec) (stop bool) {
		switch denom {
		case types.MicroCNYDenom:
			require.Equal(t, cnyExchangeRate, rate)
		case types.MicroGBPDenom:
			require.Equal(t, gbpExchangeRate, rate)
		case types.MicroKRWDenom:
			require.Equal(t, krwExchangeRate, rate)
		case types.MicroLunaDenom:
			require.Equal(t, lunaExchangeRate, rate)
		}
		return false
	})

}

func TestRewardPool(t *testing.T) {
	input := CreateTestInput(t)

	fees := sdk.NewCoins(sdk.NewCoin(types.MicroSDRDenom, sdk.NewInt(1000)))
	acc := input.AccKeeper.GetModuleAccount(input.Ctx, types.ModuleName)
	err := input.BankKeeper.SetBalances(input.Ctx, acc.GetAddress(), fees)
	if err != nil {
		panic(err) // nerver occurs
	}

	input.AccKeeper.SetModuleAccount(input.Ctx, acc)

	KFees := input.OracleKeeper.GetRewardPool(input.Ctx)
	require.Equal(t, fees, KFees)
}

func TestParams(t *testing.T) {
	input := CreateTestInput(t)

	// Test default params setting
	input.OracleKeeper.SetParams(input.Ctx, types.DefaultParams())
	params := input.OracleKeeper.GetParams(input.Ctx)
	require.NotNil(t, params)

	// Test custom params setting
	votePeriod := int64(10)
	voteThreshold := sdk.NewDecWithPrec(33, 2)
	oracleRewardBand := sdk.NewDecWithPrec(1, 2)
	rewardDistributionWindow := int64(10000000000000)
	slashFraction := sdk.NewDecWithPrec(1, 2)
	slashWindow := int64(1000)
	minValidPerWindow := sdk.NewDecWithPrec(1, 4)
	whitelist := sdk.NewDecCoins(
		sdk.NewDecCoinFromDec(types.MicroSDRDenom, types.DefaultTobinTax),
		sdk.NewDecCoinFromDec(types.MicroKRWDenom, types.DefaultTobinTax),
	)

	// Should really test validateParams, but skipping because obvious
	newParams := types.Params{
		VotePeriod:               votePeriod,
		VoteThreshold:            voteThreshold,
		RewardBand:               oracleRewardBand,
		RewardDistributionWindow: rewardDistributionWindow,
		Whitelist:                whitelist,
		SlashFraction:            slashFraction,
		SlashWindow:              slashWindow,
		MinValidPerWindow:        minValidPerWindow,
	}
	input.OracleKeeper.SetParams(input.Ctx, newParams)

	storedParams := input.OracleKeeper.GetParams(input.Ctx)
	require.NotNil(t, storedParams)
	require.Equal(t, storedParams, newParams)
}

func TestFeederDelegation(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	delegate := input.OracleKeeper.GetOracleDelegate(input.Ctx, ValAddrs[0])
	require.Equal(t, Addrs[0], delegate)

	input.OracleKeeper.SetOracleDelegate(input.Ctx, ValAddrs[0], Addrs[1])
	delegate = input.OracleKeeper.GetOracleDelegate(input.Ctx, ValAddrs[0])
	require.Equal(t, Addrs[1], delegate)
}

func TestIterateFeederDelegations(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	delegate := input.OracleKeeper.GetOracleDelegate(input.Ctx, ValAddrs[0])
	require.Equal(t, Addrs[0], delegate)

	input.OracleKeeper.SetOracleDelegate(input.Ctx, ValAddrs[0], Addrs[1])

	var delegators []sdk.ValAddress
	var delegates []sdk.AccAddress
	input.OracleKeeper.IterateOracleDelegates(input.Ctx, func(delegator sdk.ValAddress, delegate sdk.AccAddress) (stop bool) {
		delegators = append(delegators, delegator)
		delegates = append(delegates, delegate)
		return false
	})

	require.Equal(t, 1, len(delegators))
	require.Equal(t, 1, len(delegates))
	require.Equal(t, ValAddrs[0], delegators[0])
	require.Equal(t, Addrs[1], delegates[0])
}

func TestMissCounter(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	counter := input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, int64(0), counter)

	missCounter := int64(10)
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], missCounter)
	counter = input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, missCounter, counter)

	input.OracleKeeper.DeleteMissCounter(input.Ctx, ValAddrs[0])
	counter = input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, int64(0), counter)
}

func TestIterateMissCounters(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	counter := input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, int64(0), counter)

	missCounter := int64(10)
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[1], missCounter)

	var operators []sdk.ValAddress
	var missCounters []int64
	input.OracleKeeper.IterateMissCounters(input.Ctx, func(delegator sdk.ValAddress, missCounter int64) (stop bool) {
		operators = append(operators, delegator)
		missCounters = append(missCounters, missCounter)
		return false
	})

	require.Equal(t, 1, len(operators))
	require.Equal(t, 1, len(missCounters))
	require.Equal(t, ValAddrs[1], operators[0])
	require.Equal(t, missCounter, missCounters[0])
}

func TestAggregatePrevoteAddDelete(t *testing.T) {
	input := CreateTestInput(t)

	hash := types.GetAggregateVoteHash("salt", "100ukrw,1000uusd", sdk.ValAddress(Addrs[0]))
	aggregatePrevote := types.NewAggregateExchangeRatePrevote(hash, sdk.ValAddress(Addrs[0]), 0)
	input.OracleKeeper.AddAggregateExchangeRatePrevote(input.Ctx, aggregatePrevote)

	KPrevote, err := input.OracleKeeper.GetAggregateExchangeRatePrevote(input.Ctx, sdk.ValAddress(Addrs[0]))
	require.NoError(t, err)
	require.Equal(t, aggregatePrevote, KPrevote)

	input.OracleKeeper.DeleteAggregateExchangeRatePrevote(input.Ctx, aggregatePrevote)
	_, err = input.OracleKeeper.GetAggregateExchangeRatePrevote(input.Ctx, sdk.ValAddress(Addrs[0]))
	require.Error(t, err)
}

func TestAggregatePrevoteIterate(t *testing.T) {
	input := CreateTestInput(t)

	hash := types.GetAggregateVoteHash("salt", "100ukrw,1000uusd", sdk.ValAddress(Addrs[0]))
	aggregatePrevote1 := types.NewAggregateExchangeRatePrevote(hash, sdk.ValAddress(Addrs[0]), 0)
	input.OracleKeeper.AddAggregateExchangeRatePrevote(input.Ctx, aggregatePrevote1)

	hash2 := types.GetAggregateVoteHash("salt", "100ukrw,1000uusd", sdk.ValAddress(Addrs[1]))
	aggregatePrevote2 := types.NewAggregateExchangeRatePrevote(hash2, sdk.ValAddress(Addrs[1]), 0)
	input.OracleKeeper.AddAggregateExchangeRatePrevote(input.Ctx, aggregatePrevote2)

	i := 0
	bigger := bytes.Compare(Addrs[0], Addrs[1])
	input.OracleKeeper.IterateAggregateExchangeRatePrevotes(input.Ctx, func(p types.AggregateExchangeRatePrevote) (stop bool) {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, aggregatePrevote1, p)
		} else {
			require.Equal(t, aggregatePrevote2, p)
		}

		i++
		return false
	})
}

func TestAggregateVoteAddDelete(t *testing.T) {
	input := CreateTestInput(t)

	aggregateVote := types.NewAggregateExchangeRateVote(sdk.DecCoins{
		sdk.DecCoin{"foo", sdk.NewDec(-1)},
		sdk.NewDecCoinFromDec("foo", sdk.NewDec(0)),
		sdk.NewDecCoinFromDec("foo", sdk.NewDec(1)),
	}, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, aggregateVote)

	KVote, err := input.OracleKeeper.GetAggregateExchangeRateVote(input.Ctx, sdk.ValAddress(Addrs[0]))
	require.NoError(t, err)
	require.Equal(t, aggregateVote, KVote)

	input.OracleKeeper.DeleteAggregateExchangeRateVote(input.Ctx, aggregateVote)
	_, err = input.OracleKeeper.GetAggregateExchangeRateVote(input.Ctx, sdk.ValAddress(Addrs[0]))
	require.Error(t, err)
}

func TestAggregateVoteIterate(t *testing.T) {
	input := CreateTestInput(t)

	aggregateVote1 := types.NewAggregateExchangeRateVote(sdk.DecCoins{
		sdk.DecCoin{"foo", sdk.NewDec(-1)},
		sdk.NewDecCoinFromDec("foo", sdk.NewDec(0)),
		sdk.NewDecCoinFromDec("foo", sdk.NewDec(1)),
	}, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, aggregateVote1)

	aggregateVote2 := types.NewAggregateExchangeRateVote(sdk.DecCoins{
		sdk.DecCoin{"foo", sdk.NewDec(-1)},
		sdk.NewDecCoinFromDec("foo", sdk.NewDec(0)),
		sdk.NewDecCoinFromDec("foo", sdk.NewDec(1)),
	}, sdk.ValAddress(Addrs[1]))
	input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, aggregateVote2)

	i := 0
	bigger := bytes.Compare(Addrs[0], Addrs[1])
	input.OracleKeeper.IterateAggregateExchangeRateVotes(input.Ctx, func(p types.AggregateExchangeRateVote) (stop bool) {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, aggregateVote1, p)
		} else {
			require.Equal(t, aggregateVote2, p)
		}

		i++
		return false
	})
}

func TestTobinTaxGetSet(t *testing.T) {
	input := CreateTestInput(t)

	tobinTaxes := map[string]sdk.Dec{
		types.MicroSDRDenom: sdk.NewDec(1),
		types.MicroUSDDenom: sdk.NewDecWithPrec(1, 3),
		types.MicroKRWDenom: sdk.NewDecWithPrec(123, 3),
		types.MicroMNTDenom: sdk.NewDecWithPrec(1423, 4),
	}

	for denom, tobinTax := range tobinTaxes {
		input.OracleKeeper.SetTobinTax(input.Ctx, denom, tobinTax)
		factor, err := input.OracleKeeper.GetTobinTax(input.Ctx, denom)
		require.NoError(t, err)
		require.Equal(t, tobinTaxes[denom], factor)
	}

	input.OracleKeeper.IterateTobinTaxes(input.Ctx, func(denom string, tobinTax sdk.Dec) (stop bool) {
		require.Equal(t, tobinTaxes[denom], tobinTax)
		return false
	})

	input.OracleKeeper.ClearTobinTaxes(input.Ctx)
	for denom := range tobinTaxes {
		_, err := input.OracleKeeper.GetTobinTax(input.Ctx, denom)
		require.Error(t, err)
	}
}
