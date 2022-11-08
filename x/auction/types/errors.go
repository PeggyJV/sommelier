package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/auction module sentinel errors
var (
	ErrCouldNotFindSaleTokenPrice                               = sdkerrors.Register(ModuleName, 2, "could not find sale token price, need to resubmit token prices and try agian")
	ErrCouldNotFindSommTokenPrice                               = sdkerrors.Register(ModuleName, 3, "could not find usomm token price, need to resubmit token prices and try again")
	ErrLastSaleTokenPriceTooOld                                 = sdkerrors.Register(ModuleName, 4, "last sale token price update too long ago, need to resubmit token prices and try again")
	ErrLastSommTokenPriceTooOld                                 = sdkerrors.Register(ModuleName, 5, "last usomm token price update too long ago, need to resubmit token prices and try again")
	ErrAuctionStartingAmountMustBePositve                       = sdkerrors.Register(ModuleName, 6, "minimum auction sale token starting amount must be a positive amount of coins")
	ErrCannotAuctionUsomm                                       = sdkerrors.Register(ModuleName, 7, "auctioning usomm for usomm is pointless")
	ErrInvalidInitialDecreaseRate                               = sdkerrors.Register(ModuleName, 8, "initial price decrease rate must be a float less than one and greater than zero")
	ErrInvalidBlockDecreaseInterval                             = sdkerrors.Register(ModuleName, 9, "price decrease block interval cannot be 0")
	ErrUnauthorizedFundingModule                                = sdkerrors.Register(ModuleName, 10, "unauthorized funding module account")
	ErrUnauthorizedProceedsModule                               = sdkerrors.Register(ModuleName, 11, "unauthorized proceeds module account")
	ErrCannotStartTwoAuctionsForSameDenomSimultaneously         = sdkerrors.Register(ModuleName, 12, "auction for this denom is currently ongoing, cannot create another auction for the same denom until completed")
	ErrConvertingStringToDec                                    = sdkerrors.Register(ModuleName, 13, "could not convert string to dec")
	ErrAuctionNotFound                                          = sdkerrors.Register(ModuleName, 14, "auction not found")
	ErrBidAuctionDenomMismatch                                  = sdkerrors.Register(ModuleName, 15, "auction denom different from bid requested denom")
	ErrAuctionEnded                                             = sdkerrors.Register(ModuleName, 16, "auction ended")
	ErrInsufficientBid                                          = sdkerrors.Register(ModuleName, 17, "max bid amount is too small to purchase minimum sale tokens requested")
	ErrMinimumPurchaseAmountLargerThanTokensRemaining           = sdkerrors.Register(ModuleName, 18, "minimum purchase amount is larger then the number of tokens remaining for sale")
	ErrAuctionIDMustBeNonZero                                   = sdkerrors.Register(ModuleName, 19, "auction IDs must be non-zero")
	ErrInvalidStartBlock                                        = sdkerrors.Register(ModuleName, 20, "start block cannot be 0")
	ErrInvalidCurrentDecreaseRate                               = sdkerrors.Register(ModuleName, 21, "current price decrease rate must be a float less than one and greater than zero")
	ErrPriceMustBePositive                                      = sdkerrors.Register(ModuleName, 22, "price must be positive")
	ErrDenomCannotBeEmpty                                       = sdkerrors.Register(ModuleName, 23, "denom cannot be empty")
	ErrInvalidLastUpdatedBlock                                  = sdkerrors.Register(ModuleName, 24, "last updated block cannot be 0")
	ErrBidIDMustBeNonZero                                       = sdkerrors.Register(ModuleName, 25, "bid ID must be non-zero")
	ErrBidAmountMustBePositive                                  = sdkerrors.Register(ModuleName, 26, "bid amount must be positive")
	ErrBidMustBeInUsomm                                         = sdkerrors.Register(ModuleName, 27, "bid must be in usomm")
	ErrInvalidTokenBeingBidOn                                   = sdkerrors.Register(ModuleName, 28, "tokens being bid on must have the gravity prefix")
	ErrMinimumAmountMustBePositive                              = sdkerrors.Register(ModuleName, 29, "minimum amount to purchase with bid must be positive")
	ErrAddressExpected                                          = sdkerrors.Register(ModuleName, 30, "address cannot be empty")
	ErrBidUnitPriceInUsommMustBePositive                        = sdkerrors.Register(ModuleName, 31, "unit price of sale tokens in usomm must be positive")
	ErrInvalidTokenPriceDenom                                   = sdkerrors.Register(ModuleName, 32, "token price denoms must be either usomm or addresses prefixed with 'gravity'")
	ErrTokenPriceProposalAttemptsToUpdateTokenPriceMoreThanOnce = sdkerrors.Register(ModuleName, 33, "token price proposals should not attempt to update the same denom's price more than once per proposal")
	ErrTokenPriceMaxBlockAgeMustBePositive                      = sdkerrors.Register(ModuleName, 34, "price max block age must be positive")
	ErrInvalidPriceMaxBlockAgeParameterType                     = sdkerrors.Register(ModuleName, 35, "price max block age type must be uint64")
	ErrTokenPriceProposalMustHaveAtLeastOnePrice                = sdkerrors.Register(ModuleName, 36, "list of proposed token prices must be non-zero")
	ErrBidFulfilledSaleTokenAmountMustBeNonNegative             = sdkerrors.Register(ModuleName, 37, "total sale token fulfilled amount must be non-negative")
	ErrBidPaymentCannotBeNegative                               = sdkerrors.Register(ModuleName, 38, "total amount paid in usomm cannot be negative")
	ErrBidAmountIsTooSmall                                      = sdkerrors.Register(ModuleName, 39, "bid is below minimum amount")
	ErrMinimumBidParam                                          = sdkerrors.Register(ModuleName, 40, "invalid minimum bid param")
	ErrInvalidAuctionMaxBlockAgeParam                           = sdkerrors.Register(ModuleName, 41, "invalid auction max block age param")
	ErrInvalidAuctionPriceDecreaseAccelerationRateParam         = sdkerrors.Register(ModuleName, 42, "invalid auction price decrease acceleration rate param")
)
