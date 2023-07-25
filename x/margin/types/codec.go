package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgOpen{}, "margin/Open", nil)
	cdc.RegisterConcrete(&MsgClose{}, "margin/Close", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgOpen{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgClose{},
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&ProposalUpdateParams{},
		&ProposalUpdatePools{},
		&ProposalWhitelist{},
		&ProposalDewhitelist{},
	)

	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
