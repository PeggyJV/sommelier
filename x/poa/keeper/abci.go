package keeper

import (
	"sort"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	tmprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/peggyjv/sommelier/v10/x/poa/types"
)

// StakingEndBlockerFn is the callback PoA uses to drive staking's own
// EndBlocker exactly once per block. The wrapping is necessary because PoA
// owns staking's slot in the EndBlocker order: x/staking is NOT registered as
// a standalone EndBlocker (SDK 0.47's module manager would otherwise panic if
// two modules both returned non-empty validator updates).
type StakingEndBlockerFn func(sdk.Context) []abci.ValidatorUpdate

// EndBlocker runs the PoA per-block logic.
//
//  1. Run staking's EndBlocker to advance unbonding queues and emit raw
//     ValidatorUpdates against the just-written staking.LastValidatorPower.
//  2. Partition bonded-and-unjailed validators into authority and community
//     buckets, compute the boost multiplier M, and store a snapshot.
//  3. Overwrite staking.LastValidatorPower for authority validators with the
//     boosted value (so any subsequent reads see the rescaled power).
//  4. Merge raw updates with boosted entries and return the combined slice.
func EndBlocker(ctx sdk.Context, k Keeper, runStakingEndBlocker StakingEndBlockerFn) []abci.ValidatorUpdate {
	rawUpdates := runStakingEndBlocker(ctx)

	params := k.GetParams(ctx)
	if !params.Enabled {
		// Module disabled: no PoA enforcement, so safe mode (if it was set) no
		// longer applies — clear it and let consumers resume.
		k.exitSafeModeIfActive(ctx)
		// Record an empty snapshot at this height so a future slash for an
		// authority validator at this height does not hit the
		// missing-snapshot refuse path in WrappedStakingKeeper.Slash.
		k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{Height: ctx.BlockHeight()})
		k.pruneSnapshots(ctx)
		return rawUpdates
	}

	authSet := authoritySetMap(k.GetAuthoritySet(ctx))

	// Gather raw bonded powers from the just-written LastValidatorPower.
	var (
		authPower = math.ZeroInt()
		comPower  = math.ZeroInt()
		authVals  []sdk.ValAddress
	)
	k.sk.IterateLastValidatorPowers(ctx, func(op sdk.ValAddress, power int64) bool {
		v, found := k.sk.GetValidator(ctx, op)
		if !found || v.Jailed || !v.IsBonded() {
			return false
		}
		if _, isAuth := authSet[op.String()]; isAuth {
			authPower = authPower.Add(math.NewInt(power))
			authVals = append(authVals, op)
		} else {
			comPower = comPower.Add(math.NewInt(power))
		}
		return false
	})

	if authPower.IsZero() {
		if params.HaltWhenAuthorityEmpty {
			// Fail-closed: refuse to produce further blocks.
			panic(types.ErrNoBondedAuthority)
		}
		// Safe mode (Option A): keep producing blocks on community stake so
		// governance can re-seed the authority set on-chain, but flag safe mode
		// so the value-bearing modules (gravity/cork/axelarcork) freeze. No
		// boost is applied; record an empty snapshot so slashing at this height
		// passes through un-boosted.
		k.enterSafeMode(ctx)
		k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{Height: ctx.BlockHeight()})
		k.pruneSnapshots(ctx)
		return rawUpdates
	}

	// Authority set is healthy. Boosting resumes immediately below, but if we
	// were in safe mode the value-bearing freeze is held until the restored set
	// is actually securing consensus (CometBFT validator-update delay), not the
	// instant authority re-bonds — see maybeThawSafeMode.
	k.maybeThawSafeMode(ctx)

	m := types.ComputeMultiplier(authPower, comPower, params.FloorFraction)
	if m.LTE(sdk.OneDec()) {
		// No boost needed. Record an empty snapshot anyway so Slash math at
		// this height has a found-snapshot.
		k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{Height: ctx.BlockHeight()})
		k.pruneSnapshots(ctx)
		return rawUpdates
	}

	// Build snapshot and overwrite LastValidatorPower for authority validators.
	snap := types.MultiplierSnapshot{Height: ctx.BlockHeight()}
	boostedByOp := make(map[string]int64, len(authVals))
	for _, op := range authVals {
		rawPower := k.sk.GetLastValidatorPower(ctx, op)
		// Ceiling, not floor: floor() can drop the post-rescale share just
		// below the configured floor due to integer truncation.
		boosted := sdk.NewDec(rawPower).Mul(m).Ceil().TruncateInt64()
		if boosted == rawPower {
			continue
		}
		k.sk.SetLastValidatorPower(ctx, op, boosted)
		boostedByOp[op.String()] = boosted
		snap.Entries = append(snap.Entries, &types.MultiplierEntry{
			OperatorAddress: op.String(),
			Multiplier:      m.String(),
		})
	}
	k.SetMultiplierSnapshot(ctx, snap)
	k.pruneSnapshots(ctx)

	finalUpdates := mergeUpdatesWithBoost(ctx, k, rawUpdates, boostedByOp)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeAuthorityRescale,
		sdk.NewAttribute(types.AttributeMultiplier, m.String()),
		sdk.NewAttribute(types.AttributeAuthorityPower, authPower.String()),
		sdk.NewAttribute(types.AttributeCommunityPower, comPower.String()),
	))

	return finalUpdates
}

