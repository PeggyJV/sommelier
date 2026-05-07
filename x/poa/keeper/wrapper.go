package keeper

import (
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/peggyjv/sommelier/v9/x/poa/types"
)

// WrappedStakingKeeper is the rescaling adapter passed to every consumer of
// validator power (slashing, distribution, evidence, gravity-bridge,
// cork, axelarcork, pubsub, incentives). It embeds the underlying staking
// keeper and overrides the methods that return validator power so authority
// validators report their boosted ConsensusPower / Tokens.
//
// The Slash methods divide caller-supplied power by the snapshot multiplier
// at the infraction height before delegating, so authority validators are
// penalised on raw stake rather than on boosted consensus power.
type WrappedStakingKeeper struct {
	types.StakingKeeper
	poa Keeper
}

// WrappedStakingKeeper returns the wrapper. The wrapper holds a reference to
// the underlying keeper, so reads see live state.
func (k Keeper) WrappedStakingKeeper() WrappedStakingKeeper {
	return WrappedStakingKeeper{StakingKeeper: k.sk, poa: k}
}

// ----------------------------------------------------------------------------
// Internal helpers
// ----------------------------------------------------------------------------

// currentMultiplier returns the boost multiplier for `op` at the current
// block height, or 1 if the operator is not in the authority allowlist or
// no boost is in effect this block.
func (w WrappedStakingKeeper) currentMultiplier(ctx sdk.Context, op sdk.ValAddress) sdk.Dec {
	if !w.poa.IsAuthority(ctx, op) {
		return sdk.OneDec()
	}
	m := w.poa.MultiplierForValidator(ctx, op, ctx.BlockHeight())
	if m.LTE(sdk.OneDec()) {
		return sdk.OneDec()
	}
	return m
}

// boostedTokens scales `raw` by the current-block multiplier with Ceil
// semantics, matching the EndBlocker's overwrite of LastValidatorPower.
// Using Ceil consistently across all boosted reads (Tokens, ConsensusPower)
// avoids the asymmetry where one path returned floor(raw*m) and another
// returned ceil(raw*m), which under-slashed authority validators by 1 power
// unit at the boundary.
func (w WrappedStakingKeeper) boostedTokens(ctx sdk.Context, op sdk.ValAddress, raw math.Int) math.Int {
	m := w.currentMultiplier(ctx, op)
	if m.Equal(sdk.OneDec()) {
		return raw
	}
	return sdk.NewDecFromInt(raw).Mul(m).Ceil().TruncateInt()
}

// adaptValidator wraps `v` in a boostedValidator if the operator is authority
// and the current-block multiplier exceeds 1; otherwise returns `v` unchanged.
func (w WrappedStakingKeeper) adaptValidator(ctx sdk.Context, v stakingtypes.ValidatorI) stakingtypes.ValidatorI {
	if v == nil {
		return nil
	}
	op := v.GetOperator()
	m := w.currentMultiplier(ctx, op)
	if m.Equal(sdk.OneDec()) {
		return v
	}
	return boostedValidator{ValidatorI: v, multiplier: m}
}

// ----------------------------------------------------------------------------
// Validator-returning read methods
// ----------------------------------------------------------------------------

// GetLastTotalPower returns the boosted aggregate consensus power. The
// embedded staking keeper's LastTotalPower in the store reflects the raw
// total computed inside ApplyAndReturnValidatorSetUpdates BEFORE PoA
// overwrites individual LastValidatorPower entries, so we must recompute
// here to keep gravity-bridge / cork / axelarcork quorum thresholds in sync
// with the boosted per-validator powers they iterate.
func (w WrappedStakingKeeper) GetLastTotalPower(ctx sdk.Context) math.Int {
	total := math.ZeroInt()
	w.StakingKeeper.IterateLastValidatorPowers(ctx, func(_ sdk.ValAddress, power int64) bool {
		total = total.Add(math.NewInt(power))
		return false
	})
	return total
}

func (w WrappedStakingKeeper) Validator(ctx sdk.Context, op sdk.ValAddress) stakingtypes.ValidatorI {
	return w.adaptValidator(ctx, w.StakingKeeper.Validator(ctx, op))
}

func (w WrappedStakingKeeper) ValidatorByConsAddr(ctx sdk.Context, c sdk.ConsAddress) stakingtypes.ValidatorI {
	return w.adaptValidator(ctx, w.StakingKeeper.ValidatorByConsAddr(ctx, c))
}

func (w WrappedStakingKeeper) IterateBondedValidatorsByPower(ctx sdk.Context, cb func(int64, stakingtypes.ValidatorI) bool) {
	w.StakingKeeper.IterateBondedValidatorsByPower(ctx, func(i int64, v stakingtypes.ValidatorI) bool {
		return cb(i, w.adaptValidator(ctx, v))
	})
}

func (w WrappedStakingKeeper) IterateLastValidators(ctx sdk.Context, cb func(int64, stakingtypes.ValidatorI) bool) {
	w.StakingKeeper.IterateLastValidators(ctx, func(i int64, v stakingtypes.ValidatorI) bool {
		return cb(i, w.adaptValidator(ctx, v))
	})
}

func (w WrappedStakingKeeper) IterateValidators(ctx sdk.Context, cb func(int64, stakingtypes.ValidatorI) bool) {
	w.StakingKeeper.IterateValidators(ctx, func(i int64, v stakingtypes.ValidatorI) bool {
		return cb(i, w.adaptValidator(ctx, v))
	})
}

