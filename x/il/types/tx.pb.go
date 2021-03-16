// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: il/v1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
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

// MsgStoploss defines a stoploss position
type MsgCreateStoploss struct {
	// account address that owns the stoploss position
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// stoploss position details
	Stoploss *Stoploss `protobuf:"bytes,2,opt,name=stoploss,proto3" json:"stoploss,omitempty"`
}

func (m *MsgCreateStoploss) Reset()         { *m = MsgCreateStoploss{} }
func (m *MsgCreateStoploss) String() string { return proto.CompactTextString(m) }
func (*MsgCreateStoploss) ProtoMessage()    {}
func (*MsgCreateStoploss) Descriptor() ([]byte, []int) {
	return fileDescriptor_654cec8fc6a6f7e3, []int{0}
}
func (m *MsgCreateStoploss) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreateStoploss) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateStoploss.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreateStoploss) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateStoploss.Merge(m, src)
}
func (m *MsgCreateStoploss) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreateStoploss) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateStoploss.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateStoploss proto.InternalMessageInfo

func (m *MsgCreateStoploss) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *MsgCreateStoploss) GetStoploss() *Stoploss {
	if m != nil {
		return m.Stoploss
	}
	return nil
}

// MsgCreateStoplossResponse is the response type for the Msg/CreateStoploss gRPC method.
type MsgCreateStoplossResponse struct {
}

func (m *MsgCreateStoplossResponse) Reset()         { *m = MsgCreateStoplossResponse{} }
func (m *MsgCreateStoplossResponse) String() string { return proto.CompactTextString(m) }
func (*MsgCreateStoplossResponse) ProtoMessage()    {}
func (*MsgCreateStoplossResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_654cec8fc6a6f7e3, []int{1}
}
func (m *MsgCreateStoplossResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCreateStoplossResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCreateStoplossResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCreateStoplossResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCreateStoplossResponse.Merge(m, src)
}
func (m *MsgCreateStoplossResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgCreateStoplossResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCreateStoplossResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCreateStoplossResponse proto.InternalMessageInfo

// MsgDeleteStoploss removes a stoploss position
type MsgDeleteStoploss struct {
	// account address that owns the stoploss position
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	// uniswap pair hex address
	UniswapPairID string `protobuf:"bytes,2,opt,name=uniswap_pair_id,json=uniswapPairId,proto3" json:"uniswap_pair_id,omitempty"`
}

func (m *MsgDeleteStoploss) Reset()         { *m = MsgDeleteStoploss{} }
func (m *MsgDeleteStoploss) String() string { return proto.CompactTextString(m) }
func (*MsgDeleteStoploss) ProtoMessage()    {}
func (*MsgDeleteStoploss) Descriptor() ([]byte, []int) {
	return fileDescriptor_654cec8fc6a6f7e3, []int{2}
}
func (m *MsgDeleteStoploss) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDeleteStoploss) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDeleteStoploss.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDeleteStoploss) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDeleteStoploss.Merge(m, src)
}
func (m *MsgDeleteStoploss) XXX_Size() int {
	return m.Size()
}
func (m *MsgDeleteStoploss) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDeleteStoploss.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDeleteStoploss proto.InternalMessageInfo

func (m *MsgDeleteStoploss) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *MsgDeleteStoploss) GetUniswapPairID() string {
	if m != nil {
		return m.UniswapPairID
	}
	return ""
}

// MsgDeleteStoplossResponse is the response type for the Msg/DeleteStoploss gRPC method.
type MsgDeleteStoplossResponse struct {
}

