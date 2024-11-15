// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: elys/stablestake/query.proto

package types

import (
	context "context"
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types/query"
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
	return fileDescriptor_9717e1e4014b4459, []int{0}
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
	return fileDescriptor_9717e1e4014b4459, []int{1}
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

// QueryBorrowRatioRequest is request type for the Query/BorrowRatio RPC method.
type QueryBorrowRatioRequest struct {
}

func (m *QueryBorrowRatioRequest) Reset()         { *m = QueryBorrowRatioRequest{} }
func (m *QueryBorrowRatioRequest) String() string { return proto.CompactTextString(m) }
func (*QueryBorrowRatioRequest) ProtoMessage()    {}
func (*QueryBorrowRatioRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_9717e1e4014b4459, []int{2}
}
func (m *QueryBorrowRatioRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryBorrowRatioRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryBorrowRatioRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryBorrowRatioRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBorrowRatioRequest.Merge(m, src)
}
func (m *QueryBorrowRatioRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryBorrowRatioRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBorrowRatioRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBorrowRatioRequest proto.InternalMessageInfo

// QueryBorrowRatioResponse is response type for the Query/BorrowRatio RPC
// method.
type QueryBorrowRatioResponse struct {
	TotalDeposit cosmossdk_io_math.Int       `protobuf:"bytes,1,opt,name=total_deposit,json=totalDeposit,proto3,customtype=cosmossdk.io/math.Int" json:"total_deposit"`
	TotalBorrow  cosmossdk_io_math.Int       `protobuf:"bytes,2,opt,name=total_borrow,json=totalBorrow,proto3,customtype=cosmossdk.io/math.Int" json:"total_borrow"`
	BorrowRatio  cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=borrow_ratio,json=borrowRatio,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"borrow_ratio"`
}

func (m *QueryBorrowRatioResponse) Reset()         { *m = QueryBorrowRatioResponse{} }
func (m *QueryBorrowRatioResponse) String() string { return proto.CompactTextString(m) }
func (*QueryBorrowRatioResponse) ProtoMessage()    {}
func (*QueryBorrowRatioResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_9717e1e4014b4459, []int{3}
}
func (m *QueryBorrowRatioResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryBorrowRatioResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryBorrowRatioResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryBorrowRatioResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryBorrowRatioResponse.Merge(m, src)
}
func (m *QueryBorrowRatioResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryBorrowRatioResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryBorrowRatioResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryBorrowRatioResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*QueryParamsRequest)(nil), "elys.stablestake.QueryParamsRequest")
	proto.RegisterType((*QueryParamsResponse)(nil), "elys.stablestake.QueryParamsResponse")
	proto.RegisterType((*QueryBorrowRatioRequest)(nil), "elys.stablestake.QueryBorrowRatioRequest")
	proto.RegisterType((*QueryBorrowRatioResponse)(nil), "elys.stablestake.QueryBorrowRatioResponse")
}

func init() { proto.RegisterFile("elys/stablestake/query.proto", fileDescriptor_9717e1e4014b4459) }

var fileDescriptor_9717e1e4014b4459 = []byte{
	// 463 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x41, 0x6b, 0xd4, 0x40,
	0x14, 0xc7, 0x33, 0xab, 0x2e, 0x38, 0xa9, 0x20, 0x63, 0xc5, 0x18, 0xdb, 0x54, 0xa2, 0x45, 0xad,
	0x24, 0x63, 0x5b, 0xf0, 0x2c, 0xa1, 0x08, 0x8a, 0x82, 0xe6, 0xe8, 0x65, 0x99, 0xa4, 0x43, 0x1a,
	0x36, 0xc9, 0x4b, 0x33, 0xb3, 0xd6, 0x5c, 0xc5, 0x0f, 0x20, 0x78, 0xf2, 0x2b, 0xf8, 0x49, 0x7a,
	0x2c, 0x78, 0x11, 0x0f, 0x45, 0x76, 0x3d, 0xf8, 0x31, 0x24, 0x33, 0x23, 0xb6, 0xc6, 0x65, 0xf7,
	0x96, 0xcc, 0xfb, 0xbf, 0xdf, 0xfb, 0xbf, 0xff, 0x0c, 0x5e, 0xe3, 0x45, 0x2b, 0xa8, 0x90, 0x2c,
	0x29, 0xb8, 0x90, 0x6c, 0xcc, 0xe9, 0xe1, 0x84, 0x37, 0x6d, 0x58, 0x37, 0x20, 0x81, 0x5c, 0xed,
	0xaa, 0xe1, 0x99, 0xaa, 0xbb, 0x9a, 0x41, 0x06, 0xaa, 0x48, 0xbb, 0x2f, 0xad, 0x73, 0xd7, 0x32,
	0x80, 0xac, 0xe0, 0x94, 0xd5, 0x39, 0x65, 0x55, 0x05, 0x92, 0xc9, 0x1c, 0x2a, 0x61, 0xaa, 0x5b,
	0x29, 0x88, 0x12, 0x04, 0x4d, 0x98, 0x30, 0x78, 0xfa, 0x76, 0x3b, 0xe1, 0x92, 0x6d, 0xd3, 0x9a,
	0x65, 0x79, 0xa5, 0xc4, 0x46, 0xbb, 0xde, 0xf3, 0x53, 0xb3, 0x86, 0x95, 0x06, 0xe5, 0xaf, 0x62,
	0xf2, 0xba, 0x03, 0xbc, 0x52, 0x87, 0x31, 0x3f, 0x9c, 0x70, 0x21, 0xfd, 0x97, 0xf8, 0xda, 0xb9,
	0x53, 0x51, 0x43, 0x25, 0x38, 0x79, 0x8c, 0x87, 0xba, 0xd9, 0x41, 0xb7, 0xd1, 0x7d, 0x7b, 0xc7,
	0x09, 0xff, 0x5d, 0x27, 0xd4, 0x1d, 0xd1, 0xc5, 0xe3, 0xd3, 0x0d, 0x2b, 0x36, 0x6a, 0xff, 0x26,
	0xbe, 0xa1, 0x70, 0x11, 0x34, 0x0d, 0x1c, 0xc5, 0x9d, 0xbd, 0x3f, 0x93, 0x7e, 0x21, 0xec, 0xf4,
	0x6b, 0x66, 0x5e, 0x84, 0xaf, 0x48, 0x90, 0xac, 0x18, 0xed, 0xf3, 0x1a, 0x44, 0x2e, 0xd5, 0xd8,
	0xcb, 0xd1, 0x7a, 0x07, 0xff, 0x7e, 0xba, 0x71, 0x5d, 0xc7, 0x20, 0xf6, 0xc7, 0x61, 0x0e, 0xb4,
	0x64, 0xf2, 0x20, 0x7c, 0x56, 0xc9, 0x78, 0x45, 0xf5, 0xec, 0xe9, 0x16, 0xf2, 0x04, 0xeb, 0xff,
	0x51, 0xa2, 0x06, 0x38, 0x83, 0x65, 0x10, 0xb6, 0x6a, 0xd1, 0x96, 0xc8, 0x53, 0xbc, 0xa2, 0x7b,
	0x47, 0x4d, 0xe7, 0xce, 0xb9, 0xa0, 0x08, 0x77, 0x0c, 0xe1, 0x56, 0x9f, 0xf0, 0x82, 0x67, 0x2c,
	0x6d, 0xf7, 0x78, 0x1a, 0xdb, 0xc9, 0xdf, 0xad, 0x76, 0xbe, 0x0c, 0xf0, 0x25, 0xb5, 0x2a, 0xf9,
	0x80, 0xf0, 0x50, 0x07, 0x45, 0xee, 0xf6, 0x23, 0xec, 0xdf, 0x87, 0xbb, 0xb9, 0x40, 0xa5, 0xf3,
	0xf2, 0x83, 0xf7, 0x5f, 0x7f, 0x7e, 0x1a, 0xdc, 0x23, 0x9b, 0xb4, 0x93, 0x07, 0x15, 0x97, 0x47,
	0xd0, 0x8c, 0xe9, 0x9c, 0x17, 0x40, 0x3e, 0x23, 0x6c, 0x9f, 0x89, 0x9d, 0x3c, 0x98, 0x33, 0xa5,
	0x7f, 0x6d, 0xee, 0xd6, 0x32, 0x52, 0xe3, 0x6a, 0x57, 0xb9, 0x0a, 0xc8, 0xc3, 0x05, 0xae, 0x74,
	0x56, 0x81, 0x0a, 0x39, 0x7a, 0x7e, 0x3c, 0xf5, 0xd0, 0xc9, 0xd4, 0x43, 0x3f, 0xa6, 0x1e, 0xfa,
	0x38, 0xf3, 0xac, 0x93, 0x99, 0x67, 0x7d, 0x9b, 0x79, 0xd6, 0x9b, 0x47, 0x59, 0x2e, 0x0f, 0x26,
	0x49, 0x98, 0x42, 0xf9, 0x1f, 0xe0, 0xbb, 0x73, 0x48, 0xd9, 0xd6, 0x5c, 0x24, 0x43, 0xf5, 0xd4,
	0x77, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0xdf, 0x81, 0x9a, 0x29, 0x9b, 0x03, 0x00, 0x00,
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
	// BorrowRatio queries the borrow ratio compared to total deposit
	BorrowRatio(ctx context.Context, in *QueryBorrowRatioRequest, opts ...grpc.CallOption) (*QueryBorrowRatioResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/elys.stablestake.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) BorrowRatio(ctx context.Context, in *QueryBorrowRatioRequest, opts ...grpc.CallOption) (*QueryBorrowRatioResponse, error) {
	out := new(QueryBorrowRatioResponse)
	err := c.cc.Invoke(ctx, "/elys.stablestake.Query/BorrowRatio", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Parameters queries the parameters of the module.
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	// BorrowRatio queries the borrow ratio compared to total deposit
	BorrowRatio(context.Context, *QueryBorrowRatioRequest) (*QueryBorrowRatioResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Params(ctx context.Context, req *QueryParamsRequest) (*QueryParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Params not implemented")
}
func (*UnimplementedQueryServer) BorrowRatio(ctx context.Context, req *QueryBorrowRatioRequest) (*QueryBorrowRatioResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BorrowRatio not implemented")
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
		FullMethod: "/elys.stablestake.Query/Params",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_BorrowRatio_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBorrowRatioRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).BorrowRatio(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.stablestake.Query/BorrowRatio",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).BorrowRatio(ctx, req.(*QueryBorrowRatioRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var Query_serviceDesc = _Query_serviceDesc
var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "elys.stablestake.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Params",
			Handler:    _Query_Params_Handler,
		},
		{
			MethodName: "BorrowRatio",
			Handler:    _Query_BorrowRatio_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/stablestake/query.proto",
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

func (m *QueryBorrowRatioRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryBorrowRatioRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryBorrowRatioRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *QueryBorrowRatioResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryBorrowRatioResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryBorrowRatioResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.BorrowRatio.Size()
		i -= size
		if _, err := m.BorrowRatio.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.TotalBorrow.Size()
		i -= size
		if _, err := m.TotalBorrow.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size := m.TotalDeposit.Size()
		i -= size
		if _, err := m.TotalDeposit.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
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

func (m *QueryBorrowRatioRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *QueryBorrowRatioResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.TotalDeposit.Size()
	n += 1 + l + sovQuery(uint64(l))
	l = m.TotalBorrow.Size()
	n += 1 + l + sovQuery(uint64(l))
	l = m.BorrowRatio.Size()
	n += 1 + l + sovQuery(uint64(l))
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
func (m *QueryBorrowRatioRequest) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryBorrowRatioRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryBorrowRatioRequest: illegal tag %d (wire type %d)", fieldNum, wire)
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
func (m *QueryBorrowRatioResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: QueryBorrowRatioResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryBorrowRatioResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalDeposit", wireType)
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
			if err := m.TotalDeposit.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalBorrow", wireType)
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
			if err := m.TotalBorrow.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BorrowRatio", wireType)
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
			if err := m.BorrowRatio.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
