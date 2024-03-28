package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddExternalIncentive{}, "masterchef/AddExternalIncentive", nil)
	cdc.RegisterConcrete(&MsgClaimRewards{}, "masterchef/ClaimRewards", nil)
	cdc.RegisterConcrete(&MsgUpdateIncentiveParams{}, "masterchef/UpdateIncentiveParams", nil)
	cdc.RegisterConcrete(&MsgUpdatePoolMultipliers{}, "masterchef/UpdatePoolMultipliers", nil)
	cdc.RegisterConcrete(&MsgAddExternalRewardDenom{}, "masterchef/AddExternalRewardDenom", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddExternalIncentive{},
		&MsgClaimRewards{},
		&MsgUpdateIncentiveParams{},
		&MsgUpdatePoolMultipliers{},
		&MsgAddExternalRewardDenom{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
