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

type OrderPrice struct {
	BaseDenom  string                      `protobuf:"bytes,1,opt,name=base_denom,json=baseDenom,proto3" json:"base_denom,omitempty"`
	QuoteDenom string                      `protobuf:"bytes,2,opt,name=quote_denom,json=quoteDenom,proto3" json:"quote_denom,omitempty"`
	Rate       cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=rate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"rate"`
}

func (m *OrderPrice) Reset()         { *m = OrderPrice{} }
func (m *OrderPrice) String() string { return proto.CompactTextString(m) }
func (*OrderPrice) ProtoMessage()    {}
func (*OrderPrice) Descriptor() ([]byte, []int) {
	return fileDescriptor_b257e09947c1671a, []int{0}
}
func (m *OrderPrice) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OrderPrice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OrderPrice.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OrderPrice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderPrice.Merge(m, src)
}
func (m *OrderPrice) XXX_Size() int {
	return m.Size()
}
func (m *OrderPrice) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderPrice.DiscardUnknown(m)
}

var xxx_messageInfo_OrderPrice proto.InternalMessageInfo

func (m *OrderPrice) GetBaseDenom() string {
	if m != nil {
		return m.BaseDenom
	}
	return ""
}

func (m *OrderPrice) GetQuoteDenom() string {
	if m != nil {
		return m.QuoteDenom
	}
	return ""
}

type TriggerPrice struct {
	TradingAssetDenom string                      `protobuf:"bytes,1,opt,name=trading_asset_denom,json=tradingAssetDenom,proto3" json:"trading_asset_denom,omitempty"`
	Rate              cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=rate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"rate"`
}

func (m *TriggerPrice) Reset()         { *m = TriggerPrice{} }
func (m *TriggerPrice) String() string { return proto.CompactTextString(m) }
func (*TriggerPrice) ProtoMessage()    {}
func (*TriggerPrice) Descriptor() ([]byte, []int) {
	return fileDescriptor_b257e09947c1671a, []int{1}
}
func (m *TriggerPrice) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TriggerPrice) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TriggerPrice.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TriggerPrice) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TriggerPrice.Merge(m, src)
}
func (m *TriggerPrice) XXX_Size() int {
	return m.Size()
}
func (m *TriggerPrice) XXX_DiscardUnknown() {
	xxx_messageInfo_TriggerPrice.DiscardUnknown(m)
}

var xxx_messageInfo_TriggerPrice proto.InternalMessageInfo

func (m *TriggerPrice) GetTradingAssetDenom() string {
	if m != nil {
		return m.TradingAssetDenom
	}
	return ""
}

func init() {
	proto.RegisterEnum("elys.tradeshield.SpotOrderType", SpotOrderType_name, SpotOrderType_value)
	proto.RegisterEnum("elys.tradeshield.PerpetualPosition", PerpetualPosition_name, PerpetualPosition_value)
	proto.RegisterType((*OrderPrice)(nil), "elys.tradeshield.OrderPrice")
	proto.RegisterType((*TriggerPrice)(nil), "elys.tradeshield.TriggerPrice")
}

func init() { proto.RegisterFile("elys/tradeshield/order.proto", fileDescriptor_b257e09947c1671a) }

