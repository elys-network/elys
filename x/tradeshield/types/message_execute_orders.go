package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
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

	if len(msg.SpotOrderIds) == 0 && len(msg.PerpetualOrderIds) == 0 {
		return fmt.Errorf("SpotOrderIds and PerpetualOrderIds both are empty")
	}
	for _, id := range msg.SpotOrderIds {
		if id == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "spot order ID cannot be zero")
		}
	}

	for _, id := range msg.PerpetualOrderIds {
		if id == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "perpetual order ID cannot be zero")
		}
	}

	return nil
}
