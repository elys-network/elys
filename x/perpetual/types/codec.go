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
	legacy.RegisterAminoMsg(cdc, &MsgClose{}, "perpetual/MsgClose")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "perpetual/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgWhitelist{}, "perpetual/MsgWhitelist")
	legacy.RegisterAminoMsg(cdc, &MsgDewhitelist{}, "perpetual/MsgDewhitelist")
	legacy.RegisterAminoMsg(cdc, &MsgClosePositions{}, "perpetual/ClosePositions")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateStopLoss{}, "perpetual/MsgUpdateStopLoss")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateTakeProfitPrice{}, "perpetual/MsgUpdateTakeProfitPrice")
	legacy.RegisterAminoMsg(cdc, &MsgAddCollateral{}, "perpetual/MsgAddCollateral")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgOpen{},
		&MsgClose{},
		&MsgUpdateParams{},
		&MsgWhitelist{},
		&MsgDewhitelist{},
		&MsgClosePositions{},
		&MsgClosePositions{},
		&MsgUpdateStopLoss{},
		&MsgUpdateTakeProfitPrice{},
		&MsgAddCollateral{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
