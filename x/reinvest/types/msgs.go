package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgReinvestVote{}
)

// var _ codectypes.UnpackInterfacesMessage = &MsgAllocationCommit{}

const (
	TypeMsgReinvestVote = "reinvest_vote"
)

//////////////////////////
// MsgAllocationPrecommit //
//////////////////////////

// NewMsgReinvestVote return a new MsgAllocationPrecommit
func NewMsgReinvestVote(votes *ReinvestVote, signer sdk.AccAddress, val sdk.ValAddress) (*MsgReinvestVote, error) {
	if signer == nil {
		return nil, fmt.Errorf("no signer provided")
	}

	if len(votes.Cellar) == 0 {
		return nil, fmt.Errorf("need to pass in some votes")
	}

	if err := votes.ValidateBasic(); err != nil {
		return nil, err
	}

	return &MsgReinvestVote{
		ReinvestVotes: votes,
		Signer:        signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgAllocationPrecommit) Route() string { return ModuleName }

// Type implements sdk.Msg
func (m *MsgAllocationPrecommit) Type() string { return TypeMsgReinvestVote }

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
