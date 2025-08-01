// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/tokenomics/inflation_entry.proto

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

type InflationEntry struct {
	LmRewards         uint64 `protobuf:"varint,1,opt,name=lm_rewards,json=lmRewards,proto3" json:"lm_rewards,omitempty"`
	IcsStakingRewards uint64 `protobuf:"varint,2,opt,name=ics_staking_rewards,json=icsStakingRewards,proto3" json:"ics_staking_rewards,omitempty"`
	CommunityFund     uint64 `protobuf:"varint,3,opt,name=community_fund,json=communityFund,proto3" json:"community_fund,omitempty"`
	StrategicReserve  uint64 `protobuf:"varint,4,opt,name=strategic_reserve,json=strategicReserve,proto3" json:"strategic_reserve,omitempty"`
	TeamTokensVested  uint64 `protobuf:"varint,5,opt,name=team_tokens_vested,json=teamTokensVested,proto3" json:"team_tokens_vested,omitempty"`
}

func (m *InflationEntry) Reset()         { *m = InflationEntry{} }
func (m *InflationEntry) String() string { return proto.CompactTextString(m) }
func (*InflationEntry) ProtoMessage()    {}
func (*InflationEntry) Descriptor() ([]byte, []int) {
	return fileDescriptor_fbe674fab3da1fa3, []int{0}
}
func (m *InflationEntry) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InflationEntry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InflationEntry.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InflationEntry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InflationEntry.Merge(m, src)
}
func (m *InflationEntry) XXX_Size() int {
	return m.Size()
}
func (m *InflationEntry) XXX_DiscardUnknown() {
	xxx_messageInfo_InflationEntry.DiscardUnknown(m)
}

var xxx_messageInfo_InflationEntry proto.InternalMessageInfo

func (m *InflationEntry) GetLmRewards() uint64 {
	if m != nil {
		return m.LmRewards
	}
	return 0
}

func (m *InflationEntry) GetIcsStakingRewards() uint64 {
	if m != nil {
		return m.IcsStakingRewards
	}
	return 0
}

func (m *InflationEntry) GetCommunityFund() uint64 {
	if m != nil {
		return m.CommunityFund
	}
	return 0
}

func (m *InflationEntry) GetStrategicReserve() uint64 {
	if m != nil {
		return m.StrategicReserve
	}
	return 0
}

func (m *InflationEntry) GetTeamTokensVested() uint64 {
	if m != nil {
		return m.TeamTokensVested
	}
	return 0
}

func init() {
	proto.RegisterType((*InflationEntry)(nil), "elys.tokenomics.InflationEntry")
}

func init() {
	proto.RegisterFile("elys/tokenomics/inflation_entry.proto", fileDescriptor_fbe674fab3da1fa3)
}

var fileDescriptor_fbe674fab3da1fa3 = []byte{
	// 284 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0xd0, 0x41, 0x4b, 0xc3, 0x30,
	0x14, 0x07, 0xf0, 0x45, 0xa7, 0x60, 0xc0, 0xe9, 0xe2, 0xa5, 0x17, 0x83, 0x08, 0x03, 0x41, 0x6d,
	0x41, 0x0f, 0xde, 0x05, 0x05, 0xc1, 0xd3, 0x14, 0x0f, 0x5e, 0x42, 0x97, 0xbe, 0xcd, 0xb0, 0x26,
	0x19, 0xc9, 0xeb, 0x66, 0xbf, 0x85, 0x1f, 0xcb, 0xe3, 0x8e, 0x1e, 0xb5, 0xfd, 0x22, 0xd2, 0x94,
	0x55, 0x8f, 0xef, 0xff, 0xff, 0x5d, 0xde, 0x9f, 0x8e, 0x20, 0x2f, 0x7d, 0x82, 0x76, 0x0e, 0xc6,
	0x6a, 0x25, 0x7d, 0xa2, 0xcc, 0x34, 0x4f, 0x51, 0x59, 0x23, 0xc0, 0xa0, 0x2b, 0xe3, 0x85, 0xb3,
	0x68, 0xd9, 0x41, 0xc3, 0xe2, 0x3f, 0x76, 0xfa, 0x43, 0xe8, 0xe0, 0x61, 0x43, 0xef, 0x1a, 0xc9,
	0x8e, 0x29, 0xcd, 0xb5, 0x70, 0xb0, 0x4a, 0x5d, 0xe6, 0x23, 0x72, 0x42, 0xce, 0xfa, 0xe3, 0xbd,
	0x5c, 0x8f, 0xdb, 0x80, 0xc5, 0xf4, 0x48, 0x49, 0x2f, 0x3c, 0xa6, 0x73, 0x65, 0x66, 0x9d, 0xdb,
	0x0a, 0x6e, 0xa8, 0xa4, 0x7f, 0x6a, 0x9b, 0x8d, 0x1f, 0xd1, 0x81, 0xb4, 0x5a, 0x17, 0x46, 0x61,
	0x29, 0xa6, 0x85, 0xc9, 0xa2, 0xed, 0x40, 0xf7, 0xbb, 0xf4, 0xbe, 0x30, 0x19, 0x3b, 0xa7, 0x43,
	0x8f, 0x2e, 0x45, 0x98, 0x29, 0x29, 0x1c, 0x78, 0x70, 0x4b, 0x88, 0xfa, 0x41, 0x1e, 0x76, 0xc5,
	0xb8, 0xcd, 0xd9, 0x05, 0x65, 0x08, 0xa9, 0x16, 0xe1, 0x11, 0x2f, 0x96, 0xe0, 0x11, 0xb2, 0x68,
	0xa7, 0xd5, 0x4d, 0xf3, 0x1c, 0x8a, 0x97, 0x90, 0xdf, 0x3e, 0x7e, 0x56, 0x9c, 0xac, 0x2b, 0x4e,
	0xbe, 0x2b, 0x4e, 0x3e, 0x6a, 0xde, 0x5b, 0xd7, 0xbc, 0xf7, 0x55, 0xf3, 0xde, 0xeb, 0xd5, 0x4c,
	0xe1, 0x5b, 0x31, 0x89, 0xa5, 0xd5, 0x49, 0xb3, 0xcc, 0xa5, 0x01, 0x5c, 0x59, 0x37, 0x0f, 0x47,
	0xb2, 0xbc, 0x49, 0xde, 0xff, 0x4f, 0x8a, 0xe5, 0x02, 0xfc, 0x64, 0x37, 0x2c, 0x79, 0xfd, 0x1b,
	0x00, 0x00, 0xff, 0xff, 0x27, 0x6c, 0xa9, 0x23, 0x72, 0x01, 0x00, 0x00,
}

