// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/estaking/params.proto

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

type LegacyParams struct {
	StakeIncentives *LegacyIncentiveInfo `protobuf:"bytes,1,opt,name=stake_incentives,json=stakeIncentives,proto3" json:"stake_incentives,omitempty"`
	EdenCommitVal   string               `protobuf:"bytes,2,opt,name=eden_commit_val,json=edenCommitVal,proto3" json:"eden_commit_val,omitempty"`
	EdenbCommitVal  string               `protobuf:"bytes,3,opt,name=edenb_commit_val,json=edenbCommitVal,proto3" json:"edenb_commit_val,omitempty"`
	// Maximum eden reward apr for stakers - [0 - 0.3]
	MaxEdenRewardAprStakers github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=max_eden_reward_apr_stakers,json=maxEdenRewardAprStakers,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"max_eden_reward_apr_stakers"`
	EdenBoostApr            github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,6,opt,name=eden_boost_apr,json=edenBoostApr,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"eden_boost_apr"`
	// Tracking dex rewards given to stakers
	DexRewardsStakers LegacyDexRewardsTracker `protobuf:"bytes,7,opt,name=dex_rewards_stakers,json=dexRewardsStakers,proto3" json:"dex_rewards_stakers"`
}

func (m *LegacyParams) Reset()      { *m = LegacyParams{} }
func (*LegacyParams) ProtoMessage() {}
func (*LegacyParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_66041162e1ecb63b, []int{0}
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

func (m *LegacyParams) GetStakeIncentives() *LegacyIncentiveInfo {
	if m != nil {
		return m.StakeIncentives
	}
	return nil
}

func (m *LegacyParams) GetEdenCommitVal() string {
	if m != nil {
		return m.EdenCommitVal
	}
	return ""
}

func (m *LegacyParams) GetEdenbCommitVal() string {
	if m != nil {
		return m.EdenbCommitVal
	}
	return ""
}

func (m *LegacyParams) GetDexRewardsStakers() LegacyDexRewardsTracker {
	if m != nil {
		return m.DexRewardsStakers
	}
	return LegacyDexRewardsTracker{}
}

// Params defines the parameters for the module.
type Params struct {
	StakeIncentives *IncentiveInfo `protobuf:"bytes,1,opt,name=stake_incentives,json=stakeIncentives,proto3" json:"stake_incentives,omitempty"`
	EdenCommitVal   string         `protobuf:"bytes,2,opt,name=eden_commit_val,json=edenCommitVal,proto3" json:"eden_commit_val,omitempty"`
	EdenbCommitVal  string         `protobuf:"bytes,3,opt,name=edenb_commit_val,json=edenbCommitVal,proto3" json:"edenb_commit_val,omitempty"`
	// Maximum eden reward apr for stakers - [0 - 0.3]
	MaxEdenRewardAprStakers github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=max_eden_reward_apr_stakers,json=maxEdenRewardAprStakers,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"max_eden_reward_apr_stakers"`
	EdenBoostApr            github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,5,opt,name=eden_boost_apr,json=edenBoostApr,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"eden_boost_apr"`
	// Tracking dex rewards given to stakers
	DexRewardsStakers DexRewardsTracker `protobuf:"bytes,6,opt,name=dex_rewards_stakers,json=dexRewardsStakers,proto3" json:"dex_rewards_stakers"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_66041162e1ecb63b, []int{1}
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

func (m *Params) GetStakeIncentives() *IncentiveInfo {
	if m != nil {
		return m.StakeIncentives
	}
	return nil
}

func (m *Params) GetEdenCommitVal() string {
	if m != nil {
		return m.EdenCommitVal
	}
	return ""
}

func (m *Params) GetEdenbCommitVal() string {
	if m != nil {
		return m.EdenbCommitVal
	}
	return ""
}

func (m *Params) GetDexRewardsStakers() DexRewardsTracker {
	if m != nil {
		return m.DexRewardsStakers
	}
	return DexRewardsTracker{}
}

func init() {
	proto.RegisterType((*LegacyParams)(nil), "elys.estaking.LegacyParams")
	proto.RegisterType((*Params)(nil), "elys.estaking.Params")
}

func init() { proto.RegisterFile("elys/estaking/params.proto", fileDescriptor_66041162e1ecb63b) }

var fileDescriptor_66041162e1ecb63b = []byte{
	// 459 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x53, 0xcf, 0x8b, 0xd3, 0x40,
	0x18, 0x4d, 0x6c, 0x37, 0xe2, 0xb8, 0xeb, 0xae, 0x51, 0x30, 0x54, 0x4d, 0x4b, 0x0f, 0xb5, 0x97,
	0x26, 0xa0, 0x37, 0x6f, 0x1b, 0x57, 0x96, 0x05, 0x05, 0x89, 0xcb, 0x1e, 0x44, 0x18, 0x26, 0xc9,
	0x67, 0x0c, 0xf9, 0x31, 0x61, 0x66, 0xdc, 0x4d, 0xff, 0x0b, 0x8f, 0x1e, 0xfd, 0x67, 0x84, 0x3d,
	0xee, 0x4d, 0xf1, 0x50, 0xa4, 0xfd, 0x47, 0x64, 0x26, 0xe9, 0x4f, 0xea, 0x41, 0x51, 0xf0, 0xd4,
	0xf6, 0x7b, 0xef, 0x7b, 0xef, 0xf5, 0x7b, 0x0c, 0xea, 0x40, 0x36, 0xe6, 0x2e, 0x70, 0x41, 0xd2,
	0xa4, 0x88, 0xdd, 0x92, 0x30, 0x92, 0x73, 0xa7, 0x64, 0x54, 0x50, 0x73, 0x4f, 0x62, 0xce, 0x1c,
	0xeb, 0xdc, 0x8d, 0x69, 0x4c, 0x15, 0xe2, 0xca, 0x6f, 0x35, 0xa9, 0xf3, 0x70, 0x5d, 0x20, 0x29,
	0x42, 0x28, 0x44, 0x72, 0x0e, 0x0d, 0xfc, 0x68, 0x1d, 0x8e, 0xa0, 0xc2, 0x0c, 0x2e, 0x08, 0x8b,
	0x38, 0x16, 0x8c, 0x84, 0x29, 0xb0, 0x9a, 0xd8, 0xff, 0xda, 0x42, 0xbb, 0x2f, 0x20, 0x26, 0xe1,
	0xf8, 0x95, 0xca, 0x60, 0xbe, 0x44, 0x07, 0x72, 0x0b, 0xf0, 0x42, 0x92, 0x5b, 0x7a, 0x4f, 0x1f,
	0xde, 0x7c, 0xdc, 0x77, 0xd6, 0x82, 0x39, 0xf5, 0xda, 0xc9, 0x9c, 0x76, 0x52, 0xbc, 0xa3, 0xfe,
	0xbe, 0xda, 0x5d, 0xcc, 0xb8, 0x39, 0x40, 0xfb, 0x10, 0x41, 0x81, 0x43, 0x9a, 0xe7, 0x89, 0xc0,
	0xe7, 0x24, 0xb3, 0xae, 0xf5, 0xf4, 0xe1, 0x0d, 0x7f, 0x4f, 0x8e, 0x9f, 0xa9, 0xe9, 0x19, 0xc9,
	0xcc, 0x21, 0x3a, 0x90, 0x83, 0x60, 0x95, 0xd8, 0x52, 0xc4, 0x5b, 0x6a, 0xbe, 0x64, 0x66, 0xe8,
	0x7e, 0x4e, 0x2a, 0xac, 0x54, 0xeb, 0xff, 0x84, 0x49, 0xc9, 0xb0, 0x32, 0x66, 0xdc, 0xda, 0x91,
	0x4b, 0x9e, 0x73, 0x39, 0xe9, 0x6a, 0xdf, 0x27, 0xdd, 0x41, 0x9c, 0x88, 0xf7, 0x1f, 0x02, 0x27,
	0xa4, 0xb9, 0x1b, 0x52, 0x9e, 0x53, 0xde, 0x7c, 0x8c, 0x78, 0x94, 0xba, 0x62, 0x5c, 0x02, 0x77,
	0x8e, 0x20, 0xf4, 0xef, 0xe5, 0xa4, 0x7a, 0x1e, 0x41, 0xe1, 0x2b, 0xc1, 0xc3, 0x92, 0xbd, 0xae,
	0xe5, 0xcc, 0x53, 0xa4, 0xfc, 0x71, 0x40, 0x29, 0x17, 0xd2, 0xc8, 0x32, 0xfe, 0xc8, 0x60, 0x57,
	0xaa, 0x78, 0x52, 0xe4, 0xb0, 0x64, 0xe6, 0x5b, 0x74, 0x67, 0xb5, 0x92, 0x79, 0xf6, 0xeb, 0xea,
	0xce, 0x83, 0xad, 0x77, 0x3e, 0x82, 0xaa, 0x0e, 0xc7, 0x4f, 0xeb, 0x02, 0xbd, 0xb6, 0x8c, 0xe0,
	0xdf, 0x8e, 0x16, 0x40, 0x93, 0xf9, 0x69, 0xfb, 0xd3, 0xe7, 0xae, 0xd6, 0xff, 0xd2, 0x42, 0x46,
	0xd3, 0xe9, 0xf1, 0x2f, 0x3b, 0x7d, 0xb0, 0xe1, 0xf5, 0xbf, 0xb5, 0xd9, 0xfe, 0xd7, 0x6d, 0xee,
	0xfc, 0x85, 0x36, 0xcf, 0xb6, 0xb7, 0x69, 0xa8, 0x0b, 0xf7, 0x36, 0x2e, 0xfc, 0xbb, 0x3d, 0x7a,
	0xc7, 0x97, 0x53, 0x5b, 0xbf, 0x9a, 0xda, 0xfa, 0x8f, 0xa9, 0xad, 0x7f, 0x9c, 0xd9, 0xda, 0xd5,
	0xcc, 0xd6, 0xbe, 0xcd, 0x6c, 0xed, 0xcd, 0x68, 0x25, 0xad, 0x34, 0x19, 0x15, 0x20, 0x2e, 0x28,
	0x4b, 0xd5, 0x0f, 0xb7, 0x5a, 0x3e, 0x7f, 0x15, 0x3c, 0x30, 0xd4, 0x8b, 0x7f, 0xf2, 0x33, 0x00,
	0x00, 0xff, 0xff, 0x0d, 0x59, 0x39, 0x8f, 0x7c, 0x04, 0x00, 0x00,
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
	{
		size, err := m.DexRewardsStakers.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	{
		size := m.EdenBoostApr.Size()
		i -= size
		if _, err := m.EdenBoostApr.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.MaxEdenRewardAprStakers.Size()
		i -= size
		if _, err := m.MaxEdenRewardAprStakers.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.EdenbCommitVal) > 0 {
		i -= len(m.EdenbCommitVal)
		copy(dAtA[i:], m.EdenbCommitVal)
		i = encodeVarintParams(dAtA, i, uint64(len(m.EdenbCommitVal)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.EdenCommitVal) > 0 {
		i -= len(m.EdenCommitVal)
		copy(dAtA[i:], m.EdenCommitVal)
		i = encodeVarintParams(dAtA, i, uint64(len(m.EdenCommitVal)))
		i--
		dAtA[i] = 0x12
	}
	if m.StakeIncentives != nil {
		{
			size, err := m.StakeIncentives.MarshalToSizedBuffer(dAtA[:i])
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
		size, err := m.DexRewardsStakers.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	{
		size := m.EdenBoostApr.Size()
		i -= size
		if _, err := m.EdenBoostApr.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	{
		size := m.MaxEdenRewardAprStakers.Size()
		i -= size
		if _, err := m.MaxEdenRewardAprStakers.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	if len(m.EdenbCommitVal) > 0 {
		i -= len(m.EdenbCommitVal)
		copy(dAtA[i:], m.EdenbCommitVal)
		i = encodeVarintParams(dAtA, i, uint64(len(m.EdenbCommitVal)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.EdenCommitVal) > 0 {
		i -= len(m.EdenCommitVal)
		copy(dAtA[i:], m.EdenCommitVal)
		i = encodeVarintParams(dAtA, i, uint64(len(m.EdenCommitVal)))
		i--
		dAtA[i] = 0x12
	}
	if m.StakeIncentives != nil {
		{
			size, err := m.StakeIncentives.MarshalToSizedBuffer(dAtA[:i])
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
func (m *LegacyParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.StakeIncentives != nil {
		l = m.StakeIncentives.Size()
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.EdenCommitVal)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.EdenbCommitVal)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.MaxEdenRewardAprStakers.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.EdenBoostApr.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.DexRewardsStakers.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.StakeIncentives != nil {
		l = m.StakeIncentives.Size()
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.EdenCommitVal)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.EdenbCommitVal)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = m.MaxEdenRewardAprStakers.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.EdenBoostApr.Size()
	n += 1 + l + sovParams(uint64(l))
	l = m.DexRewardsStakers.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
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
				return fmt.Errorf("proto: wrong wireType = %d for field StakeIncentives", wireType)
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
			if m.StakeIncentives == nil {
				m.StakeIncentives = &LegacyIncentiveInfo{}
			}
			if err := m.StakeIncentives.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EdenCommitVal", wireType)
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
			m.EdenCommitVal = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EdenbCommitVal", wireType)
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
			m.EdenbCommitVal = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxEdenRewardAprStakers", wireType)
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
			if err := m.MaxEdenRewardAprStakers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EdenBoostApr", wireType)
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
			if err := m.EdenBoostApr.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DexRewardsStakers", wireType)
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
			if err := m.DexRewardsStakers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
				return fmt.Errorf("proto: wrong wireType = %d for field StakeIncentives", wireType)
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
			if m.StakeIncentives == nil {
				m.StakeIncentives = &IncentiveInfo{}
			}
			if err := m.StakeIncentives.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EdenCommitVal", wireType)
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
			m.EdenCommitVal = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EdenbCommitVal", wireType)
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
			m.EdenbCommitVal = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxEdenRewardAprStakers", wireType)
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
			if err := m.MaxEdenRewardAprStakers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field EdenBoostApr", wireType)
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
			if err := m.EdenBoostApr.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DexRewardsStakers", wireType)
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
			if err := m.DexRewardsStakers.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
