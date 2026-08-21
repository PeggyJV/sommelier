# Cork Wind-Down Authority Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace validator-voted cork scheduling in `x/cork` and `x/axelarcork` with a single governance-designated authority address, and drain the legacy queues in the v10 upgrade handler.

**Architecture:** A new `cork_authority` bech32 string param is added to both modules' existing `paramtypes.Subspace` params, so governance rotates it with a standard `ParameterChangeProposal` and no new message type is needed. The `GetOrchestratorValidatorAddress` delegate check is replaced by an equality test against that param. The >67% power tally is deleted; corks are stored under a new validator-free key prefix and executed directly by the EndBlocker at their target height. Both modules stop calling `x/pubsub`.

**Tech Stack:** Go 1.22, Cosmos SDK v0.47.15, CometBFT v0.37.18, gogoproto, `make proto-gen` (Docker-based buf/protoc).

**Spec:** `docs/superpowers/specs/2026-08-20-cork-wind-down-authority-design.md`

## Global Constraints

- Module path is `github.com/peggyjv/sommelier/v10` — the v9→v10 bump has already landed. Never write `sommelier/v9`.
- Initial authority address: `somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa`
- Fail-closed: an unset or malformed `cork_authority` means nothing can schedule a cork. Never add a fallback to validator voting.
- **No kill switch.** Do not add an `enabled` bool or any pause mechanism for the authority path.
- Every existing `inSafeMode` gate stays exactly where it is, now gating the authority path.
- Retired key prefixes stay **defined but unwritten**, matching `CorkForAddressKeyPrefix` and `CommitPeriodStartKey` in `x/cork/types/keys.go`.
- Retired param *fields* (`vote_threshold`, `max_corks_per_validator`) stay in `Params` and are simply unread. Do not renumber proto fields.
- Baseline before starting: `go test $(go list ./... | grep -v integration_tests)` is **23 ok, 0 FAIL**. Any new failure is yours.
- Proto changes require `make proto-gen` (needs Docker). Commit regenerated `.pb.go` alongside the `.proto` change.
- Run `go build ./app/... ./x/... ./cmd/...` after each task. `go build ./...` fails on `integration_tests` for pre-existing reasons — that is expected, not a regression.

---

## File Structure

**`x/cork`**
- `proto/cork/v2/genesis.proto` — add `cork_authority` field 3 to `Params`
- `x/cork/types/v2/params.go` — param key, default, pair, validator
- `x/cork/types/keys.go` — `AuthorityCorkKeyPrefix` + key helpers
- `x/cork/keeper/keeper.go` — authority-cork store accessors; delete cork-count and tally helpers
- `x/cork/keeper/msg_server.go` — authority check
- `x/cork/keeper/abci.go` — tally-free EndBlocker
- `x/cork/keeper/proposal_handler.go` — drop pubsub calls
- `x/cork/types/expected_keepers.go` — drop `PubsubKeeper`

**`x/axelarcork`** — same shape, files mirrored.

**Upgrade**
- `app/upgrades/v10/constants.go` — `CorkAuthorityAddress`
- `app/upgrades/v10/upgrades.go` — seed params, drain legacy queues
- `app/app.go` — keeper constructor call sites

---

### Task 1: Add `cork_authority` param to x/cork

**Files:**
- Modify: `proto/cork/v2/genesis.proto:23-32`
- Modify: `x/cork/types/v2/params.go`
- Test: `x/cork/types/v2/params_test.go`

**Interfaces:**
- Consumes: nothing
- Produces: `v2.KeyCorkAuthority []byte`; `Params.CorkAuthority string`; `v2.DefaultParams()` returns `CorkAuthority: ""`

- [ ] **Step 1: Write the failing test**

Create `x/cork/types/v2/params_test.go`:

```go
package v2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCorkAuthorityValidation(t *testing.T) {
	// Empty is allowed at the type level; the msg server is what fails closed.
	require.NoError(t, validateCorkAuthority(""))
	require.NoError(t, validateCorkAuthority("somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa"))

	require.Error(t, validateCorkAuthority("not-bech32"))
	require.Error(t, validateCorkAuthority("cosmos1qypqxpq9qcrsszg2pvxq6rs0zqg3yyc5lzv7xu"))
	require.Error(t, validateCorkAuthority(12345))
}

func TestDefaultParamsHasEmptyAuthority(t *testing.T) {
	require.Equal(t, "", DefaultParams().CorkAuthority)
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./x/cork/types/v2/... -run 'TestCorkAuthority|TestDefaultParamsHasEmpty' -v`
Expected: FAIL — `undefined: validateCorkAuthority`

- [ ] **Step 3: Add the proto field**

In `proto/cork/v2/genesis.proto`, inside `message Params`, after `max_corks_per_validator`:

```proto
    // cork_authority is the single address permitted to schedule corks.
    // Rotated by governance via ParameterChangeProposal. Empty means no
    // address can schedule a cork (fail-closed).
    string cork_authority = 3 [(gogoproto.moretags) = "yaml:\"cork_authority\""];
```

- [ ] **Step 4: Regenerate protos**

Run: `make proto-gen`
Expected: `x/cork/types/v2/genesis.pb.go` gains a `CorkAuthority string` field on `Params`.

