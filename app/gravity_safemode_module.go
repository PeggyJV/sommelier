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
// EndBlocker is re-entrant, so it resumes from current state once safe mode
// clears.
type gravitySafeModeModule struct {
	gravity.AppModule
	poa PoaSafeModeReader
}

func (m gravitySafeModeModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	if m.poa != nil && m.poa.SafeModeActive(ctx) {
		return
	}
	m.AppModule.BeginBlock(ctx, req)
}

func (m gravitySafeModeModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	if m.poa != nil && m.poa.SafeModeActive(ctx) {
		return nil
	}
	return m.AppModule.EndBlock(ctx, req)
}
