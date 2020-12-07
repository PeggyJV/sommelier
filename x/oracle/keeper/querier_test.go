package keeper

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

func TestNewQuerier(t *testing.T) {
	input := CreateTestInput(t)

	querier := NewQuerier(input.OracleKeeper)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(input.Ctx, []string{types.QueryParameters}, query)
	require.NoError(t, err)
}

func TestQueryParams(t *testing.T) {
	input := CreateTestInput(t)

	var params types.Params

	res, errRes := queryParameters(input.Ctx, input.OracleKeeper)
	require.NoError(t, errRes)

	err := json.Unmarshal(res, &params)
	require.NoError(t, err)
	require.Equal(t, input.OracleKeeper.GetParams(input.Ctx), params)
}

func TestQueryExchangeRate(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	rate := sdk.NewDec(1700)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, types.MicroSDRDenom, rate)

	// denom query params
	queryParams := types.NewQueryExchangeRateParams(types.MicroSDRDenom)
	bz, err := cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryExchangeRate}, req)
	require.NoError(t, err)

	var rrate sdk.Dec
	err = cdc.UnmarshalJSON(res, &rrate)
	require.NoError(t, err)
	require.Equal(t, rate, rrate)
}

func TestQueryExchangeRates(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	rate := sdk.NewDec(1700)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, types.MicroSDRDenom, rate)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, types.MicroUSDDenom, rate)

	res, err := querier(input.Ctx, []string{types.QueryExchangeRates}, abci.RequestQuery{})
	require.NoError(t, err)

	var rrate sdk.DecCoins
	err2 := cdc.UnmarshalJSON(res, &rrate)
	require.NoError(t, err2)
	require.Equal(t, sdk.DecCoins{
		sdk.NewDecCoinFromDec(types.MicroSDRDenom, rate),
		sdk.NewDecCoinFromDec(types.MicroUSDDenom, rate),
	}, rrate)
}

func TestQueryActives(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	rate := sdk.NewDec(1700)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, types.MicroSDRDenom, rate)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, types.MicroKRWDenom, rate)
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, types.MicroUSDDenom, rate)

	res, err := querier(input.Ctx, []string{types.QueryActives}, abci.RequestQuery{})
	require.NoError(t, err)

	targetDenoms := []string{
		types.MicroKRWDenom,
		types.MicroSDRDenom,
		types.MicroUSDDenom,
	}

	var denoms []string
	err2 := cdc.UnmarshalJSON(res, &denoms)
	require.NoError(t, err2)
	require.Equal(t, targetDenoms, denoms)
}

func TestQueryFeederDelegation(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	input.OracleKeeper.SetOracleDelegate(input.Ctx, ValAddrs[0], Addrs[1])

	queryParams := types.NewQueryFeederDelegationParams(ValAddrs[0])
	bz, err := cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryFeederDelegation}, req)
	require.NoError(t, err)

	var delegate sdk.AccAddress
	cdc.UnmarshalJSON(res, &delegate)
	require.Equal(t, Addrs[1], delegate)
}

func TestQueryAggregatePrevote(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	prevote1 := types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash{}, ValAddrs[0], 0)
	input.OracleKeeper.AddAggregateExchangeRatePrevote(input.Ctx, prevote1)
	prevote2 := types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash{}, ValAddrs[1], 0)
	input.OracleKeeper.AddAggregateExchangeRatePrevote(input.Ctx, prevote2)
	prevote3 := types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash{}, ValAddrs[2], 0)
	input.OracleKeeper.AddAggregateExchangeRatePrevote(input.Ctx, prevote3)

	// validator 0 address params
	queryParams := types.NewQueryAggregatePrevoteParams(ValAddrs[0])
	bz, err := cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryAggregatePrevote}, req)
	require.NoError(t, err)

	var prevote types.AggregateExchangeRatePrevote
	err = json.Unmarshal(res, &prevote)
	require.NoError(t, err)
	require.Equal(t, prevote1.Voter, prevote.Voter)
	require.Equal(t, prevote1.SubmitBlock, prevote.SubmitBlock)
	require.Equal(t, len(prevote1.Hash), len(prevote.Hash))

	// validator 1 address params
	queryParams = types.NewQueryAggregatePrevoteParams(ValAddrs[1])
	bz, err = cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err = querier(input.Ctx, []string{types.QueryAggregatePrevote}, req)
	require.NoError(t, err)

	err = cdc.UnmarshalJSON(res, &prevote)
	require.NoError(t, err)
	require.Equal(t, prevote2.Voter, prevote.Voter)
	require.Equal(t, prevote2.SubmitBlock, prevote.SubmitBlock)
	require.Equal(t, len(prevote2.Hash), len(prevote.Hash))
}

