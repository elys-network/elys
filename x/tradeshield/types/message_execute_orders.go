package types

import (
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgExecuteOrders{}

func NewMsgExecuteOrders(creator string, spotOrderIds []uint64, perpetualOrders []*PerpetualOrderKey) *MsgExecuteOrders {
	return &MsgExecuteOrders{
		Creator:         creator,
		SpotOrderIds:    spotOrderIds,
		PerpetualOrders: perpetualOrders,
	}
}

func (msg *MsgExecuteOrders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.SpotOrderIds) == 0 && len(msg.PerpetualOrders) == 0 {
		return fmt.Errorf("SpotOrderIds and PerpetualOrderIds both are empty")
	}
	if slices.Contains(msg.SpotOrderIds, 0) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "spot order ID cannot be zero")
	}

	for _, perpetualOrder := range msg.PerpetualOrders {
		if perpetualOrder.OwnerAddress == "" {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "perpetual order owner address cannot be empty")
		}
		if perpetualOrder.PoolId == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "perpetual order pool ID cannot be zero")
		}
		if perpetualOrder.OrderId == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "perpetual order ID cannot be zero")
		}
	}

	return nil
}
