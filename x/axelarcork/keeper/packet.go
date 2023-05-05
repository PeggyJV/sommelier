package keeper

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v3/modules/core/exported"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
)

// SendPacket wraps IBC ChannelKeeper's SendPacket function
// If the packet does not get blocked by validation, it passes to the IBC Channel keeper
func (k Keeper) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet ibcexported.PacketI) error {
	if err := k.ValidateAxelarCorkPacket(ctx, packet); err != nil {
		k.Logger(ctx).Error(fmt.Sprintf("ICS20 packet send was denied: %s", err.Error()))
		return err
	}
	return k.ics4Wrapper.SendPacket(ctx, chanCap, packet)
}

func (k Keeper) ValidateAxelarCorkPacket(ctx sdk.Context, packet ibcexported.PacketI) error {
	// check if this is a call to axelar, exit early if this isn't axelar
	channelID := packet.GetDestChannel()
	if channelID != k.GetParamSet(ctx).IbcChannel {
		return nil
	}

	// Parse the amount and denom from the packet
	var packetData transfertypes.FungibleTokenPacketData
	if err := json.Unmarshal(packet.GetData(), &packetData); err != nil {
		return err
	}

	// if we are not sending to the axelar gmp management account, we can skip
	if packetData.Receiver != k.GetParamSet(ctx).GmpAccount {
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

	var wrappedBody types.ProxyWrapper
	if err := json.Unmarshal(axelarBody.Payload, &wrappedBody); err != nil {
		return err
	}

	winningCork, ok := k.GetWinningCork(ctx, chainConfig.Id, common.HexToAddress(wrappedBody.Target))
	if !ok {
		return fmt.Errorf("no cork expected for chain %s:%d at address %s", chainConfig.Name, chainConfig.Id, axelarBody.DestinationAddress)
	}

	if !bytes.Equal(winningCork.EncodedContractCall, wrappedBody.Body) {
		return fmt.Errorf("cork body did not match expected body. received: %x, expected: %x", axelarBody.Payload, winningCork.EncodedContractCall)
	}

	// all checks have passed, delete the cork from state
	k.DeleteWinningCork(ctx, chainConfig.Id, winningCork)

	return nil
}

// WriteAcknowledgement wraps IBC ChannelKeeper's WriteAcknowledgement function
func (k Keeper) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet ibcexported.PacketI, acknowledgement ibcexported.Acknowledgement) error {
	return k.ics4Wrapper.WriteAcknowledgement(ctx, chanCap, packet, acknowledgement)
}
