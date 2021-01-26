package oracle

import (
	"github.com/peggyjv/sommelier/x/oracle/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// TODO: tally votes after voting periods are over
}
