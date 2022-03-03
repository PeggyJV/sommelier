package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "pubsub"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// note that we are using "|" as a delimiter for variable length keys as it's not a valid character for either,
// and in order to avoid shorter keys that are prefixes of existing keys we have a trailing delimiter too...it's
// not strictly necessary for e.g. the SubscriberKeyPrefix since all the fields are fixed, but keeping them all
// consistent makes it easier to read and reason about
const (
	_ = byte(iota)

	// PublisherPrefix - <prefix>|<publisher_domain>| -> Publisher
	PublisherKeyPrefix

	// SubscriberPrefix - <prefix>|<subscriber_address>| -> Subscriber
	SubscriberKeyPrefix

	// PublisherIntentByPublisherDomainKeyPrefix - <prefix>|<publisher_domain>|<subscription_id>| -> PublisherIntent
	PublisherIntentByPublisherDomainKeyPrefix

	// PublisherIntentBySubscriptionIdKeyPrefix - <prefix>|<subscription_id>|<publisher_domain>| -> PublisherIntent
	PublisherIntentBySubscriptionIdKeyPrefix

	// SubscriberIntentBySubscriberAddressKeyPrefix - <prefix>|<subscriber_address>|<subscription_id>| -> SubscriberIntent
	SubscriberIntentBySubscriberAddressKeyPrefix

	// SubscriberIntentBySubscriptionIdKeyPrefix - <prefix>|<subscription_id>|<subscriber_address>| -> SubscriberIntent
	SubscriberIntentBySubscriptionIdKeyPrefix
)

func delimiter() []byte {
	return []byte("|")
}

// GetPublishersPrefix returns a prefix for iterating all Publishers
func GetPublishersPrefix() []byte {
	return append([]byte{PublisherKeyPrefix}, delimiter()...)
}

// GetPublisherKey returns the key for a Publisher
func GetPublisherKey(publisherDomain string) []byte {
	key := bytes.Join([][]byte{{PublisherKeyPrefix}, []byte(publisherDomain)}, delimiter())
	return append(key, delimiter()...)
}

// GetSubscribersPrefix returns a prefix for iterating all Subscribers
func GetSubscribersPrefix() []byte {
	return append([]byte{SubscriberKeyPrefix}, delimiter()...)
}

// GetSubscriberKey returns the key for a Subscriber
func GetSubscriberKey(subscriberAddress sdk.AccAddress) []byte {
	key := bytes.Join([][]byte{{SubscriberKeyPrefix}, subscriberAddress.Bytes()}, delimiter())
	return append(key, delimiter()...)

}

// GetPublisherIntentsPrefix returns a prefix for iterating all PublisherIntents by choosing the domain index
// as the default
func GetPublisherIntentsPrefix() []byte {
	return append([]byte{PublisherIntentByPublisherDomainKeyPrefix}, delimiter()...)
}

// GetPublisherIntentsByPublisherDomainPrefix returns a prefix for all PublisherIntents indexed by publisher domain
func GetPublisherIntentsByPublisherDomainPrefix(publisherDomain string) []byte {
	key := bytes.Join([][]byte{{PublisherIntentByPublisherDomainKeyPrefix}, []byte(publisherDomain)}, delimiter())
	return append(key, delimiter()...)
}

// GetPublisherIntentsBySubscriptionIdPrefix returns a prefix for all PublisherIntents indexed by subscription ID
func GetPublisherIntentsBySubscriptionIdPrefix(subscriptionId string) []byte {
	key := bytes.Join([][]byte{{PublisherIntentBySubscriptionIdKeyPrefix}, []byte(subscriptionId)}, delimiter())
	return append(key, delimiter()...)
}

// GetPublisherIntentByPublisherDomainKey returns the key for a PublisherIntent indexed by domain
func GetPublisherIntentByPublisherDomainKey(publisherDomain string, subscriptionId string) []byte {
	key := bytes.Join([][]byte{{PublisherIntentByPublisherDomainKeyPrefix}, []byte(publisherDomain), []byte(subscriptionId)}, delimiter())
	return append(key, delimiter()...)
}

// GetPublisherIntentBySubscriptionIdKey returns the key for a PublisherIntent indexed by subscription ID
func GetPublisherIntentBySubscriptionIdKey(subsciptionId string, publisherDomain string) []byte {
	key := bytes.Join([][]byte{{PublisherIntentBySubscriptionIdKeyPrefix}, []byte(subsciptionId), []byte(publisherDomain)}, delimiter())
	return append(key, delimiter()...)
}

// GetSubscriberIntentsPrefix returns a prefix for iterating all SubscriberIntents by choosing the address index
// as the default
// have to tag the delimiter on the end here for simplicity in the keeper
func GetSubscriberIntentsPrefix() []byte {
	return append([]byte{SubscriberIntentBySubscriberAddressKeyPrefix}, delimiter()...)
}

// GetSubscriberIntentsBySubscriberAddressPrefix returns a prefix for all SubscriberIntents indexed by address
func GetSubscriberIntentsBySubscriberAddressPrefix(subscriberAddress sdk.AccAddress) []byte {
	key := bytes.Join([][]byte{{SubscriberIntentBySubscriberAddressKeyPrefix}, subscriberAddress.Bytes()}, delimiter())
	return append(key, delimiter()...)
}

// GetSubscriberIntentsBySubscriptionIdPrefix returns a prefix for all SubscriberIntents indexed by subscription ID
func GetSubscriberIntentsBySubscriptionIdPrefix(subscriptionId string) []byte {
	key := bytes.Join([][]byte{{SubscriberIntentBySubscriptionIdKeyPrefix}, []byte(subscriptionId)}, delimiter())
	return append(key, delimiter()...)
}

// GetSubscriberIntentBySubscriberAddressKey returns the key for a SubscriberIntent indexed by address
func GetSubscriberIntentBySubscriberAddressKey(subscriberAddress sdk.AccAddress, subscriptionId string) []byte {
	key := bytes.Join([][]byte{{SubscriberIntentBySubscriberAddressKeyPrefix}, subscriberAddress.Bytes(), []byte(subscriptionId)}, delimiter())
	return append(key, delimiter()...)
}

// GetSubscriberIntentBySubscriptionIdKey returns the key for a SubscriberIntent indexed by subscription ID
// since the last element is a fixed size we don't need a trailing delimiter
func GetSubscriberIntentBySubscriptionIdKey(subsciptionId string, subscriberAddress sdk.AccAddress) []byte {
	return bytes.Join([][]byte{{SubscriberIntentBySubscriptionIdKeyPrefix}, []byte(subsciptionId), subscriberAddress.Bytes()}, delimiter())
}
