package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgFeedPrice{}, "oracle/MsgFeedPrice")
	legacy.RegisterAminoMsg(cdc, &MsgSetPriceFeeder{}, "oracle/MsgSetPriceFeeder")
	legacy.RegisterAminoMsg(cdc, &MsgDeletePriceFeeder{}, "oracle/MsgDeletePriceFeeder")
	legacy.RegisterAminoMsg(cdc, &MsgFeedMultiplePrices{}, "oracle/MsgFeedMultiplePrices")
	legacy.RegisterAminoMsg(cdc, &MsgRemoveAssetInfo{}, "oracle/MsgRemoveAssetInfo")
	legacy.RegisterAminoMsg(cdc, &MsgAddPriceFeeders{}, "oracle/MsgAddPriceFeeders")
	legacy.RegisterAminoMsg(cdc, &MsgRemovePriceFeeders{}, "oracle/MsgRemovePriceFeeders")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "oracle/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgCreateAssetInfo{}, "oracle/CreateAssetInfo")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFeedPrice{},
		&MsgSetPriceFeeder{},
		&MsgDeletePriceFeeder{},
		&MsgFeedMultiplePrices{},
		&MsgRemoveAssetInfo{},
		&MsgAddPriceFeeders{},
		&MsgRemovePriceFeeders{},
		&MsgUpdateParams{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAssetInfo{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
