// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/margin/types.proto

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

type Position int32

const (
	Position_UNSPECIFIED Position = 0
	Position_LONG        Position = 1
	Position_SHORT       Position = 2
)

var Position_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "LONG",
	2: "SHORT",
}

var Position_value = map[string]int32{
	"UNSPECIFIED": 0,
	"LONG":        1,
	"SHORT":       2,
}

func (x Position) String() string {
	return proto.EnumName(Position_name, int32(x))
}

func (Position) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cd1c09c977f732f9, []int{0}
}

type MTP struct {
	Address                   string                                   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	CollateralAssets          []string                                 `protobuf:"bytes,2,rep,name=collateral_assets,json=collateralAssets,proto3" json:"collateral_assets,omitempty"`
	CollateralAmounts         []github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,rep,name=collateral_amounts,json=collateralAmounts,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"collateral_amounts"`
	Liabilities               github_com_cosmos_cosmos_sdk_types.Int   `protobuf:"bytes,4,opt,name=liabilities,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"liabilities"`
	InterestPaidCollaterals   []github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,5,rep,name=interest_paid_collaterals,json=interestPaidCollaterals,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"interest_paid_collaterals"`
	InterestPaidCustodys      []github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,rep,name=interest_paid_custodys,json=interestPaidCustodys,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"interest_paid_custodys"`
	InterestUnpaidCollaterals []github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,7,rep,name=interest_unpaid_collaterals,json=interestUnpaidCollaterals,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"interest_unpaid_collaterals"`
	CustodyAssets             []string                                 `protobuf:"bytes,8,rep,name=custody_assets,json=custodyAssets,proto3" json:"custody_assets,omitempty"`
	CustodyAmounts            []github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,9,rep,name=custody_amounts,json=custodyAmounts,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"custody_amounts"`
	Leverages                 []github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,10,rep,name=leverages,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"leverages"`
	MtpHealth                 github_com_cosmos_cosmos_sdk_types.Dec   `protobuf:"bytes,11,opt,name=mtp_health,json=mtpHealth,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"mtp_health"`
	Position                  Position                                 `protobuf:"varint,12,opt,name=position,proto3,enum=elys.margin.Position" json:"position,omitempty"`
	Id                        uint64                                   `protobuf:"varint,13,opt,name=id,proto3" json:"id,omitempty"`
	AmmPoolId                 uint64                                   `protobuf:"varint,14,opt,name=amm_pool_id,json=ammPoolId,proto3" json:"amm_pool_id,omitempty"`
	ConsolidateLeverage       github_com_cosmos_cosmos_sdk_types.Dec   `protobuf:"bytes,15,opt,name=consolidate_leverage,json=consolidateLeverage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"consolidate_leverage"`
	SumCollateral             github_com_cosmos_cosmos_sdk_types.Int   `protobuf:"bytes,16,opt,name=sum_collateral,json=sumCollateral,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"sum_collateral"`
}

func (m *MTP) Reset()         { *m = MTP{} }
func (m *MTP) String() string { return proto.CompactTextString(m) }
func (*MTP) ProtoMessage()    {}
func (*MTP) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd1c09c977f732f9, []int{0}
}
func (m *MTP) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MTP) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MTP.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MTP) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MTP.Merge(m, src)
}
func (m *MTP) XXX_Size() int {
	return m.Size()
}
func (m *MTP) XXX_DiscardUnknown() {
	xxx_messageInfo_MTP.DiscardUnknown(m)
}

var xxx_messageInfo_MTP proto.InternalMessageInfo

func (m *MTP) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *MTP) GetCollateralAssets() []string {
	if m != nil {
		return m.CollateralAssets
	}
	return nil
}

func (m *MTP) GetCustodyAssets() []string {
	if m != nil {
		return m.CustodyAssets
	}
	return nil
}

func (m *MTP) GetPosition() Position {
	if m != nil {
		return m.Position
	}
	return Position_UNSPECIFIED
}

func (m *MTP) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *MTP) GetAmmPoolId() uint64 {
	if m != nil {
		return m.AmmPoolId
	}
	return 0
}

type WhiteList struct {
	ValidatorList []string `protobuf:"bytes,1,rep,name=validator_list,json=validatorList,proto3" json:"validator_list,omitempty"`
}

