package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgUpdateAuthoritySet)(nil)
	_ sdk.Msg = (*MsgUpdateParams)(nil)
)

// GetSigners returns the signers (gov authority).
func (m *MsgUpdateAuthoritySet) GetSigners() []sdk.AccAddress {
	// Panic on an invalid authority rather than returning a zero-value signer.
	// ValidateBasic runs first in the tx pipeline, so this should be
	// unreachable; failing fast prevents an ambiguous signer if it isn't.
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

// ValidateBasic performs static checks on MsgUpdateAuthoritySet.
func (m *MsgUpdateAuthoritySet) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "authority: %s", err)
	}
	if len(m.Validators) == 0 {
		return ErrEmptyAuthoritySet
	}
	seen := make(map[string]struct{}, len(m.Validators))
	for _, v := range m.Validators {
		if _, err := sdk.ValAddressFromBech32(v); err != nil {
			return sdkerrors.Wrapf(errors.ErrInvalidAddress, "validator %s: %s", v, err)
		}
		if _, dup := seen[v]; dup {
			return ErrDuplicateAuthority
		}
		seen[v] = struct{}{}
	}
	return nil
}

// GetSigners returns the signers (gov authority).
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.MustAccAddressFromBech32(m.Authority)}
}

// ValidateBasic performs static checks on MsgUpdateParams.
func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return sdkerrors.Wrapf(errors.ErrInvalidAddress, "authority: %s", err)
	}
	return m.Params.Validate()
}
