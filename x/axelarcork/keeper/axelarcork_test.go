package keeper

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/golang/mock/gomock"
	"github.com/peggyjv/sommelier/v9/x/axelarcork/tests/mocks"
	"github.com/peggyjv/sommelier/v9/x/axelarcork/types"
	"github.com/stretchr/testify/require"
)

// Largest possible LegacyDec. Calculate by getting the largest 315-bit big.Int then placing a decimal point 18 digits from the right.
var maxDec = math.LegacyMustNewDecFromStr("66749594872528440074844428317798503581334516323645399060845050244444366430645.017188217565216767")

func (suite *KeeperTestSuite) TestSweepModuleAccountBalances() {
	require := suite.Require()

	tests := []struct {
		name           string
		moduleBalance  sdk.Coins
		feePool        sdk.DecCoins
		expectedSweep  sdk.Coins
		expectedPool   sdk.DecCoins
		expectTransfer bool
		expectBurn     bool
	}{
		{
			name:           "empty balance",
			moduleBalance:  sdk.NewCoins(),
			feePool:        sdk.NewDecCoins(),
			expectedSweep:  sdk.NewCoins(),
			expectedPool:   sdk.NewDecCoins(),
			expectTransfer: false,
			expectBurn:     false,
		},
		{
			name: "normal sweep single coin",
			moduleBalance: sdk.NewCoins(
				sdk.NewCoin("atom", sdk.NewInt(100)),
			),
			feePool: sdk.NewDecCoins(),
			expectedSweep: sdk.NewCoins(
				sdk.NewCoin("atom", sdk.NewInt(100)),
			),
			expectedPool: sdk.NewDecCoins(
				sdk.NewDecCoin("atom", sdk.NewInt(100)),
			),
			expectTransfer: true,
			expectBurn:     false,
		},
		{
			name: "normal sweep multiple coins",
			moduleBalance: sdk.NewCoins(
				sdk.NewCoin("atom", sdk.NewInt(100)),
				sdk.NewCoin("osmo", sdk.NewInt(200)),
			),
			feePool: sdk.NewDecCoins(),
			expectedSweep: sdk.NewCoins(
				sdk.NewCoin("atom", sdk.NewInt(100)),
				sdk.NewCoin("osmo", sdk.NewInt(200)),
			),
			expectedPool: sdk.NewDecCoins(
				sdk.NewDecCoin("atom", sdk.NewInt(100)),
				sdk.NewDecCoin("osmo", sdk.NewInt(200)),
			),
			expectTransfer: true,
			expectBurn:     false,
		},
		{
			name: "unsafe amount not swept",
			moduleBalance: sdk.NewCoins(
				sdk.NewCoin("atom", maxSafeInt.Add(sdk.OneInt())),
				sdk.NewCoin("osmo", sdk.NewInt(100)),
			),
			feePool: sdk.NewDecCoins(),
			expectedSweep: sdk.NewCoins(
				sdk.NewCoin("osmo", sdk.NewInt(100)),
			),
			expectedPool: sdk.NewDecCoins(
				sdk.NewDecCoin("osmo", sdk.NewInt(100)),
			),
			expectTransfer: true,
			expectBurn:     true,
		},
		{
			name: "only unsafe amount",
			moduleBalance: sdk.NewCoins(
				sdk.NewCoin("atom", maxSafeInt.Add(sdk.OneInt())),
			),
			feePool:        sdk.NewDecCoins(),
			expectedSweep:  sdk.NewCoins(),
			expectedPool:   sdk.NewDecCoins(),
			expectTransfer: false,
			expectBurn:     true,
		},
		{
			name: "existing pool near max, safe amounts not swept",
			moduleBalance: sdk.NewCoins(
				sdk.NewCoin("atom", sdk.NewInt(100)),
			),
			feePool: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("atom", maxDec),
			),
			expectedSweep: sdk.NewCoins(),
			expectedPool: sdk.NewDecCoins(
				sdk.NewDecCoinFromDec("atom", maxDec),
			),
			expectTransfer: false,
			expectBurn:     true,
		},
	}

	for _, tc := range tests {
		tc := tc
		suite.Run(tc.name, func() {
			// Create new mocks for each test case
			ctrl := gomock.NewController(suite.T())
			defer ctrl.Finish()

			// Reset mocks for this test case
			suite.accountKeeper = mocks.NewMockAccountKeeper(ctrl)
			suite.bankKeeper = mocks.NewMockBankKeeper(ctrl)
			suite.distributionKeeper = mocks.NewMockDistributionKeeper(ctrl)

			// Reset keeper with new mocks
			suite.axelarcorkKeeper = NewKeeper(
				suite.encCfg.Codec,
				sdk.NewKVStoreKey(types.StoreKey),
				suite.axelarcorkKeeper.paramSpace,
				suite.accountKeeper,
				suite.bankKeeper,
				suite.stakingKeeper,
				suite.transferKeeper,
				suite.distributionKeeper,
				suite.ics4wrapper,
				suite.gravityKeeper,
				suite.pubsubKeeper,
			)

			// Reset context
			ctx := suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

			// Setup module account expectation
			moduleAcct := authtypes.NewEmptyModuleAccount(types.ModuleName)
			suite.accountKeeper.EXPECT().
				GetModuleAccount(ctx, types.ModuleName).
				Return(moduleAcct)

			// Setup balance expectation
			suite.bankKeeper.EXPECT().
				GetAllBalances(ctx, moduleAcct.GetAddress()).
				Return(tc.moduleBalance)

			// Setup fee pool expectation
			feePool := distributionTypes.FeePool{
				CommunityPool: tc.feePool,
			}
			if !tc.moduleBalance.IsZero() {
				suite.distributionKeeper.EXPECT().
					GetFeePool(ctx).
					Return(feePool)
			}

			// Add burn expectation if needed
			if tc.expectBurn {
				balancesToBurn := tc.moduleBalance.Sub(tc.expectedSweep...)
				suite.bankKeeper.EXPECT().
					BurnCoins(ctx, types.ModuleName, balancesToBurn).
					Return(nil)
			}

			// Setup transfer expectation only if we expect a transfer
			if tc.expectTransfer {
				suite.bankKeeper.EXPECT().
					SendCoinsFromModuleToModule(
						ctx,
						types.ModuleName,
						distributionTypes.ModuleName,
						tc.expectedSweep,
					).Return(nil)

				suite.distributionKeeper.EXPECT().
					SetFeePool(ctx, gomock.Any()).
					Do(func(_ sdk.Context, fp distributionTypes.FeePool) {
						require.Equal(tc.expectedPool, fp.CommunityPool)
					})
			}

			// Execute
			suite.axelarcorkKeeper.SweepModuleAccountBalances(ctx)
		})
	}
}

