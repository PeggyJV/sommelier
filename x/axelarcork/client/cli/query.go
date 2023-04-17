package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/peggyjv/sommelier/v6/x/axelarcork/types"
	"github.com/spf13/cobra"
)

const (
	FlagAxelarChainID   = "axelar-chain-id"
	FlagAxelarChainName = "axelar-chain-name"
)

// AddChainFlagsToCmd adds common chain flags to a module command.
func AddChainFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String(FlagAxelarChainName, "", "The case sensitive name of the Axelar target chain")
	cmd.Flags().Uint64(FlagAxelarChainID, 0, "The Chain ID for the Axelar target EVM chain")
}

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	corkQueryCmd := &cobra.Command{
		Use:                        "axelar-cork",
		Short:                      "Querying commands for the axelar cork module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	corkQueryCmd.AddCommand([]*cobra.Command{
		queryParams(),
		queryScheduledCorks(),
		queryCellarIDs(),
		queryScheduledBlockHeights(),
		queryScheduledCorksByBlockHeight(),
		queryScheduledCorksByID(),
		queryCorkResult(),
		queryCorkResults(),
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
	AddChainFlagsToCmd(cmd)

	return cmd
}

func queryCellarIDs() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cellar-ids",
		Aliases: []string{"cids"},
		Args:    cobra.NoArgs,
		Short:   "query managed cellar ids from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			req := &types.QueryCellarIDsRequest{}
			name, err := cmd.Flags().GetString(FlagAxelarChainName)
			if err != nil {
				return err
			}
			chainID, err := cmd.Flags().GetUint64(FlagAxelarChainID)
			if err != nil {
				return err
			}

			if name != "" {
				req.ChainName = name
			}
			if chainID != 0 {
				req.ChainId = chainID
			}

			queryClient := types.NewQueryClient(ctx)

			res, err := queryClient.QueryCellarIDs(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	AddChainFlagsToCmd(cmd)

	return cmd
}

func queryScheduledCorks() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduled-corks",
		Aliases: []string{"scs"},
		Args:    cobra.NoArgs,
		Short:   "query scheduled corks from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryScheduledCorksRequest{}

			name, err := cmd.Flags().GetString(FlagAxelarChainName)
			if err != nil {
				return err
			}
			chainID, err := cmd.Flags().GetUint64(FlagAxelarChainID)
			if err != nil {
				return err
			}

			if name != "" {
				req.ChainName = name
			}
			if chainID != 0 {
				req.ChainId = chainID
			}

			res, err := queryClient.QueryScheduledCorks(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	AddChainFlagsToCmd(cmd)

	return cmd
}

func queryScheduledCorksByBlockHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduled-corks-by-block-height",
		Aliases: []string{"scbbh"},
		Args:    cobra.ExactArgs(1),
		Short:   "query scheduled corks from the chain by block height",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			height, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryScheduledCorksByBlockHeightRequest{
				BlockHeight: uint64(height),
			}

			name, err := cmd.Flags().GetString(FlagAxelarChainName)
			if err != nil {
				return err
			}
			chainID, err := cmd.Flags().GetUint64(FlagAxelarChainID)
			if err != nil {
				return err
			}

			if name != "" {
				req.ChainName = name
			}
			if chainID != 0 {
				req.ChainId = chainID
			}

			res, err := queryClient.QueryScheduledCorksByBlockHeight(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	AddChainFlagsToCmd(cmd)

	return cmd
}

func queryScheduledBlockHeights() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduled-block-heights",
		Aliases: []string{"scbhs"},
		Args:    cobra.NoArgs,
		Short:   "query scheduled cork block heights from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryScheduledBlockHeightsRequest{}

			name, err := cmd.Flags().GetString(FlagAxelarChainName)
			if err != nil {
				return err
			}
			chainID, err := cmd.Flags().GetUint64(FlagAxelarChainID)
			if err != nil {
				return err
			}

			if name != "" {
				req.ChainName = name
			}
			if chainID != 0 {
				req.ChainId = chainID
			}

			res, err := queryClient.QueryScheduledBlockHeights(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	AddChainFlagsToCmd(cmd)

	return cmd
}

func queryScheduledCorksByID() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "scheduled-corks-by-id",
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
			req := &types.QueryScheduledCorksByIDRequest{
				Id: id,
			}

			name, err := cmd.Flags().GetString(FlagAxelarChainName)
			if err != nil {
				return err
			}
			chainID, err := cmd.Flags().GetUint64(FlagAxelarChainID)
			if err != nil {
				return err
			}

			if name != "" {
				req.ChainName = name
			}
			if chainID != 0 {
				req.ChainId = chainID
			}

			res, err := queryClient.QueryScheduledCorksByID(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	AddChainFlagsToCmd(cmd)

	return cmd
}

func queryCorkResult() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cork-result",
		Aliases: []string{"cr"},
		Args:    cobra.ExactArgs(1),
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
			req := &types.QueryCorkResultRequest{
				Id: corkID,
			}

			name, err := cmd.Flags().GetString(FlagAxelarChainName)
			if err != nil {
				return err
			}
			chainID, err := cmd.Flags().GetUint64(FlagAxelarChainID)
			if err != nil {
				return err
			}

			if name != "" {
				req.ChainName = name
			}
			if chainID != 0 {
				req.ChainId = chainID
			}

			res, err := queryClient.QueryCorkResult(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	AddChainFlagsToCmd(cmd)

	return cmd
}

func queryCorkResults() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cork-results",
		Aliases: []string{"crs"},
		Args:    cobra.NoArgs,
		Short:   "query cork results from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryCorkResultsRequest{}

			name, err := cmd.Flags().GetString(FlagAxelarChainName)
			if err != nil {
				return err
			}
			chainID, err := cmd.Flags().GetUint64(FlagAxelarChainID)
			if err != nil {
				return err
			}

			if name != "" {
				req.ChainName = name
			}
			if chainID != 0 {
				req.ChainId = chainID
			}

			res, err := queryClient.QueryCorkResults(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	AddChainFlagsToCmd(cmd)

	return cmd
}
