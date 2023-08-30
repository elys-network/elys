// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/margin/pool.proto

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

type PoolAsset struct {
	Liabilities          github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,1,opt,name=liabilities,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"liabilities"`
	Custody              github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=custody,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"custody"`
	AssetBalance         github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,opt,name=asset_balance,json=assetBalance,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"asset_balance" yaml:"asset_balance"`
	UnsettledLiabilities github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,4,opt,name=unsettled_liabilities,json=unsettledLiabilities,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"unsettled_liabilities"`
	BlockInterest        github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,5,opt,name=block_interest,json=blockInterest,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"block_interest"`
	AssetDenom           string                                 `protobuf:"bytes,6,opt,name=asset_denom,json=assetDenom,proto3" json:"asset_denom,omitempty"`
}

func (m *PoolAsset) Reset()         { *m = PoolAsset{} }
func (m *PoolAsset) String() string { return proto.CompactTextString(m) }
func (*PoolAsset) ProtoMessage()    {}
func (*PoolAsset) Descriptor() ([]byte, []int) {
	return fileDescriptor_030dd771f91c2bb4, []int{0}
}
func (m *PoolAsset) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PoolAsset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PoolAsset.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PoolAsset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolAsset.Merge(m, src)
}
func (m *PoolAsset) XXX_Size() int {
	return m.Size()
}
func (m *PoolAsset) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolAsset.DiscardUnknown(m)
}

var xxx_messageInfo_PoolAsset proto.InternalMessageInfo

func (m *PoolAsset) GetAssetDenom() string {
	if m != nil {
		return m.AssetDenom
	}
	return ""
}

type Pool struct {
	AmmPoolId    uint64                                 `protobuf:"varint,1,opt,name=ammPoolId,proto3" json:"ammPoolId,omitempty"`
	Health       github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=health,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"health"`
	Enabled      bool                                   `protobuf:"varint,3,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Closed       bool                                   `protobuf:"varint,4,opt,name=closed,proto3" json:"closed,omitempty"`
	InterestRate github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=interest_rate,json=interestRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"interest_rate"`
	PoolAssets   []PoolAsset                            `protobuf:"bytes,6,rep,name=poolAssets,proto3" json:"poolAssets"`
}

func (m *Pool) Reset()         { *m = Pool{} }
func (m *Pool) String() string { return proto.CompactTextString(m) }
func (*Pool) ProtoMessage()    {}
func (*Pool) Descriptor() ([]byte, []int) {
	return fileDescriptor_030dd771f91c2bb4, []int{1}
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

func (m *Pool) GetAmmPoolId() uint64 {
	if m != nil {
		return m.AmmPoolId
	}
	return 0
}

func (m *Pool) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *Pool) GetClosed() bool {
	if m != nil {
		return m.Closed
	}
	return false
}

func (m *Pool) GetPoolAssets() []PoolAsset {
	if m != nil {
		return m.PoolAssets
	}
	return nil
}

func init() {
	proto.RegisterType((*PoolAsset)(nil), "elys.margin.PoolAsset")
	proto.RegisterType((*Pool)(nil), "elys.margin.Pool")
}

func init() { proto.RegisterFile("elys/margin/pool.proto", fileDescriptor_030dd771f91c2bb4) }

