// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/masterchef/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type MsgAddExternalIncentive struct {
	Sender         string                                 `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	RewardDenom    string                                 `protobuf:"bytes,2,opt,name=reward_denom,json=rewardDenom,proto3" json:"reward_denom,omitempty"`
	PoolId         uint64                                 `protobuf:"varint,3,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty"`
	FromBlock      uint64                                 `protobuf:"varint,4,opt,name=from_block,json=fromBlock,proto3" json:"from_block,omitempty"`
	ToBlock        uint64                                 `protobuf:"varint,5,opt,name=to_block,json=toBlock,proto3" json:"to_block,omitempty"`
	AmountPerBlock github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,6,opt,name=amount_per_block,json=amountPerBlock,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount_per_block"`
}

func (m *MsgAddExternalIncentive) Reset()         { *m = MsgAddExternalIncentive{} }
func (m *MsgAddExternalIncentive) String() string { return proto.CompactTextString(m) }
func (*MsgAddExternalIncentive) ProtoMessage()    {}
func (*MsgAddExternalIncentive) Descriptor() ([]byte, []int) {
	return fileDescriptor_2574ed545e5b2c11, []int{0}
}
func (m *MsgAddExternalIncentive) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAddExternalIncentive) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddExternalIncentive.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAddExternalIncentive) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAddExternalIncentive.Merge(m, src)
}
func (m *MsgAddExternalIncentive) XXX_Size() int {
	return m.Size()
}
func (m *MsgAddExternalIncentive) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAddExternalIncentive.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAddExternalIncentive proto.InternalMessageInfo

func (m *MsgAddExternalIncentive) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *MsgAddExternalIncentive) GetRewardDenom() string {
	if m != nil {
		return m.RewardDenom
	}
	return ""
}

func (m *MsgAddExternalIncentive) GetPoolId() uint64 {
	if m != nil {
		return m.PoolId
	}
	return 0
}

func (m *MsgAddExternalIncentive) GetFromBlock() uint64 {
	if m != nil {
		return m.FromBlock
	}
	return 0
}

func (m *MsgAddExternalIncentive) GetToBlock() uint64 {
	if m != nil {
		return m.ToBlock
	}
	return 0
}

type MsgAddExternalIncentiveResponse struct {
}

func (m *MsgAddExternalIncentiveResponse) Reset()         { *m = MsgAddExternalIncentiveResponse{} }
func (m *MsgAddExternalIncentiveResponse) String() string { return proto.CompactTextString(m) }
func (*MsgAddExternalIncentiveResponse) ProtoMessage()    {}
func (*MsgAddExternalIncentiveResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2574ed545e5b2c11, []int{1}
}
func (m *MsgAddExternalIncentiveResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgAddExternalIncentiveResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgAddExternalIncentiveResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgAddExternalIncentiveResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgAddExternalIncentiveResponse.Merge(m, src)
}
func (m *MsgAddExternalIncentiveResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgAddExternalIncentiveResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgAddExternalIncentiveResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgAddExternalIncentiveResponse proto.InternalMessageInfo

type MsgClaimRewards struct {
	Sender  string   `protobuf:"bytes,1,opt,name=sender,proto3" json:"sender,omitempty"`
	PoolIds []uint64 `protobuf:"varint,2,rep,packed,name=pool_ids,json=poolIds,proto3" json:"pool_ids,omitempty"`
}

func (m *MsgClaimRewards) Reset()         { *m = MsgClaimRewards{} }
func (m *MsgClaimRewards) String() string { return proto.CompactTextString(m) }
func (*MsgClaimRewards) ProtoMessage()    {}
func (*MsgClaimRewards) Descriptor() ([]byte, []int) {
	return fileDescriptor_2574ed545e5b2c11, []int{2}
}
func (m *MsgClaimRewards) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgClaimRewards) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgClaimRewards.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgClaimRewards) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgClaimRewards.Merge(m, src)
}
func (m *MsgClaimRewards) XXX_Size() int {
	return m.Size()
}
func (m *MsgClaimRewards) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgClaimRewards.DiscardUnknown(m)
}

var xxx_messageInfo_MsgClaimRewards proto.InternalMessageInfo

