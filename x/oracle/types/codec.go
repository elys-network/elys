package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
	groupcodec "github.com/cosmos/cosmos-sdk/x/group/codec"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgFeedPrice{}, "oracle/MsgFeedPrice")
	legacy.RegisterAminoMsg(cdc, &MsgSetPriceFeeder{}, "oracle/MsgSetPriceFeeder")
	legacy.RegisterAminoMsg(cdc, &MsgDeletePriceFeeder{}, "oracle/MsgDeletePriceFeeder")
	legacy.RegisterAminoMsg(cdc, &MsgFeedMultiplePrices{}, "oracle/MsgFeedMultiplePrices")
	legacy.RegisterAminoMsg(cdc, &MsgAddAssetInfo{}, "oracle/MsgAddAssetInfo")
	legacy.RegisterAminoMsg(cdc, &MsgRemoveAssetInfo{}, "oracle/MsgRemoveAssetInfo")
	legacy.RegisterAminoMsg(cdc, &MsgAddPriceFeeders{}, "oracle/MsgAddPriceFeeders")
	legacy.RegisterAminoMsg(cdc, &MsgRemovePriceFeeders{}, "oracle/MsgRemovePriceFeeders")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "oracle/MsgUpdateParams")
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
		&MsgUpdateParams{},
	)

	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino          = codec.NewLegacyAmino()
	ModuleAminoCdc = codec.NewAminoCodec(amino)
	ModuleCdc      = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Register all Amino interfaces and concrete types on the authz  and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterCodec(authzcodec.Amino)
	RegisterCodec(govcodec.Amino)
	RegisterCodec(groupcodec.Amino)
}
