package v2

import (
	"testing"

	_ "github.com/peggyjv/sommelier/v10/app/params" // sets the somm bech32 address prefixes via init()
	"github.com/stretchr/testify/require"
)

func TestCorkAuthorityValidation(t *testing.T) {
	// Empty is allowed at the type level; the msg server is what fails closed.
	require.NoError(t, validateCorkAuthority(""))
	require.NoError(t, validateCorkAuthority("somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa"))

	require.Error(t, validateCorkAuthority("not-bech32"))
	require.Error(t, validateCorkAuthority("cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu"))
	require.Error(t, validateCorkAuthority(12345))
}

func TestDefaultParamsHasEmptyAuthority(t *testing.T) {
	require.Equal(t, "", DefaultParams().CorkAuthority)
}
