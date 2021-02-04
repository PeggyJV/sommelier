package types

import (
	"fmt"
	"strings"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _, _, _ sdk.Msg = &MsgDelegateFeedConsent{}, &MsgOracleDataPrevote{}, &MsgOracleDataVote{}

const (
	TypeMsgDelegateFeedConsent = "delegate_feed_consent"
	TypeMsgOracleDataPrevote   = "oracle_data_prevote"
	TypeMsgOracleDataVote      = "oracle_data_vote"
)

////////////////////////////
// MsgDelegateFeedConsent //
////////////////////////////

// NewMsgDelegateFeedConsent returns a new MsgDelegateFeedConsent
func NewMsgDelegateFeedConsent(val, del sdk.AccAddress) *MsgDelegateFeedConsent {
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
	if _, err := sdk.AccAddressFromBech32(m.Validator); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if _, err := sdk.AccAddressFromBech32(m.Delegate); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDelegateFeedConsent) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgDelegateFeedConsent) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetValidator()}
}

// MustGetValidator returns the sdk.AccAddress for the validator
func (m *MsgDelegateFeedConsent) MustGetValidator() sdk.AccAddress {
	val, err := sdk.AccAddressFromBech32(m.Validator)
	if err != nil {
		panic(err)
	}
	return val
}

// MustGetDelegate returns the sdk.AccAddress for the delegate
func (m *MsgDelegateFeedConsent) MustGetDelegate() sdk.AccAddress {
	val, err := sdk.AccAddressFromBech32(m.Delegate)
	if err != nil {
		panic(err)
	}
	return val
}

//////////////////////////
// MsgOracleDataPrevote //
//////////////////////////

// NewMsgOracleDataPrevote return a new MsgOracleDataPrevote
func NewMsgOracleDataPrevote(hashes [][]byte, signer sdk.AccAddress) *MsgOracleDataPrevote {
	return &MsgOracleDataPrevote{
		Hashes: hashes,
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

	for i, hash := range m.Hashes {
		if len(hash) == 0 {
			return fmt.Errorf("hash at index %d cannot be empty", i)
		}
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgOracleDataPrevote) GetSignBytes() []byte { panic("amino support disabled") }

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
func NewMsgOracleDataVote(salt []string, data []*cdctypes.Any, signer sdk.AccAddress) *MsgOracleDataVote {
	return &MsgOracleDataVote{
		Salt:       salt,
		OracleData: data,
		Signer:     signer.String(),
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

	for _, a := range m.OracleData {
		od, err := UnpackOracleData(a)
		if err != nil {
			return sdkerrors.Wrap(ErrInvalidOracleData, err.Error())
		}

		if err = od.ValidateBasic(); err != nil {
			return sdkerrors.Wrap(ErrInvalidOracleData, err.Error())
		}
	}

	for i, salt := range m.Salt {
		if strings.TrimSpace(salt) == "" {
			return fmt.Errorf("salt string at index %d cannot be blank", i)
		}
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgOracleDataVote) GetSignBytes() []byte { panic("amino support disabled") }

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
	for _, oda := range m.OracleData {
		var od OracleData
		if err := unpacker.UnpackAny(oda, &od); err != nil {
			return err
		}
	}
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m *QueryOracleDataResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var od OracleData
	return unpacker.UnpackAny(m.OracleData, &od)
}
