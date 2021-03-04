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
	"github.com/peggyjv/sommelier/x/il/types"
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
	types.RegisterQueryServer(queryHelper, app.ILKeeper)
	suite.queryClient = types.NewQueryClient(queryHelper)

	_, _, addr := testdata.KeyTestPubAddr()
	suite.account = suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, addr)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestStoplossCRUD() {
	stoploss := types.Stoploss{
		UniswapPairID:       "0x3041cbd36888becc7bbcbc0045e3b1f144466f5f",
		LiquidityPoolShares: 10,
		MaxSlippage:         sdk.OneDec(),
		ReferencePairRatio:  sdk.OneDec(),
		ReceiverAddress:     "0x98950be0984d7cf7f5a098a6d8e53fc9c956d4bc",
	}

	suite.app.ILKeeper.SetStoplossPosition(suite.ctx, suite.account.GetAddress(), stoploss)
	suite.Require().True(suite.app.ILKeeper.HasStoplossPosition(suite.ctx, suite.account.GetAddress(), stoploss.UniswapPairID))

	res, found := suite.app.ILKeeper.GetStoplossPosition(suite.ctx, suite.account.GetAddress(), stoploss.UniswapPairID)
	suite.Require().True(found)
	suite.Require().Equal(stoploss, res)

	suite.app.ILKeeper.DeleteStoplossPosition(suite.ctx, suite.account.GetAddress(), stoploss.UniswapPairID)
	suite.Require().False(suite.app.ILKeeper.HasStoplossPosition(suite.ctx, suite.account.GetAddress(), stoploss.UniswapPairID))
}

func (suite *KeeperTestSuite) TestGetStoplossPositions() {
	stoplossPositions := []types.Stoploss{
		{
			UniswapPairID:       "pair",
			LiquidityPoolShares: 10,
			MaxSlippage:         sdk.OneDec(),
			ReferencePairRatio:  sdk.OneDec(),
		},
		{
			UniswapPairID:       "pair1",
			LiquidityPoolShares: 11,
			MaxSlippage:         sdk.OneDec(),
			ReferencePairRatio:  sdk.OneDec(),
		},
	}

	for _, position := range stoplossPositions {
		suite.app.ILKeeper.SetStoplossPosition(suite.ctx, suite.account.GetAddress(), position)
	}

	res := suite.app.ILKeeper.GetLPStoplossPositions(suite.ctx, suite.account.GetAddress())
	suite.Require().Equal(stoplossPositions, res)
}

func (suite *KeeperTestSuite) TestGetLPsStoplossPositions() {
	_, _, addr1 := testdata.KeyTestPubAddr()
	_, _, addr2 := testdata.KeyTestPubAddr()

	lpsStoplossPositions := types.LPsStoplossPositions{
		{
			Address: suite.account.GetAddress().String(),
			StoplossPositions: []types.Stoploss{
				{
					UniswapPairID:      "pair0",
					MaxSlippage:        sdk.OneDec(),
					ReferencePairRatio: sdk.OneDec(),
				},
				{
					UniswapPairID:      "pair1",
					MaxSlippage:        sdk.OneDec(),
					ReferencePairRatio: sdk.OneDec(),
				},
			},
		},
		{
			Address: addr1.String(),
			StoplossPositions: []types.Stoploss{
				{
					UniswapPairID:      "pair1",
					MaxSlippage:        sdk.OneDec(),
					ReferencePairRatio: sdk.OneDec(),
				},
				{
					UniswapPairID:      "pair2",
					MaxSlippage:        sdk.OneDec(),
					ReferencePairRatio: sdk.OneDec(),
				},
				{
					UniswapPairID:      "pair3",
					MaxSlippage:        sdk.OneDec(),
					ReferencePairRatio: sdk.OneDec(),
				},
			},
		},
		{
			Address: addr2.String(),
			StoplossPositions: []types.Stoploss{
				{
					UniswapPairID:      "pair8",
					MaxSlippage:        sdk.OneDec(),
					ReferencePairRatio: sdk.OneDec(),
				},
			},
		},
	}

	for _, lpPositions := range lpsStoplossPositions {
		addr, _ := sdk.AccAddressFromBech32(lpPositions.Address)
		for _, position := range lpPositions.StoplossPositions {
			suite.app.ILKeeper.SetStoplossPosition(suite.ctx, addr, position)
		}
	}

	res := suite.app.ILKeeper.GetLPsStoplossPositions(suite.ctx)
	suite.Require().Equal(lpsStoplossPositions.Sort(), res)
}
