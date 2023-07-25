package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgOpen = "open"

var _ sdk.Msg = &MsgOpen{}

func NewMsgOpen(creator string, collateralAsset string, collateralAmount sdk.Int, borrowAsset string, position Position, leverage sdk.Dec) *MsgOpen {
	return &MsgOpen{
		Creator:          creator,
		CollateralAsset:  collateralAsset,
		CollateralAmount: collateralAmount,
		BorrowAsset:      borrowAsset,
		Position:         position,
		Leverage:         leverage,
	}
}

func (msg *MsgOpen) Route() string {
	return RouterKey
}

func (msg *MsgOpen) Type() string {
	return TypeMsgOpen
}

func (msg *MsgOpen) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgOpen) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgOpen) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
