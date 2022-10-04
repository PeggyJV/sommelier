package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cellarfees module sentinel errors
var (
	ErrInvalidFeeAccrualAuctionThreshold = sdkerrors.Register(ModuleName, 2, "invalid fee accrual auction threshold")
	ErrInvalidRewardEmissionPeriod       = sdkerrors.Register(ModuleName, 3, "invalid reward emission period")
	ErrInvalidInitialPriceDecreaseRate   = sdkerrors.Register(ModuleName, 4, "invalid initial price decrease rate")
	ErrInvalidPriceDecreaseBlockInterval = sdkerrors.Register(ModuleName, 5, "invalid price decrease block interval")
)
