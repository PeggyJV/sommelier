// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: oracle/v1/msgs.proto

package types

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// MsgDelegateFeedConsent - struct for delegating oracle voting rights to another address.
// type MsgDelegateFeedConsent struct {
// 	Operator sdk.ValAddress `json:"operator" yaml:"operator"`
// 	Delegate sdk.AccAddress `json:"delegate" yaml:"delegate"`
// }
type MsgDelegateFeedConsent struct {
	Operator string `protobuf:"bytes,1,opt,name=operator,proto3" json:"operator,omitempty"`
	Delegate string `protobuf:"bytes,2,opt,name=delegate,proto3" json:"delegate,omitempty"`
}

func (m *MsgDelegateFeedConsent) Reset()         { *m = MsgDelegateFeedConsent{} }
func (m *MsgDelegateFeedConsent) String() string { return proto.CompactTextString(m) }
func (*MsgDelegateFeedConsent) ProtoMessage()    {}
func (*MsgDelegateFeedConsent) Descriptor() ([]byte, []int) {
	return fileDescriptor_6dda9defd295d067, []int{0}
}
func (m *MsgDelegateFeedConsent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDelegateFeedConsent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDelegateFeedConsent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDelegateFeedConsent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDelegateFeedConsent.Merge(m, src)
}
func (m *MsgDelegateFeedConsent) XXX_Size() int {
	return m.Size()
}
func (m *MsgDelegateFeedConsent) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDelegateFeedConsent.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDelegateFeedConsent proto.InternalMessageInfo

func (m *MsgDelegateFeedConsent) GetOperator() string {
	if m != nil {
		return m.Operator
	}
	return ""
}

func (m *MsgDelegateFeedConsent) GetDelegate() string {
	if m != nil {
		return m.Delegate
	}
	return ""
}

// MsgAggregateExchangeRatePrevote - struct for aggregate prevoting on the ExchangeRateVote.
// The purpose of aggregate prevote is to hide vote exchange rates with hash
// which is formatted as hex string in SHA256("{salt}:{exchange rate}{denom},...,{exchange rate}{denom}:{voter}")
// type MsgAggregateExchangeRatePrevote struct {
// 	Hash      AggregateVoteHash `json:"hash" yaml:"hash"`
// 	Feeder    sdk.AccAddress    `json:"feeder" yaml:"feeder"`
// 	Validator sdk.ValAddress    `json:"validator" yaml:"validator"`
// }
type MsgAggregateExchangeRatePrevote struct {
	Hash      []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Feeder    string `protobuf:"bytes,2,opt,name=feeder,proto3" json:"feeder,omitempty"`
	Validator string `protobuf:"bytes,3,opt,name=validator,proto3" json:"validator,omitempty"`
}

func (m *MsgAggregateExchangeRatePrevote) Reset()         { *m = MsgAggregateExchangeRatePrevote{} }
func (m *MsgAggregateExchangeRatePrevote) String() string { return proto.CompactTextString(m) }
func (*MsgAggregateExchangeRatePrevote) ProtoMessage()    {}
func (*MsgAggregateExchangeRatePrevote) Descriptor() ([]byte, []int) {
	return fileDescriptor_6dda9defd295d067, []int{1}
}
func (m *MsgAggregateExchangeRatePrevote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAggregateExchangeRatePrevote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAggregateExchangeRatePrevote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAggregateExchangeRatePrevote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAggregateExchangeRatePrevote.Merge(m, src)
}
func (m *MsgAggregateExchangeRatePrevote) XXX_Size() int {
	return m.Size()
}
func (m *MsgAggregateExchangeRatePrevote) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAggregateExchangeRatePrevote.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAggregateExchangeRatePrevote proto.InternalMessageInfo

func (m *MsgAggregateExchangeRatePrevote) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *MsgAggregateExchangeRatePrevote) GetFeeder() string {
	if m != nil {
		return m.Feeder
	}
	return ""
}

func (m *MsgAggregateExchangeRatePrevote) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

// MsgAggregateExchangeRateVote - struct for voting on the exchange rates of Luna denominated in various Terra assets.
// type MsgAggregateExchangeRateVote struct {
// 	Salt          string         `json:"salt" yaml:"salt"`
// 	ExchangeRates string         `json:"exchange_rates" yaml:"exchange_rates"` // comma separated dec coins
// 	Feeder        sdk.AccAddress `json:"feeder" yaml:"feeder"`
// 	Validator     sdk.ValAddress `json:"validator" yaml:"validator"`
// }
type MsgAggregateExchangeRateVote struct {
	Salt          string `protobuf:"bytes,1,opt,name=salt,proto3" json:"salt,omitempty"`
	ExchangeRates string `protobuf:"bytes,2,opt,name=exchange_rates,json=exchangeRates,proto3" json:"exchange_rates,omitempty"`
	Feeder        string `protobuf:"bytes,3,opt,name=feeder,proto3" json:"feeder,omitempty"`
	Validator     string `protobuf:"bytes,4,opt,name=validator,proto3" json:"validator,omitempty"`
}

