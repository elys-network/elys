// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/amm/pool.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
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

type Pool struct {
	PoolId            uint64                `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	Address           string                `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	PoolParams        PoolParams            `protobuf:"bytes,3,opt,name=pool_params,json=poolParams,proto3" json:"pool_params"`
	TotalShares       types.Coin            `protobuf:"bytes,4,opt,name=total_shares,json=totalShares,proto3" json:"total_shares"`
	PoolAssets        []PoolAsset           `protobuf:"bytes,5,rep,name=pool_assets,json=poolAssets,proto3" json:"pool_assets"`
	TotalWeight       cosmossdk_io_math.Int `protobuf:"bytes,6,opt,name=total_weight,json=totalWeight,proto3,customtype=cosmossdk.io/math.Int" json:"total_weight"`
	RebalanceTreasury string                `protobuf:"bytes,7,opt,name=rebalance_treasury,json=rebalanceTreasury,proto3" json:"rebalance_treasury,omitempty"`
}

func (m *Pool) Reset()         { *m = Pool{} }
func (m *Pool) String() string { return proto.CompactTextString(m) }
func (*Pool) ProtoMessage()    {}
func (*Pool) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ac3be9a215271f9, []int{0}
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

func (m *Pool) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *Pool) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Pool) GetPoolParams() PoolParams {
	if m != nil {
		return m.PoolParams
	}
	return PoolParams{}
}

func (m *Pool) GetTotalShares() types.Coin {
	if m != nil {
		return m.TotalShares
	}
	return types.Coin{}
}

func (m *Pool) GetPoolAssets() []PoolAsset {
	if m != nil {
		return m.PoolAssets
	}
	return nil
}

func (m *Pool) GetRebalanceTreasury() string {
	if m != nil {
		return m.RebalanceTreasury
	}
	return ""
}

type PoolExtraInfo struct {
	Tvl          cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=tvl,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"tvl"`
	LpTokenPrice cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=lp_token_price,json=lpTokenPrice,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"lp_token_price"`
}