func (m *MsgDeleteStoplossResponse) Reset()         { *m = MsgDeleteStoplossResponse{} }
func (m *MsgDeleteStoplossResponse) String() string { return proto.CompactTextString(m) }
func (*MsgDeleteStoplossResponse) ProtoMessage()    {}
func (*MsgDeleteStoplossResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_654cec8fc6a6f7e3, []int{3}
}
func (m *MsgDeleteStoplossResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDeleteStoplossResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDeleteStoplossResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDeleteStoplossResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDeleteStoplossResponse.Merge(m, src)
}
func (m *MsgDeleteStoplossResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgDeleteStoplossResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDeleteStoplossResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDeleteStoplossResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgCreateStoploss)(nil), "il.v1.MsgCreateStoploss")
	proto.RegisterType((*MsgCreateStoplossResponse)(nil), "il.v1.MsgCreateStoplossResponse")
	proto.RegisterType((*MsgDeleteStoploss)(nil), "il.v1.MsgDeleteStoploss")
	proto.RegisterType((*MsgDeleteStoplossResponse)(nil), "il.v1.MsgDeleteStoplossResponse")
}

func init() { proto.RegisterFile("il/v1/tx.proto", fileDescriptor_654cec8fc6a6f7e3) }

var fileDescriptor_654cec8fc6a6f7e3 = []byte{
	// 317 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcb, 0xcc, 0xd1, 0x2f,
	0x33, 0xd4, 0x2f, 0xa9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcd, 0xcc, 0xd1, 0x2b,
	0x33, 0x94, 0x82, 0x0a, 0x67, 0xe6, 0x40, 0x84, 0xa5, 0x44, 0xd2, 0xf3, 0xd3, 0xf3, 0xc1, 0x4c,
	0x7d, 0x10, 0x0b, 0x22, 0xaa, 0x14, 0xc5, 0x25, 0xe8, 0x5b, 0x9c, 0xee, 0x5c, 0x94, 0x9a, 0x58,
	0x92, 0x1a, 0x5c, 0x92, 0x5f, 0x90, 0x93, 0x5f, 0x5c, 0x2c, 0x24, 0xc1, 0xc5, 0x9e, 0x98, 0x92,
	0x52, 0x94, 0x5a, 0x5c, 0x2c, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0xe3, 0x0a, 0x69, 0x73,
	0x71, 0x14, 0x43, 0x55, 0x49, 0x30, 0x29, 0x30, 0x6a, 0x70, 0x1b, 0xf1, 0xeb, 0x81, 0xad, 0xd3,
	0x83, 0x69, 0x0e, 0x82, 0x2b, 0x50, 0x92, 0xe6, 0x92, 0xc4, 0x30, 0x3b, 0x28, 0xb5, 0xb8, 0x20,
	0x3f, 0xaf, 0x38, 0x55, 0x29, 0x03, 0x6c, 0xb1, 0x4b, 0x6a, 0x4e, 0x2a, 0x51, 0x16, 0x5b, 0x72,
	0xf1, 0x97, 0xe6, 0x65, 0x16, 0x97, 0x27, 0x16, 0xc4, 0x17, 0x24, 0x66, 0x16, 0xc5, 0x67, 0xa6,
	0x80, 0xed, 0xe7, 0x74, 0x12, 0x7c, 0x74, 0x4f, 0x9e, 0x37, 0x14, 0x22, 0x15, 0x90, 0x98, 0x59,
	0xe4, 0xe9, 0x12, 0xc4, 0x5b, 0x8a, 0xc4, 0x4d, 0x81, 0x3a, 0x03, 0xd5, 0x26, 0x98, 0x33, 0x8c,
	0x16, 0x32, 0x72, 0x31, 0xfb, 0x16, 0xa7, 0x0b, 0xf9, 0x70, 0xf1, 0xa1, 0x07, 0x02, 0xd4, 0x63,
	0x18, 0x5e, 0x90, 0x52, 0xc0, 0x25, 0x03, 0x33, 0x15, 0x64, 0x1a, 0xba, 0xcf, 0x10, 0x7a, 0x50,
	0x65, 0x90, 0x4d, 0xc3, 0xee, 0x46, 0x27, 0xc7, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63,
	0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f, 0xe5, 0x18, 0x6e, 0x3c, 0x96,
	0x63, 0x88, 0x52, 0x4f, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x2f, 0x48,
	0x4d, 0x4f, 0xaf, 0xcc, 0x2a, 0xd3, 0x2f, 0xce, 0xcf, 0xcd, 0x4d, 0xcd, 0xc9, 0x4c, 0x2d, 0xd2,
	0xaf, 0xd0, 0xcf, 0xcc, 0xd1, 0x2f, 0xa9, 0x2c, 0x48, 0x2d, 0x4e, 0x62, 0x03, 0xc7, 0xb6, 0x31,
	0x20, 0x00, 0x00, 0xff, 0xff, 0x75, 0xa2, 0xb2, 0x0c, 0x2c, 0x02, 0x00, 0x00,
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
	//CreateStoploss sets a new tracking stoploss position for a uniswap pair
	CreateStoploss(ctx context.Context, in *MsgCreateStoploss, opts ...grpc.CallOption) (*MsgCreateStoplossResponse, error)
	// DeleteStoploss deletes an existing stoploss position
	DeleteStoploss(ctx context.Context, in *MsgDeleteStoploss, opts ...grpc.CallOption) (*MsgDeleteStoplossResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CreateStoploss(ctx context.Context, in *MsgCreateStoploss, opts ...grpc.CallOption) (*MsgCreateStoplossResponse, error) {
	out := new(MsgCreateStoplossResponse)
	err := c.cc.Invoke(ctx, "/il.v1.Msg/CreateStoploss", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) DeleteStoploss(ctx context.Context, in *MsgDeleteStoploss, opts ...grpc.CallOption) (*MsgDeleteStoplossResponse, error) {
	out := new(MsgDeleteStoplossResponse)
	err := c.cc.Invoke(ctx, "/il.v1.Msg/DeleteStoploss", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	//CreateStoploss sets a new tracking stoploss position for a uniswap pair
	CreateStoploss(context.Context, *MsgCreateStoploss) (*MsgCreateStoplossResponse, error)
	// DeleteStoploss deletes an existing stoploss position
	DeleteStoploss(context.Context, *MsgDeleteStoploss) (*MsgDeleteStoplossResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) CreateStoploss(ctx context.Context, req *MsgCreateStoploss) (*MsgCreateStoplossResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateStoploss not implemented")
}
func (*UnimplementedMsgServer) DeleteStoploss(ctx context.Context, req *MsgDeleteStoploss) (*MsgDeleteStoplossResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStoploss not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_CreateStoploss_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateStoploss)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateStoploss(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/il.v1.Msg/CreateStoploss",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateStoploss(ctx, req.(*MsgCreateStoploss))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DeleteStoploss_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeleteStoploss)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DeleteStoploss(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/il.v1.Msg/DeleteStoploss",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DeleteStoploss(ctx, req.(*MsgDeleteStoploss))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "il.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateStoploss",
			Handler:    _Msg_CreateStoploss_Handler,
		},
		{
			MethodName: "DeleteStoploss",
			Handler:    _Msg_DeleteStoploss_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "il/v1/tx.proto",
}

func (m *MsgCreateStoploss) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateStoploss) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateStoploss) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Stoploss != nil {
		{
			size, err := m.Stoploss.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgCreateStoplossResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCreateStoplossResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCreateStoplossResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgDeleteStoploss) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDeleteStoploss) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDeleteStoploss) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.UniswapPairID) > 0 {
		i -= len(m.UniswapPairID)
		copy(dAtA[i:], m.UniswapPairID)
		i = encodeVarintTx(dAtA, i, uint64(len(m.UniswapPairID)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgDeleteStoplossResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDeleteStoplossResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDeleteStoplossResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *MsgCreateStoploss) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Stoploss != nil {
		l = m.Stoploss.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgCreateStoplossResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgDeleteStoploss) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.UniswapPairID)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgDeleteStoplossResponse) Size() (n int) {
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
func (m *MsgCreateStoploss) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgCreateStoploss: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateStoploss: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stoploss", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Stoploss == nil {
				m.Stoploss = &Stoploss{}
			}
			if err := m.Stoploss.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *MsgCreateStoplossResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgCreateStoplossResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCreateStoplossResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgDeleteStoploss) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgDeleteStoploss: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDeleteStoploss: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UniswapPairID", wireType)
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
			m.UniswapPairID = string(dAtA[iNdEx:postIndex])
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
func (m *MsgDeleteStoplossResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgDeleteStoplossResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDeleteStoplossResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
