package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestNewMsgSubmitBidRequestFormatting(t *testing.T) {
	expectedMsg := &MsgSubmitBidRequest{
		AuctionId:              uint32(1),
		MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewInt(200)),
		SaleTokenMinimumAmount: sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(1)),
		Bidder:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
		Signer:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
	}

	createdMsg, err := NewMsgSubmitBidRequest(uint32(1), sdk.NewCoin("usomm", sdk.NewInt(200)), sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(1)), sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny"))
	require.Nil(t, err)
	require.Equal(t, expectedMsg, createdMsg)
}

func TestMsgValidate(t *testing.T) {
	testCases := []struct {
		name                string
		msgSubmitBidRequest MsgSubmitBidRequest
		expPass             bool
		err                 error
	}{
		{
			name: "Happy path",
			msgSubmitBidRequest: MsgSubmitBidRequest{
				AuctionId:              uint32(1),
				MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewInt(200)),
				SaleTokenMinimumAmount: sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(1)),
				Bidder:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
				Signer:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Auction ID cannot be 0",
			msgSubmitBidRequest: MsgSubmitBidRequest{
				AuctionId:              uint32(0),
				MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewInt(200)),
				SaleTokenMinimumAmount: sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(1)),
				Bidder:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
				Signer:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrAuctionIDMustBeNonZero, "id: 0"),
		},
		{
			name: "Bid must be in usomm",
			msgSubmitBidRequest: MsgSubmitBidRequest{
				AuctionId:              uint32(1),
				MaxBidInUsomm:          sdk.NewCoin("usdc", sdk.NewInt(200)),
				SaleTokenMinimumAmount: sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(1)),
				Bidder:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
				Signer:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrBidMustBeInUsomm, "bid: %s", sdk.NewCoin("usdc", sdk.NewInt(200))),
		},
		{
			name: "Bid must be positive",
			msgSubmitBidRequest: MsgSubmitBidRequest{
				AuctionId:              uint32(1),
				MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewInt(0)),
				SaleTokenMinimumAmount: sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(1)),
				Bidder:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
				Signer:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrBidAmountMustBePositive, "bid amount in usomm: %s", sdk.NewCoin("usomm", sdk.NewInt(0))),
		},
		{
			name: "Sale token must be prefixed with gravity0x",
			msgSubmitBidRequest: MsgSubmitBidRequest{
				AuctionId:              uint32(1),
				MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewInt(200)),
				SaleTokenMinimumAmount: sdk.NewCoin("usdc", sdk.NewInt(1)),
				Bidder:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
				Signer:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrInvalidTokenBeingBidOn, "sale token: %s", sdk.NewCoin("usdc", sdk.NewInt(1))),
		},
		{
			name: "Sale token minimum amount must be positive",
			msgSubmitBidRequest: MsgSubmitBidRequest{
				AuctionId:              uint32(1),
				MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewInt(200)),
				SaleTokenMinimumAmount: sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(0)),
				Bidder:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
				Signer:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrMinimumAmountMustBePositive, "sale token amount: %s", sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(0))),
		},
		{
			name: "Bidder address must be in bech32 format",
			msgSubmitBidRequest: MsgSubmitBidRequest{
				AuctionId:              uint32(1),
				MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewInt(200)),
				SaleTokenMinimumAmount: sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(1)),
				Bidder:                 "zoidberg",
				Signer:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "decoding bech32 failed: invalid separator index -1"),
		},

		{
			name: "Signer must be the same as bidder",
			msgSubmitBidRequest: MsgSubmitBidRequest{
				AuctionId:              uint32(1),
				MaxBidInUsomm:          sdk.NewCoin("usomm", sdk.NewInt(200)),
				SaleTokenMinimumAmount: sdk.NewCoin("gravity0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", sdk.NewInt(1)),
				Bidder:                 sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String(),
				Signer:                 sdk.AccAddress("cosmos18ld4633yswcyjdklej3att6aw93nhlf7ce4v8u").String(),
			},
			expPass: false,
			err:     sdkerrors.Wrapf(ErrSignerDifferentFromBidder, "signer: %s, bidder: %s", sdk.AccAddress("cosmos18ld4633yswcyjdklej3att6aw93nhlf7ce4v8u").String(), sdk.AccAddress("cosmos16zrkzad482haunrn25ywvwy6fclh3vh7r0hcny").String()),
		},
	}

	for _, tc := range testCases {
		err := tc.msgSubmitBidRequest.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}
