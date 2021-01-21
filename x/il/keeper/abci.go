package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// iterate over each stoploss position
	//   fetch pair info from oracle
	//   calculate current_ratio/reference_pair_ratio
	//   if slippage > max_slipage
	//       send eth transaction to withdraw lp_shares liquidity for pair_id
}
