package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/peggyjv/sommelier/x/il/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	oracleQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the impermanent loss module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	oracleQueryCmd.AddCommand([]*cobra.Command{
		GetCmdQueryStoplossPositions(),
		GetCmdQueryParams(),
	}...)

	return oracleQueryCmd

}

// GetCmdQueryStoplossPositions implements the query rate command.
func GetCmdQueryStoplossPositions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stoploss [address] [il]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Query the stoploss positions for a given address", // TODO: update "Luna"
		Long: strings.TrimSpace(`
Query the current exchange rate of Luna with an asset. 
You can find the current list of active denoms by running

$ sommelier query il stoploss <address> <denom>

Or, can retrieve all the positions for an address

$ sommelier query il stoploss <address>
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			address := args[0]

			if len(args) == 1 {
				pageReq, err := client.ReadPageRequest(cmd.Flags())
				if err != nil {
					return err
				}

				req := &types.QueryStoplossPositionsRequest{
					Address:    address,
					Pagination: pageReq,
				}

				res, err := queryClient.StoplossPositions(cmd.Context(), req)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			req := &types.QueryStoplossRequest{
				Address: address,
			}

			res, err := queryClient.Stoploss(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "stoploss positions")
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current impermanent loss module parameters",
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryParametersRequest{}
			res, err := queryClient.Parameters(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}
