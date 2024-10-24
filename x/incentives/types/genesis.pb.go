// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: incentives/v1/genesis.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
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

type GenesisState struct {
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_179cfb82d3e2b395, []int{0}
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

// Params incentives parameters
type Params struct {
	// DistributionPerBlock defines the coin to be sent to the distribution module from the community pool every block
	DistributionPerBlock types.Coin `protobuf:"bytes,1,opt,name=distribution_per_block,json=distributionPerBlock,proto3" json:"distribution_per_block"`
	// IncentivesCutoffHeight defines the block height after which the incentives module will stop sending coins to the distribution module from
	// the community pool
	IncentivesCutoffHeight uint64 `protobuf:"varint,2,opt,name=incentives_cutoff_height,json=incentivesCutoffHeight,proto3" json:"incentives_cutoff_height,omitempty"`
	// ValidatorMaxDistributionPerBlock defines the maximum coins to be sent directly to voters in the last block from the community pool every block. Leftover coins remain in the community pool.
	ValidatorMaxDistributionPerBlock types.Coin `protobuf:"bytes,3,opt,name=validator_max_distribution_per_block,json=validatorMaxDistributionPerBlock,proto3" json:"validator_max_distribution_per_block"`
	// ValidatorIncentivesCutoffHeight defines the block height after which the validator incentives will be stopped
	ValidatorIncentivesCutoffHeight uint64 `protobuf:"varint,4,opt,name=validator_incentives_cutoff_height,json=validatorIncentivesCutoffHeight,proto3" json:"validator_incentives_cutoff_height,omitempty"`
	// ValidatorIncentivesMaxFraction defines the maximum fraction of the validator distribution per block that can be sent to a single validator
	ValidatorIncentivesMaxFraction github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=validator_incentives_max_fraction,json=validatorIncentivesMaxFraction,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"validator_incentives_max_fraction"`
	// ValidatorIncentivesSetSizeLimit defines the max number of validators to apportion the validator distribution per block to
	ValidatorIncentivesSetSizeLimit uint64 `protobuf:"varint,6,opt,name=validator_incentives_set_size_limit,json=validatorIncentivesSetSizeLimit,proto3" json:"validator_incentives_set_size_limit,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_179cfb82d3e2b395, []int{1}
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

func (m *Params) GetDistributionPerBlock() types.Coin {
	if m != nil {
		return m.DistributionPerBlock
	}
	return types.Coin{}
}

func (m *Params) GetIncentivesCutoffHeight() uint64 {
	if m != nil {
		return m.IncentivesCutoffHeight
	}
	return 0
}

func (m *Params) GetValidatorMaxDistributionPerBlock() types.Coin {
	if m != nil {
		return m.ValidatorMaxDistributionPerBlock
	}
	return types.Coin{}
}

func (m *Params) GetValidatorIncentivesCutoffHeight() uint64 {
	if m != nil {
		return m.ValidatorIncentivesCutoffHeight
	}
	return 0
}

func (m *Params) GetValidatorIncentivesSetSizeLimit() uint64 {
	if m != nil {
		return m.ValidatorIncentivesSetSizeLimit
	}
	return 0
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "incentives.v1.GenesisState")
	proto.RegisterType((*Params)(nil), "incentives.v1.Params")
}

func init() { proto.RegisterFile("incentives/v1/genesis.proto", fileDescriptor_179cfb82d3e2b395) }

var fileDescriptor_179cfb82d3e2b395 = []byte{
	// 441 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0xcb, 0x6a, 0xdb, 0x40,
	0x14, 0x86, 0xad, 0xd6, 0x35, 0x74, 0xda, 0x6e, 0x44, 0x1a, 0xd4, 0x14, 0x64, 0xd7, 0x2d, 0xc5,
	0x9b, 0xce, 0xe0, 0x64, 0x93, 0xb5, 0x1d, 0x7a, 0xa1, 0x0e, 0x04, 0x9b, 0x6e, 0xba, 0x19, 0x46,
	0xe3, 0x63, 0xf9, 0x34, 0x96, 0x46, 0xcc, 0x8c, 0x85, 0x9d, 0xa7, 0xe8, 0x23, 0xf4, 0x71, 0xb2,
	0xcc, 0xb2, 0x74, 0x11, 0x8a, 0xfd, 0x22, 0x45, 0x17, 0x2c, 0x17, 0x64, 0xc8, 0x4a, 0x82, 0xef,
	0xcc, 0xff, 0x7f, 0x07, 0x69, 0xc8, 0x6b, 0x8c, 0x25, 0xc4, 0x16, 0x53, 0x30, 0x2c, 0xed, 0xb3,
	0x10, 0x62, 0x30, 0x68, 0x68, 0xa2, 0x95, 0x55, 0xee, 0x8b, 0x0a, 0xd2, 0xb4, 0x7f, 0x72, 0x14,
	0xaa, 0x50, 0xe5, 0x84, 0x65, 0x6f, 0xc5, 0xd0, 0x89, 0x2f, 0x95, 0x89, 0x94, 0x61, 0x81, 0x30,
	0xc0, 0xd2, 0x7e, 0x00, 0x56, 0xf4, 0x99, 0x54, 0x18, 0x17, 0xbc, 0x3b, 0x24, 0xcf, 0x3f, 0x15,
	0xa9, 0x13, 0x2b, 0x2c, 0xb8, 0x67, 0xa4, 0x95, 0x08, 0x2d, 0x22, 0xe3, 0x39, 0x1d, 0xa7, 0xf7,
	0xec, 0xf4, 0x25, 0xfd, 0xaf, 0x85, 0x5e, 0xe5, 0x70, 0xd0, 0xbc, 0xbd, 0x6f, 0x37, 0xc6, 0xe5,
	0x68, 0xf7, 0x57, 0x93, 0xb4, 0x0a, 0xe0, 0x7e, 0x23, 0xc7, 0x53, 0x34, 0x56, 0x63, 0xb0, 0xb4,
	0xa8, 0x62, 0x9e, 0x80, 0xe6, 0xc1, 0x42, 0xc9, 0xeb, 0x32, 0xef, 0x15, 0x2d, 0x84, 0x68, 0x26,
	0x44, 0x4b, 0x21, 0x3a, 0x54, 0x18, 0x97, 0x99, 0x47, 0xfb, 0xc7, 0xaf, 0x40, 0x0f, 0xb2, 0xc3,
	0xee, 0x39, 0xf1, 0x2a, 0x0f, 0x2e, 0x97, 0x56, 0xcd, 0x66, 0x7c, 0x0e, 0x18, 0xce, 0xad, 0xf7,
	0xa8, 0xe3, 0xf4, 0x9a, 0xe3, 0xe3, 0x8a, 0x0f, 0x73, 0xfc, 0x39, 0xa7, 0xae, 0x22, 0xef, 0x52,
	0xb1, 0xc0, 0xa9, 0xb0, 0x4a, 0xf3, 0x48, 0xac, 0xf8, 0x01, 0xbd, 0xc7, 0x0f, 0xd3, 0xeb, 0xec,
	0xc2, 0x2e, 0xc5, 0xea, 0xa2, 0x4e, 0xf5, 0x2b, 0xe9, 0x56, 0x85, 0x07, 0xa5, 0x9b, 0xb9, 0x74,
	0x7b, 0x37, 0xf9, 0xa5, 0xde, 0x7e, 0x4d, 0xde, 0xd4, 0x86, 0x65, 0x8b, 0xcc, 0xb4, 0x90, 0x59,
	0xb3, 0xf7, 0xa4, 0xe3, 0xf4, 0x9e, 0x0e, 0x68, 0xe6, 0xf7, 0xe7, 0xbe, 0xfd, 0x3e, 0x44, 0x3b,
	0x5f, 0x06, 0x54, 0xaa, 0x88, 0x95, 0x1f, 0xbf, 0x78, 0x7c, 0x30, 0xd3, 0x6b, 0x66, 0xd7, 0x09,
	0x18, 0x7a, 0x01, 0x72, 0xec, 0xd7, 0x74, 0x5f, 0x8a, 0xd5, 0xc7, 0x32, 0xd5, 0x1d, 0x91, 0xb7,
	0xb5, 0xd5, 0x06, 0x2c, 0x37, 0x78, 0x03, 0x7c, 0x81, 0x11, 0x5a, 0xaf, 0x75, 0x70, 0x91, 0x09,
	0xd8, 0x09, 0xde, 0xc0, 0x28, 0x1b, 0x1b, 0x8c, 0x6e, 0x37, 0xbe, 0x73, 0xb7, 0xf1, 0x9d, 0xbf,
	0x1b, 0xdf, 0xf9, 0xb9, 0xf5, 0x1b, 0x77, 0x5b, 0xbf, 0xf1, 0x7b, 0xeb, 0x37, 0xbe, 0x9f, 0xee,
	0xf9, 0x26, 0x10, 0x86, 0xeb, 0x1f, 0x29, 0x33, 0x2a, 0x8a, 0x60, 0x81, 0xa0, 0x59, 0x7a, 0xce,
	0x56, 0x6c, 0xef, 0x16, 0xe4, 0xfe, 0x41, 0x2b, 0xff, 0x79, 0xcf, 0xfe, 0x05, 0x00, 0x00, 0xff,
	0xff, 0x89, 0x15, 0x21, 0xdd, 0x20, 0x03, 0x00, 0x00,
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
	if m.ValidatorIncentivesSetSizeLimit != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.ValidatorIncentivesSetSizeLimit))
		i--
		dAtA[i] = 0x30
	}
	{
		size := m.ValidatorIncentivesMaxFraction.Size()
		i -= size
		if _, err := m.ValidatorIncentivesMaxFraction.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.ValidatorIncentivesCutoffHeight != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.ValidatorIncentivesCutoffHeight))
		i--
		dAtA[i] = 0x20
	}
	{
		size, err := m.ValidatorMaxDistributionPerBlock.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.IncentivesCutoffHeight != 0 {
		i = encodeVarintGenesis(dAtA, i, uint64(m.IncentivesCutoffHeight))
		i--
		dAtA[i] = 0x10
	}
	{
		size, err := m.DistributionPerBlock.MarshalToSizedBuffer(dAtA[:i])
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
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.DistributionPerBlock.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.IncentivesCutoffHeight != 0 {
		n += 1 + sovGenesis(uint64(m.IncentivesCutoffHeight))
	}
	l = m.ValidatorMaxDistributionPerBlock.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.ValidatorIncentivesCutoffHeight != 0 {
		n += 1 + sovGenesis(uint64(m.ValidatorIncentivesCutoffHeight))
	}
	l = m.ValidatorIncentivesMaxFraction.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if m.ValidatorIncentivesSetSizeLimit != 0 {
		n += 1 + sovGenesis(uint64(m.ValidatorIncentivesSetSizeLimit))
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
func (m *Params) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DistributionPerBlock", wireType)
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
			if err := m.DistributionPerBlock.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncentivesCutoffHeight", wireType)
			}
			m.IncentivesCutoffHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IncentivesCutoffHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorMaxDistributionPerBlock", wireType)
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
			if err := m.ValidatorMaxDistributionPerBlock.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorIncentivesCutoffHeight", wireType)
			}
			m.ValidatorIncentivesCutoffHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ValidatorIncentivesCutoffHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorIncentivesMaxFraction", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ValidatorIncentivesMaxFraction.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorIncentivesSetSizeLimit", wireType)
			}
			m.ValidatorIncentivesSetSizeLimit = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ValidatorIncentivesSetSizeLimit |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
