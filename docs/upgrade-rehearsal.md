# Upgrade rehearsal against mainnet state

Before tagging a consensus-breaking release, replay the upgrade on a single-node
chain initialised from a **real mainnet state export**. Do not rely on unit
tests to clear an upgrade handler.

This document is the runbook. The worked example throughout is the v10 rehearsal
run on 2026-08-21 against mainnet height 27,355,689.

## Why this is not optional

Unit tests around an upgrade handler have a systematic blind spot: they build
the "before" state themselves, using the *new* binary's helpers. That state is
not the state the chain actually carries. Two defects in v10 got through a green
suite and would each have halted mainnet.

**1. The handler panicked on real state.**

```
panic: UnmarshalJSON cannot decode empty bytes
  app/upgrades/v10/upgrades.go:124
  x/upgrade.BeginBlocker -> ApplyUpgrade
```

`cork_authority` is new in v10, so the key is absent from v9 state. The handler
seeded it by round-tripping `GetParamSet` → mutate → `SetParams`, but
`GetParamSet` reads *every* registered param pair and amino-unmarshals the
missing key's empty bytes. The panic lands inside the upgrade `BeginBlocker`, so
every node stops at the upgrade height with no path forward.

Every handler test had called `SetParams(DefaultParams())` first. That writes
`cork_authority` — state no v9 chain has. **The tests were seeding away the
precondition of the bug.** See `seedV9CorkParams` in `app/upgrade_v10_test.go`
for the corrected shape, which writes the v9 key set and asserts the new key is
genuinely absent before invoking the handler.

**2. The release artifact would not start.**

The v10.0.0 binary required `GLIBC_2.34`; the validators run Debian 11 with
glibc 2.31. Moving the release job off the retired `ubuntu-20.04` runner put the
build on `ubuntu-22.04` (glibc 2.35), and the cgo-linked binary inherited its
symbol versions. Fixed by building static (`CGO_ENABLED=0` in
`.goreleaser.yaml`), which removes the dependency on whatever image built it.

A rehearsal catches this class of problem because it runs *the actual artifact*
on *the actual host OS*. Nothing in CI does.

The general lesson: **the rehearsal is the only step that exercises the real
`fromVM`, the real param store, and the real binary together.**

## Risk assessment before stopping anything

Step 1 stops a bonded validator. Measure the budget first; do not assume it.

```
sommelier query slashing params -o json
curl -s $API/cosmos/slashing/v1beta1/signing_infos/$VALCONS
```

For the v10 run:

| Fact | Value |
|---|---|
| `signed_blocks_window` | 5000 |
| `min_signed_per_window` | 0.05 |
| Max consecutive misses before jailing | 4750 |
| Block time (measured over 20k blocks) | 5.149 s |
| **Downtime budget before jailing** | **6.8 hours** |
| `slash_fraction_downtime` | 0 — jailing costs no stake |
| Actual downtime taken | **113 s** (~22 blocks) |
| Missed blocks after the run | 26 of 4750 |

### Which node

Pick a validator whose absence changes nothing:

- **Not** the Foundation validator. Under PoA it is the only member above the
  67% cork-approval threshold and the only one whose absence changes boost
  behaviour.
- An authority node with negligible stake is ideal. Removing it cannot trigger
  safe mode, which requires *all* authority members to be unbonded or jailed.

The v10 run used `sommelier-authority-3` (100 SOMM, us-west1-b).

## Step 1 — snapshot with minimal downtime

Copy the data dir while stopped, then restart immediately. Do **not** run the
export while the node is down; that extends downtime for no reason.

```bash
sudo systemctl stop sommelier
sudo cp -a /home/validator/.sommelier /home/validator/.somm-snapshot
sudo systemctl start sommelier
```

Confirm it is back **before** doing anything else — the RPC takes longer to come
up than the service takes to report `active`:

```bash
curl -s localhost:26657/status | jq -r '.result.sync_info
  | "height=\(.latest_block_height) catching_up=\(.catching_up)"'
```

Measure the data dir with `sudo`; without it `du` cannot descend into `data/`
and will report a misleadingly small number.

## Step 2 — export from the snapshot

```bash
sudo -u validator sommelier export --home /home/validator/.somm-snapshot \
  > /tmp/somm-export.json
```

## Step 3 — record the legacy state the migration must change

These counts are the rehearsal's success criteria. Record them **before**
upgrading.

```bash
jq '[.app_state.cork.scheduled_corks[]?] | length' /tmp/somm-export.json
jq '[.app_state.axelarcork.scheduled_corks.scheduled_corks[]?] | length' /tmp/somm-export.json
```

v10 run at height 27,355,689: `0` and `38`.

Also diff the param key sets, which is what surfaced defect (1) as a
*prediction* rather than a surprise:

```bash
jq -c '.app_state.cork.params'       /tmp/somm-export.json
jq -c '.app_state.axelarcork.params' /tmp/somm-export.json
```

Compare against `ParamSetPairs()` in the new binary. Any key the new binary
registers that the export lacks must be written explicitly by the handler, with
a targeted `Set` — never via a `GetParamSet` round-trip.

## Step 4 — build the rehearsal genesis

We hold only one validator's consensus key, so the exported set cannot produce
blocks. `scripts/rehearsal_genesis_surgery.py` inflates that validator to a
supermajority and keeps the books balanced:

