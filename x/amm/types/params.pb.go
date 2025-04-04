// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/amm/params.proto

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

type Params struct {
	PoolCreationFee                  cosmossdk_io_math.Int       `protobuf:"bytes,1,opt,name=pool_creation_fee,json=poolCreationFee,proto3,customtype=cosmossdk.io/math.Int" json:"pool_creation_fee"`
	SlippageTrackDuration            uint64                      `protobuf:"varint,2,opt,name=slippage_track_duration,json=slippageTrackDuration,proto3" json:"slippage_track_duration,omitempty"`
	BaseAssets                       []string                    `protobuf:"bytes,3,rep,name=base_assets,json=baseAssets,proto3" json:"base_assets,omitempty"`
	WeightBreakingFeeExponent        cosmossdk_io_math.LegacyDec `protobuf:"bytes,4,opt,name=weight_breaking_fee_exponent,json=weightBreakingFeeExponent,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"weight_breaking_fee_exponent"`
	WeightBreakingFeeMultiplier      cosmossdk_io_math.LegacyDec `protobuf:"bytes,5,opt,name=weight_breaking_fee_multiplier,json=weightBreakingFeeMultiplier,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"weight_breaking_fee_multiplier"`
	WeightBreakingFeePortion         cosmossdk_io_math.LegacyDec `protobuf:"bytes,6,opt,name=weight_breaking_fee_portion,json=weightBreakingFeePortion,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"weight_breaking_fee_portion"`
	WeightRecoveryFeePortion         cosmossdk_io_math.LegacyDec `protobuf:"bytes,7,opt,name=weight_recovery_fee_portion,json=weightRecoveryFeePortion,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"weight_recovery_fee_portion"`
	ThresholdWeightDifference        cosmossdk_io_math.LegacyDec `protobuf:"bytes,8,opt,name=threshold_weight_difference,json=thresholdWeightDifference,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"threshold_weight_difference"`
	AllowedPoolCreators              []string                    `protobuf:"bytes,9,rep,name=allowed_pool_creators,json=allowedPoolCreators,proto3" json:"allowed_pool_creators,omitempty"`
	ThresholdWeightDifferenceSwapFee cosmossdk_io_math.LegacyDec `protobuf:"bytes,10,opt,name=threshold_weight_difference_swap_fee,json=thresholdWeightDifferenceSwapFee,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"threshold_weight_difference_swap_fee"`
	LpLockupDuration                 uint64                      `protobuf:"varint,11,opt,name=lp_lockup_duration,json=lpLockupDuration,proto3" json:"lp_lockup_duration,omitempty"`
	MinSlippage                      cosmossdk_io_math.LegacyDec `protobuf:"bytes,12,opt,name=min_slippage,json=minSlippage,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"min_slippage"`
	AllowedUpfrontSwapMakers         []string                    `protobuf:"bytes,13,rep,name=allowed_upfront_swap_makers,json=allowedUpfrontSwapMakers,proto3" json:"allowed_upfront_swap_makers,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_1209ca218537a425, []int{0}
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

func (m *Params) GetSlippageTrackDuration() uint64 {
	if m != nil {
		return m.SlippageTrackDuration
	}
	return 0
}

func (m *Params) GetBaseAssets() []string {
	if m != nil {
		return m.BaseAssets
	}
	return nil
}

func (m *Params) GetAllowedPoolCreators() []string {
	if m != nil {
		return m.AllowedPoolCreators
	}
	return nil
}

func (m *Params) GetLpLockupDuration() uint64 {
	if m != nil {
		return m.LpLockupDuration
	}
	return 0
}

func (m *Params) GetAllowedUpfrontSwapMakers() []string {
	if m != nil {
		return m.AllowedUpfrontSwapMakers
	}
	return nil
}

func init() {
	proto.RegisterType((*Params)(nil), "elys.amm.Params")
}

func init() { proto.RegisterFile("elys/amm/params.proto", fileDescriptor_1209ca218537a425) }

var fileDescriptor_1209ca218537a425 = []byte{
	// 579 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0x4f, 0x6f, 0x12, 0x4f,
	0x18, 0xc7, 0xd9, 0x5f, 0xfb, 0xab, 0x30, 0xd4, 0xa8, 0xab, 0xc4, 0x6d, 0x31, 0x0b, 0x31, 0x1e,
	0x48, 0xb4, 0x6c, 0xd4, 0xc4, 0x9b, 0x07, 0x11, 0x9b, 0x34, 0x69, 0x13, 0x42, 0x6b, 0x9a, 0x78,
	0x99, 0x0c, 0xcb, 0xc3, 0xb2, 0xd9, 0x9d, 0x3f, 0xce, 0x0c, 0xa5, 0x1c, 0x7d, 0x07, 0xbe, 0x18,
	0x5f, 0x44, 0x8f, 0x8d, 0x27, 0xe3, 0xa1, 0x31, 0xf0, 0x26, 0x3c, 0x9a, 0x9d, 0x1d, 0x28, 0x86,
	0xea, 0x01, 0x6f, 0xcc, 0xf3, 0x1d, 0x3e, 0xdf, 0xef, 0xcc, 0xce, 0xf3, 0xa0, 0x0a, 0xa4, 0x13,
	0x15, 0x10, 0x4a, 0x03, 0x41, 0x24, 0xa1, 0xaa, 0x29, 0x24, 0xd7, 0xdc, 0x2d, 0x66, 0xe5, 0x26,
	0xa1, 0x74, 0xf7, 0x41, 0xc4, 0x23, 0x6e, 0x8a, 0x41, 0xf6, 0x2b, 0xd7, 0x77, 0x77, 0x42, 0xae,
	0x28, 0x57, 0x38, 0x17, 0xf2, 0x45, 0x2e, 0x3d, 0xfe, 0x59, 0x44, 0x5b, 0x1d, 0xc3, 0x72, 0x4f,
	0xd1, 0x3d, 0xc1, 0x79, 0x8a, 0x43, 0x09, 0x44, 0xc7, 0x9c, 0xe1, 0x01, 0x80, 0xe7, 0xd4, 0x9d,
	0x46, 0xa9, 0xf5, 0xf4, 0xe2, 0xaa, 0x56, 0xf8, 0x7e, 0x55, 0xab, 0xe4, 0xff, 0x55, 0xfd, 0xa4,
	0x19, 0xf3, 0x80, 0x12, 0x3d, 0x6c, 0x1e, 0x30, 0xfd, 0xf5, 0xcb, 0x1e, 0xb2, 0xd0, 0x03, 0xa6,
	0xbb, 0x77, 0x32, 0xca, 0x5b, 0x0b, 0xd9, 0x07, 0x70, 0x5f, 0xa1, 0x87, 0x2a, 0x8d, 0x85, 0x20,
	0x11, 0x60, 0x2d, 0x49, 0x98, 0xe0, 0xfe, 0x48, 0x1a, 0xd5, 0xfb, 0xaf, 0xee, 0x34, 0x36, 0xbb,
	0x95, 0xb9, 0x7c, 0x92, 0xa9, 0x6d, 0x2b, 0xba, 0x35, 0x54, 0xee, 0x11, 0x05, 0x98, 0x28, 0x05,
	0x5a, 0x79, 0x1b, 0xf5, 0x8d, 0x46, 0xa9, 0x8b, 0xb2, 0xd2, 0x1b, 0x53, 0x71, 0x25, 0x7a, 0x34,
	0x86, 0x38, 0x1a, 0x6a, 0xdc, 0x93, 0x40, 0x92, 0x98, 0x45, 0x59, 0x66, 0x0c, 0xe7, 0x82, 0x33,
	0x60, 0xda, 0xdb, 0x34, 0xe1, 0x9f, 0xdb, 0xf0, 0xd5, 0xd5, 0xf0, 0x87, 0x10, 0x91, 0x70, 0xd2,
	0x86, 0x70, 0xe9, 0x08, 0x6d, 0x08, 0xbb, 0x3b, 0x39, 0xb6, 0x65, 0xa9, 0xfb, 0x00, 0xef, 0x2c,
	0xd3, 0x3d, 0x43, 0xfe, 0x4d, 0x9e, 0x74, 0x94, 0xea, 0x58, 0xa4, 0x31, 0x48, 0xef, 0xff, 0x75,
	0x5d, 0xab, 0x2b, 0xae, 0x47, 0x0b, 0xaa, 0x2b, 0x50, 0xf5, 0x26, 0x5f, 0xc1, 0xa5, 0xb9, 0xc8,
	0xad, 0x75, 0x4d, 0xbd, 0x15, 0xd3, 0x4e, 0x8e, 0x5c, 0x72, 0x94, 0x10, 0xf2, 0x33, 0x90, 0x93,
	0xdf, 0x1c, 0x6f, 0xfd, 0xa3, 0x63, 0xd7, 0x42, 0x97, 0x1c, 0x3f, 0xa2, 0xaa, 0x1e, 0x4a, 0x50,
	0x43, 0x9e, 0xf6, 0xb1, 0xf5, 0xee, 0xc7, 0x83, 0x01, 0x48, 0x60, 0x21, 0x78, 0xc5, 0xb5, 0x3f,
	0xe7, 0x82, 0x7a, 0x6a, 0xa0, 0xed, 0x05, 0xd3, 0x7d, 0x81, 0x2a, 0x24, 0x4d, 0xf9, 0x18, 0xfa,
	0xf8, 0xfa, 0xf1, 0x73, 0xa9, 0xbc, 0x92, 0x79, 0x6d, 0xf7, 0xad, 0xd8, 0x99, 0x3f, 0x69, 0x2e,
	0x95, 0xfb, 0xc9, 0x41, 0x4f, 0xfe, 0x92, 0x13, 0xab, 0x31, 0x11, 0xa6, 0x79, 0xd0, 0xba, 0x81,
	0xeb, 0x7f, 0x0c, 0x7c, 0x3c, 0x26, 0x22, 0xeb, 0xa9, 0x67, 0xc8, 0x4d, 0x05, 0x4e, 0x79, 0x98,
	0x8c, 0xc4, 0x75, 0x3b, 0x95, 0x4d, 0x3b, 0xdd, 0x4d, 0xc5, 0xa1, 0x11, 0x16, 0x9d, 0x74, 0x82,
	0xb6, 0x69, 0xcc, 0xf0, 0xbc, 0xcd, 0xbc, 0xed, 0x75, 0x83, 0x95, 0x69, 0xcc, 0x8e, 0x2d, 0xc5,
	0x7d, 0x8d, 0xaa, 0xf3, 0xbb, 0x1b, 0x89, 0x81, 0xe4, 0x4c, 0xe7, 0x47, 0xa7, 0x24, 0x01, 0xa9,
	0xbc, 0xdb, 0xe6, 0x06, 0x3d, 0xbb, 0xe5, 0x7d, 0xbe, 0x23, 0xcb, 0x7f, 0x64, 0xf4, 0x56, 0xeb,
	0x62, 0xea, 0x3b, 0x97, 0x53, 0xdf, 0xf9, 0x31, 0xf5, 0x9d, 0xcf, 0x33, 0xbf, 0x70, 0x39, 0xf3,
	0x0b, 0xdf, 0x66, 0x7e, 0xe1, 0x43, 0x23, 0x8a, 0xf5, 0x70, 0xd4, 0x6b, 0x86, 0x9c, 0x06, 0xd9,
	0x68, 0xdb, 0x63, 0xa0, 0xc7, 0x5c, 0x26, 0x66, 0x11, 0x9c, 0x9b, 0x01, 0xa8, 0x27, 0x02, 0x54,
	0x6f, 0xcb, 0x4c, 0xb1, 0x97, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0xd2, 0xde, 0xc1, 0x53, 0x19,
	0x05, 0x00, 0x00,
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
	if len(m.AllowedUpfrontSwapMakers) > 0 {
		for iNdEx := len(m.AllowedUpfrontSwapMakers) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.AllowedUpfrontSwapMakers[iNdEx])
			copy(dAtA[i:], m.AllowedUpfrontSwapMakers[iNdEx])
			i = encodeVarintParams(dAtA, i, uint64(len(m.AllowedUpfrontSwapMakers[iNdEx])))
			i--
			dAtA[i] = 0x6a
		}
	}
	{
		size := m.MinSlippage.Size()
		i -= size
		if _, err := m.MinSlippage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x62
	if m.LpLockupDuration != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.LpLockupDuration))
		i--
		dAtA[i] = 0x58
	}
	{
		size := m.ThresholdWeightDifferenceSwapFee.Size()
		i -= size
		if _, err := m.ThresholdWeightDifferenceSwapFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	if len(m.AllowedPoolCreators) > 0 {
		for iNdEx := len(m.AllowedPoolCreators) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.AllowedPoolCreators[iNdEx])
			copy(dAtA[i:], m.AllowedPoolCreators[iNdEx])
			i = encodeVarintParams(dAtA, i, uint64(len(m.AllowedPoolCreators[iNdEx])))
			i--
			dAtA[i] = 0x4a
		}
	}
	{
		size := m.ThresholdWeightDifference.Size()
		i -= size
		if _, err := m.ThresholdWeightDifference.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.WeightRecoveryFeePortion.Size()
		i -= size
		if _, err := m.WeightRecoveryFeePortion.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.WeightBreakingFeePortion.Size()
		i -= size
		if _, err := m.WeightBreakingFeePortion.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.WeightBreakingFeeMultiplier.Size()
		i -= size
		if _, err := m.WeightBreakingFeeMultiplier.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.WeightBreakingFeeExponent.Size()
		i -= size
		if _, err := m.WeightBreakingFeeExponent.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.BaseAssets) > 0 {
		for iNdEx := len(m.BaseAssets) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.BaseAssets[iNdEx])
			copy(dAtA[i:], m.BaseAssets[iNdEx])
			i = encodeVarintParams(dAtA, i, uint64(len(m.BaseAssets[iNdEx])))
			i--
			dAtA[i] = 0x1a
		}
	}
	if m.SlippageTrackDuration != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SlippageTrackDuration))
		i--
		dAtA[i] = 0x10
	}
	{
		size := m.PoolCreationFee.Size()
		i -= size
		if _, err := m.PoolCreationFee.MarshalTo(dAtA[i:]); err != nil {
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
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.PoolCreationFee.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.SlippageTrackDuration != 0 {
		n += 1 + sovParams(uint64(m.SlippageTrackDuration))
	}
	if len(m.BaseAssets) > 0 {
		for _, s := range m.BaseAssets {
			l = len(s)
			n += 1 + l + sovParams(uint64(l))
		}
	}
	l = m.WeightBreakingFeeExponent.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.WeightBreakingFeeMultiplier.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.WeightBreakingFeePortion.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.WeightRecoveryFeePortion.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.ThresholdWeightDifference.Size()
	n += 1 + l + sovParams(uint64(l))
	if len(m.AllowedPoolCreators) > 0 {
		for _, s := range m.AllowedPoolCreators {
			l = len(s)
			n += 1 + l + sovParams(uint64(l))
		}
	}
	l = m.ThresholdWeightDifferenceSwapFee.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.LpLockupDuration != 0 {
		n += 1 + sovParams(uint64(m.LpLockupDuration))
	}
	l = m.MinSlippage.Size()
	n += 1 + l + sovParams(uint64(l))
	if len(m.AllowedUpfrontSwapMakers) > 0 {
		for _, s := range m.AllowedUpfrontSwapMakers {
			l = len(s)
			n += 1 + l + sovParams(uint64(l))
		}
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
				return fmt.Errorf("proto: wrong wireType = %d for field PoolCreationFee", wireType)
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
			if err := m.PoolCreationFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlippageTrackDuration", wireType)
			}
			m.SlippageTrackDuration = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SlippageTrackDuration |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseAssets", wireType)
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
			m.BaseAssets = append(m.BaseAssets, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WeightBreakingFeeExponent", wireType)
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
			if err := m.WeightBreakingFeeExponent.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WeightBreakingFeeMultiplier", wireType)
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
			if err := m.WeightBreakingFeeMultiplier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WeightBreakingFeePortion", wireType)
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
			if err := m.WeightBreakingFeePortion.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WeightRecoveryFeePortion", wireType)
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
			if err := m.WeightRecoveryFeePortion.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ThresholdWeightDifference", wireType)
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
			if err := m.ThresholdWeightDifference.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AllowedPoolCreators", wireType)
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
			m.AllowedPoolCreators = append(m.AllowedPoolCreators, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ThresholdWeightDifferenceSwapFee", wireType)
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
			if err := m.ThresholdWeightDifferenceSwapFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LpLockupDuration", wireType)
			}
			m.LpLockupDuration = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LpLockupDuration |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinSlippage", wireType)
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
			if err := m.MinSlippage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AllowedUpfrontSwapMakers", wireType)
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
			m.AllowedUpfrontSwapMakers = append(m.AllowedUpfrontSwapMakers, string(dAtA[iNdEx:postIndex]))
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
