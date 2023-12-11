package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateTimeBasedInflation = "create_time_based_inflation"
	TypeMsgUpdateTimeBasedInflation = "update_time_based_inflation"
	TypeMsgDeleteTimeBasedInflation = "delete_time_based_inflation"
)

var _ sdk.Msg = &MsgCreateTimeBasedInflation{}

func NewMsgCreateTimeBasedInflation(
	authority string,
	startBlockHeight uint64,
	endBlockHeight uint64,
	description string,
	inflation *InflationEntry,
) *MsgCreateTimeBasedInflation {
	return &MsgCreateTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: startBlockHeight,
		EndBlockHeight:   endBlockHeight,
		Description:      description,
		Inflation:        inflation,
	}
}

func (msg *MsgCreateTimeBasedInflation) Route() string {
	return RouterKey
}

func (msg *MsgCreateTimeBasedInflation) Type() string {
	return TypeMsgCreateTimeBasedInflation
}

func (msg *MsgCreateTimeBasedInflation) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateTimeBasedInflation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateTimeBasedInflation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateTimeBasedInflation{}

func NewMsgUpdateTimeBasedInflation(
	authority string,
	startBlockHeight uint64,
	endBlockHeight uint64,
	description string,
	inflation *InflationEntry,
) *MsgUpdateTimeBasedInflation {
	return &MsgUpdateTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: startBlockHeight,
		EndBlockHeight:   endBlockHeight,
		Description:      description,
		Inflation:        inflation,
	}
}

func (msg *MsgUpdateTimeBasedInflation) Route() string {
	return RouterKey
}

func (msg *MsgUpdateTimeBasedInflation) Type() string {
	return TypeMsgUpdateTimeBasedInflation
}

func (msg *MsgUpdateTimeBasedInflation) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateTimeBasedInflation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateTimeBasedInflation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteTimeBasedInflation{}

func NewMsgDeleteTimeBasedInflation(
	authority string,
	startBlockHeight uint64,
	endBlockHeight uint64,
) *MsgDeleteTimeBasedInflation {
	return &MsgDeleteTimeBasedInflation{
		Authority:        authority,
		StartBlockHeight: startBlockHeight,
		EndBlockHeight:   endBlockHeight,
	}
}

func (msg *MsgDeleteTimeBasedInflation) Route() string {
	return RouterKey
}

func (msg *MsgDeleteTimeBasedInflation) Type() string {
	return TypeMsgDeleteTimeBasedInflation
}

func (msg *MsgDeleteTimeBasedInflation) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgDeleteTimeBasedInflation) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteTimeBasedInflation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
