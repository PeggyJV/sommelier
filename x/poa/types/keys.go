package types

import "encoding/binary"

const (
	// ModuleName is the name of the PoA module.
	ModuleName = "poa"

	// StoreKey is the primary store key.
	StoreKey = ModuleName

	// RouterKey is the message route.
	RouterKey = ModuleName

	// QuerierRoute is the querier route.
	QuerierRoute = ModuleName
)

// KV store key prefixes.
var (
	ParamsKey                = []byte{0x01}
	AuthoritySetKey          = []byte{0x02}
	MultiplierSnapshotPrefix = []byte{0x03}
	// ActivationHeightKey stores the block height at which the PoA module first
	// became active (genesis init or v10 upgrade). Below this height no boost
	// was ever applied, so a missing snapshot is benign; at or above it a
	// missing snapshot indicates corruption and a slash must be refused.
	ActivationHeightKey = []byte{0x04}
)

// MultiplierSnapshotKey returns the store key for the snapshot at `height`,
// using big-endian encoding so iteration is height-ordered.
func MultiplierSnapshotKey(height int64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(height))
	return append(MultiplierSnapshotPrefix, bz...)
}
