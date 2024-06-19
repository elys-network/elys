package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgOpen{}, "perpetual/Open", nil)
	cdc.RegisterConcrete(&MsgClose{}, "perpetual/Close", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "perpetual/UpdateParams", nil)
	cdc.RegisterConcrete(&MsgWhitelist{}, "perpetual/Whitelist", nil)
	cdc.RegisterConcrete(&MsgDewhitelist{}, "perpetual/Dewhitelist", nil)
	cdc.RegisterConcrete(&MsgAddCollateral{}, "perpetual/AddCollateral", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgOpen{},
		&MsgClose{},
		&MsgUpdateParams{},
		&MsgWhitelist{},
		&MsgDewhitelist{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddCollateral{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