func (m *WhiteList) Reset()         { *m = WhiteList{} }
func (m *WhiteList) String() string { return proto.CompactTextString(m) }
func (*WhiteList) ProtoMessage()    {}
func (*WhiteList) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd1c09c977f732f9, []int{1}
}
func (m *WhiteList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WhiteList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WhiteList.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WhiteList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WhiteList.Merge(m, src)
}
func (m *WhiteList) XXX_Size() int {
	return m.Size()
}
func (m *WhiteList) XXX_DiscardUnknown() {
	xxx_messageInfo_WhiteList.DiscardUnknown(m)
}

var xxx_messageInfo_WhiteList proto.InternalMessageInfo

func (m *WhiteList) GetValidatorList() []string {
	if m != nil {
		return m.ValidatorList
	}
	return nil
}

func init() {
	proto.RegisterEnum("elys.margin.Position", Position_name, Position_value)
	proto.RegisterType((*MTP)(nil), "elys.margin.MTP")
	proto.RegisterType((*WhiteList)(nil), "elys.margin.WhiteList")
}

func init() { proto.RegisterFile("elys/margin/types.proto", fileDescriptor_cd1c09c977f732f9) }

var fileDescriptor_cd1c09c977f732f9 = []byte{
	// 591 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xdf, 0x4e, 0xdb, 0x3e,
	0x14, 0xc7, 0x9b, 0xf2, 0xaf, 0x39, 0xfd, 0x51, 0xfa, 0xf3, 0xd8, 0xf0, 0x36, 0x29, 0x54, 0x48,
	0x9b, 0xaa, 0x21, 0xda, 0x8d, 0x3d, 0x01, 0xff, 0x36, 0x2a, 0x15, 0x88, 0x02, 0x08, 0x69, 0xd2,
	0x14, 0x99, 0xda, 0x6a, 0x3d, 0xec, 0x38, 0x8a, 0x1d, 0x36, 0xde, 0x62, 0x8f, 0xc5, 0x25, 0x97,
	0xd3, 0x2e, 0xd0, 0x04, 0x4f, 0xb1, 0xbb, 0x29, 0x69, 0xd2, 0x86, 0xed, 0x8a, 0x5c, 0x25, 0x3e,
	0xe7, 0xe8, 0xf3, 0x3d, 0xc7, 0xfe, 0xda, 0xb0, 0xc2, 0xc4, 0x95, 0xee, 0x4a, 0x12, 0x0d, 0x79,
	0xd0, 0x35, 0x57, 0x21, 0xd3, 0x9d, 0x30, 0x52, 0x46, 0xa1, 0x7a, 0x92, 0xe8, 0x8c, 0x13, 0x2f,
	0x96, 0x87, 0x6a, 0xa8, 0xd2, 0x78, 0x37, 0xf9, 0x1b, 0x97, 0xac, 0xfd, 0xae, 0xc1, 0xcc, 0xc1,
	0x89, 0x8b, 0x30, 0x2c, 0x10, 0x4a, 0x23, 0xa6, 0x35, 0xb6, 0x5a, 0x56, 0xdb, 0xf6, 0xf2, 0x25,
	0x5a, 0x87, 0xff, 0x07, 0x4a, 0x08, 0x62, 0x58, 0x44, 0x84, 0x4f, 0xb4, 0x66, 0x46, 0xe3, 0x6a,
	0x6b, 0xa6, 0x6d, 0x7b, 0xcd, 0x69, 0x62, 0x2b, 0x8d, 0xa3, 0xcf, 0x80, 0x8a, 0xc5, 0x52, 0xc5,
	0x81, 0xd1, 0x78, 0x26, 0xa9, 0xde, 0xee, 0x5c, 0xdf, 0xae, 0x56, 0x7e, 0xde, 0xae, 0xbe, 0x1e,
	0x72, 0x33, 0x8a, 0xcf, 0x3b, 0x03, 0x25, 0xbb, 0x03, 0xa5, 0xa5, 0xd2, 0xd9, 0x67, 0x43, 0xd3,
	0x8b, 0xac, 0xff, 0x5e, 0x60, 0xbc, 0x82, 0xec, 0xd6, 0x18, 0x84, 0x5c, 0xa8, 0x0b, 0x4e, 0xce,
	0xb9, 0xe0, 0x86, 0x33, 0x8d, 0x67, 0x93, 0x4e, 0x1f, 0xcd, 0x2d, 0x22, 0xd0, 0x17, 0x78, 0xce,
	0x03, 0xc3, 0x22, 0xa6, 0x8d, 0x1f, 0x12, 0x4e, 0xfd, 0xa9, 0xa8, 0xc6, 0x73, 0xa5, 0xfa, 0x5e,
	0xc9, 0x81, 0x2e, 0xe1, 0x74, 0x67, 0x8a, 0x43, 0x14, 0x9e, 0xfd, 0xa5, 0x15, 0x6b, 0xa3, 0xe8,
	0x95, 0xc6, 0xf3, 0xa5, 0x84, 0x96, 0x1f, 0x08, 0x65, 0x2c, 0x14, 0xc0, 0xcb, 0x89, 0x4a, 0x1c,
	0xfc, 0x33, 0xd3, 0x42, 0x29, 0xa9, 0xc9, 0x26, 0x9d, 0xa6, 0xc4, 0xe2, 0x54, 0xaf, 0xa0, 0x91,
	0xcd, 0x91, 0x9b, 0xa3, 0x96, 0x9a, 0x63, 0x31, 0x8b, 0x66, 0xce, 0x38, 0x83, 0xa5, 0x49, 0x59,
	0x66, 0x0b, 0xbb, 0x54, 0x2b, 0xb9, 0x5a, 0xee, 0x89, 0x3e, 0xd8, 0x82, 0x5d, 0xb2, 0x88, 0x0c,
	0x99, 0xc6, 0xf0, 0x68, 0xe4, 0x2e, 0x1b, 0x78, 0x53, 0x00, 0x3a, 0x00, 0x90, 0x26, 0xf4, 0x47,
	0x8c, 0x08, 0x33, 0xc2, 0xf5, 0x47, 0x1b, 0x2c, 0xc5, 0x49, 0x13, 0xee, 0xa7, 0x00, 0xf4, 0x0e,
	0x6a, 0xa1, 0xd2, 0xdc, 0x70, 0x15, 0xe0, 0xff, 0x5a, 0x56, 0xbb, 0xb1, 0xf9, 0xb4, 0x53, 0xb8,
	0x94, 0x1d, 0x37, 0x4b, 0x7a, 0x93, 0x32, 0xd4, 0x80, 0x2a, 0xa7, 0x78, 0xb1, 0x65, 0xb5, 0x67,
	0xbd, 0x2a, 0xa7, 0xc8, 0x81, 0x3a, 0x91, 0xd2, 0x0f, 0x95, 0x12, 0x3e, 0xa7, 0xb8, 0x91, 0x26,
	0x6c, 0x22, 0xa5, 0xab, 0x94, 0xe8, 0x51, 0x44, 0x60, 0x79, 0xa0, 0x02, 0xad, 0x04, 0xa7, 0xc4,
	0x30, 0x3f, 0x1f, 0x05, 0x2f, 0x95, 0xea, 0xfd, 0x49, 0x81, 0xd5, 0xcf, 0x50, 0xe8, 0x14, 0x1a,
	0x3a, 0x96, 0x05, 0x1b, 0xe1, 0x66, 0xa9, 0x9b, 0xb7, 0xa8, 0x63, 0x39, 0xb5, 0xce, 0xda, 0x26,
	0xd8, 0x67, 0x23, 0x6e, 0x58, 0x9f, 0x6b, 0x93, 0xd8, 0xe8, 0x92, 0xa4, 0xba, 0x2a, 0xf2, 0x05,
	0xd7, 0x06, 0x5b, 0x63, 0x1b, 0x4d, 0xa2, 0x49, 0xd9, 0x9b, 0xb7, 0x50, 0xcb, 0xf7, 0x0c, 0x2d,
	0x41, 0xfd, 0xf4, 0xf0, 0xd8, 0xdd, 0xdb, 0xe9, 0x7d, 0xe8, 0xed, 0xed, 0x36, 0x2b, 0xa8, 0x06,
	0xb3, 0xfd, 0xa3, 0xc3, 0x8f, 0x4d, 0x0b, 0xd9, 0x30, 0x77, 0xbc, 0x7f, 0xe4, 0x9d, 0x34, 0xab,
	0xdb, 0x7b, 0xd7, 0x77, 0x8e, 0x75, 0x73, 0xe7, 0x58, 0xbf, 0xee, 0x1c, 0xeb, 0xfb, 0xbd, 0x53,
	0xb9, 0xb9, 0x77, 0x2a, 0x3f, 0xee, 0x9d, 0xca, 0xa7, 0xf5, 0x42, 0xdb, 0xc9, 0xa1, 0x6c, 0x04,
	0xcc, 0x7c, 0x55, 0xd1, 0x45, 0xba, 0xe8, 0x7e, 0x7b, 0xf0, 0xa2, 0x9e, 0xcf, 0xa7, 0xef, 0xe5,
	0xfb, 0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x49, 0x3f, 0x3a, 0xe7, 0x6d, 0x05, 0x00, 0x00,
}

