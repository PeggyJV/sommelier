package types

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper is the subset of the staking keeper that the PoA module
// depends on. The wrapper in x/poa/keeper.WrappedStakingKeeper embeds this
// interface and overrides the rescaling-relevant methods.
type StakingKeeper interface {
	GetParams(ctx sdk.Context) stakingtypes.Params
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (stakingtypes.Validator, bool)
	GetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) int64
	GetLastTotalPower(ctx sdk.Context) math.Int
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator

	IterateValidators(sdk.Context, func(int64, stakingtypes.ValidatorI) bool)
	IterateBondedValidatorsByPower(sdk.Context, func(int64, stakingtypes.ValidatorI) bool)
	IterateLastValidators(sdk.Context, func(int64, stakingtypes.ValidatorI) bool)
	IterateLastValidatorPowers(sdk.Context, func(sdk.ValAddress, int64) (stop bool))

	Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI

	Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec) math.Int
	SlashWithInfractionReason(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec, stakingtypes.Infraction) math.Int

	Jail(sdk.Context, sdk.ConsAddress)
	Unjail(sdk.Context, sdk.ConsAddress)

	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI
	MaxValidators(sdk.Context) uint32
	PowerReduction(sdk.Context) math.Int
	BondDenom(sdk.Context) string
	UnbondingTime(sdk.Context) time.Duration
	IsValidatorJailed(sdk.Context, sdk.ConsAddress) bool

	SetLastValidatorPower(sdk.Context, sdk.ValAddress, int64)
	DeleteLastValidatorPower(sdk.Context, sdk.ValAddress)

	ValidatorQueueIterator(sdk.Context, time.Time, int64) sdk.Iterator
	Hooks() stakingtypes.StakingHooks
}

// SlashingKeeper exposes the bits of x/slashing the PoA EndBlocker needs for
// snapshot retention.
type SlashingKeeper interface {
	SignedBlocksWindow(sdk.Context) int64
}
