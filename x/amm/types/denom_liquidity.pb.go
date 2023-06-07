// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/amm/denom_liquidity.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

type DenomLiquidity struct {
	Denom     string                                 `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	Liquidity github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=liquidity,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"liquidity"`
}

func (m *DenomLiquidity) Reset()         { *m = DenomLiquidity{} }
func (m *DenomLiquidity) String() string { return proto.CompactTextString(m) }
func (*DenomLiquidity) ProtoMessage()    {}
func (*DenomLiquidity) Descriptor() ([]byte, []int) {
	return fileDescriptor_0f4b59ff52ffc2fb, []int{0}
}
func (m *DenomLiquidity) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DenomLiquidity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DenomLiquidity.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DenomLiquidity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DenomLiquidity.Merge(m, src)
}
func (m *DenomLiquidity) XXX_Size() int {
	return m.Size()
}
func (m *DenomLiquidity) XXX_DiscardUnknown() {
	xxx_messageInfo_DenomLiquidity.DiscardUnknown(m)
}

var xxx_messageInfo_DenomLiquidity proto.InternalMessageInfo

func (m *DenomLiquidity) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func init() {
	proto.RegisterType((*DenomLiquidity)(nil), "elysnetwork.elys.amm.DenomLiquidity")
}

func init() { proto.RegisterFile("elys/amm/denom_liquidity.proto", fileDescriptor_0f4b59ff52ffc2fb) }

var fileDescriptor_0f4b59ff52ffc2fb = []byte{
	// 243 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0xcd, 0xa9, 0x2c,
	0xd6, 0x4f, 0xcc, 0xcd, 0xd5, 0x4f, 0x49, 0xcd, 0xcb, 0xcf, 0x8d, 0xcf, 0xc9, 0x2c, 0x2c, 0xcd,
	0x4c, 0xc9, 0x2c, 0xa9, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12, 0x01, 0xc9, 0xe7, 0xa5,
	0x96, 0x94, 0xe7, 0x17, 0x65, 0xeb, 0x81, 0xd8, 0x7a, 0x89, 0xb9, 0xb9, 0x52, 0x22, 0xe9, 0xf9,
	0xe9, 0xf9, 0x60, 0x05, 0xfa, 0x20, 0x16, 0x44, 0xad, 0x94, 0x5c, 0x72, 0x7e, 0x71, 0x6e, 0x7e,
	0xb1, 0x7e, 0x52, 0x62, 0x71, 0xaa, 0x7e, 0x99, 0x61, 0x52, 0x6a, 0x49, 0xa2, 0xa1, 0x7e, 0x72,
	0x7e, 0x66, 0x1e, 0x44, 0x5e, 0xa9, 0x84, 0x8b, 0xcf, 0x05, 0x64, 0x89, 0x0f, 0xcc, 0x0e, 0x21,
	0x11, 0x2e, 0x56, 0xb0, 0xb5, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x41, 0x10, 0x8e, 0x90, 0x0f,
	0x17, 0x27, 0xdc, 0x19, 0x12, 0x4c, 0x20, 0x19, 0x27, 0xbd, 0x13, 0xf7, 0xe4, 0x19, 0x6e, 0xdd,
	0x93, 0x57, 0x4b, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x87, 0xda, 0x06,
	0xa1, 0x74, 0x8b, 0x53, 0xb2, 0xf5, 0x4b, 0x2a, 0x0b, 0x52, 0x8b, 0xf5, 0x3c, 0xf3, 0x4a, 0x82,
	0x10, 0x06, 0x38, 0x39, 0x9d, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72,
	0x8c, 0x13, 0x1e, 0xcb, 0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x43, 0x94, 0x06,
	0x92, 0x61, 0x20, 0xaf, 0xe9, 0x42, 0xfd, 0x09, 0xe6, 0xe8, 0x57, 0x80, 0x43, 0x05, 0x6c, 0x64,
	0x12, 0x1b, 0xd8, 0x03, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x08, 0x6f, 0x4f, 0x06, 0x2e,
	0x01, 0x00, 0x00,
}

func (m *DenomLiquidity) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DenomLiquidity) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DenomLiquidity) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Liquidity.Size()
		i -= size
		if _, err := m.Liquidity.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDenomLiquidity(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintDenomLiquidity(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintDenomLiquidity(dAtA []byte, offset int, v uint64) int {
	offset -= sovDenomLiquidity(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DenomLiquidity) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovDenomLiquidity(uint64(l))
	}
	l = m.Liquidity.Size()
	n += 1 + l + sovDenomLiquidity(uint64(l))
	return n
}

func sovDenomLiquidity(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDenomLiquidity(x uint64) (n int) {
	return sovDenomLiquidity(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DenomLiquidity) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDenomLiquidity
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
			return fmt.Errorf("proto: DenomLiquidity: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DenomLiquidity: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDenomLiquidity
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
				return ErrInvalidLengthDenomLiquidity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDenomLiquidity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liquidity", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDenomLiquidity
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
				return ErrInvalidLengthDenomLiquidity
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDenomLiquidity
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Liquidity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDenomLiquidity(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDenomLiquidity
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
func skipDenomLiquidity(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDenomLiquidity
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
					return 0, ErrIntOverflowDenomLiquidity
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
					return 0, ErrIntOverflowDenomLiquidity
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
				return 0, ErrInvalidLengthDenomLiquidity
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDenomLiquidity
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDenomLiquidity
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDenomLiquidity        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDenomLiquidity          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDenomLiquidity = fmt.Errorf("proto: unexpected end of group")
)
