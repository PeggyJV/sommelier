package cmd

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	logTransferSig     = []byte("Transfer(address,address,uint256)")
	logApprovalSig     = []byte("Approval(address,address,uint256)")
	logTransferSigHash = crypto.Keccak256Hash(logTransferSig)
	logApprovalSigHash = crypto.Keccak256Hash(logApprovalSig)
)

func listenERC20Contract() *cobra.Command {
	return &cobra.Command{
		Use:     "listen-erc20 [addr]",
		Aliases: []string{"erc20"},
		Args:    cobra.ExactArgs(1),
		Short:   "listens to events from ERC20 contracts and parses using the abi",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := common.HexToAddress(args[0])
			goctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			eg, goctx := errgroup.WithContext(goctx)
			eg.Go(func() error {
				return config.ERC20ListenLoop(goctx, cancel, addr)
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

// ERC20ListenLoop is the listen for events and take action loop for the oracle feeder
func (c Config) ERC20ListenLoop(goctx context.Context, cancel context.CancelFunc, addr common.Address) error {
	fmt.Printf("Listening for events on endpoint(%s) from contract(%s)\n", c.ETHWSAddr, addr.String())
	defer cancel()

	// subscribe to eth events for the specified erc20 contract
	ethlog, ethsub, err := c.SubscribeETH(goctx, ethereum.FilterQuery{Addresses: []common.Address{addr}})
	if err != nil {
		return err
	}
	defer ethsub.Unsubscribe()

	erc20, err := abi.JSON(strings.NewReader(string(CmdABI)))
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-ethsub.Err():
			// TODO: better eth error handling here?
			return fmt.Errorf("eth subscriptiton error: %w", err)
		case vlog := <-ethlog:
			if err := handleEthLog(vlog, erc20); err != nil {
				return err
			}
		case <-goctx.Done():
			return nil
		}
	}
}

// LogTransfer parses the Transfer event from ERC20 contracts
type LogTransfer struct {
	From   common.Address
	To     common.Address
	Tokens *big.Int
}

// LogApproval parses the Approval event from ERC20 contracts
type LogApproval struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}

func handleEthLog(ethlog ethtypes.Log, erc20abi abi.ABI) error {
	topics0 := ethlog.Topics[0].Hex()
	switch topics0 {
	case logTransferSigHash.Hex():
		var te LogTransfer
		err := erc20abi.UnpackIntoInterface(&te, "Transfer", ethlog.Data)
		if err != nil {
			return err
		}
		te.From = common.HexToAddress(ethlog.Topics[1].Hex())
		te.To = common.HexToAddress(ethlog.Topics[2].Hex())

		fmt.Printf("Transfer: from(%s) to(%s) amount(%s)\n", te.From.Hex(), te.To.Hex(), te.Tokens.String())

	case logApprovalSigHash.Hex():
		var ae LogApproval
		err := erc20abi.UnpackIntoInterface(&ae, "Approval", ethlog.Data)
		if err != nil {
			return err
		}
		ae.TokenOwner = common.HexToAddress(ethlog.Topics[1].Hex())
		ae.Spender = common.HexToAddress(ethlog.Topics[2].Hex())

		fmt.Printf("Approval: from(%s) to(%s) amount(%s)\n", ae.TokenOwner.Hex(), ae.Spender.Hex(), ae.Tokens.String())
	default:
		fmt.Printf("unexpected topic(%s)\n", topics0)
	}
	return nil
}
