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
	legacy.RegisterAminoMsg(cdc, &MsgCommitClaimedRewards{}, "commitment/CommitClaimedRewards")
	legacy.RegisterAminoMsg(cdc, &MsgUncommitTokens{}, "commitment/UncommitTokens")
	legacy.RegisterAminoMsg(cdc, &MsgClaimReward{}, "commitment/ClaimReward")
	legacy.RegisterAminoMsg(cdc, &MsgVest{}, "commitment/Vest")
	legacy.RegisterAminoMsg(cdc, &MsgClaimVesting{}, "commitment/ClaimVesting")
	legacy.RegisterAminoMsg(cdc, &MsgCancelVest{}, "commitment/CancelVest")
	legacy.RegisterAminoMsg(cdc, &MsgVestNow{}, "commitment/VestNow")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateVestingInfo{}, "commitment/UpdateVestingInfo")
	legacy.RegisterAminoMsg(cdc, &MsgVestLiquid{}, "commitment/VestLiquid")
	legacy.RegisterAminoMsg(cdc, &MsgClaimRewards{}, "commitment/ClaimRewards")
	legacy.RegisterAminoMsg(cdc, &MsgStake{}, "commitment/Stake")
	legacy.RegisterAminoMsg(cdc, &MsgUnstake{}, "commitment/Unstake")
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