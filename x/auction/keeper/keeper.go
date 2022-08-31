package keeper

import (
	"bytes"
	"fmt"

	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the auction store
type Keeper struct {
	storeKey               sdk.StoreKey
	cdc                    codec.BinaryCodec
	paramSpace             paramtypes.Subspace
	bankKeeper             types.BankKeeper
	fundingModuleAccounts  map[string]string
	proceedsModuleAccounts map[string]string
}

// NewKeeper creates a new auction Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	bankKeeper types.BankKeeper, fundingModuleAccounts map[string]string, proceedsModuleAccounts map[string]string,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:               key,
		cdc:                    cdc,
		paramSpace:             paramSpace,
		bankKeeper:             bankKeeper,
		fundingModuleAccounts:  fundingModuleAccounts,
		proceedsModuleAccounts: proceedsModuleAccounts,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

////////////
// Params //
////////////

// GetParamSet returns the vote period from the parameters
func (k Keeper) GetParamSet(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSpace.GetParamSet(ctx, &p)
	return p
}

// setParams sets the parameters in the store
func (k Keeper) setParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

//////////////
// Auctions //
//////////////

// GetActiveAuctionByID returns a specific active auction
func (k Keeper) GetActiveAuctionByID(ctx sdk.Context, id uint32) (types.Auction, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetActiveAuctionKey(id))
	if len(bz) != 0 {
		return types.Auction{}, false
	}

	var auction types.Auction
	k.cdc.MustUnmarshal(bz, &auction)
	return auction, true
}

// DeleteActiveAuction deletes the active auction
func (k Keeper) deleteActiveAuction(ctx sdk.Context, id uint32) {
	ctx.KVStore(k.storeKey).Delete(types.GetActiveAuctionKey(id))
}

// GetEndedAuctionByID returns a specific active auction
func (k Keeper) GetEndedAuctionByID(ctx sdk.Context, id uint32) (types.Auction, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetEndedAuctionKey(id))
	if len(bz) != 0 {
		return types.Auction{}, false
	}

	var auction types.Auction
	k.cdc.MustUnmarshal(bz, &auction)
	return auction, true
}

// IterateAuctions iterates over all auctions in the store for a given prefix
func (k Keeper) IterateAuctions(ctx sdk.Context, auctionTypePrefix []byte, handler func(auctionID uint32, auction types.Auction) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, auctionTypePrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := bytes.NewBuffer(iter.Key())
		key.Next(1) // trim prefix byte

		auctionID := binary.BigEndian.Uint32(key.Bytes())

		var auction types.Auction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		if handler(auctionID, auction) {
			break
		}
	}
}

// GetActiveAuctions returns all active auctions
func (k Keeper) GetActiveAuctions(ctx sdk.Context) []*types.Auction {
	var auctions []*types.Auction
	k.IterateAuctions(ctx, types.GetActiveAuctionsPrefix(), func(auctionID uint32, auction types.Auction) (stop bool) {
		auctions = append(auctions, &auction)
		return false
	})

	return auctions
}

// GetEndedAuctions returns all inactive auctions (that have not been pruned)
func (k Keeper) GetEndedAuctions(ctx sdk.Context) []*types.Auction {
	var auctions []*types.Auction
	k.IterateAuctions(ctx, types.GetEndedAuctionsPrefix(), func(auctionID uint32, auction types.Auction) (stop bool) {
		auctions = append(auctions, &auction)
		return false
	})

	return auctions
}

// SetActiveAuction sets the auction specified
func (k Keeper) setActiveAuction(ctx sdk.Context, auction types.Auction) {
	bz := k.cdc.MustMarshal(&auction)
	ctx.KVStore(k.storeKey).Set(types.GetActiveAuctionKey(auction.Id), bz)
}

// SetEndedAuction sets the auction specified
func (k Keeper) setEndedAuction(ctx sdk.Context, auction types.Auction) {
	bz := k.cdc.MustMarshal(&auction)
	ctx.KVStore(k.storeKey).Set(types.GetEndedAuctionKey(auction.Id), bz)
}

