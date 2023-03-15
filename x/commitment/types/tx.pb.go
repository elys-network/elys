// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/commitment/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

type MsgCommitTokens struct {
	Creator string                                 `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Amount  github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount" yaml:"amount"`
	Denom   string                                 `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
}

func (m *MsgCommitTokens) Reset()         { *m = MsgCommitTokens{} }
func (m *MsgCommitTokens) String() string { return proto.CompactTextString(m) }
func (*MsgCommitTokens) ProtoMessage()    {}
func (*MsgCommitTokens) Descriptor() ([]byte, []int) {
	return fileDescriptor_eda78aafe7eb119c, []int{0}
}
func (m *MsgCommitTokens) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCommitTokens) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCommitTokens.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCommitTokens) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCommitTokens.Merge(m, src)
}
func (m *MsgCommitTokens) XXX_Size() int {
	return m.Size()
}
func (m *MsgCommitTokens) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCommitTokens.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCommitTokens proto.InternalMessageInfo

func (m *MsgCommitTokens) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgCommitTokens) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

type MsgCommitTokensResponse struct {
}

func (m *MsgCommitTokensResponse) Reset()         { *m = MsgCommitTokensResponse{} }
func (m *MsgCommitTokensResponse) String() string { return proto.CompactTextString(m) }
func (*MsgCommitTokensResponse) ProtoMessage()    {}
func (*MsgCommitTokensResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_eda78aafe7eb119c, []int{1}
}
func (m *MsgCommitTokensResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgCommitTokensResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgCommitTokensResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgCommitTokensResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgCommitTokensResponse.Merge(m, src)
}
func (m *MsgCommitTokensResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgCommitTokensResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgCommitTokensResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgCommitTokensResponse proto.InternalMessageInfo

type MsgUncommitTokens struct {
	Creator string                                 `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Amount  github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount" yaml:"amount"`
	Denom   string                                 `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
}

func (m *MsgUncommitTokens) Reset()         { *m = MsgUncommitTokens{} }
func (m *MsgUncommitTokens) String() string { return proto.CompactTextString(m) }
func (*MsgUncommitTokens) ProtoMessage()    {}
func (*MsgUncommitTokens) Descriptor() ([]byte, []int) {
	return fileDescriptor_eda78aafe7eb119c, []int{2}
}
func (m *MsgUncommitTokens) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUncommitTokens) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUncommitTokens.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUncommitTokens) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUncommitTokens.Merge(m, src)
}
func (m *MsgUncommitTokens) XXX_Size() int {
	return m.Size()
}
func (m *MsgUncommitTokens) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUncommitTokens.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUncommitTokens proto.InternalMessageInfo

func (m *MsgUncommitTokens) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgUncommitTokens) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

type MsgUncommitTokensResponse struct {
}

func (m *MsgUncommitTokensResponse) Reset()         { *m = MsgUncommitTokensResponse{} }
func (m *MsgUncommitTokensResponse) String() string { return proto.CompactTextString(m) }
func (*MsgUncommitTokensResponse) ProtoMessage()    {}
func (*MsgUncommitTokensResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_eda78aafe7eb119c, []int{3}
}
func (m *MsgUncommitTokensResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgUncommitTokensResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgUncommitTokensResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgUncommitTokensResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgUncommitTokensResponse.Merge(m, src)
}
func (m *MsgUncommitTokensResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgUncommitTokensResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgUncommitTokensResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgUncommitTokensResponse proto.InternalMessageInfo

type MsgWithdrawTokens struct {
	Creator string                                 `protobuf:"bytes,1,opt,name=creator,proto3" json:"creator,omitempty"`
	Amount  github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,2,opt,name=amount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"amount" yaml:"amount"`
	Denom   string                                 `protobuf:"bytes,3,opt,name=denom,proto3" json:"denom,omitempty"`
}

