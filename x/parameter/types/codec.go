package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateRewardsDataLifetime{}, "parameter/UpdateRewardsDataLifetime", nil)
	// this line is used by starport scaffolding # 2
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMinCommission{}, "parameter/MsgUpdateMinCommission")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMaxVotingPower{}, "parameter/MsgUpdateMaxVotingPower")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateMinSelfDelegation{}, "parameter/MsgUpdateMinSelfDelegation")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateTotalBlocksPerYear{}, "parameter/MsgUpdateTotalBlocksPerYear")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateMinCommission{},
		&MsgUpdateMaxVotingPower{},
		&MsgUpdateMinSelfDelegation{},
		&MsgUpdateTotalBlocksPerYear{},
		&MsgUpdateRewardsDataLifetime{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
