package integration_tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEthFromMnemonic(t *testing.T) {
	mnemonic := "receive roof marine sure lady hundred sea enact exist place bean wagon kingdom betray science photo loop funny bargain floor suspect only strike endless"
	address := "0x14fdAC734De10065093C4Ed4a83C41638378005A"

	generatedKey, err := ethereumKeyFromMnemonic(mnemonic)
	require.NoError(t, err, "error generating ethereum key")

	require.Equal(t, address, generatedKey.address)
}
