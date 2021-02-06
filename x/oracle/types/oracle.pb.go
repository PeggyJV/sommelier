// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: oracle/v1/oracle.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_tendermint_tendermint_libs_bytes "github.com/tendermint/tendermint/libs/bytes"
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

// UniswapData represents the necessary data for a given uniswap pair.
type UniswapData struct {
	Id          string                                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Reserve0    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=reserve0,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"reserve0" yaml:"reserve0"`
	Reserve1    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=reserve1,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"reserve1" yaml:"reserve1"`
	ReserveUsd  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=reserve_usd,json=reserveUsd,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"reserveUSD" yaml:"reserveUSD"`
	Token0      UniswapToken                           `protobuf:"bytes,5,opt,name=token0,proto3" json:"token0"`
	Token1      UniswapToken                           `protobuf:"bytes,6,opt,name=token1,proto3" json:"token1"`
	Token0Price github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=token0_price,json=token0Price,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"token0Price" yaml:"token0Price"`
	Token1Price github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,8,opt,name=token1_price,json=token1Price,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"token1Price" yaml:"token1Price"`
	TotalSupply github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,9,opt,name=total_supply,json=totalSupply,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"totalSupply" yaml:"totalSupply"`
}

func (m *OracleFeed) Reset()         { *m = OracleFeed{} }
func (m *OracleFeed) String() string { return proto.CompactTextString(m) }
func (*OracleFeed) ProtoMessage()    {}
func (*OracleFeed) Descriptor() ([]byte, []int) {
	return fileDescriptor_652b57db11528d07, []int{0}
}
func (m *OracleFeed) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OracleFeed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OracleFeed.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OracleFeed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OracleFeed.Merge(m, src)
}
func (m *OracleFeed) XXX_Size() int {
	return m.Size()
}
func (m *OracleFeed) XXX_DiscardUnknown() {
	xxx_messageInfo_OracleFeed.DiscardUnknown(m)
}

var xxx_messageInfo_OracleFeed proto.InternalMessageInfo

func (m *UniswapData) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UniswapData) GetToken0() UniswapToken {
	if m != nil {
		return m.Token0
	}
	return UniswapToken{}
}

func (m *UniswapData) GetToken1() UniswapToken {
	if m != nil {
		return m.Token1
	}
	return UniswapToken{}
}

// UniswapToken is the returned uniswap token representation
type UniswapToken struct {
	// TODO: what's the ID?
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// number of decimal positions of the pair token
	Decimals uint64 `protobuf:"varint,2,opt,name=decimals,proto3" json:"decimals,omitempty"`
}

func (m *UniswapToken) Reset()         { *m = UniswapToken{} }
func (m *UniswapToken) String() string { return proto.CompactTextString(m) }
func (*UniswapToken) ProtoMessage()    {}
func (*UniswapToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_652b57db11528d07, []int{1}
}
func (m *UniswapToken) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UniswapToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UniswapToken.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UniswapToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UniswapToken.Merge(m, src)
}
func (m *UniswapToken) XXX_Size() int {
	return m.Size()
}
func (m *UniswapToken) XXX_DiscardUnknown() {
	xxx_messageInfo_UniswapToken.DiscardUnknown(m)
}

var xxx_messageInfo_UniswapToken proto.InternalMessageInfo

func (m *UniswapToken) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *UniswapToken) GetDecimals() uint64 {
	if m != nil {
		return m.Decimals
	}
	return 0
}

func (m *OracleVote) GetFeed() *OracleFeed {
	if m != nil {
		return m.Feed
	}
	return nil
}

// OraclePrevote defines an array of hashed from oracle data that are used
// for the prevote phase of the oracle data feeding.
type OraclePrevote struct {
	// hex formated hash of an oracle feed
	Hash github_com_tendermint_tendermint_libs_bytes.HexBytes `protobuf:"bytes,1,opt,name=hash,proto3,casttype=github.com/tendermint/tendermint/libs/bytes.HexBytes" json:"hash,omitempty"`
}

func (m *OraclePrevote) Reset()         { *m = OraclePrevote{} }
func (m *OraclePrevote) String() string { return proto.CompactTextString(m) }
func (*OraclePrevote) ProtoMessage()    {}
func (*OraclePrevote) Descriptor() ([]byte, []int) {
	return fileDescriptor_652b57db11528d07, []int{4}
}
func (m *OraclePrevote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OraclePrevote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OraclePrevote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OraclePrevote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OraclePrevote.Merge(m, src)
}
func (m *OraclePrevote) XXX_Size() int {
	return m.Size()
}
func (m *OraclePrevote) XXX_DiscardUnknown() {
	xxx_messageInfo_OraclePrevote.DiscardUnknown(m)
}

var xxx_messageInfo_OraclePrevote proto.InternalMessageInfo

func (m *OraclePrevote) GetHash() github_com_tendermint_tendermint_libs_bytes.HexBytes {
	if m != nil {
		return m.Hash
	}
	return nil
}