func (m *PoolExtraInfo) Reset()         { *m = PoolExtraInfo{} }
func (m *PoolExtraInfo) String() string { return proto.CompactTextString(m) }
func (*PoolExtraInfo) ProtoMessage()    {}
func (*PoolExtraInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ac3be9a215271f9, []int{1}
}
func (m *PoolExtraInfo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PoolExtraInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PoolExtraInfo.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PoolExtraInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolExtraInfo.Merge(m, src)
}
func (m *PoolExtraInfo) XXX_Size() int {
	return m.Size()
}
func (m *PoolExtraInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolExtraInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PoolExtraInfo proto.InternalMessageInfo

type OraclePoolSlippageTrack struct {
	PoolId    uint64                                   `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	Timestamp uint64                                   `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Tracked   github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,3,rep,name=tracked,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"tracked"`
}

func (m *OraclePoolSlippageTrack) Reset()         { *m = OraclePoolSlippageTrack{} }
func (m *OraclePoolSlippageTrack) String() string { return proto.CompactTextString(m) }
func (*OraclePoolSlippageTrack) ProtoMessage()    {}
func (*OraclePoolSlippageTrack) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ac3be9a215271f9, []int{2}
}
func (m *OraclePoolSlippageTrack) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OraclePoolSlippageTrack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OraclePoolSlippageTrack.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OraclePoolSlippageTrack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OraclePoolSlippageTrack.Merge(m, src)
}
func (m *OraclePoolSlippageTrack) XXX_Size() int {
	return m.Size()
}
func (m *OraclePoolSlippageTrack) XXX_DiscardUnknown() {
	xxx_messageInfo_OraclePoolSlippageTrack.DiscardUnknown(m)
}

var xxx_messageInfo_OraclePoolSlippageTrack proto.InternalMessageInfo

func (m *OraclePoolSlippageTrack) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *OraclePoolSlippageTrack) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *OraclePoolSlippageTrack) GetTracked() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Tracked
	}
	return nil
}

func init() {
	proto.RegisterType((*Pool)(nil), "elys.amm.Pool")
	proto.RegisterType((*PoolExtraInfo)(nil), "elys.amm.PoolExtraInfo")
	proto.RegisterType((*OraclePoolSlippageTrack)(nil), "elys.amm.OraclePoolSlippageTrack")
}

func init() { proto.RegisterFile("elys/amm/pool.proto", fileDescriptor_3ac3be9a215271f9) }

var fileDescriptor_3ac3be9a215271f9 = []byte{
	// 547 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x4d, 0x6e, 0xd3, 0x40,
	0x14, 0x8e, 0x9b, 0xd0, 0x90, 0x49, 0x41, 0x62, 0x5a, 0x54, 0x37, 0x80, 0x13, 0x85, 0x4d, 0x36,
	0xb1, 0x69, 0x11, 0x1b, 0xd8, 0x80, 0x81, 0x45, 0x24, 0x24, 0x22, 0x37, 0x12, 0x12, 0x1b, 0x6b,
	0x62, 0x3f, 0x1c, 0x2b, 0xb6, 0x67, 0x34, 0x33, 0xfd, 0xc9, 0x2d, 0xe0, 0x1a, 0x5c, 0x80, 0x2b,
	0x74, 0xd9, 0x25, 0x62, 0x51, 0x50, 0xb2, 0xe0, 0x1a, 0x68, 0xc6, 0xe3, 0x92, 0x08, 0x81, 0x58,
	0xd9, 0xef, 0x7d, 0xf3, 0xfd, 0xbc, 0xf9, 0x41, 0xbb, 0x90, 0x2d, 0x84, 0x47, 0xf2, 0xdc, 0x63,
	0x94, 0x66, 0x2e, 0xe3, 0x54, 0x52, 0x7c, 0x53, 0x35, 0x5d, 0x92, 0xe7, 0x9d, 0xce, 0x06, 0x1c,
	0x32, 0xc2, 0x49, 0x2e, 0xca, 0x55, 0x9d, 0x83, 0x4d, 0x8c, 0x08, 0x01, 0xd2, 0x40, 0x7b, 0x09,
	0x4d, 0xa8, 0xfe, 0xf5, 0xd4, 0x9f, 0xe9, 0x3a, 0x11, 0x15, 0x39, 0x15, 0xde, 0x94, 0x08, 0xf0,
	0x4e, 0x0f, 0xa7, 0x20, 0xc9, 0xa1, 0x17, 0xd1, 0xb4, 0xa8, 0x04, 0x4b, 0x3c, 0x2c, 0x89, 0x65,
	0x51, 0x42, 0xfd, 0x9f, 0x5b, 0xa8, 0x31, 0xa6, 0x34, 0xc3, 0xfb, 0xa8, 0xa9, 0xdd, 0xd2, 0xd8,
	0xb6, 0x7a, 0xd6, 0xa0, 0x11, 0x6c, 0xab, 0x72, 0x14, 0x63, 0x1b, 0x35, 0x49, 0x1c, 0x73, 0x10,
	0xc2, 0xde, 0xea, 0x59, 0x83, 0x56, 0x50, 0x95, 0xf8, 0x19, 0x6a, 0xaf, 0x85, 0xb7, 0xeb, 0x3d,
	0x6b, 0xd0, 0x3e, 0xda, 0x73, 0xab, 0x19, 0x5d, 0xa5, 0x3b, 0xd6, 0x98, 0xdf, 0xb8, 0xb8, 0xea,
	0xd6, 0x02, 0xc4, 0xae, 0x3b, 0xd8, 0x47, 0x3b, 0x92, 0x4a, 0x92, 0x85, 0x62, 0x46, 0x38, 0x08,
	0xbb, 0xa1, 0xd9, 0x07, 0xae, 0x49, 0xa7, 0x46, 0x71, 0xcd, 0x28, 0xee, 0x4b, 0x9a, 0x16, 0x46,
	0xa2, 0xad, 0x49, 0xc7, 0x9a, 0x83, 0x9f, 0x9a, 0x00, 0x7a, 0x87, 0x84, 0x7d, 0xa3, 0x57, 0x1f,
	0xb4, 0x8f, 0x76, 0x37, 0x03, 0xbc, 0x50, 0xd8, 0xba, 0xbf, 0x6e, 0x08, 0xfc, 0xbc, 0xf2, 0x3f,
	0x83, 0x34, 0x99, 0x49, 0x7b, 0x5b, 0xcd, 0xe6, 0x3f, 0x50, 0xeb, 0xbe, 0x5d, 0x75, 0xef, 0x96,
	0x31, 0x44, 0x3c, 0x77, 0x53, 0xea, 0xe5, 0x44, 0xce, 0xdc, 0x51, 0x21, 0x8d, 0xfb, 0x3b, 0xcd,
	0xc0, 0x43, 0x84, 0x39, 0x4c, 0x49, 0x46, 0x8a, 0x08, 0x42, 0xc9, 0x81, 0x88, 0x13, 0xbe, 0xb0,
	0x9b, 0x7a, 0x8f, 0xee, 0x5c, 0x23, 0x13, 0x03, 0xf4, 0x3f, 0x59, 0xe8, 0x96, 0x0a, 0xf4, 0xfa,
	0x5c, 0x72, 0x32, 0x2a, 0x3e, 0x50, 0xfc, 0x04, 0xd5, 0xe5, 0x69, 0xa6, 0xb7, 0xbb, 0xe5, 0x3f,
	0x34, 0xce, 0xf7, 0xfe, 0x74, 0x7e, 0x03, 0x09, 0x89, 0x16, 0xaf, 0x20, 0x0a, 0xd4, 0x7a, 0x3c,
	0x42, 0xb7, 0x33, 0x16, 0x4a, 0x3a, 0x87, 0x22, 0x64, 0x3c, 0x8d, 0xa0, 0x3c, 0x97, 0xff, 0x53,
	0xd8, 0xc9, 0xd8, 0x44, 0x31, 0xc7, 0x8a, 0xd8, 0xff, 0x62, 0xa1, 0xfd, 0xb7, 0x9c, 0x44, 0x19,
	0xa8, 0x64, 0xc7, 0x59, 0xca, 0x18, 0x49, 0x60, 0xc2, 0x49, 0x34, 0xff, 0xfb, 0x85, 0xb8, 0x8f,
	0x5a, 0x32, 0xcd, 0x41, 0x48, 0x92, 0x33, 0x6d, 0xdd, 0x08, 0x7e, 0x37, 0x30, 0xa0, 0xa6, 0x54,
	0x7c, 0x88, 0xed, 0xba, 0x3e, 0x8f, 0x7f, 0x1c, 0xe9, 0x23, 0x95, 0xf8, 0xf3, 0xf7, 0xee, 0x20,
	0x49, 0xe5, 0xec, 0x64, 0xea, 0x46, 0x34, 0x37, 0xb7, 0xd3, 0x7c, 0x86, 0x22, 0x9e, 0x7b, 0x72,
	0xc1, 0x40, 0x68, 0x82, 0x08, 0x2a, 0x6d, 0xdf, 0xbf, 0x58, 0x3a, 0xd6, 0xe5, 0xd2, 0xb1, 0x7e,
	0x2c, 0x1d, 0xeb, 0xe3, 0xca, 0xa9, 0x5d, 0xae, 0x9c, 0xda, 0xd7, 0x95, 0x53, 0x7b, 0xbf, 0x2e,
	0xa6, 0x6e, 0xc2, 0xb0, 0x00, 0x79, 0x46, 0xf9, 0x5c, 0x17, 0xde, 0xb9, 0x7e, 0x57, 0x5a, 0x72,
	0xba, 0xad, 0x9f, 0xc0, 0xe3, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x84, 0x5b, 0xd7, 0x7f, 0xab,
	0x03, 0x00, 0x00,
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
	if len(m.RebalanceTreasury) > 0 {
		i -= len(m.RebalanceTreasury)
		copy(dAtA[i:], m.RebalanceTreasury)
		i = encodeVarintPool(dAtA, i, uint64(len(m.RebalanceTreasury)))
		i--
		dAtA[i] = 0x3a
	}
	{
		size := m.TotalWeight.Size()
		i -= size
		if _, err := m.TotalWeight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if len(m.PoolAssets) > 0 {
		for iNdEx := len(m.PoolAssets) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PoolAssets[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPool(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	{
		size, err := m.TotalShares.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size, err := m.PoolParams.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintPool(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x12
	}
	if m.PoolId != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *PoolExtraInfo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PoolExtraInfo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PoolExtraInfo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.LpTokenPrice.Size()
		i -= size
		if _, err := m.LpTokenPrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.Tvl.Size()
		i -= size
		if _, err := m.Tvl.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *OraclePoolSlippageTrack) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OraclePoolSlippageTrack) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OraclePoolSlippageTrack) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Tracked) > 0 {
		for iNdEx := len(m.Tracked) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Tracked[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPool(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.Timestamp != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.Timestamp))
		i--
		dAtA[i] = 0x10
	}
	if m.PoolId != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.PoolId))
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
func (m *Pool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovPool(uint64(m.PoolId))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	l = m.PoolParams.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.TotalShares.Size()
	n += 1 + l + sovPool(uint64(l))
	if len(m.PoolAssets) > 0 {
		for _, e := range m.PoolAssets {
			l = e.Size()
			n += 1 + l + sovPool(uint64(l))
		}
	}
	l = m.TotalWeight.Size()
	n += 1 + l + sovPool(uint64(l))
	l = len(m.RebalanceTreasury)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	return n
}

func (m *PoolExtraInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Tvl.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.LpTokenPrice.Size()
	n += 1 + l + sovPool(uint64(l))
	return n
}

func (m *OraclePoolSlippageTrack) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovPool(uint64(m.PoolId))
	}
	if m.Timestamp != 0 {
		n += 1 + sovPool(uint64(m.Timestamp))
	}
	if len(m.Tracked) > 0 {
		for _, e := range m.Tracked {
			l = e.Size()
			n += 1 + l + sovPool(uint64(l))
		}
	}
	return n
}

func sovPool(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPool(x uint64) (n int) {
	return sovPool(uint64((x << 1) ^ uint64((int64(x) >> 63))))
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
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolParams", wireType)
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
			if err := m.PoolParams.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalShares", wireType)
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
			if err := m.TotalShares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolAssets", wireType)
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
			m.PoolAssets = append(m.PoolAssets, PoolAsset{})
			if err := m.PoolAssets[len(m.PoolAssets)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalWeight", wireType)
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
			if err := m.TotalWeight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RebalanceTreasury", wireType)
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
			m.RebalanceTreasury = string(dAtA[iNdEx:postIndex])
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
func (m *PoolExtraInfo) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: PoolExtraInfo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PoolExtraInfo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tvl", wireType)
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
			if err := m.Tvl.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LpTokenPrice", wireType)
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
			if err := m.LpTokenPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *OraclePoolSlippageTrack) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OraclePoolSlippageTrack: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OraclePoolSlippageTrack: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			m.Timestamp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Timestamp |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Tracked", wireType)
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
			m.Tracked = append(m.Tracked, types.Coin{})
			if err := m.Tracked[len(m.Tracked)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
