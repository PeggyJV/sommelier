package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	kr "github.com/cosmos/cosmos-sdk/crypto/keyring"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/machinebox/graphql"
	"github.com/peggyjv/sommelier/app"
	"github.com/peggyjv/sommelier/app/params"
	oracle "github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/spf13/cobra"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

var (
	// FeederHome is the home directory for the oracle feeder
	FeederHome string
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	FeederHome = filepath.Join(userHomeDir, ".oracle-feeder")
}

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	encodingConfig := app.MakeEncodingConfig()
	keyring, err := kr.New("oracle-feeder", "test", FeederHome, os.Stdin)
	if err != nil {
		panic(err)
	}
	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(FeederHome).
		WithKeyring(keyring)

	rootCmd := &cobra.Command{
		Use:   "oracle-feeder",
		Short: "data feeder program for the sommelier oracle",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			return client.SetCmdClientContextHandler(initClientCtx, cmd)
		},
	}

	initRootCmd(rootCmd, encodingConfig)

	return rootCmd, encodingConfig
}

// Execute executes the root command.
func Execute(rootCmd *cobra.Command) error {
	// Create and set a client.Context on the command's Context. During the pre-run
	// of the root command, a default initialized client.Context is provided to
	// seed child command execution with values such as AccountRetriver, Keyring,
	// and a Tendermint RPC. This requires the use of a pointer reference when
	// getting and setting the client.Context. Ideally, we utilize
	// https://github.com/spf13/cobra/pull/1118.
	ctx := context.Background()
	ctx = context.WithValue(ctx, client.ClientContextKey, &client.Context{})

	executor := tmcli.PrepareBaseCmd(rootCmd, "", FeederHome)
	return executor.ExecuteContext(ctx)
}

func initRootCmd(rootCmd *cobra.Command, encodingConfig params.EncodingConfig) {
	authclient.Codec = encodingConfig.Marshaler

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		keysCmd(),
	)
}

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
	bz, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(bz))
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
