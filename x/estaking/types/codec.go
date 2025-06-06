package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "estaking/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawReward{}, "estaking/MsgWithdrawReward")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawElysStakingRewards{}, "estaking/MsgWithdrawElysStakingRewards")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawAllRewards{}, "estaking/MsgWithdrawAllRewards")
	legacy.RegisterAminoMsg(cdc, &MsgUnjailGovernor{}, "estaking/MsgUnjailGovernor")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgWithdrawReward{},
		&MsgWithdrawElysStakingRewards{},
		&MsgWithdrawAllRewards{},
		&MsgUnjailGovernor{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
