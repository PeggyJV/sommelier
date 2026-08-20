package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default PoA genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:       DefaultParams(),
		AuthoritySet: nil,
	}
}

// Validate performs genesis-state validation.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
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
