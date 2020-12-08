package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/peggyjv/sommelier/x/oracle/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryExchangeRate:
			return queryExchangeRate(ctx, req, keeper)
		case types.QueryExchangeRates:
			return queryExchangeRates(ctx, keeper)
		case types.QueryActives:
			return queryActives(ctx, keeper)
		case types.QueryParameters:
			return queryParameters(ctx, keeper)
		case types.QueryFeederDelegation:
			return queryFeederDelegation(ctx, req, keeper)
		case types.QueryMissCounter:
			return queryMissCounter(ctx, req, keeper)
		case types.QueryAggregatePrevote:
			return queryAggregatePrevote(ctx, req, keeper)
		case types.QueryAggregateVote:
			return queryAggregateVote(ctx, req, keeper)
		case types.QueryVoteTargets:
			return queryVoteTargets(ctx, keeper)
		case types.QueryTobinTax:
			return queryTobinTax(ctx, req, keeper)
		case types.QueryTobinTaxes:
			return queryTobinTaxes(ctx, keeper)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func queryExchangeRate(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryExchangeRateParams
	if err := json.Unmarshal(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	rate, err := keeper.GetLunaExchangeRate(ctx, params.Denom)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnknowDenom, params.Denom)
	}
	return []byte(rate.String()), nil
}

func queryExchangeRates(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var rates sdk.DecCoins
	keeper.IterateLunaExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		rates = append(rates, sdk.NewDecCoinFromDec(denom, rate))
		return false
	})
	return []byte(rates.String()), nil
}

func queryActives(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var denoms []string
	keeper.IterateLunaExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
		denoms = append(denoms, denom)
		return false
	})
	bz, err := json.MarshalIndent(denoms, "", "  ")
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

// TODO: investigate param returns
func queryParameters(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	bz, err := json.Marshal(keeper.GetParams(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryFeederDelegation(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryFeederDelegationParams
	if err := json.Unmarshal(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []byte(keeper.GetOracleDelegate(ctx, params.Validator).String()), nil
}

func queryMissCounter(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryMissCounterParams
	if err := json.Unmarshal(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []byte(fmt.Sprint(keeper.GetMissCounter(ctx, params.Validator))), nil
}

func queryAggregatePrevote(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryAggregatePrevoteParams
	if err := json.Unmarshal(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	aggregateExchangeRatePrevote, err := keeper.GetAggregateExchangeRatePrevote(ctx, params.Validator)
	if err != nil {
		return nil, err
	}
	bz, err := json.MarshalIndent(aggregateExchangeRatePrevote, "", "  ")
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

// TODO: make this print { "voter": "cosmos1u94xef3cp9thkcpxecuvhtpwnmg8mhlja8hzkd", "rates": "92992.1103umon/foo,3014.893umosn/foosd" }
func queryAggregateVote(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryAggregateVoteParams
	if err := json.Unmarshal(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	vote, err := keeper.GetAggregateExchangeRateVote(ctx, params.Validator)
	if err != nil {
		return nil, err
	}
	bz, err := json.MarshalIndent(vote, "", "  ")
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryVoteTargets(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	bz, err := json.MarshalIndent(keeper.GetVoteTargets(ctx), "", "  ")
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryTobinTax(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, error) {
	var params types.QueryTobinTaxParams
	if err := json.Unmarshal(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	tobinTax, err := keeper.GetTobinTax(ctx, params.Denom)
	if err != nil {
		return nil, err
	}
	return []byte(tobinTax.String()), nil
}

func queryTobinTaxes(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	var denoms sdk.DecCoins
	keeper.IterateTobinTaxes(ctx, func(denom string, tobinTax sdk.Dec) (stop bool) {
		denoms = append(denoms, sdk.NewDecCoinFromDec(denom, tobinTax))
		return false
	})
	return []byte(denoms.String()), nil
}
