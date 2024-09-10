package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateStopLoss = "update_stop_loss"

var _ sdk.Msg = &MsgUpdateStopLoss{}

func NewMsgUpdateStopLoss(creator string, id uint64, price sdk.Dec) *MsgUpdateStopLoss {
	return &MsgUpdateStopLoss{
		Creator: creator,
		Id:      id,
		Price:   price,
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

	if msg.Price.IsNegative() {
		return fmt.Errorf("stop loss price cannot be negative")
	}
	return nil
}
