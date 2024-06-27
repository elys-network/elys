package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgOpen{}, "leveragelp/Open", nil)
	cdc.RegisterConcrete(&MsgClose{}, "leveragelp/Close", nil)
	cdc.RegisterConcrete(&MsgUpdateParams{}, "leveragelp/UpdateParams", nil)
	cdc.RegisterConcrete(&MsgUpdatePools{}, "leveragelp/UpdatePools", nil)
	cdc.RegisterConcrete(&MsgWhitelist{}, "leveragelp/Whitelist", nil)
	cdc.RegisterConcrete(&MsgDewhitelist{}, "leveragelp/Dewhitelist", nil)
	cdc.RegisterConcrete(&MsgClaimRewards{}, "leveragelp/ClaimRewards", nil)
	cdc.RegisterConcrete(&MsgUpdateStopLoss{}, "leveragelp/UpdateStopLoss", nil)
	cdc.RegisterConcrete(&MsgAddCollateral{}, "leveragelp/AddCollateral", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgOpen{},
		&MsgClose{},
		&MsgUpdateParams{},
		&MsgUpdatePools{},
		&MsgWhitelist{},
		&MsgDewhitelist{},
		&MsgClaimRewards{},
		&MsgUpdateStopLoss{},
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
