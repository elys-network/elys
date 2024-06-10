package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreatePortfolio = "create_portfolio"

var _ sdk.Msg = &MsgCreatePortfolio{}

func NewMsgCreatePortfolio(creator string, user string) *MsgCreatePortfolio {
	return &MsgCreatePortfolio{
		Creator: creator,
		User:    user,
	}
}

func (msg *MsgCreatePortfolio) Route() string {
	return RouterKey
}

func (msg *MsgCreatePortfolio) Type() string {
	return TypeMsgCreatePortfolio
}

func (msg *MsgCreatePortfolio) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePortfolio) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePortfolio) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
