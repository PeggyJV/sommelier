package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default PoA genesis state, with the module
// DISABLED.
//
// DefaultParams has Enabled=true, which is correct for the v10 upgrade handler
// that turns PoA on alongside a seeded authority set. It is not correct as a
// genesis default: Enabled=true with an empty AuthoritySet makes the first
// EndBlocker see authPower == 0 and enter authority-empty safe mode, where
// MsgUpdateAuthoritySet and MsgUpdateParams are both frozen -- a chain started
// from default genesis would brick itself with no on-chain recovery.
//
// PoA is therefore opt-in: a genesis must name an authority set and enable the
// module together, which Validate enforces.
func DefaultGenesis() *GenesisState {
	params := DefaultParams()
	params.Enabled = false
	return &GenesisState{
		Params:       params,
		AuthoritySet: nil,
	}
}

// Validate performs genesis-state validation.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	// Enabling PoA without an authority set is self-bricking: safe mode
	// activates on the first block and its own recovery messages are frozen
	// there. Reject at import rather than at the first EndBlocker.
	if gs.Params.Enabled && len(gs.AuthoritySet) == 0 {
		return ErrNoBondedAuthority
	}
	seen := make(map[string]struct{}, len(gs.AuthoritySet))
	for _, v := range gs.AuthoritySet {
		if _, err := sdk.ValAddressFromBech32(v.OperatorAddress); err != nil {
			return err
		}
		if _, dup := seen[v.OperatorAddress]; dup {
			return ErrDuplicateAuthority
		}
		seen[v.OperatorAddress] = struct{}{}
	}
	if gs.ActivationHeight < 0 || gs.ActivationTime < 0 || gs.SafeModeThawHeight < 0 {
		return ErrInvalidGenesis
	}
	// Snapshot multipliers are stored as strings and are load-bearing for slash
	// normalisation; a malformed one would be treated as a corrupt entry at
	// slash time and silently skip the slash. Reject at import instead.
	for _, s := range gs.MultiplierSnapshots {
		if s.Height < 0 {
			return ErrInvalidGenesis
		}
		for _, e := range s.Entries {
			if e == nil {
				continue
			}
			if _, err := sdk.ValAddressFromBech32(e.OperatorAddress); err != nil {
				return err
			}
			if _, err := sdk.NewDecFromStr(e.Multiplier); err != nil {
				return err
			}
		}
	}
	return nil
}
