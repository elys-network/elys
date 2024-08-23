package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClosePositions = "close_positions"

var _ sdk.Msg = &MsgClosePositions{}

func NewMsgClosePositions(creator sdk.AccAddress, liquidate []*PositionRequest, stoploss []*PositionRequest) *MsgClosePositions {
	return &MsgClosePositions{
		Creator:   creator.String(),
		Liquidate: liquidate,
		StopLoss:  stoploss,
	}
}

func (msg *MsgClosePositions) Route() string {
	return RouterKey
}

func (msg *MsgClosePositions) Type() string {
	return TypeMsgClosePositions
}

func (msg *MsgClosePositions) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClosePositions) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClosePositions) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(msg.Liquidate)+len(msg.StopLoss) == 0 {
		return fmt.Errorf("no liquidate or stoploss position requests")
	}
	positionRequests := make(map[uint64]bool)
	for _, liquidation := range msg.Liquidate {
		_, err = sdk.AccAddressFromBech32(liquidation.Address)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid liquidation address (%s)", err)
		}
		if positionRequests[liquidation.Id] {
			return fmt.Errorf("repeated liquidation id (%d)", liquidation.Id)
		} else {
			positionRequests[liquidation.Id] = true
		}
	}
	for _, stoploss := range msg.StopLoss {
		_, err = sdk.AccAddressFromBech32(stoploss.Address)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid stoploss address (%s)", err)
		}
		if positionRequests[stoploss.Id] {
			return fmt.Errorf("repeated stoploss id (%d)", stoploss.Id)
		} else {
			positionRequests[stoploss.Id] = true
		}
	}
	return nil
}
