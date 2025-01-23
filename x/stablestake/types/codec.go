package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgBond{}, "stablestake/MsgBond")
	legacy.RegisterAminoMsg(cdc, &MsgUnbond{}, "stablestake/MsgUnbond")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "stablestake/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgAddPool{}, "stablestake/MsgAddPool")
	legacy.RegisterAminoMsg(cdc, &MsgUpdatePool{}, "stablestake/MsgUpdatePool")

	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBond{},
		&MsgUnbond{},
		&MsgUpdateParams{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
