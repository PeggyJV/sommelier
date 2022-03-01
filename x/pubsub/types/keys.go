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
	PublisherKeyPrefix

	// SubscriberPrefix - <prefix><subscriber_address> -> Subscriber
	SubscriberKeyPrefix

	// PublisherIntentPrefix - <prefix><publisher_domain><subscription_id> -> PublisherIntent
	PublisherIntentKeyPrefix

	// SubscriberIntentPrefix - <prefix><subscriber_address><subscription_id> -> SubscriberIntent
	SubscriberIntentKeyPrefix
)

// GetPublisherKey returns the key for a Publisher
func GetPublisherKey(publisherDomain string) []byte {
	return append([]byte{PublisherKeyPrefix}, []byte(publisherDomain)...)
}

// GetSubscriberKey returns the key for a Subscriber
func GetSubscriberKey(subscriberAddress sdk.AccAddress) []byte {
	return append([]byte{SubscriberKeyPrefix}, subscriberAddress.Bytes()...)
}

// GetPublisherIntentKey returns the key for a PublisherIntent
func GetPublisherIntentKey(publisherDomain string, subscriptionId string) []byte {
	return bytes.Join([][]byte{{PublisherIntentKeyPrefix}, []byte(publisherDomain), []byte(subscriptionId)}, []byte{})
}

// GetSubscriberIntentKey returns the key for a SubscriberIntent
func GetSubscriberIntentKey(subscriberAddress sdk.AccAddress, subscriptionId string) []byte {
	return bytes.Join([][]byte{{SubscriberIntentKeyPrefix}, subscriberAddress.Bytes(), []byte(subscriptionId)}, []byte{})
}
