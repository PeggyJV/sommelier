package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/allocation/types"
)

// HandleUpdateManagedCellarsProposal is a handler for executing a passed community spend proposal
func HandleUpdateManagedCellarsProposal(ctx sdk.Context, k Keeper, p types.ManagedCellarsUpdateProposal) error {
	panic("implement this") // todo: implement cellars proposal
}
