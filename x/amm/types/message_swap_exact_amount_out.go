package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSwapExactAmountOut = "swap_exact_amount_out"

var _ sdk.Msg = &MsgSwapExactAmountOut{}

func NewMsgSwapExactAmountOut(sender, recipient string, tokenOut sdk.Coin, tokenInMaxAmount math.Int, swapRoutePoolIds []uint64, swapRouteDenoms []string, discount sdk.Dec) *MsgSwapExactAmountOut {
	if len(swapRoutePoolIds) != len(swapRouteDenoms) {
		return nil // or raise an error as the input lists should have the same length
	}

	var routes []SwapAmountOutRoute
	for i := 0; i < len(swapRoutePoolIds); i++ {
		route := SwapAmountOutRoute{
			PoolId:       swapRoutePoolIds[i],
			TokenInDenom: swapRouteDenoms[i],
		}
		routes = append(routes, route)
	}

	return &MsgSwapExactAmountOut{
		Sender:           sender,
		Recipient:        recipient,
		Routes:           routes,
		TokenOut:         tokenOut,
		TokenInMaxAmount: tokenInMaxAmount,
		Discount:         discount,
	}
}

func (msg *MsgSwapExactAmountOut) Route() string {
	return RouterKey
}

func (msg *MsgSwapExactAmountOut) Type() string {
	return TypeMsgSwapExactAmountOut
}

func (msg *MsgSwapExactAmountOut) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgSwapExactAmountOut) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwapExactAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
