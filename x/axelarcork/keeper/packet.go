package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
)

func (k Keeper) ValidateAxelarPacket(ctx sdk.Context, packet ibcexported.PacketI) error {
	params := k.GetParamSet(ctx)
	if !params.Enabled {
		return nil
	}

	// check if this is a call to axelar, exit early if this isn't axelar
	channelID := packet.GetDestChannel()
	if channelID != params.IbcChannel {
		return nil
	}

	// Parse the data from the packet
	var packetData transfertypes.FungibleTokenPacketData
	packetDataBz := packet.GetData()
	if err := json.Unmarshal(packetDataBz, &packetData); err != nil {
		return err
	}

	// if we are not sending to the axelar gmp management account, we can skip
	if packetData.Receiver != params.GmpAccount {
		return nil
	}

	// if the memo field is empty, we can pass the message along
	if packetData.Memo == "" {
		return nil
	}

	var axelarBody types.AxelarBody
	if err := json.Unmarshal([]byte(packetData.Memo), &axelarBody); err != nil {
		return err
	}

	// get the destination chain configuration
	chainConfig, ok := k.GetChainConfigurationByName(ctx, axelarBody.DestinationChain)
	if !ok {
		return fmt.Errorf("configuration not found for chain %s", axelarBody.DestinationChain)
	}

	if chainConfig.ProxyAddress != axelarBody.DestinationAddress {
		return fmt.Errorf("msg cannot bypass the proxy. expected addr %s, received %s", chainConfig.ProxyAddress, axelarBody.DestinationAddress)
	}

	// Validate logic call
	if targetContract, nonce, _, callData, err := types.DecodeLogicCallArgs(axelarBody.Payload); err == nil {
		if nonce == 0 {
			return fmt.Errorf("nonce cannot be zero")
		}

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
	if newProxyContract, targets, err := types.DecodeUpgradeArgs(axelarBody.Payload); err == nil {
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

		// all checks have passed, delete the upgrade data from state
		k.DeleteAxelarProxyUpgradeData(ctx, chainConfig.Id)

		return nil
	}

	return fmt.Errorf("invalid payload: %s", axelarBody.Payload)
}
