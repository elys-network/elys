// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/margin/pool.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

type PoolAsset struct {
	Liabilities           github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,1,opt,name=liabilities,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"liabilities"`
	Custody               github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=custody,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"custody"`
	TakeProfitLiabilities github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,opt,name=take_profit_liabilities,json=takeProfitLiabilities,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"take_profit_liabilities"`
	TakeProfitCustody     github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,4,opt,name=take_profit_custody,json=takeProfitCustody,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"take_profit_custody"`
	AssetBalance          github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,5,opt,name=asset_balance,json=assetBalance,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"asset_balance"`
	BlockBorrowInterest   github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=block_borrow_interest,json=blockBorrowInterest,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"block_borrow_interest"`
	AssetDenom            string                                 `protobuf:"bytes,7,opt,name=asset_denom,json=assetDenom,proto3" json:"asset_denom,omitempty"`
}

func (m *PoolAsset) Reset()         { *m = PoolAsset{} }
func (m *PoolAsset) String() string { return proto.CompactTextString(m) }
func (*PoolAsset) ProtoMessage()    {}
func (*PoolAsset) Descriptor() ([]byte, []int) {
	return fileDescriptor_030dd771f91c2bb4, []int{0}
}
func (m *PoolAsset) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PoolAsset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PoolAsset.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PoolAsset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolAsset.Merge(m, src)
}
func (m *PoolAsset) XXX_Size() int {
	return m.Size()
}
func (m *PoolAsset) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolAsset.DiscardUnknown(m)
}

var xxx_messageInfo_PoolAsset proto.InternalMessageInfo

func (m *PoolAsset) GetAssetDenom() string {
	if m != nil {
		return m.AssetDenom
	}
	return ""
}

