package types

type AxelarBody struct {
	DestinationChain   string `json:"destination_chain"`
	DestinationAddress string `json:"destination_address"`
	Payload            []byte `json:"payload"`
	Type               int    `json:"type"`
	Fee                struct {
		Amount    string `json:"amount"`
		Recipient string `json:"recipient"`
	} `json:"fee"`
}
