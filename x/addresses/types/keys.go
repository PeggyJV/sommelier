package types

import "bytes"

const (
	// ModuleName defines the module name
	ModuleName = "addresses"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_addresses"
)
const (
	_ = byte(iota)

	// <prefix>
	CosmosToEvmMapPrefix

	// <prefix>
	EvmToCosmosMapPrefix
)

func GetCosmosToEvmMapPrefix() []byte {
	return []byte{CosmosToEvmMapPrefix}
}

func GetCosmosToEvmMapKey(cosmosAddr []byte) []byte {
	return bytes.Join([][]byte{{CosmosToEvmMapPrefix}, cosmosAddr}, []byte{})
}

func GetEvmToCosmosMapPrefix() []byte {
	return []byte{EvmToCosmosMapPrefix}
}

func GetEvmToCosmosMapKey(evmAddr []byte) []byte {
	return bytes.Join([][]byte{{EvmToCosmosMapPrefix}, evmAddr}, []byte{})
}
