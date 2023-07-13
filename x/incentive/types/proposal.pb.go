// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/incentive/proposal.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

type PoolMultipliers struct {
	PoolId     uint64                                 `protobuf:"varint,1,opt,name=poolId,proto3" json:"poolId,omitempty"`
	Multiplier github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=multiplier,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"multiplier"`
}

func (m *PoolMultipliers) Reset()         { *m = PoolMultipliers{} }
func (m *PoolMultipliers) String() string { return proto.CompactTextString(m) }
func (*PoolMultipliers) ProtoMessage()    {}
func (*PoolMultipliers) Descriptor() ([]byte, []int) {
	return fileDescriptor_ccb8f1511f38d9ad, []int{0}
}
func (m *PoolMultipliers) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PoolMultipliers) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PoolMultipliers.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PoolMultipliers) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PoolMultipliers.Merge(m, src)
}
func (m *PoolMultipliers) XXX_Size() int {
	return m.Size()
}
func (m *PoolMultipliers) XXX_DiscardUnknown() {
	xxx_messageInfo_PoolMultipliers.DiscardUnknown(m)
}

var xxx_messageInfo_PoolMultipliers proto.InternalMessageInfo

func (m *PoolMultipliers) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

type ProposalUpdatePoolMultipliers struct {
	Title           string            `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Description     string            `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	PoolMultipliers []PoolMultipliers `protobuf:"bytes,3,rep,name=poolMultipliers,proto3" json:"poolMultipliers"`
}

func (m *ProposalUpdatePoolMultipliers) Reset()         { *m = ProposalUpdatePoolMultipliers{} }
func (m *ProposalUpdatePoolMultipliers) String() string { return proto.CompactTextString(m) }
func (*ProposalUpdatePoolMultipliers) ProtoMessage()    {}
func (*ProposalUpdatePoolMultipliers) Descriptor() ([]byte, []int) {
	return fileDescriptor_ccb8f1511f38d9ad, []int{1}
}
func (m *ProposalUpdatePoolMultipliers) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProposalUpdatePoolMultipliers) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProposalUpdatePoolMultipliers.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProposalUpdatePoolMultipliers) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProposalUpdatePoolMultipliers.Merge(m, src)
}
func (m *ProposalUpdatePoolMultipliers) XXX_Size() int {
	return m.Size()
}
func (m *ProposalUpdatePoolMultipliers) XXX_DiscardUnknown() {
	xxx_messageInfo_ProposalUpdatePoolMultipliers.DiscardUnknown(m)
}

var xxx_messageInfo_ProposalUpdatePoolMultipliers proto.InternalMessageInfo

func (m *ProposalUpdatePoolMultipliers) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *ProposalUpdatePoolMultipliers) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ProposalUpdatePoolMultipliers) GetPoolMultipliers() []PoolMultipliers {
	if m != nil {
		return m.PoolMultipliers
	}
	return nil
}

func init() {
	proto.RegisterType((*PoolMultipliers)(nil), "elys.incentive.PoolMultipliers")
	proto.RegisterType((*ProposalUpdatePoolMultipliers)(nil), "elys.incentive.ProposalUpdatePoolMultipliers")
}

func init() { proto.RegisterFile("elys/incentive/proposal.proto", fileDescriptor_ccb8f1511f38d9ad) }

var fileDescriptor_ccb8f1511f38d9ad = []byte{
	// 315 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x51, 0x41, 0x4b, 0xfb, 0x30,
	0x1c, 0x6d, 0xfe, 0xdb, 0x7f, 0xb0, 0x0c, 0x1c, 0x94, 0x21, 0x73, 0xb0, 0xac, 0xec, 0x20, 0xbb,
	0x2c, 0x05, 0xfd, 0x06, 0xc3, 0x83, 0x1e, 0xd4, 0x51, 0xf0, 0xe2, 0x45, 0xb6, 0x36, 0xd4, 0xb0,
	0xb4, 0xbf, 0xd0, 0x64, 0x6a, 0xbf, 0x85, 0x1f, 0xc3, 0x8f, 0xb2, 0xe3, 0x8e, 0xe2, 0x61, 0x48,
	0xfb, 0x45, 0xa4, 0x4d, 0x19, 0x5d, 0x4f, 0xc9, 0xcb, 0xfb, 0xe5, 0xe5, 0xbd, 0x17, 0x3c, 0x66,
	0x22, 0x55, 0x2e, 0x8f, 0x7d, 0x16, 0x6b, 0xfe, 0xc6, 0x5c, 0x99, 0x80, 0x04, 0xb5, 0x12, 0x54,
	0x26, 0xa0, 0xc1, 0x3e, 0x2b, 0x68, 0x7a, 0xa4, 0x47, 0x83, 0x10, 0x42, 0x28, 0x29, 0xb7, 0xd8,
	0x99, 0xa9, 0xd1, 0x85, 0x0f, 0x2a, 0x02, 0xf5, 0x62, 0x08, 0x03, 0x0c, 0x35, 0x4d, 0x71, 0x7f,
	0x09, 0x20, 0xee, 0xb7, 0x42, 0x73, 0x29, 0x38, 0x4b, 0x94, 0x7d, 0x8e, 0x3b, 0x12, 0x40, 0xdc,
	0x05, 0x43, 0xe4, 0xa0, 0x59, 0xdb, 0xab, 0x90, 0xfd, 0x80, 0x71, 0x74, 0x1c, 0x1b, 0xfe, 0x73,
	0xd0, 0xac, 0xbb, 0xa0, 0xbb, 0xc3, 0xc4, 0xfa, 0x39, 0x4c, 0x2e, 0x43, 0xae, 0x5f, 0xb7, 0x6b,
	0xea, 0x43, 0x54, 0xe9, 0x57, 0xcb, 0x5c, 0x05, 0x1b, 0x57, 0xa7, 0x92, 0x29, 0x7a, 0xc3, 0x7c,
	0xaf, 0xa6, 0x30, 0xfd, 0x42, 0x78, 0xbc, 0xac, 0xe2, 0x3c, 0xc9, 0x60, 0xa5, 0x59, 0xd3, 0xc9,
	0x00, 0xff, 0xd7, 0x5c, 0x0b, 0x56, 0x1a, 0xe9, 0x7a, 0x06, 0xd8, 0x0e, 0xee, 0x05, 0x4c, 0xf9,
	0x09, 0x97, 0x9a, 0x43, 0x6c, 0x8c, 0x78, 0xf5, 0x23, 0xfb, 0x11, 0xf7, 0xe5, 0xa9, 0xd4, 0xb0,
	0xe5, 0xb4, 0x66, 0xbd, 0xab, 0x09, 0x3d, 0xed, 0x8b, 0x36, 0x5e, 0x5c, 0xb4, 0x8b, 0x3c, 0x5e,
	0xf3, 0xf6, 0xe2, 0x76, 0x97, 0x11, 0xb4, 0xcf, 0x08, 0xfa, 0xcd, 0x08, 0xfa, 0xcc, 0x89, 0xb5,
	0xcf, 0x89, 0xf5, 0x9d, 0x13, 0xeb, 0x99, 0xd6, 0x82, 0x17, 0xda, 0xf3, 0x98, 0xe9, 0x77, 0x48,
	0x36, 0x25, 0x70, 0x3f, 0x6a, 0x3f, 0x57, 0x96, 0xb0, 0xee, 0x94, 0xb5, 0x5f, 0xff, 0x05, 0x00,
	0x00, 0xff, 0xff, 0x78, 0x16, 0x33, 0xc0, 0xd8, 0x01, 0x00, 0x00,
}

func (m *PoolMultipliers) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PoolMultipliers) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PoolMultipliers) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.Multiplier.Size()
		i -= size
		if _, err := m.Multiplier.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintProposal(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.PoolId != 0 {
		i = encodeVarintProposal(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ProposalUpdatePoolMultipliers) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProposalUpdatePoolMultipliers) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProposalUpdatePoolMultipliers) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PoolMultipliers) > 0 {
		for iNdEx := len(m.PoolMultipliers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.PoolMultipliers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintProposal(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintProposal(dAtA []byte, offset int, v uint64) int {
	offset -= sovProposal(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PoolMultipliers) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PoolId != 0 {
		n += 1 + sovProposal(uint64(m.PoolId))
	}
	l = m.Multiplier.Size()
	n += 1 + l + sovProposal(uint64(l))
	return n
}

func (m *ProposalUpdatePoolMultipliers) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	if len(m.PoolMultipliers) > 0 {
		for _, e := range m.PoolMultipliers {
			l = e.Size()
			n += 1 + l + sovProposal(uint64(l))
		}
	}
	return n
}

func sovProposal(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProposal(x uint64) (n int) {
	return sovProposal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PoolMultipliers) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposal
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
			return fmt.Errorf("proto: PoolMultipliers: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PoolMultipliers: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PoolId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Multiplier", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Multiplier.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProposal
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
func (m *ProposalUpdatePoolMultipliers) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposal
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
			return fmt.Errorf("proto: ProposalUpdatePoolMultipliers: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProposalUpdatePoolMultipliers: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolMultipliers", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PoolMultipliers = append(m.PoolMultipliers, PoolMultipliers{})
			if err := m.PoolMultipliers[len(m.PoolMultipliers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProposal
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
func skipProposal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProposal
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
					return 0, ErrIntOverflowProposal
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
					return 0, ErrIntOverflowProposal
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
				return 0, ErrInvalidLengthProposal
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProposal
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProposal
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProposal        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProposal          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProposal = fmt.Errorf("proto: unexpected end of group")
)
