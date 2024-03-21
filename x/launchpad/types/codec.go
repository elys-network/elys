package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgBuyElys{}, "launchpad/BuyElys", nil)
	cdc.RegisterConcrete(&MsgReturnElys{}, "launchpad/ReturnElys", nil)
	cdc.RegisterConcrete(&MsgWithdrawRaised{}, "launchpad/WithdrawRaised", nil)
	cdc.RegisterConcrete(&MsgDepositElysToken{}, "launchpad/DepositElysToken", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBuyElys{},
		&MsgReturnElys{},
		&MsgWithdrawRaised{},
		&MsgDepositElysToken{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
