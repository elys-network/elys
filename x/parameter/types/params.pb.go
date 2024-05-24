// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/parameter/params.proto

package types

import (
	fmt "fmt"
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

// Params defines the parameters for the module.
type Params struct {
	MinCommissionRate       github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,1,opt,name=min_commission_rate,json=minCommissionRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"min_commission_rate"`
	MaxVotingPower          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=max_voting_power,json=maxVotingPower,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"max_voting_power"`
	MinSelfDelegation       github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,opt,name=min_self_delegation,json=minSelfDelegation,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"min_self_delegation"`
	BrokerAddress           string                                 `protobuf:"bytes,4,opt,name=broker_address,json=brokerAddress,proto3" json:"broker_address,omitempty"`
	TotalBlocksPerYear      int64                                  `protobuf:"varint,5,opt,name=total_blocks_per_year,json=totalBlocksPerYear,proto3" json:"total_blocks_per_year,omitempty"`
	WasmMaxLabelSize        github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=wasm_max_label_size,json=wasmMaxLabelSize,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"wasm_max_label_size"`
	WasmMaxSize             github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,7,opt,name=wasm_max_size,json=wasmMaxSize,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"wasm_max_size"`
	WasmMaxProposalWasmSize github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,8,opt,name=wasm_max_proposal_wasm_size,json=wasmMaxProposalWasmSize,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"wasm_max_proposal_wasm_size"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_b61780a5be327c2b, []int{0}
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

func (m *Params) GetBrokerAddress() string {
	if m != nil {
		return m.BrokerAddress
	}
	return ""
}

func (m *Params) GetTotalBlocksPerYear() int64 {
	if m != nil {
		return m.TotalBlocksPerYear
	}
	return 0
}

type LegacyParams struct {
	MinCommissionRate github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,1,opt,name=min_commission_rate,json=minCommissionRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"min_commission_rate"`
	MaxVotingPower    github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=max_voting_power,json=maxVotingPower,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"max_voting_power"`
	MinSelfDelegation github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,3,opt,name=min_self_delegation,json=minSelfDelegation,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"min_self_delegation"`
	BrokerAddress     string                                 `protobuf:"bytes,4,opt,name=broker_address,json=brokerAddress,proto3" json:"broker_address,omitempty"`
}

func (m *LegacyParams) Reset()      { *m = LegacyParams{} }
func (*LegacyParams) ProtoMessage() {}
func (*LegacyParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_b61780a5be327c2b, []int{1}
}
func (m *LegacyParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LegacyParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_LegacyParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *LegacyParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LegacyParams.Merge(m, src)
}
func (m *LegacyParams) XXX_Size() int {
	return m.Size()
}
func (m *LegacyParams) XXX_DiscardUnknown() {
	xxx_messageInfo_LegacyParams.DiscardUnknown(m)
}

var xxx_messageInfo_LegacyParams proto.InternalMessageInfo

func (m *LegacyParams) GetBrokerAddress() string {
	if m != nil {
		return m.BrokerAddress
	}
	return ""
}

func init() {
	proto.RegisterType((*Params)(nil), "elys.parameter.Params")
	proto.RegisterType((*LegacyParams)(nil), "elys.parameter.LegacyParams")
}

func init() { proto.RegisterFile("elys/parameter/params.proto", fileDescriptor_b61780a5be327c2b) }

var fileDescriptor_b61780a5be327c2b = []byte{
	// 431 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0xd3, 0x4f, 0x6b, 0xd4, 0x40,
	0x18, 0x06, 0xf0, 0xc4, 0xb6, 0xab, 0x8e, 0x74, 0xa9, 0xa9, 0x62, 0xb0, 0x90, 0x2d, 0x82, 0xd2,
	0x4b, 0x13, 0xc4, 0x9b, 0x37, 0xd7, 0x1e, 0x14, 0x14, 0x96, 0x14, 0xfc, 0x07, 0x3a, 0x4c, 0xb2,
	0x6f, 0xe3, 0xb0, 0x33, 0x79, 0xc3, 0xbc, 0xa3, 0xbb, 0xdb, 0x4f, 0xe1, 0xd1, 0xa3, 0xe0, 0x97,
	0xe9, 0xb1, 0x47, 0xf1, 0x50, 0x64, 0xf7, 0x8b, 0xc8, 0xcc, 0xae, 0x61, 0xaf, 0xe6, 0x94, 0xc9,
	0x3c, 0xc9, 0xef, 0x79, 0x09, 0x13, 0x76, 0x00, 0x6a, 0x4e, 0x59, 0x23, 0x8c, 0xd0, 0x60, 0xc1,
	0xac, 0x56, 0x94, 0x36, 0x06, 0x2d, 0x46, 0x7d, 0x17, 0xa6, 0x6d, 0x78, 0xff, 0x4e, 0x85, 0x15,
	0xfa, 0x28, 0x73, 0xab, 0xd5, 0x53, 0x0f, 0x7e, 0xee, 0xb0, 0xde, 0xc8, 0xbf, 0x16, 0x7d, 0x62,
	0xfb, 0x5a, 0xd6, 0xbc, 0x44, 0xad, 0x25, 0x91, 0xc4, 0x9a, 0x1b, 0x61, 0x21, 0x0e, 0x0f, 0xc3,
	0xa3, 0x9b, 0xc3, 0xf4, 0xe2, 0x6a, 0x10, 0xfc, 0xbe, 0x1a, 0x3c, 0xaa, 0xa4, 0xfd, 0xfc, 0xa5,
	0x48, 0x4b, 0xd4, 0x59, 0x89, 0xa4, 0x91, 0xd6, 0x97, 0x63, 0x1a, 0x4f, 0x32, 0x3b, 0x6f, 0x80,
	0xd2, 0x13, 0x28, 0xf3, 0xdb, 0x5a, 0xd6, 0xcf, 0x5b, 0x29, 0x17, 0x16, 0xa2, 0x77, 0x6c, 0x4f,
	0x8b, 0x19, 0xff, 0x8a, 0x56, 0xd6, 0x15, 0x6f, 0x70, 0x0a, 0x26, 0xbe, 0xd6, 0x09, 0xef, 0x6b,
	0x31, 0x7b, 0xe3, 0x99, 0x91, 0x53, 0xfe, 0x4d, 0x4e, 0xa0, 0xce, 0xf8, 0x18, 0x14, 0x54, 0xc2,
	0x4a, 0xac, 0xe3, 0xad, 0xff, 0xc6, 0x5f, 0xd6, 0xd6, 0x4f, 0x7e, 0x0a, 0xea, 0xec, 0xa4, 0x85,
	0xa2, 0x87, 0xac, 0x5f, 0x18, 0x9c, 0x80, 0xe1, 0x62, 0x3c, 0x36, 0x40, 0x14, 0x6f, 0x3b, 0x3a,
	0xdf, 0x5d, 0xed, 0x3e, 0x5b, 0x6d, 0x46, 0x8f, 0xd9, 0x5d, 0x8b, 0x56, 0x28, 0x5e, 0x28, 0x2c,
	0x27, 0xc4, 0x1b, 0x30, 0x7c, 0x0e, 0xc2, 0xc4, 0x3b, 0x87, 0xe1, 0xd1, 0x56, 0x1e, 0xf9, 0x70,
	0xe8, 0xb3, 0x11, 0x98, 0xf7, 0x20, 0x4c, 0xf4, 0x91, 0xed, 0x4f, 0x05, 0x69, 0xee, 0x3e, 0x8c,
	0x12, 0x05, 0x28, 0x4e, 0xf2, 0x1c, 0xe2, 0x5e, 0xa7, 0xc9, 0xf7, 0x1c, 0xf5, 0x5a, 0xcc, 0x5e,
	0x39, 0xe8, 0x54, 0x9e, 0x43, 0x94, 0xb3, 0xdd, 0x96, 0xf7, 0xf0, 0xf5, 0x4e, 0xf0, 0xad, 0x35,
	0xec, 0x4d, 0xc5, 0x0e, 0x5a, 0xb3, 0x31, 0xd8, 0x20, 0x09, 0xc5, 0xfd, 0x8e, 0x6f, 0xb8, 0xd1,
	0xa9, 0xe1, 0xde, 0xba, 0x61, 0xb4, 0x06, 0xdf, 0x0a, 0xd2, 0xae, 0xed, 0xe9, 0xf6, 0xf7, 0x1f,
	0x83, 0x60, 0xf8, 0xe2, 0x62, 0x91, 0x84, 0x97, 0x8b, 0x24, 0xfc, 0xb3, 0x48, 0xc2, 0x6f, 0xcb,
	0x24, 0xb8, 0x5c, 0x26, 0xc1, 0xaf, 0x65, 0x12, 0x7c, 0x48, 0x37, 0x0a, 0xdc, 0x81, 0x3f, 0xae,
	0xc1, 0x4e, 0xd1, 0x4c, 0xfc, 0x4d, 0x36, 0xdb, 0xf8, 0x39, 0x7c, 0x59, 0xd1, 0xf3, 0xc7, 0xfe,
	0xc9, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb8, 0x5a, 0x52, 0x23, 0x3b, 0x03, 0x00, 0x00,
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
	{
		size := m.WasmMaxProposalWasmSize.Size()
		i -= size
		if _, err := m.WasmMaxProposalWasmSize.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.WasmMaxSize.Size()
		i -= size
		if _, err := m.WasmMaxSize.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.WasmMaxLabelSize.Size()
		i -= size
		if _, err := m.WasmMaxLabelSize.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if m.TotalBlocksPerYear != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.TotalBlocksPerYear))
		i--
		dAtA[i] = 0x28
	}
	if len(m.BrokerAddress) > 0 {
		i -= len(m.BrokerAddress)
		copy(dAtA[i:], m.BrokerAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.BrokerAddress)))
		i--
		dAtA[i] = 0x22
	}
	{
		size := m.MinSelfDelegation.Size()
		i -= size
		if _, err := m.MinSelfDelegation.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.MaxVotingPower.Size()
		i -= size
		if _, err := m.MaxVotingPower.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.MinCommissionRate.Size()
		i -= size
		if _, err := m.MinCommissionRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *LegacyParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LegacyParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *LegacyParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.BrokerAddress) > 0 {
		i -= len(m.BrokerAddress)
		copy(dAtA[i:], m.BrokerAddress)
		i = encodeVarintParams(dAtA, i, uint64(len(m.BrokerAddress)))
		i--
		dAtA[i] = 0x22
	}
	{
		size := m.MinSelfDelegation.Size()
		i -= size
		if _, err := m.MinSelfDelegation.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.MaxVotingPower.Size()
		i -= size
		if _, err := m.MaxVotingPower.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.MinCommissionRate.Size()
		i -= size
		if _, err := m.MinCommissionRate.MarshalTo(dAtA[i:]); err != nil {
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
	l = m.MinCommissionRate.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MaxVotingPower.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MinSelfDelegation.Size()
	n += 1 + l + sovParams(uint64(l))
	l = len(m.BrokerAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.TotalBlocksPerYear != 0 {
		n += 1 + sovParams(uint64(m.TotalBlocksPerYear))
	}
	l = m.WasmMaxLabelSize.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.WasmMaxSize.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.WasmMaxProposalWasmSize.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func (m *LegacyParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.MinCommissionRate.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MaxVotingPower.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MinSelfDelegation.Size()
	n += 1 + l + sovParams(uint64(l))
	l = len(m.BrokerAddress)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field MinCommissionRate", wireType)
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
			if err := m.MinCommissionRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxVotingPower", wireType)
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
			if err := m.MaxVotingPower.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinSelfDelegation", wireType)
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
			if err := m.MinSelfDelegation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BrokerAddress", wireType)
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
			m.BrokerAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalBlocksPerYear", wireType)
			}
			m.TotalBlocksPerYear = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TotalBlocksPerYear |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WasmMaxLabelSize", wireType)
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
			if err := m.WasmMaxLabelSize.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WasmMaxSize", wireType)
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
			if err := m.WasmMaxSize.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field WasmMaxProposalWasmSize", wireType)
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
			if err := m.WasmMaxProposalWasmSize.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
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
func (m *LegacyParams) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: LegacyParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LegacyParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinCommissionRate", wireType)
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
			if err := m.MinCommissionRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxVotingPower", wireType)
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
			if err := m.MaxVotingPower.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinSelfDelegation", wireType)
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
			if err := m.MinSelfDelegation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BrokerAddress", wireType)
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
			m.BrokerAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
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
