package types

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

var (
	_ sdk.Msg = &MsgDelegateFeedConsent{}
	_ sdk.Msg = &MsgOracleDataPrevote{}
	_ sdk.Msg = &MsgOracleDataVote{}
)

var _ codectypes.UnpackInterfacesMessage = &MsgOracleDataVote{}

const (
	TypeMsgDelegateFeedConsent = "delegate_feed_consent"
	TypeMsgOracleDataPrevote   = "oracle_data_prevote"
	TypeMsgOracleDataVote      = "oracle_data_vote"
)

////////////////////////////
// MsgDelegateFeedConsent //
////////////////////////////

// NewMsgDelegateFeedConsent returns a new MsgDelegateFeedConsent
func NewMsgDelegateFeedConsent(del sdk.AccAddress, val sdk.ValAddress) *MsgDelegateFeedConsent {
	if del == nil || val == nil {
		return nil
	}

	return &MsgDelegateFeedConsent{
		Validator: val.String(),
		Delegate:  del.String(),
	}
}

// Route implements sdk.Msg
func (m *MsgDelegateFeedConsent) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgDelegateFeedConsent) Type() string { return TypeMsgDelegateFeedConsent }

// ValidateBasic implements sdk.Msg
func (m *MsgDelegateFeedConsent) ValidateBasic() error {
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
func (m *MsgDelegateFeedConsent) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgDelegateFeedConsent) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.MustGetValidator())}
}

// MustGetValidator returns the sdk.ValAddress for the validator
func (m *MsgDelegateFeedConsent) MustGetValidator() sdk.ValAddress {
	validatorAddr, err := sdk.ValAddressFromBech32(m.Validator)
	if err != nil {
		panic(err)
	}

	return validatorAddr
}

// MustGetDelegate returns the sdk.AccAddress for the delegate
func (m *MsgDelegateFeedConsent) MustGetDelegate() sdk.AccAddress {
	delegatorAddr, err := sdk.AccAddressFromBech32(m.Delegate)
	if err != nil {
		panic(err)
	}

	return delegatorAddr
}

//////////////////////////
// MsgOracleDataPrevote //
//////////////////////////

// NewMsgOracleDataPrevote return a new MsgOracleDataPrevote
func NewMsgOracleDataPrevote(hashes []tmbytes.HexBytes, signer sdk.AccAddress) *MsgOracleDataPrevote {
	if signer == nil {
		return nil
	}

	return &MsgOracleDataPrevote{
		Prevote: &OraclePrevote{
			Hashes: hashes,
		},
		Signer: signer.String(),
	}
}

// Route implements sdk.Msg
func (m *MsgOracleDataPrevote) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgOracleDataPrevote) Type() string { return TypeMsgOracleDataPrevote }

// ValidateBasic implements sdk.Msg
func (m *MsgOracleDataPrevote) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.Prevote == nil || len(m.Prevote.Hashes) == 0 {
		return fmt.Errorf("empty prevote hashes")
	}

	for i, hash := range m.Prevote.Hashes {
		if len(hash) == 0 {
			return fmt.Errorf("hash at index %d cannot be empty", i)
		}
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgOracleDataPrevote) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgOracleDataPrevote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgOracleDataPrevote) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

///////////////////////
// MsgOracleDataVote //
///////////////////////

// NewMsgOracleDataVote return a new MsgOracleDataPrevote
func NewMsgOracleDataVote(vote *OracleVote, signer sdk.AccAddress) *MsgOracleDataVote {
	if signer == nil {
		return nil
	}

	return &MsgOracleDataVote{
		Vote:   vote,
		Signer: signer.String(),
	}
}

// Route implements sdk.Msg
func (m *MsgOracleDataVote) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgOracleDataVote) Type() string { return TypeMsgOracleDataVote }

// ValidateBasic implements sdk.Msg
func (m *MsgOracleDataVote) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	return m.Vote.Validate()
}

// GetSignBytes implements sdk.Msg
func (m *MsgOracleDataVote) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgOracleDataVote) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{m.MustGetSigner()} }

// MustGetSigner returns the signer address
func (m *MsgOracleDataVote) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m *MsgOracleDataVote) UnpackInterfaces(unpacker codectypes.AnyUnpacker) (err error) {
	return m.Vote.UnpackInterfaces(unpacker)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m *QueryOracleDataResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return unpacker.UnpackAny(m.OracleData, new(OracleData))
}
