package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/peggyjv/sommelier/x/oracle/types"
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
		queryParams(),
		queryDelegeateAddress(),
		queryValidatorAddress(),
		queryOracleDataPrevote(),
		queryOracleDataVote(),
		queryVotePeriod(),
		queryMissCounter(),
		queryOracleData(),
	}...)

	return oracleQueryCmd

}

func queryParams() *cobra.Command {
	return &cobra.Command{
		Use:     "parameters",
		Aliases: []string{"params"},
		Args:    cobra.ExactArgs(0),
		Short:   "query oracle params from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(ctx).QueryParams(
				context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}
}

func queryDelegeateAddress() *cobra.Command {
	return &cobra.Command{
		Use:     "delegate-address [validator-address]",
		Aliases: []string{"del"},
		Args:    cobra.ExactArgs(1),
		Short:   "query delegeate address from the chain given validators address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(ctx).QueryDelegeateAddress(
				context.Background(), &types.QueryDelegeateAddressRequest{Validator: args[0]})
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
			res, err := types.NewQueryClient(ctx).QueryValidatorAddress(
				context.Background(), &types.QueryValidatorAddressRequest{Delegate: args[0]})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}
}

func queryOracleDataPrevote() *cobra.Command {
	return &cobra.Command{
		Use:     "oracle-prevote [signer]",
		Aliases: []string{"prevote", "pv"},
		Args:    cobra.ExactArgs(1),
		Short:   "query oracle data prevote from the chain, by either validator or delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(ctx).QueryOracleDataPrevote(
				context.Background(), &types.QueryOracleDataPrevoteRequest{Validator: args[0]})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}
}

func queryOracleDataVote() *cobra.Command {
	return &cobra.Command{
		Use:     "oracle-vote [signer]",
		Aliases: []string{"vote"},
		Args:    cobra.ExactArgs(1),
		Short:   "query oracle data vote from the chain, by either validator or delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(ctx).QueryOracleDataVote(
				context.Background(), &types.QueryOracleDataVoteRequest{Validator: args[0]})
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
		Args:    cobra.ExactArgs(0),
		Short:   "query vote period data from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(ctx).QueryVotePeriod(
				context.Background(), &types.QueryVotePeriodRequest{})
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
			res, err := types.NewQueryClient(ctx).QueryMissCounter(
				context.Background(), &types.QueryMissCounterRequest{Validator: args[0]})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}
}

func queryOracleData() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oracle-data",
		Aliases: []string{"od"},
		Args:    cobra.ExactArgs(0),
		Short:   "query consensus oracle data from the chain given its type",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.GetClientContextFromCmd(cmd)
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			typ, err := cmd.Flags().GetString("type")
			if err != nil {
				return err
			}
			res, err := types.NewQueryClient(ctx).OracleData(
				context.Background(), &types.QueryOracleDataRequest{Type: typ})
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}
	cmd.Flags().StringP("type", "t", types.UniswapDataType, "type of oracle data to fetch")
	return cmd
}