- [ ] **Step 5: Wire the param**

In `x/cork/types/v2/params.go`, add to the key block:

```go
	KeyCorkAuthority        = []byte("corkauthority")
```

Add to `DefaultParams()`:

```go
		CorkAuthority:        "",
```

Add to `ParamSetPairs()`:

```go
		paramtypes.NewParamSetPair(KeyCorkAuthority, &p.CorkAuthority, validateCorkAuthority),
```

Add to `ValidateBasic()` before `return nil`:

```go
	if err := validateCorkAuthority(p.CorkAuthority); err != nil {
		return err
	}
```

Add the validator function:

```go
// validateCorkAuthority accepts the empty string (fail-closed: no address can
// schedule) or a well-formed somm1 account address.
func validateCorkAuthority(i interface{}) error {
	authority, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if authority == "" {
		return nil
	}
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		return fmt.Errorf("invalid cork authority address %q: %w", authority, err)
	}
	return nil
}
```

- [ ] **Step 6: Run test to verify it passes**

Run: `go test ./x/cork/types/v2/... -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add proto/cork/v2/genesis.proto x/cork/types/v2/genesis.pb.go x/cork/types/v2/params.go x/cork/types/v2/params_test.go
git commit -m "cork: add cork_authority param"
```

---

### Task 2: Add authority cork key prefix and store accessors to x/cork

**Files:**
- Modify: `x/cork/types/keys.go:26-52` (const block), `:55-100` (helpers)
- Modify: `x/cork/keeper/keeper.go`
- Test: `x/cork/keeper/authority_cork_test.go`

**Interfaces:**
- Consumes: Task 1
- Produces:
  - `types.AuthorityCorkKeyPrefix byte`
  - `types.GetAuthorityCorkKeyPrefix() []byte`
  - `types.GetAuthorityCorkKeyByBlockHeightPrefix(blockHeight uint64) []byte`
  - `types.GetAuthorityCorkKey(blockHeight uint64, id []byte, contract common.Address) []byte`
  - `Keeper.SetAuthorityCork(ctx sdk.Context, blockHeight uint64, cork v2types.Cork) []byte`
  - `Keeper.DeleteAuthorityCork(ctx sdk.Context, blockHeight uint64, id []byte, contract common.Address)`
  - `Keeper.IterateAuthorityCorksByBlockHeight(ctx sdk.Context, blockHeight uint64, cb func(blockHeight uint64, id []byte, contract common.Address, cork v2types.Cork) (stop bool))`

- [ ] **Step 1: Write the failing test**

Create `x/cork/keeper/authority_cork_test.go`:

```go
package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	v2types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
)

func TestAuthorityCorkRoundTrip(t *testing.T) {
	ctx, k := setupKeeperForTest(t)

	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	cork := v2types.Cork{
		TargetContractAddress: contract.String(),
		EncodedContractCall:   []byte{0xde, 0xad, 0xbe, 0xef},
	}

	id := k.SetAuthorityCork(ctx, 100, cork)
	require.NotEmpty(t, id)

	var seen []v2types.Cork
	k.IterateAuthorityCorksByBlockHeight(ctx, 100, func(_ uint64, _ []byte, _ common.Address, c v2types.Cork) bool {
		seen = append(seen, c)
		return false
	})
	require.Len(t, seen, 1)
	require.Equal(t, cork.EncodedContractCall, seen[0].EncodedContractCall)

	// A different height must not see it.
	var other int
	k.IterateAuthorityCorksByBlockHeight(ctx, 101, func(_ uint64, _ []byte, _ common.Address, _ v2types.Cork) bool {
		other++
		return false
	})
	require.Zero(t, other)

	k.DeleteAuthorityCork(ctx, 100, id, contract)

	var after int
	k.IterateAuthorityCorksByBlockHeight(ctx, 100, func(_ uint64, _ []byte, _ common.Address, _ v2types.Cork) bool {
		after++
		return false
	})
	require.Zero(t, after)
}
```

Note: `setupKeeperForTest` is the existing helper in `x/cork/keeper/test_common.go`. Read that file first and match its exact signature — if it returns `(sdk.Context, Keeper)` in a different order or name, adapt the call, do not change the helper.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./x/cork/keeper/... -run TestAuthorityCorkRoundTrip -v`
Expected: FAIL — `k.SetAuthorityCork undefined`

- [ ] **Step 3: Add the key prefix**

In `x/cork/types/keys.go`, append to the end of the `const` iota block (a new prefix byte — never reuse or renumber an existing one):

```go
	// AuthorityCorkKeyPrefix - <prefix><block_height><cork_id><contract_address> -> <cork>
	// Corks scheduled by the cork authority. Replaces ScheduledCorkKeyPrefix,
	// which is retained above but no longer written.
	AuthorityCorkKeyPrefix
```

Add the helpers after `GetScheduledCorkKey`:

```go
func GetAuthorityCorkKeyPrefix() []byte {
	return []byte{AuthorityCorkKeyPrefix}
}

func GetAuthorityCorkKeyByBlockHeightPrefix(blockHeight uint64) []byte {
	return append(GetAuthorityCorkKeyPrefix(), sdk.Uint64ToBigEndian(blockHeight)...)
}

