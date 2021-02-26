package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// OracleHandler defines a type that is passed to the oracle keeper to archive custom handling of
// oracle data processing. It returns the aggregated data and an error.
type OracleHandler func(ctx sdk.Context, oracleDataInput []OracleData) (OracleData, error)

// UniswapDataHandler averages a collection of uniswap pairs oracle data.
// CONTRACT: input length must be > 0.
func UniswapDataHandler(oracleDataInputs []OracleData) (OracleData, error) {
	var uniswapDataAggregated *UniswapPair

	if len(oracleDataInputs) == 0 {
		return nil, nil
	}

	for i, od := range oracleDataInputs {
		up, ok := od.(*UniswapPair)
		if !ok {
			return nil, sdkerrors.Wrapf(ErrInvalidOracleData, "invalid oracle data %T at index %d", od, i)
		}

		// set up the fixed fields and zero out the
		if i == 0 {
			uniswapDataAggregated = NewUniswapPair(up.ID, up.Token0, up.Token1)
		}

		uniswapDataAggregated.Reserve0 = uniswapDataAggregated.Reserve0.Add(up.Reserve0)
		uniswapDataAggregated.Reserve1 = uniswapDataAggregated.Reserve1.Add(up.Reserve1)
		uniswapDataAggregated.ReserveUSD = uniswapDataAggregated.ReserveUSD.Add(up.ReserveUSD)
		uniswapDataAggregated.Token0Price = uniswapDataAggregated.Token0Price.Add(up.Token0Price)
		uniswapDataAggregated.Token1Price = uniswapDataAggregated.Token1Price.Add(up.Token1Price)
		uniswapDataAggregated.TotalSupply = uniswapDataAggregated.TotalSupply.Add(up.TotalSupply)
	}

	inputs := sdk.NewDec(int64(len(oracleDataInputs)))

	// division by the number of inputs
	uniswapDataAggregated.Reserve0 = uniswapDataAggregated.Reserve0.Quo(inputs)
	uniswapDataAggregated.Reserve1 = uniswapDataAggregated.Reserve1.Quo(inputs)
	uniswapDataAggregated.ReserveUSD = uniswapDataAggregated.ReserveUSD.Quo(inputs)
	uniswapDataAggregated.Token0Price = uniswapDataAggregated.Token0Price.Quo(inputs)
	uniswapDataAggregated.Token1Price = uniswapDataAggregated.Token1Price.Quo(inputs)
	uniswapDataAggregated.TotalSupply = uniswapDataAggregated.TotalSupply.Quo(inputs)

	return uniswapDataAggregated, nil
}
