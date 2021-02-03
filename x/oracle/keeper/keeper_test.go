package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
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

func (suite *KeeperTestSuite) TestDelegateAddresses() {
	ctx := suite.ctx
	k := suite.app.OracleKeeper
	t := suite.T()

	k.SetValidatorDelegateAddress(ctx, Addrs[0], Addrs[1])
	k.SetValidatorDelegateAddress(ctx, Addrs[2], Addrs[3])

	require.Equal(t, Addrs[0].String(), k.GetValidatorAddressFromDelegate(ctx, Addrs[1]).String())
	require.Equal(t, Addrs[2].String(), k.GetValidatorAddressFromDelegate(ctx, Addrs[3]).String())
	// TODO: test nil

	require.Equal(t, Addrs[1].String(), k.GetDelegateAddressFromValidator(ctx, Addrs[0]).String())
	require.Equal(t, Addrs[3].String(), k.GetDelegateAddressFromValidator(ctx, Addrs[2]).String())
	// TODO: test nil

	require.True(t, k.IsDelegateAddress(ctx, Addrs[1]))
	require.True(t, k.IsDelegateAddress(ctx, Addrs[3]))
	require.False(t, k.IsDelegateAddress(ctx, Addrs[4]))
}
func (suite *KeeperTestSuite) TestPrevotes() {
	ctx := suite.ctx
	k := suite.app.OracleKeeper
	t := suite.T()
	udjson := GetTestUniswapData().CannonicalJSON()
	zro := types.DataHash("0", udjson, Addrs[0])
	one := types.DataHash("1", udjson, Addrs[1])
	k.SetOracleDataPrevote(ctx, Addrs[0], types.NewMsgOracleDataPrevote([][]byte{zro}, Addrs[0]))
	k.SetOracleDataPrevote(ctx, Addrs[1], types.NewMsgOracleDataPrevote([][]byte{one}, Addrs[1]))

	require.Equal(t, zro, k.GetOracleDataPrevote(ctx, Addrs[0]).Hashes[0])
	require.Equal(t, one, k.GetOracleDataPrevote(ctx, Addrs[1]).Hashes[0])

	require.True(t, k.HasOracleDataPrevote(ctx, Addrs[0]))
	require.True(t, k.HasOracleDataPrevote(ctx, Addrs[1]))

	pvs := []*types.MsgOracleDataPrevote{}
	k.IterateOracleDataPrevotes(ctx, func(val sdk.AccAddress, msg *types.MsgOracleDataPrevote) bool {
		pvs = append(pvs, msg)
		return false
	})
	require.Equal(t, 2, len(pvs))

	k.DeleteOracleDataPrevote(ctx, Addrs[0])
	k.DeleteOracleDataPrevote(ctx, Addrs[1])

	require.False(t, k.HasOracleDataPrevote(ctx, Addrs[0]))
	require.False(t, k.HasOracleDataPrevote(ctx, Addrs[1]))
}
func (suite *KeeperTestSuite) TestVotes() {
	ctx := suite.ctx
	k := suite.app.OracleKeeper
	t := suite.T()

	uda, err := types.PackOracleData(testuniswapdata)
	require.NoError(t, err)
	zro := types.NewMsgOracleDataVote([]string{"0"}, []*codectypes.Any{uda}, Addrs[0])
	one := types.NewMsgOracleDataVote([]string{"1"}, []*codectypes.Any{uda}, Addrs[1])

	k.SetOracleDataVote(ctx, Addrs[0], zro)
	k.SetOracleDataVote(ctx, Addrs[1], one)

	require.Equal(t, zro, k.GetOracleDataVote(ctx, Addrs[0]))
	require.Equal(t, one, k.GetOracleDataVote(ctx, Addrs[1]))

	require.True(t, k.HasOracleDataVote(ctx, Addrs[0]))
	require.True(t, k.HasOracleDataVote(ctx, Addrs[1]))

	pvs := []*types.MsgOracleDataVote{}
	k.IterateOracleDataVotes(ctx, func(val sdk.AccAddress, msg *types.MsgOracleDataVote) bool {
		pvs = append(pvs, msg)
		return false
	})
	require.Equal(t, 2, len(pvs))

	k.DeleteOracleDataVote(ctx, Addrs[0])
	k.DeleteOracleDataVote(ctx, Addrs[1])

	require.False(t, k.HasOracleDataVote(ctx, Addrs[0]))
	require.False(t, k.HasOracleDataVote(ctx, Addrs[1]))
}
func (suite *KeeperTestSuite) TestOracleData() {
	ctx := suite.ctx
	k := suite.app.OracleKeeper
	t := suite.T()

	k.SetOracleData(ctx, testuniswapdata)
	require.True(t, k.HasOracleData(ctx, types.UniswapDataType))

	require.Equal(t, testuniswapdata, k.GetOracleData(ctx, types.UniswapDataType))

	k.DeleteOracleData(ctx, types.UniswapDataType)

	require.False(t, k.HasOracleData(ctx, types.UniswapDataType))
}
func (suite *KeeperTestSuite) TestVotePeriod() {
	// SetVotePeriodStart
	// GetVotePeriodStart
	// HasVotePeriodStart
}
func (suite *KeeperTestSuite) TestMissCounters() {
	// IncrementMissCounter
	// GetMissCounter
	// SetMissCounter
	// HasMissCounter
	// DeleteMissCounter
	// IterateMissCounters
}
func (suite *KeeperTestSuite) TestParams() {
	// GetParamSet
	// SetParams
}