// authoritySetMap is a O(1) lookup keyed by operator bech32.
func authoritySetMap(vals []sdk.ValAddress) map[string]struct{} {
	m := make(map[string]struct{}, len(vals))
	for _, v := range vals {
		m[v.String()] = struct{}{}
	}
	return m
}

// mergeUpdatesWithBoost rewrites Power on raw updates that correspond to
// boosted authority entries, and appends explicit entries for authority
// validators whose power changed but were absent from raw.
func mergeUpdatesWithBoost(ctx sdk.Context, k Keeper, raw []abci.ValidatorUpdate, boostedByOp map[string]int64) []abci.ValidatorUpdate {
	if len(boostedByOp) == 0 {
		return raw
	}

	// Build a pubkey -> operator lookup for the boosted set so we can match
	// against entries in `raw` by pubkey bytes.
	pkToOp := make(map[string]string, len(boostedByOp))
	opPk := make(map[string]abci.ValidatorUpdate, len(boostedByOp))
	for opStr := range boostedByOp {
		op, err := sdk.ValAddressFromBech32(opStr)
		if err != nil {
			continue
		}
		v, found := k.sk.GetValidator(ctx, op)
		if !found {
			continue
		}
		pk, err := v.TmConsPublicKey()
		if err != nil {
			continue
		}
		key, ok := consPubKeyIndex(pk)
		if !ok {
			// Unrecognised key type. Skip rather than indexing it: an empty
			// key would collide every such validator onto one map slot, and a
			// raw update with an equally unrecognised key would then match it
			// and be assigned an unrelated validator's boosted power.
			continue
		}
		pkToOp[key] = opStr
		opPk[opStr] = abci.ValidatorUpdate{PubKey: pk, Power: boostedByOp[opStr]}
	}

	out := make([]abci.ValidatorUpdate, 0, len(raw)+len(boostedByOp))
	seen := make(map[string]struct{}, len(boostedByOp))

	for _, u := range raw {
		if key, ok := consPubKeyIndex(u.PubKey); ok {
			if op, matched := pkToOp[key]; matched {
				// Override with boosted value.
				u.Power = boostedByOp[op]
				seen[op] = struct{}{}
			}
		}
		out = append(out, u)
	}

	// Append boosted entries for authority validators not present in raw, in
	// deterministic operator-address order: Go map iteration is randomized and
	// the ABCI ValidatorUpdates slice must be byte-identical across all nodes.
	missing := make([]string, 0, len(opPk))
	for op := range opPk {
		if _, already := seen[op]; already {
			continue
		}
		missing = append(missing, op)
	}
	sort.Strings(missing)
	for _, op := range missing {
		out = append(out, opPk[op])
	}
	return out
}

