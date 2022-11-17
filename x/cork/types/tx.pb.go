// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cork/v2/tx.proto

package types

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
	return fileDescriptor_271bdc677f232222, []int{0}
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
	// cork ID
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *MsgScheduleCorkResponse) Reset()         { *m = MsgScheduleCorkResponse{} }
func (m *MsgScheduleCorkResponse) String() string { return proto.CompactTextString(m) }
func (*MsgScheduleCorkResponse) ProtoMessage()    {}
func (*MsgScheduleCorkResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_271bdc677f232222, []int{1}
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

func (m *MsgScheduleCorkResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func init() {
	proto.RegisterType((*MsgScheduleCorkRequest)(nil), "cork.v2.MsgScheduleCorkRequest")
	proto.RegisterType((*MsgScheduleCorkResponse)(nil), "cork.v2.MsgScheduleCorkResponse")
}

func init() { proto.RegisterFile("cork/v2/tx.proto", fileDescriptor_271bdc677f232222) }

var fileDescriptor_271bdc677f232222 = []byte{
	// 280 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0xc1, 0x4a, 0xc3, 0x30,
	0x1c, 0xc6, 0x9b, 0x6e, 0x4c, 0x96, 0x4d, 0x91, 0x1c, 0x66, 0xd9, 0x21, 0x76, 0x3b, 0x55, 0x90,
	0x06, 0xaa, 0x4f, 0xa0, 0x20, 0x5e, 0x76, 0xb0, 0x5e, 0xc4, 0x8b, 0xd0, 0xf6, 0x4f, 0x1a, 0xdb,
	0x2e, 0xb5, 0x49, 0xcb, 0xf6, 0x16, 0x3e, 0x96, 0xc7, 0x1d, 0x3d, 0x4a, 0xfb, 0x22, 0xb2, 0x38,
	0x45, 0xd0, 0x1d, 0xf3, 0xfb, 0x7e, 0xe4, 0x4b, 0x3e, 0x7c, 0x1c, 0xcb, 0x2a, 0x63, 0x4d, 0xc0,
	0xf4, 0xca, 0x2f, 0x2b, 0xa9, 0x25, 0x39, 0xd8, 0x12, 0xbf, 0x09, 0xa6, 0xe4, 0x3b, 0x32, 0xc0,
	0x84, 0xf3, 0x06, 0x4f, 0x16, 0x8a, 0xdf, 0xc7, 0x29, 0x24, 0x75, 0x0e, 0xd7, 0xb2, 0xca, 0x42,
	0x78, 0xa9, 0x41, 0x69, 0x32, 0xc3, 0xfd, 0xad, 0xe7, 0x20, 0x17, 0x79, 0xa3, 0xe0, 0xd0, 0xdf,
	0xdd, 0xe2, 0x1b, 0xc7, 0x44, 0x64, 0x86, 0xc7, 0x51, 0x2e, 0xe3, 0xec, 0x29, 0x05, 0xc1, 0x53,
	0xed, 0xd8, 0x2e, 0xf2, 0xfa, 0xe1, 0xc8, 0xb0, 0x5b, 0x83, 0xc8, 0x04, 0x0f, 0x94, 0xe0, 0x4b,
	0xa8, 0x9c, 0x9e, 0x8b, 0xbc, 0x61, 0xb8, 0x3b, 0xcd, 0xcf, 0xf0, 0xc9, 0x9f, 0x5e, 0x55, 0xca,
	0xa5, 0x02, 0x72, 0x84, 0x6d, 0x91, 0x98, 0xda, 0x61, 0x68, 0x8b, 0x24, 0x78, 0xc0, 0xbd, 0x85,
	0xe2, 0xe4, 0x0e, 0x8f, 0x7f, 0xeb, 0xe4, 0xf4, 0xe7, 0x45, 0xff, 0x7f, 0x60, 0xea, 0xee, 0x17,
	0xbe, 0x9a, 0xae, 0x6e, 0xde, 0x5a, 0x8a, 0x36, 0x2d, 0x45, 0x1f, 0x2d, 0x45, 0xaf, 0x1d, 0xb5,
	0x36, 0x1d, 0xb5, 0xde, 0x3b, 0x6a, 0x3d, 0x9e, 0x73, 0xa1, 0xd3, 0x3a, 0xf2, 0x63, 0x59, 0xb0,
	0x12, 0x38, 0x5f, 0x3f, 0x37, 0x4c, 0xc9, 0xa2, 0x80, 0x5c, 0x40, 0xc5, 0x9a, 0x4b, 0xb6, 0x32,
	0x23, 0x32, 0xbd, 0x2e, 0x41, 0x45, 0x03, 0xb3, 0xe5, 0xc5, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x21, 0x64, 0x2a, 0x01, 0x7c, 0x01, 0x00, 0x00,
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
	ScheduleCork(ctx context.Context, in *MsgScheduleCorkRequest, opts ...grpc.CallOption) (*MsgScheduleCorkResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
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
	ScheduleCork(context.Context, *MsgScheduleCorkRequest) (*MsgScheduleCorkResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) ScheduleCork(ctx context.Context, req *MsgScheduleCorkRequest) (*MsgScheduleCorkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ScheduleCork not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
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
			MethodName: "ScheduleCork",
			Handler:    _Msg_ScheduleCork_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "cork/v2/tx.proto",
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
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
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
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
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
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
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
			m.Id = string(dAtA[iNdEx:postIndex])
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
