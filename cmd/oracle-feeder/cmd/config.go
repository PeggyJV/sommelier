package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config show --home %s
$ %s cfg list`, AppName, FeederHome, AppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath := path.Join(FeederHome, "config.yaml")
			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				if _, err := os.Stat(FeederHome); os.IsNotExist(err) {
					return fmt.Errorf("home path does not exist: %s", FeederHome)
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
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s config init --home %s
$ %s cfg i`, AppName, FeederHome, AppName)),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath := path.Join(FeederHome, "config.yaml")

			// If the config doesn't exist...
			if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
				// And the config folder doesn't exist...
				// And the home folder doesn't exist
				if _, err := os.Stat(FeederHome); os.IsNotExist(err) {
					// Create the home folder
					if err = os.Mkdir(FeederHome, os.ModePerm); err != nil {
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
	ChainRPC        string `yaml:"chain-rpc" json:"chain-rpc"`
	ChainID         string `yaml:"chain-id" json:"chain-id"`

	graphClient *graphql.Client
}

// initConfig reads in config file and ENV variables if set.
func initConfig(cmd *cobra.Command) error {
	config = &Config{}
	cfgPath := path.Join(FeederHome, "config.yaml")
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
		}
	}
	return nil
}

func validateConfig(c *Config) error {
	c.graphClient = graphql.NewClient(c.UniswapSubgraph)
	return nil
}

func defaultConfig() []byte {
	return Config{
		UniswapSubgraph: "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2",
		ChainRPC:        "http://localhost:26657",
		ChainID:         "sommelier-test",
		SigningKey:      "mykey",
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
