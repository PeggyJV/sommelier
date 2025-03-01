// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: addresses/v1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type MsgAddAddressMapping struct {
	EvmAddress string `protobuf:"bytes,1,opt,name=evm_address,json=evmAddress,proto3" json:"evm_address,omitempty"`
	Signer     string `protobuf:"bytes,2,opt,name=signer,proto3" json:"signer,omitempty"`
}

func (m *MsgAddAddressMapping) Reset()         { *m = MsgAddAddressMapping{} }
func (m *MsgAddAddressMapping) String() string { return proto.CompactTextString(m) }
func (*MsgAddAddressMapping) ProtoMessage()    {}
func (*MsgAddAddressMapping) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbc33d4b2b06ba95, []int{0}
}
func (m *MsgAddAddressMapping) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAddAddressMapping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddAddressMapping.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAddAddressMapping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAddAddressMapping.Merge(m, src)
}
func (m *MsgAddAddressMapping) XXX_Size() int {
	return m.Size()
}
func (m *MsgAddAddressMapping) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAddAddressMapping.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAddAddressMapping proto.InternalMessageInfo

func (m *MsgAddAddressMapping) GetEvmAddress() string {
	if m != nil {
		return m.EvmAddress
	}
	return ""
}

func (m *MsgAddAddressMapping) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

type MsgAddAddressMappingResponse struct {
}

