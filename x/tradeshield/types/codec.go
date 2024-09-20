package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreatePendingSpotOrder{}, "tradeshield/CreatePendingSpotOrder")
	legacy.RegisterAminoMsg(cdc, &MsgUpdatePendingSpotOrder{}, "tradeshield/UpdatePendingSpotOrder")
	legacy.RegisterAminoMsg(cdc, &MsgDeletePendingSpotOrder{}, "tradeshield/DeletePendingSpotOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCreatePendingPerpetualOrder{}, "tradeshield/CreatePendingPerpetualOrder")
	legacy.RegisterAminoMsg(cdc, &MsgUpdatePendingPerpetualOrder{}, "tradeshield/UpdatePendingPerpetualOrder")
	legacy.RegisterAminoMsg(cdc, &MsgDeletePendingPerpetualOrder{}, "tradeshield/DeletePendingPerpetualOrder")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "tradeshield/UpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgExecuteOrders{}, "tradeshield/ExecuteOrders")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePendingSpotOrder{},
		&MsgUpdatePendingSpotOrder{},
		&MsgDeletePendingSpotOrder{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePendingPerpetualOrder{},
		&MsgUpdatePendingPerpetualOrder{},
		&MsgDeletePendingPerpetualOrder{},
		&MsgUpdateParams{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgExecuteOrders{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
