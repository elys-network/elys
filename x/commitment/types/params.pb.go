// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/commitment/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

type EarnType int32

const (
	EarnType_ALL_PROGRAM       EarnType = 0
	EarnType_USDC_PROGRAM      EarnType = 1
	EarnType_ELYS_PROGRAM      EarnType = 2
	EarnType_EDEN_PROGRAM      EarnType = 3
	EarnType_EDENB_PROGRAM     EarnType = 4
	EarnType_LP_MINING_PROGRAM EarnType = 5
)

var EarnType_name = map[int32]string{
	0: "ALL_PROGRAM",
	1: "USDC_PROGRAM",
	2: "ELYS_PROGRAM",
	3: "EDEN_PROGRAM",
	4: "EDENB_PROGRAM",
	5: "LP_MINING_PROGRAM",
}

var EarnType_value = map[string]int32{
	"ALL_PROGRAM":       0,
	"USDC_PROGRAM":      1,
	"ELYS_PROGRAM":      2,
	"EDEN_PROGRAM":      3,
	"EDENB_PROGRAM":     4,
	"LP_MINING_PROGRAM": 5,
}

func (x EarnType) String() string {
	return proto.EnumName(EarnType_name, int32(x))
}

func (EarnType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_92e317feaf73ff7e, []int{0}
}

// Params defines the parameters for the module.
type Params struct {
	VestingInfos        []*VestingInfo                           `protobuf:"bytes,1,rep,name=vesting_infos,json=vestingInfos,proto3" json:"vesting_infos,omitempty"`
	TotalCommitted      github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=total_committed,json=totalCommitted,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_committed"`
	NumberOfCommitments uint64                                   `protobuf:"varint,3,opt,name=number_of_commitments,json=numberOfCommitments,proto3" json:"number_of_commitments,omitempty"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_92e317feaf73ff7e, []int{0}
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

func (m *Params) GetVestingInfos() []*VestingInfo {
	if m != nil {
		return m.VestingInfos
	}
	return nil
}

func (m *Params) GetTotalCommitted() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.TotalCommitted
	}
	return nil
}

func (m *Params) GetNumberOfCommitments() uint64 {
	if m != nil {
		return m.NumberOfCommitments
	}
	return 0
}

type LegacyParams struct {
	VestingInfos   []*VestingInfo                           `protobuf:"bytes,1,rep,name=vesting_infos,json=vestingInfos,proto3" json:"vesting_infos,omitempty"`
	TotalCommitted github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=total_committed,json=totalCommitted,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_committed"`
}

func (m *LegacyParams) Reset()      { *m = LegacyParams{} }
func (*LegacyParams) ProtoMessage() {}
func (*LegacyParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_92e317feaf73ff7e, []int{1}
}
func (m *LegacyParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LegacyParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LegacyParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LegacyParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LegacyParams.Merge(m, src)
}
func (m *LegacyParams) XXX_Size() int {
	return m.Size()
}
func (m *LegacyParams) XXX_DiscardUnknown() {
	xxx_messageInfo_LegacyParams.DiscardUnknown(m)
}

var xxx_messageInfo_LegacyParams proto.InternalMessageInfo

func (m *LegacyParams) GetVestingInfos() []*VestingInfo {
	if m != nil {
		return m.VestingInfos
	}
	return nil
}

func (m *LegacyParams) GetTotalCommitted() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.TotalCommitted
	}
	return nil
}