// GetBondedValidatorsByPower returns concrete []stakingtypes.Validator (gravity-bridge
// signature). We mutate the Tokens field on each returned copy and re-sort.
func (w WrappedStakingKeeper) GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator {
	raw := w.StakingKeeper.GetBondedValidatorsByPower(ctx)
	out := make([]stakingtypes.Validator, len(raw))
	for i, v := range raw {
		op, err := sdk.ValAddressFromBech32(v.OperatorAddress)
		if err == nil {
			v.Tokens = w.boostedTokens(ctx, op, v.Tokens)
		}
		out[i] = v
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Tokens.GT(out[j].Tokens)
	})
	return out
}

// GetValidator returns the concrete validator with Tokens scaled.
func (w WrappedStakingKeeper) GetValidator(ctx sdk.Context, op sdk.ValAddress) (stakingtypes.Validator, bool) {
	v, found := w.StakingKeeper.GetValidator(ctx, op)
	if !found {
		return v, false
	}
	v.Tokens = w.boostedTokens(ctx, op, v.Tokens)
	return v, true
}

// GetAllValidators returns the full validator set with authority Tokens
// rescaled. Consumed by slashing and distribution.
func (w WrappedStakingKeeper) GetAllValidators(ctx sdk.Context) []stakingtypes.Validator {
	raw := w.StakingKeeper.GetAllValidators(ctx)
	out := make([]stakingtypes.Validator, len(raw))
	for i, v := range raw {
		op, err := sdk.ValAddressFromBech32(v.OperatorAddress)
		if err == nil {
			v.Tokens = w.boostedTokens(ctx, op, v.Tokens)
		}
		out[i] = v
	}
	return out
}

// ----------------------------------------------------------------------------
// Slash normalisation
// ----------------------------------------------------------------------------

// Slash converts caller-supplied (boosted) consensus power back to raw stake
// for authority validators before delegating to the underlying keeper. If a
// snapshot at the infraction height is missing for an authority validator,
// the slash is REFUSED (returns 0) to avoid silently over-slashing — see
// design spec §3.4 and Codex review item 3.
func (w WrappedStakingKeeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) math.Int {
	val := w.StakingKeeper.ValidatorByConsAddr(ctx, consAddr) // raw lookup, NOT adapted
	if val == nil {
		return w.StakingKeeper.Slash(ctx, consAddr, infractionHeight, power, slashFactor)
	}
	op := val.GetOperator()
	if !w.poa.IsAuthority(ctx, op) {
		return w.StakingKeeper.Slash(ctx, consAddr, infractionHeight, power, slashFactor)
	}
	m, snapFound := w.poa.MultiplierForValidatorWithStatus(ctx, op, infractionHeight)
	if !snapFound {
		ctx.Logger().Error("poa: missing multiplier snapshot for authority slash; SKIPPING slash",
			"operator", op.String(), "infraction_height", infractionHeight)
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeSlashSkippedNoSnapshot,
			sdk.NewAttribute(types.AttributeOperator, op.String()),
			sdk.NewAttribute(types.AttributeInfractionHeight, sdk.NewInt(infractionHeight).String()),
		))
		return math.ZeroInt()
	}
	if m.LTE(sdk.OneDec()) {
		return w.StakingKeeper.Slash(ctx, consAddr, infractionHeight, power, slashFactor)
	}
	rawPower := sdk.NewDec(power).Quo(m).TruncateInt64()
	return w.StakingKeeper.Slash(ctx, consAddr, infractionHeight, rawPower, slashFactor)
}

// SlashWithInfractionReason mirrors Slash for the v0.47 SDK API that includes
// the infraction reason.
func (w WrappedStakingKeeper) SlashWithInfractionReason(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec, infraction stakingtypes.Infraction) math.Int {
	val := w.StakingKeeper.ValidatorByConsAddr(ctx, consAddr)
	if val == nil {
		return w.StakingKeeper.SlashWithInfractionReason(ctx, consAddr, infractionHeight, power, slashFactor, infraction)
	}
	op := val.GetOperator()
	if !w.poa.IsAuthority(ctx, op) {
		return w.StakingKeeper.SlashWithInfractionReason(ctx, consAddr, infractionHeight, power, slashFactor, infraction)
	}
	m, snapFound := w.poa.MultiplierForValidatorWithStatus(ctx, op, infractionHeight)
	if !snapFound {
		ctx.Logger().Error("poa: missing multiplier snapshot for authority slash; SKIPPING slash",
			"operator", op.String(), "infraction_height", infractionHeight, "infraction", infraction)
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeSlashSkippedNoSnapshot,
			sdk.NewAttribute(types.AttributeOperator, op.String()),
			sdk.NewAttribute(types.AttributeInfractionHeight, sdk.NewInt(infractionHeight).String()),
		))
		return math.ZeroInt()
	}
	if m.LTE(sdk.OneDec()) {
		return w.StakingKeeper.SlashWithInfractionReason(ctx, consAddr, infractionHeight, power, slashFactor, infraction)
	}
	rawPower := sdk.NewDec(power).Quo(m).TruncateInt64()
	return w.StakingKeeper.SlashWithInfractionReason(ctx, consAddr, infractionHeight, rawPower, slashFactor, infraction)
}
