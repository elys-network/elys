// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/incentive/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
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

type MsgSetWithdrawAddress struct {
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	WithdrawAddress  string `protobuf:"bytes,2,opt,name=withdraw_address,json=withdrawAddress,proto3" json:"withdraw_address,omitempty"`
}

func (m *MsgSetWithdrawAddress) Reset()         { *m = MsgSetWithdrawAddress{} }
func (m *MsgSetWithdrawAddress) String() string { return proto.CompactTextString(m) }
func (*MsgSetWithdrawAddress) ProtoMessage()    {}
func (*MsgSetWithdrawAddress) Descriptor() ([]byte, []int) {
	return fileDescriptor_59dc3bedfb1cce84, []int{0}
}
func (m *MsgSetWithdrawAddress) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetWithdrawAddress) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetWithdrawAddress.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetWithdrawAddress) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetWithdrawAddress.Merge(m, src)
}
func (m *MsgSetWithdrawAddress) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetWithdrawAddress) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetWithdrawAddress.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetWithdrawAddress proto.InternalMessageInfo

func (m *MsgSetWithdrawAddress) GetDelegatorAddress() string {
	if m != nil {
		return m.DelegatorAddress
	}
	return ""
}

func (m *MsgSetWithdrawAddress) GetWithdrawAddress() string {
	if m != nil {
		return m.WithdrawAddress
	}
	return ""
}

type MsgSetWithdrawAddressResponse struct {
}

