package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/peggyjv/sommelier/v4/x/pubsub/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group pubsub queries under a subcommand
	pubsubQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	pubsubQueryCmd.AddCommand([]*cobra.Command{
		CmdQueryParams(),
		CmdQueryPublisher(),
		CmdQueryPublishers(),
		CmdQuerySubscriber(),
		CmdQuerySubscribers(),
		CmdQueryPublisherIntent(),
		CmdQueryPublisherIntents(),
		CmdQueryPublisherIntentsByPublisherDomain(),
		CmdQueryPublisherIntentsBySubscriptionID(),
		CmdQuerySubscriberIntent(),
		CmdQuerySubscriberIntents(),
		CmdQuerySubscriberIntentsBySubscriberAddress(),
		CmdQuerySubscriberIntentsBySubscriptionID(),
		CmdQuerySubscriberIntentsByPublisherDomain(),
		CmdQueryDefaultSubscription(),
		CmdQueryDefaultSubscriptions(),
	}...)

	return pubsubQueryCmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query pubsub params",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryParamsRequest{}

			res, err := queryClient.Params(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryPublisher() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publisher",
		Args:  cobra.ExactArgs(1),
		Short: "Query publisher by publisher domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryPublisherRequest{
				PublisherDomain: args[0],
			}

			res, err := queryClient.QueryPublisher(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryPublishers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publishers",
		Args:  cobra.NoArgs,
		Short: "Query publishers",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryPublishersRequest{}

			res, err := queryClient.QueryPublishers(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQuerySubscriber() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriber",
		Args:  cobra.ExactArgs(1),
		Short: "Query subscriber by subscriber address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QuerySubscriberRequest{
				SubscriberAddress: args[0],
			}

			res, err := queryClient.QuerySubscriber(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQuerySubscribers() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscribers",
		Args:  cobra.NoArgs,
		Short: "Query subscribers",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QuerySubscribersRequest{}

			res, err := queryClient.QuerySubscribers(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryPublisherIntent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publisher-intent",
		Args:  cobra.ExactArgs(2),
		Short: "Query publisher intent by publisher domain and subscription ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryPublisherIntentRequest{
				PublisherDomain: args[0],
				SubscriptionId:  args[1],
			}

			res, err := queryClient.QueryPublisherIntent(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryPublisherIntents() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publisher-intents",
		Args:  cobra.NoArgs,
		Short: "Query publisher intents",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryPublisherIntentsRequest{}

			res, err := queryClient.QueryPublisherIntents(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryPublisherIntentsByPublisherDomain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publisher-intents-by-publisher-domain",
		Args:  cobra.ExactArgs(1),
		Short: "Query publisher intents by publisher domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryPublisherIntentsByPublisherDomainRequest{
				PublisherDomain: args[0],
			}

			res, err := queryClient.QueryPublisherIntentsByPublisherDomain(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryPublisherIntentsBySubscriptionID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publisher-intents-by-subscription-id",
		Args:  cobra.ExactArgs(1),
		Short: "Query publisher intents by subscription ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryPublisherIntentsBySubscriptionIDRequest{
				SubscriptionId: args[0],
			}

			res, err := queryClient.QueryPublisherIntentsBySubscriptionID(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQuerySubscriberIntent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriber-intent",
		Args:  cobra.ExactArgs(2),
		Short: "Query subscriber intent by subscriber address and subscription ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QuerySubscriberIntentRequest{
				SubscriberAddress: args[0],
				SubscriptionId:    args[1],
			}

			res, err := queryClient.QuerySubscriberIntent(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQuerySubscriberIntents() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriber-intents",
		Args:  cobra.NoArgs,
		Short: "Query subscriber intents",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QuerySubscriberIntentsRequest{}

			res, err := queryClient.QuerySubscriberIntents(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQuerySubscriberIntentsBySubscriberAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriber-intents-by-subscriber-address",
		Args:  cobra.ExactArgs(1),
		Short: "Query subscriber intents by subscriber address",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QuerySubscriberIntentsBySubscriberAddressRequest{
				SubscriberAddress: args[0],
			}

			res, err := queryClient.QuerySubscriberIntentsBySubscriberAddress(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQuerySubscriberIntentsBySubscriptionID() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriber-intents-by-subscription-id",
		Args:  cobra.ExactArgs(1),
		Short: "Query subscriber intents by subscription ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QuerySubscriberIntentsBySubscriptionIDRequest{
				SubscriptionId: args[0],
			}

			res, err := queryClient.QuerySubscriberIntentsBySubscriptionID(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQuerySubscriberIntentsByPublisherDomain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriber-intents-by-publisher-domain",
		Args:  cobra.ExactArgs(1),
		Short: "Query subscriber intents by publisher domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QuerySubscriberIntentsByPublisherDomainRequest{
				PublisherDomain: args[0],
			}

			res, err := queryClient.QuerySubscriberIntentsByPublisherDomain(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryDefaultSubscription() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default-subscription",
		Args:  cobra.ExactArgs(1),
		Short: "Query default subscription by subscription ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryDefaultSubscriptionRequest{
				SubscriptionId: args[0],
			}

			res, err := queryClient.QueryDefaultSubscription(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryDefaultSubscriptions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "default-subscriptions",
		Args:  cobra.NoArgs,
		Short: "Query default subscriptions",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(ctx)
			req := &types.QueryDefaultSubscriptionsRequest{}

			res, err := queryClient.QueryDefaultSubscriptions(cmd.Context(), req)
			if err != nil {
				return err
			}

			return ctx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
