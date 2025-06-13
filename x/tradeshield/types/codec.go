package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateSpotOrder{}, "tradeshield/MsgCreateSpotOrder")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateSpotOrder{}, "tradeshield/MsgUpdateSpotOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCancelSpotOrder{}, "tradeshield/MsgCancelSpotOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCancelSpotOrders{}, "tradeshield/MsgCancelSpotOrders")
	legacy.RegisterAminoMsg(cdc, &MsgCancelAllSpotOrders{}, "tradeshield/MsgCancelAllSpotOrders")

	legacy.RegisterAminoMsg(cdc, &MsgCreatePerpetualOpenOrder{}, "tradeshield/MsgCreatePerpetualOpenOrder")
	// TODO: Use Msg... Structure with v2, when close perpetual position is implemented
	legacy.RegisterAminoMsg(cdc, &MsgCreatePerpetualCloseOrder{}, "tradeshield/CreatePerpetualCloseOrder")
	legacy.RegisterAminoMsg(cdc, &MsgUpdatePerpetualOrder{}, "tradeshield/MsgUpdatePerpetualOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCancelPerpetualOrder{}, "tradeshield/MsgCancelPerpetualOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCancelPerpetualOrders{}, "tradeshield/MsgCancelPerpetualOrders")
	legacy.RegisterAminoMsg(cdc, &MsgCancelAllPerpetualOrders{}, "tradeshield/MsgCancelAllPerpetualOrders")

	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "tradeshield/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgExecuteOrders{}, "tradeshield/MsgExecuteOrders")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateSpotOrder{},
		&MsgUpdateSpotOrder{},
		&MsgCancelSpotOrder{},
		&MsgCancelSpotOrders{},
		&MsgCancelAllSpotOrders{},

		&MsgCreatePerpetualOpenOrder{},
		&MsgCreatePerpetualCloseOrder{},
		&MsgUpdatePerpetualOrder{},
		&MsgCancelPerpetualOrder{},
		&MsgCancelPerpetualOrders{},
		&MsgCancelAllPerpetualOrders{},

		&MsgUpdateParams{},
		&MsgExecuteOrders{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
