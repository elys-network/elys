package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePendingSpotOrder = "create_pending_spot_order"
	TypeMsgUpdatePendingSpotOrder = "update_pending_spot_order"
	TypeMsgDeletePendingSpotOrder = "delete_pending_spot_order"
)

var _ sdk.Msg = &MsgCreatePendingSpotOrder{}

// func NewMsgCreatePendingSpotOrder(ownerAddress string, orderType SpotOrderType,
// 	orderPrice OrderPrice, orderAmount sdk.Coin,
// 	orderTargetDenom string, status Status, date Date) *MsgCreatePendingSpotOrder {
// 	return &MsgCreatePendingSpotOrder{
// 		OrderType:        orderType,
// 		OrderPrice:       &orderPrice,
// 		OrderAmount:      &orderAmount,
// 		OwnerAddress:     ownerAddress,
// 		OrderTargetDenom: orderTargetDenom,
// 		Status:           status,
// 		Date:             &date,
// 	}
// }

func NewMsgCreatePendingSpotOrder(ownerAddress string) *MsgCreatePendingSpotOrder {
	return &MsgCreatePendingSpotOrder{
		OwnerAddress: ownerAddress,
	}
}

func (msg *MsgCreatePendingSpotOrder) Route() string {
	return RouterKey
}

func (msg *MsgCreatePendingSpotOrder) Type() string {
	return TypeMsgCreatePendingSpotOrder
}

func (msg *MsgCreatePendingSpotOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePendingSpotOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePendingSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePendingSpotOrder{}

func NewMsgUpdatePendingSpotOrder(creator string, id uint64, orderPrice *OrderPrice) *MsgUpdatePendingSpotOrder {
	return &MsgUpdatePendingSpotOrder{
		OrderId:      id,
		OwnerAddress: creator,
		OrderPrice:   orderPrice,
	}
}

func (msg *MsgUpdatePendingSpotOrder) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePendingSpotOrder) Type() string {
	return TypeMsgUpdatePendingSpotOrder
}

func (msg *MsgUpdatePendingSpotOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePendingSpotOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePendingSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeletePendingSpotOrder{}

func NewMsgDeletePendingSpotOrder(creator string, id uint64) *MsgDeletePendingSpotOrder {
	return &MsgDeletePendingSpotOrder{
		OrderId:      id,
		OwnerAddress: creator,
	}
}
func (msg *MsgDeletePendingSpotOrder) Route() string {
	return RouterKey
}

func (msg *MsgDeletePendingSpotOrder) Type() string {
	return TypeMsgDeletePendingSpotOrder
}

func (msg *MsgDeletePendingSpotOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeletePendingSpotOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeletePendingSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}
