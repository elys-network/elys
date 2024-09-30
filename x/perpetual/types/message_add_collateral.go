package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddCollateral{}

func NewMsgAddCollateral(creator string, amount sdkmath.Int, id uint64) *MsgAddCollateral {
	return &MsgAddCollateral{
		Creator: creator,
		Amount:  amount,
		Id:      id,
	}
}

func (msg *MsgAddCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
