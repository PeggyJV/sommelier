package app

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	v10 "github.com/peggyjv/sommelier/v10/app/upgrades/v10"
	axelarcorktypes "github.com/peggyjv/sommelier/v10/x/axelarcork/types"
	corktypes "github.com/peggyjv/sommelier/v10/x/cork/types"
	corkv2types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
	poatypes "github.com/peggyjv/sommelier/v10/x/poa/types"
)

// Exercises the real v10 handler against a real app: keeper wiring, store-key
// resolution through app.GetKey, and params written through the actual param
// subspaces. The unit test in app/upgrades/v10 covers DrainLegacyCorkQueues in
// isolation and cannot catch a mis-wired call site.
func TestV10UpgradeHandlerSeedsAuthorityAndDrainsLegacyQueues(t *testing.T) {
	// isCheckTx=true skips InitChain, which would otherwise panic on an empty
	// validator set. The stores are still mounted and the keepers real, which
	// is what this test is about; params are set explicitly below.
	app := Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{Height: 10})

	app.CorkKeeper.SetParams(ctx, corkv2types.DefaultParams())
	app.AxelarCorkKeeper.SetParams(ctx, axelarcorktypes.DefaultParams())
	seedBondedAuthorityValidator(t, app, ctx)

	corkStore := ctx.KVStore(app.GetKey(corktypes.StoreKey))
	axelarStore := ctx.KVStore(app.GetKey(axelarcorktypes.StoreKey))

	val := sdk.ValAddress([]byte("12345678901234567890"))
	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	id := make([]byte, 32)
	for i := range id {
		id[i] = byte(i)
	}

	// Seed the retired queues the way a live chain carries them into the
	// upgrade. Mainnet currently holds 38 such entries on arbitrum.
	seeded := 0
	for _, h := range []uint64{100, 101} {
		corkStore.Set(corktypes.GetScheduledCorkKey(h, id, val, contract), []byte{0x01})
		seeded++
	}
	corkStore.Set(corktypes.GetValidatorCorkCountKey(val), []byte{0x02})
	seeded++

	for _, chainID := range []uint64{42161, 10} {
		axelarStore.Set(axelarcorktypes.GetScheduledAxelarCorkKey(chainID, 100, id, val, contract), []byte{0x01})
		seeded++
	}
	axelarStore.Set(axelarcorktypes.GetValidatorAxelarCorkCountKey(val), []byte{0x02})
	seeded++

	require.Equal(t, 6, seeded)

	handler := v10.CreateUpgradeHandler(
		app.mm,
		app.configurator,
		app.PoaKeeper,
		app.CorkKeeper,
		app.AxelarCorkKeeper,
		app.GetKey(corktypes.StoreKey),
		app.GetKey(axelarcorktypes.StoreKey),
	)

	_, err := handler(ctx, upgradetypes.Plan{Name: v10.UpgradeName}, app.mm.GetVersionMap())
	require.NoError(t, err)

	// Authority seeded in both modules, through the real param subspaces.
	require.Equal(t, v10.CorkAuthorityAddress, app.CorkKeeper.GetParamSet(ctx).CorkAuthority)
	require.Equal(t, v10.CorkAuthorityAddress, app.AxelarCorkKeeper.GetParamSet(ctx).CorkAuthority)

	// Both retired queues and both counter prefixes emptied.
	for _, c := range []struct {
		name   string
		store  sdk.KVStore
		prefix []byte
	}{
		{"cork scheduled", corkStore, []byte{corktypes.ScheduledCorkKeyPrefix}},
		{"cork counters", corkStore, []byte{corktypes.ValidatorCorkCountKey}},
		{"axelarcork scheduled", axelarStore, []byte{axelarcorktypes.ScheduledCorkKeyPrefix}},
		{"axelarcork counters", axelarStore, []byte{axelarcorktypes.ValidatorAxelarCorkCountKey}},
	} {
		n := 0
		it := sdk.KVStorePrefixIterator(c.store, c.prefix)
		for ; it.Valid(); it.Next() {
			n++
		}
		it.Close()
		require.Zerof(t, n, "%s must be drained; leftovers are permanently unreachable state", c.name)
	}

	// The PoA authority set is seeded and non-empty, or the chain would enter
	// safe mode on the next block.
	require.NotEmpty(t, app.PoaKeeper.GetAuthoritySet(ctx))
}

