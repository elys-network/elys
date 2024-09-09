package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgVestLiquid = "vest_liquid"

var _ sdk.Msg = &MsgVestLiquid{}

func NewMsgVestLiquid(creator string, amount math.Int, denom string) *MsgVestLiquid {
	return &MsgVestLiquid{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgVestLiquid) Route() string {
	return RouterKey
}

func (msg *MsgVestLiquid) Type() string {
	return TypeMsgVestLiquid
}

func (msg *MsgVestLiquid) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVestLiquid) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVestLiquid) ValidateBasic() error {
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
