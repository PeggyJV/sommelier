package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	def := DefaultGenesisState()
	testCases := []struct {
		desc     string
		genState *GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: &def,
			valid:    true,
		},
		{
			desc:     "valid genesis state",
			genState: &GenesisState{},
			valid:    true,
		},
		{
			desc: "invalid address mapping - invalid cosmos address",
			genState: &GenesisState{
				AddressMappings: []*AddressMapping{
					{
						CosmosAddress: "invalid_cosmos_address",
						EvmAddress:    "0x1234567890123456789012345678901234567890",
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid address mapping - invalid evm address",
			genState: &GenesisState{
				AddressMappings: []*AddressMapping{
					{
						CosmosAddress: cosmosAddress1,
						EvmAddress:    "invalid_evm_address",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicate address mappings",
			genState: &GenesisState{
				AddressMappings: []*AddressMapping{
					{
						CosmosAddress: cosmosAddress1,
						EvmAddress:    "0x1234567890123456789012345678901234567890",
					},
					{
						CosmosAddress: cosmosAddress1,
						EvmAddress:    "0x1234567890123456789012345678901234567890",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicate cosmos address",
			genState: &GenesisState{
				AddressMappings: []*AddressMapping{
					{
						CosmosAddress: cosmosAddress1,
						EvmAddress:    "0x1234567890123456789012345678901234567890",
					},
					{
						CosmosAddress: cosmosAddress1,
						EvmAddress:    "0x0987654321098765432109876543210987654321",
					},
				},
			},
			valid: true,
		},
	}

	for _, tc := range testCases {
		tc := tc // create a local copy
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestValidateGenesis(t *testing.T) {
	testCases := []struct {
		desc     string
		genState GenesisState
		valid    bool
	}{
		{
			desc:     "default genesis state",
			genState: DefaultGenesisState(),
			valid:    true,
		},
		{
			desc: "custom genesis state",
			genState: GenesisState{
				Params: DefaultParams(),
				AddressMappings: []*AddressMapping{
					{
						CosmosAddress: "cosmos1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqnrql8a",
						EvmAddress:    "0x1234567890123456789012345678901234567890",
					},
				},
			},
			valid: true,
		},
	}

	for _, tc := range testCases {
		tc := tc // create a local copy
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestDefaultGenesisState(t *testing.T) {
	defaultGenesis := DefaultGenesisState()
	require.NotNil(t, defaultGenesis.Params)
	require.Equal(t, DefaultParams(), defaultGenesis.Params)
	require.Empty(t, defaultGenesis.AddressMappings)
}

func TestGenesisState_ValidateAddressMappings(t *testing.T) {
	testCases := []struct {
		desc     string
		mappings []*AddressMapping
		valid    bool
	}{
		{
			desc:     "empty mappings",
			mappings: []*AddressMapping{},
			valid:    true,
		},
		{
			desc: "valid mappings",
			mappings: []*AddressMapping{
				{
					CosmosAddress: cosmosAddress1,
					EvmAddress:    "0x1234567890123456789012345678901234567890",
				},
				{
					CosmosAddress: cosmosAddress2,
					EvmAddress:    "0x0987654321098765432109876543210987654321",
				},
			},
			valid: true,
		},
		{
			desc: "invalid cosmos address",
			mappings: []*AddressMapping{
				{
					CosmosAddress: "invalid_cosmos_address",
					EvmAddress:    "0x1234567890123456789012345678901234567890",
				},
			},
			valid: false,
		},
		{
			desc: "invalid evm address",
			mappings: []*AddressMapping{
				{
					CosmosAddress: cosmosAddress1,
					EvmAddress:    "invalid_evm_address",
				},
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		tc := tc // create a local copy
		t.Run(tc.desc, func(t *testing.T) {
			genState := GenesisState{
				Params:          DefaultParams(),
				AddressMappings: tc.mappings,
			}
			err := genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
