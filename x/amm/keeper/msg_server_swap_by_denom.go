package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k msgServer) SwapByDenom(goCtx context.Context, msg *types.MsgSwapByDenom) (*types.MsgSwapByDenomResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return k.Keeper.SwapByDenom(ctx, msg)
}

func (k Keeper) SwapByDenom(ctx sdk.Context, msg *types.MsgSwapByDenom) (*types.MsgSwapByDenomResponse, error) {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// retrieve base currency denom
	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	inRoute, outRoute, _, spotPrice, _, _, _, slippage, weightBonus, _, err := k.CalcSwapEstimationByDenom(ctx, msg.Amount, msg.DenomIn, msg.DenomOut, baseCurrency, msg.Sender, osmomath.ZeroBigDec(), 0)
	if err != nil {
		return nil, err
	}

	// swap to token out with exact amount in using in route
	if inRoute != nil {
		// check min amount denom is equals to denom out
		if msg.MinAmount.Denom != msg.DenomOut {
			return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "min amount denom %s is not equals to denom out %s", msg.MinAmount.Denom, msg.DenomOut)
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
				Recipient:         msg.Recipient,
				Routes:            route,
				TokenIn:           msg.Amount,
				TokenOutMinAmount: msg.MinAmount.Amount,
			},
		)
		if err != nil {
			return nil, err
		}

		return &types.MsgSwapByDenomResponse{
			Amount:      sdk.NewCoin(msg.DenomOut, res.TokenOutAmount),
			InRoute:     inRoute,
			OutRoute:    nil,
			SpotPrice:   spotPrice.Dec(),
			SwapFee:     res.SwapFee,
			Discount:    res.Discount,
			Recipient:   res.Recipient,
			Slippage:    slippage.Dec(),
			WeightBonus: weightBonus.Dec(),
		}, nil
	}

	// swap to token in with exact amount out using out route
	if outRoute != nil {
		// check max amount denom is equals to denom out
		if msg.MaxAmount.Denom != msg.DenomOut {
			return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "max amount denom %s is not equals to denom out %s", msg.MaxAmount.Denom, msg.DenomOut)
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
			},
		)
		if err != nil {
			return nil, err
		}

		return &types.MsgSwapByDenomResponse{
			Amount:    sdk.NewCoin(msg.DenomIn, res.TokenInAmount),
			InRoute:   nil,
			OutRoute:  outRoute,
			SpotPrice: spotPrice.Dec(),
			SwapFee:   res.SwapFee,
			Discount:  res.Discount,
			Recipient: res.Recipient,
		}, nil
	}

	// otherwise throw an error
	return nil, errorsmod.Wrapf(types.ErrInvalidSwapMsgType, "neither inRoute nor outRoute are available")
}
