package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	Address, _    = abi.NewType("address", "", nil)
	AddressArr, _ = abi.NewType("address[]", "", nil)
	Bytes, _      = abi.NewType("bytes", "", nil)
	Uint256, _    = abi.NewType("uint256", "", nil)

	HourInBlocks                     = uint64((60 * 60) / 12)
	DefaultExecutableHeightThreshold = 72 * HourInBlocks
)

type UpgradeData struct {
	Payload                  []byte
	ExecutionThresholdHeight uint64
}

func EncodeExecuteArgs(targetContract string, nonce uint64, deadline uint64, callData []byte) ([]byte, error) {
	return abi.Arguments{
		{Type: Uint256},
		{Type: Address},
		{Type: Uint256},
		{Type: Uint256},
		{Type: Bytes},
	}.Pack(big.NewInt(0), targetContract, big.NewInt(int64(nonce)), big.NewInt(int64(deadline)), callData)
}

func EncodeUpgradeArgs(newAxelarProxy string, targets []string) ([]byte, error) {
	return abi.Arguments{
		{Type: Uint256},
		{Type: Address},
		{Type: AddressArr},
	}.Pack(big.NewInt(1), newAxelarProxy, targets)
}
