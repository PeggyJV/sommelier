// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: allocation/v1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/cosmos/cosmos-sdk/codec/types"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/regen-network/cosmos-proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// MsgDelegateAllocations defines sdk.Msg for delegating allocation rights from a validator
// to another address, must be signed by an active validator
type MsgDelegateAllocations struct {
	// delegate account address
	Delegate string `protobuf:"bytes,1,opt,name=delegate,proto3" json:"delegate,omitempty"`
	// validator operator address
	Validator string `protobuf:"bytes,2,opt,name=validator,proto3" json:"validator,omitempty"`
}

func (m *MsgDelegateAllocations) Reset()         { *m = MsgDelegateAllocations{} }
func (m *MsgDelegateAllocations) String() string { return proto.CompactTextString(m) }
func (*MsgDelegateAllocations) ProtoMessage()    {}
func (*MsgDelegateAllocations) Descriptor() ([]byte, []int) {
	return fileDescriptor_194e979be0693c53, []int{0}
}
func (m *MsgDelegateAllocations) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDelegateAllocations) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDelegateAllocations.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDelegateAllocations) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDelegateAllocations.Merge(m, src)
}
func (m *MsgDelegateAllocations) XXX_Size() int {
	return m.Size()
}
func (m *MsgDelegateAllocations) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDelegateAllocations.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDelegateAllocations proto.InternalMessageInfo

func (m *MsgDelegateAllocations) GetDelegate() string {
	if m != nil {
		return m.Delegate
	}
	return ""
}

func (m *MsgDelegateAllocations) GetValidator() string {
	if m != nil {
		return m.Validator
	}
	return ""
}

// MsgDelegateAllocationsResponse is the response type for the Msg/DelegateAllocations gRPC method.
type MsgDelegateAllocationsResponse struct {
}

