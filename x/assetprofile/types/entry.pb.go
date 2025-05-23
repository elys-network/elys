// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/assetprofile/entry.proto

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

type Entry struct {
	BaseDenom                string   `protobuf:"bytes,1,opt,name=base_denom,json=baseDenom,proto3" json:"base_denom,omitempty"`
	Decimals                 uint64   `protobuf:"varint,2,opt,name=decimals,proto3" json:"decimals,omitempty"`
	Denom                    string   `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
	Path                     string   `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
	IbcChannelId             string   `protobuf:"bytes,5,opt,name=ibc_channel_id,json=ibcChannelId,proto3" json:"ibc_channel_id,omitempty"`
	IbcCounterpartyChannelId string   `protobuf:"bytes,6,opt,name=ibc_counterparty_channel_id,json=ibcCounterpartyChannelId,proto3" json:"ibc_counterparty_channel_id,omitempty"`
	DisplayName              string   `protobuf:"bytes,7,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	DisplaySymbol            string   `protobuf:"bytes,8,opt,name=display_symbol,json=displaySymbol,proto3" json:"display_symbol,omitempty"`
	Network                  string   `protobuf:"bytes,9,opt,name=network,proto3" json:"network,omitempty"`
	Address                  string   `protobuf:"bytes,10,opt,name=address,proto3" json:"address,omitempty"`
	ExternalSymbol           string   `protobuf:"bytes,11,opt,name=external_symbol,json=externalSymbol,proto3" json:"external_symbol,omitempty"`
	TransferLimit            string   `protobuf:"bytes,12,opt,name=transfer_limit,json=transferLimit,proto3" json:"transfer_limit,omitempty"`
	Permissions              []string `protobuf:"bytes,13,rep,name=permissions,proto3" json:"permissions,omitempty"`
	UnitDenom                string   `protobuf:"bytes,14,opt,name=unit_denom,json=unitDenom,proto3" json:"unit_denom,omitempty"`
	IbcCounterpartyDenom     string   `protobuf:"bytes,15,opt,name=ibc_counterparty_denom,json=ibcCounterpartyDenom,proto3" json:"ibc_counterparty_denom,omitempty"`
	IbcCounterpartyChainId   string   `protobuf:"bytes,16,opt,name=ibc_counterparty_chain_id,json=ibcCounterpartyChainId,proto3" json:"ibc_counterparty_chain_id,omitempty"`
	Authority                string   `protobuf:"bytes,17,opt,name=authority,proto3" json:"authority,omitempty"`
	CommitEnabled            bool     `protobuf:"varint,18,opt,name=commit_enabled,json=commitEnabled,proto3" json:"commit_enabled,omitempty"`
	WithdrawEnabled          bool     `protobuf:"varint,19,opt,name=withdraw_enabled,json=withdrawEnabled,proto3" json:"withdraw_enabled,omitempty"`
}

func (m *Entry) Reset()         { *m = Entry{} }
func (m *Entry) String() string { return proto.CompactTextString(m) }
func (*Entry) ProtoMessage()    {}
func (*Entry) Descriptor() ([]byte, []int) {
	return fileDescriptor_138d3d7a28f399dd, []int{0}
}
func (m *Entry) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Entry) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Entry.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Entry) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Entry.Merge(m, src)
}
func (m *Entry) XXX_Size() int {
	return m.Size()
}
func (m *Entry) XXX_DiscardUnknown() {
	xxx_messageInfo_Entry.DiscardUnknown(m)
}

var xxx_messageInfo_Entry proto.InternalMessageInfo

func (m *Entry) GetBaseDenom() string {
	if m != nil {
		return m.BaseDenom
	}
	return ""
}

func (m *Entry) GetDecimals() uint64 {
	if m != nil {
		return m.Decimals
	}
	return 0
}

func (m *Entry) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func (m *Entry) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Entry) GetIbcChannelId() string {
	if m != nil {
		return m.IbcChannelId
	}
	return ""
}

func (m *Entry) GetIbcCounterpartyChannelId() string {
	if m != nil {
		return m.IbcCounterpartyChannelId
	}
	return ""
}

func (m *Entry) GetDisplayName() string {
	if m != nil {
		return m.DisplayName
	}
	return ""
}

func (m *Entry) GetDisplaySymbol() string {
	if m != nil {
		return m.DisplaySymbol
	}
	return ""
}

func (m *Entry) GetNetwork() string {
	if m != nil {
		return m.Network
	}
	return ""
}

func (m *Entry) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Entry) GetExternalSymbol() string {
	if m != nil {
		return m.ExternalSymbol
	}
	return ""
}

func (m *Entry) GetTransferLimit() string {
	if m != nil {
		return m.TransferLimit
	}
	return ""
}

func (m *Entry) GetPermissions() []string {
	if m != nil {
		return m.Permissions
	}
	return nil
}

func (m *Entry) GetUnitDenom() string {
	if m != nil {
		return m.UnitDenom
	}
	return ""
}

func (m *Entry) GetIbcCounterpartyDenom() string {
	if m != nil {
		return m.IbcCounterpartyDenom
	}
	return ""
}

func (m *Entry) GetIbcCounterpartyChainId() string {
	if m != nil {
		return m.IbcCounterpartyChainId
	}
	return ""
}

func (m *Entry) GetAuthority() string {
	if m != nil {
		return m.Authority
	}
	return ""
}

func (m *Entry) GetCommitEnabled() bool {
	if m != nil {
		return m.CommitEnabled
	}
	return false
}

func (m *Entry) GetWithdrawEnabled() bool {
	if m != nil {
		return m.WithdrawEnabled
	}
	return false
}

func init() {
	proto.RegisterType((*Entry)(nil), "elys.assetprofile.Entry")
}

func init() { proto.RegisterFile("elys/assetprofile/entry.proto", fileDescriptor_138d3d7a28f399dd) }

var fileDescriptor_138d3d7a28f399dd = []byte{
	// 494 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x93, 0xcd, 0x72, 0xd3, 0x3e,
	0x14, 0xc5, 0xe3, 0x7f, 0x93, 0x36, 0xb9, 0xf9, 0x6a, 0xf5, 0xef, 0x74, 0xc4, 0x47, 0x3d, 0x81,
	0xa1, 0x43, 0x58, 0x90, 0x2c, 0x28, 0x0b, 0x16, 0x6c, 0x80, 0x2e, 0x3a, 0xc3, 0x74, 0x11, 0x76,
	0x6c, 0x3c, 0xb2, 0x7d, 0x8b, 0x35, 0xd8, 0x92, 0x47, 0x52, 0x48, 0xfd, 0x16, 0xbc, 0x12, 0x3b,
	0x96, 0x5d, 0xb2, 0x64, 0x92, 0x17, 0x61, 0x24, 0xd9, 0x21, 0x50, 0x76, 0xbe, 0xbf, 0x73, 0x8e,
	0x46, 0xd2, 0xb1, 0xe0, 0x14, 0xf3, 0x4a, 0xcf, 0x99, 0xd6, 0x68, 0x4a, 0x25, 0xaf, 0x79, 0x8e,
	0x73, 0x14, 0x46, 0x55, 0xb3, 0x52, 0x49, 0x23, 0xc9, 0x91, 0x95, 0x67, 0xbb, 0xf2, 0xe3, 0x6f,
	0x1d, 0xe8, 0x5c, 0x58, 0x0b, 0x39, 0x05, 0x88, 0x99, 0xc6, 0x28, 0x45, 0x21, 0x0b, 0x1a, 0x4c,
	0x82, 0x69, 0x6f, 0xd1, 0xb3, 0xe4, 0x9d, 0x05, 0xe4, 0x3e, 0x74, 0x53, 0x4c, 0x78, 0xc1, 0x72,
	0x4d, 0xff, 0x9b, 0x04, 0xd3, 0xf6, 0x62, 0x3b, 0x93, 0x63, 0xe8, 0xf8, 0xd4, 0x9e, 0x4b, 0xf9,
	0x81, 0x10, 0x68, 0x97, 0xcc, 0x64, 0xb4, 0xed, 0xa0, 0xfb, 0x26, 0x4f, 0x60, 0xc4, 0xe3, 0x24,
	0x4a, 0x32, 0x26, 0x04, 0xe6, 0x11, 0x4f, 0x69, 0xc7, 0xa9, 0x03, 0x1e, 0x27, 0x6f, 0x3d, 0xbc,
	0x4c, 0xc9, 0x6b, 0x78, 0xe0, 0x5c, 0x72, 0x29, 0x0c, 0xaa, 0x92, 0x29, 0x53, 0xed, 0x46, 0xf6,
	0x5d, 0x84, 0xda, 0xc8, 0x8e, 0xe3, 0x77, 0xfc, 0x11, 0x0c, 0x52, 0xae, 0xcb, 0x9c, 0x55, 0x91,
	0x60, 0x05, 0xd2, 0x03, 0xe7, 0xef, 0xd7, 0xec, 0x8a, 0x15, 0x48, 0xce, 0x60, 0xd4, 0x58, 0x74,
	0x55, 0xc4, 0x32, 0xa7, 0x5d, 0x67, 0x1a, 0xd6, 0xf4, 0x83, 0x83, 0x84, 0xc2, 0x81, 0x40, 0xb3,
	0x92, 0xea, 0x33, 0xed, 0x39, 0xbd, 0x19, 0xad, 0xc2, 0xd2, 0x54, 0xa1, 0xd6, 0x14, 0xbc, 0x52,
	0x8f, 0xe4, 0x29, 0x8c, 0xf1, 0xc6, 0xa0, 0x12, 0x2c, 0x6f, 0xd6, 0xee, 0x3b, 0xc7, 0xa8, 0xc1,
	0xf5, 0xe2, 0x67, 0x30, 0x32, 0x8a, 0x09, 0x7d, 0x8d, 0x2a, 0xca, 0x79, 0xc1, 0x0d, 0x1d, 0xf8,
	0x3d, 0x34, 0xf4, 0xbd, 0x85, 0x64, 0x02, 0xfd, 0x12, 0x55, 0xc1, 0xb5, 0xe6, 0x52, 0x68, 0x3a,
	0x9c, 0xec, 0xd9, 0xc3, 0xec, 0x20, 0xdb, 0xdc, 0x52, 0x70, 0x53, 0x37, 0x37, 0xf2, 0xcd, 0x59,
	0xe2, 0x9b, 0x3b, 0x87, 0x93, 0x3b, 0xb7, 0xe9, 0xad, 0x63, 0x67, 0x3d, 0xfe, 0xeb, 0x22, 0x7d,
	0xea, 0x15, 0xdc, 0xfb, 0x57, 0x07, 0x5c, 0xd8, 0x06, 0x0e, 0x5d, 0xf0, 0xe4, 0x6e, 0x03, 0x5c,
	0x5c, 0xa6, 0xe4, 0x21, 0xf4, 0xd8, 0xd2, 0x64, 0x52, 0x71, 0x53, 0xd1, 0x23, 0xbf, 0x9d, 0x2d,
	0xb0, 0xc7, 0x4e, 0x64, 0x51, 0x70, 0x13, 0xa1, 0x60, 0x71, 0x8e, 0x29, 0x25, 0x93, 0x60, 0xda,
	0x5d, 0x0c, 0x3d, 0xbd, 0xf0, 0x90, 0x3c, 0x83, 0xc3, 0x15, 0x37, 0x59, 0xaa, 0xd8, 0x6a, 0x6b,
	0xfc, 0xdf, 0x19, 0xc7, 0x0d, 0xaf, 0xad, 0x6f, 0xae, 0xbe, 0xaf, 0xc3, 0xe0, 0x76, 0x1d, 0x06,
	0x3f, 0xd7, 0x61, 0xf0, 0x75, 0x13, 0xb6, 0x6e, 0x37, 0x61, 0xeb, 0xc7, 0x26, 0x6c, 0x7d, 0x3c,
	0xff, 0xc4, 0x4d, 0xb6, 0x8c, 0x67, 0x89, 0x2c, 0xe6, 0xf6, 0xdf, 0x7f, 0x5e, 0xd7, 0xe7, 0x86,
	0xf9, 0x97, 0x97, 0xf3, 0x9b, 0x3f, 0x1f, 0x8b, 0xa9, 0x4a, 0xd4, 0xf1, 0xbe, 0x7b, 0x2d, 0x2f,
	0x7e, 0x05, 0x00, 0x00, 0xff, 0xff, 0xe6, 0x10, 0x5f, 0x69, 0x4e, 0x03, 0x00, 0x00,
}

func (m *Entry) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Entry) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Entry) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.WithdrawEnabled {
		i--
		if m.WithdrawEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x98
	}
	if m.CommitEnabled {
		i--
		if m.CommitEnabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x90
	}
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x8a
	}
	if len(m.IbcCounterpartyChainId) > 0 {
		i -= len(m.IbcCounterpartyChainId)
		copy(dAtA[i:], m.IbcCounterpartyChainId)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.IbcCounterpartyChainId)))
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x82
	}
	if len(m.IbcCounterpartyDenom) > 0 {
		i -= len(m.IbcCounterpartyDenom)
		copy(dAtA[i:], m.IbcCounterpartyDenom)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.IbcCounterpartyDenom)))
		i--
		dAtA[i] = 0x7a
	}
	if len(m.UnitDenom) > 0 {
		i -= len(m.UnitDenom)
		copy(dAtA[i:], m.UnitDenom)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.UnitDenom)))
		i--
		dAtA[i] = 0x72
	}
	if len(m.Permissions) > 0 {
		for iNdEx := len(m.Permissions) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Permissions[iNdEx])
			copy(dAtA[i:], m.Permissions[iNdEx])
			i = encodeVarintEntry(dAtA, i, uint64(len(m.Permissions[iNdEx])))
			i--
			dAtA[i] = 0x6a
		}
	}
	if len(m.TransferLimit) > 0 {
		i -= len(m.TransferLimit)
		copy(dAtA[i:], m.TransferLimit)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.TransferLimit)))
		i--
		dAtA[i] = 0x62
	}
	if len(m.ExternalSymbol) > 0 {
		i -= len(m.ExternalSymbol)
		copy(dAtA[i:], m.ExternalSymbol)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.ExternalSymbol)))
		i--
		dAtA[i] = 0x5a
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.Network) > 0 {
		i -= len(m.Network)
		copy(dAtA[i:], m.Network)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.Network)))
		i--
		dAtA[i] = 0x4a
	}
	if len(m.DisplaySymbol) > 0 {
		i -= len(m.DisplaySymbol)
		copy(dAtA[i:], m.DisplaySymbol)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.DisplaySymbol)))
		i--
		dAtA[i] = 0x42
	}
	if len(m.DisplayName) > 0 {
		i -= len(m.DisplayName)
		copy(dAtA[i:], m.DisplayName)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.DisplayName)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.IbcCounterpartyChannelId) > 0 {
		i -= len(m.IbcCounterpartyChannelId)
		copy(dAtA[i:], m.IbcCounterpartyChannelId)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.IbcCounterpartyChannelId)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.IbcChannelId) > 0 {
		i -= len(m.IbcChannelId)
		copy(dAtA[i:], m.IbcChannelId)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.IbcChannelId)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Path) > 0 {
		i -= len(m.Path)
		copy(dAtA[i:], m.Path)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.Path)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Decimals != 0 {
		i = encodeVarintEntry(dAtA, i, uint64(m.Decimals))
		i--
		dAtA[i] = 0x10
	}
	if len(m.BaseDenom) > 0 {
		i -= len(m.BaseDenom)
		copy(dAtA[i:], m.BaseDenom)
		i = encodeVarintEntry(dAtA, i, uint64(len(m.BaseDenom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintEntry(dAtA []byte, offset int, v uint64) int {
	offset -= sovEntry(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Entry) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.BaseDenom)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	if m.Decimals != 0 {
		n += 1 + sovEntry(uint64(m.Decimals))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	l = len(m.Path)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	l = len(m.IbcChannelId)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	l = len(m.IbcCounterpartyChannelId)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	l = len(m.DisplayName)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	l = len(m.DisplaySymbol)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	l = len(m.Network)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	l = len(m.ExternalSymbol)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	l = len(m.TransferLimit)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	if len(m.Permissions) > 0 {
		for _, s := range m.Permissions {
			l = len(s)
			n += 1 + l + sovEntry(uint64(l))
		}
	}
	l = len(m.UnitDenom)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	l = len(m.IbcCounterpartyDenom)
	if l > 0 {
		n += 1 + l + sovEntry(uint64(l))
	}
	l = len(m.IbcCounterpartyChainId)
	if l > 0 {
		n += 2 + l + sovEntry(uint64(l))
	}
	l = len(m.Authority)
	if l > 0 {
		n += 2 + l + sovEntry(uint64(l))
	}
	if m.CommitEnabled {
		n += 3
	}
	if m.WithdrawEnabled {
		n += 3
	}
	return n
}

func sovEntry(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEntry(x uint64) (n int) {
	return sovEntry(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Entry) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEntry
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
			return fmt.Errorf("proto: Entry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Entry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BaseDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BaseDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Decimals", wireType)
			}
			m.Decimals = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Decimals |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Path", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Path = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IbcChannelId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IbcChannelId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IbcCounterpartyChannelId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IbcCounterpartyChannelId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DisplayName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DisplayName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DisplaySymbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DisplaySymbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Network", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Network = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExternalSymbol", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExternalSymbol = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransferLimit", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TransferLimit = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Permissions", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Permissions = append(m.Permissions, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UnitDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UnitDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IbcCounterpartyDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IbcCounterpartyDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 16:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field IbcCounterpartyChainId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.IbcCounterpartyChainId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 17:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
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
				return ErrInvalidLengthEntry
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEntry
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 18:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommitEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.CommitEnabled = bool(v != 0)
		case 19:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field WithdrawEnabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEntry
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.WithdrawEnabled = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipEntry(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEntry
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
func skipEntry(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEntry
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
					return 0, ErrIntOverflowEntry
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
					return 0, ErrIntOverflowEntry
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
				return 0, ErrInvalidLengthEntry
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEntry
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEntry
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEntry        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEntry          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEntry = fmt.Errorf("proto: unexpected end of group")
)