var fileDescriptor_030dd771f91c2bb4 = []byte{
	// 451 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x93, 0x41, 0x8b, 0xd3, 0x40,
	0x14, 0xc7, 0x9b, 0x6d, 0xcd, 0x6e, 0x5f, 0xb7, 0x1e, 0x86, 0x5a, 0x82, 0x48, 0xba, 0xf4, 0x20,
	0x0b, 0xb2, 0x09, 0xe8, 0x4d, 0xbc, 0x58, 0xd6, 0xc5, 0x82, 0x87, 0x65, 0xc4, 0x8b, 0x97, 0x32,
	0x49, 0x1e, 0x6d, 0xe8, 0x24, 0xaf, 0x64, 0xa6, 0x68, 0xbf, 0x85, 0x5f, 0xc5, 0x6f, 0xb1, 0xc7,
	0x3d, 0x8a, 0x87, 0x22, 0xed, 0x27, 0xd0, 0x4f, 0x20, 0x33, 0x4d, 0x6a, 0x3c, 0x6e, 0x4e, 0x99,
	0xff, 0x9f, 0x99, 0x5f, 0xde, 0xfc, 0xdf, 0x3c, 0x18, 0xa2, 0xdc, 0xa8, 0x30, 0x13, 0xc5, 0x3c,
	0xcd, 0xc3, 0x15, 0x91, 0x0c, 0x56, 0x05, 0x69, 0x62, 0x3d, 0xe3, 0x07, 0x07, 0xff, 0xe9, 0x60,
	0x4e, 0x73, 0xb2, 0x7e, 0x68, 0x56, 0x87, 0x2d, 0xe3, 0xdf, 0x6d, 0xe8, 0xde, 0x12, 0xc9, 0xb7,
	0x4a, 0xa1, 0x66, 0xb7, 0xd0, 0x93, 0xa9, 0x88, 0x52, 0x99, 0xea, 0x14, 0x95, 0xe7, 0x5c, 0x38,
	0x97, 0xdd, 0x49, 0x70, 0xb7, 0x1d, 0xb5, 0x7e, 0x6e, 0x47, 0xcf, 0xe7, 0xa9, 0x5e, 0xac, 0xa3,
	0x20, 0xa6, 0x2c, 0x8c, 0x49, 0x65, 0xa4, 0xca, 0xcf, 0x95, 0x4a, 0x96, 0xa1, 0xde, 0xac, 0x50,
	0x05, 0xd3, 0x5c, 0xf3, 0x3a, 0x82, 0xbd, 0x87, 0xd3, 0x78, 0xad, 0x34, 0x25, 0x1b, 0xef, 0xa4,
	0x11, 0xad, 0x3a, 0xce, 0x96, 0xd0, 0x17, 0xa6, 0xc8, 0x59, 0x24, 0xa4, 0xc8, 0x63, 0xf4, 0xda,
	0x96, 0x77, 0xf3, 0x30, 0xde, 0x9f, 0xed, 0x68, 0xb0, 0x11, 0x99, 0x7c, 0x3d, 0xfe, 0x0f, 0x36,
	0xe6, 0xe7, 0x56, 0x4f, 0x0e, 0x92, 0xc5, 0xf0, 0x64, 0x9d, 0x2b, 0xd4, 0x5a, 0x62, 0x32, 0xab,
	0x47, 0xd2, 0x69, 0x74, 0x89, 0xc1, 0x11, 0xf6, 0xa1, 0x96, 0xcd, 0x27, 0x78, 0x1c, 0x49, 0x8a,
	0x97, 0xb3, 0x34, 0xd7, 0x58, 0xa0, 0xd2, 0xde, 0xa3, 0x46, 0xf4, 0xbe, 0xa5, 0x4c, 0x4b, 0x08,
	0x1b, 0x41, 0xef, 0x70, 0xb7, 0x04, 0x73, 0xca, 0x3c, 0xd7, 0x30, 0x39, 0x58, 0xeb, 0xda, 0x38,
	0xe3, 0xef, 0x27, 0xd0, 0x31, 0x3d, 0x67, 0xcf, 0xa0, 0x2b, 0xb2, 0xcc, 0x2c, 0xa7, 0x89, 0x6d,
	0x76, 0x87, 0xff, 0x33, 0xd8, 0x0d, 0xb8, 0x0b, 0x14, 0x52, 0x2f, 0x1a, 0x74, 0xee, 0x1a, 0x63,
	0x5e, 0x9e, 0x66, 0x1e, 0x9c, 0x62, 0x2e, 0x22, 0x89, 0x89, 0x6d, 0xd9, 0x19, 0xaf, 0x24, 0x1b,
	0x82, 0x1b, 0x4b, 0x52, 0x98, 0xd8, 0x58, 0xcf, 0x78, 0xa9, 0xd8, 0x47, 0xe8, 0x57, 0x91, 0xcc,
	0x0a, 0xa1, 0xb1, 0x41, 0x2e, 0xa6, 0x80, 0xf3, 0x0a, 0xc2, 0x85, 0x46, 0xf6, 0x06, 0x60, 0x55,
	0x3d, 0x74, 0xe5, 0xb9, 0x17, 0xed, 0xcb, 0xde, 0xcb, 0x61, 0x50, 0x9b, 0x90, 0xe0, 0x38, 0x07,
	0x93, 0x8e, 0xf9, 0x13, 0xaf, 0xed, 0x9f, 0xbc, 0xbb, 0xdb, 0xf9, 0xce, 0xfd, 0xce, 0x77, 0x7e,
	0xed, 0x7c, 0xe7, 0xdb, 0xde, 0x6f, 0xdd, 0xef, 0xfd, 0xd6, 0x8f, 0xbd, 0xdf, 0xfa, 0xfc, 0xa2,
	0x56, 0x8d, 0xa1, 0x5d, 0xe5, 0xa8, 0xbf, 0x50, 0xb1, 0xb4, 0x22, 0xfc, 0x5a, 0x8d, 0xa5, 0x2d,
	0x2b, 0x72, 0xed, 0xd4, 0xbd, 0xfa, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x2b, 0xf1, 0x6f, 0x74, 0xb2,
	0x03, 0x00, 0x00,
}

