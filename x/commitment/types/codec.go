package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCommitClaimedRewards{}, "commitment/MsgCommitClaimedRewards")
	legacy.RegisterAminoMsg(cdc, &MsgUncommitTokens{}, "commitment/MsgUncommitTokens")
	legacy.RegisterAminoMsg(cdc, &MsgVest{}, "commitment/MsgVest")
	legacy.RegisterAminoMsg(cdc, &MsgClaimVesting{}, "commitment/MsgClaimVesting")
	legacy.RegisterAminoMsg(cdc, &MsgCancelVest{}, "commitment/MsgCancelVest")
	legacy.RegisterAminoMsg(cdc, &MsgVestNow{}, "commitment/MsgVestNow")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateVestingInfo{}, "commitment/MsgUpdateVestingInfo")
	legacy.RegisterAminoMsg(cdc, &MsgVestLiquid{}, "commitment/MsgVestLiquid")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateEnableVestNow{}, "commitment/MsgUpdateEnableVestNow")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateAirdropParams{}, "commitment/MsgUpdateAirdropParams")
	legacy.RegisterAminoMsg(cdc, &MsgStake{}, "commitment/MsgStake")
	legacy.RegisterAminoMsg(cdc, &MsgUnstake{}, "commitment/MsgUnstake")
	legacy.RegisterAminoMsg(cdc, &MsgClaimAirdrop{}, "commitment/MsgClaimAirdrop")
	legacy.RegisterAminoMsg(cdc, &MsgClaimKol{}, "commitment/MsgClaimKol")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCommitClaimedRewards{},
		&MsgUncommitTokens{},
		&MsgVest{},
		&MsgClaimVesting{},
		&MsgCancelVest{},
		&MsgVestNow{},
		&MsgUpdateVestingInfo{},
		&MsgVestLiquid{},
		&MsgUpdateEnableVestNow{},
		&MsgUpdateAirdropParams{},
		&MsgStake{},
		&MsgUnstake{},
		&MsgClaimAirdrop{},
		&MsgClaimKol{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