func init() {
	proto.RegisterType((*UniswapData)(nil), "oracle.v1.UniswapData")
	proto.RegisterType((*UniswapToken)(nil), "oracle.v1.UniswapToken")
	proto.RegisterType((*OracleVote)(nil), "oracle.v1.OracleVote")
	proto.RegisterType((*OraclePrevote)(nil), "oracle.v1.OraclePrevote")
}

func init() { proto.RegisterFile("oracle/v1/oracle.proto", fileDescriptor_652b57db11528d07) }

var fileDescriptor_652b57db11528d07 = []byte{
	// 441 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x93, 0x41, 0x8f, 0x93, 0x40,
	0x14, 0xc7, 0xa1, 0xae, 0xd8, 0x0e, 0x9b, 0x35, 0x12, 0xa3, 0x93, 0x3d, 0xc0, 0x86, 0x83, 0xd9,
	0xc4, 0x08, 0x3b, 0x7a, 0xdb, 0x8b, 0x49, 0x83, 0x47, 0x13, 0x05, 0x7b, 0xf1, 0xd2, 0xcc, 0xc2,
	0x04, 0xc7, 0x85, 0x1d, 0xc2, 0x50, 0x94, 0x6f, 0xe1, 0xc7, 0xea, 0xb1, 0xde, 0x8c, 0x07, 0x62,
	0xda, 0x5b, 0x8f, 0xfd, 0x04, 0x06, 0x18, 0xa7, 0x34, 0xf5, 0xd2, 0xed, 0x69, 0xde, 0xfb, 0xcf,
	0xcb, 0xff, 0xf7, 0x5e, 0x5e, 0x1e, 0x78, 0xc6, 0x72, 0x1c, 0x26, 0xc4, 0x2d, 0x91, 0xdb, 0x45,
	0x4e, 0x96, 0xb3, 0x82, 0x19, 0x23, 0x91, 0x95, 0xe8, 0xfc, 0x69, 0xcc, 0x62, 0xd6, 0xaa, 0x6e,
	0x13, 0x75, 0x05, 0xf6, 0x4f, 0x0d, 0xe8, 0x93, 0x3b, 0xca, 0xbf, 0xe1, 0xcc, 0xc3, 0x05, 0x36,
	0xce, 0xc0, 0x80, 0x46, 0x50, 0xbd, 0x50, 0x2f, 0x47, 0xfe, 0x80, 0x46, 0x06, 0x05, 0xc3, 0x9c,
	0x70, 0x92, 0x97, 0xe4, 0x0a, 0x0e, 0x1a, 0x75, 0xfc, 0x7e, 0x5e, 0x5b, 0xca, 0xef, 0xda, 0x7a,
	0x11, 0xd3, 0xe2, 0xcb, 0xec, 0xc6, 0x09, 0x59, 0xea, 0x86, 0x8c, 0xa7, 0x8c, 0x8b, 0xe7, 0x15,
	0x8f, 0x6e, 0xdd, 0xa2, 0xca, 0x08, 0x77, 0x3c, 0x12, 0xae, 0x6b, 0x4b, 0x3a, 0x6c, 0x6a, 0xeb,
	0x71, 0x85, 0xd3, 0xe4, 0xda, 0xfe, 0xa7, 0xd8, 0xbe, 0xfc, 0xec, 0xa1, 0x10, 0x7c, 0x70, 0x24,
	0x0a, 0xed, 0xa1, 0xd0, 0x16, 0x85, 0x8c, 0x1c, 0xe8, 0x22, 0x9e, 0xce, 0x78, 0x04, 0x4f, 0x5a,
	0xda, 0xc7, 0x83, 0x69, 0x40, 0x98, 0x4c, 0x02, 0x6f, 0x53, 0x5b, 0x4f, 0x76, 0x78, 0x93, 0xc0,
	0xb3, 0x7d, 0x59, 0xc0, 0x23, 0xe3, 0x2d, 0xd0, 0x0a, 0x76, 0x4b, 0xee, 0xae, 0xe0, 0xc3, 0x0b,
	0xf5, 0x52, 0x7f, 0xfd, 0xdc, 0x91, 0xbb, 0x71, 0xc4, 0x06, 0x3e, 0x35, 0xff, 0xe3, 0xb3, 0xa6,
	0x8f, 0x75, 0x6d, 0x89, 0x72, 0x5f, 0xbc, 0xd2, 0x00, 0x41, 0xed, 0x10, 0x03, 0x24, 0x0c, 0x90,
	0x51, 0x82, 0xd3, 0xce, 0x6a, 0x9a, 0xe5, 0x34, 0x24, 0xf0, 0x51, 0x3b, 0x76, 0x70, 0xf0, 0xd8,
	0x7a, 0xe7, 0xf2, 0xa1, 0x31, 0xd9, 0xd4, 0x96, 0xd1, 0xcd, 0xdd, 0x13, 0x6d, 0xbf, 0x5f, 0x22,
	0xb9, 0x48, 0x70, 0x87, 0x47, 0x71, 0xd1, 0xff, 0xb8, 0x68, 0x87, 0x8b, 0x7a, 0xdc, 0x02, 0x27,
	0x53, 0x3e, 0xcb, 0xb2, 0xa4, 0x82, 0xa3, 0xfb, 0x73, 0x0b, 0x9c, 0x04, 0xad, 0x49, 0x9f, 0x2b,
	0xc5, 0x96, 0xbb, 0xcd, 0xae, 0xc1, 0x69, 0x7f, 0x1f, 0x7b, 0x37, 0x75, 0x0e, 0x86, 0x11, 0x09,
	0x69, 0x8a, 0x13, 0xde, 0xde, 0xd4, 0x89, 0x2f, 0xf3, 0xf1, 0xbb, 0xf9, 0xd2, 0x54, 0x17, 0x4b,
	0x53, 0xfd, 0xb3, 0x34, 0xd5, 0x1f, 0x2b, 0x53, 0x59, 0xac, 0x4c, 0xe5, 0xd7, 0xca, 0x54, 0x3e,
	0xbf, 0xec, 0xf5, 0x9b, 0x91, 0x38, 0xae, 0xbe, 0x96, 0x2e, 0x67, 0x69, 0x4a, 0x12, 0x4a, 0x72,
	0xf7, 0xbb, 0xb8, 0xfb, 0xae, 0xf1, 0x1b, 0xad, 0xbd, 0xee, 0x37, 0x7f, 0x03, 0x00, 0x00, 0xff,
	0xff, 0xb0, 0x57, 0x11, 0xba, 0x18, 0x04, 0x00, 0x00,
}

