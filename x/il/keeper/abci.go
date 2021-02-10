package keeper

import (
	"fmt"
	"time"

	"github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	peggytypes "github.com/althea-net/peggy/module/x/peggy/types"

	"github.com/peggyjv/sommelier/x/il/types"
	oracletypes "github.com/peggyjv/sommelier/x/oracle/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	params := k.GetParams(ctx)

	// TODO: this should fetch an oracle data per identifier
	oracleData := k.oracleKeeper.GetOracleData(ctx, oracletypes.UniswapDataType)
	if oracleData == nil {
		return
	}

	uniswapData, ok := oracleData.(*oracletypes.UniswapData)
	if !ok {
		return
	}

	pairs, err := uniswapData.Parse()
	if err != nil {
		// this shouldn't happen
		panic(err)
	}

	pairMap := make(map[string]oracletypes.UniswapPairParsed)
	for _, pair := range pairs {
		pairMap[pair.ID] = pair
	}

	k.IterateStoplossPositions(ctx, func(address sdk.AccAddress, stoploss types.Stoploss) (stop bool) {
		//   TODO: fetch pair info from oracle given the pair address
		pair, ok := pairMap[stoploss.UniswapPairId]
		if !ok {
			// pair not found, continue with the next position
			k.Logger(ctx).Error(
				"pair not found on provided oracle data",
				"pair ID", stoploss.UniswapPairId,
			)
			return false
		}

		if pair.TotalSupply.String() == "" {
			k.Logger(ctx).Error(
				"failed to parse float from total supply",
				"value", fmt.Sprintf("%v", pair.TotalSupply),
			)

			return false
		}

		if pair.ReserveUsd.String() == "" {
			k.Logger(ctx).Error(
				"failed to parse float from reserve usd",
				"value", fmt.Sprintf("%v", pair.ReserveUsd),
			)

			return false
		}

		positionShares := sdk.NewDec(stoploss.LiquidityPoolShares)

		// Calculate the current USD value of the position so that we can calculate the impermanent loss
		usdValueOfPosition := positionShares.Mul(pair.ReserveUsd).Quo(pair.TotalSupply)

		currentSlippage := usdValueOfPosition.Quo(stoploss.ReferencePairRatio)

		if currentSlippage.LTE(stoploss.MaxSlippage) {
			// threshold not met, continue with the next position
			return false
		}

		// Since the current slipage is grater than max slippage we now withdraw the liquidity of the LP for the uniswap pair
		// that is suffering impermanent loss

		uniswapERC20Pair := peggytypes.NewERC20Token(uint64(stoploss.LiquidityPoolShares), stoploss.UniswapPairId)
		uniswapCoin := uniswapERC20Pair.PeggyCoin()

		// NOTE: denom must always match coin
		fee := sdk.NewInt64Coin(uniswapCoin.Denom, 0)

		// send eth transaction to withdraw lp_shares liquidity for pair_id
		_, err = k.ethBridgeKeeper.AddToOutgoingPool(ctx, address, params.ContractAddress, uniswapCoin, fee)
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
