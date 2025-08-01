// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/perpetual/genesis.proto

package types

import (
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

// GenesisState defines the perpetual module's genesis state.
type GenesisState struct {
	Params           Params             `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	PoolList         []Pool             `protobuf:"bytes,2,rep,name=pool_list,json=poolList,proto3" json:"pool_list"`
	MtpList          []MTP              `protobuf:"bytes,3,rep,name=mtp_list,json=mtpList,proto3" json:"mtp_list"`
	AddressWhitelist []string           `protobuf:"bytes,4,rep,name=address_whitelist,json=addressWhitelist,proto3" json:"address_whitelist,omitempty"`
	PerpetualCounter []PerpetualCounter `protobuf:"bytes,5,rep,name=perpetual_counter,json=perpetualCounter,proto3" json:"perpetual_counter"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_2b712f25285e9e2d, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetPoolList() []Pool {
	if m != nil {
		return m.PoolList
	}
	return nil
}

func (m *GenesisState) GetMtpList() []MTP {
	if m != nil {
		return m.MtpList
	}
	return nil
}

func (m *GenesisState) GetAddressWhitelist() []string {
	if m != nil {
		return m.AddressWhitelist
	}
	return nil
}

func (m *GenesisState) GetPerpetualCounter() []PerpetualCounter {
	if m != nil {
		return m.PerpetualCounter
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "elys.perpetual.GenesisState")
}

func init() { proto.RegisterFile("elys/perpetual/genesis.proto", fileDescriptor_2b712f25285e9e2d) }

var fileDescriptor_2b712f25285e9e2d = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0xc1, 0x4a, 0x02, 0x41,
	0x18, 0xc7, 0x77, 0xd5, 0x4c, 0xc7, 0x08, 0xdd, 0x24, 0x6c, 0x8b, 0x6d, 0xe9, 0x24, 0x44, 0x3b,
	0x64, 0x82, 0xf7, 0x3a, 0x74, 0xa8, 0x40, 0x34, 0x08, 0xba, 0xc8, 0xaa, 0xc3, 0xba, 0xb4, 0xeb,
	0x0c, 0x33, 0x9f, 0x99, 0x6f, 0xd1, 0xcb, 0xf4, 0x0e, 0x1e, 0x3d, 0x76, 0x8a, 0x70, 0x5f, 0x24,
	0x76, 0x76, 0x14, 0x1a, 0x6f, 0x33, 0xf3, 0xfb, 0xff, 0xfe, 0xdf, 0x07, 0x83, 0xce, 0x48, 0xb4,
	0x10, 0x98, 0x11, 0xce, 0x08, 0xcc, 0xfc, 0x08, 0x07, 0x64, 0x4a, 0x44, 0x28, 0x3c, 0xc6, 0x29,
	0x50, 0xeb, 0x30, 0xa5, 0xde, 0x96, 0xda, 0xf5, 0x80, 0x06, 0x54, 0x22, 0x9c, 0x9e, 0xb2, 0x94,
	0x7d, 0xaa, 0x75, 0x30, 0x9f, 0xfb, 0xb1, 0xaa, 0xb0, 0x4f, 0x74, 0x48, 0x69, 0xa4, 0x90, 0xad,
	0x21, 0x58, 0x30, 0xa2, 0xb4, 0x8b, 0xaf, 0x1c, 0x3a, 0xb8, 0xcf, 0x76, 0xe9, 0x83, 0x0f, 0xc4,
	0x6a, 0xa3, 0x62, 0xd6, 0xdb, 0x30, 0x5d, 0xb3, 0x59, 0x69, 0x1d, 0x7b, 0xff, 0x77, 0xf3, 0xba,
	0x92, 0xde, 0x16, 0x96, 0x3f, 0xe7, 0x46, 0x4f, 0x65, 0xad, 0x0e, 0x2a, 0xa7, 0x03, 0x07, 0x51,
	0x28, 0xa0, 0x91, 0x73, 0xf3, 0xcd, 0x4a, 0xab, 0xbe, 0x23, 0x52, 0x1a, 0x29, 0xad, 0x94, 0x86,
	0x1f, 0x43, 0x01, 0x56, 0x1b, 0x95, 0x62, 0x60, 0x99, 0x97, 0x97, 0xde, 0x91, 0xee, 0x3d, 0x3d,
	0x77, 0x95, 0xb6, 0x1f, 0x03, 0x93, 0xd6, 0x25, 0xaa, 0xf9, 0xe3, 0x31, 0x27, 0x42, 0x0c, 0xe6,
	0x93, 0x10, 0x88, 0xd4, 0x0b, 0x6e, 0xbe, 0x59, 0xee, 0x55, 0x15, 0x78, 0xd9, 0xbc, 0x5b, 0x7d,
	0x54, 0xdb, 0x96, 0x0d, 0x46, 0x74, 0x36, 0x05, 0xc2, 0x1b, 0x7b, 0x72, 0x96, 0xbb, 0xb3, 0xe3,
	0xe6, 0x74, 0x97, 0xe5, 0xd4, 0xe0, 0x2a, 0xd3, 0xdf, 0x1f, 0x96, 0x6b, 0xc7, 0x5c, 0xad, 0x1d,
	0xf3, 0x77, 0xed, 0x98, 0x9f, 0x89, 0x63, 0xac, 0x12, 0xc7, 0xf8, 0x4e, 0x1c, 0xe3, 0xf5, 0x3a,
	0x08, 0x61, 0x32, 0x1b, 0x7a, 0x23, 0x1a, 0xe3, 0xb4, 0xfd, 0x6a, 0x4a, 0x60, 0x4e, 0xf9, 0x9b,
	0xbc, 0xe0, 0xf7, 0x0e, 0xfe, 0xd0, 0xbf, 0x62, 0x58, 0x94, 0x7f, 0x71, 0xf3, 0x17, 0x00, 0x00,
	0xff, 0xff, 0x26, 0xa7, 0x2c, 0x66, 0x25, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PerpetualCounter) > 0 {
		for iNdEx := len(m.PerpetualCounter) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PerpetualCounter[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.AddressWhitelist) > 0 {
		for iNdEx := len(m.AddressWhitelist) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.AddressWhitelist[iNdEx])
			copy(dAtA[i:], m.AddressWhitelist[iNdEx])
			i = encodeVarintGenesis(dAtA, i, uint64(len(m.AddressWhitelist[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.MtpList) > 0 {
		for iNdEx := len(m.MtpList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MtpList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.PoolList) > 0 {
		for iNdEx := len(m.PoolList) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PoolList[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.PoolList) > 0 {
		for _, e := range m.PoolList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.MtpList) > 0 {
		for _, e := range m.MtpList {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.AddressWhitelist) > 0 {
		for _, s := range m.AddressWhitelist {
			l = len(s)
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.PerpetualCounter) > 0 {
		for _, e := range m.PerpetualCounter {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PoolList = append(m.PoolList, Pool{})
			if err := m.PoolList[len(m.PoolList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MtpList", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MtpList = append(m.MtpList, MTP{})
			if err := m.MtpList[len(m.MtpList)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddressWhitelist", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AddressWhitelist = append(m.AddressWhitelist, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PerpetualCounter", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PerpetualCounter = append(m.PerpetualCounter, PerpetualCounter{})
			if err := m.PerpetualCounter[len(m.PerpetualCounter)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
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
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