func (m *MsgWithdrawTokens) Reset()         { *m = MsgWithdrawTokens{} }
func (m *MsgWithdrawTokens) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawTokens) ProtoMessage()    {}
func (*MsgWithdrawTokens) Descriptor() ([]byte, []int) {
	return fileDescriptor_eda78aafe7eb119c, []int{4}
}
func (m *MsgWithdrawTokens) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawTokens) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawTokens.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawTokens) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawTokens.Merge(m, src)
}
func (m *MsgWithdrawTokens) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawTokens) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawTokens.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawTokens proto.InternalMessageInfo

func (m *MsgWithdrawTokens) GetCreator() string {
	if m != nil {
		return m.Creator
	}
	return ""
}

func (m *MsgWithdrawTokens) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

type MsgWithdrawTokensResponse struct {
}

func (m *MsgWithdrawTokensResponse) Reset()         { *m = MsgWithdrawTokensResponse{} }
func (m *MsgWithdrawTokensResponse) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawTokensResponse) ProtoMessage()    {}
func (*MsgWithdrawTokensResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_eda78aafe7eb119c, []int{5}
}
func (m *MsgWithdrawTokensResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawTokensResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawTokensResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawTokensResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawTokensResponse.Merge(m, src)
}
func (m *MsgWithdrawTokensResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawTokensResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawTokensResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawTokensResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgCommitTokens)(nil), "elysnetwork.elys.commitment.MsgCommitTokens")
	proto.RegisterType((*MsgCommitTokensResponse)(nil), "elysnetwork.elys.commitment.MsgCommitTokensResponse")
	proto.RegisterType((*MsgUncommitTokens)(nil), "elysnetwork.elys.commitment.MsgUncommitTokens")
	proto.RegisterType((*MsgUncommitTokensResponse)(nil), "elysnetwork.elys.commitment.MsgUncommitTokensResponse")
	proto.RegisterType((*MsgWithdrawTokens)(nil), "elysnetwork.elys.commitment.MsgWithdrawTokens")
	proto.RegisterType((*MsgWithdrawTokensResponse)(nil), "elysnetwork.elys.commitment.MsgWithdrawTokensResponse")
}

func init() { proto.RegisterFile("elys/commitment/tx.proto", fileDescriptor_eda78aafe7eb119c) }

