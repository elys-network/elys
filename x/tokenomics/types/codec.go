package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgCreateAirdrop{}, "tokenomics/MsgCreateAirdrop")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateAirdrop{}, "tokenomics/MsgUpdateAirdrop")
	legacy.RegisterAminoMsg(cdc, &MsgDeleteAirdrop{}, "tokenomics/MsgDeleteAirdrop")
	legacy.RegisterAminoMsg(cdc, &MsgClaimAirdrop{}, "tokenomics/MsgClaimAirdrop")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateGenesisInflation{}, "tokenomics/MsgUpdateGenesisInflation")
	legacy.RegisterAminoMsg(cdc, &MsgCreateTimeBasedInflation{}, "tokenomics/MsgCreateTimeBasedInflation")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateTimeBasedInflation{}, "tokenomics/MsgUpdateTimeBasedInflation")
	legacy.RegisterAminoMsg(cdc, &MsgDeleteTimeBasedInflation{}, "tokenomics/MsgDeleteTimeBasedInflation")
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateAirdrop{},
		&MsgUpdateAirdrop{},
		&MsgDeleteAirdrop{},
		&MsgClaimAirdrop{},
		&MsgCreateTimeBasedInflation{},
		&MsgUpdateTimeBasedInflation{},
		&MsgDeleteTimeBasedInflation{},
		&MsgUpdateGenesisInflation{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
