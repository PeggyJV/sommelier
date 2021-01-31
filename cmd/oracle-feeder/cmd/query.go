package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/machinebox/graphql"
	oracle "github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// queryCmd represents the keys command
func queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "queries that can be run by the oracle-feeder",
	}

	cmd.AddCommand(queryUniswapData())

	return cmd
}

func queryUniswapData() *cobra.Command {
	return &cobra.Command{
		Use:     "uniswap-data",
		Aliases: []string{"ud"},
		Short:   "queries uniswap data for the transactions",
		RunE: func(cmd *cobra.Command, args []string) error {
			out, err := config.GetPairs(context.Background(), 100, 0)
			if err != nil {
				log.Fatal(err)
			}
			bz, err := json.MarshalIndent(out, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(bz))
			return nil
		},
	}
}

func queryParams() *cobra.Command {
	return &cobra.Command{
		Use:     "parameters",
		Aliases: []string{"params"},
		Args:    cobra.ExactArgs(0),
		Short:   "query oracle params from the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			clnt, err := config.OracleQueryClient()
			if err != nil {
				return err
			}
			defer config.StopGRPCClient()
			res, err := clnt.QueryParams(context.Background(), &oracle.QueryParamsRequest{})
			if err != nil {
				return err
			}
			bz, err := json.Marshal(res.Params)
			if err != nil {
				return err
			}
			fmt.Println(string(bz))
			return nil
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
			clnt, err := config.OracleQueryClient()
			if err != nil {
				return err
			}
			defer config.StopGRPCClient()
			res, err := clnt.QueryDelegeateAddress(context.Background(), &oracle.QueryDelegeateAddressRequest{Validator: args[0]})
			if err != nil {
				return err
			}
			fmt.Println(res.Delegate)
			return nil
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
			clnt, err := config.OracleQueryClient()
			if err != nil {
				return err
			}
			defer config.StopGRPCClient()
			res, err := clnt.QueryValidatorAddress(context.Background(), &oracle.QueryValidatorAddressRequest{Delegate: args[0]})
			if err != nil {
				return err
			}
			fmt.Println(res.Validator)
			return nil
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
			clnt, err := config.OracleQueryClient()
			if err != nil {
				return err
			}
			defer config.StopGRPCClient()
			res, err := clnt.QueryOracleDataPrevote(context.Background(), &oracle.QueryOracleDataPrevoteRequest{Validator: args[0]})
			if err != nil {
				return err
			}
			bz, err := json.Marshal(res.Hashes)
			if err != nil {
				return err
			}
			fmt.Println(string(bz))
			return nil
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
			ctx := client.GetClientContextFromCmd(cmd)
			clnt, err := config.OracleQueryClient()
			if err != nil {
				return err
			}
			defer config.StopGRPCClient()
			res, err := clnt.QueryOracleDataVote(context.Background(), &oracle.QueryOracleDataVoteRequest{Validator: args[0]})
			if err != nil {
				return err
			}
			bz, err := ctx.JSONMarshaler.MarshalJSON(res)
			if err != nil {
				return err
			}
			fmt.Println(string(bz))
			return nil
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
			clnt, err := config.OracleQueryClient()
			if err != nil {
				return err
			}
			defer config.StopGRPCClient()
			res, err := clnt.QueryVotePeriod(context.Background(), &oracle.QueryVotePeriodRequest{})
			if err != nil {
				return err
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return err
			}
			fmt.Println(string(bz))
			return nil
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
			clnt, err := config.OracleQueryClient()
			if err != nil {
				return err
			}
			defer config.StopGRPCClient()
			res, err := clnt.QueryMissCounter(context.Background(), &oracle.QueryMissCounterRequest{Validator: args[0]})
			if err != nil {
				return err
			}
			fmt.Println(res.MissCounter)
			return nil
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
			clnt, err := config.OracleQueryClient()
			if err != nil {
				return err
			}
			defer config.StopGRPCClient()
			typ, err := cmd.Flags().GetString("type")
			if err != nil {
				return err
			}
			res, err := clnt.OracleData(context.Background(), &oracle.QueryOracleDataRequest{Type: typ})
			if err != nil {
				return err
			}
			bz, err := ctx.JSONMarshaler.MarshalJSON(res.OracleData)
			if err != nil {
				return err
			}
			fmt.Println(string(bz))
			return nil
		},
	}
	cmd.Flags().StringP("type", "t", oracle.UniswapDataType, "type of oracle data to fetch")
	return cmd
}

// GetPairs returns the top N pairs from the Uniswap Subgraph
func (c Config) GetPairs(ctx context.Context, first, skip int) (*oracle.UniswapData, error) {
	req := graphql.NewRequest(fmt.Sprintf(`{ 
		pairs(first: %d, skip: %d, orderBy: volumeUSD, orderDirection: desc) {
			id
			reserveUSD
			totalSupply
			reserve0
			reserve1
			token0Price
			token1Price
			token0 {
				id
				decimals
			}
			token1 {
				id
				decimals
			}
		}
	}`, first, skip))
	out := &oracle.UniswapData{}
	err := c.graphClient.Run(ctx, req, out)
	return out, err
}

// OracleQueryClient returns a queryClient for the oracle GRPC methods
// CONTRACT: must stop the GRPC Client when done with the QueryClient
func (c Config) OracleQueryClient() (oracle.QueryClient, error) {
	err := c.StartGRPCClient()
	return oracle.NewQueryClient(c.grpcConn), err
}

// StartGRPCClient starts the GRCPClient
func (c Config) StartGRPCClient() error {
	grpcConn, err := grpc.Dial(c.ChainGRPC, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.grpcConn = grpcConn
	return nil
}

// StopGRPCClient stops the GRPCClient
func (c Config) StopGRPCClient() error {
	return c.grpcConn.Close()
}
