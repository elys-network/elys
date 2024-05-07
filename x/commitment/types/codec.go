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
	cdc.RegisterConcrete(&MsgClaimReward{}, "commitment/ClaimReward", nil)
	cdc.RegisterConcrete(&MsgVest{}, "commitment/Vest", nil)
	cdc.RegisterConcrete(&MsgClaimVesting{}, "commitment/ClaimVesting", nil)
	cdc.RegisterConcrete(&MsgCancelVest{}, "commitment/CancelVest", nil)
	cdc.RegisterConcrete(&MsgVestNow{}, "commitment/VestNow", nil)
	cdc.RegisterConcrete(&MsgUpdateVestingInfo{}, "commitment/UpdateVestingInfo", nil)
	cdc.RegisterConcrete(&MsgVestLiquid{}, "commitment/VestLiquid", nil)
	cdc.RegisterConcrete(&MsgClaimRewards{}, "commitment/ClaimRewards", nil)
	cdc.RegisterConcrete(&MsgStake{}, "commitment/Stake", nil)
	cdc.RegisterConcrete(&MsgUnstake{}, "commitment/Unstake", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCommitClaimedRewards{},
		&MsgUncommitTokens{},
		&MsgClaimReward{},
		&MsgVest{},
		&MsgClaimVesting{},
		&MsgCancelVest{},
		&MsgVestNow{},
		&MsgUpdateVestingInfo{},
		&MsgVestLiquid{},
		&MsgClaimRewards{},
		&MsgStake{},
		&MsgUnstake{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
