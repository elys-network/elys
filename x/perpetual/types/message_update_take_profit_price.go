package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateTakeProfitPrice{}

func NewMsgUpdateTakeProfitPrice(creator string, id uint64, price sdkmath.LegacyDec) *MsgUpdateTakeProfitPrice {
	return &MsgUpdateTakeProfitPrice{
		Creator: creator,
		Id:      id,
		Price:   price,
	}
}

func (msg *MsgUpdateTakeProfitPrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if err = CheckLegacyDecNilAndNegative(msg.Price, "price"); err != nil {
		return err
	}
	return nil
}
