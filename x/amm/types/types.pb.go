// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/amm/types.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
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

type FilterType int32

const (
	FilterType_FilterAll           FilterType = 0
	FilterType_FilterPerpetual     FilterType = 1
	FilterType_FilterFixedWeight   FilterType = 2
	FilterType_FilterDynamicWeight FilterType = 3
	FilterType_FilterLeverage      FilterType = 4
)

var FilterType_name = map[int32]string{
	0: "FilterAll",
	1: "FilterPerpetual",
	2: "FilterFixedWeight",
	3: "FilterDynamicWeight",
	4: "FilterLeverage",
}

var FilterType_value = map[string]int32{
	"FilterAll":           0,
	"FilterPerpetual":     1,
	"FilterFixedWeight":   2,
	"FilterDynamicWeight": 3,
	"FilterLeverage":      4,
}

func (x FilterType) String() string {
	return proto.EnumName(FilterType_name, int32(x))
}

func (FilterType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_5382be80734a352d, []int{0}
}

type QueryEarnPoolRequest struct {
	PoolIds    []uint64           `protobuf:"varint,1,rep,packed,name=pool_ids,json=poolIds,proto3" json:"pool_ids,omitempty"`
	FilterType FilterType         `protobuf:"varint,2,opt,name=filter_type,json=filterType,proto3,enum=elys.amm.FilterType" json:"filter_type,omitempty"`
	Pagination *query.PageRequest `protobuf:"bytes,3,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryEarnPoolRequest) Reset()         { *m = QueryEarnPoolRequest{} }
func (m *QueryEarnPoolRequest) String() string { return proto.CompactTextString(m) }
func (*QueryEarnPoolRequest) ProtoMessage()    {}
func (*QueryEarnPoolRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_5382be80734a352d, []int{0}
}
func (m *QueryEarnPoolRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryEarnPoolRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryEarnPoolRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryEarnPoolRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEarnPoolRequest.Merge(m, src)
}
func (m *QueryEarnPoolRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryEarnPoolRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEarnPoolRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEarnPoolRequest proto.InternalMessageInfo

func (m *QueryEarnPoolRequest) GetPoolIds() []uint64 {
	if m != nil {
		return m.PoolIds
	}
	return nil
}

func (m *QueryEarnPoolRequest) GetFilterType() FilterType {
	if m != nil {
		return m.FilterType
	}
	return FilterType_FilterAll
}

func (m *QueryEarnPoolRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

type EarnPool struct {
	Assets      []PoolAsset                             `protobuf:"bytes,1,rep,name=assets,proto3" json:"assets"`
	PoolRatio   string                                  `protobuf:"bytes,2,opt,name=pool_ratio,json=poolRatio,proto3" json:"pool_ratio,omitempty"`
	RewardsApr  github_com_cosmos_cosmos_sdk_types.Dec  `protobuf:"bytes,3,opt,name=rewards_apr,json=rewardsApr,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"rewards_apr"`
	BorrowApr   github_com_cosmos_cosmos_sdk_types.Dec  `protobuf:"bytes,4,opt,name=borrow_apr,json=borrowApr,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"borrow_apr"`
	LeverageLp  github_com_cosmos_cosmos_sdk_types.Dec  `protobuf:"bytes,5,opt,name=leverage_lp,json=leverageLp,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"leverage_lp"`
	Perpetual   github_com_cosmos_cosmos_sdk_types.Dec  `protobuf:"bytes,6,opt,name=perpetual,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"perpetual"`
	Tvl         github_com_cosmos_cosmos_sdk_types.Dec  `protobuf:"bytes,7,opt,name=tvl,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"tvl"`
	Rewards     github_com_cosmos_cosmos_sdk_types.Dec  `protobuf:"bytes,8,opt,name=rewards,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"rewards"`
	PoolId      uint64                                  `protobuf:"varint,9,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	TotalShares github_com_cosmos_cosmos_sdk_types.Coin `protobuf:"bytes,10,opt,name=total_shares,json=totalShares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Coin" json:"total_shares"`
}

func (m *EarnPool) Reset()         { *m = EarnPool{} }
func (m *EarnPool) String() string { return proto.CompactTextString(m) }
func (*EarnPool) ProtoMessage()    {}
func (*EarnPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_5382be80734a352d, []int{1}
}
func (m *EarnPool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EarnPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EarnPool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EarnPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EarnPool.Merge(m, src)
}
func (m *EarnPool) XXX_Size() int {
	return m.Size()
}
func (m *EarnPool) XXX_DiscardUnknown() {
	xxx_messageInfo_EarnPool.DiscardUnknown(m)
}

var xxx_messageInfo_EarnPool proto.InternalMessageInfo

func (m *EarnPool) GetAssets() []PoolAsset {
	if m != nil {
		return m.Assets
	}
	return nil
}

func (m *EarnPool) GetPoolRatio() string {
	if m != nil {
		return m.PoolRatio
	}
	return ""
}

func (m *EarnPool) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

type QueryEarnPoolResponse struct {
	Pools []EarnPool `protobuf:"bytes,1,rep,name=pools,proto3" json:"pools"`
}

func (m *QueryEarnPoolResponse) Reset()         { *m = QueryEarnPoolResponse{} }
func (m *QueryEarnPoolResponse) String() string { return proto.CompactTextString(m) }
func (*QueryEarnPoolResponse) ProtoMessage()    {}
func (*QueryEarnPoolResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_5382be80734a352d, []int{2}
}
func (m *QueryEarnPoolResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryEarnPoolResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryEarnPoolResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryEarnPoolResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryEarnPoolResponse.Merge(m, src)
}
func (m *QueryEarnPoolResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryEarnPoolResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryEarnPoolResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryEarnPoolResponse proto.InternalMessageInfo

func (m *QueryEarnPoolResponse) GetPools() []EarnPool {
	if m != nil {
		return m.Pools
	}
	return nil
}

func init() {
	proto.RegisterEnum("elys.amm.FilterType", FilterType_name, FilterType_value)
	proto.RegisterType((*QueryEarnPoolRequest)(nil), "elys.amm.QueryEarnPoolRequest")
	proto.RegisterType((*EarnPool)(nil), "elys.amm.EarnPool")
	proto.RegisterType((*QueryEarnPoolResponse)(nil), "elys.amm.QueryEarnPoolResponse")
}

func init() { proto.RegisterFile("elys/amm/types.proto", fileDescriptor_5382be80734a352d) }

var fileDescriptor_5382be80734a352d = []byte{
	// 621 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0x4d, 0x4f, 0x13, 0x41,
	0x18, 0xc7, 0xbb, 0xb4, 0xb4, 0xdd, 0x67, 0x15, 0xeb, 0x50, 0xc2, 0x42, 0xe2, 0xd2, 0x70, 0xc0,
	0x86, 0x84, 0xdd, 0x80, 0xf1, 0x2e, 0x15, 0x51, 0x13, 0x0c, 0xb8, 0x9a, 0x98, 0x78, 0x69, 0xa6,
	0xed, 0xb0, 0x6c, 0xd8, 0xdd, 0x19, 0x66, 0xa6, 0x40, 0xbf, 0x85, 0xdf, 0x44, 0x3f, 0x06, 0x47,
	0x8e, 0xc6, 0x03, 0x31, 0xf0, 0x45, 0xcc, 0xbc, 0x2c, 0x45, 0x4f, 0xa6, 0xa7, 0xee, 0xf3, 0xf6,
	0xeb, 0xff, 0x99, 0xfd, 0xef, 0x40, 0x9b, 0x64, 0x13, 0x11, 0xe1, 0x3c, 0x8f, 0xe4, 0x84, 0x11,
	0x11, 0x32, 0x4e, 0x25, 0x45, 0x4d, 0x95, 0x0d, 0x71, 0x9e, 0xaf, 0xb6, 0x13, 0x9a, 0x50, 0x9d,
	0x8c, 0xd4, 0x93, 0xa9, 0xaf, 0x6e, 0x0e, 0xa9, 0xc8, 0xa9, 0x88, 0x06, 0x58, 0x90, 0xe8, 0x6c,
	0x4c, 0xf8, 0x24, 0x3a, 0xdf, 0x1e, 0x10, 0x89, 0xb7, 0x23, 0x86, 0x93, 0xb4, 0xc0, 0x32, 0xa5,
	0x85, 0xed, 0x5d, 0xb9, 0xff, 0x07, 0x46, 0x69, 0xd6, 0xc7, 0x42, 0x10, 0x69, 0x4b, 0xc1, 0x43,
	0x4c, 0x09, 0x18, 0xd2, 0xd4, 0x8e, 0xae, 0xff, 0x70, 0xa0, 0xfd, 0x51, 0xd1, 0xdf, 0x60, 0x5e,
	0x1c, 0x51, 0x9a, 0xc5, 0xe4, 0x6c, 0x4c, 0x84, 0x44, 0x2b, 0xd0, 0xd4, 0xb0, 0x74, 0x24, 0x7c,
	0xa7, 0x53, 0xed, 0xd6, 0xe2, 0x86, 0x8a, 0xdf, 0x8f, 0x04, 0x7a, 0x09, 0xde, 0x71, 0x9a, 0x49,
	0xc2, 0xfb, 0x6a, 0x21, 0x7f, 0xae, 0xe3, 0x74, 0x17, 0x76, 0xda, 0x61, 0xb9, 0x50, 0xb8, 0xaf,
	0x8b, 0x9f, 0x27, 0x8c, 0xc4, 0x70, 0x7c, 0xff, 0x8c, 0xf6, 0x01, 0xa6, 0xca, 0xfd, 0x6a, 0xc7,
	0xe9, 0x7a, 0x3b, 0x1b, 0xa1, 0xd1, 0x17, 0x2a, 0x7d, 0xa1, 0x5e, 0x33, 0xb4, 0x2a, 0xc3, 0x23,
	0x9c, 0x10, 0xab, 0x26, 0x7e, 0x30, 0xb9, 0xfe, 0x7d, 0x1e, 0x9a, 0xa5, 0x5a, 0xb4, 0x0d, 0x75,
	0xbd, 0xae, 0x11, 0xe9, 0xed, 0x2c, 0x4e, 0x65, 0xa8, 0xfa, 0xae, 0xaa, 0xf5, 0x6a, 0x57, 0x37,
	0x6b, 0x95, 0xd8, 0x36, 0xa2, 0x67, 0x00, 0x7a, 0x33, 0xae, 0x70, 0x5a, 0xbd, 0x1b, 0xbb, 0x2a,
	0x13, 0xab, 0x04, 0x3a, 0x04, 0x8f, 0x93, 0x0b, 0xcc, 0x47, 0xa2, 0x8f, 0x19, 0xd7, 0x3a, 0xdd,
	0x5e, 0xa8, 0x08, 0xbf, 0x6e, 0xd6, 0x36, 0x92, 0x54, 0x9e, 0x8c, 0x07, 0xe1, 0x90, 0xe6, 0x91,
	0x3d, 0x59, 0xf3, 0xb3, 0x25, 0x46, 0xa7, 0xf6, 0xfd, 0xee, 0x91, 0x61, 0x0c, 0x16, 0xb1, 0xcb,
	0x38, 0xfa, 0x00, 0x30, 0xa0, 0x9c, 0xd3, 0x0b, 0xcd, 0xab, 0xcd, 0xc4, 0x73, 0x0d, 0x41, 0xe1,
	0x0e, 0xc1, 0xcb, 0xc8, 0x39, 0xe1, 0x38, 0x21, 0xfd, 0x8c, 0xf9, 0xf3, 0xb3, 0xe9, 0x2b, 0x11,
	0x07, 0x0c, 0x1d, 0x80, 0xcb, 0x08, 0x67, 0x44, 0x8e, 0x71, 0xe6, 0xd7, 0x67, 0x93, 0x77, 0x0f,
	0x40, 0xaf, 0xa0, 0x2a, 0xcf, 0x33, 0xbf, 0x31, 0x13, 0x47, 0x8d, 0xa2, 0x77, 0xd0, 0xb0, 0xa7,
	0xe7, 0x37, 0x67, 0xa2, 0x94, 0xe3, 0x68, 0x19, 0x1a, 0xd6, 0xc3, 0xbe, 0xdb, 0x71, 0xba, 0xb5,
	0xb8, 0x6e, 0x2c, 0x8c, 0x72, 0x78, 0x24, 0xa9, 0xc4, 0x59, 0x5f, 0x9c, 0x60, 0x4e, 0x84, 0x0f,
	0xda, 0x8c, 0x2b, 0x7f, 0x99, 0xb1, 0xb4, 0xe1, 0x6b, 0x9a, 0x16, 0xbd, 0xc8, 0x4a, 0x78, 0xfe,
	0x1f, 0x12, 0xd4, 0x40, 0xec, 0x69, 0xfe, 0x27, 0x8d, 0x5f, 0x7f, 0x0b, 0x4b, 0xff, 0x7c, 0x63,
	0x82, 0xd1, 0x42, 0x10, 0x14, 0xc2, 0xbc, 0x52, 0x54, 0x9a, 0x17, 0x4d, 0xcd, 0x5b, 0xb6, 0x5a,
	0xef, 0x9a, 0xb6, 0x4d, 0x09, 0x30, 0xfd, 0xb8, 0xd0, 0x63, 0x70, 0x4d, 0xb4, 0x9b, 0x65, 0xad,
	0x0a, 0x5a, 0x84, 0x27, 0x26, 0x3c, 0x2a, 0x5f, 0x46, 0xcb, 0x41, 0x4b, 0xf0, 0xd4, 0x24, 0xf7,
	0xd3, 0x4b, 0x32, 0xfa, 0x42, 0xd2, 0xe4, 0x44, 0xb6, 0xe6, 0xd0, 0x32, 0x2c, 0x9a, 0xf4, 0xde,
	0xa4, 0xc0, 0x79, 0x3a, 0xb4, 0x85, 0x2a, 0x42, 0xb0, 0x60, 0x0a, 0x07, 0xd6, 0x20, 0xad, 0x5a,
	0xaf, 0x77, 0x75, 0x1b, 0x38, 0xd7, 0xb7, 0x81, 0xf3, 0xfb, 0x36, 0x70, 0xbe, 0xdd, 0x05, 0x95,
	0xeb, 0xbb, 0xa0, 0xf2, 0xf3, 0x2e, 0xa8, 0x7c, 0xed, 0x3e, 0x38, 0x0e, 0x25, 0x7d, 0xab, 0x20,
	0xf2, 0x82, 0xf2, 0x53, 0x1d, 0x44, 0x97, 0xd3, 0x4b, 0x6f, 0x50, 0xd7, 0xd7, 0xcd, 0x8b, 0x3f,
	0x01, 0x00, 0x00, 0xff, 0xff, 0xbe, 0x8d, 0xb7, 0x91, 0x0d, 0x05, 0x00, 0x00,
}

func (m *QueryEarnPoolRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryEarnPoolRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryEarnPoolRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.FilterType != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.FilterType))
		i--
		dAtA[i] = 0x10
	}
	if len(m.PoolIds) > 0 {
		dAtA3 := make([]byte, len(m.PoolIds)*10)
		var j2 int
		for _, num := range m.PoolIds {
			for num >= 1<<7 {
				dAtA3[j2] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j2++
			}
			dAtA3[j2] = uint8(num)
			j2++
		}
		i -= j2
		copy(dAtA[i:], dAtA3[:j2])
		i = encodeVarintTypes(dAtA, i, uint64(j2))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *EarnPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EarnPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EarnPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.TotalShares.Size()
		i -= size
		if _, err := m.TotalShares.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	if m.PoolId != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x48
	}
	{
		size := m.Rewards.Size()
		i -= size
		if _, err := m.Rewards.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.Tvl.Size()
		i -= size
		if _, err := m.Tvl.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.Perpetual.Size()
		i -= size
		if _, err := m.Perpetual.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.LeverageLp.Size()
		i -= size
		if _, err := m.LeverageLp.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.BorrowApr.Size()
		i -= size
		if _, err := m.BorrowApr.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.RewardsApr.Size()
		i -= size
		if _, err := m.RewardsApr.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.PoolRatio) > 0 {
		i -= len(m.PoolRatio)
		copy(dAtA[i:], m.PoolRatio)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.PoolRatio)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Assets) > 0 {
		for iNdEx := len(m.Assets) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Assets[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *QueryEarnPoolResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryEarnPoolResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryEarnPoolResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Pools) > 0 {
		for iNdEx := len(m.Pools) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Pools[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintTypes(dAtA []byte, offset int, v uint64) int {
	offset -= sovTypes(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryEarnPoolRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.PoolIds) > 0 {
		l = 0
		for _, e := range m.PoolIds {
			l += sovTypes(uint64(e))
		}
		n += 1 + sovTypes(uint64(l)) + l
	}
	if m.FilterType != 0 {
		n += 1 + sovTypes(uint64(m.FilterType))
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	return n
}

func (m *EarnPool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Assets) > 0 {
		for _, e := range m.Assets {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	l = len(m.PoolRatio)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.RewardsApr.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.BorrowApr.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.LeverageLp.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.Perpetual.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.Tvl.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.Rewards.Size()
	n += 1 + l + sovTypes(uint64(l))
	if m.PoolId != 0 {
		n += 1 + sovTypes(uint64(m.PoolId))
	}
	l = m.TotalShares.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func (m *QueryEarnPoolResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Pools) > 0 {
		for _, e := range m.Pools {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryEarnPoolRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: QueryEarnPoolRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryEarnPoolRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTypes
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.PoolIds = append(m.PoolIds, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTypes
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthTypes
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthTypes
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.PoolIds) == 0 {
					m.PoolIds = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTypes
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.PoolIds = append(m.PoolIds, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolIds", wireType)
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FilterType", wireType)
			}
			m.FilterType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FilterType |= FilterType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *EarnPool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: EarnPool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EarnPool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Assets", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Assets = append(m.Assets, PoolAsset{})
			if err := m.Assets[len(m.Assets)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolRatio", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PoolRatio = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardsApr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardsApr.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BorrowApr", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BorrowApr.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LeverageLp", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.LeverageLp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Perpetual", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Perpetual.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tvl", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Tvl.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rewards", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Rewards.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalShares", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalShares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func (m *QueryEarnPoolResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTypes
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
			return fmt.Errorf("proto: QueryEarnPoolResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryEarnPoolResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pools", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
				return ErrInvalidLengthTypes
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTypes
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pools = append(m.Pools, EarnPool{})
			if err := m.Pools[len(m.Pools)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTypes(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTypes
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
func skipTypes(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
					return 0, ErrIntOverflowTypes
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
				return 0, ErrInvalidLengthTypes
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTypes
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTypes
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTypes        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTypes          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTypes = fmt.Errorf("proto: unexpected end of group")
)
