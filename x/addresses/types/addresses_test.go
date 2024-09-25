package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddressMapping_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		mapping AddressMapping
		wantErr bool
	}{
		{
			name: "valid mapping",
			mapping: AddressMapping{
				CosmosAddress: "cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu",
				EvmAddress:    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			},
			wantErr: false,
		},
		{
			name: "invalid cosmos address",
			mapping: AddressMapping{
				CosmosAddress: "invalid",
				EvmAddress:    "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
			},
			wantErr: true,
		},
		{
			name: "invalid evm address",
			mapping: AddressMapping{
				CosmosAddress: "cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu",
				EvmAddress:    "invalid",
			},
			wantErr: true,
		},
		{
			name: "both addresses invalid",
			mapping: AddressMapping{
				CosmosAddress: "invalid",
				EvmAddress:    "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.mapping.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
