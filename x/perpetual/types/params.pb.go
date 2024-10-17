// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/perpetual/params.proto

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
	LeverageMax                                    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,1,opt,name=leverage_max,json=leverageMax,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"leverage_max"`
	BorrowInterestRateMax                          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=borrow_interest_rate_max,json=borrowInterestRateMax,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"borrow_interest_rate_max"`
	BorrowInterestRateMin                          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=borrow_interest_rate_min,json=borrowInterestRateMin,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"borrow_interest_rate_min"`
	MinBorrowInterestAmount                        github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,4,opt,name=min_borrow_interest_amount,json=minBorrowInterestAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"min_borrow_interest_amount"`
	BorrowInterestRateIncrease                     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=borrow_interest_rate_increase,json=borrowInterestRateIncrease,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"borrow_interest_rate_increase"`
	BorrowInterestRateDecrease                     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=borrow_interest_rate_decrease,json=borrowInterestRateDecrease,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"borrow_interest_rate_decrease"`
	HealthGainFactor                               github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=health_gain_factor,json=healthGainFactor,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"health_gain_factor"`
	MaxOpenPositions                               int64                                  `protobuf:"varint,8,opt,name=max_open_positions,json=maxOpenPositions,proto3" json:"max_open_positions,omitempty"`
	PoolOpenThreshold                              github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,9,opt,name=pool_open_threshold,json=poolOpenThreshold,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pool_open_threshold"`
	ForceCloseFundPercentage                       github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,10,opt,name=force_close_fund_percentage,json=forceCloseFundPercentage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"force_close_fund_percentage"`
	ForceCloseFundAddress                          string                                 `protobuf:"bytes,11,opt,name=force_close_fund_address,json=forceCloseFundAddress,proto3" json:"force_close_fund_address,omitempty"`
	IncrementalBorrowInterestPaymentFundPercentage github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,12,opt,name=incremental_borrow_interest_payment_fund_percentage,json=incrementalBorrowInterestPaymentFundPercentage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"incremental_borrow_interest_payment_fund_percentage"`
	IncrementalBorrowInterestPaymentFundAddress    string                                 `protobuf:"bytes,13,opt,name=incremental_borrow_interest_payment_fund_address,json=incrementalBorrowInterestPaymentFundAddress,proto3" json:"incremental_borrow_interest_payment_fund_address,omitempty"`
	SafetyFactor                                   github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,14,opt,name=safety_factor,json=safetyFactor,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"safety_factor"`
	IncrementalBorrowInterestPaymentEnabled        bool                                   `protobuf:"varint,15,opt,name=incremental_borrow_interest_payment_enabled,json=incrementalBorrowInterestPaymentEnabled,proto3" json:"incremental_borrow_interest_payment_enabled,omitempty"`
	WhitelistingEnabled                            bool                                   `protobuf:"varint,16,opt,name=whitelisting_enabled,json=whitelistingEnabled,proto3" json:"whitelisting_enabled,omitempty"`
	TakeProfitBorrowInterestRateMin                github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,17,opt,name=take_profit_borrow_interest_rate_min,json=takeProfitBorrowInterestRateMin,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"take_profit_borrow_interest_rate_min"`
	PerpetualSwapFee                               github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,18,opt,name=perpetual_swap_fee,json=perpetualSwapFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"perpetual_swap_fee"`
	MaxLimitOrder                                  int64                                  `protobuf:"varint,19,opt,name=max_limit_order,json=maxLimitOrder,proto3" json:"max_limit_order,omitempty"`
	FixedFundingRate                               github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,20,opt,name=fixed_funding_rate,json=fixedFundingRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"fixed_funding_rate"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_032040a5fba82242, []int{0}
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

func (m *Params) GetForceCloseFundAddress() string {
	if m != nil {
		return m.ForceCloseFundAddress
	}
	return ""
}

func (m *Params) GetIncrementalBorrowInterestPaymentFundAddress() string {
	if m != nil {
		return m.IncrementalBorrowInterestPaymentFundAddress
	}
	return ""
}

func (m *Params) GetIncrementalBorrowInterestPaymentEnabled() bool {
	if m != nil {
		return m.IncrementalBorrowInterestPaymentEnabled
	}
	return false
}

func (m *Params) GetWhitelistingEnabled() bool {
	if m != nil {
		return m.WhitelistingEnabled
	}
	return false
}

func (m *Params) GetMaxLimitOrder() int64 {
	if m != nil {
		return m.MaxLimitOrder
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "elys.perpetual.Params")
}

func init() { proto.RegisterFile("elys/perpetual/params.proto", fileDescriptor_032040a5fba82242) }

var fileDescriptor_032040a5fba82242 = []byte{
	// 701 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x95, 0xcf, 0x4f, 0x1b, 0x39,
	0x14, 0xc7, 0x33, 0x0b, 0xcb, 0x82, 0xf9, 0x6d, 0x40, 0x6b, 0x81, 0x36, 0x41, 0xab, 0x15, 0x8b,
	0xc4, 0x92, 0xec, 0x8a, 0xc3, 0x4a, 0xbd, 0x91, 0xd2, 0xb4, 0x48, 0xad, 0x48, 0x43, 0x4f, 0x15,
	0xaa, 0xe5, 0xcc, 0xbc, 0x4c, 0xac, 0xcc, 0xd8, 0x53, 0xdb, 0x69, 0x92, 0x43, 0xff, 0x87, 0x1e,
	0x2b, 0xf5, 0xd2, 0x3f, 0x87, 0x23, 0xc7, 0xaa, 0x07, 0x54, 0xc1, 0xa5, 0x7f, 0x46, 0x65, 0x4f,
	0x12, 0x05, 0x02, 0x2d, 0x1d, 0x71, 0x4a, 0xc6, 0xcf, 0xfe, 0x7c, 0xdf, 0xd7, 0x7a, 0xcf, 0x0f,
	0x6d, 0x40, 0xd4, 0xd3, 0xa5, 0x04, 0x54, 0x02, 0xa6, 0xcd, 0xa2, 0x52, 0xc2, 0x14, 0x8b, 0x75,
	0x31, 0x51, 0xd2, 0x48, 0xbc, 0x60, 0x83, 0xc5, 0x61, 0x70, 0x7d, 0x35, 0x94, 0xa1, 0x74, 0xa1,
	0x92, 0xfd, 0x97, 0xee, 0xfa, 0xf3, 0xeb, 0x02, 0x9a, 0xaa, 0xba, 0x63, 0xf8, 0x39, 0x9a, 0x8b,
	0xe0, 0x0d, 0x28, 0x16, 0x02, 0x8d, 0x59, 0x97, 0x78, 0x9b, 0xde, 0xf6, 0x4c, 0xb9, 0x78, 0x7a,
	0x5e, 0xc8, 0x7d, 0x3e, 0x2f, 0x6c, 0x85, 0xdc, 0x34, 0xdb, 0xf5, 0xa2, 0x2f, 0xe3, 0x92, 0x2f,
	0x75, 0x2c, 0x75, 0xff, 0x67, 0x57, 0x07, 0xad, 0x92, 0xe9, 0x25, 0xa0, 0x8b, 0x07, 0xe0, 0xd7,
	0x66, 0x07, 0x8c, 0x67, 0xac, 0x8b, 0x43, 0x44, 0xea, 0x52, 0x29, 0xd9, 0xa1, 0x5c, 0x18, 0x50,
	0xa0, 0x0d, 0x55, 0xcc, 0xa4, 0xf8, 0x5f, 0x32, 0xe1, 0xd7, 0x52, 0xde, 0x61, 0x1f, 0x57, 0x63,
	0xe6, 0xfb, 0x42, 0x5c, 0x90, 0x89, 0x7b, 0x13, 0xe2, 0x02, 0xb7, 0xd0, 0x7a, 0xcc, 0x05, 0xbd,
	0x2e, 0xc6, 0x62, 0xd9, 0x16, 0x86, 0x4c, 0xfe, 0xb4, 0xd4, 0xa1, 0x30, 0xb5, 0xdf, 0x63, 0x2e,
	0xca, 0x57, 0xd4, 0xf6, 0x1d, 0x0e, 0xbf, 0x46, 0x7f, 0xdc, 0xe8, 0x8a, 0x0b, 0x5f, 0x01, 0xd3,
	0x40, 0x7e, 0xcd, 0x64, 0x6d, 0x7d, 0xdc, 0xda, 0x61, 0x9f, 0x78, 0xab, 0x64, 0x00, 0x7d, 0xc9,
	0xa9, 0xfb, 0x92, 0x3c, 0xe8, 0x13, 0xf1, 0x09, 0xc2, 0x4d, 0x60, 0x91, 0x69, 0xd2, 0x90, 0x71,
	0x41, 0x1b, 0xcc, 0x37, 0x52, 0x91, 0xdf, 0x32, 0xe9, 0x2c, 0xa5, 0xa4, 0xc7, 0x8c, 0x8b, 0x8a,
	0xe3, 0xe0, 0x7f, 0x10, 0x8e, 0x59, 0x97, 0xca, 0x04, 0x04, 0x4d, 0xa4, 0xe6, 0x86, 0x4b, 0xa1,
	0xc9, 0xf4, 0xa6, 0xb7, 0x3d, 0x51, 0x5b, 0x8a, 0x59, 0xf7, 0x28, 0x01, 0x51, 0x1d, 0xac, 0xe3,
	0x57, 0x68, 0x25, 0x91, 0x32, 0x4a, 0xb7, 0x9b, 0xa6, 0x02, 0xdd, 0x94, 0x51, 0x40, 0x66, 0x32,
	0x25, 0xb3, 0x6c, 0x51, 0x96, 0xff, 0x62, 0x00, 0xc2, 0x31, 0xda, 0x68, 0x48, 0xe5, 0x03, 0xf5,
	0x23, 0xa9, 0x81, 0x36, 0xda, 0x22, 0xa0, 0x09, 0x28, 0x1f, 0x84, 0x61, 0x21, 0x10, 0x94, 0x49,
	0x87, 0x38, 0xe4, 0x43, 0x4b, 0xac, 0xb4, 0x45, 0x50, 0x1d, 0xf2, 0xf0, 0xff, 0x88, 0x8c, 0xc9,
	0xb1, 0x20, 0x50, 0xa0, 0x35, 0x99, 0xb5, 0x5a, 0xb5, 0xb5, 0xab, 0x67, 0xf7, 0xd3, 0x20, 0xfe,
	0xe0, 0xa1, 0x3d, 0x57, 0x65, 0xb1, 0x25, 0x45, 0x63, 0xf5, 0x9e, 0xb0, 0x9e, 0x8d, 0x8c, 0x19,
	0x98, 0xcb, 0x64, 0xa0, 0x38, 0x22, 0x75, 0xb5, 0x11, 0xaa, 0xa9, 0xce, 0x35, 0x5b, 0x80, 0xfe,
	0xbd, 0x73, 0x72, 0x03, 0xbb, 0xf3, 0xce, 0xee, 0xce, 0x5d, 0x94, 0x06, 0x97, 0x70, 0x8c, 0xe6,
	0x35, 0x6b, 0x80, 0xe9, 0x0d, 0x6a, 0x72, 0x21, 0x93, 0xbb, 0xb9, 0x14, 0xd2, 0xaf, 0xc7, 0x13,
	0xb4, 0x73, 0x97, 0xdc, 0x41, 0xb0, 0x7a, 0x04, 0x01, 0x59, 0xdc, 0xf4, 0xb6, 0xa7, 0x6b, 0x7f,
	0xff, 0x28, 0xed, 0x47, 0xe9, 0x76, 0xfc, 0x1f, 0x5a, 0xed, 0x34, 0xb9, 0x81, 0x88, 0x6b, 0xc3,
	0x45, 0x38, 0xc4, 0x2c, 0x39, 0xcc, 0xca, 0x68, 0x6c, 0x70, 0xe4, 0x2d, 0xfa, 0xcb, 0xb0, 0x16,
	0xd0, 0x44, 0xc9, 0x06, 0x37, 0xf4, 0xd6, 0x67, 0x74, 0x39, 0x93, 0xf9, 0x82, 0x65, 0x57, 0x1d,
	0xba, 0x7c, 0xe3, 0x83, 0x7a, 0x82, 0xf0, 0x70, 0x46, 0x51, 0xdd, 0x61, 0x09, 0x6d, 0x00, 0x10,
	0x9c, 0xad, 0xfb, 0x87, 0xa4, 0xe3, 0x0e, 0x4b, 0x2a, 0x00, 0x78, 0x0b, 0x2d, 0xda, 0xee, 0x8f,
	0x78, 0xcc, 0x0d, 0x95, 0x2a, 0x00, 0x45, 0x56, 0x5c, 0xeb, 0xcf, 0xc7, 0xac, 0xfb, 0xd4, 0xae,
	0x1e, 0xd9, 0x45, 0x9b, 0x45, 0x83, 0x77, 0x21, 0x70, 0x35, 0x63, 0x2f, 0xce, 0xda, 0x26, 0xab,
	0xd9, 0xb2, 0x70, 0xa4, 0x4a, 0x0a, 0xb2, 0x36, 0x1f, 0x4c, 0xbe, 0xff, 0x58, 0xc8, 0x95, 0x9f,
	0x9c, 0x5e, 0xe4, 0xbd, 0xb3, 0x8b, 0xbc, 0xf7, 0xe5, 0x22, 0xef, 0xbd, 0xbb, 0xcc, 0xe7, 0xce,
	0x2e, 0xf3, 0xb9, 0x4f, 0x97, 0xf9, 0xdc, 0xcb, 0xe2, 0x08, 0xd9, 0x4e, 0xed, 0x5d, 0x01, 0xa6,
	0x23, 0x55, 0xcb, 0x7d, 0x94, 0xba, 0x23, 0x13, 0xde, 0xa9, 0xd4, 0xa7, 0xdc, 0xec, 0xde, 0xfb,
	0x16, 0x00, 0x00, 0xff, 0xff, 0x13, 0x28, 0x52, 0x50, 0x00, 0x08, 0x00, 0x00,
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
		size := m.FixedFundingRate.Size()
		i -= size
		if _, err := m.FixedFundingRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1
	i--
	dAtA[i] = 0xa2
	if m.MaxLimitOrder != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxLimitOrder))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x98
	}
	{
		size := m.PerpetualSwapFee.Size()
		i -= size
		if _, err := m.PerpetualSwapFee.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1
	i--
	dAtA[i] = 0x92
	{
		size := m.TakeProfitBorrowInterestRateMin.Size()
		i -= size
		if _, err := m.TakeProfitBorrowInterestRateMin.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1
	i--
	dAtA[i] = 0x8a
	if m.WhitelistingEnabled {
		i--
		if m.WhitelistingEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x80
	}
	if m.IncrementalBorrowInterestPaymentEnabled {
		i--
		if m.IncrementalBorrowInterestPaymentEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x78
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
	dAtA[i] = 0x72
	if len(m.IncrementalBorrowInterestPaymentFundAddress) > 0 {
		i -= len(m.IncrementalBorrowInterestPaymentFundAddress)
		copy(dAtA[i:], m.IncrementalBorrowInterestPaymentFundAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.IncrementalBorrowInterestPaymentFundAddress)))
		i--
		dAtA[i] = 0x6a
	}
	{
		size := m.IncrementalBorrowInterestPaymentFundPercentage.Size()
		i -= size
		if _, err := m.IncrementalBorrowInterestPaymentFundPercentage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x62
	if len(m.ForceCloseFundAddress) > 0 {
		i -= len(m.ForceCloseFundAddress)
		copy(dAtA[i:], m.ForceCloseFundAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ForceCloseFundAddress)))
		i--
		dAtA[i] = 0x5a
	}
	{
		size := m.ForceCloseFundPercentage.Size()
		i -= size
		if _, err := m.ForceCloseFundPercentage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	{
		size := m.PoolOpenThreshold.Size()
		i -= size
		if _, err := m.PoolOpenThreshold.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	if m.MaxOpenPositions != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxOpenPositions))
		i--
		dAtA[i] = 0x40
	}
	{
		size := m.HealthGainFactor.Size()
		i -= size
		if _, err := m.HealthGainFactor.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.BorrowInterestRateDecrease.Size()
		i -= size
		if _, err := m.BorrowInterestRateDecrease.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.BorrowInterestRateIncrease.Size()
		i -= size
		if _, err := m.BorrowInterestRateIncrease.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.MinBorrowInterestAmount.Size()
		i -= size
		if _, err := m.MinBorrowInterestAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.BorrowInterestRateMin.Size()
		i -= size
		if _, err := m.BorrowInterestRateMin.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.BorrowInterestRateMax.Size()
		i -= size
		if _, err := m.BorrowInterestRateMax.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
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
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.LeverageMax.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.BorrowInterestRateMax.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.BorrowInterestRateMin.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MinBorrowInterestAmount.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.BorrowInterestRateIncrease.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.BorrowInterestRateDecrease.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.HealthGainFactor.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.MaxOpenPositions != 0 {
		n += 1 + sovParams(uint64(m.MaxOpenPositions))
	}
	l = m.PoolOpenThreshold.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.ForceCloseFundPercentage.Size()
	n += 1 + l + sovParams(uint64(l))
	l = len(m.ForceCloseFundAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.IncrementalBorrowInterestPaymentFundPercentage.Size()
	n += 1 + l + sovParams(uint64(l))
	l = len(m.IncrementalBorrowInterestPaymentFundAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.SafetyFactor.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.IncrementalBorrowInterestPaymentEnabled {
		n += 2
	}
	if m.WhitelistingEnabled {
		n += 3
	}
	l = m.TakeProfitBorrowInterestRateMin.Size()
	n += 2 + l + sovParams(uint64(l))
	l = m.PerpetualSwapFee.Size()
	n += 2 + l + sovParams(uint64(l))
	if m.MaxLimitOrder != 0 {
		n += 2 + sovParams(uint64(m.MaxLimitOrder))
	}
	l = m.FixedFundingRate.Size()
	n += 2 + l + sovParams(uint64(l))
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BorrowInterestRateMax", wireType)
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
			if err := m.BorrowInterestRateMax.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BorrowInterestRateMin", wireType)
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
			if err := m.BorrowInterestRateMin.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinBorrowInterestAmount", wireType)
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
			if err := m.MinBorrowInterestAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BorrowInterestRateIncrease", wireType)
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
			if err := m.BorrowInterestRateIncrease.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BorrowInterestRateDecrease", wireType)
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
			if err := m.BorrowInterestRateDecrease.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
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
		case 8:
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
		case 9:
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
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ForceCloseFundPercentage", wireType)
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
			if err := m.ForceCloseFundPercentage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ForceCloseFundAddress", wireType)
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
			m.ForceCloseFundAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncrementalBorrowInterestPaymentFundPercentage", wireType)
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
			if err := m.IncrementalBorrowInterestPaymentFundPercentage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncrementalBorrowInterestPaymentFundAddress", wireType)
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
			m.IncrementalBorrowInterestPaymentFundAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 14:
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
		case 15:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IncrementalBorrowInterestPaymentEnabled", wireType)
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
			m.IncrementalBorrowInterestPaymentEnabled = bool(v != 0)
		case 16:
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
		case 17:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TakeProfitBorrowInterestRateMin", wireType)
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
			if err := m.TakeProfitBorrowInterestRateMin.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 18:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PerpetualSwapFee", wireType)
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
			if err := m.PerpetualSwapFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 19:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxLimitOrder", wireType)
			}
			m.MaxLimitOrder = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxLimitOrder |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 20:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FixedFundingRate", wireType)
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
			if err := m.FixedFundingRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
