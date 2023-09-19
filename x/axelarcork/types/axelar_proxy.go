package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	Address, _    = abi.NewType("address", "", nil)
	AddressArr, _ = abi.NewType("address[]", "", nil)
	Bytes, _      = abi.NewType("bytes", "", nil)
	Uint256, _    = abi.NewType("uint256", "", nil)

	HourInBlocks                     = uint64((60 * 60) / 12)
	DefaultExecutableHeightThreshold = 72 * HourInBlocks

	LogicCallMsgID = big.NewInt(0)
	UpgradeMsgID   = big.NewInt(1)
)

func EncodeLogicCallArgs(targetContract string, nonce uint64, deadline uint64, callData []byte) ([]byte, error) {
	return abi.Arguments{
		{Type: Uint256},
		{Type: Address},
		{Type: Uint256},
		{Type: Uint256},
		{Type: Bytes},
	}.Pack(LogicCallMsgID, common.HexToAddress(targetContract), big.NewInt(int64(nonce)), big.NewInt(int64(deadline)), callData)
}

func DecodeLogicCallArgs(data []byte) (targetContract string, nonce uint64, deadline uint64, callData []byte, err error) {
	args, err := abi.Arguments{
		{Type: Uint256},
		{Type: Address},
		{Type: Uint256},
		{Type: Uint256},
		{Type: Bytes},
	}.Unpack(data)
	if err != nil {
		return
	}

	targetContract = args[1].(common.Address).String()
	nonce = args[2].(*big.Int).Uint64()
	deadline = args[3].(*big.Int).Uint64()
	callData = args[4].([]byte)

	return
}

func EncodeUpgradeArgs(newAxelarProxy string, targets []string) ([]byte, error) {
	targetAddrs := []common.Address{}
	for _, target := range targets {
		targetAddrs = append(targetAddrs, common.HexToAddress(target))
	}

	return abi.Arguments{
		{Type: Uint256},
		{Type: Address},
		{Type: AddressArr},
	}.Pack(UpgradeMsgID, common.HexToAddress(newAxelarProxy), targetAddrs)
}

func DecodeUpgradeArgs(data []byte) (newAxelarProxy string, targets []string, err error) {
	args, err := abi.Arguments{
		{Type: Uint256},
		{Type: Address},
		{Type: AddressArr},
	}.Unpack(data)
	if err != nil {
		return
	}

	newAxelarProxy = args[1].(common.Address).String()
	targets = []string{}
	for _, target := range args[2].([]common.Address) {
		targets = append(targets, target.String())
	}

	return
}
