package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (msg *MsgCreatePendingSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
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

func (msg *MsgUpdatePendingSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
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

func (msg *MsgDeletePendingSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}
