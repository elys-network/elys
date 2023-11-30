package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUnstake = "unstake"

var _ sdk.Msg = &MsgUnstake{}

func NewMsgUnstake(address string, amount sdk.Int, asset string, validatorAddress string) *MsgUnstake {
	return &MsgUnstake{
		Creator:          address,
		Address:          address,
		Amount:           amount,
		Asset:            asset,
		ValidatorAddress: validatorAddress,
	}
}

func (msg *MsgUnstake) Route() string {
	return RouterKey
}

func (msg *MsgUnstake) Type() string {
	return TypeMsgUnstake
}

func (msg *MsgUnstake) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUnstake) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnstake) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
