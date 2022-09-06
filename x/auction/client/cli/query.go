package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/peggyjv/sommelier/v4/x/auction/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	auctionQueryCmd := &cobra.Command{
		Use:                        "auction",
		Short:                      "Querying commands for the auction module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	auctionQueryCmd.AddCommand([]*cobra.Command{
		queryParams(),
		queryActiveAuction(),
		queryEndedAuction(),
		queryActiveAuctions(),
		queryActiveAuctionsByDenom(),
		queryEndedAuctions(),
		queryEndedAuctionsByDenom(),
		queryBid(),
		queryBidsByAuction(),
	}...)

	return auctionQueryCmd

}

func queryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parameters",
		Aliases: []string{"params"},
		Args:    cobra.NoArgs,
		Short:   "query auction params",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryParamsRequest{}

			res, err := queryClient.QueryParams(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryActiveAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "active-auction",
		Aliases: []string{"aa"},
		Args:    cobra.ExactArgs(1),
		Short:   "query an active auction",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			auctionID, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryActiveAuctionRequest{
				AuctionId: uint32(auctionID),
			}

			res, err := queryClient.QueryActiveAuction(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryEndedAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ended-auction",
		Aliases: []string{"ea"},
		Args:    cobra.ExactArgs(1),
		Short:   "query an ended auction",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			auctionID, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryEndedAuctionRequest{
				AuctionId: uint32(auctionID),
			}

			res, err := queryClient.QueryEndedAuction(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryActiveAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "active-auctions",
		Aliases: []string{"aas"},
		Args:    cobra.NoArgs,
		Short:   "query active auctions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryActiveAuctionsRequest{}

			res, err := queryClient.QueryActiveAuctions(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryActiveAuctionsByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "active-auctions-by-denom",
		Aliases: []string{"aad"},
		Args:    cobra.ExactArgs(1),
		Short:   "query the active auctions by a denom",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryActiveAuctionsRequest{}

			res, err := queryClient.QueryActiveAuctions(cmd.Context(), req)
			if err != nil {
				return err
			}

			denom := args[0]
			for _, auction := range res.GetAuctions() {
				if auction.StartingTokensForSale.Denom == denom {
					return ctx.PrintProto(auction)
				}
			}

			return fmt.Errorf("no active auction for denom: %s", denom)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryEndedAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ended-auctions",
		Aliases: []string{"eas"},
		Args:    cobra.NoArgs,
		Short:   "query ended auctions",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryEndedAuctionsRequest{}

			res, err := queryClient.QueryEndedAuctions(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryEndedAuctionsByDenom() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ended-auctions-by-denom",
		Aliases: []string{"ead"},
		Args:    cobra.ExactArgs(1),
		Short:   "query the ended auctions by a denom",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryEndedAuctionsRequest{}

			res, err := queryClient.QueryEndedAuctions(cmd.Context(), req)
			if err != nil {
				return err
			}

			denom := args[0]
			endedAuctions := make([]*types.Auction, 0, len(res.GetAuctions()))

			for _, auction := range res.GetAuctions() {
				if auction.StartingTokensForSale.Denom == denom {
					endedAuctions = append(endedAuctions, auction)
				}
			}

			if len(endedAuctions) == 0 {
				return fmt.Errorf("no ended auction for denom: %s", denom)
			}

			return ctx.PrintProto(&types.QueryEndedAuctionsResponse{Auctions: endedAuctions})
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bid",
		Aliases: []string{"b"},
		Args:    cobra.ExactArgs(2),
		Short:   "query bid by its auction id and its bid id",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			auctionID, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}

			bidID, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryBidRequest{
				AuctionId: uint32(auctionID),
				BidId:     bidID,
			}

			res, err := queryClient.QueryBid(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryBidsByAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bids-by-auction",
		Aliases: []string{"ba"},
		Args:    cobra.ExactArgs(1),
		Short:   "query the bids by an auction",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			auctionID, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryBidsByAuctionRequest{
				AuctionId: uint32(auctionID),
			}

			res, err := queryClient.QueryBidsByAuction(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
