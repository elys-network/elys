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
	legacy.RegisterAminoMsg(cdc, &MsgFeedPrice{}, "oracle/FeedPrice")
	legacy.RegisterAminoMsg(cdc, &MsgSetPriceFeeder{}, "oracle/SetPriceFeeder")
	legacy.RegisterAminoMsg(cdc, &MsgDeletePriceFeeder{}, "oracle/DeletePriceFeeder")
	legacy.RegisterAminoMsg(cdc, &MsgFeedMultiplePrices{}, "oracle/FeedMultiplePrices")
	legacy.RegisterAminoMsg(cdc, &MsgAddAssetInfo{}, "oracle/AddAssetInfo")
	legacy.RegisterAminoMsg(cdc, &MsgRemoveAssetInfo{}, "oracle/RemoveAssetInfo")
	legacy.RegisterAminoMsg(cdc, &MsgAddPriceFeeders{}, "oracle/AddPriceFeeders")
	legacy.RegisterAminoMsg(cdc, &MsgRemovePriceFeeders{}, "oracle/RemovePriceFeeders")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "oracle/UpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgCreateAssetInfo{}, "oracle/CreateAssetInfo")
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

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAssetInfo{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
	cryptocodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)

	// Register all Amino interfaces and concrete types on the authz  and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterCodec(authzcodec.Amino)
	RegisterCodec(govcodec.Amino)
	RegisterCodec(groupcodec.Amino)
}
