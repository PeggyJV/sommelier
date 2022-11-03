package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	auctionTypes "github.com/peggyjv/sommelier/v4/x/auction/types"
)

func (suite *KeeperTestSuite) mockGetModuleAccount(ctx sdk.Context) {
	suite.accountKeeper.EXPECT().GetModuleAccount(ctx, auctionTypes.ModuleName).Return(authtypes.NewEmptyModuleAccount("mock"))
}

func (suite *KeeperTestSuite) mockGetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string, expectedOutput sdk.Coin) {
	suite.bankKeeper.EXPECT().GetBalance(ctx, addr, denom).Return(expectedOutput)
}

func (suite *KeeperTestSuite) mockSendCoinsFromModuleToModule(ctx sdk.Context, sender string, receiver string, amt sdk.Coins) {
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToModule(ctx, sender, receiver, amt).Return(nil)
}

func (suite *KeeperTestSuite) mockSendCoinsFromAccountToModule(ctx sdk.Context, senderAcct sdk.AccAddress, receiverModule string, amt sdk.Coins) {
	suite.bankKeeper.EXPECT().SendCoinsFromAccountToModule(ctx, senderAcct, receiverModule, amt).Return(nil)
}

func (suite *KeeperTestSuite) mockSendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, receiverAcct sdk.AccAddress, amt sdk.Coins) {
	suite.bankKeeper.EXPECT().SendCoinsFromModuleToAccount(ctx, senderModule, receiverAcct, amt).Return(nil)
}