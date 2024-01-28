package v1

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	corktypes "github.com/peggyjv/sommelier/v7/x/cork/types"
)

var (
	_ sdk.Msg = &MsgScheduleCorkRequest{}
)

const (
	TypeMsgSubmitCorkRequest   = "cork_submit"
	TypeMsgScheduleCorkRequest = "cork_schedule"
)

////////////////////////////
// MsgScheduleCorkRequest //
////////////////////////////

// NewMsgScheduleCorkRequest return a new MsgScheduleCorkRequest
func NewMsgScheduleCorkRequest(body []byte, address common.Address, blockHeight uint64, signer sdk.AccAddress) (*MsgScheduleCorkRequest, error) {
	return &MsgScheduleCorkRequest{
		Cork: &Cork{
			EncodedContractCall:   body,
			TargetContractAddress: address.String(),
		},
		BlockHeight: blockHeight,
		Signer:      signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (m *MsgScheduleCorkRequest) Route() string { return corktypes.ModuleName }

// Type implements sdk.Msg
func (m *MsgScheduleCorkRequest) Type() string { return TypeMsgScheduleCorkRequest }

// ValidateBasic implements sdk.Msg
func (m *MsgScheduleCorkRequest) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Signer); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}

	if m.BlockHeight == 0 {
		return fmt.Errorf("block height must be non-zero")
	}

	return m.Cork.ValidateBasic()
}

// GetSignBytes implements sdk.Msg
func (m *MsgScheduleCorkRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgScheduleCorkRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.MustGetSigner()}
}

// MustGetSigner returns the signer address
func (m *MsgScheduleCorkRequest) MustGetSigner() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return addr
}
