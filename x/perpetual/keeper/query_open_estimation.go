package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OpenEstimation(goCtx context.Context, req *types.QueryOpenEstimationRequest) (*types.QueryOpenEstimationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return k.HandleOpenEstimation(ctx, req)
}

func (k Keeper) HandleOpenEstimation(ctx sdk.Context, req *types.QueryOpenEstimationRequest) (*types.QueryOpenEstimationResponse, error) {
	_, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, status.Error(codes.NotFound, "pool not found")
	}

	ammPool, err := k.GetAmmPool(ctx, req.PoolId)
	if err != nil {
		return nil, err
	}
	if req.Leverage.LTE(math.LegacyOneDec()) {
		return nil, status.Error(codes.InvalidArgument, "leverage must be greater than one")
	}

	snapshot := k.amm.GetPoolWithAccountedBalance(ctx, req.PoolId)
	tradingAssetLiquidity, err := snapshot.GetAmmPoolBalance(req.TradingAsset)
	if err != nil {
		return nil, err
	}
	availableLiquidity := sdk.NewCoin(req.TradingAsset, tradingAssetLiquidity)
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

	tradingAssetPrice, _, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, req.TradingAsset)
	if err != nil {
		return nil, err
	}
	useLimitPrice := !req.LimitPrice.IsNil() && !req.LimitPrice.IsZero()

	assetPriceAtOpen := tradingAssetPrice

	var limitPriceDenomRatio osmomath.BigDec
	if useLimitPrice {
		assetPriceAtOpen = req.LimitPrice
		limitPriceDenomRatio, err = k.ConvertPriceToAssetUsdcDenomRatio(ctx, req.TradingAsset, assetPriceAtOpen)
		if err != nil {
			return nil, err
		}
	}

	if req.TakeProfitPrice.IsPositive() {
		if req.Position == types.Position_LONG && req.TakeProfitPrice.LTE(assetPriceAtOpen) {
			return nil, status.Error(codes.InvalidArgument, "take profit price cannot be less than equal to trading price for long")
		}
		if req.Position == types.Position_SHORT && req.TakeProfitPrice.GTE(assetPriceAtOpen) {
			return nil, status.Error(codes.InvalidArgument, "take profit price cannot be greater than equal to trading price for short")
		}
	}

	// retrieve denom in decimals
	entry, found = k.assetProfileKeeper.GetEntryByDenom(ctx, req.Collateral.Denom)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", req.Collateral.Denom)
	}

	custodyAsset := req.TradingAsset
	liabilitiesAsset := baseCurrency
	if req.Position == types.Position_SHORT {
		liabilitiesAsset = req.TradingAsset
		custodyAsset = baseCurrency
	}
	mtp := types.NewMTP(ctx, "", req.Collateral.Denom, req.TradingAsset, liabilitiesAsset, custodyAsset, req.Position, req.TakeProfitPrice, req.PoolId)

	proxyLeverage := req.Leverage
	if req.Position == types.Position_SHORT {
		proxyLeverage = req.Leverage.Add(math.LegacyOneDec())
	}
	leveragedAmount := proxyLeverage.MulInt(req.Collateral.Amount).TruncateInt()
	// LONG: if collateral asset is trading asset then custodyAmount = leveragedAmount else if it collateral asset is usdc, we swap it to trading asset below
	// SHORT: collateralAsset is always usdc, and custody has to be in usdc, so custodyAmount = leveragedAmount
	custodyAmount := leveragedAmount
	slippage := osmomath.ZeroBigDec()
	mtp.Collateral = req.Collateral.Amount
	eta := proxyLeverage.Sub(math.LegacyOneDec())
	liabilities := eta.MulInt(req.Collateral.Amount).TruncateInt()
	weightBreakingFee := osmomath.ZeroBigDec()
	if req.Position == types.Position_LONG {
		//getting custody
		if mtp.CollateralAsset == baseCurrency {
			leveragedAmtTokenIn := sdk.NewCoin(mtp.CollateralAsset, leveragedAmount)
			custodyAmount, slippage, weightBreakingFee, err = k.EstimateSwapGivenIn(ctx, leveragedAmtTokenIn, mtp.CustodyAsset, ammPool, req.Address)
			if err != nil {
				return nil, err
			}
			if useLimitPrice {
				// leveragedAmount is collateral asset which is base currency, custodyAmount has to be in trading asset
				custodyAmount = osmomath.BigDecFromSDKInt(leveragedAmount).Quo(limitPriceDenomRatio).Dec().TruncateInt()
			}
		}

		//getting Liabilities
		if mtp.CollateralAsset != baseCurrency {
			amountIn := eta.MulInt(req.Collateral.Amount).TruncateInt()
			liabilities, slippage, weightBreakingFee, err = k.EstimateSwapGivenOut(ctx, sdk.NewCoin(req.Collateral.Denom, amountIn), baseCurrency, ammPool, req.Address)
			if err != nil {
				return nil, err
			}
			if useLimitPrice {
				// liabilities needs to be in base currency
				liabilities = osmomath.BigDecFromSDKInt(amountIn).Mul(limitPriceDenomRatio).Dec().TruncateInt()
			}
		}

	}
	//getting Liabilities
	if req.Position == types.Position_SHORT {
		// Collateral will be in base currency
		amountOut := eta.MulInt(req.Collateral.Amount).TruncateInt()
		tokenOut := sdk.NewCoin(baseCurrency, amountOut)
		liabilities, slippage, weightBreakingFee, err = k.EstimateSwapGivenOut(ctx, tokenOut, mtp.LiabilitiesAsset, ammPool, mtp.Address)
		if err != nil {
			return nil, err
		}
		if useLimitPrice {
			// liabilities needs to be in trading asset
			liabilities = osmomath.BigDecFromSDKInt(amountOut).Quo(limitPriceDenomRatio).Dec().TruncateInt()
		}
	}
	mtp.Liabilities = liabilities
	mtp.Custody = custodyAmount

	mtp.TakeProfitCustody, err = k.CalcMTPTakeProfitCustody(ctx, *mtp)
	if err != nil {
		return nil, err
	}
	mtp.TakeProfitLiabilities, err = k.CalcMTPTakeProfitLiability(ctx, *mtp)
	if err != nil {
		return nil, err
	}
	mtp.TakeProfitPrice = req.TakeProfitPrice
	err = k.GetAndSetOpenPrice(ctx, mtp)
	if err != nil {
		return nil, err
	}
	executionPrice := mtp.OpenPrice

	err = mtp.UpdateMTPTakeProfitBorrowFactor()
	if err != nil {
		return nil, err
	}
	hourlyInterestRate := math.LegacyZeroDec()
	blocksPerYear := math.LegacyNewDec(int64(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear))
	blocksPerSecond := blocksPerYear.QuoInt64(86400 * 365)                                      // in seconds
	startBlock := ctx.BlockHeight() - math.LegacyNewDec(3600).Mul(blocksPerSecond).RoundInt64() // block height 1 hour ago
	if startBlock > 0 {
		hourlyInterestRate = k.GetBorrowInterestRate(ctx, uint64(startBlock), uint64(ctx.BlockTime().Unix()-3600), req.PoolId, mtp.TakeProfitBorrowFactor)
	}

	liquidationPrice, err := k.GetLiquidationPrice(ctx, *mtp)
	if err != nil {
		return nil, err
	}
	priceImpact := tradingAssetPrice.Sub(executionPrice).Quo(tradingAssetPrice)

	estimatedPnLAmount := math.ZeroInt()
	if req.TakeProfitPrice.IsPositive() {
		estimatedPnLAmount, err = k.GetEstimatedPnL(ctx, *mtp, baseCurrency, true)
		if err != nil {
			return nil, err
		}
	}

	borrowInterestRate := k.GetBorrowInterestRate(ctx, mtp.LastInterestCalcBlock, mtp.LastInterestCalcTime, req.PoolId, mtp.TakeProfitBorrowFactor)

	longRate, shortRate := k.GetFundingRate(ctx, uint64(ctx.BlockHeight()), uint64(ctx.BlockTime().Unix()), req.PoolId)
	fundingRate := longRate
	if req.Position == types.Position_SHORT {
		fundingRate = shortRate
	}

	positionSize := mtp.Custody
	positionAsset := mtp.CustodyAsset
	if req.Position == types.Position_SHORT {
		positionSize = mtp.Liabilities
		positionAsset = mtp.LiabilitiesAsset
	}

	effectiveLeverage, err := k.GetEffectiveLeverage(ctx, *mtp)
	if err != nil {
		return nil, err
	}

	return &types.QueryOpenEstimationResponse{
		Position:           req.Position,
		EffectiveLeverage:  effectiveLeverage,
		TradingAsset:       req.TradingAsset,
		Collateral:         req.Collateral,
		HourlyInterestRate: hourlyInterestRate,
		PositionSize:       sdk.NewCoin(positionAsset, positionSize),
		OpenPrice:          mtp.OpenPrice,
		TakeProfitPrice:    req.TakeProfitPrice,
		LiquidationPrice:   liquidationPrice,
		EstimatedPnl:       sdk.Coin{Denom: baseCurrency, Amount: estimatedPnLAmount},
		AvailableLiquidity: availableLiquidity,
		Slippage:           slippage.Dec(),
		PriceImpact:        priceImpact,
		BorrowInterestRate: borrowInterestRate,
		FundingRate:        fundingRate,
		Custody:            sdk.NewCoin(mtp.CustodyAsset, mtp.Custody),
		Liabilities:        sdk.NewCoin(mtp.LiabilitiesAsset, mtp.Liabilities),
		LimitPrice:         req.LimitPrice,
		WeightBreakingFee:  weightBreakingFee.Dec(),
	}, nil
}
