package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgBuyElys        = "BuyElys"
	TypeMsgReturnElys     = "ReturnElys"
	TypeMsgWithdrawRaised = "WithdrawRaised"
)

var (
	_ sdk.Msg = &MsgBuyElys{}
	_ sdk.Msg = &MsgReturnElys{}
	_ sdk.Msg = &MsgWithdrawRaised{}
)

func NewMsgBuyElys(sender string, spendingToken string, tokenAmount sdk.Int) *MsgBuyElys {
	return &MsgBuyElys{
		Sender:        sender,
		SpendingToken: spendingToken,
		TokenAmount:   tokenAmount,
	}
}

func (msg *MsgBuyElys) Route() string {
	return RouterKey
}

func (msg *MsgBuyElys) Type() string {
	return TypeMsgBuyElys
}

func (msg *MsgBuyElys) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgBuyElys) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBuyElys) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

func NewMsgReturnElys(sender string, orderId uint64, returnElysAmount math.Int) *MsgReturnElys {
	return &MsgReturnElys{
		Sender:           sender,
		OrderId:          orderId,
		ReturnElysAmount: returnElysAmount,
	}
}

func (msg *MsgReturnElys) Route() string {
	return RouterKey
}

func (msg *MsgReturnElys) Type() string {
	return TypeMsgReturnElys
}

func (msg *MsgReturnElys) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgReturnElys) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgReturnElys) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

func NewMsgWithdrawRaised(sender string, coins sdk.Coins) *MsgWithdrawRaised {
	return &MsgWithdrawRaised{
		Sender: sender,
		Coins:  coins,
	}
}

func (msg *MsgWithdrawRaised) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawRaised) Type() string {
	return TypeMsgWithdrawRaised
}

func (msg *MsgWithdrawRaised) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawRaised) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawRaised) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
