// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/amm/pool.proto

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

type LegacyPool struct {
	PoolId            uint64                                 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	Address           string                                 `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	PoolParams        LegacyPoolParams                       `protobuf:"bytes,3,opt,name=pool_params,json=poolParams,proto3" json:"pool_params"`
	TotalShares       types.Coin                             `protobuf:"bytes,4,opt,name=total_shares,json=totalShares,proto3" json:"total_shares"`
	PoolAssets        []PoolAsset                            `protobuf:"bytes,5,rep,name=pool_assets,json=poolAssets,proto3" json:"pool_assets"`
	TotalWeight       github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=total_weight,json=totalWeight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"total_weight"`
	RebalanceTreasury string                                 `protobuf:"bytes,7,opt,name=rebalance_treasury,json=rebalanceTreasury,proto3" json:"rebalance_treasury,omitempty"`
}

func (m *LegacyPool) Reset()         { *m = LegacyPool{} }
func (m *LegacyPool) String() string { return proto.CompactTextString(m) }
func (*LegacyPool) ProtoMessage()    {}
func (*LegacyPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ac3be9a215271f9, []int{0}
}
func (m *LegacyPool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LegacyPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LegacyPool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LegacyPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LegacyPool.Merge(m, src)
}
func (m *LegacyPool) XXX_Size() int {
	return m.Size()
}
func (m *LegacyPool) XXX_DiscardUnknown() {
	xxx_messageInfo_LegacyPool.DiscardUnknown(m)
}

var xxx_messageInfo_LegacyPool proto.InternalMessageInfo

func (m *LegacyPool) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *LegacyPool) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *LegacyPool) GetPoolParams() LegacyPoolParams {
	if m != nil {
		return m.PoolParams
	}
	return LegacyPoolParams{}
}

func (m *LegacyPool) GetTotalShares() types.Coin {
	if m != nil {
		return m.TotalShares
	}
	return types.Coin{}
}

func (m *LegacyPool) GetPoolAssets() []PoolAsset {
	if m != nil {
		return m.PoolAssets
	}
	return nil
}

func (m *LegacyPool) GetRebalanceTreasury() string {
	if m != nil {
		return m.RebalanceTreasury
	}
	return ""
}

