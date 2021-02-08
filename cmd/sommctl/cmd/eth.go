package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	logTransferSigHash = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	logApprovalSigHash = crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))
)

func queryContractByteCode() *cobra.Command {
	return &cobra.Command{
		Use:   "query-contract-bytecode [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "query a given contract for its bytecode",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := common.HexToAddress(args[0])
			cl, err := config.NewETHRPCClient()
			if err != nil {
				return err
			}
			bz, err := cl.CodeAt(context.Background(), addr, nil)
			if err != nil {

			}
			fmt.Println(hex.EncodeToString(bz))
			return nil
		},
	}
}

var etherscanAddrs = map[string]string{
	"mainnet": "https://api.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=YourApiKeyToken",
	"goerli":  "https://api-goerli.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=YourApiKeyToken",
	"rinkeby": "https://api-rinkeby.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=YourApiKeyToken",
	"ropsten": "https://api-ropsten.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=YourApiKeyToken",
	"kovan":   "https://api-kovan.etherscan.io/api?module=contract&action=getabi&address=%s&apikey=YourApiKeyToken",
}

func queryERC20Contract() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-contract-abi [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "query a given contract for its abi from etherscan",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}

func queryContractABI() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query-contract-abi [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "query a given contract for its abi from etherscan",
		RunE: func(cmd *cobra.Command, args []string) error {
			net, err := cmd.Flags().GetString("network")
			if err != nil {
				return err
			}
			if val, ok := etherscanAddrs[net]; ok {
				net = val
			} else {
				return fmt.Errorf("%s network not supported, try (mainnet|goerli|rinkeby|ropsten|kovan)", net)
			}

			out, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			// TODO: parse contract address
			res, err := http.Get(fmt.Sprintf(net, args[0]))
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			bz, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			var parse parseAbiResp
			if err = json.Unmarshal(bz, &parse); err != nil {
				return err
			}

			ctabi, err := abi.JSON(strings.NewReader(parse.Result))
			if err != nil {
				return err
			}

			switch out {
			case "abijson":
				fmt.Println(parse.Result)
			case "eventsig":
				for k, v := range ctabi.Events {
					fmt.Printf("%s: %s\n", k, v.Sig)
				}
			case "events":
				for k, v := range ctabi.Events {
					fmt.Printf("%s: %s\n", k, v)
				}
			case "methodsig":
				for k, v := range ctabi.Methods {
					fmt.Printf("%s: %s\n", k, v.Sig)
				}
			case "methods":
				for k, v := range ctabi.Methods {
					fmt.Printf("%s: %s\n", k, v)
				}
			default:
				fmt.Println("invalid output type, printing json...")
				fmt.Println(parse.Result)
			}
			return nil
		},
	}
	cmd.Flags().StringP("network", "n", "mainnet", "network to query (mainnet|goerli|rinkeby|ropsten|kovan)")
	cmd.Flags().StringP("output", "o", "abijson", "output format (abijson|events|methods|eventsig|methodsig)")
	return cmd
}

type parseAbiResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

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

	erc20, err := abi.JSON(strings.NewReader(string(Erc20ABI)))
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

func handleEthLog(lg ethtypes.Log, erc20abi abi.ABI) error {
	fmt.Printf("TX(%s)\n", lg.TxHash.Hex())
	topics0 := lg.Topics[0].Hex()
	switch topics0 {
	case logTransferSigHash.Hex():
		amt, err := erc20abi.Unpack("Transfer", lg.Data)
		if err != nil {
			return err
		}

		from := common.BigToAddress(lg.Topics[1].Big()).Hex()
		to := common.BigToAddress(lg.Topics[2].Big()).Hex()

		fmt.Printf("Transfer: from(%s) to(%s) amount(%v)\n", from, to, amt[0])

	case logApprovalSigHash.Hex():
		amt, err := erc20abi.Unpack("Approval", lg.Data)
		if err != nil {
			return err
		}
		from := common.BigToAddress(lg.Topics[1].Big()).Hex()
		to := common.BigToAddress(lg.Topics[2].Big()).Hex()

		fmt.Printf("Approval: from(%v) to(%v) amount(%v)\n", from, to, amt[0])
	default:
		fmt.Printf("unexpected topic(%s)\n", topics0)
	}
	return nil
}
