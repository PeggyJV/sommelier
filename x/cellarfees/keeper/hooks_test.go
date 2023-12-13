package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
	cellarfeesTypes "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
)

const gravityFeeDenom = "gravity0x1111111111111111111111111111111111111111"

func (suite *KeeperTestSuite) SetupHooksTests(ctx sdk.Context, cellarfeesKeeper Keeper) {
	cellarfeesKeeper.SetParams(ctx, cellarfeesTypes.DefaultParams())
	cellarfeesKeeper.SetFeeAccrualCounters(ctx, cellarfeesTypes.DefaultFeeAccrualCounters())

	// mocks
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, feesAccount.GetName()).Return(feesAccount).Times(1)
}

func (suite *KeeperTestSuite) TestHooksRecipientNotFeesAccountDoesNothing() {
	ctx, cellarfeesKeeper, require := suite.ctx, suite.cellarfeesKeeper, suite.Require()
	suite.SetupHooksTests(ctx, cellarfeesKeeper)

	hooks := Hooks{k: cellarfeesKeeper}
	event := gravitytypes.SendToCosmosEvent{
		CosmosReceiver: "fakeaddress",
	}

	require.NotPanics(func() { hooks.AfterSendToCosmosEvent(ctx, event) })
	require.Equal(cellarfeesTypes.DefaultFeeAccrualCounters(), cellarfeesKeeper.GetFeeAccrualCounters(ctx))
}

func (suite *KeeperTestSuite) TestHooksEventAmountZeroDoesNothing() {
	ctx, cellarfeesKeeper, require := suite.ctx, suite.cellarfeesKeeper, suite.Require()
	suite.SetupHooksTests(ctx, cellarfeesKeeper)

	hooks := Hooks{k: cellarfeesKeeper}
	event := gravitytypes.SendToCosmosEvent{
		CosmosReceiver: feesAccount.GetAddress().String(),
		Amount:         sdk.ZeroInt(),
	}

	require.NotPanics(func() { hooks.AfterSendToCosmosEvent(ctx, event) })
	require.Equal(cellarfeesTypes.DefaultFeeAccrualCounters(), cellarfeesKeeper.GetFeeAccrualCounters(ctx))
}

func (suite *KeeperTestSuite) TestHooksUnapprovedCellarDoesNothing() {
	ctx, cellarfeesKeeper, require := suite.ctx, suite.cellarfeesKeeper, suite.Require()
	suite.SetupHooksTests(ctx, cellarfeesKeeper)

	hooks := Hooks{k: cellarfeesKeeper}
	event := gravitytypes.SendToCosmosEvent{
		CosmosReceiver: feesAccount.GetAddress().String(),
		Amount:         sdk.OneInt(),
		EthereumSender: "0x0000000000000000000000000000000000000000",
	}
	cellarID := common.HexToAddress(event.EthereumSender)

	// mocks
	suite.corkKeeper.EXPECT().HasCellarID(ctx, cellarID).Return(false).Times(1)

	require.NotPanics(func() { hooks.AfterSendToCosmosEvent(ctx, event) })
	require.Equal(cellarfeesTypes.DefaultFeeAccrualCounters(), cellarfeesKeeper.GetFeeAccrualCounters(ctx))
}

func (suite *KeeperTestSuite) TestHooksDenomIsUsommDoesNothing() {
	ctx, cellarfeesKeeper, require := suite.ctx, suite.cellarfeesKeeper, suite.Require()
	suite.SetupHooksTests(ctx, cellarfeesKeeper)

	hooks := Hooks{k: cellarfeesKeeper}
	event := gravitytypes.SendToCosmosEvent{
		CosmosReceiver: feesAccount.GetAddress().String(),
		Amount:         sdk.OneInt(),
		EthereumSender: "0x0000000000000000000000000000000000000000",
		TokenContract:  "0x1111111111111111111111111111111111111111",
	}
	cellarID := common.HexToAddress(event.EthereumSender)

	// mocks
	suite.corkKeeper.EXPECT().HasCellarID(ctx, cellarID).Return(true)
	suite.gravityKeeper.EXPECT().ERC20ToDenomLookup(ctx, common.HexToAddress(event.TokenContract)).Return(true, "usomm").Times(1)

	require.NotPanics(func() { hooks.AfterSendToCosmosEvent(ctx, event) })
	require.Equal(cellarfeesTypes.DefaultFeeAccrualCounters(), cellarfeesKeeper.GetFeeAccrualCounters(ctx))
}
