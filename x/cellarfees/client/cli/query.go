package cli

import (
	"context"
	"fmt"

	// "strings"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/v2.

	cellarfeestypes "github.com/peggyjv/sommelier/v7/x/cellarfees/types"
	types "github.com/peggyjv/sommelier/v7/x/cellarfees/types/v2"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group cellarfees queries under a subcommand
	cmd := &cobra.Command{
		Use:                        cellarfeestypes.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", cellarfeestypes.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdQueryModuleAccounts())
	cmd.AddCommand(CmdQueryLastRewardSupplyPeak())
	cmd.AddCommand(CmdQueryAPY())
	cmd.AddCommand(CmdQueryFeeTokenBalance())
	cmd.AddCommand(CmdQueryFeeTokenBalances())

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.QueryParams(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryModuleAccounts() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "module-accounts",
		Aliases: []string{"ma"},
		Short:   "shows the module accounts of the module",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.QueryModuleAccounts(
				context.Background(), &types.QueryModuleAccountsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryLastRewardSupplyPeak() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "last-reward-supply-peak",
		Aliases: []string{"lrsp"},
		Short:   "shows the previous SOMM reward supply peak amount used to calculate rewards per block",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.QueryLastRewardSupplyPeak(
				context.Background(), &types.QueryLastRewardSupplyPeakRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryAPY() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apy",
		Args:  cobra.NoArgs,
		Short: "query cellarfees APY",
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

func CmdQueryFeeTokenBalance() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fee-token-balance",
		Aliases: []string{"ftb"},
		Args:    cobra.ExactArgs(1),
		Short:   "query a fee tokens balance and its USD value in the cellarfees module",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			args := cmd.Flags().Args()

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryFeeTokenBalanceRequest{
				Denom: args[0],
			}

			res, err := queryClient.QueryFeeTokenBalance(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryFeeTokenBalances() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fee-token-balances",
		Aliases: []string{"ftbs"},
		Args:    cobra.NoArgs,
		Short:   "query all fee token balances and their USD values in the cellarfees module",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryFeeTokenBalancesRequest{}

			res, err := queryClient.QueryFeeTokenBalances(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