func (m *OracleFeed) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OracleFeed) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OracleFeed) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.TotalSupply.Size()
		i -= size
		if _, err := m.TotalSupply.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.Token1Price.Size()
		i -= size
		if _, err := m.Token1Price.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.Token0Price.Size()
		i -= size
		if _, err := m.Token0Price.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size, err := m.Token1.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size, err := m.Token0.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.ReserveUsd.Size()
		i -= size
		if _, err := m.ReserveUsd.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.Reserve1.Size()
		i -= size
		if _, err := m.Reserve1.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Reserve0.Size()
		i -= size
		if _, err := m.Reserve0.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintOracle(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *UniswapToken) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UniswapToken) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UniswapToken) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Decimals != 0 {
		i = encodeVarintOracle(dAtA, i, uint64(m.Decimals))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintOracle(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintOracle(dAtA []byte, offset int, v uint64) int {
	offset -= sovOracle(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *OracleFeed) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovOracle(uint64(l))
	}
	l = m.Reserve0.Size()
	n += 1 + l + sovOracle(uint64(l))
	l = m.Reserve1.Size()
	n += 1 + l + sovOracle(uint64(l))
	l = m.ReserveUsd.Size()
	n += 1 + l + sovOracle(uint64(l))
	l = m.Token0.Size()
	n += 1 + l + sovOracle(uint64(l))
	l = m.Token1.Size()
	n += 1 + l + sovOracle(uint64(l))
	l = m.Token0Price.Size()
	n += 1 + l + sovOracle(uint64(l))
	l = m.Token1Price.Size()
	n += 1 + l + sovOracle(uint64(l))
	l = m.TotalSupply.Size()
	n += 1 + l + sovOracle(uint64(l))
	return n
}

func (m *OracleVote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovOracle(uint64(l))
	}
	if m.Decimals != 0 {
		n += 1 + sovOracle(uint64(m.Decimals))
	}
	return n
}

func sovOracle(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozOracle(x uint64) (n int) {
	return sovOracle(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *OracleFeed) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOracle
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
			return fmt.Errorf("proto: OracleFeed: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OracleFeed: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
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
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Reserve0", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
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
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Reserve0.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Reserve1", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
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
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Reserve1.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReserveUSD", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
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
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ReserveUsd.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token0", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
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
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Token0.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token1", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
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
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Token1.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token0Price", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
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
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Token0Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token1Price", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
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
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Token1Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalSupply", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
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
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalSupply.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOracle(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOracle
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthOracle
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
func (m *UniswapToken) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOracle
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
			return fmt.Errorf("proto: UniswapToken: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UniswapToken: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
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
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Decimals", wireType)
			}
			m.Decimals = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Decimals |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipOracle(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthOracle
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthOracle
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
func skipOracle(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowOracle
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
					return 0, ErrIntOverflowOracle
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
					return 0, ErrIntOverflowOracle
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
				return 0, ErrInvalidLengthOracle
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupOracle
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthOracle
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthOracle        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowOracle          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupOracle = fmt.Errorf("proto: unexpected end of group")
)
