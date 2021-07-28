package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

var (
	_ sdk.Msg = &MsgDelegateDecisions{}
	_ sdk.Msg = &MsgDecisionPrecommit{}
	_ sdk.Msg = &MsgDecisionCommit{}
)

const (
	TypeMsgDelegateDecisions = "delegate_decisions"
	TypeMsgDecisionPrecommit = "decision_precommit"
	TypeMsgDecisionCommit    = "decision_commit"
)

//////////////////////////
// MsgDelegateDecisions //
//////////////////////////

// NewMsgDelegateFeedConsent returns a new MsgDelegateFeedConsent
func NewMsgDelegateDecisions(del sdk.AccAddress, val sdk.ValAddress) *MsgDelegateDecisions {
	if del == nil || val == nil {
		return nil
	}

	return &MsgDelegateDecisions{
		Validator: val.String(),
		Delegate:  del.String(),
	}
}

// Route implements sdk.Msg
func (m *MsgDelegateDecisions) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgDelegateDecisions) Type() string { return TypeMsgDelegateDecisions }

// ValidateBasic implements sdk.Msg
func (m *MsgDelegateDecisions) ValidateBasic() error {
	validatorAddr, err := sdk.ValAddressFromBech32(m.Validator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	delegatorAddr, err := sdk.AccAddressFromBech32(m.Delegate)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if sdk.AccAddress(validatorAddr).Equals(delegatorAddr) {
		return sdkerrors.Wrap(stakingtypes.ErrBadValidatorAddr, "delegate address cannot match the delegator address")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDelegateDecisions) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgDelegateDecisions) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.MustGetValidator())}
}

// MustGetValidator returns the sdk.ValAddress for the validator
func (m *MsgDelegateDecisions) MustGetValidator() sdk.ValAddress {
	validatorAddr, err := sdk.ValAddressFromBech32(m.Validator)
	if err != nil {
		panic(err)
	}

	return validatorAddr
}

// MustGetDelegate returns the sdk.AccAddress for the delegate
func (m *MsgDelegateDecisions) MustGetDelegate() sdk.AccAddress {
	delegatorAddr, err := sdk.AccAddressFromBech32(m.Delegate)
	if err != nil {
		panic(err)
	}

	return delegatorAddr
}

//////////////////////////
// MsgDecisionPrecommit //
//////////////////////////

// NewMsgDelegateFeedConsent returns a new MsgDelegateFeedConsent
func NewMsgDecisionPrecommit(hash tmbytes.HexBytes, signer sdk.AccAddress) *MsgDecisionPrecommit {
	if signer == nil {
		return nil
	}

	return &MsgDecisionPrecommit{
		Hash:   hash,
		Signer: signer.String(),
	}
}

// Route implements sdk.Msg
func (m *MsgDecisionPrecommit) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgDecisionPrecommit) Type() string { return TypeMsgDecisionPrecommit }

// ValidateBasic implements sdk.Msg
func (m *MsgDecisionPrecommit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if len(m.Hash) == 0 {
		return fmt.Errorf("empty precommit hash")
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDecisionPrecommit) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgDecisionPrecommit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgDecisionPrecommit) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

///////////////////////
// MsgDecisionCommit //
///////////////////////

// NewMsgDecisionCommit return a new MsgDecisionCommit
func NewMsgDecisionCommit(decisions []*Decision, signer sdk.AccAddress) *MsgDecisionCommit {
	if signer == nil {
		return nil
	}

	return &MsgDecisionCommit{
		Decisions: decisions,
		Signer:    signer.String(),
	}
}

// Route implements sdk.Msg
func (m *MsgDecisionCommit) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgDecisionCommit) Type() string { return TypeMsgDecisionCommit }

// ValidateBasic implements sdk.Msg
func (m *MsgDecisionCommit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	for _, d := range m.Decisions {
		if err := d.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDecisionCommit) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgDecisionCommit) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.MustGetSigner()} }

// MustGetSigner returns the signer address
func (m *MsgDecisionCommit) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
