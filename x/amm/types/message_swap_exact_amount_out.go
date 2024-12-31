package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSwapExactAmountOut{}

func NewMsgSwapExactAmountOut(sender, recipient string, tokenOut sdk.Coin, tokenInMaxAmount sdkmath.Int, swapRoutePoolIds []uint64, swapRouteDenoms []string) *MsgSwapExactAmountOut {
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
	}
}

func (msg *MsgSwapExactAmountOut) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.Recipient != "" {
		if _, err = sdk.AccAddressFromBech32(msg.Recipient); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
		}
	}
	for _, route := range msg.Routes {
		if err = sdk.ValidateDenom(route.TokenInDenom); err != nil {
			return err
		}
	}
	if err = msg.TokenOut.Validate(); err != nil {
		return err
	}
	if msg.TokenOut.IsZero() {
		return errors.New("token in is zero")
	}
	return nil
}