// consPubKeyIndex builds a map key that identifies a consensus public key by
// both type and bytes, and reports whether the type is one we recognise.
//
// Keying on bytes alone would let an unrecognised type collapse to the empty
// string, colliding distinct validators onto one slot; including the type tag
// also stops an ed25519 key from ever aliasing a secp256k1 key.
func consPubKeyIndex(pk tmprotocrypto.PublicKey) (string, bool) {
	if b := pk.GetEd25519(); len(b) > 0 {
		return "ed25519:" + string(b), true
	}
	if b := pk.GetSecp256K1(); len(b) > 0 {
		return "secp256k1:" + string(b), true
	}
	return "", false
}

// pruneSnapshots deletes per-block multiplier snapshots older than
// max(unbonding_blocks, evidence_max_age_blocks) + signed_blocks_window.
// The retention covers both unbonding-driven slashing and evidence-based
// slashing whose infraction height can lag by `evidence.max_age_num_blocks`.
func (k Keeper) pruneSnapshots(ctx sdk.Context) {
	unbonding := k.sk.UnbondingTime(ctx)
	// Round UP and add a block so the estimate never lands short.
	blockNanos := k.estimateBlockNanos(ctx)
	unbondingBlocks := (unbonding.Nanoseconds()+blockNanos-1)/blockNanos + 1

	var evidenceBlocks int64
	if cp := ctx.ConsensusParams(); cp != nil && cp.Evidence != nil {
		evidenceBlocks = cp.Evidence.MaxAgeNumBlocks
	}
	longest := unbondingBlocks
	if evidenceBlocks > longest {
		longest = evidenceBlocks
	}

	var signedWindow int64
	if k.slashingKeeper != nil {
		signedWindow = k.slashingKeeper.SignedBlocksWindow(ctx)
	}
	retention := longest + signedWindow
	if ctx.BlockHeight() > retention {
		k.PruneSnapshotsBefore(ctx, ctx.BlockHeight()-retention)
	}
}

const (
	// fallbackBlockNanos is used until enough blocks have elapsed since
	// activation to measure the real rate. Deliberately small (1s) so the
	// fallback OVER-retains: the failure mode of under-retention is that
	// authority slashes are silently skipped (see rawSlashPower's refuse path),
	// which is far worse than keeping extra snapshots.
	fallbackBlockNanos int64 = 1_000_000_000
	// minBlockNanos floors the measured rate so a clock anomaly (or a chain
	// that genuinely runs sub-100ms) cannot collapse retention to nothing.
	minBlockNanos int64 = 100_000_000
	// blockRateSampleBlocks is how many blocks must elapse post-activation
	// before the measured rate is trusted over the fallback.
	blockRateSampleBlocks int64 = 1000
	// blockRateSafetyNumer/Denom shade the measured rate DOWN by 20%, so the
	// derived block count is ~25% higher than measured. Block times shorten
	// under load; retention must survive that without a migration.
	blockRateSafetyNumer int64 = 4
	blockRateSafetyDenom int64 = 5
)

// estimateBlockNanos returns the chain's average block time in nanoseconds,
// measured from the activation height/time to now, shaded down for safety.
//
// This replaces a compiled-in block-time constant. The constant's "conservative
// lower bound" argument only held if the assumed value really was a lower
// bound: if actual block times ever fell below it (load, a consensus param
// change, a chain relaunch), retention would silently under-cover the slashable
// window and authority slashes would start being refused with no signal beyond
// a log line.
func (k Keeper) estimateBlockNanos(ctx sdk.Context) int64 {
	activationHeight, okH := k.GetActivationHeight(ctx)
	activationTime, okT := k.GetActivationTime(ctx)
	if !okH || !okT {
		return fallbackBlockNanos
	}

	elapsedBlocks := ctx.BlockHeight() - activationHeight
	elapsedNanos := ctx.BlockTime().UnixNano() - activationTime
	if elapsedBlocks < blockRateSampleBlocks || elapsedNanos <= 0 {
		return fallbackBlockNanos
	}

	measured := elapsedNanos / elapsedBlocks
	shaded := measured * blockRateSafetyNumer / blockRateSafetyDenom
	if shaded < minBlockNanos {
		return minBlockNanos
	}
	return shaded
}
