package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCommitTokens{}, "commitment/CommitTokens", nil)
	cdc.RegisterConcrete(&MsgUncommitTokens{}, "commitment/UncommitTokens", nil)
	cdc.RegisterConcrete(&MsgWithdrawTokens{}, "commitment/WithdrawTokens", nil)
	cdc.RegisterConcrete(&MsgDepositTokens{}, "commitment/DepositTokens", nil)
	cdc.RegisterConcrete(&MsgVest{}, "commitment/Vest", nil)
	cdc.RegisterConcrete(&MsgCancelVest{}, "commitment/CancelVest", nil)
	cdc.RegisterConcrete(&MsgVestNow{}, "commitment/VestNow", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCommitTokens{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUncommitTokens{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdrawTokens{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDepositTokens{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgVest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCancelVest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgVestNow{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
