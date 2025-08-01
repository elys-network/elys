// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/leveragelp/params.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

// Params defines the parameters for the module.
type LegacyParams struct {
	LeverageMax         cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=leverage_max,json=leverageMax,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"leverage_max"`
	MaxOpenPositions    int64                       `protobuf:"varint,2,opt,name=max_open_positions,json=maxOpenPositions,proto3" json:"max_open_positions,omitempty"`
	PoolOpenThreshold   cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=pool_open_threshold,json=poolOpenThreshold,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"pool_open_threshold"`
	SafetyFactor        cosmossdk_io_math.LegacyDec `protobuf:"bytes,4,opt,name=safety_factor,json=safetyFactor,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"safety_factor"`
	WhitelistingEnabled bool                        `protobuf:"varint,5,opt,name=whitelisting_enabled,json=whitelistingEnabled,proto3" json:"whitelisting_enabled,omitempty"`
	EpochLength         int64                       `protobuf:"varint,6,opt,name=epoch_length,json=epochLength,proto3" json:"epoch_length,omitempty"`
	FallbackEnabled     bool                        `protobuf:"varint,7,opt,name=fallback_enabled,json=fallbackEnabled,proto3" json:"fallback_enabled,omitempty"`
	NumberPerBlock      int64                       `protobuf:"varint,8,opt,name=number_per_block,json=numberPerBlock,proto3" json:"number_per_block,omitempty"`
}

