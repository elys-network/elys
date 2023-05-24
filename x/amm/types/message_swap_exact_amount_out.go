package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSwapExactAmountOut = "swap_exact_amount_out"

var _ sdk.Msg = &MsgSwapExactAmountOut{}

func NewMsgSwapExactAmountOut(creator string, tokenOut sdk.Coin, tokenOutMaxAmount sdk.Uint, swapRoutePoolIds []uint64, swapRouteDenoms []string) *MsgSwapExactAmountOut {
	return &MsgSwapExactAmountOut{
		Creator:           creator,
		TokenOut:          tokenOut,
		TokenOutMaxAmount: tokenOutMaxAmount,
		SwapRoutePoolIds:  swapRoutePoolIds,
		SwapRouteDenoms:   swapRouteDenoms,
	}
}

func (msg *MsgSwapExactAmountOut) Route() string {
	return RouterKey
}

func (msg *MsgSwapExactAmountOut) Type() string {
	return TypeMsgSwapExactAmountOut
}

func (msg *MsgSwapExactAmountOut) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSwapExactAmountOut) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwapExactAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