// BeginAuction starts a new auction for a single denomination
func (k Keeper) BeginAuction(ctx sdk.Context,
	startingAmount sdk.Coin,
	initialDecreaseRate float32,
	blockDecreaseInterval uint32,
	fundingModuleAccount string,
	proceedsModuleAccount string) error {
	// Verify sale token price freshness is acceptable
	lastSaleTokenPrice, found := k.GetTokenPrice(ctx, startingAmount.Denom)
	if !found {
		return types.ErrCouldNotFindSaleTokenPrice
	} else if lastSaleTokenPrice.LastUpdatedBlock == 0 || (uint64(ctx.BlockHeight())-lastSaleTokenPrice.LastUpdatedBlock) > k.GetParamSet(ctx).PriceMaxBlockAge {
		return types.ErrLastSaleTokenPriceUpdateTooLongAgo
	}

	// Verify somm token price freshness is acceptable
	lastSommTokenPrice, found := k.GetTokenPrice(ctx, "usomm")
	if !found {
		return types.ErrCouldNotFindSommTokenPrice
	} else if lastSommTokenPrice.LastUpdatedBlock == 0 || (uint64(ctx.BlockHeight())-lastSommTokenPrice.LastUpdatedBlock) > k.GetParamSet(ctx).PriceMaxBlockAge {
		return types.ErrLastSommTokenPriceUpdateTooLongAgo
	}

	// Calculate somm per sale token price

	// Starting price is amount of usomm required for 1 of starting denom
	saleTokenUSDPriceFloat, err := lastSaleTokenPrice.UsdPrice.Float64()

	if err != nil {
		return types.ErrConvertingTokenPriceToFloat
	}

	usommTokenUSDPriceFloat, err := lastSommTokenPrice.UsdPrice.Float64()

	if err != nil {
		return types.ErrConvertingTokenPriceToFloat
	}

	saleTokenPriceInUsomm := saleTokenUSDPriceFloat / usommTokenUSDPriceFloat
	saleTokenPriceInUsommDec, err := sdk.NewDecFromStr(fmt.Sprintf("%f", saleTokenPriceInUsomm))
	if err != nil {
		return types.ErrConvertingStringToDec
	}

	// Validate starting amount
	if !startingAmount.Amount.IsPositive() {
		return types.ErrAuctionStartingAmountMustBePositve
	}

	if startingAmount.Denom == "" {
		return types.ErrAuctionDenomInvalid
	}

	if startingAmount.Denom == "usomm" {
		return types.ErrCannotAuctionUsomm
	}

	// Validate initial decrease rate
	if initialDecreaseRate <= 0 || initialDecreaseRate >= 1 {
		return types.ErrInvalidInitialDecreaseRate
	}

	// Validate block decrease interval
	if blockDecreaseInterval == 0 {
		return types.ErrInvalidBlockDecreaeInterval
	}

	// Validate funding module
	if _, ok := k.fundingModuleAccounts[authtypes.NewModuleAddress(fundingModuleAccount).String()]; !ok {
		return types.ErrUnauthorizedFundingModule
	}

	// Validate proceeds module
	if _, ok := k.proceedsModuleAccounts[authtypes.NewModuleAddress(proceedsModuleAccount).String()]; !ok {
		return types.ErrUnauthorizedFundingModule
	}

	// Validate no ongoing auction for denom
	for _, auction := range k.GetActiveAuctions(ctx) {
		if auction.AmountRemaining.Denom == startingAmount.Denom {
			return types.ErrCannotStartTwoAuctionsForSameDenomSimultaneously
		}
	}

	// Transfer the coins
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, fundingModuleAccount, types.ModuleName, sdk.Coins{startingAmount}); err != nil {
		return err
	}

	auctionID := k.GetLastAuctionID(ctx) + 1

	// Add auction to active auctions
	k.setActiveAuction(ctx, types.Auction{
		Id:                      auctionID,
		StartingAmount:          startingAmount,
		StartBlock:              uint64(ctx.BlockHeight()),
		EndBlock:                0,
		InitialDecreaseRate:     initialDecreaseRate,
		CurrentDecreaseRate:     initialDecreaseRate,
		BlockDecreaseInterval:   blockDecreaseInterval,
		CurrentUnitPriceInUsomm: saleTokenPriceInUsommDec,
		AmountRemaining:         startingAmount,
		FundingModuleAccount:    fundingModuleAccount,
		ProceedsModuleAccount:   proceedsModuleAccount,
	})

	// Update last auction id
	k.setLastAuctionID(ctx, auctionID)

	// Emit event that auction has started
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeNewAuction,
				sdk.NewAttribute(types.AttributeKeyAuctionID, fmt.Sprint(k.GetLastAuctionID(ctx))),
				sdk.NewAttribute(types.AttributeKeyStartBlock, fmt.Sprint(ctx.BlockHeight())),
				sdk.NewAttribute(types.AttributeKeyInitialDecreaseRate, fmt.Sprintf("%f", initialDecreaseRate)),
				sdk.NewAttribute(types.AttributeKeyBlockDecreaseInterval, fmt.Sprint(blockDecreaseInterval)),
				sdk.NewAttribute(types.AttributeKeyStartingDenom, startingAmount.Denom),
				sdk.NewAttribute(types.AttributeKeyStartingAmount, startingAmount.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyStartingUsommPrice, fmt.Sprintf("%f", saleTokenPriceInUsomm)),
			),
		},
	)

	return nil
}

