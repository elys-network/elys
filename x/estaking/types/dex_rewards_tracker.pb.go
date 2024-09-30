// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/estaking/dex_rewards_tracker.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
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

// DexRewardsTracker is used for tracking rewards for stakers and LPs, all
// amount here is in USDC
type DexRewardsTracker struct {
	// Number of blocks since start of epoch (distribution epoch)
	NumBlocks int64 `protobuf:"varint,1,opt,name=num_blocks,json=numBlocks,proto3" json:"num_blocks,omitempty"`
	// Accumulated amount at distribution epoch - recalculated at every
	// distribution epoch
	Amount cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=amount,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"amount"`
}

func (m *DexRewardsTracker) Reset()         { *m = DexRewardsTracker{} }
func (m *DexRewardsTracker) String() string { return proto.CompactTextString(m) }
func (*DexRewardsTracker) ProtoMessage()    {}
func (*DexRewardsTracker) Descriptor() ([]byte, []int) {
	return fileDescriptor_061875a2058b444a, []int{0}
}
func (m *DexRewardsTracker) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DexRewardsTracker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DexRewardsTracker.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DexRewardsTracker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DexRewardsTracker.Merge(m, src)
}
func (m *DexRewardsTracker) XXX_Size() int {
	return m.Size()
}
func (m *DexRewardsTracker) XXX_DiscardUnknown() {
	xxx_messageInfo_DexRewardsTracker.DiscardUnknown(m)
}

var xxx_messageInfo_DexRewardsTracker proto.InternalMessageInfo

func (m *DexRewardsTracker) GetNumBlocks() int64 {
	if m != nil {
		return m.NumBlocks
	}
	return 0
}

type LegacyDexRewardsTracker struct {
	// Number of blocks since start of epoch (distribution epoch)
	NumBlocks cosmossdk_io_math.Int `protobuf:"bytes,1,opt,name=num_blocks,json=numBlocks,proto3,customtype=cosmossdk.io/math.Int" json:"num_blocks"`
	// Accumulated amount at distribution epoch - recalculated at every
	// distribution epoch
	Amount cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=amount,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"amount"`
}

func (m *LegacyDexRewardsTracker) Reset()         { *m = LegacyDexRewardsTracker{} }
func (m *LegacyDexRewardsTracker) String() string { return proto.CompactTextString(m) }
func (*LegacyDexRewardsTracker) ProtoMessage()    {}
func (*LegacyDexRewardsTracker) Descriptor() ([]byte, []int) {
	return fileDescriptor_061875a2058b444a, []int{1}
}
func (m *LegacyDexRewardsTracker) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LegacyDexRewardsTracker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LegacyDexRewardsTracker.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LegacyDexRewardsTracker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LegacyDexRewardsTracker.Merge(m, src)
}
func (m *LegacyDexRewardsTracker) XXX_Size() int {
	return m.Size()
}
func (m *LegacyDexRewardsTracker) XXX_DiscardUnknown() {
	xxx_messageInfo_LegacyDexRewardsTracker.DiscardUnknown(m)
}

var xxx_messageInfo_LegacyDexRewardsTracker proto.InternalMessageInfo

func init() {
	proto.RegisterType((*DexRewardsTracker)(nil), "elys.estaking.DexRewardsTracker")
	proto.RegisterType((*LegacyDexRewardsTracker)(nil), "elys.estaking.LegacyDexRewardsTracker")
}

func init() {
	proto.RegisterFile("elys/estaking/dex_rewards_tracker.proto", fileDescriptor_061875a2058b444a)
}

