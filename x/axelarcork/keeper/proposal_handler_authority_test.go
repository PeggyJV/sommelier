package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	_ "github.com/peggyjv/sommelier/v10/app/params"
	"github.com/peggyjv/sommelier/v10/x/axelarcork/types"
)

func proposalFixture(t *testing.T, safeMode bool) (Keeper, sdk.Context) {
	t.Helper()

	k, ctx, _, ctrl := setupCorkKeeper(t)
	t.Cleanup(ctrl.Finish)

	k.SetPoaKeeper(stubPoaKeeper{active: safeMode})

	params := types.DefaultParams()
	params.Enabled = true
	params.CorkAuthority = testCorkAuthority
	k.SetParams(ctx, params)

	k.SetChainConfiguration(ctx, testChainArbitrum, types.ChainConfiguration{
		Name:         "arbitrum",
		Id:           testChainArbitrum,
		ProxyAddress: "0x9999999999999999999999999999999999999999",
	})

	return k, ctx
}

// Adding a managed cellar no longer consults x/pubsub. PublisherDomain remains
// on the proposal for wire compatibility but is not validated against an
// approved-publisher list, so an unknown domain must not block the addition.
func TestAxelarAddManagedCellarsNoLongerRequiresPublisher(t *testing.T) {
	k, ctx := proposalFixture(t, false)

	contract := common.HexToAddress("0x5555555555555555555555555555555555555555")
	err := HandleAddManagedCellarsProposal(ctx, k, types.AddAxelarManagedCellarIDsProposal{
		ChainId:         testChainArbitrum,
		CellarIds:       &types.CellarIDSet{ChainId: testChainArbitrum, Ids: []string{contract.String()}},
		PublisherDomain: "not-an-approved-publisher.example",
	})
	require.NoError(t, err)
	require.True(t, k.HasCellarID(ctx, testChainArbitrum, contract))
}

func TestAxelarRemoveManagedCellarsNoLongerTouchesPubsub(t *testing.T) {
	k, ctx := proposalFixture(t, false)

	contract := common.HexToAddress("0x6666666666666666666666666666666666666666")
	k.SetCellarIDs(ctx, testChainArbitrum, types.CellarIDSet{
		ChainId: testChainArbitrum,
		Ids:     []string{contract.String()},
	})
	require.True(t, k.HasCellarID(ctx, testChainArbitrum, contract))

	err := HandleRemoveManagedCellarsProposal(ctx, k, types.RemoveAxelarManagedCellarIDsProposal{
		ChainId:   testChainArbitrum,
		CellarIds: &types.CellarIDSet{ChainId: testChainArbitrum, Ids: []string{contract.String()}},
	})
	require.NoError(t, err)
	require.False(t, k.HasCellarID(ctx, testChainArbitrum, contract))
}

// The safe-mode freeze on staging new call targets must survive the decoupling.
func TestAxelarAddManagedCellarsStillFrozenInSafeMode(t *testing.T) {
	k, ctx := proposalFixture(t, true)

	contract := common.HexToAddress("0x7777777777777777777777777777777777777777")
	err := HandleAddManagedCellarsProposal(ctx, k, types.AddAxelarManagedCellarIDsProposal{
		ChainId:         testChainArbitrum,
		CellarIds:       &types.CellarIDSet{ChainId: testChainArbitrum, Ids: []string{contract.String()}},
		PublisherDomain: "somm.finance",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "safe mode")
	require.False(t, k.HasCellarID(ctx, testChainArbitrum, contract))
}
