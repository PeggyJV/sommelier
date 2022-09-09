package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"

	moduletestutil "github.com/peggyjv/sommelier/v4/testutil"
	auctionKeeper "github.com/peggyjv/sommelier/v4/x/auction/keeper"
	auctionTypes "github.com/peggyjv/sommelier/v4/x/auction/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	auctiontestutil "github.com/peggyjv/sommelier/v4/x/auction/testutil"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

var (
	permissionedFunder   = authtypes.NewEmptyModuleAccount("permissionedFunder")
	permissionedReciever = authtypes.NewEmptyModuleAccount("permissionedReciever")
)

type KeeperTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	auctionKeeper auctionKeeper.Keeper
	bankKeeper    *auctiontestutil.MockBankKeeper

	queryClient auctionTypes.QueryClient

	encCfg moduletestutil.TestEncodingConfig
}

func (suite *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(auctionTypes.StoreKey)
	tkey := sdk.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContext(key, tkey)
	ctx := testCtx.WithBlockHeader(tmproto.Header{Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// gomock initializations
	ctrl := gomock.NewController(suite.T())
	suite.bankKeeper = auctiontestutil.NewMockBankKeeper(ctrl)
	suite.ctx = ctx

	params := paramskeeper.NewKeeper(
		encCfg.Codec,
		codec.NewLegacyAmino(),
		key,
		tkey,
	)

	subSpace, found := params.GetSubspace(auctionTypes.ModuleName)
	suite.Assertions.True(found)

	suite.auctionKeeper = auctionKeeper.NewKeeper(
		encCfg.Codec,
		key,
		subSpace,
		suite.bankKeeper,
		map[string]bool{permissionedFunder.GetName(): true},
		map[string]bool{permissionedReciever.GetName(): true},
	)

	auctionTypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	auctionTypes.RegisterQueryServer(queryHelper, suite.auctionKeeper)
	queryClient := auctionTypes.NewQueryClient(queryHelper)

	suite.queryClient = queryClient
	suite.encCfg = encCfg
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
