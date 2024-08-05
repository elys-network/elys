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
	legacy.RegisterAminoMsg(cdc, &MsgCreateAirdrop{}, "elys/tokenomics/MsgCreateAirdrop")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateAirdrop{}, "elys/tokenomics/MsgUpdateAirdrop")
	legacy.RegisterAminoMsg(cdc, &MsgDeleteAirdrop{}, "elys/tokenomics/MsgDeleteAirdrop")
	legacy.RegisterAminoMsg(cdc, &MsgClaimAirdrop{}, "elys/tokenomics/MsgClaimAirdrop")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateGenesisInflation{}, "elys/tokenomics/MsgUpdateGenesisInflation")
	legacy.RegisterAminoMsg(cdc, &MsgCreateTimeBasedInflation{}, "elys/tokenomics/MsgCreateTimeBasedInflation")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateTimeBasedInflation{}, "elys/tokenomics/MsgUpdateTimeBasedInflation")
	legacy.RegisterAminoMsg(cdc, &MsgDeleteTimeBasedInflation{}, "elys/tokenomics/MsgDeleteTimeBasedInflation")
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