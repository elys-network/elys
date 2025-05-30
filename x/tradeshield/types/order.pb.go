// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/tradeshield/order.proto

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

// Spot order
type SpotOrderType int32

const (
	SpotOrderType_STOPLOSS  SpotOrderType = 0
	SpotOrderType_LIMITSELL SpotOrderType = 1
	SpotOrderType_LIMITBUY  SpotOrderType = 2
	SpotOrderType_MARKETBUY SpotOrderType = 3
)

var SpotOrderType_name = map[int32]string{
	0: "STOPLOSS",
	1: "LIMITSELL",
	2: "LIMITBUY",
	3: "MARKETBUY",
}

var SpotOrderType_value = map[string]int32{
	"STOPLOSS":  0,
	"LIMITSELL": 1,
	"LIMITBUY":  2,
	"MARKETBUY": 3,
}

func (x SpotOrderType) String() string {
	return proto.EnumName(SpotOrderType_name, int32(x))
}

func (SpotOrderType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b257e09947c1671a, []int{0}
}

type PerpetualPosition int32

const (
	PerpetualPosition_UNSPECIFIED PerpetualPosition = 0
	PerpetualPosition_LONG        PerpetualPosition = 1
	PerpetualPosition_SHORT       PerpetualPosition = 2
)

var PerpetualPosition_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "LONG",
	2: "SHORT",
}

var PerpetualPosition_value = map[string]int32{
	"UNSPECIFIED": 0,
	"LONG":        1,
	"SHORT":       2,
}

func (x PerpetualPosition) String() string {
	return proto.EnumName(PerpetualPosition_name, int32(x))
}

func (PerpetualPosition) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b257e09947c1671a, []int{1}
}

type LegacyOrderPriceV1 struct {
	BaseDenom  string                      `protobuf:"bytes,1,opt,name=base_denom,json=baseDenom,proto3" json:"base_denom,omitempty"`
	QuoteDenom string                      `protobuf:"bytes,2,opt,name=quote_denom,json=quoteDenom,proto3" json:"quote_denom,omitempty"`
	Rate       cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=rate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"rate"`
}

func (m *LegacyOrderPriceV1) Reset()         { *m = LegacyOrderPriceV1{} }
func (m *LegacyOrderPriceV1) String() string { return proto.CompactTextString(m) }
func (*LegacyOrderPriceV1) ProtoMessage()    {}
func (*LegacyOrderPriceV1) Descriptor() ([]byte, []int) {
	return fileDescriptor_b257e09947c1671a, []int{0}
}
func (m *LegacyOrderPriceV1) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LegacyOrderPriceV1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LegacyOrderPriceV1.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LegacyOrderPriceV1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LegacyOrderPriceV1.Merge(m, src)
}
func (m *LegacyOrderPriceV1) XXX_Size() int {
	return m.Size()
}
func (m *LegacyOrderPriceV1) XXX_DiscardUnknown() {
	xxx_messageInfo_LegacyOrderPriceV1.DiscardUnknown(m)
}

var xxx_messageInfo_LegacyOrderPriceV1 proto.InternalMessageInfo

func (m *LegacyOrderPriceV1) GetBaseDenom() string {
	if m != nil {
		return m.BaseDenom
	}
	return ""
}

func (m *LegacyOrderPriceV1) GetQuoteDenom() string {
	if m != nil {
		return m.QuoteDenom
	}
	return ""
}

type LegacyTriggerPriceV1 struct {
	TradingAssetDenom string                      `protobuf:"bytes,1,opt,name=trading_asset_denom,json=tradingAssetDenom,proto3" json:"trading_asset_denom,omitempty"`
	Rate              cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=rate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"rate"`
}

func (m *LegacyTriggerPriceV1) Reset()         { *m = LegacyTriggerPriceV1{} }
func (m *LegacyTriggerPriceV1) String() string { return proto.CompactTextString(m) }
func (*LegacyTriggerPriceV1) ProtoMessage()    {}
func (*LegacyTriggerPriceV1) Descriptor() ([]byte, []int) {
	return fileDescriptor_b257e09947c1671a, []int{1}
}
func (m *LegacyTriggerPriceV1) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LegacyTriggerPriceV1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LegacyTriggerPriceV1.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LegacyTriggerPriceV1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LegacyTriggerPriceV1.Merge(m, src)
}
func (m *LegacyTriggerPriceV1) XXX_Size() int {
	return m.Size()
}
func (m *LegacyTriggerPriceV1) XXX_DiscardUnknown() {
	xxx_messageInfo_LegacyTriggerPriceV1.DiscardUnknown(m)
}

var xxx_messageInfo_LegacyTriggerPriceV1 proto.InternalMessageInfo

func (m *LegacyTriggerPriceV1) GetTradingAssetDenom() string {
	if m != nil {
		return m.TradingAssetDenom
	}
	return ""
}

