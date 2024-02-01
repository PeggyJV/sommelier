package types

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"
)

func TestEncodingMetadata(t *testing.T) {
	axelarMemo := AxelarBody{
		DestinationChain:   "arbitrum",
		DestinationAddress: "0xEe75bA2C81C04DcA4b0ED6d1B7077c188FEde4d2",
		Payload:            []byte("1234"),
		Type:               PureMessage,
		Fee: &Fee{
			Amount:    "100",
			Recipient: "axelar1aythygn6z5thymj6tmzfwekzh05ewg3l7d6y89",
		},
	}

	memoBz, err := json.Marshal(axelarMemo)
	require.NoError(t, err)

	var body AxelarBody
	err = json.Unmarshal(memoBz, &body)
	require.NoError(t, err)
}

func TestEncodingDecodingLogicCalls(t *testing.T) {
	targetContract := "0x1111111111111111111111111111111111111111"
	nonce := uint64(1)
	deadline := uint64(1000000000)
	callData := []byte("testdata")

	// Can encode and decode valid logic calls
	result, err := EncodeLogicCallArgs(targetContract, nonce, deadline, callData)
	require.NoError(t, err)
	actualTargetContract, actualNonce, actualDeadline, actualCallData, err := DecodeLogicCallArgs(result)
	require.NoError(t, err)
	require.Equal(t, targetContract, actualTargetContract)
	require.Equal(t, nonce, actualNonce)
	require.Equal(t, deadline, actualDeadline)
	require.Equal(t, callData, actualCallData)

	// Decoding logic call as upgrade caught
	_, _, err = DecodeUpgradeArgs(result)
	require.Error(t, err)
	require.Equal(t, err.Error(), "invalid upgrade args")

	// Specifically using the wrong msgID in a logic call errors
	wrongMsgID := result
	upgradeMsgIDBytes, err := abi.Arguments{{Type: Uint256}}.Pack(UpgradeMsgID)
	require.NoError(t, err)
	wrongMsgID = bytes.Join([][]byte{upgradeMsgIDBytes, wrongMsgID[len(upgradeMsgIDBytes):]}, []byte{})
	_, _, _, _, err = DecodeLogicCallArgs(wrongMsgID) //nolint:dogsled
	require.Error(t, err)
	require.Equal(t, err.Error(), "invalid logic call args")

	// Can encode and decode valid upgrade calls
	targets := []string{targetContract, "0x2222222222222222222222222222222222222222"}

	result, err = EncodeUpgradeArgs(targetContract, targets)
	require.NoError(t, err)
	actualNewAxelarProxy, actualTargets, err := DecodeUpgradeArgs(result)
	require.NoError(t, err)
	require.Equal(t, targetContract, actualNewAxelarProxy)
	require.Equal(t, targets, actualTargets)

	// Decoding upgrade call as logic call caught
	_, _, _, _, err = DecodeLogicCallArgs(result) //nolint:dogsled
	require.Error(t, err)

	// Specifically using the wrong msgID in an upgrade call errors
	wrongMsgID, err = EncodeUpgradeArgs(targetContract, targets)
	require.NoError(t, err)
	logicCallMsgIDBytes, err := abi.Arguments{{Type: Uint256}}.Pack(LogicCallMsgID)
	require.NoError(t, err)
	wrongMsgID = bytes.Join([][]byte{logicCallMsgIDBytes, wrongMsgID[len(logicCallMsgIDBytes):]}, []byte{})
	_, _, err = DecodeUpgradeArgs(wrongMsgID)
	require.Error(t, err)
	require.Equal(t, err.Error(), "invalid upgrade args")
}
