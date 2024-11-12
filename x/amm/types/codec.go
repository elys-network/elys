package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreatePool{}, "amm/MsgCreatePool")
	legacy.RegisterAminoMsg(cdc, &MsgJoinPool{}, "amm/MsgJoinPool")
	legacy.RegisterAminoMsg(cdc, &MsgExitPool{}, "amm/MsgExitPool")
	legacy.RegisterAminoMsg(cdc, &MsgSwapExactAmountIn{}, "amm/MsgSwapExactAmountIn")
	legacy.RegisterAminoMsg(cdc, &MsgSwapExactAmountOut{}, "amm/MsgSwapExactAmountOut")
	legacy.RegisterAminoMsg(cdc, &MsgSwapByDenom{}, "amm/MsgSwapByDenom")
	legacy.RegisterAminoMsg(cdc, &MsgFeedMultipleExternalLiquidity{}, "amm/MsgFeedMultipleExternalLiquidity")
	legacy.RegisterAminoMsg(cdc, &MsgUpdatePoolParams{}, "amm/MsgUpdatePoolParams")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "amm/MsgUpdateParams")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePool{},
		&MsgJoinPool{},
		&MsgExitPool{},
		&MsgSwapExactAmountIn{},
		&MsgSwapExactAmountOut{},
		&MsgSwapByDenom{},
		&MsgFeedMultipleExternalLiquidity{},
		&MsgUpdateParams{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