// Regression: the real upgrade's fromVM comes from on-chain v9 state, which has
// no x/poa entry because the module does not exist there. RunMigrations then
// runs poa's InitGenesis with DefaultGenesis, whose AuthoritySet is nil, and
// InitGenesis calls SetAuthoritySet unconditionally.
//
// The test above passes app.mm.GetVersionMap(), which DOES contain poa, so
// RunMigrations skips InitGenesis and the wipe never happens. That is a false
// pass: it does not model the upgrade being shipped.
func TestV10UpgradeHandlerKeepsAuthoritySetThroughRunMigrations(t *testing.T) {
	app := Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{Height: 10})

	app.CorkKeeper.SetParams(ctx, corkv2types.DefaultParams())
	app.AxelarCorkKeeper.SetParams(ctx, axelarcorktypes.DefaultParams())
	seedBondedAuthorityValidator(t, app, ctx)

	handler := v10.CreateUpgradeHandler(
		app.mm,
		app.configurator,
		app.PoaKeeper,
		app.CorkKeeper,
		app.AxelarCorkKeeper,
		app.GetKey(corktypes.StoreKey),
		app.GetKey(axelarcorktypes.StoreKey),
	)

	// Model a v9 chain: every module at its current version EXCEPT poa, which
	// does not exist on v9.
	fromVM := app.mm.GetVersionMap()
	delete(fromVM, poatypes.ModuleName)

	_, err := handler(ctx, upgradetypes.Plan{Name: v10.UpgradeName}, fromVM)
	require.NoError(t, err)

	require.NotEmpty(t, app.PoaKeeper.GetAuthoritySet(ctx),
		"authority set was wiped by RunMigrations running poa InitGenesis with "+
			"DefaultGenesis; the chain would enter authority-empty safe mode on the "+
			"next block with no on-chain recovery")
	require.Len(t, app.PoaKeeper.GetAuthoritySet(ctx), len(v10.DefaultAuthorityValidators))
}

// seedBondedAuthorityValidator makes the first configured authority validator
// exist as bonded and unjailed, which is what the chain looks like at the real
// upgrade height. The handler refuses to seed an allowlist with no usable
// validator, since that would enter authority-empty safe mode with no on-chain
// recovery.
func seedBondedAuthorityValidator(t *testing.T, app *SommelierApp, ctx sdk.Context) {
	t.Helper()

	opAddr, err := sdk.ValAddressFromBech32(v10.DefaultAuthorityValidators[0])
	require.NoError(t, err)

	val, err := stakingtypes.NewValidator(
		opAddr,
		ed25519.GenPrivKey().PubKey(),
		stakingtypes.Description{Moniker: "authority-0"},
	)
	require.NoError(t, err)
	val.Status = stakingtypes.Bonded
	val.Jailed = false
	val.Tokens = sdk.NewInt(1_000_000)
	val.DelegatorShares = sdk.NewDecFromInt(val.Tokens)

	app.StakingKeeper.SetValidator(ctx, val)
}

// A valid-looking but non-live authority address must abort the upgrade, not
// seed an allowlist whose members carry no power.
//
// The empty-slice guard catches a forgotten constant; it cannot catch a typo, a
// stale operator address, or every listed validator being jailed at the upgrade
// height. Any of those puts the chain into authority-empty safe mode on the
// next block, where MsgUpdateAuthoritySet and MsgUpdateParams are both frozen
// and there is no on-chain recovery. Halting at the upgrade height is
// recoverable; that state is not.
func TestV10UpgradeHandlerRefusesWhenNoAuthorityValidatorIsLive(t *testing.T) {
	app := Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{Height: 10})

	app.CorkKeeper.SetParams(ctx, corkv2types.DefaultParams())
	app.AxelarCorkKeeper.SetParams(ctx, axelarcorktypes.DefaultParams())
	// Deliberately seed NO validator.

	handler := v10.CreateUpgradeHandler(
		app.mm,
		app.configurator,
		app.PoaKeeper,
		app.CorkKeeper,
		app.AxelarCorkKeeper,
		app.GetKey(corktypes.StoreKey),
		app.GetKey(axelarcorktypes.StoreKey),
	)

	fromVM := app.mm.GetVersionMap()
	delete(fromVM, poatypes.ModuleName)

	_, err := handler(ctx, upgradetypes.Plan{Name: v10.UpgradeName}, fromVM)
	require.Error(t, err)
	require.Contains(t, err.Error(), "bonded and unjailed")
}

// A jailed authority validator does not count as live.
func TestV10UpgradeHandlerRefusesWhenAuthorityValidatorIsJailed(t *testing.T) {
	app := Setup(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{Height: 10})

	app.CorkKeeper.SetParams(ctx, corkv2types.DefaultParams())
	app.AxelarCorkKeeper.SetParams(ctx, axelarcorktypes.DefaultParams())

	opAddr, err := sdk.ValAddressFromBech32(v10.DefaultAuthorityValidators[0])
	require.NoError(t, err)
	val, err := stakingtypes.NewValidator(opAddr, ed25519.GenPrivKey().PubKey(),
		stakingtypes.Description{Moniker: "jailed-authority"})
	require.NoError(t, err)
	val.Status = stakingtypes.Bonded
	val.Jailed = true
	app.StakingKeeper.SetValidator(ctx, val)

	handler := v10.CreateUpgradeHandler(
		app.mm, app.configurator, app.PoaKeeper, app.CorkKeeper, app.AxelarCorkKeeper,
		app.GetKey(corktypes.StoreKey), app.GetKey(axelarcorktypes.StoreKey),
	)
	fromVM := app.mm.GetVersionMap()
	delete(fromVM, poatypes.ModuleName)

	_, err = handler(ctx, upgradetypes.Plan{Name: v10.UpgradeName}, fromVM)
	require.Error(t, err, "a jailed authority validator must not satisfy the liveness check")
}
