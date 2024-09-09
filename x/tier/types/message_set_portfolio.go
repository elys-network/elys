package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetPortfolio = "set_portfolio"

var _ sdk.Msg = &MsgSetPortfolio{}

func NewMsgSetPortfolio(creator string, user string) *MsgSetPortfolio {
	return &MsgSetPortfolio{
		Creator: creator,
		User:    user,
	}
}

func (msg *MsgSetPortfolio) Route() string {
	return RouterKey
}

func (msg *MsgSetPortfolio) Type() string {
	return TypeMsgSetPortfolio
}

func (msg *MsgSetPortfolio) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetPortfolio) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetPortfolio) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.User)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid user address (%s)", err)
	}
	return nil
}
