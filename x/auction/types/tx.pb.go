// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: auction/v1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/gogoproto/gogoproto"
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

type MsgSubmitBidRequest struct {
	AuctionId              uint32     `protobuf:"varint,1,opt,name=auction_id,json=auctionId,proto3" json:"auction_id,omitempty"`
	Signer                 string     `protobuf:"bytes,2,opt,name=signer,proto3" json:"signer,omitempty"`
	MaxBidInUsomm          types.Coin `protobuf:"bytes,3,opt,name=max_bid_in_usomm,json=maxBidInUsomm,proto3" json:"max_bid_in_usomm"`
	SaleTokenMinimumAmount types.Coin `protobuf:"bytes,4,opt,name=sale_token_minimum_amount,json=saleTokenMinimumAmount,proto3" json:"sale_token_minimum_amount"`
}

func (m *MsgSubmitBidRequest) Reset()         { *m = MsgSubmitBidRequest{} }
func (m *MsgSubmitBidRequest) String() string { return proto.CompactTextString(m) }
func (*MsgSubmitBidRequest) ProtoMessage()    {}
func (*MsgSubmitBidRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a60fe804de30894a, []int{0}
}
func (m *MsgSubmitBidRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSubmitBidRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitBidRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSubmitBidRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSubmitBidRequest.Merge(m, src)
}
func (m *MsgSubmitBidRequest) XXX_Size() int {
	return m.Size()
}
func (m *MsgSubmitBidRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSubmitBidRequest.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSubmitBidRequest proto.InternalMessageInfo

func (m *MsgSubmitBidRequest) GetAuctionId() uint32 {
	if m != nil {
		return m.AuctionId
	}
	return 0
}

func (m *MsgSubmitBidRequest) GetSigner() string {
	if m != nil {
		return m.Signer
	}
	return ""
}

func (m *MsgSubmitBidRequest) GetMaxBidInUsomm() types.Coin {
	if m != nil {
		return m.MaxBidInUsomm
	}
	return types.Coin{}
}

func (m *MsgSubmitBidRequest) GetSaleTokenMinimumAmount() types.Coin {
	if m != nil {
		return m.SaleTokenMinimumAmount
	}
	return types.Coin{}
}

type MsgSubmitBidResponse struct {
	Bid *Bid `protobuf:"bytes,1,opt,name=bid,proto3" json:"bid,omitempty"`
}

func (m *MsgSubmitBidResponse) Reset()         { *m = MsgSubmitBidResponse{} }
func (m *MsgSubmitBidResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSubmitBidResponse) ProtoMessage()    {}
func (*MsgSubmitBidResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a60fe804de30894a, []int{1}
}
func (m *MsgSubmitBidResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSubmitBidResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSubmitBidResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSubmitBidResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSubmitBidResponse.Merge(m, src)
}
func (m *MsgSubmitBidResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSubmitBidResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSubmitBidResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSubmitBidResponse proto.InternalMessageInfo

func (m *MsgSubmitBidResponse) GetBid() *Bid {
	if m != nil {
		return m.Bid
	}
	return nil
}

func init() {
	proto.RegisterType((*MsgSubmitBidRequest)(nil), "auction.v1.MsgSubmitBidRequest")
	proto.RegisterType((*MsgSubmitBidResponse)(nil), "auction.v1.MsgSubmitBidResponse")
}

func init() { proto.RegisterFile("auction/v1/tx.proto", fileDescriptor_a60fe804de30894a) }

var fileDescriptor_a60fe804de30894a = []byte{
	// 423 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x3f, 0x6f, 0xd4, 0x30,
	0x14, 0x3f, 0x73, 0xa5, 0xd2, 0xf9, 0x54, 0x81, 0xdc, 0xaa, 0xe4, 0x4e, 0x22, 0x3d, 0x3a, 0x9d,
	0x3a, 0xd8, 0xca, 0x31, 0x20, 0xd8, 0x08, 0x0b, 0x37, 0x1c, 0x43, 0x80, 0xa5, 0x42, 0xb2, 0xf2,
	0xc7, 0x32, 0x86, 0xda, 0x0e, 0x67, 0x27, 0x4a, 0x37, 0xc4, 0x27, 0xe0, 0xa3, 0xf4, 0x63, 0x74,
	0xec, 0xc8, 0x84, 0xd0, 0xdd, 0xd0, 0x8f, 0x01, 0x72, 0xe2, 0x40, 0x91, 0x40, 0xea, 0x96, 0xf7,
	0xfb, 0xf3, 0xf2, 0xde, 0xef, 0x19, 0xee, 0xa7, 0x55, 0x6e, 0x85, 0x56, 0xa4, 0x8e, 0x88, 0x6d,
	0x70, 0xb9, 0xd6, 0x56, 0x23, 0xe8, 0x41, 0x5c, 0x47, 0xd3, 0xe0, 0x86, 0xa0, 0x87, 0x5b, 0xd5,
	0x34, 0xcc, 0xb5, 0x91, 0xda, 0x90, 0x2c, 0x35, 0x8c, 0xd4, 0x51, 0xc6, 0x6c, 0x1a, 0x91, 0x5c,
	0x8b, 0x9e, 0x7f, 0xe0, 0x79, 0x69, 0xb8, 0x33, 0x4b, 0xc3, 0x3d, 0x31, 0xe9, 0x08, 0xda, 0x56,
	0xa4, 0x2b, 0x3c, 0x75, 0xc0, 0x35, 0xd7, 0x1d, 0xee, 0xbe, 0x3a, 0xf4, 0xf8, 0x27, 0x80, 0xfb,
	0x2b, 0xc3, 0x5f, 0x57, 0x99, 0x14, 0x36, 0x16, 0x45, 0xc2, 0x3e, 0x55, 0xcc, 0x58, 0xf4, 0x10,
	0xf6, 0x93, 0x52, 0x51, 0x04, 0x60, 0x06, 0xe6, 0x7b, 0xc9, 0xc8, 0x23, 0xcb, 0x02, 0x1d, 0xc2,
	0x5d, 0x23, 0xb8, 0x62, 0xeb, 0xe0, 0xce, 0x0c, 0xcc, 0x47, 0x89, 0xaf, 0xd0, 0x4b, 0x78, 0x5f,
	0xa6, 0x0d, 0xcd, 0x44, 0x41, 0x85, 0xa2, 0x95, 0xd1, 0x52, 0x06, 0xc3, 0x19, 0x98, 0x8f, 0x17,
	0x13, 0xec, 0xa7, 0x71, 0x3b, 0x61, 0xbf, 0x13, 0x7e, 0xa1, 0x85, 0x8a, 0x77, 0x2e, 0xbf, 0x1f,
	0x0d, 0x92, 0x3d, 0x99, 0x36, 0xb1, 0x28, 0x96, 0xea, 0xad, 0x73, 0xa1, 0x53, 0x38, 0x31, 0xe9,
	0x19, 0xa3, 0x56, 0x7f, 0x64, 0x8a, 0x4a, 0xa1, 0x84, 0xac, 0x24, 0x4d, 0xa5, 0xae, 0x94, 0x0d,
	0x76, 0x6e, 0xd7, 0xf2, 0xd0, 0x75, 0x78, 0xe3, 0x1a, 0xac, 0x3a, 0xff, 0xf3, 0xd6, 0xfe, 0x6c,
	0xfc, 0xe5, 0xfa, 0xe2, 0xc4, 0x8f, 0x7c, 0xfc, 0x14, 0x1e, 0xfc, 0x1d, 0x80, 0x29, 0xb5, 0x32,
	0x0c, 0x3d, 0x82, 0xc3, 0xcc, 0xaf, 0x3e, 0x5e, 0xdc, 0xc3, 0x7f, 0xee, 0x86, 0x9d, 0xca, 0x71,
	0x8b, 0x77, 0x70, 0xb8, 0x32, 0x1c, 0xbd, 0x82, 0xa3, 0xdf, 0x76, 0x74, 0x74, 0x53, 0xf9, 0x8f,
	0x64, 0xa7, 0xb3, 0xff, 0x0b, 0xba, 0x3f, 0x4f, 0xef, 0x7e, 0xbe, 0xbe, 0x38, 0x01, 0xf1, 0xf2,
	0x72, 0x13, 0x82, 0xab, 0x4d, 0x08, 0x7e, 0x6c, 0x42, 0xf0, 0x75, 0x1b, 0x0e, 0xae, 0xb6, 0xe1,
	0xe0, 0xdb, 0x36, 0x1c, 0x9c, 0x12, 0x2e, 0xec, 0xfb, 0x2a, 0xc3, 0xb9, 0x96, 0xa4, 0x64, 0x9c,
	0x9f, 0x7f, 0xa8, 0x89, 0x0b, 0x8d, 0x9d, 0x09, 0xb6, 0x26, 0xf5, 0x13, 0xd2, 0xf4, 0xef, 0x89,
	0xd8, 0xf3, 0x92, 0x99, 0x6c, 0xb7, 0x3d, 0xf6, 0xe3, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xdc,
	0x29, 0xba, 0x84, 0x93, 0x02, 0x00, 0x00,
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
	SubmitBid(ctx context.Context, in *MsgSubmitBidRequest, opts ...grpc.CallOption) (*MsgSubmitBidResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SubmitBid(ctx context.Context, in *MsgSubmitBidRequest, opts ...grpc.CallOption) (*MsgSubmitBidResponse, error) {
	out := new(MsgSubmitBidResponse)
	err := c.cc.Invoke(ctx, "/auction.v1.Msg/SubmitBid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	SubmitBid(context.Context, *MsgSubmitBidRequest) (*MsgSubmitBidResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) SubmitBid(ctx context.Context, req *MsgSubmitBidRequest) (*MsgSubmitBidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitBid not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_SubmitBid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSubmitBidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SubmitBid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auction.v1.Msg/SubmitBid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SubmitBid(ctx, req.(*MsgSubmitBidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "auction.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitBid",
			Handler:    _Msg_SubmitBid_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auction/v1/tx.proto",
}

func (m *MsgSubmitBidRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSubmitBidRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSubmitBidRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.SaleTokenMinimumAmount.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size, err := m.MaxBidInUsomm.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Signer) > 0 {
		i -= len(m.Signer)
		copy(dAtA[i:], m.Signer)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signer)))
		i--
		dAtA[i] = 0x12
	}
	if m.AuctionId != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.AuctionId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MsgSubmitBidResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSubmitBidResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSubmitBidResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Bid != nil {
		{
			size, err := m.Bid.MarshalToSizedBuffer(dAtA[:i])
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
func (m *MsgSubmitBidRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AuctionId != 0 {
		n += 1 + sovTx(uint64(m.AuctionId))
	}
	l = len(m.Signer)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.MaxBidInUsomm.Size()
	n += 1 + l + sovTx(uint64(l))
	l = m.SaleTokenMinimumAmount.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgSubmitBidResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Bid != nil {
		l = m.Bid.Size()
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
func (m *MsgSubmitBidRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSubmitBidRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSubmitBidRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuctionId", wireType)
			}
			m.AuctionId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AuctionId |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxBidInUsomm", wireType)
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
			if err := m.MaxBidInUsomm.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SaleTokenMinimumAmount", wireType)
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
			if err := m.SaleTokenMinimumAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *MsgSubmitBidResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSubmitBidResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSubmitBidResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bid", wireType)
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
			if m.Bid == nil {
				m.Bid = &Bid{}
			}
			if err := m.Bid.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
