# Cork wind-down authority

Date: 2026-08-20
Status: approved, targeted at the v10 upgrade
Modules: `x/cork`, `x/axelarcork`, `x/pubsub` (decoupling only)

## Problem

Scheduling a cork today requires the signer to be a registered Gravity
orchestrator delegate of a bonded validator (`GetOrchestratorValidatorAddress`),
and a cork executes only when validators holding **more than 67%** of consensus
power schedule the byte-identical cork (`corkVoteThresholdStr = "0.67"`,
`x/cork/keeper/keeper.go:18,331`; `x/axelarcork/keeper/keeper.go:98,411`).

Sommelier is winding down its cellars: 36 managed Ethereum cellars and 8 Axelar
destination chains. Driving that wind-down through a validator vote means
coordinating a supermajority for every cork, on a validator set that is itself
being restructured by the v10 PoA power floor. The strategist system the vote
was built for is being retired.

## Goal

Replace validator-voted cork scheduling with a single governance-designated
authority address, and delete as much of the legacy strategist machinery as
possible, while keeping the chain and its remaining modules operating normally.

Initial authority address: `somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa`

## Decisions

These were settled during design and are not open questions.

| Question | Decision |
|---|---|
| End state | Cellars close; the chain keeps operating. The mechanism must be safe to leave in place indefinitely. |
| Replace or coexist | Replace validator voting entirely. No dormant second path. |
| Authority designation | A single address in module `Params`, rotated by governance. |
| Target scope | The managed-cellar allowlist constraint is **kept**. The authority cannot target arbitrary contracts. |
| Axelarcork messages | All four (`ScheduleCork`, `RelayCork`, `BumpCorkGas`, `CancelScheduledCork`) become authority-only. |
| Scheduling model | The scheduled-height queue is **kept**. Only the tally is removed. |
| `x/pubsub` | Decoupled from cork, module left standing with state intact. Not deleted in this upgrade. |

## Design

### 1. Authorization

A new param `cork_authority`, a bech32 `somm1…` string:

* `proto/cork/v2/genesis.proto` — `Params` field 3
* `proto/axelarcork/v1/genesis.proto` — `Params` field 8

Both modules already use `paramtypes.Subspace`, so rotation requires **no new
message type**: a standard `ParameterChangeProposal` moves the authority. That
is the recovery path if the key is lost or compromised.

The authorization check replacing `GetOrchestratorValidatorAddress` is an
equality test against the param, applied in:

* `x/cork/keeper/msg_server.go` — `ScheduleCork`
* `x/axelarcork/keeper/msg_server.go` — `ScheduleCork`, `RelayCork`,
  `BumpCorkGas`, `CancelScheduledCork`

**Fail-closed.** An unset or malformed `cork_authority` means nothing can
schedule a cork. There is deliberately no fallback to validator voting: a
fallback would defeat the purpose and leave a dormant second path. The accepted
cost is that a bad parameter change disables cork scheduling until governance
corrects it.

`vote_threshold` is already marked `Deprecated` in the proto and is unread by
the keeper. Leave the field in place; removing it would renumber the wire
format for no benefit.

### 2. Cork storage and execution

Current key:

```
ScheduledCorkKeyPrefix | block_height | cork_id | val_address | contract_address -> Cork
```

The validator dimension disappears. New prefix:

```
AuthorityCorkKeyPrefix | block_height | cork_id | contract_address -> Cork
```

`MsgScheduleCorkRequest` keeps its name, minus validator-derived fields.

The EndBlocker becomes: iterate corks due at the current height, submit each via
`submitContractCall`, delete. No tally, no quorum, no approval bookkeeping.

Deleted:

* the power tally inside `GetApprovedScheduledCorks`
* `ValidatorCorkCountKey` and its increment/decrement
* `MaxCorksPerValidator` enforcement
* the `corkVoteThresholdStr` constant

The `max_corks_per_validator` and `vote_threshold` param *fields* remain in
`Params` and are simply unread, for the same wire-compatibility reason given
above. Only their enforcement is removed.

`ScheduledCorkKeyPrefix` stays **defined but unwritten**, matching the existing
convention in `x/cork/types/keys.go` for `CorkForAddressKeyPrefix` and
`CommitPeriodStartKey`.

`CorkResult` records become vestigial. Keep the `cork-result` query endpoints
serving historical records and stop writing new ones, rather than breaking API
consumers mid-wind-down.

### 3. Axelarcork

Same param, same check, all four messages.

New key, retaining the chain dimension:

