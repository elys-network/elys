package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	// this line is used by starport scaffolding # 1
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
	groupcodec "github.com/cosmos/cosmos-sdk/x/group/codec"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "elys/estaking/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawReward{}, "elys/estaking/MsgWithdrawReward")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawElysStakingRewards{}, "elys/estaking/MsgWithdrawElysStakingRewards")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawAllRewards{}, "elys/estaking/MsgWithdrawAllRewards")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgWithdrawReward{},
		&MsgWithdrawElysStakingRewards{},
		&MsgWithdrawAllRewards{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
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