func TestQueryAggregateVote(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	vote1 := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{"", sdk.OneDec()}}, ValAddrs[0])
	input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, vote1)
	vote2 := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{"", sdk.OneDec()}}, ValAddrs[1])
	input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, vote2)
	vote3 := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{"", sdk.OneDec()}}, ValAddrs[2])
	input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, vote3)

	// validator 0 address params
	queryParams := types.NewQueryAggregateVoteParams(ValAddrs[0])
	bz, err := cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryAggregateVote}, req)
	require.NoError(t, err)

	var vote types.AggregateExchangeRateVote
	err = cdc.UnmarshalJSON(res, &vote)
	require.NoError(t, err)
	require.Equal(t, vote1, vote)

	// validator 1 address params
	queryParams = types.NewQueryAggregateVoteParams(ValAddrs[1])
	bz, err = cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req = abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err = querier(input.Ctx, []string{types.QueryAggregateVote}, req)
	require.NoError(t, err)

	err = cdc.UnmarshalJSON(res, &vote)
	require.NoError(t, err)
	require.Equal(t, vote2, vote)
}

func TestQueryVoteTargets(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	// clear tobin taxes
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)

	voteTargets := []string{"denom", "denom2", "denom3"}
	for _, target := range voteTargets {
		input.OracleKeeper.SetTobinTax(input.Ctx, target, sdk.OneDec())
	}

	req := abci.RequestQuery{
		Path: "",
		Data: nil,
	}

	res, err := querier(input.Ctx, []string{types.QueryVoteTargets}, req)
	require.NoError(t, err)

	var voteTargetsRes []string
	err2 := cdc.UnmarshalJSON(res, &voteTargetsRes)
	require.NoError(t, err2)
	require.Equal(t, voteTargets, voteTargetsRes)
}

func TestQueryTobinTaxes(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	// clear tobin taxes
	input.OracleKeeper.ClearTobinTaxes(input.Ctx)

	tobinTaxes := types.DenomList{{types.MicroKRWDenom, sdk.OneDec()}, {types.MicroSDRDenom, sdk.NewDecWithPrec(123, 2)}}
	for _, item := range tobinTaxes {
		input.OracleKeeper.SetTobinTax(input.Ctx, item.Name, item.TobinTax)
	}

	req := abci.RequestQuery{
		Path: "",
		Data: nil,
	}

	res, err := querier(input.Ctx, []string{types.QueryTobinTaxes}, req)
	require.NoError(t, err)

	var tobinTaxesRes types.DenomList
	err2 := cdc.UnmarshalJSON(res, &tobinTaxesRes)
	require.NoError(t, err2)
	require.Equal(t, tobinTaxes, tobinTaxesRes)
}

func TestQueryTobinTax(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	input := CreateTestInput(t)
	querier := NewQuerier(input.OracleKeeper)

	denom := types.Denom{types.MicroKRWDenom, sdk.OneDec()}
	input.OracleKeeper.SetTobinTax(input.Ctx, denom.Name, denom.TobinTax)

	queryParams := types.NewQueryTobinTaxParams(types.MicroKRWDenom)
	bz, err := cdc.MarshalJSON(queryParams)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: "",
		Data: bz,
	}

	res, err := querier(input.Ctx, []string{types.QueryTobinTax}, req)
	require.NoError(t, err)

	var tobinTaxRes sdk.Dec
	cdc.UnmarshalJSON(res, &tobinTaxRes)
	require.Equal(t, denom.TobinTax, tobinTaxRes)
}
