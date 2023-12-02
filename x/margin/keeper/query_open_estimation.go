package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OpenEstimation(goCtx context.Context, req *types.QueryOpenEstimationRequest) (*types.QueryOpenEstimationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// calculate min collateral
	minCollateral, err := k.CalcMinCollateral(ctx, req.Leverage)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrCalcMinCollateral, "error calculating min collateral: %s", err.Error())
	}

	// get swap fee param
	swapFee := k.GetSwapFee(ctx)

	entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	leveragedAmount := sdk.NewDecFromBigInt(req.Collateral.Amount.BigInt()).Mul(req.Leverage).TruncateInt()
	leveragedCoin := sdk.NewCoin(req.Collateral.Denom, leveragedAmount)

	_, _, positionSize, openPrice, swapFee, discount, availableLiquidity, err := k.amm.CalcSwapEstimationByDenom(ctx, leveragedCoin, req.Collateral.Denom, req.TradingAsset, baseCurrency, req.Discount, swapFee)
	if err != nil {
		return nil, err
	}

	// calculate estimated pnl
	// estimated_pnl = leveraged_amount * (open_price - take_profit_price)
	estimatedPnL := sdk.NewDecFromBigInt(leveragedAmount.BigInt()).Mul(req.TakeProfitPrice.Sub(openPrice)).TruncateInt()

	if leveragedAmount.IsZero() {
		return nil, sdkerrors.Wrapf(types.ErrAmountTooLow, "leveraged amount is zero")
	}

	// calculate liquidation price
	// liquidation_price = -collateral_amount / leveraged_amount_value + open_price_value
	liquidationPrice := sdk.NewDecFromBigInt(req.Collateral.Amount.Neg().BigInt()).Quo(sdk.NewDecFromBigInt(leveragedAmount.BigInt())).Add(openPrice)

	return &types.QueryOpenEstimationResponse{
		Position:           req.Position,
		Leverage:           req.Leverage,
		TradingAsset:       req.TradingAsset,
		Collateral:         req.Collateral,
		MinCollateral:      sdk.NewCoin(req.Collateral.Denom, minCollateral),
		ValidCollateral:    req.Collateral.Amount.GTE(minCollateral),
		PositionSize:       positionSize,
		SwapFee:            swapFee,
		Discount:           discount,
		OpenPrice:          openPrice,
		TakeProfitPrice:    req.TakeProfitPrice,
		LiquidationPrice:   liquidationPrice,
		EstimatedPnl:       estimatedPnL,
		AvailableLiquidity: availableLiquidity,
	}, nil
}
