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

// TODO(bolten): fill out keys

const (
	_ = byte(iota)

	// PublisherPrefix - <prefix><publisher_domain> -> Publisher
	PublisherPrefix

	// SubscriberPrefix - <prefix><subscriber_address> -> Subscriber
	SubscriberPrefix

	// PublisherIntentPrefix - <prefix><publisher_domain><subscription_id> -> PublisherIntent
	PublisherIntentPrefix

	// SubscriberIntentPrefix - <prefix><subscriber_address><subscription_id> -> SubscriberIntent
	SubscriberIntentPrefix
)

// GetPublisherKey returns the key for a Publisher
func GetPublisherKey(publisherDomain string) []byte {
	return append([]byte{PublisherPrefix}, []byte(publisherDomain)...)
}

// GetSubscriberKey returns the key for a Subscriber
func GetSubscriberKey(subscriberAddress sdk.Address) []byte {
	return append([]byte{SubscriberPrefix}, subscriberAddress.Bytes()...)
}

// GetPublisherIntentKey returns the key for a PublisherIntent
func GetPublisherIntentkey(publisherDomain string, subscriptionId string) []byte {
	return bytes.Join([][]byte{{PublisherIntentPrefix}, []byte(publisherDomain), []byte(subscriptionId)}, []byte{})
}

// GetSubscriberIntentKey returns the key for a SubscriberIntent
func GetSubscriberIntentKey(subscriberAddress sdk.Address, subscriptionId string) []byte {
	return bytes.Join([][]byte{{SubscriberIntentPrefix}, subscriberAddress.Bytes(), []byte(subscriptionId)}, []byte{})
}
