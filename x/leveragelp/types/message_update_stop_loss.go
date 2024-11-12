package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateStopLoss{}

func NewMsgUpdateStopLoss(creator string, position uint64, price sdkmath.LegacyDec) *MsgUpdateStopLoss {
	return &MsgUpdateStopLoss{
		Creator:  creator,
		Position: position,
		Price:    price,
	}
}

func (msg *MsgUpdateStopLoss) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Price.IsNegative() {
		return fmt.Errorf("stop loss price cannot be negative")
	}
	return nil
}
