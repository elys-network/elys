// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/tokenomics/query.proto

package types

import (
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// QueryParamsRequest is request type for the Query/Params RPC method.
type QueryParamsRequest struct {
}

func (m *QueryParamsRequest) Reset()         { *m = QueryParamsRequest{} }
func (m *QueryParamsRequest) String() string { return proto.CompactTextString(m) }
func (*QueryParamsRequest) ProtoMessage()    {}
func (*QueryParamsRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4e816dc251df320b, []int{0}
}
func (m *QueryParamsRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsRequest.Merge(m, src)
}
func (m *QueryParamsRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsRequest proto.InternalMessageInfo

// QueryParamsResponse is response type for the Query/Params RPC method.
type QueryParamsResponse struct {
	// params holds all the parameters of this module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
}

func (m *QueryParamsResponse) Reset()         { *m = QueryParamsResponse{} }
func (m *QueryParamsResponse) String() string { return proto.CompactTextString(m) }
func (*QueryParamsResponse) ProtoMessage()    {}
func (*QueryParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4e816dc251df320b, []int{1}
}
func (m *QueryParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryParamsResponse.Merge(m, src)
}
func (m *QueryParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryParamsResponse proto.InternalMessageInfo

func (m *QueryParamsResponse) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

type QueryGetAirdropRequest struct {
	Intent string `protobuf:"bytes,1,opt,name=intent,proto3" json:"intent,omitempty"`
}

func (m *QueryGetAirdropRequest) Reset()         { *m = QueryGetAirdropRequest{} }
func (m *QueryGetAirdropRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetAirdropRequest) ProtoMessage()    {}
func (*QueryGetAirdropRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4e816dc251df320b, []int{2}
}
func (m *QueryGetAirdropRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetAirdropRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetAirdropRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetAirdropRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetAirdropRequest.Merge(m, src)
}
func (m *QueryGetAirdropRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetAirdropRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetAirdropRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetAirdropRequest proto.InternalMessageInfo

func (m *QueryGetAirdropRequest) GetIntent() string {
	if m != nil {
		return m.Intent
	}
	return ""
}

type QueryGetAirdropResponse struct {
	Airdrop Airdrop `protobuf:"bytes,1,opt,name=airdrop,proto3" json:"airdrop"`
}

func (m *QueryGetAirdropResponse) Reset()         { *m = QueryGetAirdropResponse{} }
func (m *QueryGetAirdropResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetAirdropResponse) ProtoMessage()    {}
func (*QueryGetAirdropResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4e816dc251df320b, []int{3}
}
func (m *QueryGetAirdropResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetAirdropResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetAirdropResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetAirdropResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetAirdropResponse.Merge(m, src)
}
func (m *QueryGetAirdropResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetAirdropResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetAirdropResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetAirdropResponse proto.InternalMessageInfo

func (m *QueryGetAirdropResponse) GetAirdrop() Airdrop {
	if m != nil {
		return m.Airdrop
	}
	return Airdrop{}
}

type QueryAllAirdropRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllAirdropRequest) Reset()         { *m = QueryAllAirdropRequest{} }
func (m *QueryAllAirdropRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAllAirdropRequest) ProtoMessage()    {}
func (*QueryAllAirdropRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4e816dc251df320b, []int{4}
}
func (m *QueryAllAirdropRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllAirdropRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllAirdropRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllAirdropRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllAirdropRequest.Merge(m, src)
}
func (m *QueryAllAirdropRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllAirdropRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllAirdropRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllAirdropRequest proto.InternalMessageInfo

func (m *QueryAllAirdropRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

type QueryAllAirdropResponse struct {
	Airdrop    []Airdrop           `protobuf:"bytes,1,rep,name=airdrop,proto3" json:"airdrop"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllAirdropResponse) Reset()         { *m = QueryAllAirdropResponse{} }
func (m *QueryAllAirdropResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAllAirdropResponse) ProtoMessage()    {}
func (*QueryAllAirdropResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4e816dc251df320b, []int{5}
}
func (m *QueryAllAirdropResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllAirdropResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllAirdropResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllAirdropResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllAirdropResponse.Merge(m, src)
}
func (m *QueryAllAirdropResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllAirdropResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllAirdropResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllAirdropResponse proto.InternalMessageInfo

func (m *QueryAllAirdropResponse) GetAirdrop() []Airdrop {
	if m != nil {
		return m.Airdrop
	}
	return nil
}

func (m *QueryAllAirdropResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryParamsRequest)(nil), "elysnetwork.elys.tokenomics.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "elysnetwork.elys.tokenomics.QueryParamsResponse")
	proto.RegisterType((*QueryGetAirdropRequest)(nil), "elysnetwork.elys.tokenomics.QueryGetAirdropRequest")
	proto.RegisterType((*QueryGetAirdropResponse)(nil), "elysnetwork.elys.tokenomics.QueryGetAirdropResponse")
	proto.RegisterType((*QueryAllAirdropRequest)(nil), "elysnetwork.elys.tokenomics.QueryAllAirdropRequest")
	proto.RegisterType((*QueryAllAirdropResponse)(nil), "elysnetwork.elys.tokenomics.QueryAllAirdropResponse")
}

func init() { proto.RegisterFile("elys/tokenomics/query.proto", fileDescriptor_4e816dc251df320b) }

var fileDescriptor_4e816dc251df320b = []byte{
	// 502 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xb1, 0x6e, 0x13, 0x41,
	0x10, 0x86, 0xbd, 0x09, 0x38, 0x62, 0xe9, 0x96, 0x28, 0x20, 0x03, 0x07, 0x5a, 0x42, 0x82, 0x10,
	0xd9, 0x75, 0x12, 0x44, 0xef, 0x08, 0x11, 0xd1, 0x05, 0x57, 0x88, 0x06, 0xd6, 0x66, 0x75, 0x9c,
	0x72, 0xde, 0xb9, 0xdc, 0xae, 0x01, 0x0b, 0xd1, 0xf0, 0x04, 0x48, 0x54, 0xbc, 0x00, 0xa2, 0xe0,
	0x41, 0x52, 0x46, 0xa2, 0xa1, 0x42, 0xc8, 0xe6, 0x2d, 0x68, 0xd0, 0xed, 0x8e, 0x15, 0xc7, 0x67,
	0xd9, 0x07, 0xdd, 0xd9, 0x33, 0xff, 0x3f, 0xdf, 0xdc, 0xfc, 0x3a, 0x7a, 0x55, 0xa7, 0x03, 0x2b,
	0x1d, 0x1c, 0x6a, 0x03, 0xbd, 0xa4, 0x6b, 0xe5, 0x51, 0x5f, 0xe7, 0x03, 0x91, 0xe5, 0xe0, 0x80,
	0xf9, 0xa2, 0xd1, 0xee, 0x0d, 0xe4, 0x87, 0xa2, 0x78, 0x16, 0xa7, 0x8d, 0x8d, 0xd5, 0x18, 0x62,
	0xf0, 0x7d, 0xb2, 0x78, 0x0a, 0x92, 0xc6, 0xb5, 0x18, 0x20, 0x4e, 0xb5, 0x54, 0x59, 0x22, 0x95,
	0x31, 0xe0, 0x94, 0x4b, 0xc0, 0x58, 0xac, 0xde, 0xed, 0x82, 0xed, 0x81, 0x95, 0x1d, 0x65, 0x75,
	0x98, 0x24, 0x5f, 0x6f, 0x77, 0xb4, 0x53, 0xdb, 0x32, 0x53, 0x71, 0x62, 0x7c, 0xf3, 0xd8, 0x69,
	0x9a, 0x2c, 0x53, 0xb9, 0xea, 0x8d, 0x9d, 0xae, 0x4f, 0x57, 0x55, 0x92, 0xbf, 0xcc, 0x21, 0x0b,
	0x65, 0xbe, 0x4a, 0xd9, 0x93, 0xc2, 0xfe, 0xc0, 0x6b, 0xda, 0xfa, 0xa8, 0xaf, 0xad, 0xe3, 0x4f,
	0xe9, 0xa5, 0x33, 0xff, 0xda, 0x0c, 0x8c, 0xd5, 0xac, 0x45, 0xeb, 0xc1, 0xfb, 0x0a, 0xb9, 0x49,
	0xee, 0x5c, 0xdc, 0xb9, 0x25, 0xe6, 0xec, 0x2d, 0x82, 0x78, 0xef, 0xdc, 0xf1, 0xcf, 0x1b, 0xb5,
	0x36, 0x0a, 0x79, 0x93, 0xae, 0x79, 0xe7, 0x7d, 0xed, 0x5a, 0x01, 0x04, 0x67, 0xb2, 0x35, 0x5a,
	0x4f, 0x8c, 0xd3, 0xc6, 0x79, 0xf3, 0x0b, 0x6d, 0xfc, 0xc5, 0x9f, 0xd3, 0xcb, 0x25, 0x05, 0xf2,
	0x3c, 0xa4, 0x2b, 0xb8, 0x0d, 0x02, 0xad, 0xcf, 0x05, 0x42, 0x39, 0x12, 0x8d, 0xa5, 0xfc, 0x05,
	0x22, 0xb5, 0xd2, 0x74, 0x0a, 0xe9, 0x11, 0xa5, 0xa7, 0x6f, 0x1b, 0x47, 0x6c, 0x88, 0x70, 0x1a,
	0x51, 0x9c, 0x46, 0x84, 0x10, 0xe0, 0x69, 0xc4, 0x81, 0x8a, 0x35, 0x6a, 0xdb, 0x13, 0x4a, 0xfe,
	0x95, 0xe0, 0x0e, 0x93, 0x23, 0x66, 0xed, 0xb0, 0xfc, 0x9f, 0x3b, 0xb0, 0xfd, 0x33, 0xa4, 0x4b,
	0x9e, 0x74, 0x73, 0x21, 0x69, 0x40, 0x98, 0x44, 0xdd, 0xf9, 0xb3, 0x4c, 0xcf, 0x7b, 0x54, 0xf6,
	0x99, 0xd0, 0x7a, 0x38, 0x21, 0x93, 0x73, 0x91, 0xca, 0xf9, 0x69, 0x34, 0xab, 0x0b, 0x02, 0x03,
	0xbf, 0xf7, 0xe1, 0xfb, 0xef, 0x4f, 0x4b, 0x1b, 0x6c, 0x5d, 0x16, 0xdd, 0x5b, 0x28, 0x95, 0xb3,
	0xa3, 0xcd, 0xbe, 0x11, 0xba, 0x82, 0x6f, 0x82, 0xed, 0x2e, 0x9e, 0x55, 0x0a, 0x5b, 0xe3, 0xfe,
	0xbf, 0x89, 0x10, 0xf2, 0x81, 0x87, 0x6c, 0x32, 0x31, 0x1f, 0x12, 0x8f, 0x22, 0xdf, 0x85, 0x04,
	0xbf, 0x67, 0x5f, 0x08, 0xa5, 0xe8, 0xd5, 0x4a, 0xd3, 0x2a, 0xc4, 0xa5, 0x2c, 0x56, 0x21, 0x2e,
	0xa7, 0x8b, 0x6f, 0x79, 0xe2, 0x4d, 0x76, 0xbb, 0x12, 0xf1, 0xde, 0xe3, 0xe3, 0x61, 0x44, 0x4e,
	0x86, 0x11, 0xf9, 0x35, 0x8c, 0xc8, 0xc7, 0x51, 0x54, 0x3b, 0x19, 0x45, 0xb5, 0x1f, 0xa3, 0xa8,
	0xf6, 0x4c, 0xc6, 0x89, 0x7b, 0xd5, 0xef, 0x88, 0x2e, 0xf4, 0x66, 0x58, 0xbd, 0x9d, 0x34, 0x73,
	0x83, 0x4c, 0xdb, 0x4e, 0xdd, 0x7f, 0x5f, 0x76, 0xff, 0x06, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x2d,
	0x32, 0xc1, 0x38, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Parameters queries the parameters of the module.
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	// Queries a list of Airdrop items.
	Airdrop(ctx context.Context, in *QueryGetAirdropRequest, opts ...grpc.CallOption) (*QueryGetAirdropResponse, error)
	AirdropAll(ctx context.Context, in *QueryAllAirdropRequest, opts ...grpc.CallOption) (*QueryAllAirdropResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/elysnetwork.elys.tokenomics.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Airdrop(ctx context.Context, in *QueryGetAirdropRequest, opts ...grpc.CallOption) (*QueryGetAirdropResponse, error) {
	out := new(QueryGetAirdropResponse)
	err := c.cc.Invoke(ctx, "/elysnetwork.elys.tokenomics.Query/Airdrop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) AirdropAll(ctx context.Context, in *QueryAllAirdropRequest, opts ...grpc.CallOption) (*QueryAllAirdropResponse, error) {
	out := new(QueryAllAirdropResponse)
	err := c.cc.Invoke(ctx, "/elysnetwork.elys.tokenomics.Query/AirdropAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Queries a list of Airdrop items.
	Airdrop(context.Context, *QueryGetAirdropRequest) (*QueryGetAirdropResponse, error)
	AirdropAll(context.Context, *QueryAllAirdropRequest) (*QueryAllAirdropResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) Airdrop(ctx context.Context, req *QueryGetAirdropRequest) (*QueryGetAirdropResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Airdrop not implemented")
}
func (*UnimplementedQueryServer) AirdropAll(ctx context.Context, req *QueryAllAirdropRequest) (*QueryAllAirdropResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AirdropAll not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elysnetwork.elys.tokenomics.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Airdrop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetAirdropRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Airdrop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elysnetwork.elys.tokenomics.Query/Airdrop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Airdrop(ctx, req.(*QueryGetAirdropRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_AirdropAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllAirdropRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).AirdropAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elysnetwork.elys.tokenomics.Query/AirdropAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).AirdropAll(ctx, req.(*QueryAllAirdropRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elysnetwork.elys.tokenomics.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Airdrop",
			Handler:    _Query_Airdrop_Handler,
		},
		{
			MethodName: "AirdropAll",
			Handler:    _Query_AirdropAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/tokenomics/query.proto",
}

func (m *QueryParamsRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryGetAirdropRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetAirdropRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetAirdropRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Intent) > 0 {
		i -= len(m.Intent)
		copy(dAtA[i:], m.Intent)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Intent)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryGetAirdropResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetAirdropResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetAirdropResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Airdrop.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *QueryAllAirdropRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllAirdropRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllAirdropRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryAllAirdropResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllAirdropResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllAirdropResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Pagination != nil {
		{
			size, err := m.Pagination.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintQuery(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Airdrop) > 0 {
		for iNdEx := len(m.Airdrop) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Airdrop[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintQuery(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryParamsRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryGetAirdropRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Intent)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryGetAirdropResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Airdrop.Size()
	n += 1 + l + sovQuery(uint64(l))
	return n
}

func (m *QueryAllAirdropRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryAllAirdropResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Airdrop) > 0 {
		for _, e := range m.Airdrop {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	if m.Pagination != nil {
		l = m.Pagination.Size()
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryParamsRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryParamsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryGetAirdropRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryGetAirdropRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetAirdropRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Intent", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Intent = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryGetAirdropResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryGetAirdropResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetAirdropResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Airdrop", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Airdrop.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryAllAirdropRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryAllAirdropRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllAirdropRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageRequest{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func (m *QueryAllAirdropResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
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
			return fmt.Errorf("proto: QueryAllAirdropResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllAirdropResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Airdrop", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Airdrop = append(m.Airdrop, Airdrop{})
			if err := m.Airdrop[len(m.Airdrop)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pagination", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
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
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Pagination == nil {
				m.Pagination = &query.PageResponse{}
			}
			if err := m.Pagination.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
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
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
					return 0, ErrIntOverflowQuery
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
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
