package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgOpen{}, "perpetual/MsgOpen")
	legacy.RegisterAminoMsg(cdc, &MsgBrokerOpen{}, "perpetual/MsgBrokerOpen")
	legacy.RegisterAminoMsg(cdc, &MsgClose{}, "perpetual/MsgClose")
	legacy.RegisterAminoMsg(cdc, &MsgBrokerClose{}, "perpetual/MsgBrokerClose")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "perpetual/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgWhitelist{}, "perpetual/MsgWhitelist")
	legacy.RegisterAminoMsg(cdc, &MsgDewhitelist{}, "perpetual/MsgDewhitelist")
	legacy.RegisterAminoMsg(cdc, &MsgAddCollateral{}, "perpetual/MsgAddCollateral")
	legacy.RegisterAminoMsg(cdc, &MsgBrokerAddCollateral{}, "perpetual/MsgBrokerAddCollateral")
	legacy.RegisterAminoMsg(cdc, &MsgClosePositions{}, "perpetual/ClosePositions")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateStopLoss{}, "perpetual/MsgUpdateStopLoss")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateTakeProfitPrice{}, "perpetual/MsgUpdateTakeProfitPrice")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgOpen{},
		&MsgBrokerOpen{},
		&MsgClose{},
		&MsgBrokerClose{},
		&MsgUpdateParams{},
		&MsgWhitelist{},
		&MsgDewhitelist{},
		&MsgAddCollateral{},
		&MsgBrokerAddCollateral{},
		&MsgClosePositions{},
		&MsgUpdateStopLoss{},
		&MsgUpdateTakeProfitPrice{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
