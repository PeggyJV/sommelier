package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/peggyjv/sommelier/x/il/types"
)

func (suite *KeeperTestSuite) TestStoploss() {
	var (
		req         *types.QueryStoplossRequest
		expStoploss types.Stoploss
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryStoplossRequest{}
			},
			false,
		},
		{
			"stoploss not found",
			func() {
				req = &types.QueryStoplossRequest{
					Address: suite.account.GetAddress().String(),
				}
			},
			false,
		},
		{
			"success",
			func() {
				expStoploss = types.Stoploss{
					UniswapPairId:       "pair",
					LiquidityPoolShares: 10,
					MaxSlippage:         sdk.OneDec(),
					ReferencePairRatio:  sdk.OneDec(),
				}
				suite.app.ILKeeper.SetStoplossPosition(suite.ctx, suite.account.GetAddress(), expStoploss)

				req = &types.QueryStoplossRequest{
					Address:     suite.account.GetAddress().String(),
					UniswapPair: "pair",
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)

			res, err := suite.queryClient.Stoploss(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expStoploss, res.Stoploss)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestStoplossPositions() {
	var (
		req    *types.QueryStoplossPositionsRequest
		expRes *types.QueryStoplossPositionsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryStoplossPositionsRequest{}
			},
			false,
		},
		{
			"stoploss empty",
			func() {
				req = &types.QueryStoplossPositionsRequest{
					Address: suite.account.GetAddress().String(),
				}
				expRes = &types.QueryStoplossPositionsResponse{
					Pagination: &query.PageResponse{
						Total: 0,
					},
				}
			},
			true,
		},
		{
			"success",
			func() {
				expRes = &types.QueryStoplossPositionsResponse{
					Pagination: &query.PageResponse{
						Total: 2,
					},
					StoplossPositions: []types.Stoploss{
						{
							UniswapPairId:       "pair",
							LiquidityPoolShares: 10,
							MaxSlippage:         sdk.OneDec(),
							ReferencePairRatio:  sdk.OneDec(),
						},
						{
							UniswapPairId:       "pair1",
							LiquidityPoolShares: 11,
							MaxSlippage:         sdk.OneDec(),
							ReferencePairRatio:  sdk.OneDec(),
						},
					},
				}

				for _, position := range expRes.StoplossPositions {
					suite.app.ILKeeper.SetStoplossPosition(suite.ctx, suite.account.GetAddress(), position)
				}

				req = &types.QueryStoplossPositionsRequest{
					Address: suite.account.GetAddress().String(),
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)

			res, err := suite.queryClient.StoplossPositions(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes, res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestQueryParams() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	expParams := types.DefaultParams()

	res, err := suite.queryClient.Parameters(ctx, &types.QueryParametersRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(expParams, res.Params)
}
