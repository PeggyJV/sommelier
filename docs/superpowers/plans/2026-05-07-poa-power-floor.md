# Sommelier PoA Power Floor — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a new `x/poa` module that guarantees a binary-specified authority validator set holds ≥67% of consensus power on every block, while keeping community delegation and rewards working.

**Architecture:** A new Cosmos SDK module wraps `stakingkeeper.Keeper`, rescales validator power for an "authority allowlist" via a multiplier in EndBlocker (after `x/staking`), and is wired into every downstream consumer of validator power (gravity-bridge, slashing, distribution, evidence, cork, axelarcork, pubsub, incentives). `Slash` is normalised so authority validators are penalised on raw stake. Shipped via v10 upgrade.

**Tech Stack:** Go 1.22, Cosmos SDK v0.47.15, CometBFT v0.37.18, gravity-bridge v6, gogoproto, buf.

**Spec:** [`docs/superpowers/specs/2026-05-07-poa-power-floor-design.md`](../specs/2026-05-07-poa-power-floor-design.md)

---

## File Structure

### Files created

```
proto/poa/v1/poa.proto          # AuthoritySet entry, MultiplierSnapshot
proto/poa/v1/genesis.proto      # GenesisState
proto/poa/v1/params.proto       # Params
proto/poa/v1/tx.proto           # MsgUpdateAuthoritySet, MsgUpdateParams
proto/poa/v1/query.proto        # AuthoritySet, Params, EffectivePower

x/poa/types/keys.go             # ModuleName, StoreKey, KV prefixes
x/poa/types/codec.go            # RegisterLegacyAminoCodec, RegisterInterfaces
x/poa/types/errors.go
x/poa/types/events.go
x/poa/types/expected_keepers.go # interface for the underlying staking keeper
x/poa/types/params.go           # Params, validation, defaults
x/poa/types/genesis.go          # default genesis, validate
x/poa/types/multiplier.go       # ComputeMultiplier (pure function)
x/poa/types/multiplier_test.go
x/poa/types/msgs.go             # MsgUpdateAuthoritySet, MsgUpdateParams + ValidateBasic
x/poa/types/msgs_test.go
x/poa/types/poa.pb.go           # generated
x/poa/types/genesis.pb.go       # generated
x/poa/types/params.pb.go        # generated
x/poa/types/tx.pb.go            # generated
x/poa/types/query.pb.go         # generated

x/poa/keeper/keeper.go          # Keeper struct, NewKeeper
x/poa/keeper/authority.go       # SetAuthoritySet/GetAuthoritySet/IsAuthority
x/poa/keeper/authority_test.go
x/poa/keeper/params.go
x/poa/keeper/snapshot.go        # SetMultiplierSnapshot/GetMultiplierSnapshot/PruneSnapshots
x/poa/keeper/snapshot_test.go
x/poa/keeper/wrapper.go         # WrappedStakingKeeper (the rescaling adapter)
x/poa/keeper/wrapper_test.go
x/poa/keeper/abci.go            # EndBlocker: compute, overwrite, emit, snapshot, prune
x/poa/keeper/abci_test.go
x/poa/keeper/msg_server.go      # gov msg handlers
x/poa/keeper/grpc_query.go
x/poa/keeper/genesis.go         # InitGenesis/ExportGenesis
x/poa/keeper/test_helpers.go    # simapp setup with PoA wired

x/poa/module.go                 # AppModule, AppModuleBasic, EndBlocker entry
x/poa/abci.go                   # thin wrapper for module.go EndBlocker

app/upgrades/v10/constants.go   # Upgrade name, store upgrades, default authority list
app/upgrades/v10/upgrades.go    # Upgrade handler

docs/poa.md                     # Operator-facing brief on PoA semantics
```

### Files modified

```
app/app.go                      # register x/poa, wrap StakingKeeper for downstream consumers, EndBlocker order
app/upgrades.go                 # register v10 upgrade
proto/buf.yaml or buf.work.yaml # add poa proto (only if buf modules require explicit listing)
go.mod                          # nothing new expected; verify
Makefile                        # if proto-gen target lists modules, add poa
```

### Module boundaries

- `types/multiplier.go` is pure math (no SDK deps beyond `math.Int`/`sdk.Dec`) so it can be unit-tested without simapp.
- `keeper/wrapper.go` is the only file other modules' keepers ever see; it is the durable surface for the 67% invariant.
- `keeper/abci.go` is the imperative side: it both rewrites staking storage and emits ABCI updates. Keep it small; helpers live in `wrapper.go` and `snapshot.go`.

---

## Task Sequencing

Tasks 0–8 build the module in isolation (no app wiring). Task 9 wires it. Tasks 10–11 add cross-module tests and the upgrade. Task 12 ships docs.

Each task ends with a green test run and a commit.

---

## Task 0: Audit downstream `StakingKeeper` interfaces

**Files:** none modified — produces `x/poa/types/expected_keepers.go` content for Task 2.

The wrapper must satisfy every `StakingKeeper`-shaped interface consumed in the app. Discovering these via compile errors during Task 9 is slow and prone to omissions.

- [ ] **Step 1: Enumerate every consumer interface**

Run: `grep -rn "type StakingKeeper interface" x/ third_party/ vendor/github.com/peggyjv/gravity-bridge third_party/proto 2>/dev/null` and inspect:
- `x/cork/types/expected_keepers.go`
- `x/axelarcork/types/expected_keepers.go`
- `x/pubsub/types/expected_keepers.go`
- `x/incentives/types/expected_keepers.go`
- `github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types` (use `go doc` if vendored)
- Cosmos SDK `x/slashing/types`, `x/distribution/types`, `x/evidence/types`, `x/genutil/types`

- [ ] **Step 2: Produce a checklist of method signatures**

Write the union as a single Go file or markdown table; commit nothing in this task — it feeds Task 2 Step 3.

- [ ] **Step 3: Identify which methods need rescaling vs pass-through**

| Method | Behavior |
|--------|----------|
| `GetLastValidatorPower(op)` | rescaled (rely on overwritten staking storage) |
| `GetLastTotalPower()` | rescaled (same store) |
| `GetBondedValidatorsByPower()` | **rescaled** — returns `[]Validator` whose `Tokens` field must reflect boost |
| `IterateBondedValidatorsByPower(cb)` | **rescaled** — yields a `ValidatorI` adapter |
| `IterateLastValidators(cb)` | **rescaled** — yields adapter |
| `IterateValidators(cb)` | **rescaled** — yields adapter |
| `IterateLastValidatorPowers(cb)` | **rescaled** — yields boosted power int |
| `Validator(op) ValidatorI` | **rescaled** — returns adapter |
| `ValidatorByConsAddr(cons) ValidatorI` | **rescaled** — returns adapter |
| `GetValidator(op) (Validator, bool)` | **rescaled** — returns Validator with boosted Tokens |
| `Slash` / `SlashWithInfractionReason` | normalised (divide by snapshot multiplier) |
| `Jail` / `Unjail` | pass-through |
| `Delegation`, `MaxValidators`, `PowerReduction`, `BondDenom`, `UnbondingTime`, `IsValidatorJailed`, `GetParams`, `Hooks`, `ValidatorQueueIterator`, `SetLastValidatorPower`, `DeleteLastValidatorPower` | pass-through |

- [ ] **Step 4: Commit nothing — feeds Task 2 directly.**

---

## Task 1: Scaffold proto definitions

**Files:**
- Create: `proto/poa/v1/poa.proto`, `proto/poa/v1/params.proto`, `proto/poa/v1/genesis.proto`, `proto/poa/v1/tx.proto`, `proto/poa/v1/query.proto`
- Modify: `Makefile` if proto-gen target enumerates modules

- [ ] **Step 1: Author `params.proto`**

