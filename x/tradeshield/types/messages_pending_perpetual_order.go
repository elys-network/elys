package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewMsgCreatePendingPerpetualOrder(creator string) *MsgCreatePendingPerpetualOrder {
	return &MsgCreatePendingPerpetualOrder{
		OwnerAddress: creator,
	}
}

func (msg *MsgCreatePendingPerpetualOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePendingPerpetualOrder{}

func NewMsgUpdatePendingPerpetualOrder(creator string, id uint64, orderPrice *OrderPrice) *MsgUpdatePendingPerpetualOrder {
	return &MsgUpdatePendingPerpetualOrder{
		OrderId:      id,
		OwnerAddress: creator,
		OrderPrice:   orderPrice,
	}
}

func (msg *MsgUpdatePendingPerpetualOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCancelPerpetualOrders{}

func NewMsgCancelPerpetualOrders(creator string, ids []uint64) *MsgCancelPerpetualOrders {
	return &MsgCancelPerpetualOrders{
		OrderIds:     ids,
		OwnerAddress: creator,
	}
}

func (msg *MsgCancelPerpetualOrders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}
