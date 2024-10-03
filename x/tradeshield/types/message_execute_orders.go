package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgExecuteOrders{}

func NewMsgExecuteOrders(creator string, spotOrderIds []uint64, perpetualOrderIds []uint64) *MsgExecuteOrders {
	return &MsgExecuteOrders{
		Creator:           creator,
		SpotOrderIds:      spotOrderIds,
		PerpetualOrderIds: perpetualOrderIds,
	}
}

func (msg *MsgExecuteOrders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