func GetAuthorityCorkKey(blockHeight uint64, id []byte, contract common.Address) []byte {
	return bytes.Join(
		[][]byte{GetAuthorityCorkKeyPrefix(), sdk.Uint64ToBigEndian(blockHeight), id, contract.Bytes()},
		[]byte{},
	)
}
```

- [ ] **Step 4: Add the keeper accessors**

In `x/cork/keeper/keeper.go`, after the existing scheduled-cork accessors:

```go
// SetAuthorityCork stores a cork scheduled by the cork authority for execution
// at blockHeight and returns its ID.
func (k Keeper) SetAuthorityCork(ctx sdk.Context, blockHeight uint64, cork types.Cork) []byte {
	id := cork.IDHash(blockHeight)
	bz := k.cdc.MustMarshal(&cork)
	ctx.KVStore(k.storeKey).Set(
		corktypes.GetAuthorityCorkKey(blockHeight, id, common.HexToAddress(cork.TargetContractAddress)),
		bz,
	)
	return id
}

// DeleteAuthorityCork removes a scheduled authority cork.
func (k Keeper) DeleteAuthorityCork(ctx sdk.Context, blockHeight uint64, id []byte, contract common.Address) {
	ctx.KVStore(k.storeKey).Delete(corktypes.GetAuthorityCorkKey(blockHeight, id, contract))
}

// IterateAuthorityCorksByBlockHeight walks authority corks targeting blockHeight.
func (k Keeper) IterateAuthorityCorksByBlockHeight(
	ctx sdk.Context,
	blockHeight uint64,
	cb func(blockHeight uint64, id []byte, contract common.Address, cork types.Cork) (stop bool),
) {
	prefix := corktypes.GetAuthorityCorkKeyByBlockHeightPrefix(blockHeight)
	iter := sdk.KVStorePrefixIterator(ctx.KVStore(k.storeKey), prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		// key layout after the 1-byte prefix: height(8) | id(32) | contract(20)
		idStart := 1 + 8
		contractStart := idStart + 32
		id := key[idStart:contractStart]
		contract := common.BytesToAddress(key[contractStart:])

		var cork types.Cork
		k.cdc.MustUnmarshal(iter.Value(), &cork)

		if cb(blockHeight, id, contract, cork) {
			break
		}
	}
}
```

Before writing this, confirm the ID width by reading `Cork.IDHash` in `x/cork/types/v2/cork.go`. If it is not 32 bytes, correct `contractStart` accordingly — do not guess.

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./x/cork/keeper/... -run TestAuthorityCorkRoundTrip -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add x/cork/types/keys.go x/cork/keeper/keeper.go x/cork/keeper/authority_cork_test.go
git commit -m "cork: add authority cork key prefix and store accessors"
```

---

### Task 3: Switch cork ScheduleCork to authority authorization

**Files:**
- Modify: `x/cork/keeper/msg_server.go:20-60`
- Test: `x/cork/keeper/msg_server_authority_test.go`

**Interfaces:**
- Consumes: Tasks 1, 2
- Produces: `ScheduleCork` accepts only `params.CorkAuthority`

- [ ] **Step 1: Write the failing test**

Create `x/cork/keeper/msg_server_authority_test.go`:

```go
package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const testAuthority = "somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa"

func TestScheduleCorkRejectsNonAuthority(t *testing.T) {
	// authority set, signer is someone else -> unauthorized
}

func TestScheduleCorkRejectsEmptyAuthority(t *testing.T) {
	// authority == "" -> unauthorized even for a well-formed signer (fail-closed)
}

func TestScheduleCorkAcceptsAuthority(t *testing.T) {
	// signer == authority, cellar allowlisted, future height -> cork stored
}

func TestScheduleCorkRejectsUnmanagedCellar(t *testing.T) {
	// signer == authority but target not in allowlist -> ErrUnmanagedCellarAddress
}

var _ = require.New
```

Fill each body against the existing test harness in `x/cork/keeper/test_common.go`. Read that file and mirror how existing msg-server tests construct a context, set params, and add a cellar ID. Do not invent new harness helpers.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./x/cork/keeper/... -run TestScheduleCork -v`
Expected: FAIL — the non-authority case still succeeds via the delegate path.

- [ ] **Step 3: Replace the authorization check**

In `x/cork/keeper/msg_server.go`, delete these lines:

```go
	signer := msg.MustGetSigner()
	validatorAddr := k.gravityKeeper.GetOrchestratorValidatorAddress(ctx, signer)
	if validatorAddr == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "signer %s is not a delegate", signer.String())
	}

	params := k.GetParamSet(ctx)
	validatorCorkCount := k.GetValidatorCorkCount(ctx, validatorAddr)
	if validatorCorkCount >= params.MaxCorksPerValidator {
		return nil, corktypes.ErrValidatorCorkCapacityReached
	}
```

Replace with:

```go
	signer := msg.MustGetSigner()
	params := k.GetParamSet(ctx)
	// Fail-closed: an unset authority means no address may schedule a cork.
	if params.CorkAuthority == "" || signer.String() != params.CorkAuthority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
			"signer %s is not the cork authority", signer.String())
	}
