package types

import (
	"testing"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	evmAddress1    = "0x1111111111111111111111111111111111111111"
	cosmosAddress1 = "cosmos154d0p9xhrruhxvazumej9nq29afeura2alje4u"
	evmAddress2    = "0x2222222222222222222222222222222222222222"
	cosmosAddress2 = "cosmos1y6d5kasehecexf09ka6y0ggl0pxzt6dgk0gnl9"
)

func TestNewMsgAddAddressMappingFormatting(t *testing.T) {
	expectedMsg := &MsgAddAddressMapping{
		EvmAddress: evmAddress1,
		Signer:     cosmosAddress1,
	}

	cosmosAccount1, err := sdk.AccAddressFromBech32(cosmosAddress1)
	require.NoError(t, err)
	createdMsg, err := NewMsgAddAddressMapping(common.HexToAddress(evmAddress1), cosmosAccount1)
	require.NoError(t, err)
	require.Equal(t, expectedMsg, createdMsg)

	// Test with nil Cosmos address
	_, err = NewMsgAddAddressMapping(common.HexToAddress(evmAddress1), nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid address")
}

func TestMsgAddAddressMappingValidate(t *testing.T) {
	testCases := []struct {
		name                 string
		msgAddAddressMapping MsgAddAddressMapping
		expPass              bool
		err                  error
	}{
		{
			name: "Happy path",
			msgAddAddressMapping: MsgAddAddressMapping{
				EvmAddress: evmAddress1,
				Signer:     cosmosAddress1,
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Invalid signer address",
			msgAddAddressMapping: MsgAddAddressMapping{
				EvmAddress: evmAddress1,
				Signer:     "sladjflaksjfd",
			},
			expPass: false,
			err:     errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "decoding bech32 failed: invalid separator index -1"),
		},
		{
			name: "Invalid EVM address",
			msgAddAddressMapping: MsgAddAddressMapping{
				EvmAddress: "lasjfdlsdf",
				Signer:     cosmosAddress1,
			},
			expPass: false,
			err:     errorsmod.Wrap(ErrInvalidEvmAddress, "lasjfdlsdf is not a valid hex address"),
		},
		{
			name: "Empty EVM address",
			msgAddAddressMapping: MsgAddAddressMapping{
				EvmAddress: "",
				Signer:     cosmosAddress1,
			},
			expPass: false,
			err:     errorsmod.Wrap(ErrInvalidEvmAddress, " is not a valid hex address"),
		},
		{
			name: "Empty signer address",
			msgAddAddressMapping: MsgAddAddressMapping{
				EvmAddress: evmAddress1,
				Signer:     "",
			},
			expPass: false,
			err:     errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "empty address string is not allowed"),
		},
	}

	for _, tc := range testCases {
		err := tc.msgAddAddressMapping.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error(), tc.name)
		}
	}
}

func TestNewMsgRemoveAddressMappingFormatting(t *testing.T) {
	expectedMsg := &MsgRemoveAddressMapping{
		Signer: cosmosAddress1,
	}

	cosmosAccount1, err := sdk.AccAddressFromBech32(cosmosAddress1)
	require.NoError(t, err)
	createdMsg, err := NewMsgRemoveAddressMapping(cosmosAccount1)
	require.NoError(t, err)
	require.Equal(t, expectedMsg, createdMsg)

	// Test with nil Cosmos address
	_, err = NewMsgRemoveAddressMapping(nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid address")
}

func TestMsgRemoveAddressMappingValidate(t *testing.T) {
	testCases := []struct {
		name                    string
		msgRemoveAddressMapping MsgRemoveAddressMapping
		expPass                 bool
		err                     error
	}{
		{
			name: "Happy path",
			msgRemoveAddressMapping: MsgRemoveAddressMapping{
				Signer: cosmosAddress1,
			},
			expPass: true,
			err:     nil,
		},
		{
			name: "Invalid signer address",
			msgRemoveAddressMapping: MsgRemoveAddressMapping{
				Signer: "sladjflaksjfd",
			},
			expPass: false,
			err:     errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "decoding bech32 failed: invalid separator index -1"),
		},
		{
			name: "Empty signer address",
			msgRemoveAddressMapping: MsgRemoveAddressMapping{
				Signer: "",
			},
			expPass: false,
			err:     errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "empty address string is not allowed"),
		},
	}

	for _, tc := range testCases {
		err := tc.msgRemoveAddressMapping.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error(), tc.name)
		}
	}
}

func TestMsgAddAddressMappingGetSigners(t *testing.T) {
	msg := MsgAddAddressMapping{
		EvmAddress: evmAddress1,
		Signer:     cosmosAddress1,
	}

	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, cosmosAddress1, signers[0].String())

	// Test with invalid signer address
	msg.Signer = "invalid_address"
	require.Panics(t, func() { msg.GetSigners() })
}

func TestMsgRemoveAddressMappingGetSigners(t *testing.T) {
	msg := MsgRemoveAddressMapping{
		Signer: cosmosAddress1,
	}

	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, cosmosAddress1, signers[0].String())

	// Test with invalid signer address
	msg.Signer = "invalid_address"
	require.Panics(t, func() { msg.GetSigners() })
}

func TestMsgAddAddressMappingType(t *testing.T) {
	msg := MsgAddAddressMapping{}
	require.Equal(t, "add_address_mapping", msg.Type())
}

func TestMsgRemoveAddressMappingType(t *testing.T) {
	msg := MsgRemoveAddressMapping{}
	require.Equal(t, "remove_address_mapping", msg.Type())
}

func TestMsgAddAddressMappingGetSignBytes(t *testing.T) {
	msg := MsgAddAddressMapping{
		EvmAddress: evmAddress1,
		Signer:     cosmosAddress1,
	}

	signBytes := msg.GetSignBytes()
	require.NotEmpty(t, signBytes)
}

func TestMsgRemoveAddressMappingGetSignBytes(t *testing.T) {
	msg := MsgRemoveAddressMapping{
		Signer: cosmosAddress1,
	}

	signBytes := msg.GetSignBytes()
	require.NotEmpty(t, signBytes)
}
