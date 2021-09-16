package integration_tests

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/common/hexutil"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func createMnemonic() (string, error) {
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

func createMemoryKeyFromMnemonic(name string, mnemonic string, passphrase string, hdPath *hd.BIP44Params) (*keyring.Info, *keyring.Keyring, error) {
	kb, err := keyring.New("testnet", keyring.BackendMemory, "", nil)
	if err != nil {
		return nil, nil, err
	}

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return nil, nil, err
	}

	var pathString string
	if hdPath == nil {
		pathString = sdk.FullFundraiserPath
	} else {
		pathString = hdPath.String()
	}
	account, err := kb.NewAccount(name, mnemonic, passphrase, pathString, algo)
	if err != nil {
		return nil, nil, err
	}

	return &account, &kb, nil
}

func ethereumKeyFromMnemonic(mnemonic string) (*ethereumKey, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		return nil, err
	}

	path, err := hdwallet.ParseDerivationPath(DerivationPath)
	if err != nil {
		return nil, err
	}

	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := wallet.PrivateKeyBytes(account)
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := wallet.PublicKeyBytes(account)
	if err != nil {
		return nil, err
	}

	return &ethereumKey{
		privateKey: hexutil.Encode(privateKeyBytes),
		publicKey:  hexutil.Encode(publicKeyBytes),
		address:    account.Address.String(),
	}, nil
}
