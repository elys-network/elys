// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/parameter/params.proto

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

type LegacyParams struct {
	MinCommissionRate   cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=min_commission_rate,json=minCommissionRate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"min_commission_rate"`
	MaxVotingPower      cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=max_voting_power,json=maxVotingPower,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"max_voting_power"`
	MinSelfDelegation   cosmossdk_io_math.Int       `protobuf:"bytes,3,opt,name=min_self_delegation,json=minSelfDelegation,proto3,customtype=cosmossdk.io/math.Int" json:"min_self_delegation"`
	TotalBlocksPerYear  uint64                      `protobuf:"varint,4,opt,name=total_blocks_per_year,json=totalBlocksPerYear,proto3" json:"total_blocks_per_year,omitempty"`
	RewardsDataLifetime uint64                      `protobuf:"varint,5,opt,name=rewards_data_lifetime,json=rewardsDataLifetime,proto3" json:"rewards_data_lifetime,omitempty"`
}

func (m *LegacyParams) Reset()         { *m = LegacyParams{} }
func (m *LegacyParams) String() string { return proto.CompactTextString(m) }
func (*LegacyParams) ProtoMessage()    {}
func (*LegacyParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_b61780a5be327c2b, []int{0}
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

func (m *LegacyParams) GetTotalBlocksPerYear() uint64 {
	if m != nil {
		return m.TotalBlocksPerYear
	}
	return 0
}

func (m *LegacyParams) GetRewardsDataLifetime() uint64 {
	if m != nil {
		return m.RewardsDataLifetime
	}
	return 0
}

// Params defines the parameters for the module.
type Params struct {
	MinCommissionRate   cosmossdk_io_math.LegacyDec `protobuf:"bytes,1,opt,name=min_commission_rate,json=minCommissionRate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"min_commission_rate"`
	MaxVotingPower      cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=max_voting_power,json=maxVotingPower,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"max_voting_power"`
	MinSelfDelegation   cosmossdk_io_math.Int       `protobuf:"bytes,3,opt,name=min_self_delegation,json=minSelfDelegation,proto3,customtype=cosmossdk.io/math.Int" json:"min_self_delegation"`
	TotalBlocksPerYear  uint64                      `protobuf:"varint,4,opt,name=total_blocks_per_year,json=totalBlocksPerYear,proto3" json:"total_blocks_per_year,omitempty"`
	RewardsDataLifetime uint64                      `protobuf:"varint,5,opt,name=rewards_data_lifetime,json=rewardsDataLifetime,proto3" json:"rewards_data_lifetime,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_b61780a5be327c2b, []int{1}
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

func (m *Params) GetTotalBlocksPerYear() uint64 {
	if m != nil {
		return m.TotalBlocksPerYear
	}
	return 0
}

func (m *Params) GetRewardsDataLifetime() uint64 {
	if m != nil {
		return m.RewardsDataLifetime
	}
	return 0
}

func init() {
	proto.RegisterType((*LegacyParams)(nil), "elys.parameter.LegacyParams")
	proto.RegisterType((*Params)(nil), "elys.parameter.Params")
}

func init() { proto.RegisterFile("elys/parameter/params.proto", fileDescriptor_b61780a5be327c2b) }

var fileDescriptor_b61780a5be327c2b = []byte{
	// 401 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x93, 0xc1, 0x6a, 0x14, 0x41,
	0x10, 0x86, 0x77, 0xdc, 0x18, 0xb0, 0x91, 0xa0, 0x13, 0x17, 0xc6, 0x04, 0x26, 0x21, 0xa7, 0x80,
	0x64, 0x86, 0xe8, 0x1b, 0xac, 0x7b, 0x30, 0x90, 0xc3, 0xb2, 0x82, 0xa0, 0x39, 0x34, 0xb5, 0xb3,
	0xb5, 0x93, 0x66, 0xa7, 0xbb, 0x86, 0xee, 0xd2, 0xdd, 0x79, 0x07, 0x0f, 0xe2, 0xb3, 0xf8, 0x10,
	0x39, 0x06, 0x4f, 0xe2, 0x21, 0xc8, 0xee, 0x8b, 0xc8, 0x74, 0x0f, 0x4b, 0xc0, 0x9b, 0x27, 0x0f,
	0xde, 0xaa, 0xf8, 0xba, 0xbf, 0x2a, 0x28, 0x7e, 0x71, 0x88, 0x55, 0xe3, 0xf2, 0x1a, 0x2c, 0x68,
	0x64, 0xb4, 0xa1, 0x72, 0x59, 0x6d, 0x89, 0x29, 0xde, 0x6b, 0x61, 0xb6, 0x85, 0x07, 0xcf, 0x4a,
	0x2a, 0xc9, 0xa3, 0xbc, 0xad, 0xc2, 0xab, 0x83, 0xe7, 0x05, 0x39, 0x4d, 0x4e, 0x06, 0x10, 0x9a,
	0x80, 0x4e, 0xbe, 0xf6, 0xc5, 0xe3, 0x4b, 0x2c, 0xa1, 0x68, 0xc6, 0xde, 0x1b, 0x83, 0xd8, 0xd7,
	0xca, 0xc8, 0x82, 0xb4, 0x56, 0xce, 0x29, 0x32, 0xd2, 0x02, 0x63, 0x12, 0x1d, 0x47, 0xa7, 0x8f,
	0x86, 0xe7, 0x37, 0x77, 0x47, 0xbd, 0x9f, 0x77, 0x47, 0x87, 0xc1, 0xe1, 0x66, 0x8b, 0x4c, 0x51,
	0xae, 0x81, 0xaf, 0xb3, 0x20, 0x19, 0x61, 0xf1, 0xfd, 0xdb, 0x99, 0xe8, 0x46, 0x8c, 0xb0, 0x98,
	0x3c, 0xd5, 0xca, 0xbc, 0xde, 0xca, 0x26, 0xc0, 0x18, 0x5f, 0x89, 0x27, 0x1a, 0x56, 0xf2, 0x13,
	0xb1, 0x32, 0xa5, 0xac, 0x69, 0x89, 0x36, 0x79, 0xf0, 0xb7, 0xfe, 0x3d, 0x0d, 0xab, 0x77, 0xde,
	0x34, 0x6e, 0x45, 0xf1, 0x55, 0xd8, 0xdf, 0x61, 0x35, 0x97, 0x33, 0xac, 0xb0, 0x04, 0x56, 0x64,
	0x92, 0xbe, 0xf7, 0xbf, 0xe8, 0xfc, 0x83, 0x3f, 0xfd, 0x17, 0x86, 0xef, 0x99, 0x2f, 0x0c, 0xfb,
	0xcd, 0xdf, 0x62, 0x35, 0x1f, 0x6d, 0x2d, 0xf1, 0xb9, 0x18, 0x30, 0x31, 0x54, 0x72, 0x5a, 0x51,
	0xb1, 0x70, 0xb2, 0x46, 0x2b, 0x1b, 0x04, 0x9b, 0xec, 0x1c, 0x47, 0xa7, 0x3b, 0x93, 0xd8, 0xc3,
	0xa1, 0x67, 0x63, 0xb4, 0xef, 0x11, 0x6c, 0xfc, 0x52, 0x0c, 0x2c, 0x2e, 0xc1, 0xce, 0x9c, 0x9c,
	0x01, 0x83, 0xac, 0xd4, 0x1c, 0x59, 0x69, 0x4c, 0x1e, 0xfa, 0x2f, 0xfb, 0x1d, 0x1c, 0x01, 0xc3,
	0x65, 0x87, 0x4e, 0x3e, 0xf7, 0xc5, 0xee, 0xff, 0x73, 0xfc, 0x2b, 0xe7, 0x18, 0xbe, 0xb9, 0x59,
	0xa7, 0xd1, 0xed, 0x3a, 0x8d, 0x7e, 0xad, 0xd3, 0xe8, 0xcb, 0x26, 0xed, 0xdd, 0x6e, 0xd2, 0xde,
	0x8f, 0x4d, 0xda, 0xfb, 0x90, 0x95, 0x8a, 0xaf, 0x3f, 0x4e, 0xb3, 0x82, 0x74, 0xde, 0x26, 0xf1,
	0xcc, 0x20, 0x2f, 0xc9, 0x2e, 0x7c, 0x93, 0xaf, 0xee, 0xa5, 0x96, 0x9b, 0x1a, 0xdd, 0x74, 0xd7,
	0x87, 0xee, 0xd5, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb1, 0xa9, 0x12, 0x98, 0xd4, 0x03, 0x00,
	0x00,
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
	if m.RewardsDataLifetime != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.RewardsDataLifetime))
		i--
		dAtA[i] = 0x28
	}
	if m.TotalBlocksPerYear != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.TotalBlocksPerYear))
		i--
		dAtA[i] = 0x20
	}
	{
		size := m.MinSelfDelegation.Size()
		i -= size
		if _, err := m.MinSelfDelegation.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.MaxVotingPower.Size()
		i -= size
		if _, err := m.MaxVotingPower.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.MinCommissionRate.Size()
		i -= size
		if _, err := m.MinCommissionRate.MarshalTo(dAtA[i:]); err != nil {
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
	if m.RewardsDataLifetime != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.RewardsDataLifetime))
		i--
		dAtA[i] = 0x28
	}
	if m.TotalBlocksPerYear != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.TotalBlocksPerYear))
		i--
		dAtA[i] = 0x20
	}
	{
		size := m.MinSelfDelegation.Size()
		i -= size
		if _, err := m.MinSelfDelegation.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.MaxVotingPower.Size()
		i -= size
		if _, err := m.MaxVotingPower.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.MinCommissionRate.Size()
		i -= size
		if _, err := m.MinCommissionRate.MarshalTo(dAtA[i:]); err != nil {
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
	l = m.MinCommissionRate.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MaxVotingPower.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MinSelfDelegation.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.TotalBlocksPerYear != 0 {
		n += 1 + sovParams(uint64(m.TotalBlocksPerYear))
	}
	if m.RewardsDataLifetime != 0 {
		n += 1 + sovParams(uint64(m.RewardsDataLifetime))
	}
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.MinCommissionRate.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MaxVotingPower.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MinSelfDelegation.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.TotalBlocksPerYear != 0 {
		n += 1 + sovParams(uint64(m.TotalBlocksPerYear))
	}
	if m.RewardsDataLifetime != 0 {
		n += 1 + sovParams(uint64(m.RewardsDataLifetime))
	}
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
				return fmt.Errorf("proto: wrong wireType = %d for field MinCommissionRate", wireType)
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
			if err := m.MinCommissionRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxVotingPower", wireType)
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
			if err := m.MaxVotingPower.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinSelfDelegation", wireType)
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
			if err := m.MinSelfDelegation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalBlocksPerYear", wireType)
			}
			m.TotalBlocksPerYear = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalBlocksPerYear |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardsDataLifetime", wireType)
			}
			m.RewardsDataLifetime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RewardsDataLifetime |= uint64(b&0x7F) << shift
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
				return fmt.Errorf("proto: wrong wireType = %d for field MinCommissionRate", wireType)
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
			if err := m.MinCommissionRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxVotingPower", wireType)
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
			if err := m.MaxVotingPower.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinSelfDelegation", wireType)
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
			if err := m.MinSelfDelegation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalBlocksPerYear", wireType)
			}
			m.TotalBlocksPerYear = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalBlocksPerYear |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardsDataLifetime", wireType)
			}
			m.RewardsDataLifetime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RewardsDataLifetime |= uint64(b&0x7F) << shift
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
