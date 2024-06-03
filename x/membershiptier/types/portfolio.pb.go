// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/membershiptier/portfolio.proto

package types

import (
	fmt "fmt"
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

type Portfolio struct {
	Assetkey     string     `protobuf:"bytes,2,opt,name=assetkey,proto3" json:"assetkey,omitempty"`
	Token        types.Coin `protobuf:"bytes,1,opt,name=token,proto3" json:"token"`
	Creator      string     `protobuf:"bytes,3,opt,name=creator,proto3" json:"creator,omitempty"`
	MinimumToday string     `protobuf:"bytes,4,opt,name=minimum_today,json=minimumToday,proto3" json:"minimum_today,omitempty"`
}

func (m *Portfolio) Reset()         { *m = Portfolio{} }
func (m *Portfolio) String() string { return proto.CompactTextString(m) }
func (*Portfolio) ProtoMessage()    {}
func (*Portfolio) Descriptor() ([]byte, []int) {
	return fileDescriptor_de2fe2eb41556835, []int{0}
}
func (m *Portfolio) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Portfolio) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Portfolio.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Portfolio) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Portfolio.Merge(m, src)
}
func (m *Portfolio) XXX_Size() int {
	return m.Size()
}
func (m *Portfolio) XXX_DiscardUnknown() {
	xxx_messageInfo_Portfolio.DiscardUnknown(m)
}

var xxx_messageInfo_Portfolio proto.InternalMessageInfo

func (m *Portfolio) GetAssetkey() string {
	if m != nil {
		return m.Assetkey
	}
	return ""
}

func (m *Portfolio) GetToken() types.Coin {
	if m != nil {
		return m.Token
	}
	return types.Coin{}
}

func (m *Portfolio) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *Portfolio) GetMinimumToday() string {
	if m != nil {
		return m.MinimumToday
	}
	return ""
}

func init() {
	proto.RegisterType((*Portfolio)(nil), "elys.membershiptier.Portfolio")
}

func init() {
	proto.RegisterFile("elys/membershiptier/portfolio.proto", fileDescriptor_de2fe2eb41556835)
}

var fileDescriptor_de2fe2eb41556835 = []byte{
	// 284 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xbf, 0x4e, 0xf3, 0x30,
	0x14, 0xc5, 0xe3, 0xef, 0x2b, 0x7f, 0x6a, 0x60, 0x09, 0x0c, 0x21, 0x83, 0xa9, 0xe8, 0xd2, 0x05,
	0x5b, 0xa5, 0xe2, 0x05, 0xca, 0x8c, 0x84, 0x2a, 0x26, 0x16, 0x94, 0x84, 0x4b, 0x6a, 0xa5, 0xce,
	0x8d, 0xec, 0x5b, 0x20, 0x6f, 0xc1, 0x0b, 0xf0, 0x3e, 0x1d, 0x3b, 0x32, 0x21, 0x94, 0xbc, 0x08,
	0x4a, 0xd2, 0x22, 0xc1, 0xe6, 0x73, 0xcf, 0x4f, 0x3f, 0xc9, 0x87, 0x0f, 0x61, 0x51, 0x3a, 0x65,
	0xc0, 0xc4, 0x60, 0xdd, 0x5c, 0x17, 0xa4, 0xc1, 0xaa, 0x02, 0x2d, 0x3d, 0xe1, 0x42, 0xa3, 0x2c,
	0x2c, 0x12, 0xfa, 0xc7, 0x0d, 0x24, 0x7f, 0x43, 0xe1, 0x49, 0x8a, 0x29, 0xb6, 0xbd, 0x6a, 0x5e,
	0x1d, 0x1a, 0x8a, 0x04, 0x9d, 0x41, 0xa7, 0xe2, 0xc8, 0x81, 0x7a, 0x1e, 0xc7, 0x40, 0xd1, 0x58,
	0x25, 0xa8, 0xf3, 0xae, 0x3f, 0x7f, 0x67, 0xbc, 0x7f, 0xbb, 0xd5, 0xfb, 0x21, 0xdf, 0x8f, 0x9c,
	0x03, 0xca, 0xa0, 0x0c, 0xfe, 0x0d, 0xd8, 0xa8, 0x3f, 0xfb, 0xc9, 0xfe, 0x15, 0xdf, 0x21, 0xcc,
	0x20, 0x0f, 0xd8, 0x80, 0x8d, 0x0e, 0x2e, 0x4f, 0x65, 0x67, 0x96, 0x8d, 0x59, 0x6e, 0xcc, 0xf2,
	0x1a, 0x75, 0x3e, 0xed, 0xad, 0x3e, 0xcf, 0xbc, 0x59, 0x47, 0xfb, 0x01, 0xdf, 0x4b, 0x2c, 0x44,
	0x84, 0x36, 0xf8, 0xdf, 0x1a, 0xb7, 0xd1, 0x1f, 0xf2, 0x23, 0xa3, 0x73, 0x6d, 0x96, 0xe6, 0x81,
	0xf0, 0x31, 0x2a, 0x83, 0x5e, 0xdb, 0x1f, 0x6e, 0x8e, 0x77, 0xcd, 0x6d, 0x7a, 0xb3, 0xaa, 0x04,
	0x5b, 0x57, 0x82, 0x7d, 0x55, 0x82, 0xbd, 0xd5, 0xc2, 0x5b, 0xd7, 0xc2, 0xfb, 0xa8, 0x85, 0x77,
	0x3f, 0x49, 0x35, 0xcd, 0x97, 0xb1, 0x4c, 0xd0, 0xa8, 0x66, 0x8f, 0x8b, 0x1c, 0xe8, 0x05, 0x6d,
	0xd6, 0x06, 0xf5, 0xfa, 0x77, 0x43, 0x2a, 0x0b, 0x70, 0xf1, 0x6e, 0xfb, 0xeb, 0xc9, 0x77, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x06, 0x0a, 0x41, 0x56, 0x67, 0x01, 0x00, 0x00,
}

