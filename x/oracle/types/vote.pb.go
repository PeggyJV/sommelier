// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: oracle/v1/vote.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
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

// ExchangeRatePrevote - struct to store a validator's prevote on the rate of Luna in the denom asset
// type ExchangeRatePrevote struct {
// 	Hash        VoteHash       `json:"hash"`  // Vote hex hash to protect centralize data source problem
// 	Denom       string         `json:"denom"` // Ticker name of target fiat currency
// 	Voter       sdk.ValAddress `json:"voter"` // Voter val address
// 	SubmitBlock int64          `json:"submit_block"`
// }
type ExchangeRatePrevote struct {
	Hash        []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Denom       string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
	Voter       string `protobuf:"bytes,3,opt,name=voter,proto3" json:"voter,omitempty"`
	SubmitBlock int64  `protobuf:"varint,4,opt,name=submit_block,json=submitBlock,proto3" json:"submit_block,omitempty"`
}

func (m *ExchangeRatePrevote) Reset()         { *m = ExchangeRatePrevote{} }
func (m *ExchangeRatePrevote) String() string { return proto.CompactTextString(m) }
func (*ExchangeRatePrevote) ProtoMessage()    {}
func (*ExchangeRatePrevote) Descriptor() ([]byte, []int) {
	return fileDescriptor_92bdd9e640ece7e2, []int{0}
}
func (m *ExchangeRatePrevote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExchangeRatePrevote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExchangeRatePrevote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExchangeRatePrevote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExchangeRatePrevote.Merge(m, src)
}
func (m *ExchangeRatePrevote) XXX_Size() int {
	return m.Size()
}
func (m *ExchangeRatePrevote) XXX_DiscardUnknown() {
	xxx_messageInfo_ExchangeRatePrevote.DiscardUnknown(m)
}

var xxx_messageInfo_ExchangeRatePrevote proto.InternalMessageInfo

func (m *ExchangeRatePrevote) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *ExchangeRatePrevote) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *ExchangeRatePrevote) GetVoter() string {
	if m != nil {
		return m.Voter
	}
	return ""
}

func (m *ExchangeRatePrevote) GetSubmitBlock() int64 {
	if m != nil {
		return m.SubmitBlock
	}
	return 0
}

// ExchangeRateVote - struct to store a validator's vote on the rate of Luna in the denom asset
// type ExchangeRateVote struct {
// 	ExchangeRate sdk.Dec        `json:"exchange_rate"` // ExchangeRate of Luna in target fiat currency
// 	Denom        string         `json:"denom"`         // Ticker name of target fiat currency
// 	Voter        sdk.ValAddress `json:"voter"`         // voter val address of validator
// }
type ExchangeRateVote struct {
	ExchangeRate github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,1,opt,name=exchange_rate,json=exchangeRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"exchange_rate"`
	Denom        string                                 `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
	Voter        string                                 `protobuf:"bytes,3,opt,name=voter,proto3" json:"voter,omitempty"`
}

func (m *ExchangeRateVote) Reset()         { *m = ExchangeRateVote{} }
func (m *ExchangeRateVote) String() string { return proto.CompactTextString(m) }
func (*ExchangeRateVote) ProtoMessage()    {}
func (*ExchangeRateVote) Descriptor() ([]byte, []int) {
	return fileDescriptor_92bdd9e640ece7e2, []int{1}
}
func (m *ExchangeRateVote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExchangeRateVote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExchangeRateVote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExchangeRateVote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExchangeRateVote.Merge(m, src)
}
func (m *ExchangeRateVote) XXX_Size() int {
	return m.Size()
}
func (m *ExchangeRateVote) XXX_DiscardUnknown() {
	xxx_messageInfo_ExchangeRateVote.DiscardUnknown(m)
}

var xxx_messageInfo_ExchangeRateVote proto.InternalMessageInfo

func (m *ExchangeRateVote) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *ExchangeRateVote) GetVoter() string {
	if m != nil {
		return m.Voter
	}
	return ""
}

