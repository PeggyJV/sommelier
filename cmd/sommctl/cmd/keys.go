package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/spf13/cobra"
)

var (
	flagCoinType        = "coin-type"
	flagCoinTypeHelp    = "coin type number for HD derivation"
	flagCoinTypeDefault = uint32(118)
	defaultpass         = "secretsecretsecret"
)

// keysCmd represents the keys command
func keysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "keys",
		Aliases: []string{"k"},
		Short:   "manage keys held by sommctl for each chain",
	}

	sommCmd := &cobra.Command{
		Use:     "somm",
		Aliases: []string{"s"},
		Short:   "manage somm keys",
	}

	ethCmd := &cobra.Command{
		Use:     "eth",
		Aliases: []string{"e"},
		Short:   "manage eth keys",
	}

	cmd.AddCommand(
		sommCmd,
		ethCmd,
	)

	ethCmd.AddCommand(
		keysAddETHCmd(),
		keysRestoreETHCmd(),
		keysDeleteETHCmd(),
		keysListETHCmd(),
		keysShowETHCmd(),
		keysExportETHCmd(),
	)

	sommCmd.AddCommand(
		keysAddSommCmd(),
		keysRestoreSommCmd(),
		keysDeleteSommCmd(),
		keysListSommCmd(),
		keysShowSommCmd(),
		keysExportSommCmd(),
	)

	return cmd
}

func ethKeystore(cmd *cobra.Command) (*keystore.KeyStore, error) {
	home, err := cmd.Flags().GetString("home")
	if err != nil {
		return nil, err
	}
	return keystore.NewKeyStore(path.Join(home, "keyring-eth"), keystore.StandardScryptN, keystore.StandardScryptP), nil
}

func keysAddETHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"a"},
		Short:   "adds a new key to the eth keystore",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := ethKeystore(cmd)
			acct, err := ks.NewAccount(defaultpass)
			if err != nil {
				return err
			}
			fmt.Println(acct.Address.String())
			return nil
		},
	}
	return cmd
}

func keysRestoreETHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "restore [file]",
		Aliases: []string{"r"},
		Short:   "restores a key to the keychain from a json file",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := ethKeystore(cmd)
			if err != nil {
				return err
			}
			f, err := os.Open(args[0])
			if err != nil {
				return err
			}
			bz, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			acct, err := ks.Import(bz, defaultpass, defaultpass)
			if err != nil {
				return err
			}
			fmt.Println(acct.Address.String())
			return nil
		},
	}
	return cmd
}

func keysDeleteETHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [index]",
		Aliases: []string{"d"},
		Short:   "deletes a key from the eth keystore",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := ethKeystore(cmd)
			if err != nil {
				return err
			}
			i, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			accts := ks.Accounts()
			if len(accts) < int(i-1) {
				return fmt.Errorf("account %d not in keystore, only %d accounts", i, len(accts))
			}
			return ks.Delete(accts[i], defaultpass)
		},
	}
	return cmd
}

func keysListETHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "lists keys in the keychain",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := ethKeystore(cmd)
			if err != nil {
				return err
			}
			for i, acct := range ks.Accounts() {
				fmt.Println(i, acct.Address.String())
			}
			return nil
		},
	}
	return cmd
}

func keysShowETHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [index]",
		Aliases: []string{"s"},
		Short:   "shows a key from the keychain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := ethKeystore(cmd)
			if err != nil {
				return err
			}
			i, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			accts := ks.Accounts()
			if len(accts) < int(i-1) {
				return fmt.Errorf("account %d not in keystore, only %d accounts", i, len(accts))
			}
			fmt.Println(accts[i].Address.String())
			return nil
		},
	}
	return cmd
}

func keysExportETHCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "export [index]",
		Aliases: []string{"e"},
		Short:   "exports a key from the keychain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ks, err := ethKeystore(cmd)
			if err != nil {
				return err
			}
			i, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			accts := ks.Accounts()
			if len(accts) < int(i-1) {
				return fmt.Errorf("account %d not in keystore, only %d accounts", i, len(accts))
			}
			bz, err := ks.Export(accts[i], defaultpass, defaultpass)
			if err != nil {
				return err
			}
			fmt.Println(string(bz))
			return nil
		},
	}
	return cmd
}

// keysAddSommCmd respresents the `keys add` command
func keysAddSommCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add [name]",
		Aliases: []string{"a"},
		Short:   "adds a key to the keychain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}
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

// keysRestoreSommCmd respresents the `keys add` command
func keysRestoreSommCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "restore [name] [mnemonic]",
		Aliases: []string{"r"},
		Short:   "restores a mnemonic to the keychain",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}
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

// keysDeleteSommCmd respresents the `keys delete` command
func keysDeleteSommCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete [name]",
		Aliases: []string{"d"},
		Short:   "deletes a key from the keychain",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}
			if err := ctx.Keyring.Delete(args[0]); err != nil {
				return err
			}

			fmt.Printf("Key %s deleted for good...\n", args[0])
			return nil
		},
	}

	return cmd
}

// keysListSommCmd respresents the `keys list` command
func keysListSommCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Short:   "lists keys from the keychain associated with a particular chain",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}
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

// keysShowSommCmd respresents the `keys show` command
func keysShowSommCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [name]",
		Aliases: []string{"s"},
		Short:   "shows a key from the keychain",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}
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

// keysExportSommCmd respresents the `keys export` command
func keysExportSommCmd() *cobra.Command {
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
