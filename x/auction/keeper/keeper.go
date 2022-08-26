package keeper

import (
	"fmt"
	"bytes"
	"strconv"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
	"github.com/tendermint/tendermint/libs/log"
	"encoding/binary"
)

// Keeper of the auction store
type Keeper struct {
	storeKey                 sdk.StoreKey
	cdc                      codec.BinaryCodec
	paramSpace               paramtypes.Subspace
	bankKeeper               types.BankKeeper
	fundingModuleAccounts    map[string]string
	proceedsModuleAccounts   map[string]string
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
		storeKey:   key,
		cdc:        cdc,
		paramSpace: paramSpace,
		bankKeeper: bankKeeper,
		fundingModuleAccounts: fundingModuleAccounts,
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

// GetActiveAuctionById returns a specific active auction
func (k Keeper) GetActiveAuctionById(ctx sdk.Context, id uint32) (types.Auction, bool) {
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

// DeleteEndedAuction deletes the ended auction
func (k Keeper) deleteEndedAuction(ctx sdk.Context, id uint32) {
	ctx.KVStore(k.storeKey).Delete(types.GetEndedAuctionKey(id))
}

// GetEndedAuctionById returns a specific active auction
func (k Keeper) GetEndedAuctionById(ctx sdk.Context, id uint32) (types.Auction, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetEndedAuctionKey(id))
	if len(bz) != 0 {
		return types.Auction{}, false
	}

	var auction types.Auction
	k.cdc.MustUnmarshal(bz, &auction)
	return auction, true
}

// iterateAuctions iterates over all auctions in the store for a given prefix
func (k Keeper) IterateAuctions(ctx sdk.Context, auctionTypePrefix []byte, handler func(auctionId uint32, auction types.Auction) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, auctionTypePrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := bytes.NewBuffer(iter.Key())
		key.Next(1) // trim prefix byte

		auctionId := binary.BigEndian.Uint32(key.Bytes())

		var auction types.Auction
		k.cdc.MustUnmarshal(iter.Value(), &auction)
		if handler(auctionId, auction) {
			break
		}
	}
}

// GetActiveAuctions returns all active auctions
func (k Keeper) GetActiveAuctions(ctx sdk.Context) []*types.Auction {
	var auctions []*types.Auction
	k.IterateAuctions(ctx, types.GetActiveAuctionsPrefix(), func(auctionId uint32, auction types.Auction) (stop bool) {
		auctions = append(auctions, &auction)
		return false
	})

	return auctions
}

// GetEndedAuctions returns all inactive auctions (that have not been pruned)
func (k Keeper) GetEndedAuctions(ctx sdk.Context) []*types.Auction {
	var auctions []*types.Auction
	k.IterateAuctions(ctx, types.GetEndedAuctionsPrefix(), func(auctionId uint32, auction types.Auction) (stop bool) {
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
	fundingModule string,
	proceedsModule string) error {
	// Verify token price freshness is acceptable
	lastUpdate := k.GetLastTokenPriceUpdateBlock(ctx)
	if lastUpdate == 0 || (uint64(ctx.BlockHeight()) - lastUpdate) > k.GetParamSet(ctx).PriceMaxBlockAge {
		return fmt.Errorf("cannot perform action, last price update was too long ago, submit new prices and try again")
	}

	// Validate starting amount
	if !startingAmount.Amount.IsPositive() {
		return fmt.Errorf("minimum amount must be a positive amount of coins")
	}

	if startingAmount.Denom == "" {
		return fmt.Errorf("denom cannot be empty")
	}

	if startingAmount.Denom == "usomm" {
		return fmt.Errorf("auction pointless")
	}

	// Validate initial decrease rate
	if initialDecreaseRate <= 0 || initialDecreaseRate >= 1 {
		return fmt.Errorf("initial decrease rate must be a float less than one and greater than zero")
	}

	// Validate block decrease interval
	if blockDecreaseInterval == 0 {
		return fmt.Errorf("block decrease interval cannot be 0")
	}

	// Validate funding module
	if _, ok := k.fundingModuleAccounts[authtypes.NewModuleAddress(fundingModule).String()]; !ok {
		return fmt.Errorf("unauthorized funding module")
	}

	// Validate proceeds module 
	if _, ok := k.proceedsModuleAccounts[authtypes.NewModuleAddress(proceedsModule).String()]; !ok {
		return fmt.Errorf("unauthorized proceeds module")
	}

	// Validate no ongoing auction for denom
	for _, auction := range k.GetActiveAuctions(ctx) {
		if auction.CurrentPrice.Denom == startingAmount.Denom {
			return fmt.Errorf("auction for this denom is currently ongoing, cannot create another auction for the same denom until completed")
		}
	}

	// Transfer the coins
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, fundingModule, types.ModuleName, sdk.Coins{startingAmount}); err != nil {
		return err
	}

	currentSaleTokenPrice, found := k.GetTokenPrice(ctx, startingAmount.Denom)

	if !found {
		return fmt.Errorf("no price data found for starting amount denom")
	}
	
	// Calculate starting price
	sommPrice, found := k.GetTokenPrice(ctx, "usomm")

	if !found {
		return fmt.Errorf("no price data found for starting amount denom")
	}

	// Starting price is amount of usomm required for 1 of starting denom
	salePriceFloat := currentSaleTokenPrice.UsdPrice.MustFloat64()
	sommPriceFloat := sommPrice.UsdPrice.MustFloat64()

	saleTokenPriceInUsomm := int64(salePriceFloat / sommPriceFloat)

	if saleTokenPriceInUsomm < 1 {
		saleTokenPriceInUsomm = 1
	}

	// Add auction to active auctions
	k.setActiveAuction(ctx, types.Auction{
		Id: k.GetLastAuctionId(ctx) + 1, 
		StartingAmount: startingAmount,
		StartBlock: uint64(ctx.BlockHeight()),
		EndBlock: 0,
		InitialDecreaseRate: initialDecreaseRate,
		CurrentDecreaseRate: initialDecreaseRate,
		BlockDecreaseInterval: blockDecreaseInterval,
		CurrentPrice: sdk.NewCoin(currentSaleTokenPrice.Denom, sdk.NewInt(saleTokenPriceInUsomm)), 
		AmountRemaining: startingAmount,
		ProceedsModuleAccount: proceedsModule,
	})

	// Update last auction id
	k.setLastAuctionId(ctx, k.GetLastAuctionId(ctx) + 1)

	// Emit event that auction has started
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeAuction,
				sdk.NewAttribute(types.AttributeKeyAuctionId, string(k.GetLastAuctionId(ctx))),
				sdk.NewAttribute(types.AttributeKeyStartBlock, string(ctx.BlockHeight())),
				sdk.NewAttribute(types.AttributeKeyInitialDecreaseRate, strconv.FormatFloat(float64(initialDecreaseRate), 'E', -1, 32)),
				sdk.NewAttribute(types.AttributeKeyBlockDecreaseInterval, string(blockDecreaseInterval)),
				sdk.NewAttribute(types.AttributeKeyStartingDenom, startingAmount.Denom),
				sdk.NewAttribute(types.AttributeKeyStartingAmount,startingAmount.Amount.String()),
				sdk.NewAttribute(types.AttributeKeyStartingUsommPrice, string(saleTokenPriceInUsomm)),
			),
		},
	)

	return nil
}