// AggregateExchangeRatePrevote - struct to store a validator's aggregate prevote on the rate of Luna in the denom asset
// type AggregateExchangeRatePrevote struct {
// 	Hash        AggregateVoteHash `json:"hash"`  // Vote hex hash to protect centralize data source problem
// 	Voter       sdk.ValAddress    `json:"voter"` // Voter val address
// 	SubmitBlock int64             `json:"submit_block"`
// }
type AggregateExchangeRatePrevote struct {
	Hash        []byte `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	Voter       string `protobuf:"bytes,2,opt,name=voter,proto3" json:"voter,omitempty"`
	SubmitBlock int64  `protobuf:"varint,3,opt,name=submit_block,json=submitBlock,proto3" json:"submit_block,omitempty"`
}

func (m *AggregateExchangeRatePrevote) Reset()         { *m = AggregateExchangeRatePrevote{} }
func (m *AggregateExchangeRatePrevote) String() string { return proto.CompactTextString(m) }
func (*AggregateExchangeRatePrevote) ProtoMessage()    {}
func (*AggregateExchangeRatePrevote) Descriptor() ([]byte, []int) {
	return fileDescriptor_92bdd9e640ece7e2, []int{2}
}
func (m *AggregateExchangeRatePrevote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AggregateExchangeRatePrevote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AggregateExchangeRatePrevote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AggregateExchangeRatePrevote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AggregateExchangeRatePrevote.Merge(m, src)
}
func (m *AggregateExchangeRatePrevote) XXX_Size() int {
	return m.Size()
}
func (m *AggregateExchangeRatePrevote) XXX_DiscardUnknown() {
	xxx_messageInfo_AggregateExchangeRatePrevote.DiscardUnknown(m)
}

var xxx_messageInfo_AggregateExchangeRatePrevote proto.InternalMessageInfo

func (m *AggregateExchangeRatePrevote) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *AggregateExchangeRatePrevote) GetVoter() string {
	if m != nil {
		return m.Voter
	}
	return ""
}

func (m *AggregateExchangeRatePrevote) GetSubmitBlock() int64 {
	if m != nil {
		return m.SubmitBlock
	}
	return 0
}

// ExchangeRateTuple - struct to represent a exchange rate of Luna in the denom asset
// type ExchangeRateTuple struct {
// 	Denom        string  `json:"denom"`
// 	ExchangeRate sdk.Dec `json:"exchange_rate"`
// }
type ExchangeRateTuple struct {
	Denom        string                                 `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	ExchangeRate github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=exchange_rate,json=exchangeRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"exchange_rate"`
}

func (m *ExchangeRateTuple) Reset()         { *m = ExchangeRateTuple{} }
func (m *ExchangeRateTuple) String() string { return proto.CompactTextString(m) }
func (*ExchangeRateTuple) ProtoMessage()    {}
func (*ExchangeRateTuple) Descriptor() ([]byte, []int) {
	return fileDescriptor_92bdd9e640ece7e2, []int{3}
}
func (m *ExchangeRateTuple) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExchangeRateTuple) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExchangeRateTuple.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExchangeRateTuple) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExchangeRateTuple.Merge(m, src)
}
func (m *ExchangeRateTuple) XXX_Size() int {
	return m.Size()
}
func (m *ExchangeRateTuple) XXX_DiscardUnknown() {
	xxx_messageInfo_ExchangeRateTuple.DiscardUnknown(m)
}

var xxx_messageInfo_ExchangeRateTuple proto.InternalMessageInfo

func (m *ExchangeRateTuple) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

// AggregateExchangeRateVote - struct to store a validator's aggregate vote on the rate of Luna in the denom asset
// type AggregateExchangeRateVote struct {
// 	ExchangeRateTuples ExchangeRateTuples `json:"exchange_rate_tuples"` // ExchangeRates of Luna in target fiat currencies
// 	Voter              sdk.ValAddress     `json:"voter"`                // voter val address of validator
// }
type AggregateExchangeRateVote struct {
	ExchangeRateTuples []ExchangeRateTuple `protobuf:"bytes,1,rep,name=exchange_rate_tuples,json=exchangeRateTuples,proto3" json:"exchange_rate_tuples"`
	Voter              string              `protobuf:"bytes,2,opt,name=voter,proto3" json:"voter,omitempty"`
}

func (m *AggregateExchangeRateVote) Reset()         { *m = AggregateExchangeRateVote{} }
func (m *AggregateExchangeRateVote) String() string { return proto.CompactTextString(m) }
func (*AggregateExchangeRateVote) ProtoMessage()    {}
func (*AggregateExchangeRateVote) Descriptor() ([]byte, []int) {
	return fileDescriptor_92bdd9e640ece7e2, []int{4}
}
func (m *AggregateExchangeRateVote) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AggregateExchangeRateVote) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AggregateExchangeRateVote.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AggregateExchangeRateVote) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AggregateExchangeRateVote.Merge(m, src)
}
func (m *AggregateExchangeRateVote) XXX_Size() int {
	return m.Size()
}
func (m *AggregateExchangeRateVote) XXX_DiscardUnknown() {
	xxx_messageInfo_AggregateExchangeRateVote.DiscardUnknown(m)
}

var xxx_messageInfo_AggregateExchangeRateVote proto.InternalMessageInfo

func (m *AggregateExchangeRateVote) GetExchangeRateTuples() []ExchangeRateTuple {
	if m != nil {
		return m.ExchangeRateTuples
	}
	return nil
}

func (m *AggregateExchangeRateVote) GetVoter() string {
	if m != nil {
		return m.Voter
	}
	return ""
}

func init() {
	proto.RegisterType((*ExchangeRatePrevote)(nil), "oracle.v1.ExchangeRatePrevote")
	proto.RegisterType((*ExchangeRateVote)(nil), "oracle.v1.ExchangeRateVote")
	proto.RegisterType((*AggregateExchangeRatePrevote)(nil), "oracle.v1.AggregateExchangeRatePrevote")
	proto.RegisterType((*ExchangeRateTuple)(nil), "oracle.v1.ExchangeRateTuple")
	proto.RegisterType((*AggregateExchangeRateVote)(nil), "oracle.v1.AggregateExchangeRateVote")
}

func init() { proto.RegisterFile("oracle/v1/vote.proto", fileDescriptor_92bdd9e640ece7e2) }

var fileDescriptor_92bdd9e640ece7e2 = []byte{
	// 386 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x52, 0x4f, 0x6b, 0xa3, 0x40,
	0x14, 0x77, 0x62, 0x76, 0x21, 0x93, 0x2c, 0xec, 0xba, 0x1e, 0xdc, 0x25, 0x18, 0xd7, 0xc3, 0x22,
	0x2c, 0xeb, 0x90, 0xdd, 0x4f, 0xb0, 0xb2, 0xb9, 0x17, 0x1b, 0x7a, 0xe8, 0x25, 0xa8, 0x79, 0x8c,
	0x69, 0x34, 0x23, 0xce, 0x44, 0x92, 0x4b, 0xaf, 0xbd, 0x16, 0xfa, 0xa5, 0x72, 0xcc, 0xb1, 0xf4,
	0x10, 0x4a, 0xf2, 0x45, 0xca, 0x68, 0x28, 0xb6, 0x09, 0x85, 0xd0, 0x93, 0xef, 0xfd, 0xf4, 0xbd,
	0xdf, 0x1f, 0x1f, 0xd6, 0x59, 0x1e, 0x44, 0x09, 0x90, 0xa2, 0x4f, 0x0a, 0x26, 0xc0, 0xcd, 0x72,
	0x26, 0x98, 0xd6, 0xaa, 0x50, 0xb7, 0xe8, 0x7f, 0xd7, 0x29, 0xa3, 0xac, 0x44, 0x89, 0xac, 0xaa,
	0x0f, 0xec, 0x05, 0xfe, 0x3a, 0x58, 0x44, 0x71, 0x30, 0xa3, 0xe0, 0x07, 0x02, 0xce, 0x72, 0x90,
	0xd3, 0x9a, 0x86, 0x9b, 0x71, 0xc0, 0x63, 0x03, 0x59, 0xc8, 0xe9, 0xf8, 0x65, 0xad, 0xe9, 0xf8,
	0xc3, 0x18, 0x66, 0x2c, 0x35, 0x1a, 0x16, 0x72, 0x5a, 0x7e, 0xd5, 0x48, 0x54, 0x4e, 0xe4, 0x86,
	0x5a, 0xa1, 0x65, 0xa3, 0xfd, 0xc0, 0x1d, 0x3e, 0x0f, 0xd3, 0x89, 0x18, 0x85, 0x09, 0x8b, 0xa6,
	0x46, 0xd3, 0x42, 0x8e, 0xea, 0xb7, 0x2b, 0xcc, 0x93, 0x90, 0x7d, 0x87, 0xf0, 0xe7, 0x3a, 0xf5,
	0x85, 0xe4, 0x3d, 0xc7, 0x9f, 0x60, 0x8f, 0x8d, 0xf2, 0x40, 0x40, 0x29, 0xa0, 0xe5, 0xb9, 0xab,
	0x4d, 0x4f, 0x79, 0xd8, 0xf4, 0x7e, 0xd2, 0x89, 0x88, 0xe7, 0xa1, 0x1b, 0xb1, 0x94, 0x44, 0x8c,
	0xa7, 0x8c, 0xef, 0x1f, 0xbf, 0xf9, 0x78, 0x4a, 0xc4, 0x32, 0x03, 0xee, 0xfe, 0x87, 0xc8, 0xef,
	0x40, 0x6d, 0xf1, 0x29, 0xc2, 0xed, 0x29, 0xee, 0xfe, 0xa3, 0x34, 0x07, 0x1a, 0x08, 0x38, 0x21,
	0x98, 0x6a, 0x53, 0xe3, 0xad, 0x08, 0xd4, 0xc3, 0x08, 0xae, 0xf1, 0x97, 0x3a, 0xc7, 0x70, 0x9e,
	0x25, 0x35, 0xb5, 0xa8, 0xae, 0xf6, 0x20, 0x98, 0xc6, 0xfb, 0x83, 0xb1, 0x6f, 0x10, 0xfe, 0x76,
	0xd4, 0x6d, 0xf9, 0x2f, 0x86, 0x58, 0x7f, 0x41, 0x39, 0x12, 0x52, 0x1f, 0x37, 0x90, 0xa5, 0x3a,
	0xed, 0x3f, 0x5d, 0xf7, 0xf9, 0xb4, 0xdc, 0x03, 0x13, 0x5e, 0x53, 0xea, 0xf2, 0x35, 0x78, 0xfd,
	0x82, 0x1f, 0x0f, 0xcb, 0x1b, 0xac, 0xb6, 0x26, 0x5a, 0x6f, 0x4d, 0xf4, 0xb8, 0x35, 0xd1, 0xed,
	0xce, 0x54, 0xd6, 0x3b, 0x53, 0xb9, 0xdf, 0x99, 0xca, 0xe5, 0xaf, 0x9a, 0xb3, 0x0c, 0x28, 0x5d,
	0x5e, 0x15, 0x84, 0xb3, 0x34, 0x85, 0x64, 0x02, 0x39, 0x59, 0x90, 0xfd, 0xd9, 0x97, 0x16, 0xc3,
	0x8f, 0xe5, 0x51, 0xff, 0x7d, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x09, 0xc5, 0xc4, 0x69, 0x0d, 0x03,
	0x00, 0x00,
}

func (m *ExchangeRatePrevote) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExchangeRatePrevote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExchangeRatePrevote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SubmitBlock != 0 {
		i = encodeVarintVote(dAtA, i, uint64(m.SubmitBlock))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Voter) > 0 {
		i -= len(m.Voter)
		copy(dAtA[i:], m.Voter)
		i = encodeVarintVote(dAtA, i, uint64(len(m.Voter)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintVote(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintVote(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ExchangeRateVote) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExchangeRateVote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExchangeRateVote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Voter) > 0 {
		i -= len(m.Voter)
		copy(dAtA[i:], m.Voter)
		i = encodeVarintVote(dAtA, i, uint64(len(m.Voter)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintVote(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x12
	}
	{
		size := m.ExchangeRate.Size()
		i -= size
		if _, err := m.ExchangeRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintVote(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *AggregateExchangeRatePrevote) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AggregateExchangeRatePrevote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AggregateExchangeRatePrevote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SubmitBlock != 0 {
		i = encodeVarintVote(dAtA, i, uint64(m.SubmitBlock))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Voter) > 0 {
		i -= len(m.Voter)
		copy(dAtA[i:], m.Voter)
		i = encodeVarintVote(dAtA, i, uint64(len(m.Voter)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintVote(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ExchangeRateTuple) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExchangeRateTuple) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExchangeRateTuple) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.ExchangeRate.Size()
		i -= size
		if _, err := m.ExchangeRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintVote(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintVote(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AggregateExchangeRateVote) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AggregateExchangeRateVote) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AggregateExchangeRateVote) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Voter) > 0 {
		i -= len(m.Voter)
		copy(dAtA[i:], m.Voter)
		i = encodeVarintVote(dAtA, i, uint64(len(m.Voter)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.ExchangeRateTuples) > 0 {
		for iNdEx := len(m.ExchangeRateTuples) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.ExchangeRateTuples[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintVote(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintVote(dAtA []byte, offset int, v uint64) int {
	offset -= sovVote(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ExchangeRatePrevote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovVote(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovVote(uint64(l))
	}
	l = len(m.Voter)
	if l > 0 {
		n += 1 + l + sovVote(uint64(l))
	}
	if m.SubmitBlock != 0 {
		n += 1 + sovVote(uint64(m.SubmitBlock))
	}
	return n
}

func (m *ExchangeRateVote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ExchangeRate.Size()
	n += 1 + l + sovVote(uint64(l))
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovVote(uint64(l))
	}
	l = len(m.Voter)
	if l > 0 {
		n += 1 + l + sovVote(uint64(l))
	}
	return n
}

func (m *AggregateExchangeRatePrevote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovVote(uint64(l))
	}
	l = len(m.Voter)
	if l > 0 {
		n += 1 + l + sovVote(uint64(l))
	}
	if m.SubmitBlock != 0 {
		n += 1 + sovVote(uint64(m.SubmitBlock))
	}
	return n
}

func (m *ExchangeRateTuple) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovVote(uint64(l))
	}
	l = m.ExchangeRate.Size()
	n += 1 + l + sovVote(uint64(l))
	return n
}

func (m *AggregateExchangeRateVote) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ExchangeRateTuples) > 0 {
		for _, e := range m.ExchangeRateTuples {
			l = e.Size()
			n += 1 + l + sovVote(uint64(l))
		}
	}
	l = len(m.Voter)
	if l > 0 {
		n += 1 + l + sovVote(uint64(l))
	}
	return n
}

func sovVote(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozVote(x uint64) (n int) {
	return sovVote(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ExchangeRatePrevote) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVote
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
			return fmt.Errorf("proto: ExchangeRatePrevote: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExchangeRatePrevote: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append(m.Hash[:0], dAtA[iNdEx:postIndex]...)
			if m.Hash == nil {
				m.Hash = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
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
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Voter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
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
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Voter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubmitBlock", wireType)
			}
			m.SubmitBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SubmitBlock |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipVote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthVote
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthVote
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
func (m *ExchangeRateVote) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVote
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
			return fmt.Errorf("proto: ExchangeRateVote: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExchangeRateVote: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExchangeRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
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
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ExchangeRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
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
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Voter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
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
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Voter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthVote
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthVote
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
func (m *AggregateExchangeRatePrevote) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVote
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
			return fmt.Errorf("proto: AggregateExchangeRatePrevote: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AggregateExchangeRatePrevote: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = append(m.Hash[:0], dAtA[iNdEx:postIndex]...)
			if m.Hash == nil {
				m.Hash = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Voter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
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
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Voter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubmitBlock", wireType)
			}
			m.SubmitBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SubmitBlock |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipVote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthVote
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthVote
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
func (m *ExchangeRateTuple) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVote
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
			return fmt.Errorf("proto: ExchangeRateTuple: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExchangeRateTuple: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
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
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExchangeRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
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
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ExchangeRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthVote
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthVote
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
func (m *AggregateExchangeRateVote) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowVote
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
			return fmt.Errorf("proto: AggregateExchangeRateVote: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AggregateExchangeRateVote: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExchangeRateTuples", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
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
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExchangeRateTuples = append(m.ExchangeRateTuples, ExchangeRateTuple{})
			if err := m.ExchangeRateTuples[len(m.ExchangeRateTuples)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Voter", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowVote
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
				return ErrInvalidLengthVote
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthVote
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Voter = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipVote(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthVote
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthVote
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
func skipVote(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowVote
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
					return 0, ErrIntOverflowVote
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
					return 0, ErrIntOverflowVote
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
				return 0, ErrInvalidLengthVote
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupVote
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthVote
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthVote        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowVote          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupVote = fmt.Errorf("proto: unexpected end of group")
)