```proto
syntax = "proto3";
package poa.v1;
option go_package = "github.com/peggyjv/sommelier/v9/x/poa/types";

import "gogoproto/gogo.proto";
import "amino/amino.proto";

message Params {
  string floor_fraction = 1 [
    (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
    (gogoproto.nullable) = false
  ];
  bool enabled = 2;
  bool halt_when_authority_empty = 3;
}
```

- [ ] **Step 2: Author `poa.proto`**

```proto
syntax = "proto3";
package poa.v1;
option go_package = "github.com/peggyjv/sommelier/v9/x/poa/types";

message AuthorityValidator {
  // operator address (sdk.ValAddress bech32)
  string operator_address = 1;
}

// AuthoritySetWrapper is the on-disk representation written under AuthoritySetKey.
message AuthoritySetWrapper {
  repeated AuthorityValidator validators = 1;
}

message MultiplierEntry {
  string operator_address = 1;
  // dec fixed-point string
  string multiplier = 2;
}

message MultiplierSnapshot {
  int64 height = 1;
  repeated MultiplierEntry entries = 2;
}
```

- [ ] **Step 3: Author `genesis.proto`**

```proto
syntax = "proto3";
package poa.v1;
option go_package = "github.com/peggyjv/sommelier/v9/x/poa/types";

import "gogoproto/gogo.proto";
import "poa/v1/params.proto";
import "poa/v1/poa.proto";

message GenesisState {
  Params params = 1 [(gogoproto.nullable) = false];
  repeated AuthorityValidator authority_set = 2 [(gogoproto.nullable) = false];
}
```

- [ ] **Step 4: Author `tx.proto`**

```proto
syntax = "proto3";
package poa.v1;
option go_package = "github.com/peggyjv/sommelier/v9/x/poa/types";

import "cosmos/msg/v1/msg.proto";
import "gogoproto/gogo.proto";
import "poa/v1/params.proto";

service Msg {
  rpc UpdateAuthoritySet(MsgUpdateAuthoritySet) returns (MsgUpdateAuthoritySetResponse);
  rpc UpdateParams(MsgUpdateParams) returns (MsgUpdateParamsResponse);
}

message MsgUpdateAuthoritySet {
  option (cosmos.msg.v1.signer) = "authority";
  string authority = 1;
  repeated string validators = 2;  // bech32 ValAddresses
}
message MsgUpdateAuthoritySetResponse {}

message MsgUpdateParams {
  option (cosmos.msg.v1.signer) = "authority";
  string authority = 1;
  Params params = 2 [(gogoproto.nullable) = false];
}
message MsgUpdateParamsResponse {}
```

- [ ] **Step 5: Author `query.proto`**

```proto
syntax = "proto3";
package poa.v1;
option go_package = "github.com/peggyjv/sommelier/v9/x/poa/types";

import "google/api/annotations.proto";
import "gogoproto/gogo.proto";
import "poa/v1/params.proto";

service Query {
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse) {
    option (google.api.http).get = "/sommelier/poa/v1/params";
  }
  rpc AuthoritySet(QueryAuthoritySetRequest) returns (QueryAuthoritySetResponse) {
    option (google.api.http).get = "/sommelier/poa/v1/authority_set";
  }
  rpc EffectivePower(QueryEffectivePowerRequest) returns (QueryEffectivePowerResponse) {
    option (google.api.http).get = "/sommelier/poa/v1/effective_power/{operator_address}";
  }
}

message QueryParamsRequest {}
message QueryParamsResponse { Params params = 1 [(gogoproto.nullable) = false]; }
message QueryAuthoritySetRequest {}
message QueryAuthoritySetResponse { repeated string validators = 1; }
message QueryEffectivePowerRequest { string operator_address = 1; }
message QueryEffectivePowerResponse { int64 power = 1; bool is_authority = 2; }
```

- [ ] **Step 6: Generate code**

Run: `make proto-gen` (or follow the existing repo target).
Expected: new `*.pb.go` files in `x/poa/types/`. Verify with `git status`.

- [ ] **Step 7: Commit**

```bash
git add proto/poa x/poa/types/*.pb.go
git commit -m "poa: scaffold proto definitions for x/poa module"
```

---

## Task 2: Module skeleton (compiles, not yet wired)

**Files:**
- Create: `x/poa/types/keys.go`, `x/poa/types/codec.go`, `x/poa/types/errors.go`, `x/poa/types/events.go`, `x/poa/types/expected_keepers.go`, `x/poa/types/params.go`, `x/poa/types/genesis.go`, `x/poa/keeper/keeper.go`, `x/poa/module.go`

- [ ] **Step 1: `types/keys.go`**

```go
package types

const (
    ModuleName = "poa"
    StoreKey   = ModuleName
    RouterKey  = ModuleName
    QuerierRoute = ModuleName
)

var (
    ParamsKey            = []byte{0x01}
    AuthoritySetKey      = []byte{0x02}
    MultiplierSnapshotPrefix = []byte{0x03} // followed by big-endian int64 height
)

func MultiplierSnapshotKey(height int64) []byte {
    bz := make([]byte, 8)
    // big-endian for natural sort order
    binary.BigEndian.PutUint64(bz, uint64(height))
    return append(MultiplierSnapshotPrefix, bz...)
}
```

(Add `import "encoding/binary"`.)

- [ ] **Step 2: `types/params.go`**

```go
package types

import (
    "fmt"
    sdk "github.com/cosmos/cosmos-sdk/types"
    paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
    KeyFloorFraction          = []byte("FloorFraction")
    KeyEnabled                = []byte("Enabled")
    KeyHaltWhenAuthorityEmpty = []byte("HaltWhenAuthorityEmpty")

    DefaultFloorFraction = sdk.MustNewDecFromStr("0.670000000000000001")
)

func ParamKeyTable() paramtypes.KeyTable { return paramtypes.NewKeyTable().RegisterParamSet(&Params{}) }

func DefaultParams() Params {
    return Params{
        FloorFraction:           DefaultFloorFraction,
        Enabled:                 true,
        HaltWhenAuthorityEmpty:  true,
    }
}

func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
    return paramtypes.ParamSetPairs{
        paramtypes.NewParamSetPair(KeyFloorFraction, &p.FloorFraction, validateFloor),
        paramtypes.NewParamSetPair(KeyEnabled, &p.Enabled, validateBool),
        paramtypes.NewParamSetPair(KeyHaltWhenAuthorityEmpty, &p.HaltWhenAuthorityEmpty, validateBool),
    }
}

func (p Params) Validate() error { return validateFloor(p.FloorFraction) }

func validateFloor(i interface{}) error {
    v, ok := i.(sdk.Dec)
    if !ok { return fmt.Errorf("invalid type for floor: %T", i) }
    if !v.GT(sdk.ZeroDec()) || v.GTE(sdk.OneDec()) {
        return fmt.Errorf("floor must be in (0,1), got %s", v)
    }
    if v.LT(sdk.MustNewDecFromStr("0.5")) {
        return fmt.Errorf("floor must be > 0.5, got %s", v)
    }
    return nil
}
func validateBool(i interface{}) error { _, ok := i.(bool); if !ok { return fmt.Errorf("bool expected") }; return nil }
```

- [ ] **Step 3: `types/expected_keepers.go`**

Document the staking-keeper interface union the wrapper will implement. List every method consumed by gravity-bridge, slashing, distribution, cork, axelarcork, pubsub, incentives, evidence (audit each `expected_keepers.go` and copy method signatures into one interface).