var fileDescriptor_b257e09947c1671a = []byte{
	// 422 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0xed, 0xb4, 0xa0, 0x66, 0xda, 0x0a, 0xd7, 0x70, 0x28, 0x05, 0x1c, 0x54, 0x2e, 0xa8,
	0x12, 0x5e, 0x2a, 0x0e, 0x88, 0x63, 0x43, 0x02, 0x18, 0xdc, 0xda, 0x8a, 0xdd, 0x03, 0x5c, 0x22,
	0xc7, 0x1e, 0x39, 0xab, 0x26, 0x1e, 0xb3, 0xbb, 0xa1, 0xf8, 0x01, 0xb8, 0xf3, 0x58, 0x3d, 0xf6,
	0x88, 0x38, 0x54, 0x28, 0x79, 0x11, 0xb4, 0xeb, 0x56, 0xa2, 0x67, 0x6e, 0xfe, 0xff, 0x6f, 0xc6,
	0xff, 0xec, 0xee, 0xc0, 0x63, 0x9c, 0x35, 0x92, 0x29, 0x91, 0x15, 0x28, 0xa7, 0x1c, 0x67, 0x05,
	0x23, 0x51, 0xa0, 0xf0, 0x6b, 0x41, 0x8a, 0x5c, 0x47, 0x53, 0xff, 0x1f, 0xba, 0xf7, 0xa0, 0xa4,
	0x92, 0x0c, 0x64, 0xfa, 0xab, 0xad, 0xdb, 0xf3, 0x72, 0x92, 0x73, 0x92, 0x6c, 0x92, 0x49, 0x64,
	0xdf, 0x0e, 0x27, 0xa8, 0xb2, 0x43, 0x96, 0x13, 0xaf, 0xae, 0xf9, 0xc3, 0x96, 0x8f, 0xdb, 0xc6,
	0x56, 0xb4, 0x68, 0xff, 0x87, 0x0d, 0x10, 0xe9, 0xc8, 0x58, 0xf0, 0x1c, 0xdd, 0x27, 0x00, 0xfa,
	0x27, 0xe3, 0x02, 0x2b, 0x9a, 0xef, 0xda, 0x4f, 0xed, 0xe7, 0xdd, 0x51, 0x57, 0x3b, 0x03, 0x6d,
	0xb8, 0x3d, 0xd8, 0xfc, 0xba, 0x20, 0x75, 0xc3, 0x3b, 0x86, 0x83, 0xb1, 0xda, 0x82, 0xd7, 0xb0,
	0x2e, 0x32, 0x85, 0xbb, 0x6b, 0x9a, 0xf4, 0x9f, 0x5d, 0x5c, 0xf5, 0xac, 0xdf, 0x57, 0xbd, 0x47,
	0x6d, 0xa4, 0x2c, 0xce, 0x7c, 0x4e, 0x6c, 0x9e, 0xa9, 0xa9, 0x1f, 0x62, 0x99, 0xe5, 0xcd, 0x00,
	0xf3, 0x91, 0x69, 0xd8, 0x3f, 0x87, 0xad, 0x54, 0xf0, 0xb2, 0xbc, 0x19, 0xc4, 0x87, 0xfb, 0xfa,
	0xdc, 0xbc, 0x2a, 0xc7, 0x99, 0x94, 0xa8, 0x6e, 0x4d, 0xb4, 0x73, 0x8d, 0x8e, 0x34, 0xf9, 0xbf,
	0xe0, 0x83, 0x00, 0xb6, 0x93, 0x9a, 0x94, 0xb9, 0x83, 0xb4, 0xa9, 0xd1, 0xdd, 0x82, 0x8d, 0x24,
	0x8d, 0xe2, 0x30, 0x4a, 0x12, 0xc7, 0x72, 0xb7, 0xa1, 0x1b, 0x06, 0xc7, 0x41, 0x9a, 0x0c, 0xc3,
	0xd0, 0xb1, 0x35, 0x34, 0xb2, 0x7f, 0xfa, 0xd9, 0xe9, 0x68, 0x78, 0x7c, 0x34, 0xfa, 0x34, 0x34,
	0x72, 0xed, 0xe0, 0x0d, 0xec, 0xc4, 0x28, 0x6a, 0x54, 0x8b, 0x6c, 0x16, 0x93, 0xe4, 0x8a, 0x53,
	0xe5, 0xde, 0x83, 0xcd, 0xd3, 0x93, 0x24, 0x1e, 0xbe, 0x0d, 0xde, 0x05, 0xc3, 0x81, 0x63, 0xb9,
	0x1b, 0xb0, 0x1e, 0x46, 0x27, 0xef, 0x1d, 0xdb, 0xed, 0xc2, 0x9d, 0xe4, 0x43, 0x34, 0x4a, 0x9d,
	0x4e, 0xff, 0xe3, 0xc5, 0xd2, 0xb3, 0x2f, 0x97, 0x9e, 0xfd, 0x67, 0xe9, 0xd9, 0x3f, 0x57, 0x9e,
	0x75, 0xb9, 0xf2, 0xac, 0x5f, 0x2b, 0xcf, 0xfa, 0xf2, 0xb2, 0xe4, 0x6a, 0xba, 0x98, 0xf8, 0x39,
	0xcd, 0x99, 0x5e, 0x87, 0x17, 0x15, 0xaa, 0x73, 0x12, 0x67, 0x46, 0xb0, 0xef, 0xb7, 0x76, 0x47,
	0x35, 0x35, 0xca, 0xc9, 0x5d, 0xf3, 0xb2, 0xaf, 0xfe, 0x06, 0x00, 0x00, 0xff, 0xff, 0xb7, 0xa2,
	0x9c, 0x9a, 0x5c, 0x02, 0x00, 0x00,
}

func (m *OrderPrice) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OrderPrice) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OrderPrice) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *TriggerPrice) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TriggerPrice) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TriggerPrice) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *OrderPrice) Size() (n int) {
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

func (m *TriggerPrice) Size() (n int) {
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
func (m *OrderPrice) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: OrderPrice: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OrderPrice: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *TriggerPrice) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: TriggerPrice: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TriggerPrice: illegal tag %d (wire type %d)", fieldNum, wire)
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
