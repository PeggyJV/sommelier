package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/sommelier/x/allocation/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	allocationQueryCmd := &cobra.Command{
		Use:                        "allocation",
		Short:                      "Querying commands for the allocation module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	allocationQueryCmd.AddCommand([]*cobra.Command{
		queryParams(),
		queryAllocationPrecommit(),
		queryAllocationCommit(),
		queryVotePeriod(),
	}...)

	return allocationQueryCmd

}

func queryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "parameters",
		Aliases: []string{"params"},
		Args:    cobra.NoArgs,
		Short:   "query allocation params from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			chorus, _ := sdk.ValAddressFromBech32("sommvaloper15urq2dtp9qce4fyc85m6upwm9xul30499el64g")
			fmt.Printf("\nchorus: %s\n", sdk.AccAddress(chorus.Bytes()).String())

			boubou, _ := sdk.ValAddressFromBech32("sommvaloper15055v5dh2hsn6fmy6ucday3pwn8m6gcmfz3kqr")
			fmt.Printf("boubou: %s\n", sdk.AccAddress(boubou.Bytes()).String())

			simply, _ := sdk.ValAddressFromBech32("sommvaloper1cgdlryczzgrk7d4kkeawqg7t6ldz4x84yu305c")
			fmt.Printf("simply: %s\n", sdk.AccAddress(simply.Bytes()).String())

			figment, _ := sdk.ValAddressFromBech32("sommvaloper1lexs4myxfp7k6n685qp6tw6mddkr2wetwddm2p")
			fmt.Printf("figment: %s\n", sdk.AccAddress(figment.Bytes()).String())

			tendermint, _ := sdk.ValAddressFromBech32("sommvaloper1ejqsr74xw6syh9nmukmqqtnnup4znwjmrkdmm0")
			fmt.Printf("tendermint: %s\n", sdk.AccAddress(tendermint.Bytes()).String())

			standard, _ := sdk.ValAddressFromBech32("sommvaloper173xq5ys8m7pvs2hesuz4ccx9mjuz2pthkq0zfj")
			fmt.Printf("standard: %s\n", sdk.AccAddress(standard.Bytes()).String())

			stakecito, _ := sdk.ValAddressFromBech32("sommvaloper1qe8uuf5x69c526h4nzxwv4ltftr73v7qeh9gwf")
			fmt.Printf("stakecito: %s\n", sdk.AccAddress(stakecito.Bytes()).String())

			imperator, _ := sdk.ValAddressFromBech32("sommvaloper1nm3xrar9j7dw6e9ua77ernarpc3axxr77alscy")
			fmt.Printf("imperator: %s\n", sdk.AccAddress(imperator.Bytes()).String())

			blockscape, _ := sdk.ValAddressFromBech32("sommvaloper1ju7p97r3atsqlpruy3a9dr25ltdc7qcjr6z6ff")
			fmt.Printf("blockscape: %s\n", sdk.AccAddress(blockscape.Bytes()).String())

			rbf, _ := sdk.ValAddressFromBech32("sommvaloper1thl5syhmscgnj7whdyrydw3w6vy80044ty5hup")
			fmt.Printf("rbf: %s\n\n", sdk.AccAddress(rbf.Bytes()).String())

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

func queryAllocationPrecommit() *cobra.Command {
	return &cobra.Command{
		Use:     "allocation-precommit [signer]",
		Aliases: []string{"precommit", "pc"},
		Args:    cobra.ExactArgs(1),
		Short:   "query allocation precommit from the chain, by either validator or delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryAllocationPrecommitRequest{Validator: args[0]}

			res, err := queryClient.QueryAllocationPrecommit(cmd.Context(), req)
			if err != nil {
				return err
			}
			return ctx.PrintProto(res)
		},
	}
}

func queryAllocationCommit() *cobra.Command {
	return &cobra.Command{
		Use:     "allocation-commit [signer]",
		Aliases: []string{"commit"},
		Args:    cobra.ExactArgs(1),
		Short:   "query allocation commit from the chain, by either validator or delegate address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryAllocationCommitRequest{Validator: args[0]}

			res, err := queryClient.QueryAllocationCommit(cmd.Context(), req)
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
		Args:    cobra.NoArgs,
		Short:   "query vote period data from the chain",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryCommitPeriodRequest{}

			res, err := queryClient.QueryCommitPeriod(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}
}
