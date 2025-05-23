// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cellarfees/v1/cellarfees.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
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

type FeeAccrualCounter struct {
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Count uint64 `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (m *FeeAccrualCounter) Reset()         { *m = FeeAccrualCounter{} }
func (m *FeeAccrualCounter) String() string { return proto.CompactTextString(m) }
func (*FeeAccrualCounter) ProtoMessage()    {}
func (*FeeAccrualCounter) Descriptor() ([]byte, []int) {
	return fileDescriptor_34c89ca12b610c1b, []int{0}
}
func (m *FeeAccrualCounter) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FeeAccrualCounter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FeeAccrualCounter.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FeeAccrualCounter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeeAccrualCounter.Merge(m, src)
}
func (m *FeeAccrualCounter) XXX_Size() int {
	return m.Size()
}
func (m *FeeAccrualCounter) XXX_DiscardUnknown() {
	xxx_messageInfo_FeeAccrualCounter.DiscardUnknown(m)
}

var xxx_messageInfo_FeeAccrualCounter proto.InternalMessageInfo

func (m *FeeAccrualCounter) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *FeeAccrualCounter) GetCount() uint64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type FeeAccrualCounters struct {
	Counters []FeeAccrualCounter `protobuf:"bytes,1,rep,name=counters,proto3" json:"counters"`
}

func (m *FeeAccrualCounters) Reset()         { *m = FeeAccrualCounters{} }
func (m *FeeAccrualCounters) String() string { return proto.CompactTextString(m) }
func (*FeeAccrualCounters) ProtoMessage()    {}
func (*FeeAccrualCounters) Descriptor() ([]byte, []int) {
	return fileDescriptor_34c89ca12b610c1b, []int{1}
}
func (m *FeeAccrualCounters) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FeeAccrualCounters) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FeeAccrualCounters.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FeeAccrualCounters) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeeAccrualCounters.Merge(m, src)
}
func (m *FeeAccrualCounters) XXX_Size() int {
	return m.Size()
}
func (m *FeeAccrualCounters) XXX_DiscardUnknown() {
	xxx_messageInfo_FeeAccrualCounters.DiscardUnknown(m)
}

var xxx_messageInfo_FeeAccrualCounters proto.InternalMessageInfo

func (m *FeeAccrualCounters) GetCounters() []FeeAccrualCounter {
	if m != nil {
		return m.Counters
	}
	return nil
}

func init() {
	proto.RegisterType((*FeeAccrualCounter)(nil), "cellarfees.v1.FeeAccrualCounter")
	proto.RegisterType((*FeeAccrualCounters)(nil), "cellarfees.v1.FeeAccrualCounters")
}

func init() { proto.RegisterFile("cellarfees/v1/cellarfees.proto", fileDescriptor_34c89ca12b610c1b) }

var fileDescriptor_34c89ca12b610c1b = []byte{
	// 233 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0x4e, 0xcd, 0xc9,
	0x49, 0x2c, 0x4a, 0x4b, 0x4d, 0x2d, 0xd6, 0x2f, 0x33, 0xd4, 0x47, 0xf0, 0xf4, 0x0a, 0x8a, 0xf2,
	0x4b, 0xf2, 0x85, 0x78, 0x91, 0x44, 0xca, 0x0c, 0xa5, 0x44, 0xd2, 0xf3, 0xd3, 0xf3, 0xc1, 0x32,
	0xfa, 0x20, 0x16, 0x44, 0x91, 0x92, 0x3d, 0x97, 0xa0, 0x5b, 0x6a, 0xaa, 0x63, 0x72, 0x72, 0x51,
	0x69, 0x62, 0x8e, 0x73, 0x7e, 0x69, 0x5e, 0x49, 0x6a, 0x91, 0x90, 0x08, 0x17, 0x6b, 0x4a, 0x6a,
	0x5e, 0x7e, 0xae, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x84, 0x03, 0x12, 0x4d, 0x06, 0x29,
	0x90, 0x60, 0x52, 0x60, 0xd4, 0x60, 0x09, 0x82, 0x70, 0x94, 0x22, 0xb8, 0x84, 0x30, 0x0c, 0x28,
	0x16, 0x72, 0xe2, 0xe2, 0x48, 0x86, 0xb2, 0x25, 0x18, 0x15, 0x98, 0x35, 0xb8, 0x8d, 0x14, 0xf4,
	0x50, 0x9c, 0xa3, 0x87, 0xa1, 0xc9, 0x89, 0xe5, 0xc4, 0x3d, 0x79, 0x86, 0x20, 0xb8, 0x3e, 0x27,
	0x9f, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63,
	0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18, 0x6e, 0x3c, 0x96, 0x63, 0x88, 0x32, 0x4a, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x2f, 0x48, 0x4d, 0x4f, 0xaf, 0xcc, 0x2a, 0xd3, 0x2f,
	0xce, 0xcf, 0xcd, 0x4d, 0xcd, 0xc9, 0x4c, 0x2d, 0xd2, 0x2f, 0x33, 0xd7, 0xaf, 0x40, 0x0a, 0x0d,
	0xfd, 0x92, 0xca, 0x82, 0xd4, 0xe2, 0x24, 0x36, 0xb0, 0x7f, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x9f, 0x35, 0xe6, 0x97, 0x36, 0x01, 0x00, 0x00,
}

func (m *FeeAccrualCounter) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FeeAccrualCounter) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FeeAccrualCounter) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Count != 0 {
		i = encodeVarintCellarfees(dAtA, i, uint64(m.Count))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintCellarfees(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *FeeAccrualCounters) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FeeAccrualCounters) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FeeAccrualCounters) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Counters) > 0 {
		for iNdEx := len(m.Counters) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Counters[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintCellarfees(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintCellarfees(dAtA []byte, offset int, v uint64) int {
	offset -= sovCellarfees(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *FeeAccrualCounter) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovCellarfees(uint64(l))
	}
	if m.Count != 0 {
		n += 1 + sovCellarfees(uint64(m.Count))
	}
	return n
}

func (m *FeeAccrualCounters) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Counters) > 0 {
		for _, e := range m.Counters {
			l = e.Size()
			n += 1 + l + sovCellarfees(uint64(l))
		}
	}
	return n
}

func sovCellarfees(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCellarfees(x uint64) (n int) {
	return sovCellarfees(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *FeeAccrualCounter) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCellarfees
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
			return fmt.Errorf("proto: FeeAccrualCounter: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FeeAccrualCounter: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCellarfees
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
				return ErrInvalidLengthCellarfees
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCellarfees
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Count", wireType)
			}
			m.Count = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCellarfees
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Count |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipCellarfees(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCellarfees
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
func (m *FeeAccrualCounters) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCellarfees
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
			return fmt.Errorf("proto: FeeAccrualCounters: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FeeAccrualCounters: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Counters", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCellarfees
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
				return ErrInvalidLengthCellarfees
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthCellarfees
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Counters = append(m.Counters, FeeAccrualCounter{})
			if err := m.Counters[len(m.Counters)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCellarfees(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCellarfees
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
func skipCellarfees(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCellarfees
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
					return 0, ErrIntOverflowCellarfees
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
					return 0, ErrIntOverflowCellarfees
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
				return 0, ErrInvalidLengthCellarfees
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCellarfees
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCellarfees
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCellarfees        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCellarfees          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCellarfees = fmt.Errorf("proto: unexpected end of group")
)
