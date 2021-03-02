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
	ilQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the impermanent loss module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ilQueryCmd.AddCommand([]*cobra.Command{
		GetCmdQueryStoploss(),
		GetCmdQueryParams(),
	}...)

	return ilQueryCmd

}

// GetCmdQueryStoploss implements the query rate command.
func GetCmdQueryStoploss() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stoploss [address] [[uniswap_pair]]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Query the stoploss positions for a given address",
		Long: strings.TrimSpace(`
Query the stoploss positions for a given address and uniswap pair

$ sommelier query il stoploss <address> <uniswap_pair>

Or, you can retrieve all the positions for the address

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
				Address:     address,
				UniswapPair: args[1],
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

			req := &types.QueryParamsRequest{}
			res, err := queryClient.Params(cmd.Context(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}
