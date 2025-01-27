package v2

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestValidateProceedsPortion(t *testing.T) {
	params := DefaultParams()

	params.ProceedsPortion = sdk.MustNewDecFromStr("-0.1")

	err := ValidateProceedsPortion(params.ProceedsPortion)
	require.Error(t, err)

	params.ProceedsPortion = sdk.MustNewDecFromStr("1.1")

	err = ValidateProceedsPortion(params.ProceedsPortion)
	require.Error(t, err)

	params.ProceedsPortion = sdk.MustNewDecFromStr("0")
	err = ValidateProceedsPortion(params.ProceedsPortion)
	require.NoError(t, err)

	params.ProceedsPortion = sdk.MustNewDecFromStr("1")
	err = ValidateProceedsPortion(params.ProceedsPortion)
	require.NoError(t, err)

	params.ProceedsPortion = sdk.MustNewDecFromStr("0.5")
	err = ValidateProceedsPortion(params.ProceedsPortion)
	require.NoError(t, err)
}
