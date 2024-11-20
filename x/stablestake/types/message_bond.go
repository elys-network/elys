package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgBond{}

func NewMsgBond(creator string, amount math.Int) *MsgBond {
	return &MsgBond{
		Creator: creator,
		Amount:  amount,
	}
}

func (msg *MsgBond) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Amount.IsNil() {
		return fmt.Errorf("amount cannot be nil")
	}
	if !msg.Amount.IsPositive() {
		return fmt.Errorf("amount should be positive: " + msg.Amount.String())
	}
	return nil
}