```go
package types

import (
    "time"
    "cosmossdk.io/math"
    sdk "github.com/cosmos/cosmos-sdk/types"
    stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper is the underlying staking module functionality the wrapper depends on.
// We only depend on the concrete keeper through this interface to keep tests fast.
type StakingKeeper interface {
    GetParams(ctx sdk.Context) stakingtypes.Params
    GetValidator(ctx sdk.Context, addr sdk.ValAddress) (stakingtypes.Validator, bool)
    GetLastValidatorPower(ctx sdk.Context, operator sdk.ValAddress) int64
    GetLastTotalPower(ctx sdk.Context) math.Int
    GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
    IterateValidators(sdk.Context, func(int64, stakingtypes.ValidatorI) bool)
    IterateBondedValidatorsByPower(sdk.Context, func(int64, stakingtypes.ValidatorI) bool)
    IterateLastValidators(sdk.Context, func(int64, stakingtypes.ValidatorI) bool)
    IterateLastValidatorPowers(sdk.Context, func(sdk.ValAddress, int64) bool)
    Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI
    ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingtypes.ValidatorI
    Slash(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec) math.Int
    SlashWithInfractionReason(sdk.Context, sdk.ConsAddress, int64, int64, sdk.Dec, stakingtypes.Infraction) math.Int
    Jail(sdk.Context, sdk.ConsAddress)
    Unjail(sdk.Context, sdk.ConsAddress)
    Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) stakingtypes.DelegationI
    MaxValidators(sdk.Context) uint32
    PowerReduction(sdk.Context) math.Int
    BondDenom(sdk.Context) string
    UnbondingTime(sdk.Context) time.Duration
    IsValidatorJailed(sdk.Context, sdk.ConsAddress) bool
    SetLastValidatorPower(sdk.Context, sdk.ValAddress, int64) // for EndBlocker overwrite
    DeleteLastValidatorPower(sdk.Context, sdk.ValAddress)
    ValidatorQueueIterator(sdk.Context, time.Time, int64) sdk.Iterator
    Hooks() stakingtypes.StakingHooks
}
```

(If any method is missing, add it as you discover compile errors — this is the union surface.)

- [ ] **Step 4: `types/codec.go`** registering `MsgUpdateAuthoritySet` and `MsgUpdateParams`. Mirror `x/cork/types/codec.go`.

- [ ] **Step 5: `types/genesis.go`**

```go
func DefaultGenesis() *GenesisState { return &GenesisState{Params: DefaultParams()} }
func (gs GenesisState) Validate() error { return gs.Params.Validate() }
```

- [ ] **Step 6: `keeper/keeper.go` minimal**

```go
type Keeper struct {
    cdc        codec.BinaryCodec
    storeKey   storetypes.StoreKey
    paramSpace paramtypes.Subspace
    sk         types.StakingKeeper
    authority  string // gov module address (bech32)
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, ps paramtypes.Subspace, sk types.StakingKeeper, authority string) Keeper {
    if !ps.HasKeyTable() { ps = ps.WithKeyTable(types.ParamKeyTable()) }
    return Keeper{cdc, key, ps, sk, authority}
}

func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) { k.paramSpace.GetParamSet(ctx, &p); return }
func (k Keeper) SetParams(ctx sdk.Context, p types.Params)  { k.paramSpace.SetParamSet(ctx, &p) }
```

- [ ] **Step 7: `module.go`** mirror `x/cork/module.go` but minimal (no msg server / query server bound yet — just enough to compile and register module name).

- [ ] **Step 8: Verify build**

Run: `go build ./x/poa/...`
Expected: success.

- [ ] **Step 9: Commit**

```bash
git add x/poa/
git commit -m "poa: scaffold module skeleton (types, keeper, module)"
```

---

## Task 3: Multiplier math — pure function with TDD

**Files:**
- Create: `x/poa/types/multiplier.go`, `x/poa/types/multiplier_test.go`

- [ ] **Step 1: Write failing tests first**

```go
package types_test

import (
    "testing"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/peggyjv/sommelier/v9/x/poa/types"
    "github.com/stretchr/testify/require"
)

func TestComputeMultiplier_BoostNeeded(t *testing.T) {
    // authority=100, community=300, floor=0.67 → M = (0.67/0.33)*(300/100) ≈ 6.0909
    f := sdk.MustNewDecFromStr("0.67")
    m := types.ComputeMultiplier(sdk.NewInt(100), sdk.NewInt(300), f)
    require.True(t, m.GT(sdk.OneDec()))
    // resulting authority share check:
    boosted := sdk.NewDecFromInt(sdk.NewInt(100)).Mul(m)
    total := boosted.Add(sdk.NewDecFromInt(sdk.NewInt(300)))
    share := boosted.Quo(total)
    require.True(t, share.GTE(f), "share %s < floor %s", share, f)
}

func TestComputeMultiplier_AlreadyAboveFloor(t *testing.T) {
    f := sdk.MustNewDecFromStr("0.67")
    // authority=900, community=100 → already 90% → multiplier clamps to 1
    m := types.ComputeMultiplier(sdk.NewInt(900), sdk.NewInt(100), f)
    require.Equal(t, sdk.OneDec(), m)
}

func TestComputeMultiplier_ZeroCommunity(t *testing.T) {
    f := sdk.MustNewDecFromStr("0.67")
    m := types.ComputeMultiplier(sdk.NewInt(500), sdk.ZeroInt(), f)
    require.Equal(t, sdk.OneDec(), m)
}

func TestComputeMultiplier_ZeroAuthority(t *testing.T) {
    f := sdk.MustNewDecFromStr("0.67")
    m := types.ComputeMultiplier(sdk.ZeroInt(), sdk.NewInt(100), f)
    require.True(t, m.IsZero())
}

func TestComputeMultiplier_FloorEdges(t *testing.T) {
    f := sdk.MustNewDecFromStr("0.5")
    // authority=community=100, floor=0.5 → already at floor, M=1
    m := types.ComputeMultiplier(sdk.NewInt(100), sdk.NewInt(100), f)
    require.Equal(t, sdk.OneDec(), m)
}
```

- [ ] **Step 2: Run tests — expect compile failure**

Run: `go test ./x/poa/types/ -run TestComputeMultiplier`
Expected: undefined: types.ComputeMultiplier

- [ ] **Step 3: Implement `multiplier.go`**

```go
package types

import (
    "cosmossdk.io/math"
    sdk "github.com/cosmos/cosmos-sdk/types"
)

// ComputeMultiplier returns the per-authority-validator power multiplier
// M such that (B*M) / (B*M + C) >= floor.
//   M = (floor / (1-floor)) * (C / B)
// Returns 1 if B is already large enough or C == 0.
// Returns 0 if B == 0 (caller decides halt vs. pass-through).
func ComputeMultiplier(authorityPower, communityPower math.Int, floor sdk.Dec) sdk.Dec {
    if authorityPower.IsZero() {
        return sdk.ZeroDec()
    }
    if communityPower.IsZero() {
        return sdk.OneDec()
    }
    // current authority share
    total := authorityPower.Add(communityPower)
    share := sdk.NewDecFromInt(authorityPower).Quo(sdk.NewDecFromInt(total))
    if share.GTE(floor) {
        return sdk.OneDec()
    }
    // M = (f / (1-f)) * (C / B)
    ratio := floor.Quo(sdk.OneDec().Sub(floor))
    cb := sdk.NewDecFromInt(communityPower).Quo(sdk.NewDecFromInt(authorityPower))
    return ratio.Mul(cb)
}
```

- [ ] **Step 4: Run tests**

Run: `go test ./x/poa/types/ -run TestComputeMultiplier -v`
Expected: PASS for all five.

- [ ] **Step 5: Commit**

```bash
git add x/poa/types/multiplier.go x/poa/types/multiplier_test.go
git commit -m "poa: ComputeMultiplier with unit tests"
```

---

## Task 4: Authority set storage

**Files:**
- Create: `x/poa/keeper/authority.go`, `x/poa/keeper/authority_test.go`, `x/poa/keeper/test_helpers.go`

- [ ] **Step 1: `keeper/test_helpers.go`** — minimal in-memory keeper for unit tests. Use `testutil.DefaultContextWithDB` style (look at `x/cork/keeper/test_common.go`). Provide `NewTestKeeper(t) (Keeper, sdk.Context)` returning a keeper with a mock `StakingKeeper` and an in-memory store.

- [ ] **Step 2: Write failing tests**

