package types

import (
	errorsmod "cosmossdk.io/errors"
)

// x/cellarfees module sentinel errors
var (
	ErrInvalidFeeAccrualAuctionThreshold = errorsmod.Register(ModuleName, 2, "invalid fee accrual auction threshold")
	ErrInvalidFeeAccrualCounters         = errorsmod.Register(ModuleName, 6, "invalid fee accrual counters")
)
