// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/leveragelp/types.proto

package types

import (
	fmt "fmt"
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

type PositionRequest struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Id      uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *PositionRequest) Reset()         { *m = PositionRequest{} }
func (m *PositionRequest) String() string { return proto.CompactTextString(m) }
func (*PositionRequest) ProtoMessage()    {}
func (*PositionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_992d513dd201f55b, []int{0}
}
func (m *PositionRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PositionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PositionRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PositionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PositionRequest.Merge(m, src)
}
func (m *PositionRequest) XXX_Size() int {
	return m.Size()
}
func (m *PositionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PositionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PositionRequest proto.InternalMessageInfo

func (m *PositionRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *PositionRequest) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type Position struct {
	Address           string                                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Collateral        types.Coin                             `protobuf:"bytes,2,opt,name=collateral,proto3" json:"collateral"`
	Liabilities       github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,opt,name=liabilities,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"liabilities"`
	InterestPaid      github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,4,opt,name=interest_paid,json=interestPaid,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"interest_paid"`
	Leverage          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=leverage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"leverage"`
	LeveragedLpAmount github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=leveraged_lp_amount,json=leveragedLpAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"leveraged_lp_amount"`
	PositionHealth    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,7,opt,name=position_health,json=positionHealth,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"position_health"`
	Id                uint64                                 `protobuf:"varint,8,opt,name=id,proto3" json:"id,omitempty"`
	AmmPoolId         uint64                                 `protobuf:"varint,9,opt,name=amm_pool_id,json=ammPoolId,proto3" json:"amm_pool_id,omitempty"`
	StopLossPrice     github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,10,opt,name=stop_loss_price,json=stopLossPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"stop_loss_price"`
}

func (m *Position) Reset()         { *m = Position{} }
func (m *Position) String() string { return proto.CompactTextString(m) }
func (*Position) ProtoMessage()    {}
func (*Position) Descriptor() ([]byte, []int) {
	return fileDescriptor_992d513dd201f55b, []int{1}
}
func (m *Position) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Position) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Position.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Position) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Position.Merge(m, src)
}
func (m *Position) XXX_Size() int {
	return m.Size()
}
func (m *Position) XXX_DiscardUnknown() {
	xxx_messageInfo_Position.DiscardUnknown(m)
}

var xxx_messageInfo_Position proto.InternalMessageInfo

func (m *Position) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Position) GetCollateral() types.Coin {
	if m != nil {
		return m.Collateral
	}
	return types.Coin{}
}

func (m *Position) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Position) GetAmmPoolId() uint64 {
	if m != nil {
		return m.AmmPoolId
	}
	return 0
}

type AddressId struct {
	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Id      uint64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *AddressId) Reset()         { *m = AddressId{} }
func (m *AddressId) String() string { return proto.CompactTextString(m) }
func (*AddressId) ProtoMessage()    {}
func (*AddressId) Descriptor() ([]byte, []int) {
	return fileDescriptor_992d513dd201f55b, []int{2}
}
func (m *AddressId) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AddressId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AddressId.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AddressId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddressId.Merge(m, src)
}
func (m *AddressId) XXX_Size() int {
	return m.Size()
}
func (m *AddressId) XXX_DiscardUnknown() {
	xxx_messageInfo_AddressId.DiscardUnknown(m)
}

var xxx_messageInfo_AddressId proto.InternalMessageInfo

func (m *AddressId) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *AddressId) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type PositionAndInterest struct {
	Position            *Position                              `protobuf:"bytes,1,opt,name=position,proto3" json:"position,omitempty"`
	InterestRateHour    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=interest_rate_hour,json=interestRateHour,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate_hour"`
	InterestRateHourUsd github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=interest_rate_hour_usd,json=interestRateHourUsd,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate_hour_usd"`
}

