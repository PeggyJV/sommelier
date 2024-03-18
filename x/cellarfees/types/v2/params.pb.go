// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cellarfees/v2/params.proto

package v2

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
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

// Params defines the parameters for the module.
type Params struct {
	// Emission rate factor. Specifically, the number of blocks over which to distribute
	// some amount of staking rewards.
	RewardEmissionPeriod uint64 `protobuf:"varint,2,opt,name=reward_emission_period,json=rewardEmissionPeriod,proto3" json:"reward_emission_period,omitempty"`
	// The initial rate at which auctions should decrease their denom's price in SOMM
	InitialPriceDecreaseRate github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=initial_price_decrease_rate,json=initialPriceDecreaseRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"initial_price_decrease_rate"`
	// Number of blocks between auction price decreases
	PriceDecreaseBlockInterval uint64 `protobuf:"varint,4,opt,name=price_decrease_block_interval,json=priceDecreaseBlockInterval,proto3" json:"price_decrease_block_interval,omitempty"`
	// The interval between starting auctions
	AuctionInterval uint64 `protobuf:"varint,5,opt,name=auction_interval,json=auctionInterval,proto3" json:"auction_interval,omitempty"`
	// A fee token's total USD value threshold, based on it's auction.v1.TokenPrice, above which an auction is triggered
	AuctionThresholdUsdValue github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=auction_threshold_usd_value,json=auctionThresholdUsdValue,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"auction_threshold_usd_value"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_270f377f75209ba6, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetRewardEmissionPeriod() uint64 {
	if m != nil {
		return m.RewardEmissionPeriod
	}
	return 0
}

func (m *Params) GetPriceDecreaseBlockInterval() uint64 {
	if m != nil {
		return m.PriceDecreaseBlockInterval
	}
	return 0
}

func (m *Params) GetAuctionInterval() uint64 {
	if m != nil {
		return m.AuctionInterval
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "cellarfees.v2.Params")
}

func init() { proto.RegisterFile("cellarfees/v2/params.proto", fileDescriptor_270f377f75209ba6) }

var fileDescriptor_270f377f75209ba6 = []byte{
	// 381 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x41, 0x6b, 0xdb, 0x30,
	0x1c, 0xc5, 0xed, 0xc4, 0x0b, 0x9b, 0x61, 0x2c, 0x98, 0x30, 0x4c, 0xc2, 0x9c, 0xb0, 0xc3, 0xc8,
	0x0e, 0xb3, 0x20, 0xdb, 0x18, 0xec, 0xb6, 0x90, 0x1d, 0xb6, 0x4b, 0x43, 0x68, 0x7b, 0xe8, 0x45,
	0x28, 0xf2, 0xbf, 0x8e, 0x1a, 0x29, 0x32, 0x92, 0xec, 0x36, 0xdf, 0xa2, 0xa7, 0xd2, 0x63, 0x3f,
	0x4e, 0x8e, 0x39, 0x96, 0x1e, 0x42, 0x49, 0xbe, 0x48, 0xb1, 0xe3, 0xa4, 0xa1, 0xc7, 0x9e, 0x24,
	0xf4, 0xde, 0xfb, 0x3d, 0xfd, 0x91, 0xdc, 0x26, 0x05, 0xce, 0x89, 0x3a, 0x07, 0xd0, 0x28, 0xeb,
	0xa1, 0x84, 0x28, 0x22, 0x74, 0x98, 0x28, 0x69, 0xa4, 0xf7, 0xfe, 0x59, 0x0b, 0xb3, 0x5e, 0xb3,
	0x11, 0xcb, 0x58, 0x16, 0x0a, 0xca, 0x77, 0x5b, 0xd3, 0xe7, 0x9b, 0xaa, 0x5b, 0x1b, 0x16, 0x29,
	0xef, 0x87, 0xfb, 0x51, 0xc1, 0x25, 0x51, 0x11, 0x06, 0xc1, 0xb4, 0x66, 0x72, 0x86, 0x13, 0x50,
	0x4c, 0x46, 0x7e, 0xa5, 0x63, 0x77, 0x9d, 0x51, 0x63, 0xab, 0xfe, 0x2d, 0xc5, 0x61, 0xa1, 0x79,
	0xc2, 0x6d, 0xb1, 0x19, 0x33, 0x8c, 0x70, 0x9c, 0x28, 0x46, 0x01, 0x47, 0x40, 0x15, 0x10, 0x0d,
	0x58, 0x11, 0x03, 0x7e, 0xb5, 0x63, 0x77, 0xdf, 0xf5, 0xc3, 0xc5, 0xaa, 0x6d, 0x3d, 0xac, 0xda,
	0x5f, 0x62, 0x66, 0x26, 0xe9, 0x38, 0xa4, 0x52, 0x20, 0x2a, 0xb5, 0x90, 0xba, 0x5c, 0xbe, 0xe9,
	0x68, 0x8a, 0xcc, 0x3c, 0x01, 0x1d, 0x0e, 0x80, 0x8e, 0xfc, 0x12, 0x39, 0xcc, 0x89, 0x83, 0x12,
	0x38, 0x22, 0x06, 0xbc, 0x3f, 0xee, 0xa7, 0x17, 0x35, 0x63, 0x2e, 0xe9, 0x14, 0xb3, 0x99, 0x01,
	0x95, 0x11, 0xee, 0x3b, 0xc5, 0x5d, 0x9b, 0xc9, 0x61, 0xb2, 0x9f, 0x5b, 0xfe, 0x95, 0x0e, 0xef,
	0xab, 0x5b, 0x27, 0x29, 0x35, 0xf9, 0x7c, 0xfb, 0xd4, 0x9b, 0x22, 0xf5, 0xa1, 0x3c, 0xdf, 0x5b,
	0x85, 0xdb, 0xda, 0x59, 0xcd, 0x44, 0x81, 0x9e, 0x48, 0x1e, 0xe1, 0x54, 0x47, 0x38, 0x23, 0x3c,
	0x05, 0xbf, 0xf6, 0xba, 0xe1, 0x4a, 0xe4, 0xf1, 0x8e, 0x78, 0xa2, 0xa3, 0xd3, 0x9c, 0xf7, 0xdb,
	0xb9, 0xbd, 0x6b, 0x5b, 0xff, 0x9d, 0xb7, 0x76, 0xbd, 0xd2, 0x3f, 0x5a, 0xac, 0x03, 0x7b, 0xb9,
	0x0e, 0xec, 0xc7, 0x75, 0x60, 0x5f, 0x6f, 0x02, 0x6b, 0xb9, 0x09, 0xac, 0xfb, 0x4d, 0x60, 0x9d,
	0xfd, 0x3c, 0xe8, 0x49, 0x20, 0x8e, 0xe7, 0x17, 0x19, 0xd2, 0x52, 0x08, 0xe0, 0x0c, 0x14, 0xca,
	0x7e, 0xa1, 0x2b, 0x74, 0xf0, 0x2b, 0x8a, 0x5e, 0x94, 0xf5, 0xc6, 0xb5, 0xe2, 0xc1, 0xbf, 0x3f,
	0x05, 0x00, 0x00, 0xff, 0xff, 0x87, 0x16, 0xc6, 0xed, 0x33, 0x02, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.AuctionThresholdUsdValue.Size()
		i -= size
		if _, err := m.AuctionThresholdUsdValue.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if m.AuctionInterval != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.AuctionInterval))
		i--
		dAtA[i] = 0x28
	}
	if m.PriceDecreaseBlockInterval != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.PriceDecreaseBlockInterval))
		i--
		dAtA[i] = 0x20
	}
	{
		size := m.InitialPriceDecreaseRate.Size()
		i -= size
		if _, err := m.InitialPriceDecreaseRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.RewardEmissionPeriod != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.RewardEmissionPeriod))
		i--
		dAtA[i] = 0x10
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.RewardEmissionPeriod != 0 {
		n += 1 + sovParams(uint64(m.RewardEmissionPeriod))
	}
	l = m.InitialPriceDecreaseRate.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.PriceDecreaseBlockInterval != 0 {
		n += 1 + sovParams(uint64(m.PriceDecreaseBlockInterval))
	}
	if m.AuctionInterval != 0 {
		n += 1 + sovParams(uint64(m.AuctionInterval))
	}
	l = m.AuctionThresholdUsdValue.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardEmissionPeriod", wireType)
			}
			m.RewardEmissionPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RewardEmissionPeriod |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InitialPriceDecreaseRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.InitialPriceDecreaseRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PriceDecreaseBlockInterval", wireType)
			}
			m.PriceDecreaseBlockInterval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PriceDecreaseBlockInterval |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuctionInterval", wireType)
			}
			m.AuctionInterval = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AuctionInterval |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AuctionThresholdUsdValue", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AuctionThresholdUsdValue.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
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
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
					return 0, ErrIntOverflowParams
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
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)
