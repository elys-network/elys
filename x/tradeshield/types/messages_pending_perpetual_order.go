package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreatePendingPerpetualOrder{}

func NewMsgCreatePendingPerpetualOrder(creator string, order string) *MsgCreatePendingPerpetualOrder {
	return &MsgCreatePendingPerpetualOrder{
		Creator: creator,
		Order:   order,
	}
}

func (msg *MsgCreatePendingPerpetualOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePendingPerpetualOrder{}

func NewMsgUpdatePendingPerpetualOrder(creator string, id uint64, order string) *MsgUpdatePendingPerpetualOrder {
	return &MsgUpdatePendingPerpetualOrder{
		Id:      id,
		Creator: creator,
		Order:   order,
	}
}

func (msg *MsgUpdatePendingPerpetualOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeletePendingPerpetualOrder{}

func NewMsgDeletePendingPerpetualOrder(creator string, id uint64) *MsgDeletePendingPerpetualOrder {
	return &MsgDeletePendingPerpetualOrder{
		Id:      id,
		Creator: creator,
	}
}

func (msg *MsgDeletePendingPerpetualOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
