package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	corkQueryCmd := &cobra.Command{
		Use:                        "cork",
		Short:                      "Querying commands for the cork module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	corkQueryCmd.AddCommand([]*cobra.Command{
		queryParams(),
		querySubmittedCorks(),
		queryCommitPeriod(),
		queryScheduledCorks(),
		queryCellarIDs(),
		queryScheduledBlockHeights(),
		queryScheduledCorksByBlockHeight(),
	}...)

	return corkQueryCmd

}

func queryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parameters",
		Aliases: []string{"params"},
		Args:    cobra.NoArgs,
		Short:   "query cork params from the chain",
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

func querySubmittedCorks() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "corks",
		Aliases: []string{"cs"},
		Args:    cobra.NoArgs,
		Short:   "query submitted corks from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QuerySubmittedCorksRequest{}

			res, err := queryClient.QuerySubmittedCorks(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryCommitPeriod() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "commit-period",
		Aliases: []string{"cp"},
		Args:    cobra.NoArgs,
		Short:   "query commit period data from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryCommitPeriodRequest{}

			res, err := queryClient.QueryCommitPeriod(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryCellarIDs() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cellar-ids",
		Aliases: []string{"cids"},
		Args:    cobra.NoArgs,
		Short:   "query managed cellar ids from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryCellarIDsRequest{}

			res, err := queryClient.QueryCellarIDs(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryScheduledCorks() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduled-corks",
		Aliases: []string{"scs"},
		Args:    cobra.NoArgs,
		Short:   "query scheduled corks from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryScheduledCorksRequest{}

			res, err := queryClient.QueryScheduledCorks(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryScheduledCorksByBlockHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduled-corks-by-block-height",
		Aliases: []string{"scbbh"},
		Args:    cobra.ExactArgs(1),
		Short:   "query scheduled corks from the chain by block height",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			height, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryScheduledCorksByBlockHeightRequest{
				BlockHeight: uint64(height),
			}

			res, err := queryClient.QueryScheduledCorksByBlockHeight(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryScheduledBlockHeights() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduled-block-heights",
		Aliases: []string{"scbhs"},
		Args:    cobra.NoArgs,
		Short:   "query scheduled cork block heights from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryScheduledBlockHeightsRequest{}

			res, err := queryClient.QueryScheduledBlockHeights(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
