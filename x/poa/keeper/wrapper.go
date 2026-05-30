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

// currentMultiplier returns the boost multiplier in effect for `op` for reads
// during this block, or 1 if no boost applies.
//
// The boost is sourced purely from the latest committed snapshot (height <=
// current), NOT from the live authority allowlist. This keeps every boosted
// read consistent with the LastValidatorPower the staking store holds: both
// reflect the last EndBlocker's decision and ignore mid-block authority-set
// mutations, which only take effect when the next EndBlocker recomputes powers.
// Gating on the live allowlist instead would (a) de-boost a validator the
// instant it is removed, before its consensus power actually changes at
// EndBlock, and (b) return no boost during BeginBlock/DeliverTx because the
// current height's snapshot is not written until that height's EndBlocker.
func (w WrappedStakingKeeper) currentMultiplier(ctx sdk.Context, op sdk.ValAddress) sdk.Dec {
	m, found := w.poa.LatestMultiplierForValidator(ctx, op, ctx.BlockHeight())
	if !found || m.LTE(sdk.OneDec()) {
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

// rawSlashPower converts caller-supplied (boosted) consensus power back to the
// raw stake that should actually be slashed for operator `op`, using the boost
// recorded at the infraction height. It returns refuse=true when the slash must
// be skipped to avoid over-slashing.
//
// Authority/boost status is evaluated at the infraction height via the
// snapshot, NEVER against the current allowlist: the live set may differ from
// the set in effect when the infraction occurred (e.g. delayed evidence after
// an authority-set change), and gating on current membership would over- or
// under-slash.
//
// Missing-snapshot handling keys off the PoA activation height:
//   - infractionHeight < activation (or activation unset): a pre-PoA height
//     where no boost was ever applied — slash the raw power unchanged.
//   - infractionHeight >= activation: an empty snapshot is written every block
//     and retention covers the full slashable window, so a gap here means
//     corruption — refuse rather than risk over-slashing on boosted power.
func (w WrappedStakingKeeper) rawSlashPower(ctx sdk.Context, op sdk.ValAddress, infractionHeight, power int64) (int64, bool) {
	m, snapFound := w.poa.MultiplierForValidatorWithStatus(ctx, op, infractionHeight)
	if snapFound {
		if m.LTE(sdk.OneDec()) {
			return power, false
		}
		return sdk.NewDec(power).Quo(m).TruncateInt64(), false
	}
	activation, ok := w.poa.GetActivationHeight(ctx)
	if !ok || infractionHeight < activation {
		return power, false
	}
	return 0, true
}

// emitSlashSkipped logs and emits the event for a refused slash.
func (w WrappedStakingKeeper) emitSlashSkipped(ctx sdk.Context, op sdk.ValAddress, infractionHeight int64) {
	ctx.Logger().Error("poa: missing multiplier snapshot for post-activation authority slash; SKIPPING slash",
		"operator", op.String(), "infraction_height", infractionHeight)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeSlashSkippedNoSnapshot,
		sdk.NewAttribute(types.AttributeOperator, op.String()),
		sdk.NewAttribute(types.AttributeInfractionHeight, sdk.NewInt(infractionHeight).String()),
	))
}

// Slash normalizes boosted consensus power back to raw stake before delegating
// to the underlying keeper. See rawSlashPower for the snapshot/activation rules.
func (w WrappedStakingKeeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) math.Int {
	val := w.StakingKeeper.ValidatorByConsAddr(ctx, consAddr) // raw lookup, NOT adapted
	if val == nil {
		return w.StakingKeeper.Slash(ctx, consAddr, infractionHeight, power, slashFactor)
	}
	op := val.GetOperator()
	rawPower, refuse := w.rawSlashPower(ctx, op, infractionHeight, power)
	if refuse {
		w.emitSlashSkipped(ctx, op, infractionHeight)
		return math.ZeroInt()
	}
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
	rawPower, refuse := w.rawSlashPower(ctx, op, infractionHeight, power)
	if refuse {
		w.emitSlashSkipped(ctx, op, infractionHeight)
		return math.ZeroInt()
	}
	return w.StakingKeeper.SlashWithInfractionReason(ctx, consAddr, infractionHeight, rawPower, slashFactor, infraction)
}