// FinishAuction completes an auction by sending relevant funds to destination addresses and updates state
func (k Keeper) FinishAuction(ctx sdk.Context, auction *types.Auction) error {
	// Figure out how many funds we have left over, if any, to send
	// Since we can only have 1 auction per denom active at a time, we can just query the balance
	saleTokenBalance := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.ModuleName), auction.StartingAmount.Denom)

	if saleTokenBalance.Amount.IsPositive() {
		// Send remaining funds to their appropriate destination module
		if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, auction.FundingModuleAccount, sdk.Coins{saleTokenBalance}); err != nil {
			return err
		}
	}

	// Calculate amount of usomm proceeds to send back from total bids for auction
	bids := k.GetBidsByAuctionID(ctx, auction.Id)
	usommProceeds := sdk.NewDec(int64(0))

	for _, bid := range bids {
		usommProceeds.Add(bid.TotalFulfilledSaleTokenAmount.Amount.ToDec().Mul(bid.UnitPriceOfSaleTokenInUsomm))
	}

	// Send proceeds to their appropriate destination module
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, auction.ProceedsModuleAccount, sdk.Coins{sdk.NewCoin("usomm", sdk.NewInt(usommProceeds.TruncateInt64()))}); err != nil {
		return err
	}

	// Remove auction from active list
	k.deleteActiveAuction(ctx, auction.Id)

	// Move auction to ended auctions list with updated fields
	auction.EndBlock = uint64(ctx.BlockHeight())
	k.setEndedAuction(ctx, *auction)

	// Emit event that auction has finished
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeAuctionFinished,
				sdk.NewAttribute(types.AttributeKeyAuctionID, fmt.Sprint(auction.Id)),
				sdk.NewAttribute(types.AttributeKeyStartBlock, fmt.Sprint(auction.StartBlock)),
				sdk.NewAttribute(types.AttributeKeyEndBlock, fmt.Sprint(auction.EndBlock)),
				sdk.NewAttribute(types.AttributeKeyInitialDecreaseRate, fmt.Sprintf("%f", auction.InitialDecreaseRate)),
				sdk.NewAttribute(types.AttributeKeyCurrentDecreaseRate, fmt.Sprintf("%f", auction.CurrentDecreaseRate)),
				sdk.NewAttribute(types.AttributeKeyBlockDecreaseInterval, fmt.Sprint(auction.BlockDecreaseInterval)),
				sdk.NewAttribute(types.AttributeKeyStartingDenom, auction.StartingAmount.Denom),
				sdk.NewAttribute(types.AttributeKeyStartingAmount, auction.StartingAmount.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyAmountRemaining, auction.AmountRemaining.String()),
			),
		},
	)

	return nil
}

//////////////
//   Bids   //
//////////////

// IterateBids iterates over all bids in the store
func (k Keeper) IterateBids(ctx sdk.Context, handler func(auctionID uint32, bidID uint64, bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBidsByAuctionPrefix())
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := bytes.NewBuffer(iter.Key())
		key.Next(1)                                       // trim prefix byte
		auctionID := binary.BigEndian.Uint32(key.Next(4)) // trim auction bytes

		bidID := binary.BigEndian.Uint64(key.Bytes())

		var bid types.Bid
		k.cdc.MustUnmarshal(iter.Value(), &bid)
		if handler(auctionID, bidID, bid) {
			break
		}
	}
}

// IterateBidsByAuction iterates over all bids in the store for a given auction
func (k Keeper) IterateBidsByAuction(ctx sdk.Context, auctionID uint32, handler func(auctionID uint32, bidID uint64, bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBidsByAuctionIDPrefix(auctionID))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := bytes.NewBuffer(iter.Key())
		key.Next(1)                                       // trim prefix byte
		auctionID := binary.BigEndian.Uint32(key.Next(4)) // trim auction bytes

		bidID := binary.BigEndian.Uint64(key.Bytes())

		var bid types.Bid
		k.cdc.MustUnmarshal(iter.Value(), &bid)
		if handler(auctionID, bidID, bid) {
			break
		}
	}
}

