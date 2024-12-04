package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateEnableVestNow{}

func NewMsgUpdateEnableVestNow(govAddress string, enableVestNow bool) MsgUpdateEnableVestNow {
	return MsgUpdateEnableVestNow{
		Authority:     govAddress,
		EnableVestNow: enableVestNow,
	}

}

func (msg *MsgUpdateEnableVestNow) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	return nil
}
