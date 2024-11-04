package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgExecuteOrders = "execute_orders"

var _ sdk.Msg = &MsgExecuteOrders{}

func NewMsgExecuteOrders(creator string, spotOrderIds []uint64, perpetualOrderIds []uint64) *MsgExecuteOrders {
	return &MsgExecuteOrders{
		Creator:           creator,
		SpotOrderIds:      spotOrderIds,
		PerpetualOrderIds: perpetualOrderIds,
	}
}

func (msg *MsgExecuteOrders) Route() string {
	return RouterKey
}

func (msg *MsgExecuteOrders) Type() string {
	return TypeMsgExecuteOrders
}

func (msg *MsgExecuteOrders) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgExecuteOrders) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgExecuteOrders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	for _, id := range msg.SpotOrderIds {
		if id == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "spot order ID cannot be zero")
		}
	}

	for _, id := range msg.PerpetualOrderIds {
		if id == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "perpetual order ID cannot be zero")
		}
	}

	return nil
}
