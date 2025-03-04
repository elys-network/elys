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
	MarketId     uint64                      `protobuf:"varint,1,opt,name=market_id,json=marketId,proto3" json:"market_id,omitempty"`
	OrderType    OrderType                   `protobuf:"varint,2,opt,name=order_type,json=orderType,proto3,enum=elys.clob.OrderType" json:"order_type,omitempty"`
	Price        cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=price,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"price"`
	BlockHeight  uint64                      `protobuf:"varint,4,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	Owner        string                      `protobuf:"bytes,5,opt,name=owner,proto3" json:"owner,omitempty"`
	SubAccountId uint64                      `protobuf:"varint,6,opt,name=sub_account_id,json=subAccountId,proto3" json:"sub_account_id,omitempty"`
	// There are two kinds of selling that can happen. 1. I don't own any futures
	// and I am placing a limit order. This is 0.
	// 2. I own some future already and I am selling it, in this case this is > 0
	// Similarly for buy side
	PerpetualId uint64                      `protobuf:"varint,7,opt,name=perpetual_id,json=perpetualId,proto3" json:"perpetual_id,omitempty"`
	Collateral  cosmossdk_io_math.Int       `protobuf:"bytes,8,opt,name=collateral,proto3,customtype=cosmossdk.io/math.Int" json:"collateral"`
	Leverage    cosmossdk_io_math.LegacyDec `protobuf:"bytes,9,opt,name=leverage,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"leverage"`
	Filled      cosmossdk_io_math.LegacyDec `protobuf:"bytes,10,opt,name=filled,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"filled"`
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

func (m *PerpetualOrder) GetSubAccountId() uint64 {
	if m != nil {
		return m.SubAccountId
	}
	return 0
}

func (m *PerpetualOrder) GetPerpetualId() uint64 {
	if m != nil {
		return m.PerpetualId
	}
	return 0
}

func init() {
	proto.RegisterEnum("elys.clob.OrderType", OrderType_name, OrderType_value)
	proto.RegisterType((*PerpetualOrder)(nil), "elys.clob.PerpetualOrder")
}

func init() { proto.RegisterFile("elys/clob/order.proto", fileDescriptor_b1c73c2a583cf493) }