type Pool struct {
	AmmPoolId                            uint64                                 `protobuf:"varint,1,opt,name=amm_pool_id,json=ammPoolId,proto3" json:"amm_pool_id,omitempty"`
	Health                               github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=health,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"health"`
	Enabled                              bool                                   `protobuf:"varint,3,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Closed                               bool                                   `protobuf:"varint,4,opt,name=closed,proto3" json:"closed,omitempty"`
	BorrowInterestRate                   github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=borrow_interest_rate,json=borrowInterestRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"borrow_interest_rate"`
	PoolAssetsLong                       []PoolAsset                            `protobuf:"bytes,6,rep,name=pool_assets_long,json=poolAssetsLong,proto3" json:"pool_assets_long"`
	PoolAssetsShort                      []PoolAsset                            `protobuf:"bytes,7,rep,name=pool_assets_short,json=poolAssetsShort,proto3" json:"pool_assets_short"`
	LastHeightBorrowInterestRateComputed int64                                  `protobuf:"varint,8,opt,name=last_height_borrow_interest_rate_computed,json=lastHeightBorrowInterestRateComputed,proto3" json:"last_height_borrow_interest_rate_computed,omitempty"`
}

func (m *Pool) Reset()         { *m = Pool{} }
func (m *Pool) String() string { return proto.CompactTextString(m) }
func (*Pool) ProtoMessage()    {}
func (*Pool) Descriptor() ([]byte, []int) {
	return fileDescriptor_030dd771f91c2bb4, []int{1}
}
func (m *Pool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Pool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Pool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Pool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pool.Merge(m, src)
}
func (m *Pool) XXX_Size() int {
	return m.Size()
}
func (m *Pool) XXX_DiscardUnknown() {
	xxx_messageInfo_Pool.DiscardUnknown(m)
}

var xxx_messageInfo_Pool proto.InternalMessageInfo

func (m *Pool) GetAmmPoolId() uint64 {
	if m != nil {
		return m.AmmPoolId
	}
	return 0
}

func (m *Pool) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *Pool) GetClosed() bool {
	if m != nil {
		return m.Closed
	}
	return false
}

func (m *Pool) GetPoolAssetsLong() []PoolAsset {
	if m != nil {
		return m.PoolAssetsLong
	}
	return nil
}

func (m *Pool) GetPoolAssetsShort() []PoolAsset {
	if m != nil {
		return m.PoolAssetsShort
	}
	return nil
}

func (m *Pool) GetLastHeightBorrowInterestRateComputed() int64 {
	if m != nil {
		return m.LastHeightBorrowInterestRateComputed
	}
	return 0
}

func init() {
	proto.RegisterType((*PoolAsset)(nil), "elys.margin.PoolAsset")
	proto.RegisterType((*Pool)(nil), "elys.margin.Pool")
}

func init() { proto.RegisterFile("elys/margin/pool.proto", fileDescriptor_030dd771f91c2bb4) }

var fileDescriptor_030dd771f91c2bb4 = []byte{
	// 541 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xcb, 0x6e, 0xd3, 0x4c,
	0x14, 0xc7, 0xe3, 0x2f, 0xf9, 0x72, 0x99, 0x70, 0xeb, 0xb4, 0x0d, 0x16, 0x0b, 0x27, 0xaa, 0x10,
	0x0a, 0x42, 0xb5, 0x25, 0x78, 0x02, 0xdc, 0x52, 0x25, 0x52, 0x17, 0x91, 0xbb, 0x40, 0x62, 0xc1,
	0x30, 0xb6, 0xa7, 0xb6, 0x95, 0xb1, 0x8f, 0xe5, 0x99, 0xa8, 0xe4, 0x2d, 0x78, 0xac, 0xb2, 0xeb,
	0x12, 0xb1, 0xa8, 0x50, 0xb2, 0xe6, 0x1d, 0xd0, 0x4c, 0x6c, 0x6a, 0x2e, 0x0b, 0xe4, 0x55, 0x32,
	0xc7, 0x73, 0x7e, 0xe7, 0x7f, 0xfe, 0x73, 0x66, 0xd0, 0x88, 0xf1, 0xb5, 0x70, 0x52, 0x5a, 0x44,
	0x49, 0xe6, 0xe4, 0x00, 0xdc, 0xce, 0x0b, 0x90, 0x80, 0x87, 0x2a, 0x6e, 0xef, 0xe2, 0x4f, 0x0e,
	0x22, 0x88, 0x40, 0xc7, 0x1d, 0xf5, 0x6f, 0xb7, 0xe5, 0xe8, 0x73, 0x07, 0x0d, 0x16, 0x00, 0xfc,
	0xb5, 0x10, 0x4c, 0xe2, 0x05, 0x1a, 0xf2, 0x84, 0xfa, 0x09, 0x4f, 0x64, 0xc2, 0x84, 0x69, 0x4c,
	0x8c, 0xe9, 0xc0, 0xb5, 0xaf, 0x6f, 0xc7, 0xad, 0xaf, 0xb7, 0xe3, 0x67, 0x51, 0x22, 0xe3, 0x95,
	0x6f, 0x07, 0x90, 0x3a, 0x01, 0x88, 0x14, 0x44, 0xf9, 0x73, 0x2c, 0xc2, 0xa5, 0x23, 0xd7, 0x39,
	0x13, 0xf6, 0x3c, 0x93, 0x5e, 0x1d, 0x81, 0x67, 0xa8, 0x17, 0xac, 0x84, 0x84, 0x70, 0x6d, 0xfe,
	0xd7, 0x88, 0x56, 0xa5, 0xe3, 0x4b, 0xf4, 0x58, 0xd2, 0x25, 0x23, 0x79, 0x01, 0x97, 0x89, 0x24,
	0x75, 0x9d, 0xed, 0x46, 0xe4, 0x43, 0x85, 0x5b, 0x68, 0xda, 0x79, 0x4d, 0xf1, 0x7b, 0xb4, 0x5f,
	0xaf, 0x53, 0xa9, 0xef, 0x34, 0xaa, 0xb1, 0x77, 0x57, 0xe3, 0xa4, 0xec, 0xe3, 0x02, 0xdd, 0xa7,
	0xca, 0x6c, 0xe2, 0x53, 0x4e, 0xb3, 0x80, 0x99, 0xff, 0x37, 0x22, 0xdf, 0xd3, 0x10, 0x77, 0xc7,
	0xc0, 0x3e, 0x3a, 0xf4, 0x39, 0x04, 0x4b, 0xe2, 0x43, 0x51, 0xc0, 0x15, 0x49, 0x32, 0xc9, 0x0a,
	0x26, 0xa4, 0xd9, 0x6d, 0x04, 0xdf, 0xd7, 0x30, 0x57, 0xb3, 0xe6, 0x25, 0x0a, 0x8f, 0xd1, 0x70,
	0x27, 0x3c, 0x64, 0x19, 0xa4, 0x66, 0x4f, 0x91, 0x3d, 0xa4, 0x43, 0xa7, 0x2a, 0x72, 0xf4, 0xbd,
	0x8d, 0x3a, 0x6a, 0x96, 0xb0, 0x85, 0x86, 0x34, 0x4d, 0x89, 0x9a, 0x44, 0x92, 0x84, 0x7a, 0x8c,
	0x3a, 0xde, 0x80, 0xa6, 0xa9, 0xfa, 0x3a, 0x0f, 0xf1, 0x19, 0xea, 0xc6, 0x8c, 0x72, 0x19, 0x37,
	0x98, 0x89, 0x53, 0x16, 0x78, 0x65, 0x36, 0x36, 0x51, 0x8f, 0x65, 0xd4, 0xe7, 0x2c, 0xd4, 0x23,
	0xd0, 0xf7, 0xaa, 0x25, 0x1e, 0xa1, 0x6e, 0xc0, 0x41, 0xb0, 0x50, 0x9f, 0x5b, 0xdf, 0x2b, 0x57,
	0xf8, 0x03, 0x3a, 0xf8, 0xcd, 0x21, 0x52, 0x50, 0xd9, 0xe4, 0x0c, 0x94, 0x0e, 0xec, 0xff, 0xe2,
	0x90, 0x47, 0x25, 0xc3, 0x67, 0xe8, 0x91, 0xee, 0x5b, 0xfb, 0x22, 0x08, 0x87, 0x2c, 0x32, 0xbb,
	0x93, 0xf6, 0x74, 0xf8, 0x72, 0x64, 0xd7, 0xae, 0xa3, 0xfd, 0xf3, 0xd2, 0xb9, 0x1d, 0x55, 0xd5,
	0x7b, 0x90, 0x57, 0x01, 0x71, 0x0e, 0x59, 0x84, 0x67, 0x68, 0xaf, 0xce, 0x11, 0x31, 0x14, 0xd2,
	0xec, 0xfd, 0x03, 0xe8, 0xe1, 0x1d, 0xe8, 0x42, 0x25, 0xe1, 0xb7, 0xe8, 0x39, 0xa7, 0x42, 0x92,
	0x98, 0x25, 0x51, 0x2c, 0xc9, 0xdf, 0xfa, 0x27, 0x01, 0xa4, 0xf9, 0x4a, 0xb2, 0xd0, 0xec, 0x4f,
	0x8c, 0x69, 0xdb, 0x7b, 0xaa, 0x12, 0x66, 0x7a, 0xbf, 0xfb, 0x47, 0x8b, 0x27, 0xe5, 0x5e, 0xf7,
	0xcd, 0xf5, 0xc6, 0x32, 0x6e, 0x36, 0x96, 0xf1, 0x6d, 0x63, 0x19, 0x9f, 0xb6, 0x56, 0xeb, 0x66,
	0x6b, 0xb5, 0xbe, 0x6c, 0xad, 0xd6, 0xbb, 0x17, 0x35, 0x03, 0x95, 0xd6, 0xe3, 0x8c, 0xc9, 0x2b,
	0x28, 0x96, 0x7a, 0xe1, 0x7c, 0xac, 0x9e, 0x2a, 0xed, 0xa4, 0xdf, 0xd5, 0x2f, 0xd1, 0xab, 0x1f,
	0x01, 0x00, 0x00, 0xff, 0xff, 0xb9, 0xe1, 0x33, 0x02, 0xc6, 0x04, 0x00, 0x00,
}

func (m *PoolAsset) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PoolAsset) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PoolAsset) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AssetDenom) > 0 {
		i -= len(m.AssetDenom)
		copy(dAtA[i:], m.AssetDenom)
		i = encodeVarintPool(dAtA, i, uint64(len(m.AssetDenom)))
		i--
		dAtA[i] = 0x3a
	}
	{
		size := m.BlockBorrowInterest.Size()
		i -= size
		if _, err := m.BlockBorrowInterest.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.AssetBalance.Size()
		i -= size
		if _, err := m.AssetBalance.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.TakeProfitCustody.Size()
		i -= size
		if _, err := m.TakeProfitCustody.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.TakeProfitLiabilities.Size()
		i -= size
		if _, err := m.TakeProfitLiabilities.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Custody.Size()
		i -= size
		if _, err := m.Custody.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.Liabilities.Size()
		i -= size
		if _, err := m.Liabilities.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *Pool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Pool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Pool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.LastHeightBorrowInterestRateComputed != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.LastHeightBorrowInterestRateComputed))
		i--
		dAtA[i] = 0x40
	}
	if len(m.PoolAssetsShort) > 0 {
		for iNdEx := len(m.PoolAssetsShort) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PoolAssetsShort[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPool(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.PoolAssetsLong) > 0 {
		for iNdEx := len(m.PoolAssetsLong) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PoolAssetsLong[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPool(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	{
		size := m.BorrowInterestRate.Size()
		i -= size
		if _, err := m.BorrowInterestRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.Closed {
		i--
		if m.Closed {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if m.Enabled {
		i--
		if m.Enabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	{
		size := m.Health.Size()
		i -= size
		if _, err := m.Health.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.AmmPoolId != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.AmmPoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintPool(dAtA []byte, offset int, v uint64) int {
	offset -= sovPool(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PoolAsset) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Liabilities.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.Custody.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.TakeProfitLiabilities.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.TakeProfitCustody.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.AssetBalance.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.BlockBorrowInterest.Size()
	n += 1 + l + sovPool(uint64(l))
	l = len(m.AssetDenom)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	return n
}

func (m *Pool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AmmPoolId != 0 {
		n += 1 + sovPool(uint64(m.AmmPoolId))
	}
	l = m.Health.Size()
	n += 1 + l + sovPool(uint64(l))
	if m.Enabled {
		n += 2
	}
	if m.Closed {
		n += 2
	}
	l = m.BorrowInterestRate.Size()
	n += 1 + l + sovPool(uint64(l))
	if len(m.PoolAssetsLong) > 0 {
		for _, e := range m.PoolAssetsLong {
			l = e.Size()
			n += 1 + l + sovPool(uint64(l))
		}
	}
	if len(m.PoolAssetsShort) > 0 {
		for _, e := range m.PoolAssetsShort {
			l = e.Size()
			n += 1 + l + sovPool(uint64(l))
		}
	}
	if m.LastHeightBorrowInterestRateComputed != 0 {
		n += 1 + sovPool(uint64(m.LastHeightBorrowInterestRateComputed))
	}
	return n
}

func sovPool(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPool(x uint64) (n int) {
	return sovPool(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PoolAsset) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPool
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
			return fmt.Errorf("proto: PoolAsset: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PoolAsset: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liabilities", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Liabilities.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Custody", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Custody.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TakeProfitLiabilities", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TakeProfitLiabilities.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TakeProfitCustody", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TakeProfitCustody.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetBalance", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AssetBalance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockBorrowInterest", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BlockBorrowInterest.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AssetDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPool(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPool
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
func (m *Pool) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPool
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
			return fmt.Errorf("proto: Pool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Pool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AmmPoolId", wireType)
			}
			m.AmmPoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AmmPoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Health", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Health.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Enabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Enabled = bool(v != 0)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Closed", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Closed = bool(v != 0)
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BorrowInterestRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.BorrowInterestRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolAssetsLong", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PoolAssetsLong = append(m.PoolAssetsLong, PoolAsset{})
			if err := m.PoolAssetsLong[len(m.PoolAssetsLong)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolAssetsShort", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
				return ErrInvalidLengthPool
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPool
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PoolAssetsShort = append(m.PoolAssetsShort, PoolAsset{})
			if err := m.PoolAssetsShort[len(m.PoolAssetsShort)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LastHeightBorrowInterestRateComputed", wireType)
			}
			m.LastHeightBorrowInterestRateComputed = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LastHeightBorrowInterestRateComputed |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipPool(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPool
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
func skipPool(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPool
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
					return 0, ErrIntOverflowPool
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
					return 0, ErrIntOverflowPool
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
				return 0, ErrInvalidLengthPool
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPool
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPool
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPool        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPool          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPool = fmt.Errorf("proto: unexpected end of group")
)
