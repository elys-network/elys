// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/tier/portfolio.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types"
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

type LegacyPortfolio struct {
	Creator   string                      `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Portfolio cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=portfolio,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"portfolio"`
}

func (m *LegacyPortfolio) Reset()         { *m = LegacyPortfolio{} }
func (m *LegacyPortfolio) String() string { return proto.CompactTextString(m) }
func (*LegacyPortfolio) ProtoMessage()    {}
func (*LegacyPortfolio) Descriptor() ([]byte, []int) {
	return fileDescriptor_3d8e8fa01b8029d7, []int{0}
}
func (m *LegacyPortfolio) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LegacyPortfolio) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LegacyPortfolio.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LegacyPortfolio) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LegacyPortfolio.Merge(m, src)
}
func (m *LegacyPortfolio) XXX_Size() int {
	return m.Size()
}
func (m *LegacyPortfolio) XXX_DiscardUnknown() {
	xxx_messageInfo_LegacyPortfolio.DiscardUnknown(m)
}

var xxx_messageInfo_LegacyPortfolio proto.InternalMessageInfo

func (m *LegacyPortfolio) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

type Portfolio struct {
	Date      string                      `protobuf:"bytes,1,opt,name=date,proto3" json:"date,omitempty"`
	Creator   string                      `protobuf:"bytes,2,opt,name=creator,proto3" json:"creator,omitempty"`
	Portfolio cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=portfolio,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"portfolio"`
}

func (m *Portfolio) Reset()         { *m = Portfolio{} }
func (m *Portfolio) String() string { return proto.CompactTextString(m) }
func (*Portfolio) ProtoMessage()    {}
func (*Portfolio) Descriptor() ([]byte, []int) {
	return fileDescriptor_3d8e8fa01b8029d7, []int{1}
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

func (m *Portfolio) GetDate() string {
	if m != nil {
		return m.Date
	}
	return ""
}

func (m *Portfolio) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func init() {
	proto.RegisterType((*LegacyPortfolio)(nil), "elys.tier.LegacyPortfolio")
	proto.RegisterType((*Portfolio)(nil), "elys.tier.Portfolio")
}

func init() { proto.RegisterFile("elys/tier/portfolio.proto", fileDescriptor_3d8e8fa01b8029d7) }

var fileDescriptor_3d8e8fa01b8029d7 = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4c, 0xcd, 0xa9, 0x2c,
	0xd6, 0x2f, 0xc9, 0x4c, 0x2d, 0xd2, 0x2f, 0xc8, 0x2f, 0x2a, 0x49, 0xcb, 0xcf, 0xc9, 0xcc, 0xd7,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x04, 0x49, 0xe9, 0x81, 0xa4, 0xa4, 0x44, 0xd2, 0xf3,
	0xd3, 0xf3, 0xc1, 0xa2, 0xfa, 0x20, 0x16, 0x44, 0x81, 0x94, 0x5c, 0x72, 0x7e, 0x71, 0x6e, 0x7e,
	0xb1, 0x7e, 0x52, 0x62, 0x71, 0xaa, 0x7e, 0x99, 0x61, 0x52, 0x6a, 0x49, 0xa2, 0xa1, 0x7e, 0x72,
	0x7e, 0x66, 0x1e, 0x54, 0x5e, 0x12, 0x22, 0x1f, 0x0f, 0xd1, 0x08, 0xe1, 0x40, 0xa4, 0x94, 0xa6,
	0x31, 0x72, 0xf1, 0xfb, 0xa4, 0xa6, 0x27, 0x26, 0x57, 0x06, 0xc0, 0x6c, 0x15, 0x32, 0xe2, 0x62,
	0x4f, 0x2e, 0x4a, 0x4d, 0x2c, 0xc9, 0x2f, 0x92, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x74, 0x92, 0xb8,
	0xb4, 0x45, 0x57, 0x04, 0xaa, 0xcd, 0x31, 0x25, 0xa5, 0x28, 0xb5, 0xb8, 0x38, 0xb8, 0xa4, 0x28,
	0x33, 0x2f, 0x3d, 0x08, 0xa6, 0x50, 0xc8, 0x9f, 0x8b, 0x13, 0xee, 0x6c, 0x09, 0x26, 0xb0, 0x2e,
	0xc3, 0x13, 0xf7, 0xe4, 0x19, 0x6e, 0xdd, 0x93, 0x97, 0x86, 0xe8, 0x2c, 0x4e, 0xc9, 0xd6, 0xcb,
	0xcc, 0xd7, 0xcf, 0x4d, 0x2c, 0xc9, 0xd0, 0x83, 0xd8, 0xe8, 0x92, 0x9a, 0x7c, 0x69, 0x8b, 0x2e,
	0x17, 0xd4, 0x60, 0x97, 0xd4, 0xe4, 0x20, 0x84, 0x19, 0x4a, 0x4b, 0x18, 0xb9, 0x38, 0x11, 0x4e,
	0x12, 0xe2, 0x62, 0x49, 0x49, 0x2c, 0x49, 0x85, 0xb8, 0x27, 0x08, 0xcc, 0x46, 0x76, 0x26, 0x13,
	0x59, 0xce, 0x64, 0xa6, 0xdc, 0x99, 0x4e, 0x6e, 0x27, 0x1e, 0xc9, 0x31, 0x5e, 0x78, 0x24, 0xc7,
	0xf8, 0xe0, 0x91, 0x1c, 0xe3, 0x84, 0xc7, 0x72, 0x0c, 0x17, 0x1e, 0xcb, 0x31, 0xdc, 0x78, 0x2c,
	0xc7, 0x10, 0xa5, 0x93, 0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97, 0x9c, 0x9f, 0xab, 0x0f, 0x8a,
	0x40, 0xdd, 0xbc, 0xd4, 0x92, 0xf2, 0xfc, 0xa2, 0x6c, 0x30, 0x47, 0xbf, 0xcc, 0x5c, 0xbf, 0x02,
	0x12, 0xdb, 0x25, 0x95, 0x05, 0xa9, 0xc5, 0x49, 0x6c, 0xe0, 0xe8, 0x30, 0x06, 0x04, 0x00, 0x00,
	0xff, 0xff, 0xe6, 0x9f, 0x59, 0x10, 0x07, 0x02, 0x00, 0x00,
}

func (m *LegacyPortfolio) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LegacyPortfolio) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LegacyPortfolio) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Portfolio.Size()
		i -= size
		if _, err := m.Portfolio.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPortfolio(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintPortfolio(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
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
	{
		size := m.Portfolio.Size()
		i -= size
		if _, err := m.Portfolio.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintPortfolio(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintPortfolio(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Date) > 0 {
		i -= len(m.Date)
		copy(dAtA[i:], m.Date)
		i = encodeVarintPortfolio(dAtA, i, uint64(len(m.Date)))
		i--
		dAtA[i] = 0xa
	}
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
func (m *LegacyPortfolio) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovPortfolio(uint64(l))
	}
	l = m.Portfolio.Size()
	n += 1 + l + sovPortfolio(uint64(l))
	return n
}

func (m *Portfolio) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Date)
	if l > 0 {
		n += 1 + l + sovPortfolio(uint64(l))
	}
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovPortfolio(uint64(l))
	}
	l = m.Portfolio.Size()
	n += 1 + l + sovPortfolio(uint64(l))
	return n
}

func sovPortfolio(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozPortfolio(x uint64) (n int) {
	return sovPortfolio(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LegacyPortfolio) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: LegacyPortfolio: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LegacyPortfolio: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
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
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Portfolio", wireType)
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
			if err := m.Portfolio.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
				return fmt.Errorf("proto: wrong wireType = %d for field Date", wireType)
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
			m.Date = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
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
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Portfolio", wireType)
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
			if err := m.Portfolio.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
