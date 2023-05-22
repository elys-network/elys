package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSetWithdrawAddress{}, "incentive/SetWithdrawAddress", nil)
	cdc.RegisterConcrete(&MsgWithdrawValidatorCommission{}, "incentive/WithdrawValidatorCommission", nil)
	cdc.RegisterConcrete(&MsgWithdrawDelegatorReward{}, "incentive/WithdrawDelegatorReward", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetWithdrawAddress{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdrawValidatorCommission{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdrawDelegatorReward{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&ProposalUpdateMinCommission{},
		&ProposalUpdateMaxVotingPower{},
		&ProposalUpdateMinSelfDelegation{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
