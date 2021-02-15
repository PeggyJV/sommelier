package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/peggyjv/sommelier/x/oracle/types"
)

// DefaultOracleHandler is the default oracle handler for the uniswap oracle data type
// used on sommelier chain.
// CONTRACT: this function assumes all the data from the slice share the same type.
func (k Keeper) DefaultOracleHandler() types.OracleHandler {
	return func(_ sdk.Context, oracleDataInputs []types.OracleData) (types.OracleData, error) {
		var (
			aggregatedData types.OracleData
			err            error
		)

		if len(oracleDataInputs) == 0 {
			return nil, nil
		}

		// NOTE: we only check the first element as the rest should be of the same type
		switch oracleData := oracleDataInputs[0].(type) {
		case *types.UniswapPair:
			aggregatedData, err = types.UniswapDataHandler(oracleDataInputs)
		default:
			return nil, sdkerrors.Wrapf(types.ErrInvalidOracleData, "unsupported data type %s", oracleData)
		}

		if err != nil {
			return nil, err
		}

		return aggregatedData, nil
	}
}
