package keeper

import (
	"fmt"
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ValCellar struct {
	Val  sdktypes.ValAddress
	Cork types.Cork
}

type VoteCalculatorTestCase struct {
	title        string
	description  string
	CellarID     common.Address
	ValCellars   []ValCellar
	WinningVotes []types.Cork
}

var (
	vallAddrA, _ = sdktypes.ValAddressFromHex("24ep6yqkhpwnfdrrapu6fzmjp3xrpsgca11ab1e")

	exampleAddrA = common.HexToAddress("0xc0ffee254729296a45a3885639AC7E10F9d54979")
)

func TestGetWinningVotes(t *testing.T) {
	testCases := []VoteCalculatorTestCase{
		{title: "Single voter",
			description: "Check that a single voter returns it's vote",
			CellarID:    exampleAddrA,
			ValCellars: []ValCellar{
				{Val: vallAddrA,
					Cork: types.Cork{
						TargetContractAddress: exampleAddrA.String(),
						EncodedContractCall:   []byte{33},
					},
				},
			},
			WinningVotes: []types.Cork{
				{
					TargetContractAddress: exampleAddrA.String(),
					EncodedContractCall:   []byte{33},
				},
			},
		},
	}

	for _, test := range testCases {
		input := CreateTestEnv(t)
		ctx := input.Context
		t.Logf(test.title)

		for _, vc := range test.ValCellars {
			commit := types.Cork{
				TargetContractAddress: exampleAddrA.String(),
				EncodedContractCall:   []byte{33},
			}

			input.corkKeeper.SetCork(ctx, vc.Val, commit)
		}

		winningVotes := input.corkKeeper.GetApprovedCorks(ctx, sdktypes.MustNewDecFromStr("0.66"))
		require.Lenf(t, winningVotes, 1, test.description)
	}
}

func TestSetGetCellarIDs(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		k, _, ctx := setupCorkKeeper(t)

		cellarID := exampleAddrA
		cellars := k.GetCellarIDs(ctx)
		require.True(t, len(cellars) == 0)

		k.SetCellarIDs(ctx, types.CellarIDSet{
			Ids: []string{cellarID.String()},
		})
		cellars = k.GetCellarIDs(ctx)
		assert.True(t, len(cellars) > 0)
		assert.Contains(t, cellars[0].String(), cellarID.String())
	})
}

