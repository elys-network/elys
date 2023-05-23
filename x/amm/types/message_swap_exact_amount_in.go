package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSwapExactAmountIn = "swap_exact_amount_in"

var _ sdk.Msg = &MsgSwapExactAmountIn{}

func NewMsgSwapExactAmountIn(creator string, tokenIn sdk.Coin, tokenOutMinAmount sdk.Uint, swapRoutePoolIds []uint64, swapRouteDenoms []string) *MsgSwapExactAmountIn {
	return &MsgSwapExactAmountIn{
		Creator:           creator,
		TokenIn:           tokenIn,
		TokenOutMinAmount: tokenOutMinAmount,
		SwapRoutePoolIds:  swapRoutePoolIds,
		SwapRouteDenoms:   swapRouteDenoms,
	}
}

func (msg *MsgSwapExactAmountIn) Route() string {
	return RouterKey
}

func (msg *MsgSwapExactAmountIn) Type() string {
	return TypeMsgSwapExactAmountIn
}

func (msg *MsgSwapExactAmountIn) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSwapExactAmountIn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwapExactAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
