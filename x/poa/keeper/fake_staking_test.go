package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// fakeStakingKeeper is a minimal implementation of types.StakingKeeper used
// to exercise the WrappedStakingKeeper. Methods not exercised by tests panic
// with "not implemented" so missing coverage is loud.
type fakeStakingKeeper struct {
	validators       map[string]stakingtypes.Validator        // operator bech32 -> validator
	consToOperator   map[string]sdk.ValAddress                // consAddr string -> op
	bondedOrder      []sdk.ValAddress                          // for IterateBondedValidatorsByPower
	lastPower        map[string]int64                          // operator bech32 -> overwritten LastValidatorPower
	lastSlashPower   int64                                     // last `power` arg passed to Slash
	lastSlashConsAddr sdk.ConsAddress
	slashCalled      bool
}

func newFakeStaking() *fakeStakingKeeper {
	return &fakeStakingKeeper{
		validators:     map[string]stakingtypes.Validator{},
		consToOperator: map[string]sdk.ValAddress{},
		lastPower:      map[string]int64{},
	}
}

func (f *fakeStakingKeeper) addValidator(op sdk.ValAddress, tokens math.Int) stakingtypes.Validator {
	v := stakingtypes.Validator{
		OperatorAddress: op.String(),
		Tokens:          tokens,
		Status:          stakingtypes.Bonded,
		DelegatorShares: sdk.NewDecFromInt(tokens),
	}
	f.validators[op.String()] = v
	return v
}

func (f *fakeStakingKeeper) mapCons(cons sdk.ConsAddress, op sdk.ValAddress) {
	f.consToOperator[cons.String()] = op
}

// --- types.StakingKeeper implementation ---

func (f *fakeStakingKeeper) GetParams(ctx sdk.Context) stakingtypes.Params {
	return stakingtypes.DefaultParams()
}

func (f *fakeStakingKeeper) GetValidator(ctx sdk.Context, op sdk.ValAddress) (stakingtypes.Validator, bool) {
	v, ok := f.validators[op.String()]
	return v, ok
}

func (f *fakeStakingKeeper) GetLastValidatorPower(ctx sdk.Context, op sdk.ValAddress) int64 {
	if p, ok := f.lastPower[op.String()]; ok {
		return p
	}
	v, ok := f.validators[op.String()]
	if !ok {
		return 0
	}
	return sdk.TokensToConsensusPower(v.Tokens, sdk.DefaultPowerReduction)
}

func (f *fakeStakingKeeper) GetLastTotalPower(ctx sdk.Context) math.Int {
	total := math.ZeroInt()
	for _, v := range f.validators {
		total = total.Add(math.NewInt(sdk.TokensToConsensusPower(v.Tokens, sdk.DefaultPowerReduction)))
	}
	return total
}

func (f *fakeStakingKeeper) GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator {
	out := make([]stakingtypes.Validator, 0, len(f.bondedOrder))
	for _, op := range f.bondedOrder {
		if v, ok := f.validators[op.String()]; ok {
			out = append(out, v)
		}
	}
	return out
}

func (f *fakeStakingKeeper) IterateValidators(ctx sdk.Context, cb func(int64, stakingtypes.ValidatorI) bool) {
	i := int64(0)
	for _, v := range f.validators {
		if cb(i, v) {
			return
		}
		i++
	}
}

func (f *fakeStakingKeeper) IterateBondedValidatorsByPower(ctx sdk.Context, cb func(int64, stakingtypes.ValidatorI) bool) {
	for i, op := range f.bondedOrder {
		v, ok := f.validators[op.String()]
		if !ok {
			continue
		}
		if cb(int64(i), v) {
			return
		}
	}
}

func (f *fakeStakingKeeper) IterateLastValidators(ctx sdk.Context, cb func(int64, stakingtypes.ValidatorI) bool) {
	f.IterateBondedValidatorsByPower(ctx, cb)
}

