package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewClaim generates a Claim instance.
func NewClaim(weight int64, recipient sdk.ValAddress) Claim {
	return Claim{
		Weight:    weight,
		Recipient: recipient.String(),
	}
}
