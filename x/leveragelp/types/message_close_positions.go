package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClosePositions = "close_positions"

var _ sdk.Msg = &MsgClosePositions{}

func NewMsgClosePositions(creator string, liquidate []*PositionRequest, stoploss []*PositionRequest) *MsgClosePositions {
	return &MsgClosePositions{
		Creator:   creator,
		Liquidate: liquidate,
		Stoploss:  stoploss,
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
