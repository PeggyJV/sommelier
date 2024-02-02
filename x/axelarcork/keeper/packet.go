package keeper

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
)

func (k Keeper) ValidateAxelarPacket(ctx sdk.Context, sourceChannel string, data []byte) error {
	params := k.GetParamSet(ctx)

	// check if this is a call to axelar, exit early if this isn't axelar
	if sourceChannel != params.IbcChannel {
		return nil
	}

	k.Logger(ctx).Info("checking IBC packet against Axelar middleware validations")

	// Parse the data from the packet
	var packetData transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(data, &packetData); err != nil {
		return err
	}

	// decoding some bech32 strings so our comparisons are guaranteed to be accurate
	gmpAccountAddr, err := sdk.GetFromBech32(params.GmpAccount, "axelar")
	if err != nil {
		return fmt.Errorf("GmpAccount parameter is an invalid address: %s", params.GmpAccount)
	}

	receiverAddr, err := sdk.GetFromBech32(packetData.Receiver, "axelar")
	if err != nil {
		return fmt.Errorf("receiver in IBC packet data is an invalid address: %s", packetData.Receiver)
	}

	senderAddr, err := sdk.AccAddressFromBech32(packetData.Sender)
	if err != nil {
		return fmt.Errorf("sender in IBC packet data is an invalid address: %s", packetData.Sender)
	}

	// if we are not sending to the axelar gmp management account, we can skip
	if !bytes.Equal(receiverAddr, gmpAccountAddr) {
		k.Logger(ctx).Info("Axelar receiver is not the GMP account, allowing packet", "receiver", hex.EncodeToString(receiverAddr), "gmp account", hex.EncodeToString(gmpAccountAddr))
		return nil
	}

	// reject if the axelar module account is not the sender
	if !senderAddr.Equals(k.GetSenderAccount(ctx).GetAddress()) {
		return fmt.Errorf("sender to Axelar GMP account is not axelarcork module account: %s", packetData.Sender)
	}

	// if the memo field is empty, we can pass the message along
	if packetData.Memo == "" {
		k.Logger(ctx).Error("Axelar GMP packet memo is empty")
		return nil
	}

	var axelarBody types.AxelarBody
	if err := json.Unmarshal([]byte(packetData.Memo), &axelarBody); err != nil {
		return err
	}
	payloadBytes := axelarBody.Payload

	// shortcircuit for pure token transfer
	if axelarBody.Type == types.PureTokenTransfer {
		if len(payloadBytes) != 0 {
			return fmt.Errorf("payload must be empty for pure token transfer")
		}

		return nil
	}

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
	if targetContract, nonce, deadline, callData, err := types.DecodeLogicCallArgs(payloadBytes); err == nil {
		if nonce == 0 {
			return fmt.Errorf("nonce cannot be zero")
		}

		if deadline == 0 {
			return fmt.Errorf("deadline cannot be zero")
		}

		blockHeight, winningCork, ok := k.GetWinningAxelarCork(ctx, chainConfig.Id, common.HexToAddress(targetContract))
		if !ok {
			return fmt.Errorf("no cork expected for chain %s:%d at address %s", chainConfig.Name, chainConfig.Id, axelarBody.DestinationAddress)
		}

		if !bytes.Equal(winningCork.EncodedContractCall, callData) {
			return fmt.Errorf("cork body did not match expected body. received: %x, expected: %x", callData, winningCork.EncodedContractCall)
		}

		// all checks have passed, delete the cork from state
		k.Logger(ctx).Info("Axelar GMP message validated, deleting from state", "chain ID", chainConfig.Id, "block height", blockHeight, "contract", winningCork.TargetContractAddress)
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
		k.Logger(ctx).Info("Axelar GMP upgrade message validated, deleting from state", "chain ID", chainConfig.Id)
		k.DeleteAxelarProxyUpgradeData(ctx, chainConfig.Id)

		return nil
	}

	return fmt.Errorf("invalid payload: %s", payloadBytes)
}
