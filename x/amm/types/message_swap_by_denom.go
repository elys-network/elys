package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSwapByDenom = "swap_by_denom"

var _ sdk.Msg = &MsgSwapByDenom{}

func NewMsgSwapByDenom(sender, recipient string, amount sdk.Coin, minAmount sdk.Coin, maxAmount sdk.Coin, denomIn string, denomOut string, discount sdk.Dec) *MsgSwapByDenom {
	return &MsgSwapByDenom{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		MinAmount: minAmount,
		MaxAmount: maxAmount,
		DenomIn:   denomIn,
		DenomOut:  denomOut,
		Discount:  discount,
	}
}

func (msg *MsgSwapByDenom) Route() string {
	return RouterKey
}

func (msg *MsgSwapByDenom) Type() string {
	return TypeMsgSwapByDenom
}

func (msg *MsgSwapByDenom) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgSwapByDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwapByDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
