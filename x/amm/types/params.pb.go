// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/amm/params.proto

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

// Params defines the parameters for the module.
type Params struct {
	PoolCreationFee                  cosmossdk_io_math.Int `protobuf:"bytes,1,opt,name=pool_creation_fee,json=poolCreationFee,proto3,customtype=cosmossdk.io/math.Int" json:"pool_creation_fee"`
	SlippageTrackDuration            uint64                `protobuf:"varint,2,opt,name=slippage_track_duration,json=slippageTrackDuration,proto3" json:"slippage_track_duration,omitempty"`
	EnableBaseCurrencyPairedPoolOnly bool                  `protobuf:"varint,3,opt,name=enable_base_currency_paired_pool_only,json=enableBaseCurrencyPairedPoolOnly,proto3" json:"enable_base_currency_paired_pool_only,omitempty"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_1209ca218537a425, []int{0}
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

func (m *Params) GetSlippageTrackDuration() uint64 {
	if m != nil {
		return m.SlippageTrackDuration
	}
	return 0
}

func (m *Params) GetEnableBaseCurrencyPairedPoolOnly() bool {
	if m != nil {
		return m.EnableBaseCurrencyPairedPoolOnly
	}
	return false
}

func init() {
	proto.RegisterType((*Params)(nil), "elys.amm.Params")
}

func init() { proto.RegisterFile("elys/amm/params.proto", fileDescriptor_1209ca218537a425) }

var fileDescriptor_1209ca218537a425 = []byte{
	// 334 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x34, 0x90, 0xc1, 0x4a, 0xeb, 0x40,
	0x14, 0x86, 0x33, 0xf7, 0x96, 0x52, 0xb3, 0x11, 0x83, 0xc5, 0xda, 0x45, 0x1a, 0x04, 0x21, 0x20,
	0x4d, 0x16, 0x82, 0x0b, 0x97, 0xa9, 0x08, 0x5d, 0xb5, 0x14, 0x41, 0x70, 0x33, 0x4c, 0xd2, 0x63,
	0x1a, 0x9a, 0x99, 0x13, 0x66, 0xa6, 0x68, 0xde, 0xc2, 0xa5, 0x4b, 0x1f, 0xc2, 0x87, 0xe8, 0xb2,
	0xb8, 0x12, 0x17, 0x45, 0xda, 0x17, 0xf0, 0x11, 0x24, 0x93, 0xb8, 0x9b, 0x73, 0xbe, 0x7f, 0x3e,
	0x0e, 0xbf, 0xdd, 0x85, 0xbc, 0x54, 0x21, 0xe3, 0x3c, 0x2c, 0x98, 0x64, 0x5c, 0x05, 0x85, 0x44,
	0x8d, 0x4e, 0xa7, 0x5a, 0x07, 0x8c, 0xf3, 0xfe, 0x71, 0x8a, 0x29, 0x9a, 0x65, 0x58, 0xbd, 0x6a,
	0xde, 0x3f, 0x4d, 0x50, 0x71, 0x54, 0xb4, 0x06, 0xf5, 0x50, 0xa3, 0xb3, 0x1f, 0x62, 0xb7, 0xa7,
	0xc6, 0xe5, 0xdc, 0xdb, 0x47, 0x05, 0x62, 0x4e, 0x13, 0x09, 0x4c, 0x67, 0x28, 0xe8, 0x23, 0x40,
	0x8f, 0x78, 0xc4, 0x3f, 0x88, 0x2e, 0xd6, 0xdb, 0x81, 0xf5, 0xb5, 0x1d, 0x74, 0xeb, 0xbf, 0x6a,
	0xbe, 0x0c, 0x32, 0x0c, 0x39, 0xd3, 0x8b, 0x60, 0x2c, 0xf4, 0xc7, 0xfb, 0xd0, 0x6e, 0xa4, 0x63,
	0xa1, 0x67, 0x87, 0x95, 0x65, 0xd4, 0x48, 0x6e, 0x01, 0x9c, 0x2b, 0xfb, 0x44, 0xe5, 0x59, 0x51,
	0xb0, 0x14, 0xa8, 0x96, 0x2c, 0x59, 0xd2, 0xf9, 0x4a, 0x1a, 0xda, 0xfb, 0xe7, 0x11, 0xbf, 0x35,
	0xeb, 0xfe, 0xe1, 0xbb, 0x8a, 0xde, 0x34, 0xd0, 0x99, 0xd8, 0xe7, 0x20, 0x58, 0x9c, 0x03, 0x8d,
	0x99, 0x02, 0x9a, 0xac, 0xa4, 0x04, 0x91, 0x94, 0xb4, 0x60, 0x99, 0x84, 0x39, 0x35, 0xc7, 0xa2,
	0xc8, 0xcb, 0xde, 0x7f, 0x8f, 0xf8, 0x9d, 0x99, 0x57, 0x87, 0x23, 0xa6, 0x60, 0xd4, 0x44, 0xa7,
	0x26, 0x39, 0x45, 0xcc, 0x27, 0x22, 0x2f, 0xaf, 0x5b, 0xaf, 0x6f, 0x03, 0x2b, 0x8a, 0xd6, 0x3b,
	0x97, 0x6c, 0x76, 0x2e, 0xf9, 0xde, 0xb9, 0xe4, 0x65, 0xef, 0x5a, 0x9b, 0xbd, 0x6b, 0x7d, 0xee,
	0x5d, 0xeb, 0xc1, 0x4f, 0x33, 0xbd, 0x58, 0xc5, 0x41, 0x82, 0x3c, 0xac, 0x2a, 0x1d, 0x0a, 0xd0,
	0x4f, 0x28, 0x97, 0x66, 0x08, 0x9f, 0x4d, 0xf1, 0xba, 0x2c, 0x40, 0xc5, 0x6d, 0xd3, 0xde, 0xe5,
	0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x3f, 0x71, 0x2d, 0x45, 0x91, 0x01, 0x00, 0x00,
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
	if m.EnableBaseCurrencyPairedPoolOnly {
		i--
		if m.EnableBaseCurrencyPairedPoolOnly {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x18
	}
	if m.SlippageTrackDuration != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.SlippageTrackDuration))
		i--
		dAtA[i] = 0x10
	}
	{
		size := m.PoolCreationFee.Size()
		i -= size
		if _, err := m.PoolCreationFee.MarshalTo(dAtA[i:]); err != nil {
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
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.PoolCreationFee.Size()
	n += 1 + l + sovParams(uint64(l))
	if m.SlippageTrackDuration != 0 {
		n += 1 + sovParams(uint64(m.SlippageTrackDuration))
	}
	if m.EnableBaseCurrencyPairedPoolOnly {
		n += 2
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
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
				return fmt.Errorf("proto: wrong wireType = %d for field PoolCreationFee", wireType)
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
			if err := m.PoolCreationFee.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlippageTrackDuration", wireType)
			}
			m.SlippageTrackDuration = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SlippageTrackDuration |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnableBaseCurrencyPairedPoolOnly", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
			m.EnableBaseCurrencyPairedPoolOnly = bool(v != 0)
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