type VestingInfo struct {
	BaseDenom      string                                 `protobuf:"bytes,1,opt,name=base_denom,json=baseDenom,proto3" json:"base_denom,omitempty"`
	VestingDenom   string                                 `protobuf:"bytes,2,opt,name=vesting_denom,json=vestingDenom,proto3" json:"vesting_denom,omitempty"`
	NumBlocks      int64                                  `protobuf:"varint,3,opt,name=num_blocks,json=numBlocks,proto3" json:"num_blocks,omitempty"`
	VestNowFactor  github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,4,opt,name=vest_now_factor,json=vestNowFactor,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"vest_now_factor"`
	NumMaxVestings int64                                  `protobuf:"varint,5,opt,name=num_max_vestings,json=numMaxVestings,proto3" json:"num_max_vestings,omitempty"`
}

func (m *VestingInfo) Reset()         { *m = VestingInfo{} }
func (m *VestingInfo) String() string { return proto.CompactTextString(m) }
func (*VestingInfo) ProtoMessage()    {}
func (*VestingInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_92e317feaf73ff7e, []int{2}
}
func (m *VestingInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *VestingInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_VestingInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *VestingInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VestingInfo.Merge(m, src)
}
func (m *VestingInfo) XXX_Size() int {
	return m.Size()
}
func (m *VestingInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_VestingInfo.DiscardUnknown(m)
}

var xxx_messageInfo_VestingInfo proto.InternalMessageInfo

func (m *VestingInfo) GetBaseDenom() string {
	if m != nil {
		return m.BaseDenom
	}
	return ""
}

func (m *VestingInfo) GetVestingDenom() string {
	if m != nil {
		return m.VestingDenom
	}
	return ""
}

func (m *VestingInfo) GetNumBlocks() int64 {
	if m != nil {
		return m.NumBlocks
	}
	return 0
}

func (m *VestingInfo) GetNumMaxVestings() int64 {
	if m != nil {
		return m.NumMaxVestings
	}
	return 0
}

func init() {
	proto.RegisterEnum("elys.commitment.EarnType", EarnType_name, EarnType_value)
	proto.RegisterType((*Params)(nil), "elys.commitment.Params")
	proto.RegisterType((*LegacyParams)(nil), "elys.commitment.LegacyParams")
	proto.RegisterType((*VestingInfo)(nil), "elys.commitment.VestingInfo")
}

func init() { proto.RegisterFile("elys/commitment/params.proto", fileDescriptor_92e317feaf73ff7e) }

var fileDescriptor_92e317feaf73ff7e = []byte{
	// 546 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x53, 0x3f, 0x6f, 0xd3, 0x40,
	0x1c, 0xb5, 0x93, 0xb4, 0x22, 0x97, 0xb4, 0x49, 0x0f, 0x2a, 0xb9, 0x55, 0xeb, 0x44, 0x45, 0x42,
	0x11, 0x52, 0x6d, 0x5a, 0x36, 0xb6, 0xfc, 0xa3, 0x8a, 0x94, 0xa4, 0x91, 0x0b, 0x95, 0x60, 0x39,
	0xd9, 0xce, 0x25, 0x58, 0x89, 0xef, 0x22, 0xdf, 0xe5, 0x9f, 0xf8, 0x12, 0x8c, 0x8c, 0xcc, 0x7c,
	0x92, 0x2e, 0x48, 0x1d, 0x11, 0x43, 0x41, 0xc9, 0x37, 0x60, 0x65, 0x41, 0x77, 0xe7, 0x26, 0x11,
	0x13, 0x6b, 0x27, 0xfb, 0xf7, 0xde, 0xf3, 0xbd, 0xf7, 0x7e, 0xf2, 0x81, 0x23, 0x3c, 0x9c, 0x33,
	0xdb, 0xa7, 0x61, 0x18, 0xf0, 0x10, 0x13, 0x6e, 0x8f, 0xdc, 0xc8, 0x0d, 0x99, 0x35, 0x8a, 0x28,
	0xa7, 0x30, 0x27, 0x58, 0x6b, 0xcd, 0x1e, 0x3e, 0xe9, 0xd3, 0x3e, 0x95, 0x9c, 0x2d, 0xde, 0x94,
	0xec, 0xf0, 0xc0, 0xa7, 0x2c, 0xa4, 0x0c, 0x29, 0x42, 0x0d, 0x31, 0x65, 0xaa, 0xc9, 0xf6, 0x5c,
	0x86, 0xed, 0xc9, 0x99, 0x87, 0xb9, 0x7b, 0x66, 0xfb, 0x34, 0x20, 0x8a, 0x3f, 0xf9, 0xa3, 0x83,
	0xed, 0x8e, 0xb4, 0x84, 0x65, 0xb0, 0x33, 0xc1, 0x8c, 0x07, 0xa4, 0x8f, 0x02, 0xd2, 0xa3, 0xcc,
	0xd0, 0x8b, 0xc9, 0x52, 0xe6, 0xfc, 0xc8, 0xfa, 0x27, 0x84, 0x75, 0xad, 0x54, 0x0d, 0xd2, 0xa3,
	0x4e, 0x76, 0xb2, 0x1e, 0x18, 0xe4, 0x20, 0xc7, 0x29, 0x77, 0x87, 0x48, 0xa9, 0x39, 0xee, 0x1a,
	0x09, 0x79, 0xc8, 0x81, 0x15, 0xa7, 0x12, 0x39, 0xac, 0x38, 0x87, 0x55, 0xa5, 0x01, 0xa9, 0xbc,
	0xb8, 0xb9, 0x2b, 0x68, 0x5f, 0x7f, 0x16, 0x4a, 0xfd, 0x80, 0x7f, 0x18, 0x7b, 0xc2, 0x28, 0xae,
	0x10, 0x3f, 0x4e, 0x59, 0x77, 0x60, 0xf3, 0xf9, 0x08, 0x33, 0xf9, 0x01, 0x73, 0x76, 0xa5, 0x47,
	0xf5, 0xde, 0x02, 0x9e, 0x83, 0x7d, 0x32, 0x0e, 0x3d, 0x1c, 0x21, 0xda, 0x43, 0xeb, 0x9c, 0xcc,
	0x48, 0x16, 0xf5, 0x52, 0xca, 0x79, 0xac, 0xc8, 0xcb, 0x5e, 0x75, 0x4d, 0xbd, 0x4a, 0x7d, 0xfe,
	0x52, 0xd0, 0x4e, 0xbe, 0xe9, 0x20, 0xdb, 0xc4, 0x7d, 0xd7, 0x9f, 0x3f, 0xf0, 0x1d, 0xc4, 0x7d,
	0x7e, 0xeb, 0x20, 0xb3, 0x91, 0x0c, 0x1e, 0x03, 0x20, 0xcc, 0x50, 0x17, 0x13, 0x1a, 0x1a, 0x7a,
	0x51, 0x2f, 0xa5, 0x9d, 0xb4, 0x40, 0x6a, 0x02, 0x80, 0x4f, 0xd7, 0x6d, 0x95, 0x22, 0x21, 0x15,
	0xf7, 0x7d, 0x94, 0xe8, 0x18, 0x00, 0x32, 0x0e, 0x91, 0x37, 0xa4, 0xfe, 0x40, 0xad, 0x34, 0xe9,
	0xa4, 0xc9, 0x38, 0xac, 0x48, 0x00, 0x5e, 0x83, 0x9c, 0x90, 0x23, 0x42, 0xa7, 0xa8, 0xe7, 0xfa,
	0x9c, 0x46, 0x46, 0x4a, 0x9c, 0x52, 0xb1, 0x44, 0xa7, 0x1f, 0x77, 0x85, 0x67, 0xff, 0xd1, 0xa9,
	0x41, 0xb8, 0x23, 0xa3, 0xb4, 0xe9, 0xf4, 0xb5, 0x3c, 0x04, 0x96, 0x40, 0x5e, 0xd8, 0x86, 0xee,
	0x0c, 0xc5, 0x71, 0x98, 0xb1, 0x25, 0xcd, 0x77, 0xc9, 0x38, 0x6c, 0xb9, 0xb3, 0xb8, 0x27, 0x7b,
	0xfe, 0x11, 0x3c, 0xaa, 0xbb, 0x11, 0x79, 0x33, 0x1f, 0x61, 0x98, 0x03, 0x99, 0x72, 0xb3, 0x89,
	0x3a, 0xce, 0xe5, 0x85, 0x53, 0x6e, 0xe5, 0x35, 0x98, 0x07, 0xd9, 0xb7, 0x57, 0xb5, 0xea, 0x0a,
	0xd1, 0x05, 0x52, 0x6f, 0xbe, 0xbb, 0x5a, 0x21, 0x09, 0x89, 0xd4, 0xea, 0xed, 0x15, 0x92, 0x84,
	0x7b, 0x60, 0x47, 0x20, 0x95, 0x15, 0x94, 0x82, 0xfb, 0x60, 0xaf, 0xd9, 0x41, 0xad, 0x46, 0xbb,
	0xd1, 0xbe, 0x58, 0xc1, 0x5b, 0x95, 0xc6, 0xcd, 0xc2, 0xd4, 0x6f, 0x17, 0xa6, 0xfe, 0x6b, 0x61,
	0xea, 0x9f, 0x96, 0xa6, 0x76, 0xbb, 0x34, 0xb5, 0xef, 0x4b, 0x53, 0x7b, 0x6f, 0x6f, 0xf4, 0x16,
	0x7f, 0xcf, 0x29, 0xc1, 0x7c, 0x4a, 0xa3, 0x81, 0x1c, 0xec, 0xd9, 0xe6, 0x9d, 0x97, 0x4b, 0xf0,
	0xb6, 0xe5, 0x8d, 0x7c, 0xf9, 0x37, 0x00, 0x00, 0xff, 0xff, 0x57, 0x98, 0x9f, 0x4b, 0x13, 0x04,
	0x00, 0x00,
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
	if m.NumberOfCommitments != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.NumberOfCommitments))
		i--
		dAtA[i] = 0x18
	}
	if len(m.TotalCommitted) > 0 {
		for iNdEx := len(m.TotalCommitted) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TotalCommitted[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.VestingInfos) > 0 {
		for iNdEx := len(m.VestingInfos) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.VestingInfos[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *LegacyParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LegacyParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LegacyParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TotalCommitted) > 0 {
		for iNdEx := len(m.TotalCommitted) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TotalCommitted[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.VestingInfos) > 0 {
		for iNdEx := len(m.VestingInfos) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.VestingInfos[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *VestingInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *VestingInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *VestingInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.NumMaxVestings != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.NumMaxVestings))
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.VestNowFactor.Size()
		i -= size
		if _, err := m.VestNowFactor.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if m.NumBlocks != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.NumBlocks))
		i--
		dAtA[i] = 0x18
	}
	if len(m.VestingDenom) > 0 {
		i -= len(m.VestingDenom)
		copy(dAtA[i:], m.VestingDenom)
		i = encodeVarintParams(dAtA, i, uint64(len(m.VestingDenom)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.BaseDenom) > 0 {
		i -= len(m.BaseDenom)
		copy(dAtA[i:], m.BaseDenom)
		i = encodeVarintParams(dAtA, i, uint64(len(m.BaseDenom)))
		i--
		dAtA[i] = 0xa
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
	if len(m.VestingInfos) > 0 {
		for _, e := range m.VestingInfos {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.TotalCommitted) > 0 {
		for _, e := range m.TotalCommitted {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if m.NumberOfCommitments != 0 {
		n += 1 + sovParams(uint64(m.NumberOfCommitments))
	}
	return n
}

func (m *LegacyParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.VestingInfos) > 0 {
		for _, e := range m.VestingInfos {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if len(m.TotalCommitted) > 0 {
		for _, e := range m.TotalCommitted {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	return n
}

func (m *VestingInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.BaseDenom)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.VestingDenom)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.NumBlocks != 0 {
		n += 1 + sovParams(uint64(m.NumBlocks))
	}
	l = m.VestNowFactor.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.NumMaxVestings != 0 {
		n += 1 + sovParams(uint64(m.NumMaxVestings))
	}
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
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VestingInfos", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.VestingInfos = append(m.VestingInfos, &VestingInfo{})
			if err := m.VestingInfos[len(m.VestingInfos)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalCommitted", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TotalCommitted = append(m.TotalCommitted, types.Coin{})
			if err := m.TotalCommitted[len(m.TotalCommitted)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumberOfCommitments", wireType)
			}
			m.NumberOfCommitments = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumberOfCommitments |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
func (m *LegacyParams) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: LegacyParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LegacyParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VestingInfos", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.VestingInfos = append(m.VestingInfos, &VestingInfo{})
			if err := m.VestingInfos[len(m.VestingInfos)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalCommitted", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TotalCommitted = append(m.TotalCommitted, types.Coin{})
			if err := m.TotalCommitted[len(m.TotalCommitted)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *VestingInfo) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: VestingInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: VestingInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseDenom", wireType)
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
			m.BaseDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VestingDenom", wireType)
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
			m.VestingDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumBlocks", wireType)
			}
			m.NumBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumBlocks |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VestNowFactor", wireType)
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
			if err := m.VestNowFactor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumMaxVestings", wireType)
			}
			m.NumMaxVestings = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumMaxVestings |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
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
