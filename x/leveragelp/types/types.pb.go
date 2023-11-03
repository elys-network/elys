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
	TakeProfitPrice   github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,10,opt,name=take_profit_price,json=takeProfitPrice,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"take_profit_price"`
}

func (m *Position) Reset()         { *m = Position{} }
func (m *Position) String() string { return proto.CompactTextString(m) }
func (*Position) ProtoMessage()    {}
func (*Position) Descriptor() ([]byte, []int) {
	return fileDescriptor_992d513dd201f55b, []int{0}
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

type WhiteList struct {
	ValidatorList []string `protobuf:"bytes,1,rep,name=validator_list,json=validatorList,proto3" json:"validator_list,omitempty"`
}

func (m *WhiteList) Reset()         { *m = WhiteList{} }
func (m *WhiteList) String() string { return proto.CompactTextString(m) }
func (*WhiteList) ProtoMessage()    {}
func (*WhiteList) Descriptor() ([]byte, []int) {
	return fileDescriptor_992d513dd201f55b, []int{1}
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
	proto.RegisterType((*Position)(nil), "elys.leveragelp.Position")
	proto.RegisterType((*WhiteList)(nil), "elys.leveragelp.WhiteList")
}

func init() { proto.RegisterFile("elys/leveragelp/types.proto", fileDescriptor_992d513dd201f55b) }

var fileDescriptor_992d513dd201f55b = []byte{
	// 472 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xe3, 0x34, 0x6d, 0x93, 0x0d, 0x4d, 0xd4, 0x85, 0xc3, 0x52, 0x24, 0x37, 0xaa, 0x04,
	0xca, 0xa5, 0xb6, 0x5a, 0x1e, 0x00, 0x11, 0x38, 0x10, 0xd4, 0x83, 0x65, 0x0e, 0x95, 0x7a, 0xc0,
	0xda, 0x78, 0x97, 0x64, 0x94, 0xb5, 0x67, 0xe5, 0x9d, 0x06, 0xfa, 0x16, 0xbc, 0x00, 0xef, 0xd3,
	0x63, 0x8f, 0x88, 0x43, 0x85, 0x92, 0x17, 0x41, 0x76, 0xfe, 0x10, 0x89, 0x13, 0x39, 0x79, 0xc6,
	0xdf, 0xcc, 0x4f, 0xdf, 0xee, 0xea, 0x63, 0x2f, 0xb4, 0xb9, 0x73, 0xa1, 0xd1, 0x33, 0x5d, 0xc8,
	0xb1, 0x36, 0x36, 0xa4, 0x3b, 0xab, 0x5d, 0x60, 0x0b, 0x24, 0xe4, 0xdd, 0x52, 0x0c, 0xfe, 0x8a,
	0x27, 0xcf, 0xc6, 0x38, 0xc6, 0x4a, 0x0b, 0xcb, 0x6a, 0x39, 0x76, 0xe2, 0xa7, 0xe8, 0x32, 0x74,
	0xe1, 0x48, 0x3a, 0x1d, 0xce, 0x2e, 0x46, 0x9a, 0xe4, 0x45, 0x98, 0x22, 0xe4, 0x4b, 0xfd, 0xec,
	0xc7, 0x3e, 0x6b, 0x46, 0xe8, 0x80, 0x00, 0x73, 0x2e, 0xd8, 0xa1, 0x54, 0xaa, 0xd0, 0xce, 0x09,
	0xaf, 0xe7, 0xf5, 0x5b, 0xf1, 0xba, 0xe5, 0x6f, 0x18, 0x4b, 0xd1, 0x18, 0x49, 0xba, 0x90, 0x46,
	0xd4, 0x7b, 0x5e, 0xbf, 0x7d, 0xf9, 0x3c, 0x58, 0xb2, 0x83, 0x92, 0x1d, 0xac, 0xd8, 0xc1, 0x3b,
	0x84, 0x7c, 0xd0, 0xb8, 0x7f, 0x3c, 0xad, 0xc5, 0x5b, 0x2b, 0x3c, 0x62, 0x6d, 0x03, 0x72, 0x04,
	0x06, 0x08, 0xb4, 0x13, 0x7b, 0x25, 0x7e, 0x10, 0x94, 0x63, 0xbf, 0x1e, 0x4f, 0x5f, 0x8d, 0x81,
	0x26, 0xb7, 0xa3, 0x20, 0xc5, 0x2c, 0x5c, 0xf9, 0x5d, 0x7e, 0xce, 0x9d, 0x9a, 0xae, 0x4e, 0x3d,
	0xcc, 0x29, 0xde, 0x46, 0xf0, 0x4f, 0xec, 0x08, 0x72, 0xd2, 0x85, 0x76, 0x94, 0x58, 0x09, 0x4a,
	0x34, 0x76, 0x62, 0x3e, 0x59, 0x43, 0x22, 0x09, 0x8a, 0x7f, 0x64, 0xcd, 0xf5, 0x95, 0x8a, 0xfd,
	0xff, 0xe6, 0xbd, 0xd7, 0x69, 0xbc, 0xd9, 0xe7, 0x9f, 0xd9, 0xd3, 0x75, 0xad, 0x12, 0x63, 0x13,
	0x99, 0xe1, 0x6d, 0x4e, 0xe2, 0x60, 0x27, 0x9b, 0xc7, 0x1b, 0xd4, 0x95, 0x7d, 0x5b, 0x81, 0xf8,
	0x35, 0xeb, 0xda, 0xd5, 0xcb, 0x25, 0x13, 0x2d, 0x0d, 0x4d, 0xc4, 0xe1, 0x4e, 0x96, 0x3b, 0x6b,
	0xcc, 0x87, 0x8a, 0xc2, 0x3b, 0xac, 0x0e, 0x4a, 0x34, 0x7b, 0x5e, 0xbf, 0x11, 0xd7, 0x41, 0x71,
	0x9f, 0xb5, 0x65, 0x96, 0x25, 0x16, 0xd1, 0x24, 0xa0, 0x44, 0xab, 0x12, 0x5a, 0x32, 0xcb, 0x22,
	0x44, 0x33, 0x54, 0xfc, 0x86, 0x1d, 0x93, 0x9c, 0xea, 0xc4, 0x16, 0xf8, 0x05, 0x28, 0xb1, 0x05,
	0xa4, 0x5a, 0xb0, 0x9d, 0xac, 0x74, 0x4b, 0x50, 0x54, 0x71, 0xa2, 0x12, 0x73, 0x76, 0xc9, 0x5a,
	0xd7, 0x13, 0x20, 0x7d, 0x05, 0x8e, 0xf8, 0x4b, 0xd6, 0x99, 0x49, 0x03, 0x4a, 0x12, 0x16, 0x89,
	0x01, 0x47, 0xc2, 0xeb, 0xed, 0xf5, 0x5b, 0xf1, 0xd1, 0xe6, 0x6f, 0x39, 0x36, 0x18, 0xde, 0xcf,
	0x7d, 0xef, 0x61, 0xee, 0x7b, 0xbf, 0xe7, 0xbe, 0xf7, 0x7d, 0xe1, 0xd7, 0x1e, 0x16, 0x7e, 0xed,
	0xe7, 0xc2, 0xaf, 0xdd, 0x84, 0x5b, 0x36, 0xca, 0xfc, 0x9c, 0xe7, 0x9a, 0xbe, 0x62, 0x31, 0xad,
	0x9a, 0xf0, 0xdb, 0x3f, 0x59, 0x1b, 0x1d, 0x54, 0x29, 0x79, 0xfd, 0x27, 0x00, 0x00, 0xff, 0xff,
	0x45, 0x09, 0x38, 0x5d, 0x8b, 0x03, 0x00, 0x00,
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
		size := m.TakeProfitPrice.Size()
		i -= size
		if _, err := m.TakeProfitPrice.MarshalTo(dAtA[i:]); err != nil {
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
	l = m.TakeProfitPrice.Size()
	n += 1 + l + sovTypes(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field TakeProfitPrice", wireType)
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
			if err := m.TakeProfitPrice.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
