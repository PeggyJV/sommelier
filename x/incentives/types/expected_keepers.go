//go:generate mockgen  -destination=../testutil/expected_keepers_mocks.go -package=keeper github.com/peggyjv/sommelier/v5/x/incentives/types AccountKeeper,DistributionKeeper,BankKeeper,MintKeeper

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// DistributionKeeper defines the expected distribution keeper methods
type DistributionKeeper interface {
	GetFeePool(ctx sdk.Context) (feePool distributiontypes.FeePool)
	SetFeePool(ctx sdk.Context, feePool distributiontypes.FeePool)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	BlockedAddr(addr sdk.AccAddress) bool
}

// MintKeeper defines the expected mint keeper methods
type MintKeeper interface {
	GetParams(ctx sdk.Context) minttypes.Params
	StakingTokenSupply(ctx sdk.Context) sdk.Int
	BondedRatio(ctx sdk.Context) sdk.Dec
}