```go
func TestAuthoritySet_RoundTrip(t *testing.T) {
    k, ctx := keeper.NewTestKeeper(t)
    addrs := []sdk.ValAddress{
        sdk.ValAddress("validator-aaaaaa"),
        sdk.ValAddress("validator-bbbbbb"),
    }
    k.SetAuthoritySet(ctx, addrs)
    got := k.GetAuthoritySet(ctx)
    require.ElementsMatch(t, addrs, got)
    require.True(t, k.IsAuthority(ctx, addrs[0]))
    require.False(t, k.IsAuthority(ctx, sdk.ValAddress("notinset-xxxxxx")))
}

func TestAuthoritySet_Empty(t *testing.T) {
    k, ctx := keeper.NewTestKeeper(t)
    require.Empty(t, k.GetAuthoritySet(ctx))
    require.False(t, k.IsAuthority(ctx, sdk.ValAddress("anyone-xxxxxx")))
}
```

- [ ] **Step 3: Run tests — expect failures**

Run: `go test ./x/poa/keeper/ -run TestAuthoritySet`
Expected: undefined methods.

- [ ] **Step 4: Implement `authority.go`**

```go
func (k Keeper) SetAuthoritySet(ctx sdk.Context, vals []sdk.ValAddress) {
    store := ctx.KVStore(k.storeKey)
    list := make([]types.AuthorityValidator, len(vals))
    for i, v := range vals {
        list[i] = types.AuthorityValidator{OperatorAddress: v.String()}
    }
    bz := k.cdc.MustMarshal(&types.AuthoritySetWrapper{Validators: list})
    store.Set(types.AuthoritySetKey, bz)
}

func (k Keeper) GetAuthoritySet(ctx sdk.Context) []sdk.ValAddress {
    store := ctx.KVStore(k.storeKey)
    bz := store.Get(types.AuthoritySetKey)
    if bz == nil { return nil }
    var w types.AuthoritySetWrapper
    k.cdc.MustUnmarshal(bz, &w)
    out := make([]sdk.ValAddress, len(w.Validators))
    for i, v := range w.Validators { out[i], _ = sdk.ValAddressFromBech32(v.OperatorAddress) }
    return out
}

func (k Keeper) IsAuthority(ctx sdk.Context, val sdk.ValAddress) bool {
    for _, a := range k.GetAuthoritySet(ctx) {
        if a.Equals(val) { return true }
    }
    return false
}
```

(Add an `AuthoritySetWrapper` proto message in `poa.proto` with `repeated AuthorityValidator validators = 1;` for the on-disk representation.)

- [ ] **Step 5: Run tests**

Run: `go test ./x/poa/keeper/ -run TestAuthoritySet -v`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add x/poa/ proto/poa
git commit -m "poa: authority set storage + tests"
```

---

## Task 5: Multiplier snapshot store

**Files:**
- Create: `x/poa/keeper/snapshot.go`, `x/poa/keeper/snapshot_test.go`

- [ ] **Step 1: Write failing tests**

```go
func TestSnapshot_RoundTrip(t *testing.T) {
    k, ctx := keeper.NewTestKeeper(t)
    snap := types.MultiplierSnapshot{
        Height: 100,
        Entries: []types.MultiplierEntry{
            {OperatorAddress: "somm1...", Multiplier: "2.5"},
        },
    }
    k.SetMultiplierSnapshot(ctx, snap)
    got, found := k.GetMultiplierSnapshot(ctx, 100)
    require.True(t, found)
    require.Equal(t, snap, got)
}

func TestSnapshot_Prune(t *testing.T) {
    k, ctx := keeper.NewTestKeeper(t)
    for h := int64(1); h <= 50; h++ {
        k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{Height: h})
    }
    k.PruneSnapshotsBefore(ctx, 30)
    _, found := k.GetMultiplierSnapshot(ctx, 29)
    require.False(t, found)
    _, found = k.GetMultiplierSnapshot(ctx, 30)
    require.True(t, found)
    _, found = k.GetMultiplierSnapshot(ctx, 50)
    require.True(t, found)
}
```

- [ ] **Step 2: Run tests — expect failure**

- [ ] **Step 3: Implement `snapshot.go`**

```go
func (k Keeper) SetMultiplierSnapshot(ctx sdk.Context, s types.MultiplierSnapshot) {
    store := ctx.KVStore(k.storeKey)
    store.Set(types.MultiplierSnapshotKey(s.Height), k.cdc.MustMarshal(&s))
}

func (k Keeper) GetMultiplierSnapshot(ctx sdk.Context, height int64) (types.MultiplierSnapshot, bool) {
    bz := ctx.KVStore(k.storeKey).Get(types.MultiplierSnapshotKey(height))
    if bz == nil { return types.MultiplierSnapshot{}, false }
    var s types.MultiplierSnapshot
    k.cdc.MustUnmarshal(bz, &s)
    return s, true
}

// MultiplierForValidator returns the boost applied to `val` at `height`, or 1 if no snapshot or not present.
func (k Keeper) MultiplierForValidator(ctx sdk.Context, val sdk.ValAddress, height int64) sdk.Dec {
    s, ok := k.GetMultiplierSnapshot(ctx, height)
    if !ok { return sdk.OneDec() }
    addr := val.String()
    for _, e := range s.Entries {
        if e.OperatorAddress == addr {
            d, err := sdk.NewDecFromStr(e.Multiplier)
            if err != nil { return sdk.OneDec() }
            return d
        }
    }
    return sdk.OneDec()
}

func (k Keeper) PruneSnapshotsBefore(ctx sdk.Context, height int64) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.MultiplierSnapshotPrefix)
    iter := store.Iterator(nil, nil)
    defer iter.Close()
    for ; iter.Valid(); iter.Next() {
        h := int64(binary.BigEndian.Uint64(iter.Key()))
        if h >= height { break }
        store.Delete(iter.Key())
    }
}
```

- [ ] **Step 4: Run tests** — PASS.

- [ ] **Step 5: Commit**

```bash
git add x/poa/keeper/snapshot.go x/poa/keeper/snapshot_test.go
git commit -m "poa: multiplier snapshot storage with pruning"
```

---

## Task 6: Wrapper keeper (the rescaling adapter)

**Files:**
- Create: `x/poa/keeper/wrapper.go`, `x/poa/keeper/wrapper_test.go`, `x/poa/keeper/validator_adapter.go`

**Architecture decision (locks in spec §3.2):** The wrapper does BOTH:
- Read methods that return `ValidatorI` / `Validator` wrap them in an adapter that reports boosted `Tokens` and `ConsensusPower`. This is required because consumers like gravity-bridge and slashing call `ConsensusPower(powerReduction)` on the returned validator, which is computed from `Tokens`, not from `LastValidatorPower`.
- `Slash` / `SlashWithInfractionReason` divide power by the snapshot multiplier before delegating.
- `LastValidatorPower` storage is also overwritten by EndBlocker (Task 7) so direct reads return boosted values too. Belt-and-suspenders.

- [ ] **Step 1: Define `WrappedStakingKeeper` and constructor**

```go
type WrappedStakingKeeper struct {
    types.StakingKeeper
    poa Keeper
}

func (k Keeper) WrappedStakingKeeper() WrappedStakingKeeper {
    return WrappedStakingKeeper{StakingKeeper: k.sk, poa: k}
}

