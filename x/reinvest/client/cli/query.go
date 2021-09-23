package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/peggyjv/sommelier/x/allocation/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	allocationQueryCmd := &cobra.Command{
		Use:                        "allocation",
		Short:                      "Querying commands for the allocation module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	allocationQueryCmd.AddCommand([]*cobra.Command{
		queryParams(),
		queryAllocationPrecommit(),
		queryAllocationCommit(),
		queryVotePeriod(),
	}...)

	return allocationQueryCmd

}

func queryParams() *cobra.Command {
	return &cobra.Command{
		Use:     "parameters",
		Aliases: []string{"params"},
		Args:    cobra.NoArgs,
		Short:   "query allocation params from the chain",
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
}

func queryAllocationPrecommit() *cobra.Command {
	return &cobra.Command{
		Use:     "allocation-precommit [signer]",
		Aliases: []string{"precommit", "pc"},
		Args:    cobra.ExactArgs(1),
		Short:   "query allocation precommit from the chain, by either validator or delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryAllocationPrecommitRequest{Validator: args[0]}

			res, err := queryClient.QueryAllocationPrecommit(cmd.Context(), req)
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}
}

func queryAllocationCommit() *cobra.Command {
	return &cobra.Command{
		Use:     "allocation-commit [signer]",
		Aliases: []string{"commit"},
		Args:    cobra.ExactArgs(1),
		Short:   "query allocation commit from the chain, by either validator or delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryAllocationCommitRequest{Validator: args[0]}

			res, err := queryClient.QueryAllocationCommit(cmd.Context(), req)
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}
}

func queryVotePeriod() *cobra.Command {
	return &cobra.Command{
		Use:     "vote-period",
		Aliases: []string{"vp"},
		Args:    cobra.NoArgs,
		Short:   "query vote period data from the chain",
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
}