func (m *MsgDelegateAllocationsResponse) Reset()         { *m = MsgDelegateAllocationsResponse{} }
func (m *MsgDelegateAllocationsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgDelegateAllocationsResponse) ProtoMessage()    {}
func (*MsgDelegateAllocationsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_194e979be0693c53, []int{1}
}
func (m *MsgDelegateAllocationsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgDelegateAllocationsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgDelegateAllocationsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgDelegateAllocationsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgDelegateAllocationsResponse.Merge(m, src)
}
func (m *MsgDelegateAllocationsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgDelegateAllocationsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgDelegateAllocationsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgDelegateAllocationsResponse proto.InternalMessageInfo

// MsgAllocationPrecommit - sdk.Msg for prevoting on an array of oracle data types.
// The purpose of the prevote is to hide vote for data with hashes formatted as hex string:
// SHA256("{salt}:{data_cannonical_json}:{voter}")
type MsgAllocationPrecommit struct {
	// precommit containing the hash of the allocation precommit contents
	Precommit []*AllocationPrecommit `protobuf:"bytes,1,rep,name=precommit,proto3" json:"precommit,omitempty"`
	// signer (i.e feeder) account address
	Signer string `protobuf:"bytes,2,opt,name=signer,proto3" json:"signer,omitempty"`
}

func (m *MsgAllocationPrecommit) Reset()         { *m = MsgAllocationPrecommit{} }
func (m *MsgAllocationPrecommit) String() string { return proto.CompactTextString(m) }
func (*MsgAllocationPrecommit) ProtoMessage()    {}
func (*MsgAllocationPrecommit) Descriptor() ([]byte, []int) {
	return fileDescriptor_194e979be0693c53, []int{2}
}
func (m *MsgAllocationPrecommit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAllocationPrecommit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAllocationPrecommit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAllocationPrecommit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAllocationPrecommit.Merge(m, src)
}
func (m *MsgAllocationPrecommit) XXX_Size() int {
	return m.Size()
}
func (m *MsgAllocationPrecommit) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAllocationPrecommit.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAllocationPrecommit proto.InternalMessageInfo

func (m *MsgAllocationPrecommit) GetPrecommit() []*AllocationPrecommit {
	if m != nil {
		return m.Precommit
	}
	return nil
}

func (m *MsgAllocationPrecommit) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

// MsgAllocationPrecommitResponse is the response type for the Msg/AllocationPrecommitResponse gRPC method.
type MsgAllocationPrecommitResponse struct {
}

func (m *MsgAllocationPrecommitResponse) Reset()         { *m = MsgAllocationPrecommitResponse{} }
func (m *MsgAllocationPrecommitResponse) String() string { return proto.CompactTextString(m) }
func (*MsgAllocationPrecommitResponse) ProtoMessage()    {}
func (*MsgAllocationPrecommitResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_194e979be0693c53, []int{3}
}
func (m *MsgAllocationPrecommitResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAllocationPrecommitResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAllocationPrecommitResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAllocationPrecommitResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAllocationPrecommitResponse.Merge(m, src)
}
func (m *MsgAllocationPrecommitResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgAllocationPrecommitResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAllocationPrecommitResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAllocationPrecommitResponse proto.InternalMessageInfo

// MsgAllocationCommit - sdk.Msg for submitting arbitrary oracle data that has been prevoted on
type MsgAllocationCommit struct {
	// vote containing the oracle data feed
	Commit []*Allocation `protobuf:"bytes,1,rep,name=commit,proto3" json:"commit,omitempty"`
	// signer (i.e feeder) account address
	Signer string `protobuf:"bytes,2,opt,name=signer,proto3" json:"signer,omitempty"`
}

func (m *MsgAllocationCommit) Reset()         { *m = MsgAllocationCommit{} }
func (m *MsgAllocationCommit) String() string { return proto.CompactTextString(m) }
func (*MsgAllocationCommit) ProtoMessage()    {}
func (*MsgAllocationCommit) Descriptor() ([]byte, []int) {
	return fileDescriptor_194e979be0693c53, []int{4}
}
func (m *MsgAllocationCommit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAllocationCommit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAllocationCommit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAllocationCommit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAllocationCommit.Merge(m, src)
}
func (m *MsgAllocationCommit) XXX_Size() int {
	return m.Size()
}
func (m *MsgAllocationCommit) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAllocationCommit.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAllocationCommit proto.InternalMessageInfo

func (m *MsgAllocationCommit) GetCommit() []*Allocation {
	if m != nil {
		return m.Commit
	}
	return nil
}

func (m *MsgAllocationCommit) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

// MsgAllocationCommitResponse is the response type for the Msg/AllocationCommit gRPC method.
type MsgAllocationCommitResponse struct {
}

func (m *MsgAllocationCommitResponse) Reset()         { *m = MsgAllocationCommitResponse{} }
func (m *MsgAllocationCommitResponse) String() string { return proto.CompactTextString(m) }
func (*MsgAllocationCommitResponse) ProtoMessage()    {}
func (*MsgAllocationCommitResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_194e979be0693c53, []int{5}
}
func (m *MsgAllocationCommitResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAllocationCommitResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAllocationCommitResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAllocationCommitResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAllocationCommitResponse.Merge(m, src)
}
func (m *MsgAllocationCommitResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgAllocationCommitResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAllocationCommitResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAllocationCommitResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgDelegateAllocations)(nil), "allocation.v1.MsgDelegateAllocations")
	proto.RegisterType((*MsgDelegateAllocationsResponse)(nil), "allocation.v1.MsgDelegateAllocationsResponse")
	proto.RegisterType((*MsgAllocationPrecommit)(nil), "allocation.v1.MsgAllocationPrecommit")
	proto.RegisterType((*MsgAllocationPrecommitResponse)(nil), "allocation.v1.MsgAllocationPrecommitResponse")
	proto.RegisterType((*MsgAllocationCommit)(nil), "allocation.v1.MsgAllocationCommit")
	proto.RegisterType((*MsgAllocationCommitResponse)(nil), "allocation.v1.MsgAllocationCommitResponse")
}

func init() { proto.RegisterFile("allocation/v1/tx.proto", fileDescriptor_194e979be0693c53) }

