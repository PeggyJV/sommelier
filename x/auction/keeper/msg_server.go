package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/peggyjv/sommelier/v4/x/auction/types"
)

var _ types.MsgServer = Keeper{}



// SubmitBid implements types.MsgServer
func (k Keeper) SubmitBid(c context.Context, msg *types.MsgSubmitBidRequest) (*types.MsgSubmitBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	signer := msg.MustGetSigner()
	
	// TODO: Figure out id update, probably need to access bid datastore and find latest bid id
	// Does ordering matter here? Should bid transfer be done first?
	k.SetBid(ctx, types.Bid{
		Id: 0,
		AuctionId: msg.GetAuctionId(),
		Bidder: msg.GetBidder(),
		MaxBid: msg.GetMaxBid(),
		MinimumAmount: msg.GetMinimumAmount(),
	})

	// TODO: Check if we can fulfil entire bid or not and decide on appropriate amt

	// TODO: Transfer funds
	


	// TODO: Emit bid + fulfilled bid
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
			sdk.NewEvent(
				types.EventTypeBid,
				sdk.NewAttribute(types.AttributeKeyAuctionId, fmt.Sprint(msg.GetAuctionId())),
				sdk.NewAttribute(types.AttributeKeyBidId, "nil" /* TODO: Populate once above code is done*/),
				sdk.NewAttribute(types.AttributeKeyBidder,  fmt.Sprint(msg.GetBidder())),
				sdk.NewAttribute(types.AttributeKeyMinimumAmount, msg.GetMinimumAmount().String()),
				sdk.NewAttribute(types.AttributeKeySigner, msg.GetSigner()),
			),
			sdk.NewEvent(
				types.EventTypeFulfilledBid,
				sdk.NewAttribute(types.AttributeKeyAuctionId, fmt.Sprint(msg.GetAuctionId())),
				sdk.NewAttribute(types.AttributeKeyBidId, "nil" /* TODO: Populate once above code is done*/),
				sdk.NewAttribute(types.AttributeKeyBidder, fmt.Sprint(msg.GetBidder())),
				sdk.NewAttribute(types.AttributeKeyMinimumAmount, msg.GetMinimumAmount().String()),
				sdk.NewAttribute(types.AttributeKeySigner, msg.GetSigner()),
				sdk.NewAttribute(types.AttributeKeyFulfilledPrice, "nil" /* TODO: Populate once above code is done*/),
				sdk.NewAttribute(types.AttributeKeyFulfilledAmount, "nil"/* TODO: Populate once above code is done*/),
			),
		},
	)

	// TODO: Populate below
	return &types.MsgSubmitBidResponse{&types.FulfilledBid{}}, nil
}