//////////////
//   Bids   //
//////////////

// IterateBids iterates over all bids in the store 
func (k Keeper) IterateBids(ctx sdk.Context, handler func(bidId uint64, bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBidsByAuctionPrefix())
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := bytes.NewBuffer(iter.Key())
		key.Next(1) // trim prefix byte
		key.Next(4) // trim auction bytes

		bidId := binary.BigEndian.Uint64(key.Bytes())

		var bid types.Bid
		k.cdc.MustUnmarshal(iter.Value(), &bid)
		if handler(bidId, bid) {
			break
		}
	}
}

// IterateBidsByAuction iterates over all bids in the store for a given auction
func (k Keeper) IterateBidsByAuction(ctx sdk.Context, auctionId uint32, handler func(bidId uint64, bid types.Bid) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetBidsByAuctionIdPrefix(auctionId))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := bytes.NewBuffer(iter.Key())
		key.Next(1) // trim prefix byte
		key.Next(4) // trim auction bytes

		bidId := binary.BigEndian.Uint64(key.Bytes())

		var bid types.Bid
		k.cdc.MustUnmarshal(iter.Value(), &bid)
		if handler(bidId, bid) {
			break
		}
	}
}

// GetBids returns all stored bids (that have not been pruned)
func (k Keeper) GetBids(ctx sdk.Context) []*types.Bid {
	var bids []*types.Bid
	k.IterateBids(ctx, func(bidId uint64, bid types.Bid) (stop bool) {
		bids = append(bids, &bid)
		return false
	})

	return bids
}

// GetBidsByAuctionId returns all stored bids for an auction id (that have not been pruned)
func (k Keeper) GetBidsByAuctionId(ctx sdk.Context, auctionId uint32) []*types.Bid {
	var bids []*types.Bid
	k.IterateBidsByAuction(ctx, auctionId, func(bidId uint64, bid types.Bid) (stop bool) {
		bids = append(bids, &bid)
		return false
	})

	return bids
}

// GetBid returns a specified bid by its id (if it has not been pruned)
func (k Keeper) GetBid(ctx sdk.Context, auctionId uint32, bidId uint64) (types.Bid, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetBidKey(auctionId, bidId))
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

		denom := string(key.Bytes())

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
	k.setLastTokenPriceUpdateBlock(ctx, uint64(ctx.BlockHeight()))
}

//////////////////
//  Id storage  //
//////////////////

// setLastAuctionId sets the last auction id
func (k Keeper) setLastAuctionId(ctx sdk.Context, id uint32) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, id)

	ctx.KVStore(k.storeKey).Set(types.GetLastAuctionIdKey(), b)
}

// setLastBidId sets the last bid id
func (k Keeper) setLastBidId(ctx sdk.Context, id uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)

	ctx.KVStore(k.storeKey).Set(types.GetLastBidIdKey(), b)
}

// setLastTokenPriceUpdateBlock sets the last block a token price update was performed
func (k Keeper) setLastTokenPriceUpdateBlock(ctx sdk.Context, block uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, block)

	ctx.KVStore(k.storeKey).Set(types.GetLastBidIdKey(), b)
}


// GetLastAuctionId gets the last auction id
func (k Keeper) GetLastAuctionId(ctx sdk.Context) uint32 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetLastAuctionIdKey())
	if len(bz) != 0 {
		return 0
	}

	return binary.BigEndian.Uint32(bz)
}

// GetLastBidId gets the last bid id
func (k Keeper) GetLastBidId(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetLastBidIdKey())
	if len(bz) != 0 {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}

// GetLastTokenPriceUpdateBlock gets the last bid id
func (k Keeper) GetLastTokenPriceUpdateBlock(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetLastTokenPriceBlock())
	if len(bz) != 0 {
		return 0
	}

	return binary.BigEndian.Uint64(bz)
}