```

Then replace the storage call and the count increment:

```go
	corkID := k.SetAuthorityCork(ctx, msg.BlockHeight, *msg.Cork)
```

Delete the `k.IncrementValidatorCorkCount(ctx, validatorAddr)` line. In the event emission below, remove the `AttributeKeyValidator` attribute since there is no longer a validator.

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./x/cork/keeper/... -run TestScheduleCork -v`
Expected: PASS

- [ ] **Step 5: Build**

Run: `go build ./x/cork/...`
Expected: exit 0

- [ ] **Step 6: Commit**

```bash
git add x/cork/keeper/msg_server.go x/cork/keeper/msg_server_authority_test.go
git commit -m "cork: authorize ScheduleCork by cork authority param"
```

---

### Task 4: Rewrite cork EndBlocker without the tally

**Files:**
- Modify: `x/cork/keeper/abci.go:57-101`
- Modify: `x/cork/keeper/keeper.go:294-350` (delete `GetApprovedScheduledCorks`)
- Test: `x/cork/keeper/abci_authority_test.go`

**Interfaces:**
- Consumes: Tasks 2, 3
- Produces: EndBlocker submits every authority cork due at the current height exactly once, then deletes it

- [ ] **Step 1: Write the failing test**

Create `x/cork/keeper/abci_authority_test.go` covering:

```
TestEndBlockerSubmitsDueAuthorityCorks   — two corks at height H, both submitted, both deleted
TestEndBlockerIgnoresFutureCorks         — cork at H+1 is untouched at H
TestEndBlockerSafeModeDropsCorks         — safe mode: corks deleted, no contract call, drop event emitted
```

Mirror the existing safe-mode test style in `x/cork/keeper/safemode_test.go`; it already mocks the gravity keeper and asserts on `CreateContractCallTx`. Read it before writing.

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./x/cork/keeper/... -run TestEndBlocker -v`
Expected: FAIL — authority corks are never read by the EndBlocker.

- [ ] **Step 3: Rewrite the EndBlocker**

Replace the body of `func (k Keeper) EndBlocker(ctx sdk.Context)` in `x/cork/keeper/abci.go` with:

```go
func (k Keeper) EndBlocker(ctx sdk.Context) {
	height := uint64(ctx.BlockHeight())
	safeMode := k.inSafeMode(ctx)

	type due struct {
		id       []byte
		contract common.Address
		cork     v2types.Cork
	}
	var dueCorks []due

	// Collect first, then mutate: deleting inside the iterator invalidates it.
	k.IterateAuthorityCorksByBlockHeight(ctx, height, func(_ uint64, id []byte, contract common.Address, cork v2types.Cork) bool {
		dueCorks = append(dueCorks, due{id: id, contract: contract, cork: cork})
		return false
	})

	if len(dueCorks) == 0 {
		return
	}

	// Always delete, even in safe mode: this is the only site that removes due
	// corks, so skipping deletion would strand them in state permanently.
	for _, d := range dueCorks {
		k.DeleteAuthorityCork(ctx, height, d.id, d.contract)
	}

	if safeMode {
		// Corks due during the freeze are DROPPED, not deferred: submitting one
		// later would execute a call through a bridge secured by a set we no
		// longer trust. They must be re-scheduled after recovery.
		k.Logger(ctx).Error("x/poa safe mode active: dropping due authority corks without submitting",
			"height", fmt.Sprintf("%d", ctx.BlockHeight()),
			"dropped", len(dueCorks))
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			types.EventTypeCorksDroppedInSafeMode,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyBlockHeight, fmt.Sprintf("%d", ctx.BlockHeight())),
			sdk.NewAttribute(types.AttributeKeyDroppedCount, fmt.Sprintf("%d", len(dueCorks))),
		))
		return
	}

	for _, d := range dueCorks {
		k.submitContractCall(ctx, d.cork)
	}
}
```

- [ ] **Step 4: Delete the tally**

Delete `func (k Keeper) GetApprovedScheduledCorks` from `x/cork/keeper/keeper.go` and the `corkVoteThresholdStr` constant at `x/cork/keeper/keeper.go:18`. Delete `GetValidatorCorkCount`, `IncrementValidatorCorkCount`, and `DecrementValidatorCorkCount`. Leave `GetValidatorCorkCountKey` in `keys.go` — the migration in Task 11 needs it.

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./x/cork/keeper/... -v`
Expected: PASS, no failures

- [ ] **Step 6: Commit**

```bash
git add x/cork/keeper/abci.go x/cork/keeper/keeper.go x/cork/keeper/abci_authority_test.go
git commit -m "cork: execute authority corks directly, remove power tally"
```

---

### Task 5: Decouple x/cork from x/pubsub

**Files:**
- Modify: `x/cork/keeper/proposal_handler.go:19-86`
- Modify: `x/cork/keeper/keeper.go:20-50`
- Modify: `x/cork/types/expected_keepers.go`
- Modify: `app/app.go` (cork keeper constructor call)

**Interfaces:**
- Consumes: Task 4
- Produces: `keeper.NewKeeper(cdc, key, paramSpace, stakingKeeper, gravityKeeper)` — the `pubsubKeeper` parameter is **removed**

