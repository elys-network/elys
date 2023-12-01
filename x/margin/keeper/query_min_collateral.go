package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) MinCollateral(goCtx context.Context, req *types.QueryMinCollateralRequest) (*types.QueryMinCollateralResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	collateralAmount, err := k.CalcMinCollateral(ctx, req.Position, req.Leverage, req.TradingAsset, req.CollateralAsset)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrCalcMinCollateral, "error calculating min collateral: %s", err.Error())
	}

	swapFee := k.GetSwapFee(ctx)

	// Apply discount to swap fee if applicable
	swapFee, discount, err := k.ApplyDiscount(ctx, swapFee, req.Discount, k.GetBrokerAddress(ctx).String())
	if err != nil {
		return nil, err
	}

	return &types.QueryMinCollateralResponse{
		Position:      req.Position,
		Leverage:      req.Leverage,
		TradingAsset:  req.TradingAsset,
		MinCollateral: sdk.NewDecCoinFromDec(req.CollateralAsset, collateralAmount),
		SwapFee:       swapFee,
		Discount:      discount,
	}, nil
}
