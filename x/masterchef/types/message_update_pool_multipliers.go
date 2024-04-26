package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdatePoolMultipliers = "update_pool_multipliers"

var _ sdk.Msg = &MsgUpdatePoolMultipliers{}

func (msg *MsgUpdatePoolMultipliers) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePoolMultipliers) Type() string {
	return TypeMsgUpdatePoolMultipliers
}

func (msg *MsgUpdatePoolMultipliers) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgUpdatePoolMultipliers) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePoolMultipliers) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}
