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

// variable length string fields being stored under these keys are hashed to provide a fixed size key
const (
	_ = byte(iota)

	// PublisherPrefix - <prefix>|<publisher_domain_hash> -> Publisher
	PublisherKeyPrefix

	// SubscriberPrefix - <prefix>|<subscriber_address_bytes> -> Subscriber
	SubscriberKeyPrefix

	// PublisherIntentByPublisherDomainKeyPrefix - <prefix>|<publisher_domain_hash>|<subscription_id_hash> -> PublisherIntent
	PublisherIntentByPublisherDomainKeyPrefix

	// PublisherIntentBySubscriptionIDKeyPrefix - <prefix>|<subscription_id_hash>|<publisher_domain_hash> -> PublisherIntent
	PublisherIntentBySubscriptionIDKeyPrefix

	// SubscriberIntentBySubscriberAddressKeyPrefix - <prefix>|<subscriber_address_bytes>|<subscription_id_hash> -> SubscriberIntent
	SubscriberIntentBySubscriberAddressKeyPrefix

	// SubscriberIntentBySubscriptionIDKeyPrefix - <prefix>|<subscription_id_hash>|<subscriber_address_bytes> -> SubscriberIntent
	SubscriberIntentBySubscriptionIDKeyPrefix

	// SubscriberIntentByPublisherDomainKeyPrefix - <prefix>|<publisher_domain_hash>|<subscriber_address_bytes>|<subscription_id_hash> -> SubscriberIntent
	SubscriberIntentByPublisherDomainKeyPrefix

	// DefaultSubscriptionKeyPrefix - <prefix>|<subscription_id_hash> -> DefaultSubscription
	DefaultSubscriptionKeyPrefix
)

func delimiter() []byte {
	return []byte{}
}

// GetPublishersPrefix returns a prefix for iterating all Publishers
func GetPublishersPrefix() []byte {
	return append([]byte{PublisherKeyPrefix}, delimiter()...)
}

// GetPublisherKey returns the key for a Publisher
func GetPublisherKey(publisherDomain string) []byte {
	return bytes.Join([][]byte{{PublisherKeyPrefix}, StringHash(publisherDomain)}, delimiter())
}

// GetSubscribersPrefix returns a prefix for iterating all Subscribers
func GetSubscribersPrefix() []byte {
	return append([]byte{SubscriberKeyPrefix}, delimiter()...)
}

// GetSubscriberKey returns the key for a Subscriber
func GetSubscriberKey(subscriberAddress sdk.AccAddress) []byte {
	return bytes.Join([][]byte{{SubscriberKeyPrefix}, subscriberAddress.Bytes()}, delimiter())
}

// GetPublisherIntentsPrefix returns a prefix for iterating all PublisherIntents by choosing the domain index
// as the default
func GetPublisherIntentsPrefix() []byte {
	return append([]byte{PublisherIntentByPublisherDomainKeyPrefix}, delimiter()...)
}

// GetPublisherIntentsByPublisherDomainPrefix returns a prefix for all PublisherIntents indexed by publisher domain
func GetPublisherIntentsByPublisherDomainPrefix(publisherDomain string) []byte {
	key := bytes.Join([][]byte{{PublisherIntentByPublisherDomainKeyPrefix}, StringHash(publisherDomain)}, delimiter())
	return append(key, delimiter()...)
}

// GetPublisherIntentsBySubscriptionIDPrefix returns a prefix for all PublisherIntents indexed by subscription ID
func GetPublisherIntentsBySubscriptionIDPrefix(subscriptionID string) []byte {
	key := bytes.Join([][]byte{{PublisherIntentBySubscriptionIDKeyPrefix}, StringHash(subscriptionID)}, delimiter())
	return append(key, delimiter()...)
}

// GetPublisherIntentByPublisherDomainKey returns the key for a PublisherIntent indexed by domain
func GetPublisherIntentByPublisherDomainKey(publisherDomain string, subscriptionID string) []byte {
	return bytes.Join([][]byte{{PublisherIntentByPublisherDomainKeyPrefix}, StringHash(publisherDomain), StringHash(subscriptionID)}, delimiter())
}

