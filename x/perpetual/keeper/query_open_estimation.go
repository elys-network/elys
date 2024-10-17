package keeper

import (
	"context"
	"cosmossdk.io/math"

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

	params := k.GetParams(ctx)

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
		}
		mtp.Custody = custodyAmount

		//getting Liabilities
		if mtp.CollateralAsset != baseCurrency {
			amountIn := req.Collateral.Amount.ToLegacyDec().Mul(eta).TruncateInt()
			mtp.Liabilities, slippage, err = k.EstimateSwapGivenOut(ctx, sdk.NewCoin(req.Collateral.Denom, amountIn), baseCurrency, ammPool)
			if err != nil {
				return nil, err
			}

		} else {
			mtp.Liabilities = req.Collateral.Amount.ToLegacyDec().Mul(eta).TruncateInt()
		}
		// executionPrice = Liabilities / Custody
		executionPrice = mtp.Liabilities.ToLegacyDec().Quo(mtp.Custody.ToLegacyDec())

	}
	//getting Liabilities
	if req.Position == types.Position_SHORT {
		// Collateral will be in base currency
		liabilitiesInCollateralTokenIn := sdk.NewCoin(baseCurrency, req.Collateral.Amount.ToLegacyDec().Mul(eta).TruncateInt())
		mtp.Liabilities, slippage, err = k.EstimateSwap(ctx, liabilitiesInCollateralTokenIn, mtp.LiabilitiesAsset, ammPool)
		// executionPrice = liabilitiesInCollateral / Liabilities
		executionPrice = liabilitiesInCollateralTokenIn.Amount.ToLegacyDec().Quo(mtp.Liabilities.ToLegacyDec())
		mtp.Custody = custodyAmount
	}
	mtp.TakeProfitCustody = types.CalcMTPTakeProfitCustody(mtp)
	mtp.TakeProfitLiabilities, err = k.CalcMTPTakeProfitLiability(ctx, mtp, baseCurrency)
	mtp.TakeProfitPrice = req.TakeProfitPrice
	mtp.OpenPrice, err = k.CalOpenPrice(ctx, mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	k.UpdateMTPBorrowInterestUnpaidLiability(ctx, mtp)
	liabilityInterestTokenOut := sdk.NewCoin(mtp.LiabilitiesAsset, mtp.BorrowInterestUnpaidLiability)
	borrowInterestPaymentInCustody, _, err := k.EstimateSwapGivenOut(ctx, liabilityInterestTokenOut, mtp.CustodyAsset, ammPool)
	if err != nil {
		return nil, err
	}
	custodyAmountAfterInterest := mtp.Custody.Sub(borrowInterestPaymentInCustody)

	liquidationPrice := math.LegacyZeroDec()
	oraclePrice := math.LegacyZeroDec()
	// calculate liquidation price
	if mtp.Position == types.Position_LONG {
		// liquidation_price = (safety_factor * liabilities) / custody
		liquidationPrice = mtp.Liabilities.ToLegacyDec().Quo(params.SafetyFactor.Mul(mtp.Custody.ToLegacyDec()))
		oracleTokenPrice, found := k.oracleKeeper.GetAssetPrice(ctx, mtp.CustodyAsset)
		if found {
			oraclePrice = oracleTokenPrice.Price
		}
	}
	if mtp.Position == types.Position_SHORT {
		// liquidation_price =  Custody / (Liabilities * safety_factor)
		liquidationPrice = mtp.Custody.ToLegacyDec().Quo(mtp.Liabilities.ToLegacyDec().Mul(params.SafetyFactor))
		oracleTokenPrice, found := k.oracleKeeper.GetAssetPrice(ctx, mtp.LiabilitiesAsset)
		if found {
			oraclePrice = oracleTokenPrice.Price
		}
	}

	priceImpact := math.LegacyZeroDec()
	if !oraclePrice.IsZero() {
		priceImpact = oraclePrice.Sub(executionPrice).Quo(oraclePrice)
	}

	estimatedPnLAmount := math.ZeroInt()
	// if position is short then:
	if req.Position == types.Position_SHORT {
		// estimated_pnl = custodyAmountAfterInterest - liabilities_amount * take_profit_price - collateral_amount
		estimatedPnLAmount = custodyAmountAfterInterest.Sub(mtp.Collateral).Mul(mtp.Liabilities.ToLegacyDec().Mul(mtp.TakeProfitPrice).TruncateInt())
	} else {
		// if position is long then:
		// if collateral is not in base currency
		if types.IsTakeProfitPriceInifite(mtp) || mtp.TakeProfitPrice.IsZero() {
			estimatedPnLAmount = math.ZeroInt()
		} else {
			if req.Collateral.Denom != baseCurrency {
				// estimated_pnl = (custodyAmountAfterInterest - collateral_amount) * take_profit_price - liabilities_amount
				estimatedPnLAmount = (custodyAmountAfterInterest.Sub(mtp.Collateral)).ToLegacyDec().Mul(mtp.TakeProfitPrice).TruncateInt().Sub(mtp.Liabilities)
			} else {
				// estimated_pnl = custodyAmountAfterInterest * take_profit_price - liabilities_amount - collateral_amount
				estimatedPnLAmount = custodyAmountAfterInterest.ToLegacyDec().Mul(mtp.TakeProfitPrice).TruncateInt().Sub(mtp.Liabilities).Sub(mtp.Collateral)
			}
		}
	}

	// TODO, borrowInterestRate and fundingRate both are summed up values, discuss this for correct values
	borrowInterestRate := k.GetBorrowInterestRate(ctx, uint64(ctx.BlockHeight()), uint64(ctx.BlockTime().Unix()), req.PoolId, mtp.TakeProfitBorrowFactor)
	longRate, shortRate := k.GetFundingRate(ctx, uint64(ctx.BlockHeight()), uint64(ctx.BlockTime().Unix()), req.PoolId)
	borrowFee := borrowInterestRate.MulInt(mtp.Liabilities)
	fundingRate := longRate
	if req.Position == types.Position_SHORT {
		fundingRate = shortRate
	}
	fundingFee := fundingRate.MulInt(mtp.BorrowInterestUnpaidLiability)

	return &types.QueryOpenEstimationResponse{
		Position:           req.Position,
		Leverage:           req.Leverage,
		TradingAsset:       req.TradingAsset,
		Collateral:         req.Collateral,
		InterestAmount:     sdk.NewCoin(mtp.LiabilitiesAsset, mtp.BorrowInterestUnpaidLiability),
		PositionSize:       sdk.NewCoin(mtp.CustodyAsset, mtp.Custody),
		OpenPrice:          mtp.OpenPrice,
		TakeProfitPrice:    req.TakeProfitPrice,
		LiquidationPrice:   liquidationPrice,
		EstimatedPnl:       sdk.Coin{mtp.CustodyAsset, estimatedPnLAmount},
		AvailableLiquidity: availableLiquidity,
		Slippage:           slippage,
		PriceImpact:        priceImpact,
		BorrowInterestRate: borrowInterestRate,
		FundingRate:        fundingRate,
		BorrowFee:          sdk.NewCoin(mtp.LiabilitiesAsset, borrowFee.TruncateInt()),
		FundingFee:         sdk.NewCoin(mtp.LiabilitiesAsset, fundingFee.TruncateInt()),
	}, nil
}
