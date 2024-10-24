package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/auction module sentinel errors
var (
	ErrCouldNotFindSaleTokenPrice                               = errorsmod.Register(ModuleName, 2, "could not find sale token price, need to resubmit token prices and try again")
	ErrCouldNotFindSommTokenPrice                               = errorsmod.Register(ModuleName, 3, "could not find usomm token price, need to resubmit token prices and try again")
	ErrLastSaleTokenPriceTooOld                                 = errorsmod.Register(ModuleName, 4, "last sale token price update too long ago, need to resubmit token prices and try again")
	ErrLastSommTokenPriceTooOld                                 = errorsmod.Register(ModuleName, 5, "last usomm token price update too long ago, need to resubmit token prices and try again")
	ErrAuctionStartingAmountMustBePositve                       = errorsmod.Register(ModuleName, 6, "minimum auction sale token starting amount must be a positive amount of coins")
	ErrCannotAuctionUsomm                                       = errorsmod.Register(ModuleName, 7, "auctioning usomm for usomm is pointless")
	ErrInvalidInitialDecreaseRate                               = errorsmod.Register(ModuleName, 8, "initial price decrease rate must be a float less than one and greater than zero")
	ErrInvalidBlockDecreaseInterval                             = errorsmod.Register(ModuleName, 9, "price decrease block interval cannot be 0")
	ErrUnauthorizedFundingModule                                = errorsmod.Register(ModuleName, 10, "unauthorized funding module account")
	ErrUnauthorizedProceedsModule                               = errorsmod.Register(ModuleName, 11, "unauthorized proceeds module account")
	ErrCannotStartTwoAuctionsForSameDenomSimultaneously         = errorsmod.Register(ModuleName, 12, "auction for this denom is currently ongoing, cannot create another auction for the same denom until completed")
	ErrConvertingStringToDec                                    = errorsmod.Register(ModuleName, 13, "could not convert string to dec")
	ErrAuctionNotFound                                          = errorsmod.Register(ModuleName, 14, "auction not found")
	ErrBidAuctionDenomMismatch                                  = errorsmod.Register(ModuleName, 15, "auction denom different from bid requested denom")
	ErrAuctionEnded                                             = errorsmod.Register(ModuleName, 16, "auction ended")
	ErrInsufficientBid                                          = errorsmod.Register(ModuleName, 17, "max bid amount is too small to purchase minimum sale tokens requested")
	ErrMinimumPurchaseAmountLargerThanTokensRemaining           = errorsmod.Register(ModuleName, 18, "minimum purchase amount is larger then the number of tokens remaining for sale")
	ErrAuctionIDMustBeNonZero                                   = errorsmod.Register(ModuleName, 19, "auction IDs must be non-zero")
	ErrInvalidStartBlock                                        = errorsmod.Register(ModuleName, 20, "start block cannot be 0")
	ErrInvalidCurrentDecreaseRate                               = errorsmod.Register(ModuleName, 21, "current price decrease rate must be a float less than one and greater than zero")
	ErrPriceMustBePositive                                      = errorsmod.Register(ModuleName, 22, "price must be positive")
	ErrDenomCannotBeEmpty                                       = errorsmod.Register(ModuleName, 23, "denom cannot be empty")
	ErrInvalidLastUpdatedBlock                                  = errorsmod.Register(ModuleName, 24, "last updated block cannot be 0")
	ErrBidIDMustBeNonZero                                       = errorsmod.Register(ModuleName, 25, "bid ID must be non-zero")
	ErrBidAmountMustBePositive                                  = errorsmod.Register(ModuleName, 26, "bid amount must be positive")
	ErrBidMustBeInUsomm                                         = errorsmod.Register(ModuleName, 27, "bid must be in usomm")
	ErrInvalidTokenBeingBidOn                                   = errorsmod.Register(ModuleName, 28, "tokens being bid on must have the gravity prefix")
	ErrMinimumAmountMustBePositive                              = errorsmod.Register(ModuleName, 29, "minimum amount to purchase with bid must be positive")
	ErrAddressExpected                                          = errorsmod.Register(ModuleName, 30, "address cannot be empty")
	ErrBidUnitPriceInUsommMustBePositive                        = errorsmod.Register(ModuleName, 31, "unit price of sale tokens in usomm must be positive")
	ErrInvalidTokenPriceDenom                                   = errorsmod.Register(ModuleName, 32, "token price denoms must be either usomm or addresses prefixed with 'gravity'")
	ErrTokenPriceProposalAttemptsToUpdateTokenPriceMoreThanOnce = errorsmod.Register(ModuleName, 33, "token price proposals should not attempt to update the same denom's price more than once per proposal")
	ErrTokenPriceMaxBlockAgeMustBePositive                      = errorsmod.Register(ModuleName, 34, "price max block age must be positive")
	ErrInvalidPriceMaxBlockAgeParameterType                     = errorsmod.Register(ModuleName, 35, "price max block age type must be uint64")
	ErrTokenPriceProposalMustHaveAtLeastOnePrice                = errorsmod.Register(ModuleName, 36, "list of proposed token prices must be non-zero")
	ErrBidFulfilledSaleTokenAmountMustBeNonNegative             = errorsmod.Register(ModuleName, 37, "total sale token fulfilled amount must be non-negative")
	ErrBidPaymentCannotBeNegative                               = errorsmod.Register(ModuleName, 38, "total amount paid in usomm cannot be negative")
	ErrBidAmountIsTooSmall                                      = errorsmod.Register(ModuleName, 39, "bid is below minimum amount")
	ErrMinimumBidParam                                          = errorsmod.Register(ModuleName, 40, "invalid minimum bid param")
	ErrInvalidAuctionMaxBlockAgeParam                           = errorsmod.Register(ModuleName, 41, "invalid auction max block age param")
	ErrInvalidAuctionPriceDecreaseAccelerationRateParam         = errorsmod.Register(ModuleName, 42, "invalid auction price decrease acceleration rate param")
	ErrTokenPriceExponentTooHigh                                = errorsmod.Register(ModuleName, 43, "token price exponent too high, maximum precision of 18")
	ErrInvalidMinimumSaleTokensUSDValue                         = errorsmod.Register(ModuleName, 44, "invalid minimum sale tokens USD value")
	ErrAuctionBelowMinimumUSDValue                              = errorsmod.Register(ModuleName, 45, "auction USD value below minimum")
	ErrInvalidMinimumAuctionHeightParam                         = errorsmod.Register(ModuleName, 46, "invalid minimum auction height param")
	ErrAuctionBelowMinimumHeight                                = errorsmod.Register(ModuleName, 47, "auction block height below minimum")
)
