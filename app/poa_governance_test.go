package app

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	poatypes "github.com/peggyjv/sommelier/v10/x/poa/types"
	"github.com/stretchr/testify/require"
)

func TestGovernanceTallyDoesNotGiveAuthorityBoostToExternalDelegator(t *testing.T) {
	app, ctx, authority, community := newPoASlashApp(t)
	authorityOperator := authority.GetOperator()
	operatorAccount := sdk.AccAddress(authorityOperator)
	externalDelegator := sdk.AccAddress(bytes.Repeat([]byte{0x43}, 20))

	// Build a 10-token authority validator whose operator owns 3 shares and an
	// external account owns 7. The 10x PoA multiplier must belong to the
	// authority operator; the permissionless delegator still votes with 7.
	authority.Tokens = sdk.NewInt(10_000_000)
	authority.DelegatorShares = sdk.NewDec(10_000_000)
	app.StakingKeeper.SetValidator(ctx, authority)
	app.StakingKeeper.SetDelegation(ctx, stakingtypes.NewDelegation(
		operatorAccount, authorityOperator, sdk.NewDec(3_000_000),
	))
	app.StakingKeeper.SetDelegation(ctx, stakingtypes.NewDelegation(
		externalDelegator, authorityOperator, sdk.NewDec(7_000_000),
	))
	communityDelegation, found := app.StakingKeeper.GetDelegation(
		ctx, operatorAccount, community.GetOperator(),
	)
	if found {
		require.NoError(t, app.StakingKeeper.RemoveDelegation(ctx, communityDelegation))
	}

	app.PoaKeeper.SetMultiplierSnapshot(ctx, poatypes.MultiplierSnapshot{
		Height: ctx.BlockHeight(),
		Entries: []*poatypes.MultiplierEntry{
			{OperatorAddress: authorityOperator.String(), Multiplier: "10.0"},
		},
	})

	const proposalID uint64 = 1
	app.GovKeeper.SetVote(ctx, govtypesv1.NewVote(
		proposalID, externalDelegator,
		govtypesv1.NewNonSplitVoteOption(govtypesv1.OptionYes), "",
	))
	app.GovKeeper.SetVote(ctx, govtypesv1.NewVote(
		proposalID, operatorAccount,
		govtypesv1.NewNonSplitVoteOption(govtypesv1.OptionNoWithVeto), "",
	))

	_, _, result := app.GovKeeper.Tally(ctx, govtypesv1.Proposal{Id: proposalID})
	yes, err := sdk.NewDecFromStr(result.YesCount)
	require.NoError(t, err)
	veto, err := sdk.NewDecFromStr(result.NoWithVetoCount)
	require.NoError(t, err)
	require.Equal(t, sdk.NewDec(7_000_000), yes)
	require.Equal(t, sdk.NewDec(93_000_000), veto)
}
