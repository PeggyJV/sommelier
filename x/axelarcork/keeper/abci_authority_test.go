package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	_ "github.com/peggyjv/sommelier/v10/app/params"
	"github.com/peggyjv/sommelier/v10/x/axelarcork/types"
)

// endBlockerFixture wires an enabled module with two configured chains so the
// per-chain isolation of the EndBlocker can be exercised.
func endBlockerFixture(t *testing.T, safeMode bool) (Keeper, sdk.Context) {
	t.Helper()

	k, ctx, m, ctrl := setupCorkKeeper(t)
	t.Cleanup(ctrl.Finish)

	k.SetPoaKeeper(stubPoaKeeper{active: safeMode})

	// The EndBlocker sweeps the module account when not in safe mode. Return a
	// zero balance so the sweep short-circuits: these tests are about cork
	// movement, not fund sweeping, which has its own coverage.
	moduleAcct := authtypes.NewEmptyModuleAccount(types.ModuleName)
	m.mockAccountKeeper.EXPECT().GetModuleAccount(gomock.Any(), types.ModuleName).
		Return(moduleAcct).AnyTimes()
	m.mockBankKeeper.EXPECT().GetAllBalances(gomock.Any(), moduleAcct.GetAddress()).
		Return(sdk.Coins{}).AnyTimes()

	params := types.DefaultParams()
	params.Enabled = true
	params.CorkAuthority = testCorkAuthority
	k.SetParams(ctx, params)

	for id, name := range map[uint64]string{testChainArbitrum: "arbitrum", testChainOptimism: "optimism"} {
		k.SetChainConfiguration(ctx, id, types.ChainConfiguration{
			Name:         name,
			Id:           id,
			ProxyAddress: "0x9999999999999999999999999999999999999999",
		})
	}

	// Allowlist every cellar these tests target. The EndBlocker re-checks the
	// allowlist at execution, so an unmanaged target is dropped rather than
	// becoming relayable.
	managed := []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
		"0x3333333333333333333333333333333333333333",
		"0x4444444444444444444444444444444444444444",
		"0x5555555555555555555555555555555555555555",
	}
	for _, id := range []uint64{testChainArbitrum, testChainOptimism} {
		k.SetCellarIDs(ctx, id, types.CellarIDSet{ChainId: id, Ids: managed})
	}

	return k, ctx
}

func countWinning(k Keeper, ctx sdk.Context, chainID uint64) int {
	n := 0
	k.IterateWinningAxelarCorks(ctx, chainID, func(_ common.Address, _ uint64, _ types.AxelarCork) bool {
		n++
		return false
	})
	return n
}

// A cork due at the current height moves into the relayable (winning) queue and
// leaves the scheduled queue.
func TestEndBlockerMovesDueCorksToWinning(t *testing.T) {
	k, ctx := endBlockerFixture(t, false)

	height := uint64(ctx.BlockHeight())
	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	k.SetAuthorityAxelarCork(ctx, testChainArbitrum, height, types.AxelarCork{
		TargetContractAddress: contract.String(),
		EncodedContractCall:   []byte{0xaa},
		ChainId:               testChainArbitrum,
	})

	k.EndBlocker(ctx)

	require.Zero(t, countAt(k, ctx, testChainArbitrum, height),
		"due cork must leave the scheduled queue or it re-processes every block")
	require.Equal(t, 1, countWinning(k, ctx, testChainArbitrum),
		"due cork must become relayable")
}

// A cork on one chain must not surface on another.
func TestEndBlockerPerChainIsolation(t *testing.T) {
	k, ctx := endBlockerFixture(t, false)

	height := uint64(ctx.BlockHeight())
	contract := common.HexToAddress("0x2222222222222222222222222222222222222222")
	k.SetAuthorityAxelarCork(ctx, testChainArbitrum, height, types.AxelarCork{
		TargetContractAddress: contract.String(),
		EncodedContractCall:   []byte{0xbb},
		ChainId:               testChainArbitrum,
	})

	k.EndBlocker(ctx)

	require.Equal(t, 1, countWinning(k, ctx, testChainArbitrum))
	require.Zero(t, countWinning(k, ctx, testChainOptimism),
		"a cork scheduled for arbitrum must not become relayable on optimism")
}

