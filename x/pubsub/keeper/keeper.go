package keeper

import (
	"bytes"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v3/x/pubsub/types"
)

type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      sdk.StoreKey
	paramSpace    paramtypes.Subspace
	stakingKeeper types.StakingKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	stakingKeeper types.StakingKeeper,

) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramSpace:    paramSpace,
		stakingKeeper: stakingKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// TODO(bolten): add keeper support functions here

///////////////
// Publisher //
///////////////

func (k Keeper) SetPublisher(ctx sdk.Context, publisher types.Publisher) {
	bz := k.cdc.MustMarshal(&publisher)
	ctx.KVStore(k.storeKey).Set(types.GetPublisherKey(publisher.Domain), bz)
}

func (k Keeper) GetPublisher(ctx sdk.Context, publisherDomain string) (publisher types.Publisher, found bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetPublisherKey(publisherDomain))
	if len(bz) == 0 {
		return types.Publisher{}, false
	}

	k.cdc.MustUnmarshal(bz, &publisher)
	return publisher, true
}

func (k Keeper) IteratePublishers(ctx sdk.Context, handler func(publisher types.Publisher) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.PublisherKeyPrefix})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var publisher types.Publisher
		k.cdc.MustUnmarshal(iter.Value(), &publisher)

		if handler(publisher) {
			break
		}
	}
}

func (k Keeper) GetPublishers(ctx sdk.Context) (publishers []types.Publisher) {
	k.IteratePublishers(ctx, func(publisher types.Publisher) (stop bool) {
		publishers = append(publishers, publisher)
		return false
	})

	return
}

////////////////
// Subscriber //
////////////////

func (k Keeper) SetSubscriber(ctx sdk.Context, subscriberAddress sdk.AccAddress, subscriber types.Subscriber) {
	bz := k.cdc.MustMarshal(&subscriber)
	ctx.KVStore(k.storeKey).Set(types.GetSubscriberKey(subscriberAddress), bz)
}

func (k Keeper) GetSubscriber(ctx sdk.Context, subscriberAddress sdk.AccAddress) (subscriber types.Subscriber, found bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetSubscriberKey(subscriberAddress))
	if len(bz) == 0 {
		return types.Subscriber{}, false
	}

	k.cdc.MustUnmarshal(bz, &subscriber)
	return subscriber, true
}

func (k Keeper) IterateSubscribers(ctx sdk.Context, handler func(subscriberAddress sdk.Address, subscriber types.Subscriber) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.SubscriberKeyPrefix})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var subscriber types.Subscriber
		k.cdc.MustUnmarshal(iter.Value(), &subscriber)

		addressKey := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.SubscriberKeyPrefix}))
		subscriberAddress := sdk.AccAddress(addressKey.Bytes())

		if handler(subscriberAddress, subscriber) {
			break
		}
	}
}

func (k Keeper) GetSubscribers(ctx sdk.Context) (subscribers []types.Subscriber) {
	k.IterateSubscribers(ctx, func(subscriberAddress sdk.Address, subscriber types.Subscriber) (stop bool) {
		subscribers = append(subscribers, subscriber)
		return false
	})

	return
}

/////////////////////
// PublisherIntent //
/////////////////////

func (k Keeper) SetPublisherIntent(ctx sdk.Context, publisherIntent types.PublisherIntent) {
	bz := k.cdc.MustMarshal(&publisherIntent)
	ctx.KVStore(k.storeKey).Set(types.GetPublisherIntentKey(publisherIntent.PublisherDomain, publisherIntent.SubscriptionId), bz)
}

func (k Keeper) GetPublisherIntent(ctx sdk.Context, publisherDomain string, subscriptionId string) (publisherIntent types.PublisherIntent, found bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetPublisherIntentKey(publisherDomain, subscriptionId))
	if len(bz) == 0 {
		return types.PublisherIntent{}, false
	}

	k.cdc.MustUnmarshal(bz, &publisherIntent)
	return publisherIntent, true
}

func (k Keeper) IteratePublisherIntents(ctx sdk.Context, handler func(publisherIntent types.PublisherIntent) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.PublisherIntentKeyPrefix})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var publisherIntent types.PublisherIntent
		k.cdc.MustUnmarshal(iter.Value(), &publisherIntent)

		if handler(publisherIntent) {
			break
		}
	}
}

func (k Keeper) GetPublisherIntents(ctx sdk.Context) (publisherIntents []types.PublisherIntent) {
	k.IteratePublisherIntents(ctx, func(publisherIntent types.PublisherIntent) (stop bool) {
		publisherIntents = append(publisherIntents, publisherIntent)
		return false
	})

	return
}

//////////////////////
// SubscriberIntent //
//////////////////////

func (k Keeper) SetSubscriberIntent(ctx sdk.Context, subscriberAddress sdk.AccAddress, subscriberIntent types.SubscriberIntent) {
	bz := k.cdc.MustMarshal(&subscriberIntent)
	ctx.KVStore(k.storeKey).Set(types.GetSubscriberIntentKey(subscriberAddress, subscriberIntent.SubscsriptionId), bz)
}

func (k Keeper) GetSubscriberIntent(ctx sdk.Context, subscriberAddress sdk.AccAddress, subscriptionId string) (subscriberIntent types.SubscriberIntent, found bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetSubscriberIntentKey(subscriberAddress, subscriptionId))
	if len(bz) == 0 {
		return types.SubscriberIntent{}, false
	}

	k.cdc.MustUnmarshal(bz, &subscriberIntent)
	return subscriberIntent, true
}

func (k Keeper) IterateSubscriberIntents(ctx sdk.Context, handler func(subscriberAddress sdk.AccAddress, subscriberIntent types.SubscriberIntent) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, []byte{types.SubscriberIntentKeyPrefix})
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var subscriberIntent types.SubscriberIntent
		k.cdc.MustUnmarshal(iter.Value(), &subscriberIntent)

		addressKey := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), []byte{types.SubscriberIntentKeyPrefix}))
		subscriberAddress := sdk.AccAddress(addressKey.Next(20))

		if handler(subscriberAddress, subscriberIntent) {
			break
		}
	}
}

func (k Keeper) GetSubscriberIntents(ctx sdk.Context) (subscriberIntents []types.SubscriberIntent) {
	k.IterateSubscriberIntents(ctx, func(subscriberAddress sdk.AccAddress, subscriberIntent types.SubscriberIntent) (stop bool) {
		subscriberIntents = append(subscriberIntents, subscriberIntent)
		return false
	})

	return
}
