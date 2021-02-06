package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/x/oracle/types"
)

// DefaultOracleHandler is the default oracle handler for the uniswap oracle data type
// used on sommelier chain.
func (k Keeper) DefaultOracleHandler() types.OracleHandler {
	return func(ctx sdk.Context, oracleDataInputs []types.OracleData) error {

		return nil
	}
}
