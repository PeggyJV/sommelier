package keeper

import (
	"fmt"
	"time"

	"github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bridgetypes "github.com/althea-net/peggy/module/x/peggy/types"

	"github.com/peggyjv/sommelier/x/il/types"
	oracletypes "github.com/peggyjv/sommelier/x/oracle/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	params := k.GetParams(ctx)
	pairMap := make(map[string]*oracletypes.UniswapPair)

	k.IterateStoplossPositions(ctx, func(address sdk.AccAddress, stoploss types.Stoploss) (stop bool) {
		pair, ok := pairMap[stoploss.UniswapPairId]
		if !ok {
			oracleData, _ := k.oracleKeeper.GetLatestAggregatedOracleData(ctx, oracletypes.UniswapDataType, stoploss.UniswapPairId)
			if oracleData == nil {
				// pair not found, continue with the next position
				k.Logger(ctx).Debug(
					"pair not found on provided oracle data",
					"pair-id", stoploss.UniswapPairId,
					"height", fmt.Sprintf("%d", ctx.BlockHeight()),
				)
				return false
			}

			// panic if it is not a uniswap pair type
			pair = oracleData.(*oracletypes.UniswapPair)
			pairMap[stoploss.UniswapPairId] = pair
		}

		// check if total supply is 0 to avoid panics
		if pair.TotalSupply.IsZero() {
			k.Logger(ctx).Debug(
				"0 total supply for pair",
				"pair-id", stoploss.UniswapPairId,
				"height", fmt.Sprintf("%d", ctx.BlockHeight()),
			)
			return false
		}

		positionShares := sdk.NewDec(stoploss.LiquidityPoolShares)

		// Calculate the current USD value of the position so that we can calculate the impermanent loss
		usdValueOfPosition := positionShares.Mul(pair.ReserveUSD).Quo(pair.TotalSupply)

		currentSlippage := usdValueOfPosition.Quo(stoploss.ReferencePairRatio)

		if currentSlippage.LTE(stoploss.MaxSlippage) {
			// threshold not met, continue with the next position
			return false
		}

		// Since the current slipage is grater than max slippage we now withdraw the liquidity of the LP for the uniswap pair
		// that is suffering impermanent loss

		uniswapERC20Pair := bridgetypes.NewERC20Token(uint64(stoploss.LiquidityPoolShares), stoploss.UniswapPairId)

		// TODO: fill the missing fields
		call := &bridgetypes.OutgoingLogicCall{
			Transfers:            []*bridgetypes.ERC20Token{uniswapERC20Pair},
			Fees:                 nil, // TODO: who pays the fee
			LogicContractAddress: params.ContractAddress,
			Payload:              nil,
			Timeout:              0,
			InvalidationId:       nil,
			InvalidationNonce:    0,
		}

		// send eth transaction to withdraw lp_shares liquidity for pair_id
		k.ethBridgeKeeper.SetOutgoingLogicCall(ctx, call)

		// log and emit metrics
		k.Logger(ctx).Info("stoploss executed", "pair", stoploss.UniswapPairId, "address", address.String())

		defer func() {
			// amount metric
			telemetry.SetGaugeWithLabels(
				[]string{"stoploss", "execution"},
				float32(stoploss.LiquidityPoolShares),
				[]metrics.Label{telemetry.NewLabel("pair", stoploss.UniswapPairId)},
			)

			// counter metric
			telemetry.IncrCounterWithLabels(
				[]string{"stoploss", "execution"},
				1,
				[]metrics.Label{telemetry.NewLabel("pair", stoploss.UniswapPairId)},
			)
		}()

		return false
	})
}