func (m *PositionAndInterest) Reset()         { *m = PositionAndInterest{} }
func (m *PositionAndInterest) String() string { return proto.CompactTextString(m) }
func (*PositionAndInterest) ProtoMessage()    {}
func (*PositionAndInterest) Descriptor() ([]byte, []int) {
	return fileDescriptor_992d513dd201f55b, []int{3}
}
func (m *PositionAndInterest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PositionAndInterest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PositionAndInterest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PositionAndInterest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PositionAndInterest.Merge(m, src)
}
func (m *PositionAndInterest) XXX_Size() int {
	return m.Size()
}
func (m *PositionAndInterest) XXX_DiscardUnknown() {
	xxx_messageInfo_PositionAndInterest.DiscardUnknown(m)
}

var xxx_messageInfo_PositionAndInterest proto.InternalMessageInfo

func (m *PositionAndInterest) GetPosition() *Position {
	if m != nil {
		return m.Position
	}
	return nil
}

func init() {
	proto.RegisterType((*PositionRequest)(nil), "elys.leveragelp.PositionRequest")
	proto.RegisterType((*Position)(nil), "elys.leveragelp.Position")
	proto.RegisterType((*AddressId)(nil), "elys.leveragelp.AddressId")
	proto.RegisterType((*PositionAndInterest)(nil), "elys.leveragelp.PositionAndInterest")
}

func init() { proto.RegisterFile("elys/leveragelp/types.proto", fileDescriptor_992d513dd201f55b) }

