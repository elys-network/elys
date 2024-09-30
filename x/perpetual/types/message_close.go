package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgClose{}

func NewMsgClose(creator string, id uint64, amount math.Int) *MsgClose {
	return &MsgClose{
		Creator: creator,
		Id:      id,
		Amount:  amount,
	}
}

func (msg *MsgClose) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Amount.IsNil() {
		return ErrInvalidAmount
	}

	if msg.Amount.IsNegative() {
		return ErrInvalidAmount
	}
	return nil
}