func (m *MsgAddAddressMappingResponse) Reset()         { *m = MsgAddAddressMappingResponse{} }
func (m *MsgAddAddressMappingResponse) String() string { return proto.CompactTextString(m) }
func (*MsgAddAddressMappingResponse) ProtoMessage()    {}
func (*MsgAddAddressMappingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbc33d4b2b06ba95, []int{1}
}
func (m *MsgAddAddressMappingResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAddAddressMappingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddAddressMappingResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAddAddressMappingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAddAddressMappingResponse.Merge(m, src)
}
func (m *MsgAddAddressMappingResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgAddAddressMappingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAddAddressMappingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAddAddressMappingResponse proto.InternalMessageInfo

type MsgRemoveAddressMapping struct {
	Signer string `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
}

func (m *MsgRemoveAddressMapping) Reset()         { *m = MsgRemoveAddressMapping{} }
func (m *MsgRemoveAddressMapping) String() string { return proto.CompactTextString(m) }
func (*MsgRemoveAddressMapping) ProtoMessage()    {}
func (*MsgRemoveAddressMapping) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbc33d4b2b06ba95, []int{2}
}
func (m *MsgRemoveAddressMapping) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRemoveAddressMapping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRemoveAddressMapping.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRemoveAddressMapping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRemoveAddressMapping.Merge(m, src)
}
func (m *MsgRemoveAddressMapping) XXX_Size() int {
	return m.Size()
}
func (m *MsgRemoveAddressMapping) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRemoveAddressMapping.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRemoveAddressMapping proto.InternalMessageInfo

func (m *MsgRemoveAddressMapping) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

type MsgRemoveAddressMappingResponse struct {
}

func (m *MsgRemoveAddressMappingResponse) Reset()         { *m = MsgRemoveAddressMappingResponse{} }
func (m *MsgRemoveAddressMappingResponse) String() string { return proto.CompactTextString(m) }
func (*MsgRemoveAddressMappingResponse) ProtoMessage()    {}
func (*MsgRemoveAddressMappingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dbc33d4b2b06ba95, []int{3}
}
func (m *MsgRemoveAddressMappingResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRemoveAddressMappingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRemoveAddressMappingResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRemoveAddressMappingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRemoveAddressMappingResponse.Merge(m, src)
}
func (m *MsgRemoveAddressMappingResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgRemoveAddressMappingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRemoveAddressMappingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRemoveAddressMappingResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgAddAddressMapping)(nil), "addresses.v1.MsgAddAddressMapping")
	proto.RegisterType((*MsgAddAddressMappingResponse)(nil), "addresses.v1.MsgAddAddressMappingResponse")
	proto.RegisterType((*MsgRemoveAddressMapping)(nil), "addresses.v1.MsgRemoveAddressMapping")
	proto.RegisterType((*MsgRemoveAddressMappingResponse)(nil), "addresses.v1.MsgRemoveAddressMappingResponse")
}

func init() { proto.RegisterFile("addresses/v1/tx.proto", fileDescriptor_dbc33d4b2b06ba95) }

var fileDescriptor_dbc33d4b2b06ba95 = []byte{
	// 313 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4d, 0x4c, 0x49, 0x29,
	0x4a, 0x2d, 0x2e, 0x4e, 0x2d, 0xd6, 0x2f, 0x33, 0xd4, 0x2f, 0xa9, 0xd0, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0xe2, 0x81, 0x0b, 0xeb, 0x95, 0x19, 0x4a, 0x89, 0x27, 0xe7, 0x17, 0xe7, 0xe6, 0x17,
	0xeb, 0xe7, 0x16, 0xa7, 0x83, 0x54, 0xe5, 0x16, 0xa7, 0x43, 0x94, 0x29, 0xc5, 0x70, 0x89, 0xf8,
	0x16, 0xa7, 0x3b, 0xa6, 0xa4, 0x38, 0x42, 0x94, 0xfb, 0x26, 0x16, 0x14, 0x64, 0xe6, 0xa5, 0x0b,
	0xc9, 0x73, 0x71, 0xa7, 0x96, 0xe5, 0xc6, 0x43, 0x0d, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c,
	0xe2, 0x4a, 0x2d, 0xcb, 0x85, 0xaa, 0x13, 0x12, 0xe3, 0x62, 0x2b, 0xce, 0x4c, 0xcf, 0x4b, 0x2d,
	0x92, 0x60, 0x02, 0xcb, 0x41, 0x79, 0x56, 0xdc, 0x4d, 0xcf, 0x37, 0x68, 0x41, 0x39, 0x4a, 0x72,
	0x5c, 0x32, 0xd8, 0x4c, 0x0f, 0x4a, 0x2d, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x55, 0xb2, 0xe3, 0x12,
	0xf7, 0x2d, 0x4e, 0x0f, 0x4a, 0xcd, 0xcd, 0x2f, 0x4b, 0x45, 0x73, 0x00, 0xc2, 0x7c, 0x46, 0xdc,
	0xe6, 0x2b, 0x72, 0xc9, 0xe3, 0xd0, 0x0f, 0xb3, 0xc2, 0xe8, 0x3e, 0x23, 0x17, 0xb3, 0x6f, 0x71,
	0xba, 0x50, 0x32, 0x97, 0x20, 0xa6, 0x2f, 0x95, 0xf4, 0x90, 0x43, 0x49, 0x0f, 0x9b, 0x5b, 0xa5,
	0xb4, 0x08, 0xab, 0x81, 0x59, 0x26, 0x94, 0xc3, 0x25, 0x82, 0xd5, 0x33, 0xaa, 0x18, 0x66, 0x60,
	0x53, 0x26, 0xa5, 0x4b, 0x94, 0x32, 0x98, 0x6d, 0x52, 0xac, 0x0d, 0xcf, 0x37, 0x68, 0x31, 0x3a,
	0x79, 0x9f, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e,
	0xcb, 0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x61, 0x7a, 0x66, 0x49,
	0x46, 0x69, 0x92, 0x5e, 0x72, 0x7e, 0xae, 0x7e, 0x41, 0x6a, 0x7a, 0x7a, 0x65, 0x56, 0x99, 0x7e,
	0x71, 0x7e, 0x6e, 0x6e, 0x6a, 0x4e, 0x66, 0x6a, 0x91, 0x7e, 0x99, 0xa5, 0x7e, 0x85, 0x3e, 0x22,
	0xf1, 0x94, 0x54, 0x16, 0xa4, 0x16, 0x27, 0xb1, 0x81, 0x93, 0x85, 0x31, 0x20, 0x00, 0x00, 0xff,
	0xff, 0x89, 0x6f, 0x52, 0xde, 0x56, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// Adds a mapping between the cosmos address of the sender and the provided EVM address
	AddAddressMapping(ctx context.Context, in *MsgAddAddressMapping, opts ...grpc.CallOption) (*MsgAddAddressMappingResponse, error)
	// Removes the mapping containing the cosmos address of the sender
	RemoveAddressMapping(ctx context.Context, in *MsgRemoveAddressMapping, opts ...grpc.CallOption) (*MsgRemoveAddressMappingResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) AddAddressMapping(ctx context.Context, in *MsgAddAddressMapping, opts ...grpc.CallOption) (*MsgAddAddressMappingResponse, error) {
	out := new(MsgAddAddressMappingResponse)
	err := c.cc.Invoke(ctx, "/addresses.v1.Msg/AddAddressMapping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RemoveAddressMapping(ctx context.Context, in *MsgRemoveAddressMapping, opts ...grpc.CallOption) (*MsgRemoveAddressMappingResponse, error) {
	out := new(MsgRemoveAddressMappingResponse)
	err := c.cc.Invoke(ctx, "/addresses.v1.Msg/RemoveAddressMapping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// Adds a mapping between the cosmos address of the sender and the provided EVM address
	AddAddressMapping(context.Context, *MsgAddAddressMapping) (*MsgAddAddressMappingResponse, error)
	// Removes the mapping containing the cosmos address of the sender
	RemoveAddressMapping(context.Context, *MsgRemoveAddressMapping) (*MsgRemoveAddressMappingResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) AddAddressMapping(ctx context.Context, req *MsgAddAddressMapping) (*MsgAddAddressMappingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAddressMapping not implemented")
}
func (*UnimplementedMsgServer) RemoveAddressMapping(ctx context.Context, req *MsgRemoveAddressMapping) (*MsgRemoveAddressMappingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveAddressMapping not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_AddAddressMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddAddressMapping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddAddressMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/addresses.v1.Msg/AddAddressMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddAddressMapping(ctx, req.(*MsgAddAddressMapping))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemoveAddressMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveAddressMapping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RemoveAddressMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/addresses.v1.Msg/RemoveAddressMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RemoveAddressMapping(ctx, req.(*MsgRemoveAddressMapping))
	}
	return interceptor(ctx, in, info, handler)
}

var Msg_serviceDesc = _Msg_serviceDesc
var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "addresses.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddAddressMapping",
			Handler:    _Msg_AddAddressMapping_Handler,
		},
		{
			MethodName: "RemoveAddressMapping",
			Handler:    _Msg_RemoveAddressMapping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "addresses/v1/tx.proto",
}

func (m *MsgAddAddressMapping) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddAddressMapping) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddAddressMapping) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.EvmAddress) > 0 {
		i -= len(m.EvmAddress)
		copy(dAtA[i:], m.EvmAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.EvmAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgAddAddressMappingResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddAddressMappingResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddAddressMappingResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgRemoveAddressMapping) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRemoveAddressMapping) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRemoveAddressMapping) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgRemoveAddressMappingResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRemoveAddressMappingResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRemoveAddressMappingResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgAddAddressMapping) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.EvmAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgAddAddressMappingResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgRemoveAddressMapping) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgRemoveAddressMappingResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgAddAddressMapping) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgAddAddressMapping: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddAddressMapping: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EvmAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.EvmAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgAddAddressMappingResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgAddAddressMappingResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddAddressMappingResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgRemoveAddressMapping) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgRemoveAddressMapping: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRemoveAddressMapping: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signer", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signer = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgRemoveAddressMappingResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgRemoveAddressMappingResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRemoveAddressMappingResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
