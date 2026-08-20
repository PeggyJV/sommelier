package keeper

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	_ "github.com/peggyjv/sommelier/v10/app/params"
	v2types "github.com/peggyjv/sommelier/v10/x/cork/types/v2"
)

// Adding a managed cellar no longer consults x/pubsub. The publisher-domain
// field remains on the proposal for wire compatibility but is not validated,
// so an unknown domain must not block the addition.
func TestAddManagedCellarsNoLongerRequiresPublisher(t *testing.T) {
	k, ctx, _, ctrl := setupCorkKeeper(t)
	defer ctrl.Finish()

	k.SetPoaKeeper(stubPoaKeeper{active: false})

	contract := common.HexToAddress("0x5555555555555555555555555555555555555555")
	err := HandleAddManagedCellarsProposal(ctx, k, v2types.AddManagedCellarIDsProposal{
		CellarIds:       &v2types.CellarIDSet{Ids: []string{contract.String()}},
		PublisherDomain: "not-an-approved-publisher.example",
	})
	require.NoError(t, err)
	require.True(t, k.HasCellarID(ctx, contract))
}

// Removal likewise no longer touches x/pubsub.
func TestRemoveManagedCellarsNoLongerTouchesPubsub(t *testing.T) {
	k, ctx, _, ctrl := setupCorkKeeper(t)
	defer ctrl.Finish()

	k.SetPoaKeeper(stubPoaKeeper{active: false})

	contract := common.HexToAddress("0x6666666666666666666666666666666666666666")
	k.SetCellarIDs(ctx, v2types.CellarIDSet{Ids: []string{contract.String()}})
	require.True(t, k.HasCellarID(ctx, contract))

	err := HandleRemoveManagedCellarsProposal(ctx, k, v2types.RemoveManagedCellarIDsProposal{
		CellarIds: &v2types.CellarIDSet{Ids: []string{contract.String()}},
	})
	require.NoError(t, err)
	require.False(t, k.HasCellarID(ctx, contract))
}

// The safe-mode freeze on adding call targets must survive the decoupling.
func TestAddManagedCellarsStillFrozenInSafeMode(t *testing.T) {
	k, ctx, _, ctrl := setupCorkKeeper(t)
	defer ctrl.Finish()

	k.SetPoaKeeper(stubPoaKeeper{active: true})

	contract := common.HexToAddress("0x7777777777777777777777777777777777777777")
	err := HandleAddManagedCellarsProposal(ctx, k, v2types.AddManagedCellarIDsProposal{
		CellarIds:       &v2types.CellarIDSet{Ids: []string{contract.String()}},
		PublisherDomain: "somm.finance",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "safe mode")
	require.False(t, k.HasCellarID(ctx, contract))
}
