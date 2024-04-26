package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddExternalIncentive = "add_external_incentive"

var _ sdk.Msg = &MsgAddExternalIncentive{}

func (msg *MsgAddExternalIncentive) Route() string {
	return RouterKey
}

func (msg *MsgAddExternalIncentive) Type() string {
	return TypeMsgAddExternalIncentive
}

func (msg *MsgAddExternalIncentive) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgAddExternalIncentive) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddExternalIncentive) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}
