//go:generate mockgen  -destination=../testutil/expected_keepers_mocks.go -package=keeper github.com/peggyjv/sommelier/v9/x/incentives/types AccountKeeper,DistributionKeeper,BankKeeper,MintKeeper

package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// DistributionKeeper defines the expected distribution keeper methods
type DistributionKeeper interface {
	GetFeePool(ctx sdk.Context) (feePool distributiontypes.FeePool)
	SetFeePool(ctx sdk.Context, feePool distributiontypes.FeePool)
	GetValidatorOutstandingRewards(ctx sdk.Context, valAddr sdk.ValAddress) (rewards distributiontypes.ValidatorOutstandingRewards)
	SetValidatorOutstandingRewards(ctx sdk.Context, valAddr sdk.ValAddress, rewards distributiontypes.ValidatorOutstandingRewards)
	GetValidatorCurrentRewards(ctx sdk.Context, valAddr sdk.ValAddress) (rewards distributiontypes.ValidatorCurrentRewards)
	SetValidatorCurrentRewards(ctx sdk.Context, valAddr sdk.ValAddress, rewards distributiontypes.ValidatorCurrentRewards)
	GetValidatorAccumulatedCommission(ctx sdk.Context, valAddr sdk.ValAddress) (commission distributiontypes.ValidatorAccumulatedCommission)
	SetValidatorAccumulatedCommission(ctx sdk.Context, valAddr sdk.ValAddress, commission distributiontypes.ValidatorAccumulatedCommission)
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
	StakingTokenSupply(ctx sdk.Context) math.Int
	BondedRatio(ctx sdk.Context) sdk.Dec
}

type StakingKeeper interface {
	ValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) stakingtypes.ValidatorI
}
