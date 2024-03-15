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
)

func TestNewMsgAddAddressMappingFormatting(t *testing.T) {
	expectedMsg := &MsgAddAddressMapping{
		EvmAddress: evmAddress1,
		Signer:     cosmosAddress1,
	}

	cosmosAccount1, err := sdk.AccAddressFromBech32(cosmosAddress1)
	require.NoError(t, err)
	createdMsg, err := NewMsgAddAddressMapping(common.HexToAddress(evmAddress1), cosmosAccount1)
	require.Nil(t, err)
	require.Equal(t, expectedMsg, createdMsg)
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
	}

	for _, tc := range testCases {
		err := tc.msgAddAddressMapping.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
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
	require.Nil(t, err)
	require.Equal(t, expectedMsg, createdMsg)
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
	}

	for _, tc := range testCases {
		err := tc.msgRemoveAddressMapping.ValidateBasic()
		if tc.expPass {
			require.NoError(t, err, tc.name)
			require.Nil(t, err)
		} else {
			require.Error(t, err, tc.name)
			require.Equal(t, tc.err.Error(), err.Error())
		}
	}
}
