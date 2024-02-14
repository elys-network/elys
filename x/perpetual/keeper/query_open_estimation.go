package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
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
		return nil, errorsmod.Wrapf(types.ErrCalcMinCollateral, "error calculating min collateral: %s", err.Error())
	}

	// get swap fee param
	swapFee := k.GetSwapFee(ctx)

	// retrieve base currency denom
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// retrieve denom in decimals
	entry, found = k.assetProfileKeeper.GetEntryByDenom(ctx, req.Collateral.Denom)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", req.Collateral.Denom)
	}
	decimals := entry.Decimals

	leveragedAmount := sdk.NewDecFromBigInt(req.Collateral.Amount.BigInt()).Mul(req.Leverage).TruncateInt()
	leveragedCoin := sdk.NewCoin(req.Collateral.Denom, leveragedAmount)

	_, _, positionSize, openPrice, swapFee, discount, availableLiquidity, weightBonus, priceImpact, err := k.amm.CalcSwapEstimationByDenom(ctx, leveragedCoin, req.Collateral.Denom, req.TradingAsset, baseCurrency, req.Discount, swapFee, decimals)
	if err != nil {
		return nil, err
	}

	// invert openPrice
	openPrice = sdk.OneDec().Quo(openPrice)

	// calculate estimated pnl
	// estimated_pnl = leveraged_amount * (take_profit_price - open_price) - leveraged_amount
	estimatedPnL := sdk.NewDecFromBigInt(leveragedAmount.BigInt())
	estimatedPnL = estimatedPnL.Mul(req.TakeProfitPrice.Sub(openPrice))
	estimatedPnL = estimatedPnL.Sub(sdk.NewDecFromBigInt(leveragedAmount.BigInt()))
	estimatedPnLInt := estimatedPnL.TruncateInt()

	if leveragedAmount.IsZero() {
		return nil, errorsmod.Wrapf(types.ErrAmountTooLow, "leveraged amount is zero")
	}

	// calculate liquidation price
	// liquidation_price = -collateral_amount / leveraged_amount_value + open_price_value
	liquidationPrice := sdk.NewDecFromBigInt(req.Collateral.Amount.Neg().BigInt())
	liquidationPrice = liquidationPrice.Quo(sdk.NewDecFromBigInt(leveragedAmount.BigInt()))
	liquidationPrice = liquidationPrice.Add(openPrice)

	// get pool rates
	poolId, err := k.GetBestPool(ctx, req.Collateral.Denom, req.TradingAsset)
	if err != nil {
		return nil, err
	}
	borrowInterestRate := sdk.ZeroDec()
	fundingRate := sdk.ZeroDec()
	pool, found := k.GetPool(ctx, poolId)
	if found {
		borrowInterestRate = pool.BorrowInterestRate
		fundingRate = pool.FundingRate
	}

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
		EstimatedPnl:       estimatedPnLInt,
		AvailableLiquidity: availableLiquidity,
		WeightBalanceRatio: weightBonus,
		PriceImpact:        priceImpact,
		BorrowInterestRate: borrowInterestRate,
		FundingRate:        fundingRate,
	}, nil
}
