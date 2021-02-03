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
				if string(att.Key) == oracle.AttributeKeyValidator &&
					string(att.Value) == c.Val.String() {
					if err := c.SubmitOracleDataVote(); err != nil {
						return err
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
			prevote = true
		}
	}
	if prevote {
		fmt.Println("Submitting Oracle Data Prevote", c.height)
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
	return config.BroadcastTx(c.Ctx, c, msg)
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
	return config.BroadcastTx(c.Ctx, c,
		oracle.NewMsgOracleDataPrevote([][]byte{c.Hash}, c.Addr))
}

// BroadcastTx broadcasts a transaction from the oracle
func (c Config) BroadcastTx(ctx client.Context, coord *Coordinator, msgs ...sdk.Msg) error {
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

// Current Round State
// Last Uniswap Query
// Last Uniswap Query Hash
// Last Uniswap Query Hash Salt
// Channel for listening to chain

// Core loop
//   listen over channel for new blocks and other events
//   update round state every block
//     on round start
//       query uniswap data and store
//       compute hash and send prevote
//     on prevote event recieved
//       send vote

// OracleFeederLoop is the listen for events and take action loop for the oracle feeder
func (c Config) OracleFeederLoop(goctx context.Context, cancel context.CancelFunc, ctx client.Context) error {
	var (
		txEventsChan, blEventsChan <-chan ctypes.ResultEvent
		txCancel, blCancel         context.CancelFunc
		err                        error
		coord                      = &Coordinator{Ctx: ctx}
	)
	defer cancel()
	key, err := ctx.Keyring.Key(config.SigningKey)
	if err != nil || key == nil {
		return fmt.Errorf("configured key(%s) not found in keyring", config.SigningKey)
	}

	coord.Addr = key.GetAddress()

	val, err := GetValFromDel(ctx, coord.Addr)
	if err != nil || val == nil {
		return fmt.Errorf("configured address(%s) not delegated to a validator", coord.Addr.String())
	}

	coord.Val = val

	cl, err := c.NewRPCClient(5 * time.Second)
	if err != nil {
		return err
	}

	if err := cl.Start(); err != nil {
		return err
	}

	txEventsChan, txCancel, err = c.Subscribe(goctx, cl, txEvents)
	if err != nil {
		return err
	}
	defer txCancel()

	blEventsChan, blCancel, err = c.Subscribe(goctx, cl, blEvents)
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
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case sig := <-sigCh:
			close(sigCh)
			cancel()
			return fmt.Errorf("Exiting feeder loop, recieved stop signal(%s)", sig.String())
		case <-ctx.Done():
			close(sigCh)
			return nil
		}
	}
}