func (m *MsgSetWithdrawAddressResponse) Reset()         { *m = MsgSetWithdrawAddressResponse{} }
func (m *MsgSetWithdrawAddressResponse) String() string { return proto.CompactTextString(m) }
func (*MsgSetWithdrawAddressResponse) ProtoMessage()    {}
func (*MsgSetWithdrawAddressResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_59dc3bedfb1cce84, []int{1}
}
func (m *MsgSetWithdrawAddressResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgSetWithdrawAddressResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgSetWithdrawAddressResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgSetWithdrawAddressResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgSetWithdrawAddressResponse.Merge(m, src)
}
func (m *MsgSetWithdrawAddressResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgSetWithdrawAddressResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgSetWithdrawAddressResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgSetWithdrawAddressResponse proto.InternalMessageInfo

type MsgWithdrawValidatorCommission struct {
	ValidatorAddress string `protobuf:"bytes,1,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
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

func (m *MsgWithdrawValidatorCommission) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

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

type MsgWithdrawDelegatorReward struct {
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
	ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
}

func (m *MsgWithdrawDelegatorReward) Reset()         { *m = MsgWithdrawDelegatorReward{} }
func (m *MsgWithdrawDelegatorReward) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawDelegatorReward) ProtoMessage()    {}
func (*MsgWithdrawDelegatorReward) Descriptor() ([]byte, []int) {
	return fileDescriptor_59dc3bedfb1cce84, []int{4}
}
func (m *MsgWithdrawDelegatorReward) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawDelegatorReward) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawDelegatorReward.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawDelegatorReward) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawDelegatorReward.Merge(m, src)
}
func (m *MsgWithdrawDelegatorReward) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawDelegatorReward) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawDelegatorReward.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawDelegatorReward proto.InternalMessageInfo

func (m *MsgWithdrawDelegatorReward) GetDelegatorAddress() string {
	if m != nil {
		return m.DelegatorAddress
	}
	return ""
}

func (m *MsgWithdrawDelegatorReward) GetValidatorAddress() string {
	if m != nil {
		return m.ValidatorAddress
	}
	return ""
}

type MsgWithdrawDelegatorRewardResponse struct {
	// Since: cosmos-sdk 0.46
	Amount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
}

func (m *MsgWithdrawDelegatorRewardResponse) Reset()         { *m = MsgWithdrawDelegatorRewardResponse{} }
func (m *MsgWithdrawDelegatorRewardResponse) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawDelegatorRewardResponse) ProtoMessage()    {}
func (*MsgWithdrawDelegatorRewardResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_59dc3bedfb1cce84, []int{5}
}
func (m *MsgWithdrawDelegatorRewardResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawDelegatorRewardResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawDelegatorRewardResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawDelegatorRewardResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawDelegatorRewardResponse.Merge(m, src)
}
func (m *MsgWithdrawDelegatorRewardResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawDelegatorRewardResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawDelegatorRewardResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawDelegatorRewardResponse proto.InternalMessageInfo

func (m *MsgWithdrawDelegatorRewardResponse) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

func init() {
	proto.RegisterType((*MsgSetWithdrawAddress)(nil), "elysnetwork.elys.incentive.MsgSetWithdrawAddress")
	proto.RegisterType((*MsgSetWithdrawAddressResponse)(nil), "elysnetwork.elys.incentive.MsgSetWithdrawAddressResponse")
	proto.RegisterType((*MsgWithdrawValidatorCommission)(nil), "elysnetwork.elys.incentive.MsgWithdrawValidatorCommission")
	proto.RegisterType((*MsgWithdrawValidatorCommissionResponse)(nil), "elysnetwork.elys.incentive.MsgWithdrawValidatorCommissionResponse")
	proto.RegisterType((*MsgWithdrawDelegatorReward)(nil), "elysnetwork.elys.incentive.MsgWithdrawDelegatorReward")
	proto.RegisterType((*MsgWithdrawDelegatorRewardResponse)(nil), "elysnetwork.elys.incentive.MsgWithdrawDelegatorRewardResponse")
}

func init() { proto.RegisterFile("elys/incentive/tx.proto", fileDescriptor_59dc3bedfb1cce84) }

var fileDescriptor_59dc3bedfb1cce84 = []byte{
	// 456 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0x4f, 0x6b, 0x13, 0x41,
	0x14, 0xcf, 0x34, 0x50, 0x70, 0x3c, 0x58, 0x17, 0xa5, 0x75, 0xc5, 0x49, 0xc9, 0x41, 0x22, 0xd2,
	0x19, 0x53, 0x41, 0xd0, 0x83, 0x60, 0xea, 0xc1, 0x4b, 0x2e, 0x11, 0x14, 0xbc, 0xc8, 0x64, 0x77,
	0xd8, 0x0e, 0xcd, 0xce, 0x0b, 0xfb, 0xa6, 0xbb, 0xed, 0xd5, 0x9b, 0x37, 0xfd, 0x00, 0x7e, 0x01,
	0x3f, 0x84, 0xe7, 0x1e, 0x7b, 0xf4, 0xa4, 0x92, 0x7c, 0x11, 0x99, 0xfd, 0x47, 0x2c, 0xbb, 0xb1,
	0xed, 0x29, 0x93, 0xf7, 0xde, 0xef, 0xcf, 0xec, 0xfb, 0x31, 0x74, 0x5b, 0xcd, 0x4e, 0x51, 0x68,
	0x13, 0x28, 0x63, 0x75, 0xaa, 0x84, 0x3d, 0xe1, 0xf3, 0x04, 0x2c, 0x78, 0xbe, 0x6b, 0x18, 0x65,
	0x33, 0x48, 0x8e, 0xb8, 0x3b, 0xf3, 0x7a, 0xc8, 0x67, 0x01, 0x60, 0x0c, 0x28, 0xa6, 0x12, 0x95,
	0x48, 0x87, 0x53, 0x65, 0xe5, 0x50, 0x04, 0xa0, 0x4d, 0x81, 0xf5, 0xef, 0x44, 0x10, 0x41, 0x7e,
	0x14, 0xee, 0x54, 0x54, 0xfb, 0x40, 0xef, 0x8e, 0x31, 0x7a, 0xab, 0xec, 0x7b, 0x6d, 0x0f, 0xc3,
	0x44, 0x66, 0xaf, 0xc2, 0x30, 0x51, 0x88, 0xde, 0x63, 0x7a, 0x3b, 0x54, 0x33, 0x15, 0x49, 0x0b,
	0xc9, 0x47, 0x59, 0x14, 0x77, 0xc8, 0x2e, 0x19, 0xdc, 0x98, 0x6c, 0xd5, 0x8d, 0x6a, 0xf8, 0x11,
	0xdd, 0xca, 0x4a, 0x7c, 0x3d, 0xbb, 0x91, 0xcf, 0xde, 0xca, 0xfe, 0xe5, 0xed, 0xf7, 0xe8, 0x83,
	0x46, 0xc1, 0x89, 0xc2, 0x39, 0x18, 0x54, 0xfd, 0x31, 0x65, 0x63, 0x8c, 0xaa, 0xee, 0x3b, 0x39,
	0xd3, 0xa1, 0x93, 0x3a, 0x80, 0x38, 0xd6, 0x88, 0x1a, 0x8c, 0xb3, 0x96, 0x56, 0xe5, 0x8b, 0xd6,
	0xea, 0x46, 0xa5, 0x37, 0xa0, 0x0f, 0xd7, 0xd3, 0xd5, 0xc2, 0x29, 0xf5, 0x57, 0x26, 0x5f, 0x57,
	0x77, 0x9c, 0xa8, 0x4c, 0x26, 0xe1, 0xd5, 0xbe, 0x47, 0xa3, 0xc3, 0x8d, 0x16, 0x87, 0x9f, 0x09,
	0xed, 0xb7, 0x0b, 0x57, 0xf6, 0xbc, 0x80, 0x6e, 0xca, 0x18, 0x8e, 0x8d, 0xdd, 0x21, 0xbb, 0xdd,
	0xc1, 0xcd, 0xfd, 0x7b, 0xbc, 0x58, 0x38, 0x77, 0x0b, 0xe7, 0xe5, 0xc2, 0xf9, 0x01, 0x68, 0x33,
	0x7a, 0x72, 0xf6, 0xab, 0xd7, 0xf9, 0xfe, 0xbb, 0x37, 0x88, 0xb4, 0x3d, 0x3c, 0x9e, 0xf2, 0x00,
	0x62, 0x51, 0xa6, 0xa3, 0xf8, 0xd9, 0xc3, 0xf0, 0x48, 0xd8, 0xd3, 0xb9, 0xc2, 0x1c, 0x80, 0x93,
	0x92, 0x7a, 0xff, 0x47, 0x97, 0x76, 0xc7, 0x18, 0x79, 0x9f, 0x08, 0xf5, 0x1a, 0x42, 0x31, 0xe4,
	0xed, 0x01, 0xe4, 0x8d, 0x6b, 0xf5, 0x9f, 0x5f, 0x19, 0x52, 0xdf, 0xf8, 0x1b, 0xa1, 0xf7, 0xd7,
	0xe5, 0xe0, 0xc5, 0x7f, 0xa8, 0xd7, 0x60, 0xfd, 0xd1, 0xf5, 0xb1, 0xb5, 0xbf, 0xaf, 0x84, 0x6e,
	0xb7, 0xc5, 0xe5, 0xd9, 0x25, 0xf9, 0x2f, 0xe0, 0xfc, 0x97, 0xd7, 0xc3, 0x55, 0x9e, 0x46, 0x6f,
	0xce, 0x16, 0x8c, 0x9c, 0x2f, 0x18, 0xf9, 0xb3, 0x60, 0xe4, 0xcb, 0x92, 0x75, 0xce, 0x97, 0xac,
	0xf3, 0x73, 0xc9, 0x3a, 0x1f, 0xf8, 0x4a, 0x18, 0x1c, 0xef, 0x5e, 0x29, 0x92, 0xff, 0x11, 0x27,
	0xab, 0xcf, 0x8d, 0x0b, 0xc6, 0x74, 0x33, 0x7f, 0x20, 0x9e, 0xfe, 0x0d, 0x00, 0x00, 0xff, 0xff,
	0x88, 0x4e, 0xc1, 0x6d, 0x8d, 0x04, 0x00, 0x00,
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
	SetWithdrawAddress(ctx context.Context, in *MsgSetWithdrawAddress, opts ...grpc.CallOption) (*MsgSetWithdrawAddressResponse, error)
	WithdrawValidatorCommission(ctx context.Context, in *MsgWithdrawValidatorCommission, opts ...grpc.CallOption) (*MsgWithdrawValidatorCommissionResponse, error)
	WithdrawDelegatorReward(ctx context.Context, in *MsgWithdrawDelegatorReward, opts ...grpc.CallOption) (*MsgWithdrawDelegatorRewardResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) SetWithdrawAddress(ctx context.Context, in *MsgSetWithdrawAddress, opts ...grpc.CallOption) (*MsgSetWithdrawAddressResponse, error) {
	out := new(MsgSetWithdrawAddressResponse)
	err := c.cc.Invoke(ctx, "/elysnetwork.elys.incentive.Msg/SetWithdrawAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) WithdrawValidatorCommission(ctx context.Context, in *MsgWithdrawValidatorCommission, opts ...grpc.CallOption) (*MsgWithdrawValidatorCommissionResponse, error) {
	out := new(MsgWithdrawValidatorCommissionResponse)
	err := c.cc.Invoke(ctx, "/elysnetwork.elys.incentive.Msg/WithdrawValidatorCommission", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) WithdrawDelegatorReward(ctx context.Context, in *MsgWithdrawDelegatorReward, opts ...grpc.CallOption) (*MsgWithdrawDelegatorRewardResponse, error) {
	out := new(MsgWithdrawDelegatorRewardResponse)
	err := c.cc.Invoke(ctx, "/elysnetwork.elys.incentive.Msg/WithdrawDelegatorReward", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	SetWithdrawAddress(context.Context, *MsgSetWithdrawAddress) (*MsgSetWithdrawAddressResponse, error)
	WithdrawValidatorCommission(context.Context, *MsgWithdrawValidatorCommission) (*MsgWithdrawValidatorCommissionResponse, error)
	WithdrawDelegatorReward(context.Context, *MsgWithdrawDelegatorReward) (*MsgWithdrawDelegatorRewardResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) SetWithdrawAddress(ctx context.Context, req *MsgSetWithdrawAddress) (*MsgSetWithdrawAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetWithdrawAddress not implemented")
}
func (*UnimplementedMsgServer) WithdrawValidatorCommission(ctx context.Context, req *MsgWithdrawValidatorCommission) (*MsgWithdrawValidatorCommissionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawValidatorCommission not implemented")
}
func (*UnimplementedMsgServer) WithdrawDelegatorReward(ctx context.Context, req *MsgWithdrawDelegatorReward) (*MsgWithdrawDelegatorRewardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawDelegatorReward not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_SetWithdrawAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetWithdrawAddress)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetWithdrawAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elysnetwork.elys.incentive.Msg/SetWithdrawAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetWithdrawAddress(ctx, req.(*MsgSetWithdrawAddress))
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
		FullMethod: "/elysnetwork.elys.incentive.Msg/WithdrawValidatorCommission",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawValidatorCommission(ctx, req.(*MsgWithdrawValidatorCommission))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_WithdrawDelegatorReward_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawDelegatorReward)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WithdrawDelegatorReward(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elysnetwork.elys.incentive.Msg/WithdrawDelegatorReward",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawDelegatorReward(ctx, req.(*MsgWithdrawDelegatorReward))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elysnetwork.elys.incentive.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetWithdrawAddress",
			Handler:    _Msg_SetWithdrawAddress_Handler,
		},
		{
			MethodName: "WithdrawValidatorCommission",
			Handler:    _Msg_WithdrawValidatorCommission_Handler,
		},
		{
			MethodName: "WithdrawDelegatorReward",
			Handler:    _Msg_WithdrawDelegatorReward_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/incentive/tx.proto",
}

func (m *MsgSetWithdrawAddress) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetWithdrawAddress) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetWithdrawAddress) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.WithdrawAddress) > 0 {
		i -= len(m.WithdrawAddress)
		copy(dAtA[i:], m.WithdrawAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.WithdrawAddress)))
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

func (m *MsgSetWithdrawAddressResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSetWithdrawAddressResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSetWithdrawAddressResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.ValidatorAddress) > 0 {
		i -= len(m.ValidatorAddress)
		copy(dAtA[i:], m.ValidatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.ValidatorAddress)))
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

func (m *MsgWithdrawDelegatorReward) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawDelegatorReward) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawDelegatorReward) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
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

func (m *MsgWithdrawDelegatorRewardResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawDelegatorRewardResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawDelegatorRewardResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
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
func (m *MsgSetWithdrawAddress) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DelegatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.WithdrawAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgSetWithdrawAddressResponse) Size() (n int) {
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
	l = len(m.ValidatorAddress)
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

func (m *MsgWithdrawDelegatorReward) Size() (n int) {
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
	return n
}

func (m *MsgWithdrawDelegatorRewardResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgSetWithdrawAddress) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSetWithdrawAddress: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetWithdrawAddress: illegal tag %d (wire type %d)", fieldNum, wire)
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
				return fmt.Errorf("proto: wrong wireType = %d for field WithdrawAddress", wireType)
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
			m.WithdrawAddress = string(dAtA[iNdEx:postIndex])
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
func (m *MsgSetWithdrawAddressResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgSetWithdrawAddressResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSetWithdrawAddressResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgWithdrawDelegatorReward) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgWithdrawDelegatorReward: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawDelegatorReward: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgWithdrawDelegatorRewardResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgWithdrawDelegatorRewardResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawDelegatorRewardResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
