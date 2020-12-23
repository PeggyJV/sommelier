package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/peggyjv/sommelier/x/oracle/types"

	"github.com/cosmos/cosmos-sdk/client"
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
		Short: "Query the current Luna exchange rate w.r.t an asset",
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

			if len(args) == 0 {
				res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryExchangeRates), nil)
				if err != nil {
					return err
				}

				rates, err := sdk.ParseDecCoins(string(res))
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(&types.QueryExchangeRatesResponse{Rates: rates})
			}

			bz, err := json.Marshal(types.NewQueryExchangeRateParams(args[0]))
			if err != nil {
				return err
			}
			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryExchangeRate), bz)
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil

		},
	}
	return cmd
}

// GetCmdQueryActives implements the query actives command.
func GetCmdQueryActives() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "actives",
		Args:  cobra.NoArgs,
		Short: "Query the active list of Terra assets recognized by the oracle",
		Long: strings.TrimSpace(`
Query the active list of Terra assets recognized by the types.

$ sommelier query oracle actives
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryActives), nil)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(string(res))
		},
	}

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current Oracle params",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters), nil)
			if err != nil {
				return err
			}
			return clientCtx.PrintString(string(res))
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

			bz, err := json.Marshal(types.NewQueryFeederDelegationParams(validator))
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryFeederDelegation), bz)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(string(res))
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

			bz, err := json.Marshal(types.NewQueryMissCounterParams(validator))
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMissCounter), bz)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(string(res))
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

			bz, err := json.Marshal(types.NewQueryAggregatePrevoteParams(validator))
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAggregatePrevote), bz)
			if err != nil {
				return err
			}

			var aggregatePrevote types.AggregateExchangeRatePrevote
			clientCtx.JSONMarshaler.MustUnmarshalJSON(res, &aggregatePrevote)
			return clientCtx.PrintProto(&aggregatePrevote)
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

			bz, err := json.Marshal(types.NewQueryAggregateVoteParams(validator))
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAggregateVote), bz)
			if err != nil {
				return err
			}

			var aggregateVote types.AggregateExchangeRateVote
			clientCtx.JSONMarshaler.MustUnmarshalJSON(res, &aggregateVote)
			return clientCtx.PrintProto(&aggregateVote)
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

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryVoteTargets), nil)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(string(res))
		},
	}

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

			if len(args) == 0 {
				res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTobinTaxes), nil)
				if err != nil {
					return err
				}

				out, err := sdk.ParseDecCoins(string(res))
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(&types.QueryTobinTaxesResponse{Rates: out})
			}

			bz, err := json.Marshal(types.NewQueryTobinTaxParams(args[0]))
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTobinTax), bz)
			if err != nil {
				return err
			}

			return clientCtx.PrintString(string(res))
		},
	}

	return cmd
}
