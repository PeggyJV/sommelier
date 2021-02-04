package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	proto "github.com/gogo/protobuf/proto"
)

// RegisterInterfaces registers the oracle proto files
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgDelegateFeedConsent{},
		&MsgOracleDataPrevote{},
		&MsgOracleDataVote{},
	)
	registry.RegisterInterface(
		"oracle.v1.OracleData",
		(*OracleData)(nil),
		&UniswapData{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// PackOracleData constructs a new Any packed with the given OracleData value. It returns
// an error if the oracle data can't be casted to a protobuf message or if the concrete
// implemention is not registered to the protobuf codec.
func PackOracleData(oracleData OracleData) (*codectypes.Any, error) {
	msg, ok := oracleData.(proto.Message)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", oracleData)
	}
	anyoracleData, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}
	return anyoracleData, nil
}

// UnpackOracleData Unpacks an Any into a OracleData. It returns an error if the
// client state can't be Unpacked into a OracleData.
func UnpackOracleData(any *codectypes.Any) (OracleData, error) {
	if any == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}
	oracleData, ok := any.GetCachedValue().(OracleData)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot Unpack Any into OracleData %T", any)
	}
	return oracleData, nil
}
