// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package stablestake

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
	Bond(ctx context.Context, in *MsgBond, opts ...grpc.CallOption) (*MsgBondResponse, error)
	Unbond(ctx context.Context, in *MsgUnbond, opts ...grpc.CallOption) (*MsgUnbondResponse, error)
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
	AddPool(ctx context.Context, in *MsgAddPool, opts ...grpc.CallOption) (*MsgAddPoolResponse, error)
	UpdatePool(ctx context.Context, in *MsgUpdatePool, opts ...grpc.CallOption) (*MsgUpdatePoolResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) Bond(ctx context.Context, in *MsgBond, opts ...grpc.CallOption) (*MsgBondResponse, error) {
	out := new(MsgBondResponse)
	err := c.cc.Invoke(ctx, "/elys.stablestake.Msg/Bond", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Unbond(ctx context.Context, in *MsgUnbond, opts ...grpc.CallOption) (*MsgUnbondResponse, error) {
	out := new(MsgUnbondResponse)
	err := c.cc.Invoke(ctx, "/elys.stablestake.Msg/Unbond", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, "/elys.stablestake.Msg/UpdateParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) AddPool(ctx context.Context, in *MsgAddPool, opts ...grpc.CallOption) (*MsgAddPoolResponse, error) {
	out := new(MsgAddPoolResponse)
	err := c.cc.Invoke(ctx, "/elys.stablestake.Msg/AddPool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdatePool(ctx context.Context, in *MsgUpdatePool, opts ...grpc.CallOption) (*MsgUpdatePoolResponse, error) {
	out := new(MsgUpdatePoolResponse)
	err := c.cc.Invoke(ctx, "/elys.stablestake.Msg/UpdatePool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	Bond(context.Context, *MsgBond) (*MsgBondResponse, error)
	Unbond(context.Context, *MsgUnbond) (*MsgUnbondResponse, error)
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	AddPool(context.Context, *MsgAddPool) (*MsgAddPoolResponse, error)
	UpdatePool(context.Context, *MsgUpdatePool) (*MsgUpdatePoolResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) Bond(context.Context, *MsgBond) (*MsgBondResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bond not implemented")
}
func (UnimplementedMsgServer) Unbond(context.Context, *MsgUnbond) (*MsgUnbondResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unbond not implemented")
}
func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) AddPool(context.Context, *MsgAddPool) (*MsgAddPoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPool not implemented")
}
func (UnimplementedMsgServer) UpdatePool(context.Context, *MsgUpdatePool) (*MsgUpdatePoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePool not implemented")
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

func _Msg_Bond_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgBond)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Bond(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.stablestake.Msg/Bond",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Bond(ctx, req.(*MsgBond))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Unbond_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUnbond)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Unbond(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.stablestake.Msg/Unbond",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Unbond(ctx, req.(*MsgUnbond))
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
		FullMethod: "/elys.stablestake.Msg/UpdateParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddPool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddPool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddPool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.stablestake.Msg/AddPool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddPool(ctx, req.(*MsgAddPool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdatePool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdatePool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdatePool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.stablestake.Msg/UpdatePool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdatePool(ctx, req.(*MsgUpdatePool))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "elys.stablestake.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Bond",
			Handler:    _Msg_Bond_Handler,
		},
		{
			MethodName: "Unbond",
			Handler:    _Msg_Unbond_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
		{
			MethodName: "AddPool",
			Handler:    _Msg_AddPool_Handler,
		},
		{
			MethodName: "UpdatePool",
			Handler:    _Msg_UpdatePool_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/stablestake/tx.proto",
}
