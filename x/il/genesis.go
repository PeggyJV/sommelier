package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/x/il/keeper"
	"github.com/peggyjv/sommelier/x/il/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {

}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	return types.GenesisState{}
}
