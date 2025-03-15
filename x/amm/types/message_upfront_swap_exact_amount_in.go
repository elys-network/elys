package types

import (
	"errors"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpFrontSwapExactAmountIn{}

func NewMsgUpFrontSwapExactAmountIn(sender string, tokenIn sdk.Coin, tokenOutMinAmount sdkmath.Int, swapRoutePoolIds []uint64, swapRouteDenoms []string) *MsgUpFrontSwapExactAmountIn {
	if len(swapRoutePoolIds) != len(swapRouteDenoms) {
		return nil // or raise an error as the input lists should have the same length
	}

	var routes []SwapAmountInRoute
	for i := range swapRoutePoolIds {
		route := SwapAmountInRoute{
			PoolId:        swapRoutePoolIds[i],
			TokenOutDenom: swapRouteDenoms[i],
		}
		routes = append(routes, route)
	}

	return &MsgUpFrontSwapExactAmountIn{
		Sender:            sender,
		Routes:            routes,
		TokenIn:           tokenIn,
		TokenOutMinAmount: tokenOutMinAmount,
	}
}

func (msg *MsgUpFrontSwapExactAmountIn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	for _, route := range msg.Routes {
		if err = sdk.ValidateDenom(route.TokenOutDenom); err != nil {
			return err
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
