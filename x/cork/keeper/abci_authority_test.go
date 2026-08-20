package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"

	corktypes "github.com/peggyjv/sommelier/v10/x/cork/types"
	types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
)

func authorityTestCork(call string, contract common.Address) types.Cork {
	return types.Cork{
		EncodedContractCall:   []byte(call),
		TargetContractAddress: contract.String(),
	}
}

func countAuthorityCorksAt(k Keeper, ctx sdk.Context, height uint64) int {
	n := 0
	k.IterateAuthorityCorksByBlockHeight(ctx, height, func(_ uint64, _ []byte, _ common.Address, _ types.Cork) bool {
		n++
		return false
	})
	return n
}

// Every authority cork due at the current height is submitted to the bridge
// exactly once, then removed from the queue.
func (suite *KeeperTestSuite) TestEndBlockerSubmitsDueAuthorityCorks() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	corkKeeper.SetPoaKeeper(stubPoaKeeper{active: false})

	height := uint64(ctx.BlockHeight())
	contractA := common.HexToAddress("0x1111111111111111111111111111111111111111")
	contractB := common.HexToAddress("0x2222222222222222222222222222222222222222")
	corkKeeper.SetAuthorityCork(ctx, height, authorityTestCork("callA", contractA))
	corkKeeper.SetAuthorityCork(ctx, height, authorityTestCork("callB", contractB))
	require.Equal(2, countAuthorityCorksAt(corkKeeper, ctx, height))

	// Exactly two submissions, no more: Times(2) fails on a third call, and
	// gomock's controller fails the test if fewer than two occur.
	suite.gravityKeeper.EXPECT().
		CreateContractCallTx(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(2)

	corkKeeper.EndBlocker(ctx)

	require.Zero(countAuthorityCorksAt(corkKeeper, ctx, height),
		"due corks must be deleted after submission or they re-execute every block")
}

// A cork targeting a future height is neither submitted nor deleted.
func (suite *KeeperTestSuite) TestEndBlockerIgnoresFutureAuthorityCorks() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	corkKeeper.SetPoaKeeper(stubPoaKeeper{active: false})

	future := uint64(ctx.BlockHeight()) + 1
	contract := common.HexToAddress("0x3333333333333333333333333333333333333333")
	corkKeeper.SetAuthorityCork(ctx, future, authorityTestCork("later", contract))

	// No CreateContractCallTx expectation: any submission fails the test.
	corkKeeper.EndBlocker(ctx)

	require.Equal(1, countAuthorityCorksAt(corkKeeper, ctx, future),
		"a cork scheduled for a later height must survive this block untouched")
}

// Safe mode drops due corks: they are deleted (never stranded) but not
// submitted, and the drop is reported as an event.
func (suite *KeeperTestSuite) TestEndBlockerSafeModeDropsAuthorityCorks() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	corkKeeper.SetPoaKeeper(stubPoaKeeper{active: true})

	height := uint64(ctx.BlockHeight())
	contract := common.HexToAddress("0x4444444444444444444444444444444444444444")
	corkKeeper.SetAuthorityCork(ctx, height, authorityTestCork("frozen", contract))

	// No CreateContractCallTx expectation is registered: if the EndBlocker
	// submitted the cork, gomock would fail this test.
	corkKeeper.EndBlocker(ctx)

	require.Zero(countAuthorityCorksAt(corkKeeper, ctx, height),
		"cork must be deleted even in safe mode, or it is stranded in state forever")

	found := false
	for _, e := range ctx.EventManager().Events() {
		if e.Type == corktypes.EventTypeCorksDroppedInSafeMode {
			found = true
		}
	}
	require.True(found, "expected a corks_dropped_safe_mode event")
}

// An empty queue must be a clean no-op in both modes.
func (suite *KeeperTestSuite) TestEndBlockerNoDueAuthorityCorksIsNoop() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	corkKeeper.SetPoaKeeper(stubPoaKeeper{active: false})

	// No expectations registered: any bridge call or event fails the test.
	require.NotPanics(func() { corkKeeper.EndBlocker(ctx) })

	for _, e := range ctx.EventManager().Events() {
		require.NotEqual(corktypes.EventTypeCorksDroppedInSafeMode, e.Type,
			"must not report drops when there was nothing due")
	}
}

// Removing a cellar from the allowlist must stop corks ALREADY queued against
// it, not just block new ones.
//
// The allowlist is checked when a cork is scheduled. Without a re-check at
// execution, a compromised authority key can queue calls far into the future
// and neither rotating the authority nor removing the target cellar stops them
// reaching the bridge -- so the documented de-escalation path (remove the
// cellar ID) does not actually work against a staged attack.
func (suite *KeeperTestSuite) TestEndBlockerDropsCorksForDelistedCellars() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	corkKeeper.SetPoaKeeper(stubPoaKeeper{active: false})

	height := uint64(ctx.BlockHeight())
	allowed := common.HexToAddress("0x1111111111111111111111111111111111111111")
	delisted := common.HexToAddress("0x2222222222222222222222222222222222222222")

	corkKeeper.SetAuthorityCork(ctx, height, authorityTestCork("stillOk", allowed))
	corkKeeper.SetAuthorityCork(ctx, height, authorityTestCork("revoked", delisted))

	// Only the first cellar remains managed.
	corkKeeper.SetCellarIDs(ctx, types.CellarIDSet{Ids: []string{allowed.String()}})

	// Exactly one submission: the delisted cork must not reach the bridge.
	suite.gravityKeeper.EXPECT().
		CreateContractCallTx(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)

	corkKeeper.EndBlocker(ctx)

	require.Zero(countAuthorityCorksAt(corkKeeper, ctx, height),
		"both corks must leave the queue; a dropped cork must not be stranded")
}
