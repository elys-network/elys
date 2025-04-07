package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSwapExactAmountIn{}

func NewMsgSwapExactAmountIn(sender, recipient string, tokenIn sdk.Coin, tokenOutMinAmount sdkmath.Int, swapRoutePoolIds []uint64, swapRouteDenoms []string) *MsgSwapExactAmountIn {
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
		Recipient:         recipient,
		Routes:            routes,
		TokenIn:           tokenIn,
		TokenOutMinAmount: tokenOutMinAmount,
	}
}

func (msg *MsgSwapExactAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.Recipient != "" {
		if _, err = sdk.AccAddressFromBech32(msg.Recipient); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
		}
	}
	for i, route := range msg.Routes {
		if err = sdk.ValidateDenom(route.TokenOutDenom); err != nil {
			return err
		}

		// Ensure no route has the same input and output denomination
		if i > 0 && msg.Routes[i-1].TokenOutDenom == route.TokenOutDenom {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "route %d has the same input and output denom as the previous route", i)
		}

		//  Ensure TokenIn denomination matches the first route
		if len(msg.Routes) > 0 && msg.TokenIn.Denom != msg.Routes[0].TokenOutDenom {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token in denom does not match the first route's input denom")
		}

		// Ensure no circular swaps
		if len(msg.Routes) > 0 && msg.TokenIn.Denom == msg.Routes[len(msg.Routes)-1].TokenOutDenom {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "circular swap detected: token in denom matches the last route's token out denom")
		}

		// Ensure all TokenOutDenom values are unique
		seenDenoms := make(map[string]bool)
		for i, route := range msg.Routes {
			if seenDenoms[route.TokenOutDenom] {
				return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "duplicate TokenOutDenom found in route %d", i)
			}
			seenDenoms[route.TokenOutDenom] = true
		}
	}
	if err = msg.TokenIn.Validate(); err != nil {
		return err
	}
	if msg.TokenIn.IsZero() {
		return errors.New("token in is zero")
	}
	return nil
}