func (m *Portfolio) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Portfolio) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Portfolio) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.MinimumToday) > 0 {
		i -= len(m.MinimumToday)
		copy(dAtA[i:], m.MinimumToday)
		i = encodeVarintPortfolio(dAtA, i, uint64(len(m.MinimumToday)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintPortfolio(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Assetkey) > 0 {
		i -= len(m.Assetkey)
		copy(dAtA[i:], m.Assetkey)
		i = encodeVarintPortfolio(dAtA, i, uint64(len(m.Assetkey)))
		i--
		dAtA[i] = 0x12
	}
	{
		size, err := m.Token.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintPortfolio(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintPortfolio(dAtA []byte, offset int, v uint64) int {
	offset -= sovPortfolio(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Portfolio) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Token.Size()
	n += 1 + l + sovPortfolio(uint64(l))
	l = len(m.Assetkey)
	if l > 0 {
		n += 1 + l + sovPortfolio(uint64(l))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovPortfolio(uint64(l))
	}
	l = len(m.MinimumToday)
	if l > 0 {
		n += 1 + l + sovPortfolio(uint64(l))
	}
	return n
}

func sovPortfolio(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPortfolio(x uint64) (n int) {
	return sovPortfolio(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Portfolio) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowPortfolio
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
			return fmt.Errorf("proto: Portfolio: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Portfolio: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPortfolio
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
				return ErrInvalidLengthPortfolio
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthPortfolio
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
				return fmt.Errorf("proto: wrong wireType = %d for field Assetkey", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPortfolio
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
				return ErrInvalidLengthPortfolio
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPortfolio
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Assetkey = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPortfolio
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
				return ErrInvalidLengthPortfolio
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPortfolio
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinimumToday", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowPortfolio
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
				return ErrInvalidLengthPortfolio
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthPortfolio
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MinimumToday = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipPortfolio(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthPortfolio
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
func skipPortfolio(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowPortfolio
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
					return 0, ErrIntOverflowPortfolio
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
					return 0, ErrIntOverflowPortfolio
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
				return 0, ErrInvalidLengthPortfolio
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupPortfolio
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthPortfolio
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthPortfolio        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowPortfolio          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupPortfolio = fmt.Errorf("proto: unexpected end of group")
)
