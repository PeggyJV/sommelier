package oracle

import (
	"testing"

	"github.com/peggyjv/sommelier/x/oracle/types"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/oracle/keeper"
)

func TestExportInitGenesis(t *testing.T) {
	input, _ := setup(t)

	input.OracleKeeper.SetOracleDelegate(input.Ctx, keeper.ValAddrs[0], keeper.Addrs[1])
	input.OracleKeeper.SetLunaExchangeRate(input.Ctx, "denom", sdk.NewDec(123))
	input.OracleKeeper.AddAggregateExchangeRatePrevote(input.Ctx, types.NewAggregateExchangeRatePrevote([]byte{123}, sdk.ValAddress{}, int64(2)))
	input.OracleKeeper.AddAggregateExchangeRateVote(input.Ctx, types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{Denom: "foo", ExchangeRate: sdk.NewDec(123)}}, sdk.ValAddress{}))
	input.OracleKeeper.SetTobinTax(input.Ctx, "denom", sdk.NewDecWithPrec(123, 3))
	input.OracleKeeper.SetTobinTax(input.Ctx, "denom2", sdk.NewDecWithPrec(123, 3))
	genesis := ExportGenesis(input.Ctx, input.OracleKeeper)

	newInput := keeper.CreateTestInput(t)
	InitGenesis(newInput.Ctx, newInput.OracleKeeper, genesis)
	newGenesis := ExportGenesis(newInput.Ctx, newInput.OracleKeeper)

	require.Equal(t, genesis, newGenesis)
}
