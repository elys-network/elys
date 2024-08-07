package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddEntry = "add_entry"

var _ sdk.Msg = &MsgAddEntry{}

func NewMsgAddEntry(creator string, baseDenom string, decimals uint64, denom string, path string, ibcChannelId string, ibcCounterpartyChannelId string, displayName string, displaySymbol string, network string, address string, externalSymbol string, transferLimit string, permissions []string, unitDenom string, ibcCounterpartyDenom string, ibcCounterpartyChainId string, commitEnabled bool, withdrawEnabled bool) *MsgAddEntry {
	return &MsgAddEntry{
		Creator:                  creator,
		BaseDenom:                baseDenom,
		Decimals:                 decimals,
		Denom:                    denom,
		Path:                     path,
		IbcChannelId:             ibcChannelId,
		IbcCounterpartyChannelId: ibcCounterpartyChannelId,
		DisplayName:              displayName,
		DisplaySymbol:            displaySymbol,
		Network:                  network,
		Address:                  address,
		ExternalSymbol:           externalSymbol,
		TransferLimit:            transferLimit,
		Permissions:              permissions,
		UnitDenom:                unitDenom,
		IbcCounterpartyDenom:     ibcCounterpartyDenom,
		IbcCounterpartyChainId:   ibcCounterpartyChainId,
		CommitEnabled:            commitEnabled,
		WithdrawEnabled:          withdrawEnabled,
	}
}

func (msg *MsgAddEntry) Route() string {
	return RouterKey
}

func (msg *MsgAddEntry) Type() string {
	return TypeMsgAddEntry
}

func (msg *MsgAddEntry) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddEntry) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddEntry) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
