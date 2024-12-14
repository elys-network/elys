// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package commitment

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
	// CommitClaimedRewards commit the tokens on claimed store to committed
	CommitClaimedRewards(ctx context.Context, in *MsgCommitClaimedRewards, opts ...grpc.CallOption) (*MsgCommitClaimedRewardsResponse, error)
	// UncommitTokens uncommits the tokens from committed store and make it liquid
	// immediately
	UncommitTokens(ctx context.Context, in *MsgUncommitTokens, opts ...grpc.CallOption) (*MsgUncommitTokensResponse, error)
	// Vest converts user's commitment to vesting - start with unclaimed rewards
	// and if it's not enough deduct from committed bucket mainly utilized for
	// Eden
	Vest(ctx context.Context, in *MsgVest, opts ...grpc.CallOption) (*MsgVestResponse, error)
	// VestNow provides functionality to get the token immediately but lower
	// amount than original e.g. user can burn 1000 ueden and get 800 uelys when
	// the ratio is 80%
	VestNow(ctx context.Context, in *MsgVestNow, opts ...grpc.CallOption) (*MsgVestNowResponse, error)
	// VestLiquid converts user's balance to vesting to be utilized for normal
	// tokens vesting like ATOM vesting
	VestLiquid(ctx context.Context, in *MsgVestLiquid, opts ...grpc.CallOption) (*MsgVestLiquidResponse, error)
	// CancelVest cancel the user's vesting and the user reject to get vested
	// tokens
	CancelVest(ctx context.Context, in *MsgCancelVest, opts ...grpc.CallOption) (*MsgCancelVestResponse, error)
	// ClaimVesting claims already vested amount
	ClaimVesting(ctx context.Context, in *MsgClaimVesting, opts ...grpc.CallOption) (*MsgClaimVestingResponse, error)
	// UpdateVestingInfo add/update specific vesting info by denom on Params
	UpdateVestingInfo(ctx context.Context, in *MsgUpdateVestingInfo, opts ...grpc.CallOption) (*MsgUpdateVestingInfoResponse, error)
	// UpdateEnableVestNow add/update enable vest now on Params
	UpdateEnableVestNow(ctx context.Context, in *MsgUpdateEnableVestNow, opts ...grpc.CallOption) (*MsgUpdateEnableVestNowResponse, error)
	Stake(ctx context.Context, in *MsgStake, opts ...grpc.CallOption) (*MsgStakeResponse, error)
	Unstake(ctx context.Context, in *MsgUnstake, opts ...grpc.CallOption) (*MsgUnstakeResponse, error)
	ClaimAirdrop(ctx context.Context, in *MsgClaimAirdrop, opts ...grpc.CallOption) (*MsgClaimAirdropResponse, error)
	ClaimKol(ctx context.Context, in *MsgClaimKol, opts ...grpc.CallOption) (*MsgClaimKolResponse, error)
}

type msgClient struct {
	cc grpc.ClientConnInterface
}

