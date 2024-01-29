package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
)

func (k Keeper) ValidateAxelarPacket(ctx sdk.Context, sourceChannel string, data []byte) error {
	params := k.GetParamSet(ctx)
	if !params.Enabled {
		return nil
	}

	// check if this is a call to axelar, exit early if this isn't axelar
	if sourceChannel != params.IbcChannel {
		return nil
	}

	// Parse the data from the packet
	var packetData transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(data, &packetData); err != nil {
		return err
	}

	// decoding some bech32 strings so our comparisons are guaranteed to be accurate
	gmpAccountAddr, err := sdk.AccAddressFromBech32(params.GmpAccount)
	if err != nil {
		return fmt.Errorf("GmpAccount parameter is an invalid address: %s", params.GmpAccount)
	}

	receiverAddr, err := sdk.AccAddressFromBech32(packetData.Receiver)
	if err != nil {
		return fmt.Errorf("receiver in IBC packet data is an invalid address: %s", packetData.Receiver)
	}

	senderAddr, err := sdk.AccAddressFromBech32(packetData.Sender)
	if err != nil {
		return fmt.Errorf("sender in IBC packet data is an invalid address: %s", packetData.Sender)
	}

	// if we are not sending to the axelar gmp management account, we can skip
	if !receiverAddr.Equals(gmpAccountAddr) {
		return nil
	}

	// reject if the axelar module account is not the sender
	if !senderAddr.Equals(k.GetSenderAccount(ctx).GetAddress()) {
		return fmt.Errorf("sender to Axelar GMP account is not axelarcork module account: %s", packetData.Sender)
	}

	// if the memo field is empty, we can pass the message along
	if packetData.Memo == "" {
		return nil
	}

	var axelarBody types.AxelarBody
	if err := json.Unmarshal([]byte(packetData.Memo), &axelarBody); err != nil {
		return err
	}
	payloadBytes := intsToBytes(axelarBody.Payload)

	// get the destination chain configuration
	chainConfig, ok := k.GetChainConfigurationByName(ctx, axelarBody.DestinationChain)
	if !ok {
		return fmt.Errorf("configuration not found for chain %s", axelarBody.DestinationChain)
	}

	// decoding some EVM addresses here so our comparisons are guaranteed to be accurate
	if !common.IsHexAddress(chainConfig.ProxyAddress) {
		return fmt.Errorf("proxy address in chain config is not valid, chain ID %d, address %s", chainConfig.Id, chainConfig.ProxyAddress)
	}
	proxyAddr := common.HexToAddress(chainConfig.ProxyAddress)

	if !common.IsHexAddress(axelarBody.DestinationAddress) {
		return fmt.Errorf("axelar destination address is not a valid EVM address: %s", axelarBody.DestinationAddress)
	}
	axelarDestinationAddr := common.HexToAddress(axelarBody.DestinationAddress)

	if !bytes.Equal(axelarDestinationAddr.Bytes(), proxyAddr.Bytes()) {
		return fmt.Errorf("msg cannot bypass the proxy. expected addr %s, received %s", chainConfig.ProxyAddress, axelarBody.DestinationAddress)
	}

	// Validate logic call
	if targetContract, nonce, _, callData, err := types.DecodeLogicCallArgs(payloadBytes); err == nil {
		if nonce == 0 {
			return fmt.Errorf("nonce cannot be zero")
		}

		// TODO(bolten): is there any validation on the deadline worth doing?

		blockHeight, winningCork, ok := k.GetWinningAxelarCork(ctx, chainConfig.Id, common.HexToAddress(targetContract))
		if !ok {
			return fmt.Errorf("no cork expected for chain %s:%d at address %s", chainConfig.Name, chainConfig.Id, axelarBody.DestinationAddress)
		}

		if !bytes.Equal(winningCork.EncodedContractCall, callData) {
			return fmt.Errorf("cork body did not match expected body. received: %x, expected: %x", callData, winningCork.EncodedContractCall)
		}

		// all checks have passed, delete the cork from state
		k.DeleteWinningAxelarCorkByBlockheight(ctx, chainConfig.Id, blockHeight, winningCork)

		return nil
	}

	// Validate upgrade
	if newProxyContract, targets, err := types.DecodeUpgradeArgs(payloadBytes); err == nil {
		if !common.IsHexAddress(newProxyContract) {
			return fmt.Errorf("invalid proxy address %s", newProxyContract)
		}

		if len(targets) == 0 {
			return fmt.Errorf("no targets provided")
		}

		for _, target := range targets {
			if !common.IsHexAddress(target) {
				return fmt.Errorf("invalid target %s", target)
			}
		}

		upgradeData, ok := k.GetAxelarProxyUpgradeData(ctx, chainConfig.Id)
		if !ok {
			return fmt.Errorf("no upgrade data expected for chain %s:%d", chainConfig.Name, chainConfig.Id)
		}

		if !bytes.Equal(upgradeData.Payload, payloadBytes) {
			return fmt.Errorf("upgrade data did not match expected data. received: %s, expected: %s", payloadBytes, upgradeData.Payload)
		}

		// all checks have passed, delete the upgrade data from state
		k.DeleteAxelarProxyUpgradeData(ctx, chainConfig.Id)

		return nil
	}

	return fmt.Errorf("invalid payload: %s", payloadBytes)
}

func intsToBytes(payload []int) []byte {
	bytePayload := make([]byte, len(payload))
	for i, b := range payload {
		bytePayload[i] = byte(b)
	}
	return bytePayload
}