var fileDescriptor_992d513dd201f55b = []byte{
	// 541 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0x9b, 0xd2, 0x6d, 0xad, 0xcb, 0x56, 0x70, 0x11, 0x0a, 0x43, 0xca, 0xa6, 0x1e, 0xd0,
	0x2e, 0x4b, 0x34, 0xd0, 0x4e, 0x1c, 0x50, 0x0b, 0x87, 0x15, 0xed, 0x50, 0x05, 0x01, 0x12, 0x42,
	0x44, 0x6e, 0xfc, 0xd4, 0x5a, 0x73, 0xf2, 0x42, 0xec, 0x0c, 0xf6, 0x29, 0xe0, 0xca, 0x37, 0xda,
	0x71, 0x47, 0xc4, 0x61, 0x42, 0xed, 0x17, 0x41, 0x49, 0xe3, 0x52, 0x6d, 0xd2, 0xa4, 0xe6, 0x94,
	0x58, 0xcf, 0xef, 0xe7, 0xff, 0xf3, 0xfb, 0x3f, 0x93, 0xa7, 0x20, 0x2f, 0x94, 0x27, 0xe1, 0x1c,
	0x52, 0x36, 0x01, 0x99, 0x78, 0xfa, 0x22, 0x01, 0xe5, 0x26, 0x29, 0x6a, 0xa4, 0x9d, 0x3c, 0xe8,
	0xfe, 0x0f, 0xee, 0x3e, 0x9a, 0xe0, 0x04, 0x8b, 0x98, 0x97, 0xff, 0x2d, 0xb6, 0xed, 0x3a, 0x21,
	0xaa, 0x08, 0x95, 0x37, 0x66, 0x0a, 0xbc, 0xf3, 0xa3, 0x31, 0x68, 0x76, 0xe4, 0x85, 0x28, 0xe2,
	0x45, 0xbc, 0xf7, 0x92, 0x74, 0x46, 0xa8, 0x84, 0x16, 0x18, 0xfb, 0xf0, 0x35, 0x03, 0xa5, 0xa9,
	0x4d, 0xb6, 0x18, 0xe7, 0x29, 0x28, 0x65, 0x5b, 0xfb, 0xd6, 0x41, 0xcb, 0x37, 0x4b, 0xba, 0x43,
	0xea, 0x82, 0xdb, 0xf5, 0x7d, 0xeb, 0xa0, 0xe1, 0xd7, 0x05, 0xef, 0xfd, 0xda, 0x20, 0x4d, 0x93,
	0x7d, 0x47, 0xda, 0x2b, 0x42, 0x42, 0x94, 0x92, 0x69, 0x48, 0x99, 0x2c, 0xd2, 0xdb, 0xcf, 0x9f,
	0xb8, 0x0b, 0x61, 0x6e, 0x2e, 0xcc, 0x2d, 0x85, 0xb9, 0xaf, 0x51, 0xc4, 0x83, 0xc6, 0xe5, 0xf5,
	0x5e, 0xcd, 0x5f, 0x49, 0xa1, 0x23, 0xd2, 0x96, 0x82, 0x8d, 0x85, 0x14, 0x5a, 0x80, 0xb2, 0xef,
	0xe5, 0xf8, 0x81, 0x9b, 0x6f, 0xfb, 0x73, 0xbd, 0xf7, 0x6c, 0x22, 0xf4, 0x34, 0x1b, 0xbb, 0x21,
	0x46, 0x5e, 0x59, 0xec, 0xe2, 0x73, 0xa8, 0xf8, 0x59, 0x79, 0x65, 0xc3, 0x58, 0xfb, 0xab, 0x08,
	0xfa, 0x8e, 0x6c, 0x8b, 0x58, 0x43, 0x0a, 0x4a, 0x07, 0x09, 0x13, 0xdc, 0x6e, 0x54, 0x62, 0xde,
	0x37, 0x90, 0x11, 0x13, 0x9c, 0xbe, 0x25, 0x4d, 0xd3, 0x0f, 0x7b, 0x63, 0x6d, 0xde, 0x1b, 0x08,
	0xfd, 0x65, 0x3e, 0xfd, 0x42, 0xba, 0xe6, 0x9f, 0x07, 0x32, 0x09, 0x58, 0x84, 0x59, 0xac, 0xed,
	0xcd, 0x4a, 0x32, 0x1f, 0x2e, 0x51, 0xa7, 0x49, 0xbf, 0x00, 0xd1, 0x8f, 0xa4, 0x93, 0x94, 0x9d,
	0x0b, 0xa6, 0xc0, 0xa4, 0x9e, 0xda, 0x5b, 0x95, 0x24, 0xef, 0x18, 0xcc, 0x49, 0x41, 0x29, 0x3d,
	0xd2, 0x34, 0x1e, 0xa1, 0x0e, 0x69, 0xb3, 0x28, 0x0a, 0x12, 0x44, 0x19, 0x08, 0x6e, 0xb7, 0x8a,
	0x40, 0x8b, 0x45, 0xd1, 0x08, 0x51, 0x0e, 0x39, 0xfd, 0x40, 0x3a, 0x4a, 0x63, 0x12, 0x48, 0x54,
	0x2a, 0x48, 0x52, 0x11, 0x82, 0x4d, 0x2a, 0x09, 0xd9, 0xce, 0x31, 0xa7, 0xa8, 0xd4, 0x28, 0x87,
	0xf4, 0x8e, 0x49, 0xab, 0xbf, 0xf0, 0xdf, 0x90, 0xaf, 0x61, 0xe9, 0x1f, 0x75, 0xd2, 0x35, 0x96,
	0xee, 0xc7, 0x7c, 0x58, 0xf6, 0x97, 0x1e, 0x93, 0xa6, 0x29, 0xb4, 0x40, 0xe4, 0x0e, 0xbe, 0x31,
	0x81, 0xee, 0x72, 0x90, 0x96, 0x5b, 0xe9, 0x67, 0x42, 0x97, 0x3e, 0x4b, 0x99, 0x86, 0x60, 0x8a,
	0x59, 0x5a, 0x1c, 0xb7, 0x7e, 0x81, 0x0f, 0x0c, 0xc9, 0x67, 0x1a, 0x4e, 0x30, 0x4b, 0x69, 0x48,
	0x1e, 0xdf, 0xa6, 0x07, 0x99, 0xe2, 0x15, 0x46, 0x24, 0x3f, 0xa1, 0x7b, 0xf3, 0x84, 0xf7, 0x8a,
	0x0f, 0x86, 0x97, 0x33, 0xc7, 0xba, 0x9a, 0x39, 0xd6, 0xdf, 0x99, 0x63, 0xfd, 0x9c, 0x3b, 0xb5,
	0xab, 0xb9, 0x53, 0xfb, 0x3d, 0x77, 0x6a, 0x9f, 0xbc, 0x15, 0x6c, 0x7e, 0x17, 0x87, 0x31, 0xe8,
	0x6f, 0x98, 0x9e, 0x15, 0x0b, 0xef, 0xfb, 0xad, 0x97, 0x6b, 0xbc, 0x59, 0xbc, 0x39, 0x2f, 0xfe,
	0x05, 0x00, 0x00, 0xff, 0xff, 0x0c, 0x55, 0xcc, 0xb2, 0xd9, 0x04, 0x00, 0x00,
}

