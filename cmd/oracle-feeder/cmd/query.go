package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
	oracle "github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/spf13/cobra"
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
