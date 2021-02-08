package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/cobra"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"golang.org/x/sync/errgroup"
)

func queryContractABI() *cobra.Command {
	return &cobra.Command{
		Use:     "query-abi",
		Aliases: []string{"abi"},
		Short:   "queries the configed contract address for it's abi",
		RunE: func(cmd *cobra.Command, args []string) error {
			// cl, err := config.NewETHRPCClient()
			// if err != nil {
			// 	return err
			// }
			return nil
		},
	}
}

func startGravityOrchestrator() *cobra.Command {
	return &cobra.Command{
		Use:     "start-gravity-orchestrator",
		Aliases: []string{"orch"},
		Short:   "orchestrates the interaction between eth based chains and gravity zone",
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
				return config.GravityOrchestratorLoop(goctx, cancel, ctx)
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

// GravityCoordinator helps coordinate feeding of the oracle
type GravityCoordinator struct {
	Ctx  client.Context
	Addr sdk.AccAddress
	Val  sdk.AccAddress
	Salt string
	Hash []byte

	height int64

	sync.Mutex
}

// GravityOrchestratorLoop is the listen for events and take action loop for the oracle feeder
func (c Config) GravityOrchestratorLoop(goctx context.Context, cancel context.CancelFunc, ctx client.Context) error {
	var (
		txEventsChan, blEventsChan <-chan ctypes.ResultEvent
		txCancel, blCancel         context.CancelFunc
		coord                      = &GravityCoordinator{Ctx: ctx}
	)
	defer cancel()

	// TODO: reenable sanity checks once we get closer to prod
	// key, err := ctx.Keyring.Key(config.SommKey)
	// if err != nil || key == nil {
	// 	return fmt.Errorf("configured key(%s) not found in keyring", config.SommKey)
	// }

	// coord.Addr = key.GetAddress()

	// // TODO: have this get the validator from the delegate address using the gravity RPCs
	// val, err := GetValFromDel(ctx, coord.Addr)
	// if err != nil || val == nil {
	// 	return fmt.Errorf("configured address(%s) not delegated to a validator", coord.Addr.String())
	// }

	// coord.Val = val

	// TODO: have this make a new GravityRPC client
	cl, err := c.NewSommRPCClient(5 * time.Second)
	if err != nil {
		return err
	}

	if err := cl.Start(); err != nil {
		return err
	}

	// subscribe to tx events on somm
	txEventsChan, txCancel, err = c.SubscribeSomm(goctx, cl, txEvents)
	if err != nil {
		return err
	}
	defer txCancel()

	// subscribe to block events on somm
	blEventsChan, blCancel, err = c.SubscribeSomm(goctx, cl, blEvents)
	if err != nil {
		cancel()
		return err
	}
	defer blCancel()

	contractAddress := common.HexToAddress(c.GravityAddr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// subscribe to eth events
	ethlog, ethsub, err := c.SubscribeETH(goctx, query)
	if err != nil {
		return err
	}
	defer ethsub.Unsubscribe()

	for {
		select {
		case err := <-ethsub.Err():
			// TODO: better eth error handling here?
			return err
		case vlog := <-ethlog:
			if err := coord.handleEthLog(vlog); err != nil {
				return err
			}
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

// SubscribeETH subscribes to events from the configured contract addresss
func (c Config) SubscribeETH(goctx context.Context, query ethereum.FilterQuery) (chan ethtypes.Log, ethereum.Subscription, error) {
	cl, err := c.NewETHWSClient()
	if err != nil {
		return nil, nil, err
	}

	logs := make(chan ethtypes.Log)
	sub, err := cl.SubscribeFilterLogs(goctx, query, logs)
	return logs, sub, err
}

func (c *GravityCoordinator) handleTx(eve ctypes.ResultEvent) error {
	tx, ok := eve.Data.(tmtypes.EventDataTx)
	if !ok {
		return fmt.Errorf("tx event not tx data")
	}
	fmt.Printf("handling tx at height %d\n", tx.Height)
	return nil
}
func (c *GravityCoordinator) handleBlock(eve ctypes.ResultEvent) error {
	bl, ok := eve.Data.(tmtypes.EventDataNewBlock)
	if !ok {
		return fmt.Errorf("block event not block data")
	}
	fmt.Printf("handling new block %d\n", bl.Block.Height)
	return nil
}
func (c *GravityCoordinator) handleEthLog(ethlog ethtypes.Log) error {
	bz, err := ethlog.MarshalJSON()
	if err != nil {
		return err
	}
	fmt.Println("ETH EVENT")
	fmt.Println(string(bz))
	return nil
}
