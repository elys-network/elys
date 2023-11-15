// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/incentive/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
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

// MsgWithdrawDelegatorReward represents delegation withdrawal to a delegator
// from a single validator.
type MsgWithdrawRewards struct {
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	Denom            string `protobuf:"bytes,2,opt,name=denom,proto3" json:"denom,omitempty"`
	WithdrawType     int64  `protobuf:"varint,3,opt,name=withdraw_type,json=withdrawType,proto3" json:"withdraw_type,omitempty"`
}

func (m *MsgWithdrawRewards) Reset()         { *m = MsgWithdrawRewards{} }
func (m *MsgWithdrawRewards) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawRewards) ProtoMessage()    {}
func (*MsgWithdrawRewards) Descriptor() ([]byte, []int) {
	return fileDescriptor_59dc3bedfb1cce84, []int{0}
}
func (m *MsgWithdrawRewards) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawRewards) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawRewards.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawRewards) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawRewards.Merge(m, src)
}
func (m *MsgWithdrawRewards) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawRewards) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawRewards.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawRewards proto.InternalMessageInfo

// MsgWithdrawDelegatorRewardResponse defines the Msg/WithdrawDelegatorReward response type.
type MsgWithdrawRewardsResponse struct {
}

func (m *MsgWithdrawRewardsResponse) Reset()         { *m = MsgWithdrawRewardsResponse{} }
func (m *MsgWithdrawRewardsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawRewardsResponse) ProtoMessage()    {}
func (*MsgWithdrawRewardsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_59dc3bedfb1cce84, []int{1}
}
func (m *MsgWithdrawRewardsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawRewardsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawRewardsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawRewardsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawRewardsResponse.Merge(m, src)
}
func (m *MsgWithdrawRewardsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawRewardsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawRewardsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawRewardsResponse proto.InternalMessageInfo

// MsgWithdrawValidatorCommission withdraws the full commission to the validator
// address.
type MsgWithdrawValidatorCommission struct {
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
	Denom            string `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
}

func (m *MsgWithdrawValidatorCommission) Reset()         { *m = MsgWithdrawValidatorCommission{} }
func (m *MsgWithdrawValidatorCommission) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawValidatorCommission) ProtoMessage()    {}
func (*MsgWithdrawValidatorCommission) Descriptor() ([]byte, []int) {
	return fileDescriptor_59dc3bedfb1cce84, []int{2}
}
func (m *MsgWithdrawValidatorCommission) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawValidatorCommission) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawValidatorCommission.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawValidatorCommission) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawValidatorCommission.Merge(m, src)
}
func (m *MsgWithdrawValidatorCommission) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawValidatorCommission) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawValidatorCommission.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawValidatorCommission proto.InternalMessageInfo

// MsgWithdrawValidatorCommissionResponse defines the Msg/WithdrawValidatorCommission response type.
type MsgWithdrawValidatorCommissionResponse struct {
}

func (m *MsgWithdrawValidatorCommissionResponse) Reset() {
	*m = MsgWithdrawValidatorCommissionResponse{}
}
func (m *MsgWithdrawValidatorCommissionResponse) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawValidatorCommissionResponse) ProtoMessage()    {}
func (*MsgWithdrawValidatorCommissionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_59dc3bedfb1cce84, []int{3}
}
func (m *MsgWithdrawValidatorCommissionResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawValidatorCommissionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawValidatorCommissionResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawValidatorCommissionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawValidatorCommissionResponse.Merge(m, src)
}
func (m *MsgWithdrawValidatorCommissionResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawValidatorCommissionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawValidatorCommissionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawValidatorCommissionResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgWithdrawRewards)(nil), "elys.incentive.MsgWithdrawRewards")
	proto.RegisterType((*MsgWithdrawRewardsResponse)(nil), "elys.incentive.MsgWithdrawRewardsResponse")
	proto.RegisterType((*MsgWithdrawValidatorCommission)(nil), "elys.incentive.MsgWithdrawValidatorCommission")
	proto.RegisterType((*MsgWithdrawValidatorCommissionResponse)(nil), "elys.incentive.MsgWithdrawValidatorCommissionResponse")
}

func init() { proto.RegisterFile("elys/incentive/tx.proto", fileDescriptor_59dc3bedfb1cce84) }

var fileDescriptor_59dc3bedfb1cce84 = []byte{
	// 395 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4f, 0xcd, 0xa9, 0x2c,
	0xd6, 0xcf, 0xcc, 0x4b, 0x4e, 0xcd, 0x2b, 0xc9, 0x2c, 0x4b, 0xd5, 0x2f, 0xa9, 0xd0, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x03, 0x49, 0xe8, 0xc1, 0x25, 0xa4, 0x44, 0xd2, 0xf3, 0xd3, 0xf3,
	0xc1, 0x52, 0xfa, 0x20, 0x16, 0x44, 0x95, 0x94, 0x64, 0x72, 0x7e, 0x71, 0x6e, 0x7e, 0x71, 0x3c,
	0x44, 0x02, 0xc2, 0x81, 0x48, 0x29, 0x2d, 0x60, 0xe4, 0x12, 0xf2, 0x2d, 0x4e, 0x0f, 0xcf, 0x2c,
	0xc9, 0x48, 0x29, 0x4a, 0x2c, 0x0f, 0x4a, 0x2d, 0x4f, 0x2c, 0x4a, 0x29, 0x16, 0x72, 0xe5, 0x12,
	0x4c, 0x49, 0xcd, 0x49, 0x4d, 0x4f, 0x2c, 0xc9, 0x2f, 0x8a, 0x4f, 0x4c, 0x49, 0x29, 0x4a, 0x2d,
	0x2e, 0x96, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x74, 0x92, 0xb8, 0xb4, 0x45, 0x57, 0x04, 0x6a, 0x86,
	0x23, 0x44, 0x26, 0xb8, 0xa4, 0x28, 0x33, 0x2f, 0x3d, 0x48, 0x00, 0xae, 0x05, 0x2a, 0x2e, 0x24,
	0xc2, 0xc5, 0x9a, 0x92, 0x9a, 0x97, 0x9f, 0x2b, 0xc1, 0x04, 0xd2, 0x1a, 0x04, 0xe1, 0x08, 0x29,
	0x73, 0xf1, 0x96, 0x43, 0xed, 0x8b, 0x2f, 0xa9, 0x2c, 0x48, 0x95, 0x60, 0x56, 0x60, 0xd4, 0x60,
	0x0e, 0xe2, 0x81, 0x09, 0x86, 0x54, 0x16, 0xa4, 0x5a, 0x71, 0x74, 0x2c, 0x90, 0x67, 0x78, 0xb1,
	0x40, 0x9e, 0x41, 0x49, 0x86, 0x4b, 0x0a, 0xd3, 0x85, 0x41, 0xa9, 0xc5, 0x05, 0xf9, 0x79, 0xc5,
	0xa9, 0x4a, 0xe7, 0x18, 0xb9, 0xe4, 0x90, 0xa4, 0xc3, 0x12, 0x73, 0x32, 0x53, 0x40, 0x4e, 0x70,
	0xce, 0xcf, 0xcd, 0xcd, 0x2c, 0x2e, 0xce, 0xcc, 0xcf, 0xa3, 0x96, 0x67, 0x5c, 0xb9, 0x04, 0xcb,
	0x60, 0xa6, 0xc3, 0x8d, 0x61, 0x22, 0x64, 0x0c, 0x5c, 0x0b, 0x46, 0x98, 0x30, 0x23, 0x85, 0x09,
	0x92, 0x77, 0x35, 0xb8, 0xd4, 0xf0, 0xfb, 0x07, 0xe6, 0x75, 0xa3, 0x0f, 0x8c, 0x5c, 0xcc, 0xbe,
	0xc5, 0xe9, 0x42, 0x89, 0x5c, 0xfc, 0xe8, 0xf1, 0xa7, 0xa4, 0x87, 0x9a, 0x30, 0xf4, 0x30, 0x43,
	0x50, 0x4a, 0x8b, 0xb0, 0x1a, 0x98, 0x55, 0x42, 0xad, 0x8c, 0x5c, 0xd2, 0xf8, 0x82, 0x58, 0x0f,
	0x8f, 0x59, 0x58, 0xd4, 0x4b, 0x99, 0x91, 0xa6, 0x1e, 0xe6, 0x0e, 0x27, 0x8f, 0x13, 0x8f, 0xe4,
	0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f,
	0xe5, 0x18, 0x6e, 0x3c, 0x96, 0x63, 0x88, 0xd2, 0x4b, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b,
	0xce, 0xcf, 0xd5, 0x07, 0x99, 0xad, 0x9b, 0x97, 0x5a, 0x52, 0x9e, 0x5f, 0x94, 0x0d, 0xe6, 0xe8,
	0x57, 0x20, 0x67, 0x9e, 0xca, 0x82, 0xd4, 0xe2, 0x24, 0x36, 0x70, 0xfa, 0x37, 0x06, 0x04, 0x00,
	0x00, 0xff, 0xff, 0x68, 0xce, 0x9f, 0x2a, 0x5b, 0x03, 0x00, 0x00,
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
	// WithdrawDelegatorReward defines a method to withdraw rewards of delegator
	// from a single validator.
	WithdrawRewards(ctx context.Context, in *MsgWithdrawRewards, opts ...grpc.CallOption) (*MsgWithdrawRewardsResponse, error)
	// WithdrawValidatorCommission defines a method to withdraw the
	// full commission to the validator address.
	WithdrawValidatorCommission(ctx context.Context, in *MsgWithdrawValidatorCommission, opts ...grpc.CallOption) (*MsgWithdrawValidatorCommissionResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) WithdrawRewards(ctx context.Context, in *MsgWithdrawRewards, opts ...grpc.CallOption) (*MsgWithdrawRewardsResponse, error) {
	out := new(MsgWithdrawRewardsResponse)
	err := c.cc.Invoke(ctx, "/elys.incentive.Msg/WithdrawRewards", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) WithdrawValidatorCommission(ctx context.Context, in *MsgWithdrawValidatorCommission, opts ...grpc.CallOption) (*MsgWithdrawValidatorCommissionResponse, error) {
	out := new(MsgWithdrawValidatorCommissionResponse)
	err := c.cc.Invoke(ctx, "/elys.incentive.Msg/WithdrawValidatorCommission", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// WithdrawDelegatorReward defines a method to withdraw rewards of delegator
	// from a single validator.
	WithdrawRewards(context.Context, *MsgWithdrawRewards) (*MsgWithdrawRewardsResponse, error)
	// WithdrawValidatorCommission defines a method to withdraw the
	// full commission to the validator address.
	WithdrawValidatorCommission(context.Context, *MsgWithdrawValidatorCommission) (*MsgWithdrawValidatorCommissionResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) WithdrawRewards(ctx context.Context, req *MsgWithdrawRewards) (*MsgWithdrawRewardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawRewards not implemented")
}
func (*UnimplementedMsgServer) WithdrawValidatorCommission(ctx context.Context, req *MsgWithdrawValidatorCommission) (*MsgWithdrawValidatorCommissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawValidatorCommission not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_WithdrawRewards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawRewards)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WithdrawRewards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.incentive.Msg/WithdrawRewards",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawRewards(ctx, req.(*MsgWithdrawRewards))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_WithdrawValidatorCommission_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawValidatorCommission)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WithdrawValidatorCommission(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.incentive.Msg/WithdrawValidatorCommission",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawValidatorCommission(ctx, req.(*MsgWithdrawValidatorCommission))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elys.incentive.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WithdrawRewards",
			Handler:    _Msg_WithdrawRewards_Handler,
		},
		{
			MethodName: "WithdrawValidatorCommission",
			Handler:    _Msg_WithdrawValidatorCommission_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/incentive/tx.proto",
}

func (m *MsgWithdrawRewards) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawRewards) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawRewards) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.WithdrawType != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.WithdrawType))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DelegatorAddress) > 0 {
		i -= len(m.DelegatorAddress)
		copy(dAtA[i:], m.DelegatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.DelegatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgWithdrawRewardsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawRewardsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawRewardsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgWithdrawValidatorCommission) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawValidatorCommission) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawValidatorCommission) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ValidatorAddress)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.DelegatorAddress) > 0 {
		i -= len(m.DelegatorAddress)
		copy(dAtA[i:], m.DelegatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.DelegatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgWithdrawValidatorCommissionResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawValidatorCommissionResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawValidatorCommissionResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *MsgWithdrawRewards) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DelegatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.WithdrawType != 0 {
		n += 1 + sovTx(uint64(m.WithdrawType))
	}
	return n
}

func (m *MsgWithdrawRewardsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgWithdrawValidatorCommission) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DelegatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.ValidatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgWithdrawValidatorCommissionResponse) Size() (n int) {
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
func (m *MsgWithdrawRewards) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgWithdrawRewards: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawRewards: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegatorAddress", wireType)
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
			m.DelegatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
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
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field WithdrawType", wireType)
			}
			m.WithdrawType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.WithdrawType |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
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
func (m *MsgWithdrawRewardsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgWithdrawRewardsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawRewardsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgWithdrawValidatorCommission) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgWithdrawValidatorCommission: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawValidatorCommission: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegatorAddress", wireType)
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
			m.DelegatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorAddress", wireType)
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
			m.ValidatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
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
			m.Denom = string(dAtA[iNdEx:postIndex])
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
func (m *MsgWithdrawValidatorCommissionResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgWithdrawValidatorCommissionResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawValidatorCommissionResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
