package keeper

import (
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) TrackPositionTimeout(ctx sdk.Context, currentEthHeight, timeout uint64) {
	k.IterateExecutionQueue(ctx, func(executedHeight uint64, address sdk.AccAddress, pairID string) bool {

		if currentEthHeight < executedHeight+timeout {
			// break the iteration as all future iterations will won't timeout
			return true
		}

		// since the position timeouted on ethereum, we need to reenable it with the
		// in order to continue tracking IL on EndBlock
		stoploss, found := k.GetStoplossPosition(ctx, address, pairID)
		if !found {
			panic("stoploss should be present")
		}

		stoploss.Executed = false
		// set the updated position
		k.SetStoplossPosition(ctx, address, stoploss)

		// delete executed position from queue
		k.DeleteExecutedPosition(ctx, executedHeight, address)

		k.Logger(ctx).Info(
			"prev executed position reenabled due to timeout",
			"pair-id", pairID,
			"receiver-address", stoploss.ReceiverAddress,
		)

		defer func() {
			// amount metric
			telemetry.SetGaugeWithLabels(
				[]string{"stoploss", "timeout"},
				float32(stoploss.LiquidityPoolShares),
				[]metrics.Label{telemetry.NewLabel("pair", stoploss.UniswapPairID)},
			)

			// counter metric
			telemetry.IncrCounterWithLabels(
				[]string{"stoploss", "timeout"},
				1,
				[]metrics.Label{telemetry.NewLabel("pair", stoploss.UniswapPairID)},
			)
		}()

		return false
	})
}
