package types

// pubsub module event types
const (
	EventTypeAddPublisherIntent     = "add_publisher_intent"
	EventTypeAddSubscriberIntent    = "add_subscriber_intent"
	EventTypeAddSubscriber          = "add_subscriber"
	EventTypeRemovePublisherIntent  = "remove_publisher_intent"
	EventTypeRemoveSubscriberIntent = "remove_subscriber_intent"
	EventTypeRemoveSubscriber       = "remove_subscriber"
	EventTypeRemovePublisher        = "remove_publisher"

	AttributeKeyPublisherDomain   = "publisher_domain"
	AttributeKeyPublisherAddress  = "publisher_address"
	AttributeKeySubscriptionId    = "subscription_id"
	AttributeKeySubscriberAddress = "subscriber_address"

	AttributeValueCategory = ModuleName
)
