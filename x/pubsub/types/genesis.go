package types

const DefaultParamspace = ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:               DefaultParams(),
		Publishers:           []*Publisher{},
		Subscribers:          []*Subscriber{},
		PublisherIntents:     []*PublisherIntent{},
		SubscriberIntents:    []*SubscriberIntent{},
		DefaultSubscriptions: []*DefaultSubscription{},
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	for _, publisher := range gs.Publishers {
		if err := publisher.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, subscriber := range gs.Subscribers {
		if err := subscriber.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, publisherIntent := range gs.PublisherIntents {
		if err := publisherIntent.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, subscriberIntent := range gs.SubscriberIntents {
		if err := subscriberIntent.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, defaultSubscription := range gs.DefaultSubscriptions {
		if err := defaultSubscription.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}
