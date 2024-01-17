package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgFeedPrice{}, "oracle/FeedPrice", nil)
	cdc.RegisterConcrete(&MsgSetPriceFeeder{}, "oracle/SetPriceFeeder", nil)
	cdc.RegisterConcrete(&MsgDeletePriceFeeder{}, "oracle/DeletePriceFeeder", nil)
	cdc.RegisterConcrete(&MsgFeedMultiplePrices{}, "oracle/FeedMultiplePrices", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFeedPrice{},
		&MsgSetPriceFeeder{},
		&MsgDeletePriceFeeder{},
		&MsgFeedMultiplePrices{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&ProposalAddAssetInfo{},
		&ProposalRemoveAssetInfo{},
		&ProposalAddPriceFeeders{},
		&ProposalRemovePriceFeeders{},
	)

	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
