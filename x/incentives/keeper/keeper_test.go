package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/golang/mock/gomock"
	moduletestutil "github.com/peggyjv/sommelier/v4/testutil"
	incentivestestutil "github.com/peggyjv/sommelier/v4/x/incentives/testutil"
	incentivesTypes "github.com/peggyjv/sommelier/v4/x/incentives/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx               sdk.Context
	incentivesKeeper  Keeper
	distributionKeepr *incentivestestutil.MockDistributionKeeper
	bankKeeper        *incentivestestutil.MockBankKeeper

	// queryClient cellarfeesTypes.QueryClient

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

	suite.distributionKeepr = incentivestestutil.NewMockDistributionKeeper(ctrl)
	suite.bankKeeper = incentivestestutil.NewMockBankKeeper(ctrl)
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
		suite.distributionKeepr,
		suite.bankKeeper,
	)

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
