package cmd

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	oracle "github.com/peggyjv/sommelier/x/oracle/types"
	"github.com/spf13/cobra"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"golang.org/x/sync/errgroup"
)

var (
	txEvents = "tm.event='Tx'"
	blEvents = "tm.event='NewBlock'"
)

func startOracleFeederCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "start",
		Aliases: []string{"feed"},
		Short:   "feeds the oracle with new uniswap data",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			if err := eg.Wait(); err != nil {
				return err
			}
			return nil
		},
	}
}

// Coordinator helps coordinate feeding of the oracle
type Coordinator struct {
	Ctx  client.Context
	UD   *oracle.UniswapData
	Addr sdk.AccAddress
	Val  sdk.AccAddress
	Salt string
	Hash []byte

	height int64

	sync.Mutex
}

func (c *Coordinator) handleTx(txEvent ctypes.ResultEvent) error {
	tx, ok := txEvent.Data.(tmtypes.EventDataTx)
	if !ok {
		return fmt.Errorf("tx event not tx data")
	}
	for _, ev := range tx.Result.Events {
		if ev.Type == oracle.EventTypeOracleDataPrevote {
			for _, att := range ev.Attributes {
				if string(att.Key) == oracle.AttributeKeyValidator {
					config.log.Debug("prevote detected", "validator", string(att.Value))
					if string(att.Value) == c.Val.String() {
						config.log.Info("submitting oracle data vote", "signer", c.Addr.String(), "height", tx.Height)
						if err := c.SubmitOracleDataVote(); err != nil {
							return err
						}
					}
				}
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
	c.height = bl.Block.Height
	prevote := false
	for _, ev := range bl.ResultBeginBlock.Events {
		if ev.Type == oracle.EventTypeVotePeriod {
			config.log.Info("new vote period beginning", "height", bl.Block.Height)
			prevote = true
		}
	}
	if prevote {
		config.log.Info("submitting oracle data prevote", "signer", c.Addr.String(), "height", c.height)
		if err := c.SubmitOracleDataPrevote(); err != nil {
			return err
		}
	}
	return nil
}

// SubmitOracleDataVote is called to send the vote
func (c *Coordinator) SubmitOracleDataVote() (err error) {
	od, err := oracle.PackOracleData(c.UD)
	if err != nil {
		return
	}
	msg := oracle.NewMsgOracleDataVote([]string{c.Salt}, []*cdctypes.Any{od}, c.Addr)
	if err = msg.ValidateBasic(); err != nil {
		return
	}
	return config.BroadcastTxSomm(c.Ctx, c, msg)
}

// SubmitOracleDataPrevote is called to send the prevote
func (c *Coordinator) SubmitOracleDataPrevote() error {
	ud, err := config.GetPairs(context.Background(), 100, 0)
	if err != nil {
		return err
	}
	c.UD = ud
	c.Salt = genrandstr(6)
	c.Hash = oracle.DataHash(c.Salt, ud.CannonicalJSON(), c.Val)
	return config.BroadcastTxSomm(c.Ctx, c,
		oracle.NewMsgOracleDataPrevote([][]byte{c.Hash}, c.Addr))
}

// BroadcastTxSomm broadcasts a transaction from the oracle
func (c Config) BroadcastTxSomm(ctx client.Context, coord *Coordinator, msgs ...sdk.Msg) error {
	ctx = ctx.WithFromAddress(coord.Addr)
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
func (c Config) OracleFeederLoop(goctx context.Context, cancel context.CancelFunc, ctx client.Context) error {
	var (
		txEventsChan, blEventsChan <-chan ctypes.ResultEvent
		txCancel, blCancel         context.CancelFunc
		coord                      = &Coordinator{Ctx: ctx}
	)
	defer cancel()
	key, err := ctx.Keyring.Key(config.SommKey)
	if err != nil || key == nil {
		return fmt.Errorf("configured key(%s) not found in keyring", config.SommKey)
	}

	coord.Addr = key.GetAddress()

	val, err := GetValFromDel(ctx, coord.Addr)
	if err != nil || val == nil {
		return fmt.Errorf("configured address(%s) not delegated to a validator", coord.Addr.String())
	}

	coord.Val = val

	cl, err := c.NewSommRPCClient(5 * time.Second)
	if err != nil {
		return err
	}

	if err := cl.Start(); err != nil {
		return err
	}

	// subscribe to tx events
	txEventsChan, txCancel, err = c.SubscribeSomm(goctx, cl, txEvents)
	if err != nil {
		return err
	}

	defer txCancel()

	// subscribe to block events
	blEventsChan, blCancel, err = c.SubscribeSomm(goctx, cl, blEvents)
	if err != nil {
		cancel()
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
			return fmt.Errorf("Exiting feeder loop, recieved stop signal(%s)", sig.String())
		case <-ctx.Done():
			return nil
		}
	}
}
