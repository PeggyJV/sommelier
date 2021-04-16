package simulation

import (
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding distribution type
func DecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	// TODO: reimplement store decoding
	return func(kvA, kvB kv.Pair) string { return "" }
}
