package types

import (
	"math/big"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// Payload represents a gravity logic call
type Payload interface {
	GetEncodedCall() ([]byte, error)
}

const testUniswapLiquidityABI = `[{"inputs":[{"internalType":"address","name":"_uni_router","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"tokenA","type":"address"},{"internalType":"address","name":"tokenB","type":"address"},{"internalType":"uint256","name":"liquidity","type":"uint256"},{"internalType":"uint256","name":"amountAMin","type":"uint256"},{"internalType":"uint256","name":"amountBMin","type":"uint256"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"deadline","type":"uint256"}],"name":"redeemLiquidity","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"token","type":"address"},{"internalType":"uint256","name":"liquidity","type":"uint256"},{"internalType":"uint256","name":"amountTokenMin","type":"uint256"},{"internalType":"uint256","name":"amountETHMin","type":"uint256"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"deadline","type":"uint256"}],"name":"redeemLiquidityETH","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_to","type":"address"},{"internalType":"uint256","name":"_a","type":"uint256"},{"internalType":"uint256","name":"_b","type":"uint256"},{"internalType":"address","name":"state_tokenContract","type":"address"}],"name":"transferTokens","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

func NewRedeemLiquidityETHCall(token string, amount, amountTokenMin, amountEthMin uint64, to string, deadline int64) RedeemLiquidityETH {
	return RedeemLiquidityETH{
		Token:          common.HexToAddress(token),
		Liquidity:      new(big.Int).SetUint64(amount),
		AmountTokenMin: new(big.Int).SetUint64(amountTokenMin),
		AmountETHMin:   new(big.Int).SetUint64(amountEthMin),
		To:             common.HexToAddress(to),
		Deadline:       big.NewInt(deadline),
	}
}

// RedeemLiquidityETH represents an ETH liquidity removal call
type RedeemLiquidityETH struct {
	Token          common.Address
	Liquidity      *big.Int
	AmountTokenMin *big.Int
	AmountETHMin   *big.Int
	To             common.Address
	Deadline       *big.Int
}

// GetEncodedCall gets the checkpoint signature from the given outgoing tx batch
func (b RedeemLiquidityETH) GetEncodedCall() ([]byte, error) {

	abi, err := abi.JSON(strings.NewReader(testUniswapLiquidityABI))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "bad ABI definition in code")
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	return abi.Pack(
		"redeemLiquidityETH",
		b.Token,          // address token,
		b.Liquidity,      // uint256 liquidity,
		b.AmountTokenMin, // uint256 amountTokenMin,
		b.AmountETHMin,   // uint256 amountETHMin,
		b.To,             // address to,
		b.Deadline,       // uint256 deadline
	)
}

const simpleLogicBatch = `[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"_logicContract","type":"address"},{"indexed":false,"internalType":"address","name":"_tokenContract","type":"address"},{"indexed":false,"internalType":"bool","name":"_success","type":"bool"},{"indexed":false,"internalType":"bytes","name":"_returnData","type":"bytes"}],"name":"LogicCallEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[{"internalType":"uint256[]","name":"_amounts","type":"uint256[]"},{"internalType":"bytes[]","name":"_payloads","type":"bytes[]"},{"internalType":"address","name":"_logicContract","type":"address"},{"internalType":"address","name":"_tokenContract","type":"address"}],"name":"logicBatch","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

type SimpleLogicBatch struct {
	Amounts       []*big.Int
	Payloads      []Payload
	LogicContract common.Address
	TokenContract common.Address
}

// GetEncodedCall gets the checkpoint signature from the given outgoing tx batch
func (b SimpleLogicBatch) GetEncodedCall() ([]byte, error) {

	abi, err := abi.JSON(strings.NewReader(simpleLogicBatch))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "bad ABI definition in code")
	}

	var payload = [][]byte{}
	for _, p := range b.Payloads {
		bz, err := p.GetEncodedCall()
		if err != nil {
			return nil, err
		}
		payload = append(payload, bz)
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	return abi.Pack(
		"logicBatch",
		b.Amounts,       // uint256[] memory _amounts,
		payload,         // bytes[] memory _ payloads,
		b.LogicContract, // address _logicContract,
		b.TokenContract, // address _tokenContract
	)
}
