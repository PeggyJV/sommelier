package types

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/stretchr/testify/require"
)

func TestNewSetTokenPricesProposal(t *testing.T) {
	expectedMsg := &SetTokenPricesProposal{
		Title:       "Planet Express",
		Description: "Why not zoidberg",
		TokenPrices: []*ProposedTokenPrice{
			{
				Denom:    "usomm",
				UsdPrice: sdk.MustNewDecFromStr("4.2"),
			},
		},
	}

	createdMsg := NewSetTokenPricesProposal("Planet Express", "Why not zoidberg", []*ProposedTokenPrice{
		{
			Denom:    "usomm",
			UsdPrice: sdk.MustNewDecFromStr("4.2"),
		},
	})
	require.Equal(t, expectedMsg, createdMsg)
}

func TestTokenPriceProposalValidate(t *testing.T) {
	testCases := []struct {
		name                   string
		setTokenPricesProposal SetTokenPricesProposal
		expPass                bool
		err                    error
	}{
		{
			name: "Happy path",
			setTokenPricesProposal: SetTokenPricesProposal{
				Title:       "Planet Express",
				Description: "Why not zoidberg",
				TokenPrices: []*ProposedTokenPrice{
					{
						Denom:    "usomm",
						UsdPrice: sdk.MustNewDecFromStr("4.2"),
					},
				},
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Gov validate basic canary 1 -- title cannot be empty",
			setTokenPricesProposal: SetTokenPricesProposal{
				Title:       "",
				Description: "Why not zoidberg",
				TokenPrices: []*ProposedTokenPrice{
					{
						Denom:    "usomm",
						UsdPrice: sdk.MustNewDecFromStr("4.2"),
					},
				},
			},
			expPass: false,
			err:     errorsmod.Wrap(govtypes.ErrInvalidProposalContent, "proposal title cannot be blank"),
		},
		{
			name: "Gov validate basic canary 2 -- description cannot be empty",
			setTokenPricesProposal: SetTokenPricesProposal{
				Title:       "Planet Express",
				Description: "",
				TokenPrices: []*ProposedTokenPrice{
					{
						Denom:    "usomm",
						UsdPrice: sdk.MustNewDecFromStr("4.2"),
					},
				},
			},
			expPass: false,
			err:     errorsmod.Wrap(govtypes.ErrInvalidProposalContent, "proposal description cannot be blank"),
		},
		{
			name: "Token price proposal must have at least one token price",
			setTokenPricesProposal: SetTokenPricesProposal{
				Title:       "Planet Express",
				Description: "Why not zoidberg",
				TokenPrices: []*ProposedTokenPrice{},
			},
			expPass: false,
			err:     errorsmod.Wrapf(ErrTokenPriceProposalMustHaveAtLeastOnePrice, "prices: %v", []*ProposedTokenPrice{}),
		},
		{
			name: "Cannot have duplicate denoms in token price proposal",
			setTokenPricesProposal: SetTokenPricesProposal{
				Title:       "Planet Express",
				Description: "Why not zoidberg",
				TokenPrices: []*ProposedTokenPrice{
					{
						Denom:    "usomm",
						UsdPrice: sdk.MustNewDecFromStr("4.2"),
					},
					{
						Denom:    "usomm",
						UsdPrice: sdk.MustNewDecFromStr("7.8"),
					},
				},
			},
			expPass: false,
			err:     errorsmod.Wrapf(ErrTokenPriceProposalAttemptsToUpdateTokenPriceMoreThanOnce, "denom: usomm"),
		},
		{
			name: "Token price validate basic canary 1 -- cannot have empty denom",
			setTokenPricesProposal: SetTokenPricesProposal{
				Title:       "Planet Express",
				Description: "Why not zoidberg",
				TokenPrices: []*ProposedTokenPrice{
					{
						Denom:    "",
						UsdPrice: sdk.MustNewDecFromStr("4.2"),
					},
				},
			},
			expPass: false,
			err:     errorsmod.Wrapf(ErrDenomCannotBeEmpty, "price denom: "),
		},
		{
			name: "Token price validate basic canary 2 -- price must be positive",
			setTokenPricesProposal: SetTokenPricesProposal{
				Title:       "Planet Express",
				Description: "Why not zoidberg",
				TokenPrices: []*ProposedTokenPrice{
					{
						Denom:    "usomm",
						UsdPrice: sdk.MustNewDecFromStr("0.0"),
					},
				},
			},
			expPass: false,
			err:     errorsmod.Wrapf(ErrPriceMustBePositive, "usd price: %s", sdk.MustNewDecFromStr("0.0").String()),
		},
	}

	for _, tc := range testCases {
		err := tc.setTokenPricesProposal.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}
