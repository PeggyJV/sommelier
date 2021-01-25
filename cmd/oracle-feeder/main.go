package main

import (
	"context"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
	oracle "github.com/peggyjv/sommelier/x/oracle/types"
)

// Config is the configuration for the graphql querier
type Config struct {
	Client *graphql.Client
}

func main() {
	clnt := graphql.NewClient("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2")
	// clnt.Log = func(s string) { fmt.Println(s) }
	config := Config{clnt}
	out, err := config.GetPairs(context.Background(), 100, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
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
	err := c.Client.Run(ctx, req, out)
	return out, err
}