var fileDescriptor_194e979be0693c53 = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0xcd, 0x4a, 0xc3, 0x40,
	0x18, 0x6c, 0x5a, 0x08, 0x76, 0x45, 0x90, 0x14, 0x4a, 0x1b, 0x75, 0x29, 0x01, 0xa1, 0x08, 0x66,
	0x69, 0x7d, 0x01, 0xff, 0x2e, 0x1e, 0x0a, 0x92, 0xa3, 0x17, 0x4d, 0xd2, 0x75, 0x8d, 0x26, 0xf9,
	0x42, 0x76, 0x1b, 0xda, 0xb7, 0xf0, 0x9d, 0xbc, 0x78, 0xec, 0xd1, 0xa3, 0xb4, 0x2f, 0x22, 0xe6,
	0xa7, 0x69, 0x74, 0x6d, 0x3d, 0xce, 0x37, 0xf3, 0xcd, 0x0c, 0xdf, 0xb2, 0xa8, 0x6d, 0xfb, 0x3e,
	0xb8, 0xb6, 0xf0, 0x20, 0x24, 0xc9, 0x80, 0x88, 0xa9, 0x19, 0xc5, 0x20, 0x40, 0xdb, 0x2b, 0xe7,
	0x66, 0x32, 0xd0, 0xbb, 0x2e, 0xf0, 0x00, 0xf8, 0x7d, 0x4a, 0x92, 0x0c, 0x64, 0x4a, 0xbd, 0xcb,
	0x00, 0x98, 0x4f, 0x49, 0x8a, 0x9c, 0xc9, 0x23, 0xb1, 0xc3, 0x59, 0x4e, 0xe1, 0xaa, 0xf9, 0x9a,
	0x65, 0xca, 0x1b, 0x16, 0x6a, 0x8f, 0x38, 0xbb, 0xa6, 0x3e, 0x65, 0xb6, 0xa0, 0x17, 0x2b, 0x9a,
	0x6b, 0x3a, 0xda, 0x19, 0xe7, 0xe3, 0x8e, 0xd2, 0x53, 0xfa, 0x4d, 0x6b, 0x85, 0xb5, 0x43, 0xd4,
	0x4c, 0x6c, 0xdf, 0x1b, 0xdb, 0x02, 0xe2, 0x4e, 0x3d, 0x25, 0xcb, 0x81, 0xd1, 0x43, 0x58, 0xee,
	0x69, 0x51, 0x1e, 0x41, 0xc8, 0xa9, 0x11, 0xa7, 0xa9, 0x25, 0x73, 0x1b, 0x53, 0x17, 0x82, 0xc0,
	0x13, 0xda, 0x39, 0x6a, 0x46, 0x05, 0xe8, 0x28, 0xbd, 0x46, 0x7f, 0x77, 0x68, 0x98, 0x95, 0x43,
	0x98, 0x92, 0x35, 0xab, 0x5c, 0xd2, 0xda, 0x48, 0xe5, 0x1e, 0x0b, 0x69, 0x51, 0x2c, 0x47, 0x79,
	0x2b, 0xd9, 0x72, 0xd1, 0xea, 0x01, 0xb5, 0x2a, 0x8a, 0xab, 0xcc, 0x70, 0x80, 0xd4, 0x4a, 0x9f,
	0xee, 0x9f, 0x7d, 0x2c, 0x75, 0x4b, 0x87, 0x23, 0x74, 0x20, 0x49, 0x28, 0x0a, 0x0c, 0xdf, 0xea,
	0xa8, 0x31, 0xe2, 0x4c, 0x7b, 0x41, 0x2d, 0xd9, 0x8b, 0x1c, 0xff, 0x08, 0x96, 0x1f, 0x59, 0x3f,
	0xfd, 0x97, 0xac, 0x08, 0xfd, 0x0e, 0x93, 0x3d, 0x84, 0x24, 0x4c, 0x22, 0x93, 0x85, 0x6d, 0x38,
	0xb1, 0xe6, 0xa0, 0xfd, 0x5f, 0xf7, 0x35, 0x36, 0x59, 0x64, 0x1a, 0xfd, 0x64, 0xbb, 0xa6, 0xc8,
	0xb8, 0xbc, 0x79, 0x5f, 0x60, 0x65, 0xbe, 0xc0, 0xca, 0xe7, 0x02, 0x2b, 0xaf, 0x4b, 0x5c, 0x9b,
	0x2f, 0x71, 0xed, 0x63, 0x89, 0x6b, 0x77, 0x84, 0x79, 0xe2, 0x69, 0xe2, 0x98, 0x2e, 0x04, 0x24,
	0xa2, 0x8c, 0xcd, 0x9e, 0x13, 0xc2, 0x21, 0x08, 0xa8, 0xef, 0xd1, 0x98, 0x4c, 0xd7, 0x7e, 0x07,
	0x11, 0xb3, 0x88, 0x72, 0x47, 0x4d, 0x3f, 0xc9, 0xd9, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x0e,
	0xd0, 0xe6, 0x9e, 0xa3, 0x03, 0x00, 0x00,
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
	// DelegateAllocations defines a message that delegates the allocating to an account address.
	DelegateAllocations(ctx context.Context, in *MsgDelegateAllocations, opts ...grpc.CallOption) (*MsgDelegateAllocationsResponse, error)
	// OracleDataPrevote defines a message that commits a hash of a oracle data feed before the data is actually submitted.
	AllocationPrecommit(ctx context.Context, in *MsgAllocationPrecommit, opts ...grpc.CallOption) (*MsgAllocationPrecommitResponse, error)
	// OracleDataVote defines a message to submit the actual oracle data that was committed by the feeder through the prevote.
	AllocationCommit(ctx context.Context, in *MsgAllocationCommit, opts ...grpc.CallOption) (*MsgAllocationCommitResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) DelegateAllocations(ctx context.Context, in *MsgDelegateAllocations, opts ...grpc.CallOption) (*MsgDelegateAllocationsResponse, error) {
	out := new(MsgDelegateAllocationsResponse)
	err := c.cc.Invoke(ctx, "/allocation.v1.Msg/DelegateAllocations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) AllocationPrecommit(ctx context.Context, in *MsgAllocationPrecommit, opts ...grpc.CallOption) (*MsgAllocationPrecommitResponse, error) {
	out := new(MsgAllocationPrecommitResponse)
	err := c.cc.Invoke(ctx, "/allocation.v1.Msg/AllocationPrecommit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) AllocationCommit(ctx context.Context, in *MsgAllocationCommit, opts ...grpc.CallOption) (*MsgAllocationCommitResponse, error) {
	out := new(MsgAllocationCommitResponse)
	err := c.cc.Invoke(ctx, "/allocation.v1.Msg/AllocationCommit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// DelegateAllocations defines a message that delegates the allocating to an account address.
	DelegateAllocations(context.Context, *MsgDelegateAllocations) (*MsgDelegateAllocationsResponse, error)
	// OracleDataPrevote defines a message that commits a hash of a oracle data feed before the data is actually submitted.
	AllocationPrecommit(context.Context, *MsgAllocationPrecommit) (*MsgAllocationPrecommitResponse, error)
	// OracleDataVote defines a message to submit the actual oracle data that was committed by the feeder through the prevote.
	AllocationCommit(context.Context, *MsgAllocationCommit) (*MsgAllocationCommitResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) DelegateAllocations(ctx context.Context, req *MsgDelegateAllocations) (*MsgDelegateAllocationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelegateAllocations not implemented")
}
func (*UnimplementedMsgServer) AllocationPrecommit(ctx context.Context, req *MsgAllocationPrecommit) (*MsgAllocationPrecommitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllocationPrecommit not implemented")
}
func (*UnimplementedMsgServer) AllocationCommit(ctx context.Context, req *MsgAllocationCommit) (*MsgAllocationCommitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AllocationCommit not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_DelegateAllocations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDelegateAllocations)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DelegateAllocations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/allocation.v1.Msg/DelegateAllocations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DelegateAllocations(ctx, req.(*MsgDelegateAllocations))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AllocationPrecommit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAllocationPrecommit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AllocationPrecommit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/allocation.v1.Msg/AllocationPrecommit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AllocationPrecommit(ctx, req.(*MsgAllocationPrecommit))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AllocationCommit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAllocationCommit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AllocationCommit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/allocation.v1.Msg/AllocationCommit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AllocationCommit(ctx, req.(*MsgAllocationCommit))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "allocation.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DelegateAllocations",
			Handler:    _Msg_DelegateAllocations_Handler,
		},
		{
			MethodName: "AllocationPrecommit",
			Handler:    _Msg_AllocationPrecommit_Handler,
		},
		{
			MethodName: "AllocationCommit",
			Handler:    _Msg_AllocationCommit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "allocation/v1/tx.proto",
}

func (m *MsgDelegateAllocations) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDelegateAllocations) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDelegateAllocations) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Validator) > 0 {
		i -= len(m.Validator)
		copy(dAtA[i:], m.Validator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Validator)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Delegate) > 0 {
		i -= len(m.Delegate)
		copy(dAtA[i:], m.Delegate)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Delegate)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgDelegateAllocationsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDelegateAllocationsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDelegateAllocationsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgAllocationPrecommit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAllocationPrecommit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAllocationPrecommit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.Precommit) > 0 {
		for iNdEx := len(m.Precommit) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Precommit[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *MsgAllocationPrecommitResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAllocationPrecommitResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAllocationPrecommitResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgAllocationCommit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAllocationCommit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAllocationCommit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.Commit) > 0 {
		for iNdEx := len(m.Commit) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Commit[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *MsgAllocationCommitResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAllocationCommitResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAllocationCommitResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *MsgDelegateAllocations) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Delegate)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Validator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgDelegateAllocationsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgAllocationPrecommit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Precommit) > 0 {
		for _, e := range m.Precommit {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgAllocationPrecommitResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgAllocationCommit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Commit) > 0 {
		for _, e := range m.Commit {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgAllocationCommitResponse) Size() (n int) {
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
func (m *MsgDelegateAllocations) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgDelegateAllocations: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDelegateAllocations: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Delegate", wireType)
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
			m.Delegate = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validator", wireType)
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
			m.Validator = string(dAtA[iNdEx:postIndex])
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
func (m *MsgDelegateAllocationsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgDelegateAllocationsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDelegateAllocationsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgAllocationPrecommit) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgAllocationPrecommit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAllocationPrecommit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Precommit", wireType)
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
			m.Precommit = append(m.Precommit, &AllocationPrecommit{})
			if err := m.Precommit[len(m.Precommit)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *MsgAllocationPrecommitResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgAllocationPrecommitResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAllocationPrecommitResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgAllocationCommit) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgAllocationCommit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAllocationCommit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Commit", wireType)
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
			m.Commit = append(m.Commit, &Allocation{})
			if err := m.Commit[len(m.Commit)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *MsgAllocationCommitResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgAllocationCommitResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAllocationCommitResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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