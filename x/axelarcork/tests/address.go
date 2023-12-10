package tests

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// AccAddress returns a random account address
func AccAddress() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	return sdk.AccAddress(addr)
}

func AccAddressFromBech32(t *testing.T, addr string) sdk.AccAddress {
	t.Helper()
	a, err := sdk.AccAddressFromBech32(addr)
	require.NoError(t, err)

	return a
}
