package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/machinebox/graphql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tmlog "github.com/tendermint/tendermint/libs/log"
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
	ETHWSAddr   string `yaml:"eth-ws-addr" json:"eth-ws-addr"`
	ETHRPCAddr  string `yaml:"eth-rpc-addr" json:"eth-rpc-addr"`
	GraphAddr   string `yaml:"graph-add" json:"graph-add"`
	SommKey     string `yaml:"somm-key" json:"somm-key"`
	EthKey      string `yaml:"eth-key" json:"eth-key"`
	GravityAddr string `yaml:"gravity-addr" json:"gravity-addr"`
	SommGRPC    string `yaml:"somm-grpc" json:"somm-grpc"`
	SommRPC     string `yaml:"somm-rpc" json:"somm-rpc"`
	ChainID     string `yaml:"chain-id" json:"chain-id"`
	GasPrices   string `yaml:"gas-prices" json:"gas-prices"`

	graphClient *graphql.Client
	grpcConn    *grpc.ClientConn
	gasPrices   sdk.DecCoin
	log         tmlog.Logger
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
			if ll, err := cmd.Flags().GetString("log-level"); err == nil {
				tmlog.AllowLevel(ll)
			}
			config.log = tmlog.NewTMLogger(tmlog.NewSyncWriter(os.Stdout))
		}
	}
	return nil
}

func validateConfig(c *Config) error {
	// gaddr := fmt.Sprintf("%s/subgraphs/name/davekaj/uniswap", c.GraphAddr)
	// fmt.Println(gaddr)
	c.graphClient = graphql.NewClient(c.GraphAddr)
	// c.graphClient.Log = func(s string) { fmt.Println(s) }
	coin, err := sdk.ParseDecCoin(c.GasPrices)
	if err != nil {
		return err
	}
	c.gasPrices = coin
	return nil
}

func defaultConfig() []byte {
	return Config{
		ETHRPCAddr:  "https://eth-goerli.alchemyapi.io/v2/vvvb-ZvZY_CX6D8ZwvMV8VH7E6vCGAbY",
		ETHWSAddr:   "wss://eth-goerli.ws.alchemyapi.io/v2/vvvb-ZvZY_CX6D8ZwvMV8VH7E6vCGAbY",
		GraphAddr:   "https://api.thegraph.com/",
		GravityAddr: "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D", // uniswap router
		SommGRPC:    "http://localhost:9090",
		SommRPC:     "http://localhost:26657",
		ChainID:     "somm",
		SommKey:     "feeder",
		EthKey:      "feeder",
		GasPrices:   "0.025stake",
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

// SubscribeSomm returns channel of events from the sommlier chain given a query (TX or BLOCK)
func (c Config) SubscribeSomm(ctx context.Context, clnt *rpchttp.HTTP, query string) (<-chan ctypes.ResultEvent, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ch, err := clnt.Subscribe(ctx, fmt.Sprintf("feeder-%s", genrandstr(8)), query, 1000)
	return ch, cancel, err
}

// NewSommRPCClient returns a new instance of the tendermint RPC client
func (c Config) NewSommRPCClient(timeout time.Duration) (*rpchttp.HTTP, error) {
	httpClient, err := libclient.DefaultHTTPClient(c.SommRPC)
	if err != nil {
		return nil, err
	}
	httpClient.Timeout = timeout
	return rpchttp.NewWithClient(c.SommRPC, "/websocket", httpClient)
}

// NewETHWSClient returns a new ethereum websocket client
func (c Config) NewETHWSClient() (*ethclient.Client, error) {
	return ethclient.Dial(c.ETHWSAddr)
}

// NewETHRPCClient returns a new ethereum rpc client
func (c Config) NewETHRPCClient() (*ethclient.Client, error) {
	return ethclient.Dial(c.ETHRPCAddr)
}
