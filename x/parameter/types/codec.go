package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// this line is used by starport scaffolding # 1
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgUpdateMinCommission{}, "parameter/MsgUpdateMinCommission", nil)
	cdc.RegisterConcrete(&MsgUpdateMaxVotingPower{}, "parameter/MsgUpdateMaxVotingPower", nil)
	cdc.RegisterConcrete(&MsgUpdateMinSelfDelegation{}, "parameter/MsgUpdateMinSelfDelegation", nil)
	cdc.RegisterConcrete(&MsgUpdateBrokerAddress{}, "parameter/MsgUpdateBrokerAddress", nil)
	cdc.RegisterConcrete(&MsgUpdateTotalBlocksPerYear{}, "parameter/MsgUpdateTotalBlocksPerYear", nil)
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
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
