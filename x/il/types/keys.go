package types

const (
	// ModuleName is the name of the impermanent loss module
	ModuleName = "il"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the oracle module
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the oracle module
	QuerierRoute = ModuleName
)

// Keys for oracle store
// Items are stored with the following key: values
//
// - 0x01<denom_Bytes><valAddress_Bytes>: ExchangeRatePrevote
//
// - 0x02<denom_Bytes><valAddress_Bytes>: ExchangeRateVote
//
// - 0x03<denom_Bytes>: sdk.Dec
//
// - 0x04<valAddress_Bytes>: accAddress
//
// - 0x05<valAddress_Bytes>: int64
//
// - 0x06<valAddress_Bytes>: AggregateExchangeRatePrevote
//
// - 0x07<valAddress_Bytes>: AggregateExchangeRateVote
//
// - 0x08<denom_Bytes>: sdk.Dec
var (
	// Keys for store prefixes
	PrevoteKey                      = []byte{0x01} // prefix for each key to a prevote
	VoteKey                         = []byte{0x02} // prefix for each key to a vote
	ExchangeRateKey                 = []byte{0x03} // prefix for each key to a rate
	FeederDelegationKey             = []byte{0x04} // prefix for each key to a feeder delegation
	MissCounterKey                  = []byte{0x05} // prefix for each key to a miss counter
	AggregateExchangeRatePrevoteKey = []byte{0x06} // prefix for each key to a aggregate prevote
	AggregateExchangeRateVoteKey    = []byte{0x07} // prefix for each key to a aggregate vote
	TobinTaxKey                     = []byte{0x08} // prefix for each key to a tobin tax
)
