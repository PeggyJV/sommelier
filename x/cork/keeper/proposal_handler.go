package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v3/x/cork/types"
)

// HandleAddManagedCellarsProposal is a handler for executing a passed community cellar addition proposal
func HandleAddManagedCellarsProposal(ctx sdk.Context, k Keeper, p types.AddManagedCellarsProposal) error {
	var IDMap map[common.Address]bool

	for _, existingID := range k.GetCellarIDs(ctx) {
		IDMap[existingID] = true
	}

	for _, cellarID := range p.CellarIds.Ids {
		IDMap[common.HexToAddress(cellarID)] = true
	}

	var outputCellarIDs types.CellarIDSet
	for key := range IDMap {
		outputCellarIDs.Ids = append(outputCellarIDs.Ids, key.Hex())
	}
	k.SetCellarIDs(ctx, outputCellarIDs)

	return nil
}

// HandleRemoveManagedCellarsProposal is a handler for executing a passed community cellar removal proposal
func HandleRemoveManagedCellarsProposal(ctx sdk.Context, k Keeper, p types.RemoveManagedCellarsProposal) error {
	var outputCellarIDs types.CellarIDSet

	for _, existingID := range k.GetCellarIDs(ctx) {
		found := false
		for _, inputID := range p.CellarIds.Ids {
			if existingID == common.HexToAddress(inputID) {
				found = true
			}
		}

		if !found {
			outputCellarIDs.Ids = append(outputCellarIDs.Ids, existingID.Hex())
		}
	}
	k.SetCellarIDs(ctx, outputCellarIDs)

	return nil
}
