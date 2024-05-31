// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/membershiptier/query.proto

package types

import (
	context "context"
	fmt "fmt"
	query "github.com/cosmos/cosmos-sdk/types/query"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
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
	return fileDescriptor_92aae422e277a0e0, []int{0}
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
	return fileDescriptor_92aae422e277a0e0, []int{1}
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

type QueryGetPortfolioRequest struct {
	User      string `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	AssetType string `protobuf:"bytes,2,opt,name=assetType,proto3" json:"assetType,omitempty"`
}

func (m *QueryGetPortfolioRequest) Reset()         { *m = QueryGetPortfolioRequest{} }
func (m *QueryGetPortfolioRequest) String() string { return proto.CompactTextString(m) }
func (*QueryGetPortfolioRequest) ProtoMessage()    {}
func (*QueryGetPortfolioRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_92aae422e277a0e0, []int{2}
}
func (m *QueryGetPortfolioRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetPortfolioRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetPortfolioRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetPortfolioRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetPortfolioRequest.Merge(m, src)
}
func (m *QueryGetPortfolioRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetPortfolioRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetPortfolioRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetPortfolioRequest proto.InternalMessageInfo

func (m *QueryGetPortfolioRequest) GetUser() string {
	if m != nil {
		return m.User
	}
	return ""
}

func (m *QueryGetPortfolioRequest) GetAssetType() string {
	if m != nil {
		return m.AssetType
	}
	return ""
}

type QueryGetPortfolioResponse struct {
	Portfolio []Portfolio `protobuf:"bytes,1,rep,name=portfolio,proto3" json:"portfolio"`
}

func (m *QueryGetPortfolioResponse) Reset()         { *m = QueryGetPortfolioResponse{} }
func (m *QueryGetPortfolioResponse) String() string { return proto.CompactTextString(m) }
func (*QueryGetPortfolioResponse) ProtoMessage()    {}
func (*QueryGetPortfolioResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_92aae422e277a0e0, []int{3}
}
func (m *QueryGetPortfolioResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryGetPortfolioResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryGetPortfolioResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryGetPortfolioResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryGetPortfolioResponse.Merge(m, src)
}
func (m *QueryGetPortfolioResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryGetPortfolioResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryGetPortfolioResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryGetPortfolioResponse proto.InternalMessageInfo

func (m *QueryGetPortfolioResponse) GetPortfolio() []Portfolio {
	if m != nil {
		return m.Portfolio
	}
	return nil
}

type QueryAllPortfolioRequest struct {
	Pagination *query.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllPortfolioRequest) Reset()         { *m = QueryAllPortfolioRequest{} }
func (m *QueryAllPortfolioRequest) String() string { return proto.CompactTextString(m) }
func (*QueryAllPortfolioRequest) ProtoMessage()    {}
func (*QueryAllPortfolioRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_92aae422e277a0e0, []int{4}
}
func (m *QueryAllPortfolioRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllPortfolioRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllPortfolioRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllPortfolioRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllPortfolioRequest.Merge(m, src)
}
func (m *QueryAllPortfolioRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllPortfolioRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllPortfolioRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllPortfolioRequest proto.InternalMessageInfo

func (m *QueryAllPortfolioRequest) GetPagination() *query.PageRequest {
	if m != nil {
		return m.Pagination
	}
	return nil
}

type QueryAllPortfolioResponse struct {
	Portfolio  []Portfolio         `protobuf:"bytes,1,rep,name=portfolio,proto3" json:"portfolio"`
	Pagination *query.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

func (m *QueryAllPortfolioResponse) Reset()         { *m = QueryAllPortfolioResponse{} }
func (m *QueryAllPortfolioResponse) String() string { return proto.CompactTextString(m) }
func (*QueryAllPortfolioResponse) ProtoMessage()    {}
func (*QueryAllPortfolioResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_92aae422e277a0e0, []int{5}
}
func (m *QueryAllPortfolioResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryAllPortfolioResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryAllPortfolioResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryAllPortfolioResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryAllPortfolioResponse.Merge(m, src)
}
func (m *QueryAllPortfolioResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryAllPortfolioResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryAllPortfolioResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryAllPortfolioResponse proto.InternalMessageInfo

func (m *QueryAllPortfolioResponse) GetPortfolio() []Portfolio {
	if m != nil {
		return m.Portfolio
	}
	return nil
}

func (m *QueryAllPortfolioResponse) GetPagination() *query.PageResponse {
	if m != nil {
		return m.Pagination
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryParamsRequest)(nil), "elys.membershiptier.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "elys.membershiptier.QueryParamsResponse")
	proto.RegisterType((*QueryGetPortfolioRequest)(nil), "elys.membershiptier.QueryGetPortfolioRequest")
	proto.RegisterType((*QueryGetPortfolioResponse)(nil), "elys.membershiptier.QueryGetPortfolioResponse")
	proto.RegisterType((*QueryAllPortfolioRequest)(nil), "elys.membershiptier.QueryAllPortfolioRequest")
	proto.RegisterType((*QueryAllPortfolioResponse)(nil), "elys.membershiptier.QueryAllPortfolioResponse")
}

func init() { proto.RegisterFile("elys/membershiptier/query.proto", fileDescriptor_92aae422e277a0e0) }

var fileDescriptor_92aae422e277a0e0 = []byte{
	// 528 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0x3f, 0x6f, 0x13, 0x31,
	0x18, 0xc6, 0xe3, 0xb4, 0x44, 0x8a, 0xcb, 0xe4, 0x76, 0x08, 0xa1, 0xba, 0x46, 0x87, 0x44, 0xa3,
	0xa2, 0xd8, 0x34, 0x99, 0x58, 0x10, 0xcd, 0xd0, 0x2e, 0x20, 0x85, 0x88, 0x89, 0x05, 0xf9, 0xaa,
	0x97, 0xeb, 0x89, 0xbb, 0xf3, 0xf5, 0xec, 0x00, 0x51, 0xd5, 0x85, 0x1d, 0x09, 0x89, 0x6f, 0xc0,
	0xc4, 0xcc, 0xa7, 0xe8, 0x58, 0x89, 0x85, 0x09, 0xa1, 0x84, 0x2f, 0xc0, 0x37, 0x40, 0x67, 0x3b,
	0x49, 0x8f, 0xba, 0x7f, 0x90, 0xba, 0x9d, 0xec, 0xe7, 0x79, 0xdf, 0xdf, 0xeb, 0xf7, 0xd1, 0xe1,
	0x0d, 0x88, 0xc7, 0x92, 0x25, 0x90, 0x04, 0x90, 0xcb, 0x83, 0x28, 0x53, 0x11, 0xe4, 0xec, 0x70,
	0x04, 0xf9, 0x98, 0x66, 0xb9, 0x50, 0x82, 0xac, 0x16, 0x02, 0x5a, 0x16, 0x34, 0xd7, 0x42, 0x11,
	0x0a, 0x7d, 0xcf, 0x8a, 0x2f, 0x23, 0x6d, 0xae, 0x87, 0x42, 0x84, 0x31, 0x30, 0x9e, 0x45, 0x8c,
	0xa7, 0xa9, 0x50, 0x5c, 0x45, 0x22, 0x95, 0xf6, 0x76, 0x6b, 0x5f, 0xc8, 0x44, 0x48, 0x16, 0x70,
	0x09, 0xa6, 0x03, 0x7b, 0xbb, 0x1d, 0x80, 0xe2, 0xdb, 0x2c, 0xe3, 0x61, 0x94, 0x6a, 0xb1, 0xd5,
	0xb6, 0x5c, 0x54, 0x19, 0xcf, 0x79, 0x32, 0xab, 0x76, 0xcf, 0xa9, 0x10, 0xb9, 0x7a, 0x2d, 0xe2,
	0xc8, 0x02, 0xf9, 0x6b, 0x98, 0x3c, 0x2f, 0x1a, 0x0d, 0xb4, 0x73, 0x08, 0x87, 0x23, 0x90, 0xca,
	0x1f, 0xe0, 0xd5, 0xd2, 0xa9, 0xcc, 0x44, 0x2a, 0x81, 0x3c, 0xc2, 0x35, 0xd3, 0xa1, 0x81, 0x5a,
	0xa8, 0xbd, 0xd2, 0xbd, 0x4b, 0x1d, 0x93, 0x53, 0x63, 0xea, 0x2f, 0x9f, 0xfc, 0xdc, 0xa8, 0x0c,
	0xad, 0xc1, 0x7f, 0x8a, 0x1b, 0xba, 0xe2, 0x1e, 0xa8, 0xc1, 0x0c, 0xc1, 0x76, 0x23, 0x04, 0x2f,
	0x8f, 0x24, 0xe4, 0xba, 0x68, 0x7d, 0xa8, 0xbf, 0xc9, 0x3a, 0xae, 0x73, 0x29, 0x41, 0xbd, 0x18,
	0x67, 0xd0, 0xa8, 0xea, 0x8b, 0xc5, 0x81, 0xff, 0x0a, 0xdf, 0x71, 0x54, 0xb3, 0x94, 0x7d, 0x5c,
	0x9f, 0x4f, 0xd9, 0x40, 0xad, 0xa5, 0xf6, 0x4a, 0xd7, 0x73, 0x83, 0xce, 0x54, 0x96, 0x75, 0x61,
	0xf3, 0x03, 0x8b, 0xbb, 0x13, 0xc7, 0xe7, 0x70, 0x77, 0x31, 0x5e, 0x6c, 0xc3, 0xbe, 0xc4, 0x7d,
	0x6a, 0x56, 0x47, 0x8b, 0xd5, 0x51, 0x13, 0x0e, 0xbb, 0x3a, 0x3a, 0xe0, 0x21, 0x58, 0xef, 0xf0,
	0x8c, 0xd3, 0xff, 0x8a, 0xec, 0x14, 0xe5, 0x26, 0x37, 0x37, 0x05, 0xd9, 0x2b, 0x91, 0x56, 0x35,
	0xe9, 0xe6, 0x95, 0xa4, 0x06, 0xe0, 0x2c, 0x6a, 0xf7, 0xcf, 0x12, 0xbe, 0xa5, 0x51, 0xc9, 0x47,
	0x84, 0x6b, 0x66, 0xc1, 0x64, 0xd3, 0x89, 0x73, 0x3e, 0x4d, 0xcd, 0xf6, 0xd5, 0x42, 0xd3, 0xd3,
	0x7f, 0xf8, 0xe1, 0xfb, 0xef, 0xcf, 0xd5, 0x2d, 0xd2, 0x66, 0x85, 0xa3, 0x93, 0x82, 0x7a, 0x27,
	0xf2, 0x37, 0xec, 0xe2, 0xa8, 0x93, 0x6f, 0x08, 0xd7, 0xe7, 0x2f, 0x40, 0x3a, 0x17, 0x77, 0x72,
	0x04, 0xaf, 0x49, 0xaf, 0x2b, 0xb7, 0x78, 0xbb, 0x1a, 0xef, 0x09, 0x79, 0x7c, 0x0d, 0xbc, 0x99,
	0x99, 0x1d, 0x15, 0x79, 0x3e, 0x66, 0x47, 0xf3, 0xf4, 0x1e, 0x93, 0x2f, 0x08, 0xdf, 0x9e, 0x57,
	0xdf, 0x89, 0xe3, 0xcb, 0xb8, 0x1d, 0x09, 0xbc, 0x8c, 0xdb, 0x95, 0x25, 0xbf, 0xa7, 0xb9, 0x3b,
	0xe4, 0xc1, 0x7f, 0x70, 0xf7, 0x9f, 0x9d, 0x4c, 0x3c, 0x74, 0x3a, 0xf1, 0xd0, 0xaf, 0x89, 0x87,
	0x3e, 0x4d, 0xbd, 0xca, 0xe9, 0xd4, 0xab, 0xfc, 0x98, 0x7a, 0x95, 0x97, 0xbd, 0x30, 0x52, 0x07,
	0xa3, 0x80, 0xee, 0x8b, 0xc4, 0x51, 0xf0, 0xfd, 0xbf, 0x25, 0xd5, 0x38, 0x03, 0x19, 0xd4, 0xf4,
	0xff, 0xa6, 0xf7, 0x37, 0x00, 0x00, 0xff, 0xff, 0x93, 0x33, 0x6f, 0x25, 0x4e, 0x05, 0x00, 0x00,
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
	// Queries a list of Portfolio items.
	Portfolio(ctx context.Context, in *QueryGetPortfolioRequest, opts ...grpc.CallOption) (*QueryGetPortfolioResponse, error)
	PortfolioAll(ctx context.Context, in *QueryAllPortfolioRequest, opts ...grpc.CallOption) (*QueryAllPortfolioResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/elys.membershiptier.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Portfolio(ctx context.Context, in *QueryGetPortfolioRequest, opts ...grpc.CallOption) (*QueryGetPortfolioResponse, error) {
	out := new(QueryGetPortfolioResponse)
	err := c.cc.Invoke(ctx, "/elys.membershiptier.Query/Portfolio", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) PortfolioAll(ctx context.Context, in *QueryAllPortfolioRequest, opts ...grpc.CallOption) (*QueryAllPortfolioResponse, error) {
	out := new(QueryAllPortfolioResponse)
	err := c.cc.Invoke(ctx, "/elys.membershiptier.Query/PortfolioAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// Queries a list of Portfolio items.
	Portfolio(context.Context, *QueryGetPortfolioRequest) (*QueryGetPortfolioResponse, error)
	PortfolioAll(context.Context, *QueryAllPortfolioRequest) (*QueryAllPortfolioResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) Portfolio(ctx context.Context, req *QueryGetPortfolioRequest) (*QueryGetPortfolioResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Portfolio not implemented")
}
func (*UnimplementedQueryServer) PortfolioAll(ctx context.Context, req *QueryAllPortfolioRequest) (*QueryAllPortfolioResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PortfolioAll not implemented")
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
		FullMethod: "/elys.membershiptier.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Portfolio_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryGetPortfolioRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Portfolio(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.membershiptier.Query/Portfolio",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Portfolio(ctx, req.(*QueryGetPortfolioRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_PortfolioAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryAllPortfolioRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).PortfolioAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.membershiptier.Query/PortfolioAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).PortfolioAll(ctx, req.(*QueryAllPortfolioRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elys.membershiptier.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "Portfolio",
			Handler:    _Query_Portfolio_Handler,
		},
		{
			MethodName: "PortfolioAll",
			Handler:    _Query_PortfolioAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/membershiptier/query.proto",
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

func (m *QueryGetPortfolioRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetPortfolioRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetPortfolioRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.AssetType) > 0 {
		i -= len(m.AssetType)
		copy(dAtA[i:], m.AssetType)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.AssetType)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.User) > 0 {
		i -= len(m.User)
		copy(dAtA[i:], m.User)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.User)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryGetPortfolioResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryGetPortfolioResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryGetPortfolioResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Portfolio) > 0 {
		for iNdEx := len(m.Portfolio) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Portfolio[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryAllPortfolioRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllPortfolioRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllPortfolioRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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

func (m *QueryAllPortfolioResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryAllPortfolioResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryAllPortfolioResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
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
	if len(m.Portfolio) > 0 {
		for iNdEx := len(m.Portfolio) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Portfolio[iNdEx].MarshalToSizedBuffer(dAtA[:i])
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

func (m *QueryGetPortfolioRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.User)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	l = len(m.AssetType)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryGetPortfolioResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Portfolio) > 0 {
		for _, e := range m.Portfolio {
			l = e.Size()
			n += 1 + l + sovQuery(uint64(l))
		}
	}
	return n
}

func (m *QueryAllPortfolioRequest) Size() (n int) {
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

func (m *QueryAllPortfolioResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Portfolio) > 0 {
		for _, e := range m.Portfolio {
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
func (m *QueryGetPortfolioRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryGetPortfolioRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetPortfolioRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field User", wireType)
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
			m.User = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AssetType", wireType)
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
			m.AssetType = string(dAtA[iNdEx:postIndex])
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
func (m *QueryGetPortfolioResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryGetPortfolioResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryGetPortfolioResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Portfolio", wireType)
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
			m.Portfolio = append(m.Portfolio, Portfolio{})
			if err := m.Portfolio[len(m.Portfolio)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *QueryAllPortfolioRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAllPortfolioRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllPortfolioRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *QueryAllPortfolioResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryAllPortfolioResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryAllPortfolioResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Portfolio", wireType)
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
			m.Portfolio = append(m.Portfolio, Portfolio{})
			if err := m.Portfolio[len(m.Portfolio)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