func (f *fakeStakingKeeper) IterateLastValidatorPowers(ctx sdk.Context, cb func(sdk.ValAddress, int64) bool) {
	for _, op := range f.bondedOrder {
		v, ok := f.validators[op.String()]
		if !ok {
			continue
		}
		var power int64
		if p, has := f.lastPower[op.String()]; has {
			power = p
		} else {
			power = sdk.TokensToConsensusPower(v.Tokens, sdk.DefaultPowerReduction)
		}
		if cb(op, power) {
			return
		}
	}
}

func (f *fakeStakingKeeper) Validator(ctx sdk.Context, op sdk.ValAddress) stakingtypes.ValidatorI {
	v, ok := f.validators[op.String()]
	if !ok {
		return nil
	}
	return v
}

func (f *fakeStakingKeeper) ValidatorByConsAddr(ctx sdk.Context, c sdk.ConsAddress) stakingtypes.ValidatorI {
	op, ok := f.consToOperator[c.String()]
	if !ok {
		return nil
	}
	return f.Validator(ctx, op)
}

func (f *fakeStakingKeeper) Slash(ctx sdk.Context, c sdk.ConsAddress, infractionHeight int64, power int64, _ sdk.Dec) math.Int {
	f.slashCalled = true
	f.lastSlashPower = power
	f.lastSlashConsAddr = c
	return math.NewInt(power)
}

func (f *fakeStakingKeeper) SlashWithInfractionReason(ctx sdk.Context, c sdk.ConsAddress, h int64, power int64, d sdk.Dec, _ stakingtypes.Infraction) math.Int {
	return f.Slash(ctx, c, h, power, d)
}

func (f *fakeStakingKeeper) Jail(ctx sdk.Context, _ sdk.ConsAddress)    {}
func (f *fakeStakingKeeper) Unjail(ctx sdk.Context, _ sdk.ConsAddress)  {}

func (f *fakeStakingKeeper) Delegation(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) stakingtypes.DelegationI {
	return nil
}
func (f *fakeStakingKeeper) MaxValidators(_ sdk.Context) uint32          { return 100 }
func (f *fakeStakingKeeper) PowerReduction(_ sdk.Context) math.Int       { return sdk.DefaultPowerReduction }
func (f *fakeStakingKeeper) BondDenom(_ sdk.Context) string              { return "usomm" }
func (f *fakeStakingKeeper) UnbondingTime(_ sdk.Context) time.Duration   { return 21 * 24 * time.Hour }
func (f *fakeStakingKeeper) IsValidatorJailed(_ sdk.Context, _ sdk.ConsAddress) bool { return false }
func (f *fakeStakingKeeper) SetLastValidatorPower(_ sdk.Context, op sdk.ValAddress, power int64) {
	f.lastPower[op.String()] = power
}
func (f *fakeStakingKeeper) DeleteLastValidatorPower(_ sdk.Context, op sdk.ValAddress) {
	delete(f.lastPower, op.String())
}
func (f *fakeStakingKeeper) ValidatorQueueIterator(_ sdk.Context, _ time.Time, _ int64) sdk.Iterator {
	return nil
}
func (f *fakeStakingKeeper) Hooks() stakingtypes.StakingHooks { return nil }

func (f *fakeStakingKeeper) GetAllValidators(_ sdk.Context) []stakingtypes.Validator {
	out := make([]stakingtypes.Validator, 0, len(f.validators))
	for _, v := range f.validators {
		out = append(out, v)
	}
	return out
}

func (f *fakeStakingKeeper) IterateDelegations(_ sdk.Context, _ sdk.AccAddress, _ func(int64, stakingtypes.DelegationI) bool) {
}
func (f *fakeStakingKeeper) GetAllSDKDelegations(_ sdk.Context) []stakingtypes.Delegation {
	return nil
}
func (f *fakeStakingKeeper) GetAllDelegatorDelegations(_ sdk.Context, _ sdk.AccAddress) []stakingtypes.Delegation {
	return nil
}
