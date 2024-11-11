package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgBrokerClose(creator string, id uint64, amount math.Int, owner string) *MsgBrokerClose {
	return &MsgBrokerClose{
		Creator: creator,
		Id:      id,
		Amount:  amount,
		Owner:   owner,
	}
}

func (msg *MsgBrokerClose) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if msg.Amount.IsNil() {
		return ErrInvalidAmount
	}

	if msg.Amount.IsNegative() {
		return ErrInvalidAmount
	}

	return nil
}
