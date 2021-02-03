package keeper

import (
	"time"

	"github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	peggytypes "github.com/althea-net/peggy/module/x/peggy/types"
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
			// threshold not met, continue with the next position
			return false
		}

		// Since the current slipage is grater than max slippage we now withdraw the liquidity of the LP for the uniswap pair
		// that is suffering impermanent loss

		// peggy-0x... ?
		uniswapERC20Pair := peggytypes.NewERC20Token(uint64(stoploss.LiquidityPoolShares), stoploss.UniswapPairId)
		uniswapCoin := uniswapERC20Pair.PeggyCoin()

		// NOTE: denom must always match coin
		fee := sdk.NewInt64Coin(uniswapCoin.Denom, 0)

		// send eth transaction to withdraw lp_shares liquidity for pair_id
		_, err := k.ethBridgeKeeper.AddToOutgoingPool(ctx, address, params.ContractAddress, uniswapCoin, fee)
		if err != nil {
			// FIXME: figure out why this might error and ensure correctness
			k.Logger(ctx).Error(
				"failed to send to eth bridge's outgoing pool",
				"sender-address", address.String(),
				"uniswap-coin", uniswapCoin.String(),
				"error", err.Error(),
			)

			// continue with the next position
			return false
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
