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
	if goldHash != hex.EncodeToString(ourHash) {
		t.Fail()
	}
}