func init() {
	proto.RegisterEnum("elys.tradeshield.SpotOrderType", SpotOrderType_name, SpotOrderType_value)
	proto.RegisterEnum("elys.tradeshield.PerpetualPosition", PerpetualPosition_name, PerpetualPosition_value)
	proto.RegisterType((*LegacyOrderPriceV1)(nil), "elys.tradeshield.LegacyOrderPriceV1")
	proto.RegisterType((*LegacyTriggerPriceV1)(nil), "elys.tradeshield.LegacyTriggerPriceV1")
}

func init() { proto.RegisterFile("elys/tradeshield/order.proto", fileDescriptor_b257e09947c1671a) }

var fileDescriptor_b257e09947c1671a = []byte{
	// 444 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x92, 0xc1, 0x6e, 0xd3, 0x4e,
	0x10, 0xc6, 0xbd, 0x69, 0xff, 0x7f, 0x35, 0x5b, 0x2a, 0xdc, 0xa5, 0x87, 0x52, 0xc0, 0x41, 0x3d,
	0xa1, 0x4a, 0xf5, 0x2a, 0xaa, 0x84, 0xc4, 0xb1, 0x21, 0x06, 0x2c, 0x9c, 0xda, 0x8a, 0x5d, 0x24,
	0xb8, 0x44, 0x8e, 0x3d, 0x72, 0x56, 0x4d, 0xbc, 0x66, 0x77, 0x53, 0xc8, 0x43, 0x20, 0x71, 0xe7,
	0x35, 0x78, 0x88, 0x1e, 0x2b, 0x4e, 0x88, 0x43, 0x85, 0x92, 0x17, 0x41, 0xbb, 0x1b, 0x04, 0xbd,
	0x73, 0xdb, 0x6f, 0x7e, 0x33, 0x3b, 0x9f, 0x46, 0x1f, 0x7e, 0x08, 0xd3, 0x85, 0xa4, 0x4a, 0xe4,
	0x25, 0xc8, 0x09, 0x83, 0x69, 0x49, 0xb9, 0x28, 0x41, 0xf8, 0x8d, 0xe0, 0x8a, 0x13, 0x57, 0x53,
	0xff, 0x2f, 0x7a, 0xb0, 0x57, 0xf1, 0x8a, 0x1b, 0x48, 0xf5, 0xcb, 0xf6, 0x1d, 0x78, 0x05, 0x97,
	0x33, 0x2e, 0xe9, 0x38, 0x97, 0x40, 0x2f, 0xbb, 0x63, 0x50, 0x79, 0x97, 0x16, 0x9c, 0xd5, 0x6b,
	0x7e, 0xdf, 0xf2, 0x91, 0x1d, 0xb4, 0xc2, 0xa2, 0xc3, 0x2f, 0x08, 0x93, 0x08, 0xaa, 0xbc, 0x58,
	0xc4, 0x7a, 0x71, 0x22, 0x58, 0x01, 0x6f, 0xba, 0xe4, 0x11, 0xc6, 0xfa, 0xb3, 0x51, 0x09, 0x35,
	0x9f, 0xed, 0xa3, 0xc7, 0xe8, 0x49, 0x7b, 0xd8, 0xd6, 0x95, 0xbe, 0x2e, 0x90, 0x0e, 0xde, 0x7e,
	0x3f, 0xe7, 0xea, 0x37, 0x6f, 0x19, 0x8e, 0x4d, 0xc9, 0x36, 0x04, 0x78, 0x53, 0xe4, 0x0a, 0xf6,
	0x37, 0x34, 0xe9, 0x75, 0xaf, 0x6e, 0x3a, 0xce, 0x8f, 0x9b, 0xce, 0x03, 0xbb, 0x5a, 0x96, 0x17,
	0x3e, 0xe3, 0x74, 0x96, 0xab, 0x89, 0x6f, 0x77, 0xf7, 0xa1, 0xf8, 0xf6, 0xf5, 0x18, 0xaf, 0x9d,
	0xf5, 0xa1, 0x18, 0x9a, 0xf1, 0xc3, 0x4f, 0x08, 0xef, 0xd9, 0x8e, 0x4c, 0xb0, 0xaa, 0xfa, 0xe3,
	0xcf, 0xc7, 0xf7, 0xf4, 0x59, 0x58, 0x5d, 0x8d, 0x72, 0x29, 0x41, 0xdd, 0x32, 0xba, 0xbb, 0x46,
	0xa7, 0x9a, 0xfc, 0x4b, 0x3f, 0x47, 0x21, 0xde, 0x49, 0x1b, 0xae, 0xcc, 0xa9, 0xb2, 0x45, 0x03,
	0xe4, 0x0e, 0xde, 0x4a, 0xb3, 0x38, 0x89, 0xe2, 0x34, 0x75, 0x1d, 0xb2, 0x83, 0xdb, 0x51, 0x38,
	0x08, 0xb3, 0x34, 0x88, 0x22, 0x17, 0x69, 0x68, 0x64, 0xef, 0xfc, 0xad, 0xdb, 0xd2, 0x70, 0x70,
	0x3a, 0x7c, 0x1d, 0x18, 0xb9, 0x71, 0xf4, 0x0c, 0xef, 0x26, 0x20, 0x1a, 0x50, 0xf3, 0x7c, 0x9a,
	0x70, 0xc9, 0x14, 0xe3, 0x35, 0xb9, 0x8b, 0xb7, 0xcf, 0xcf, 0xd2, 0x24, 0x78, 0x1e, 0xbe, 0x08,
	0x83, 0xbe, 0xeb, 0x90, 0x2d, 0xbc, 0x19, 0xc5, 0x67, 0x2f, 0x5d, 0x44, 0xda, 0xf8, 0xbf, 0xf4,
	0x55, 0x3c, 0xcc, 0xdc, 0x56, 0x6f, 0x70, 0xb5, 0xf4, 0xd0, 0xf5, 0xd2, 0x43, 0x3f, 0x97, 0x1e,
	0xfa, 0xbc, 0xf2, 0x9c, 0xeb, 0x95, 0xe7, 0x7c, 0x5f, 0x79, 0xce, 0xbb, 0x93, 0x8a, 0xa9, 0xc9,
	0x7c, 0xec, 0x17, 0x7c, 0x46, 0x75, 0x76, 0x8e, 0x6b, 0x50, 0x1f, 0xb8, 0xb8, 0x30, 0x82, 0x5e,
	0x3e, 0xa5, 0x1f, 0x6f, 0x65, 0x4d, 0x2d, 0x1a, 0x90, 0xe3, 0xff, 0x4d, 0x12, 0x4e, 0x7e, 0x05,
	0x00, 0x00, 0xff, 0xff, 0x40, 0xc9, 0xfe, 0x9a, 0x8c, 0x02, 0x00, 0x00,
}

