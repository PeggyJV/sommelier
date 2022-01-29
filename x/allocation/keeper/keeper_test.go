package keeper

import (
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/x/allocation/types"
	"github.com/stretchr/testify/require"
)

type ValCellar struct {
	Val    sdktypes.ValAddress
	Cellar types.Cellar
}

type VoteCalculatorTestCase struct {
	title        string
	description  string
	CellarID     common.Address
	ValCellars   []ValCellar
	WinningVotes []types.Cellar
}

var (
	vallAddrA, _ = sdktypes.ValAddressFromHex("24ep6yqkhpwnfdrrapu6fzmjp3xrpsgca11ab1e")

	cellarAddrA = common.HexToAddress("0xc0ffee254729296a45a3885639AC7E10F9d54979")
)

func TestGetWinningVotes(t *testing.T) {
	testCases := []VoteCalculatorTestCase{
		{title: "Single voter",
			description: "Check that a single voter returns it's vote",
			CellarID:    cellarAddrA,
			ValCellars: []ValCellar{
				{Val: vallAddrA,
					Cellar: types.Cellar{Id: cellarAddrA.String(), TickRanges: []*types.TickRange{
						{Upper: 100, Lower: -100, Weight: 30},
					},
					},
				},
			},
			WinningVotes: []types.Cellar{
				{
					Id: cellarAddrA.String(),
					TickRanges: []*types.TickRange{
						{Upper: 100, Lower: -100, Weight: 30},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		input := CreateTestEnv(t)
		ctx := input.Context
		t.Logf(test.title)

		for _, vc := range test.ValCellars {
			commit := types.Allocation{
				Vote: &types.RebalanceVote{
					Cellar:       &vc.Cellar,
					CurrentPrice: 100,
				},
				Salt: "testsalt",
			}

			if _, found := input.AllocationKeeper.GetCellarByID(ctx, common.HexToAddress(vc.Cellar.Id)); !found {
				input.AllocationKeeper.SetCellar(ctx, vc.Cellar)
			}

			input.AllocationKeeper.SetAllocationCommit(ctx, vc.Val, cellarAddrA, commit)
		}

		winningVotes := input.AllocationKeeper.GetWinningVotes(ctx, sdktypes.MustNewDecFromStr("0.66"))
		require.Lenf(t, winningVotes, 1, test.description)
	}
}

func TestHashingPreCommitsAndCommits(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context

	testCellar := common.HexToAddress("0x6ea5992aB4A78D5720bD12A089D13c073d04B55d")

	commit := types.Allocation{
		Vote: &types.RebalanceVote{
			Cellar: &types.Cellar{
				Id: testCellar.String(),
				TickRanges: []*types.TickRange{
					{Upper: 200, Lower: 100, Weight: 10},
					{Upper: 300, Lower: 200, Weight: 20},
					{Upper: 400, Lower: 300, Weight: 30},
					{Upper: 500, Lower: 400, Weight: 40},
				},
			},
			CurrentPrice: 100,
		},
		Salt: "testsalt",
	}

	testAcc, err := sdktypes.AccAddressFromHex("beefface")
	require.NoError(t, err, "unable to parse acc addr")
	testVal := sdktypes.ValAddress(testAcc)

	preCommitMsg, err := types.NewMsgAllocationPrecommit(*commit.Vote, commit.Salt, testAcc, testVal)
	require.NoError(t, err, "can't make precommit message")

	// store precommit
	input.AllocationKeeper.SetAllocationPrecommit(ctx, testVal, testCellar, *preCommitMsg.Precommit[0])

	// retrieve precommit
	pc, found := input.AllocationKeeper.GetAllocationPrecommit(ctx, testVal, common.HexToAddress(commit.Vote.Cellar.Id))
	require.True(t, found, "didn't find precommit")
	require.Equal(t, preCommitMsg.Precommit[0].Hash, pc.Hash, "bytes unequal after retrieving precommit")

	commitHash, err := commit.Vote.Cellar.Hash(commit.Salt, testVal)
	require.NoError(t, err, "couldn't hash commit")
	require.Equal(t, pc.Hash, commitHash, "hashes don't match")
}
