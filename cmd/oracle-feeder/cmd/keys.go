package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/go-bip39"
	"github.com/spf13/cobra"
)

var (
	flagCoinType        = "coin-type"
	flagCoinTypeHelp    = "coin type number for HD derivation"
	flagCoinTypeDefault = uint32(118)
)

// keysCmd represents the keys command
func keysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"k"},
		Short:   "manage keys held by the relayer for each chain",
	}

	cmd.AddCommand(
		keysAddCmd(),
		keysRestoreCmd(),
		keysDeleteCmd(),
		keysListCmd(),
		keysShowCmd(),
		keysExportCmd(),
	)

	return cmd
}

// keysAddCmd respresents the `keys add` command
func keysAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [name]",
		Aliases: []string{"a"},
		Short:   "adds a key to the keychain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.GetClientContextFromCmd(cmd)
			mnemonic, err := CreateMnemonic()
			if err != nil {
				return err
			}

			coinType, _ := cmd.Flags().GetUint32(flagCoinType)

			info, err := ctx.Keyring.NewAccount(args[0], mnemonic, "", hd.CreateHDPath(coinType, 0, 0).String(), hd.Secp256k1)
			if err != nil {
				return err
			}

			ko := keyOutput{Mnemonic: mnemonic, Address: info.GetAddress().String()}

			out, err := json.Marshal(&ko)
			if err != nil {
				return err
			}

			fmt.Println(string(out))
			return nil
		},
	}
	cmd.Flags().Uint32(flagCoinType, flagCoinTypeDefault, flagCoinTypeHelp)
	return cmd
}

type keyOutput struct {
	Mnemonic string `json:"mnemonic" yaml:"mnemonic"`
	Address  string `json:"address" yaml:"address"`
}

// keysRestoreCmd respresents the `keys add` command
func keysRestoreCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "restore [name] [mnemonic]",
		Aliases: []string{"r"},
		Short:   "restores a mnemonic to the keychain",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.GetClientContextFromCmd(cmd)
			coinType, _ := cmd.Flags().GetUint32(flagCoinType)
			info, err := ctx.Keyring.NewAccount(
				args[0], args[1], "",
				hd.CreateHDPath(coinType, 0, 0).String(), hd.Secp256k1,
			)

			if err != nil {
				return err
			}
			fmt.Println(info.GetAddress().String())
			return nil
		},
	}

	cmd.Flags().Uint32(flagCoinType, flagCoinTypeDefault, flagCoinTypeHelp)
	return cmd
}

// keysDeleteCmd respresents the `keys delete` command
func keysDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [name]",
		Aliases: []string{"d"},
		Short:   "deletes a key from the keychain",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.GetClientContextFromCmd(cmd)
			if err := ctx.Keyring.Delete(args[0]); err != nil {
				return err
			}

			fmt.Printf("Key %s deleted for good...\n", args[0])
			return nil
		},
	}

	return cmd
}

// keysListCmd respresents the `keys list` command
func keysListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "lists keys from the keychain associated with a particular chain",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := client.GetClientContextFromCmd(cmd)
			infos, err := ctx.Keyring.List()
			if err != nil {
				return err
			}

			kos, err := cryptokeyring.Bech32KeysOutput(infos)
			if err != nil {
				return err
			}

			bz, err := ctx.LegacyAmino.MarshalJSON(kos)
			if err != nil {
				return err
			}

			fmt.Println(string(bz))
			return nil
		},
	}
	return cmd
}

// keysShowCmd respresents the `keys show` command
func keysShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [name]",
		Aliases: []string{"s"},
		Short:   "shows a key from the keychain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := client.GetClientContextFromCmd(cmd)
			info, err := ctx.Keyring.Key(args[0])
			if err != nil {
				return err
			}
			fmt.Println(info.GetAddress().String())
			return nil
		},
	}

	return cmd
}

// keysExportCmd respresents the `keys export` command
func keysExportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export [name]",
		Aliases: []string{"e"},
		Short:   "exports a privkey from the keychain associated with a particular chain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			panic("TODO IMPLEMENT")
		},
	}

	return cmd
}

// CreateMnemonic creates a new mnemonic
func CreateMnemonic() (string, error) {
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}