// boostedTokens returns v.Tokens scaled by the current-block multiplier
// for op if op is in the authority set, else v.Tokens unchanged.
func (w WrappedStakingKeeper) boostedTokens(ctx sdk.Context, op sdk.ValAddress, raw math.Int) math.Int {
    if !w.poa.IsAuthority(ctx, op) {
        return raw
    }
    m := w.poa.MultiplierForValidator(ctx, op, ctx.BlockHeight())
    if m.LTE(sdk.OneDec()) {
        return raw
    }
    return sdk.NewDecFromInt(raw).Mul(m).TruncateInt()
}
```

- [ ] **Step 2: Write failing tests**

```go
func TestWrappedKeeper_Validator_BoostsAuthorityTokens(t *testing.T) {
    k, ctx, mock := keeper.NewTestKeeperWithMockStaking(t)
    auth := sdk.ValAddress("auth-validator-xxxx")
    com  := sdk.ValAddress("com-validator-xxxxx")
    k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
    mock.SetValidatorTokens(auth, sdk.NewInt(100))
    mock.SetValidatorTokens(com, sdk.NewInt(300))
    // snapshot at current height records 5x boost for auth
    k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
        Height: ctx.BlockHeight(),
        Entries: []types.MultiplierEntry{{OperatorAddress: auth.String(), Multiplier: "5.0"}},
    })

    w := k.WrappedStakingKeeper()
    vAuth := w.Validator(ctx, auth)
    require.Equal(t, sdk.NewInt(500), vAuth.GetTokens())

    vCom := w.Validator(ctx, com)
    require.Equal(t, sdk.NewInt(300), vCom.GetTokens())
}

func TestWrappedKeeper_IterateBondedValidatorsByPower_AppliesBoost(t *testing.T) {
    k, ctx, mock := keeper.NewTestKeeperWithMockStaking(t)
    auth := sdk.ValAddress("auth-validator-xxxx")
    com  := sdk.ValAddress("com-validator-xxxxx")
    k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
    mock.SetValidatorTokens(auth, sdk.NewInt(100))
    mock.SetValidatorTokens(com, sdk.NewInt(300))
    mock.SetBondedOrder([]sdk.ValAddress{com, auth})
    k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
        Height: ctx.BlockHeight(),
        Entries: []types.MultiplierEntry{{OperatorAddress: auth.String(), Multiplier: "5.0"}},
    })
    w := k.WrappedStakingKeeper()
    var seen []math.Int
    w.IterateBondedValidatorsByPower(ctx, func(_ int64, v stakingtypes.ValidatorI) bool {
        seen = append(seen, v.GetTokens())
        return false
    })
    require.Equal(t, []math.Int{sdk.NewInt(300), sdk.NewInt(500)}, seen)
}

func TestWrappedKeeper_Slash_AuthorityNormalises(t *testing.T) {
    k, ctx, mock := keeper.NewTestKeeperWithMockStaking(t)
    auth := sdk.ValAddress("auth-validator-xxxx")
    cons := sdk.ConsAddress("auth-cons-aaaaaaaaaa")
    k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
    mock.MapConsToOperator(cons, auth)
    // Snapshot says auth had multiplier 5 at height 50
    k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{
        Height:  50,
        Entries: []types.MultiplierEntry{{OperatorAddress: auth.String(), Multiplier: "5.0"}},
    })

    w := k.WrappedStakingKeeper()
    // Caller passes boosted power 500 (5x100). Wrapper must call sk.Slash with 100.
    w.Slash(ctx, cons, 50, 500, sdk.MustNewDecFromStr("0.05"))
    require.EqualValues(t, 100, mock.LastSlashPower)
}

func TestWrappedKeeper_Slash_CommunityPassThrough(t *testing.T) {
    k, ctx, mock := keeper.NewTestKeeperWithMockStaking(t)
    com := sdk.ValAddress("com-validator-xxxxx")
    cons := sdk.ConsAddress("com-cons-aaaaaaaaaaa")
    mock.MapConsToOperator(cons, com)
    w := k.WrappedStakingKeeper()
    w.Slash(ctx, cons, 100, 300, sdk.MustNewDecFromStr("0.05"))
    require.EqualValues(t, 300, mock.LastSlashPower)
}
```

- [ ] **Step 3: Run tests — expect undefined**

- [ ] **Step 4: Implement `validator_adapter.go`**

The adapter wraps a `stakingtypes.ValidatorI` and overrides token/power methods.

```go
package keeper

