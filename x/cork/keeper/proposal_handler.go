package keeper

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
)

// HandleAddManagedCellarsProposal is a handler for executing a passed community cellar addition proposal
func HandleAddManagedCellarsProposal(ctx sdk.Context, k Keeper, p types.AddManagedCellarIDsProposal) error {
	cellarIDs := k.GetCellarIDs(ctx)

	for _, proposedCellarID := range p.CellarIds.Ids {
		found := false
		for _, id := range cellarIDs {
			if id == common.HexToAddress(proposedCellarID) {
				found = true
			}
		}
		if !found {
			cellarIDs = append(cellarIDs, common.HexToAddress(proposedCellarID))
		}
	}

	idStrings := make([]string, len(cellarIDs))
	for i, cid := range cellarIDs {
		idStrings[i] = cid.String()
	}

	sort.Strings(idStrings)
	k.SetCellarIDs(ctx, types.CellarIDSet{Ids: idStrings})

	return nil
}

// HandleRemoveManagedCellarsProposal is a handler for executing a passed community cellar removal proposal
func HandleRemoveManagedCellarsProposal(ctx sdk.Context, k Keeper, p types.RemoveManagedCellarIDsProposal) error {
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

// HandleScheduledCorkProposal is a handler for executing a passed scheduled cork proposal
func HandleScheduledCorkProposal(ctx sdk.Context, k Keeper, p types.ScheduledCorkProposal) error {
	cork := types.Cork{
		EncodedContractCall:   []byte(p.ContractCallProtoJson),
		TargetContractAddress: p.TargetContractAddress,
	}
	k.stakingKeeper.IterateLastValidators(ctx, func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
		k.SetScheduledCork(ctx, p.BlockHeight, validator.GetOperator(), cork)
		return false
	})

	return nil
}
