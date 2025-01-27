package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/cellarfees module sentinel errors
var (
	// Codes 2 and 6 were deleted during v2 module upgrade

	ErrInvalidRewardEmissionPeriod       = errorsmod.Register(ModuleName, 3, "invalid reward emission period")
	ErrInvalidInitialPriceDecreaseRate   = errorsmod.Register(ModuleName, 4, "invalid initial price decrease rate")
	ErrInvalidPriceDecreaseBlockInterval = errorsmod.Register(ModuleName, 5, "invalid price decrease block interval")
	ErrInvalidLastRewardSupplyPeak       = errorsmod.Register(ModuleName, 7, "invalid last reward supply peak")
	ErrInvalidAuctionInterval            = errorsmod.Register(ModuleName, 8, "invalid interval blocks between auctions")
	ErrInvalidAuctionThresholdUsdValue   = errorsmod.Register(ModuleName, 9, "invalid auction threshold USD value")
	ErrInvalidProceedsPortion            = errorsmod.Register(ModuleName, 10, "invalid proceeds portion")
)
