package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

var (
	_ sdk.Msg = &MsgDelegateAllocations{}
	_ sdk.Msg = &MsgAllocationPrecommit{}
	_ sdk.Msg = &MsgAllocationCommit{}
)

// var _ codectypes.UnpackInterfacesMessage = &MsgAllocationCommit{}

const (
	TypeMsgDelegateAllocations = "delegate_allocations"
	TypeMsgAllocationPrecommit = "allocation_precommit"
	TypeMsgAllocationCommit    = "allocation_commit"
)

////////////////////////////
// MsgDelegateAllocations //
////////////////////////////

// NewMsgDelegateAllocations returns a new MsgDelegateAllocations
func NewMsgDelegateAllocations(del sdk.AccAddress, val sdk.ValAddress) *MsgDelegateAllocations {
	if del == nil || val == nil {
		return nil
	}

	return &MsgDelegateAllocations{
		Validator: val.String(),
		Delegate:  del.String(),
	}
}

// Route implements sdk.Msg
func (m *MsgDelegateAllocations) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgDelegateAllocations) Type() string { return TypeMsgDelegateAllocations }

// ValidateBasic implements sdk.Msg
func (m *MsgDelegateAllocations) ValidateBasic() error {
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
func (m *MsgDelegateAllocations) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgDelegateAllocations) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.MustGetValidator())}
}

// MustGetValidator returns the sdk.ValAddress for the validator
func (m *MsgDelegateAllocations) MustGetValidator() sdk.ValAddress {
	validatorAddr, err := sdk.ValAddressFromBech32(m.Validator)
	if err != nil {
		panic(err)
	}

	return validatorAddr
}

// MustGetDelegate returns the sdk.AccAddress for the delegate
func (m *MsgDelegateAllocations) MustGetDelegate() sdk.AccAddress {
	delegatorAddr, err := sdk.AccAddressFromBech32(m.Delegate)
	if err != nil {
		panic(err)
	}

	return delegatorAddr
}

//////////////////////////
// MsgAllocationPrecommit //
//////////////////////////

// NewMsgAllocationPrecommit return a new MsgAllocationPrecommit
func NewMsgAllocationPrecommit(hash tmbytes.HexBytes, signer sdk.AccAddress) *MsgAllocationPrecommit {
	if signer == nil {
		return nil
	}

	return &MsgAllocationPrecommit{
		Precommit: []*AllocationPrecommit{{Hash: hash}},
		Signer:    signer.String(),
	}
}

// Route implements sdk.Msg
func (m *MsgAllocationPrecommit) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAllocationPrecommit) Type() string { return TypeMsgAllocationPrecommit }

// ValidateBasic implements sdk.Msg
func (m *MsgAllocationPrecommit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.Precommit == nil {
		return fmt.Errorf("empty precommit")
	}
	for _, pc := range m.Precommit {
		if len(pc.Hash) == 0 {
			return fmt.Errorf("empty precommit hash")
		} else if pc.CellarId == "" {
			return fmt.Errorf("empty precommit cellar id")
		}
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgAllocationPrecommit) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgAllocationPrecommit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgAllocationPrecommit) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

///////////////////////
// MsgAllocationCommit //
///////////////////////

// NewMsgAllocationCommit return a new MsgAllocationPrecommit
func NewMsgAllocationCommit(commits []*Allocation, signer sdk.AccAddress) *MsgAllocationCommit {
	if signer == nil {
		return nil
	}

	return &MsgAllocationCommit{
		Commit: commits,
		Signer: signer.String(),
	}
}

// Route implements sdk.Msg
func (m *MsgAllocationCommit) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAllocationCommit) Type() string { return TypeMsgAllocationCommit }

// ValidateBasic implements sdk.Msg
func (m *MsgAllocationCommit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	for _, commit := range m.Commit {
		if err := commit.ValidateBasic(); err != nil {
			return err
		}
	}

	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgAllocationCommit) GetSignBytes() []byte {
	panic("amino support disabled")
}

// GetSigners implements sdk.Msg
func (m *MsgAllocationCommit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgAllocationCommit) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}

// // UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
// func (m *MsgAllocationCommit) UnpackInterfaces(unpacker codectypes.AnyUnpacker) (err error) {
// 	return m.Vote.UnpackInterfaces(unpacker)
// }

// // UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
// func (m *QueryOracleDataResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
// 	return unpacker.UnpackAny(m.OracleData, new(OracleData))
// }
