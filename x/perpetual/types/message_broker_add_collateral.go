package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgBrokerAddCollateral = "broker_add_collateral"

var _ sdk.Msg = &MsgBrokerAddCollateral{}

func NewMsgBrokerAddCollateral(creator string, amount sdk.Int, id int32, owner string) *MsgBrokerAddCollateral {
	return &MsgBrokerAddCollateral{
		Creator: creator,
		Amount:  amount,
		Id:      id,
		Owner:   owner,
	}
}

func (msg *MsgBrokerAddCollateral) Route() string {
	return RouterKey
}

func (msg *MsgBrokerAddCollateral) Type() string {
	return TypeMsgBrokerAddCollateral
}

func (msg *MsgBrokerAddCollateral) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBrokerAddCollateral) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBrokerAddCollateral) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
