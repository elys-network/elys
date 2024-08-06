package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
	groupcodec "github.com/cosmos/cosmos-sdk/x/group/codec"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateWasmConfig{}, "parameter/MsgUpdateWasmConfig")
	// this line is used by starport scaffolding # 2
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMinCommission{}, "parameter/MsgUpdateMinCommission")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMaxVotingPower{}, "parameter/MsgUpdateMaxVotingPower")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMinSelfDelegation{}, "parameter/MsgUpdateMinSelfDelegation")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateBrokerAddress{}, "parameter/MsgUpdateBrokerAddress")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateTotalBlocksPerYear{}, "parameter/MsgUpdateTotalBlocksPerYear")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateMinCommission{},
		&MsgUpdateMaxVotingPower{},
		&MsgUpdateMinSelfDelegation{},
		&MsgUpdateBrokerAddress{},
		&MsgUpdateTotalBlocksPerYear{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateWasmConfig{},
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
