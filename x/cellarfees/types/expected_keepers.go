package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/ethereum/go-ethereum/common"
	auctiontypes "github.com/peggyjv/sommelier/v7/x/auction/types"
)

// AccountKeeper defines the expected account keeper.
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) types.ModuleAccountI
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	BlockedAddr(addr sdk.AccAddress) bool
}

// CorkKeeper defines the minimum interface needed to check registered Cellars
type CorkKeeper interface {
	GetCellarIDs(ctx sdk.Context) (cellars []common.Address)
	HasCellarID(ctx sdk.Context, address common.Address) (found bool)
}

// GravityKeeper defines the expected gravity keeper methods
type GravityKeeper interface {
	ERC20ToDenomLookup(ctx sdk.Context, tokenContract common.Address) (bool, string)
}

// AuctionKeeper defines the expected auction keeper methods
type AuctionKeeper interface {
	GetActiveAuctions(ctx sdk.Context) []*auctiontypes.Auction
	GetTokenPrice(ctx sdk.Context, denom string) (auctiontypes.TokenPrice, bool)
	GetTokenPrices(ctx sdk.Context) []auctiontypes.TokenPrice
	BeginAuction(ctx sdk.Context,
		startingTokensForSale sdk.Coin,
		initialPriceDecreaseRate sdk.Dec,
		priceDecreaseBlockInterval uint64,
		fundingModuleAccount string,
		proceedsModuleAccount string) error
}

// MintKeeper defines the expected mint keeper methods
type MintKeeper interface {
	GetParams(ctx sdk.Context) minttypes.Params
	StakingTokenSupply(ctx sdk.Context) math.Int
	BondedRatio(ctx sdk.Context) sdk.Dec
}
