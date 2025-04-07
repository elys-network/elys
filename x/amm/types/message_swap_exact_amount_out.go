package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
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
	seenDenoms := make(map[string]bool)
	for i, route := range msg.Routes {
		if err = sdk.ValidateDenom(route.TokenInDenom); err != nil {
			return err
		}

		// Ensure no route has the same input and output denomination
		if i < (len(msg.Routes)-1) && msg.Routes[i+1].TokenInDenom == route.TokenInDenom {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "route %d has the same input and output denom as the previous route", i)
		}

		// Ensure all TokenInDenom values are unique
		if seenDenoms[route.TokenInDenom] {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "duplicate TokenInDenom found in route %d", i)
		}
		seenDenoms[route.TokenInDenom] = true
	}
	if err = msg.TokenOut.Validate(); err != nil {
		return err
	}
	if msg.TokenOut.IsZero() {
		return errors.New("token in is zero")
	}

	// Ensure no circular swaps
	if len(msg.Routes) > 0 && msg.TokenOut.Denom == msg.Routes[0].TokenInDenom {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "circular swap detected: token in denom matches the last route's token out denom")
	}

	return nil
}
