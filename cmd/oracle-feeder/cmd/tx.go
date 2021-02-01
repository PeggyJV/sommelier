package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// txCmd represents the keys command
func txCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "transaction",
		Aliases: []string{"tx"},
		Short:   "transactions that can be run by the oracle-feeder",
	}

	cmd.AddCommand(txFeedOracle())

	return cmd
}

func txFeedOracle() *cobra.Command {
	return &cobra.Command{
		Use:     "feed-oracle",
		Aliases: []string{"feed"},
		Short:   "feeds the oracle with new uniswap data",
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

type Coordinator struct {
	// Current Round State
	// Last Uniswap Query
	// Last Uniswap Query Hash
	// Last Uniswap Query Hash Salt
	// Channel for listening to chain

	// Core loop
	//   listen over channel for new blocks and other events
	//   update round state every block
	//     on round start
	//       query uniswap data and store
	//       compute hash and send prevote
	//     on prevote event recieved
	//       send vote
}
