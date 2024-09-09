package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgBond = "stake"

var _ sdk.Msg = &MsgBond{}

func NewMsgBond(creator string, amount math.Int) *MsgBond {
	return &MsgBond{
		Creator: creator,
		Amount:  amount,
	}
}

func (msg *MsgBond) Route() string {
	return RouterKey
}

func (msg *MsgBond) Type() string {
	return TypeMsgBond
}

func (msg *MsgBond) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBond) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if !msg.Amount.IsPositive() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Bond amount should be positive: "+msg.Amount.String())
	}
	return nil
}
