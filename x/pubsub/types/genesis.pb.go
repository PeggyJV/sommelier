// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pubsub/v1/genesis.proto

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

// GenesisState defines the pubsub module's genesis state.
type GenesisState struct {
	Params            Params              `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	Publishers        []*Publisher        `protobuf:"bytes,2,rep,name=publishers,proto3" json:"publishers,omitempty"`
	Subscribers       []*Subscriber       `protobuf:"bytes,3,rep,name=subscribers,proto3" json:"subscribers,omitempty"`
	PublisherIntents  []*PublisherIntent  `protobuf:"bytes,4,rep,name=publisher_intents,json=publisherIntents,proto3" json:"publisher_intents,omitempty"`
	SubscriberIntents []*SubscriberIntent `protobuf:"bytes,5,rep,name=subscriber_intents,json=subscriberIntents,proto3" json:"subscriber_intents,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c358cacca624823, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetPublishers() []*Publisher {
	if m != nil {
		return m.Publishers
	}
	return nil
}

func (m *GenesisState) GetSubscribers() []*Subscriber {
	if m != nil {
		return m.Subscribers
	}
	return nil
}

func (m *GenesisState) GetPublisherIntents() []*PublisherIntent {
	if m != nil {
		return m.PublisherIntents
	}
	return nil
}

func (m *GenesisState) GetSubscriberIntents() []*SubscriberIntent {
	if m != nil {
		return m.SubscriberIntents
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "pubsub.v1.GenesisState")
}

func init() { proto.RegisterFile("pubsub/v1/genesis.proto", fileDescriptor_0c358cacca624823) }

var fileDescriptor_0c358cacca624823 = []byte{
	// 309 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xb1, 0x4f, 0x3a, 0x31,
	0x1c, 0xc5, 0xef, 0x80, 0x1f, 0xc9, 0xaf, 0x38, 0x48, 0x83, 0x4a, 0x30, 0xa9, 0xc4, 0x89, 0xa9,
	0x0d, 0x60, 0xe2, 0xce, 0x82, 0x3a, 0x99, 0x63, 0x73, 0x31, 0x94, 0x34, 0xa5, 0x86, 0xa3, 0xcd,
	0x7d, 0x7b, 0x17, 0xf9, 0x2f, 0xfc, 0x9f, 0x5c, 0x18, 0x19, 0x9d, 0x8c, 0xe1, 0xfe, 0x11, 0x43,
	0xef, 0x3c, 0x4e, 0x73, 0x5b, 0xf3, 0xde, 0xfb, 0xbc, 0xd7, 0xe4, 0x8b, 0x2e, 0x4c, 0xcc, 0x21,
	0xe6, 0x2c, 0x19, 0x32, 0x29, 0xd6, 0x02, 0x14, 0x50, 0x13, 0x69, 0xab, 0xf1, 0xff, 0xcc, 0xa0,
	0xc9, 0xb0, 0xd7, 0x91, 0x5a, 0x6a, 0xa7, 0xb2, 0xc3, 0x2b, 0x0b, 0xf4, 0xce, 0x8f, 0xa4, 0x99,
	0x47, 0xf3, 0x10, 0x2a, 0xf4, 0xac, 0xc2, 0xe9, 0xd7, 0xef, 0x35, 0x74, 0x32, 0xcd, 0x26, 0x66,
	0x76, 0x6e, 0x05, 0x66, 0xa8, 0x99, 0x81, 0x5d, 0xbf, 0xef, 0x0f, 0x5a, 0xa3, 0x36, 0x2d, 0x26,
	0xe9, 0xa3, 0x33, 0x26, 0x8d, 0xed, 0xe7, 0x95, 0x17, 0xe4, 0x31, 0x7c, 0x83, 0x90, 0x89, 0xf9,
	0x4a, 0xc1, 0x52, 0x44, 0xd0, 0xad, 0xf5, 0xeb, 0x83, 0xd6, 0xa8, 0x53, 0x86, 0x7e, 0xcc, 0xa0,
	0x94, 0xc3, 0xb7, 0xa8, 0x05, 0x31, 0x87, 0x45, 0xa4, 0xf8, 0x01, 0xab, 0x3b, 0xec, 0xac, 0x84,
	0xcd, 0x0a, 0x37, 0x28, 0x27, 0xf1, 0x14, 0xb5, 0x8b, 0x9a, 0x67, 0xb5, 0xb6, 0x62, 0x6d, 0xa1,
	0xdb, 0x70, 0x78, 0xaf, 0x6a, 0xf5, 0xde, 0x45, 0x82, 0x53, 0xf3, 0x5b, 0x00, 0xfc, 0x80, 0xf0,
	0xb1, 0xb7, 0x68, 0xfa, 0xe7, 0x9a, 0x2e, 0x2b, 0x3f, 0x92, 0x57, 0xb5, 0xe1, 0x8f, 0x02, 0x93,
	0xbb, 0xed, 0x9e, 0xf8, 0xbb, 0x3d, 0xf1, 0xbf, 0xf6, 0xc4, 0x7f, 0x4b, 0x89, 0xb7, 0x4b, 0x89,
	0xf7, 0x91, 0x12, 0xef, 0x89, 0x4a, 0x65, 0x97, 0x31, 0xa7, 0x0b, 0x1d, 0x32, 0x23, 0xa4, 0xdc,
	0xbc, 0x24, 0x0c, 0x74, 0x18, 0x8a, 0x95, 0x12, 0x11, 0x4b, 0xc6, 0xec, 0x35, 0xbf, 0x07, 0xb3,
	0x1b, 0x23, 0x80, 0x37, 0xdd, 0x59, 0xc6, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xf1, 0xb9, 0x59,
	0xf6, 0x02, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SubscriberIntents) > 0 {
		for iNdEx := len(m.SubscriberIntents) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SubscriberIntents[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.PublisherIntents) > 0 {
		for iNdEx := len(m.PublisherIntents) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PublisherIntents[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Subscribers) > 0 {
		for iNdEx := len(m.Subscribers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Subscribers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Publishers) > 0 {
		for iNdEx := len(m.Publishers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Publishers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.Publishers) > 0 {
		for _, e := range m.Publishers {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Subscribers) > 0 {
		for _, e := range m.Subscribers {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.PublisherIntents) > 0 {
		for _, e := range m.PublisherIntents {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.SubscriberIntents) > 0 {
		for _, e := range m.SubscriberIntents {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Publishers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Publishers = append(m.Publishers, &Publisher{})
			if err := m.Publishers[len(m.Publishers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subscribers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Subscribers = append(m.Subscribers, &Subscriber{})
			if err := m.Subscribers[len(m.Subscribers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PublisherIntents", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PublisherIntents = append(m.PublisherIntents, &PublisherIntent{})
			if err := m.PublisherIntents[len(m.PublisherIntents)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubscriberIntents", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SubscriberIntents = append(m.SubscriberIntents, &SubscriberIntent{})
			if err := m.SubscriberIntents[len(m.SubscriberIntents)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)