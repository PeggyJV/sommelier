package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
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
		queryDelegateAddress(),
		queryValidatorAddress(),
		queryOracleDataPrevote(),
		queryOracleDataVote(),
		queryVotePeriod(),
		queryMissCounter(),
		queryAggregatedOracleData(),
	}...)

	return oracleQueryCmd

}

func queryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parameters",
		Aliases: []string{"params"},
		Args:    cobra.NoArgs,
		Short:   "query oracle params from the chain",
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

func queryDelegateAddress() *cobra.Command {
	cmd := &cobra.Command{
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
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func queryValidatorAddress() *cobra.Command {
	cmd := &cobra.Command{
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
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func queryOracleDataPrevote() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oracle-prevote [signer]",
		Aliases: []string{"prevote", "pv"},
		Args:    cobra.ExactArgs(1),
		Short:   "query oracle data prevote from the chain, by either validator or delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryOracleDataPrevoteRequest{Validator: args[0]}

			res, err := queryClient.QueryOracleDataPrevote(cmd.Context(), req)
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func queryOracleDataVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oracle-vote [signer]",
		Aliases: []string{"vote"},
		Args:    cobra.ExactArgs(1),
		Short:   "query oracle data vote from the chain, by either validator or delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryOracleDataVoteRequest{Validator: args[0]}

			res, err := queryClient.QueryOracleDataVote(cmd.Context(), req)
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func queryVotePeriod() *cobra.Command {
	cmd := &cobra.Command{
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
			req := &types.QueryVotePeriodRequest{}

			res, err := queryClient.QueryVotePeriod(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func queryMissCounter() *cobra.Command {
	cmd := &cobra.Command{
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
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func queryAggregatedOracleData() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "oracle-data [id]",
		Aliases: []string{"od"},
		Args:    cobra.ExactArgs(1),
		Short:   "query aggregated oracle data from the chain given its type",
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
	cmd.Flags().StringP("type", "t", types.UniswapDataType, "type of oracle data to fetch")
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
