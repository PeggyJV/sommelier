package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisValidate(t *testing.T) {
	delAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	valAddr := sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address())

	testCases := []struct {
		name     string
		genState GenesisState
		expPass  bool
	}{
		{
			name:     "default",
			genState: DefaultGenesisState(),
			expPass:  true,
		},
		{
			name: "invalid feeder delegator",
			genState: GenesisState{
				Params: DefaultParams(),
				FeederDelegations: []MsgDelegateAllocations{
					{
						Delegate:  "",
						Validator: valAddr.String(),
					},
				},
			},
			expPass: false,
		},
		{
			name: "invalid feeder validator",
			genState: GenesisState{
				Params: DefaultParams(),
				FeederDelegations: []MsgDelegateAllocations{
					{
						Delegate:  delAddr.String(),
						Validator: "",
					},
				},
			},
			expPass: false,
		},
		{
			name: "equal feeder addresses",
			genState: GenesisState{
				Params: DefaultParams(),
				FeederDelegations: []MsgDelegateAllocations{
					{
						Delegate:  delAddr.String(),
						Validator: sdk.ValAddress(delAddr).String(),
					},
				},
			},
			expPass: false,
		},
		{
			name: "dup feeder delegation",
			genState: GenesisState{
				Params: DefaultParams(),
				FeederDelegations: []MsgDelegateAllocations{
					{
						Delegate:  delAddr.String(),
						Validator: valAddr.String(),
					},
					{
						Delegate:  delAddr.String(),
						Validator: valAddr.String(),
					},
				},
			},
			expPass: false,
		},
	}

	for _, tc := range testCases {

		err := tc.genState.Validate()
		if tc.expPass {
			require.NoError(t, err, tc.name)
		} else {
			require.Error(t, err, tc.name)
		}
	}
}