var fileDescriptor_061875a2058b444a = []byte{
	// 280 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4f, 0xcd, 0xa9, 0x2c,
	0xd6, 0x4f, 0x2d, 0x2e, 0x49, 0xcc, 0xce, 0xcc, 0x4b, 0xd7, 0x4f, 0x49, 0xad, 0x88, 0x2f, 0x4a,
	0x2d, 0x4f, 0x2c, 0x4a, 0x29, 0x8e, 0x2f, 0x29, 0x4a, 0x4c, 0xce, 0x4e, 0x2d, 0xd2, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x05, 0x29, 0xd4, 0x83, 0x29, 0x94, 0x12, 0x49, 0xcf, 0x4f, 0xcf,
	0x07, 0xcb, 0xe8, 0x83, 0x58, 0x10, 0x45, 0x4a, 0xf9, 0x5c, 0x82, 0x2e, 0xa9, 0x15, 0x41, 0x10,
	0x03, 0x42, 0x20, 0xfa, 0x85, 0x64, 0xb9, 0xb8, 0xf2, 0x4a, 0x73, 0xe3, 0x93, 0x72, 0xf2, 0x93,
	0xb3, 0x8b, 0x25, 0x18, 0x15, 0x18, 0x35, 0x98, 0x83, 0x38, 0xf3, 0x4a, 0x73, 0x9d, 0xc0, 0x02,
	0x42, 0xd6, 0x5c, 0x6c, 0x89, 0xb9, 0xf9, 0xa5, 0x79, 0x25, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0x9c,
	0x4e, 0xca, 0x27, 0xee, 0xc9, 0x33, 0xdc, 0xba, 0x27, 0x2f, 0x9d, 0x9c, 0x5f, 0x9c, 0x9b, 0x5f,
	0x5c, 0x9c, 0x92, 0xad, 0x97, 0x99, 0xaf, 0x9f, 0x9b, 0x58, 0x92, 0xa1, 0xe7, 0x93, 0x9a, 0x9e,
	0x98, 0x5c, 0xe9, 0x92, 0x9a, 0x1c, 0x04, 0xd5, 0xa2, 0x34, 0x85, 0x91, 0x4b, 0x1c, 0x26, 0x8a,
	0x6e, 0xaf, 0x0d, 0x86, 0xbd, 0x9c, 0x4e, 0xb2, 0x50, 0xc3, 0x45, 0x31, 0x0d, 0xf7, 0xcc, 0x2b,
	0xa1, 0x96, 0xb3, 0x9c, 0xbc, 0x57, 0x3c, 0x92, 0x63, 0x3c, 0xf1, 0x48, 0x8e, 0xf1, 0xc2, 0x23,
	0x39, 0xc6, 0x07, 0x8f, 0xe4, 0x18, 0x27, 0x3c, 0x96, 0x63, 0xb8, 0xf0, 0x58, 0x8e, 0xe1, 0xc6,
	0x63, 0x39, 0x86, 0x28, 0xdd, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0x7d,
	0x50, 0xa8, 0xea, 0xe6, 0xa5, 0x96, 0x94, 0xe7, 0x17, 0x65, 0x83, 0x39, 0xfa, 0x15, 0x88, 0xd8,
	0x28, 0xa9, 0x2c, 0x48, 0x2d, 0x4e, 0x62, 0x03, 0x87, 0xad, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff,
	0x61, 0x8f, 0xd7, 0x1c, 0xab, 0x01, 0x00, 0x00,
}

func (this *DexRewardsTracker) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*DexRewardsTracker)
	if !ok {
		that2, ok := that.(DexRewardsTracker)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.NumBlocks != that1.NumBlocks {
		return false
	}
	if !this.Amount.Equal(that1.Amount) {
		return false
	}
	return true
}
func (this *LegacyDexRewardsTracker) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*LegacyDexRewardsTracker)
	if !ok {
		that2, ok := that.(LegacyDexRewardsTracker)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.NumBlocks.Equal(that1.NumBlocks) {
		return false
	}
	if !this.Amount.Equal(that1.Amount) {
		return false
	}
	return true
}
func (m *DexRewardsTracker) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DexRewardsTracker) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DexRewardsTracker) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDexRewardsTracker(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.NumBlocks != 0 {
		i = encodeVarintDexRewardsTracker(dAtA, i, uint64(m.NumBlocks))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *LegacyDexRewardsTracker) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LegacyDexRewardsTracker) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LegacyDexRewardsTracker) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDexRewardsTracker(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.NumBlocks.Size()
		i -= size
		if _, err := m.NumBlocks.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintDexRewardsTracker(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintDexRewardsTracker(dAtA []byte, offset int, v uint64) int {
	offset -= sovDexRewardsTracker(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *DexRewardsTracker) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.NumBlocks != 0 {
		n += 1 + sovDexRewardsTracker(uint64(m.NumBlocks))
	}
	l = m.Amount.Size()
	n += 1 + l + sovDexRewardsTracker(uint64(l))
	return n
}

func (m *LegacyDexRewardsTracker) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.NumBlocks.Size()
	n += 1 + l + sovDexRewardsTracker(uint64(l))
	l = m.Amount.Size()
	n += 1 + l + sovDexRewardsTracker(uint64(l))
	return n
}

func sovDexRewardsTracker(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDexRewardsTracker(x uint64) (n int) {
	return sovDexRewardsTracker(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *DexRewardsTracker) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDexRewardsTracker
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
			return fmt.Errorf("proto: DexRewardsTracker: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DexRewardsTracker: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumBlocks", wireType)
			}
			m.NumBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDexRewardsTracker
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumBlocks |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDexRewardsTracker
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
				return ErrInvalidLengthDexRewardsTracker
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDexRewardsTracker
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDexRewardsTracker(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDexRewardsTracker
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
func (m *LegacyDexRewardsTracker) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDexRewardsTracker
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
			return fmt.Errorf("proto: LegacyDexRewardsTracker: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LegacyDexRewardsTracker: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumBlocks", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDexRewardsTracker
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
				return ErrInvalidLengthDexRewardsTracker
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDexRewardsTracker
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.NumBlocks.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDexRewardsTracker
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
				return ErrInvalidLengthDexRewardsTracker
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthDexRewardsTracker
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipDexRewardsTracker(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDexRewardsTracker
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
func skipDexRewardsTracker(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDexRewardsTracker
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
					return 0, ErrIntOverflowDexRewardsTracker
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
					return 0, ErrIntOverflowDexRewardsTracker
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
				return 0, ErrInvalidLengthDexRewardsTracker
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDexRewardsTracker
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDexRewardsTracker
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDexRewardsTracker        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDexRewardsTracker          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDexRewardsTracker = fmt.Errorf("proto: unexpected end of group")
)
