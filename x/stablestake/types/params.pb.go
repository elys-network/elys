// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/stablestake/params.proto

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

// Params defines the parameters for the module.
type Params struct {
	DepositDenom                             string                                 `protobuf:"bytes,1,opt,name=deposit_denom,json=depositDenom,proto3" json:"deposit_denom,omitempty"`
	RedemptionRate                           github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=redemption_rate,json=redemptionRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"redemption_rate"`
	EpochLength                              int64                                  `protobuf:"varint,3,opt,name=epoch_length,json=epochLength,proto3" json:"epoch_length,omitempty"`
	InterestRate                             github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=interest_rate,json=interestRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate"`
	InterestRateMax                          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=interest_rate_max,json=interestRateMax,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate_max"`
	InterestRateMin                          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=interest_rate_min,json=interestRateMin,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate_min"`
	InterestRateIncrease                     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=interest_rate_increase,json=interestRateIncrease,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate_increase"`
	InterestRateDecrease                     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,8,opt,name=interest_rate_decrease,json=interestRateDecrease,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate_decrease"`
	IncrementalInterestPaymentFundPercentage github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,9,opt,name=incremental_interest_payment_fund_percentage,json=incrementalInterestPaymentFundPercentage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"incremental_interest_payment_fund_percentage"`
	IncrementalInterestPaymentFundAddress    string                                 `protobuf:"bytes,10,opt,name=incremental_interest_payment_fund_address,json=incrementalInterestPaymentFundAddress,proto3" json:"incremental_interest_payment_fund_address,omitempty"`
	IncrementalInterestPaymentEnabled        bool                                   `protobuf:"varint,11,opt,name=incremental_interest_payment_enabled,json=incrementalInterestPaymentEnabled,proto3" json:"incremental_interest_payment_enabled,omitempty"`
	HealthGainFactor                         github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,12,opt,name=health_gain_factor,json=healthGainFactor,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"health_gain_factor"`
	TotalValue                               github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,13,opt,name=total_value,json=totalValue,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"total_value"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_d2a5512d8c510b7b, []int{0}
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

func (m *Params) GetDepositDenom() string {
	if m != nil {
		return m.DepositDenom
	}
	return ""
}

func (m *Params) GetEpochLength() int64 {
	if m != nil {
		return m.EpochLength
	}
	return 0
}

func (m *Params) GetIncrementalInterestPaymentFundAddress() string {
	if m != nil {
		return m.IncrementalInterestPaymentFundAddress
	}
	return ""
}

func (m *Params) GetIncrementalInterestPaymentEnabled() bool {
	if m != nil {
		return m.IncrementalInterestPaymentEnabled
	}
	return false
}

func init() {
	proto.RegisterType((*Params)(nil), "elys.stablestake.Params")
}

func init() { proto.RegisterFile("elys/stablestake/params.proto", fileDescriptor_d2a5512d8c510b7b) }

var fileDescriptor_d2a5512d8c510b7b = []byte{
	// 513 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xdf, 0x6a, 0x13, 0x4f,
	0x14, 0xc7, 0xb3, 0xbf, 0xf6, 0x17, 0xd3, 0x49, 0x62, 0xeb, 0x50, 0x64, 0x10, 0xdc, 0xa4, 0xfe,
	0x23, 0x82, 0xcd, 0x0a, 0xde, 0x79, 0x67, 0x89, 0x95, 0x88, 0xd2, 0xb0, 0x82, 0x4a, 0x11, 0x96,
	0xc9, 0xce, 0xe9, 0x66, 0xc9, 0xee, 0xcc, 0x32, 0x33, 0xd1, 0xe4, 0x25, 0xc4, 0x4b, 0x2f, 0x7d,
	0x9c, 0x5e, 0xf6, 0x52, 0xbc, 0x28, 0x92, 0x3c, 0x82, 0x2f, 0x20, 0x33, 0xbb, 0xd6, 0x04, 0xb5,
	0xe2, 0xea, 0xd5, 0xee, 0x7e, 0xcf, 0xd9, 0xcf, 0xe7, 0x30, 0x30, 0x07, 0x5d, 0x85, 0x64, 0xa6,
	0x3c, 0xa5, 0xe9, 0x30, 0x01, 0xa5, 0xe9, 0x18, 0xbc, 0x8c, 0x4a, 0x9a, 0xaa, 0x6e, 0x26, 0x85,
	0x16, 0x78, 0xcb, 0x94, 0xbb, 0x4b, 0xe5, 0x2b, 0xdb, 0x91, 0x88, 0x84, 0x2d, 0x7a, 0xe6, 0x2d,
	0xef, 0xbb, 0xf6, 0xa5, 0x86, 0xaa, 0x03, 0xfb, 0x23, 0xbe, 0x8e, 0x9a, 0x0c, 0x32, 0xa1, 0x62,
	0x1d, 0x30, 0xe0, 0x22, 0x25, 0x4e, 0xdb, 0xe9, 0x6c, 0xf8, 0x8d, 0x22, 0xec, 0x99, 0x0c, 0xbf,
	0x40, 0x9b, 0x12, 0x18, 0xa4, 0x99, 0x8e, 0x05, 0x0f, 0x24, 0xd5, 0x40, 0xfe, 0x33, 0x6d, 0x7b,
	0xdd, 0xe3, 0xd3, 0x56, 0xe5, 0xd3, 0x69, 0xeb, 0x56, 0x14, 0xeb, 0xd1, 0x64, 0xd8, 0x0d, 0x45,
	0xea, 0x85, 0x42, 0xa5, 0x42, 0x15, 0x8f, 0x5d, 0xc5, 0xc6, 0x9e, 0x9e, 0x65, 0xa0, 0xba, 0x3d,
	0x08, 0xfd, 0x8b, 0xdf, 0x31, 0x3e, 0xd5, 0x80, 0x77, 0x50, 0x03, 0x32, 0x11, 0x8e, 0x82, 0x04,
	0x78, 0xa4, 0x47, 0x64, 0xad, 0xed, 0x74, 0xd6, 0xfc, 0xba, 0xcd, 0x9e, 0xd8, 0x08, 0x3f, 0x43,
	0xcd, 0x98, 0x6b, 0x90, 0xa0, 0x74, 0x6e, 0x5e, 0x2f, 0x65, 0x6e, 0x7c, 0x83, 0x58, 0xef, 0x21,
	0xba, 0xb4, 0x02, 0x0d, 0x52, 0x3a, 0x25, 0xff, 0x97, 0x02, 0x6f, 0x2e, 0x83, 0x9f, 0xd2, 0xe9,
	0x4f, 0xd8, 0x31, 0x27, 0xd5, 0x7f, 0xc0, 0x8e, 0x39, 0x66, 0xe8, 0xf2, 0x2a, 0x3b, 0xe6, 0xa1,
	0x04, 0xaa, 0x80, 0x5c, 0x28, 0x25, 0xd8, 0x5e, 0x16, 0xf4, 0x0b, 0xd6, 0x8f, 0x16, 0x06, 0x85,
	0xa5, 0xf6, 0xf7, 0x96, 0x5e, 0xc1, 0xc2, 0x6f, 0x1d, 0x74, 0xc7, 0x8e, 0x9f, 0x02, 0xd7, 0x34,
	0x09, 0xce, 0x94, 0x19, 0x9d, 0x99, 0x28, 0x38, 0x9a, 0x70, 0x16, 0x64, 0x20, 0x43, 0x53, 0x8f,
	0x80, 0x6c, 0x94, 0x92, 0x77, 0x96, 0x1c, 0xfd, 0x42, 0x31, 0xc8, 0x0d, 0xfb, 0x13, 0xce, 0x06,
	0x67, 0x7c, 0xfc, 0x12, 0xdd, 0xfe, 0xfd, 0x3c, 0x94, 0x31, 0x09, 0x4a, 0x11, 0x64, 0xaf, 0xc9,
	0xcd, 0xf3, 0xe1, 0x0f, 0xf2, 0x66, 0x7c, 0x80, 0x6e, 0x9c, 0x4b, 0x06, 0x6e, 0x6e, 0x2c, 0x23,
	0xf5, 0xb6, 0xd3, 0xa9, 0xf9, 0x3b, 0xbf, 0x86, 0x3e, 0xcc, 0x1b, 0xf1, 0x2b, 0x84, 0x47, 0x40,
	0x13, 0x3d, 0x0a, 0x22, 0x1a, 0xf3, 0xe0, 0x88, 0x86, 0x5a, 0x48, 0xd2, 0x28, 0x75, 0x40, 0x5b,
	0x39, 0xe9, 0x11, 0x8d, 0xf9, 0xbe, 0xe5, 0xe0, 0x03, 0x54, 0xd7, 0xc2, 0x0c, 0xfa, 0x9a, 0x26,
	0x13, 0x20, 0xcd, 0x3f, 0xc6, 0xf6, 0xb9, 0xf6, 0x91, 0x45, 0x3c, 0x37, 0x84, 0xfb, 0xeb, 0xef,
	0x3f, 0xb4, 0x2a, 0x7b, 0x8f, 0x8f, 0xe7, 0xae, 0x73, 0x32, 0x77, 0x9d, 0xcf, 0x73, 0xd7, 0x79,
	0xb7, 0x70, 0x2b, 0x27, 0x0b, 0xb7, 0xf2, 0x71, 0xe1, 0x56, 0x0e, 0xef, 0x2e, 0x31, 0xcd, 0x0a,
	0xdb, 0xe5, 0xa0, 0xdf, 0x08, 0x39, 0xb6, 0x1f, 0xde, 0x74, 0x65, 0xe1, 0x59, 0xc3, 0xb0, 0x6a,
	0x17, 0xd9, 0xbd, 0xaf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x1b, 0x8d, 0xfb, 0x9c, 0x11, 0x05, 0x00,
	0x00,
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
		size := m.TotalValue.Size()
		i -= size
		if _, err := m.TotalValue.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x6a
	{
		size := m.HealthGainFactor.Size()
		i -= size
		if _, err := m.HealthGainFactor.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x62
	if m.IncrementalInterestPaymentEnabled {
		i--
		if m.IncrementalInterestPaymentEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x58
	}
	if len(m.IncrementalInterestPaymentFundAddress) > 0 {
		i -= len(m.IncrementalInterestPaymentFundAddress)
		copy(dAtA[i:], m.IncrementalInterestPaymentFundAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.IncrementalInterestPaymentFundAddress)))
		i--
		dAtA[i] = 0x52
	}
	{
		size := m.IncrementalInterestPaymentFundPercentage.Size()
		i -= size
		if _, err := m.IncrementalInterestPaymentFundPercentage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.InterestRateDecrease.Size()
		i -= size
		if _, err := m.InterestRateDecrease.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.InterestRateIncrease.Size()
		i -= size
		if _, err := m.InterestRateIncrease.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.InterestRateMin.Size()
		i -= size
		if _, err := m.InterestRateMin.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.InterestRateMax.Size()
		i -= size
		if _, err := m.InterestRateMax.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.InterestRate.Size()
		i -= size
		if _, err := m.InterestRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if m.EpochLength != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.EpochLength))
		i--
		dAtA[i] = 0x18
	}
	{
		size := m.RedemptionRate.Size()
		i -= size
		if _, err := m.RedemptionRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.DepositDenom) > 0 {
		i -= len(m.DepositDenom)
		copy(dAtA[i:], m.DepositDenom)
		i = encodeVarintParams(dAtA, i, uint64(len(m.DepositDenom)))
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
	l = len(m.DepositDenom)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.RedemptionRate.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.EpochLength != 0 {
		n += 1 + sovParams(uint64(m.EpochLength))
	}
	l = m.InterestRate.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.InterestRateMax.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.InterestRateMin.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.InterestRateIncrease.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.InterestRateDecrease.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.IncrementalInterestPaymentFundPercentage.Size()
	n += 1 + l + sovParams(uint64(l))
	l = len(m.IncrementalInterestPaymentFundAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.IncrementalInterestPaymentEnabled {
		n += 2
	}
	l = m.HealthGainFactor.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.TotalValue.Size()
	n += 1 + l + sovParams(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field DepositDenom", wireType)
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
			m.DepositDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RedemptionRate", wireType)
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
			if err := m.RedemptionRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
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
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRate", wireType)
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
			if err := m.InterestRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRateMax", wireType)
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
			if err := m.InterestRateMax.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRateMin", wireType)
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
			if err := m.InterestRateMin.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRateIncrease", wireType)
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
			if err := m.InterestRateIncrease.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRateDecrease", wireType)
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
			if err := m.InterestRateDecrease.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncrementalInterestPaymentFundPercentage", wireType)
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
			if err := m.IncrementalInterestPaymentFundPercentage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncrementalInterestPaymentFundAddress", wireType)
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
			m.IncrementalInterestPaymentFundAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncrementalInterestPaymentEnabled", wireType)
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
			m.IncrementalInterestPaymentEnabled = bool(v != 0)
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HealthGainFactor", wireType)
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
			if err := m.HealthGainFactor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalValue", wireType)
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
			if err := m.TotalValue.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