//////////////
// FIXTURES //
//////////////

var (
	ConsPubKeys = []crypto.PubKey{
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
	}

	AccPubKeys = []crypto.PubKey{
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
		secp256k1.GenPrivKey().PubKey(),
	}

	Addrs = []sdk.AccAddress{
		sdk.AccAddress(AccPubKeys[0].Address()),
		sdk.AccAddress(AccPubKeys[1].Address()),
		sdk.AccAddress(AccPubKeys[2].Address()),
		sdk.AccAddress(AccPubKeys[3].Address()),
		sdk.AccAddress(AccPubKeys[4].Address()),
	}

	ValAddrs = []sdk.ValAddress{
		sdk.ValAddress(Addrs[0]),
		sdk.ValAddress(Addrs[1]),
		sdk.ValAddress(Addrs[2]),
		sdk.ValAddress(Addrs[3]),
		sdk.ValAddress(Addrs[4]),
	}

	InitTokens = sdk.TokensFromConsensusPower(200)
	InitCoins  = sdk.NewCoins(sdk.NewCoin("stake", InitTokens))
)

// GetTestUniswapData returns some realistic uniswap data for testing
func GetTestUniswapData() (out *types.UniswapData) {
	return testuniswapdata
}

var testuniswapdata = &types.UniswapData{
	Pairs: []types.UniswapPair{
		types.UniswapPair{
			Id:         "0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc",
			Reserve0:   "148681992.765143",
			Reserve1:   "97709.503398661101176213",
			ReserveUsd: "297632095.4398610329641308505948537",
			Token0: types.UniswapToken{
				Decimals: "6",
				Id:       "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
			},
			Token1: types.UniswapToken{
				Decimals: "18",
				Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			},
			Token0Price: "1521.673814659673802831294654725591",
			Token1Price: "0.0006571710641045975606382036411013578",
			TotalSupply: "2.754869216896965436",
		},
		types.UniswapPair{
			Id:         "0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852",
			Reserve0:   "64052.148928752869718841",
			Reserve1:   "97675312.070397",
			ReserveUsd: "195116448.3284569661435357469623931",
			Token0: types.UniswapToken{
				Decimals: "18",
				Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			},
			Token1: types.UniswapToken{
				Decimals: "6",
				Id:       "0xdac17f958d2ee523a2206206994597c13d831ec7",
			},
			Token0Price: "0.0006557660023915655416435420299954565",
			Token1Price: "1524.934193527904663377910319564834",
			TotalSupply: "1.80992106067496882",
		},
		types.UniswapPair{
			Id:         "0xa478c2975ab1ea89e8196811f51a7b7ade33eb11",
			Reserve0:   "69453224.061579510781012891",
			Reserve1:   "45584.711379804929448746",
			ReserveUsd: "138883455.9382328581978198800889056",
			Token0: types.UniswapToken{
				Decimals: "18",
				Id:       "0x6b175474e89094c44da98b954eedeac495271d0f",
			},
			Token1: types.UniswapToken{
				Decimals: "18",
				Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			},
			Token0Price: "1523.607849195440530703001829222234",
			Token1Price: "0.0006563368655051639699748790273756903",
			TotalSupply: "1387139.630260982742563912",
		},
		types.UniswapPair{
			Id:         "0xbb2b8038a1640196fbe3e38816f3e67cba72d940",
			Reserve0:   "3677.00380811",
			Reserve1:   "87751.610676879397734011",
			ReserveUsd: "267295991.4968480122080814610109883",
			Token0: types.UniswapToken{
				Decimals: "8",
				Id:       "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
			},
			Token1: types.UniswapToken{
				Decimals: "18",
				Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			},
			Token0Price: "0.04190240816945835175733332440293728",
			Token1Price: "23.86497682796369034462936135803174",
			TotalSupply: "0.159515800042228218",
		},
		types.UniswapPair{
			Id:         "0xd3d2e2692501a5c9ca623199d38826e513033a17",
			Reserve0:   "6727806.05368655342316",
			Reserve1:   "83188.043794789543616929",
			ReserveUsd: "253452781.802854374271358539030316",
			Token0: types.UniswapToken{
				Decimals: "18",
				Id:       "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984",
			},
			Token1: types.UniswapToken{
				Decimals: "18",
				Id:       "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			},
			Token0Price: "80.87467557577002901879232854141164",
			Token1Price: "0.01236481003331034057706013245120429",
			TotalSupply: "381065.480646149140512768",
		},
	},
}
