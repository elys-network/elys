package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateEntry = "create_entry"
	TypeMsgUpdateEntry = "update_entry"
	TypeMsgDeleteEntry = "delete_entry"
)

var _ sdk.Msg = &MsgCreateEntry{}

func NewMsgCreateEntry(
	authority string,
	baseDenom string,
	decimals uint64,
	denom string,
	path string,
	ibcChannelId string,
	ibcCounterpartyChannelId string,
	displayName string,
	displaySymbol string,
	network string,
	address string,
	externalSymbol string,
	transferLimit string,
	permissions []string,
	unitDenom string,
	ibcCounterpartyDenom string,
	ibcCounterpartyChainId string,
) *MsgCreateEntry {
	return &MsgCreateEntry{
		Authority:                authority,
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
	}
}

func (msg *MsgCreateEntry) Route() string {
	return RouterKey
}

func (msg *MsgCreateEntry) Type() string {
	return TypeMsgCreateEntry
}

func (msg *MsgCreateEntry) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateEntry) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateEntry) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateEntry{}

func NewMsgUpdateEntry(
	authority string,
	baseDenom string,
	decimals uint64,
	denom string,
	path string,
	ibcChannelId string,
	ibcCounterpartyChannelId string,
	displayName string,
	displaySymbol string,
	network string,
	address string,
	externalSymbol string,
	transferLimit string,
	permissions []string,
	unitDenom string,
	ibcCounterpartyDenom string,
	ibcCounterpartyChainId string,
) *MsgUpdateEntry {
	return &MsgUpdateEntry{
		Authority:                authority,
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
	}
}

func (msg *MsgUpdateEntry) Route() string {
	return RouterKey
}

func (msg *MsgUpdateEntry) Type() string {
	return TypeMsgUpdateEntry
}

func (msg *MsgUpdateEntry) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateEntry) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateEntry) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteEntry{}

func NewMsgDeleteEntry(
	authority string,
	baseDenom string,
) *MsgDeleteEntry {
	return &MsgDeleteEntry{
		Authority: authority,
		BaseDenom: baseDenom,
	}
}

func (msg *MsgDeleteEntry) Route() string {
	return RouterKey
}

func (msg *MsgDeleteEntry) Type() string {
	return TypeMsgDeleteEntry
}

func (msg *MsgDeleteEntry) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgDeleteEntry) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteEntry) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