- [ ] **Step 1: Write the failing test**

Add to `x/cork/keeper/proposal_handler_test.go` (create if absent):

```go
func TestAddManagedCellarsNoLongerRequiresPublisher(t *testing.T) {
	// A proposal with an unknown PublisherDomain must now succeed.
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./x/cork/keeper/... -run TestAddManagedCellarsNoLonger -v`
Expected: FAIL — "not an approved publisher"

- [ ] **Step 3: Strip the pubsub calls**

In `x/cork/keeper/proposal_handler.go`:
- In `HandleAddManagedCellarsProposal`, delete the `GetPublisher` lookup and its error return, and delete the `SetDefaultSubscription` block.
- In `HandleRemoveManagedCellarsProposal`, delete the `DeleteDefaultSubscription` block.
- Remove the now-unused `pubsubtypes` import.

In `x/cork/keeper/keeper.go`, remove the `pubsubKeeper` struct field, the constructor parameter, and its assignment.

In `x/cork/types/expected_keepers.go`, delete the `PubsubKeeper` interface.

In `app/app.go`, drop the `app.PubsubKeeper` argument from the cork `NewKeeper` call.

- [ ] **Step 4: Run tests and build**

Run: `go test ./x/cork/... && go build ./app/... ./x/...`
Expected: PASS, exit 0

- [ ] **Step 5: Commit**

```bash
git add x/cork app/app.go
git commit -m "cork: decouple from x/pubsub"
```

---

### Task 6: Add `cork_authority` param to x/axelarcork

**Files:**
- Modify: `proto/axelarcork/v1/genesis.proto:20-28`
- Modify: `x/axelarcork/types/params.go`
- Test: `x/axelarcork/types/params_test.go`

**Interfaces:**
- Consumes: nothing (independent of Tasks 1-5)
- Produces: `types.KeyCorkAuthority`; `Params.CorkAuthority string`; default `""`

- [ ] **Step 1-7:** Identical in shape to Task 1, with these substitutions:
  - proto field number is **8** (fields 1-7 are taken)
  - package is `x/axelarcork/types`, not `x/cork/types/v2`
  - test file is `x/axelarcork/types/params_test.go`

Add to `proto/axelarcork/v1/genesis.proto` inside `message Params`:

```proto
  // cork_authority is the single address permitted to schedule, relay, bump,
  // and cancel corks. Rotated by governance via ParameterChangeProposal.
  // Empty means no address may act (fail-closed).
  string cork_authority = 8 [(gogoproto.moretags) = "yaml:\"cork_authority\""];
```

Add the same `validateCorkAuthority` function, key `KeyCorkAuthority = []byte("corkauthority")`, param pair, `DefaultParams()` entry (`CorkAuthority: ""`), and `ValidateBasic()` call as Task 1.

- [ ] **Commit**

```bash
git add proto/axelarcork/v1/genesis.proto x/axelarcork/types/genesis.pb.go x/axelarcork/types/params.go x/axelarcork/types/params_test.go
git commit -m "axelarcork: add cork_authority param"
```

---

### Task 7: Add authority cork key prefix and store accessors to x/axelarcork

**Files:**
- Modify: `x/axelarcork/types/keys.go:29-110`
- Modify: `x/axelarcork/keeper/keeper.go`
- Test: `x/axelarcork/keeper/authority_cork_test.go`

**Interfaces:**
- Consumes: Task 6
- Produces:
  - `types.AuthorityCorkKeyPrefix byte`
  - `types.GetAuthorityCorkKeyPrefix(chainID uint64) []byte`
  - `types.GetAuthorityCorkKeyByBlockHeightPrefix(chainID, blockHeight uint64) []byte`
  - `types.GetAuthorityCorkKey(chainID, blockHeight uint64, id []byte, contract common.Address) []byte`
  - `Keeper.SetAuthorityAxelarCork(ctx sdk.Context, chainID, blockHeight uint64, cork types.AxelarCork) []byte`
  - `Keeper.DeleteAuthorityAxelarCork(ctx sdk.Context, chainID, blockHeight uint64, id []byte, contract common.Address)`
  - `Keeper.IterateAuthorityAxelarCorksByBlockHeight(ctx sdk.Context, chainID, blockHeight uint64, cb func(id []byte, contract common.Address, cork types.AxelarCork) (stop bool))`

Same structure as Task 2, with the chain ID as the leading component after the prefix byte:

```go
func GetAuthorityCorkKeyPrefix(chainID uint64) []byte {
	return append([]byte{AuthorityCorkKeyPrefix}, sdk.Uint64ToBigEndian(chainID)...)
}
```

Key layout: `prefix(1) | chain_id(8) | block_height(8) | cork_id(32) | contract(20)`. Adjust offsets in the iterator accordingly.

- [ ] Write round-trip test (mirror Task 2, adding a second chain ID to prove isolation)
- [ ] Run, verify it fails
- [ ] Implement keys and accessors
- [ ] Run, verify it passes
- [ ] **Commit**

```bash
git add x/axelarcork/types/keys.go x/axelarcork/keeper/keeper.go x/axelarcork/keeper/authority_cork_test.go
git commit -m "axelarcork: add authority cork key prefix and store accessors"
```

