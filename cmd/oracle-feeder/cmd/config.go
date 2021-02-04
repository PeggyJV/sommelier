package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/machinebox/graphql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	libclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"cfg"},
		Short:   "manage configuration file",
	}

	cmd.AddCommand(
		configShowCmd(),
		configInitCmd(),
	)

	return cmd
}

// Command for printing current configuration
func configShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show",
		Aliases: []string{"s", "list", "l"},
		Short:   "Prints current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString("home")
			if err != nil {
				return err
			}
			cfgPath := path.Join(home, "config.yaml")
			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				if _, err := os.Stat(home); os.IsNotExist(err) {
					return fmt.Errorf("home path does not exist: %s", home)
				}
				return fmt.Errorf("config does not exist: %s", cfgPath)
			}

			out, err := yaml.Marshal(config)
			if err != nil {
				return err
			}

			fmt.Println(string(out))
			return nil
		},
	}

	return cmd
}

// Command for inititalizing an empty config at the --home location
func configInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "Creates a default home directory at path defined by --home",
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := cmd.Flags().GetString("home")
			if err != nil {
				return err
			}
			cfgPath := path.Join(home, "config.yaml")

			// If the config doesn't exist...
			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				// And the config folder doesn't exist...
				// And the home folder doesn't exist
				if _, err := os.Stat(home); os.IsNotExist(err) {
					// Create the home folder
					if err = os.Mkdir(home, os.ModePerm); err != nil {
						return err
					}
				}

				// Then create the file...
				f, err := os.Create(cfgPath)
				if err != nil {
					return err
				}
				defer f.Close()

				// And write the default config to that location...
				if _, err = f.Write(defaultConfig()); err != nil {
					return err
				}

				// And return no error...
				return nil
			}

			// Otherwise, the config file exists, and an error is returned...
			return fmt.Errorf("config already exists: %s", cfgPath)
		},
	}
	return cmd
}

// Config is the configuration for the graphql querier
type Config struct {
	UniswapSubgraph string `yaml:"uniswap-subgraph" json:"uniswap-subgraph"`
	SigningKey      string `yaml:"signing-key" json:"signing-key"`
	ChainGRPC       string `yaml:"chain-grpc" json:"chain-grpc"`
	ChainRPC        string `yaml:"chain-rpc" json:"chain-rpc"`
	ChainID         string `yaml:"chain-id" json:"chain-id"`
	GasPrices       string `yaml:"gas-prices" json:"gas-prices"`

	graphClient *graphql.Client
	grpcConn    *grpc.ClientConn
	gasPrices   sdk.DecCoin
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	config = &Config{}
	home, err := cmd.Flags().GetString("home")
	if err != nil {
		return err
	}
	cfgPath := path.Join(home, "config.yaml")
	if _, err := os.Stat(cfgPath); err == nil {
		viper.SetConfigFile(cfgPath)
		if err := viper.ReadInConfig(); err == nil {
			// read the config file bytes
			file, err := ioutil.ReadFile(viper.ConfigFileUsed())
			if err != nil {
				fmt.Println("Error reading file:", err)
				os.Exit(1)
			}

			// unmarshall them into the struct
			if err = yaml.Unmarshal(file, config); err != nil {
				fmt.Println("Error unmarshalling config:", err)
				os.Exit(1)
			}

			// ensure config has []*relayer.Chain used for all chain operations
			if err = validateConfig(config); err != nil {
				fmt.Println("Error parsing chain config:", err)
				os.Exit(1)
			}

			// TODO: set logger
		}
	}
	return nil
}

func validateConfig(c *Config) error {
	c.graphClient = graphql.NewClient(c.UniswapSubgraph)
	coin, err := sdk.ParseDecCoin(c.GasPrices)
	if err != nil {
		return err
	}
	c.gasPrices = coin

	return nil
}

func defaultConfig() []byte {
	return Config{
		UniswapSubgraph: "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2",
		ChainGRPC:       "http://localhost:9090",
		ChainRPC:        "http://localhost:26657",
		ChainID:         "somm",
		SigningKey:      "feeder",
		GasPrices:       "0.025stake",
	}.MustYAML()
}

// MustYAML returns the yaml string representation of the Paths
func (c Config) MustYAML() []byte {
	out, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	return out
}

// Subscribe returns channel of events from the sommlier chain given a query (TX or BLOCK)
func (c Config) Subscribe(ctx context.Context, clnt *rpchttp.HTTP, query string) (<-chan ctypes.ResultEvent, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ch, err := clnt.Subscribe(ctx, fmt.Sprintf("feeder-%s", genrandstr(8)), query, 1000)
	return ch, cancel, err
}

// NewRPCClient returns a new instance of the tendermint RPC client
func (c Config) NewRPCClient(timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(c.ChainRPC)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = timeout
	return rpchttp.NewWithClient(c.ChainRPC, "/websocket", httpClient)
}
