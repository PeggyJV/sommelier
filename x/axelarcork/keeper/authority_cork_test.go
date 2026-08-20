package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/sommelier/v10/x/axelarcork/types"
)

const (
	testChainArbitrum = uint64(42161)
	testChainOptimism = uint64(10)
)

// countAt reports how many authority corks are queued for a chain at a height.
func countAt(k Keeper, ctx sdk.Context, chainID, height uint64) int {
	n := 0
	k.IterateAuthorityAxelarCorksByBlockHeight(ctx, chainID, height,
		func(_ []byte, _ common.Address, _ types.AxelarCork) bool {
			n++
			return false
		})
	return n
}

func TestAuthorityAxelarCorkRoundTrip(t *testing.T) {
	k, ctx, _, ctrl := setupCorkKeeper(t)
	defer ctrl.Finish()

	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	cork := types.AxelarCork{
		TargetContractAddress: contract.String(),
		EncodedContractCall:   []byte{0xde, 0xad, 0xbe, 0xef},
		ChainId:               testChainArbitrum,
	}

	id := k.SetAuthorityAxelarCork(ctx, testChainArbitrum, 100, cork)
	require.NotEmpty(t, id)

	var seenID []byte
	var seenContract common.Address
	var seen []types.AxelarCork
	k.IterateAuthorityAxelarCorksByBlockHeight(ctx, testChainArbitrum, 100,
		func(gotID []byte, gotContract common.Address, c types.AxelarCork) bool {
			seen = append(seen, c)
			seenID, seenContract = gotID, gotContract
			return false
		})
	require.Len(t, seen, 1)
	require.Equal(t, cork.EncodedContractCall, seen[0].EncodedContractCall)

	// The iterator recovers id and contract by slicing the raw store key. Assert
	// them: the EndBlocker feeds both back into DeleteAuthorityAxelarCork, so an
	// offset error that still round-trips the VALUE would leave corks
	// undeletable and re-processing every block.
	require.Equal(t, id, seenID)
	require.Equal(t, contract, seenContract)

	// A different height on the same chain must not see it.
	require.Zero(t, countAt(k, ctx, testChainArbitrum, 101))

	// A different chain at the same height must not see it either.
	require.Zero(t, countAt(k, ctx, testChainOptimism, 100))

	k.DeleteAuthorityAxelarCork(ctx, testChainArbitrum, 100, id, contract)
	require.Zero(t, countAt(k, ctx, testChainArbitrum, 100))
}

func TestAuthorityAxelarCorkChainIsolation(t *testing.T) {
	k, ctx, _, ctrl := setupCorkKeeper(t)
	defer ctrl.Finish()

	contract := common.HexToAddress("0x2222222222222222222222222222222222222222")
	for _, chainID := range []uint64{testChainArbitrum, testChainOptimism} {
		k.SetAuthorityAxelarCork(ctx, chainID, 200, types.AxelarCork{
			TargetContractAddress: contract.String(),
			EncodedContractCall:   []byte{0x01},
			ChainId:               chainID,
		})
	}

	require.Equal(t, 1, countAt(k, ctx, testChainArbitrum, 200))
	require.Equal(t, 1, countAt(k, ctx, testChainOptimism, 200))
}
