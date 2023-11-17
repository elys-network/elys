// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/margin/params.proto

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
	LeverageMax                              github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,1,opt,name=leverage_max,json=leverageMax,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"leverage_max"`
	InterestRateMax                          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=interest_rate_max,json=interestRateMax,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate_max"`
	InterestRateMin                          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=interest_rate_min,json=interestRateMin,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate_min"`
	InterestRateIncrease                     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=interest_rate_increase,json=interestRateIncrease,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate_increase"`
	InterestRateDecrease                     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=interest_rate_decrease,json=interestRateDecrease,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate_decrease"`
	HealthGainFactor                         github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=health_gain_factor,json=healthGainFactor,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"health_gain_factor"`
	EpochLength                              int64                                  `protobuf:"varint,7,opt,name=epoch_length,json=epochLength,proto3" json:"epoch_length,omitempty"`
	RemovalQueueThreshold                    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,8,opt,name=removal_queue_threshold,json=removalQueueThreshold,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"removal_queue_threshold"`
	MaxOpenPositions                         int64                                  `protobuf:"varint,9,opt,name=max_open_positions,json=maxOpenPositions,proto3" json:"max_open_positions,omitempty"`
	PoolOpenThreshold                        github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,10,opt,name=pool_open_threshold,json=poolOpenThreshold,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"pool_open_threshold"`
	ForceCloseFundPercentage                 github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,11,opt,name=force_close_fund_percentage,json=forceCloseFundPercentage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"force_close_fund_percentage"`
	ForceCloseFundAddress                    string                                 `protobuf:"bytes,12,opt,name=force_close_fund_address,json=forceCloseFundAddress,proto3" json:"force_close_fund_address,omitempty"`
	IncrementalInterestPaymentFundPercentage github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,13,opt,name=incremental_interest_payment_fund_percentage,json=incrementalInterestPaymentFundPercentage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"incremental_interest_payment_fund_percentage"`
	IncrementalInterestPaymentFundAddress    string                                 `protobuf:"bytes,14,opt,name=incremental_interest_payment_fund_address,json=incrementalInterestPaymentFundAddress,proto3" json:"incremental_interest_payment_fund_address,omitempty"`
	SqModifier                               github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,15,opt,name=sq_modifier,json=sqModifier,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"sq_modifier"`
	SafetyFactor                             github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,16,opt,name=safety_factor,json=safetyFactor,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"safety_factor"`
	IncrementalInterestPaymentEnabled        bool                                   `protobuf:"varint,17,opt,name=incremental_interest_payment_enabled,json=incrementalInterestPaymentEnabled,proto3" json:"incremental_interest_payment_enabled,omitempty"`
	WhitelistingEnabled                      bool                                   `protobuf:"varint,18,opt,name=whitelisting_enabled,json=whitelistingEnabled,proto3" json:"whitelisting_enabled,omitempty"`
	InvariantCheckEpoch                      string                                 `protobuf:"bytes,19,opt,name=invariant_check_epoch,json=invariantCheckEpoch,proto3" json:"invariant_check_epoch,omitempty"`
	BrokerAddress                            string                                 `protobuf:"bytes,20,opt,name=broker_address,json=brokerAddress,proto3" json:"broker_address,omitempty"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_f427d3667a99d828, []int{0}
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

func (m *Params) GetEpochLength() int64 {
	if m != nil {
		return m.EpochLength
	}
	return 0
}

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

func (m *Params) GetWhitelistingEnabled() bool {
	if m != nil {
		return m.WhitelistingEnabled
	}
	return false
}

func (m *Params) GetInvariantCheckEpoch() string {
	if m != nil {
		return m.InvariantCheckEpoch
	}
	return ""
}

func (m *Params) GetBrokerAddress() string {
	if m != nil {
		return m.BrokerAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*Params)(nil), "elys.margin.Params")
}

func init() { proto.RegisterFile("elys/margin/params.proto", fileDescriptor_f427d3667a99d828) }