type Pool struct {
	PoolId            uint64                                 `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	Address           string                                 `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	PoolParams        PoolParams                             `protobuf:"bytes,3,opt,name=pool_params,json=poolParams,proto3" json:"pool_params"`
	TotalShares       types.Coin                             `protobuf:"bytes,4,opt,name=total_shares,json=totalShares,proto3" json:"total_shares"`
	PoolAssets        []PoolAsset                            `protobuf:"bytes,5,rep,name=pool_assets,json=poolAssets,proto3" json:"pool_assets"`
	TotalWeight       github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=total_weight,json=totalWeight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"total_weight"`
	RebalanceTreasury string                                 `protobuf:"bytes,7,opt,name=rebalance_treasury,json=rebalanceTreasury,proto3" json:"rebalance_treasury,omitempty"`
}

func (m *Pool) Reset()         { *m = Pool{} }
func (m *Pool) String() string { return proto.CompactTextString(m) }
func (*Pool) ProtoMessage()    {}
func (*Pool) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ac3be9a215271f9, []int{1}
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
	Tvl          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,1,opt,name=tvl,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"tvl"`
	LpTokenPrice github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=lp_token_price,json=lpTokenPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"lp_token_price"`
}

func (m *PoolExtraInfo) Reset()         { *m = PoolExtraInfo{} }
func (m *PoolExtraInfo) String() string { return proto.CompactTextString(m) }
func (*PoolExtraInfo) ProtoMessage()    {}
func (*PoolExtraInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ac3be9a215271f9, []int{2}
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
	return fileDescriptor_3ac3be9a215271f9, []int{3}
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
	proto.RegisterType((*LegacyPool)(nil), "elys.amm.LegacyPool")
	proto.RegisterType((*Pool)(nil), "elys.amm.Pool")
	proto.RegisterType((*PoolExtraInfo)(nil), "elys.amm.PoolExtraInfo")
	proto.RegisterType((*OraclePoolSlippageTrack)(nil), "elys.amm.OraclePoolSlippageTrack")
}

func init() { proto.RegisterFile("elys/amm/pool.proto", fileDescriptor_3ac3be9a215271f9) }

var fileDescriptor_3ac3be9a215271f9 = []byte{
	// 561 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x8e, 0x9b, 0x90, 0xd0, 0x4d, 0x41, 0x62, 0x5b, 0xa9, 0x6e, 0x84, 0x9c, 0x28, 0x07, 0x94,
	0x4b, 0x6c, 0x5a, 0x6e, 0x70, 0xa1, 0x01, 0x0e, 0x91, 0x90, 0x08, 0x6e, 0x24, 0x24, 0x2e, 0xd6,
	0xc6, 0x1e, 0x1c, 0x2b, 0xb6, 0x77, 0xb5, 0xbb, 0xfd, 0xc9, 0x9d, 0x07, 0xe0, 0x2d, 0x90, 0x78,
	0x01, 0x5e, 0xa1, 0xc7, 0x1e, 0x11, 0x87, 0x82, 0x92, 0x97, 0xe0, 0x88, 0x76, 0xbd, 0x69, 0x12,
	0x89, 0x22, 0xca, 0xb9, 0x27, 0x7b, 0xe6, 0xdb, 0xf9, 0xe6, 0x9b, 0x6f, 0x56, 0x8b, 0xb6, 0x21,
	0x9d, 0x0a, 0x8f, 0x64, 0x99, 0xc7, 0x28, 0x4d, 0x5d, 0xc6, 0xa9, 0xa4, 0xf8, 0xae, 0x4a, 0xba,
	0x24, 0xcb, 0x1a, 0x8d, 0x35, 0x38, 0x60, 0x84, 0x93, 0x4c, 0x14, 0xa7, 0x1a, 0x7b, 0xeb, 0x18,
	0x11, 0x02, 0xa4, 0x81, 0x76, 0x62, 0x1a, 0x53, 0xfd, 0xeb, 0xa9, 0x3f, 0x93, 0x75, 0x42, 0x2a,
	0x32, 0x2a, 0xbc, 0x11, 0x11, 0xe0, 0x9d, 0xec, 0x8f, 0x40, 0x92, 0x7d, 0x2f, 0xa4, 0x49, 0xbe,
	0x20, 0x2c, 0xf0, 0xa0, 0x28, 0x2c, 0x82, 0x02, 0x6a, 0x7f, 0x2c, 0x23, 0xf4, 0x1a, 0x62, 0x12,
	0x4e, 0x07, 0x94, 0xa6, 0x78, 0x17, 0xd5, 0x74, 0xcf, 0x24, 0xb2, 0xad, 0x96, 0xd5, 0xa9, 0xf8,
	0x55, 0x15, 0xf6, 0x23, 0x6c, 0xa3, 0x1a, 0x89, 0x22, 0x0e, 0x42, 0xd8, 0x1b, 0x2d, 0xab, 0xb3,
	0xe9, 0x2f, 0x42, 0x7c, 0x88, 0xea, 0x2b, 0x23, 0xd8, 0xe5, 0x96, 0xd5, 0xa9, 0x1f, 0x34, 0xdc,
	0xc5, 0xa4, 0xee, 0x92, 0x7d, 0xa0, 0x4f, 0xf4, 0x2a, 0xe7, 0x97, 0xcd, 0x92, 0x8f, 0xd8, 0x55,
	0x06, 0xf7, 0xd0, 0x96, 0xa4, 0x92, 0xa4, 0x81, 0x18, 0x13, 0x0e, 0xc2, 0xae, 0x68, 0x8e, 0x3d,
	0xd7, 0x28, 0x55, 0x63, 0xb9, 0x66, 0x2c, 0xf7, 0x05, 0x4d, 0x72, 0x43, 0x51, 0xd7, 0x45, 0x47,
	0xba, 0x06, 0x3f, 0x35, 0x32, 0xb4, 0x5b, 0xc2, 0xbe, 0xd3, 0x2a, 0x77, 0xea, 0x07, 0xdb, 0x4b,
	0x19, 0x4a, 0xc0, 0xa1, 0xc2, 0x56, 0xfb, 0xeb, 0x84, 0xc0, 0x6f, 0x17, 0xfd, 0x4f, 0x21, 0x89,
	0xc7, 0xd2, 0xae, 0xaa, 0x09, 0x7b, 0xae, 0x3a, 0xf7, 0xfd, 0xb2, 0xf9, 0x28, 0x4e, 0xe4, 0xf8,
	0x78, 0xe4, 0x86, 0x34, 0x33, 0xde, 0x99, 0x4f, 0x57, 0x44, 0x13, 0x4f, 0x4e, 0x19, 0x08, 0xb7,
	0x9f, 0x4b, 0x23, 0xe7, 0x9d, 0xa6, 0xc0, 0x5d, 0x84, 0x39, 0x8c, 0x48, 0x4a, 0xf2, 0x10, 0x02,
	0xc9, 0x81, 0x88, 0x63, 0x3e, 0xb5, 0x6b, 0xda, 0xba, 0x07, 0x57, 0xc8, 0xd0, 0x00, 0xed, 0x5f,
	0x1b, 0xa8, 0xf2, 0xbf, 0x0b, 0x78, 0xf6, 0xa7, 0x05, 0xec, 0xac, 0x4f, 0x7e, 0x6b, 0xfd, 0xb5,
	0xd6, 0x7f, 0xb6, 0xd0, 0x3d, 0xa5, 0xf0, 0xd5, 0x99, 0xe4, 0xa4, 0x9f, 0x7f, 0xa0, 0xf8, 0x39,
	0x2a, 0xcb, 0x93, 0x54, 0xfb, 0x7f, 0x33, 0x29, 0x2f, 0x21, 0xf4, 0x55, 0x29, 0x1e, 0xa2, 0xfb,
	0x29, 0x0b, 0x24, 0x9d, 0x40, 0x1e, 0x30, 0x9e, 0x84, 0x50, 0xec, 0xec, 0xc6, 0x64, 0x5b, 0x29,
	0x1b, 0x2a, 0x92, 0x81, 0xe2, 0x68, 0x7f, 0xb5, 0xd0, 0xee, 0x1b, 0x4e, 0xc2, 0x14, 0x94, 0xde,
	0xa3, 0x34, 0x61, 0x8c, 0xc4, 0x30, 0xe4, 0x24, 0x9c, 0x5c, 0x7f, 0x6f, 0x1e, 0xa2, 0x4d, 0x99,
	0x64, 0x20, 0x24, 0xc9, 0x98, 0x56, 0x51, 0xf1, 0x97, 0x09, 0x0c, 0xa8, 0x26, 0x55, 0x3d, 0x44,
	0x76, 0x59, 0xaf, 0xed, 0x2f, 0x9b, 0x7f, 0xac, 0xc4, 0x7f, 0xf9, 0xd1, 0xec, 0xfc, 0x83, 0x78,
	0x55, 0x20, 0xfc, 0x05, 0x77, 0xaf, 0x77, 0x3e, 0x73, 0xac, 0x8b, 0x99, 0x63, 0xfd, 0x9c, 0x39,
	0xd6, 0xa7, 0xb9, 0x53, 0xba, 0x98, 0x3b, 0xa5, 0x6f, 0x73, 0xa7, 0xf4, 0x7e, 0x95, 0x4c, 0x5d,
	0x98, 0x6e, 0x0e, 0xf2, 0x94, 0xf2, 0x89, 0x0e, 0xbc, 0x33, 0xfd, 0x0a, 0x6a, 0xca, 0x51, 0x55,
	0x3f, 0x58, 0x4f, 0x7e, 0x07, 0x00, 0x00, 0xff, 0xff, 0xb9, 0x95, 0x53, 0xf5, 0x59, 0x05, 0x00,
	0x00,
}

func (m *LegacyPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LegacyPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LegacyPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *LegacyPool) Size() (n int) {
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
func (m *LegacyPool) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: LegacyPool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LegacyPool: illegal tag %d (wire type %d)", fieldNum, wire)
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
