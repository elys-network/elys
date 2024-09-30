package keeper

import (
	"context"
	sdkmath "cosmossdk.io/math"

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

	// get swap fee param
	swapFee := k.GetSwapFee(ctx)

	// retrieve base currency denom
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	switch req.Position {
	case types.Position_LONG:
		if err := types.CheckLongAssets(req.Collateral.Denom, req.TradingAsset, baseCurrency); err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		if err := types.CheckShortAssets(req.Collateral.Denom, req.TradingAsset, baseCurrency); err != nil {
			return nil, err
		}
	default:
		return nil, errorsmod.Wrap(types.ErrInvalidPosition, req.Position.String())
	}

	// retrieve denom in decimals
	entry, found = k.assetProfileKeeper.GetEntryByDenom(ctx, req.Collateral.Denom)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", req.Collateral.Denom)
	}
	decimals := entry.Decimals

	// if collateral amount is not in base currency then convert it
	collateralAmountInBaseCurrency := req.Collateral
	if req.Collateral.Denom != baseCurrency {
		var err error
		_, _, collateralAmountInBaseCurrency, _, _, _, _, _, _, _, err = k.amm.CalcSwapEstimationByDenom(ctx, req.Collateral, req.Collateral.Denom, baseCurrency, baseCurrency, req.Discount, swapFee, uint64(0))
		if err != nil {
			return nil, err
		}
	}

	leveragedAmount := sdkmath.LegacyNewDecFromBigInt(collateralAmountInBaseCurrency.Amount.BigInt()).Mul(req.Leverage).TruncateInt()
	leveragedCoin := sdk.NewCoin(baseCurrency, leveragedAmount)

	_, _, positionSize, openPrice, swapFee, discount, availableLiquidity, slippage, weightBonus, priceImpact, err := k.amm.CalcSwapEstimationByDenom(ctx, leveragedCoin, baseCurrency, req.TradingAsset, baseCurrency, req.Discount, swapFee, decimals)
	if err != nil {
		return nil, err
	}

	if req.Position == types.Position_SHORT {
		positionSize = sdk.NewCoin(req.Collateral.Denom, leveragedAmount)
	}

	// invert openPrice if collateral is not in base currency
	openPrice = sdkmath.LegacyOneDec().Quo(openPrice)

	// calculate min collateral
	minCollateral, err := k.CalcMinCollateral(ctx, req.Leverage, openPrice, decimals)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrCalcMinCollateral, "error calculating min collateral: %s", err.Error())
	}

	// check req.TakeProfitPrice not zero to prevent division by zero
	if req.TakeProfitPrice.IsZero() {
		return nil, errorsmod.Wrapf(types.ErrAmountTooLow, "take profit price is zero")
	}

	// calculate liabilities amount
	liabilitiesAmountDec := sdkmath.LegacyNewDecFromBigInt(collateralAmountInBaseCurrency.Amount.BigInt()).Mul(req.Leverage.Sub(sdkmath.LegacyOneDec()))

	// calculate estimated pnl
	// estimated_pnl = custody_amount - (liability_amount + collateral_amount) / take_profit_price
	estimatedPnL := sdkmath.LegacyNewDecFromBigInt(positionSize.Amount.BigInt())
	estimatedPnL = estimatedPnL.Sub(liabilitiesAmountDec.Add(sdkmath.LegacyNewDecFromBigInt(req.Collateral.Amount.BigInt())).Quo(req.TakeProfitPrice))
	estimatedPnLDenom := req.TradingAsset

	// if position is short then estimated pnl is custody_amount / open_price - (liability_amount + collateral_amount) / take_profit_price
	if req.Position == types.Position_SHORT {
		estimatedPnL = liabilitiesAmountDec.Add(sdkmath.LegacyNewDecFromBigInt(req.Collateral.Amount.BigInt())).Quo(req.TakeProfitPrice)
		estimatedPnL = estimatedPnL.Sub(sdkmath.LegacyNewDecFromBigInt(positionSize.Amount.BigInt()).Quo(openPrice))
		estimatedPnLDenom = baseCurrency
	}

	// calculate liquidation price
	// liquidation_price = open_price_value - collateral_amount / custody_amount
	liquidationPrice := openPrice.Sub(
		sdkmath.LegacyNewDecFromBigInt(collateralAmountInBaseCurrency.Amount.BigInt()).Quo(sdkmath.LegacyNewDecFromBigInt(positionSize.Amount.BigInt())),
	)

	// if position is short then liquidation price is open price + collateral amount / (custody amount / open price)
	if req.Position == types.Position_SHORT {
		positionSizeInTradingAsset := sdkmath.LegacyNewDecFromBigInt(positionSize.Amount.BigInt()).Quo(openPrice)
		liquidationPrice = openPrice.Add(
			sdkmath.LegacyNewDecFromBigInt(collateralAmountInBaseCurrency.Amount.BigInt()).Quo(positionSizeInTradingAsset),
		)
	}

	// get pool rates
	poolId, err := k.GetBestPool(ctx, baseCurrency, req.TradingAsset)
	if err != nil {
		return nil, err
	}
	borrowInterestRate := sdkmath.LegacyZeroDec()
	fundingRate := sdkmath.LegacyZeroDec()
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
		EstimatedPnl:       estimatedPnL.TruncateInt(),
		EstimatedPnlDenom:  estimatedPnLDenom,
		AvailableLiquidity: availableLiquidity,
		Slippage:           slippage,
		WeightBalanceRatio: weightBonus,
		PriceImpact:        priceImpact,
		BorrowInterestRate: borrowInterestRate,
		FundingRate:        fundingRate,
	}, nil
}
