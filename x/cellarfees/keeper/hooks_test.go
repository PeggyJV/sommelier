package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
	auctionTypes "github.com/peggyjv/sommelier/v7/x/auction/types"
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

func (suite *KeeperTestSuite) TestHooksCountAccruesNoAuction() {
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
	expectedCounters := cellarfeesTypes.FeeAccrualCounters{
		Counters: []cellarfeesTypes.FeeAccrualCounter{
			{
				Denom: gravityFeeDenom,
				Count: 1,
			},
		},
	}

	// mocks
	suite.corkKeeper.EXPECT().HasCellarID(ctx, cellarID).Return(true)
	suite.gravityKeeper.EXPECT().ERC20ToDenomLookup(ctx, common.HexToAddress(event.TokenContract)).Return(false, gravityFeeDenom).Times(1)
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

	require.NotPanics(func() { hooks.AfterSendToCosmosEvent(ctx, event) })
	require.Equal(expectedCounters, cellarfeesKeeper.GetFeeAccrualCounters(ctx))
}

func (suite *KeeperTestSuite) TestHooksCountAccruesAuctionStarts() {
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
	cellarfeesKeeper.SetFeeAccrualCounters(ctx, cellarfeesTypes.FeeAccrualCounters{
		Counters: []cellarfeesTypes.FeeAccrualCounter{
			{
				Denom: gravityFeeDenom,
				Count: 1,
			},
		},
	})
	expectedCounters := cellarfeesTypes.FeeAccrualCounters{
		Counters: []cellarfeesTypes.FeeAccrualCounter{
			{
				Denom: gravityFeeDenom,
				Count: 0,
			},
		},
	}

	// mocks
	suite.corkKeeper.EXPECT().HasCellarID(ctx, cellarID).Return(true)
	suite.gravityKeeper.EXPECT().ERC20ToDenomLookup(ctx, common.HexToAddress(event.TokenContract)).Return(false, gravityFeeDenom)
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeesTypes.ModuleName).Return(feesAccount)
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), gravityFeeDenom).Return(sdk.NewCoin(gravityFeeDenom, event.Amount))
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	require.NotPanics(func() { hooks.AfterSendToCosmosEvent(ctx, event) })
	require.Equal(expectedCounters, cellarfeesKeeper.GetFeeAccrualCounters(ctx))
}

func (suite *KeeperTestSuite) TestHooksCountAccruesAuctionDoesNotStart() {
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
	cellarfeesKeeper.SetFeeAccrualCounters(ctx, cellarfeesTypes.FeeAccrualCounters{
		Counters: []cellarfeesTypes.FeeAccrualCounter{
			{
				Denom: gravityFeeDenom,
				Count: 1,
			},
		},
	})
	expectedCounters := cellarfeesTypes.FeeAccrualCounters{
		Counters: []cellarfeesTypes.FeeAccrualCounter{
			{
				Denom: gravityFeeDenom,
				Count: 2,
			},
		},
	}

	// mocks
	suite.corkKeeper.EXPECT().HasCellarID(ctx, cellarID).Return(true)
	suite.gravityKeeper.EXPECT().ERC20ToDenomLookup(ctx, common.HexToAddress(event.TokenContract)).Return(false, gravityFeeDenom)
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, cellarfeesTypes.ModuleName).Return(feesAccount)
	suite.auctionKeeper.EXPECT().GetActiveAuctions(ctx).Return([]*auctionTypes.Auction{})
	suite.bankKeeper.EXPECT().GetBalance(ctx, feesAccount.GetAddress(), gravityFeeDenom).Return(sdk.NewCoin(gravityFeeDenom, event.Amount))
	suite.auctionKeeper.EXPECT().BeginAuction(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(auctionTypes.ErrDenomCannotBeEmpty).Times(1)

	require.NotPanics(func() { hooks.AfterSendToCosmosEvent(ctx, event) })
	require.Equal(expectedCounters, cellarfeesKeeper.GetFeeAccrualCounters(ctx))
}