func (m *LegacyOrderPriceV1) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LegacyOrderPriceV1) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LegacyOrderPriceV1) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Rate.Size()
		i -= size
		if _, err := m.Rate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOrder(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.QuoteDenom) > 0 {
		i -= len(m.QuoteDenom)
		copy(dAtA[i:], m.QuoteDenom)
		i = encodeVarintOrder(dAtA, i, uint64(len(m.QuoteDenom)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.BaseDenom) > 0 {
		i -= len(m.BaseDenom)
		copy(dAtA[i:], m.BaseDenom)
		i = encodeVarintOrder(dAtA, i, uint64(len(m.BaseDenom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *LegacyTriggerPriceV1) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LegacyTriggerPriceV1) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LegacyTriggerPriceV1) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Rate.Size()
		i -= size
		if _, err := m.Rate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOrder(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.TradingAssetDenom) > 0 {
		i -= len(m.TradingAssetDenom)
		copy(dAtA[i:], m.TradingAssetDenom)
		i = encodeVarintOrder(dAtA, i, uint64(len(m.TradingAssetDenom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintOrder(dAtA []byte, offset int, v uint64) int {
	offset -= sovOrder(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *LegacyOrderPriceV1) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.BaseDenom)
	if l > 0 {
		n += 1 + l + sovOrder(uint64(l))
	}
	l = len(m.QuoteDenom)
	if l > 0 {
		n += 1 + l + sovOrder(uint64(l))
	}
	l = m.Rate.Size()
	n += 1 + l + sovOrder(uint64(l))
	return n
}

func (m *LegacyTriggerPriceV1) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.TradingAssetDenom)
	if l > 0 {
		n += 1 + l + sovOrder(uint64(l))
	}
	l = m.Rate.Size()
	n += 1 + l + sovOrder(uint64(l))
	return n
}

func sovOrder(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozOrder(x uint64) (n int) {
	return sovOrder(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LegacyOrderPriceV1) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOrder
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
			return fmt.Errorf("proto: LegacyOrderPriceV1: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LegacyOrderPriceV1: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrder
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
				return ErrInvalidLengthOrder
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BaseDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QuoteDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrder
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
				return ErrInvalidLengthOrder
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.QuoteDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrder
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
				return ErrInvalidLengthOrder
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Rate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOrder(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthOrder
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
func (m *LegacyTriggerPriceV1) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOrder
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
			return fmt.Errorf("proto: LegacyTriggerPriceV1: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LegacyTriggerPriceV1: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TradingAssetDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrder
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
				return ErrInvalidLengthOrder
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TradingAssetDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Rate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrder
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
				return ErrInvalidLengthOrder
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOrder
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Rate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOrder(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthOrder
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
func skipOrder(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowOrder
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
					return 0, ErrIntOverflowOrder
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
					return 0, ErrIntOverflowOrder
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
				return 0, ErrInvalidLengthOrder
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupOrder
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthOrder
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthOrder        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowOrder          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupOrder = fmt.Errorf("proto: unexpected end of group")
)
