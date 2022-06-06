package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ValCellar struct {
	Val  sdk.ValAddress
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
	vallAddrA, _ = sdk.ValAddressFromHex("24ep6yqkhpwnfdrrapu6fzmjp3xrpsgca11ab1e")

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

		winningVotes := input.corkKeeper.GetApprovedCorks(ctx, sdk.MustNewDecFromStr("0.66"))
		require.Lenf(t, winningVotes, 1, test.description)
	}
}

func TestSetGetCellarIDs(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		k, ctx, _, _ := setupCorkKeeper(t)

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

func TestSetCorkGetCork_Unit(t *testing.T) {
	testCases := []struct {
		name string
		test func()
	}{
		{
			name: "todo",
			test: func() {
				t.Log("Declare test case parameters")
				cellarID := exampleAddrA
				valCellar := ValCellar{
					Val: vallAddrA,
					Cork: types.Cork{
						TargetContractAddress: cellarID.String(),
						EncodedContractCall:   []byte{33},
					},
				}

				k, ctx, _, _ := setupCorkKeeper(t)

				t.Log("Set corks")
				vc := valCellar
				commit := types.Cork{
					TargetContractAddress: exampleAddrA.String(),
					EncodedContractCall:   []byte{33},
				}

				k.SetCork(
					ctx,
					/* val */ vc.Val,
					/* cork */ commit,
				)

				// TODO: test getter after k.SetCork
				// contract :=
				// k.GetCork(ctx, vc.Val, contract)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.test()
		})
	}
}

func TestGetWinningVotes_Unit(t *testing.T) {
	testCases := []struct {
		name string
		test func()
	}{
		{
			name: "single voter",
			test: func() {
				t.Log("Declare test case parameters")
				testParams := VoteCalculatorTestCase{
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
				}

				fmt.Println(testParams)
				k, ctx, mocks, _ := setupCorkKeeper(t)
				fmt.Println(mocks)

				for _, vc := range testParams.ValCellars {
					commit := types.Cork{
						TargetContractAddress: exampleAddrA.String(),
						EncodedContractCall:   []byte{33},
					}

					k.SetCork(ctx, vc.Val, commit)
				}

				totalPower := sdk.NewInt(100)
				mocks.mockStakingKeeper.
					EXPECT().GetLastTotalPower(ctx).
					Return(totalPower)

				mockValidator := mocks.mockValidator
				mocks.mockStakingKeeper.
					EXPECT().Validator(
					ctx,
					/* val */ testParams.ValCellars[0].Val).
					Return(mockValidator)
				winningVotes := k.GetApprovedCorks(
					ctx, sdk.MustNewDecFromStr("0.66"))
				require.Lenf(t, winningVotes, 1, testParams.description)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tc.test()
		})
	}
}
