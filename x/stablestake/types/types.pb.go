// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/stablestake/types.proto

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

type BalanceBorrowed struct {
	UsdAmount  cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=usd_amount,json=usdAmount,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"usd_amount"`
	Percentage cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=percentage,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"percentage"`
}

func (m *BalanceBorrowed) Reset()         { *m = BalanceBorrowed{} }
func (m *BalanceBorrowed) String() string { return proto.CompactTextString(m) }
func (*BalanceBorrowed) ProtoMessage()    {}
func (*BalanceBorrowed) Descriptor() ([]byte, []int) {
	return fileDescriptor_7179d85642fbc30a, []int{0}
}
func (m *BalanceBorrowed) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BalanceBorrowed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BalanceBorrowed.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BalanceBorrowed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BalanceBorrowed.Merge(m, src)
}
func (m *BalanceBorrowed) XXX_Size() int {
	return m.Size()
}
func (m *BalanceBorrowed) XXX_DiscardUnknown() {
	xxx_messageInfo_BalanceBorrowed.DiscardUnknown(m)
}

var xxx_messageInfo_BalanceBorrowed proto.InternalMessageInfo

type InterestBlock struct {
	InterestRate cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=interest_rate,json=interestRate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"interest_rate"`
	BlockTime    int64                       `protobuf:"varint,2,opt,name=block_time,json=blockTime,proto3" json:"block_time,omitempty"`
	BlockHeight  uint64                      `protobuf:"varint,3,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	PoolId       uint64                      `protobuf:"varint,4,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
}

func (m *InterestBlock) Reset()         { *m = InterestBlock{} }
func (m *InterestBlock) String() string { return proto.CompactTextString(m) }
func (*InterestBlock) ProtoMessage()    {}
func (*InterestBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_7179d85642fbc30a, []int{1}
}
func (m *InterestBlock) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InterestBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InterestBlock.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InterestBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InterestBlock.Merge(m, src)
}
func (m *InterestBlock) XXX_Size() int {
	return m.Size()
}
func (m *InterestBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_InterestBlock.DiscardUnknown(m)
}

var xxx_messageInfo_InterestBlock proto.InternalMessageInfo

func (m *InterestBlock) GetBlockTime() int64 {
	if m != nil {
		return m.BlockTime
	}
	return 0
}

func (m *InterestBlock) GetBlockHeight() uint64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *InterestBlock) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

type LegacyInterestBlock struct {
	InterestRate cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=interest_rate,json=interestRate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"interest_rate"`
	BlockTime    int64                       `protobuf:"varint,2,opt,name=block_time,json=blockTime,proto3" json:"block_time,omitempty"`
	BlockHeight  uint64                      `protobuf:"varint,3,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
}

func (m *LegacyInterestBlock) Reset()         { *m = LegacyInterestBlock{} }
func (m *LegacyInterestBlock) String() string { return proto.CompactTextString(m) }
func (*LegacyInterestBlock) ProtoMessage()    {}
func (*LegacyInterestBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_7179d85642fbc30a, []int{2}
}
func (m *LegacyInterestBlock) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LegacyInterestBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LegacyInterestBlock.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LegacyInterestBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LegacyInterestBlock.Merge(m, src)
}
func (m *LegacyInterestBlock) XXX_Size() int {
	return m.Size()
}
func (m *LegacyInterestBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_LegacyInterestBlock.DiscardUnknown(m)
}

var xxx_messageInfo_LegacyInterestBlock proto.InternalMessageInfo

func (m *LegacyInterestBlock) GetBlockTime() int64 {
	if m != nil {
		return m.BlockTime
	}
	return 0
}

func (m *LegacyInterestBlock) GetBlockHeight() uint64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func init() {
	proto.RegisterType((*BalanceBorrowed)(nil), "elys.stablestake.BalanceBorrowed")
	proto.RegisterType((*InterestBlock)(nil), "elys.stablestake.InterestBlock")
	proto.RegisterType((*LegacyInterestBlock)(nil), "elys.stablestake.LegacyInterestBlock")
}