func (m *InflationEntry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InflationEntry) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InflationEntry) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.TeamTokensVested != 0 {
		i = encodeVarintInflationEntry(dAtA, i, uint64(m.TeamTokensVested))
		i--
		dAtA[i] = 0x28
	}
	if m.StrategicReserve != 0 {
		i = encodeVarintInflationEntry(dAtA, i, uint64(m.StrategicReserve))
		i--
		dAtA[i] = 0x20
	}
	if m.CommunityFund != 0 {
		i = encodeVarintInflationEntry(dAtA, i, uint64(m.CommunityFund))
		i--
		dAtA[i] = 0x18
	}
	if m.IcsStakingRewards != 0 {
		i = encodeVarintInflationEntry(dAtA, i, uint64(m.IcsStakingRewards))
		i--
		dAtA[i] = 0x10
	}
	if m.LmRewards != 0 {
		i = encodeVarintInflationEntry(dAtA, i, uint64(m.LmRewards))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintInflationEntry(dAtA []byte, offset int, v uint64) int {
	offset -= sovInflationEntry(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *InflationEntry) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.LmRewards != 0 {
		n += 1 + sovInflationEntry(uint64(m.LmRewards))
	}
	if m.IcsStakingRewards != 0 {
		n += 1 + sovInflationEntry(uint64(m.IcsStakingRewards))
	}
	if m.CommunityFund != 0 {
		n += 1 + sovInflationEntry(uint64(m.CommunityFund))
	}
	if m.StrategicReserve != 0 {
		n += 1 + sovInflationEntry(uint64(m.StrategicReserve))
	}
	if m.TeamTokensVested != 0 {
		n += 1 + sovInflationEntry(uint64(m.TeamTokensVested))
	}
	return n
}

func sovInflationEntry(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozInflationEntry(x uint64) (n int) {
	return sovInflationEntry(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *InflationEntry) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowInflationEntry
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
			return fmt.Errorf("proto: InflationEntry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InflationEntry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LmRewards", wireType)
			}
			m.LmRewards = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInflationEntry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LmRewards |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IcsStakingRewards", wireType)
			}
			m.IcsStakingRewards = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInflationEntry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.IcsStakingRewards |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommunityFund", wireType)
			}
			m.CommunityFund = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInflationEntry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CommunityFund |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field StrategicReserve", wireType)
			}
			m.StrategicReserve = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInflationEntry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.StrategicReserve |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TeamTokensVested", wireType)
			}
			m.TeamTokensVested = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowInflationEntry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TeamTokensVested |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipInflationEntry(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthInflationEntry
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
func skipInflationEntry(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowInflationEntry
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
					return 0, ErrIntOverflowInflationEntry
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
					return 0, ErrIntOverflowInflationEntry
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
				return 0, ErrInvalidLengthInflationEntry
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupInflationEntry
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthInflationEntry
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthInflationEntry        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowInflationEntry          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupInflationEntry = fmt.Errorf("proto: unexpected end of group")
)
