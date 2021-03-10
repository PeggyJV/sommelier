package types

import (
	"encoding/hex"
	"testing"

	oracle "github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestRedeemLiquidtiyGold1(t *testing.T) {
	var (
		erc20Addr = "0x835973768750b3ED2D5c3EF5AdcD5eDb44d12aD4"
	)

	src := RedeemLiquidityCall{
		Pair: oracle.UniswapPair{},
	}
	// TODO: get from params
	ourHash, err := src.GetCheckpoint("foo")
	require.NoError(t, err)

	// hash from bridge contract
	goldHash := "0xa3a7ee0a363b8ad2514e7ee8f110d7449c0d88f3b0913c28c1751e6e0079a9b2"[2:]
	// The function used to compute the "gold hash" above is in /solidity/test/updateValsetAndSubmitBatch.ts
	// Be aware that every time that you run the above .ts file, it will use a different tokenContractAddress and thus compute
	// a different hash.
	assert.Equal(t, goldHash, hex.EncodeToString(ourHash))
}
