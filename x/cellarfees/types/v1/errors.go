package v1

import (
	errorsmod "cosmossdk.io/errors"
)

// x/cellarfees module sentinel errors
var (
	ErrInvalidFeeAccrualAuctionThreshold = errorsmod.Register("cellarfees", 2, "invalid fee accrual auction threshold")
	ErrInvalidRewardEmissionPeriod       = errorsmod.Register("cellarfees", 3, "invalid reward emission period")
	ErrInvalidInitialPriceDecreaseRate   = errorsmod.Register("cellarfees", 4, "invalid initial price decrease rate")
	ErrInvalidPriceDecreaseBlockInterval = errorsmod.Register("cellarfees", 5, "invalid price decrease block interval")
	ErrInvalidFeeAccrualCounters         = errorsmod.Register("cellarfees", 6, "invalid fee accrual counters")
	ErrInvalidLastRewardSupplyPeak       = errorsmod.Register("cellarfees", 7, "invalid last reward supply peak")
	ErrInvalidAuctionInterval            = errorsmod.Register("cellarfees", 8, "invalid interval blocks between auctions")
)
