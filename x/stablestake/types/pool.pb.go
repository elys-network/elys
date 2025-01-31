// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/stablestake/pool.proto

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
	DepositDenom         string                      `protobuf:"bytes,1,opt,name=deposit_denom,json=depositDenom,proto3" json:"deposit_denom,omitempty"`
	RedemptionRate       cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=redemption_rate,json=redemptionRate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"redemption_rate"`
	InterestRate         cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=interest_rate,json=interestRate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"interest_rate"`
	InterestRateMax      cosmossdk_io_math.LegacyDec `protobuf:"bytes,4,opt,name=interest_rate_max,json=interestRateMax,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"interest_rate_max"`
	InterestRateMin      cosmossdk_io_math.LegacyDec `protobuf:"bytes,5,opt,name=interest_rate_min,json=interestRateMin,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"interest_rate_min"`
	InterestRateIncrease cosmossdk_io_math.LegacyDec `protobuf:"bytes,6,opt,name=interest_rate_increase,json=interestRateIncrease,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"interest_rate_increase"`
	InterestRateDecrease cosmossdk_io_math.LegacyDec `protobuf:"bytes,7,opt,name=interest_rate_decrease,json=interestRateDecrease,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"interest_rate_decrease"`
	HealthGainFactor     cosmossdk_io_math.LegacyDec `protobuf:"bytes,8,opt,name=health_gain_factor,json=healthGainFactor,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"health_gain_factor"`
	TotalValue           cosmossdk_io_math.Int       `protobuf:"bytes,9,opt,name=total_value,json=totalValue,proto3,customtype=cosmossdk.io/math.Int" json:"total_value"`
	MaxLeverageRatio     cosmossdk_io_math.LegacyDec `protobuf:"bytes,10,opt,name=max_leverage_ratio,json=maxLeverageRatio,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"max_leverage_ratio"`
	PoolId               uint64                      `protobuf:"varint,11,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
}

func (m *Pool) Reset()         { *m = Pool{} }
func (m *Pool) String() string { return proto.CompactTextString(m) }
func (*Pool) ProtoMessage()    {}
func (*Pool) Descriptor() ([]byte, []int) {
	return fileDescriptor_f438111f72086e2f, []int{0}
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

func (m *Pool) GetDepositDenom() string {
	if m != nil {
		return m.DepositDenom
	}
	return ""
}

func (m *Pool) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

type AmmPool struct {
	Id               uint64                                   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	TotalLiabilities github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=total_liabilities,json=totalLiabilities,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"total_liabilities"`
}

func (m *AmmPool) Reset()         { *m = AmmPool{} }
func (m *AmmPool) String() string { return proto.CompactTextString(m) }
func (*AmmPool) ProtoMessage()    {}
func (*AmmPool) Descriptor() ([]byte, []int) {
	return fileDescriptor_f438111f72086e2f, []int{1}
}
func (m *AmmPool) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AmmPool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AmmPool.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AmmPool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AmmPool.Merge(m, src)
}
func (m *AmmPool) XXX_Size() int {
	return m.Size()
}
func (m *AmmPool) XXX_DiscardUnknown() {
	xxx_messageInfo_AmmPool.DiscardUnknown(m)
}

var xxx_messageInfo_AmmPool proto.InternalMessageInfo

func (m *AmmPool) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *AmmPool) GetTotalLiabilities() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.TotalLiabilities
	}
	return nil
}

func init() {
	proto.RegisterType((*Pool)(nil), "elys.stablestake.Pool")
	proto.RegisterType((*AmmPool)(nil), "elys.stablestake.AmmPool")
}

func init() { proto.RegisterFile("elys/stablestake/pool.proto", fileDescriptor_f438111f72086e2f) }

