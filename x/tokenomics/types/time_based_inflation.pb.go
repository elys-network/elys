// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/tokenomics/time_based_inflation.proto

package types

import (
	fmt "fmt"
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

type TimeBasedInflation struct {
	StartBlockHeight uint64          `protobuf:"varint,1,opt,name=start_block_height,json=startBlockHeight,proto3" json:"start_block_height,omitempty"`
	EndBlockHeight   uint64          `protobuf:"varint,2,opt,name=end_block_height,json=endBlockHeight,proto3" json:"end_block_height,omitempty"`
	Description      string          `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Inflation        *InflationEntry `protobuf:"bytes,4,opt,name=inflation,proto3" json:"inflation,omitempty"`
	Authority        string          `protobuf:"bytes,5,opt,name=authority,proto3" json:"authority,omitempty"`
}

func (m *TimeBasedInflation) Reset()         { *m = TimeBasedInflation{} }
func (m *TimeBasedInflation) String() string { return proto.CompactTextString(m) }
func (*TimeBasedInflation) ProtoMessage()    {}
func (*TimeBasedInflation) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4fe870416dff3fa, []int{0}
}
func (m *TimeBasedInflation) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TimeBasedInflation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TimeBasedInflation.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TimeBasedInflation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TimeBasedInflation.Merge(m, src)
}
func (m *TimeBasedInflation) XXX_Size() int {
	return m.Size()
}
func (m *TimeBasedInflation) XXX_DiscardUnknown() {
	xxx_messageInfo_TimeBasedInflation.DiscardUnknown(m)
}

var xxx_messageInfo_TimeBasedInflation proto.InternalMessageInfo

func (m *TimeBasedInflation) GetStartBlockHeight() uint64 {
	if m != nil {
		return m.StartBlockHeight
	}
	return 0
}

func (m *TimeBasedInflation) GetEndBlockHeight() uint64 {
	if m != nil {
		return m.EndBlockHeight
	}
	return 0
}

func (m *TimeBasedInflation) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *TimeBasedInflation) GetInflation() *InflationEntry {
	if m != nil {
		return m.Inflation
	}
	return nil
}

func (m *TimeBasedInflation) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func init() {
	proto.RegisterType((*TimeBasedInflation)(nil), "elys.tokenomics.TimeBasedInflation")
}

func init() {
	proto.RegisterFile("elys/tokenomics/time_based_inflation.proto", fileDescriptor_a4fe870416dff3fa)
}

var fileDescriptor_a4fe870416dff3fa = []byte{
	// 296 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0x4f, 0x4b, 0xf3, 0x40,
	0x10, 0xc6, 0xbb, 0xef, 0x5b, 0x85, 0x6e, 0x41, 0xcb, 0x9e, 0x82, 0xc8, 0x1a, 0x04, 0x21, 0x88,
	0x26, 0x50, 0xc1, 0x9b, 0x97, 0x82, 0xa0, 0xe0, 0xa9, 0x78, 0xf2, 0x12, 0xf2, 0x67, 0x6c, 0x96,
	0x36, 0xbb, 0x65, 0x77, 0xaa, 0xe6, 0x5b, 0xf8, 0xb1, 0x3c, 0xf6, 0xe8, 0x51, 0x92, 0xa3, 0x5f,
	0x42, 0x76, 0xc5, 0xa4, 0xed, 0x71, 0x9f, 0xf9, 0x3d, 0x33, 0xf3, 0xec, 0xd0, 0x73, 0x58, 0x54,
	0x26, 0x42, 0x35, 0x07, 0xa9, 0x4a, 0x91, 0x99, 0x08, 0x45, 0x09, 0x71, 0x9a, 0x18, 0xc8, 0x63,
	0x21, 0x9f, 0x17, 0x09, 0x0a, 0x25, 0xc3, 0xa5, 0x56, 0xa8, 0xd8, 0xa1, 0x65, 0xc3, 0x8e, 0x3d,
	0x3a, 0xdb, 0x35, 0xb7, 0x8e, 0x18, 0x24, 0xea, 0xea, 0xd7, 0x77, 0xfa, 0x4d, 0x28, 0x7b, 0x14,
	0x25, 0x4c, 0x6c, 0xd7, 0xfb, 0x3f, 0x84, 0x5d, 0x50, 0x66, 0x30, 0xd1, 0x18, 0xa7, 0x0b, 0x95,
	0xcd, 0xe3, 0x02, 0xc4, 0xac, 0x40, 0x8f, 0xf8, 0x24, 0xe8, 0x4f, 0x47, 0xae, 0x32, 0xb1, 0x85,
	0x3b, 0xa7, 0xb3, 0x80, 0x8e, 0x40, 0xe6, 0xdb, 0xec, 0x3f, 0xc7, 0x1e, 0x80, 0xcc, 0x37, 0x49,
	0x9f, 0x0e, 0x73, 0x30, 0x99, 0x16, 0x4b, 0x3b, 0xc6, 0xfb, 0xef, 0x93, 0x60, 0x30, 0xdd, 0x94,
	0xd8, 0x0d, 0x1d, 0xb4, 0x9b, 0x7a, 0x7d, 0x9f, 0x04, 0xc3, 0xf1, 0x49, 0xb8, 0x13, 0x2e, 0x6c,
	0x17, 0xbd, 0xb5, 0x51, 0xa6, 0x9d, 0x83, 0x1d, 0xd3, 0x41, 0xb2, 0xc2, 0x42, 0x69, 0x81, 0x95,
	0xb7, 0xe7, 0xda, 0x77, 0xc2, 0xe4, 0xe1, 0xa3, 0xe6, 0x64, 0x5d, 0x73, 0xf2, 0x55, 0x73, 0xf2,
	0xde, 0xf0, 0xde, 0xba, 0xe1, 0xbd, 0xcf, 0x86, 0xf7, 0x9e, 0xc6, 0x33, 0x81, 0xc5, 0x2a, 0x0d,
	0x33, 0x55, 0x46, 0x76, 0xda, 0xa5, 0x04, 0x7c, 0x55, 0x7a, 0xee, 0x1e, 0xd1, 0xcb, 0x75, 0xf4,
	0xb6, 0x75, 0x88, 0x6a, 0x09, 0x26, 0xdd, 0x77, 0x5f, 0x78, 0xf5, 0x13, 0x00, 0x00, 0xff, 0xff,
	0xc6, 0xee, 0x4f, 0xc1, 0xa8, 0x01, 0x00, 0x00,
}

func (m *TimeBasedInflation) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TimeBasedInflation) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TimeBasedInflation) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintTimeBasedInflation(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0x2a
	}
	if m.Inflation != nil {
		{
			size, err := m.Inflation.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTimeBasedInflation(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintTimeBasedInflation(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x1a
	}
	if m.EndBlockHeight != 0 {
		i = encodeVarintTimeBasedInflation(dAtA, i, uint64(m.EndBlockHeight))
		i--
		dAtA[i] = 0x10
	}
	if m.StartBlockHeight != 0 {
		i = encodeVarintTimeBasedInflation(dAtA, i, uint64(m.StartBlockHeight))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintTimeBasedInflation(dAtA []byte, offset int, v uint64) int {
	offset -= sovTimeBasedInflation(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TimeBasedInflation) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.StartBlockHeight != 0 {
		n += 1 + sovTimeBasedInflation(uint64(m.StartBlockHeight))
	}
	if m.EndBlockHeight != 0 {
		n += 1 + sovTimeBasedInflation(uint64(m.EndBlockHeight))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovTimeBasedInflation(uint64(l))
	}
	if m.Inflation != nil {
		l = m.Inflation.Size()
		n += 1 + l + sovTimeBasedInflation(uint64(l))
	}
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovTimeBasedInflation(uint64(l))
	}
	return n
}

func sovTimeBasedInflation(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTimeBasedInflation(x uint64) (n int) {
	return sovTimeBasedInflation(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TimeBasedInflation) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTimeBasedInflation
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
			return fmt.Errorf("proto: TimeBasedInflation: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TimeBasedInflation: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StartBlockHeight", wireType)
			}
			m.StartBlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTimeBasedInflation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StartBlockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EndBlockHeight", wireType)
			}
			m.EndBlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTimeBasedInflation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EndBlockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTimeBasedInflation
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
				return ErrInvalidLengthTimeBasedInflation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTimeBasedInflation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Inflation", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTimeBasedInflation
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
				return ErrInvalidLengthTimeBasedInflation
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTimeBasedInflation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Inflation == nil {
				m.Inflation = &InflationEntry{}
			}
			if err := m.Inflation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTimeBasedInflation
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
				return ErrInvalidLengthTimeBasedInflation
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTimeBasedInflation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTimeBasedInflation(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTimeBasedInflation
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
func skipTimeBasedInflation(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTimeBasedInflation
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
					return 0, ErrIntOverflowTimeBasedInflation
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
					return 0, ErrIntOverflowTimeBasedInflation
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
				return 0, ErrInvalidLengthTimeBasedInflation
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTimeBasedInflation
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTimeBasedInflation
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTimeBasedInflation        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTimeBasedInflation          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTimeBasedInflation = fmt.Errorf("proto: unexpected end of group")
)
