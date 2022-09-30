package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/cellarfees module sentinel errors
var (
	ErrInvalidAuctionBlockDelay          = sdkerrors.Register(ModuleName, 2, "invalid auction block delay")
	ErrInvalidRewardEmissionPeriod       = sdkerrors.Register(ModuleName, 3, "invalid reward emission period")
	ErrInvalidInitialPriceDecreaseRate   = sdkerrors.Register(ModuleName, 4, "invalid initial price decrease rate")
	ErrInvalidPriceDecreaseBlockInterval = sdkerrors.Register(ModuleName, 5, "invalid price decrease block interval")
)
