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
	//vallAddrB, _ = sdktypes.ValAddressFromHex("1wr4386xp9u0mtk8u56hdf5zuurga0hb01dface")
	//vallAddrC, _ = sdktypes.ValAddressFromHex("1wr4386xp9u0mtk8u56hdf5zuurga0hdeadbeef")
	//vallAddrD, _ = sdktypes.ValAddressFromHex("1wr4386xp9u0mtk8u56hdf5zuurga0hf005ba11")

	cellarAddrA = common.HexToAddress("0xc0ffee254729296a45a3885639AC7E10F9d54979")
)

func TestGetWinningVotes(t *testing.T) {
	testCases := []VoteCalculatorTestCase{
		{"Single voter",
			"Check that a single voter returns it's vote",
			cellarAddrA,
			[]ValCellar{
				{vallAddrA,
					types.Cellar{cellarAddrA.String(), []*types.TickRange{
						{100, -100, 30},
					},
					},
				},
			},
			[]types.Cellar{
				{
					cellarAddrA.String(),
					[]*types.TickRange{
						{100,
							-100,
							30},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		input := CreateTestEnv(t)
		ctx := input.Context

		for _, vc := range test.ValCellars {
			commit := types.Allocation{
				Cellar: &vc.Cellar,
				Salt:   "testsalt",
			}

			if _, found := input.AllocationKeeper.GetCellarByID(ctx, common.HexToAddress(vc.Cellar.Id)); !found {
				input.AllocationKeeper.SetCellar(ctx, vc.Cellar)
			}

			input.AllocationKeeper.SetAllocationCommit(ctx, vc.Val, cellarAddrA, commit)
		}

		winningVotes := input.AllocationKeeper.GetWinningVotes(ctx, sdktypes.MustNewDecFromStr("0.66"))
		require.Lenf(t, winningVotes, 1, "require that winning votes contains only one cellar")
	}
}

func TestHashingPreCommitsAndCommits(t *testing.T) {
	input := CreateTestEnv(t)
	ctx := input.Context

	testCellar := common.HexToAddress("0x6ea5992aB4A78D5720bD12A089D13c073d04B55d")

	commit := types.Allocation{
		Cellar: &types.Cellar{
			Id: testCellar.String(),
			TickRanges: []*types.TickRange{
				{200, 100, 10},
				{300, 200, 20},
				{400, 300, 30},
				{500, 400, 40},
			},
		},
		Salt: "testsalt",
	}

	testAcc, err := sdktypes.AccAddressFromHex("beefface")
	require.NoError(t, err, "unable to parse acc addr")
	testVal := sdktypes.ValAddress(testAcc)

	preCommitMsg, err := types.NewMsgAllocationPrecommit(*commit.Cellar, commit.Salt, testAcc, testVal)
	require.NoError(t, err, "can't make precommit message")

	// store precommit
	input.AllocationKeeper.SetAllocationPrecommit(ctx, testVal, testCellar, *preCommitMsg.Precommit[0])

	// retrieve precommit
	pc, found := input.AllocationKeeper.GetAllocationPrecommit(ctx, testVal, common.HexToAddress(commit.Cellar.Id))
	require.True(t, found, "didn't find precommit")
	require.Equal(t, preCommitMsg.Precommit[0].Hash, pc.Hash, "bytes unequal after retrieving precommit")

	commitHash, err := commit.Cellar.Hash(commit.Salt, testVal)
	require.NoError(t, err, "couldn't hash commit")
	require.Equal(t, pc.Hash, commitHash, "hashes don't match")
}