import (
    "cosmossdk.io/math"
    sdk "github.com/cosmos/cosmos-sdk/types"
    stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type boostedValidator struct {
    stakingtypes.ValidatorI
    boostedTokens math.Int
}

func (b boostedValidator) GetTokens() math.Int { return b.boostedTokens }
func (b boostedValidator) GetBondedTokens() math.Int {
    if b.IsBonded() { return b.boostedTokens }
    return math.ZeroInt()
}
func (b boostedValidator) ConsensusPower(r math.Int) int64 {
    if !b.IsBonded() { return 0 }
    return sdk.TokensToConsensusPower(b.boostedTokens, r)
}
func (b boostedValidator) PotentialConsensusPower(r math.Int) int64 {
    return sdk.TokensToConsensusPower(b.boostedTokens, r)
}

// adaptValidator wraps v in a boostedValidator if op is an authority and the
// current-block multiplier > 1; otherwise returns v unchanged.
func (w WrappedStakingKeeper) adaptValidator(ctx sdk.Context, v stakingtypes.ValidatorI) stakingtypes.ValidatorI {
    if v == nil { return nil }
    op := v.GetOperator()
    boosted := w.boostedTokens(ctx, op, v.GetTokens())
    if boosted.Equal(v.GetTokens()) { return v }
    return boostedValidator{ValidatorI: v, boostedTokens: boosted}
}
```

- [ ] **Step 5: Implement wrapper methods in `wrapper.go`**

```go
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

// GetBondedValidatorsByPower must return concrete []stakingtypes.Validator (gravity-bridge signature),
// so we mutate the Tokens field of returned copies.
func (w WrappedStakingKeeper) GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator {
    raw := w.StakingKeeper.GetBondedValidatorsByPower(ctx)
    out := make([]stakingtypes.Validator, len(raw))
    for i, v := range raw {
        op, _ := sdk.ValAddressFromBech32(v.OperatorAddress)
        v.Tokens = w.boostedTokens(ctx, op, v.Tokens)
        out[i] = v
    }
    sort.SliceStable(out, func(i, j int) bool { return out[i].Tokens.GT(out[j].Tokens) })
    return out
}

func (w WrappedStakingKeeper) GetValidator(ctx sdk.Context, op sdk.ValAddress) (stakingtypes.Validator, bool) {
    v, found := w.StakingKeeper.GetValidator(ctx, op)
    if !found { return v, false }
    v.Tokens = w.boostedTokens(ctx, op, v.Tokens)
    return v, true
}

func (w WrappedStakingKeeper) IterateLastValidatorPowers(ctx sdk.Context, cb func(sdk.ValAddress, int64) bool) {
    // LastValidatorPower is overwritten by EndBlocker, so the underlying call already returns boosted ints.
    w.StakingKeeper.IterateLastValidatorPowers(ctx, cb)
}

// Slash normalises power for authority validators using the snapshot at infractionHeight.
func (w WrappedStakingKeeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) math.Int {
    val := w.StakingKeeper.ValidatorByConsAddr(ctx, consAddr) // raw, not adapted
    if val == nil {
        return w.StakingKeeper.Slash(ctx, consAddr, infractionHeight, power, slashFactor)
    }
    op := val.GetOperator()
    if !w.poa.IsAuthority(ctx, op) {
        return w.StakingKeeper.Slash(ctx, consAddr, infractionHeight, power, slashFactor)
    }
    m, snapFound := w.poa.MultiplierForValidatorWithStatus(ctx, op, infractionHeight)
    if !snapFound {
        // Authority validator but no snapshot at infractionHeight (snapshot pruned or never written).
        // Refuse silent over-slashing: log and skip the power scaling, treating power as already raw.
        ctx.Logger().Error("poa: missing multiplier snapshot for authority slash; using raw power",
            "operator", op.String(), "infraction_height", infractionHeight)
        return w.StakingKeeper.Slash(ctx, consAddr, infractionHeight, power, slashFactor)
    }
    if m.LTE(sdk.OneDec()) {
        return w.StakingKeeper.Slash(ctx, consAddr, infractionHeight, power, slashFactor)
    }
    raw := sdk.NewDec(power).Quo(m).TruncateInt64()
    return w.StakingKeeper.Slash(ctx, consAddr, infractionHeight, raw, slashFactor)
}

// SlashWithInfractionReason mirrors the same normalisation. (Implementation analogous.)
```

`MultiplierForValidatorWithStatus` is a new helper next to `MultiplierForValidator` in `keeper/snapshot.go`:
```go
func (k Keeper) MultiplierForValidatorWithStatus(ctx sdk.Context, val sdk.ValAddress, height int64) (sdk.Dec, bool) {
    s, ok := k.GetMultiplierSnapshot(ctx, height)
    if !ok { return sdk.OneDec(), false }
    addr := val.String()
    for _, e := range s.Entries {
        if e.OperatorAddress == addr {
            d, err := sdk.NewDecFromStr(e.Multiplier)
            if err != nil { return sdk.OneDec(), true }
            return d, true
        }
    }
    // snapshot exists, validator not boosted on that block
    return sdk.OneDec(), true
}
```

All other interface methods (`Jail`, `Hooks`, `GetParams`, etc.) are inherited from the embedded `types.StakingKeeper`.

- [ ] **Step 5: Run tests** — PASS.

- [ ] **Step 6: Commit**

```bash
git add x/poa/keeper/wrapper.go x/poa/keeper/wrapper_test.go
git commit -m "poa: wrapped staking keeper with Slash power normalisation"
```

---

## Task 7: EndBlocker — overwrite LastValidatorPower & emit ABCI updates

**Files:**
- Create: `x/poa/keeper/abci.go`, `x/poa/keeper/abci_test.go`, `x/poa/abci.go`

- [ ] **Step 1: Write failing test**

```go
func TestEndBlocker_BoostsAuthorityToFloor(t *testing.T) {
    k, ctx, mock := keeper.NewTestKeeperWithMockStaking(t)
    auth1 := sdk.ValAddress("auth1-aaaaaaaaaa")
    auth2 := sdk.ValAddress("auth2-aaaaaaaaaa")
    com1  := sdk.ValAddress("com1-aaaaaaaaaaaa")
    com2  := sdk.ValAddress("com2-aaaaaaaaaaaa")
    k.SetAuthoritySet(ctx, []sdk.ValAddress{auth1, auth2})
    k.SetParams(ctx, types.DefaultParams())

    // raw powers — authority 200, community 400 → 33% authority, needs boost
    mock.SetLastPower(auth1, 100)
    mock.SetLastPower(auth2, 100)
    mock.SetLastPower(com1, 200)
    mock.SetLastPower(com2, 200)
    mock.SetBondedUnjailed([]sdk.ValAddress{auth1, auth2, com1, com2})

    updates := keeper.EndBlocker(ctx, k)

    // Each authority should have been written with boosted power
    p1 := mock.GetLastPower(auth1)
    p2 := mock.GetLastPower(auth2)
    pc1 := mock.GetLastPower(com1)
    pc2 := mock.GetLastPower(com2)
    total := p1 + p2 + pc1 + pc2
    authShare := sdk.NewDec(p1 + p2).Quo(sdk.NewDec(total))
    require.True(t, authShare.GTE(sdk.MustNewDecFromStr("0.67")), "got share %s", authShare)

    // Updates should contain only the changed (boosted) authority entries
    require.Len(t, updates, 2)
}

func TestEndBlocker_AlreadyAboveFloor_NoOp(t *testing.T) {
    k, ctx, mock := keeper.NewTestKeeperWithMockStaking(t)
    auth := sdk.ValAddress("auth-validator-xxxx")
    k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
    k.SetParams(ctx, types.DefaultParams())
    mock.SetLastPower(auth, 1000)
    com := sdk.ValAddress("com-validator-xxxxx")
    mock.SetLastPower(com, 100)
    mock.SetBondedUnjailed([]sdk.ValAddress{auth, com})

    updates := keeper.EndBlocker(ctx, k)
    require.Empty(t, updates)
    require.Equal(t, int64(1000), mock.GetLastPower(auth))
}

func TestEndBlocker_HaltOnEmptyAuthority(t *testing.T) {
    k, ctx, mock := keeper.NewTestKeeperWithMockStaking(t)
    com := sdk.ValAddress("com-validator-xxxxx")
    k.SetAuthoritySet(ctx, []sdk.ValAddress{sdk.ValAddress("offline-aaaaaaa")}) // not bonded
    k.SetParams(ctx, types.DefaultParams())
    mock.SetLastPower(com, 100)
    mock.SetBondedUnjailed([]sdk.ValAddress{com})
    require.Panics(t, func() { keeper.EndBlocker(ctx, k) })
}

func TestEndBlocker_WritesSnapshot(t *testing.T) {
    k, ctx, mock := keeper.NewTestKeeperWithMockStaking(t)
    auth := sdk.ValAddress("auth-validator-xxxx")
    com  := sdk.ValAddress("com-validator-xxxxx")
    k.SetAuthoritySet(ctx, []sdk.ValAddress{auth})
    k.SetParams(ctx, types.DefaultParams())
    mock.SetLastPower(auth, 100)
    mock.SetLastPower(com, 300)
    mock.SetBondedUnjailed([]sdk.ValAddress{auth, com})
    ctx = ctx.WithBlockHeight(42)
    keeper.EndBlocker(ctx, k)
    snap, ok := k.GetMultiplierSnapshot(ctx, 42)
    require.True(t, ok)
    require.Len(t, snap.Entries, 1)
    require.Equal(t, auth.String(), snap.Entries[0].OperatorAddress)
}
```

- [ ] **Step 2: Run tests — expect failure**

- [ ] **Step 3: Implement `keeper/abci.go`**

```go
package keeper

import (
    abci "github.com/cometbft/cometbft/abci/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "cosmossdk.io/math"

    "github.com/peggyjv/sommelier/v9/x/poa/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) []abci.ValidatorUpdate {
    params := k.GetParams(ctx)
    if !params.Enabled {
        return nil
    }

    // 1. Gather raw powers from staking's just-written LastValidatorPower.
    auth := authoritySetMap(k.GetAuthoritySet(ctx))
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
        if _, isAuth := auth[op.String()]; isAuth {
            authPower = authPower.Add(math.NewInt(power))
            authVals = append(authVals, op)
        } else {
            comPower = comPower.Add(math.NewInt(power))
        }
        return false
    })

    if authPower.IsZero() {
        if params.HaltWhenAuthorityEmpty {
            panic(types.ErrNoBondedAuthority)
        }
        return nil
    }

    m := types.ComputeMultiplier(authPower, comPower, params.FloorFraction)
    if m.LTE(sdk.OneDec()) {
        // No boost needed; still record an empty snapshot to enable consistent slashing math.
        k.SetMultiplierSnapshot(ctx, types.MultiplierSnapshot{Height: ctx.BlockHeight()})
        k.pruneSnapshots(ctx, params)
        return nil
    }

    var updates []abci.ValidatorUpdate
    snap := types.MultiplierSnapshot{Height: ctx.BlockHeight()}
    for _, op := range authVals {
        rawPower := k.sk.GetLastValidatorPower(ctx, op)
        boosted := sdk.NewDec(rawPower).Mul(m).TruncateInt64()
        if boosted == rawPower {
            continue
        }
        v, _ := k.sk.GetValidator(ctx, op)
        pk, err := v.TmConsPublicKey()
        if err != nil { panic(err) }

        k.sk.SetLastValidatorPower(ctx, op, boosted)
        updates = append(updates, abci.ValidatorUpdate{PubKey: pk, Power: boosted})
        snap.Entries = append(snap.Entries, types.MultiplierEntry{
            OperatorAddress: op.String(),
            Multiplier:      m.String(),
        })
    }

    k.SetMultiplierSnapshot(ctx, snap)
    k.pruneSnapshots(ctx, params)

    ctx.EventManager().EmitEvent(sdk.NewEvent(
        types.EventTypeAuthorityRescale,
        sdk.NewAttribute(types.AttributeMultiplier, m.String()),
        sdk.NewAttribute(types.AttributeAuthorityPower, authPower.String()),
        sdk.NewAttribute(types.AttributeCommunityPower, comPower.String()),
    ))

    return updates
}

