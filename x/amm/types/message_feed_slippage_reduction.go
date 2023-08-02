package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgFeedSlippageReduction = "feed_slippage_reduction"

var _ sdk.Msg = &MsgFeedSlippageReduction{}

func NewMsgFeedSlippageReduction(sender string, poolId uint64, reduction sdk.Dec) *MsgFeedSlippageReduction {
	return &MsgFeedSlippageReduction{
		Sender:            sender,
		PoolId:            poolId,
		SlippageReduction: reduction,
	}
}

func (msg *MsgFeedSlippageReduction) Route() string {
	return RouterKey
}

func (msg *MsgFeedSlippageReduction) Type() string {
	return TypeMsgFeedSlippageReduction
}

func (msg *MsgFeedSlippageReduction) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgFeedSlippageReduction) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFeedSlippageReduction) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid Sender address (%s)", err)
	}

	if msg.SlippageReduction.IsNegative() || msg.SlippageReduction.GT(sdk.OneDec()) {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid slippage reduction parameter (%s)", err)
	}
	return nil
}