func (m *PositionRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PositionRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PositionRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *Position) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Position) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Position) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.StopLossPrice.Size()
		i -= size
		if _, err := m.StopLossPrice.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x52
	if m.AmmPoolId != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.AmmPoolId))
		i--
		dAtA[i] = 0x48
	}
	if m.Id != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x40
	}
	{
		size := m.PositionHealth.Size()
		i -= size
		if _, err := m.PositionHealth.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.LeveragedLpAmount.Size()
		i -= size
		if _, err := m.LeveragedLpAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.Leverage.Size()
		i -= size
		if _, err := m.Leverage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.InterestPaid.Size()
		i -= size
		if _, err := m.InterestPaid.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.Liabilities.Size()
		i -= size
		if _, err := m.Liabilities.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.Collateral.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *AddressId) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AddressId) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AddressId) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintTypes(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *PositionAndInterest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PositionAndInterest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PositionAndInterest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.InterestRateHourUsd.Size()
		i -= size
		if _, err := m.InterestRateHourUsd.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.InterestRateHour.Size()
		i -= size
		if _, err := m.InterestRateHour.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.Position != nil {
		{
			size, err := m.Position.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTypes(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
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
func (m *PositionRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.Id != 0 {
		n += 1 + sovTypes(uint64(m.Id))
	}
	return n
}

func (m *Position) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.Collateral.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.Liabilities.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.InterestPaid.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.Leverage.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.LeveragedLpAmount.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.PositionHealth.Size()
	n += 1 + l + sovTypes(uint64(l))
	if m.Id != 0 {
		n += 1 + sovTypes(uint64(m.Id))
	}
	if m.AmmPoolId != 0 {
		n += 1 + sovTypes(uint64(m.AmmPoolId))
	}
	l = m.StopLossPrice.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func (m *AddressId) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if m.Id != 0 {
		n += 1 + sovTypes(uint64(m.Id))
	}
	return n
}

func (m *PositionAndInterest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Position != nil {
		l = m.Position.Size()
		n += 1 + l + sovTypes(uint64(l))
	}
	l = m.InterestRateHour.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.InterestRateHourUsd.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PositionRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: PositionRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PositionRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
func (m *Position) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: Position: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Position: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Collateral", wireType)
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
			if err := m.Collateral.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liabilities", wireType)
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
			if err := m.Liabilities.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestPaid", wireType)
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
			if err := m.InterestPaid.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Leverage", wireType)
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
			if err := m.Leverage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LeveragedLpAmount", wireType)
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
			if err := m.LeveragedLpAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PositionHealth", wireType)
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
			if err := m.PositionHealth.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AmmPoolId", wireType)
			}
			m.AmmPoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StopLossPrice", wireType)
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
			if err := m.StopLossPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *AddressId) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: AddressId: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AddressId: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
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
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
func (m *PositionAndInterest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: PositionAndInterest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PositionAndInterest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Position", wireType)
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
			if m.Position == nil {
				m.Position = &Position{}
			}
			if err := m.Position.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRateHour", wireType)
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
			if err := m.InterestRateHour.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRateHourUsd", wireType)
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
			if err := m.InterestRateHourUsd.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
