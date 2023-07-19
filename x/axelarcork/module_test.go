package axelarcork_test

import (
	"encoding/json"
	"testing"

	transfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v3/modules/core/04-channel/types"
	"github.com/golang/mock/gomock"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/tests"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
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
		setup.Mocks.ICS4WrapperMock.EXPECT().SendPacket(ctx, nil, packet).
			Return(nil),
	)

	require.NoError(t, acMiddleware.SendPacket(ctx, nil, packet))
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
		setup.Mocks.ICS4WrapperMock.EXPECT().SendPacket(ctx, nil, packet).
			Return(nil),
	)

	require.NoError(t, acMiddleware.SendPacket(ctx, nil, packet))
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
		setup.Mocks.ICS4WrapperMock.EXPECT().SendPacket(ctx, nil, packet).
			Return(nil),
	)

	require.NoError(t, acMiddleware.SendPacket(ctx, nil, packet))
}

func TestSendPacket_EmptyPayload(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	setup := tests.NewTestSetup(t, ctl)
	ctx := setup.Initializer.Ctx
	acMiddleware := setup.AxelarCorkMiddleware

	// Test data
	acBody := types.AxelarBody{
		DestinationChain:   "sommelier",
		DestinationAddress: "test-addr",
		Payload:            nil,
		Type:               0,
		Fee:                nil,
	}
	packet := transferPacket(t, tests.TestGMPAccount.String(), acBody)

	// expect error for non-existent
	require.Error(t, acMiddleware.SendPacket(ctx, nil, packet))
}