var fileDescriptor_b1c73c2a583cf493 = []byte{
	// 516 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x53, 0x4d, 0x6e, 0xda, 0x40,
	0x18, 0xc5, 0x09, 0xa1, 0x78, 0x1a, 0x21, 0x34, 0x82, 0xca, 0x21, 0x92, 0x43, 0xab, 0x2e, 0x68,
	0x2b, 0x6c, 0xb5, 0x39, 0x01, 0x04, 0xb7, 0xb5, 0x02, 0x0d, 0x32, 0x64, 0x91, 0x6e, 0x2c, 0x7b,
	0xfc, 0xd5, 0x58, 0x18, 0x8f, 0x35, 0x1e, 0x9a, 0x72, 0x8b, 0x4a, 0xbd, 0x47, 0x57, 0x39, 0x44,
	0x96, 0x51, 0x56, 0x55, 0x17, 0x51, 0x05, 0x17, 0xa9, 0x3c, 0x43, 0x10, 0x4a, 0x76, 0xd9, 0xcd,
	0xbc, 0xf7, 0x7d, 0xef, 0x7d, 0x6f, 0x7e, 0x50, 0x1d, 0xe2, 0x45, 0x66, 0x92, 0x98, 0xfa, 0x26,
	0x65, 0x01, 0x30, 0x23, 0x65, 0x94, 0x53, 0xac, 0xe6, 0xb0, 0x91, 0xc3, 0x8d, 0x5a, 0x48, 0x43,
	0x2a, 0x50, 0x33, 0x5f, 0xc9, 0x82, 0xc6, 0x01, 0xa1, 0xd9, 0x8c, 0x66, 0xae, 0x24, 0xe4, 0x46,
	0x52, 0xaf, 0x7e, 0x17, 0x51, 0x65, 0x08, 0x2c, 0x05, 0x3e, 0xf7, 0xe2, 0xb3, 0x5c, 0x14, 0x1f,
	0x22, 0x75, 0xe6, 0xb1, 0x29, 0x70, 0x37, 0x0a, 0x34, 0xa5, 0xa9, 0xb4, 0x8a, 0x4e, 0x59, 0x02,
	0x76, 0x80, 0x8f, 0x11, 0x12, 0xd6, 0x2e, 0x5f, 0xa4, 0xa0, 0xed, 0x34, 0x95, 0x56, 0xe5, 0x43,
	0xcd, 0xd8, 0x0c, 0x60, 0x08, 0x89, 0xf1, 0x22, 0x05, 0x47, 0xa5, 0xf7, 0x4b, 0xfc, 0x09, 0xed,
	0xa5, 0x2c, 0x22, 0xa0, 0xed, 0x36, 0x95, 0x96, 0xda, 0x7d, 0x7f, 0x7d, 0x77, 0x54, 0xf8, 0x7b,
	0x77, 0x74, 0x28, 0x27, 0xc9, 0x82, 0xa9, 0x11, 0x51, 0x73, 0xe6, 0xf1, 0x89, 0xd1, 0x87, 0xd0,
	0x23, 0x8b, 0x1e, 0x90, 0xdb, 0xab, 0x36, 0x5a, 0x0f, 0xda, 0x03, 0xe2, 0xc8, 0x7e, 0xfc, 0x12,
	0xed, 0xfb, 0x31, 0x25, 0x53, 0x77, 0x02, 0x51, 0x38, 0xe1, 0x5a, 0x51, 0x4c, 0xf7, 0x5c, 0x60,
	0x9f, 0x05, 0x84, 0x0d, 0xb4, 0x47, 0x2f, 0x13, 0x60, 0xda, 0x9e, 0xf0, 0xd2, 0x6e, 0xaf, 0xda,
	0xb5, 0xb5, 0x50, 0x27, 0x08, 0x18, 0x64, 0xd9, 0x88, 0xb3, 0x28, 0x09, 0x1d, 0x59, 0x86, 0x5f,
	0xa3, 0x4a, 0x36, 0xf7, 0x5d, 0x8f, 0x10, 0x3a, 0x4f, 0x44, 0xe4, 0x92, 0x10, 0xdd, 0xcf, 0xe6,
	0x7e, 0x47, 0x82, 0x76, 0x90, 0x1b, 0xa7, 0xf7, 0xa7, 0x94, 0xd7, 0x3c, 0x93, 0xc6, 0x1b, 0xcc,
	0x0e, 0xf0, 0x29, 0x42, 0x84, 0xc6, 0xb1, 0xc7, 0x81, 0x79, 0xb1, 0x56, 0x16, 0xee, 0xef, 0xd6,
	0x49, 0xeb, 0x8f, 0x93, 0xda, 0x09, 0xdf, 0xca, 0x68, 0x27, 0xdc, 0xd9, 0x6a, 0xc7, 0x03, 0x54,
	0x8e, 0xe1, 0x3b, 0x30, 0x2f, 0x04, 0x4d, 0x7d, 0xea, 0xa1, 0x6d, 0x24, 0xb0, 0x8d, 0x4a, 0xdf,
	0xa2, 0x38, 0x86, 0x40, 0x43, 0x4f, 0x15, 0x5b, 0x0b, 0xbc, 0xfd, 0xa5, 0x20, 0x75, 0x73, 0xc9,
	0xb8, 0x81, 0x5e, 0x9c, 0x39, 0x3d, 0xcb, 0x71, 0xc7, 0x17, 0x43, 0xcb, 0x3d, 0xff, 0x32, 0x1a,
	0x5a, 0x27, 0xf6, 0x47, 0xdb, 0xea, 0x55, 0x0b, 0x58, 0x43, 0xb5, 0x2d, 0xae, 0x6f, 0x0f, 0xec,
	0xb1, 0xdb, 0x3d, 0xbf, 0xa8, 0x2a, 0xf8, 0x00, 0xd5, 0x1f, 0x31, 0x23, 0xab, 0xdf, 0xaf, 0xee,
	0x3c, 0xa0, 0x06, 0x1d, 0xe7, 0xd4, 0x92, 0x5d, 0xbb, 0x0f, 0xbc, 0xd6, 0x94, 0x68, 0x2b, 0x76,
	0x4f, 0xae, 0x97, 0xba, 0x72, 0xb3, 0xd4, 0x95, 0x7f, 0x4b, 0x5d, 0xf9, 0xb9, 0xd2, 0x0b, 0x37,
	0x2b, 0xbd, 0xf0, 0x67, 0xa5, 0x17, 0xbe, 0xbe, 0x09, 0x23, 0x3e, 0x99, 0xfb, 0x06, 0xa1, 0x33,
	0x33, 0x7f, 0xa6, 0xed, 0x04, 0xf8, 0x25, 0x65, 0x53, 0xb1, 0x31, 0x7f, 0xc8, 0xdf, 0x94, 0x3f,
	0xe6, 0xcc, 0x2f, 0x89, 0x2f, 0x71, 0xfc, 0x3f, 0x00, 0x00, 0xff, 0xff, 0xe1, 0x1d, 0xaf, 0x18,
	0x67, 0x03, 0x00, 0x00,
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
	dAtA[i] = 0x52
	{
		size := m.Leverage.Size()
		i -= size
		if _, err := m.Leverage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOrder(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x4a
	{
		size := m.Collateral.Size()
		i -= size
		if _, err := m.Collateral.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOrder(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	if m.PerpetualId != 0 {
		i = encodeVarintOrder(dAtA, i, uint64(m.PerpetualId))
		i--
		dAtA[i] = 0x38
	}
	if m.SubAccountId != 0 {
		i = encodeVarintOrder(dAtA, i, uint64(m.SubAccountId))
		i--
		dAtA[i] = 0x30
	}
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
	if m.SubAccountId != 0 {
		n += 1 + sovOrder(uint64(m.SubAccountId))
	}
	if m.PerpetualId != 0 {
		n += 1 + sovOrder(uint64(m.PerpetualId))
	}
	l = m.Collateral.Size()
	n += 1 + l + sovOrder(uint64(l))
	l = m.Leverage.Size()
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SubAccountId", wireType)
			}
			m.SubAccountId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrder
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SubAccountId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PerpetualId", wireType)
			}
			m.PerpetualId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOrder
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PerpetualId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Collateral", wireType)
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
			if err := m.Collateral.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Leverage", wireType)
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
			if err := m.Leverage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
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
