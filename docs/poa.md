# Sommelier Proof-of-Authority Power Floor

This document describes the PoA semantics introduced in the **v10** upgrade.
Operator-facing reference; protocol details are in
[`docs/superpowers/specs/2026-05-07-poa-power-floor-design.md`](superpowers/specs/2026-05-07-poa-power-floor-design.md).

## What changed

A new on-chain module, `x/poa`, guarantees that a binary-specified set of
"authority" validators always controls **at least 67%** of consensus voting
power. Community members may continue to register validators, accept
delegations, and earn block rewards on their actual stake — but their
combined consensus power is capped at less than 33%.

The 67% floor is enforced by rescaling validator power in the EndBlocker:
authority validators report a *boosted* `ConsensusPower` to CometBFT and to
every internal consumer of validator power (gravity-bridge, slashing,
distribution, evidence, cork, axelarcork, pubsub, incentives).

## Configuration

Module params (`x/poa/types/params.go`):

| Param | Default | Effect |
|---|---|---|
| `floor_fraction` | `0.670000000000000001` | Minimum authority share. The trailing decimal cushion guarantees a strict supermajority after integer rounding. |
| `enabled` | `true` | When false, EndBlocker is a no-op. Emergency disable. |
| `halt_when_authority_empty` | `false` | Behavior when no authority validator is bonded and unjailed. `false` (default) → **safe mode**: keep producing blocks but freeze the value-bearing modules (see [Safe mode](#safe-mode-authority-empty)). `true` → halt the chain via panic. |

Authority allowlist:
- Seeded at the v10 upgrade from `app/upgrades/v10/constants.go::DefaultAuthorityValidators`.
- Mutable post-upgrade via `MsgUpdateAuthoritySet` (gov-only).

Params can be updated via `MsgUpdateParams` (gov-only).

## How a community delegator's stake translates to consensus power

Power is rescaled per block. If the authority bucket holds raw consensus
power `B` and the community bucket holds raw consensus power `C`, the
multiplier applied to each authority validator is:

```
M = floor / (1 - floor) * C / B          when authority is below floor
M = 1                                    when authority is already at or above floor
```

Worked example with `floor = 0.67`, three authority validators with raw
powers `100, 100, 100` (B=300) and one community validator with raw power
`900` (C=900):

```
share_before = 300 / 1200 = 0.25         (below floor)
M            = 0.67/0.33 * 900/300       ≈ 6.0909
authority_after = 100 * 6.0909           ≈ 609 each (boosted)
total_after  = 3*609 + 900 = 2727
share_after  = 1827 / 2727               ≈ 0.67   (at floor)
```

A community delegator who delegates a large stake to a community validator
will see their validator's *delegation share* allocations behave normally
(token-share math is computed on raw stake), but the validator's
`ConsensusPower` will not exceed the 33% community ceiling.

## Slashing semantics

Authority validators are penalised on **raw stake**, not on their boosted
consensus power. When CometBFT reports an infraction:

1. The wrapper's `Slash` looks up the multiplier snapshot at the infraction
   height.
2. It divides the caller-supplied (boosted) power by that multiplier to
   recover raw stake.
3. The underlying staking keeper is then slashed on the raw value.

This means a 5% downtime fraction on an authority validator burns 5% of
their *actual* bonded tokens, not 5% of their boosted CometBFT power.

If a multiplier snapshot is missing for an authority infraction (for
example after state-sync if the snapshot was pruned before the evidence
arrived), the slash is **refused** rather than executed on possibly-boosted
input. A `slash_skipped_no_snapshot` event is emitted.

## Reward semantics

Block rewards are allocated by `x/distribution` from CometBFT's
`LastCommitInfo`, which carries the **boosted** voting power. Authority
validators therefore receive boosted gross rewards proportional to their
inflated consensus weight.

This is the deliberate incentive alignment of the PoA design: authority
validators bear the security responsibility of the chain (and the
operational requirements that come with it) and earn rewards proportional
to that responsibility. Slashing exposure is on raw stake, but reward
accrual is on boosted weight.

Community validators earn rewards proportional to their unboosted
consensus share.

## Slashing semantics

Because authority validators report *boosted* consensus power, a naive slash
would burn `slash_factor × boosted_power`, i.e. up to `M×` their real stake —
an over-slash. To prevent this, `x/poa` wraps the staking keeper and, at slash
time, converts the caller-supplied (boosted) power back to raw stake using the
boost multiplier recorded **at the infraction height**, then delegates to the
underlying keeper. Net effect: an authority validator is penalised on its
**real stake**, not on its boosted consensus power.

Key rules:
- **Infraction-height membership, not current.** Whether a validator was
  boosted is read from the per-height snapshot at the infraction height, never
  from the live allowlist. Evidence can arrive after an authority-set change;
  using current membership would over- or under-slash.
- **Per-block snapshots.** Every block writes a multiplier snapshot (an *empty*
  one when no boost is applied), so for any height at or after activation a
  snapshot is always present. Snapshots are retained to cover the full slashable
  window: `max(unbonding_blocks, evidence_max_age_blocks) + signed_blocks_window`
  (computed with a conservative lower-bound block time).
- **Activation height.** The height at which PoA went live (recorded at the v10
  upgrade / genesis init) is the boundary for missing-snapshot handling:
  - infraction height **below** activation → pre-PoA, no boost was ever applied
    → slash passes through against raw power.
  - infraction height **at or above** activation with a missing snapshot →
    treated as corruption → slash is **refused** (`slash_skipped_no_snapshot`)
    rather than risk over-slashing.
- **Tombstoning is the real deterrent.** Double-sign tombstoning and jailing
  happen in the evidence/slashing modules independent of the burn amount, so a
  misbehaving authority validator is still permanently ejected from consensus
  even when the token penalty is small or refused.

### Snapshots and chain restarts

Multiplier snapshots are **not** included in `ExportGenesis` (exporting the full
retained window would bloat genesis). On an export/restart the activation height
is re-recorded at the new initial height, so any infraction height before it —
including all pre-restart heights — passes through as un-boosted. This is safe
because CometBFT evidence cannot reference pre-restart heights after a genesis
restart.

## Economic security considerations

The slash normalisation above is correct accounting, but it has a deliberate
consequence operators must understand: **slashing is a weak economic deterrent
for boosted authority validators.** A validator wielding 67% of consensus power
via boost may have only a small real self-stake, so `slash_factor × real_stake`
is a small absolute penalty relative to the consensus damage it could do. The
same applies to rewards — distribution allocates by boosted power, so authority
validators earn rewards proportional to their boosted share, not their real
stake.

This is inherent to hybrid PoA and is intended: authority validators are trusted,
binary-/governance-designated operators, and the primary enforcement is
**tombstoning + governance removal**, not the token burn. To keep the burn
meaningful, operators SHOULD require a **minimum real self-stake** for authority
validators as an off-chain admission criterion (there is currently no on-chain
`min_self_stake` param enforcing this; adding one is a possible future change).

## Failure modes

### Authority validators get jailed (liveness)
Boost only applies to **bonded, unjailed** authority validators. As authority
validators are jailed (downtime) or tombstoned (double-sign), the remaining
authority validators must absorb the entire 67% floor, so their multiplier `M`
climbs. If the bonded authority set shrinks to zero, the chain enters safe mode
(or halts, depending on `halt_when_authority_empty`; see below).

Operational runbook:
- **Monitor** the bonded+unjailed authority count, the live `multiplier`
  attribute on the `authority_rescale` event (a rising multiplier warns the set
  is thinning), and the `authority_safe_mode_entered` event.
- **Re-seed quickly** via `MsgUpdateAuthoritySet` (gov) to add healthy
  validators, and unjail recoverable ones, before the set collapses.
- Keep a standing governance process (and signer availability) so an emergency
  `MsgUpdateAuthoritySet` / `MsgUpdateParams` can pass on a short voting period.
  With safe mode (the default), governance still runs because the chain keeps
  producing blocks — recovery is on-chain.

### All authority validators are jailed or unbonded
Behavior depends on `halt_when_authority_empty`:

- **Default (`false`) → safe mode.** The chain keeps producing blocks on
  community stake so governance can re-seed the authority set on-chain, while
  the value-bearing modules freeze (see [Safe mode](#safe-mode-authority-empty)).
  Recovery: pass a `MsgUpdateAuthoritySet` to restore a bonded authority set; on
  the next block boosting resumes, safe mode clears, and the frozen modules
  resume.
- **`true` → halt.** The chain halts via panic. The security guarantee is
  broken, so production of further blocks is refused; recovery requires an
  off-chain coordinated restart (governance cannot run on a halted chain).

## Safe mode (authority-empty)

When the bonded, unjailed authority set is empty and `halt_when_authority_empty`
is `false` (default), the chain enters **safe mode**: it keeps producing blocks
on community stake — so governance can recover the authority set on-chain —
while freezing every module that would commit a trust-bearing action under the
untrusted, community-only validator set.

This protects bridge/cellar funds from an attack that knocks the authority
validators offline (e.g. a DoS that downtime-jails them all): the worst outcome
is frozen value-bearing operations, never a fund movement signed by the
community set.

Frozen while in safe mode (txs, module BeginBlock/EndBlock, and legacy gov
proposals — the latter run in gov EndBlock and bypass the ante/msg servers):

| Module | Frozen |
|---|---|
| gravity-bridge | `MsgSendToEthereum`, `MsgSubmitEthereumEvent`, `MsgSubmitEthereumTxConfirmation` (and the same wrapped in an authz `MsgExec`), rejected by the ante handler; the module's BeginBlock/EndBlock (signer-set/batch creation, attestation observation, non-signing slashing) are no-op'd; the community-pool Ethereum spend gov proposal is rejected. `MsgDelegateKeys`, `MsgCancelSendToEthereum`, `MsgEthereumHeightVote` stay enabled. |
| cork | `MsgScheduleCork` and the `ScheduledCork` gov proposal; scheduled-cork execution in EndBlock |
| axelarcork | `MsgScheduleAxelarCork`, `MsgRelayAxelarCork`, `MsgRelayAxelarProxyUpgrade`, `MsgBumpAxelarCorkGas` msgs and the `AxelarScheduledCork` / `AxelarCommunityPoolSpend` / `UpgradeAxelarProxyContract` gov proposals; cork tally / fund sweep in EndBlock |

Pending items (a queued send-to-Ethereum, a cork scheduled for a future height)
are **not** dropped — they stay in module state and resume once safe mode
clears. Governance (`MsgUpdateAuthoritySet`), staking, and bank txs are not
frozen, so the recovery path stays open. Entry/exit emit the
`authority_safe_mode_entered` / `authority_safe_mode_exited` events.

**Thaw delay.** When the authority set is restored, boosting resumes
immediately but the value-bearing freeze is held for 2 more blocks — CometBFT
only applies the re-boosted validator set two blocks after it is emitted, so
thawing earlier would let the frozen modules act in a block still secured by
the old community-only set.

### A community validator grows large enough to exceed 33%
Cannot happen by construction. As long as the authority set is healthy,
the multiplier `M` is recomputed each block to keep the authority share at
or above `floor_fraction`. The community ceiling is `1 - floor_fraction`
by definition.

## Governance: changing the authority set

Submit a `MsgUpdateAuthoritySet` through gov. The message body lists the
new authority validators by operator address. The signer must be the chain's
gov authority (the `somm1...` address derived from the gov module account).

Validation:
- Empty list rejected.
- Duplicate operator addresses rejected.
- Malformed bech32 rejected.

Once the gov proposal passes, the new allowlist takes effect on the next
block. The boost is recomputed in the EndBlocker; CometBFT receives
adjusted ABCI ValidatorUpdates that reflect the new partition.

## Diagnostics

| Event | When | Attributes |
|---|---|---|
| `authority_rescale` | EndBlocker each block when boost is applied | `multiplier`, `authority_power`, `community_power` |
| `slash_skipped_no_snapshot` | Slash refused: snapshot missing for an at/after-activation infraction height (treated as corruption) | `operator`, `infraction_height` |
| `authority_safe_mode_entered` | Authority set became empty; chain entered safe mode (value-bearing modules frozen) | — |
| `authority_safe_mode_exited` | Authority set restored; safe mode cleared | — |
| `authority_set_updated` | After successful `MsgUpdateAuthoritySet` | `size` |
| `params_updated` | After successful `MsgUpdateParams` | `floor_fraction` |

Queries:
- `/sommelier/poa/v1/params` — current params
- `/sommelier/poa/v1/authority_set` — current allowlist
- `/sommelier/poa/v1/effective_power/{operator}` — current `LastValidatorPower` (boosted if authority) and `is_authority` flag
