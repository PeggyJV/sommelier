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
| `halt_when_authority_empty` | `true` | Panic in EndBlocker if no authority validator is bonded and unjailed. |

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

## Failure modes

### All authority validators are jailed or unbonded
With `halt_when_authority_empty=true` (default), the chain halts via panic.
This is the correct PoA failure mode: the security guarantee is broken,
production of further blocks is refused. Recovery requires governance
intervention to update or unjail authority validators.

### A community validator grows large enough to exceed 33%
Cannot happen by construction. As long as the authority set is healthy,
the multiplier `M` is recomputed each block to keep the authority share at
or above `floor_fraction`. The community ceiling is `1 - floor_fraction`
by definition.

## Governance: changing the authority set

Submit a `MsgUpdateAuthoritySet` through gov. The message body lists the
new authority validators by operator address. The signer must be the chain's
gov authority (typically `cosmos1...gov` derived from the gov module account).

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
| `slash_skipped_no_snapshot` | Authority slash refused due to missing snapshot | `operator`, `infraction_height` |
| `authority_set_updated` | After successful `MsgUpdateAuthoritySet` | `size` |
| `params_updated` | After successful `MsgUpdateParams` | `floor_fraction` |

Queries:
- `/sommelier/poa/v1/params` — current params
- `/sommelier/poa/v1/authority_set` — current allowlist
- `/sommelier/poa/v1/effective_power/{operator}` — current `LastValidatorPower` (boosted if authority) and `is_authority` flag
