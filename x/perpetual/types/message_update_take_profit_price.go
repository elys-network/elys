package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateTakeProfitPrice = "update_take_profit_price"

var _ sdk.Msg = &MsgUpdateTakeProfitPrice{}

func NewMsgUpdateTakeProfitPrice(creator string, id uint64, price sdk.Dec) *MsgUpdateTakeProfitPrice {
	return &MsgUpdateTakeProfitPrice{
		Creator: creator,
		Id:      id,
		Price:   price,
	}
}

func (msg *MsgUpdateTakeProfitPrice) Route() string {
	return RouterKey
}

func (msg *MsgUpdateTakeProfitPrice) Type() string {
	return TypeMsgUpdateTakeProfitPrice
}

func (msg *MsgUpdateTakeProfitPrice) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateTakeProfitPrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateTakeProfitPrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Price.IsNegative() {
		return fmt.Errorf("take profit price cannot be negative")
	}
	return nil
}