func TestAddToCommunityPoolIfSafe(t *testing.T) {
	tests := []struct {
		name          string
		communityPool sdk.DecCoins
		coinToAdd     sdk.Coin
		expectSuccess bool
		expectedPool  sdk.DecCoins
	}{
		{
			name:          "normal addition",
			communityPool: sdk.NewDecCoins(sdk.NewDecCoin("atom", sdk.NewInt(100))),
			coinToAdd:     sdk.NewCoin("atom", sdk.NewInt(50)),
			expectSuccess: true,
			expectedPool:  sdk.NewDecCoins(sdk.NewDecCoin("atom", sdk.NewInt(150))),
		},
		{
			name:          "new denomination",
			communityPool: sdk.NewDecCoins(sdk.NewDecCoin("atom", sdk.NewInt(100))),
			coinToAdd:     sdk.NewCoin("osmo", sdk.NewInt(50)),
			expectSuccess: true,
			expectedPool: sdk.NewDecCoins(
				sdk.NewDecCoin("atom", sdk.NewInt(100)),
				sdk.NewDecCoin("osmo", sdk.NewInt(50)),
			),
		},
		{
			name:          "amount too large",
			communityPool: sdk.NewDecCoins(sdk.NewDecCoin("atom", sdk.NewInt(100))),
			coinToAdd:     sdk.NewCoin("atom", math.NewIntFromBigInt(maxSafeInt.BigInt()).Add(math.OneInt())),
			expectSuccess: false,
			expectedPool:  sdk.NewDecCoins(sdk.NewDecCoin("atom", sdk.NewInt(100))),
		},
		{
			name:          "existing amount near max, small addition causes overflow",
			communityPool: sdk.NewDecCoins(sdk.NewDecCoinFromDec("atom", maxDec)),
			coinToAdd:     sdk.NewCoin("atom", sdk.OneInt()),
			expectSuccess: false,
			expectedPool:  sdk.NewDecCoins(sdk.NewDecCoinFromDec("atom", maxDec)),
		},
		{
			name: "multiple denominations, overflow in one denom",
			communityPool: sdk.NewDecCoins(
				sdk.NewDecCoin("atom", sdk.NewInt(100)),
				sdk.NewDecCoinFromDec("osmo", maxDec),
			),
			coinToAdd:     sdk.NewCoin("osmo", sdk.OneInt()),
			expectSuccess: false,
			expectedPool: sdk.NewDecCoins(
				sdk.NewDecCoin("atom", sdk.NewInt(100)),
				sdk.NewDecCoinFromDec("osmo", maxDec),
			),
		},
		{
			name:          "zero coin addition",
			communityPool: sdk.NewDecCoins(sdk.NewDecCoin("atom", sdk.NewInt(100))),
			coinToAdd:     sdk.NewCoin("atom", sdk.ZeroInt()),
			expectSuccess: true,
			expectedPool:  sdk.NewDecCoins(sdk.NewDecCoin("atom", sdk.NewInt(100))),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pool := tc.communityPool
			success := AddToCommunityPoolIfSafe(&pool, tc.coinToAdd)
			require.Equal(t, tc.expectSuccess, success)
			require.Equal(t, tc.expectedPool, pool)
		})
	}
}