---

### Task 8: Switch all four axelarcork messages to authority authorization

**Files:**
- Modify: `x/axelarcork/keeper/msg_server.go:25,75,154,221,251`
- Test: `x/axelarcork/keeper/msg_server_authority_test.go`

**Interfaces:**
- Consumes: Tasks 6, 7
- Produces: `ScheduleCork`, `RelayCork`, `BumpCorkGas`, `CancelScheduledCork` all require `params.CorkAuthority`

- [ ] **Step 1: Write the failing test**

One accept case and one reject case per message, plus one empty-authority reject per message — 12 cases total. Table-driven is fine and preferred here.

- [ ] **Step 2: Run, verify it fails**

Run: `go test ./x/axelarcork/keeper/... -run Authority -v`

- [ ] **Step 3: Replace the check in all four handlers**

In each of `ScheduleCork`, `RelayCork`, `BumpCorkGas`, and `CancelScheduledCork`, delete the delegate lookup:

```go
	signer := msg.MustGetSigner()
	validatorAddr := k.gravityKeeper.GetOrchestratorValidatorAddress(ctx, signer)
	if validatorAddr == nil {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "signer %s is not a delegate", signer.String())
	}
```

and replace with:

```go
	signer := msg.MustGetSigner()
	if params.CorkAuthority == "" || signer.String() != params.CorkAuthority {
		return nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized,
			"signer %s is not the cork authority", signer.String())
	}
```

`params` is already fetched in each handler via `k.GetParamSet(ctx)`; if a handler does not fetch it, add the fetch above the check.

In `ScheduleCork`, also replace `SetScheduledAxelarCork` with `SetAuthorityAxelarCork` and delete the `IncrementValidatorAxelarCorkCount` call and the `MaxAxelarCorksPerValidator` capacity check.

- [ ] **Step 4: Run, verify it passes**
- [ ] **Step 5: Commit**

```bash
git add x/axelarcork/keeper/msg_server.go x/axelarcork/keeper/msg_server_authority_test.go
git commit -m "axelarcork: authorize all four messages by cork authority param"
```

---

### Task 9: Rewrite axelarcork EndBlocker without the tally

**Files:**
- Modify: `x/axelarcork/keeper/abci.go:23-90`
- Modify: `x/axelarcork/keeper/keeper.go:400-430` (delete `GetApprovedScheduledAxelarCorks`)
- Test: `x/axelarcork/keeper/abci_authority_test.go`

**Interfaces:**
- Consumes: Tasks 7, 8
- Produces: due authority corks move into the existing `WinningAxelarCork` queue with no tally

- [ ] **Step 1: Write the failing test**

```
TestEndBlockerMovesDueCorksToWinning  — cork due at H lands in WinningAxelarCork
TestEndBlockerPerChainIsolation       — a cork on chain 10 does not appear for chain 42161
TestEndBlockerSafeModeDropsCorks      — safe mode: deleted, not marked relayable, drop event emitted
TestEndBlockerTimeoutSweepStillRuns   — an unrelayed winning cork past cork_timeout_blocks is still swept
```

- [ ] **Step 2: Run, verify it fails**

- [ ] **Step 3: Rewrite**

Inside the existing `IterateChainConfigurations` callback, replace the `GetApprovedScheduledAxelarCorks` call with a collect-then-mutate pass over `IterateAuthorityAxelarCorksByBlockHeight` for `config.Id` at the current height. Keep the existing safe-mode branch, the `SetWinningAxelarCork` call, the event emission, and the timed-out-cork sweep **exactly as they are**. Delete only the tally.

Delete `GetApprovedScheduledAxelarCorks` and the `CorkVoteThresholdStr` constant. Delete the validator cork-count helpers; leave the key helper for Task 11.

- [ ] **Step 4: Run, verify it passes**
- [ ] **Step 5: Commit**

```bash
git add x/axelarcork/keeper/abci.go x/axelarcork/keeper/keeper.go x/axelarcork/keeper/abci_authority_test.go
git commit -m "axelarcork: mark due authority corks relayable, remove power tally"
```

---

### Task 10: Decouple x/axelarcork from x/pubsub

**Files:**
- Modify: `x/axelarcork/keeper/proposal_handler.go:24-113`
- Modify: `x/axelarcork/keeper/keeper.go:38-60`
- Modify: `x/axelarcork/types/expected_keepers.go:96`
- Modify: `app/app.go`

**Interfaces:**
- Consumes: Task 9
- Produces: axelarcork `NewKeeper` no longer takes `pubsubKeeper`

Same shape as Task 5.

- [ ] Write failing test for publisher-free `AddManagedCellarsProposal`
- [ ] Run, verify it fails
- [ ] Strip `GetPublisher` / `SetDefaultSubscription` / `DeleteDefaultSubscription`, the keeper field, the interface, and the `app.go` argument
- [ ] Run `go test ./x/axelarcork/... && go build ./app/... ./x/...`
- [ ] **Commit**

```bash
git add x/axelarcork app/app.go
git commit -m "axelarcork: decouple from x/pubsub"
```

---

### Task 11: Drain the legacy queues in the v10 upgrade handler

