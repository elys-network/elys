package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateStopLoss = "update_stop_loss"

var _ sdk.Msg = &MsgUpdateStopLoss{}

func NewMsgUpdateStopLoss(creator string, position uint64, price sdk.Dec) *MsgUpdateStopLoss {
	return &MsgUpdateStopLoss{
		Creator:  creator,
		Position: position,
		Price:    price,
	}
}

func (msg *MsgUpdateStopLoss) Route() string {
	return RouterKey
}

func (msg *MsgUpdateStopLoss) Type() string {
	return TypeMsgUpdateStopLoss
}

func (msg *MsgUpdateStopLoss) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateStopLoss) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateStopLoss) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