func (m *LegacyParams) Reset()         { *m = LegacyParams{} }
func (m *LegacyParams) String() string { return proto.CompactTextString(m) }
func (*LegacyParams) ProtoMessage()    {}
func (*LegacyParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_36c27f46b597fbee, []int{0}
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

func (m *LegacyParams) GetMaxOpenPositions() int64 {
	if m != nil {
		return m.MaxOpenPositions
	}
	return 0
}

func (m *LegacyParams) GetWhitelistingEnabled() bool {
	if m != nil {
		return m.WhitelistingEnabled
	}
	return false
}

func (m *LegacyParams) GetEpochLength() int64 {
	if m != nil {
		return m.EpochLength
	}
	return 0
}

func (m *LegacyParams) GetFallbackEnabled() bool {
	if m != nil {
		return m.FallbackEnabled
	}
	return false
}

func (m *LegacyParams) GetNumberPerBlock() int64 {
	if m != nil {
		return m.NumberPerBlock
	}
	return 0
}

// Params defines the parameters for the module.
type Params struct {
	LeverageMax         cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=leverage_max,json=leverageMax,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"leverage_max"`
	MaxOpenPositions    int64                       `protobuf:"varint,2,opt,name=max_open_positions,json=maxOpenPositions,proto3" json:"max_open_positions,omitempty"`
	PoolOpenThreshold   cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=pool_open_threshold,json=poolOpenThreshold,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"pool_open_threshold"`
	SafetyFactor        cosmossdk_io_math.LegacyDec `protobuf:"bytes,4,opt,name=safety_factor,json=safetyFactor,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"safety_factor"`
	WhitelistingEnabled bool                        `protobuf:"varint,5,opt,name=whitelisting_enabled,json=whitelistingEnabled,proto3" json:"whitelisting_enabled,omitempty"`
	EpochLength         int64                       `protobuf:"varint,6,opt,name=epoch_length,json=epochLength,proto3" json:"epoch_length,omitempty"`
	FallbackEnabled     bool                        `protobuf:"varint,7,opt,name=fallback_enabled,json=fallbackEnabled,proto3" json:"fallback_enabled,omitempty"`
	NumberPerBlock      int64                       `protobuf:"varint,8,opt,name=number_per_block,json=numberPerBlock,proto3" json:"number_per_block,omitempty"`
	EnabledPools        []uint64                    `protobuf:"varint,9,rep,packed,name=enabled_pools,json=enabledPools,proto3" json:"enabled_pools,omitempty"`
	ExitBuffer          cosmossdk_io_math.LegacyDec `protobuf:"bytes,10,opt,name=exit_buffer,json=exitBuffer,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"exit_buffer"`
	StopLossEnabled     bool                        `protobuf:"varint,11,opt,name=stop_loss_enabled,json=stopLossEnabled,proto3" json:"stop_loss_enabled,omitempty"`
	LiabilitiesFactor   cosmossdk_io_math.LegacyDec `protobuf:"bytes,12,opt,name=liabilities_factor,json=liabilitiesFactor,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"liabilities_factor"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_36c27f46b597fbee, []int{1}
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

func (m *Params) GetMaxOpenPositions() int64 {
	if m != nil {
		return m.MaxOpenPositions
	}
	return 0
}

func (m *Params) GetWhitelistingEnabled() bool {
	if m != nil {
		return m.WhitelistingEnabled
	}
	return false
}

func (m *Params) GetEpochLength() int64 {
	if m != nil {
		return m.EpochLength
	}
	return 0
}

func (m *Params) GetFallbackEnabled() bool {
	if m != nil {
		return m.FallbackEnabled
	}
	return false
}

func (m *Params) GetNumberPerBlock() int64 {
	if m != nil {
		return m.NumberPerBlock
	}
	return 0
}

func (m *Params) GetEnabledPools() []uint64 {
	if m != nil {
		return m.EnabledPools
	}
	return nil
}

func (m *Params) GetStopLossEnabled() bool {
	if m != nil {
		return m.StopLossEnabled
	}
	return false
}

func init() {
	proto.RegisterType((*LegacyParams)(nil), "elys.leveragelp.LegacyParams")
	proto.RegisterType((*Params)(nil), "elys.leveragelp.Params")
}

func init() { proto.RegisterFile("elys/leveragelp/params.proto", fileDescriptor_36c27f46b597fbee) }

var fileDescriptor_36c27f46b597fbee = []byte{
	// 521 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x94, 0xc1, 0x6e, 0xd3, 0x4c,
	0x14, 0x85, 0xe3, 0x3f, 0x6d, 0xfe, 0x76, 0xe2, 0xd2, 0x74, 0xda, 0x85, 0x29, 0xc8, 0x0d, 0x65,
	0x13, 0x10, 0x8d, 0x55, 0x58, 0xb0, 0x8f, 0x0a, 0xab, 0x20, 0xa2, 0xa8, 0x62, 0xc1, 0x66, 0x18,
	0x3b, 0xd7, 0xf6, 0x28, 0x63, 0xcf, 0xc8, 0x33, 0x69, 0x9d, 0xb7, 0xe0, 0x61, 0x78, 0x88, 0x2e,
	0x2b, 0x56, 0x88, 0x45, 0x85, 0x92, 0x97, 0x80, 0x1d, 0xf2, 0x4c, 0x5c, 0xb2, 0xce, 0x16, 0x76,
	0x9e, 0x73, 0xe6, 0x7e, 0xf6, 0xbd, 0xc7, 0xba, 0xe8, 0x31, 0xf0, 0xb9, 0x0a, 0x38, 0x5c, 0x41,
	0x41, 0x13, 0xe0, 0x32, 0x90, 0xb4, 0xa0, 0x99, 0xea, 0xcb, 0x42, 0x68, 0x81, 0xf7, 0x2b, 0xb7,
	0xff, 0xc7, 0x3d, 0x3e, 0x4a, 0x44, 0x22, 0x8c, 0x17, 0x54, 0x4f, 0xf6, 0xda, 0xf1, 0xc3, 0x48,
	0xa8, 0x4c, 0x28, 0x62, 0x0d, 0x7b, 0xb0, 0xd6, 0xe9, 0xcf, 0x26, 0x72, 0x87, 0x90, 0xd0, 0x68,
	0x3e, 0x32, 0x60, 0x7c, 0x89, 0xdc, 0x9a, 0x47, 0x32, 0x5a, 0x7a, 0x4e, 0xd7, 0xe9, 0xed, 0x0e,
	0xce, 0x6f, 0xee, 0x4e, 0x1a, 0xdf, 0xef, 0x4e, 0x1e, 0xd9, 0x62, 0x35, 0x99, 0xf6, 0x99, 0x08,
	0x32, 0xaa, 0xd3, 0xbe, 0xad, 0xbe, 0x80, 0xe8, 0xeb, 0x97, 0x33, 0xb4, 0x62, 0x5f, 0x40, 0x34,
	0x6e, 0xd7, 0x98, 0x77, 0xb4, 0xc4, 0x2f, 0x10, 0xce, 0x68, 0x49, 0x84, 0x84, 0x9c, 0x48, 0xa1,
	0x98, 0x66, 0x22, 0x57, 0xde, 0x7f, 0x5d, 0xa7, 0xd7, 0x1c, 0x77, 0x32, 0x5a, 0xbe, 0x97, 0x90,
	0x8f, 0x6a, 0x1d, 0x53, 0x74, 0x28, 0x85, 0xe0, 0xf6, 0xba, 0x4e, 0x0b, 0x50, 0xa9, 0xe0, 0x13,
	0xaf, 0xb9, 0xe9, 0xa7, 0x1c, 0x54, 0xb4, 0xea, 0x15, 0x97, 0x35, 0x0b, 0x7f, 0x40, 0x7b, 0x8a,
	0xc6, 0xa0, 0xe7, 0x24, 0xa6, 0x91, 0x16, 0x85, 0xb7, 0xb5, 0x29, 0xdc, 0xb5, 0x9c, 0xb7, 0x06,
	0x83, 0xcf, 0xd1, 0xd1, 0x75, 0xca, 0x34, 0x70, 0xa6, 0x34, 0xcb, 0x13, 0x02, 0x39, 0x0d, 0x39,
	0x4c, 0xbc, 0xed, 0xae, 0xd3, 0xdb, 0x19, 0x1f, 0xae, 0x7b, 0x6f, 0xac, 0x85, 0x9f, 0x20, 0x17,
	0xa4, 0x88, 0x52, 0xc2, 0x21, 0x4f, 0x74, 0xea, 0xb5, 0xcc, 0x54, 0xda, 0x46, 0x1b, 0x1a, 0x09,
	0x3f, 0x43, 0x9d, 0x98, 0x72, 0x1e, 0xd2, 0x68, 0x7a, 0x4f, 0xfc, 0xdf, 0x10, 0xf7, 0x6b, 0xbd,
	0xa6, 0xf5, 0x50, 0x27, 0x9f, 0x65, 0x21, 0x14, 0x44, 0x42, 0x41, 0x42, 0x2e, 0xa2, 0xa9, 0xb7,
	0x63, 0x88, 0x0f, 0xac, 0x3e, 0x82, 0x62, 0x50, 0xa9, 0xa7, 0xbf, 0xb6, 0x51, 0xeb, 0x5f, 0xe8,
	0x7f, 0x57, 0xe8, 0xf8, 0x29, 0xda, 0x5b, 0xb1, 0x48, 0x35, 0x1f, 0xe5, 0xed, 0x76, 0x9b, 0xbd,
	0xad, 0xb1, 0xbb, 0x12, 0x47, 0x95, 0x86, 0xc7, 0xa8, 0x0d, 0x25, 0xd3, 0x24, 0x9c, 0xc5, 0x31,
	0x14, 0x1e, 0xda, 0x74, 0x4a, 0xa8, 0xa2, 0x0c, 0x0c, 0x04, 0x3f, 0x47, 0x07, 0x4a, 0x0b, 0x49,
	0xb8, 0x50, 0xea, 0xbe, 0x9d, 0xb6, 0x6d, 0xa7, 0x32, 0x86, 0x42, 0xa9, 0xba, 0x9d, 0x4f, 0x08,
	0x73, 0x46, 0x43, 0xc6, 0x99, 0x66, 0xa0, 0xea, 0xb0, 0xdc, 0x8d, 0xff, 0x84, 0x35, 0x98, 0x4d,
	0x6c, 0x30, 0xbc, 0x59, 0xf8, 0xce, 0xed, 0xc2, 0x77, 0x7e, 0x2c, 0x7c, 0xe7, 0xf3, 0xd2, 0x6f,
	0xdc, 0x2e, 0xfd, 0xc6, 0xb7, 0xa5, 0xdf, 0xf8, 0xf8, 0x32, 0x61, 0x3a, 0x9d, 0x85, 0xfd, 0x48,
	0x64, 0x41, 0xb5, 0x5d, 0xcf, 0x72, 0xd0, 0xd7, 0xa2, 0x98, 0x9a, 0x43, 0x70, 0xf5, 0x3a, 0x28,
	0xd7, 0xb7, 0xb1, 0x9e, 0x4b, 0x50, 0x61, 0xcb, 0xec, 0xd2, 0x57, 0xbf, 0x03, 0x00, 0x00, 0xff,
	0xff, 0xa7, 0xfb, 0x8c, 0xee, 0xad, 0x05, 0x00, 0x00,
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
	if m.NumberPerBlock != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.NumberPerBlock))
		i--
		dAtA[i] = 0x40
	}
	if m.FallbackEnabled {
		i--
		if m.FallbackEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x38
	}
	if m.EpochLength != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.EpochLength))
		i--
		dAtA[i] = 0x30
	}
	if m.WhitelistingEnabled {
		i--
		if m.WhitelistingEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.SafetyFactor.Size()
		i -= size
		if _, err := m.SafetyFactor.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.PoolOpenThreshold.Size()
		i -= size
		if _, err := m.PoolOpenThreshold.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.MaxOpenPositions != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxOpenPositions))
		i--
		dAtA[i] = 0x10
	}
	{
		size := m.LeverageMax.Size()
		i -= size
		if _, err := m.LeverageMax.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
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
	{
		size := m.LiabilitiesFactor.Size()
		i -= size
		if _, err := m.LiabilitiesFactor.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x62
	if m.StopLossEnabled {
		i--
		if m.StopLossEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x58
	}
	{
		size := m.ExitBuffer.Size()
		i -= size
		if _, err := m.ExitBuffer.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	if len(m.EnabledPools) > 0 {
		dAtA2 := make([]byte, len(m.EnabledPools)*10)
		var j1 int
		for _, num := range m.EnabledPools {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintParams(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0x4a
	}
	if m.NumberPerBlock != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.NumberPerBlock))
		i--
		dAtA[i] = 0x40
	}
	if m.FallbackEnabled {
		i--
		if m.FallbackEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x38
	}
	if m.EpochLength != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.EpochLength))
		i--
		dAtA[i] = 0x30
	}
	if m.WhitelistingEnabled {
		i--
		if m.WhitelistingEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.SafetyFactor.Size()
		i -= size
		if _, err := m.SafetyFactor.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.PoolOpenThreshold.Size()
		i -= size
		if _, err := m.PoolOpenThreshold.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.MaxOpenPositions != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxOpenPositions))
		i--
		dAtA[i] = 0x10
	}
	{
		size := m.LeverageMax.Size()
		i -= size
		if _, err := m.LeverageMax.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
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
func (m *LegacyParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.LeverageMax.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.MaxOpenPositions != 0 {
		n += 1 + sovParams(uint64(m.MaxOpenPositions))
	}
	l = m.PoolOpenThreshold.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.SafetyFactor.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.WhitelistingEnabled {
		n += 2
	}
	if m.EpochLength != 0 {
		n += 1 + sovParams(uint64(m.EpochLength))
	}
	if m.FallbackEnabled {
		n += 2
	}
	if m.NumberPerBlock != 0 {
		n += 1 + sovParams(uint64(m.NumberPerBlock))
	}
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.LeverageMax.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.MaxOpenPositions != 0 {
		n += 1 + sovParams(uint64(m.MaxOpenPositions))
	}
	l = m.PoolOpenThreshold.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.SafetyFactor.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.WhitelistingEnabled {
		n += 2
	}
	if m.EpochLength != 0 {
		n += 1 + sovParams(uint64(m.EpochLength))
	}
	if m.FallbackEnabled {
		n += 2
	}
	if m.NumberPerBlock != 0 {
		n += 1 + sovParams(uint64(m.NumberPerBlock))
	}
	if len(m.EnabledPools) > 0 {
		l = 0
		for _, e := range m.EnabledPools {
			l += sovParams(uint64(e))
		}
		n += 1 + sovParams(uint64(l)) + l
	}
	l = m.ExitBuffer.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.StopLossEnabled {
		n += 2
	}
	l = m.LiabilitiesFactor.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
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
				return fmt.Errorf("proto: wrong wireType = %d for field LeverageMax", wireType)
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
			if err := m.LeverageMax.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxOpenPositions", wireType)
			}
			m.MaxOpenPositions = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxOpenPositions |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolOpenThreshold", wireType)
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
			if err := m.PoolOpenThreshold.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SafetyFactor", wireType)
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
			if err := m.SafetyFactor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field WhitelistingEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
			m.WhitelistingEnabled = bool(v != 0)
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EpochLength", wireType)
			}
			m.EpochLength = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EpochLength |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FallbackEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
			m.FallbackEnabled = bool(v != 0)
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumberPerBlock", wireType)
			}
			m.NumberPerBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumberPerBlock |= int64(b&0x7F) << shift
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
				return fmt.Errorf("proto: wrong wireType = %d for field LeverageMax", wireType)
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
			if err := m.LeverageMax.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxOpenPositions", wireType)
			}
			m.MaxOpenPositions = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxOpenPositions |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolOpenThreshold", wireType)
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
			if err := m.PoolOpenThreshold.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SafetyFactor", wireType)
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
			if err := m.SafetyFactor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field WhitelistingEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
			m.WhitelistingEnabled = bool(v != 0)
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EpochLength", wireType)
			}
			m.EpochLength = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EpochLength |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FallbackEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
			m.FallbackEnabled = bool(v != 0)
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumberPerBlock", wireType)
			}
			m.NumberPerBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumberPerBlock |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowParams
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
				m.EnabledPools = append(m.EnabledPools, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowParams
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
					return ErrInvalidLengthParams
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthParams
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
				if elementCount != 0 && len(m.EnabledPools) == 0 {
					m.EnabledPools = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowParams
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
					m.EnabledPools = append(m.EnabledPools, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field EnabledPools", wireType)
			}
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExitBuffer", wireType)
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
			if err := m.ExitBuffer.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StopLossEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
			m.StopLossEnabled = bool(v != 0)
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LiabilitiesFactor", wireType)
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
			if err := m.LiabilitiesFactor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
