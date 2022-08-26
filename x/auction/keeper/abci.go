package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

// BeginBlocker is called at the beginning of every block
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	for _, auction := range k.GetActiveAuctions(ctx) {
		// TODO Step function for auction price updates
		auction.Id = uint32(0)

	}


	// TODO: do we need to do anything with proposal voting periods here?
}

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// End Auctions that have no funds left
	for _, auction := range k.GetActiveAuctions(ctx) {
		if auction.AmountRemaining.Amount.IsZero() {
			// Figure out how many funds we have to send
			supply := k.bankKeeper.GetSupply(ctx, auction.StartingAmount.Denom)

			// Send proceeds to their appropriate destination module
			if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, auction.ProceedsModuleAccount, sdk.Coins{supply}); err != nil {
				// TODO: Panic here or no?
				panic(err)
			}

			// Move auction to ended auctions list with updated fields
			k.setEndedAuction(ctx, types.Auction{
				Id: auction.Id,
				StartingAmount: auction.StartingAmount,
				StartBlock: auction.StartBlock,
				EndBlock: uint64(ctx.BlockHeight()),
				InitialDecreaseRate: auction.InitialDecreaseRate,
				CurrentDecreaseRate: auction.CurrentDecreaseRate,
				BlockDecreaseInterval: auction.BlockDecreaseInterval,
				CurrentPrice: auction.CurrentPrice,
				AmountRemaining: auction.AmountRemaining,
				ProceedsModuleAccount: auction.ProceedsModuleAccount,
			})

			// Remove auction from active list
			k.deleteActiveAuction(ctx, auction.Id)
		}
	}

	// TODO: anything else? Trim down old auctions and bids maybe?
}
