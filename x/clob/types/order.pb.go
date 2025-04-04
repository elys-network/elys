// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/clob/order.proto

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

type OrderType int32

const (
	OrderType_ORDER_TYPE_UNSPECIFIED OrderType = 0
	OrderType_ORDER_TYPE_LIMIT_BUY   OrderType = 1
	OrderType_ORDER_TYPE_LIMIT_SELL  OrderType = 2
	OrderType_ORDER_TYPE_MARKET_BUY  OrderType = 3
	OrderType_ORDER_TYPE_MARKET_SELL OrderType = 4
)

var OrderType_name = map[int32]string{
	0: "ORDER_TYPE_UNSPECIFIED",
	1: "ORDER_TYPE_LIMIT_BUY",
	2: "ORDER_TYPE_LIMIT_SELL",
	3: "ORDER_TYPE_MARKET_BUY",
	4: "ORDER_TYPE_MARKET_SELL",
}

var OrderType_value = map[string]int32{
	"ORDER_TYPE_UNSPECIFIED": 0,
	"ORDER_TYPE_LIMIT_BUY":   1,
	"ORDER_TYPE_LIMIT_SELL":  2,
	"ORDER_TYPE_MARKET_BUY":  3,
	"ORDER_TYPE_MARKET_SELL": 4,
}

func (x OrderType) String() string {
	return proto.EnumName(OrderType_name, int32(x))
}

func (OrderType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_b1c73c2a583cf493, []int{0}
}

// key = market_id + is_long + price + block_height
type PerpetualOrder struct {
	MarketId    uint64                      `protobuf:"varint,1,opt,name=market_id,json=marketId,proto3" json:"market_id,omitempty"`
	OrderType   OrderType                   `protobuf:"varint,2,opt,name=order_type,json=orderType,proto3,enum=elys.clob.OrderType" json:"order_type,omitempty"`
	Price       cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=price,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"price"`
	BlockHeight uint64                      `protobuf:"varint,4,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	Owner       string                      `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
	Amount      cosmossdk_io_math.Int       `protobuf:"bytes,6,opt,name=amount,proto3,customtype=cosmossdk.io/math.Int" json:"amount"`
	Filled      cosmossdk_io_math.Int       `protobuf:"bytes,7,opt,name=filled,proto3,customtype=cosmossdk.io/math.Int" json:"filled"`
}

func (m *PerpetualOrder) Reset()         { *m = PerpetualOrder{} }
func (m *PerpetualOrder) String() string { return proto.CompactTextString(m) }
func (*PerpetualOrder) ProtoMessage()    {}
func (*PerpetualOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_b1c73c2a583cf493, []int{0}
}
func (m *PerpetualOrder) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PerpetualOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PerpetualOrder.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PerpetualOrder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PerpetualOrder.Merge(m, src)
}
func (m *PerpetualOrder) XXX_Size() int {
	return m.Size()
}
func (m *PerpetualOrder) XXX_DiscardUnknown() {
	xxx_messageInfo_PerpetualOrder.DiscardUnknown(m)
}

var xxx_messageInfo_PerpetualOrder proto.InternalMessageInfo

func (m *PerpetualOrder) GetMarketId() uint64 {
	if m != nil {
		return m.MarketId
	}
	return 0
}

func (m *PerpetualOrder) GetOrderType() OrderType {
	if m != nil {
		return m.OrderType
	}
	return OrderType_ORDER_TYPE_UNSPECIFIED
}

func (m *PerpetualOrder) GetBlockHeight() uint64 {
	if m != nil {
		return m.BlockHeight
	}
	return 0
}

func (m *PerpetualOrder) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func init() {
	proto.RegisterEnum("elys.clob.OrderType", OrderType_name, OrderType_value)
	proto.RegisterType((*PerpetualOrder)(nil), "elys.clob.PerpetualOrder")
}

func init() { proto.RegisterFile("elys/clob/order.proto", fileDescriptor_b1c73c2a583cf493) }

