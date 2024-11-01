package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateSpotOrder{}, "tradeshield/CreateSpotOrder")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateSpotOrder{}, "tradeshield/UpdateSpotOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCancelSpotOrder{}, "tradeshield/CancelSpotOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCancelSpotOrders{}, "tradeshield/CancelSpotOrders")

	legacy.RegisterAminoMsg(cdc, &MsgCreatePerpetualOpenOrder{}, "tradeshield/CreatePerpetualOpenOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCreatePerpetualCloseOrder{}, "tradeshield/CreatePerpetualCloseOrder")
	legacy.RegisterAminoMsg(cdc, &MsgUpdatePerpetualOrder{}, "tradeshield/UpdatePerpetualOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCancelPerpetualOrder{}, "tradeshield/CancelPerpetualOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCancelPerpetualOrders{}, "tradeshield/CancelPerpetualOrders")

	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "tradeshield/UpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgExecuteOrders{}, "tradeshield/ExecuteOrders")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateSpotOrder{},
		&MsgUpdateSpotOrder{},
		&MsgCancelSpotOrder{},
		&MsgCancelSpotOrders{},

		&MsgCreatePerpetualOpenOrder{},
		&MsgCreatePerpetualCloseOrder{},
		&MsgUpdatePerpetualOrder{},
		&MsgCancelPerpetualOrder{},
		&MsgCancelPerpetualOrders{},

		&MsgUpdateParams{},
		&MsgExecuteOrders{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
