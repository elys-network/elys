package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSwapExactAmountIn = "swap_exact_amount_in"

var _ sdk.Msg = &MsgSwapExactAmountIn{}

func NewMsgSwapExactAmountIn(sender string, tokenIn sdk.Coin, tokenOutMinAmount math.Int, swapRoutePoolIds []uint64, swapRouteDenoms []string) *MsgSwapExactAmountIn {
	if len(swapRoutePoolIds) != len(swapRouteDenoms) {
		return nil // or raise an error as the input lists should have the same length
	}

	var routes []SwapAmountInRoute
	for i := 0; i < len(swapRoutePoolIds); i++ {
		route := SwapAmountInRoute{
			PoolId:        swapRoutePoolIds[i],
			TokenOutDenom: swapRouteDenoms[i],
		}
		routes = append(routes, route)
	}

	return &MsgSwapExactAmountIn{
		Sender:            sender,
		Routes:            routes,
		TokenIn:           tokenIn,
		TokenOutMinAmount: tokenOutMinAmount,
	}
}

func (msg *MsgSwapExactAmountIn) Route() string {
	return RouterKey
}

func (msg *MsgSwapExactAmountIn) Type() string {
	return TypeMsgSwapExactAmountIn
}

func (msg *MsgSwapExactAmountIn) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgSwapExactAmountIn) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwapExactAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
