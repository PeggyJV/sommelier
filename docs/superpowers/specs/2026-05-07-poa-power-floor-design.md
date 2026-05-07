# Sommelier Proof-of-Authority Power Floor — Design

**Date:** 2026-05-07
**Author:** Zaki (with Claude)
**Status:** Implemented in PR [#340](https://github.com/PeggyJV/sommelier/pull/340)

## 1. Summary

Convert Sommelier from a permissionless proof-of-stake chain to a hybrid proof-of-authority chain in which a binary-specified set of "authority" validators always controls **at least 67% of consensus voting power**. Community members may continue to register validators and accept delegations, but their combined voting power is capped at <33% of consensus. Stake economics (delegations, rewards, unbonding) remain functionally unchanged.

Authority is enforced not by replacing `x/staking`, but by a new `x/poa` module that:

1. Maintains an allowlist of authority validators (compiled-in default; mutable by governance).
2. Wraps the staking keeper to expose **rescaled** validator power to every downstream consumer (consensus updates, gravity-bridge, slashing, distribution, cork, pubsub, etc.).
3. Normalizes power back to raw stake at the slashing boundary so authority validators are not over-slashed.

## 2. Goals & non-goals

### Goals

- Authority validators collectively hold ≥67% of CometBFT voting power every block.
- Gravity-bridge orchestrator signing weights match consensus weights (no security divergence).
- Community validators may join, be delegated to, and earn rewards proportional to their actual stake.
- Authority allowlist is changeable via governance (`MsgUpdateAuthoritySet`) without a binary upgrade.
- Slashing (downtime, double-sign) penalizes authority validators based on their **actual** bonded stake, not their boosted power.
- The chain halts cleanly if zero authority validators are bonded-and-unjailed (correct PoA failure mode).

### Non-goals

- Removing community delegations or community validators.
- Replacing `x/staking`, `x/distribution`, `x/slashing`.
- Changing inflation, fee, or reward mechanics.
- Implementing per-validator hardcoded weights inside the authority set (authority validators' relative power tracks their actual stake; only the *aggregate* binary share is fixed at ≥67%).

## 3. Architecture

```
                       ┌────────────────────────────────────┐
                       │  CometBFT (consumes ValidatorUpdates)
                       └────────────────▲───────────────────┘
                                        │ merged (boosted) updates
                                        │
   ┌──────────────────────┐         ┌───┴────────────────────┐
   │  x/staking           │         │  x/poa AppModule       │
   │  AppModule (no-op    │         │  EndBlock              │
   │  EndBlock; PoA owns) │◀────────│   1. invoke staking    │
   │                      │ closure │      EndBlocker        │
   └─────────▲────────────┘         │   2. compute M         │
             │                      │   3. overwrite         │
             │ embedded             │      LastValidatorPower│
             │                      │   4. merge updates     │
   ┌─────────┴────────────┐         └───────────▲────────────┘
   │ poakeeper.Keeper     │                     │ keeper interface
   │  (raw concrete + sk) │                     │
   └─────────▲────────────┘         ┌───────────┴────────────┐
             │ exposes              │  Gravity, Slashing,    │
             │                      │  Distribution, Cork,   │
   ┌─────────┴────────────┐         │  Pubsub, Axelarcork,   │
   │ WrappedStakingKeeper │◀────────│  Incentives, Evidence  │
   │  - boosted reads     │         └────────────────────────┘
   │  - normalised Slash  │
   └──────────────────────┘
```

### 3.1 New module: `x/poa`

**State:**

- `AuthoritySet`: ordered list of `sdk.ValAddress` (operator addresses) recognised as authority.
- `Params`:
  - `floor_fraction` (`sdk.Dec`, default `"0.670000000000000001"`) — minimum aggregate share required for the authority set.
  - `enabled` (`bool`, default `true`) — feature flag for emergency disable.
  - `halt_when_authority_empty` (`bool`, default `true`) — panic in EndBlocker when no bonded-and-unjailed authority validator exists (see §3.5).

**Messages (gov-only authority):**

- `MsgUpdateAuthoritySet { authority: gov, validators: []ValAddress }` — replace the allowlist.
- `MsgUpdateParams { authority: gov, params: Params }` — update floor or toggle enabled.

**Queries:**

- `AuthoritySet`, `Params`, `EffectivePower(operator)` — last computed (boosted or raw) power.

**Genesis / initial population:** seeded from a compiled-in default `DefaultAuthorityValidators` slice in the binary. Initial seeding happens via the v10 upgrade handler (see §5).

### 3.2 Keeper wrapper

`x/poa` exposes a `StakingKeeper`-shaped wrapper that satisfies the union of all interfaces consumed today:
`stakingtypes.StakingKeeper` (used by slashing, distribution), the gravity-bridge `StakingKeeper`, and Sommelier's per-module `StakingKeeper` interfaces (cork, pubsub, axelarcork, incentives).

**Read methods (apply multiplier to authority validators):**

- `GetLastValidatorPower(ctx, operator) → int64` — pass-through; boost is realised via the EndBlocker's `LastValidatorPower` overwrite.
- `GetLastTotalPower(ctx) → math.Int` — recomputes the boosted total by re-iterating overwritten `LastValidatorPower` (the staking module's own `LastTotalPower` slot is set inside `ApplyAndReturnValidatorSetUpdates` *before* PoA's overwrite, so it reflects the pre-rescale total).
- `GetBondedValidatorsByPower(ctx) → []Validator`, `GetAllValidators(ctx) → []Validator`, `GetValidator(...)` — return concrete validators with `Tokens` field rescaled (`Ceil(rawTokens * M)`).
- `IterateBondedValidatorsByPower`, `IterateLastValidators`, `IterateValidators` — yield a `boostedValidator` adapter overriding `GetTokens` / `GetBondedTokens` / `GetConsensusPower`. `IterateLastValidatorPowers` is pass-through (the underlying store already holds boosted values).
- `Validator(addr)`, `ValidatorByConsAddr(addr)` — return the `boostedValidator` adapter.

**Pass-through methods (unchanged):**

- `GetParams`, `Delegation`, `MaxValidators`, `IsValidatorJailed`, `Jail`, `Unjail`, `PowerReduction`, `BondDenom`, `UnbondingTime`, `ValidatorQueueIterator`, `Hooks`, `SetLastValidatorPower`, `DeleteLastValidatorPower`, `IterateDelegations`, `GetAllSDKDelegations`, `GetAllDelegatorDelegations`.

The adapter intentionally does NOT override `TokensFromShares*` / `SharesFromTokens*`: delegation share allocation must operate on raw tokens so boosted authority validators do not dilute community delegators' shares. Boost is a property of consensus power and slashing exposure, not a property of token ownership.

**Normalised-write method:**

- `Slash(ctx, consAddr, infractionHeight, power, slashFactor) math.Int`
  - If `consAddr` belongs to an authority validator, divide `power` by that validator's effective multiplier captured at `infractionHeight` (see §3.4) before delegating to `stakingKeeper.Slash`.
  - Else pass through unchanged.

### 3.3 Power rescaling math

Let:

- `B = Σ raw_power(v)` for v in authority set, bonded & unjailed
- `C = Σ raw_power(v)` for v in community set, bonded & unjailed
- `f = floor_fraction` (≈ 0.67)

Required scaling factor for the authority bucket:

```
M = (f / (1 - f)) * (C / B)     if B > 0
```

If `M ≤ 1` (authority already controls ≥ f), the multiplier is clamped to 1 (no boost needed). Otherwise each authority validator `v` reports power `floor(raw_power(v) * M)`. Community powers are unchanged.

**Edge cases:**

- `B = 0` (no bonded authority validators): module emits a critical-level event `AuthoritySetUnavailable` and returns the staking-module's raw updates verbatim. CometBFT will continue producing blocks based on community alone — *but* in practice gravity-bridge cellar txs and slashing semantics depend on the invariant; ops should treat this as an immediate halt-worthy alert. (See §3.5 for halt option.)
- `C = 0`: no rescale needed; authority already at 100%.
- Rounding: ceil on each authority power so the post-rescale authority share is ≥ f even after integer truncation. The wrapper applies the same Ceil semantics to `GetTokens` / `GetConsensusPower` to keep reads consistent with the consensus power written to the store.

### 3.4 EndBlocker

PoA owns the only EndBlocker that returns validator updates for the chain. SDK v0.47's `module.Manager.EndBlock` panics if more than one registered module returns a non-empty `[]abci.ValidatorUpdate`. To avoid this, `x/staking` is registered in the module manager (so its `InitGenesis` / `BeginBlock` / `ExportGenesis` continue to run) but its `AppModule` is wrapped (`app/staking_endblocker_noop.go`) to make `EndBlock` a no-op. PoA's `AppModule` instead receives a closure over the production `*stakingkeeper.Keeper` and invokes `staking.EndBlocker` exactly once per block from inside its own EndBlocker.

App ordering (relevant slice):

```
... → crisis → gov → poa → staking(no-op) → ica → ... → community modules → ...
```

Per block, PoA's EndBlocker:

1. Calls the staking-end-blocker closure to advance unbonding queues and produce raw `abci.ValidatorUpdate` against the just-written `staking.LastValidatorPower`.
2. Iterates `LastValidatorPower`, partitioning bonded-and-unjailed validators into authority (A) and community (C) buckets.
3. Computes `M = floor / (1 - floor) * C / B`, clamped to 1 when authority already exceeds the floor.
4. For each authority validator, computes `Ceil(rawPower * M)` and:
   a. Overwrites `staking.LastValidatorPower[v]` with the boosted value.
   b. Records the per-validator multiplier in the block's snapshot.
5. Merges raw updates with boosted entries (`mergeUpdatesWithBoost`): rewrites `Power` on raw entries that match a boosted authority pubkey and appends explicit entries for boosted authorities not present in raw. Returns the merged slice as the chain's only `[]abci.ValidatorUpdate` for the block.
6. Persists `(block_height → per_authority_multiplier_snapshot)` to PoA state. Retention covers `max(unbonding_blocks, evidence.MaxAgeNumBlocks) + slashing.SignedBlocksWindow`; the unbonding-blocks term uses a conservative lower bound on block time so retention errs toward over-keeping (missing snapshots cause authority slashes to be skipped, so over-retention is the safe direction).
7. Emits the `authority_rescale` telemetry event with `multiplier`, `authority_power`, `community_power`.

Reads from the wrapper layer (§3.2) provide the durable surface for the floor invariant. The `LastValidatorPower` overwrite is belt-and-suspenders: any consumer that queries staking directly via `GetLastValidatorPower` also sees the boosted value, and `WrappedStakingKeeper.GetLastTotalPower` recomputes the boosted aggregate by re-iterating the store (the staking module's `LastTotalPower` slot is set inside `ApplyAndReturnValidatorSetUpdates` *before* PoA's overwrite, so it reflects the pre-rescale total).

### 3.5 Halt-on-empty-authority

Add a param `halt_when_authority_empty` (default `true`). When the bonded-and-unjailed authority set is empty in EndBlocker, panic with a descriptive message. This is the correct PoA failure mode: the security guarantee is broken, do not produce more blocks.

## 4. Module wiring changes (`app/app.go`)

For each consumer currently passed `app.StakingKeeper`, replace with `app.PoaKeeper.WrappedStakingKeeper()`:

- `slashingkeeper.NewKeeper(... stakingKeeper ...)` → `app.PoaKeeper.WrappedStakingKeeper()`
- `distrkeeper.NewKeeper(...)` → wrapped
- `evidencekeeper.NewKeeper(...)` → wrapped
- `gravitykeeper.NewKeeper(...)` → wrapped
- `corkkeeper.NewKeeper(...)` → wrapped
- `axelarcorkkeeper.NewKeeper(...)` → wrapped
- `pubsubkeeper.NewKeeper(...)` → wrapped
- `incentiveskeeper.NewKeeper(...)` → wrapped
- Staking hooks (distribution, slashing) — registered against the **real** staking keeper (unchanged); hooks fire on real stake changes which is correct behavior.

`x/staking` itself continues to receive the real keeper for its internal bookkeeping.

EndBlocker order: `poatypes.ModuleName` immediately after `govtypes.ModuleName` and before `stakingtypes.ModuleName` in `mm.SetOrderEndBlockers`. `stakingtypes.ModuleName` remains in the slice (SDK requires every registered module to appear) but its EndBlock is the no-op variant from `app/staking_endblocker_noop.go`; PoA invokes `staking.EndBlocker` itself.

InitGenesis order: `poatypes.ModuleName` after `staking` and `slashing`, before `gravity`.

## 5. Migration / upgrade

Ship as a `v10` upgrade (`app/upgrades/v10`). Upgrade handler:

1. Adds `poatypes.StoreKey` to the `StoreUpgrades.Added` slice.
2. Initialises PoA params with the default floor.
3. Writes the initial `AuthoritySet` from a compiled-in slice `DefaultAuthorityValidators` declared in `app/upgrades/v10/constants.go`. (Operator: provide the list before tagging the release.)
4. Runs the standard module-version migration map.

Rollback note: removing `x/poa` would require a follow-up upgrade restoring all wired consumers to the raw staking keeper.

## 6. Testing strategy

### Unit

- `keeper/multiplier_test.go`: math correctness (rounding, clamp at 1, B=0, C=0, mixed jailed states).
- `keeper/wrapper_test.go`: each wrapped read method returns boosted power for authority and raw for community; pass-through methods unchanged.
- `keeper/slash_test.go`: `Slash` on authority validator with stale infraction height uses the snapshot multiplier.

### Integration

- `x/poa` simapp test: bond 3 authority + 5 community validators with mixed stake; run 10 blocks; assert `staking.GetLastTotalPower` reports authority share ≥ 67% every block.
- Gravity-bridge validator-set update (`SignerSetTx`): authority signers' weights sum ≥ 67% of total.
- Slashing: double-sign by an authority validator burns tokens equal to `slashFraction × actual_bonded`, not boosted.
- Downtime jail of all authority validators ⇒ chain halt panic (with `halt_when_authority_empty=true`).

### End-to-end (`integration_tests/`)

- New scenario: community delegator delegates large stake to a community validator; assert each block's community share `< (1 - floor_fraction)` (deterministic against the configured floor param).
- Governance proposal to add/remove an authority validator; assert allowlist update lands in the next block and rescaling reflects it.

## 7. Risks & open questions

| # | Risk | Mitigation |
|---|------|------------|
| 1 | A wrapped consumer slipped through and reads raw staking. | Grep audit in plan; add a CI lint that flags `app.StakingKeeper` references outside `app.go` and `x/staking` wiring. |
| 2 | Slashing math wrong for historical infractions when authority set changes between infraction and slash. | Multiplier snapshot map keyed by `infraction_height`, retained per §3.4. |
| 3 | Staking's diff-based ABCI emission fights with our overwrite. | Document and integration-test the BaseApp last-writer-wins ordering; assert deterministic output in simapp tests. |
| 4 | IBC light client assumptions about validator set churn. | Authority rescale changes voting power abruptly when set changes; document for relayer operators. Mitigated because authority set changes are infrequent (gov-mediated). |
| 5 | A community validator's stake grows large enough that even with M=1 it exceeds 33%. | Cannot happen by construction: when authority set is healthy, `M` is chosen so authority is ≥67%, hence community ≤33%. The only failure is `B=0`. |
| 6 | Multiplier snapshot retention bloats state. | Prune snapshots older than `UnbondingTime + SignedBlocksWindow`. |

## 8. Out-of-scope follow-ups

- Per-authority-validator weight overrides (e.g., "ECC always 30%, others split remainder").
- Authority validators exempt from inflation/MEV/slashing changes.
- A UI in the explorer showing raw vs. effective power side by side.

---

**Approvals required before implementation:**

- [ ] Zaki: confirms approach A and §3 design.
- [ ] Spec reviewer (automated).
