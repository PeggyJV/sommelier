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
	return nil
}
