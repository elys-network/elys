// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package leveragelp

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
	Open(ctx context.Context, in *MsgOpen, opts ...grpc.CallOption) (*MsgOpenResponse, error)
	Close(ctx context.Context, in *MsgClose, opts ...grpc.CallOption) (*MsgCloseResponse, error)
	ClaimRewards(ctx context.Context, in *MsgClaimRewards, opts ...grpc.CallOption) (*MsgClaimRewardsResponse, error)
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
	AddPool(ctx context.Context, in *MsgAddPool, opts ...grpc.CallOption) (*MsgAddPoolResponse, error)
	RemovePool(ctx context.Context, in *MsgRemovePool, opts ...grpc.CallOption) (*MsgRemovePoolResponse, error)
	Whitelist(ctx context.Context, in *MsgWhitelist, opts ...grpc.CallOption) (*MsgWhitelistResponse, error)
	Dewhitelist(ctx context.Context, in *MsgDewhitelist, opts ...grpc.CallOption) (*MsgDewhitelistResponse, error)
	UpdateStopLoss(ctx context.Context, in *MsgUpdateStopLoss, opts ...grpc.CallOption) (*MsgUpdateStopLossResponse, error)
	ClosePositions(ctx context.Context, in *MsgClosePositions, opts ...grpc.CallOption) (*MsgClosePositionsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) Open(ctx context.Context, in *MsgOpen, opts ...grpc.CallOption) (*MsgOpenResponse, error) {
	out := new(MsgOpenResponse)
	err := c.cc.Invoke(ctx, "/elys.leveragelp.Msg/Open", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Close(ctx context.Context, in *MsgClose, opts ...grpc.CallOption) (*MsgCloseResponse, error) {
	out := new(MsgCloseResponse)
	err := c.cc.Invoke(ctx, "/elys.leveragelp.Msg/Close", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ClaimRewards(ctx context.Context, in *MsgClaimRewards, opts ...grpc.CallOption) (*MsgClaimRewardsResponse, error) {
	out := new(MsgClaimRewardsResponse)
	err := c.cc.Invoke(ctx, "/elys.leveragelp.Msg/ClaimRewards", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, "/elys.leveragelp.Msg/UpdateParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) AddPool(ctx context.Context, in *MsgAddPool, opts ...grpc.CallOption) (*MsgAddPoolResponse, error) {
	out := new(MsgAddPoolResponse)
	err := c.cc.Invoke(ctx, "/elys.leveragelp.Msg/AddPool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RemovePool(ctx context.Context, in *MsgRemovePool, opts ...grpc.CallOption) (*MsgRemovePoolResponse, error) {
	out := new(MsgRemovePoolResponse)
	err := c.cc.Invoke(ctx, "/elys.leveragelp.Msg/RemovePool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Whitelist(ctx context.Context, in *MsgWhitelist, opts ...grpc.CallOption) (*MsgWhitelistResponse, error) {
	out := new(MsgWhitelistResponse)
	err := c.cc.Invoke(ctx, "/elys.leveragelp.Msg/Whitelist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Dewhitelist(ctx context.Context, in *MsgDewhitelist, opts ...grpc.CallOption) (*MsgDewhitelistResponse, error) {
	out := new(MsgDewhitelistResponse)
	err := c.cc.Invoke(ctx, "/elys.leveragelp.Msg/Dewhitelist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateStopLoss(ctx context.Context, in *MsgUpdateStopLoss, opts ...grpc.CallOption) (*MsgUpdateStopLossResponse, error) {
	out := new(MsgUpdateStopLossResponse)
	err := c.cc.Invoke(ctx, "/elys.leveragelp.Msg/UpdateStopLoss", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ClosePositions(ctx context.Context, in *MsgClosePositions, opts ...grpc.CallOption) (*MsgClosePositionsResponse, error) {
	out := new(MsgClosePositionsResponse)
	err := c.cc.Invoke(ctx, "/elys.leveragelp.Msg/ClosePositions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	Open(context.Context, *MsgOpen) (*MsgOpenResponse, error)
	Close(context.Context, *MsgClose) (*MsgCloseResponse, error)
	ClaimRewards(context.Context, *MsgClaimRewards) (*MsgClaimRewardsResponse, error)
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	AddPool(context.Context, *MsgAddPool) (*MsgAddPoolResponse, error)
	RemovePool(context.Context, *MsgRemovePool) (*MsgRemovePoolResponse, error)
	Whitelist(context.Context, *MsgWhitelist) (*MsgWhitelistResponse, error)
	Dewhitelist(context.Context, *MsgDewhitelist) (*MsgDewhitelistResponse, error)
	UpdateStopLoss(context.Context, *MsgUpdateStopLoss) (*MsgUpdateStopLossResponse, error)
	ClosePositions(context.Context, *MsgClosePositions) (*MsgClosePositionsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) Open(context.Context, *MsgOpen) (*MsgOpenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Open not implemented")
}
func (UnimplementedMsgServer) Close(context.Context, *MsgClose) (*MsgCloseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}
func (UnimplementedMsgServer) ClaimRewards(context.Context, *MsgClaimRewards) (*MsgClaimRewardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClaimRewards not implemented")
}
func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) AddPool(context.Context, *MsgAddPool) (*MsgAddPoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPool not implemented")
}
func (UnimplementedMsgServer) RemovePool(context.Context, *MsgRemovePool) (*MsgRemovePoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemovePool not implemented")
}
func (UnimplementedMsgServer) Whitelist(context.Context, *MsgWhitelist) (*MsgWhitelistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Whitelist not implemented")
}
func (UnimplementedMsgServer) Dewhitelist(context.Context, *MsgDewhitelist) (*MsgDewhitelistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Dewhitelist not implemented")
}
func (UnimplementedMsgServer) UpdateStopLoss(context.Context, *MsgUpdateStopLoss) (*MsgUpdateStopLossResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStopLoss not implemented")
}
func (UnimplementedMsgServer) ClosePositions(context.Context, *MsgClosePositions) (*MsgClosePositionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClosePositions not implemented")
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

func _Msg_Open_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgOpen)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Open(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.leveragelp.Msg/Open",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Open(ctx, req.(*MsgOpen))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClose)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.leveragelp.Msg/Close",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Close(ctx, req.(*MsgClose))
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
		FullMethod: "/elys.leveragelp.Msg/ClaimRewards",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimRewards(ctx, req.(*MsgClaimRewards))
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
		FullMethod: "/elys.leveragelp.Msg/UpdateParams",
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
		FullMethod: "/elys.leveragelp.Msg/AddPool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddPool(ctx, req.(*MsgAddPool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemovePool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemovePool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RemovePool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.leveragelp.Msg/RemovePool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RemovePool(ctx, req.(*MsgRemovePool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Whitelist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWhitelist)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Whitelist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.leveragelp.Msg/Whitelist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Whitelist(ctx, req.(*MsgWhitelist))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Dewhitelist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDewhitelist)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Dewhitelist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.leveragelp.Msg/Dewhitelist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Dewhitelist(ctx, req.(*MsgDewhitelist))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateStopLoss_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateStopLoss)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateStopLoss(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.leveragelp.Msg/UpdateStopLoss",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateStopLoss(ctx, req.(*MsgUpdateStopLoss))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClosePositions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClosePositions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClosePositions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.leveragelp.Msg/ClosePositions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClosePositions(ctx, req.(*MsgClosePositions))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "elys.leveragelp.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Open",
			Handler:    _Msg_Open_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _Msg_Close_Handler,
		},
		{
			MethodName: "ClaimRewards",
			Handler:    _Msg_ClaimRewards_Handler,
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
			MethodName: "RemovePool",
			Handler:    _Msg_RemovePool_Handler,
		},
		{
			MethodName: "Whitelist",
			Handler:    _Msg_Whitelist_Handler,
		},
		{
			MethodName: "Dewhitelist",
			Handler:    _Msg_Dewhitelist_Handler,
		},
		{
			MethodName: "UpdateStopLoss",
			Handler:    _Msg_UpdateStopLoss_Handler,
		},
		{
			MethodName: "ClosePositions",
			Handler:    _Msg_ClosePositions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/leveragelp/tx.proto",
}
