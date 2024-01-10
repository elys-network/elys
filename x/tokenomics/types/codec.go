package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateAirdrop{}, "tokenomics/CreateAirdrop", nil)
	cdc.RegisterConcrete(&MsgUpdateAirdrop{}, "tokenomics/UpdateAirdrop", nil)
	cdc.RegisterConcrete(&MsgDeleteAirdrop{}, "tokenomics/DeleteAirdrop", nil)
	cdc.RegisterConcrete(&MsgClaimAirdrop{}, "tokenomics/ClaimAirdrop", nil)
	cdc.RegisterConcrete(&MsgUpdateGenesisInflation{}, "tokenomics/UpdateGenesisInflation", nil)
	cdc.RegisterConcrete(&MsgCreateTimeBasedInflation{}, "tokenomics/CreateTimeBasedInflation", nil)
	cdc.RegisterConcrete(&MsgUpdateTimeBasedInflation{}, "tokenomics/UpdateTimeBasedInflation", nil)
	cdc.RegisterConcrete(&MsgDeleteTimeBasedInflation{}, "tokenomics/DeleteTimeBasedInflation", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAirdrop{},
		&MsgUpdateAirdrop{},
		&MsgDeleteAirdrop{},
		&MsgClaimAirdrop{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateGenesisInflation{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateTimeBasedInflation{},
		&MsgUpdateTimeBasedInflation{},
		&MsgDeleteTimeBasedInflation{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