func authoritySetMap(vals []sdk.ValAddress) map[string]struct{} {
    m := make(map[string]struct{}, len(vals))
    for _, v := range vals { m[v.String()] = struct{}{} }
    return m
}

func (k Keeper) pruneSnapshots(ctx sdk.Context, _ types.Params) {
    // retention_blocks = ceil(unbonding_time / avg_block_time) + signed_blocks_window
    unbonding := k.sk.UnbondingTime(ctx)
    const avgBlockNanos = 6 * 1_000_000_000 // 6s; conservative for sommelier
    unbondingBlocks := int64(unbonding.Nanoseconds()/avgBlockNanos) + 1
    signedWindow := k.slashingKeeper.SignedBlocksWindow(ctx) // wire SlashingKeeper into Keeper struct
    retention := unbondingBlocks + signedWindow
    if ctx.BlockHeight() > retention {
        k.PruneSnapshotsBefore(ctx, ctx.BlockHeight()-retention)
    }
}
```

(Add a `SlashingKeeper` interface dependency to `types/expected_keepers.go` exposing `SignedBlocksWindow(ctx) int64`, and inject it in `NewKeeper`. Update Task 2 Step 6 and Task 9 keeper construction accordingly.)

```go
```

- [ ] **Step 4: Add `types/errors.go` entry**

```go
ErrNoBondedAuthority = sdkerrors.Register(ModuleName, 1, "no bonded, unjailed authority validator: refusing to advance block")
```

- [ ] **Step 5: Add `types/events.go`**

```go
const (
    EventTypeAuthorityRescale = "authority_rescale"
    AttributeMultiplier       = "multiplier"
    AttributeAuthorityPower   = "authority_power"
    AttributeCommunityPower   = "community_power"
)
```

- [ ] **Step 6: Wire `module.go` EndBlocker**

```go
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
    return keeper.EndBlocker(ctx, am.keeper)
}
```

- [ ] **Step 7: Run tests** — PASS.

- [ ] **Step 8: Commit**

```bash
git add x/poa
git commit -m "poa: EndBlocker rescales authority validators to >=floor share"
```

---

## Task 8: Gov messages, query server, msg server

**Files:**
- Create: `x/poa/keeper/msg_server.go`, `x/poa/keeper/grpc_query.go`, `x/poa/types/msgs.go`, `x/poa/types/msgs_test.go`, `x/poa/keeper/genesis.go`

- [ ] **Step 1: `types/msgs.go` — ValidateBasic for both messages**

```go
func (m MsgUpdateAuthoritySet) GetSigners() []sdk.AccAddress { /* parse m.Authority */ }
func (m MsgUpdateAuthoritySet) ValidateBasic() error {
    if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil { return err }
    if len(m.Validators) == 0 { return types.ErrEmptyAuthoritySet }
    seen := make(map[string]bool)
    for _, v := range m.Validators {
        if _, err := sdk.ValAddressFromBech32(v); err != nil { return err }
        if seen[v] { return types.ErrDuplicateAuthority }
        seen[v] = true
    }
    return nil
}
// Similar for MsgUpdateParams
```

- [ ] **Step 2: Tests for `ValidateBasic`** covering empty list, duplicate, malformed bech32.

- [ ] **Step 3: `keeper/msg_server.go`**

```go
type msgServer struct{ Keeper }
func NewMsgServerImpl(k Keeper) types.MsgServer { return msgServer{k} }

func (s msgServer) UpdateAuthoritySet(goCtx context.Context, msg *types.MsgUpdateAuthoritySet) (*types.MsgUpdateAuthoritySetResponse, error) {
    if msg.Authority != s.authority { return nil, sdkerrors.Wrap(govtypes.ErrInvalidSigner, msg.Authority) }
    ctx := sdk.UnwrapSDKContext(goCtx)
    addrs := make([]sdk.ValAddress, len(msg.Validators))
    for i, v := range msg.Validators { addrs[i], _ = sdk.ValAddressFromBech32(v) }
    s.SetAuthoritySet(ctx, addrs)
    return &types.MsgUpdateAuthoritySetResponse{}, nil
}

