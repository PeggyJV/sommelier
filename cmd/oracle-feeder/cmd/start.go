package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	oracle "github.com/peggyjv/sommelier/x/oracle/types"
)

var (
	txEvents = "tm.event='Tx'"
	blEvents = "tm.event='NewBlock'"
)

func startOracleFeederCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "start",
		Aliases: []string{"feed"},
		Args:    cobra.NoArgs,
		Short:   "feeds the oracle with new uniswap data",
		RunE: func(cmd *cobra.Command, _ []string) error {
			config.log.Info("starting oracle feeder")
			ctx, err := config.GetClientContext(cmd)
			if err != nil {
				return err
			}

			goctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			eg, goctx := errgroup.WithContext(goctx)
			if err != nil {
				return err
			}

			eg.Go(func() error {
				return config.OracleFeederLoop(goctx, cancel, ctx)
			})

			eg.Go(func() error {
				return stopLoop(goctx, cancel)
			})

			return eg.Wait()
		},
	}
}

// Coordinator helps coordinate feeding of the oracle
type Coordinator struct {
	clientCtx     client.Context
	uniswapPair   *oracle.UniswapPair
	delegatorAddr sdk.AccAddress
	validatorAddr sdk.AccAddress
	salt          string
	hash          tmbytes.HexBytes

	height int64

	sync.Mutex
}

func (c *Coordinator) handleTx(txEvent ctypes.ResultEvent) error {
	tx, ok := txEvent.Data.(tmtypes.EventDataTx)
	if !ok {
		return fmt.Errorf("tx event not tx data")
	}
	config.log.Debug("transaction detected", "height", tx.Height)
	for _, ev := range tx.Result.Events {
		if ev.Type != oracle.EventTypeOracleDataPrevote {
			continue
		}

		for _, att := range ev.Attributes {
			if string(att.Key) != oracle.AttributeKeyValidator {
				continue
			}

			config.log.Debug("prevote detected", "validator", string(att.Value))
			if string(att.Value) != c.validatorAddr.String() {
				continue
			}

			config.log.Info("submitting oracle data vote", "signer", c.delegatorAddr.String(), "height", tx.Height)
			if err := c.SubmitOracleDataVote(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Coordinator) handleBlock(blockEvent ctypes.ResultEvent) error {
	bl, ok := blockEvent.Data.(tmtypes.EventDataNewBlock)
	if !ok {
		return fmt.Errorf("block event not block data")
	}
	config.log.Debug("new block detected", "height", bl.Block.Height)
	c.height = bl.Block.Height
	prevote := false
	for _, ev := range bl.ResultBeginBlock.Events {
		if ev.Type != oracle.EventTypeVotePeriod {
			continue
		}

		config.log.Info("new vote period beginning", "height", bl.Block.Height)
		prevote = true
		break
	}

	if !prevote {
		return nil
	}

	config.log.Info("submitting oracle data prevote", "signer", c.delegatorAddr.String(), "height", c.height)

	return c.SubmitOracleDataPrevote()
}

// SubmitOracleDataVote is called to send the vote
func (c *Coordinator) SubmitOracleDataVote() error {
	od, err := oracle.PackOracleData(c.uniswapPair)
	if err != nil {
		return err
	}
	msg := oracle.NewMsgOracleDataVote([]string{c.salt}, []*cdctypes.Any{od}, c.delegatorAddr)

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return config.BroadcastTx(c.clientCtx, c, msg)
}

// SubmitOracleDataPrevote is called to send the prevote
func (c *Coordinator) SubmitOracleDataPrevote() error {
	pairs, err := config.GetPairs(context.Background(), 100, 0)
	if err != nil {
		// Failing here on parsing data with decimals
		return err
	}

	// marshal and sort
	jsonBz, err := json.Marshal(pairs)
	if err != nil {
		return err
	}

	jsonBz, err = sdk.SortJSON(jsonBz)
	if err != nil {
		return err
	}

	for _, pair := range ud.OracleData {
		bz, err := pair.MarshalJSON()
		if err != nil {
			return err
		}

		fmt.Println(bz)
	}

	// TODO: fix, consider changing the oracle DataHash to use an array of oracle data instead of the json
	// OR marshal the oracle data slice to json here and sort

	// c.uniswapPair = ud
	// c.salt = genrandstr(6)

	// c.hash = oracle.DataHash(c.salt, ud.CannonicalJSON(), c.validatorAddr)

	msg := oracle.NewMsgOracleDataPrevote([]tmbytes.HexBytes{c.hash}, c.delegatorAddr)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return config.BroadcastTx(c.clientCtx, c, msg)
}

// BroadcastTx broadcasts a transaction from the oracle
func (c Config) BroadcastTx(ctx client.Context, coord *Coordinator, msgs ...sdk.Msg) error {
	ctx = ctx.WithFromAddress(coord.delegatorAddr)
	return tx.BroadcastTx(ctx, c.TxFactory(ctx), msgs...)
}

// TxFactory returns a factory for sending transactions
func (c Config) TxFactory(ctx client.Context) tx.Factory {
	return tx.Factory{}.
		WithTxConfig(ctx.TxConfig).
		WithAccountRetriever(ctx.AccountRetriever).
		WithKeybase(ctx.Keyring).
		WithChainID(config.ChainID).
		WithSimulateAndExecute(true).
		WithGasAdjustment(1.5).
		WithGasPrices(config.GasPrices).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
}

// OracleFeederLoop is the listen for events and take action loop for the oracle feeder
func (c *Config) OracleFeederLoop(goctx context.Context, cancel context.CancelFunc, clientCtx client.Context) error {
	var (
		txEventsChan, blEventsChan <-chan ctypes.ResultEvent
		txCancel, blCancel         context.CancelFunc
	)

	coord := &Coordinator{clientCtx: clientCtx}

	defer cancel()
	key, err := clientCtx.Keyring.Key(config.SigningKey)
	if err != nil || key == nil {
		return fmt.Errorf("configured key %s not found in keyring", config.SigningKey)
	}

	coord.delegatorAddr = key.GetAddress()

	coord.validatorAddr, err = GetValFromDel(clientCtx, coord.delegatorAddr)
	if err != nil || coord.validatorAddr == nil {
		return fmt.Errorf("configured address %s not delegated to a validator: %w", coord.delegatorAddr, err)
	}

	timeout := 5 * time.Second

	rpcClient, err := c.NewRPCClient(timeout)
	if err != nil {
		return err
	}

	if err := rpcClient.Start(); err != nil {
		return err
	}

	// subscribe to tx events
	txEventsChan, txCancel, err = c.Subscribe(goctx, rpcClient, txEvents)
	if err != nil {
		return err
	}

	defer txCancel()

	// subscribe to block events
	blEventsChan, blCancel, err = c.Subscribe(goctx, rpcClient, blEvents)
	if err != nil {
		return err
	}

	defer blCancel()

	for {
		select {
		case txEvent := <-txEventsChan:
			if err = coord.handleTx(txEvent); err != nil {
				return err
			}
		case blockEvent := <-blEventsChan:
			if err = coord.handleBlock(blockEvent); err != nil {
				return err
			}
		case <-goctx.Done():
			return nil
		}
	}
}

func genrandstr(s int) string {
	b := make([]byte, s)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// stopLoop waits for a SIGINT or SIGTERM and returns an error
func stopLoop(ctx context.Context, cancel context.CancelFunc) error {
	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigCh:
			cancel()
			return fmt.Errorf("erxiting feeder loop, received stop signal %s", sig.String())
		case <-ctx.Done():
			return nil
		}
	}
}