func NewMsgClient(cc grpc.ClientConnInterface) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) CommitClaimedRewards(ctx context.Context, in *MsgCommitClaimedRewards, opts ...grpc.CallOption) (*MsgCommitClaimedRewardsResponse, error) {
	out := new(MsgCommitClaimedRewardsResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/CommitClaimedRewards", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UncommitTokens(ctx context.Context, in *MsgUncommitTokens, opts ...grpc.CallOption) (*MsgUncommitTokensResponse, error) {
	out := new(MsgUncommitTokensResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/UncommitTokens", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Vest(ctx context.Context, in *MsgVest, opts ...grpc.CallOption) (*MsgVestResponse, error) {
	out := new(MsgVestResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/Vest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) VestNow(ctx context.Context, in *MsgVestNow, opts ...grpc.CallOption) (*MsgVestNowResponse, error) {
	out := new(MsgVestNowResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/VestNow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) VestLiquid(ctx context.Context, in *MsgVestLiquid, opts ...grpc.CallOption) (*MsgVestLiquidResponse, error) {
	out := new(MsgVestLiquidResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/VestLiquid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) CancelVest(ctx context.Context, in *MsgCancelVest, opts ...grpc.CallOption) (*MsgCancelVestResponse, error) {
	out := new(MsgCancelVestResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/CancelVest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ClaimVesting(ctx context.Context, in *MsgClaimVesting, opts ...grpc.CallOption) (*MsgClaimVestingResponse, error) {
	out := new(MsgClaimVestingResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/ClaimVesting", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateVestingInfo(ctx context.Context, in *MsgUpdateVestingInfo, opts ...grpc.CallOption) (*MsgUpdateVestingInfoResponse, error) {
	out := new(MsgUpdateVestingInfoResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/UpdateVestingInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) UpdateEnableVestNow(ctx context.Context, in *MsgUpdateEnableVestNow, opts ...grpc.CallOption) (*MsgUpdateEnableVestNowResponse, error) {
	out := new(MsgUpdateEnableVestNowResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/UpdateEnableVestNow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Stake(ctx context.Context, in *MsgStake, opts ...grpc.CallOption) (*MsgStakeResponse, error) {
	out := new(MsgStakeResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/Stake", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) Unstake(ctx context.Context, in *MsgUnstake, opts ...grpc.CallOption) (*MsgUnstakeResponse, error) {
	out := new(MsgUnstakeResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/Unstake", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ClaimAirdrop(ctx context.Context, in *MsgClaimAirdrop, opts ...grpc.CallOption) (*MsgClaimAirdropResponse, error) {
	out := new(MsgClaimAirdropResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/ClaimAirdrop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *msgClient) ClaimKol(ctx context.Context, in *MsgClaimKol, opts ...grpc.CallOption) (*MsgClaimKolResponse, error) {
	out := new(MsgClaimKolResponse)
	err := c.cc.Invoke(ctx, "/elys.commitment.Msg/ClaimKol", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
// All implementations must embed UnimplementedMsgServer
// for forward compatibility
type MsgServer interface {
	// CommitClaimedRewards commit the tokens on claimed store to committed
	CommitClaimedRewards(context.Context, *MsgCommitClaimedRewards) (*MsgCommitClaimedRewardsResponse, error)
	// UncommitTokens uncommits the tokens from committed store and make it liquid
	// immediately
	UncommitTokens(context.Context, *MsgUncommitTokens) (*MsgUncommitTokensResponse, error)
	// Vest converts user's commitment to vesting - start with unclaimed rewards
	// and if it's not enough deduct from committed bucket mainly utilized for
	// Eden
	Vest(context.Context, *MsgVest) (*MsgVestResponse, error)
	// VestNow provides functionality to get the token immediately but lower
	// amount than original e.g. user can burn 1000 ueden and get 800 uelys when
	// the ratio is 80%
	VestNow(context.Context, *MsgVestNow) (*MsgVestNowResponse, error)
	// VestLiquid converts user's balance to vesting to be utilized for normal
	// tokens vesting like ATOM vesting
	VestLiquid(context.Context, *MsgVestLiquid) (*MsgVestLiquidResponse, error)
	// CancelVest cancel the user's vesting and the user reject to get vested
	// tokens
	CancelVest(context.Context, *MsgCancelVest) (*MsgCancelVestResponse, error)
	// ClaimVesting claims already vested amount
	ClaimVesting(context.Context, *MsgClaimVesting) (*MsgClaimVestingResponse, error)
	// UpdateVestingInfo add/update specific vesting info by denom on Params
	UpdateVestingInfo(context.Context, *MsgUpdateVestingInfo) (*MsgUpdateVestingInfoResponse, error)
	// UpdateEnableVestNow add/update enable vest now on Params
	UpdateEnableVestNow(context.Context, *MsgUpdateEnableVestNow) (*MsgUpdateEnableVestNowResponse, error)
	Stake(context.Context, *MsgStake) (*MsgStakeResponse, error)
	Unstake(context.Context, *MsgUnstake) (*MsgUnstakeResponse, error)
	ClaimAirdrop(context.Context, *MsgClaimAirdrop) (*MsgClaimAirdropResponse, error)
	ClaimKol(context.Context, *MsgClaimKol) (*MsgClaimKolResponse, error)
	mustEmbedUnimplementedMsgServer()
}

// UnimplementedMsgServer must be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (UnimplementedMsgServer) CommitClaimedRewards(context.Context, *MsgCommitClaimedRewards) (*MsgCommitClaimedRewardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommitClaimedRewards not implemented")
}
func (UnimplementedMsgServer) UncommitTokens(context.Context, *MsgUncommitTokens) (*MsgUncommitTokensResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UncommitTokens not implemented")
}
func (UnimplementedMsgServer) Vest(context.Context, *MsgVest) (*MsgVestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Vest not implemented")
}
func (UnimplementedMsgServer) VestNow(context.Context, *MsgVestNow) (*MsgVestNowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VestNow not implemented")
}
func (UnimplementedMsgServer) VestLiquid(context.Context, *MsgVestLiquid) (*MsgVestLiquidResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VestLiquid not implemented")
}
func (UnimplementedMsgServer) CancelVest(context.Context, *MsgCancelVest) (*MsgCancelVestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelVest not implemented")
}
func (UnimplementedMsgServer) ClaimVesting(context.Context, *MsgClaimVesting) (*MsgClaimVestingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClaimVesting not implemented")
}
func (UnimplementedMsgServer) UpdateVestingInfo(context.Context, *MsgUpdateVestingInfo) (*MsgUpdateVestingInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateVestingInfo not implemented")
}
func (UnimplementedMsgServer) UpdateEnableVestNow(context.Context, *MsgUpdateEnableVestNow) (*MsgUpdateEnableVestNowResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEnableVestNow not implemented")
}
func (UnimplementedMsgServer) Stake(context.Context, *MsgStake) (*MsgStakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stake not implemented")
}
func (UnimplementedMsgServer) Unstake(context.Context, *MsgUnstake) (*MsgUnstakeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unstake not implemented")
}
func (UnimplementedMsgServer) ClaimAirdrop(context.Context, *MsgClaimAirdrop) (*MsgClaimAirdropResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClaimAirdrop not implemented")
}
func (UnimplementedMsgServer) ClaimKol(context.Context, *MsgClaimKol) (*MsgClaimKolResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClaimKol not implemented")
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

func _Msg_CommitClaimedRewards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCommitClaimedRewards)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CommitClaimedRewards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/CommitClaimedRewards",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CommitClaimedRewards(ctx, req.(*MsgCommitClaimedRewards))
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
		FullMethod: "/elys.commitment.Msg/UncommitTokens",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UncommitTokens(ctx, req.(*MsgUncommitTokens))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Vest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgVest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Vest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/Vest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Vest(ctx, req.(*MsgVest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_VestNow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgVestNow)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).VestNow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/VestNow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).VestNow(ctx, req.(*MsgVestNow))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_VestLiquid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgVestLiquid)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).VestLiquid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/VestLiquid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).VestLiquid(ctx, req.(*MsgVestLiquid))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_CancelVest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgCancelVest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).CancelVest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/CancelVest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).CancelVest(ctx, req.(*MsgCancelVest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClaimVesting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimVesting)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimVesting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/ClaimVesting",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimVesting(ctx, req.(*MsgClaimVesting))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateVestingInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateVestingInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateVestingInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/UpdateVestingInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateVestingInfo(ctx, req.(*MsgUpdateVestingInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_UpdateEnableVestNow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUpdateEnableVestNow)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).UpdateEnableVestNow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/UpdateEnableVestNow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).UpdateEnableVestNow(ctx, req.(*MsgUpdateEnableVestNow))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Stake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgStake)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Stake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/Stake",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Stake(ctx, req.(*MsgStake))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Unstake_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgUnstake)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Unstake(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/Unstake",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Unstake(ctx, req.(*MsgUnstake))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClaimAirdrop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimAirdrop)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimAirdrop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/ClaimAirdrop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimAirdrop(ctx, req.(*MsgClaimAirdrop))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_ClaimKol_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgClaimKol)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).ClaimKol(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/elys.commitment.Msg/ClaimKol",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).ClaimKol(ctx, req.(*MsgClaimKol))
	}
	return interceptor(ctx, in, info, handler)
}

// Msg_ServiceDesc is the grpc.ServiceDesc for Msg service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Msg_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "elys.commitment.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CommitClaimedRewards",
			Handler:    _Msg_CommitClaimedRewards_Handler,
		},
		{
			MethodName: "UncommitTokens",
			Handler:    _Msg_UncommitTokens_Handler,
		},
		{
			MethodName: "Vest",
			Handler:    _Msg_Vest_Handler,
		},
		{
			MethodName: "VestNow",
			Handler:    _Msg_VestNow_Handler,
		},
		{
			MethodName: "VestLiquid",
			Handler:    _Msg_VestLiquid_Handler,
		},
		{
			MethodName: "CancelVest",
			Handler:    _Msg_CancelVest_Handler,
		},
		{
			MethodName: "ClaimVesting",
			Handler:    _Msg_ClaimVesting_Handler,
		},
		{
			MethodName: "UpdateVestingInfo",
			Handler:    _Msg_UpdateVestingInfo_Handler,
		},
		{
			MethodName: "UpdateEnableVestNow",
			Handler:    _Msg_UpdateEnableVestNow_Handler,
		},
		{
			MethodName: "Stake",
			Handler:    _Msg_Stake_Handler,
		},
		{
			MethodName: "Unstake",
			Handler:    _Msg_Unstake_Handler,
		},
		{
			MethodName: "ClaimAirdrop",
			Handler:    _Msg_ClaimAirdrop_Handler,
		},
		{
			MethodName: "ClaimKol",
			Handler:    _Msg_ClaimKol_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "elys/commitment/tx.proto",
}
