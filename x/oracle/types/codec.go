package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCoinRatesData{}, "oracle/CoinRatesData", nil)
	cdc.RegisterConcrete(&MsgCreateAssetInfo{}, "oracle/CreateAssetInfo", nil)
	cdc.RegisterConcrete(&MsgUpdateAssetInfo{}, "oracle/UpdateAssetInfo", nil)
	cdc.RegisterConcrete(&MsgDeleteAssetInfo{}, "oracle/DeleteAssetInfo", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCoinRatesData{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAssetInfo{},
		&MsgUpdateAssetInfo{},
		&MsgDeleteAssetInfo{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