func TestIsSafeToAdd(t *testing.T) {
	if maxDec.BigInt().BitLen() != 315 {
		t.Fatalf("maxDec should have 315 bits but it's %d", maxDec.BigInt().BitLen())
	}

	maxInt := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
	quo := new(big.Int).Quo(maxDec.BigInt(), maxInt)
	t.Logf("quo: %s", quo.String())

	tests := []struct {
		name     string
		a        math.LegacyDec
		b        math.LegacyDec
		expected bool
	}{
		{
			name:     "small numbers",
			a:        sdk.NewDec(100),
			b:        sdk.NewDec(200),
			expected: true,
		},
		{
			name:     "zero addition 1",
			a:        sdk.NewDec(100),
			b:        sdk.ZeroDec(),
			expected: true,
		},
		{
			name:     "zero addition 2",
			a:        sdk.ZeroDec(),
			b:        sdk.NewDec(100),
			expected: true,
		},
		{
			name:     "zero addition 3",
			a:        sdk.ZeroDec(),
			b:        sdk.ZeroDec(),
			expected: true,
		},
		{
			name:     "max dec plus itself",
			a:        maxDec,
			b:        maxDec,
			expected: false, // Adding two max decimals will exceed 315 bits
		},
		{
			name:     "max dec plus zero is safe",
			a:        maxDec,
			b:        sdk.ZeroDec(),
			expected: true,
		},
		{
			name:     "half max dec plus itself is safe",
			a:        maxDec.QuoInt(sdk.NewInt(2)),
			b:        maxDec.QuoInt(sdk.NewInt(2)),
			expected: true, // Half + half should be safe
		},
		{
			name:     "max dec plus smallest dec isn't safe",
			a:        maxDec,
			b:        math.LegacyOneDec().QuoInt(sdk.NewInt(1000000000000000000)),
			expected: false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			result := IsSafeToAdd(tc.a, tc.b)
			require.Equal(t, tc.expected, result)
		})
	}
}
