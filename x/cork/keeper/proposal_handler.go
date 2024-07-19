package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	corktypes "github.com/peggyjv/sommelier/v7/x/cork/types"
	types "github.com/peggyjv/sommelier/v7/x/cork/types/v2"
	pubsubtypes "github.com/peggyjv/sommelier/v7/x/pubsub/types"
)

func NewEthereumSubscriptionID(address common.Address) string {
	return fmt.Sprintf("1:%s", address.String())
}

// HandleAddManagedCellarsProposal is a handler for executing a passed community cellar addition proposal
func HandleAddManagedCellarsProposal(ctx sdk.Context, k Keeper, p types.AddManagedCellarIDsProposal) error {
	_, publisherFound := k.pubsubKeeper.GetPublisher(ctx, p.PublisherDomain)
	if !publisherFound {
		return fmt.Errorf("not an approved publisher: %s", p.PublisherDomain)
	}

	cellarAddresses := k.GetCellarIDs(ctx)

	for _, proposedCellarID := range p.CellarIds.Ids {
		proposedCellarAddress := common.HexToAddress(proposedCellarID)
		found := false
		for _, id := range cellarAddresses {
			if id == proposedCellarAddress {
				found = true
			}
		}
		if !found {
			cellarAddresses = append(cellarAddresses, proposedCellarAddress)
			subscriptionID := NewEthereumSubscriptionID(proposedCellarAddress)
			defaultSubscription := pubsubtypes.DefaultSubscription{
				SubscriptionId:  subscriptionID,
				PublisherDomain: p.PublisherDomain,
			}
			k.pubsubKeeper.SetDefaultSubscription(ctx, defaultSubscription)
		}
	}

	idStrings := make([]string, len(cellarAddresses))
	for i, cid := range cellarAddresses {
		idStrings[i] = cid.String()
	}

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
			outputCellarIDs.Ids = append(outputCellarIDs.Ids, existingID.String())
		}
	}
	k.SetCellarIDs(ctx, outputCellarIDs)

	for _, cellarToDelete := range p.CellarIds.Ids {
		subscriptionID := NewEthereumSubscriptionID(common.HexToAddress(cellarToDelete))
		k.pubsubKeeper.DeleteDefaultSubscription(ctx, subscriptionID)
	}

	return nil
}

// HandleScheduledCorkProposal is a handler for executing a passed scheduled cork proposal
func HandleScheduledCorkProposal(ctx sdk.Context, k Keeper, p types.ScheduledCorkProposal) error {
	if !k.HasCellarID(ctx, common.HexToAddress(p.TargetContractAddress)) {
		return errorsmod.Wrapf(corktypes.ErrUnmanagedCellarAddress, "id: %s", p.TargetContractAddress)
	}

	return nil
}
