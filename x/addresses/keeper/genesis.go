package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/sommelier/v7/x/addresses/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k Keeper, gs types.GenesisState) {
	if err := gs.Validate(); err != nil {
		panic(fmt.Sprintf("failed to validate %s genesis state: %s", types.ModuleName, err))
	}

	k.setParams(ctx, gs.Params)

	for _, mapping := range gs.AddressMappings {
		cosmosAcc, err := sdk.AccAddressFromBech32(mapping.CosmosAddress)
		if err != nil {
			panic(err)
		}

		if !common.IsHexAddress(mapping.EvmAddress) {
			panic(fmt.Sprintf("invalid EVM address %s", mapping.EvmAddress))
		}

		evmAddr := common.HexToAddress(mapping.EvmAddress).Bytes()

		k.SetAddressMapping(ctx, cosmosAcc.Bytes(), evmAddr)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	mappings := []*types.AddressMapping{}
	k.IterateAddressMappings(ctx, func(cosmosAddr []byte, evmAddr []byte) bool {
		mappings = append(mappings, &types.AddressMapping{
			CosmosAddress: sdk.AccAddress(cosmosAddr).String(),
			EvmAddress:    common.BytesToAddress(evmAddr).Hex(),
		})

		return false
	})

	return types.GenesisState{
		Params:          k.GetParamSet(ctx),
		AddressMappings: mappings,
	}
}