func (s msgServer) UpdateParams(...) {...}
```

- [ ] **Step 4: `keeper/grpc_query.go`** — implement Params, AuthoritySet, EffectivePower.

- [ ] **Step 5: `keeper/genesis.go`**

```go
func (k Keeper) InitGenesis(ctx sdk.Context, gs types.GenesisState) {
    k.SetParams(ctx, gs.Params)
    addrs := make([]sdk.ValAddress, len(gs.AuthoritySet))
    for i, v := range gs.AuthoritySet { addrs[i], _ = sdk.ValAddressFromBech32(v.OperatorAddress) }
    k.SetAuthoritySet(ctx, addrs)
}
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState { /* mirror */ }
```

- [ ] **Step 6: Wire msg/query servers in `module.go`** (`RegisterServices`).

- [ ] **Step 7: Run tests**

Run: `go test ./x/poa/...`
Expected: PASS.

- [ ] **Step 8: Commit**

```bash
git add x/poa
git commit -m "poa: gov msg server, gRPC queries, genesis"
```

---

## Task 9: Wire into `app/app.go`

**Files:**
- Modify: `app/app.go`

This is the surgical task. Each consumer must be passed the wrapped keeper. **Do this one consumer at a time and run `go build ./...` between each** to keep error surface manageable.

- [ ] **Step 1: Add module imports and store key**

Add `poatypes.StoreKey` to `keys` map in `NewSommelierApp`. Add `poatypes.ModuleName` to `maccPerms` only if needed (no module account expected).

- [ ] **Step 2: Construct `app.PoaKeeper` AFTER `app.StakingKeeper`, BEFORE any consumer**

```go
app.PoaKeeper = poakeeper.NewKeeper(
    appCodec, keys[poatypes.StoreKey],
    app.GetSubspace(poatypes.ModuleName),
    app.StakingKeeper, // raw, the wrapper is called downstream
    authority,
)
wrappedSk := app.PoaKeeper.WrappedStakingKeeper()
```

- [ ] **Step 3: Replace each downstream `app.StakingKeeper` with `wrappedSk`**

Edit, build, repeat for **consumers that must see boosted power**:
- `slashingkeeper.NewKeeper(..., wrappedSk, ...)`
- `distrkeeper.NewKeeper(..., wrappedSk, ...)` — note: distribution rewards are still calculated from real bonded ratio because distribution allocates rewards by `ConsensusPower` from validators it iterates, but reward *math* (inflation, fee splits) flows from mint/bank which use raw values; in practice rewarding by boosted power is acceptable since boosted relative weights == raw relative weights for community↔community comparisons and authority validators are already incentivised separately. Document this in `docs/poa.md` (Task 12).
- `evidencekeeper.NewKeeper(..., &wrappedSk, ...)`
- `gravitykeeper.NewKeeper(..., wrappedSk, ...)`
- `corkkeeper.NewKeeper(..., wrappedSk, ...)`
- `axelarcorkkeeper.NewKeeper(..., wrappedSk, ...)`
- `pubsubkeeper.NewKeeper(..., wrappedSk, ...)`
- `incentiveskeeper.NewKeeper(..., wrappedSk, ...)`

**Consumers that intentionally keep the RAW `app.StakingKeeper` (no boost):**
- `mintkeeper` — inflation/`BondedRatio`/`StakingTokenSupply` must reflect actual bonded tokens, not boosted.
- `crisiskeeper` — module invariants check real economic state.
- `genutilkeeper` — operates on genesis transactions and real validator metadata.
- `ibckeeper` light client — IBC light clients verify against the actual CometBFT validator set, which already reflects boosted powers via the ABCI `ValidatorUpdates` we emit. Pass raw to avoid double-counting.
- `authkeeper`, `bankkeeper` — no power semantics.
- `x/staking` itself — owns the underlying keeper; never wrapped.

Run `go build ./...` after each substitution. If interface mismatch, add the missing method to `WrappedStakingKeeper` (delegate to embedded).

- [ ] **Step 4: Add `x/poa` to module manager AND module basics**

Append `poa.AppModuleBasic{}` to the `module.NewBasicManager(...)` call.

```go
app.mm = module.NewManager(
    ...
    poa.NewAppModule(appCodec, app.PoaKeeper),
)
```

Also add `poa.NewAppModule(appCodec, app.PoaKeeper)` to the `simulationManager` constructor (`app.sm = module.NewSimulationManager(...)`).

- [ ] **Step 5: EndBlocker order**

Insert `poatypes.ModuleName` immediately after `stakingtypes.ModuleName` in `mm.SetOrderEndBlockers(...)`.

- [ ] **Step 6: BeginBlocker order**

Append `poatypes.ModuleName` to `SetOrderBeginBlockers(...)` (no-op begin block, but registration is required).

- [ ] **Step 7: InitGenesis & ExportGenesis order**

Insert `poatypes.ModuleName` after `stakingtypes.ModuleName` and `slashingtypes.ModuleName`, before `gravitytypes.ModuleName` in `SetOrderInitGenesis(...)`. Append the same module name to `SetOrderExportGenesis(...)` if the project uses an explicit export ordering — Cosmos SDK 0.47 panics on missing modules in export.

- [ ] **Step 8: Param subspace**

Add `paramsKeeper.Subspace(poatypes.ModuleName)` in `initParamsKeeper`.

- [ ] **Step 9: Build and run existing tests**

Run: `go build ./... && go test ./app/...`
Expected: PASS.

- [ ] **Step 10: Lint sentinel**

Add a comment near `app.StakingKeeper` definition warning future contributors:
```go
// WARNING: For all power-sensitive consumers, pass app.PoaKeeper.WrappedStakingKeeper(),
// not app.StakingKeeper directly. See docs/poa.md.
```

- [ ] **Step 11: Commit**

```bash
git add app/app.go
git commit -m "app: wire x/poa and route validator-power consumers through wrapped keeper"
```

---

## Task 10: Integration tests

**Files:**
- Create: `x/poa/keeper/integration_test.go`

- [ ] **Step 1: simapp test — bond mixed validators, run blocks, assert floor**

```go
func TestPoa_FloorMaintainedAcrossBlocks(t *testing.T) {
    app, ctx := setupSimApp(t) // helper that creates 3 authority + 5 community vals with mixed stake
    for i := 0; i < 10; i++ {
        ctx = ctx.WithBlockHeight(int64(i+1))
        app.EndBlocker(ctx, abci.RequestEndBlock{Height: ctx.BlockHeight()})
        share := authorityShare(t, app, ctx)
        require.True(t, share.GTE(sdk.MustNewDecFromStr("0.67")),
            "block %d: authority share %s below floor", i, share)
    }
}
```

- [ ] **Step 2: gravity-bridge SignerSetTx weight check**

Inspect the SignerSetTx the gravity module would produce next; assert authority signers' summed weight `>= 67%`. (Use `app.GravityKeeper.CreateSignerSetTx(ctx)` semantics — exact API depends on gravity v6.)

- [ ] **Step 3: Slashing integration**

Trigger a downtime slash on an authority validator; assert burned tokens == `slashFraction * actual_bonded`, NOT `slashFraction * boosted_bonded`.

- [ ] **Step 4: Halt-on-empty integration**

Jail every authority validator; assert next `EndBlocker` panics.

- [ ] **Step 5: Run** `go test ./x/poa/keeper/ -run TestPoa_ -v`

- [ ] **Step 6: Commit**

```bash
git add x/poa/keeper/integration_test.go
git commit -m "poa: integration tests for floor invariant, gravity, slashing, halt"
```

---

## Task 11: v10 upgrade

**Files:**
- Create: `app/upgrades/v10/constants.go`, `app/upgrades/v10/upgrades.go`
- Modify: `app/upgrades.go` (or wherever `setupUpgradeHandlers` lives)

- [ ] **Step 1: `constants.go`**

```go
package v10

import (
    storetypes "github.com/cosmos/cosmos-sdk/store/types"
    poatypes "github.com/peggyjv/sommelier/v9/x/poa/types"
)

const Name = "v10"

var Upgrade = upgrades.Upgrade{
    Name: Name,
    StoreUpgrades: storetypes.StoreUpgrades{
        Added: []string{poatypes.StoreKey},
    },
}

// DefaultAuthorityValidators is the compiled-in initial authority set.
// Operator: replace with actual operator addresses before tagging the v10 release.
// The upgrade handler refuses to run if this slice is empty (see upgrades.go).
var DefaultAuthorityValidators = []string{
    // "sommvaloper1...",
}
```

- [ ] **Step 2: `upgrades.go`**

```go
func CreateUpgradeHandler(
    mm *module.Manager,
    cfg module.Configurator,
    keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
    return func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
        if len(DefaultAuthorityValidators) == 0 {
            return fromVM, fmt.Errorf("v10 upgrade refuses to run with empty DefaultAuthorityValidators: chain would halt next block (halt_when_authority_empty=true)")
        }

        // Initialise PoA params and authority set.
        keepers.PoaKeeper.SetParams(ctx, poatypes.DefaultParams())

        addrs := make([]sdk.ValAddress, 0, len(DefaultAuthorityValidators))
        for _, s := range DefaultAuthorityValidators {
            v, err := sdk.ValAddressFromBech32(s)
            if err != nil { return fromVM, fmt.Errorf("invalid authority %s: %w", s, err) }
            addrs = append(addrs, v)
        }
        keepers.PoaKeeper.SetAuthoritySet(ctx, addrs)

        return mm.RunMigrations(ctx, cfg, fromVM)
    }
}
```

- [ ] **Step 3: Register in `app/upgrades.go`**

Mirror the v9 registration block.

- [ ] **Step 4: Test upgrade in simapp**

```go
func TestUpgrade_v10(t *testing.T) {
    app, ctx := setupSimApp(t)
    plan := upgradetypes.Plan{Name: v10.Name, Height: ctx.BlockHeight() + 1}
    require.NoError(t, app.UpgradeKeeper.ScheduleUpgrade(ctx, plan))
    // advance to upgrade height; assert PoA params populated
    // ...
}
```

- [ ] **Step 5: Commit**

```bash
git add app/upgrades/v10 app/upgrades.go
git commit -m "v10: upgrade handler that registers x/poa and seeds authority set"
```

---

## Task 12: Operator docs

**Files:**
- Create: `docs/poa.md`

- [ ] **Step 1: Write a 1-page operator brief** covering:
  - What changed in v10 (PoA semantics, 67% floor).
  - Why a community validator with 50% stake will only show ~30% consensus power.
  - How slashing works (slash on raw stake, not boosted).
  - How to propose adding/removing an authority validator (`MsgUpdateAuthoritySet`).
  - Warning about chain halt if all authority validators are jailed.

- [ ] **Step 2: Commit**

```bash
git add docs/poa.md
git commit -m "docs: PoA operator brief"
```

---

## Verification gate (run before declaring done)

- [ ] `go build ./...` — clean
- [ ] `go test ./x/poa/...` — all green
- [ ] `go test ./app/...` — all green (including upgrade test)
- [ ] `go vet ./...` — clean
- [ ] Manually grep: `git grep -n 'app\.StakingKeeper' app/` should only show the construction site, the staking module wiring, and the warning comment.
- [ ] Run an existing integration test from `integration_tests/` to confirm gravity-bridge still functions end-to-end.

---

## Deferred (intentionally not in this plan)

- Per-authority weight overrides (spec §8).
- CLI commands for the new gov messages — implement when first used; can be done via `--generate-only` JSON in the meantime.
- Metrics dashboards / Prometheus exporters for the new events.
