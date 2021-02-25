package keeper_test

import (
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/stretchr/testify/suite"

	"github.com/peggyjv/sommelier/app"
	"github.com/peggyjv/sommelier/x/oracle/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	app     *app.SommelierApp
	account authtypes.AccountI

	queryClient types.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	checkTx := false
	app := app.Setup(checkTx)

	suite.ctx = app.BaseApp.NewContext(checkTx, tmproto.Header{Height: 1})
	suite.app = app

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.OracleKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	_, _, addr := testdata.KeyTestPubAddr()
	suite.account = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr)
}
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