**Files:**
- Create: `app/upgrades/v10/migrations.go`
- Modify: `app/upgrades/v10/upgrades.go:23-60`
- Test: `app/upgrades/v10/migrations_test.go`

**Interfaces:**
- Consumes: Tasks 2, 4, 7, 9
- Produces: `DrainLegacyCorkQueues(ctx sdk.Context, corkStoreKey, axelarcorkStoreKey storetypes.StoreKey) int`
- Produces: `deletePrefix(store storetypes.KVStore, prefix []byte) int`

Note on the signature: the drain takes **store keys, not keepers**, and iterates
raw prefixes. The typed iterators (`IterateScheduledCorks`,
`GetApprovedScheduledCorks`) are deleted in Tasks 4 and 9, so the migration
cannot depend on them. Param seeding in Step 4 does use the keepers — both
`corkkeeper.Keeper` and `axelarcorkkeeper.Keeper` expose `GetParamSet(ctx)` and
`SetParams(ctx, params)`. `CreateUpgradeHandler` therefore takes both keepers
**and** both store keys; `app/app.go` has `keys[corktypes.StoreKey]` and
`keys[axelarcorktypes.StoreKey]` in scope at the registration site.

**This is the highest-consequence task in the plan.** Once the tally is gone, nothing else deletes from the legacy queues. Any cork left behind is stranded in state permanently.

- [ ] **Step 1: Write the failing test**

```go
func TestDrainLegacyCorkQueuesRemovesEverything(t *testing.T) {
	// Seed: 3 legacy scheduled corks in x/cork at differing heights and validators.
	// Seed: 2 legacy scheduled corks per chain across 2 axelarcork chains.
	// Seed: validator cork-count keys for each validator involved.
	// Run DrainLegacyCorkQueues.
	// Assert: iterating both legacy prefixes yields zero entries.
	// Assert: no ValidatorCorkCountKey / ValidatorAxelarCorkCountKey entries remain.
	// Assert: the returned count equals the number of corks seeded.
	// Assert: authority-cork prefixes are untouched (drain must not migrate).
}

func TestDrainLegacyCorkQueuesIsIdempotent(t *testing.T) {
	// Running it twice on an already-drained store returns 0 and does not panic.
}
```

- [ ] **Step 2: Run, verify it fails**

Run: `go test ./app/upgrades/v10/... -run TestDrainLegacy -v`
Expected: FAIL — `undefined: DrainLegacyCorkQueues`

- [ ] **Step 3: Implement the drain**

Create `app/upgrades/v10/migrations.go`. Iterate raw store prefixes rather than typed helpers — the typed iterators are being deleted, and raw iteration is what makes this robust:

```go
// DrainLegacyCorkQueues deletes every cork left in the pre-v10
// validator-scheduled queues, along with the per-validator cork counters.
//
// This MUST run in the v10 handler. The power tally in each module's EndBlocker
// was the only site that deleted from these prefixes; with the tally removed,
// anything left behind becomes permanently undeletable state.
//
// Corks are DROPPED rather than migrated to the authority queue: they were
// scheduled under validator consent that no longer carries meaning, and
// re-scheduling under the authority key costs one transaction.
func DrainLegacyCorkQueues(ctx sdk.Context, corkStoreKey, axelarcorkStoreKey storetypes.StoreKey) int {
	drained := 0
	drained += deletePrefix(ctx.KVStore(corkStoreKey), []byte{corktypes.ScheduledCorkKeyPrefix})
	drained += deletePrefix(ctx.KVStore(corkStoreKey), []byte{corktypes.ValidatorCorkCountKey})
	drained += deletePrefix(ctx.KVStore(axelarcorkStoreKey), []byte{axelarcorktypes.ScheduledCorkKeyPrefix})
	drained += deletePrefix(ctx.KVStore(axelarcorkStoreKey), []byte{axelarcorktypes.ValidatorAxelarCorkCountKey})
	return drained
}

// deletePrefix removes every key under prefix. Keys are collected before
// deletion because deleting through a live iterator is undefined behaviour.
func deletePrefix(store storetypes.KVStore, prefix []byte) int {
	var keys [][]byte
	iter := sdk.KVStorePrefixIterator(store, prefix)
	for ; iter.Valid(); iter.Next() {
		key := make([]byte, len(iter.Key()))
		copy(key, iter.Key())
		keys = append(keys, key)
	}
	iter.Close()

	for _, k := range keys {
		store.Delete(k)
	}
	return len(keys)
}
```

The store keys must be threaded into `CreateUpgradeHandler`. Update its signature and the `app.go` call site at `app/app.go:1207-1213`.

- [ ] **Step 4: Seed the authority params in the same handler**

In `app/upgrades/v10/upgrades.go`, after the PoA seeding block, before `RunMigrations`:

```go
		corkParams := corkkeeper.GetParamSet(ctx)
		corkParams.CorkAuthority = CorkAuthorityAddress
		corkkeeper.SetParams(ctx, corkParams)

		axelarcorkParams := axelarcorkkeeper.GetParamSet(ctx)
		axelarcorkParams.CorkAuthority = CorkAuthorityAddress
		axelarcorkkeeper.SetParams(ctx, axelarcorkParams)

		drained := DrainLegacyCorkQueues(ctx, corkStoreKey, axelarcorkStoreKey)
		ctx.Logger().Info("v10 upgrade: cork authority seeded and legacy queues drained",
			"authority", CorkAuthorityAddress, "drained_keys", drained)
```

