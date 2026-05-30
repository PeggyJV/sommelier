package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	types "github.com/peggyjv/sommelier/v9/x/cork/types/v2"
)

type stubPoaKeeper struct{ active bool }

func (s stubPoaKeeper) SafeModeActive(sdk.Context) bool { return s.active }

// When PoA is in authority-empty safe mode, ScheduleCork is rejected up front,
// before any orchestrator/cellar validation.
func (suite *KeeperTestSuite) TestScheduleCorkFrozenInSafeMode() {
	require := suite.Require()

	suite.corkKeeper.SetPoaKeeper(stubPoaKeeper{active: true})
	_, err := suite.corkKeeper.ScheduleCork(
		sdk.WrapSDKContext(suite.ctx),
		&types.MsgScheduleCorkRequest{},
	)
	require.Error(err)
	require.Contains(err.Error(), "safe mode")

	// Sanity: the EndBlocker is a no-op in safe mode and does not panic.
	suite.corkKeeper.EndBlocker(suite.ctx)
}
