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

func createMemoryKey(passphrase string) (mnemonic string, info *keyring.Info, err error) {
	mnemonic, err = createMnemonic()
	if err != nil {
		return "", nil, err
	}

	account, err := createMemoryKeyFromMnemonic(mnemonic, passphrase)
	if err != nil {
		return "", nil, err
	}

	return mnemonic, account, nil
}

func createMemoryKeyFromMnemonic(mnemonic string, passphrase string) (*keyring.Info, error) {
	kb, err := keyring.New("testnet", keyring.BackendMemory, "", nil)
	if err != nil {
		return nil, err
	}

	keyringAlgos, _ := kb.SupportedAlgorithms()
	algo, err := keyring.NewSigningAlgoFromString(string(hd.Secp256k1Type), keyringAlgos)
	if err != nil {
		return nil, err
	}

	account, err := kb.NewAccount("", mnemonic, passphrase, sdk.FullFundraiserPath, algo)
	if err != nil {
		return nil, err
	}

	return &account, nil
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
