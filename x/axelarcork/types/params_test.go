package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	// The SDK's bech32 prefix config is process-global and defaults to "cosmos";
	// this installs the "somm" prefixes the authority address is encoded with.
	_ "github.com/peggyjv/sommelier/v10/app/params"
)

func TestCorkAuthorityValidation(t *testing.T) {
	// Empty is allowed at the type level. Fail-closed enforcement lives in the
	// msg server, not here: the param must be settable to empty so governance
	// can revoke the authority without the param validator refusing the write.
	require.NoError(t, validateCorkAuthority(""))
	require.NoError(t, validateCorkAuthority("somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa"))

	require.Error(t, validateCorkAuthority("not-bech32"))
	require.Error(t, validateCorkAuthority("cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu"))
	require.Error(t, validateCorkAuthority(12345))
}

func TestDefaultParamsHasEmptyAuthority(t *testing.T) {
	require.Equal(t, "", DefaultParams().CorkAuthority)
}

func TestValidateBasicRejectsMalformedAuthority(t *testing.T) {
	p := DefaultParams()
	p.CorkAuthority = "not-bech32"
	require.Error(t, p.ValidateBasic())
}
