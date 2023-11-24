package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k msgServer) SwapByDenom(goCtx context.Context, msg *types.MsgSwapByDenom) (*types.MsgSwapByDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	inRoute, outRoute, _, spotPrice, err := k.CalcSwapEstimationByDenom(ctx, msg.Amount, msg.DenomIn, msg.DenomOut, baseCurrency)
	if err != nil {
		return nil, err
	}

	// swap to token out with exact amount in using in route
	if inRoute != nil {
		// check min amount denom is equals to denom out
		if msg.MinAmount.Denom != msg.DenomOut {
			return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "min amount denom %s is not equals to denom out %s", msg.MinAmount.Denom, msg.DenomOut)
		}

		// convert route []*types.SwapAmountInRoute to []types.SwapAmountInRoute
		route := make([]types.SwapAmountInRoute, len(inRoute))
		for i, r := range inRoute {
			route[i] = *r
		}

		res, err := k.SwapExactAmountIn(
			ctx,
			&types.MsgSwapExactAmountIn{
				Sender:            msg.Sender,
				Routes:            route,
				TokenIn:           msg.Amount,
				TokenOutMinAmount: msg.MinAmount.Amount,
				Discount:          msg.Discount,
			},
		)
		if err != nil {
			return nil, err
		}

		return &types.MsgSwapByDenomResponse{
			Amount:    sdk.NewCoin(msg.DenomOut, res.TokenOutAmount),
			InRoute:   inRoute,
			OutRoute:  nil,
			SpotPrice: spotPrice,
			Discount:  msg.Discount,
		}, nil
	}

	// swap to token in with exact amount out using out route
	if outRoute != nil {
		// check max amount denom is equals to denom out
		if msg.MaxAmount.Denom != msg.DenomOut {
			return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "max amount denom %s is not equals to denom out %s", msg.MaxAmount.Denom, msg.DenomOut)
		}

		// convert route []*types.SwapAmountOutRoute to []types.SwapAmountOutRoute
		route := make([]types.SwapAmountOutRoute, len(outRoute))
		for i, r := range outRoute {
			route[i] = *r
		}

		res, err := k.SwapExactAmountOut(
			ctx,
			&types.MsgSwapExactAmountOut{
				Sender:           msg.Sender,
				Routes:           route,
				TokenInMaxAmount: msg.MaxAmount.Amount,
				TokenOut:         msg.Amount,
				Discount:         msg.Discount,
			},
		)
		if err != nil {
			return nil, err
		}

		return &types.MsgSwapByDenomResponse{
			Amount:    sdk.NewCoin(msg.DenomOut, res.TokenInAmount),
			InRoute:   nil,
			OutRoute:  outRoute,
			SpotPrice: spotPrice,
			Discount:  msg.Discount,
		}, nil
	}

	// otherwise throw an error
	return nil, sdkerrors.Wrapf(types.ErrInvalidSwapMsgType, "neither inRoute nor outRoute are available")
}
