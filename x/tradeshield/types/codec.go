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
	legacy.RegisterAminoMsg(cdc, &MsgCreateSpotOrder{}, "tradeshield/CreateSpotOrder")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateSpotOrder{}, "tradeshield/UpdateSpotOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCancelSpotOrders{}, "tradeshield/CancelSpotOrders")

	legacy.RegisterAminoMsg(cdc, &MsgCreatePerpetualOpenOrder{}, "tradeshield/CreatePerpetualOpenOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCreatePerpetualCloseOrder{}, "tradeshield/CreatePerpetualCloseOrder")
	legacy.RegisterAminoMsg(cdc, &MsgUpdatePerpetualOrder{}, "tradeshield/UpdatePerpetualOrder")
	legacy.RegisterAminoMsg(cdc, &MsgCancelPerpetualOrders{}, "tradeshield/CancelPerpetualOrders")

	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "tradeshield/UpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgExecuteOrders{}, "tradeshield/ExecuteOrders")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateSpotOrder{},
		&MsgUpdateSpotOrder{},
		&MsgCancelSpotOrders{},

		&MsgCreatePerpetualOpenOrder{},
		&MsgCreatePerpetualCloseOrder{},
		&MsgUpdatePerpetualOrder{},
		&MsgCancelPerpetualOrders{},

		&MsgUpdateParams{},
		&MsgExecuteOrders{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
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
