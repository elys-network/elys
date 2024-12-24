package types

import (
	errorsmod "cosmossdk.io/errors"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSwapByDenom{}

func NewMsgSwapByDenom(sender, recipient string, amount sdk.Coin, minAmount sdk.Coin, maxAmount sdk.Coin, denomIn string, denomOut string) *MsgSwapByDenom {
	return &MsgSwapByDenom{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		MinAmount: minAmount,
		MaxAmount: maxAmount,
		DenomIn:   denomIn,
		DenomOut:  denomOut,
	}
}

func (msg *MsgSwapByDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.Recipient != "" {
		if _, err = sdk.AccAddressFromBech32(msg.Recipient); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
		}
	}
	if err = msg.Amount.Validate(); err != nil {
		return err
	}
	if err = sdk.ValidateDenom(msg.DenomIn); err != nil {
		return err
	}
	if err = sdk.ValidateDenom(msg.DenomOut); err != nil {
		return err
	}
	if msg.DenomIn == msg.DenomOut {
		return errors.New("demom in cannot be same as denom out")
	}
	return nil
}