var fileDescriptor_b1c73c2a583cf493 = []byte{
	// 460 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x52, 0xcd, 0x6e, 0xd3, 0x40,
	0x18, 0xf4, 0x36, 0x3f, 0xe0, 0x05, 0x55, 0xd1, 0x2a, 0x41, 0xdb, 0x54, 0x72, 0x03, 0xa7, 0x00,
	0x8a, 0x2d, 0xe8, 0x13, 0x34, 0x89, 0x01, 0x8b, 0x94, 0x46, 0x4e, 0x7a, 0x28, 0x17, 0xcb, 0x59,
	0x2f, 0x8e, 0x15, 0xdb, 0x6b, 0xad, 0x37, 0x2a, 0x7e, 0x03, 0x8e, 0x48, 0xbc, 0x4a, 0x1f, 0xa2,
	0xc7, 0xaa, 0x27, 0xc4, 0xa1, 0x42, 0xc9, 0x8b, 0x20, 0xef, 0x9a, 0xaa, 0x6a, 0x6f, 0xbd, 0xed,
	0x37, 0xf3, 0xcd, 0x7c, 0xa3, 0xd5, 0xc0, 0x0e, 0x8d, 0x8b, 0xdc, 0x22, 0x31, 0x5b, 0x58, 0x8c,
	0x07, 0x94, 0x9b, 0x19, 0x67, 0x82, 0x21, 0xbd, 0x84, 0xcd, 0x12, 0xee, 0xb6, 0x43, 0x16, 0x32,
	0x89, 0x5a, 0xe5, 0x4b, 0x2d, 0x74, 0xf7, 0x08, 0xcb, 0x13, 0x96, 0x7b, 0x8a, 0x50, 0x83, 0xa2,
	0x5e, 0xfd, 0xa8, 0xc1, 0xdd, 0x29, 0xe5, 0x19, 0x15, 0x6b, 0x3f, 0x3e, 0x29, 0x4d, 0xd1, 0x3e,
	0xd4, 0x13, 0x9f, 0xaf, 0xa8, 0xf0, 0xa2, 0x00, 0x83, 0x1e, 0xe8, 0xd7, 0xdd, 0xa7, 0x0a, 0x70,
	0x02, 0x74, 0x08, 0xa1, 0x3c, 0xed, 0x89, 0x22, 0xa3, 0x78, 0xa7, 0x07, 0xfa, 0xbb, 0xef, 0xdb,
	0xe6, 0x6d, 0x00, 0x53, 0x5a, 0xcc, 0x8b, 0x8c, 0xba, 0x3a, 0xfb, 0xff, 0x44, 0x1f, 0x61, 0x23,
	0xe3, 0x11, 0xa1, 0xb8, 0xd6, 0x03, 0x7d, 0x7d, 0xf8, 0xee, 0xf2, 0xe6, 0x40, 0xfb, 0x73, 0x73,
	0xb0, 0xaf, 0x92, 0xe4, 0xc1, 0xca, 0x8c, 0x98, 0x95, 0xf8, 0x62, 0x69, 0x4e, 0x68, 0xe8, 0x93,
	0x62, 0x4c, 0xc9, 0xf5, 0xc5, 0x00, 0x56, 0x41, 0xc7, 0x94, 0xb8, 0x4a, 0x8f, 0x5e, 0xc2, 0xe7,
	0x8b, 0x98, 0x91, 0x95, 0xb7, 0xa4, 0x51, 0xb8, 0x14, 0xb8, 0x2e, 0xd3, 0x3d, 0x93, 0xd8, 0x27,
	0x09, 0x21, 0x13, 0x36, 0xd8, 0x79, 0x4a, 0x39, 0x6e, 0xc8, 0x5b, 0xf8, 0xfa, 0x62, 0xd0, 0xae,
	0x8c, 0x8e, 0x82, 0x80, 0xd3, 0x3c, 0x9f, 0x09, 0x1e, 0xa5, 0xa1, 0xab, 0xd6, 0xd0, 0x08, 0x36,
	0xfd, 0x84, 0xad, 0x53, 0x81, 0x9b, 0x52, 0xf0, 0xb6, 0x0a, 0xd7, 0x79, 0x18, 0xce, 0x49, 0xc5,
	0x9d, 0x58, 0x4e, 0x2a, 0xdc, 0x4a, 0x5a, 0x9a, 0x7c, 0x8b, 0xe2, 0x98, 0x06, 0xf8, 0xc9, 0x23,
	0x4c, 0x94, 0xf4, 0xcd, 0x2f, 0x00, 0xf5, 0xdb, 0xef, 0x43, 0x5d, 0xf8, 0xe2, 0xc4, 0x1d, 0xdb,
	0xae, 0x37, 0x3f, 0x9b, 0xda, 0xde, 0xe9, 0x97, 0xd9, 0xd4, 0x1e, 0x39, 0x1f, 0x1c, 0x7b, 0xdc,
	0xd2, 0x10, 0x86, 0xed, 0x3b, 0xdc, 0xc4, 0x39, 0x76, 0xe6, 0xde, 0xf0, 0xf4, 0xac, 0x05, 0xd0,
	0x1e, 0xec, 0x3c, 0x60, 0x66, 0xf6, 0x64, 0xd2, 0xda, 0xb9, 0x47, 0x1d, 0x1f, 0xb9, 0x9f, 0x6d,
	0xa5, 0xaa, 0xdd, 0xbb, 0x55, 0x51, 0x52, 0x56, 0x1f, 0x8e, 0x2e, 0x37, 0x06, 0xb8, 0xda, 0x18,
	0xe0, 0xef, 0xc6, 0x00, 0x3f, 0xb7, 0x86, 0x76, 0xb5, 0x35, 0xb4, 0xdf, 0x5b, 0x43, 0xfb, 0xfa,
	0x3a, 0x8c, 0xc4, 0x72, 0xbd, 0x30, 0x09, 0x4b, 0xac, 0xb2, 0x00, 0x83, 0x94, 0x8a, 0x73, 0xc6,
	0x57, 0x72, 0xb0, 0xbe, 0xab, 0x9e, 0x96, 0x35, 0xc9, 0x17, 0x4d, 0x59, 0xb6, 0xc3, 0x7f, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x64, 0xd1, 0x49, 0xbd, 0xc1, 0x02, 0x00, 0x00,
}

