package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
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

	var delAddrs []sdk.AccAddress
	var valAddrs []sdk.AccAddress
	k.IterateDelegateAddresses(ctx, func(del, val sdk.AccAddress) bool {
		delAddrs = append(delAddrs, del)
		valAddrs = append(valAddrs, val)
		return false
	})
	require.Equal(t, []sdk.AccAddress{Addrs[1], Addrs[3]}, delAddrs)
	require.Equal(t, []sdk.AccAddress{Addrs[0], Addrs[2]}, valAddrs)
}
func (suite *KeeperTestSuite) TestPrevotes() {

	// SetOracleDataPrevote
	// GetOracleDataPrevote
	// DeleteOracleDataPrevote
	// HasOracleDataPrevote
	// IterateOracleDataPrevotes

}
func (suite *KeeperTestSuite) TestVotes() {
	// SetOracleDataVote
	// GetOracleDataVote
	// DeleteOracleDataVote
	// HasOracleDataVote
	// IterateOracleDataVotes
}
func (suite *KeeperTestSuite) TestOracleData() {
	// SetOracleData
	// GetOracleData
	// DeleteOracleData
	// HasOracleData
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
	_ = json.Unmarshal([]byte(`{"pairs":[{"id":"0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc","reserve0":"149168502.127692","reserve1":"98914.728571090173772608","reserveUSD":"298465076.2726070256620710867523966","token0":{"decimals":"6","id":"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"},"token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token0Price":"1508.051473047154536981217254417823","token1Price":"0.0006631073394195288720366013950579345","totalSupply":"2.776928042758313339"},{"id":"0x0d4a11d5eeaac28ec3f61d100daf4d40471f1852","reserve0":"64388.028406657774747257","reserve1":"97110195.980325","reserveUSD":"194286280.6014220933152549537931338","token0":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token1":{"decimals":"6","id":"0xdac17f958d2ee523a2206206994597c13d831ec7"},"token0Price":"0.0006630408656543449235650418867565894","token1Price":"1508.202664119526393031771820896357","totalSupply":"1.809722368239841107"},{"id":"0xa478c2975ab1ea89e8196811f51a7b7ade33eb11","reserve0":"69666576.790413028896587687","reserve1":"46112.974337848761758699","reserveUSD":"139139424.218685770965493537968932","token0":{"decimals":"18","id":"0x6b175474e89094c44da98b954eedeac495271d0f"},"token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token0Price":"1510.780377773893936581451470552918","token1Price":"0.0006619095764756230959440895211021876","totalSupply":"1397637.606496882484285926"},{"id":"0xbb2b8038a1640196fbe3e38816f3e67cba72d940","reserve0":"3702.00263419","reserve1":"87578.554975281578154217","reserveUSD":"263970026.5503400767005757489435989","token0":{"decimals":"8","id":"0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"},"token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token0Price":"0.04227065216176338844401836197546882","token1Price":"23.65707527229888347304731067895859","totalSupply":"0.159912392338326095"},{"id":"0xd3d2e2692501a5c9ca623199d38826e513033a17","reserve0":"6723450.906791595454903831","reserve1":"83341.073955279932726161","reserveUSD":"251486860.9475119406522283429396433","token0":{"decimals":"18","id":"0x1f9840a85d5af5bf1d1762f925bdaddc4201f984"},"token1":{"decimals":"18","id":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},"token0Price":"80.67391728596320300811372353718798","token1Price":"0.01239558005414959861899222059722421","totalSupply":"381314.494595180173933409"}]}`), out)
	return
}
