package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/golang/mock/gomock"
	"github.com/peggyjv/sommelier/v5/app/params"
	moduletestutil "github.com/peggyjv/sommelier/v5/testutil"
	incentivestestutil "github.com/peggyjv/sommelier/v5/x/incentives/testutil"
	incentivesTypes "github.com/peggyjv/sommelier/v5/x/incentives/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx                sdk.Context
	incentivesKeeper   Keeper
	distributionKeeper *incentivestestutil.MockDistributionKeeper
	bankKeeper         *incentivestestutil.MockBankKeeper
	mintKeeper         *incentivestestutil.MockMintKeeper

	encCfg moduletestutil.TestEncodingConfig
}

func (suite *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(incentivesTypes.StoreKey)
	tkey := sdk.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContext(key, tkey)
	ctx := testCtx.WithBlockHeader(tmproto.Header{Height: 5, Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// gomock initializations
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	suite.distributionKeeper = incentivestestutil.NewMockDistributionKeeper(ctrl)
	suite.bankKeeper = incentivestestutil.NewMockBankKeeper(ctrl)
	suite.mintKeeper = incentivestestutil.NewMockMintKeeper(ctrl)
	suite.ctx = ctx

	params := paramskeeper.NewKeeper(
		encCfg.Codec,
		codec.NewLegacyAmino(),
		key,
		tkey,
	)

	params.Subspace(incentivesTypes.ModuleName)
	subSpace, found := params.GetSubspace(incentivesTypes.ModuleName)
	suite.Assertions.True(found)

	suite.incentivesKeeper = NewKeeper(
		encCfg.Codec,
		key,
		subSpace,
		suite.distributionKeeper,
		suite.bankKeeper,
		suite.mintKeeper,
	)

	suite.encCfg = encCfg
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestGetAPY() {
	ctx, incentivesKeeper := suite.ctx, suite.incentivesKeeper
	require := suite.Require()
	distributionPerBlock := sdk.NewCoin(params.BaseCoinUnit, sdk.OneInt())
	blocksPerYear := 365 * 6
	bondedRatio := sdk.MustNewDecFromStr("0.2")
	stakingTotalSupply := sdk.NewInt(10_000_000)

	incentivesKeeper.SetParams(ctx, incentivesTypes.DefaultParams())
	suite.mintKeeper.EXPECT().GetParams(ctx).Return(mintTypes.Params{
		BlocksPerYear: uint64(blocksPerYear),
		MintDenom:     params.BaseCoinUnit,
	})
	suite.mintKeeper.EXPECT().BondedRatio(ctx).Return(bondedRatio)
	suite.mintKeeper.EXPECT().StakingTokenSupply(ctx).Return(stakingTotalSupply)

	// incentives disabled
	require.Equal(sdk.ZeroDec(), incentivesKeeper.GetAPY(ctx))

	// incentives enabled
	incentivesKeeper.SetParams(ctx, incentivesTypes.Params{
		DistributionPerBlock:   distributionPerBlock,
		IncentivesCutoffHeight: 1000,
	})
	expected := distributionPerBlock.Amount.Mul(sdk.NewInt(int64(blocksPerYear))).ToDec().Quo(stakingTotalSupply.ToDec().Mul(bondedRatio))
	require.Equal(expected, incentivesKeeper.GetAPY(ctx))
}
