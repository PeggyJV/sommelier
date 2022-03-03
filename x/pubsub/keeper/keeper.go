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

////////////
// Params //
////////////

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

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
	iter := sdk.KVStorePrefixIterator(store, types.GetPublishersPrefix())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var publisher types.Publisher
		k.cdc.MustUnmarshal(iter.Value(), &publisher)

		if handler(publisher) {
			break
		}
	}
}

func (k Keeper) GetPublishers(ctx sdk.Context) (publishers []*types.Publisher) {
	k.IteratePublishers(ctx, func(publisher types.Publisher) (stop bool) {
		publishers = append(publishers, &publisher)
		return false
	})

	return
}

func (k Keeper) DeletePublisher(ctx sdk.Context, publisher types.Publisher) {
	ctx.KVStore(k.storeKey).Delete(types.GetPublisherKey(publisher.Domain))

	for _, publisherIntent := range k.GetPublisherIntentsByPublisherDomain(ctx, publisher.Domain) {
		k.DeletePublisherIntent(ctx, *publisherIntent)
	}
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

func (k Keeper) IterateSubscribers(ctx sdk.Context, handler func(subscriberAddress sdk.AccAddress, subscriber types.Subscriber) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetSubscribersPrefix())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var subscriber types.Subscriber
		k.cdc.MustUnmarshal(iter.Value(), &subscriber)

		addressKey := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), types.GetSubscribersPrefix()))
		subscriberAddress := sdk.AccAddress(addressKey.Bytes())

		if handler(subscriberAddress, subscriber) {
			break
		}
	}
}

func (k Keeper) GetSubscribers(ctx sdk.Context) (subscribers []*types.Subscriber) {
	k.IterateSubscribers(ctx, func(subscriberAddress sdk.AccAddress, subscriber types.Subscriber) (stop bool) {
		subscribers = append(subscribers, &subscriber)
		return false
	})

	return
}

func (k Keeper) DeleteSubscriber(ctx sdk.Context, subscriberAddress sdk.AccAddress, subscriber types.Subscriber) {
	ctx.KVStore(k.storeKey).Delete(types.GetSubscriberKey(subscriberAddress))
	for _, subscriberIntent := range k.GetSubscriberIntentsBySubscriberAddress(ctx, subscriberAddress) {
		k.DeleteSubscriberIntent(ctx, subscriberAddress, *subscriberIntent)
	}
}

/////////////////////
// PublisherIntent //
/////////////////////

func (k Keeper) SetPublisherIntent(ctx sdk.Context, publisherIntent types.PublisherIntent) {
	bz := k.cdc.MustMarshal(&publisherIntent)
	ctx.KVStore(k.storeKey).Set(types.GetPublisherIntentByPublisherDomainKey(publisherIntent.PublisherDomain, publisherIntent.SubscriptionId), bz)
	ctx.KVStore(k.storeKey).Set(types.GetPublisherIntentBySubscriptionIdKey(publisherIntent.SubscriptionId, publisherIntent.PublisherDomain), bz)
}

func (k Keeper) GetPublisherIntent(ctx sdk.Context, publisherDomain string, subscriptionId string) (publisherIntent types.PublisherIntent, found bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetPublisherIntentByPublisherDomainKey(publisherDomain, subscriptionId))
	if len(bz) == 0 {
		return types.PublisherIntent{}, false
	}

	k.cdc.MustUnmarshal(bz, &publisherIntent)
	return publisherIntent, true
}

func (k Keeper) IteratePublisherIntents(ctx sdk.Context, handler func(publisherIntent types.PublisherIntent) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetPublisherIntentsPrefix())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var publisherIntent types.PublisherIntent
		k.cdc.MustUnmarshal(iter.Value(), &publisherIntent)

		if handler(publisherIntent) {
			break
		}
	}
}

func (k Keeper) GetPublisherIntents(ctx sdk.Context) (publisherIntents []*types.PublisherIntent) {
	return k.getPublisherIntentsByKeyPrefix(ctx, types.GetPublisherIntentsPrefix())
}

func (k Keeper) GetPublisherIntentsByPublisherDomain(ctx sdk.Context, publisherDomain string) (publisherIntents []*types.PublisherIntent) {
	return k.getPublisherIntentsByKeyPrefix(ctx, types.GetPublisherIntentsByPublisherDomainPrefix(publisherDomain))
}

func (k Keeper) GetPublisherIntentsBySubscriptionId(ctx sdk.Context, subscriptionId string) (publisherIntents []*types.PublisherIntent) {
	return k.getPublisherIntentsByKeyPrefix(ctx, types.GetPublisherIntentsBySubscriptionIdPrefix(subscriptionId))
}

func (k Keeper) getPublisherIntentsByKeyPrefix(ctx sdk.Context, keyPrefix []byte) (publisherIntents []*types.PublisherIntent) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, keyPrefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var publisherIntent types.PublisherIntent
		k.cdc.MustUnmarshal(iter.Value(), &publisherIntent)
		publisherIntents = append(publisherIntents, &publisherIntent)
	}

	return
}

