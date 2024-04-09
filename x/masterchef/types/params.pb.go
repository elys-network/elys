// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/masterchef/params.proto

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
	LpIncentives *IncentiveInfo `protobuf:"bytes,1,opt,name=lp_incentives,json=lpIncentives,proto3" json:"lp_incentives,omitempty"`
	// Dex revenue percent for lps, `100 - reward_portion_for_lps - reward_portion_for_stakers = revenue percent for protocol`.
	RewardPortionForLps github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=reward_portion_for_lps,json=rewardPortionForLps,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"reward_portion_for_lps"`
	// Tracking dex rewards given to LPs
	DexRewardsLps DexRewardsTracker `protobuf:"bytes,3,opt,name=dex_rewards_lps,json=dexRewardsLps,proto3" json:"dex_rewards_lps"`
	// Maximum eden reward apr for lps - [0 - 0.3]
	MaxEdenRewardAprLps   github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=max_eden_reward_apr_lps,json=maxEdenRewardAprLps,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"max_eden_reward_apr_lps"`
	SupportedRewardDenoms []*SupportedRewardDenom                `protobuf:"bytes,5,rep,name=supported_reward_denoms,json=supportedRewardDenoms,proto3" json:"supported_reward_denoms,omitempty"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_bc83f5a7f5a55e20, []int{0}
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

func (m *Params) GetLpIncentives() *IncentiveInfo {
	if m != nil {
		return m.LpIncentives
	}
	return nil
}

func (m *Params) GetDexRewardsLps() DexRewardsTracker {
	if m != nil {
		return m.DexRewardsLps
	}
	return DexRewardsTracker{}
}

func (m *Params) GetSupportedRewardDenoms() []*SupportedRewardDenom {
	if m != nil {
		return m.SupportedRewardDenoms
	}
	return nil
}

type SupportedRewardDenom struct {
	Denom     string                                 `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
	MinAmount github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=min_amount,json=minAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"min_amount"`
}

func (m *SupportedRewardDenom) Reset()         { *m = SupportedRewardDenom{} }
func (m *SupportedRewardDenom) String() string { return proto.CompactTextString(m) }
func (*SupportedRewardDenom) ProtoMessage()    {}
func (*SupportedRewardDenom) Descriptor() ([]byte, []int) {
	return fileDescriptor_bc83f5a7f5a55e20, []int{1}
}
func (m *SupportedRewardDenom) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SupportedRewardDenom) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_SupportedRewardDenom.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *SupportedRewardDenom) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SupportedRewardDenom.Merge(m, src)
}
func (m *SupportedRewardDenom) XXX_Size() int {
	return m.Size()
}
func (m *SupportedRewardDenom) XXX_DiscardUnknown() {
	xxx_messageInfo_SupportedRewardDenom.DiscardUnknown(m)
}

var xxx_messageInfo_SupportedRewardDenom proto.InternalMessageInfo

func (m *SupportedRewardDenom) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

func init() {
	proto.RegisterType((*Params)(nil), "elys.masterchef.Params")
	proto.RegisterType((*SupportedRewardDenom)(nil), "elys.masterchef.SupportedRewardDenom")
}

func init() { proto.RegisterFile("elys/masterchef/params.proto", fileDescriptor_bc83f5a7f5a55e20) }

var fileDescriptor_bc83f5a7f5a55e20 = []byte{
	// 454 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0x4f, 0x6b, 0xdb, 0x30,
	0x18, 0xc6, 0xed, 0x25, 0x2d, 0x44, 0x5d, 0x29, 0x78, 0xd9, 0x1a, 0xc2, 0x70, 0x42, 0x60, 0x23,
	0x97, 0xda, 0xd0, 0xdd, 0x76, 0x6b, 0x96, 0x0d, 0x02, 0x1b, 0x04, 0x6f, 0xa7, 0xc1, 0x10, 0xaa,
	0xf5, 0x26, 0x35, 0xb1, 0xfe, 0x20, 0x29, 0xab, 0xcb, 0xbe, 0xc4, 0x8e, 0x3b, 0xee, 0xe3, 0xf4,
	0xb6, 0x1e, 0xc7, 0x0e, 0x65, 0x24, 0x5f, 0x64, 0x58, 0x72, 0xb2, 0x90, 0xf4, 0x94, 0x93, 0x2d,
	0x3d, 0x8f, 0x7f, 0xcf, 0x23, 0xeb, 0x45, 0xcf, 0x21, 0xbf, 0xd1, 0x31, 0x23, 0xda, 0x80, 0x4a,
	0xaf, 0x60, 0x12, 0x4b, 0xa2, 0x08, 0xd3, 0x91, 0x54, 0xc2, 0x88, 0xe0, 0xa4, 0x54, 0xa3, 0xff,
	0x6a, 0xbb, 0x39, 0x15, 0x53, 0x61, 0xb5, 0xb8, 0x7c, 0x73, 0xb6, 0x76, 0x67, 0x1b, 0x92, 0xf1,
	0x14, 0xb8, 0xc9, 0xbe, 0x42, 0x65, 0x68, 0xef, 0xa4, 0x08, 0x91, 0x57, 0x5a, 0x7f, 0x5b, 0xa3,
	0x50, 0x60, 0x05, 0xd7, 0x44, 0x51, 0x8d, 0x8d, 0x22, 0x33, 0x50, 0xce, 0xd9, 0xfb, 0x55, 0x43,
	0x87, 0x63, 0x5b, 0x2f, 0x78, 0x83, 0x8e, 0x73, 0x89, 0xd7, 0x31, 0xba, 0xe5, 0x77, 0xfd, 0xfe,
	0xd1, 0x79, 0x18, 0x6d, 0x15, 0x8e, 0x46, 0x2b, 0xcb, 0x88, 0x4f, 0x44, 0xf2, 0x38, 0x97, 0xeb,
	0x0d, 0x1d, 0xa4, 0xe8, 0x99, 0xcb, 0xc1, 0x52, 0x28, 0x93, 0x09, 0x8e, 0x27, 0x42, 0xe1, 0x5c,
	0xea, 0xd6, 0xa3, 0xae, 0xdf, 0x6f, 0x0c, 0xa2, 0xdb, 0xfb, 0x8e, 0xf7, 0xe7, 0xbe, 0xf3, 0x72,
	0x9a, 0x99, 0xab, 0xf9, 0x65, 0x94, 0x0a, 0x16, 0xa7, 0x42, 0x33, 0xa1, 0xab, 0xc7, 0x99, 0xa6,
	0xb3, 0xd8, 0xdc, 0x48, 0xd0, 0xd1, 0x10, 0xd2, 0xe4, 0x89, 0xa3, 0x8d, 0x1d, 0xec, 0x9d, 0x50,
	0xef, 0xa5, 0x0e, 0xc6, 0xe8, 0x64, 0xf3, 0x40, 0x25, 0xbd, 0x66, 0xbb, 0xf6, 0x76, 0xba, 0x0e,
	0xa1, 0x48, 0x9c, 0xed, 0x93, 0x22, 0xe9, 0x0c, 0xd4, 0xa0, 0x5e, 0x36, 0x48, 0x8e, 0xe9, 0x5a,
	0x28, 0x89, 0x14, 0x9d, 0x32, 0x52, 0x60, 0xa0, 0xc0, 0x2b, 0x2c, 0x26, 0xd2, 0xf5, 0xae, 0xef,
	0xd7, 0x9b, 0x91, 0xe2, 0x2d, 0x05, 0xee, 0x32, 0x2e, 0xa4, 0xed, 0xfd, 0x05, 0x9d, 0xea, 0xb9,
	0x2c, 0x7f, 0x0c, 0xd0, 0x55, 0x0c, 0x05, 0x2e, 0x98, 0x6e, 0x1d, 0x74, 0x6b, 0xfd, 0xa3, 0xf3,
	0x17, 0x3b, 0xfd, 0x3f, 0xae, 0xfc, 0x0e, 0x34, 0x2c, 0xdd, 0xc9, 0x53, 0xfd, 0xc0, 0xae, 0x7e,
	0x5d, 0xff, 0xf1, 0xb3, 0xe3, 0xf5, 0xbe, 0xa1, 0xe6, 0x43, 0x1f, 0x05, 0x4d, 0x74, 0x60, 0xb3,
	0xec, 0xb5, 0x36, 0x12, 0xb7, 0x08, 0x3e, 0x20, 0xc4, 0x32, 0x8e, 0x09, 0x13, 0x73, 0x6e, 0xf6,
	0xb8, 0xa3, 0x11, 0x37, 0x49, 0x83, 0x65, 0xfc, 0xc2, 0x02, 0x06, 0xa3, 0xdb, 0x45, 0xe8, 0xdf,
	0x2d, 0x42, 0xff, 0xef, 0x22, 0xf4, 0xbf, 0x2f, 0x43, 0xef, 0x6e, 0x19, 0x7a, 0xbf, 0x97, 0xa1,
	0xf7, 0x39, 0xde, 0x80, 0x95, 0x87, 0x3c, 0xe3, 0x60, 0xae, 0x85, 0x9a, 0xd9, 0x45, 0x5c, 0x6c,
	0x0e, 0xab, 0x25, 0x5f, 0x1e, 0xda, 0x01, 0x7d, 0xf5, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xb2, 0xc4,
	0x87, 0x7f, 0x4e, 0x03, 0x00, 0x00,
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
	if len(m.SupportedRewardDenoms) > 0 {
		for iNdEx := len(m.SupportedRewardDenoms) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.SupportedRewardDenoms[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintParams(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	{
		size := m.MaxEdenRewardAprLps.Size()
		i -= size
		if _, err := m.MaxEdenRewardAprLps.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size, err := m.DexRewardsLps.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.RewardPortionForLps.Size()
		i -= size
		if _, err := m.RewardPortionForLps.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.LpIncentives != nil {
		{
			size, err := m.LpIncentives.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintParams(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SupportedRewardDenom) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SupportedRewardDenom) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SupportedRewardDenom) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.MinAmount.Size()
		i -= size
		if _, err := m.MinAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Denom)))
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
	if m.LpIncentives != nil {
		l = m.LpIncentives.Size()
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.RewardPortionForLps.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.DexRewardsLps.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.MaxEdenRewardAprLps.Size()
	n += 1 + l + sovParams(uint64(l))
	if len(m.SupportedRewardDenoms) > 0 {
		for _, e := range m.SupportedRewardDenoms {
			l = e.Size()
			n += 1 + l + sovParams(uint64(l))
		}
	}
	return n
}

func (m *SupportedRewardDenom) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.MinAmount.Size()
	n += 1 + l + sovParams(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field LpIncentives", wireType)
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
			if m.LpIncentives == nil {
				m.LpIncentives = &IncentiveInfo{}
			}
			if err := m.LpIncentives.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardPortionForLps", wireType)
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
			if err := m.RewardPortionForLps.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DexRewardsLps", wireType)
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
			if err := m.DexRewardsLps.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxEdenRewardAprLps", wireType)
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
			if err := m.MaxEdenRewardAprLps.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SupportedRewardDenoms", wireType)
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
			m.SupportedRewardDenoms = append(m.SupportedRewardDenoms, &SupportedRewardDenom{})
			if err := m.SupportedRewardDenoms[len(m.SupportedRewardDenoms)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *SupportedRewardDenom) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: SupportedRewardDenom: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SupportedRewardDenom: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
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
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinAmount", wireType)
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
			if err := m.MinAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
