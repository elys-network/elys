// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/amm/pool_asset.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

type PoolAsset struct {
	Token                  types.Coin                  `protobuf:"bytes,1,opt,name=token,proto3" json:"token"`
	Weight                 cosmossdk_io_math.Int       `protobuf:"bytes,2,opt,name=weight,proto3,customtype=cosmossdk.io/math.Int" json:"weight"`
	ExternalLiquidityRatio cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=external_liquidity_ratio,json=externalLiquidityRatio,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"external_liquidity_ratio"`
}

func (m *PoolAsset) Reset()         { *m = PoolAsset{} }
func (m *PoolAsset) String() string { return proto.CompactTextString(m) }
func (*PoolAsset) ProtoMessage()    {}
func (*PoolAsset) Descriptor() ([]byte, []int) {
	return fileDescriptor_152994f419c8cb00, []int{0}
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

func (m *PoolAsset) GetToken() types.Coin {
	if m != nil {
		return m.Token
	}
	return types.Coin{}
}

func init() {
	proto.RegisterType((*PoolAsset)(nil), "elys.amm.PoolAsset")
}

func init() { proto.RegisterFile("elys/amm/pool_asset.proto", fileDescriptor_152994f419c8cb00) }

var fileDescriptor_152994f419c8cb00 = []byte{
	// 332 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0x4d, 0x4a, 0xc3, 0x40,
	0x14, 0xc7, 0x13, 0x3f, 0x8a, 0x8d, 0xbb, 0xa0, 0x92, 0x56, 0x48, 0x8b, 0xab, 0x82, 0x74, 0x86,
	0x2a, 0x1e, 0xc0, 0xb4, 0x9b, 0x42, 0x17, 0x92, 0xa5, 0x9b, 0x30, 0x49, 0x87, 0x74, 0xc8, 0xc7,
	0xab, 0x99, 0x57, 0xdb, 0xdc, 0xc2, 0xc3, 0x78, 0x88, 0x2e, 0x8b, 0x2b, 0x71, 0x51, 0xa4, 0xbd,
	0x82, 0x07, 0x90, 0xc9, 0xa4, 0x20, 0xb8, 0x9b, 0x37, 0xbf, 0xf7, 0xff, 0x0d, 0xfc, 0xc7, 0x6a,
	0xf1, 0xb4, 0x94, 0x94, 0x65, 0x19, 0x9d, 0x03, 0xa4, 0x01, 0x93, 0x92, 0x23, 0x99, 0x17, 0x80,
	0x60, 0x9f, 0x29, 0x44, 0x58, 0x96, 0xb5, 0x2f, 0x62, 0x88, 0xa1, 0xba, 0xa4, 0xea, 0xa4, 0x79,
	0xdb, 0x8d, 0x40, 0x66, 0x20, 0x69, 0xc8, 0x24, 0xa7, 0xaf, 0x83, 0x90, 0x23, 0x1b, 0xd0, 0x08,
	0x44, 0x5e, 0xf3, 0x96, 0xe6, 0x81, 0x0e, 0xea, 0x41, 0xa3, 0x9b, 0x1f, 0xd3, 0x6a, 0x3e, 0x01,
	0xa4, 0x8f, 0xea, 0x39, 0xfb, 0xc1, 0x3a, 0x45, 0x48, 0x78, 0xee, 0x98, 0x5d, 0xb3, 0x77, 0x7e,
	0xd7, 0x22, 0xf5, 0xae, 0x12, 0x93, 0x5a, 0x4c, 0x86, 0x20, 0x72, 0xef, 0x64, 0xbd, 0xed, 0x18,
	0xbe, 0xde, 0xb6, 0x87, 0x56, 0x63, 0xc9, 0x45, 0x3c, 0x43, 0xe7, 0xa8, 0x6b, 0xf6, 0x9a, 0xde,
	0xad, 0x82, 0x5f, 0xdb, 0xce, 0xa5, 0x8e, 0xcb, 0x69, 0x42, 0x04, 0xd0, 0x8c, 0xe1, 0x8c, 0x8c,
	0x73, 0xfc, 0x78, 0xef, 0x5b, 0xb5, 0x77, 0x9c, 0xa3, 0x5f, 0x47, 0xed, 0xc4, 0x72, 0xf8, 0x0a,
	0x79, 0x91, 0xb3, 0x34, 0x48, 0xc5, 0xcb, 0x42, 0x4c, 0x05, 0x96, 0x41, 0xc1, 0x50, 0x80, 0x73,
	0x5c, 0x69, 0x07, 0xb5, 0xf6, 0xfa, 0xbf, 0x76, 0xc2, 0x63, 0x16, 0x95, 0x23, 0x1e, 0xfd, 0x91,
	0x8f, 0x78, 0xe4, 0x5f, 0x1d, 0x94, 0x93, 0x83, 0xd1, 0x57, 0x42, 0xcf, 0x5b, 0xef, 0x5c, 0x73,
	0xb3, 0x73, 0xcd, 0xef, 0x9d, 0x6b, 0xbe, 0xed, 0x5d, 0x63, 0xb3, 0x77, 0x8d, 0xcf, 0xbd, 0x6b,
	0x3c, 0xf7, 0x62, 0x81, 0xb3, 0x45, 0x48, 0x22, 0xc8, 0xa8, 0xaa, 0xbd, 0x9f, 0x73, 0x5c, 0x42,
	0x91, 0x54, 0x03, 0x5d, 0x55, 0x1f, 0x84, 0xe5, 0x9c, 0xcb, 0xb0, 0x51, 0x35, 0x78, 0xff, 0x1b,
	0x00, 0x00, 0xff, 0xff, 0x6e, 0x7f, 0x1e, 0xf5, 0xb9, 0x01, 0x00, 0x00,
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
	{
		size := m.ExternalLiquidityRatio.Size()
		i -= size
		if _, err := m.ExternalLiquidityRatio.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPoolAsset(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Weight.Size()
		i -= size
		if _, err := m.Weight.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPoolAsset(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.Token.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintPoolAsset(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintPoolAsset(dAtA []byte, offset int, v uint64) int {
	offset -= sovPoolAsset(v)
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
	l = m.Token.Size()
	n += 1 + l + sovPoolAsset(uint64(l))
	l = m.Weight.Size()
	n += 1 + l + sovPoolAsset(uint64(l))
	l = m.ExternalLiquidityRatio.Size()
	n += 1 + l + sovPoolAsset(uint64(l))
	return n
}

func sovPoolAsset(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPoolAsset(x uint64) (n int) {
	return sovPoolAsset(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PoolAsset) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPoolAsset
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
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPoolAsset
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
				return ErrInvalidLengthPoolAsset
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPoolAsset
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Token.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Weight", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPoolAsset
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
				return ErrInvalidLengthPoolAsset
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPoolAsset
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Weight.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExternalLiquidityRatio", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPoolAsset
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
				return ErrInvalidLengthPoolAsset
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPoolAsset
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ExternalLiquidityRatio.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPoolAsset(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPoolAsset
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
func skipPoolAsset(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPoolAsset
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
					return 0, ErrIntOverflowPoolAsset
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
					return 0, ErrIntOverflowPoolAsset
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
				return 0, ErrInvalidLengthPoolAsset
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPoolAsset
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPoolAsset
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPoolAsset        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPoolAsset          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPoolAsset = fmt.Errorf("proto: unexpected end of group")
)
