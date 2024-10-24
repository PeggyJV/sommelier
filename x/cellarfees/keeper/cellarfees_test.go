package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	accounttypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/peggyjv/sommelier/v8/app/params"
	auctiontypes "github.com/peggyjv/sommelier/v8/x/auction/types"
	cellarfeestypes "github.com/peggyjv/sommelier/v8/x/cellarfees/types"
	cellarfeestypesv2 "github.com/peggyjv/sommelier/v8/x/cellarfees/types/v2"
)

func (suite *KeeperTestSuite) TestGetFeesAccount() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	expectedAddress := "somm1hqf42j6zxfnth4xpdse05wpnjjrgc864vwujxx"
	account := accounttypes.ModuleAccount{
		Name:        cellarfeestypes.ModuleName,
		BaseAccount: accounttypes.NewBaseAccountWithAddress(sdk.MustAccAddressFromBech32(expectedAddress)),
	}
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(&account)

	address := cellarfeesKeeper.GetFeesAccount(ctx).GetAddress().String()

	require.Equal(expectedAddress, address)
}

func (suite *KeeperTestSuite) TestGetFeeBalance() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	expectedDenom := "testdenom"
	expectedBalances := sdk.Coins{sdk.NewCoin(expectedDenom, sdk.NewInt(1000000))}
	account := accounttypes.ModuleAccount{
		Name:        cellarfeestypes.ModuleName,
		BaseAccount: accounttypes.NewBaseAccountWithAddress(sdk.MustAccAddressFromBech32("somm1hqf42j6zxfnth4xpdse05wpnjjrgc864vwujxx")),
	}
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeestypes.ModuleName).Return(&account)
	suite.bankKeeper.EXPECT().GetAllBalances(ctx, account.GetAddress()).Return(expectedBalances)

	balance, found := cellarfeesKeeper.GetFeeBalance(ctx, expectedDenom)

	require.True(found)
	require.Equal(expectedBalances[0], balance)
}

func (suite *KeeperTestSuite) TestGetBalanceUsdValue() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	// Test case 1: Non-zero exponent
	balance := sdk.NewCoin("testdenom", sdk.NewInt(1000000))
	usdPrice := sdk.NewDec(100)
	tokenPrice := auctiontypes.TokenPrice{
		Exponent: 6,
		UsdPrice: usdPrice,
	}

	expectedUsdValue := usdPrice.Mul(sdk.NewDecFromInt(balance.Amount).Quo(sdk.NewDec(10).Power(tokenPrice.Exponent)))
	usdValue := cellarfeesKeeper.GetBalanceUsdValue(ctx, balance, tokenPrice)

	require.Equal(expectedUsdValue, usdValue)

	// Test case 2: Zero exponent
	zeroExponentBalance := sdk.NewCoin("zerodenom", sdk.NewInt(500))
	zeroExponentUsdPrice := sdk.NewDec(50)
	zeroExponentTokenPrice := auctiontypes.TokenPrice{
		Exponent: 0,
		UsdPrice: zeroExponentUsdPrice,
	}

	expectedZeroExponentUsdValue := zeroExponentUsdPrice.Mul(sdk.NewDecFromInt(zeroExponentBalance.Amount))
	zeroExponentUsdValue := cellarfeesKeeper.GetBalanceUsdValue(ctx, zeroExponentBalance, zeroExponentTokenPrice)

	require.Equal(expectedZeroExponentUsdValue, zeroExponentUsdValue)
}

func (suite *KeeperTestSuite) TestGetEmission() {
	ctx, cellarfeesKeeper := suite.ctx, suite.cellarfeesKeeper
	require := suite.Require()

	remainingRewardsSupply := sdk.NewInt(1000000)
	previousSupplyPeak := sdk.NewInt(500000)
	cellarfeesParams := cellarfeestypesv2.DefaultParams()
	cellarfeesParams.RewardEmissionPeriod = 10
	cellarfeesKeeper.SetParams(ctx, cellarfeesParams)
	cellarfeesKeeper.SetLastRewardSupplyPeak(ctx, previousSupplyPeak)

	expectedEmissionAmount := remainingRewardsSupply.Quo(sdk.NewInt(int64(cellarfeesParams.RewardEmissionPeriod)))
	emission := cellarfeesKeeper.GetEmission(ctx, remainingRewardsSupply)

	require.Equal(sdk.NewCoins(sdk.NewCoin(params.BaseCoinUnit, expectedEmissionAmount)), emission)
}
