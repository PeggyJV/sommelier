package types

import (
	"bytes"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	// ModuleName is the module name constant used in many places
	ModuleName = "axelarcork"

	// StoreKey is the store key string for oracle
	StoreKey = ModuleName

	// RouterKey is the message route for oracle
	RouterKey = ModuleName

	// QuerierRoute is the querier route for oracle
	QuerierRoute = ModuleName
)

// Keys for cork store, with <prefix><key> -> <value>
const (
	_ = byte(iota)

	// CorkForAddressKeyPrefix - <prefix><chain_id><val_address><address> -> <cork>
	CorkForAddressKeyPrefix // key for corks

	// CellarIDsKeyPrefix - <prefix><chain_id> -> []string
	CellarIDsKeyPrefix

	// ScheduledCorkKeyPrefix - <prefix><chain_id><block_height><val_address><address> -> <cork>
	ScheduledCorkKeyPrefix

	// CorkResultPrefix - <prefix><chain_id><id> -> AxelarCorkResult
	CorkResultPrefix

	// ChainConfigurationPrefix - <prefix><chain_id> -> ChainConfiguration
	ChainConfigurationPrefix

	// WinningCorkPrefix - <prefix><chain_id> -> AxelarCork
	WinningCorkPrefix

	// AxelarContractCallNoncePrefix - <prefix><chain_id><contract_address> -> <nonce>
	AxelarContractCallNoncePrefix

	// AxelarProxyUpgradeDataPrefix - <prefix><chain_id> -> <payload>
	AxelarProxyUpgradeDataPrefix
)

// GetCorkValidatorKeyPrefix returns the key prefix for cork commits for a validator
func GetCorkValidatorKeyPrefix(chainID uint64, val sdk.ValAddress) []byte {
	cid := make([]byte, 8)
	binary.BigEndian.PutUint64(cid, chainID)
	return bytes.Join([][]byte{{CorkForAddressKeyPrefix}, cid, val.Bytes()}, []byte{})
}

func MakeCellarIDsKey(chainID uint64) []byte {
	cid := make([]byte, 8)
	binary.BigEndian.PutUint64(cid, chainID)
	return bytes.Join([][]byte{{CellarIDsKeyPrefix}, cid}, []byte{})
}

func GetScheduledAxelarCorkKeyPrefix(chainID uint64) []byte {
	cid := make([]byte, 8)
	binary.BigEndian.PutUint64(cid, chainID)
	return bytes.Join([][]byte{{ScheduledCorkKeyPrefix}, cid}, []byte{})
}

func GetScheduledAxelarCorkKeyByBlockHeightPrefix(chainID uint64, blockHeight uint64) []byte {
	return append(GetScheduledAxelarCorkKeyPrefix(chainID), sdk.Uint64ToBigEndian(blockHeight)...)
}

func GetScheduledAxelarCorkKey(chainID uint64, blockHeight uint64, id []byte, val sdk.ValAddress, contract common.Address) []byte {
	blockHeightBytes := sdk.Uint64ToBigEndian(blockHeight)
	return bytes.Join([][]byte{GetScheduledAxelarCorkKeyPrefix(chainID), blockHeightBytes, id, val.Bytes(), contract.Bytes()}, []byte{})
}

func GetAxelarCorkResultPrefix(chainID uint64) []byte {
	cid := make([]byte, 8)
	binary.BigEndian.PutUint64(cid, chainID)
	return bytes.Join([][]byte{{CorkResultPrefix}, cid}, []byte{})
}

func GetAxelarCorkResultKey(chainID uint64, id []byte) []byte {
	return append(GetAxelarCorkResultPrefix(chainID), id...)
}

func ChainConfigurationKey(chainID uint64) []byte {
	cid := make([]byte, 8)
	binary.BigEndian.PutUint64(cid, chainID)
	return bytes.Join([][]byte{{ChainConfigurationPrefix}, cid}, []byte{})
}

func GetWinningAxelarCorkKeyPrefix(chainID uint64) []byte {
	cid := make([]byte, 8)
	binary.BigEndian.PutUint64(cid, chainID)
	return bytes.Join([][]byte{{WinningCorkPrefix}, cid}, []byte{})
}

func GetWinningAxelarCorkKey(chainID uint64, blockheight uint64, address common.Address) []byte {
	bh := make([]byte, 8)
	binary.BigEndian.PutUint64(bh, blockheight)
	return bytes.Join([][]byte{GetWinningAxelarCorkKeyPrefix(chainID), bh, address.Bytes()}, []byte{})
}

func GetAxelarContractCallNonceKey(chainID uint64, contractAddress common.Address) []byte {
	cid := make([]byte, 8)
	binary.BigEndian.PutUint64(cid, chainID)
	return bytes.Join([][]byte{{AxelarContractCallNoncePrefix}, cid, contractAddress.Bytes()}, []byte{})
}

func GetAxelarProxyUpgradeDataKey(chainID uint64) []byte {
	cid := make([]byte, 8)
	binary.BigEndian.PutUint64(cid, chainID)
	return bytes.Join([][]byte{{AxelarProxyUpgradeDataPrefix}, cid}, []byte{})
}