```bash
scripts/rehearsal_genesis_surgery.py /tmp/somm-export.json /tmp/rehearsal-genesis.json \
  --valoper sommvaloper1... --delegator somm1...
```

It deliberately leaves every module under test untouched. Verify that:

```bash
jq '[.app_state.axelarcork.scheduled_corks.scheduled_corks[]?] | length' \
  /tmp/rehearsal-genesis.json   # still 38
```

## Step 5 — run the chain on the old binary

Set up an isolated home. **Every port must move**, or the rehearsal node fights
the live validator for them:

| Setting | Rehearsal value |
|---|---|
| `config.toml` p2p `laddr` | `tcp://0.0.0.0:36656` |
| `config.toml` rpc `laddr` | `tcp://127.0.0.1:36657` |
| `config.toml` `persistent_peers`, `seeds` | empty |
| `config.toml` `pex` | `false` |
| `config.toml` `[statesync] enable` | **`false`** |
| `app.toml` `[api]`, `[grpc]`, `[grpc-web]` `enable` | **`false`** |

Two of these bit during the v10 run and are worth calling out:

- **State sync left enabled** made the node try to sync from mainnet RPCs and
  reject its own genesis (`header belongs to another chain "sommelier-3"`).
- **gRPC left enabled** made the node bind `:9090`, which the live validator
  owns. The rehearsal node exited with `bind: address already in use`. The live
  node held the port and was unaffected, but disable gRPC rather than relying on
  that.

Then start the **old** binary and confirm blocks commit. Getting this far
already proves a full mainnet export initialises cleanly.

## Step 6 — pass the upgrade proposal

Submit a `MsgSoftwareUpgrade` at roughly current height + 150, vote yes, wait.

Practical notes from the v10 run:

- Use `http://` for `curl` and `tcp://` for `--node`. Mixing them yields an
  empty height and a proposal targeting height 60.
- Give the voting period enough room (`--voting-period 180s`). A 30s period
  expired before the vote landed and the first proposal was rejected.
- Extract the proposal id from the **tx result**, not from
  `query gov proposals | tail`, which pages and will show old proposals.
- Proposal ids restart from the exported counter, so the id will not match
  mainnet's next id.

## Step 7 — swap in the new binary

The old binary halts with `UPGRADE "<name>" NEEDED at height: N`. That is the
expected stop. Kill it and start the new binary against the same home.

**Use the real release artifact**, downloaded from the release page and checksum
-verified — not a locally built binary. The glibc defect existed only in the
artifact and would have been invisible to a local build.

```bash
curl -sL -o v.tar.gz https://github.com/PeggyJV/sommelier/releases/download/vX/....tar.gz
sha256sum v.tar.gz          # compare against SHA256SUMS
file ./sommelier            # expect "statically linked"
objdump -T ./sommelier | grep -o 'GLIBC_[0-9.]*' | sort -uV | tail   # expect empty
./sommelier version --long  # expect the tag and commit you released
```

## Step 8 — assert

The upgrade log should show the handler's own markers. From the v10 run:

```
applying upgrade "v10" at height: 27356219
migrating module gravity from version 6 to version 7
adding a new module: poa
v10 upgrade: authority validator liveness verified  bonded_and_unjailed=4 configured=4
v10 upgrade: PoA params and authority set initialised  authority_count=4 activation_height=27356219
v10 upgrade: cork authority seeded and legacy queues drained  drained_keys=38
```

`drained_keys` **must equal the Step 3 counts**. Then verify via queries and a
fresh export:

| Assertion | v10 result |
|---|---|
| Blocks commit past the upgrade height | yes |
| `cork params` → `cork_authority` | set |
| `axelarcork params` → `cork_authority` | set |
| `poa authority-set` | 4 configured validators |
| `cork scheduled-corks` | `[]` |
| `axelarcork scheduled-corks 42161` | `[]` |
| Fresh export: both legacy queues | `0`, `0` |
| Fresh export: `poa.activation_height` | upgrade height |
| Release artifact runs on upgraded state | 42 blocks, 0 panics |

## Step 9 — clean up

```bash
sudo rm -rf /home/validator/.somm-snapshot
rm -rf /tmp/rehearsal
```

Leaving a stale duplicate chain dir on a validator is an operational trap. The
v10 run reclaimed 16G.

Then re-check the validator you borrowed:

```bash
systemctl is-active sommelier
curl -s $API/cosmos/slashing/v1beta1/signing_infos/$VALCONS | jq .val_signing_info
```

## Known pre-existing failure: `validate-genesis`

`sommelier validate-genesis` **fails on Sommelier mainnet exports**, before and
after any upgrade:

```
Error: connection identifier is not in the format: `connection-{N}`: invalid identifier
```

The cause is a `connection-localhost` entry in IBC connection genesis, an
artifact of an older ibc-go. Confirm it is pre-existing rather than caused by
your change by running the **old** binary against the **pre-upgrade** export; if
it fails identically, it is not your regression.

Operationally this means **a genesis restart from an export would be rejected
today**, independent of any upgrade. Worth fixing separately; it is not an
upgrade blocker.

## What this rehearsal does not cover

State it explicitly rather than implying full coverage:

- **Multi-validator consensus.** One node with a synthetic supermajority does
  not exercise a real validator set upgrading together.
- **Authorised-actor paths needing keys you do not hold.** The v10 run could not
  exercise scheduling a cork from the cork authority.
- **Real network conditions** — peers, mempool pressure, cosmovisor mechanics.
