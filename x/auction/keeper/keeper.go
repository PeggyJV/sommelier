package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the auction store
type Keeper struct {
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	paramSpace paramtypes.Subspace
	bankKeeper types.BankKeeper
}

// NewKeeper creates a new auction Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	bankKeeper types.BankKeeper,
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
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

//////////////
// Auctions //
//////////////

// GetAllAuctions returns all stored auctions
func (k Keeper) GetAllAuctions(ctx sdk.Context) []*types.Auction {
	// TODO: Fill in
	return nil
}

// GetAuctionById returns a specific auction
func (k Keeper) GetAuctionsById(ctx sdk.Context, id uint32) *types.Auction {
	// TODO: Fill in
	return nil
}

// GetAllActiveAuctions returns all active auctions
func (k Keeper) GetAllActiveAuctions(ctx sdk.Context) []*types.Auction {
	// TODO: Fill in
	return nil
}


// GetAllInactiveAuctions returns all inactive auctions (that have not been pruned)
func (k Keeper) GetAllInactiveAuctions(ctx sdk.Context, id uint32) []*types.Auction {
	// TODO: Fill in
	return nil
}


// SetAuctions sets the auctions specified
func (k Keeper) SetAuctions(ctx sdk.Context, auctions []*types.Auction) {
	// TODO: Fill in
}

// BeginAuction starts a new auction for a single denomination
func (k Keeper) BeginAuction(ctx sdk.Context,
	startingAmount sdk.Coin,
	initialDecreaseRate float32,
	blockDecreaseInterval uint16,
	fundingModuleAccount authtypes.ModuleAccountI,
	proceeedsModuleAccount authtypes.AccountI) error {
	// TODO: Verify inputs as first step, return error if problematic

	// TODO: Fill in

	return nil
}

//////////////
//   Bids   //
//////////////

// GetBids returns all stored bids (that have not been pruned)
func (k Keeper) GetBids(ctx sdk.Context) []*types.Bid {
	store := ctx.KVStore(k.storeKey)

	
	// TODO: Fill in
	return nil
}

// SetBids sets the bids specified
func (k Keeper) SetBids(ctx sdk.Context, bids []*types.Bid) {
	// TODO: Fill in
}

// SetBid sets the bid specified
func (k Keeper) SetBid(ctx sdk.Context, bids types.Bid) {
	// TODO: Fill in
}

// GetBid returns a specified bid by its id (if it has not been pruned)
func (k Keeper) GetBid(ctx sdk.Context, id uint64) *types.Bid {
	// TODO: Fill in
	return nil
}

/////////////////
// TokenPrices //
/////////////////

// GetTokenPrices returns all stored token prices
func (k Keeper) GetTokenPrices(ctx sdk.Context) []*types.TokenPrice {
	// TODO: Fill in
	return nil
}

// SetTokenPrices sets the token prices specified
func (k Keeper) SetTokenPrices(ctx sdk.Context, tokenPrices []*types.TokenPrice) {
	// TODO: Fill in
}