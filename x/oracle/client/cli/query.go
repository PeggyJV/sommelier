package cli

import (
	"context"
	"strings"

	"github.com/peggyjv/sommelier/x/oracle/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	oracleQueryCmd := &cobra.Command{
		Use:                        "oracle",
		Short:                      "Querying commands for the oracle module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	oracleQueryCmd.AddCommand([]*cobra.Command{
		GetCmdQueryExchangeRates(),
		GetCmdQueryActives(),
		GetCmdQueryParams(),
		GetCmdQueryFeederDelegation(),
		GetCmdQueryMissCounter(),
		GetCmdQueryAggregatePrevote(),
		GetCmdQueryAggregateVote(),
		GetCmdQueryVoteTargets(),
		GetCmdQueryTobinTaxes(),
	}...)

	return oracleQueryCmd

}

// GetCmdQueryExchangeRates implements the query rate command.
func GetCmdQueryExchangeRates() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exchange-rates [denom]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query the current Luna exchange rate w.r.t an asset", // TODO: update "Luna"
		Long: strings.TrimSpace(`
Query the current exchange rate of Luna with an asset. 
You can find the current list of active denoms by running

$ sommelier query oracle exchange-rates 

Or, can filter with denom

$ sommelier query oracle exchange-rates ukrw
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			if len(args) == 0 {
				// pageReq, err := client.ReadPageRequest(cmd.Flags())
				// if err != nil {
				// 	return err
				// }

				req := &types.QueryExchangeRatesRequest{
					// Pagination: pageReq,
				}

				res, err := queryClient.ExchangeRates(context.Background(), req)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			req := &types.QueryExchangeRateRequest{
				Denom: args[0],
			}

			res, err := queryClient.ExchangeRate(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "exchange rates")
	return cmd
}

// GetCmdQueryActives implements the query actives command.
func GetCmdQueryActives() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actives",
		Args:  cobra.NoArgs,
		Short: "Query the active list of Sommelier assets recognized by the oracle",
		Long: strings.TrimSpace(`
Query the active list of Sommelier assets recognized by the types.

$ sommelier query oracle actives
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// pageReq, err := client.ReadPageRequest(cmd.Flags())
			// if err != nil {
			// 	return err
			// }

			req := &types.QueryActivesRequest{
				// Pagination: pageReq,
			}

			res, err := queryClient.Actives(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "actives")
	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current oracle module parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryParametersRequest{}
			res, err := queryClient.Parameters(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

// GetCmdQueryFeederDelegation implements the query feeder delegation command
func GetCmdQueryFeederDelegation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feeder [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the oracle feeder delegate account",
		Long: strings.TrimSpace(`
Query the account the validator's oracle voting right is delegated to.

$ sommelier query oracle feeder terravaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			validator, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryFeederDelegationRequest{
				Validator: validator.String(),
			}

			res, err := queryClient.FeederDelegation(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

// GetCmdQueryMissCounter implements the query miss counter of the validator command
func GetCmdQueryMissCounter() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "miss [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the # of the miss count",
		Long: strings.TrimSpace(`
Query the # of vote periods missed in this oracle slash window.

$ sommelier query oracle miss terravaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryMissCounterRequest{
				Validator: validator.String(),
			}

			res, err := queryClient.MissCounter(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

// GetCmdQueryAggregatePrevote implements the query aggregate prevote of the validator command
func GetCmdQueryAggregatePrevote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-prevote [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query outstanding oracle aggregate prevote, filtered by voter address.",
		Long: strings.TrimSpace(`
Query outstanding oracle aggregate prevote, filtered by voter address.

$ sommelier query oracle aggregate-prevote terravaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryAggregatePrevoteRequest{
				Validator: validator.String(),
			}

			res, err := queryClient.AggregatePrevote(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

// GetCmdQueryAggregateVote implements the query aggregate prevote of the validator command
func GetCmdQueryAggregateVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aggregate-vote [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query outstanding oracle aggregate vote, filtered by voter address.",
		Long: strings.TrimSpace(`
Query outstanding oracle aggregate vote, filtered by voter address.

$ sommelier query oracle aggregate-vote terravaloper...
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryAggregateVoteRequest{
				Validator: validator.String(),
			}

			res, err := queryClient.AggregateVote(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}

// GetCmdQueryVoteTargets implements the query params command.
func GetCmdQueryVoteTargets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote-targets",
		Args:  cobra.NoArgs,
		Short: "Query the current Oracle vote targets",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// pageReq, err := client.ReadPageRequest(cmd.Flags())
			// if err != nil {
			// 	return err
			// }

			queryClient := types.NewQueryClient(clientCtx)

			req := &types.QueryVoteTargetsRequest{
				// Pagination: pageReq,
			}

			res, err := queryClient.VoteTargets(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "vote targets")
	return cmd
}

// GetCmdQueryTobinTaxes implements the query params command.
func GetCmdQueryTobinTaxes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tobin-taxes [denom]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Query the current Oracle tobin taxes.",
		Long: strings.TrimSpace(`
Query the current Oracle tobin taxes.

$ sommelier query oracle tobin-taxes

Or, can filter with denom

$ sommelier query oracle tobin-taxes ukrw

Or, can 
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			if len(args) == 0 {
				// pageReq, err := client.ReadPageRequest(cmd.Flags())
				// if err != nil {
				// 	return err
				// }

				req := &types.QueryTobinTaxesRequest{
					// Pagination: pageReq,
				}

				res, err := queryClient.TobinTaxes(context.Background(), req)
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			req := &types.QueryTobinTaxRequest{
				Denom: args[0],
			}

			res, err := queryClient.TobinTax(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	return cmd
}
