package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	// this line is used by starport scaffolding # 1
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateParams{}, "estaking/UpdateParams", nil)
	cdc.RegisterConcrete(&MsgWithdrawReward{}, "estaking/WithdrawReward", nil)
	cdc.RegisterConcrete(&MsgWithdrawElysStakingRewards{}, "estaking/WithdrawElysStakingRewards", nil)
	cdc.RegisterConcrete(&MsgWithdrawAllRewards{}, "estaking/MsgWithdrawAllRewards", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgWithdrawReward{},
		&MsgWithdrawElysStakingRewards{},
		&MsgWithdrawAllRewards{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
