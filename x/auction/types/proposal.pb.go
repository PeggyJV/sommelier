// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: auction/v1/proposal.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

type SetTokenPricesProposal struct {
	Title       string                `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description string                `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	TokenPrices []*ProposedTokenPrice `protobuf:"bytes,3,rep,name=token_prices,json=tokenPrices,proto3" json:"token_prices,omitempty"`
}

func (m *SetTokenPricesProposal) Reset()         { *m = SetTokenPricesProposal{} }
func (m *SetTokenPricesProposal) String() string { return proto.CompactTextString(m) }
func (*SetTokenPricesProposal) ProtoMessage()    {}
func (*SetTokenPricesProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_b838951ecd76ffbc, []int{0}
}
func (m *SetTokenPricesProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SetTokenPricesProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SetTokenPricesProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SetTokenPricesProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetTokenPricesProposal.Merge(m, src)
}
func (m *SetTokenPricesProposal) XXX_Size() int {
	return m.Size()
}
func (m *SetTokenPricesProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_SetTokenPricesProposal.DiscardUnknown(m)
}

var xxx_messageInfo_SetTokenPricesProposal proto.InternalMessageInfo

func (m *SetTokenPricesProposal) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *SetTokenPricesProposal) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *SetTokenPricesProposal) GetTokenPrices() []*ProposedTokenPrice {
	if m != nil {
		return m.TokenPrices
	}
	return nil
}

type SetTokenPricesProposalWithDeposit struct {
	Title       string                `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description string                `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	TokenPrices []*ProposedTokenPrice `protobuf:"bytes,3,rep,name=token_prices,json=tokenPrices,proto3" json:"token_prices,omitempty"`
	Deposit     string                `protobuf:"bytes,4,opt,name=deposit,proto3" json:"deposit,omitempty"`
}

func (m *SetTokenPricesProposalWithDeposit) Reset()         { *m = SetTokenPricesProposalWithDeposit{} }
func (m *SetTokenPricesProposalWithDeposit) String() string { return proto.CompactTextString(m) }
func (*SetTokenPricesProposalWithDeposit) ProtoMessage()    {}
func (*SetTokenPricesProposalWithDeposit) Descriptor() ([]byte, []int) {
	return fileDescriptor_b838951ecd76ffbc, []int{1}
}
func (m *SetTokenPricesProposalWithDeposit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SetTokenPricesProposalWithDeposit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SetTokenPricesProposalWithDeposit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SetTokenPricesProposalWithDeposit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetTokenPricesProposalWithDeposit.Merge(m, src)
}
func (m *SetTokenPricesProposalWithDeposit) XXX_Size() int {
	return m.Size()
}
func (m *SetTokenPricesProposalWithDeposit) XXX_DiscardUnknown() {
	xxx_messageInfo_SetTokenPricesProposalWithDeposit.DiscardUnknown(m)
}

var xxx_messageInfo_SetTokenPricesProposalWithDeposit proto.InternalMessageInfo

func (m *SetTokenPricesProposalWithDeposit) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *SetTokenPricesProposalWithDeposit) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *SetTokenPricesProposalWithDeposit) GetTokenPrices() []*ProposedTokenPrice {
	if m != nil {
		return m.TokenPrices
	}
	return nil
}

func (m *SetTokenPricesProposalWithDeposit) GetDeposit() string {
	if m != nil {
		return m.Deposit
	}
	return ""
}

func init() {
	proto.RegisterType((*SetTokenPricesProposal)(nil), "auction.v1.SetTokenPricesProposal")
	proto.RegisterType((*SetTokenPricesProposalWithDeposit)(nil), "auction.v1.SetTokenPricesProposalWithDeposit")
}

func init() { proto.RegisterFile("auction/v1/proposal.proto", fileDescriptor_b838951ecd76ffbc) }

