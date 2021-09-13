package integration_tests

import (
	"fmt"
	"io/ioutil"
	"os"

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
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gravitytypes "github.com/peggyjv/gravity-bridge/module/x/gravity/types"
	"github.com/peggyjv/sommelier/app"
	"github.com/peggyjv/sommelier/app/params"
	"github.com/peggyjv/sommelier/x/allocation/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	keyringPassphrase = "testpassphrase"
	keyringAppName    = "testnet"
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
	sommNodes     []*validator
}

func newChain() (*chain, error) {
	tmpDir, err := ioutil.TempDir("", "somm-e2e-testnet")
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

func (c *chain) createAndInitValidators(count int) error {
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

func (c *chain) createAndInitOrchestrators(count int) error {
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
		info, err := createMemoryKeyFromMnemonic(mnemonics[i], "", hdPath)
		if err != nil {
			return err
		}

		orchestrator.keyInfo = *info
		orchestrator.mnemonic = mnemonics[i]

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

func (c *chain) clientContext(nodeURI string, val validator) (*client.Context, error) {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := sdkTypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil),
		&stakingtypes.MsgCreateValidator{},
		&gravitytypes.MsgDelegateKeys{},
		&types.MsgAllocationCommit{},
		&types.MsgAllocationCommitResponse{},
		&types.MsgAllocationPrecommit{},
		&types.MsgAllocationPrecommitResponse{},
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
	//
	//std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//simapp.ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//simapp.ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//ibc.AppModuleBasic{}.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//transfer.AppModuleBasic{}.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//transfer.AppModuleBasic{}.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	rpcClient, err := rpchttp.New(nodeURI, "/websocket")
	if err != nil {
		return nil, err
	}

	kb, err := keyring.New(keyringAppName, keyring.BackendTest, val.configDir(), nil)
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
		WithBroadcastMode(flags.BroadcastBlock).
		WithKeyring(kb).
		WithOutputFormat("json").
		WithFrom("val").
		WithFromName("val").
		WithFromAddress(val.keyInfo.GetAddress()).
		WithSkipConfirmation(true)

	return &clientContext, nil
}

func (c *chain) sendMsgs(clientCtx client.Context, msgs ...sdk.Msg) (*sdk.TxResponse, error) {
	txf := tx.Factory{}.
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithChainID(c.id).
		WithTxConfig(clientCtx.TxConfig).
		WithGasAdjustment(1.2).
		//WithGasPrices("").
		WithKeybase(clientCtx.Keyring).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	from := clientCtx.GetFromAddress()

	if err := txf.AccountRetriever().EnsureExists(clientCtx, from); err != nil {
		return nil, err
	}

	initNum, initSeq := txf.AccountNumber(), txf.Sequence()
	if initNum == 0 || initSeq == 0 {
		num, seq, err := txf.AccountRetriever().GetAccountNumberSequence(clientCtx, from)
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

	txb, err := txf.BuildUnsignedTx(msgs...)
	if err != nil {
		return nil, err
	}

	err = tx.Sign(txf, "val", txb, false)
	if err != nil {
		return nil, err
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		return nil, err
	}

	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return nil, err
	}

	return res, nil
}