func (m *PerpetualOrder) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PerpetualOrder) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PerpetualOrder) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Filled.Size()
		i -= size
		if _, err := m.Filled.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOrder(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOrder(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintOrder(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x2a
	}
	if m.BlockHeight != 0 {
		i = encodeVarintOrder(dAtA, i, uint64(m.BlockHeight))
		i--
		dAtA[i] = 0x20
	}
	{
		size := m.Price.Size()
		i -= size
		if _, err := m.Price.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOrder(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if m.OrderType != 0 {
		i = encodeVarintOrder(dAtA, i, uint64(m.OrderType))
		i--
		dAtA[i] = 0x10
	}
	if m.MarketId != 0 {
		i = encodeVarintOrder(dAtA, i, uint64(m.MarketId))
		i--
		dAtA[i] = 0x8
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
func (m *PerpetualOrder) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MarketId != 0 {
		n += 1 + sovOrder(uint64(m.MarketId))
	}
	if m.OrderType != 0 {
		n += 1 + sovOrder(uint64(m.OrderType))
	}
	l = m.Price.Size()
	n += 1 + l + sovOrder(uint64(l))
	if m.BlockHeight != 0 {
		n += 1 + sovOrder(uint64(m.BlockHeight))
	}
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovOrder(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovOrder(uint64(l))
	l = m.Filled.Size()
	n += 1 + l + sovOrder(uint64(l))
	return n
}

func sovOrder(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozOrder(x uint64) (n int) {
	return sovOrder(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PerpetualOrder) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: PerpetualOrder: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PerpetualOrder: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MarketId", wireType)
			}
			m.MarketId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrder
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MarketId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OrderType", wireType)
			}
			m.OrderType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrder
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OrderType |= OrderType(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Price", wireType)
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
			if err := m.Price.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHeight", wireType)
			}
			m.BlockHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrder
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.BlockHeight |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
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
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
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
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Filled", wireType)
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
			if err := m.Filled.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
