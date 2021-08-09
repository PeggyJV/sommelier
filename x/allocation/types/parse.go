package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TruncateDec splits a decimal into the integer and decimal components and then
// truncates the decimals in case it has a precision larger than the max allowed
// one (18).
func TruncateDec(decStr string) (sdk.Dec, error) {
	dec := strings.Split(decStr, ".")
	if len(dec) != 2 {
		value, ok := sdk.NewIntFromString(decStr)
		if !ok {
			return sdk.Dec{}, sdk.ErrInvalidDecimalStr
		}

		return sdk.NewDecFromInt(value), nil
	}

	if len(dec[1]) > sdk.Precision {
		dec[1] = dec[1][0:sdk.Precision]
	}

	return sdk.NewDecFromStr(strings.Join(dec, "."))
}

// MustTruncateDec is a util function that panics on TruncateDec error.
