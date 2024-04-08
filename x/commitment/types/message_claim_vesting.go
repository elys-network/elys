package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimVesting = "claim_vesting"

var _ sdk.Msg = &MsgClaimVesting{}

func NewMsgClaimVesting(sender string) *MsgClaimVesting {
	return &MsgClaimVesting{
		Sender: sender,
	}
}

func (msg *MsgClaimVesting) Route() string {
	return RouterKey
}

func (msg *MsgClaimVesting) Type() string {
	return TypeMsgClaimVesting
}

func (msg *MsgClaimVesting) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimVesting) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimVesting) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