// A cork targeting a later height is left alone.
func TestEndBlockerIgnoresFutureCorks(t *testing.T) {
	k, ctx := endBlockerFixture(t, false)

	future := uint64(ctx.BlockHeight()) + 1
	contract := common.HexToAddress("0x3333333333333333333333333333333333333333")
	k.SetAuthorityAxelarCork(ctx, testChainArbitrum, future, types.AxelarCork{
		TargetContractAddress: contract.String(),
		EncodedContractCall:   []byte{0xcc},
		ChainId:               testChainArbitrum,
	})

	k.EndBlocker(ctx)

	require.Equal(t, 1, countAt(k, ctx, testChainArbitrum, future),
		"a cork scheduled for a later height must survive untouched")
	require.Zero(t, countWinning(k, ctx, testChainArbitrum))
}

// Safe mode deletes due corks without marking them relayable, and reports it.
func TestEndBlockerSafeModeDropsCorks(t *testing.T) {
	k, ctx := endBlockerFixture(t, true)

	height := uint64(ctx.BlockHeight())
	contract := common.HexToAddress("0x4444444444444444444444444444444444444444")
	k.SetAuthorityAxelarCork(ctx, testChainArbitrum, height, types.AxelarCork{
		TargetContractAddress: contract.String(),
		EncodedContractCall:   []byte{0xdd},
		ChainId:               testChainArbitrum,
	})

	k.EndBlocker(ctx)

	require.Zero(t, countAt(k, ctx, testChainArbitrum, height),
		"cork must be deleted even in safe mode, or it is stranded in state forever")
	require.Zero(t, countWinning(k, ctx, testChainArbitrum),
		"safe mode must not mark corks relayable")

	found := false
	for _, e := range ctx.EventManager().Events() {
		if e.Type == types.EventTypeAxelarCorksDroppedInSafeMode {
			found = true
		}
	}
	require.True(t, found, "expected an axelar_corks_dropped_safe_mode event")
}

// The timed-out-cork sweep is a liveness mechanism, not a strategist one, and
// must survive the rewrite. It matters more now that relaying depends on a
// single key: an unrelayed cork must still expire.
func TestEndBlockerTimeoutSweepStillRuns(t *testing.T) {
	k, ctx := endBlockerFixture(t, false)

	params := k.GetParamSet(ctx)
	params.CorkTimeoutBlocks = 1
	k.SetParams(ctx, params)

	scheduledAt := uint64(ctx.BlockHeight())
	contract := common.HexToAddress("0x5555555555555555555555555555555555555555")
	k.SetWinningAxelarCork(ctx, testChainArbitrum, scheduledAt, types.AxelarCork{
		TargetContractAddress: contract.String(),
		EncodedContractCall:   []byte{0xee},
		ChainId:               testChainArbitrum,
	})
	require.Equal(t, 1, countWinning(k, ctx, testChainArbitrum))

	// Advance past the timeout and run the EndBlocker again.
	later := ctx.WithBlockHeight(int64(scheduledAt) + 5)
	k.EndBlocker(later)

	require.Zero(t, countWinning(k, later, testChainArbitrum),
		"an unrelayed cork past cork_timeout_blocks must be swept")
}

// Removing a cellar must stop corks already queued against it, not just block
// new scheduling. Without this, a compromised authority key can stage calls far
// into the future that no revocation reaches.
func TestEndBlockerDropsCorksForDelistedCellars(t *testing.T) {
	k, ctx := endBlockerFixture(t, false)

	height := uint64(ctx.BlockHeight())
	delisted := common.HexToAddress("0x2222222222222222222222222222222222222222")
	k.SetAuthorityAxelarCork(ctx, testChainArbitrum, height, types.AxelarCork{
		TargetContractAddress: delisted.String(),
		EncodedContractCall:   []byte{0x99},
		ChainId:               testChainArbitrum,
	})

	// Narrow the allowlist so the queued target is no longer managed.
	k.SetCellarIDs(ctx, testChainArbitrum, types.CellarIDSet{
		ChainId: testChainArbitrum,
		Ids:     []string{"0x1111111111111111111111111111111111111111"},
	})

	k.EndBlocker(ctx)

	require.Zero(t, countAt(k, ctx, testChainArbitrum, height),
		"the cork must leave the scheduled queue rather than being stranded")
	require.Zero(t, countWinning(k, ctx, testChainArbitrum),
		"a delisted cellar's cork must not become relayable")
}
