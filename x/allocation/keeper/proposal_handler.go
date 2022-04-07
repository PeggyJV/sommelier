package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v4/x/allocation/types"
)

// HandleAddManagedCellarsProposal is a handler for executing a passed community cellar addition proposal
func HandleAddManagedCellarsProposal(ctx sdk.Context, k Keeper, p types.AddManagedCellarsProposal) error {
	for _, cellarID := range p.CellarIds {
		cellar := types.Cellar{
			Id:         cellarID,
			TickRanges: nil,
		}
		if _, ok := k.GetCellarByID(ctx, cellar.Address()); ok {
			return fmt.Errorf("cellar with id %v already exists for proposal %v", cellarID, p)
		}

		k.SetCellar(ctx, cellar)
	}

	return nil
}

// HandleRemoveManagedCellarsProposal is a handler for executing a passed community cellar removal proposal
func HandleRemoveManagedCellarsProposal(ctx sdk.Context, k Keeper, p types.RemoveManagedCellarsProposal) error {

	// todo: should we do any checks here? Is there any circumstance where a cellar can't/shouldn't be removed?
	for _, cellarID := range p.CellarIds {
		cellarAddr := common.HexToAddress(cellarID)
		if _, ok := k.GetCellarByID(ctx, cellarAddr); !ok {
			return fmt.Errorf("cellar with id %v not found for proposal %v", cellarAddr, p)
		}
		k.DeleteCellar(ctx, cellarAddr)
	}

	return nil
}
