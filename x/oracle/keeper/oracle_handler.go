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
	return func(ctx sdk.Context, oracleDataInputs []types.OracleData) (types.OracleData, error) {
		var (
			aggregatedData types.OracleData
			err            error
		)

		if len(oracleDataInputs) == 0 {
			return nil, nil
		}

		switch oracleData := oracleDataInputs[0].(type) {
		case *types.UniswapPair:
			aggregatedData, err = UniswapDataHandler(oracleDataInputs)
		default:
			return nil, sdkerrors.Wrapf(types.ErrInvalidOracleData, "unsupported data type %s", oracleData)
		}

		if err != nil {
			return nil, err
		}

		return aggregatedData, nil
	}
}

// UniswapDataHandler averages a collection of uniswap pairs oracle data
func UniswapDataHandler(oracleDataInputs []types.OracleData) (types.OracleData, error) {
	var uniswapDataAggregated *types.UniswapPair

	for i, od := range oracleDataInputs {
		up, ok := od.(*types.UniswapPair)
		if !ok {
			return nil, sdkerrors.Wrapf(types.ErrInvalidOracleData, "invalid oracle data %T at index %d", od, i)
		}

		// set up the fixed fields
		if i == 0 {
			uniswapDataAggregated = &types.UniswapPair{
				Id:     up.Id,
				Token0: up.Token0,
				Token1: up.Token1,
			}
		}

		// TODO: add pair data
	}

	// TODO: division

	return uniswapDataAggregated, nil
}