```
AuthorityCorkKeyPrefix | chain_id | block_height | cork_id | contract_address -> Cork
```

Axelarcork is two-stage: scheduled → relayable (`WinningAxelarCork`) → relayed
over Axelar GMP by `RelayCork`. The relayable stage is **not** a vote artifact
and is retained. The EndBlocker keeps its per-chain iteration, but the tally is
replaced by a direct move: corks due at this height go straight into the
`WinningAxelarCork` queue.

`cork_timeout_blocks` and the timed-out-cork sweep are **unchanged**. These are
liveness mechanisms, not strategist mechanisms, and they matter more once
relaying depends on a single key: a cork that is never relayed must still
expire.

### 4. Pubsub decoupling

* `HandleAddManagedCellarsProposal` — drop the `GetPublisher` check and the
  `SetDefaultSubscription` call (both modules)
* `HandleRemoveManagedCellarsProposal` — drop `DeleteDefaultSubscription`
  (both modules)
* Remove the `pubsubKeeper` field and the `PubsubKeeper` expected-keeper
  interface from both keepers' constructors

`x/pubsub` remains registered, with its state and its four proposal handlers
intact and no callers. No store migration. Deleting it is a clean follow-up.

### 5. Migration

**This is the part that must not be got wrong.**

`GetApprovedScheduledCorks` is currently the *only* site that deletes from the
legacy scheduled-cork queue. Once the tally is removed, any cork sitting in that
queue at upgrade height becomes permanently undeletable, with its scheduler's
cork-count quota consumed forever. The existing code comments in both
EndBlockers call out precisely this hazard for the safe-mode case.

The v10 upgrade handler must therefore **drain both legacy queues**:

1. Iterate every scheduled cork under `x/cork`'s `ScheduledCorkKeyPrefix` and
   delete it.
2. For each chain configuration, iterate every scheduled cork under
   `x/axelarcork`'s `ScheduledCorkKeyPrefix` and delete it.
3. Delete all `ValidatorCorkCountKey` / `ValidatorAxelarCorkCountKey` entries.

Drop these corks rather than migrating them into the new queue. They were
scheduled under validator consent that no longer carries meaning, and
re-scheduling under the authority key costs one transaction.

### 6. Safe mode interaction

Every existing `inSafeMode` gate stays exactly where it is, now gating the
authority path instead of the validator path. This is deliberate: a cork
ultimately executes through a bridge secured by the validator set, so an
untrusted set makes an authority-scheduled cork no safer than a voted one.

## Testing

Weighted for the fact that this ships in v10 alongside the PoA power floor and
the Gravity GHSA-4vf2-m5pw-3r3r fixes.

**Unit**
* authority accepted; non-authority rejected; empty and malformed params
  rejected — across all five message types in both modules
* allowlist constraint still enforced for the authority (unmanaged cellar
  rejected)

**Store**
* new key round-trips
* legacy prefix is never written

**EndBlocker**
* corks due at height submit exactly once, then are deleted
* nothing due is a no-op
* safe mode drops rather than defers
* axelarcork: due cork lands in the `WinningAxelarCork` queue
* axelarcork: timeout sweep still expires unrelayed corks

**Migration**
* against exported mainnet state at a recent height: run the drain, assert both
  legacy queues are empty and no cork-count keys remain

**Integration**
* authority tx → `CreateContractCallTx` end to end
* axelarcork: schedule → relayable → relay

**Testnet**
* full v10 upgrade rehearsal with corks in flight across the upgrade boundary

## Risks

**Concentration.** A single key gains unilateral control over every cork on 36
Ethereum cellars and 8 Axelar chains. The allowlist constraint bounds the blast
radius to already-approved cellars — a compromised key cannot make the bridge
call arbitrary contracts — but within that set its power is total. Governance
rotation is the only remedy, on a 48-hour voting period.

**Release concentration.** v10 will carry three independent large changes: a
novel consensus-affecting module that has never run in production, a security
fix, and this authorization rewrite. If the upgrade misbehaves at the upgrade
height, attribution will be hard. This was raised during design and shipping in
v10 was chosen deliberately; the testnet rehearsal above is the mitigation and
should not be cut for schedule.

**Fail-closed parameter.** A malformed `cork_authority` disables cork scheduling
chain-wide until a governance proposal completes.

## Out of scope

* Deleting `x/pubsub` (follow-up upgrade)
* Removing `x/cork` / `x/axelarcork` entirely
* Any change to Gravity, auction, cellarfees, or incentives
* The `x/poa` authority validator set, which is a separate mechanism with a
  separate purpose
