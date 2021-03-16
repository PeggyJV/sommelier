package keeper

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/armon/go-metrics"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bridgetypes "github.com/althea-net/peggy/module/x/peggy/types"

	"github.com/peggyjv/sommelier/x/il/types"
	oracletypes "github.com/peggyjv/sommelier/x/oracle/types"
)

// BeginBlock recreates stoploss positions that may have timeout.
// CONTRACT: this logic assumes that upon execution of the transaction on Ethereum
// the bridge relayers submit a receipt msg to cosmos so that we can delete
// the executed stoploss from the queue.
func (k Keeper) BeginBlock(_ sdk.Context) {
	// iterate over the queue of executed stoploss positions in asc order of eth block heights
	//  if current ethereum blockHeight < executed ethereum block + block timeout:
	//  	return true // break
	// 	recreate stoploss position with same fields
	// 	delete position from queue
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// check if the latest eth block height
	ethHeight := k.ethBridgeKeeper.GetLastObservedEthereumBlockHeight(ctx)
	if ethHeight.EthereumBlockHeight == 0 {
		k.Logger(ctx).Debug("tracked ethereum height is 0")
		return
	}

	params := k.GetParams(ctx)

	invalidationID := k.GetInvalidationID(ctx)

	// variable for the original value to check if we need to store the new invalidationID value
	originalInvalidationID := invalidationID
	deadline := ctx.BlockTime().Unix() + int64(params.EthTimeoutTimestamp)*int64(time.Second)

	pairMap := make(map[string]*oracletypes.UniswapPair)

	// TODO: fill up contracts
	batchedPositions := types.SimpleLogicBatch{
		Amounts:       []*big.Int{},
		Payloads:      []types.Payload{},
		LogicContract: common.Address{},
		TokenContract: common.Address{},
	}

	var transfers, fees []*bridgetypes.ERC20Token

	k.IterateStoplossPositions(ctx, func(address sdk.AccAddress, stoploss types.Stoploss) (stop bool) {
		pair, ok := pairMap[stoploss.UniswapPairID]
		if !ok {
			oracleData, _ := k.oracleKeeper.GetLatestAggregatedOracleData(ctx, oracletypes.UniswapDataType, stoploss.UniswapPairID)
			if oracleData == nil {
				// pair not found, continue with the next position
				k.Logger(ctx).Debug(
					"pair not found on provided oracle data",
					"pair-id", stoploss.UniswapPairID,
					"height", fmt.Sprintf("%d", ctx.BlockHeight()),
				)
				return false
			}

			// panic if it is not a uniswap pair type
			pair = oracleData.(*oracletypes.UniswapPair)

			// set the pair to the map in care there another stoploss position with the same pair
			pairMap[stoploss.UniswapPairID] = pair
		}

		// check if total supply is 0 to avoid panics
		if pair.TotalSupply.IsZero() {
			k.Logger(ctx).Debug(
				"0 total supply for pair",
				"pair-id", stoploss.UniswapPairID,
				"height", fmt.Sprintf("%d", ctx.BlockHeight()),
			)
			return false
		}

		positionShares := sdk.NewDec(int64(stoploss.LiquidityPoolShares))

		// Calculate the current USD value of the position so that we can calculate the impermanent loss
		usdValueOfPosition := positionShares.Mul(pair.ReserveUSD).Quo(pair.TotalSupply)

		currentSlippage := usdValueOfPosition.Quo(stoploss.ReferencePairRatio)

		if currentSlippage.LTE(stoploss.MaxSlippage) {
			// threshold not met, continue with the next position
			return false
		}

		// Since the current slipage is grater than max slippage we now withdraw the liquidity of the LP for the uniswap pair
		// that is suffering impermanent loss

		uniswapERC20Pair := bridgetypes.NewERC20Token(stoploss.LiquidityPoolShares, stoploss.UniswapPairID)

		// TODO: who pays the fee?
		ethFee := bridgetypes.NewERC20Token(0, stoploss.UniswapPairID)

		var redeemCall types.Payload

		if stoploss.RedeemEth {
			redeemCall = types.NewRedeemLiquidityETHCall(
				pair.ID,
				stoploss.LiquidityPoolShares,
				0, 0, // min values are 0
				stoploss.ReceiverAddress,
				deadline,
			)
		} else {
			redeemCall = types.NewRedeemLiquidityCall(
				pair.Token0.ID,
				pair.Token1.ID,
				stoploss.LiquidityPoolShares,
				0, 0, // min values are 0
				stoploss.ReceiverAddress,
				deadline,
			)
		}

		batchedPositions.Payloads = append(batchedPositions.Payloads, redeemCall)
		// TODO: add amounts
		batchedPositions.Amounts = append(batchedPositions.Amounts, big.NewInt(0))

		transfers = append(transfers, uniswapERC20Pair)
		fees = append(fees, ethFee)

		defer func() {
			// amount metric
			telemetry.SetGaugeWithLabels(
				[]string{"stoploss", "execution"},
				float32(stoploss.LiquidityPoolShares),
				[]metrics.Label{telemetry.NewLabel("pair", stoploss.UniswapPairID)},
			)

			// counter metric
			telemetry.IncrCounterWithLabels(
				[]string{"stoploss", "execution"},
				1,
				[]metrics.Label{telemetry.NewLabel("pair", stoploss.UniswapPairID)},
			)
		}()

		// TODO: technically we should remove the stoploss position now that is has been executed, but we still need to account
		// for the cases when the outgoing txs to ethereum fail for some reason.
		k.DeleteStoplossPosition(ctx, address, stoploss.UniswapPairID)

		// log and emit metrics
		k.Logger(ctx).Debug("stoploss batched", "pair", stoploss.UniswapPairID, "receiver-address", stoploss.ReceiverAddress)

		return false
	})

	// return if there is no payload
	if len(batchedPositions.Payloads) == 0 {
		k.Logger(ctx).Debug("no payload for batched execution")
		return
	}

	// increment the invalidation ID counter
	invalidationID++

	// encode the simple logic batch ABI
	payload, err := batchedPositions.GetEncodedCall()
	if err != nil {
		panic(fmt.Errorf("sommelier contract ABI payload pack failed: %w", err))
	}

	// NOTE: by setting the invalidation nonce always to 0 and the invalidation ID to an increasing
	// counter, will prevent a logic call to get invalidated unless the outgoing Ethereum transaction
	// times out

	call := &bridgetypes.OutgoingLogicCall{
		Transfers:            transfers,
		Fees:                 fees,
		LogicContractAddress: params.ContractAddress,
		Payload:              payload,
		Timeout:              ethHeight.EthereumBlockHeight + params.EthTimeoutBlocks,
		InvalidationId:       sdk.Uint64ToBigEndian(invalidationID), // TODO: should this be hex?
		InvalidationNonce:    0,
	}

	// send eth transaction to withdraw lp_shares liquidity for pair_id
	k.ethBridgeKeeper.SetOutgoingLogicCall(ctx, call)

	// Set the new invalidation
	if originalInvalidationID != invalidationID {
		k.SetInvalidationID(ctx, invalidationID)
	}

	// log and emit metrics
	k.Logger(ctx).Debug("stoploss batch txs executed", "invalidation-id", strconv.FormatUint(invalidationID, 64))
}
