package keeper

import (
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestOrganize(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// TODO: Do this initialization in the test_utils file?
	bd := input.StakingKeeper.GetParams(ctx).BondDenom
	for i := range []int{0, 1, 2} {
		acc := input.AccKeeper.NewAccount(ctx, authtypes.NewBaseAccount(Addrs[i], AccPubKeys[i], uint64(i), 0))
		input.BankKeeper.SetBalances(ctx, acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(bd, amt)))
		input.AccKeeper.SetAccount(ctx, acc)
	}

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], PubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], PubKeys[1], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], PubKeys[2], amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	sdrBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(17), types.MicroSDRDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(10), types.MicroSDRDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(6), types.MicroSDRDenom, ValAddrs[2]), power),
	}
	krwBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1000), types.MicroKRWDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1300), types.MicroKRWDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(2000), types.MicroKRWDenom, ValAddrs[2]), power),
	}

	for _, vote := range sdrBallot {
		input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
	}
	for _, vote := range krwBallot {
		input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
	}

	// organize votes by denom
	ballotMap := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx)

	// sort each ballot for comparison
	sort.Sort(sdrBallot)
	sort.Sort(krwBallot)
	sort.Sort(ballotMap[types.MicroSDRDenom])
	sort.Sort(ballotMap[types.MicroKRWDenom])

	require.Equal(t, sdrBallot, ballotMap[types.MicroSDRDenom])
	require.Equal(t, krwBallot, ballotMap[types.MicroKRWDenom])
}

func TestOrganizeAggregate(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// TODO: Do this initialization in the test_utils file?
	bd := input.StakingKeeper.GetParams(ctx).BondDenom
	for i := range []int{0, 1, 2} {
		acc := input.AccKeeper.NewAccount(ctx, authtypes.NewBaseAccount(Addrs[i], AccPubKeys[i], uint64(i), 0))
		input.BankKeeper.SetBalances(ctx, acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(bd, amt)))
		input.AccKeeper.SetAccount(ctx, acc)
	}

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], PubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], PubKeys[1], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], PubKeys[2], amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	sdrBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(17), types.MicroSDRDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(10), types.MicroSDRDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(6), types.MicroSDRDenom, ValAddrs[2]), power),
	}
	krwBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1000), types.MicroKRWDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1300), types.MicroKRWDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(2000), types.MicroKRWDenom, ValAddrs[2]), power),
	}

	for i := range sdrBallot {
		input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
			{Denom: sdrBallot[i].Denom, ExchangeRate: sdrBallot[i].ExchangeRate},
			{Denom: krwBallot[i].Denom, ExchangeRate: krwBallot[i].ExchangeRate},
		}, ValAddrs[i]))
	}

	// organize votes by denom
	ballotMap := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx)

	// sort each ballot for comparison
	sort.Sort(sdrBallot)
	sort.Sort(krwBallot)
	sort.Sort(ballotMap[types.MicroSDRDenom])
	sort.Sort(ballotMap[types.MicroKRWDenom])

	require.Equal(t, sdrBallot, ballotMap[types.MicroSDRDenom])
	require.Equal(t, krwBallot, ballotMap[types.MicroKRWDenom])
}

func TestDuplicateVote(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// TODO: Do this initialization in the test_utils file?
	bd := input.StakingKeeper.GetParams(ctx).BondDenom
	for i := range []int{0, 1, 2} {
		acc := input.AccKeeper.NewAccount(ctx, authtypes.NewBaseAccount(Addrs[i], AccPubKeys[i], uint64(i), 0))
		input.BankKeeper.SetBalances(ctx, acc.GetAddress(), sdk.NewCoins(sdk.NewCoin(bd, amt)))
		input.AccKeeper.SetAccount(ctx, acc)
	}

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], PubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], PubKeys[1], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], PubKeys[2], amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	sdrBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(17), types.MicroSDRDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(10), types.MicroSDRDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(6), types.MicroSDRDenom, ValAddrs[2]), power),
	}
	krwBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1000), types.MicroKRWDenom, ValAddrs[0]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(1300), types.MicroKRWDenom, ValAddrs[1]), power),
		types.NewVoteForTally(types.NewExchangeRateVote(sdk.NewDec(2000), types.MicroKRWDenom, ValAddrs[2]), power),
	}

	for i := range sdrBallot {

		// this will be ignored
		for _, vote := range sdrBallot {
			input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
		}
		for _, vote := range krwBallot {
			input.OracleKeeper.AddExchangeRateVote(input.Ctx, vote.ExchangeRateVote)
		}

		input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
			{Denom: sdrBallot[i].Denom, ExchangeRate: sdrBallot[i].ExchangeRate},
			{Denom: krwBallot[i].Denom, ExchangeRate: krwBallot[i].ExchangeRate},
		}, ValAddrs[i]))
	}

	// organize votes by denom
	ballotMap := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx)

	// sort each ballot for comparison
	sort.Sort(sdrBallot)
	sort.Sort(krwBallot)
	sort.Sort(ballotMap[types.MicroSDRDenom])
	sort.Sort(ballotMap[types.MicroKRWDenom])

	require.Equal(t, sdrBallot, ballotMap[types.MicroSDRDenom])
	require.Equal(t, krwBallot, ballotMap[types.MicroKRWDenom])
}
