// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cork/v2/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/codec/types"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/regen-network/cosmos-proto"
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

// MsgSubmitCorkRequest - sdk.Msg for submitting calls to Ethereum through the gravity bridge contract
type MsgSubmitCorkRequest struct {
	// the cork to send across the bridge
	Cork *Cork `protobuf:"bytes,1,opt,name=cork,proto3" json:"cork,omitempty"`
	// signer account address
	Signer string `protobuf:"bytes,2,opt,name=signer,proto3" json:"signer,omitempty"`
}

func (m *MsgSubmitCorkRequest) Reset()         { *m = MsgSubmitCorkRequest{} }
func (m *MsgSubmitCorkRequest) String() string { return proto.CompactTextString(m) }
func (*MsgSubmitCorkRequest) ProtoMessage()    {}
func (*MsgSubmitCorkRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_271bdc677f232222, []int{0}
}
func (m *MsgSubmitCorkRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSubmitCorkRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitCorkRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSubmitCorkRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSubmitCorkRequest.Merge(m, src)
}
func (m *MsgSubmitCorkRequest) XXX_Size() int {
	return m.Size()
}
func (m *MsgSubmitCorkRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSubmitCorkRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSubmitCorkRequest proto.InternalMessageInfo

func (m *MsgSubmitCorkRequest) GetCork() *Cork {
	if m != nil {
		return m.Cork
	}
	return nil
}

func (m *MsgSubmitCorkRequest) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

type MsgSubmitCorkResponse struct {
}

func (m *MsgSubmitCorkResponse) Reset()         { *m = MsgSubmitCorkResponse{} }
func (m *MsgSubmitCorkResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSubmitCorkResponse) ProtoMessage()    {}
func (*MsgSubmitCorkResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_271bdc677f232222, []int{1}
}
func (m *MsgSubmitCorkResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSubmitCorkResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitCorkResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSubmitCorkResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSubmitCorkResponse.Merge(m, src)
}
func (m *MsgSubmitCorkResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSubmitCorkResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSubmitCorkResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSubmitCorkResponse proto.InternalMessageInfo

// MsgScheduleCorkRequest - sdk.Msg for scheduling a cork request for on or after a specific block height
type MsgScheduleCorkRequest struct {
	// the scheduled cork
	Cork *Cork `protobuf:"bytes,1,opt,name=cork,proto3" json:"cork,omitempty"`
	// the block height that must be reached
	BlockHeight uint64 `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	// signer account address
	Signer string `protobuf:"bytes,3,opt,name=signer,proto3" json:"signer,omitempty"`
}

func (m *MsgScheduleCorkRequest) Reset()         { *m = MsgScheduleCorkRequest{} }
func (m *MsgScheduleCorkRequest) String() string { return proto.CompactTextString(m) }
func (*MsgScheduleCorkRequest) ProtoMessage()    {}
func (*MsgScheduleCorkRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_271bdc677f232222, []int{2}
}
func (m *MsgScheduleCorkRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgScheduleCorkRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgScheduleCorkRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgScheduleCorkRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgScheduleCorkRequest.Merge(m, src)
}
func (m *MsgScheduleCorkRequest) XXX_Size() int {
	return m.Size()
}
func (m *MsgScheduleCorkRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgScheduleCorkRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MsgScheduleCorkRequest proto.InternalMessageInfo

func (m *MsgScheduleCorkRequest) GetCork() *Cork {
	if m != nil {
		return m.Cork
	}
	return nil
}

func (m *MsgScheduleCorkRequest) GetBlockHeight() uint64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *MsgScheduleCorkRequest) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

type MsgScheduleCorkResponse struct {
}

func (m *MsgScheduleCorkResponse) Reset()         { *m = MsgScheduleCorkResponse{} }
func (m *MsgScheduleCorkResponse) String() string { return proto.CompactTextString(m) }
func (*MsgScheduleCorkResponse) ProtoMessage()    {}
func (*MsgScheduleCorkResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_271bdc677f232222, []int{3}
}
func (m *MsgScheduleCorkResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgScheduleCorkResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgScheduleCorkResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgScheduleCorkResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgScheduleCorkResponse.Merge(m, src)
}
func (m *MsgScheduleCorkResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgScheduleCorkResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgScheduleCorkResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgScheduleCorkResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgSubmitCorkRequest)(nil), "cork.v2.MsgSubmitCorkRequest")
	proto.RegisterType((*MsgSubmitCorkResponse)(nil), "cork.v2.MsgSubmitCorkResponse")
	proto.RegisterType((*MsgScheduleCorkRequest)(nil), "cork.v2.MsgScheduleCorkRequest")
	proto.RegisterType((*MsgScheduleCorkResponse)(nil), "cork.v2.MsgScheduleCorkResponse")
}

func init() { proto.RegisterFile("cork/v2/tx.proto", fileDescriptor_271bdc677f232222) }

var fileDescriptor_271bdc677f232222 = []byte{
	// 336 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x51, 0x4d, 0x4e, 0xc2, 0x40,
	0x18, 0x65, 0x84, 0x60, 0x1c, 0x30, 0x31, 0x13, 0xe5, 0xa7, 0x89, 0x23, 0xb0, 0x62, 0x61, 0x3a,
	0x49, 0xf5, 0x04, 0x9a, 0x18, 0x13, 0xc3, 0x02, 0xdc, 0xb9, 0x21, 0xb6, 0x8e, 0xd3, 0x4a, 0xcb,
	0x57, 0x3b, 0xd3, 0x06, 0x6e, 0xe1, 0x25, 0xbc, 0x8b, 0x4b, 0x96, 0x2e, 0x0d, 0xbd, 0x88, 0xe9,
	0xb4, 0x2a, 0x88, 0x2c, 0x5c, 0xbe, 0x9f, 0xbc, 0xf7, 0x66, 0x3e, 0x7c, 0xe0, 0x40, 0x34, 0x61,
	0x89, 0xc5, 0xd4, 0xcc, 0x0c, 0x23, 0x50, 0x40, 0x76, 0x33, 0xc6, 0x4c, 0x2c, 0xa3, 0xed, 0x80,
	0x0c, 0x40, 0x8e, 0x35, 0xcd, 0x72, 0x90, 0x7b, 0x8c, 0xb6, 0x00, 0x10, 0x3e, 0x67, 0x1a, 0xd9,
	0xf1, 0x23, 0xbb, 0x9f, 0xce, 0x0b, 0x89, 0x7c, 0x05, 0xea, 0x18, 0xcd, 0xf5, 0x86, 0xf8, 0x70,
	0x20, 0xc5, 0x6d, 0x6c, 0x07, 0x9e, 0xba, 0x84, 0x68, 0x32, 0xe2, 0xcf, 0x31, 0x97, 0x8a, 0x74,
	0x71, 0x25, 0x73, 0xb5, 0x50, 0x07, 0xf5, 0x6b, 0xd6, 0xbe, 0x59, 0x34, 0x9b, 0xda, 0xa3, 0x25,
	0xd2, 0xc0, 0x55, 0xe9, 0x89, 0x29, 0x8f, 0x5a, 0x3b, 0x1d, 0xd4, 0xdf, 0x1b, 0x15, 0xa8, 0xd7,
	0xc4, 0x47, 0xbf, 0x22, 0x65, 0x08, 0x53, 0xc9, 0x7b, 0x09, 0x6e, 0x64, 0x82, 0xe3, 0xf2, 0x87,
	0xd8, 0xe7, 0xff, 0x6c, 0xeb, 0xe2, 0xba, 0xed, 0x83, 0x33, 0x19, 0xbb, 0xdc, 0x13, 0xae, 0xd2,
	0x9d, 0x95, 0x51, 0x4d, 0x73, 0xd7, 0x9a, 0x5a, 0x19, 0x54, 0x5e, 0x1b, 0xd4, 0xc6, 0xcd, 0x8d,
	0xde, 0x7c, 0x92, 0xf5, 0x8a, 0x70, 0x79, 0x20, 0x05, 0xb9, 0xc1, 0xf8, 0x67, 0x30, 0x39, 0xfe,
	0x1e, 0xf0, 0xd7, 0xdf, 0x18, 0x74, 0x9b, 0x9c, 0x87, 0x92, 0x21, 0xae, 0xaf, 0x96, 0x91, 0x93,
	0x35, 0xff, 0xe6, 0xf3, 0x8d, 0xce, 0x76, 0x43, 0x1e, 0x79, 0x71, 0xf5, 0xb6, 0xa4, 0x68, 0xb1,
	0xa4, 0xe8, 0x63, 0x49, 0xd1, 0x4b, 0x4a, 0x4b, 0x8b, 0x94, 0x96, 0xde, 0x53, 0x5a, 0xba, 0x3b,
	0x15, 0x9e, 0x72, 0x63, 0xdb, 0x74, 0x20, 0x60, 0x21, 0x17, 0x62, 0xfe, 0x94, 0x30, 0x09, 0x41,
	0xc0, 0x7d, 0x8f, 0x47, 0x2c, 0x39, 0x67, 0x33, 0x7d, 0x6e, 0xa6, 0xe6, 0x21, 0x97, 0x76, 0x55,
	0x5f, 0xfd, 0xec, 0x33, 0x00, 0x00, 0xff, 0xff, 0xbb, 0x9d, 0x4a, 0x57, 0x5c, 0x02, 0x00, 0x00,
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
	SubmitCork(ctx context.Context, in *MsgSubmitCorkRequest, opts ...grpc.CallOption) (*MsgSubmitCorkResponse, error)
	ScheduleCork(ctx context.Context, in *MsgScheduleCorkRequest, opts ...grpc.CallOption) (*MsgScheduleCorkResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SubmitCork(ctx context.Context, in *MsgSubmitCorkRequest, opts ...grpc.CallOption) (*MsgSubmitCorkResponse, error) {
	out := new(MsgSubmitCorkResponse)
	err := c.cc.Invoke(ctx, "/cork.v2.Msg/SubmitCork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ScheduleCork(ctx context.Context, in *MsgScheduleCorkRequest, opts ...grpc.CallOption) (*MsgScheduleCorkResponse, error) {
	out := new(MsgScheduleCorkResponse)
	err := c.cc.Invoke(ctx, "/cork.v2.Msg/ScheduleCork", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	SubmitCork(context.Context, *MsgSubmitCorkRequest) (*MsgSubmitCorkResponse, error)
	ScheduleCork(context.Context, *MsgScheduleCorkRequest) (*MsgScheduleCorkResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) SubmitCork(ctx context.Context, req *MsgSubmitCorkRequest) (*MsgSubmitCorkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitCork not implemented")
}
func (*UnimplementedMsgServer) ScheduleCork(ctx context.Context, req *MsgScheduleCorkRequest) (*MsgScheduleCorkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScheduleCork not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_SubmitCork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitCorkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitCork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cork.v2.Msg/SubmitCork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitCork(ctx, req.(*MsgSubmitCorkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ScheduleCork_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgScheduleCorkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ScheduleCork(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cork.v2.Msg/ScheduleCork",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ScheduleCork(ctx, req.(*MsgScheduleCorkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cork.v2.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitCork",
			Handler:    _Msg_SubmitCork_Handler,
		},
		{
			MethodName: "ScheduleCork",
			Handler:    _Msg_ScheduleCork_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cork/v2/tx.proto",
}

func (m *MsgSubmitCorkRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSubmitCorkRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSubmitCorkRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if m.Cork != nil {
		{
			size, err := m.Cork.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSubmitCorkResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSubmitCorkResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSubmitCorkResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgScheduleCorkRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgScheduleCorkRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgScheduleCorkRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0x1a
	}
	if m.BlockHeight != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x10
	}
	if m.Cork != nil {
		{
			size, err := m.Cork.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTx(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgScheduleCorkResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgScheduleCorkResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgScheduleCorkResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *MsgSubmitCorkRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Cork != nil {
		l = m.Cork.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgSubmitCorkResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgScheduleCorkRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Cork != nil {
		l = m.Cork.Size()
		n += 1 + l + sovTx(uint64(l))
	}
	if m.BlockHeight != 0 {
		n += 1 + sovTx(uint64(m.BlockHeight))
	}
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgScheduleCorkResponse) Size() (n int) {
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
func (m *MsgSubmitCorkRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSubmitCorkRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSubmitCorkRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cork", wireType)
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
			if m.Cork == nil {
				m.Cork = &Cork{}
			}
			if err := m.Cork.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *MsgSubmitCorkResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSubmitCorkResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSubmitCorkResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgScheduleCorkRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgScheduleCorkRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgScheduleCorkRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cork", wireType)
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
			if m.Cork == nil {
				m.Cork = &Cork{}
			}
			if err := m.Cork.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
		case 3:
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
func (m *MsgScheduleCorkResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgScheduleCorkResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgScheduleCorkResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
