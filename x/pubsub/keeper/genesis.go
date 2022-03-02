package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	k.SetParams(ctx, gs.Params)

	for _, publisher := range gs.Publishers {
		k.SetPublisher(ctx, *publisher)
	}

	for _, subscriber := range gs.Subscribers {
		addr, err := sdk.AccAddressFromBech32(subscriber.Address)
		if err != nil {
			panic(fmt.Errorf("invalid subscriber address in genesis state: %s", err.Error()))
		}
		k.SetSubscriber(ctx, addr, *subscriber)
	}

	for _, publisherIntent := range gs.PublisherIntents {
		k.SetPublisherIntent(ctx, *publisherIntent)
	}

	for _, subscriberIntent := range gs.SubscriberIntents {
		addr, err := sdk.AccAddressFromBech32(subscriberIntent.SubscriberAddress)
		if err != nil {
			panic(fmt.Errorf("invalid subscriber address in genesis state: %s", err.Error()))
		}
		k.SetSubscriberIntent(ctx, addr, *subscriberIntent)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var publishers []*types.Publisher
	k.IteratePublishers(ctx, func(publisher types.Publisher) (stop bool) {
		publishers = append(publishers, &publisher)
		return false
	})

	var subscribers []*types.Subscriber
	k.IterateSubscribers(ctx, func(subscriberAddress sdk.AccAddress, subscriber types.Subscriber) (stop bool) {
		subscribers = append(subscribers, &subscriber)
		return false
	})

	var publisherIntents []*types.PublisherIntent
	k.IteratePublisherIntents(ctx, func(publisherIntent types.PublisherIntent) (stop bool) {
		publisherIntents = append(publisherIntents, &publisherIntent)
		return false
	})

	var subscriberIntents []*types.SubscriberIntent
	k.IterateSubscriberIntents(ctx, func(subscriberAddress sdk.AccAddress, subscriberIntent types.SubscriberIntent) (stop bool) {
		subscriberIntents = append(subscriberIntents, &subscriberIntent)
		return false
	})

	return types.GenesisState{
		Params:            k.GetParams(ctx),
		Publishers:        publishers,
		Subscribers:       subscribers,
		PublisherIntents:  publisherIntents,
		SubscriberIntents: subscriberIntents,
	}
}
