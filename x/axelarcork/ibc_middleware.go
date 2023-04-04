package axelarcork

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	porttypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
)

var _ porttypes.Middleware = &IBCMiddleware{}

type IBCMiddleware struct {
	app    porttypes.IBCModule
	keeper keeper.Keeper
}

func NewIBCMiddleware(k keeper.Keeper, app porttypes.IBCModule) IBCMiddleware {
	return IBCMiddleware{
		app:    app,
		keeper: k,
	}
}

func (im IBCMiddleware) OnChanOpenInit(ctx sdk.Context, order types.Order, connectionHops []string, portID string, channelID string, channelCap *capabilitytypes.Capability, counterparty types.Counterparty, version string) error {
	return im.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, version)
}

func (im IBCMiddleware) OnChanOpenTry(ctx sdk.Context, order types.Order, connectionHops []string, portID, channelID string, channelCap *capabilitytypes.Capability, counterparty types.Counterparty, counterpartyVersion string) (version string, err error) {
	return im.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, channelCap, counterparty, counterpartyVersion)
}

func (im IBCMiddleware) OnChanOpenAck(ctx sdk.Context, portID, channelID string, counterpartyChannelID string, counterpartyVersion string) error {
	return im.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion)
}

func (im IBCMiddleware) OnChanOpenConfirm(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanOpenConfirm(ctx, portID, channelID)
}

func (im IBCMiddleware) OnChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanCloseInit(ctx, portID, channelID)
}

func (im IBCMiddleware) OnChanCloseConfirm(ctx sdk.Context, portID, channelID string) error {
	return im.app.OnChanCloseConfirm(ctx, portID, channelID)
}

func (im IBCMiddleware) OnRecvPacket(ctx sdk.Context, packet types.Packet, relayer sdk.AccAddress) exported.Acknowledgement {
	return im.app.OnRecvPacket(ctx, packet, relayer)
}

func (im IBCMiddleware) OnAcknowledgementPacket(ctx sdk.Context, packet types.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
	return im.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
}

func (im IBCMiddleware) OnTimeoutPacket(ctx sdk.Context, packet types.Packet, relayer sdk.AccAddress) error {
	return im.app.OnTimeoutPacket(ctx, packet, relayer)
}

func (im IBCMiddleware) SendPacket(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI) error {
	//TODO implement me
	panic("implement me")
}

func (im IBCMiddleware) WriteAcknowledgement(ctx sdk.Context, chanCap *capabilitytypes.Capability, packet exported.PacketI, ack exported.Acknowledgement) error {
	//TODO implement me
	panic("implement me")
}