func (m *MsgAggregateExchangeRateVote) Reset()         { *m = MsgAggregateExchangeRateVote{} }
func (m *MsgAggregateExchangeRateVote) String() string { return proto.CompactTextString(m) }
func (*MsgAggregateExchangeRateVote) ProtoMessage()    {}
func (*MsgAggregateExchangeRateVote) Descriptor() ([]byte, []int) {
	return fileDescriptor_6dda9defd295d067, []int{2}
}
func (m *MsgAggregateExchangeRateVote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAggregateExchangeRateVote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAggregateExchangeRateVote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAggregateExchangeRateVote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAggregateExchangeRateVote.Merge(m, src)
}
func (m *MsgAggregateExchangeRateVote) XXX_Size() int {
	return m.Size()
}
func (m *MsgAggregateExchangeRateVote) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAggregateExchangeRateVote.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAggregateExchangeRateVote proto.InternalMessageInfo

func (m *MsgAggregateExchangeRateVote) GetSalt() string {
	if m != nil {
		return m.Salt
	}
	return ""
}

func (m *MsgAggregateExchangeRateVote) GetExchangeRates() string {
	if m != nil {
		return m.ExchangeRates
	}
	return ""
}

func (m *MsgAggregateExchangeRateVote) GetFeeder() string {
	if m != nil {
		return m.Feeder
	}
	return ""
}

func (m *MsgAggregateExchangeRateVote) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

func init() {
	proto.RegisterType((*MsgDelegateFeedConsent)(nil), "oracle.v1.MsgDelegateFeedConsent")
	proto.RegisterType((*MsgAggregateExchangeRatePrevote)(nil), "oracle.v1.MsgAggregateExchangeRatePrevote")
	proto.RegisterType((*MsgAggregateExchangeRateVote)(nil), "oracle.v1.MsgAggregateExchangeRateVote")
}

func init() { proto.RegisterFile("oracle/v1/msgs.proto", fileDescriptor_6dda9defd295d067) }

var fileDescriptor_6dda9defd295d067 = []byte{
	// 299 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xcf, 0x4a, 0x03, 0x31,
	0x10, 0xc6, 0x1b, 0x5b, 0x8a, 0x0d, 0xea, 0x21, 0x48, 0x29, 0x52, 0xa2, 0x14, 0x04, 0x41, 0xd8,
	0x50, 0x7c, 0x02, 0xff, 0xd4, 0x5b, 0xa1, 0xec, 0xc1, 0x83, 0x17, 0x49, 0x77, 0xc7, 0xec, 0xea,
	0xee, 0x66, 0xc9, 0xc4, 0xa5, 0x7d, 0x0a, 0x7d, 0x2c, 0x8f, 0x3d, 0x7a, 0x94, 0xdd, 0x17, 0x11,
	0xd3, 0xa8, 0x45, 0xe8, 0x6d, 0xbe, 0xf9, 0x32, 0xf3, 0xfd, 0xc2, 0xd0, 0x43, 0x6d, 0x64, 0x94,
	0x81, 0xa8, 0xc6, 0x22, 0x47, 0x85, 0x41, 0x69, 0xb4, 0xd5, 0xac, 0xb7, 0xee, 0x06, 0xd5, 0x78,
	0x34, 0xa3, 0xfd, 0x29, 0xaa, 0x1b, 0xc8, 0x40, 0x49, 0x0b, 0xb7, 0x00, 0xf1, 0xb5, 0x2e, 0x10,
	0x0a, 0xcb, 0x8e, 0xe8, 0xae, 0x2e, 0xc1, 0x48, 0xab, 0xcd, 0x80, 0x9c, 0x90, 0xb3, 0x5e, 0xf8,
	0xab, 0xbf, 0xbd, 0xd8, 0x8f, 0x0c, 0x76, 0xd6, 0xde, 0x8f, 0x1e, 0x3d, 0xd3, 0xe3, 0x29, 0xaa,
	0x4b, 0xa5, 0x8c, 0xd3, 0x93, 0x45, 0x94, 0xc8, 0x42, 0x41, 0x28, 0x2d, 0xcc, 0x0c, 0x54, 0xda,
	0x02, 0x63, 0xb4, 0x93, 0x48, 0x4c, 0xdc, 0xda, 0xbd, 0xd0, 0xd5, 0xac, 0x4f, 0xbb, 0x8f, 0x00,
	0x31, 0x18, 0xbf, 0xd0, 0x2b, 0x36, 0xa4, 0xbd, 0x4a, 0x66, 0x69, 0xec, 0x38, 0xda, 0xce, 0xfa,
	0x6b, 0x8c, 0x5e, 0x09, 0x1d, 0x6e, 0x4b, 0xbb, 0xf3, 0x51, 0x28, 0x33, 0xeb, 0x7f, 0xe0, 0x6a,
	0x76, 0x4a, 0x0f, 0xc0, 0xbf, 0x7b, 0x30, 0xd2, 0x02, 0xfa, 0xc8, 0x7d, 0xd8, 0x98, 0xc6, 0x0d,
	0xa2, 0xf6, 0x76, 0xa2, 0xce, 0x3f, 0xa2, 0xab, 0xc9, 0x7b, 0xcd, 0xc9, 0xaa, 0xe6, 0xe4, 0xb3,
	0xe6, 0xe4, 0xad, 0xe1, 0xad, 0x55, 0xc3, 0x5b, 0x1f, 0x0d, 0x6f, 0xdd, 0x9f, 0xab, 0xd4, 0x26,
	0x2f, 0xf3, 0x20, 0xd2, 0xb9, 0x28, 0x41, 0xa9, 0xe5, 0x53, 0x25, 0x50, 0xe7, 0x39, 0x64, 0x29,
	0x18, 0xb1, 0x10, 0xfe, 0x54, 0x76, 0x59, 0x02, 0xce, 0xbb, 0xee, 0x52, 0x17, 0x5f, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x30, 0xca, 0x62, 0x8f, 0xc1, 0x01, 0x00, 0x00,
}