// GetBids returns all stored bids (that have not been pruned)
func (k Keeper) GetBids(ctx sdk.Context) []*types.Bid {
	var bids []*types.Bid
	k.IterateBids(ctx, func(auctionID uint32, bidID uint64, bid types.Bid) (stop bool) {
		bids = append(bids, &bid)
		return false
	})

	return bids
}

// GetBidsByAuctionID returns all stored bids for an auction id (that have not been pruned)
func (k Keeper) GetBidsByAuctionID(ctx sdk.Context, auctionID uint32) []*types.Bid {
	var bids []*types.Bid
	k.IterateBidsByAuction(ctx, auctionID, func(auctionId uint32, bidID uint64, bid types.Bid) (stop bool) {
		bids = append(bids, &bid)
		return false
	})

	return bids
}

// GetBid returns a specified bid by its id (if it has not been pruned)
func (k Keeper) GetBid(ctx sdk.Context, auctionID uint32, bidID uint64) (types.Bid, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetBidKey(auctionID, bidID))
	if len(bz) != 0 {
		return types.Bid{}, false
	}

	var bid types.Bid
	k.cdc.MustUnmarshal(bz, &bid)
	return bid, true
}

// SetBid sets the bid specified
func (k Keeper) setBid(ctx sdk.Context, bid types.Bid) {
	bz := k.cdc.MustMarshal(&bid)
	ctx.KVStore(k.storeKey).Set(types.GetBidKey(bid.GetAuctionId(), bid.GetId()), bz)
}

/////////////////
// TokenPrices //
/////////////////

// IterateTokenPrices iterates over all token prices in the store
func (k Keeper) IterateTokenPrices(ctx sdk.Context, handler func(denom string, tokenPrice types.TokenPrice) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetTokenPricesPrefix())
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := bytes.NewBuffer(iter.Key())
		key.Next(1) // trim prefix byte

		denom := key.String()

		var tokenPrice types.TokenPrice
		k.cdc.MustUnmarshal(iter.Value(), &tokenPrice)
		if handler(denom, tokenPrice) {
			break
		}
	}
}

// GetTokenPrices returns all stored token prices
func (k Keeper) GetTokenPrices(ctx sdk.Context) []*types.TokenPrice {
	var tokenPrices []*types.TokenPrice
	k.IterateTokenPrices(ctx, func(denom string, tokenPrice types.TokenPrice) (stop bool) {
		tokenPrices = append(tokenPrices, &tokenPrice)
		return false
	})

	return tokenPrices
}

// GetTokenPrice returns the stored token price
func (k Keeper) GetTokenPrice(ctx sdk.Context, denom string) (types.TokenPrice, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetTokenPriceKey(denom))
	if len(bz) != 0 {
		return types.TokenPrice{}, false
	}

	var tokenPrice types.TokenPrice
	k.cdc.MustUnmarshal(bz, &tokenPrice)
	return tokenPrice, true
}

// SetTokenPrice sets the token price specified
func (k Keeper) setTokenPrice(ctx sdk.Context, tokenPrice types.TokenPrice) {
	bz := k.cdc.MustMarshal(&tokenPrice)
	ctx.KVStore(k.storeKey).Set(types.GetTokenPriceKey(tokenPrice.GetDenom()), bz)
}

//////////////////
//  Id storage  //
//////////////////

// setLastAuctionID sets the last auction id
func (k Keeper) setLastAuctionID(ctx sdk.Context, id uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, id)

	ctx.KVStore(k.storeKey).Set(types.GetLastAuctionIDKey(), b)
}

// setLastBidID sets the last bid id
func (k Keeper) setLastBidID(ctx sdk.Context, id uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)

	ctx.KVStore(k.storeKey).Set(types.GetLastBidIDKey(), b)
}

// GetLastAuctionID gets the last auction id
func (k Keeper) GetLastAuctionID(ctx sdk.Context) uint32 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetLastAuctionIDKey())
	if len(bz) != 0 {
		return 0
	}

	return binary.BigEndian.Uint32(bz)
}

// GetLastBidID gets the last bid id
func (k Keeper) GetLastBidID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetLastBidIDKey())
	if len(bz) != 0 {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}
