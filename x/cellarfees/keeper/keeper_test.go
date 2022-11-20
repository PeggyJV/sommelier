package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/golang/mock/gomock"
	cellarfeestestutil "github.com/peggyjv/sommelier/v4/x/cellarfees/testutil"

	moduletestutil "github.com/peggyjv/sommelier/v4/testutil"
	cellarfeesTypes "github.com/peggyjv/sommelier/v4/x/cellarfees/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

var (
	feesAccount = authtypes.NewEmptyModuleAccount("cellarfees")
)

type KeeperTestSuite struct {
	suite.Suite

	ctx              sdk.Context
	cellarfeesKeeper Keeper
	accountKeeper    *cellarfeestestutil.MockAccountKeeper
	bankKeeper       *cellarfeestestutil.MockBankKeeper
	corkKeeper       *cellarfeestestutil.MockCorkKeeper
	gravityKeeper    *cellarfeestestutil.MockGravityKeeper
	auctionKeeper    *cellarfeestestutil.MockAuctionKeeper

	queryClient cellarfeesTypes.QueryClient

	encCfg moduletestutil.TestEncodingConfig
}

func (suite *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(cellarfeesTypes.StoreKey)
	tkey := sdk.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContext(key, tkey)
	ctx := testCtx.WithBlockHeader(tmproto.Header{Height: 5, Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// gomock initializations
	ctrl := gomock.NewController(suite.T())
	suite.bankKeeper = cellarfeestestutil.NewMockBankKeeper(ctrl)
	suite.accountKeeper = cellarfeestestutil.NewMockAccountKeeper(ctrl)
	suite.corkKeeper = cellarfeestestutil.NewMockCorkKeeper(ctrl)
	suite.gravityKeeper = cellarfeestestutil.NewMockGravityKeeper(ctrl)
	suite.auctionKeeper = cellarfeestestutil.NewMockAuctionKeeper(ctrl)
	suite.ctx = ctx

	params := paramskeeper.NewKeeper(
		encCfg.Codec,
		codec.NewLegacyAmino(),
		key,
		tkey,
	)

	params.Subspace(cellarfeesTypes.ModuleName)
	subSpace, found := params.GetSubspace(cellarfeesTypes.ModuleName)
	suite.Assertions.True(found)

	suite.cellarfeesKeeper = NewKeeper(
		encCfg.Codec,
		key,
		subSpace,
		suite.accountKeeper,
		suite.bankKeeper,
		suite.corkKeeper,
		suite.gravityKeeper,
		suite.auctionKeeper,
	)

	cellarfeesTypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	cellarfeesTypes.RegisterQueryServer(queryHelper, suite.cellarfeesKeeper)
	queryClient := cellarfeesTypes.NewQueryClient(queryHelper)

	suite.queryClient = queryClient
	suite.encCfg = encCfg
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestKeeperGettingSettingFeeAccrualCounters() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	expected := cellarfeesTypes.DefaultFeeAccrualCounters()
	cellarfeesKeeper.SetFeeAccrualCounters(ctx, expected)

	require.Equal(expected, cellarfeesKeeper.GetFeeAccrualCounters(ctx))
}

func (suite *KeeperTestSuite) TestKeeperGettingSettingLastRewardSupplyPeak() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	expected := sdk.NewInt(10 ^ 19 - 1)
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, expected)

	require.Equal(expected, cellarfeesKeeper.GetLastRewardSupplyPeak(ctx))
}
