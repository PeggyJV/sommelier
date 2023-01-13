package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/peggyjv/sommelier/v4/x/incentives/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd() *cobra.Command {
	incentiveQueryCmd := &cobra.Command{
		Use:                        "incentives",
		Short:                      "Querying commands for the incentives module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	incentiveQueryCmd.AddCommand([]*cobra.Command{
		QueryParams(),
		QueryAPY(),
	}...)

	return incentiveQueryCmd

}

func QueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parameters",
		Aliases: []string{"params"},
		Args:    cobra.NoArgs,
		Short:   "query incentive params",
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

func QueryAPY() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "Gets APY of SOMM rewards from incentives",
		Aliases: []string{"apy"},
		Args:    cobra.NoArgs,
		Short:   "query incentives APY",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryAPYRequest{}

			res, err := queryClient.QueryAPY(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
