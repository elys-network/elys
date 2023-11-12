package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCommitClaimedRewards{}, "commitment/CommitClaimedRewards", nil)
	cdc.RegisterConcrete(&MsgUncommitTokens{}, "commitment/UncommitTokens", nil)
	cdc.RegisterConcrete(&MsgWithdrawTokens{}, "commitment/WithdrawTokens", nil)
	cdc.RegisterConcrete(&MsgCommitLiquidTokens{}, "commitment/CommitLiquidTokens", nil)
	cdc.RegisterConcrete(&MsgVest{}, "commitment/Vest", nil)
	cdc.RegisterConcrete(&MsgCancelVest{}, "commitment/CancelVest", nil)
	cdc.RegisterConcrete(&MsgVestNow{}, "commitment/VestNow", nil)
	cdc.RegisterConcrete(&MsgUpdateVestingInfo{}, "commitment/UpdateVestingInfo", nil)
	cdc.RegisterConcrete(&MsgVestLiquid{}, "commitment/VestLiquid", nil)
	cdc.RegisterConcrete(&MsgClaimRewards{}, "commitment/ClaimRewards", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCommitClaimedRewards{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUncommitTokens{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdrawTokens{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCommitLiquidTokens{},
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
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateVestingInfo{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgVestLiquid{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgClaimRewards{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
