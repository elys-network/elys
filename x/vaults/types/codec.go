package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "vaults/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgAddVault{}, "vaults/MsgAddVault")
	legacy.RegisterAminoMsg(cdc, &MsgDeposit{}, "vaults/MsgDeposit")
	legacy.RegisterAminoMsg(cdc, &MsgWithdraw{}, "vaults/MsgWithdraw")
	legacy.RegisterAminoMsg(cdc, &MsgPerformAction{}, "vaults/MsgPerformAction")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateVaultCoins{}, "vaults/MsgUpdateVaultCoins")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateVaultFees{}, "vaults/MsgUpdateVaultFees")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateVaultLockupPeriod{}, "vaults/MsgUpdateVaultLockupPeriod")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateVaultMaxAmountUsd{}, "vaults/MsgUpdateVaultMaxAmountUsd")

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgAddVault{},
		&MsgDeposit{},
		&MsgWithdraw{},
		&MsgPerformAction{},
		&MsgUpdateVaultCoins{},
		&MsgUpdateVaultFees{},
		&MsgUpdateVaultLockupPeriod{},
		&MsgUpdateVaultMaxAmountUsd{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
