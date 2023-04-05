package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCoinRatesData{}, "oracle/CoinRatesData", nil)
	cdc.RegisterConcrete(&MsgCreateAssetInfo{}, "oracle/CreateAssetInfo", nil)
	cdc.RegisterConcrete(&MsgUpdateAssetInfo{}, "oracle/UpdateAssetInfo", nil)
	cdc.RegisterConcrete(&MsgDeleteAssetInfo{}, "oracle/DeleteAssetInfo", nil)
	cdc.RegisterConcrete(&MsgFeedPrice{}, "oracle/FeedPrice", nil)
	cdc.RegisterConcrete(&MsgCreatePriceFeeder{}, "oracle/CreatePriceFeeder", nil)
	cdc.RegisterConcrete(&MsgUpdatePriceFeeder{}, "oracle/UpdatePriceFeeder", nil)
	cdc.RegisterConcrete(&MsgDeletePriceFeeder{}, "oracle/DeletePriceFeeder", nil)
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

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&ProposalAddAssetInfo{},
		&ProposalRemoveAssetInfo{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFeedPrice{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePriceFeeder{},
		&MsgUpdatePriceFeeder{},
		&MsgDeletePriceFeeder{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
