package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CloseEstimation(goCtx context.Context, req *types.QueryCloseEstimationRequest) (res *types.QueryCloseEstimationResponse, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return k.HandleCloseEstimation(ctx, req)
}

func (k Keeper) HandleCloseEstimation(ctx sdk.Context, req *types.QueryCloseEstimationRequest) (res *types.QueryCloseEstimationResponse, err error) {
	if req.CloseAmount.IsNegative() {
		return nil, status.Error(codes.InvalidArgument, "invalid close_amount")
	}
	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	mtp, err := k.GetMTP(ctx, req.PoolId, address, req.PositionId)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	pool, found := k.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return &types.QueryCloseEstimationResponse{}, fmt.Errorf("perpetual pool %d not found", mtp.AmmPoolId)
	}

	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	err = k.UpdateMTPBorrowInterestUnpaidLiability(ctx, &mtp)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}
	_, _, _, err = k.UpdateFundingFee(ctx, &mtp, &pool)
	if err != nil {
		return nil, err
	}
	unpaidInterestLiability := mtp.BorrowInterestUnpaidLiability

	tradingAssetPrice, tradingAssetPriceDenomRatio, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	borrowInterestPaymentInCustody, err := mtp.GetBorrowInterestAmountAsCustodyAsset(tradingAssetPriceDenomRatio)
	if err != nil {
		return nil, err
	}
	mtp.Custody = mtp.Custody.Sub(borrowInterestPaymentInCustody)

	maxCloseAmount := mtp.Custody
	if mtp.Position == types.Position_SHORT {
		maxCloseAmount = mtp.Liabilities
	}

	closingRatio := sdkmath.LegacyOneDec()
	if req.CloseAmount.IsPositive() && req.CloseAmount.LT(maxCloseAmount) {
		closingRatio = req.CloseAmount.ToLegacyDec().Quo(maxCloseAmount.ToLegacyDec())
	}

	repayAmount, payingLiabilities, slippage, _, weightBreakingFee, _, _, _, executionPrice, err := k.CalcRepayAmount(ctx, &mtp, &ammPool, closingRatio)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	// need to make sure mtp.Custody has been used to unpaid liability
	returnAmount, err := k.CalcReturnAmount(mtp, repayAmount, closingRatio)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	mtp.Liabilities = mtp.Liabilities.Sub(payingLiabilities)
	mtp.Custody = mtp.Custody.ToLegacyDec().Mul(sdkmath.LegacyOneDec().Sub(closingRatio)).TruncateInt()
	mtp.Collateral = mtp.Collateral.ToLegacyDec().Mul(sdkmath.LegacyOneDec().Sub(closingRatio)).TruncateInt()

	liquidationPrice, err := k.GetLiquidationPrice(ctx, mtp)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	priceImpact := tradingAssetPrice.Sub(executionPrice).Quo(tradingAssetPrice)

	positionSize := mtp.Custody
	positionAsset := mtp.CustodyAsset
	if mtp.Position == types.Position_SHORT {
		positionSize = mtp.Liabilities
		positionAsset = mtp.LiabilitiesAsset
	}
	return &types.QueryCloseEstimationResponse{
		Position:                      mtp.Position,
		PositionSize:                  sdk.NewCoin(positionAsset, positionSize),
		Liabilities:                   sdk.NewCoin(mtp.LiabilitiesAsset, mtp.Liabilities),
		Custody:                       sdk.NewCoin(mtp.CustodyAsset, mtp.Custody),
		Collateral:                    sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral),
		PriceImpact:                   priceImpact,
		LiquidationPrice:              liquidationPrice,
		MaxCloseAmount:                maxCloseAmount,
		ClosingPrice:                  executionPrice,
		BorrowInterestUnpaidLiability: sdk.NewCoin(mtp.LiabilitiesAsset, unpaidInterestLiability),
		ReturningAmount:               sdk.NewCoin(mtp.CustodyAsset, returnAmount),
		PayingLiabilities:             sdk.NewCoin(mtp.LiabilitiesAsset, payingLiabilities),
		WeightBreakingFee:             weightBreakingFee.Dec(),
		Slippage:                      slippage.Dec(),
	}, nil
}
