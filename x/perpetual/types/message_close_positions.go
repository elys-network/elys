package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgClosePositions{}

func NewMsgClosePositions(creator string, liquidate []PositionRequest, stopLoss []PositionRequest, takeProfit []PositionRequest) *MsgClosePositions {
	return &MsgClosePositions{
		Creator:    creator,
		Liquidate:  liquidate,
		StopLoss:   stopLoss,
		TakeProfit: takeProfit,
	}
}
func (msg *MsgClosePositions) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Liquidate) == 0 && len(msg.StopLoss) == 0 && len(msg.TakeProfit) == 0 {
		return errors.New("liquidate, stop loss, take profit all are empty")
	}

	for _, position := range msg.Liquidate {
		_, err = sdk.AccAddressFromBech32(position.Address)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s) for id %d", err.Error(), position.Id)
		}
	}
	for _, position := range msg.StopLoss {
		_, err = sdk.AccAddressFromBech32(position.Address)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s) for id %d", err.Error(), position.Id)
		}
	}
	for _, position := range msg.TakeProfit {
		_, err = sdk.AccAddressFromBech32(position.Address)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address (%s) for id %d", err.Error(), position.Id)
		}
	}
	return nil
}