func init() { proto.RegisterFile("elys/stablestake/types.proto", fileDescriptor_7179d85642fbc30a) }

var fileDescriptor_7179d85642fbc30a = []byte{
	// 372 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xd4, 0x92, 0x3d, 0x6a, 0xe3, 0x40,
	0x1c, 0xc5, 0x35, 0x6b, 0xe3, 0xc5, 0xb3, 0x36, 0xbb, 0x68, 0x17, 0x56, 0xeb, 0xdd, 0x95, 0xbd,
	0xae, 0xdc, 0x58, 0xda, 0x90, 0x13, 0x44, 0xb8, 0x88, 0x43, 0x8a, 0x44, 0x84, 0x14, 0x69, 0xc4,
	0x68, 0xf4, 0x47, 0x12, 0xfa, 0x18, 0xa1, 0x19, 0xe1, 0xf8, 0x16, 0xb9, 0x49, 0x9a, 0x90, 0x03,
	0xa4, 0x72, 0x69, 0x52, 0x85, 0x14, 0x26, 0xd8, 0x17, 0x09, 0x9a, 0x71, 0xc0, 0x90, 0xce, 0x5d,
	0x3a, 0xbd, 0xf7, 0xd3, 0xbc, 0x79, 0x03, 0x0f, 0xff, 0x81, 0x74, 0xce, 0x6d, 0x2e, 0x88, 0x9f,
	0x02, 0x17, 0x24, 0x01, 0x5b, 0xcc, 0x0b, 0xe0, 0x56, 0x51, 0x32, 0xc1, 0xf4, 0x6f, 0x35, 0xb5,
	0x76, 0x68, 0xef, 0x47, 0xc8, 0x42, 0x26, 0xa1, 0x5d, 0x7f, 0xa9, 0xff, 0x7a, 0xbf, 0x28, 0xe3,
	0x19, 0xe3, 0x9e, 0x02, 0x4a, 0x28, 0x34, 0xbc, 0x47, 0xf8, 0xab, 0x43, 0x52, 0x92, 0x53, 0x70,
	0x58, 0x59, 0xb2, 0x19, 0x04, 0xfa, 0x19, 0xc6, 0x15, 0x0f, 0x3c, 0x92, 0xb1, 0x2a, 0x17, 0x06,
	0x1a, 0xa0, 0x51, 0xdb, 0x39, 0x58, 0xac, 0xfa, 0xda, 0xf3, 0xaa, 0xff, 0x5b, 0x9d, 0xe6, 0x41,
	0x62, 0xc5, 0xcc, 0xce, 0x88, 0x88, 0xac, 0x53, 0x08, 0x09, 0x9d, 0x4f, 0x80, 0x3e, 0xde, 0x8d,
	0xf1, 0x36, 0x7c, 0x02, 0xd4, 0x6d, 0x57, 0x3c, 0x38, 0x92, 0x19, 0xfa, 0x39, 0xc6, 0x05, 0x94,
	0x14, 0x72, 0x41, 0x42, 0x30, 0x3e, 0xed, 0x9b, 0xb8, 0x13, 0x32, 0x7c, 0x40, 0xb8, 0x3b, 0xcd,
	0x05, 0x94, 0xc0, 0x85, 0x93, 0x32, 0x9a, 0xe8, 0x97, 0xb8, 0x1b, 0x6f, 0x0d, 0xaf, 0x24, 0x02,
	0xf6, 0x6f, 0xde, 0x79, 0xcb, 0x71, 0x89, 0x00, 0xfd, 0x2f, 0xc6, 0x7e, 0x7d, 0x81, 0x27, 0xe2,
	0x4c, 0x95, 0x6f, 0xb8, 0x6d, 0xe9, 0x5c, 0xc4, 0x19, 0xe8, 0xff, 0x70, 0x47, 0xe1, 0x08, 0xe2,
	0x30, 0x12, 0x46, 0x63, 0x80, 0x46, 0x4d, 0xf7, 0x8b, 0xf4, 0x8e, 0xa5, 0xa5, 0xff, 0xc4, 0x9f,
	0x0b, 0xc6, 0x52, 0x2f, 0x0e, 0x8c, 0xa6, 0xa4, 0xad, 0x5a, 0x4e, 0x83, 0xe1, 0x2d, 0xc2, 0xdf,
	0x55, 0x89, 0x0f, 0xf2, 0x14, 0xe7, 0x64, 0xb1, 0x36, 0xd1, 0x72, 0x6d, 0xa2, 0x97, 0xb5, 0x89,
	0x6e, 0x36, 0xa6, 0xb6, 0xdc, 0x98, 0xda, 0xd3, 0xc6, 0xd4, 0xae, 0xfe, 0x87, 0xb1, 0x88, 0x2a,
	0xdf, 0xa2, 0x2c, 0xb3, 0xeb, 0x5d, 0x8e, 0x73, 0x10, 0x33, 0x56, 0x26, 0x52, 0xd8, 0xd7, 0xef,
	0x47, 0xec, 0xb7, 0xe4, 0x04, 0x0f, 0x5f, 0x03, 0x00, 0x00, 0xff, 0xff, 0x9f, 0xc4, 0x4b, 0x2d,
	0xe5, 0x02, 0x00, 0x00,
}