func (m *MsgClaimRewards) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *MsgClaimRewards) GetPoolIds() []uint64 {
	if m != nil {
		return m.PoolIds
	}
	return nil
}

type MsgClaimRewardsResponse struct {
}

func (m *MsgClaimRewardsResponse) Reset()         { *m = MsgClaimRewardsResponse{} }
func (m *MsgClaimRewardsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgClaimRewardsResponse) ProtoMessage()    {}
func (*MsgClaimRewardsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2574ed545e5b2c11, []int{3}
}
func (m *MsgClaimRewardsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgClaimRewardsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgClaimRewardsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgClaimRewardsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgClaimRewardsResponse.Merge(m, src)
}
func (m *MsgClaimRewardsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgClaimRewardsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgClaimRewardsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgClaimRewardsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgAddExternalIncentive)(nil), "elys.masterchef.MsgAddExternalIncentive")
	proto.RegisterType((*MsgAddExternalIncentiveResponse)(nil), "elys.masterchef.MsgAddExternalIncentiveResponse")
	proto.RegisterType((*MsgClaimRewards)(nil), "elys.masterchef.MsgClaimRewards")
	proto.RegisterType((*MsgClaimRewardsResponse)(nil), "elys.masterchef.MsgClaimRewardsResponse")
}

func init() { proto.RegisterFile("elys/masterchef/tx.proto", fileDescriptor_2574ed545e5b2c11) }