func (m *MTP) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MTP) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MTP) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.SumCollateral.Size()
		i -= size
		if _, err := m.SumCollateral.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1
	i--
	dAtA[i] = 0x82
	{
		size := m.ConsolidateLeverage.Size()
		i -= size
		if _, err := m.ConsolidateLeverage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x7a
	if m.AmmPoolId != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.AmmPoolId))
		i--
		dAtA[i] = 0x70
	}
	if m.Id != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x68
	}
	if m.Position != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.Position))
		i--
		dAtA[i] = 0x60
	}
	{
		size := m.MtpHealth.Size()
		i -= size
		if _, err := m.MtpHealth.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x5a
	if len(m.Leverages) > 0 {
		for iNdEx := len(m.Leverages) - 1; iNdEx >= 0; iNdEx-- {
			{
				size := m.Leverages[iNdEx].Size()
				i -= size
				if _, err := m.Leverages[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x52
		}
	}
	if len(m.CustodyAmounts) > 0 {
		for iNdEx := len(m.CustodyAmounts) - 1; iNdEx >= 0; iNdEx-- {
			{
				size := m.CustodyAmounts[iNdEx].Size()
				i -= size
				if _, err := m.CustodyAmounts[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x4a
		}
	}
	if len(m.CustodyAssets) > 0 {
		for iNdEx := len(m.CustodyAssets) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.CustodyAssets[iNdEx])
			copy(dAtA[i:], m.CustodyAssets[iNdEx])
			i = encodeVarintTypes(dAtA, i, uint64(len(m.CustodyAssets[iNdEx])))
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.InterestUnpaidCollaterals) > 0 {
		for iNdEx := len(m.InterestUnpaidCollaterals) - 1; iNdEx >= 0; iNdEx-- {
			{
				size := m.InterestUnpaidCollaterals[iNdEx].Size()
				i -= size
				if _, err := m.InterestUnpaidCollaterals[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.InterestPaidCustodys) > 0 {
		for iNdEx := len(m.InterestPaidCustodys) - 1; iNdEx >= 0; iNdEx-- {
			{
				size := m.InterestPaidCustodys[iNdEx].Size()
				i -= size
				if _, err := m.InterestPaidCustodys[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.InterestPaidCollaterals) > 0 {
		for iNdEx := len(m.InterestPaidCollaterals) - 1; iNdEx >= 0; iNdEx-- {
			{
				size := m.InterestPaidCollaterals[iNdEx].Size()
				i -= size
				if _, err := m.InterestPaidCollaterals[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	{
		size := m.Liabilities.Size()
		i -= size
		if _, err := m.Liabilities.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.CollateralAmounts) > 0 {
		for iNdEx := len(m.CollateralAmounts) - 1; iNdEx >= 0; iNdEx-- {
			{
				size := m.CollateralAmounts[iNdEx].Size()
				i -= size
				if _, err := m.CollateralAmounts[iNdEx].MarshalTo(dAtA[i:]); err != nil {
					return 0, err
				}
				i = encodeVarintTypes(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.CollateralAssets) > 0 {
		for iNdEx := len(m.CollateralAssets) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.CollateralAssets[iNdEx])
			copy(dAtA[i:], m.CollateralAssets[iNdEx])
			i = encodeVarintTypes(dAtA, i, uint64(len(m.CollateralAssets[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
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

func (m *WhiteList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WhiteList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WhiteList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ValidatorList) > 0 {
		for iNdEx := len(m.ValidatorList) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ValidatorList[iNdEx])
			copy(dAtA[i:], m.ValidatorList[iNdEx])
			i = encodeVarintTypes(dAtA, i, uint64(len(m.ValidatorList[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
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
func (m *MTP) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovTypes(uint64(l))
	}
	if len(m.CollateralAssets) > 0 {
		for _, s := range m.CollateralAssets {
			l = len(s)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if len(m.CollateralAmounts) > 0 {
		for _, e := range m.CollateralAmounts {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	l = m.Liabilities.Size()
	n += 1 + l + sovTypes(uint64(l))
	if len(m.InterestPaidCollaterals) > 0 {
		for _, e := range m.InterestPaidCollaterals {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if len(m.InterestPaidCustodys) > 0 {
		for _, e := range m.InterestPaidCustodys {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if len(m.InterestUnpaidCollaterals) > 0 {
		for _, e := range m.InterestUnpaidCollaterals {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if len(m.CustodyAssets) > 0 {
		for _, s := range m.CustodyAssets {
			l = len(s)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if len(m.CustodyAmounts) > 0 {
		for _, e := range m.CustodyAmounts {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	if len(m.Leverages) > 0 {
		for _, e := range m.Leverages {
			l = e.Size()
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	l = m.MtpHealth.Size()
	n += 1 + l + sovTypes(uint64(l))
	if m.Position != 0 {
		n += 1 + sovTypes(uint64(m.Position))
	}
	if m.Id != 0 {
		n += 1 + sovTypes(uint64(m.Id))
	}
	if m.AmmPoolId != 0 {
		n += 1 + sovTypes(uint64(m.AmmPoolId))
	}
	l = m.ConsolidateLeverage.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.SumCollateral.Size()
	n += 2 + l + sovTypes(uint64(l))
	return n
}

func (m *WhiteList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ValidatorList) > 0 {
		for _, s := range m.ValidatorList {
			l = len(s)
			n += 1 + l + sovTypes(uint64(l))
		}
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MTP) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MTP: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MTP: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralAssets", wireType)
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
			m.CollateralAssets = append(m.CollateralAssets, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CollateralAmounts", wireType)
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
			var v github_com_cosmos_cosmos_sdk_types.Int
			m.CollateralAmounts = append(m.CollateralAmounts, v)
			if err := m.CollateralAmounts[len(m.CollateralAmounts)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
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
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestPaidCollaterals", wireType)
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
			var v github_com_cosmos_cosmos_sdk_types.Int
			m.InterestPaidCollaterals = append(m.InterestPaidCollaterals, v)
			if err := m.InterestPaidCollaterals[len(m.InterestPaidCollaterals)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestPaidCustodys", wireType)
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
			var v github_com_cosmos_cosmos_sdk_types.Int
			m.InterestPaidCustodys = append(m.InterestPaidCustodys, v)
			if err := m.InterestPaidCustodys[len(m.InterestPaidCustodys)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestUnpaidCollaterals", wireType)
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
			var v github_com_cosmos_cosmos_sdk_types.Int
			m.InterestUnpaidCollaterals = append(m.InterestUnpaidCollaterals, v)
			if err := m.InterestUnpaidCollaterals[len(m.InterestUnpaidCollaterals)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CustodyAssets", wireType)
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
			m.CustodyAssets = append(m.CustodyAssets, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field CustodyAmounts", wireType)
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
			var v github_com_cosmos_cosmos_sdk_types.Int
			m.CustodyAmounts = append(m.CustodyAmounts, v)
			if err := m.CustodyAmounts[len(m.CustodyAmounts)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Leverages", wireType)
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
			var v github_com_cosmos_cosmos_sdk_types.Dec
			m.Leverages = append(m.Leverages, v)
			if err := m.Leverages[len(m.Leverages)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MtpHealth", wireType)
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
			if err := m.MtpHealth.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Position", wireType)
			}
			m.Position = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Position |= Position(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 13:
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
		case 14:
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
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConsolidateLeverage", wireType)
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
			if err := m.ConsolidateLeverage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 16:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SumCollateral", wireType)
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
			if err := m.SumCollateral.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *WhiteList) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: WhiteList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WhiteList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorList", wireType)
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
			m.ValidatorList = append(m.ValidatorList, string(dAtA[iNdEx:postIndex]))
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
