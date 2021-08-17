package types

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestContractCallTxCheckpoint(t *testing.T) {

	call := Cellar{
		Id: "0x0000000000",
		TickRanges: []*TickRange{
			// TODO: replace with actual values from the rust code
			{-189780, -192480, 160},
			{-192480, -197880, 680},
			{-197880, -200640, 160},
		},
	}

	ourHash := call.GetCheckpoint()

	// hash from rust code
	goldHash := "0x21942107e0af59eee8101c5c940a932e8fdf243c6e9aff26b750b7b2eda62464"[2:]
	fmt.Printf("%x\n", ourHash)
	testHash := hex.EncodeToString(ourHash)
	if goldHash != testHash {
		t.Errorf("gold hash is not equal to generated hash:\n gold hash: %v\n test hash: %v", goldHash, testHash)
	}
}
