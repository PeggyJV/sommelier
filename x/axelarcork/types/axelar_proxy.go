package types

import (
	fmt "fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
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

func (ccn AxelarContractCallNonce) ValidateBasic() error {
	if ccn.Nonce == 0 {
		return fmt.Errorf("nonce cannot be zero")
	}

	if ccn.ContractAddress == "" {
		return fmt.Errorf("contract address cannot be empty")
	}

	if !common.IsHexAddress(ccn.ContractAddress) {
		return errorsmod.Wrapf(ErrInvalidEVMAddress, "%s", ccn.ContractAddress)
	}

	if ccn.ChainId == 0 {
		return fmt.Errorf("chain ID cannot be zero")
	}

	return nil
}

func (ud AxelarUpgradeData) ValidateBasic() error {
	if ud.ChainId == 0 {
		return fmt.Errorf("chain ID cannot be zero")
	}

	if ud.ExecutableHeightThreshold < 0 {
		return fmt.Errorf("height threshold cannot be negative")
	}

	if len(ud.Payload) == 0 {
		return fmt.Errorf("payload cannot be empty")
	}

	return nil
}

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

	// Unpack(data) error will catch overflow errors, but this will handle unambiguously invalid data
	if msgID, ok := args[0].(*big.Int); !ok || msgID.Cmp(LogicCallMsgID) != 0 {
		return "", 0, 0, nil, fmt.Errorf("invalid logic call args")
	}

	// Collin: Do we need to check for zero values in case these assertions somehow fail?
	targetContractArg, _ := args[1].(common.Address)
	targetContract = targetContractArg.String()
	nonceArg, _ := args[2].(*big.Int)
	nonce = nonceArg.Uint64()
	deadlineArg, _ := args[3].(*big.Int)
	deadline = deadlineArg.Uint64()
	callData, _ = args[4].([]byte)

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

	// Unpack(data) error will catch overflow errors, but this will handle unambiguously invalid data
	if msgID, ok := args[0].(*big.Int); !ok || msgID.Cmp(UpgradeMsgID) != 0 {
		return "", nil, fmt.Errorf("invalid upgrade args")
	}

	// Collin: Do we need to check for zero values in case these assertions somehow fail?
	newAxelarProxyArg, _ := args[1].(common.Address)
	newAxelarProxy = newAxelarProxyArg.String()
	targetsArg, _ := args[2].([]common.Address)

	targets = []string{}
	for _, target := range targetsArg {
		targets = append(targets, target.String())
	}

	return
}
