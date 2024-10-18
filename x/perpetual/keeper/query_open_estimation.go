package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
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

	tradingAssetLiquidity, err := ammPool.GetAmmPoolBalance(req.TradingAsset)
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

	tradingAssetPrice, err := k.GetAssetPriceByDenom(ctx, req.TradingAsset)
	if err != nil {
		return nil, err
	}
	if req.Position == types.Position_LONG && req.TakeProfitPrice.LTE(tradingAssetPrice) {
		return nil, status.Error(codes.InvalidArgument, "take profit price cannot be less than equal to trading price for long")
	}
	if req.Position == types.Position_SHORT && req.TakeProfitPrice.GTE(tradingAssetPrice) {
		return nil, status.Error(codes.InvalidArgument, "take profit price cannot be greater than equal to trading price for short")
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
	mtp := types.NewMTP("", req.Collateral.Denom, req.TradingAsset, liabilitiesAsset, custodyAsset, req.Position, req.TakeProfitPrice, req.PoolId)

	blocksPerYear := sdk.NewDec(k.parameterKeeper.GetParams(ctx).TotalBlocksPerYear)
	avgBlockTime := blocksPerYear.Quo(math.LegacyNewDec(86400 * 365)).TruncateInt().Uint64()
	mtp.LastInterestCalcBlock = uint64(ctx.BlockHeight()) - 1
	mtp.LastInterestCalcTime = uint64(ctx.BlockTime().Unix()) - avgBlockTime

	leveragedAmount := req.Collateral.Amount.ToLegacyDec().Mul(req.Leverage).TruncateInt()
	// LONG: if collateral asset is trading asset then custodyAmount = leveragedAmount else if it collateral asset is usdc, we swap it to trading asset below
	// SHORT: collateralAsset is always usdc, and custody has to be in usdc, so custodyAmount = leveragedAmount
	custodyAmount := leveragedAmount
	slippage := math.LegacyZeroDec()
	mtp.Collateral = req.Collateral.Amount
	executionPrice := math.LegacyZeroDec()
	eta := req.Leverage.Sub(math.LegacyOneDec())
	if req.Position == types.Position_LONG {
		//getting custody
		if mtp.CollateralAsset == baseCurrency {
			leveragedAmtTokenIn := sdk.NewCoin(mtp.CollateralAsset, leveragedAmount)
			custodyAmount, slippage, err = k.EstimateSwap(ctx, leveragedAmtTokenIn, mtp.CustodyAsset, ammPool)
			if err != nil {
				return nil, err
			}

			mtp.Liabilities = req.Collateral.Amount.ToLegacyDec().Mul(eta).TruncateInt()
			// executionPrice = leveragedAmount / Custody
			executionPrice = leveragedAmount.ToLegacyDec().Quo(custodyAmount.ToLegacyDec())
		}
		mtp.Custody = custodyAmount

		//getting Liabilities
		if mtp.CollateralAsset != baseCurrency {
			amountIn := req.Collateral.Amount.ToLegacyDec().Mul(eta).TruncateInt()
			mtp.Liabilities, slippage, err = k.EstimateSwapGivenOut(ctx, sdk.NewCoin(req.Collateral.Denom, amountIn), baseCurrency, ammPool)
			if err != nil {
				return nil, err
			}

			// executionPrice = (Liabilities in base currency) / Custody
			executionPrice = mtp.Liabilities.ToLegacyDec().Quo(amountIn.ToLegacyDec())
		}

	}
	//getting Liabilities
	if req.Position == types.Position_SHORT {
		mtp.Custody = custodyAmount
		// Collateral will be in base currency
		custodyTokenIn := sdk.NewCoin(baseCurrency, mtp.Custody)
		mtp.Liabilities, slippage, err = k.EstimateSwap(ctx, custodyTokenIn, mtp.LiabilitiesAsset, ammPool)
		if err != nil {
			return nil, err
		}

		executionPrice = mtp.Custody.ToLegacyDec().Quo(mtp.Liabilities.ToLegacyDec())
	}
	mtp.TakeProfitCustody = types.CalcMTPTakeProfitCustody(*mtp)
	mtp.TakeProfitLiabilities, err = k.CalcMTPTakeProfitLiability(ctx, mtp, baseCurrency)
	mtp.TakeProfitPrice = req.TakeProfitPrice
	mtp.OpenPrice, err = k.CalOpenPrice(ctx, mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	k.UpdateMTPBorrowInterestUnpaidLiability(ctx, mtp)

	liquidationPrice := k.GetLiquidationPrice(ctx, *mtp)
	priceImpact := tradingAssetPrice.Sub(executionPrice).Quo(tradingAssetPrice)

	estimatedPnLAmount, err := k.GetEstimatedPnL(ctx, *mtp, baseCurrency, true)
	if err != nil {
		return nil, err
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

	return &types.QueryOpenEstimationResponse{
		Position:           req.Position,
		Leverage:           req.Leverage,
		TradingAsset:       req.TradingAsset,
		Collateral:         req.Collateral,
		InterestAmount:     sdk.NewCoin(mtp.LiabilitiesAsset, mtp.BorrowInterestUnpaidLiability),
		PositionSize:       sdk.NewCoin(positionAsset, positionSize),
		OpenPrice:          mtp.OpenPrice,
		TakeProfitPrice:    req.TakeProfitPrice,
		LiquidationPrice:   liquidationPrice,
		EstimatedPnl:       sdk.Coin{mtp.CustodyAsset, estimatedPnLAmount},
		AvailableLiquidity: availableLiquidity,
		Slippage:           slippage,
		PriceImpact:        priceImpact,
		BorrowInterestRate: borrowInterestRate,
		FundingRate:        fundingRate,
	}, nil
}
