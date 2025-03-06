package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "clob/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgCreatPerpetualMarket{}, "clob/MsgCreatPerpetualMarket")
	legacy.RegisterAminoMsg(cdc, &MsgPlaceLimitOrder{}, "clob/MsgPlaceLimitOrder")
	legacy.RegisterAminoMsg(cdc, &MsgPlaceMarketOrder{}, "clob/MsgPlaceMarketOrder")
	legacy.RegisterAminoMsg(cdc, &MsgDeposit{}, "clob/MsgDeposit")
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgPlaceLimitOrder{},
		&MsgPlaceMarketOrder{},
		&MsgCreatPerpetualMarket{},
		&MsgDeposit{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
