package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgFeedPrice{}, "oracle/FeedPrice", nil)
	cdc.RegisterConcrete(&MsgSetPriceFeeder{}, "oracle/SetPriceFeeder", nil)
	cdc.RegisterConcrete(&MsgDeletePriceFeeder{}, "oracle/DeletePriceFeeder", nil)
	cdc.RegisterConcrete(&MsgFeedMultiplePrices{}, "oracle/FeedMultiplePrices", nil)
	cdc.RegisterConcrete(&MsgAddAssetInfo{}, "oracle/AddAssetInfo", nil)
	cdc.RegisterConcrete(&MsgRemoveAssetInfo{}, "oracle/RemoveAssetInfo", nil)
	cdc.RegisterConcrete(&MsgAddPriceFeeders{}, "oracle/AddPriceFeeders", nil)
	cdc.RegisterConcrete(&MsgRemovePriceFeeders{}, "oracle/RemovePriceFeeders", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFeedPrice{},
		&MsgSetPriceFeeder{},
		&MsgDeletePriceFeeder{},
		&MsgFeedMultiplePrices{},
		&MsgAddAssetInfo{},
		&MsgRemoveAssetInfo{},
		&MsgAddPriceFeeders{},
		&MsgRemovePriceFeeders{},
	)

	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