var fileDescriptor_eda78aafe7eb119c = []byte{
	// 367 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x48, 0xcd, 0xa9, 0x2c,
	0xd6, 0x4f, 0xce, 0xcf, 0xcd, 0xcd, 0x2c, 0xc9, 0x4d, 0xcd, 0x2b, 0xd1, 0x2f, 0xa9, 0xd0, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x06, 0xc9, 0xe4, 0xa5, 0x96, 0x94, 0xe7, 0x17, 0x65, 0xeb,
	0x81, 0xd8, 0x7a, 0x08, 0x55, 0x52, 0x22, 0xe9, 0xf9, 0xe9, 0xf9, 0x60, 0x75, 0xfa, 0x20, 0x16,
	0x44, 0x8b, 0xd2, 0x2c, 0x46, 0x2e, 0x7e, 0xdf, 0xe2, 0x74, 0x67, 0xb0, 0xba, 0x90, 0xfc, 0xec,
	0xd4, 0xbc, 0x62, 0x21, 0x09, 0x2e, 0xf6, 0xe4, 0xa2, 0xd4, 0xc4, 0x92, 0xfc, 0x22, 0x09, 0x46,
	0x05, 0x46, 0x0d, 0xce, 0x20, 0x18, 0x57, 0x28, 0x9c, 0x8b, 0x2d, 0x31, 0x37, 0xbf, 0x34, 0xaf,
	0x44, 0x82, 0x09, 0x24, 0xe1, 0x64, 0x7f, 0xe2, 0x9e, 0x3c, 0xc3, 0xad, 0x7b, 0xf2, 0x6a, 0xe9,
	0x99, 0x25, 0x19, 0xa5, 0x49, 0x20, 0xeb, 0xf4, 0x93, 0xf3, 0x8b, 0x73, 0xf3, 0x8b, 0xa1, 0x94,
	0x6e, 0x71, 0x4a, 0xb6, 0x7e, 0x49, 0x65, 0x41, 0x6a, 0xb1, 0x9e, 0x67, 0x5e, 0xc9, 0xa7, 0x7b,
	0xf2, 0xbc, 0x95, 0x89, 0xb9, 0x39, 0x56, 0x4a, 0x10, 0x53, 0x94, 0x82, 0xa0, 0xc6, 0x09, 0x89,
	0x70, 0xb1, 0xa6, 0xa4, 0xe6, 0xe5, 0xe7, 0x4a, 0x30, 0x83, 0x2d, 0x84, 0x70, 0x94, 0x24, 0xb9,
	0xc4, 0xd1, 0xdc, 0x16, 0x94, 0x5a, 0x5c, 0x90, 0x9f, 0x57, 0x9c, 0xaa, 0x34, 0x87, 0x91, 0x4b,
	0xd0, 0xb7, 0x38, 0x3d, 0x34, 0x2f, 0x79, 0x50, 0xba, 0x5c, 0x9a, 0x4b, 0x12, 0xc3, 0x75, 0xe8,
	0x6e, 0x0f, 0xcf, 0x2c, 0xc9, 0x48, 0x29, 0x4a, 0x2c, 0x1f, 0x9c, 0x6e, 0x47, 0x75, 0x1d, 0xcc,
	0xed, 0x46, 0x6f, 0x99, 0xb8, 0x98, 0x7d, 0x8b, 0xd3, 0x85, 0x8a, 0xb8, 0x78, 0x50, 0xd2, 0x8c,
	0x8e, 0x1e, 0x9e, 0xb4, 0xa7, 0x87, 0x16, 0x8b, 0x52, 0x26, 0xa4, 0xa8, 0x86, 0xd9, 0x2d, 0x54,
	0xc1, 0xc5, 0x87, 0x16, 0xdf, 0x7a, 0x84, 0xcc, 0x41, 0x55, 0x2f, 0x65, 0x46, 0x9a, 0x7a, 0x64,
	0x9b, 0xd1, 0x62, 0x8b, 0xa0, 0xcd, 0xa8, 0xea, 0x09, 0xdb, 0x8c, 0x3d, 0xbc, 0x9d, 0x3c, 0x4f,
	0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18,
	0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x21, 0x4a, 0x1f, 0x29, 0xf6, 0x41, 0xe6, 0xe9,
	0x42, 0x0d, 0x07, 0x73, 0xf4, 0x2b, 0x50, 0x0a, 0x08, 0x50, 0x52, 0x48, 0x62, 0x03, 0xe7, 0x78,
	0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x5d, 0x76, 0x5a, 0x0a, 0x40, 0x04, 0x00, 0x00,
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
	CommitTokens(ctx context.Context, in *MsgCommitTokens, opts ...grpc.CallOption) (*MsgCommitTokensResponse, error)
	UncommitTokens(ctx context.Context, in *MsgUncommitTokens, opts ...grpc.CallOption) (*MsgUncommitTokensResponse, error)
	WithdrawTokens(ctx context.Context, in *MsgWithdrawTokens, opts ...grpc.CallOption) (*MsgWithdrawTokensResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CommitTokens(ctx context.Context, in *MsgCommitTokens, opts ...grpc.CallOption) (*MsgCommitTokensResponse, error) {
	out := new(MsgCommitTokensResponse)
	err := c.cc.Invoke(ctx, "/elysnetwork.elys.commitment.Msg/CommitTokens", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UncommitTokens(ctx context.Context, in *MsgUncommitTokens, opts ...grpc.CallOption) (*MsgUncommitTokensResponse, error) {
	out := new(MsgUncommitTokensResponse)
	err := c.cc.Invoke(ctx, "/elysnetwork.elys.commitment.Msg/UncommitTokens", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) WithdrawTokens(ctx context.Context, in *MsgWithdrawTokens, opts ...grpc.CallOption) (*MsgWithdrawTokensResponse, error) {
	out := new(MsgWithdrawTokensResponse)
	err := c.cc.Invoke(ctx, "/elysnetwork.elys.commitment.Msg/WithdrawTokens", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	CommitTokens(context.Context, *MsgCommitTokens) (*MsgCommitTokensResponse, error)
	UncommitTokens(context.Context, *MsgUncommitTokens) (*MsgUncommitTokensResponse, error)
	WithdrawTokens(context.Context, *MsgWithdrawTokens) (*MsgWithdrawTokensResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) CommitTokens(ctx context.Context, req *MsgCommitTokens) (*MsgCommitTokensResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommitTokens not implemented")
}
func (*UnimplementedMsgServer) UncommitTokens(ctx context.Context, req *MsgUncommitTokens) (*MsgUncommitTokensResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UncommitTokens not implemented")
}
func (*UnimplementedMsgServer) WithdrawTokens(ctx context.Context, req *MsgWithdrawTokens) (*MsgWithdrawTokensResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawTokens not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_CommitTokens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCommitTokens)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CommitTokens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elysnetwork.elys.commitment.Msg/CommitTokens",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CommitTokens(ctx, req.(*MsgCommitTokens))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UncommitTokens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUncommitTokens)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UncommitTokens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elysnetwork.elys.commitment.Msg/UncommitTokens",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UncommitTokens(ctx, req.(*MsgUncommitTokens))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_WithdrawTokens_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawTokens)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WithdrawTokens(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elysnetwork.elys.commitment.Msg/WithdrawTokens",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawTokens(ctx, req.(*MsgWithdrawTokens))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elysnetwork.elys.commitment.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CommitTokens",
			Handler:    _Msg_CommitTokens_Handler,
		},
		{
			MethodName: "UncommitTokens",
			Handler:    _Msg_UncommitTokens_Handler,
		},
		{
			MethodName: "WithdrawTokens",
			Handler:    _Msg_WithdrawTokens_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/commitment/tx.proto",
}

func (m *MsgCommitTokens) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCommitTokens) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCommitTokens) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgCommitTokensResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgCommitTokensResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgCommitTokensResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgUncommitTokens) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUncommitTokens) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUncommitTokens) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgUncommitTokensResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgUncommitTokensResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgUncommitTokensResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *MsgWithdrawTokens) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawTokens) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawTokens) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	{
		size := m.Amount.Size()
		i -= size
		if _, err := m.Amount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Creator) > 0 {
		i -= len(m.Creator)
		copy(dAtA[i:], m.Creator)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Creator)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgWithdrawTokensResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawTokensResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawTokensResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
func (m *MsgCommitTokens) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovTx(uint64(l))
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgCommitTokensResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgUncommitTokens) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovTx(uint64(l))
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgUncommitTokensResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *MsgWithdrawTokens) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Creator)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = m.Amount.Size()
	n += 1 + l + sovTx(uint64(l))
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgWithdrawTokensResponse) Size() (n int) {
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
func (m *MsgCommitTokens) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgCommitTokens: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCommitTokens: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
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
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
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
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *MsgCommitTokensResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgCommitTokensResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgCommitTokensResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgUncommitTokens) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgUncommitTokens: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUncommitTokens: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
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
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
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
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *MsgUncommitTokensResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgUncommitTokensResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgUncommitTokensResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *MsgWithdrawTokens) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgWithdrawTokens: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawTokens: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Creator", wireType)
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
			m.Creator = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
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
			if err := m.Amount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
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
func (m *MsgWithdrawTokensResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgWithdrawTokensResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawTokensResponse: illegal tag %d (wire type %d)", fieldNum, wire)
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