var fileDescriptor_f438111f72086e2f = []byte{
	// 560 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0x80, 0xe3, 0x34, 0x24, 0x74, 0xd3, 0x9f, 0xd4, 0x2a, 0xe0, 0xb6, 0x92, 0x13, 0x95, 0x4b,
	0x24, 0x14, 0xbb, 0x81, 0x27, 0x20, 0x44, 0xa0, 0xa0, 0x20, 0x21, 0x1f, 0x7a, 0xa8, 0x84, 0xac,
	0xb5, 0x3d, 0x38, 0xab, 0xd8, 0xbb, 0x91, 0x77, 0x1b, 0x9c, 0xb7, 0x40, 0xe2, 0xc2, 0x33, 0x70,
	0xe6, 0x21, 0x7a, 0xac, 0x38, 0x21, 0x0e, 0x05, 0x25, 0x2f, 0x82, 0x76, 0xbd, 0xa8, 0x89, 0xca,
	0xc9, 0xea, 0xc9, 0x9e, 0x9d, 0x9d, 0xef, 0x1b, 0x5b, 0x9a, 0x41, 0x27, 0x90, 0x2c, 0xb8, 0xcb,
	0x05, 0x0e, 0x12, 0xe0, 0x02, 0x4f, 0xc1, 0x9d, 0x31, 0x96, 0x38, 0xb3, 0x8c, 0x09, 0x66, 0xb6,
	0x64, 0xd2, 0x59, 0x4b, 0x1e, 0x1f, 0xc6, 0x2c, 0x66, 0x2a, 0xe9, 0xca, 0xb7, 0xe2, 0xde, 0xf1,
	0x51, 0xc8, 0x78, 0xca, 0xb8, 0x5f, 0x24, 0x8a, 0x40, 0xa7, 0xec, 0x22, 0x72, 0x03, 0xcc, 0xc1,
	0x9d, 0xf7, 0x03, 0x10, 0xb8, 0xef, 0x86, 0x8c, 0xd0, 0x22, 0x7f, 0xfa, 0xb5, 0x81, 0x6a, 0xef,
	0x19, 0x4b, 0xcc, 0xa7, 0x68, 0x37, 0x82, 0x19, 0xe3, 0x44, 0xf8, 0x11, 0x50, 0x96, 0x5a, 0x46,
	0xc7, 0xe8, 0x6e, 0x7b, 0x3b, 0xfa, 0x70, 0x28, 0xcf, 0xcc, 0x0b, 0xb4, 0x9f, 0x41, 0x04, 0xe9,
	0x4c, 0x10, 0x46, 0xfd, 0x0c, 0x0b, 0xb0, 0xaa, 0xf2, 0xda, 0xa0, 0x7f, 0x75, 0xd3, 0xae, 0xfc,
	0xba, 0x69, 0x9f, 0x14, 0x3a, 0x1e, 0x4d, 0x1d, 0xc2, 0xdc, 0x14, 0x8b, 0x89, 0x33, 0x86, 0x18,
	0x87, 0x8b, 0x21, 0x84, 0x3f, 0xbe, 0xf7, 0x90, 0xee, 0x6d, 0x08, 0xa1, 0xb7, 0x77, 0x4b, 0xf2,
	0xb0, 0x00, 0xf3, 0x1c, 0xed, 0x12, 0x2a, 0x20, 0x03, 0x2e, 0x0a, 0xf2, 0x56, 0x59, 0xf2, 0xce,
	0x3f, 0x8e, 0xe2, 0x7e, 0x40, 0x07, 0x1b, 0x5c, 0x3f, 0xc5, 0xb9, 0x55, 0x2b, 0xcb, 0xde, 0x5f,
	0x67, 0xbf, 0xc3, 0xf9, 0x7f, 0xf0, 0x84, 0x5a, 0x0f, 0xee, 0x07, 0x4f, 0xa8, 0x19, 0xa3, 0xc7,
	0x9b, 0x78, 0x42, 0xc3, 0x0c, 0x30, 0x07, 0xab, 0x5e, 0xd6, 0x71, 0xb8, 0xee, 0x18, 0x69, 0xdc,
	0x5d, 0x51, 0x04, 0x5a, 0xd4, 0xb8, 0x17, 0xd1, 0x50, 0xe3, 0x4c, 0x1f, 0x99, 0x13, 0xc0, 0x89,
	0x98, 0xf8, 0x31, 0x26, 0xd4, 0xff, 0x88, 0x43, 0xc1, 0x32, 0xeb, 0x61, 0x59, 0x49, 0xab, 0x80,
	0xbd, 0xc1, 0x84, 0xbe, 0x56, 0x28, 0x73, 0x8c, 0x9a, 0x82, 0x09, 0x9c, 0xf8, 0x73, 0x9c, 0x5c,
	0x82, 0xb5, 0xad, 0xc8, 0xcf, 0x34, 0xf9, 0xd1, 0x5d, 0xf2, 0x88, 0x8a, 0x35, 0xe6, 0x88, 0x0a,
	0x0f, 0xa9, 0xfa, 0x73, 0x59, 0x2e, 0xdb, 0x4d, 0x71, 0xee, 0x27, 0x30, 0x87, 0x0c, 0xc7, 0x20,
	0xff, 0x0d, 0x61, 0x16, 0x2a, 0xdd, 0x6e, 0x8a, 0xf3, 0xb1, 0x66, 0x79, 0x12, 0x65, 0x3e, 0x41,
	0x0d, 0x39, 0xf2, 0x3e, 0x89, 0xac, 0x66, 0xc7, 0xe8, 0xd6, 0xbc, 0xba, 0x0c, 0x47, 0xd1, 0xe9,
	0x17, 0x03, 0x35, 0x5e, 0xa6, 0xa9, 0x9a, 0xce, 0x3d, 0x54, 0x25, 0x91, 0x1a, 0xc9, 0x9a, 0x57,
	0x25, 0x91, 0x99, 0xa3, 0x83, 0xe2, 0x1b, 0x13, 0x82, 0x03, 0x92, 0x10, 0x41, 0x80, 0x5b, 0xd5,
	0xce, 0x56, 0xb7, 0xf9, 0xfc, 0xc8, 0xd1, 0x3a, 0x39, 0xf2, 0x8e, 0x1e, 0x79, 0xe7, 0x15, 0x23,
	0x74, 0x70, 0x26, 0xfb, 0xfd, 0xf6, 0xbb, 0xdd, 0x8d, 0x89, 0x98, 0x5c, 0x06, 0x4e, 0xc8, 0x52,
	0xbd, 0x2d, 0xf4, 0xa3, 0xc7, 0xa3, 0xa9, 0x2b, 0x16, 0x33, 0xe0, 0xaa, 0x80, 0x7b, 0x2d, 0x65,
	0x19, 0xdf, 0x4a, 0x06, 0x6f, 0xaf, 0x96, 0xb6, 0x71, 0xbd, 0xb4, 0x8d, 0x3f, 0x4b, 0xdb, 0xf8,
	0xbc, 0xb2, 0x2b, 0xd7, 0x2b, 0xbb, 0xf2, 0x73, 0x65, 0x57, 0x2e, 0xce, 0xd6, 0xa8, 0x72, 0x71,
	0xf5, 0x28, 0x88, 0x4f, 0x2c, 0x9b, 0xaa, 0xc0, 0xcd, 0x37, 0x96, 0x9c, 0x72, 0x04, 0x75, 0xb5,
	0x83, 0x5e, 0xfc, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x15, 0x50, 0x4d, 0xfb, 0x05, 0x05, 0x00, 0x00,
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
	if m.PoolId != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x58
	}
	{
		size := m.MaxLeverageRatio.Size()
		i -= size
		if _, err := m.MaxLeverageRatio.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	{
		size := m.TotalValue.Size()
		i -= size
		if _, err := m.TotalValue.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.HealthGainFactor.Size()
		i -= size
		if _, err := m.HealthGainFactor.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.InterestRateDecrease.Size()
		i -= size
		if _, err := m.InterestRateDecrease.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.InterestRateIncrease.Size()
		i -= size
		if _, err := m.InterestRateIncrease.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.InterestRateMin.Size()
		i -= size
		if _, err := m.InterestRateMin.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.InterestRateMax.Size()
		i -= size
		if _, err := m.InterestRateMax.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.InterestRate.Size()
		i -= size
		if _, err := m.InterestRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.RedemptionRate.Size()
		i -= size
		if _, err := m.RedemptionRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.DepositDenom) > 0 {
		i -= len(m.DepositDenom)
		copy(dAtA[i:], m.DepositDenom)
		i = encodeVarintPool(dAtA, i, uint64(len(m.DepositDenom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AmmPool) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AmmPool) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AmmPool) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TotalLiabilities) > 0 {
		for iNdEx := len(m.TotalLiabilities) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.TotalLiabilities[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintPool(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Id != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.Id))
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
	l = len(m.DepositDenom)
	if l > 0 {
		n += 1 + l + sovPool(uint64(l))
	}
	l = m.RedemptionRate.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.InterestRate.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.InterestRateMax.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.InterestRateMin.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.InterestRateIncrease.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.InterestRateDecrease.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.HealthGainFactor.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.TotalValue.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.MaxLeverageRatio.Size()
	n += 1 + l + sovPool(uint64(l))
	if m.PoolId != 0 {
		n += 1 + sovPool(uint64(m.PoolId))
	}
	return n
}

func (m *AmmPool) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovPool(uint64(m.Id))
	}
	if len(m.TotalLiabilities) > 0 {
		for _, e := range m.TotalLiabilities {
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositDenom", wireType)
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
			m.DepositDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RedemptionRate", wireType)
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
			if err := m.RedemptionRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRate", wireType)
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
			if err := m.InterestRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRateMax", wireType)
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
			if err := m.InterestRateMax.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRateMin", wireType)
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
			if err := m.InterestRateMin.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRateIncrease", wireType)
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
			if err := m.InterestRateIncrease.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRateDecrease", wireType)
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
			if err := m.InterestRateDecrease.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HealthGainFactor", wireType)
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
			if err := m.HealthGainFactor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalValue", wireType)
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
			if err := m.TotalValue.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxLeverageRatio", wireType)
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
			if err := m.MaxLeverageRatio.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
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
func (m *AmmPool) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: AmmPool: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AmmPool: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalLiabilities", wireType)
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
			m.TotalLiabilities = append(m.TotalLiabilities, types.Coin{})
			if err := m.TotalLiabilities[len(m.TotalLiabilities)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
