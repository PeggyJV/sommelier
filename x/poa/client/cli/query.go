package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/peggyjv/sommelier/v10/x/poa/types"
)

// GetQueryCmd returns the cli query commands for the PoA module.
func GetQueryCmd() *cobra.Command {
	poaQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the proof-of-authority module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	poaQueryCmd.AddCommand([]*cobra.Command{
		queryParams(),
		queryAuthoritySet(),
		queryEffectivePower(),
		querySafeMode(),
	}...)

	return poaQueryCmd
}

func queryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parameters",
		Aliases: []string{"params"},
		Args:    cobra.NoArgs,
		Short:   "query poa params from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, err := types.NewQueryClient(ctx).Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryAuthoritySet() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "authority-set",
		Aliases: []string{"authorities"},
		Args:    cobra.NoArgs,
		Short:   "query the authority validator allowlist",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, err := types.NewQueryClient(ctx).AuthoritySet(cmd.Context(), &types.QueryAuthoritySetRequest{})
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryEffectivePower() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "effective-power [operator-address]",
		Aliases: []string{"power"},
		Args:    cobra.ExactArgs(1),
		Short:   "query a validator's effective (post-boost) consensus power",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, err := types.NewQueryClient(ctx).EffectivePower(cmd.Context(), &types.QueryEffectivePowerRequest{
				OperatorAddress: args[0],
			})
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// querySafeMode is the module's operational health check: while safe mode is
// active, gravity / cork / axelarcork are frozen. Without this the only way to
// observe that state was to scrape events or node logs.
func querySafeMode() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "safe-mode",
		Aliases: []string{"safemode"},
		Args:    cobra.NoArgs,
		Short:   "report whether value-bearing modules are frozen by authority-empty safe mode",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			res, err := types.NewQueryClient(ctx).SafeMode(cmd.Context(), &types.QuerySafeModeRequest{})
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
