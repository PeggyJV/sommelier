package types

import (
	"math/big"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/peggyjv/sommelier/x/il/types/contract"
	oracletypes "github.com/peggyjv/sommelier/x/oracle/types"
)

// RedeemLiquidityCall represents an ETH liquidity removal call
type RedeemLiquidityCall struct {
	Pair     oracletypes.UniswapPair
	Stoploss Stoploss
}

// GetCheckpoint gets the checkpoint signature from the given outgoing tx batch
func (b RedeemLiquidityCall) GetCheckpoint(deadline int64) ([]byte, error) {

	abi, err := abi.JSON(strings.NewReader(contract.ContractABI))
	if err != nil {
		return nil, sdkerrors.Wrap(err, "bad ABI definition in code")
	}

	// the methodName needs to be the same as the 'name' above in the checkpointAbiJson
	// but other than that it's a constant that has no impact on the output. This is because
	// it gets encoded as a function name which we must then discard.
	abiEncodedCall, err := abi.Pack(
		"redeemLiquidity",
		common.HexToAddress(b.Pair.Token0.ID),             //	address tokenA
		common.HexToAddress(b.Pair.Token1.ID),             //	address tokenB
		big.NewInt(int64(b.Stoploss.LiquidityPoolShares)), //	uint256 liquidity
		big.NewInt(0), //	uint256 amountAMin
		big.NewInt(0), //	uint256 amountBMin
		common.HexToAddress(b.Stoploss.ReceiverAddress), // address to
		big.NewInt(deadline),                            // uint256 deadline
	)

	// this should never happen outside of test since any case that could crash on encoding
	// should be filtered above.
	if err != nil {
		return nil, sdkerrors.Wrap(err, "packing checkpoint")
	}

	// we hash the resulting encoded bytes discarding the first 4 bytes these 4 bytes are the constant
	// method name 'checkpoint'. If you where to replace the checkpoint constant in this code you would
	// then need to adjust how many bytes you truncate off the front to get the output of abi.encode()
	return crypto.Keccak256Hash(abiEncodedCall[4:]).Bytes(), nil
}
