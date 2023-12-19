package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgFeedMultipleExternalLiquidity = "feed_multiple_external_liquidity"

var _ sdk.Msg = &MsgFeedMultipleExternalLiquidity{}

func NewMsgFeedMultipleExternalLiquidity(sender string) *MsgFeedMultipleExternalLiquidity {
	return &MsgFeedMultipleExternalLiquidity{
		Sender: sender,
	}
}

func (msg *MsgFeedMultipleExternalLiquidity) Route() string {
	return RouterKey
}

func (msg *MsgFeedMultipleExternalLiquidity) Type() string {
	return TypeMsgSwapExactAmountOut
}

func (msg *MsgFeedMultipleExternalLiquidity) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgFeedMultipleExternalLiquidity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFeedMultipleExternalLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