Add to `app/upgrades/v10/constants.go`:

```go
// CorkAuthorityAddress is the address seeded as the sole cork authority for
// x/cork and x/axelarcork at the v10 upgrade. Rotated afterwards by governance
// via ParameterChangeProposal.
const CorkAuthorityAddress = "somm1lcsjy2d5s33h0sddd8lpuqvwyz5ruz7ju4aeqa"
```

Add a guard at the top of the handler, mirroring the existing `DefaultAuthorityValidators` guard:

```go
		if _, err := sdk.AccAddressFromBech32(CorkAuthorityAddress); err != nil {
			return vm, fmt.Errorf("v10: invalid CorkAuthorityAddress %q: %w", CorkAuthorityAddress, err)
		}
```

- [ ] **Step 5: Run, verify it passes**

Run: `go test ./app/upgrades/v10/... -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add app/upgrades/v10 app/app.go
git commit -m "v10: seed cork authority and drain legacy cork queues"
```

---

### Task 12: Full-suite regression and integration coverage

**Files:**
- Modify: `integration_tests/` — cork and axelarcork scenarios
- Test: whole repo

**Interfaces:**
- Consumes: all prior tasks
- Produces: green suite at or above the 23 ok / 0 FAIL baseline

- [ ] **Step 1: Run the full unit suite**

Run: `go test $(go list ./... | grep -v integration_tests)`
Expected: 0 FAIL. Compare the `ok` count against the 23 baseline; it should be ≥ 23.

- [ ] **Step 2: Update integration tests**

The existing integration tests schedule corks from validator delegates. Find them:

Run: `grep -rn "ScheduleCork\|MsgScheduleCorkRequest" integration_tests/`

Rewrite each to sign from the authority account, and set `cork_authority` in the test genesis (`integration_tests/genesis.go`).

- [ ] **Step 3: Add an authority-rejection integration case**

Assert that a cork signed by a validator delegate is now rejected with "is not the cork authority".

- [ ] **Step 4: Run integration tests**

Run: `go test ./integration_tests/... -v` (requires Docker)
Expected: PASS

- [ ] **Step 5: Vet everything**

Run: `go vet ./...`
Expected: exit 0

- [ ] **Step 6: Commit**

```bash
git add integration_tests
git commit -m "test: drive cork integration tests from the cork authority"
```

---

### Task 13: Mainnet-state migration rehearsal

**Files:**
- Create: `docs/superpowers/plans/2026-08-20-v10-migration-rehearsal.md` (results log)

**Interfaces:**
- Consumes: Task 11
- Produces: evidence that the drain is correct against real state

This is the check that the spec calls out as non-negotiable, and it cannot be satisfied by unit tests.

- [ ] **Step 1: Export mainnet state**

On `sommelier-authority-1`:

```
sudo systemctl stop sommelier
sudo -u validator /home/validator/.local/bin/sommelier export --home /home/validator/.sommelier > /tmp/somm-export.json
sudo systemctl start sommelier
```

Take the node out of the active set first, or use a non-validating node — do not stop a bonded validator without checking the effect on the PoA authority set.

- [ ] **Step 2: Count legacy queue entries in the export**

```
jq '[.app_state.cork.scheduled_corks[]] | length' /tmp/somm-export.json
jq '[.app_state.axelarcork.scheduled_corks[]] | length' /tmp/somm-export.json
```

Record both numbers. These are what the drain must remove.

- [ ] **Step 3: Run the upgrade against the export**

Start a single-node chain from the export with the v9.4.0 binary, submit the v10 plan at a near height, swap to the v10 binary, and let it upgrade.

- [ ] **Step 4: Assert the drain**

After the upgrade height, export again and assert both counts are zero, and that `cork_authority` is set in both modules' params.

- [ ] **Step 5: Exercise the authority path**

Schedule a cork from the authority key against a real managed cellar, targeting a near-future height. Confirm it executes and that a delegate-signed cork is rejected.

- [ ] **Step 6: Record results and commit**

Write the observed counts and outcomes to the results log and commit it.

---

## Self-Review Notes

**Spec coverage:** Section 1 (authorization) → Tasks 1, 3, 6, 8. Section 2 (cork storage/execution) → Tasks 2, 4. Section 3 (axelarcork) → Tasks 7, 8, 9. Section 4 (pubsub) → Tasks 5, 10. Section 5 (migration) → Task 11. Section 6 (safe mode) → preserved in Tasks 4 and 9, tested in both. Testing section → Tasks 12, 13.

**Known deviation:** the spec lists a `CorkResult` decision (keep queries, stop writing). That is implicit in Task 4's deletion of the tally, which was the only writer. No separate task is needed, but the implementer should confirm the `cork-result` gRPC query still compiles and returns historical records.

**Open risk carried from the spec:** v10 ships this alongside the PoA power floor and the Gravity GHSA fixes. Task 13 is the mitigation and must not be cut for schedule.
