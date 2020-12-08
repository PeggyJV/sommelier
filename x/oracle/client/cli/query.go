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

	oracleQueryCmd.AddCommand([]*cobra.Command(
		GetCmdQueryExchangeRates(),
		GetCmdQueryActives(),
		GetCmdQueryParams(),
		GetCmdQueryFeederDelegation(),
		GetCmdQueryMissCounter(),
		GetCmdQueryAggregatePrevote(),
		GetCmdQueryAggregateVote(),
		GetCmdQueryVoteTargets(),
		GetCmdQueryTobinTaxes(),
	)...)

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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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
				return clientCtx.PrintOutput(&types.QueryExchangeRatesResponse{Rates: rates})
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryActives), nil)
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters), nil)
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
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

			var delegate sdk.AccAddress
			clientCtx.JSONMarshaler.MustUnmarshalJSON(res, &delegate)
			return clientCtx.PrintOutput(delegate)
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			params := types.NewQueryMissCounterParams(validator)
			bz, err := clientCtx.JSONMarshaler.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMissCounter), bz)
			if err != nil {
				return err
			}

			var missCounter int64
			clientCtx.JSONMarshaler.MustUnmarshalJSON(res, &missCounter)
			return clientCtx.PrintOutput(sdk.NewInt(missCounter))
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			params := types.NewQueryAggregatePrevoteParams(validator)
			bz, err := clientCtx.JSONMarshaler.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAggregatePrevote), bz)
			if err != nil {
				return err
			}

			var aggregatePrevote types.AggregateExchangeRatePrevote
			clientCtx.JSONMarshaler.MustUnmarshalJSON(res, &aggregatePrevote)
			return clientCtx.PrintOutput(aggregatePrevote)
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			valString := args[0]
			validator, err := sdk.ValAddressFromBech32(valString)
			if err != nil {
				return err
			}

			params := types.NewQueryAggregateVoteParams(validator)
			bz, err := clientCtx.JSONMarshaler.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAggregateVote), bz)
			if err != nil {
				return err
			}

			var aggregateVote types.AggregateExchangeRateVote
			clientCtx.JSONMarshaler.MustUnmarshalJSON(res, &aggregateVote)
			return clientCtx.PrintOutput(aggregateVote)
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryVoteTargets), nil)
			if err != nil {
				return err
			}

			var voteTargets Denoms
			clientCtx.JSONMarshaler.MustUnmarshalJSON(res, &voteTargets)
			return clientCtx.PrintOutput(voteTargets)
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			if len(args) == 0 {
				res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTobinTaxes), nil)
				if err != nil {
					return err
				}

				var tobinTaxes types.DenomList
				clientCtx.JSONMarshaler.MustUnmarshalJSON(res, &tobinTaxes)
				return clientCtx.PrintOutput(tobinTaxes)
			}

			denom := args[0]
			params := types.NewQueryTobinTaxParams(denom)

			bz, err := clientCtx.JSONMarshaler.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTobinTax), bz)
			if err != nil {
				return err
			}

			var tobinTax sdk.Dec
			clientCtx.JSONMarshaler.MustUnmarshalJSON(res, &tobinTax)
			return clientCtx.PrintOutput(tobinTax)
		},
	}

	return cmd
}
