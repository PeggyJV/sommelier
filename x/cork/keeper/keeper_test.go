package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A validator-cork tuple for the cellar at address, Cork.TargetContractAddress
type ValCellar struct {
	Val  sdk.ValAddress
	Cork types.Cork
}

var (
	sampleValHex     = "24ep6yqkhpwnfdrrapu6fzmjp3xrpsgca11ab1e"
	sampleValAddr, _ = sdk.ValAddressFromHex(sampleValHex)

	sampleCellarHex  = "0xc0ffee254729296a45a3885639AC7E10F9d54979"
	sampleCellarAddr = common.HexToAddress(sampleCellarHex)
)

func TestCellarIDs_SetGetHas(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		k, ctx, _, _ := setupCorkKeeper(t)

		cellarID := sampleCellarAddr
		cellarAddrs := k.GetCellarIDs(ctx)
		assert.Len(t, cellarAddrs, 0)
		assert.Equal(t, false, k.HasCellarID(ctx, cellarID))

		k.SetCellarIDs(ctx, types.CellarIDSet{
			Ids: []string{cellarID.String()},
		})
		cellarAddrs = k.GetCellarIDs(ctx)
		assert.Len(t, cellarAddrs, 1)
		assert.Contains(t, cellarAddrs[0].String(), cellarID.String())
		assert.Equal(t, true, k.HasCellarID(ctx, cellarID))
	})
}

func TestSetCorkGetCork_Unit(t *testing.T) {
	testCases := []struct {
		name string
		test func()
	}{
		{
			name: "set cork and get cork - happy",
			test: func() {
				cellarID := sampleCellarAddr
				valCork := types.ValidatorCork{
					Validator: sampleValAddr.String(),
					Cork: &types.Cork{
						TargetContractAddress: cellarID.String(),
						EncodedContractCall:   []byte{33},
					},
				}

				k, ctx, _, _ := setupCorkKeeper(t)

				valAddr, err := sdk.ValAddressFromBech32(valCork.Validator)
				require.NoError(t, err)

				t.Log("Verify no cork exists for the validator")
				cork, found := k.GetCork(ctx, valAddr, sampleCellarAddr)
				assert.Equal(t, false, found)
				assert.Error(t, cork.ValidateBasic())
				assert.Contains(t, cork.ValidateBasic().Error(),
					"cork has an empty contract call body")

				t.Log("Set corks")
				commit := types.Cork{
					TargetContractAddress: sampleCellarAddr.String(),
					EncodedContractCall:   []byte{33},
				}
				require.NoError(t, err)
				k.SetCork(
					ctx,
					/* val */ valAddr,
					/* cork */ commit,
				)

				t.Log("test getter after k.SetCork")
				cork, found = k.GetCork(ctx, valAddr, sampleCellarAddr)
				assert.NoError(t, cork.ValidateBasic())
				assert.Equal(t, true, found)
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
	type TestCase struct {
		name         string
		description  string
		CellarID     common.Address
		ValCorkPairs []types.ValidatorCork
		WinningVotes []types.Cork
	}

	tc := TestCase{
		name:        "single vote - happy",
		description: "Check that a single voter returns it's vote",
		CellarID:    sampleCellarAddr,
		ValCorkPairs: []types.ValidatorCork{
			{
				Validator: sampleValAddr.String(),
				Cork: &types.Cork{
					TargetContractAddress: sampleCellarAddr.String(),
					EncodedContractCall:   []byte{33},
				},
			},
		},
		WinningVotes: []types.Cork{
			{
				TargetContractAddress: sampleCellarAddr.String(),
				EncodedContractCall:   []byte{33},
			},
		},
	}
	t.Run(tc.name, func(t *testing.T) {
		k, ctx, mocks, _ := setupCorkKeeper(t)

		for _, vc := range tc.ValCorkPairs {
			valAddr, err := sdk.ValAddressFromBech32(vc.Validator)
			require.NoError(t, err)
			commit := types.Cork{
				TargetContractAddress: sampleCellarAddr.String(),
				EncodedContractCall:   []byte{33},
			}
			k.SetCork(ctx, valAddr, commit)
		}

		totalPower := sdk.NewInt(100)
		mocks.mockStakingKeeper.
			EXPECT().GetLastTotalPower(ctx).
			Return(totalPower)

		valAddr, err := sdk.ValAddressFromHex("24C0FFEE254729296A45A3885639AC7E10F9D549")
		assert.NoError(t, err)
		mocks.mockStakingKeeper.
			EXPECT().Validator(ctx, valAddr).Return(mocks.mockValidator)
		mocks.mockStakingKeeper.
			EXPECT().PowerReduction(ctx).Return(totalPower)
		mocks.mockValidator.
			EXPECT().GetConsensusPower(totalPower).Return(int64(100))
		winningVotes := k.GetApprovedCorks(ctx,
			/*threshold=*/ sdk.MustNewDecFromStr("0.66"))
		assert.Lenf(t, winningVotes, 1, tc.description)

		encodedPayloadForContract := tc.ValCorkPairs[0].Cork.EncodedContractCall
		assert.EqualValues(t,
			encodedPayloadForContract,
			winningVotes[0].EncodedContractCall)
	})

}
