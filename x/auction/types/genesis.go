package types

const DefaultParamspace = ModuleName

// DefaultGenesisState get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:        DefaultParams(),
		Auctions:      []*Auction{},
		Bids:          []*Bid{},
		TokenPrices:   []*TokenPrice{},
		LastAuctionId: 0,
		LastBidId:     0,
	}
}

// Validate performs a basic stateless validation of the genesis fields.
func (gs GenesisState) Validate() error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}

	for _, auction := range gs.Auctions {
		if err := auction.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, bid := range gs.Bids {
		if err := bid.ValidateBasic(); err != nil {
			return err
		}
	}

	for _, tokenPrice := range gs.TokenPrices {
		if err := tokenPrice.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}