var fileDescriptor_b838951ecd76ffbc = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x51, 0x3d, 0x4f, 0xc3, 0x30,
	0x10, 0xad, 0x29, 0x1f, 0xc2, 0x65, 0x8a, 0x2a, 0xe4, 0x76, 0xb0, 0x4a, 0xa7, 0x4e, 0xb1, 0x0a,
	0x13, 0x23, 0x88, 0x85, 0xad, 0x2a, 0x48, 0x48, 0x2c, 0x55, 0x9b, 0x9c, 0x52, 0x43, 0xd2, 0xb3,
	0xe2, 0x6b, 0x44, 0xff, 0x05, 0x12, 0x7f, 0x86, 0x9f, 0xc0, 0xd8, 0x91, 0x11, 0x25, 0x7f, 0x04,
	0x25, 0x6e, 0x94, 0x0e, 0xec, 0x6c, 0x7e, 0xf7, 0x7c, 0xef, 0xbd, 0xd3, 0xe3, 0xbd, 0xf9, 0x3a,
	0x20, 0x8d, 0x2b, 0x95, 0x8d, 0x95, 0x49, 0xd1, 0xa0, 0x9d, 0xc7, 0xbe, 0x49, 0x91, 0xd0, 0xe3,
	0x3b, 0xca, 0xcf, 0xc6, 0xfd, 0x5e, 0x80, 0x36, 0x41, 0x3b, 0xab, 0x18, 0xe5, 0x80, 0xfb, 0xd6,
	0x17, 0x7b, 0x0a, 0xf5, 0x86, 0x63, 0xba, 0x11, 0x46, 0xe8, 0x36, 0xca, 0x97, 0x9b, 0x0e, 0x3f,
	0x18, 0x3f, 0x7f, 0x00, 0x7a, 0xc4, 0x57, 0x58, 0x4d, 0x52, 0x1d, 0x80, 0x9d, 0xec, 0x7c, 0xbd,
	0x2e, 0x3f, 0x22, 0x4d, 0x31, 0x08, 0x36, 0x60, 0xa3, 0xd3, 0xa9, 0x03, 0xde, 0x80, 0x77, 0x42,
	0xb0, 0x41, 0xaa, 0x4d, 0xa9, 0x2d, 0x0e, 0x2a, 0x6e, 0x7f, 0xe4, 0xdd, 0xf0, 0x33, 0x2a, 0xe5,
	0x66, 0xa6, 0xd2, 0x13, 0xed, 0x41, 0x7b, 0xd4, 0xb9, 0x94, 0x7e, 0x73, 0x80, 0xef, 0x3c, 0x20,
	0x6c, 0x6c, 0xa7, 0x1d, 0x6a, 0x22, 0x0c, 0x3f, 0x19, 0xbf, 0xf8, 0x3b, 0xd5, 0x93, 0xa6, 0xe5,
	0x1d, 0x18, 0xb4, 0x9a, 0xfe, 0x31, 0xa0, 0x27, 0xf8, 0x49, 0xe8, 0x52, 0x88, 0xc3, 0xca, 0xa0,
	0x86, 0xb7, 0xf7, 0x5f, 0xb9, 0x64, 0xdb, 0x5c, 0xb2, 0x9f, 0x5c, 0xb2, 0xf7, 0x42, 0xb6, 0xb6,
	0x85, 0x6c, 0x7d, 0x17, 0xb2, 0xf5, 0xac, 0x22, 0x4d, 0xcb, 0xf5, 0xc2, 0x0f, 0x30, 0x51, 0x06,
	0xa2, 0x68, 0xf3, 0x92, 0x29, 0x8b, 0x49, 0x02, 0xb1, 0x86, 0x54, 0x65, 0xd7, 0xea, 0xad, 0x6e,
	0x4c, 0xd1, 0xc6, 0x80, 0x5d, 0x1c, 0x57, 0x15, 0x5d, 0xfd, 0x06, 0x00, 0x00, 0xff, 0xff, 0x03,
	0x83, 0x59, 0x3e, 0x16, 0x02, 0x00, 0x00,
}

func (m *SetTokenPricesProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SetTokenPricesProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SetTokenPricesProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TokenPrices) > 0 {
		for iNdEx := len(m.TokenPrices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TokenPrices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintProposal(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SetTokenPricesProposalWithDeposit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SetTokenPricesProposalWithDeposit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SetTokenPricesProposalWithDeposit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Deposit) > 0 {
		i -= len(m.Deposit)
		copy(dAtA[i:], m.Deposit)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Deposit)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.TokenPrices) > 0 {
		for iNdEx := len(m.TokenPrices) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TokenPrices[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintProposal(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintProposal(dAtA []byte, offset int, v uint64) int {
	offset -= sovProposal(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SetTokenPricesProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	if len(m.TokenPrices) > 0 {
		for _, e := range m.TokenPrices {
			l = e.Size()
			n += 1 + l + sovProposal(uint64(l))
		}
	}
	return n
}

func (m *SetTokenPricesProposalWithDeposit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	if len(m.TokenPrices) > 0 {
		for _, e := range m.TokenPrices {
			l = e.Size()
			n += 1 + l + sovProposal(uint64(l))
		}
	}
	l = len(m.Deposit)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	return n
}

func sovProposal(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProposal(x uint64) (n int) {
	return sovProposal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SetTokenPricesProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposal
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
			return fmt.Errorf("proto: SetTokenPricesProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SetTokenPricesProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenPrices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenPrices = append(m.TokenPrices, &ProposedTokenPrice{})
			if err := m.TokenPrices[len(m.TokenPrices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProposal
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
func (m *SetTokenPricesProposalWithDeposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposal
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
			return fmt.Errorf("proto: SetTokenPricesProposalWithDeposit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SetTokenPricesProposalWithDeposit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenPrices", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TokenPrices = append(m.TokenPrices, &ProposedTokenPrice{})
			if err := m.TokenPrices[len(m.TokenPrices)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deposit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Deposit = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProposal
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
func skipProposal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProposal
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
					return 0, ErrIntOverflowProposal
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
					return 0, ErrIntOverflowProposal
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
				return 0, ErrInvalidLengthProposal
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProposal
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProposal
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProposal        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProposal          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProposal = fmt.Errorf("proto: unexpected end of group")
)
