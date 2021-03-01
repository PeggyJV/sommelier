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
	pairs       []types.OracleData
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

	suite.pairs = []types.OracleData{
		&types.UniswapPair{
			ID:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
			Reserve0:   sdk.MustNewDecFromStr("148681992.765143"),
			Reserve1:   sdk.MustNewDecFromStr("97709.503398661101176213"),
			ReserveUSD: sdk.MustNewDecFromStr("297632095.439861032964130850"),
			Token0: types.UniswapToken{
				Decimals: 6,
				ID:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
			},
			Token1: types.UniswapToken{
				Decimals: 18,
				ID:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			},
			Token0Price: sdk.MustNewDecFromStr("1521.673814659673802831"),
			Token1Price: sdk.MustNewDecFromStr("0.000657171064104597"),
			TotalSupply: sdk.MustNewDecFromStr("2.754869216896965436"),
		},
		&types.UniswapPair{
			ID:         "0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852",
			Reserve0:   sdk.NewDec(100),
			Reserve1:   sdk.NewDec(100),
			ReserveUSD: sdk.NewDec(100),
			Token0: types.UniswapToken{
				Decimals: 18,
				ID:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			},
			Token1: types.UniswapToken{
				Decimals: 6,
				ID:       "0xdac17f958d2ee523a2206206994597c13d831ec7",
			},
			Token0Price: sdk.NewDec(100),
			Token1Price: sdk.NewDec(100),
			TotalSupply: sdk.NewDec(100),
		},
	}
}
func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestAggregateData() {

	for _, pair := range suite.pairs {
		suite.app.OracleKeeper.SetAggregatedOracleData(suite.ctx, 10, pair)
	}

	oracleData := suite.app.OracleKeeper.GetAggregatedOracleData(suite.ctx, 10, types.UniswapDataType, suite.pairs[0].GetID())
	suite.Require().NotNil(oracleData)
	suite.Require().Equal(suite.pairs[0], oracleData)
}
