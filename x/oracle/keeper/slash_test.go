package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestSlashAndResetMissCounters(t *testing.T) {
	stakingAmt := sdk.TokensFromConsensusPower(10)
	input := CreateTestInput(t)

	votePeriodsPerWindow := sdk.NewDec(input.OracleKeeper.SlashWindow(input.Ctx)).QuoInt64(input.OracleKeeper.VotePeriod(input.Ctx)).TruncateInt64()
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidVotes := input.OracleKeeper.MinValidPerWindow(input.Ctx).MulInt64(votePeriodsPerWindow).TruncateInt64()
	// Case 1, no slash
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], votePeriodsPerWindow-minValidVotes)
	input.OracleKeeper.SlashAndResetMissCounters(input.Ctx)
	staking.EndBlocker(input.Ctx, input.StakingKeeper)

	validator, _ := input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// Case 2, slash
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], votePeriodsPerWindow-minValidVotes+1)
	input.OracleKeeper.SlashAndResetMissCounters(input.Ctx)
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	require.Equal(t, stakingAmt.Sub(slashFraction.MulInt(stakingAmt).TruncateInt()), validator.GetBondedTokens())
	require.True(t, validator.IsJailed())

	// Case 3, slash unbonded validator
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	validator.Status = stakingtypes.Unbonded
	validator.Jailed = false
	validator.Tokens = stakingAmt
	input.StakingKeeper.SetValidator(input.Ctx, validator)

	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], votePeriodsPerWindow-minValidVotes+1)
	input.OracleKeeper.SlashAndResetMissCounters(input.Ctx)
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	require.Equal(t, stakingAmt, validator.Tokens)
	require.False(t, validator.IsJailed())

	// Case 4, slash jailed validator
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	validator.Status = stakingtypes.Bonded
	validator.Jailed = true
	validator.Tokens = stakingAmt
	input.StakingKeeper.SetValidator(input.Ctx, validator)

	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], votePeriodsPerWindow-minValidVotes+1)
	input.OracleKeeper.SlashAndResetMissCounters(input.Ctx)
	validator, _ = input.StakingKeeper.GetValidator(input.Ctx, ValAddrs[0])
	require.Equal(t, stakingAmt, validator.Tokens)
}
