package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/auction module sentinel errors
var (
	ErrSignerDifferentFromBidder                        = sdkerrors.Register(ModuleName, 2, "signer is different from bidder")
	ErrCouldNotFindSaleTokenPrice                       = sdkerrors.Register(ModuleName, 3, "could not find sale token price, need to resubmit token prices and try agian")
	ErrCouldNotFindSommTokenPrice                       = sdkerrors.Register(ModuleName, 4, "could not find somm token price, need to resubmit token prices and try again")
	ErrLastSaleTokenPriceUpdateTooLongAgo               = sdkerrors.Register(ModuleName, 5, "last sale token price update too long ago, need to resubmit token prices and try again")
	ErrLastSommTokenPriceUpdateTooLongAgo               = sdkerrors.Register(ModuleName, 6, "last somm token price update too long ago, need to resubmit token prices and try again")
	ErrAuctionStartinAmountMustBePositve                = sdkerrors.Register(ModuleName, 7, "minimum auction amount must be a positive amount of coins")
	ErrAuctionDenomInvalid                              = sdkerrors.Register(ModuleName, 8, "action denom must be non empty")
	ErrCannotAuctionUsomm                               = sdkerrors.Register(ModuleName, 9, "auctioning usomm for usomm is pointless")
	ErrInvalidInitialDecreaseRate                       = sdkerrors.Register(ModuleName, 10, "initial decrease rate must be a float less than one and greater than zero")
	ErrInvalidBlockDecreaeInterval                      = sdkerrors.Register(ModuleName, 11, "block decrease interval cannot be 0")
	ErrUnauthorizedFundingModule                        = sdkerrors.Register(ModuleName, 12, "unauthorized funding module account")
	ErrUnauthorizedProceedsModule                       = sdkerrors.Register(ModuleName, 13, "unauthorized proceeds module account")
	ErrCannotStartTwoAuctionsForSameDenomSimultaneously = sdkerrors.Register(ModuleName, 14, "auction for this denom is currently ongoing, cannot create another auction for the same denom until completed")
	ErrConvertingTokenPriceToFloat                      = sdkerrors.Register(ModuleName, 15, "could not convert token price to float")
	ErrConvertingStringToDec                            = sdkerrors.Register(ModuleName, 16, "could not convert string to dec")
	ErrAuctionNotFound                                  = sdkerrors.Register(ModuleName, 17, "auction not found")
	ErrInvalidAddress                                   = sdkerrors.Register(ModuleName, 18, "invalid address")
	ErrBidAuctionDenomMismatch                          = sdkerrors.Register(ModuleName, 19, "auction denom different from bid requested denom")
	ErrAuctionEnded                                     = sdkerrors.Register(ModuleName, 20, "auction ended")
	ErrMinimumPurchaseLargerThanBid                     = sdkerrors.Register(ModuleName, 21, "minimum purchase is larger than allocated bid amount")
	ErrBidSmallerThanMinimumPurchasePrice               = sdkerrors.Register(ModuleName, 22, "bid smaller than minimum purchase price")
	ErrInsufficientBid                                  = sdkerrors.Register(ModuleName, 23, "bid amount is too small to purchase any tokens on sale, please increase and try again")
	ErrMinimumPurchaseAmountLargerThanTokensRemaining   = sdkerrors.Register(ModuleName, 24, "minimum purchase amount is larger then the number of tokens remaining for sale")
)
