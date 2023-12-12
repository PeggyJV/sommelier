package keeper

import (
	"fmt"
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/peggyjv/sommelier/v7/app/params"
	auctionTypes "github.com/peggyjv/sommelier/v7/x/auction/types"
)

// Happy path test for proposal handler
func (suite *KeeperTestSuite) TestHappPathForProposalHandler() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	tokenPrices := make([]*auctionTypes.ProposedTokenPrice, 3)
	tokenPrices[0] = &auctionTypes.ProposedTokenPrice{
		Denom:    "gravity0x3506424f91fd33084466f402d5d97f05f8e3b4af",
		UsdPrice: sdk.MustNewDecFromStr("2.5"),
	}
	tokenPrices[1] = &auctionTypes.ProposedTokenPrice{
		Denom:    "gravity0x5a98fcbea516cf06857215779fd812ca3bef1b32",
		UsdPrice: sdk.MustNewDecFromStr("1.7"),
	}
	tokenPrices[2] = &auctionTypes.ProposedTokenPrice{
		Denom:    params.BaseCoinUnit,
		UsdPrice: sdk.MustNewDecFromStr("1000.0"),
	}

	proposal := auctionTypes.SetTokenPricesProposal{
		Title:       "Super cool and exciting token update proposal",
		Description: "NYC style pizza >>> Chicago style pizza",
		TokenPrices: tokenPrices,
	}

	err := HandleSetTokenPricesProposal(ctx, auctionKeeper, proposal)
	require.Nil(err)

	// Verify token prices set
	foundTokenPrices := auctionKeeper.GetTokenPrices(ctx)
	require.Len(foundTokenPrices, 3)

	require.Equal(tokenPrices[0].Denom, foundTokenPrices[0].Denom)
	require.Equal(tokenPrices[0].UsdPrice, foundTokenPrices[0].UsdPrice)

	require.Equal(tokenPrices[1].Denom, foundTokenPrices[1].Denom)
	require.Equal(tokenPrices[1].UsdPrice, foundTokenPrices[1].UsdPrice)

	require.Equal(tokenPrices[2].Denom, foundTokenPrices[2].Denom)
	require.Equal(tokenPrices[2].UsdPrice, foundTokenPrices[2].UsdPrice)

}

// Unhappy path test for proposal handler
func (suite *KeeperTestSuite) TestUnhappPathForProposalHandler() {
	ctx, auctionKeeper := suite.ctx, suite.auctionKeeper
	require := suite.Require()

	tests := []struct {
		name          string
		proposal      auctionTypes.SetTokenPricesProposal
		expectedError error
	}{
		{
			name: "Validate basic canary 1 -- Govtypes validate abstract canary -- title length cannot be 0",
			proposal: auctionTypes.SetTokenPricesProposal{
				Title:       "",
				Description: "Description",
				TokenPrices: []*auctionTypes.ProposedTokenPrice{},
			},
			expectedError: errorsmod.Wrap(govTypes.ErrInvalidProposalContent, "proposal title cannot be blank"),
		},
		{
			name: "Validate basic canary 2 -- cannot have non usomm & non gravity denom",
			proposal: auctionTypes.SetTokenPricesProposal{
				Title:       "Title",
				Description: "Description",
				TokenPrices: []*auctionTypes.ProposedTokenPrice{
					{
						Denom:    "weth",
						UsdPrice: sdk.MustNewDecFromStr("17.0"),
					},
				},
			},
			expectedError: errorsmod.Wrapf(auctionTypes.ErrInvalidTokenPriceDenom, "denom: weth"),
		},
		{
			name: "Cannot attempt to update prices twice for a denom in one proposal",
			proposal: auctionTypes.SetTokenPricesProposal{
				Title:       "Title",
				Description: "Description",
				TokenPrices: []*auctionTypes.ProposedTokenPrice{
					{
						Denom:    "gravity0x761d38e5ddf6ccf6cf7c55759d5210750b5d60f3",
						UsdPrice: sdk.MustNewDecFromStr("0.01"),
					},
					{
						Denom:    "gravity0x761d38e5ddf6ccf6cf7c55759d5210750b5d60f3",
						UsdPrice: sdk.MustNewDecFromStr("0.02"),
					},
				},
			},
			expectedError: errorsmod.Wrapf(auctionTypes.ErrTokenPriceProposalAttemptsToUpdateTokenPriceMoreThanOnce, "denom: gravity0x761d38e5ddf6ccf6cf7c55759d5210750b5d60f3"),
		},
	}

	for _, tc := range tests {
		tc := tc // Redefine variable here due to passing it to function literal below (scopelint)
		suite.T().Run(fmt.Sprint(tc.name), func(t *testing.T) {
			err := HandleSetTokenPricesProposal(ctx, auctionKeeper, tc.proposal)
			require.Equal(tc.expectedError.Error(), err.Error())
		})
	}
}