var fileDescriptor_f427d3667a99d828 = []byte{
	// 678 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x95, 0xcf, 0x6e, 0xd3, 0x4a,
	0x14, 0xc6, 0xe3, 0xdb, 0xde, 0xde, 0x76, 0x92, 0xfe, 0x9b, 0xa4, 0x17, 0x0b, 0xa4, 0xb4, 0x45,
	0x14, 0x05, 0xd1, 0x26, 0x02, 0x16, 0x48, 0xec, 0xe8, 0x3f, 0x54, 0x89, 0xaa, 0x69, 0x60, 0x81,
	0x2a, 0xc4, 0x68, 0x62, 0x9f, 0xd8, 0xa3, 0xd8, 0x33, 0xee, 0xcc, 0xa4, 0x4d, 0x5e, 0x02, 0xb1,
	0x64, 0xc9, 0xe3, 0x74, 0xd9, 0x25, 0x62, 0x51, 0xa1, 0xf6, 0x15, 0x78, 0x00, 0xe4, 0xb1, 0x9d,
	0xa6, 0x54, 0x2d, 0x92, 0x61, 0x95, 0xf8, 0x9c, 0x33, 0xbf, 0xef, 0xfb, 0xc6, 0x1e, 0x0d, 0xb2,
	0x21, 0x18, 0xa8, 0x46, 0x48, 0xa5, 0xc7, 0x78, 0x23, 0xa2, 0x92, 0x86, 0xaa, 0x1e, 0x49, 0xa1,
	0x05, 0x2e, 0xc6, 0x9d, 0x7a, 0xd2, 0xb9, 0x5b, 0xf1, 0x84, 0x27, 0x4c, 0xbd, 0x11, 0xff, 0x4b,
	0x46, 0xee, 0xff, 0x28, 0xa1, 0x89, 0xa6, 0x59, 0x83, 0xf7, 0x51, 0x29, 0x80, 0x23, 0x90, 0xd4,
	0x03, 0x12, 0xd2, 0xbe, 0x6d, 0x2d, 0x59, 0xb5, 0xa9, 0xf5, 0xfa, 0xc9, 0xd9, 0x62, 0xe1, 0xdb,
	0xd9, 0xe2, 0x43, 0x8f, 0x69, 0xbf, 0xd7, 0xae, 0x3b, 0x22, 0x6c, 0x38, 0x42, 0x85, 0x42, 0xa5,
	0x3f, 0x6b, 0xca, 0xed, 0x36, 0xf4, 0x20, 0x02, 0x55, 0xdf, 0x04, 0xa7, 0x55, 0xcc, 0x18, 0xbb,
	0xb4, 0x8f, 0x0f, 0xd0, 0x3c, 0xe3, 0x1a, 0x24, 0x28, 0x4d, 0x24, 0xd5, 0x09, 0xf7, 0x9f, 0x5c,
	0xdc, 0xd9, 0x0c, 0xd4, 0xa2, 0xfa, 0x06, 0x36, 0xe3, 0xf6, 0xd8, 0x5f, 0x60, 0x33, 0x8e, 0x5d,
	0xf4, 0xff, 0x55, 0x36, 0xe3, 0x8e, 0x04, 0xaa, 0xc0, 0x1e, 0xcf, 0x25, 0x50, 0x19, 0x15, 0xd8,
	0x49, 0x59, 0xd7, 0x55, 0x5c, 0x48, 0x55, 0xfe, 0xfd, 0x73, 0x95, 0xcd, 0x94, 0x85, 0xdf, 0x23,
	0xec, 0x03, 0x0d, 0xb4, 0x4f, 0x3c, 0xca, 0x38, 0xe9, 0x50, 0x47, 0x0b, 0x69, 0x4f, 0xe4, 0x52,
	0x98, 0x4b, 0x48, 0xaf, 0x28, 0xe3, 0xdb, 0x86, 0x83, 0x97, 0x51, 0x09, 0x22, 0xe1, 0xf8, 0x24,
	0x00, 0xee, 0x69, 0xdf, 0xfe, 0x6f, 0xc9, 0xaa, 0x8d, 0xb5, 0x8a, 0xa6, 0xf6, 0xda, 0x94, 0x70,
	0x07, 0xdd, 0x91, 0x10, 0x8a, 0x23, 0x1a, 0x90, 0xc3, 0x1e, 0xf4, 0x80, 0x68, 0x5f, 0x82, 0xf2,
	0x45, 0xe0, 0xda, 0x93, 0xb9, 0x5c, 0x2c, 0xa4, 0xb8, 0xfd, 0x98, 0xf6, 0x36, 0x83, 0xe1, 0x55,
	0x84, 0x43, 0xda, 0x27, 0x22, 0x02, 0x4e, 0x22, 0xa1, 0x98, 0x66, 0x82, 0x2b, 0x7b, 0xca, 0x18,
	0x9a, 0x0b, 0x69, 0x7f, 0x2f, 0x02, 0xde, 0xcc, 0xea, 0xf8, 0x03, 0x2a, 0x47, 0x42, 0x04, 0xc9,
	0xf8, 0xa5, 0x23, 0x94, 0xcb, 0xd1, 0x7c, 0x8c, 0x8a, 0xf9, 0x97, 0x6e, 0x42, 0x74, 0xaf, 0x23,
	0xa4, 0x03, 0xc4, 0x09, 0x84, 0x02, 0xd2, 0xe9, 0x71, 0x97, 0x44, 0x20, 0x1d, 0xe0, 0x9a, 0x7a,
	0x60, 0x17, 0x73, 0xe9, 0xd8, 0x06, 0xb9, 0x11, 0x13, 0xb7, 0x7b, 0xdc, 0x6d, 0x0e, 0x79, 0xf8,
	0x39, 0xb2, 0xaf, 0xc9, 0x51, 0xd7, 0x95, 0xa0, 0x94, 0x5d, 0x8a, 0xb5, 0x5a, 0x0b, 0x57, 0xd7,
	0xbe, 0x4c, 0x9a, 0xf8, 0xa3, 0x85, 0x56, 0xcd, 0xd7, 0x1d, 0xc6, 0xa4, 0x80, 0x0c, 0xbf, 0xc8,
	0x88, 0x0e, 0xe2, 0xd2, 0x35, 0xe7, 0xd3, 0xb9, 0x9c, 0xd7, 0x46, 0x34, 0x76, 0x52, 0x89, 0x66,
	0xa2, 0xf0, 0x4b, 0x92, 0x77, 0xe8, 0xd1, 0xef, 0xfd, 0x64, 0xd1, 0x66, 0x4c, 0xb4, 0x95, 0xdb,
	0xe1, 0x59, 0xd4, 0x3d, 0x54, 0x54, 0x87, 0x24, 0x14, 0x2e, 0xeb, 0x30, 0x90, 0xf6, 0x6c, 0xae,
	0x20, 0x48, 0x1d, 0xee, 0xa6, 0x04, 0xfc, 0x06, 0x4d, 0x2b, 0xda, 0x01, 0x3d, 0xc8, 0x4e, 0xd5,
	0x5c, 0x2e, 0x64, 0x29, 0x81, 0xa4, 0x27, 0x6a, 0x0f, 0x3d, 0xb8, 0x35, 0x3f, 0x70, 0xda, 0x0e,
	0xc0, 0xb5, 0xe7, 0x97, 0xac, 0xda, 0x64, 0x6b, 0xf9, 0xe6, 0xe8, 0x5b, 0xc9, 0x20, 0x7e, 0x82,
	0x2a, 0xc7, 0x3e, 0xd3, 0x10, 0x30, 0xa5, 0x19, 0xf7, 0x86, 0x00, 0x6c, 0x00, 0xe5, 0xd1, 0x5e,
	0xb6, 0xe4, 0x29, 0x5a, 0x60, 0xfc, 0x88, 0x4a, 0x46, 0xb9, 0x26, 0x8e, 0x0f, 0x4e, 0x97, 0x98,
	0x13, 0x6d, 0x97, 0xcd, 0x7e, 0x97, 0x87, 0xcd, 0x8d, 0xb8, 0xb7, 0x15, 0xb7, 0xf0, 0x0a, 0x9a,
	0x69, 0x4b, 0xd1, 0x05, 0x39, 0x7c, 0x39, 0x15, 0x33, 0x3c, 0x9d, 0x54, 0xd3, 0x97, 0xf0, 0x62,
	0xfc, 0xf3, 0x97, 0xc5, 0xc2, 0xfa, 0xd6, 0xc9, 0x79, 0xd5, 0x3a, 0x3d, 0xaf, 0x5a, 0xdf, 0xcf,
	0xab, 0xd6, 0xa7, 0x8b, 0x6a, 0xe1, 0xf4, 0xa2, 0x5a, 0xf8, 0x7a, 0x51, 0x2d, 0x1c, 0x3c, 0x1e,
	0xd9, 0xb4, 0xf8, 0xfa, 0x5a, 0xe3, 0xa0, 0x8f, 0x85, 0xec, 0x9a, 0x87, 0x46, 0x3f, 0xbb, 0xe7,
	0xcc, 0xee, 0xb5, 0x27, 0xcc, 0x25, 0xf6, 0xec, 0x67, 0x00, 0x00, 0x00, 0xff, 0xff, 0x9d, 0x03,
	0x75, 0xf4, 0x03, 0x07, 0x00, 0x00,
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
	if len(m.BrokerAddress) > 0 {
		i -= len(m.BrokerAddress)
		copy(dAtA[i:], m.BrokerAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.BrokerAddress)))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0xa2
	}
	if len(m.InvariantCheckEpoch) > 0 {
		i -= len(m.InvariantCheckEpoch)
		copy(dAtA[i:], m.InvariantCheckEpoch)
		i = encodeVarintParams(dAtA, i, uint64(len(m.InvariantCheckEpoch)))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x9a
	}
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
		dAtA[i] = 0x90
	}
	if m.IncrementalInterestPaymentEnabled {
		i--
		if m.IncrementalInterestPaymentEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x88
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
	dAtA[i] = 0x1
	i--
	dAtA[i] = 0x82
	{
		size := m.SqModifier.Size()
		i -= size
		if _, err := m.SqModifier.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x7a
	if len(m.IncrementalInterestPaymentFundAddress) > 0 {
		i -= len(m.IncrementalInterestPaymentFundAddress)
		copy(dAtA[i:], m.IncrementalInterestPaymentFundAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.IncrementalInterestPaymentFundAddress)))
		i--
		dAtA[i] = 0x72
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
	dAtA[i] = 0x6a
	if len(m.ForceCloseFundAddress) > 0 {
		i -= len(m.ForceCloseFundAddress)
		copy(dAtA[i:], m.ForceCloseFundAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ForceCloseFundAddress)))
		i--
		dAtA[i] = 0x62
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
	dAtA[i] = 0x5a
	{
		size := m.PoolOpenThreshold.Size()
		i -= size
		if _, err := m.PoolOpenThreshold.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	if m.MaxOpenPositions != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxOpenPositions))
		i--
		dAtA[i] = 0x48
	}
	{
		size := m.RemovalQueueThreshold.Size()
		i -= size
		if _, err := m.RemovalQueueThreshold.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	if m.EpochLength != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.EpochLength))
		i--
		dAtA[i] = 0x38
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
	dAtA[i] = 0x32
	{
		size := m.InterestRateDecrease.Size()
		i -= size
		if _, err := m.InterestRateDecrease.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.InterestRateIncrease.Size()
		i -= size
		if _, err := m.InterestRateIncrease.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.InterestRateMin.Size()
		i -= size
		if _, err := m.InterestRateMin.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.InterestRateMax.Size()
		i -= size
		if _, err := m.InterestRateMax.MarshalTo(dAtA[i:]); err != nil {
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
	l = m.InterestRateMax.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.InterestRateMin.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.InterestRateIncrease.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.InterestRateDecrease.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.HealthGainFactor.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.EpochLength != 0 {
		n += 1 + sovParams(uint64(m.EpochLength))
	}
	l = m.RemovalQueueThreshold.Size()
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
	l = m.IncrementalInterestPaymentFundPercentage.Size()
	n += 1 + l + sovParams(uint64(l))
	l = len(m.IncrementalInterestPaymentFundAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.SqModifier.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.SafetyFactor.Size()
	n += 2 + l + sovParams(uint64(l))
	if m.IncrementalInterestPaymentEnabled {
		n += 3
	}
	if m.WhitelistingEnabled {
		n += 3
	}
	l = len(m.InvariantCheckEpoch)
	if l > 0 {
		n += 2 + l + sovParams(uint64(l))
	}
	l = len(m.BrokerAddress)
	if l > 0 {
		n += 2 + l + sovParams(uint64(l))
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
		case 3:
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
		case 4:
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
		case 5:
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
		case 6:
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
		case 7:
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
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RemovalQueueThreshold", wireType)
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
			if err := m.RemovalQueueThreshold.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
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
		case 10:
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
		case 11:
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
		case 12:
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
		case 13:
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
		case 14:
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
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SqModifier", wireType)
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
			if err := m.SqModifier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 16:
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
		case 17:
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
		case 18:
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
		case 19:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InvariantCheckEpoch", wireType)
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
			m.InvariantCheckEpoch = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 20:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BrokerAddress", wireType)
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
			m.BrokerAddress = string(dAtA[iNdEx:postIndex])
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
