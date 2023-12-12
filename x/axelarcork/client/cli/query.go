package cli

import (
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/v7/x/axelarcork/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	corkQueryCmd := &cobra.Command{
		Use:                        "axelarcork",
		Short:                      "Querying commands for the axelar cork module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	corkQueryCmd.AddCommand([]*cobra.Command{
		queryParams(),
		queryScheduledCorks(),
		queryCellarIDs(),
		queryCellarIDsByChainID(),
		queryScheduledBlockHeights(),
		queryScheduledCorksByBlockHeight(),
		queryScheduledCorksByID(),
		queryCorkResult(),
		queryCorkResults(),
		queryChainConfigurations(),
		queryAxelarContractCallNonces(),
		queryAxelayProxyUpgradeData(),
	}...)

	return corkQueryCmd

}

func queryParams() *cobra.Command {
	cmd := &cobra.Command{
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

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryCellarIDs() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cellar-ids",
		Aliases: []string{"cids"},
		Args:    cobra.NoArgs,
		Short:   "query all managed cellar ids from all chains",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			req := &types.QueryCellarIDsRequest{}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryCellarIDs(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryCellarIDsByChainID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cellar-ids-by-chain-id [chain-id]",
		Args:  cobra.ExactArgs(1),
		Short: "query managed cellar ids from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			chainID, err := math.ParseUint(args[0])
			if err != nil {
				return err
			}

			req := &types.QueryCellarIDsByChainIDRequest{
				ChainId: chainID.Uint64(),
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryCellarIDsByChainID(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryScheduledCorks() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduled-corks [chain-id]",
		Aliases: []string{"scs"},
		Args:    cobra.ExactArgs(1),
		Short:   "query scheduled corks from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			chainID, err := math.ParseUint(args[0])
			if err != nil {
				return err
			}

			req := &types.QueryScheduledCorksRequest{
				ChainId: chainID.Uint64(),
			}

			res, err := queryClient.QueryScheduledCorks(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryScheduledCorksByBlockHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduled-corks-by-block-height [chain-id]",
		Aliases: []string{"scbbh"},
		Args:    cobra.ExactArgs(1),
		Short:   "query scheduled corks from the chain by block height",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			height, err := math.ParseUint(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			chainID, err := math.ParseUint(args[0])
			if err != nil {
				return err
			}

			req := &types.QueryScheduledCorksByBlockHeightRequest{
				BlockHeight: height.Uint64(),
				ChainId:     chainID.Uint64(),
			}

			res, err := queryClient.QueryScheduledCorksByBlockHeight(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryScheduledBlockHeights() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduled-block-heights [chain-id]",
		Aliases: []string{"scbhs"},
		Args:    cobra.ExactArgs(1),
		Short:   "query scheduled cork block heights from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			chainID, err := math.ParseUint(args[0])
			if err != nil {
				return err
			}

			req := &types.QueryScheduledBlockHeightsRequest{
				ChainId: chainID.Uint64(),
			}

			res, err := queryClient.QueryScheduledBlockHeights(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryScheduledCorksByID() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduled-corks-by-id [chain-id]",
		Aliases: []string{"scbi"},
		Args:    cobra.ExactArgs(1),
		Short:   "query scheduled corks by their cork ID",
		Long:    "query scheduled corks by their cork ID, which is the keccak256 hash of the block height, target contract address, and encoded contract call concatenated",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			id := args[0]
			// the length of a keccak256 hash string
			if len(id) != 64 {
				return sdkerrors.New("", uint32(1), "invalid ID length, must be a keccak256 hash")
			}

			queryClient := types.NewQueryClient(ctx)
			chainID, err := math.ParseUint(args[0])
			if err != nil {
				return err
			}

			req := &types.QueryScheduledCorksByIDRequest{
				Id:      id,
				ChainId: chainID.Uint64(),
			}
			res, err := queryClient.QueryScheduledCorksByID(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryCorkResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cork-result [chain-id] [cork-id]",
		Aliases: []string{"cr"},
		Args:    cobra.ExactArgs(2),
		Short:   "query cork result from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			corkID := args[0]
			// the length of a keccak256 hash string
			if len(corkID) != 64 {
				return sdkerrors.New("", uint32(1), "invalid ID length, must be a keccak256 hash")
			}

			queryClient := types.NewQueryClient(ctx)
			chainID, err := math.ParseUint(args[1])
			if err != nil {
				return err
			}

			req := &types.QueryCorkResultRequest{
				Id:      corkID,
				ChainId: chainID.Uint64(),
			}

			res, err := queryClient.QueryCorkResult(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryCorkResults() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cork-results [chain-id]",
		Aliases: []string{"crs"},
		Args:    cobra.ExactArgs(1),
		Short:   "query cork results from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			chainID, err := math.ParseUint(args[0])
			if err != nil {
				return err
			}

			req := &types.QueryCorkResultsRequest{
				ChainId: chainID.Uint64(),
			}

			res, err := queryClient.QueryCorkResults(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryChainConfigurations() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "chain-configurations",
		Aliases: []string{"cfgs"},
		Args:    cobra.NoArgs,
		Short:   "query axelar chain configurations",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryChainConfigurationsRequest{}

			res, err := queryClient.QueryChainConfigurations(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryAxelarContractCallNonces() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "axelar-contract-call-nonces",
		Aliases: []string{"accn"},
		Args:    cobra.NoArgs,
		Short:   "query axelar contract call nonces from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			req := &types.QueryAxelarContractCallNoncesRequest{}

			res, err := queryClient.QueryAxelarContractCallNonces(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryAxelayProxyUpgradeData() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "axelar-proxy-upgrade-data",
		Aliases: []string{"apud"},
		Args:    cobra.NoArgs,
		Short:   "query axelar proxy upgrade data from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)

			req := &types.QueryAxelarProxyUpgradeDataRequest{}

			res, err := queryClient.QueryAxelarProxyUpgradeData(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
