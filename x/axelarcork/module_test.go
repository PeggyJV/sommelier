package axelarcork_test

import (
	"encoding/json"
	"fmt"
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/golang/mock/gomock"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/tests"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
	"github.com/stretchr/testify/require"
)

var (
	testDenom  = "usomm"
	testAmount = "100"

	testSourcePort         = "sommelier"
	testSourceChannel      = "channel-1"
	testDestinationPort    = "axelar"
	testDestinationChannel = "channel-2"
)

func transferPacket(t *testing.T, receiver string, metadata any) channeltypes.Packet {
	t.Helper()
	transferPacket := transfertypes.FungibleTokenPacketData{
		Denom:    testDenom,
		Amount:   testAmount,
		Receiver: receiver,
	}

	if metadata != nil {
		if mStr, ok := metadata.(string); ok {
			transferPacket.Memo = mStr
		} else {
			memo, err := json.Marshal(metadata)
			require.NoError(t, err)
			transferPacket.Memo = string(memo)
		}
	}

	transferData, err := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
	require.NoError(t, err)

	return channeltypes.Packet{
		SourcePort:         testSourcePort,
		SourceChannel:      testSourceChannel,
		DestinationPort:    testDestinationPort,
		DestinationChannel: testDestinationChannel,
		Data:               transferData,
	}
}

func TestSendPacket_NoMemo(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	setup := tests.NewTestSetup(t, ctl)
	ctx := setup.Initializer.Ctx
	acMiddleware := setup.AxelarCorkMiddleware

	// Test data
	packet := transferPacket(t, tests.TestGMPAccount.String(), "")

	// Expected mocks
	gomock.InOrder(
		setup.Mocks.ICS4WrapperMock.EXPECT().SendPacket(ctx, nil, packet.SourcePort, packet.SourceChannel, packet.TimeoutHeight, packet.TimeoutTimestamp, packet.Data).
			Return(uint64(1), nil),
	)

	_, err := acMiddleware.SendPacket(ctx, nil, packet.SourcePort, packet.SourceChannel, packet.TimeoutHeight, packet.TimeoutTimestamp, packet.Data)
	require.NoError(t, err)
}

func TestSendPacket_NotAxelarChannel(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	setup := tests.NewTestSetup(t, ctl)
	ctx := setup.Initializer.Ctx
	acMiddleware := setup.AxelarCorkMiddleware

	// Test data
	packet := transferPacket(t, tests.TestGMPAccount.String(), "{}")
	packet.DestinationChannel = "channel-other"

	// Expected mocks
	gomock.InOrder(
		setup.Mocks.ICS4WrapperMock.EXPECT().SendPacket(ctx, nil, packet.SourcePort, packet.SourceChannel, packet.TimeoutHeight, packet.TimeoutTimestamp, packet.Data).
			Return(uint64(1), nil),
	)

	_, err := acMiddleware.SendPacket(ctx, nil, packet.SourcePort, packet.SourceChannel, packet.TimeoutHeight, packet.TimeoutTimestamp, packet.Data)
	require.NoError(t, err)
}

func TestSendPacket_NotGMPReceiver(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	setup := tests.NewTestSetup(t, ctl)
	ctx := setup.Initializer.Ctx
	acMiddleware := setup.AxelarCorkMiddleware

	// Test data
	packet := transferPacket(t, "cosmos16plylpsgxechajltx9yeseqexzdzut9g8vla4k", "{}")

	// Expected mocks
	gomock.InOrder(
		setup.Mocks.ICS4WrapperMock.EXPECT().SendPacket(ctx, nil, packet.SourcePort, packet.SourceChannel, packet.TimeoutHeight, packet.TimeoutTimestamp, packet.Data).
			Return(uint64(1), nil),
	)

	_, err := acMiddleware.SendPacket(ctx, nil, packet.SourcePort, packet.SourceChannel, packet.TimeoutHeight, packet.TimeoutTimestamp, packet.Data)
	require.NoError(t, err)
}

func TestSendPacket_EmptyPayload(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	setup := tests.NewTestSetup(t, ctl)
	ctx := setup.Initializer.Ctx
	acMiddleware := setup.AxelarCorkMiddleware

	params := types.Params{
		Enabled:         true,
		IbcChannel:      testDestinationChannel,
		IbcPort:         testDestinationPort,
		GmpAccount:      tests.TestGMPAccount.String(),
		ExecutorAccount: "abc123",
		TimeoutDuration: 10,
	}
	setup.Keepers.AxelarCorkKeeper.SetParams(ctx, params)

	ethChainConfig := types.ChainConfiguration{
		Name:         "Ethereum",
		Id:           1,
		ProxyAddress: "test-proxy-addr",
	}
	setup.Keepers.AxelarCorkKeeper.SetChainConfiguration(ctx, 1, ethChainConfig)

	// Test data
	acBody := types.AxelarBody{
		DestinationChain:   "Ethereum",
		DestinationAddress: "test-addr",
		Payload:            nil,
		Type:               0,
		Fee:                nil,
	}
	packet := transferPacket(t, tests.TestGMPAccount.String(), acBody)

	// Expected mocks
	gomock.InOrder(
		setup.Mocks.ICS4WrapperMock.EXPECT().SendPacket(ctx, nil, packet.SourcePort, packet.SourceChannel, packet.TimeoutHeight, packet.TimeoutTimestamp, packet.Data).
			Return(uint64(0), fmt.Errorf("mock error")),
	)

	// expect error for non-existent
	_, err := acMiddleware.SendPacket(ctx, nil, packet.SourcePort, packet.SourceChannel, packet.TimeoutHeight, packet.TimeoutTimestamp, packet.Data)
	require.Error(t, err)
}
