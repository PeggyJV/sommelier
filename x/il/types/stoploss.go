package types

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	bridgetypes "github.com/cosmos/gravity-bridge/module/x/gravity/types"
)

// LPsStoplossPositions is a slice of liquidity provider (LP) stoploss positions
// This type is used during genesis initialization and export logic.
type LPsStoplossPositions []StoplossPositions

var _ sort.Interface = LPsStoplossPositions{}

// Len implements sort.Interface for Traces
func (lps LPsStoplossPositions) Len() int { return len(lps) }

// Less implements sort.Interface for Traces
func (lps LPsStoplossPositions) Less(i, j int) bool { return lps[i].Address < lps[j].Address }

// Swap implements sort.Interface for Traces
func (lps LPsStoplossPositions) Swap(i, j int) { lps[i], lps[j] = lps[j], lps[i] }

// Sort the LPs stoploss positions slice by address
func (lps LPsStoplossPositions) Sort() LPsStoplossPositions {
	sort.Sort(lps)
	return lps
}

// Validate performs a basic validation of the stoploss positions. It checks
// 	- duplicate account address entries
//  - duplicate account stoploss pairs entries
// 	- address validation
//  - stoploss validation
func (lps LPsStoplossPositions) Validate() error {
	seenLps := make(map[string]bool)
	for _, lpStoplossPosition := range lps {
		// check duplicate account address
		if seenLps[lpStoplossPosition.Address] {
			return fmt.Errorf("duplicate stoploss positions for liquidity provider %s", lpStoplossPosition.Address)
		}

		_, err := sdk.AccAddressFromBech32(lpStoplossPosition.Address)
		if err != nil {
			return err
		}

		seenPositions := make(map[string]bool)
		for _, position := range lpStoplossPosition.StoplossPositions {
			// check duplicate stoploss pair for the account
			if seenPositions[position.UniswapPairID] {
				return fmt.Errorf("duplicated stoploss position for liquidity provider %s and pair %s", lpStoplossPosition.Address, position.UniswapPairID)
			}

			if err := position.Validate(); err != nil {
				return fmt.Errorf("invalid position for address %s and pair %s, %w", lpStoplossPosition.Address, position.UniswapPairID, err)
			}

			seenPositions[position.UniswapPairID] = true
		}
		seenLps[lpStoplossPosition.Address] = true
	}
	return nil
}

// Validate performs a basic validation of the stoploss fields
func (s Stoploss) Validate() error {
	if err := bridgetypes.ValidateEthAddress(s.UniswapPairID); err != nil {
		return fmt.Errorf("invalid uniswap pair id: %w", err)
	}

	if s.MaxSlippage.LTE(sdk.ZeroDec()) || s.MaxSlippage.GT(sdk.NewDec(1)) {
		return fmt.Errorf("max slippage must be (0,1], got %s", s.MaxSlippage)
	}
	if s.LiquidityPoolShares == 0 {
		return fmt.Errorf("liquidity pool shares must be positive, got %d", s.LiquidityPoolShares)
	}
	if s.ReferencePairRatio.LTE(sdk.ZeroDec()) || s.ReferencePairRatio.GT(sdk.NewDec(1)) {
		return fmt.Errorf("reference pair ratio must be (0,1], got %s", s.ReferencePairRatio)
	}

	if err := bridgetypes.ValidateEthAddress(s.ReceiverAddress); err != nil {
		return fmt.Errorf("invalid ethereum receiver address: %w", err)
	}

	return nil
}
