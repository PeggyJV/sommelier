package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	moduletestutil "github.com/peggyjv/sommelier/v4/testutil"
	corktestutil "github.com/peggyjv/sommelier/v4/x/cork/testutil"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

var (
	sampleCellarHex  = "0xc0ffee254729296a45a3885639AC7E10F9d54979"
	sampleCellarAddr = common.HexToAddress(sampleCellarHex)
)

type KeeperTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	corkKeeper    Keeper
	gravityKeeper *corktestutil.MockGravityKeeper
	stakingKeeper *corktestutil.MockStakingKeeper
	validator     *corktestutil.MockValidatorI

	queryClient types.QueryClient

	encCfg moduletestutil.TestEncodingConfig
}

func (suite *KeeperTestSuite) SetupTest() {
	key := sdk.NewKVStoreKey(types.StoreKey)
	tkey := sdk.NewTransientStoreKey("transient_test")
	testCtx := testutil.DefaultContext(key, tkey)
	ctx := testCtx.WithBlockHeader(tmproto.Header{Height: 5, Time: tmtime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// gomock initializations
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	suite.gravityKeeper = corktestutil.NewMockGravityKeeper(ctrl)
	suite.stakingKeeper = corktestutil.NewMockStakingKeeper(ctrl)
	suite.validator = corktestutil.NewMockValidatorI(ctrl)
	suite.ctx = ctx

	params := paramskeeper.NewKeeper(
		encCfg.Codec,
		codec.NewLegacyAmino(),
		key,
		tkey,
	)

	params.Subspace(types.ModuleName)
	subSpace, found := params.GetSubspace(types.ModuleName)
	suite.Assertions.True(found)

	suite.corkKeeper = NewKeeper(
		encCfg.Codec,
		key,
		subSpace,
		suite.stakingKeeper,
		suite.gravityKeeper,
	)

	types.RegisterInterfaces(encCfg.InterfaceRegistry)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, suite.corkKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	suite.queryClient = queryClient
	suite.encCfg = encCfg
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetGetCellarIDsHappyPath() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	cellarIDSet := types.CellarIDSet{
		Ids: []string{sampleCellarHex},
	}
	expected := []common.Address{}
	for _, id := range cellarIDSet.Ids {
		expected = append(expected, common.HexToAddress(id))
	}
	corkKeeper.SetCellarIDs(ctx, cellarIDSet)
	actual := corkKeeper.GetCellarIDs(ctx)

	require.Equal(expected, actual)
	require.True(corkKeeper.HasCellarID(ctx, sampleCellarAddr))
}

func (suite *KeeperTestSuite) TestSetGetDeleteScheduledCork() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	testHeight := uint64(123)
	val := []byte("testaddress")
	expectedCork := types.Cork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
	}
	expectedID := expectedCork.IDHash(testHeight)
	actualID := corkKeeper.SetScheduledCork(ctx, testHeight, val, expectedCork)
	require.Equal(expectedID, actualID)
	actualCork, found := corkKeeper.GetScheduledCork(ctx, testHeight, actualID, val, sampleCellarAddr)
	require.True(found)
	require.Equal(expectedCork, actualCork)

	actualCorks := corkKeeper.GetScheduledCorks(ctx)
	require.Equal(&expectedCork, actualCorks[0].Cork)

	actualCorks = corkKeeper.GetScheduledCorksByID(ctx, actualID)
	require.Equal(&expectedCork, actualCorks[0].Cork)
	require.Equal(expectedID, actualCorks[0].Id)

	actualHeights := corkKeeper.GetScheduledBlockHeights(ctx)
	require.Equal(actualCorks[0].BlockHeight, actualHeights[0])

	actualCorks = corkKeeper.GetScheduledCorksByBlockHeight(ctx, testHeight)
	require.Equal(&expectedCork, actualCorks[0].Cork)
	require.Equal(testHeight, actualCorks[0].BlockHeight)
	require.Equal(expectedID, actualCorks[0].Id)

	corkKeeper.DeleteScheduledCork(ctx, testHeight, expectedID, sdk.ValAddress(val), sampleCellarAddr)
	require.Empty(corkKeeper.GetScheduledCorks(ctx))
}

func (suite *KeeperTestSuite) TestGetWinningVotes() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()
	testHeight := uint64(ctx.BlockHeight())
	cork := types.Cork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
	}
	_, bytes, err := bech32.DecodeAndConvert("somm1fcl08ymkl70dhyg3vmx4hjsqvxym7dawnp0zfp")
	require.NoError(err)
	require.Equal(20, len(bytes))
	corkKeeper.SetScheduledCork(ctx, testHeight, bytes, cork)

	suite.stakingKeeper.EXPECT().GetLastTotalPower(ctx).Return(sdk.NewInt(100))
	suite.stakingKeeper.EXPECT().Validator(ctx, gomock.Any()).Return(suite.validator)
	suite.validator.EXPECT().GetConsensusPower(gomock.Any()).Return(int64(100))
	suite.stakingKeeper.EXPECT().PowerReduction(ctx).Return(sdk.OneInt())

	winningScheduledVotes := corkKeeper.GetApprovedScheduledCorks(ctx, testHeight, sdk.ZeroDec())
	results := corkKeeper.GetCorkResults(ctx)
	require.Equal(cork, winningScheduledVotes[0])
	require.Equal(&cork, results[0].Cork)
	require.True(results[0].Approved)
	require.Equal("100.000000000000000000", results[0].ApprovalPercentage)

	// scheduled cork should be deleted at the scheduled height
	require.Empty(corkKeeper.GetScheduledCorksByBlockHeight(ctx, testHeight))
}

func (suite *KeeperTestSuite) TestInvalidationNonce() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	require.Zero(corkKeeper.GetLatestInvalidationNonce(ctx))

	corkKeeper.SetLatestInvalidationNonce(ctx, uint64(5))
	require.Equal(uint64(5), corkKeeper.GetLatestInvalidationNonce(ctx))

	corkKeeper.IncrementInvalidationNonce(ctx)
	require.Equal(uint64(6), corkKeeper.GetLatestInvalidationNonce(ctx))
}

func (suite *KeeperTestSuite) TestCorkResults() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	require.Empty(corkKeeper.GetCorkResults(ctx))

	testHeight := uint64(ctx.BlockHeight())
	cork := types.Cork{
		EncodedContractCall:   []byte("testcall"),
		TargetContractAddress: sampleCellarHex,
	}
	id := cork.IDHash(testHeight)
	result := types.CorkResult{
		Cork:               &cork,
		BlockHeight:        testHeight,
		Approved:           true,
		ApprovalPercentage: "100.00",
	}
	corkKeeper.SetCorkResult(ctx, id, result)
	actualResult, found := corkKeeper.GetCorkResult(ctx, id)
	require.True(found)
	require.Equal(result, actualResult)

	results := corkKeeper.GetCorkResults(ctx)
	require.Equal(&actualResult, results[0])

	corkKeeper.DeleteCorkResult(ctx, id)
	require.Empty(corkKeeper.GetCorkResults(ctx))
}

func (suite *KeeperTestSuite) TestParamSet() {
	ctx, corkKeeper := suite.ctx, suite.corkKeeper
	require := suite.Require()

	require.Panics(func() { corkKeeper.GetParamSet(ctx) })

	params := types.DefaultParams()
	corkKeeper.setParams(ctx, params)
	require.Equal(params, corkKeeper.GetParamSet(ctx))
}