func (m *MsgDelegateFeedConsent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDelegateFeedConsent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDelegateFeedConsent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Delegate) > 0 {
		i -= len(m.Delegate)
		copy(dAtA[i:], m.Delegate)
		i = encodeVarintMsgs(dAtA, i, uint64(len(m.Delegate)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Operator) > 0 {
		i -= len(m.Operator)
		copy(dAtA[i:], m.Operator)
		i = encodeVarintMsgs(dAtA, i, uint64(len(m.Operator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgAggregateExchangeRatePrevote) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAggregateExchangeRatePrevote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAggregateExchangeRatePrevote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintMsgs(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Feeder) > 0 {
		i -= len(m.Feeder)
		copy(dAtA[i:], m.Feeder)
		i = encodeVarintMsgs(dAtA, i, uint64(len(m.Feeder)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintMsgs(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgAggregateExchangeRateVote) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAggregateExchangeRateVote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAggregateExchangeRateVote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintMsgs(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Feeder) > 0 {
		i -= len(m.Feeder)
		copy(dAtA[i:], m.Feeder)
		i = encodeVarintMsgs(dAtA, i, uint64(len(m.Feeder)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ExchangeRates) > 0 {
		i -= len(m.ExchangeRates)
		copy(dAtA[i:], m.ExchangeRates)
		i = encodeVarintMsgs(dAtA, i, uint64(len(m.ExchangeRates)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Salt) > 0 {
		i -= len(m.Salt)
		copy(dAtA[i:], m.Salt)
		i = encodeVarintMsgs(dAtA, i, uint64(len(m.Salt)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMsgs(dAtA []byte, offset int, v uint64) int {
	offset -= sovMsgs(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgDelegateFeedConsent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Operator)
	if l > 0 {
		n += 1 + l + sovMsgs(uint64(l))
	}
	l = len(m.Delegate)
	if l > 0 {
		n += 1 + l + sovMsgs(uint64(l))
	}
	return n
}

func (m *MsgAggregateExchangeRatePrevote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovMsgs(uint64(l))
	}
	l = len(m.Feeder)
	if l > 0 {
		n += 1 + l + sovMsgs(uint64(l))
	}
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovMsgs(uint64(l))
	}
	return n
}

func (m *MsgAggregateExchangeRateVote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Salt)
	if l > 0 {
		n += 1 + l + sovMsgs(uint64(l))
	}
	l = len(m.ExchangeRates)
	if l > 0 {
		n += 1 + l + sovMsgs(uint64(l))
	}
	l = len(m.Feeder)
	if l > 0 {
		n += 1 + l + sovMsgs(uint64(l))
	}
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovMsgs(uint64(l))
	}
	return n
}

func sovMsgs(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMsgs(x uint64) (n int) {
	return sovMsgs(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgDelegateFeedConsent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsgs
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgDelegateFeedConsent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDelegateFeedConsent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Operator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgs
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Operator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Delegate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgs
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Delegate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsgs(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMsgs
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMsgs
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgAggregateExchangeRatePrevote) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsgs
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgAggregateExchangeRatePrevote: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAggregateExchangeRatePrevote: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsgs
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append(m.Hash[:0], dAtA[iNdEx:postIndex]...)
			if m.Hash == nil {
				m.Hash = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Feeder", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgs
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Feeder = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgs
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsgs(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMsgs
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMsgs
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *MsgAggregateExchangeRateVote) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsgs
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MsgAggregateExchangeRateVote: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAggregateExchangeRateVote: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Salt", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgs
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Salt = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExchangeRates", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgs
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExchangeRates = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Feeder", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgs
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Feeder = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgs
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthMsgs
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgs
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsgs(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMsgs
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMsgs
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMsgs(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMsgs
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMsgs
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMsgs
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMsgs
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMsgs
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMsgs
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMsgs        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMsgs          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMsgs = fmt.Errorf("proto: unexpected end of group")
)
