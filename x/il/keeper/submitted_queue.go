package keeper

import (
	"fmt"
	"strconv"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// TrackPositionTimeout iterates over all the stoploss positions that have been
// processed during EndBlock and successfully submitted to the ethereum bridge.
// In case the position timeout height is greater or equal, we reenable the position
// to keep tracking the impermanent loss for the uniswap pair.
func (k Keeper) TrackPositionTimeout(ctx sdk.Context, currentEthHeight uint64) {
	k.IterateSubmittedQueue(ctx, func(timeoutHeight uint64, address sdk.AccAddress, pairID common.Address) bool {
		if currentEthHeight < timeoutHeight {
			// break the iteration as all future iterations won't timeout
			return true
		}

		// since the position timeout on ethereum, we need to reenable it
		// in order to continue tracking the impermanent loss on EndBlock
		stoploss, found := k.GetStoplossPosition(ctx, address, pairID.String())
		if !found {
			panic(fmt.Errorf("stoploss owned for pair id %s not present during tracking at timeout height %d", pairID, timeoutHeight))
		}

		// set the updated position with the submitted value to false
		stoploss.Submitted = false
		k.SetStoplossPosition(ctx, address, stoploss)

		// delete submitted position from this queue
		k.DeleteSubmittedPosition(ctx, timeoutHeight, address, pairID)

		k.Logger(ctx).Info(
			"prev submitted position reenabled due to timeout",
			"pair-id", pairID,
			"ethereum timeout height", strconv.FormatUint(timeoutHeight, 64),
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
