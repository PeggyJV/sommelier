package app

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// stakingNoopEndBlocker wraps staking.AppModule and overrides EndBlock to a
// no-op. The PoA module's EndBlocker invokes staking.EndBlocker directly so
// it can rescale the resulting ValidatorUpdates; SDK 0.47's module.Manager
// panics if two modules return non-empty ValidatorUpdates from EndBlock,
// which would happen if staking's own EndBlocker also ran.
//
// All other AppModule responsibilities (InitGenesis, BeginBlock,
// ExportGenesis, RegisterServices, etc.) are inherited from the embedded
// staking.AppModule unchanged.
type stakingNoopEndBlocker struct {
	staking.AppModule
}

// EndBlock is a no-op. Validator updates are emitted by x/poa.
func (stakingNoopEndBlocker) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}
