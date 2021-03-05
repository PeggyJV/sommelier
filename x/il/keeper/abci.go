package keeper

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/armon/go-metrics"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bridgetypes "github.com/althea-net/peggy/module/x/peggy/types"

	"github.com/peggyjv/sommelier/x/il/types"
	"github.com/peggyjv/sommelier/x/il/types/contract"
	oracletypes "github.com/peggyjv/sommelier/x/oracle/types"
)

// TODO: the stoploss position should be recreated if the tx on ethereum timeouts or fails.

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

	pairMap := make(map[string]*oracletypes.UniswapPair)

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

		abi, err := abi.JSON(strings.NewReader(contract.ContractABI))
		if err != nil {
			panic(fmt.Errorf("sommelier contract ABI encoding failed: %w", err))
		}

		deadline := big.NewInt(ctx.BlockTime().Unix() + 1000*int64(time.Second))

		// TODO: we should give the option to redeemLiquidityETH
		// TODO: fill the missing fields
		payload, err := abi.Pack(
			"redeemLiquidity",
			common.HexToAddress(pair.Token0.ID),             //	address tokenA
			common.HexToAddress(pair.Token1.ID),             //	address tokenB
			big.NewInt(int64(stoploss.LiquidityPoolShares)), //	uint256 liquidity
			big.NewInt(0),                                   //	uint256 amountAMin
			big.NewInt(0),                                   //	uint256 amountBMin
			common.HexToAddress(stoploss.ReceiverAddress),   // address to
			deadline, // uint256 deadline
		)

		if err != nil {
			panic(fmt.Errorf("sommelier contract ABI payload pack failed: %w", err))
		}

		// increment the invalidation ID counter
		invalidationID++

		// NOTE: by setting the invalidation nonce always to 0 and the invalidation ID to an increasing
		// counter, will prevent a logic call to get invalidated unless the outgoing Ethereum transaction
		// times out

		call := &bridgetypes.OutgoingLogicCall{
			Transfers:            []*bridgetypes.ERC20Token{uniswapERC20Pair},
			Fees:                 []*bridgetypes.ERC20Token{ethFee},
			LogicContractAddress: params.ContractAddress,
			Payload:              payload,
			Timeout:              ethHeight.EthereumBlockHeight + params.EthTimeout,
			InvalidationId:       sdk.Uint64ToBigEndian(invalidationID), // TODO: should this be hex?
			InvalidationNonce:    1,
		}

		// send eth transaction to withdraw lp_shares liquidity for pair_id
		k.ethBridgeKeeper.SetOutgoingLogicCall(ctx, call)

		// log and emit metrics
		k.Logger(ctx).Info("stoploss executed", "pair", stoploss.UniswapPairID, "receiver-address", stoploss.ReceiverAddress)

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

		return false
	})

	// Set the new invalidation
	if originalInvalidationID != invalidationID {
		k.SetInvalidationID(ctx, invalidationID)
	}
}
