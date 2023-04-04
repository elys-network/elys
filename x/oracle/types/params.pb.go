// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/oracle/params.proto

package types

import (
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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
	BandChannelSource string                                   `protobuf:"bytes,1,opt,name=band_channel_source,json=bandChannelSource,proto3" json:"band_channel_source,omitempty"`
	OracleScriptID    uint64                                   `protobuf:"varint,2,opt,name=oracle_script_id,json=oracleScriptId,proto3" json:"oracle_script_id,omitempty" yaml:"oracle_script_id"`
	Multiplier        uint64                                   `protobuf:"varint,3,opt,name=multiplier,proto3" json:"multiplier,omitempty"`
	AskCount          uint64                                   `protobuf:"varint,4,opt,name=ask_count,json=askCount,proto3" json:"ask_count,omitempty"`
	MinCount          uint64                                   `protobuf:"varint,5,opt,name=min_count,json=minCount,proto3" json:"min_count,omitempty"`
	FeeLimit          github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,6,rep,name=fee_limit,json=feeLimit,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"fee_limit"`
	PrepareGas        uint64                                   `protobuf:"varint,7,opt,name=prepare_gas,json=prepareGas,proto3" json:"prepare_gas,omitempty"`
	ExecuteGas        uint64                                   `protobuf:"varint,8,opt,name=execute_gas,json=executeGas,proto3" json:"execute_gas,omitempty"`
	ClientID          string                                   `protobuf:"bytes,9,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	BandEpoch         string                                   `protobuf:"bytes,10,opt,name=band_epoch,json=bandEpoch,proto3" json:"band_epoch,omitempty"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7d7a7bc7ab1ff79, []int{0}
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

func (m *Params) GetBandChannelSource() string {
	if m != nil {
		return m.BandChannelSource
	}
	return ""
}

func (m *Params) GetOracleScriptID() uint64 {
	if m != nil {
		return m.OracleScriptID
	}
	return 0
}

func (m *Params) GetMultiplier() uint64 {
	if m != nil {
		return m.Multiplier
	}
	return 0
}

func (m *Params) GetAskCount() uint64 {
	if m != nil {
		return m.AskCount
	}
	return 0
}

func (m *Params) GetMinCount() uint64 {
	if m != nil {
		return m.MinCount
	}
	return 0
}

func (m *Params) GetFeeLimit() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.FeeLimit
	}
	return nil
}

func (m *Params) GetPrepareGas() uint64 {
	if m != nil {
		return m.PrepareGas
	}
	return 0
}

func (m *Params) GetExecuteGas() uint64 {
	if m != nil {
		return m.ExecuteGas
	}
	return 0
}

func (m *Params) GetClientID() string {
	if m != nil {
		return m.ClientID
	}
	return ""
}

func (m *Params) GetBandEpoch() string {
	if m != nil {
		return m.BandEpoch
	}
	return ""
}

func init() {
	proto.RegisterType((*Params)(nil), "elysnetwork.elys.oracle.Params")
}

func init() { proto.RegisterFile("elys/oracle/params.proto", fileDescriptor_f7d7a7bc7ab1ff79) }

