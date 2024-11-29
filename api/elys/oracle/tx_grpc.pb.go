// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package oracle

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
	FeedPrice(ctx context.Context, in *MsgFeedPrice, opts ...grpc.CallOption) (*MsgFeedPriceResponse, error)
	FeedMultiplePrices(ctx context.Context, in *MsgFeedMultiplePrices, opts ...grpc.CallOption) (*MsgFeedMultiplePricesResponse, error)
	SetPriceFeeder(ctx context.Context, in *MsgSetPriceFeeder, opts ...grpc.CallOption) (*MsgSetPriceFeederResponse, error)
	DeletePriceFeeder(ctx context.Context, in *MsgDeletePriceFeeder, opts ...grpc.CallOption) (*MsgDeletePriceFeederResponse, error)
	// proposals
	RemoveAssetInfo(ctx context.Context, in *MsgRemoveAssetInfo, opts ...grpc.CallOption) (*MsgRemoveAssetInfoResponse, error)
	AddPriceFeeders(ctx context.Context, in *MsgAddPriceFeeders, opts ...grpc.CallOption) (*MsgAddPriceFeedersResponse, error)
	RemovePriceFeeders(ctx context.Context, in *MsgRemovePriceFeeders, opts ...grpc.CallOption) (*MsgRemovePriceFeedersResponse, error)
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
	// this line is used by starport scaffolding # proto/tx/rpc
	CreateAssetInfo(ctx context.Context, in *MsgCreateAssetInfo, opts ...grpc.CallOption) (*MsgCreateAssetInfoResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) FeedPrice(ctx context.Context, in *MsgFeedPrice, opts ...grpc.CallOption) (*MsgFeedPriceResponse, error) {
	out := new(MsgFeedPriceResponse)
	err := c.cc.Invoke(ctx, "/elys.oracle.Msg/FeedPrice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) FeedMultiplePrices(ctx context.Context, in *MsgFeedMultiplePrices, opts ...grpc.CallOption) (*MsgFeedMultiplePricesResponse, error) {
	out := new(MsgFeedMultiplePricesResponse)
	err := c.cc.Invoke(ctx, "/elys.oracle.Msg/FeedMultiplePrices", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) SetPriceFeeder(ctx context.Context, in *MsgSetPriceFeeder, opts ...grpc.CallOption) (*MsgSetPriceFeederResponse, error) {
	out := new(MsgSetPriceFeederResponse)
	err := c.cc.Invoke(ctx, "/elys.oracle.Msg/SetPriceFeeder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) DeletePriceFeeder(ctx context.Context, in *MsgDeletePriceFeeder, opts ...grpc.CallOption) (*MsgDeletePriceFeederResponse, error) {
	out := new(MsgDeletePriceFeederResponse)
	err := c.cc.Invoke(ctx, "/elys.oracle.Msg/DeletePriceFeeder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RemoveAssetInfo(ctx context.Context, in *MsgRemoveAssetInfo, opts ...grpc.CallOption) (*MsgRemoveAssetInfoResponse, error) {
	out := new(MsgRemoveAssetInfoResponse)
	err := c.cc.Invoke(ctx, "/elys.oracle.Msg/RemoveAssetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) AddPriceFeeders(ctx context.Context, in *MsgAddPriceFeeders, opts ...grpc.CallOption) (*MsgAddPriceFeedersResponse, error) {
	out := new(MsgAddPriceFeedersResponse)
	err := c.cc.Invoke(ctx, "/elys.oracle.Msg/AddPriceFeeders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) RemovePriceFeeders(ctx context.Context, in *MsgRemovePriceFeeders, opts ...grpc.CallOption) (*MsgRemovePriceFeedersResponse, error) {
	out := new(MsgRemovePriceFeedersResponse)
	err := c.cc.Invoke(ctx, "/elys.oracle.Msg/RemovePriceFeeders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, "/elys.oracle.Msg/UpdateParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CreateAssetInfo(ctx context.Context, in *MsgCreateAssetInfo, opts ...grpc.CallOption) (*MsgCreateAssetInfoResponse, error) {
	out := new(MsgCreateAssetInfoResponse)
	err := c.cc.Invoke(ctx, "/elys.oracle.Msg/CreateAssetInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	FeedPrice(context.Context, *MsgFeedPrice) (*MsgFeedPriceResponse, error)
	FeedMultiplePrices(context.Context, *MsgFeedMultiplePrices) (*MsgFeedMultiplePricesResponse, error)
	SetPriceFeeder(context.Context, *MsgSetPriceFeeder) (*MsgSetPriceFeederResponse, error)
	DeletePriceFeeder(context.Context, *MsgDeletePriceFeeder) (*MsgDeletePriceFeederResponse, error)
	// proposals
	RemoveAssetInfo(context.Context, *MsgRemoveAssetInfo) (*MsgRemoveAssetInfoResponse, error)
	AddPriceFeeders(context.Context, *MsgAddPriceFeeders) (*MsgAddPriceFeedersResponse, error)
	RemovePriceFeeders(context.Context, *MsgRemovePriceFeeders) (*MsgRemovePriceFeedersResponse, error)
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	// this line is used by starport scaffolding # proto/tx/rpc
	CreateAssetInfo(context.Context, *MsgCreateAssetInfo) (*MsgCreateAssetInfoResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) FeedPrice(context.Context, *MsgFeedPrice) (*MsgFeedPriceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FeedPrice not implemented")
}
func (UnimplementedMsgServer) FeedMultiplePrices(context.Context, *MsgFeedMultiplePrices) (*MsgFeedMultiplePricesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FeedMultiplePrices not implemented")
}
func (UnimplementedMsgServer) SetPriceFeeder(context.Context, *MsgSetPriceFeeder) (*MsgSetPriceFeederResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPriceFeeder not implemented")
}
func (UnimplementedMsgServer) DeletePriceFeeder(context.Context, *MsgDeletePriceFeeder) (*MsgDeletePriceFeederResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePriceFeeder not implemented")
}
func (UnimplementedMsgServer) RemoveAssetInfo(context.Context, *MsgRemoveAssetInfo) (*MsgRemoveAssetInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveAssetInfo not implemented")
}
func (UnimplementedMsgServer) AddPriceFeeders(context.Context, *MsgAddPriceFeeders) (*MsgAddPriceFeedersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPriceFeeders not implemented")
}
func (UnimplementedMsgServer) RemovePriceFeeders(context.Context, *MsgRemovePriceFeeders) (*MsgRemovePriceFeedersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemovePriceFeeders not implemented")
}
func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) CreateAssetInfo(context.Context, *MsgCreateAssetInfo) (*MsgCreateAssetInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAssetInfo not implemented")
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

func _Msg_FeedPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFeedPrice)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).FeedPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.oracle.Msg/FeedPrice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).FeedPrice(ctx, req.(*MsgFeedPrice))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_FeedMultiplePrices_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgFeedMultiplePrices)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).FeedMultiplePrices(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.oracle.Msg/FeedMultiplePrices",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).FeedMultiplePrices(ctx, req.(*MsgFeedMultiplePrices))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_SetPriceFeeder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgSetPriceFeeder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).SetPriceFeeder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.oracle.Msg/SetPriceFeeder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).SetPriceFeeder(ctx, req.(*MsgSetPriceFeeder))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_DeletePriceFeeder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeletePriceFeeder)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).DeletePriceFeeder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.oracle.Msg/DeletePriceFeeder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).DeletePriceFeeder(ctx, req.(*MsgDeletePriceFeeder))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemoveAssetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveAssetInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RemoveAssetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.oracle.Msg/RemoveAssetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RemoveAssetInfo(ctx, req.(*MsgRemoveAssetInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddPriceFeeders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddPriceFeeders)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddPriceFeeders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.oracle.Msg/AddPriceFeeders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddPriceFeeders(ctx, req.(*MsgAddPriceFeeders))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemovePriceFeeders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemovePriceFeeders)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RemovePriceFeeders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.oracle.Msg/RemovePriceFeeders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RemovePriceFeeders(ctx, req.(*MsgRemovePriceFeeders))
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
		FullMethod: "/elys.oracle.Msg/UpdateParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CreateAssetInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCreateAssetInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CreateAssetInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.oracle.Msg/CreateAssetInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CreateAssetInfo(ctx, req.(*MsgCreateAssetInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "elys.oracle.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FeedPrice",
			Handler:    _Msg_FeedPrice_Handler,
		},
		{
			MethodName: "FeedMultiplePrices",
			Handler:    _Msg_FeedMultiplePrices_Handler,
		},
		{
			MethodName: "SetPriceFeeder",
			Handler:    _Msg_SetPriceFeeder_Handler,
		},
		{
			MethodName: "DeletePriceFeeder",
			Handler:    _Msg_DeletePriceFeeder_Handler,
		},
		{
			MethodName: "RemoveAssetInfo",
			Handler:    _Msg_RemoveAssetInfo_Handler,
		},
		{
			MethodName: "AddPriceFeeders",
			Handler:    _Msg_AddPriceFeeders_Handler,
		},
		{
			MethodName: "RemovePriceFeeders",
			Handler:    _Msg_RemovePriceFeeders_Handler,
		},
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
		{
			MethodName: "CreateAssetInfo",
			Handler:    _Msg_CreateAssetInfo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/oracle/tx.proto",
}
