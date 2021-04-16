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
		queryDelegateAddress(),
		queryValidatorAddress(),
		queryAllocationPrecommit(),
		queryAllocationCommit(),
		queryVotePeriod(),
		queryMissCounter(),
		queryAggregatedAllocationData(),
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

func queryDelegateAddress() *cobra.Command {
	return &cobra.Command{
		Use:     "delegate-address [validator-address]",
		Aliases: []string{"del"},
		Args:    cobra.ExactArgs(1),
		Short:   "query delegate address from the chain given validators address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryDelegateAddressRequest{Validator: args[0]}

			res, err := queryClient.QueryDelegateAddress(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}
}

func queryValidatorAddress() *cobra.Command {
	return &cobra.Command{
		Use:     "validator-address [delegate-address]",
		Aliases: []string{"val"},
		Args:    cobra.ExactArgs(1),
		Short:   "query validator address from the chain given the address that validator delegated to",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryValidatorAddressRequest{Delegate: args[0]}

			res, err := queryClient.QueryValidatorAddress(cmd.Context(), req)
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

func queryMissCounter() *cobra.Command {
	return &cobra.Command{
		Use:     "miss-counter [signer]",
		Aliases: []string{"mc"},
		Args:    cobra.ExactArgs(1),
		Short:   "query miss counter for a validator from the chain given its address or the delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryMissCounterRequest{
				Validator: args[0],
			}

			res, err := queryClient.QueryMissCounter(cmd.Context(), req)
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}
}

func queryAggregatedAllocationData() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "allocation-data [id]",
		Aliases: []string{"ad"},
		Args:    cobra.ExactArgs(1),
		Short:   "query aggregated allocation data from the chain given its type",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			dataType, _ := cmd.Flags().GetString("type")

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryAggregateDataRequest{
				Type: dataType,
				Id:   args[0],
			}

			res, err := queryClient.QueryAggregateData(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}
	cmd.Flags().StringP("type", "t", types.UniswapDataType, "type of allocation data to fetch")
	return cmd
}
