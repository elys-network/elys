// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/commitment/genesis.proto

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

// GenesisState defines the commitment module's genesis state.
type GenesisState struct {
	Params      Params         `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	Commitments []*Commitments `protobuf:"bytes,2,rep,name=commitments,proto3" json:"commitments,omitempty"`
	AtomStakers []*AtomStaker  `protobuf:"bytes,3,rep,name=atom_stakers,json=atomStakers,proto3" json:"atom_stakers,omitempty"`
	NftHolders  []*NftHolder   `protobuf:"bytes,4,rep,name=nft_holders,json=nftHolders,proto3" json:"nft_holders,omitempty"`
	Cadets      []*Cadet       `protobuf:"bytes,5,rep,name=cadets,proto3" json:"cadets,omitempty"`
	Governors   []*Governor    `protobuf:"bytes,6,rep,name=governors,proto3" json:"governors,omitempty"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_5600d7581efecfbc, []int{0}
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

func (m *GenesisState) GetCommitments() []*Commitments {
	if m != nil {
		return m.Commitments
	}
	return nil
}

func (m *GenesisState) GetAtomStakers() []*AtomStaker {
	if m != nil {
		return m.AtomStakers
	}
	return nil
}

func (m *GenesisState) GetNftHolders() []*NftHolder {
	if m != nil {
		return m.NftHolders
	}
	return nil
}

func (m *GenesisState) GetCadets() []*Cadet {
	if m != nil {
		return m.Cadets
	}
	return nil
}

func (m *GenesisState) GetGovernors() []*Governor {
	if m != nil {
		return m.Governors
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "elys.commitment.GenesisState")
}

func init() { proto.RegisterFile("elys/commitment/genesis.proto", fileDescriptor_5600d7581efecfbc) }

var fileDescriptor_5600d7581efecfbc = []byte{
	// 346 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0xcd, 0x4a, 0xf3, 0x40,
	0x14, 0x86, 0x93, 0xb6, 0x5f, 0xe0, 0x9b, 0x14, 0x84, 0x41, 0x34, 0xc6, 0x1a, 0xab, 0xab, 0x6e,
	0x4c, 0xa0, 0x22, 0x2e, 0x84, 0x82, 0x75, 0x51, 0xdd, 0x88, 0xa4, 0x3b, 0x37, 0x65, 0xda, 0x4e,
	0xd3, 0xd0, 0x26, 0x13, 0x66, 0x8e, 0x3f, 0xbd, 0x0b, 0x2f, 0xab, 0xcb, 0x2e, 0x5d, 0x89, 0xb4,
	0x4b, 0x6f, 0x42, 0x32, 0x49, 0x4d, 0xc9, 0xec, 0xce, 0xf0, 0x3e, 0xcf, 0x3b, 0x07, 0x0e, 0x3a,
	0xa1, 0xf3, 0x85, 0xf0, 0x46, 0x2c, 0x8a, 0x42, 0x88, 0x68, 0x0c, 0x5e, 0x40, 0x63, 0x2a, 0x42,
	0xe1, 0x26, 0x9c, 0x01, 0xc3, 0x7b, 0x69, 0xec, 0x16, 0xb1, 0xbd, 0x1f, 0xb0, 0x80, 0xc9, 0xcc,
	0x4b, 0xa7, 0x0c, 0xb3, 0x1b, 0xe5, 0x96, 0x84, 0x70, 0x12, 0xe5, 0x25, 0xf6, 0x59, 0x39, 0x2d,
	0xc6, 0x2d, 0xa2, 0xac, 0x41, 0x42, 0x3e, 0xe6, 0x2c, 0xc9, 0xe2, 0xf3, 0x9f, 0x0a, 0xaa, 0xf7,
	0xb2, 0xc5, 0xfa, 0x40, 0x80, 0xe2, 0x2b, 0x64, 0x64, 0x5f, 0x58, 0x7a, 0x53, 0x6f, 0x99, 0xed,
	0x43, 0xb7, 0xb4, 0xa8, 0xfb, 0x24, 0xe3, 0x6e, 0x6d, 0xf9, 0x75, 0xaa, 0xf9, 0x39, 0x8c, 0x3b,
	0xc8, 0xdc, 0xf9, 0xdb, 0xaa, 0x34, 0xab, 0x2d, 0xb3, 0xdd, 0x50, 0xdc, 0xbb, 0x82, 0xf1, 0x77,
	0x05, 0xdc, 0x41, 0x75, 0x02, 0x2c, 0x1a, 0x08, 0x20, 0x33, 0xca, 0x85, 0x55, 0x95, 0x05, 0xc7,
	0x4a, 0xc1, 0x2d, 0xb0, 0xa8, 0x2f, 0x19, 0xdf, 0x24, 0x7f, 0xb3, 0xc0, 0x37, 0xc8, 0x8c, 0x27,
	0x30, 0x98, 0xb2, 0xf9, 0x38, 0xd5, 0x6b, 0x52, 0xb7, 0x15, 0xfd, 0x71, 0x02, 0xf7, 0x12, 0xf1,
	0x51, 0xbc, 0x1d, 0x05, 0x76, 0x91, 0x31, 0x22, 0x63, 0x0a, 0xc2, 0xfa, 0x27, 0xbd, 0x03, 0x75,
	0xef, 0x34, 0xf6, 0x73, 0x0a, 0x5f, 0xa3, 0xff, 0x01, 0x7b, 0xa5, 0x3c, 0x66, 0x5c, 0x58, 0x86,
	0x54, 0x8e, 0x14, 0xa5, 0x97, 0x13, 0x7e, 0xc1, 0x76, 0x1f, 0x96, 0x6b, 0x47, 0x5f, 0xad, 0x1d,
	0xfd, 0x7b, 0xed, 0xe8, 0x1f, 0x1b, 0x47, 0x5b, 0x6d, 0x1c, 0xed, 0x73, 0xe3, 0x68, 0xcf, 0x5e,
	0x10, 0xc2, 0xf4, 0x65, 0x98, 0x16, 0x78, 0x69, 0xd3, 0x45, 0x4c, 0xe1, 0x8d, 0xf1, 0x99, 0x7c,
	0x78, 0xef, 0xbb, 0x07, 0x84, 0x45, 0x42, 0xc5, 0xd0, 0x90, 0xf7, 0xbb, 0xfc, 0x0d, 0x00, 0x00,
	0xff, 0xff, 0xb4, 0xe3, 0xf8, 0x22, 0x67, 0x02, 0x00, 0x00,
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
	if len(m.Governors) > 0 {
		for iNdEx := len(m.Governors) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Governors[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.Cadets) > 0 {
		for iNdEx := len(m.Cadets) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Cadets[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.NftHolders) > 0 {
		for iNdEx := len(m.NftHolders) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.NftHolders[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.AtomStakers) > 0 {
		for iNdEx := len(m.AtomStakers) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.AtomStakers[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.Commitments) > 0 {
		for iNdEx := len(m.Commitments) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Commitments[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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
	if len(m.Commitments) > 0 {
		for _, e := range m.Commitments {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.AtomStakers) > 0 {
		for _, e := range m.AtomStakers {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.NftHolders) > 0 {
		for _, e := range m.NftHolders {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Cadets) > 0 {
		for _, e := range m.Cadets {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Governors) > 0 {
		for _, e := range m.Governors {
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
				return fmt.Errorf("proto: wrong wireType = %d for field Commitments", wireType)
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
			m.Commitments = append(m.Commitments, &Commitments{})
			if err := m.Commitments[len(m.Commitments)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AtomStakers", wireType)
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
			m.AtomStakers = append(m.AtomStakers, &AtomStaker{})
			if err := m.AtomStakers[len(m.AtomStakers)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field NftHolders", wireType)
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
			m.NftHolders = append(m.NftHolders, &NftHolder{})
			if err := m.NftHolders[len(m.NftHolders)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Cadets", wireType)
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
			m.Cadets = append(m.Cadets, &Cadet{})
			if err := m.Cadets[len(m.Cadets)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Governors", wireType)
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
			m.Governors = append(m.Governors, &Governor{})
			if err := m.Governors[len(m.Governors)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
