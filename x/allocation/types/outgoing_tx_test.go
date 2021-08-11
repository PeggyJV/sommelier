package types

import (
	"encoding/hex"
	"testing"
)

func TestContractCallTxCheckpoint(t *testing.T) {

	call := Cellar{
		Id: "0x0000000000",
		TickRanges: []*TickRange{
			// TODO: replace with actual values from the rust code
			{19394, 191466, 19214124},
			{19394, 191466, 19214124},
			{19394, 191466, 19214124},
		},
	}

	ourHash := call.GetCheckpoint()

	// hash from rust code
	goldHash := "0x000000000...."[2:]
	testHash := hex.EncodeToString(ourHash)
	if goldHash != testHash {
		t.Errorf("gold hash is not equal to generated hash:\n gold hash: %v\n test hash: %v", goldHash, testHash)
	}
}
