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

	return k.HandleOpenEstimation(ctx, req)
}

func (k Keeper) HandleOpenEstimation(ctx sdk.Context, req *types.QueryOpenEstimationRequest) (*types.QueryOpenEstimationResponse, error) {
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

	leveragedAmount := sdk.NewDecFromBigInt(collateralAmountInBaseCurrency.Amount.BigInt()).Mul(req.Leverage).TruncateInt()
	leveragedCoin := sdk.NewCoin(baseCurrency, leveragedAmount)

	_, _, custody, openPrice, swapFee, discount, availableLiquidity, slippage, weightBonus, priceImpact, err := k.amm.CalcSwapEstimationByDenom(ctx, leveragedCoin, baseCurrency, req.TradingAsset, baseCurrency, req.Discount, swapFee, decimals)
	if err != nil {
		return nil, err
	}

	if req.Position == types.Position_SHORT {
		custody = sdk.NewCoin(req.Collateral.Denom, leveragedAmount)
	}

	// invert openPrice if collateral is not in base currency
	openPrice = sdk.OneDec().Quo(openPrice)

	// check req.TakeProfitPrice not zero to prevent division by zero
	if req.TakeProfitPrice.IsZero() {
		req.TakeProfitPrice = sdk.MustNewDecFromStr(types.TakeProfitPriceDefault)
	}

	// calculate liabilities amount
	liabilitiesAmountDec := sdk.NewDecFromBigInt(collateralAmountInBaseCurrency.Amount.BigInt()).Mul(req.Leverage.Sub(sdk.OneDec()))

	// calculate estimated pnl
	var estimatedPnL sdk.Dec
	estimatedPnLDenom := baseCurrency

	// if position is short then:
	if req.Position == types.Position_SHORT {
		// estimated_pnl = custody_amount - liabilities_amount * take_profit_price - collateral_amount
		estimatedPnL = sdk.NewDecFromBigInt(custody.Amount.BigInt()).Sub(liabilitiesAmountDec.Mul(req.TakeProfitPrice)).Sub(sdk.NewDecFromBigInt(req.Collateral.Amount.BigInt()))
	} else {
		// if position is long then:
		// if collateral is not in base currency
		if req.Collateral.Denom != baseCurrency {
			// estimated_pnl = (custody_amount - collateral_amount) * take_profit_price - liabilities_amount
			estimatedPnL = sdk.NewDecFromBigInt(custody.Amount.BigInt()).Sub(sdk.NewDecFromBigInt(req.Collateral.Amount.BigInt())).Mul(req.TakeProfitPrice).Sub(liabilitiesAmountDec)
		} else {
			// estimated_pnl = custody_amount * take_profit_price - liabilities_amount - collateral_amount
			estimatedPnL = sdk.NewDecFromBigInt(custody.Amount.BigInt()).Mul(req.TakeProfitPrice).Sub(liabilitiesAmountDec).Sub(sdk.NewDecFromBigInt(req.Collateral.Amount.BigInt()))
		}
	}

	// calculate liquidation price
	// liquidation_price = open_price_value - collateral_amount / custody_amount
	liquidationPrice := openPrice.Sub(
		sdk.NewDecFromBigInt(collateralAmountInBaseCurrency.Amount.BigInt()).Quo(sdk.NewDecFromBigInt(custody.Amount.BigInt())),
	)

	positionSizeInTradingAsset := custody

	// if position is short then liquidation price is open price + collateral amount / (custody amount / open price)
	if req.Position == types.Position_SHORT {
		// for short position size = liabilities
		custodyAmountInTradingAssetDec := sdk.NewDecFromBigInt(custody.Amount.BigInt()).Quo(openPrice)
		liquidationPrice = openPrice.Add(
			sdk.NewDecFromBigInt(collateralAmountInBaseCurrency.Amount.BigInt()).Quo(custodyAmountInTradingAssetDec),
		)
		positionSizeInTradingAssetDec := sdk.NewDecFromBigInt((custody.Amount.Sub(req.Collateral.Amount)).BigInt()).Quo(openPrice)
		positionSizeInTradingAsset = sdk.NewCoin(req.TradingAsset, positionSizeInTradingAssetDec.TruncateInt())
	}

	// get pool rates
	poolId, err := k.GetBestPool(ctx, baseCurrency, req.TradingAsset)
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

	borrowFee := borrowInterestRate.MulInt(leveragedAmount.Sub(collateralAmountInBaseCurrency.Amount))
	fundingFee := fundingRate.MulInt(leveragedAmount)

	// | funding rate        |   direction     |  funding fee               |
	// ----------------------------------------------------------------------
	// | +ve                 |   LONG          |  amount to pay per block   |
	// | +ve                 |   SHORT         |  amount to earn per block  |
	// | -ve                 |   LONG          |  amount to earn per block  |
	// | -ve                 |   SHORT         |  amount to pay per block   |
	if req.Position == types.Position_LONG && fundingRate.IsPositive() {
		// long pays
		fundingFee = fundingFee.Abs()
	} else if req.Position == types.Position_SHORT && fundingRate.IsPositive() {
		// short earns
		fundingFee = fundingFee.Neg()
	} else if req.Position == types.Position_LONG && fundingRate.IsNegative() {
		// long earns
		fundingFee = fundingFee.Neg()
	} else if req.Position == types.Position_SHORT && fundingRate.IsNegative() {
		// short pays
		fundingFee = fundingFee.Abs()
	}

	ammPool, err := k.GetAmmPool(ctx, poolId, "")
	if err != nil {
		return nil, err
	}
	liabilitiesAsset := baseCurrency
	custodyAsset := req.TradingAsset
	mtp := types.NewMTP("", req.Collateral.Denom, req.TradingAsset, liabilitiesAsset, custodyAsset, req.Position, req.Leverage, req.TakeProfitPrice, poolId)
	interestAmount := k.GetBorrowInterest(ctx, mtp, ammPool)

	return &types.QueryOpenEstimationResponse{
		Position:           req.Position,
		Leverage:           req.Leverage,
		TradingAsset:       req.TradingAsset,
		Collateral:         req.Collateral,
		InterestAmount:     interestAmount,
		Custody:            custody,
		PositionSize:       positionSizeInTradingAsset,
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
		BorrowFee: sdk.Coin{
			Denom:  baseCurrency,
			Amount: borrowFee.TruncateInt(),
		},
		FundingFee: sdk.Coin{
			Denom:  baseCurrency,
			Amount: fundingFee.TruncateInt(),
		},
	}, nil
}
