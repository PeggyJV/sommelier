package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/cellarfees module sentinel errors
var (
	ErrInvalidFeeAccrualAuctionThreshold = errorsmod.Register(ModuleName, 2, "invalid fee accrual auction threshold")
	ErrInvalidRewardEmissionPeriod       = errorsmod.Register(ModuleName, 3, "invalid reward emission period")
	ErrInvalidInitialPriceDecreaseRate   = errorsmod.Register(ModuleName, 4, "invalid initial price decrease rate")
	ErrInvalidPriceDecreaseBlockInterval = errorsmod.Register(ModuleName, 5, "invalid price decrease block interval")
	ErrInvalidFeeAccrualCounters         = errorsmod.Register(ModuleName, 6, "invalid fee accrual counters")
	ErrInvalidLastRewardSupplyPeak       = errorsmod.Register(ModuleName, 7, "invalid last reward supply peak")
)