var fileDescriptor_f7d7a7bc7ab1ff79 = []byte{
	// 469 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x92, 0x31, 0x8f, 0xd3, 0x30,
	0x14, 0xc7, 0x13, 0xae, 0x94, 0xc4, 0x87, 0x4e, 0x10, 0x90, 0x2e, 0x1c, 0x22, 0xa9, 0x3a, 0x15,
	0xa1, 0x26, 0x1c, 0x6c, 0x37, 0xb6, 0x77, 0x42, 0x95, 0x90, 0x40, 0xb9, 0x01, 0x89, 0x25, 0x72,
	0x9c, 0x77, 0xad, 0xd5, 0xc4, 0x8e, 0x6c, 0x07, 0xae, 0xdf, 0x82, 0x91, 0x91, 0x99, 0x4f, 0x72,
	0xe3, 0x8d, 0x4c, 0x05, 0xa5, 0x03, 0x3b, 0x9f, 0x00, 0xd9, 0x0e, 0x52, 0xc5, 0x14, 0xfb, 0xff,
	0xfb, 0xfb, 0x45, 0xef, 0xff, 0x1e, 0x0a, 0xa1, 0xda, 0xc8, 0x94, 0x0b, 0x4c, 0x2a, 0x48, 0x1b,
	0x2c, 0x70, 0x2d, 0x93, 0x46, 0x70, 0xc5, 0x83, 0x63, 0x4d, 0x18, 0xa8, 0xcf, 0x5c, 0xac, 0x13,
	0x7d, 0x4e, 0xac, 0xeb, 0xe4, 0xf1, 0x92, 0x2f, 0xb9, 0xf1, 0xa4, 0xfa, 0x64, 0xed, 0x27, 0x11,
	0xe1, 0xb2, 0xe6, 0x32, 0x2d, 0xb0, 0x84, 0xf4, 0xd3, 0x69, 0x01, 0x0a, 0x9f, 0xa6, 0x84, 0x53,
	0x66, 0xf9, 0xf8, 0xf7, 0x01, 0x1a, 0xbe, 0x37, 0xf5, 0x83, 0x04, 0x3d, 0x2a, 0x30, 0x2b, 0x73,
	0xb2, 0xc2, 0x8c, 0x41, 0x95, 0x4b, 0xde, 0x0a, 0x02, 0xa1, 0x3b, 0x72, 0x27, 0x7e, 0xf6, 0x50,
	0xa3, 0xb9, 0x25, 0x97, 0x06, 0x04, 0x1f, 0xd0, 0x03, 0xfb, 0xeb, 0x5c, 0x12, 0x41, 0x1b, 0x95,
	0xd3, 0x32, 0xbc, 0x33, 0x72, 0x27, 0x83, 0xd9, 0xb4, 0xdb, 0xc6, 0x47, 0xef, 0x0c, 0xbb, 0x34,
	0x68, 0x71, 0xfe, 0x67, 0x1b, 0x1f, 0x6f, 0x70, 0x5d, 0x9d, 0x8d, 0xff, 0x7f, 0x33, 0xce, 0x8e,
	0xf8, 0xbe, 0xb5, 0x0c, 0x22, 0x84, 0xea, 0xb6, 0x52, 0xb4, 0xa9, 0x28, 0x88, 0xf0, 0x40, 0x97,
	0xcc, 0xf6, 0x94, 0xe0, 0x29, 0xf2, 0xb1, 0x5c, 0xe7, 0x84, 0xb7, 0x4c, 0x85, 0x03, 0x83, 0x3d,
	0x2c, 0xd7, 0x73, 0x7d, 0xd7, 0xb0, 0xa6, 0xac, 0x87, 0x77, 0x2d, 0xac, 0x29, 0xb3, 0x70, 0x85,
	0xfc, 0x2b, 0x80, 0xbc, 0xa2, 0x35, 0x55, 0xe1, 0x70, 0x74, 0x30, 0x39, 0x7c, 0xf5, 0x24, 0xb1,
	0x09, 0x25, 0x3a, 0xa1, 0xa4, 0x4f, 0x28, 0x99, 0x73, 0xca, 0x66, 0x2f, 0x6f, 0xb6, 0xb1, 0xf3,
	0xfd, 0x67, 0x3c, 0x59, 0x52, 0xb5, 0x6a, 0x8b, 0x84, 0xf0, 0x3a, 0xed, 0xe3, 0xb4, 0x9f, 0xa9,
	0x2c, 0xd7, 0xa9, 0xda, 0x34, 0x20, 0xcd, 0x03, 0x99, 0x79, 0x57, 0x00, 0x6f, 0x75, 0xf1, 0x20,
	0x46, 0x87, 0x8d, 0x80, 0x06, 0x0b, 0xc8, 0x97, 0x58, 0x86, 0xf7, 0x6c, 0x13, 0xbd, 0xf4, 0x06,
	0x4b, 0x6d, 0x80, 0x6b, 0x20, 0xad, 0xb2, 0x06, 0xcf, 0x1a, 0x7a, 0x49, 0x1b, 0x9e, 0x23, 0x9f,
	0x54, 0x14, 0x98, 0xc9, 0xd5, 0xd7, 0x43, 0x98, 0xdd, 0xef, 0xb6, 0xb1, 0x37, 0x37, 0xe2, 0xe2,
	0x3c, 0xf3, 0x2c, 0x5e, 0x94, 0xc1, 0x33, 0x84, 0xcc, 0xe4, 0xa0, 0xe1, 0x64, 0x15, 0x22, 0x33,
	0x30, 0x5f, 0x2b, 0x17, 0x5a, 0x38, 0x1b, 0x7c, 0xfd, 0x16, 0x3b, 0xb3, 0x8b, 0x9b, 0x2e, 0x72,
	0x6f, 0xbb, 0xc8, 0xfd, 0xd5, 0x45, 0xee, 0x97, 0x5d, 0xe4, 0xdc, 0xee, 0x22, 0xe7, 0xc7, 0x2e,
	0x72, 0x3e, 0xbe, 0xd8, 0xeb, 0x4f, 0x6f, 0xd4, 0xb4, 0x5f, 0x2f, 0x73, 0x49, 0xaf, 0xff, 0xad,
	0xa1, 0x69, 0xb4, 0x18, 0x9a, 0xbd, 0x79, 0xfd, 0x37, 0x00, 0x00, 0xff, 0xff, 0x18, 0xde, 0x1c,
	0x52, 0xa2, 0x02, 0x00, 0x00,
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
	if len(m.BandEpoch) > 0 {
		i -= len(m.BandEpoch)
		copy(dAtA[i:], m.BandEpoch)
		i = encodeVarintParams(dAtA, i, uint64(len(m.BandEpoch)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.ClientID) > 0 {
		i -= len(m.ClientID)
		copy(dAtA[i:], m.ClientID)
		i = encodeVarintParams(dAtA, i, uint64(len(m.ClientID)))
		i--
		dAtA[i] = 0x4a
	}
	if m.ExecuteGas != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.ExecuteGas))
		i--
		dAtA[i] = 0x40
	}
	if m.PrepareGas != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.PrepareGas))
		i--
		dAtA[i] = 0x38
	}
	if len(m.FeeLimit) > 0 {
		for iNdEx := len(m.FeeLimit) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.FeeLimit[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if m.MinCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MinCount))
		i--
		dAtA[i] = 0x28
	}
	if m.AskCount != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.AskCount))
		i--
		dAtA[i] = 0x20
	}
	if m.Multiplier != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.Multiplier))
		i--
		dAtA[i] = 0x18
	}
	if m.OracleScriptID != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.OracleScriptID))
		i--
		dAtA[i] = 0x10
	}
	if len(m.BandChannelSource) > 0 {
		i -= len(m.BandChannelSource)
		copy(dAtA[i:], m.BandChannelSource)
		i = encodeVarintParams(dAtA, i, uint64(len(m.BandChannelSource)))
		i--
		dAtA[i] = 0xa
	}
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
	l = len(m.BandChannelSource)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.OracleScriptID != 0 {
		n += 1 + sovParams(uint64(m.OracleScriptID))
	}
	if m.Multiplier != 0 {
		n += 1 + sovParams(uint64(m.Multiplier))
	}
	if m.AskCount != 0 {
		n += 1 + sovParams(uint64(m.AskCount))
	}
	if m.MinCount != 0 {
		n += 1 + sovParams(uint64(m.MinCount))
	}
	if len(m.FeeLimit) > 0 {
		for _, e := range m.FeeLimit {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if m.PrepareGas != 0 {
		n += 1 + sovParams(uint64(m.PrepareGas))
	}
	if m.ExecuteGas != 0 {
		n += 1 + sovParams(uint64(m.ExecuteGas))
	}
	l = len(m.ClientID)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.BandEpoch)
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
				return fmt.Errorf("proto: wrong wireType = %d for field BandChannelSource", wireType)
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
			m.BandChannelSource = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OracleScriptID", wireType)
			}
			m.OracleScriptID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OracleScriptID |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Multiplier", wireType)
			}
			m.Multiplier = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Multiplier |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AskCount", wireType)
			}
			m.AskCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AskCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinCount", wireType)
			}
			m.MinCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MinCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeLimit", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
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
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FeeLimit = append(m.FeeLimit, types.Coin{})
			if err := m.FeeLimit[len(m.FeeLimit)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PrepareGas", wireType)
			}
			m.PrepareGas = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PrepareGas |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExecuteGas", wireType)
			}
			m.ExecuteGas = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExecuteGas |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClientID", wireType)
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
			m.ClientID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BandEpoch", wireType)
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
			m.BandEpoch = string(dAtA[iNdEx:postIndex])
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
