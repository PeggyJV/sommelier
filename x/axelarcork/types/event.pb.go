// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: axelar-cork/v1/event.proto

package types

import (
	fmt "fmt"
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

type CorkEvent struct {
	Signer      string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	Validator   string `protobuf:"bytes,2,opt,name=validator,proto3" json:"validator,omitempty"`
	Cork        string `protobuf:"bytes,3,opt,name=cork,proto3" json:"cork,omitempty"`
	BlockHeight uint64 `protobuf:"varint,4,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	ChainId     uint64 `protobuf:"varint,5,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
}

func (m *CorkEvent) Reset()         { *m = CorkEvent{} }
func (m *CorkEvent) String() string { return proto.CompactTextString(m) }
func (*CorkEvent) ProtoMessage()    {}
func (*CorkEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_0b05b5c0e9506f51, []int{0}
}
func (m *CorkEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CorkEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CorkEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CorkEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CorkEvent.Merge(m, src)
}
func (m *CorkEvent) XXX_Size() int {
	return m.Size()
}
func (m *CorkEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_CorkEvent.DiscardUnknown(m)
}

var xxx_messageInfo_CorkEvent proto.InternalMessageInfo

func (m *CorkEvent) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *CorkEvent) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

func (m *CorkEvent) GetCork() string {
	if m != nil {
		return m.Cork
	}
	return ""
}

func (m *CorkEvent) GetBlockHeight() uint64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *CorkEvent) GetChainId() uint64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func init() {
	proto.RegisterType((*CorkEvent)(nil), "axelar_cork.v1.CorkEvent")
}

func init() { proto.RegisterFile("axelar-cork/v1/event.proto", fileDescriptor_0b05b5c0e9506f51) }

var fileDescriptor_0b05b5c0e9506f51 = []byte{
	// 258 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xcd, 0x4a, 0xc3, 0x40,
	0x14, 0x85, 0x33, 0x1a, 0xab, 0x19, 0xc5, 0xc5, 0x20, 0x12, 0x8b, 0x0c, 0xd5, 0x55, 0x37, 0x66,
	0x28, 0xbe, 0x81, 0x22, 0xd8, 0x6d, 0x97, 0x6e, 0x42, 0x7e, 0x86, 0xc9, 0x98, 0x9f, 0x1b, 0x26,
	0xe3, 0xd0, 0xbe, 0x85, 0xe0, 0x4b, 0xb9, 0xec, 0xd2, 0xa5, 0x24, 0x2f, 0x22, 0xb9, 0x11, 0xba,
	0x3b, 0xe7, 0x7c, 0xf7, 0x72, 0x39, 0x97, 0xce, 0x93, 0xad, 0xac, 0x12, 0xf3, 0x90, 0x81, 0x29,
	0x85, 0x5b, 0x09, 0xe9, 0x64, 0x63, 0xa3, 0xd6, 0x80, 0x05, 0x76, 0x39, 0xb1, 0x78, 0x64, 0x91,
	0x5b, 0xcd, 0xaf, 0x14, 0x28, 0x40, 0x24, 0x46, 0x35, 0x4d, 0xdd, 0x7f, 0x11, 0x1a, 0x3c, 0x83,
	0x29, 0x5f, 0xc6, 0x4d, 0x76, 0x4d, 0x67, 0x9d, 0x56, 0x8d, 0x34, 0x21, 0x59, 0x90, 0x65, 0xb0,
	0xf9, 0x77, 0xec, 0x96, 0x06, 0x2e, 0xa9, 0x74, 0x9e, 0x58, 0x30, 0xe1, 0x11, 0xa2, 0x43, 0xc0,
	0x18, 0xf5, 0xc7, 0x23, 0xe1, 0x31, 0x02, 0xd4, 0xec, 0x8e, 0x5e, 0xa4, 0x15, 0x64, 0x65, 0x5c,
	0x48, 0xad, 0x0a, 0x1b, 0xfa, 0x0b, 0xb2, 0xf4, 0x37, 0xe7, 0x98, 0xbd, 0x62, 0xc4, 0x6e, 0xe8,
	0x59, 0x56, 0x24, 0xba, 0x89, 0x75, 0x1e, 0x9e, 0x20, 0x3e, 0x45, 0xbf, 0xce, 0x9f, 0xd6, 0xdf,
	0x3d, 0x27, 0xfb, 0x9e, 0x93, 0xdf, 0x9e, 0x93, 0xcf, 0x81, 0x7b, 0xfb, 0x81, 0x7b, 0x3f, 0x03,
	0xf7, 0xde, 0x84, 0xd2, 0xb6, 0xf8, 0x48, 0xa3, 0x0c, 0x6a, 0xd1, 0x4a, 0xa5, 0x76, 0xef, 0x4e,
	0x74, 0x50, 0xd7, 0xb2, 0xd2, 0xd2, 0x88, 0xad, 0x98, 0x4a, 0xe3, 0x3f, 0xec, 0xae, 0x95, 0x5d,
	0x3a, 0xc3, 0x9e, 0x8f, 0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1f, 0xfb, 0x01, 0xae, 0x2b, 0x01,
	0x00, 0x00,
}

func (m *CorkEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CorkEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CorkEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ChainId != 0 {
		i = encodeVarintEvent(dAtA, i, uint64(m.ChainId))
		i--
		dAtA[i] = 0x28
	}
	if m.BlockHeight != 0 {
		i = encodeVarintEvent(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Cork) > 0 {
		i -= len(m.Cork)
		copy(dAtA[i:], m.Cork)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Cork)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvent(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *CorkEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.Cork)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if m.BlockHeight != 0 {
		n += 1 + sovEvent(uint64(m.BlockHeight))
	}
	if m.ChainId != 0 {
		n += 1 + sovEvent(uint64(m.ChainId))
	}
	return n
}

func sovEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *CorkEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
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
			return fmt.Errorf("proto: CorkEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CorkEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Validator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cork", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
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
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Cork = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainId", wireType)
			}
			m.ChainId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ChainId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
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
func skipEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
					return 0, ErrIntOverflowEvent
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
				return 0, ErrInvalidLengthEvent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvent = fmt.Errorf("proto: unexpected end of group")
)