func (k Keeper) DeletePublisherIntent(ctx sdk.Context, publisherIntent types.PublisherIntent) {
	ctx.KVStore(k.storeKey).Delete(types.GetPublisherIntentByPublisherDomainKey(publisherIntent.PublisherDomain, publisherIntent.SubscriptionId))
	ctx.KVStore(k.storeKey).Delete(types.GetPublisherIntentBySubscriptionIdKey(publisherIntent.SubscriptionId, publisherIntent.PublisherDomain))

	for _, subscriberIntent := range k.GetSubscriberIntentsByPublisherDomain(ctx, publisherIntent.PublisherDomain) {
		if publisherIntent.SubscriptionId == subscriberIntent.SubscriptionId {
			k.DeleteSubscriberIntent(ctx, sdk.AccAddress(subscriberIntent.SubscriberAddress), *subscriberIntent)
		}
	}
}

//////////////////////
// SubscriberIntent //
//////////////////////

func (k Keeper) SetSubscriberIntent(ctx sdk.Context, subscriberAddress sdk.AccAddress, subscriberIntent types.SubscriberIntent) {
	bz := k.cdc.MustMarshal(&subscriberIntent)

	// TODO(bolten): we are storing three different indices to improve query processing speed at the cost of additional storage
	// is this the right trade-off, or should we use a single index and manually compose those query responses with iterators?
	ctx.KVStore(k.storeKey).Set(types.GetSubscriberIntentBySubscriberAddressKey(subscriberAddress, subscriberIntent.SubscriptionId), bz)
	ctx.KVStore(k.storeKey).Set(types.GetSubscriberIntentBySubscriptionIdKey(subscriberIntent.SubscriptionId, subscriberAddress), bz)
	ctx.KVStore(k.storeKey).Set(types.GetSubscriberIntentByPublisherDomainKey(subscriberIntent.PublisherDomain, subscriberAddress, subscriberIntent.SubscriptionId), bz)
}

func (k Keeper) GetSubscriberIntent(ctx sdk.Context, subscriberAddress sdk.AccAddress, subscriptionId string) (subscriberIntent types.SubscriberIntent, found bool) {
	bz := ctx.KVStore(k.storeKey).Get(types.GetSubscriberIntentBySubscriberAddressKey(subscriberAddress, subscriptionId))
	if len(bz) == 0 {
		return types.SubscriberIntent{}, false
	}

	k.cdc.MustUnmarshal(bz, &subscriberIntent)
	return subscriberIntent, true
}

func (k Keeper) IterateSubscriberIntents(ctx sdk.Context, handler func(subscriberAddress sdk.AccAddress, subscriberIntent types.SubscriberIntent) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetSubscriberIntentsPrefix())
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var subscriberIntent types.SubscriberIntent
		k.cdc.MustUnmarshal(iter.Value(), &subscriberIntent)

		addressKey := bytes.NewBuffer(bytes.TrimPrefix(iter.Key(), types.GetSubscriberIntentsPrefix()))
		subscriberAddress := sdk.AccAddress(addressKey.Next(20))

		if handler(subscriberAddress, subscriberIntent) {
			break
		}
	}
}

func (k Keeper) GetSubscriberIntents(ctx sdk.Context) (subscriberIntents []*types.SubscriberIntent) {
	return k.getSubscriberIntentsByKeyPrefix(ctx, types.GetSubscriberIntentsPrefix())
}

func (k Keeper) GetSubscriberIntentsBySubscriberAddress(ctx sdk.Context, subscriberAddress sdk.AccAddress) (subscriberIntents []*types.SubscriberIntent) {
	return k.getSubscriberIntentsByKeyPrefix(ctx, types.GetSubscriberIntentsBySubscriberAddressPrefix(subscriberAddress))
}

func (k Keeper) GetSubscriberIntentsBySubscriptionId(ctx sdk.Context, subscriptionId string) (subscriberIntents []*types.SubscriberIntent) {
	return k.getSubscriberIntentsByKeyPrefix(ctx, types.GetSubscriberIntentsBySubscriptionIdPrefix(subscriptionId))
}

func (k Keeper) GetSubscriberIntentsByPublisherDomain(ctx sdk.Context, publisherDomain string) (subscriberIntents []*types.SubscriberIntent) {
	return k.getSubscriberIntentsByKeyPrefix(ctx, types.GetSubscriberIntentsByPublisherDomainPrefix(publisherDomain))
}

func (k Keeper) getSubscriberIntentsByKeyPrefix(ctx sdk.Context, keyPrefix []byte) (subscriberIntents []*types.SubscriberIntent) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, keyPrefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var subscriberIntent types.SubscriberIntent
		k.cdc.MustUnmarshal(iter.Value(), &subscriberIntent)
		subscriberIntents = append(subscriberIntents, &subscriberIntent)
	}

	return
}

func (k Keeper) DeleteSubscriberIntent(ctx sdk.Context, subscriberAddress sdk.AccAddress, subscriberIntent types.SubscriberIntent) {
	ctx.KVStore(k.storeKey).Delete(types.GetSubscriberIntentBySubscriberAddressKey(subscriberAddress, subscriberIntent.SubscriptionId))
	ctx.KVStore(k.storeKey).Delete(types.GetSubscriberIntentBySubscriptionIdKey(subscriberIntent.SubscriptionId, subscriberAddress))
	ctx.KVStore(k.storeKey).Delete(types.GetSubscriberIntentByPublisherDomainKey(subscriberIntent.PublisherDomain, subscriberAddress, subscriberIntent.SubscriptionId))
}