// GetPublisherIntentBySubscriptionIDKey returns the key for a PublisherIntent indexed by subscription ID
func GetPublisherIntentBySubscriptionIDKey(subsciptionID string, publisherDomain string) []byte {
	key := bytes.Join([][]byte{{PublisherIntentBySubscriptionIDKeyPrefix}, StringHash(subsciptionID), StringHash(publisherDomain)}, delimiter())
	return append(key, delimiter()...)
}

// GetSubscriberIntentsPrefix returns a prefix for iterating all SubscriberIntents by choosing the address index
// as the default
func GetSubscriberIntentsPrefix() []byte {
	return append([]byte{SubscriberIntentBySubscriberAddressKeyPrefix}, delimiter()...)
}

// GetSubscriberIntentsBySubscriberAddressPrefix returns a prefix for all SubscriberIntents indexed by address
func GetSubscriberIntentsBySubscriberAddressPrefix(subscriberAddress sdk.AccAddress) []byte {
	key := bytes.Join([][]byte{{SubscriberIntentBySubscriberAddressKeyPrefix}, subscriberAddress.Bytes()}, delimiter())
	return append(key, delimiter()...)
}

// GetSubscriberIntentsBySubscriptionIDPrefix returns a prefix for all SubscriberIntents indexed by subscription ID
func GetSubscriberIntentsBySubscriptionIDPrefix(subscriptionID string) []byte {
	key := bytes.Join([][]byte{{SubscriberIntentBySubscriptionIDKeyPrefix}, StringHash(subscriptionID)}, delimiter())
	return append(key, delimiter()...)
}

// GetSubscriberIntentsByPublisherDomainPrefix returns a prefix for all SubscriberIntents indexed by publisher domain
func GetSubscriberIntentsByPublisherDomainPrefix(publisherDomain string) []byte {
	key := bytes.Join([][]byte{{SubscriberIntentByPublisherDomainKeyPrefix}, StringHash(publisherDomain)}, delimiter())
	return append(key, delimiter()...)
}

// GetSubscriberIntentBySubscriberAddressKey returns the key for a SubscriberIntent indexed by address
func GetSubscriberIntentBySubscriberAddressKey(subscriberAddress sdk.AccAddress, subscriptionID string) []byte {
	return bytes.Join([][]byte{{SubscriberIntentBySubscriberAddressKeyPrefix}, subscriberAddress.Bytes(), StringHash(subscriptionID)}, delimiter())
}

// GetSubscriberIntentBySubscriptionIDKey returns the key for a SubscriberIntent indexed by subscription ID
func GetSubscriberIntentBySubscriptionIDKey(subscriptionID string, subscriberAddress sdk.AccAddress) []byte {
	return bytes.Join([][]byte{{SubscriberIntentBySubscriptionIDKeyPrefix}, StringHash(subscriptionID), subscriberAddress.Bytes()}, delimiter())
}

// GetSubscriberIntentByPublisherDomainKey returns the key for a SubscriberIntent indexed by publisher domain
func GetSubscriberIntentByPublisherDomainKey(publisherDomain string, subscriberAddress sdk.AccAddress, subscriptionID string) []byte {
	return bytes.Join([][]byte{{SubscriberIntentByPublisherDomainKeyPrefix}, StringHash(publisherDomain), subscriberAddress.Bytes(), StringHash(subscriptionID)}, delimiter())
}

// GetDefaultSubscriptionPrefix returns a prefix for iterating all DefaultSubscripions
func GetDefaultSubscriptionPrefix() []byte {
	return append([]byte{DefaultSubscriptionKeyPrefix}, delimiter()...)
}

// GetDefaultSubscriptionKey returns the key for a DefaultSubscription
func GetDefaultSubscriptionKey(subsciptionID string) []byte {
	return bytes.Join([][]byte{{DefaultSubscriptionKeyPrefix}, StringHash(subsciptionID)}, delimiter())
}