var fileDescriptor_2574ed545e5b2c11 = []byte{
	// 422 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x8d, 0x93, 0xe2, 0xb6, 0x4b, 0x45, 0xd1, 0xaa, 0xa2, 0x4e, 0x24, 0x1c, 0xd7, 0x07, 0xe4,
	0x4b, 0x6d, 0x04, 0x5f, 0x40, 0x28, 0x07, 0x1f, 0x22, 0x21, 0x9f, 0x50, 0x2f, 0x96, 0xe3, 0x9d,
	0xba, 0x56, 0xec, 0x1d, 0x6b, 0x77, 0x4b, 0xd3, 0x23, 0x7f, 0xc0, 0x67, 0xf5, 0xd8, 0x0b, 0x12,
	0xe2, 0x50, 0xa1, 0xe4, 0x47, 0xd0, 0xae, 0x1d, 0x91, 0xa2, 0x46, 0xf4, 0x64, 0xcf, 0xbc, 0xa7,
	0x37, 0xfb, 0xe6, 0x0d, 0x71, 0xa0, 0xba, 0x91, 0x51, 0x9d, 0x49, 0x05, 0x22, 0xbf, 0x84, 0x8b,
	0x48, 0x2d, 0xc2, 0x46, 0xa0, 0x42, 0x7a, 0xa8, 0x91, 0xf0, 0x2f, 0x32, 0x3a, 0x2a, 0xb0, 0x40,
	0x83, 0x45, 0xfa, 0xaf, 0xa5, 0xf9, 0xdf, 0xfa, 0xe4, 0x78, 0x2a, 0x8b, 0x0f, 0x8c, 0x7d, 0x5a,
	0x28, 0x10, 0x3c, 0xab, 0x62, 0x9e, 0x03, 0x57, 0xe5, 0x57, 0xa0, 0xaf, 0x88, 0x2d, 0x81, 0x33,
	0x10, 0x8e, 0xe5, 0x59, 0xc1, 0x7e, 0xd2, 0x55, 0xf4, 0x84, 0x1c, 0x08, 0xb8, 0xce, 0x04, 0x4b,
	0x19, 0x70, 0xac, 0x9d, 0xbe, 0x41, 0x9f, 0xb7, 0xbd, 0x33, 0xdd, 0xa2, 0xc7, 0x64, 0xb7, 0x41,
	0xac, 0xd2, 0x92, 0x39, 0x03, 0xcf, 0x0a, 0x76, 0x12, 0x5b, 0x97, 0x31, 0xa3, 0xaf, 0x09, 0xb9,
	0x10, 0x58, 0xa7, 0xb3, 0x0a, 0xf3, 0xb9, 0xb3, 0x63, 0xb0, 0x7d, 0xdd, 0x99, 0xe8, 0x06, 0x1d,
	0x92, 0x3d, 0x85, 0x1d, 0xf8, 0xcc, 0x80, 0xbb, 0x0a, 0x5b, 0xe8, 0x0b, 0x79, 0x99, 0xd5, 0x78,
	0xc5, 0x55, 0xda, 0x80, 0xe8, 0x28, 0xb6, 0x9e, 0x3c, 0x09, 0x6f, 0xef, 0xc7, 0xbd, 0x5f, 0xf7,
	0xe3, 0x37, 0x45, 0xa9, 0x2e, 0xaf, 0x66, 0x61, 0x8e, 0x75, 0x94, 0xa3, 0xac, 0x51, 0x76, 0x9f,
	0x53, 0xc9, 0xe6, 0x91, 0xba, 0x69, 0x40, 0x86, 0x31, 0x57, 0xc9, 0x8b, 0x56, 0xe7, 0x33, 0x08,
	0xa3, 0xec, 0x9f, 0x90, 0xf1, 0x96, 0x15, 0x24, 0x20, 0x1b, 0xe4, 0x12, 0xfc, 0x33, 0x72, 0x38,
	0x95, 0xc5, 0xc7, 0x2a, 0x2b, 0xeb, 0xc4, 0xd8, 0x94, 0x5b, 0xb7, 0x33, 0x24, 0x7b, 0x9d, 0x75,
	0xe9, 0xf4, 0xbd, 0x81, 0xb6, 0xd0, 0x7a, 0x97, 0xfe, 0xd0, 0xec, 0x7a, 0x53, 0x65, 0x3d, 0xe0,
	0xdd, 0x0f, 0x8b, 0x0c, 0xa6, 0xb2, 0xa0, 0x82, 0x1c, 0x3d, 0x9a, 0x45, 0x10, 0xfe, 0x93, 0x67,
	0xb8, 0xe5, 0xc9, 0xa3, 0xb7, 0x4f, 0x65, 0xae, 0x67, 0xd3, 0x73, 0x72, 0xf0, 0xc0, 0x99, 0xf7,
	0x98, 0xc2, 0x26, 0x63, 0x14, 0xfc, 0x8f, 0xb1, 0xd6, 0x9e, 0xc4, 0xb7, 0x4b, 0xd7, 0xba, 0x5b,
	0xba, 0xd6, 0xef, 0xa5, 0x6b, 0x7d, 0x5f, 0xb9, 0xbd, 0xbb, 0x95, 0xdb, 0xfb, 0xb9, 0x72, 0x7b,
	0xe7, 0xd1, 0x46, 0x5a, 0x5a, 0xed, 0x94, 0x83, 0xba, 0x46, 0x31, 0x37, 0x45, 0xb4, 0x78, 0x70,
	0xd4, 0x3a, 0xba, 0x99, 0x6d, 0x2e, 0xf6, 0xfd, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc1, 0x8c,
	0x3e, 0x1a, 0xf4, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	AddExternalIncentive(ctx context.Context, in *MsgAddExternalIncentive, opts ...grpc.CallOption) (*MsgAddExternalIncentiveResponse, error)
	ClaimRewards(ctx context.Context, in *MsgClaimRewards, opts ...grpc.CallOption) (*MsgClaimRewardsResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) AddExternalIncentive(ctx context.Context, in *MsgAddExternalIncentive, opts ...grpc.CallOption) (*MsgAddExternalIncentiveResponse, error) {
	out := new(MsgAddExternalIncentiveResponse)
	err := c.cc.Invoke(ctx, "/elys.masterchef.Msg/AddExternalIncentive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ClaimRewards(ctx context.Context, in *MsgClaimRewards, opts ...grpc.CallOption) (*MsgClaimRewardsResponse, error) {
	out := new(MsgClaimRewardsResponse)
	err := c.cc.Invoke(ctx, "/elys.masterchef.Msg/ClaimRewards", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	AddExternalIncentive(context.Context, *MsgAddExternalIncentive) (*MsgAddExternalIncentiveResponse, error)
	ClaimRewards(context.Context, *MsgClaimRewards) (*MsgClaimRewardsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) AddExternalIncentive(ctx context.Context, req *MsgAddExternalIncentive) (*MsgAddExternalIncentiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddExternalIncentive not implemented")
}
func (*UnimplementedMsgServer) ClaimRewards(ctx context.Context, req *MsgClaimRewards) (*MsgClaimRewardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClaimRewards not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_AddExternalIncentive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddExternalIncentive)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddExternalIncentive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.masterchef.Msg/AddExternalIncentive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddExternalIncentive(ctx, req.(*MsgAddExternalIncentive))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClaimRewards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimRewards)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimRewards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.masterchef.Msg/ClaimRewards",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimRewards(ctx, req.(*MsgClaimRewards))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elys.masterchef.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddExternalIncentive",
			Handler:    _Msg_AddExternalIncentive_Handler,
		},
		{
			MethodName: "ClaimRewards",
			Handler:    _Msg_ClaimRewards_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/masterchef/tx.proto",
}

func (m *MsgAddExternalIncentive) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddExternalIncentive) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddExternalIncentive) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.AmountPerBlock.Size()
		i -= size
		if _, err := m.AmountPerBlock.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x32
	if m.ToBlock != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.ToBlock))
		i--
		dAtA[i] = 0x28
	}
	if m.FromBlock != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.FromBlock))
		i--
		dAtA[i] = 0x20
	}
	if m.PoolId != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.PoolId))
		i--
		dAtA[i] = 0x18
	}
	if len(m.RewardDenom) > 0 {
		i -= len(m.RewardDenom)
		copy(dAtA[i:], m.RewardDenom)
		i = encodeVarintTx(dAtA, i, uint64(len(m.RewardDenom)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgAddExternalIncentiveResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgAddExternalIncentiveResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgAddExternalIncentiveResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgClaimRewards) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgClaimRewards) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgClaimRewards) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.PoolIds) > 0 {
		dAtA2 := make([]byte, len(m.PoolIds)*10)
		var j1 int
		for _, num := range m.PoolIds {
			for num >= 1<<7 {
				dAtA2[j1] = uint8(uint64(num)&0x7f | 0x80)
				num >>= 7
				j1++
			}
			dAtA2[j1] = uint8(num)
			j1++
		}
		i -= j1
		copy(dAtA[i:], dAtA2[:j1])
		i = encodeVarintTx(dAtA, i, uint64(j1))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Sender) > 0 {
		i -= len(m.Sender)
		copy(dAtA[i:], m.Sender)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Sender)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgClaimRewardsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgClaimRewardsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgClaimRewardsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgAddExternalIncentive) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.RewardDenom)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.PoolId != 0 {
		n += 1 + sovTx(uint64(m.PoolId))
	}
	if m.FromBlock != 0 {
		n += 1 + sovTx(uint64(m.FromBlock))
	}
	if m.ToBlock != 0 {
		n += 1 + sovTx(uint64(m.ToBlock))
	}
	l = m.AmountPerBlock.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgAddExternalIncentiveResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgClaimRewards) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Sender)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if len(m.PoolIds) > 0 {
		l = 0
		for _, e := range m.PoolIds {
			l += sovTx(uint64(e))
		}
		n += 1 + sovTx(uint64(l)) + l
	}
	return n
}

func (m *MsgClaimRewardsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgAddExternalIncentive) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgAddExternalIncentive: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddExternalIncentive: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RewardDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolId", wireType)
			}
			m.PoolId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field FromBlock", wireType)
			}
			m.FromBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.FromBlock |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ToBlock", wireType)
			}
			m.ToBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ToBlock |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AmountPerBlock", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.AmountPerBlock.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgAddExternalIncentiveResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgAddExternalIncentiveResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgAddExternalIncentiveResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgClaimRewards) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgClaimRewards: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgClaimRewards: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sender = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType == 0 {
				var v uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTx
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= uint64(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.PoolIds = append(m.PoolIds, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTx
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= int(b&0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthTx
				}
				postIndex := iNdEx + packedLen
				if postIndex < 0 {
					return ErrInvalidLengthTx
				}
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA[iNdEx:postIndex] {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.PoolIds) == 0 {
					m.PoolIds = make([]uint64, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowTx
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.PoolIds = append(m.PoolIds, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field PoolIds", wireType)
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgClaimRewardsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgClaimRewardsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgClaimRewardsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)
