package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgVestNow = "vest_now"

var _ sdk.Msg = &MsgVestNow{}

func NewMsgVestNow(creator string, amount math.Int, denom string) *MsgVestNow {
	return &MsgVestNow{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgVestNow) Route() string {
	return RouterKey
}

func (msg *MsgVestNow) Type() string {
	return TypeMsgVestNow
}

func (msg *MsgVestNow) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVestNow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVestNow) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Amount.IsNil() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount can not be nil")
	}

	if msg.Amount.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount can not be negative")
	}

	return nil
}
