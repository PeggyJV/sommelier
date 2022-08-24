package auction

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v4/x/auction/keeper"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

// NewHandler returns a handler for "auction" type messages.
func NewHandler(k keeper.Keeper) sdk.Handler {
	// TODO: fill in
	return nil
}

func NewSetTokenPricesProposalHandler(k keeper.Keeper) govtypes.Handler {
	// TODO: fill in
	return nil
}

// TODO: Fill our remaining handler functions