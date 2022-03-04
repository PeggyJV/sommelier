package keeper

import (
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v3/x/cork/types"
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
