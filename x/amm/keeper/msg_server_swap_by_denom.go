package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	// retrieve base currency denom
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// retrieve denom in decimals
	entry, found = k.assetProfileKeeper.GetEntryByDenom(ctx, msg.DenomIn)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", msg.DenomIn)
	}
	decimals := entry.Decimals

	inRoute, outRoute, _, spotPrice, _, _, _, _, _, err := k.CalcSwapEstimationByDenom(ctx, msg.Amount, msg.DenomIn, msg.DenomOut, baseCurrency, msg.Discount, sdk.ZeroDec(), decimals)
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
			sdk.WrapSDKContext(ctx),
			&types.MsgSwapExactAmountIn{
				Sender:            msg.Sender,
				Recipient:         msg.Recipient,
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
			SwapFee:   res.SwapFee,
			Discount:  res.Discount,
			Recipient: res.Recipient,
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
			sdk.WrapSDKContext(ctx),
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
			SwapFee:   res.SwapFee,
			Discount:  res.Discount,
			Recipient: res.Recipient,
		}, nil
	}

	// otherwise throw an error
	return nil, errorsmod.Wrapf(types.ErrInvalidSwapMsgType, "neither inRoute nor outRoute are available")
}
