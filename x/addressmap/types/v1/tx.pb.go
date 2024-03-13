// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: addressmap/v1/tx.proto

package v1

import (
	context "context"
	fmt "fmt"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
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

type MsgAddAddressMappingRequest struct {
	EvmAddress string `protobuf:"bytes,1,opt,name=evm_address,json=evmAddress,proto3" json:"evm_address,omitempty"`
}

func (m *MsgAddAddressMappingRequest) Reset()         { *m = MsgAddAddressMappingRequest{} }
func (m *MsgAddAddressMappingRequest) String() string { return proto.CompactTextString(m) }
func (*MsgAddAddressMappingRequest) ProtoMessage()    {}
func (*MsgAddAddressMappingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2212ffde55545bd7, []int{0}
}
func (m *MsgAddAddressMappingRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAddAddressMappingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddAddressMappingRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAddAddressMappingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAddAddressMappingRequest.Merge(m, src)
}
func (m *MsgAddAddressMappingRequest) XXX_Size() int {
	return m.Size()
}
func (m *MsgAddAddressMappingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAddAddressMappingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAddAddressMappingRequest proto.InternalMessageInfo

func (m *MsgAddAddressMappingRequest) GetEvmAddress() string {
	if m != nil {
		return m.EvmAddress
	}
	return ""
}

type MsgAddAddressMappingResponse struct {
}

func (m *MsgAddAddressMappingResponse) Reset()         { *m = MsgAddAddressMappingResponse{} }
func (m *MsgAddAddressMappingResponse) String() string { return proto.CompactTextString(m) }
func (*MsgAddAddressMappingResponse) ProtoMessage()    {}
func (*MsgAddAddressMappingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2212ffde55545bd7, []int{1}
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

type MsgRemoveAddressMappingRequest struct {
}

func (m *MsgRemoveAddressMappingRequest) Reset()         { *m = MsgRemoveAddressMappingRequest{} }
func (m *MsgRemoveAddressMappingRequest) String() string { return proto.CompactTextString(m) }
func (*MsgRemoveAddressMappingRequest) ProtoMessage()    {}
func (*MsgRemoveAddressMappingRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2212ffde55545bd7, []int{2}
}
func (m *MsgRemoveAddressMappingRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgRemoveAddressMappingRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgRemoveAddressMappingRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgRemoveAddressMappingRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgRemoveAddressMappingRequest.Merge(m, src)
}
func (m *MsgRemoveAddressMappingRequest) XXX_Size() int {
	return m.Size()
}
func (m *MsgRemoveAddressMappingRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgRemoveAddressMappingRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MsgRemoveAddressMappingRequest proto.InternalMessageInfo

type MsgRemoveAddressMappingResponse struct {
}

func (m *MsgRemoveAddressMappingResponse) Reset()         { *m = MsgRemoveAddressMappingResponse{} }
func (m *MsgRemoveAddressMappingResponse) String() string { return proto.CompactTextString(m) }
func (*MsgRemoveAddressMappingResponse) ProtoMessage()    {}
func (*MsgRemoveAddressMappingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2212ffde55545bd7, []int{3}
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
	proto.RegisterType((*MsgAddAddressMappingRequest)(nil), "addressmap.v1.MsgAddAddressMappingRequest")
	proto.RegisterType((*MsgAddAddressMappingResponse)(nil), "addressmap.v1.MsgAddAddressMappingResponse")
	proto.RegisterType((*MsgRemoveAddressMappingRequest)(nil), "addressmap.v1.MsgRemoveAddressMappingRequest")
	proto.RegisterType((*MsgRemoveAddressMappingResponse)(nil), "addressmap.v1.MsgRemoveAddressMappingResponse")
}

func init() { proto.RegisterFile("addressmap/v1/tx.proto", fileDescriptor_2212ffde55545bd7) }

var fileDescriptor_2212ffde55545bd7 = []byte{
	// 270 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4b, 0x4c, 0x49, 0x29,
	0x4a, 0x2d, 0x2e, 0xce, 0x4d, 0x2c, 0xd0, 0x2f, 0x33, 0xd4, 0x2f, 0xa9, 0xd0, 0x2b, 0x28, 0xca,
	0x2f, 0xc9, 0x17, 0xe2, 0x45, 0x88, 0xeb, 0x95, 0x19, 0x2a, 0xd9, 0x71, 0x49, 0xfb, 0x16, 0xa7,
	0x3b, 0xa6, 0xa4, 0x38, 0x42, 0x84, 0x7d, 0x13, 0x0b, 0x0a, 0x32, 0xf3, 0xd2, 0x83, 0x52, 0x0b,
	0x4b, 0x53, 0x8b, 0x4b, 0x84, 0xe4, 0xb9, 0xb8, 0x53, 0xcb, 0x72, 0xe3, 0xa1, 0x7a, 0x24, 0x18,
	0x15, 0x18, 0x35, 0x38, 0x83, 0xb8, 0x52, 0xcb, 0x72, 0xa1, 0xca, 0x95, 0xe4, 0xb8, 0x64, 0xb0,
	0xeb, 0x2f, 0x2e, 0xc8, 0xcf, 0x2b, 0x4e, 0x55, 0x52, 0xe0, 0x92, 0xf3, 0x2d, 0x4e, 0x0f, 0x4a,
	0xcd, 0xcd, 0x2f, 0x4b, 0xc5, 0x6a, 0x85, 0x92, 0x22, 0x97, 0x3c, 0x4e, 0x15, 0x10, 0x43, 0x8c,
	0x5e, 0x31, 0x72, 0x31, 0xfb, 0x16, 0xa7, 0x0b, 0xe5, 0x70, 0x09, 0x62, 0xd8, 0x24, 0xa4, 0xa5,
	0x87, 0xe2, 0x23, 0x3d, 0x3c, 0xde, 0x91, 0xd2, 0x26, 0x4a, 0x2d, 0xc4, 0x56, 0xa1, 0x52, 0x2e,
	0x11, 0x6c, 0xae, 0x12, 0xd2, 0xc5, 0x34, 0x04, 0x8f, 0xff, 0xa4, 0xf4, 0x88, 0x55, 0x0e, 0xb1,
	0xd6, 0xc9, 0xff, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c,
	0xf0, 0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18, 0xa2, 0x4c, 0xd3, 0x33,
	0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0x0b, 0x52, 0xd3, 0xd3, 0x2b, 0xb3, 0xca,
	0xf4, 0x8b, 0xf3, 0x73, 0x73, 0x53, 0x73, 0x32, 0x53, 0x8b, 0xf4, 0xcb, 0xcc, 0xf5, 0x2b, 0xf4,
	0x91, 0x22, 0xbd, 0xa4, 0xb2, 0x20, 0xb5, 0x58, 0xbf, 0xcc, 0x30, 0x89, 0x0d, 0x1c, 0xf1, 0xc6,
	0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x68, 0x6f, 0xd2, 0x24, 0x12, 0x02, 0x00, 0x00,
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
	AddAddressMapping(ctx context.Context, in *MsgAddAddressMappingRequest, opts ...grpc.CallOption) (*MsgAddAddressMappingResponse, error)
	// Removes the mapping containing the cosmos address of the sender
	RemoveAddressMapping(ctx context.Context, in *MsgRemoveAddressMappingRequest, opts ...grpc.CallOption) (*MsgRemoveAddressMappingResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) AddAddressMapping(ctx context.Context, in *MsgAddAddressMappingRequest, opts ...grpc.CallOption) (*MsgAddAddressMappingResponse, error) {
	out := new(MsgAddAddressMappingResponse)
	err := c.cc.Invoke(ctx, "/addressmap.v1.Msg/AddAddressMapping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RemoveAddressMapping(ctx context.Context, in *MsgRemoveAddressMappingRequest, opts ...grpc.CallOption) (*MsgRemoveAddressMappingResponse, error) {
	out := new(MsgRemoveAddressMappingResponse)
	err := c.cc.Invoke(ctx, "/addressmap.v1.Msg/RemoveAddressMapping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// Adds a mapping between the cosmos address of the sender and the provided EVM address
	AddAddressMapping(context.Context, *MsgAddAddressMappingRequest) (*MsgAddAddressMappingResponse, error)
	// Removes the mapping containing the cosmos address of the sender
	RemoveAddressMapping(context.Context, *MsgRemoveAddressMappingRequest) (*MsgRemoveAddressMappingResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) AddAddressMapping(ctx context.Context, req *MsgAddAddressMappingRequest) (*MsgAddAddressMappingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAddressMapping not implemented")
}
func (*UnimplementedMsgServer) RemoveAddressMapping(ctx context.Context, req *MsgRemoveAddressMappingRequest) (*MsgRemoveAddressMappingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveAddressMapping not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_AddAddressMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddAddressMappingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddAddressMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/addressmap.v1.Msg/AddAddressMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddAddressMapping(ctx, req.(*MsgAddAddressMappingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemoveAddressMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveAddressMappingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RemoveAddressMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/addressmap.v1.Msg/RemoveAddressMapping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RemoveAddressMapping(ctx, req.(*MsgRemoveAddressMappingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "addressmap.v1.Msg",
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
	Metadata: "addressmap/v1/tx.proto",
}

func (m *MsgAddAddressMappingRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddAddressMappingRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddAddressMappingRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
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

func (m *MsgRemoveAddressMappingRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgRemoveAddressMappingRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgRemoveAddressMappingRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
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
func (m *MsgAddAddressMappingRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.EvmAddress)
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

func (m *MsgRemoveAddressMappingRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
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
func (m *MsgAddAddressMappingRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgAddAddressMappingRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddAddressMappingRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgRemoveAddressMappingRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgRemoveAddressMappingRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgRemoveAddressMappingRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
