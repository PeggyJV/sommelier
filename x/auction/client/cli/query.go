package cli

import (
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
		queryCurrentAuction(),
		queryEndedAuction(),
		queryCurrentAuctions(),
		queryEndedAuctions(),
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
		Short:   "query auction params from the chain",
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

func queryCurrentAuction() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "current-auction",
		Aliases: []string{"ca"},
		Args:    cobra.ExactArgs(1),
		Short:   "query an ongoing auction from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			auction_id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryCurrentAuctionRequest{
				AuctionId: uint32(auction_id),
			}

			res, err := queryClient.QueryCurrentAuction(cmd.Context(), req)
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
		Short:   "query an ended auction from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			auction_id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryEndedAuctionRequest{
				AuctionId: uint32(auction_id),
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

func queryCurrentAuctions() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "current-auctions",
		Aliases: []string{"cas"},
		Args:    cobra.NoArgs,
		Short:   "query current auctions from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryCurrentAuctionsRequest{}

			res, err := queryClient.QueryCurrentAuctions(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
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
		Short:   "query ended auctions from the chain",
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

func queryBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bid",
		Aliases: []string{"b"},
		Args:    cobra.ExactArgs(2),
		Short:   "query bid from the chain by its auction id and its bid id",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			auction_id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			bid_id, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryBidRequest{
				AuctionId: uint32(auction_id),
				BidId:     uint64(bid_id),
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
		Use:     "bids-for-auction",
		Aliases: []string{"ba"},
		Args:    cobra.ExactArgs(1),
		Short:   "query the bids for an auction on the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			auction_id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryBidsByAuctionRequest{
				AuctionId: uint32(auction_id),
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
