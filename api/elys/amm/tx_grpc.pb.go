// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package amm

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MsgClient interface {
	CreatePool(ctx context.Context, in *MsgCreatePool, opts ...grpc.CallOption) (*MsgCreatePoolResponse, error)
	JoinPool(ctx context.Context, in *MsgJoinPool, opts ...grpc.CallOption) (*MsgJoinPoolResponse, error)
	ExitPool(ctx context.Context, in *MsgExitPool, opts ...grpc.CallOption) (*MsgExitPoolResponse, error)
	SwapExactAmountIn(ctx context.Context, in *MsgSwapExactAmountIn, opts ...grpc.CallOption) (*MsgSwapExactAmountInResponse, error)
	SwapExactAmountOut(ctx context.Context, in *MsgSwapExactAmountOut, opts ...grpc.CallOption) (*MsgSwapExactAmountOutResponse, error)
	SwapByDenom(ctx context.Context, in *MsgSwapByDenom, opts ...grpc.CallOption) (*MsgSwapByDenomResponse, error)
	FeedMultipleExternalLiquidity(ctx context.Context, in *MsgFeedMultipleExternalLiquidity, opts ...grpc.CallOption) (*MsgFeedMultipleExternalLiquidityResponse, error)
	UpdatePoolParams(ctx context.Context, in *MsgUpdatePoolParams, opts ...grpc.CallOption) (*MsgUpdatePoolParamsResponse, error)
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CreatePool(ctx context.Context, in *MsgCreatePool, opts ...grpc.CallOption) (*MsgCreatePoolResponse, error) {
	out := new(MsgCreatePoolResponse)
	err := c.cc.Invoke(ctx, "/elys.amm.Msg/CreatePool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) JoinPool(ctx context.Context, in *MsgJoinPool, opts ...grpc.CallOption) (*MsgJoinPoolResponse, error) {
	out := new(MsgJoinPoolResponse)
	err := c.cc.Invoke(ctx, "/elys.amm.Msg/JoinPool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ExitPool(ctx context.Context, in *MsgExitPool, opts ...grpc.CallOption) (*MsgExitPoolResponse, error) {
	out := new(MsgExitPoolResponse)
	err := c.cc.Invoke(ctx, "/elys.amm.Msg/ExitPool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SwapExactAmountIn(ctx context.Context, in *MsgSwapExactAmountIn, opts ...grpc.CallOption) (*MsgSwapExactAmountInResponse, error) {
	out := new(MsgSwapExactAmountInResponse)
	err := c.cc.Invoke(ctx, "/elys.amm.Msg/SwapExactAmountIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SwapExactAmountOut(ctx context.Context, in *MsgSwapExactAmountOut, opts ...grpc.CallOption) (*MsgSwapExactAmountOutResponse, error) {
	out := new(MsgSwapExactAmountOutResponse)
	err := c.cc.Invoke(ctx, "/elys.amm.Msg/SwapExactAmountOut", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SwapByDenom(ctx context.Context, in *MsgSwapByDenom, opts ...grpc.CallOption) (*MsgSwapByDenomResponse, error) {
	out := new(MsgSwapByDenomResponse)
	err := c.cc.Invoke(ctx, "/elys.amm.Msg/SwapByDenom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) FeedMultipleExternalLiquidity(ctx context.Context, in *MsgFeedMultipleExternalLiquidity, opts ...grpc.CallOption) (*MsgFeedMultipleExternalLiquidityResponse, error) {
	out := new(MsgFeedMultipleExternalLiquidityResponse)
	err := c.cc.Invoke(ctx, "/elys.amm.Msg/FeedMultipleExternalLiquidity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdatePoolParams(ctx context.Context, in *MsgUpdatePoolParams, opts ...grpc.CallOption) (*MsgUpdatePoolParamsResponse, error) {
	out := new(MsgUpdatePoolParamsResponse)
	err := c.cc.Invoke(ctx, "/elys.amm.Msg/UpdatePoolParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, "/elys.amm.Msg/UpdateParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	CreatePool(context.Context, *MsgCreatePool) (*MsgCreatePoolResponse, error)
	JoinPool(context.Context, *MsgJoinPool) (*MsgJoinPoolResponse, error)
	ExitPool(context.Context, *MsgExitPool) (*MsgExitPoolResponse, error)
	SwapExactAmountIn(context.Context, *MsgSwapExactAmountIn) (*MsgSwapExactAmountInResponse, error)
	SwapExactAmountOut(context.Context, *MsgSwapExactAmountOut) (*MsgSwapExactAmountOutResponse, error)
	SwapByDenom(context.Context, *MsgSwapByDenom) (*MsgSwapByDenomResponse, error)
	FeedMultipleExternalLiquidity(context.Context, *MsgFeedMultipleExternalLiquidity) (*MsgFeedMultipleExternalLiquidityResponse, error)
	UpdatePoolParams(context.Context, *MsgUpdatePoolParams) (*MsgUpdatePoolParamsResponse, error)
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) CreatePool(context.Context, *MsgCreatePool) (*MsgCreatePoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePool not implemented")
}
func (UnimplementedMsgServer) JoinPool(context.Context, *MsgJoinPool) (*MsgJoinPoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinPool not implemented")
}
func (UnimplementedMsgServer) ExitPool(context.Context, *MsgExitPool) (*MsgExitPoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExitPool not implemented")
}
func (UnimplementedMsgServer) SwapExactAmountIn(context.Context, *MsgSwapExactAmountIn) (*MsgSwapExactAmountInResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SwapExactAmountIn not implemented")
}
func (UnimplementedMsgServer) SwapExactAmountOut(context.Context, *MsgSwapExactAmountOut) (*MsgSwapExactAmountOutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SwapExactAmountOut not implemented")
}
func (UnimplementedMsgServer) SwapByDenom(context.Context, *MsgSwapByDenom) (*MsgSwapByDenomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SwapByDenom not implemented")
}
func (UnimplementedMsgServer) FeedMultipleExternalLiquidity(context.Context, *MsgFeedMultipleExternalLiquidity) (*MsgFeedMultipleExternalLiquidityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FeedMultipleExternalLiquidity not implemented")
}
func (UnimplementedMsgServer) UpdatePoolParams(context.Context, *MsgUpdatePoolParams) (*MsgUpdatePoolParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePoolParams not implemented")
}
func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) mustEmbedUnimplementedMsgServer() {}

// UnsafeMsgServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MsgServer will
// result in compilation errors.
type UnsafeMsgServer interface {
	mustEmbedUnimplementedMsgServer()
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&Msg_ServiceDesc, srv)
}

func _Msg_CreatePool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreatePool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreatePool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.amm.Msg/CreatePool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreatePool(ctx, req.(*MsgCreatePool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_JoinPool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgJoinPool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).JoinPool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.amm.Msg/JoinPool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).JoinPool(ctx, req.(*MsgJoinPool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ExitPool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgExitPool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ExitPool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.amm.Msg/ExitPool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ExitPool(ctx, req.(*MsgExitPool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SwapExactAmountIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSwapExactAmountIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SwapExactAmountIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.amm.Msg/SwapExactAmountIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SwapExactAmountIn(ctx, req.(*MsgSwapExactAmountIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SwapExactAmountOut_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSwapExactAmountOut)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SwapExactAmountOut(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.amm.Msg/SwapExactAmountOut",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SwapExactAmountOut(ctx, req.(*MsgSwapExactAmountOut))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SwapByDenom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSwapByDenom)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SwapByDenom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.amm.Msg/SwapByDenom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SwapByDenom(ctx, req.(*MsgSwapByDenom))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_FeedMultipleExternalLiquidity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFeedMultipleExternalLiquidity)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).FeedMultipleExternalLiquidity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.amm.Msg/FeedMultipleExternalLiquidity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).FeedMultipleExternalLiquidity(ctx, req.(*MsgFeedMultipleExternalLiquidity))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdatePoolParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdatePoolParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdatePoolParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.amm.Msg/UpdatePoolParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdatePoolParams(ctx, req.(*MsgUpdatePoolParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.amm.Msg/UpdateParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "elys.amm.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePool",
			Handler:    _Msg_CreatePool_Handler,
		},
		{
			MethodName: "JoinPool",
			Handler:    _Msg_JoinPool_Handler,
		},
		{
			MethodName: "ExitPool",
			Handler:    _Msg_ExitPool_Handler,
		},
		{
			MethodName: "SwapExactAmountIn",
			Handler:    _Msg_SwapExactAmountIn_Handler,
		},
		{
			MethodName: "SwapExactAmountOut",
			Handler:    _Msg_SwapExactAmountOut_Handler,
		},
		{
			MethodName: "SwapByDenom",
			Handler:    _Msg_SwapByDenom_Handler,
		},
		{
			MethodName: "FeedMultipleExternalLiquidity",
			Handler:    _Msg_FeedMultipleExternalLiquidity_Handler,
		},
		{
			MethodName: "UpdatePoolParams",
			Handler:    _Msg_UpdatePoolParams_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/amm/tx.proto",
}