func (m *PoolAsset) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PoolAsset) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PoolAsset) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AssetDenom) > 0 {
		i -= len(m.AssetDenom)
		copy(dAtA[i:], m.AssetDenom)
		i = encodeVarintPool(dAtA, i, uint64(len(m.AssetDenom)))
		i--
		dAtA[i] = 0x32
	}
	{
		size := m.BlockInterest.Size()
		i -= size
		if _, err := m.BlockInterest.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.UnsettledLiabilities.Size()
		i -= size
		if _, err := m.UnsettledLiabilities.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size := m.AssetBalance.Size()
		i -= size
		if _, err := m.AssetBalance.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Custody.Size()
		i -= size
		if _, err := m.Custody.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.Liabilities.Size()
		i -= size
		if _, err := m.Liabilities.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
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
			dAtA[i] = 0x32
		}
	}
	{
		size := m.InterestRate.Size()
		i -= size
		if _, err := m.InterestRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if m.Closed {
		i--
		if m.Closed {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x20
	}
	if m.Enabled {
		i--
		if m.Enabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	{
		size := m.Health.Size()
		i -= size
		if _, err := m.Health.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPool(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.AmmPoolId != 0 {
		i = encodeVarintPool(dAtA, i, uint64(m.AmmPoolId))
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
func (m *PoolAsset) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Liabilities.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.Custody.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.AssetBalance.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.UnsettledLiabilities.Size()
	n += 1 + l + sovPool(uint64(l))
	l = m.BlockInterest.Size()
	n += 1 + l + sovPool(uint64(l))
	l = len(m.AssetDenom)
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
	if m.AmmPoolId != 0 {
		n += 1 + sovPool(uint64(m.AmmPoolId))
	}
	l = m.Health.Size()
	n += 1 + l + sovPool(uint64(l))
	if m.Enabled {
		n += 2
	}
	if m.Closed {
		n += 2
	}
	l = m.InterestRate.Size()
	n += 1 + l + sovPool(uint64(l))
	if len(m.PoolAssets) > 0 {
		for _, e := range m.PoolAssets {
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
func (m *PoolAsset) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: PoolAsset: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PoolAsset: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liabilities", wireType)
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
			if err := m.Liabilities.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Custody", wireType)
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
			if err := m.Custody.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetBalance", wireType)
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
			if err := m.AssetBalance.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnsettledLiabilities", wireType)
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
			if err := m.UnsettledLiabilities.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockInterest", wireType)
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
			if err := m.BlockInterest.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetDenom", wireType)
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
			m.AssetDenom = string(dAtA[iNdEx:postIndex])
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
				return fmt.Errorf("proto: wrong wireType = %d for field AmmPoolId", wireType)
			}
			m.AmmPoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Health", wireType)
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
			if err := m.Health.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Enabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
			m.Enabled = bool(v != 0)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Closed", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPool
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
			m.Closed = bool(v != 0)
		case 5:
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
		case 6:
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
