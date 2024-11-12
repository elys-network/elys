package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgAddExternalIncentive{}, "masterchef/MsgAddExternalIncentive")
	legacy.RegisterAminoMsg(cdc, &MsgClaimRewards{}, "masterchef/MsgClaimRewards")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "masterchef/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgUpdatePoolMultipliers{}, "masterchef/MsgUpdatePoolMultipliers")
	legacy.RegisterAminoMsg(cdc, &MsgAddExternalRewardDenom{}, "masterchef/MsgAddExternalRewardDenom")
	legacy.RegisterAminoMsg(cdc, &MsgTogglePoolEdenRewards{}, "masterchef/TogglePoolEdenRewards")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddExternalIncentive{},
		&MsgClaimRewards{},
		&MsgUpdateParams{},
		&MsgUpdatePoolMultipliers{},
		&MsgAddExternalRewardDenom{},
		&MsgTogglePoolEdenRewards{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
