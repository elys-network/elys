package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgOpen{}, "leveragelp/MsgOpen")
	legacy.RegisterAminoMsg(cdc, &MsgClose{}, "leveragelp/MsgClose")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "leveragelp/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgAddPool{}, "leveragelp/MsgAddPool")
	legacy.RegisterAminoMsg(cdc, &MsgRemovePool{}, "leveragelp/MsgRemovePool")
	legacy.RegisterAminoMsg(cdc, &MsgWhitelist{}, "leveragelp/MsgWhitelist")
	legacy.RegisterAminoMsg(cdc, &MsgDewhitelist{}, "leveragelp/MsgDewhitelist")
	legacy.RegisterAminoMsg(cdc, &MsgClaimRewards{}, "leveragelp/MsgClaimRewards")
	legacy.RegisterAminoMsg(cdc, &MsgClaimAllRewards{}, "leveragelp/MsgClaimAllRewards")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateStopLoss{}, "leveragelp/MsgUpdateStopLoss")
	legacy.RegisterAminoMsg(cdc, &MsgClosePositions{}, "leveragelp/MsgClosePositions")
	legacy.RegisterAminoMsg(cdc, &MsgUpdatePool{}, "leveragelp/MsgUpdatePool")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateEnabledPools{}, "leveragelp/MsgUpdateEnabledPools")

	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgOpen{},
		&MsgClose{},
		&MsgUpdateParams{},
		&MsgAddPool{},
		&MsgRemovePool{},
		&MsgWhitelist{},
		&MsgDewhitelist{},
		&MsgClaimRewards{},
		&MsgClaimAllRewards{},
		&MsgUpdateStopLoss{},
		&MsgClosePositions{},
		&MsgUpdatePool{},
		&MsgUpdateEnabledPools{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
