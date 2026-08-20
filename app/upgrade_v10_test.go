package app

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	v10 "github.com/peggyjv/sommelier/v10/app/upgrades/v10"
	axelarcorktypes "github.com/peggyjv/sommelier/v10/x/axelarcork/types"
	corktypes "github.com/peggyjv/sommelier/v10/x/cork/types"
	corkv2types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
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
