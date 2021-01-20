package eip712

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core"
)

const version = "1"

// EIP712Domain defines a EIP712 compatible domain for separation purposes
type EIP712Domain struct {
	// Name is the user readable name of signing domain, i.e. the name of the DApp or the protocol.
	Name string `json:"name"`
	// Version defines the current major version of the signing domain. Signatures from different versions are not compatible.
	Version string `json:"version"`
	// VerifyingContract is the address of the contract that will verify the signature. The user-agent may do contract specific phishing prevention
	VerifyingContract common.Address `json:"verifyingContract"`
	// ChainID is the the EIP-155 chain id.
	ChainID *big.Int `json:"chainID"`
	// Salt defines an disambiguating salt for the protocol. This can be used as a domain separator of last resort.
	Salt [32]byte `json:"salt"`
}

// Token defines a token
type Token struct {
	Symbol   string `json:"symbol"`
	Slippage float64
}

// EIP712DomainTypes defines the slice of signer types used for EIP712
var EIP712DomainTypes = []core.Type{
	{
		Name: "name",
		Type: "string",
	},
	{
		Name: "version",
		Type: "string",
	},
	{
		Name: "verifyingContract",
		Type: "address",
	},
	{
		Name: "chainId",
		Type: "uint256",
	},
	{
		Name: "salt",
		Type: "bytes32",
	},
}

func sign(contract common.Address) core.TypedData {
	data := core.TypedData{
		Types: core.Types{
			"EIP712Domain": EIP712DomainTypes,
			// TODO: other types definitions
			"MessageType": []core.Type{},
		},
		PrimaryType: "MessageType",
		Domain: core.TypedDataDomain{
			Name:              "Sommelier",
			Version:           version,
			ChainId:           math.NewHexOrDecimal256(1), // Ethereum mainnet
			VerifyingContract: contract.String(),
			Salt:              "", // TODO: add salt
		},
		Message: core.TypedDataMessage{
			"asset1": "ETH",
			"asset2": "WETH",
			// ...
		},
	}
	return data
}

/*
Input: address in the call data of the tx

Send LP token from Uniswap
Check to which pair corresponds the LP token
Ratio of the tokens
Allowed slippage


Send 712 message
sdk.Msg
someone collects the 712 msg payloads into a msg


Output:


*/