func (m *BalanceBorrowed) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BalanceBorrowed) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BalanceBorrowed) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Percentage.Size()
		i -= size
		if _, err := m.Percentage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.UsdAmount.Size()
		i -= size
		if _, err := m.UsdAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *InterestBlock) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InterestBlock) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InterestBlock) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.PoolId != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x20
	}
	if m.BlockHeight != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x18
	}
	if m.BlockTime != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.BlockTime))
		i--
		dAtA[i] = 0x10
	}
	{
		size := m.InterestRate.Size()
		i -= size
		if _, err := m.InterestRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *LegacyInterestBlock) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LegacyInterestBlock) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LegacyInterestBlock) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.BlockHeight != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x18
	}
	if m.BlockTime != 0 {
		i = encodeVarintTypes(dAtA, i, uint64(m.BlockTime))
		i--
		dAtA[i] = 0x10
	}
	{
		size := m.InterestRate.Size()
		i -= size
		if _, err := m.InterestRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTypes(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
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
func (m *BalanceBorrowed) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.UsdAmount.Size()
	n += 1 + l + sovTypes(uint64(l))
	l = m.Percentage.Size()
	n += 1 + l + sovTypes(uint64(l))
	return n
}

func (m *InterestBlock) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.InterestRate.Size()
	n += 1 + l + sovTypes(uint64(l))
	if m.BlockTime != 0 {
		n += 1 + sovTypes(uint64(m.BlockTime))
	}
	if m.BlockHeight != 0 {
		n += 1 + sovTypes(uint64(m.BlockHeight))
	}
	if m.PoolId != 0 {
		n += 1 + sovTypes(uint64(m.PoolId))
	}
	return n
}

func (m *LegacyInterestBlock) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.InterestRate.Size()
	n += 1 + l + sovTypes(uint64(l))
	if m.BlockTime != 0 {
		n += 1 + sovTypes(uint64(m.BlockTime))
	}
	if m.BlockHeight != 0 {
		n += 1 + sovTypes(uint64(m.BlockHeight))
	}
	return n
}

func sovTypes(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTypes(x uint64) (n int) {
	return sovTypes(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *BalanceBorrowed) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: BalanceBorrowed: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BalanceBorrowed: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UsdAmount", wireType)
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
			if err := m.UsdAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Percentage", wireType)
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
			if err := m.Percentage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *InterestBlock) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: InterestBlock: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InterestBlock: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRate", wireType)
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
			if err := m.InterestRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockTime", wireType)
			}
			m.BlockTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
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
func (m *LegacyInterestBlock) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: LegacyInterestBlock: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LegacyInterestBlock: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InterestRate", wireType)
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
			if err := m.InterestRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockTime", wireType)
			}
			m.BlockTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTypes
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= uint64(b&0x7F) << shift
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
