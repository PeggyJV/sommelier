package app

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/peggyjv/gravity-bridge/module/v6/x/gravity"
)

// gravitySafeModeModule wraps gravity.AppModule and turns BeginBlock/EndBlock
// into no-ops while x/poa is in authority-empty safe mode. Gravity is an
// external dependency whose per-block logic (signer-set tx creation, batch
// creation, inbound attestation observation, and slashing of validators that
// did not sign) would otherwise commit bridge state — or wrongly slash
// validators for confirmations the ante handler is itself blocking — under the
// untrusted, community-only validator set.
//
// All other AppModule responsibilities (InitGenesis, ExportGenesis,
// RegisterServices, simulation, etc.) are inherited from the embedded
// gravity.AppModule unchanged. Skipped per-block work is not lost: gravity's
// EndBlocker is re-entrant, so it resumes from current state after a full block
// begins unfrozen.
type gravitySafeModeModule struct {
	gravity.AppModule
	poa                PoaSafeModeReader
	frozenAtBeginBlock bool
}

func (m *gravitySafeModeModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	m.frozenAtBeginBlock = m.poa != nil && m.poa.SafeModeActive(ctx)
	if m.frozenAtBeginBlock {
		return
	}
	m.AppModule.BeginBlock(ctx, req)
}

func (m *gravitySafeModeModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	// Keep the freeze for the whole block. PoA may clear its persisted flag in
	// an earlier EndBlocker at the thaw height, but validators have not had an
	// unfrozen DeliverTx phase in which to submit the confirmations that gravity
	// is about to evaluate for slashing.
	if m.frozenAtBeginBlock || (m.poa != nil && m.poa.SafeModeActive(ctx)) {
		return nil
	}
	return m.AppModule.EndBlock(ctx, req)
}
