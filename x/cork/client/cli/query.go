package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/peggyjv/sommelier/v4/x/cork/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	allocationQueryCmd := &cobra.Command{
		Use:                        "cork",
		Short:                      "Querying commands for the cork module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	allocationQueryCmd.AddCommand([]*cobra.Command{
		queryParams(),
	}...)

	return allocationQueryCmd

}

func queryParams() *cobra.Command {
	return &cobra.Command{
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
}
