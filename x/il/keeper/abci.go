package keeper

import (
	"time"

	"github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/x/il/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	params := k.GetParams(ctx)

	k.IterateStoplossPositions(ctx, func(address sdk.AccAddress, stoploss types.Stoploss) (stop bool) {
		//   TODO: fetch pair info from oracle given the pair address

		//   calculate current_ratio/reference_pair_ratio
		currentSlippage := stoploss.ReferencePairRatio

		if currentSlippage.LTE(stoploss.MaxSlippage) {
			//
			return false
		}

		// Since the current slipage is grater than max slippage we now withdraw the liquidity of the LP for the uniswap pair
		// that is suffering impermanent loss

		// peggy-0x... ?
		uniswapSharesCoin := sdk.NewInt64Coin(stoploss.UniswapPairId, stoploss.LiquidityPoolShares) // TODO: What's the name of the uniswap coin
		fee := sdk.Coin{}                                                                           //                                                                             // TODO: which coin should we use as fees?
		// send eth transaction to withdraw lp_shares liquidity for pair_id
		_, err := k.ethBridgeKeeper.AddToOutgoingPool(ctx, address, params.ContractAddress, uniswapSharesCoin, fee)
		if err != nil {
			// FIXME: figure out why this might panic and ensure correctness
			panic(err)
		}

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
