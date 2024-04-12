package integration_tests

import (
	"fmt"
	"os"

	corktypes "github.com/peggyjv/sommelier/v7/x/cork/types/v2"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	sdkTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/v4/x/gravity/types"
	"github.com/peggyjv/sommelier/v7/app"
	"github.com/peggyjv/sommelier/v7/app/params"
)

const (
	keyringPassphrase = "testpassphrase"
	keyringAppName    = "testnet"
	feeAmount         = 246913560
)

var (
	encodingConfig params.EncodingConfig
	cdc            codec.Codec
)

func init() {
	encodingConfig = app.MakeEncodingConfig()

	encodingConfig.InterfaceRegistry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&stakingtypes.MsgCreateValidator{},
		&stakingtypes.MsgBeginRedelegate{},
		&gravitytypes.MsgDelegateKeys{},
	)
	encodingConfig.InterfaceRegistry.RegisterImplementations(
		(*cryptotypes.PubKey)(nil),
		&secp256k1.PubKey{},
		&ed25519.PubKey{},
	)

	cdc = encodingConfig.Marshaler
}

type chain struct {
	dataDir       string
	id            string
	validators    []*validator
	orchestrators []*orchestrator
	proposer      *proposer
}

func newChain() (*chain, error) {
	var dir string
	var err error
	if _, found := os.LookupEnv("CI"); found {
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	tmpDir, err := os.MkdirTemp(dir, "somm-e2e-testnet")
	if err != nil {
		return nil, err
	}

	return &chain{
		id:      "chain-" + tmrand.NewRand().Str(6),
		dataDir: tmpDir,
	}, nil
}

func (c *chain) configDir() string {
	return fmt.Sprintf("%s/%s", c.dataDir, c.id)
}

func (c *chain) createAndInitValidators(count int) error { //nolint:unused
	for i := 0; i < count; i++ {
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKey("val"); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chain) createAndInitValidatorsWithMnemonics(mnemonics []string) error {
	for i := 0; i < len(mnemonics); i++ {
		// create node
		node := c.createValidator(i)

		// generate genesis files
		if err := node.init(); err != nil {
			return err
		}

		c.validators = append(c.validators, node)

		// create keys
		if err := node.createKeyFromMnemonic("val", mnemonics[i], ""); err != nil {
			return err
		}
		if err := node.createNodeKey(); err != nil {
			return err
		}
		if err := node.createConsensusKey(); err != nil {
			return err
		}
	}

	return nil
}

func (c *chain) createAndInitProposerWithMnemonic(mnemonic string) error {
	hdPath := hd.CreateHDPath(sdk.CoinType, 1, 0)

	// create keys
	record, kb, err := createMemoryKeyFromMnemonic("proposer", mnemonic, "", hdPath)
	if err != nil {
		return err
	}

	proposer := proposer{}

	proposer.keyRecord = *record
	proposer.mnemonic = mnemonic
	proposer.keyring = kb

	c.proposer = &proposer

	return nil
}

func (c *chain) createAndInitOrchestrators(count int) error { //nolint:unused
	mnemonics := make([]string, count)
	for i := 0; i < count; i++ {
		mnemonic, err := createMnemonic()
		if err != nil {
			return err
		}
		mnemonics = append(mnemonics, mnemonic)
	}

	return c.createAndInitOrchestratorsWithMnemonics(mnemonics)
}

func (c *chain) createAndInitOrchestratorsWithMnemonics(mnemonics []string) error {
	hdPath := hd.CreateHDPath(sdk.CoinType, 1, 0)

	for i := 0; i < len(mnemonics); i++ {
		// create orchestrator
		orchestrator := c.createOrchestrator(i)

		// create keys
		record, kb, err := createMemoryKeyFromMnemonic("orch", mnemonics[i], "", hdPath)
		if err != nil {
			return err
		}

		orchestrator.keyRecord = *record
		orchestrator.mnemonic = mnemonics[i]
		orchestrator.keyring = kb

		c.orchestrators = append(c.orchestrators, orchestrator)
	}

	return nil
}

func (c *chain) createValidator(index int) *validator {
	return &validator{
		chain:   c,
		index:   index,
		moniker: "sommelier",
	}
}

func (c *chain) createOrchestrator(index int) *orchestrator {
	return &orchestrator{
		index: index,
	}
}

func (c *chain) clientContext(nodeURI string, kb *keyring.Keyring, fromName string, fromAddr sdk.AccAddress) (*client.Context, error) { //nolint:unparam
	amino := codec.NewLegacyAmino()
	interfaceRegistry := sdkTypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil),
		&stakingtypes.MsgCreateValidator{},
		&gravitytypes.MsgDelegateKeys{},
	)
	interfaceRegistry.RegisterImplementations((*govtypesv1beta1.Content)(nil),
		&corktypes.AddManagedCellarIDsProposal{},
		&corktypes.RemoveManagedCellarIDsProposal{},
		&corktypes.ScheduledCorkProposal{},
	)
	interfaceRegistry.RegisterImplementations((*cryptotypes.PubKey)(nil), &secp256k1.PubKey{}, &ed25519.PubKey{})

	protoCodec := codec.NewProtoCodec(interfaceRegistry)
	txCfg := sdkTx.NewTxConfig(protoCodec, sdkTx.DefaultSignModes)

	encodingConfig := params.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         protoCodec,
		TxConfig:          txCfg,
		Amino:             amino,
	}
	//sims.ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//sims.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	rpcClient, err := rpchttp.New(nodeURI, "/websocket")
	if err != nil {
		return nil, err
	}

	clientContext := client.Context{}.
		WithChainID(c.id).
		WithCodec(protoCodec).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithNodeURI(nodeURI).
		WithClient(rpcClient).
		WithBroadcastMode(flags.BroadcastSync).
		WithKeyring(*kb).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithOutputFormat("json").
		WithFrom(fromName).
		WithFromName(fromName).
		WithFromAddress(fromAddr).
		WithSkipConfirmation(true)

	return &clientContext, nil
}

func (c *chain) sendMsgs(clientCtx client.Context, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	txf := tx.Factory{}.
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithChainID(c.id).
		WithTxConfig(clientCtx.TxConfig).
		WithGasAdjustment(1.2).
		WithKeybase(clientCtx.Keyring).
		WithGas(12345678).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	fromAddr := clientCtx.GetFromAddress()

	if err := txf.AccountRetriever().EnsureExists(clientCtx, fromAddr); err != nil {
		return nil, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, fromAddr)
		if err != nil {
			return nil, err
		}

		if initNum == 0 {
			txf = txf.WithAccountNumber(num)
		}

		if initSeq == 0 {
			txf = txf.WithSequence(seq)
		}
	}

	txf = txf.WithFees("246913560usomm")

	err := tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msgs...)
	if err != nil {
		return nil, err
	}

	resBytes := []byte{}
	_, err = clientCtx.Input.Read(resBytes)
	if err != nil {
		return nil, err
	}

	var res sdk.TxResponse
	err = cdc.Unmarshal(resBytes, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
