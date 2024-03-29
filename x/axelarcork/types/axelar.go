package types

const (
	_ int64 = iota
	PureMessage
	MessageWithToken
	PureTokenTransfer
)

type AxelarBody struct {
	DestinationChain   string `json:"destination_chain"`
	DestinationAddress string `json:"destination_address"`
	Payload            []byte `json:"payload"`
	Type               int64  `json:"type"`
	Fee                *Fee   `json:"fee"`
}

type Fee struct {
	Amount    string `json:"amount"`
	Recipient string `json:"recipient"`
}
