// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package vaults

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
	// UpdateParams defines a (governance) operation for updating the module
	// parameters. The authority defaults to the x/gov module account.
	UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error)
	// Deposit defines a method for depositing tokens into a vault.
	Deposit(ctx context.Context, in *MsgDeposit, opts ...grpc.CallOption) (*MsgDepositResponse, error)
	// Withdraw defines a method for withdrawing tokens from a vault.
	Withdraw(ctx context.Context, in *MsgWithdraw, opts ...grpc.CallOption) (*MsgWithdrawResponse, error)
	// AddVault defines a method for creating a new vault.
	AddVault(ctx context.Context, in *MsgAddVault, opts ...grpc.CallOption) (*MsgAddVaultResponse, error)
	// PerformAction defines a method for performing an action on a vault.
	// rpc PerformAction(MsgPerformAction) returns (MsgPerformActionResponse);
	PerformActionJoinPool(ctx context.Context, in *MsgPerformActionJoinPool, opts ...grpc.CallOption) (*MsgPerformActionJoinPoolResponse, error)
	PerformActionExitPool(ctx context.Context, in *MsgPerformActionExitPool, opts ...grpc.CallOption) (*MsgPerformActionExitPoolResponse, error)
	PerformActionSwapByDenom(ctx context.Context, in *MsgPerformActionSwapByDenom, opts ...grpc.CallOption) (*MsgPerformActionSwapByDenomResponse, error)
	// UpdateVaultCoins defines a method for updating the coins of a vault.
	UpdateVaultCoins(ctx context.Context, in *MsgUpdateVaultCoins, opts ...grpc.CallOption) (*MsgUpdateVaultCoinsResponse, error)
	// UpdateVaultFees defines a method for updating the fees of a vault.
	UpdateVaultFees(ctx context.Context, in *MsgUpdateVaultFees, opts ...grpc.CallOption) (*MsgUpdateVaultFeesResponse, error)
	// UpdateVaultLockupPeriod defines a method for updating the lockup period of
	// a vault.
	UpdateVaultLockupPeriod(ctx context.Context, in *MsgUpdateVaultLockupPeriod, opts ...grpc.CallOption) (*MsgUpdateVaultLockupPeriodResponse, error)
	// UpdateVaultMaxAmountUsd defines a method for updating the max amount of a
	// vault.
	UpdateVaultMaxAmountUsd(ctx context.Context, in *MsgUpdateVaultMaxAmountUsd, opts ...grpc.CallOption) (*MsgUpdateVaultMaxAmountUsdResponse, error)
	// ClaimRewards defines a method for claiming rewards from a vault.
	ClaimRewards(ctx context.Context, in *MsgClaimRewards, opts ...grpc.CallOption) (*MsgClaimRewardsResponse, error)
	// UpdateVaultAllowedActions defines a method for updating the allowed actions
	// of a vault.
	UpdateVaultAllowedActions(ctx context.Context, in *MsgUpdateVaultAllowedActions, opts ...grpc.CallOption) (*MsgUpdateVaultAllowedActionsResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) UpdateParams(ctx context.Context, in *MsgUpdateParams, opts ...grpc.CallOption) (*MsgUpdateParamsResponse, error) {
	out := new(MsgUpdateParamsResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/UpdateParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Deposit(ctx context.Context, in *MsgDeposit, opts ...grpc.CallOption) (*MsgDepositResponse, error) {
	out := new(MsgDepositResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/Deposit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Withdraw(ctx context.Context, in *MsgWithdraw, opts ...grpc.CallOption) (*MsgWithdrawResponse, error) {
	out := new(MsgWithdrawResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/Withdraw", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) AddVault(ctx context.Context, in *MsgAddVault, opts ...grpc.CallOption) (*MsgAddVaultResponse, error) {
	out := new(MsgAddVaultResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/AddVault", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) PerformActionJoinPool(ctx context.Context, in *MsgPerformActionJoinPool, opts ...grpc.CallOption) (*MsgPerformActionJoinPoolResponse, error) {
	out := new(MsgPerformActionJoinPoolResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/PerformActionJoinPool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) PerformActionExitPool(ctx context.Context, in *MsgPerformActionExitPool, opts ...grpc.CallOption) (*MsgPerformActionExitPoolResponse, error) {
	out := new(MsgPerformActionExitPoolResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/PerformActionExitPool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) PerformActionSwapByDenom(ctx context.Context, in *MsgPerformActionSwapByDenom, opts ...grpc.CallOption) (*MsgPerformActionSwapByDenomResponse, error) {
	out := new(MsgPerformActionSwapByDenomResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/PerformActionSwapByDenom", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateVaultCoins(ctx context.Context, in *MsgUpdateVaultCoins, opts ...grpc.CallOption) (*MsgUpdateVaultCoinsResponse, error) {
	out := new(MsgUpdateVaultCoinsResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/UpdateVaultCoins", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateVaultFees(ctx context.Context, in *MsgUpdateVaultFees, opts ...grpc.CallOption) (*MsgUpdateVaultFeesResponse, error) {
	out := new(MsgUpdateVaultFeesResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/UpdateVaultFees", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateVaultLockupPeriod(ctx context.Context, in *MsgUpdateVaultLockupPeriod, opts ...grpc.CallOption) (*MsgUpdateVaultLockupPeriodResponse, error) {
	out := new(MsgUpdateVaultLockupPeriodResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/UpdateVaultLockupPeriod", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateVaultMaxAmountUsd(ctx context.Context, in *MsgUpdateVaultMaxAmountUsd, opts ...grpc.CallOption) (*MsgUpdateVaultMaxAmountUsdResponse, error) {
	out := new(MsgUpdateVaultMaxAmountUsdResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/UpdateVaultMaxAmountUsd", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ClaimRewards(ctx context.Context, in *MsgClaimRewards, opts ...grpc.CallOption) (*MsgClaimRewardsResponse, error) {
	out := new(MsgClaimRewardsResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/ClaimRewards", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateVaultAllowedActions(ctx context.Context, in *MsgUpdateVaultAllowedActions, opts ...grpc.CallOption) (*MsgUpdateVaultAllowedActionsResponse, error) {
	out := new(MsgUpdateVaultAllowedActionsResponse)
	err := c.cc.Invoke(ctx, "/elys.vaults.Msg/UpdateVaultAllowedActions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// UpdateParams defines a (governance) operation for updating the module
	// parameters. The authority defaults to the x/gov module account.
	UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error)
	// Deposit defines a method for depositing tokens into a vault.
	Deposit(context.Context, *MsgDeposit) (*MsgDepositResponse, error)
	// Withdraw defines a method for withdrawing tokens from a vault.
	Withdraw(context.Context, *MsgWithdraw) (*MsgWithdrawResponse, error)
	// AddVault defines a method for creating a new vault.
	AddVault(context.Context, *MsgAddVault) (*MsgAddVaultResponse, error)
	// PerformAction defines a method for performing an action on a vault.
	// rpc PerformAction(MsgPerformAction) returns (MsgPerformActionResponse);
	PerformActionJoinPool(context.Context, *MsgPerformActionJoinPool) (*MsgPerformActionJoinPoolResponse, error)
	PerformActionExitPool(context.Context, *MsgPerformActionExitPool) (*MsgPerformActionExitPoolResponse, error)
	PerformActionSwapByDenom(context.Context, *MsgPerformActionSwapByDenom) (*MsgPerformActionSwapByDenomResponse, error)
	// UpdateVaultCoins defines a method for updating the coins of a vault.
	UpdateVaultCoins(context.Context, *MsgUpdateVaultCoins) (*MsgUpdateVaultCoinsResponse, error)
	// UpdateVaultFees defines a method for updating the fees of a vault.
	UpdateVaultFees(context.Context, *MsgUpdateVaultFees) (*MsgUpdateVaultFeesResponse, error)
	// UpdateVaultLockupPeriod defines a method for updating the lockup period of
	// a vault.
	UpdateVaultLockupPeriod(context.Context, *MsgUpdateVaultLockupPeriod) (*MsgUpdateVaultLockupPeriodResponse, error)
	// UpdateVaultMaxAmountUsd defines a method for updating the max amount of a
	// vault.
	UpdateVaultMaxAmountUsd(context.Context, *MsgUpdateVaultMaxAmountUsd) (*MsgUpdateVaultMaxAmountUsdResponse, error)
	// ClaimRewards defines a method for claiming rewards from a vault.
	ClaimRewards(context.Context, *MsgClaimRewards) (*MsgClaimRewardsResponse, error)
	// UpdateVaultAllowedActions defines a method for updating the allowed actions
	// of a vault.
	UpdateVaultAllowedActions(context.Context, *MsgUpdateVaultAllowedActions) (*MsgUpdateVaultAllowedActionsResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) UpdateParams(context.Context, *MsgUpdateParams) (*MsgUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateParams not implemented")
}
func (UnimplementedMsgServer) Deposit(context.Context, *MsgDeposit) (*MsgDepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposit not implemented")
}
func (UnimplementedMsgServer) Withdraw(context.Context, *MsgWithdraw) (*MsgWithdrawResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdraw not implemented")
}
func (UnimplementedMsgServer) AddVault(context.Context, *MsgAddVault) (*MsgAddVaultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddVault not implemented")
}
func (UnimplementedMsgServer) PerformActionJoinPool(context.Context, *MsgPerformActionJoinPool) (*MsgPerformActionJoinPoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PerformActionJoinPool not implemented")
}
func (UnimplementedMsgServer) PerformActionExitPool(context.Context, *MsgPerformActionExitPool) (*MsgPerformActionExitPoolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PerformActionExitPool not implemented")
}
func (UnimplementedMsgServer) PerformActionSwapByDenom(context.Context, *MsgPerformActionSwapByDenom) (*MsgPerformActionSwapByDenomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PerformActionSwapByDenom not implemented")
}
func (UnimplementedMsgServer) UpdateVaultCoins(context.Context, *MsgUpdateVaultCoins) (*MsgUpdateVaultCoinsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVaultCoins not implemented")
}
func (UnimplementedMsgServer) UpdateVaultFees(context.Context, *MsgUpdateVaultFees) (*MsgUpdateVaultFeesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVaultFees not implemented")
}
func (UnimplementedMsgServer) UpdateVaultLockupPeriod(context.Context, *MsgUpdateVaultLockupPeriod) (*MsgUpdateVaultLockupPeriodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVaultLockupPeriod not implemented")
}
func (UnimplementedMsgServer) UpdateVaultMaxAmountUsd(context.Context, *MsgUpdateVaultMaxAmountUsd) (*MsgUpdateVaultMaxAmountUsdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVaultMaxAmountUsd not implemented")
}
func (UnimplementedMsgServer) ClaimRewards(context.Context, *MsgClaimRewards) (*MsgClaimRewardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClaimRewards not implemented")
}
func (UnimplementedMsgServer) UpdateVaultAllowedActions(context.Context, *MsgUpdateVaultAllowedActions) (*MsgUpdateVaultAllowedActionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVaultAllowedActions not implemented")
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
		FullMethod: "/elys.vaults.Msg/UpdateParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateParams(ctx, req.(*MsgUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Deposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgDeposit)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.vaults.Msg/Deposit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Deposit(ctx, req.(*MsgDeposit))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Withdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdraw)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Withdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.vaults.Msg/Withdraw",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Withdraw(ctx, req.(*MsgWithdraw))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddVault_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddVault)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddVault(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.vaults.Msg/AddVault",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddVault(ctx, req.(*MsgAddVault))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_PerformActionJoinPool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgPerformActionJoinPool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).PerformActionJoinPool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.vaults.Msg/PerformActionJoinPool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).PerformActionJoinPool(ctx, req.(*MsgPerformActionJoinPool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_PerformActionExitPool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgPerformActionExitPool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).PerformActionExitPool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.vaults.Msg/PerformActionExitPool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).PerformActionExitPool(ctx, req.(*MsgPerformActionExitPool))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_PerformActionSwapByDenom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgPerformActionSwapByDenom)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).PerformActionSwapByDenom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.vaults.Msg/PerformActionSwapByDenom",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).PerformActionSwapByDenom(ctx, req.(*MsgPerformActionSwapByDenom))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateVaultCoins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateVaultCoins)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateVaultCoins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.vaults.Msg/UpdateVaultCoins",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateVaultCoins(ctx, req.(*MsgUpdateVaultCoins))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateVaultFees_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateVaultFees)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateVaultFees(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.vaults.Msg/UpdateVaultFees",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateVaultFees(ctx, req.(*MsgUpdateVaultFees))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateVaultLockupPeriod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateVaultLockupPeriod)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateVaultLockupPeriod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.vaults.Msg/UpdateVaultLockupPeriod",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateVaultLockupPeriod(ctx, req.(*MsgUpdateVaultLockupPeriod))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateVaultMaxAmountUsd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateVaultMaxAmountUsd)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateVaultMaxAmountUsd(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.vaults.Msg/UpdateVaultMaxAmountUsd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateVaultMaxAmountUsd(ctx, req.(*MsgUpdateVaultMaxAmountUsd))
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
		FullMethod: "/elys.vaults.Msg/ClaimRewards",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimRewards(ctx, req.(*MsgClaimRewards))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateVaultAllowedActions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateVaultAllowedActions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateVaultAllowedActions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.vaults.Msg/UpdateVaultAllowedActions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateVaultAllowedActions(ctx, req.(*MsgUpdateVaultAllowedActions))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "elys.vaults.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateParams",
			Handler:    _Msg_UpdateParams_Handler,
		},
		{
			MethodName: "Deposit",
			Handler:    _Msg_Deposit_Handler,
		},
		{
			MethodName: "Withdraw",
			Handler:    _Msg_Withdraw_Handler,
		},
		{
			MethodName: "AddVault",
			Handler:    _Msg_AddVault_Handler,
		},
		{
			MethodName: "PerformActionJoinPool",
			Handler:    _Msg_PerformActionJoinPool_Handler,
		},
		{
			MethodName: "PerformActionExitPool",
			Handler:    _Msg_PerformActionExitPool_Handler,
		},
		{
			MethodName: "PerformActionSwapByDenom",
			Handler:    _Msg_PerformActionSwapByDenom_Handler,
		},
		{
			MethodName: "UpdateVaultCoins",
			Handler:    _Msg_UpdateVaultCoins_Handler,
		},
		{
			MethodName: "UpdateVaultFees",
			Handler:    _Msg_UpdateVaultFees_Handler,
		},
		{
			MethodName: "UpdateVaultLockupPeriod",
			Handler:    _Msg_UpdateVaultLockupPeriod_Handler,
		},
		{
			MethodName: "UpdateVaultMaxAmountUsd",
			Handler:    _Msg_UpdateVaultMaxAmountUsd_Handler,
		},
		{
			MethodName: "ClaimRewards",
			Handler:    _Msg_ClaimRewards_Handler,
		},
		{
			MethodName: "UpdateVaultAllowedActions",
			Handler:    _Msg_UpdateVaultAllowedActions_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/vaults/tx.proto",
}
