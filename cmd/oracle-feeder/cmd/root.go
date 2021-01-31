package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	kr "github.com/cosmos/cosmos-sdk/crypto/keyring"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/peggyjv/sommelier/app"
	"github.com/peggyjv/sommelier/app/params"
	"github.com/spf13/cobra"
	tmcli "github.com/tendermint/tendermint/libs/cli"
)

var (
	// FeederHome is the home directory for the oracle feeder
	FeederHome string
	AppName    string
	config     *Config
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	AppName = "oracle-feeder"
	FeederHome = filepath.Join(userHomeDir, ".oracle-feeder")
}

// NewRootCmd creates a new root command for simd. It is called once in the
// main function.
func NewRootCmd() (*cobra.Command, params.EncodingConfig) {
	
	encodingConfig := app.MakeEncodingConfig()
	keyring, err := kr.New(AppName, "test", FeederHome, os.Stdin)
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
		Use:   AppName,
		Short: "data feeder program for the sommelier oracle",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := initConfig(cmd); err != nil {
				return err
			}
			return client.SetCmdClientContext(cmd, initClientCtx)
		},
	}

	authclient.Codec = encodingConfig.Marshaler

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		keysCmd(),
		configCmd(),
		queryCmd(),
	)

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
