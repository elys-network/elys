package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
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

	mtp, err := k.GetMTP(ctx, address, req.PositionId)
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

	k.UpdateMTPBorrowInterestUnpaidLiability(ctx, &mtp)
	err = k.UpdateFundingFee(ctx, &mtp, &pool, ammPool)
	if err != nil {
		return nil, err
	}
	unpaidInterestLiability := mtp.BorrowInterestUnpaidLiability
	borrowInterestPaymentTokenIn := sdk.NewCoin(mtp.LiabilitiesAsset, mtp.BorrowInterestUnpaidLiability)
	borrowInterestPaymentInCustody, _, err := k.EstimateSwapGivenOut(ctx, borrowInterestPaymentTokenIn, mtp.CustodyAsset, ammPool)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}
	maxCloseAmount := mtp.Custody.Sub(borrowInterestPaymentInCustody)
	if mtp.Position == types.Position_SHORT {
		maxCloseAmount = mtp.Liabilities
	}

	closingRatio := math.LegacyOneDec()
	if req.CloseAmount.IsPositive() && req.CloseAmount.LT(maxCloseAmount) {
		closingRatio = req.CloseAmount.ToLegacyDec().Quo(maxCloseAmount.ToLegacyDec())
	}

	repayAmount, payingLiabilities, err := k.CalcRepayAmount(ctx, &mtp, &ammPool, closingRatio)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	// need to make sure mtp.Custody has been used to unpaid liability
	returnAmount, err := k.CalcReturnAmount(mtp, repayAmount, sdk.OneDec())
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	liquidationPrice := k.GetLiquidationPrice(ctx, mtp)
	executionPrice := math.LegacyZeroDec()
	// calculate liquidation price
	if mtp.Position == types.Position_LONG {
		// executionPrice = payingLiabilities / repayAmount
		executionPrice = payingLiabilities.ToLegacyDec().Quo(repayAmount.ToLegacyDec())
	}
	if mtp.Position == types.Position_SHORT {
		// executionPrice = repayAmount / payingLiabilities
		executionPrice = repayAmount.ToLegacyDec().Quo(payingLiabilities.ToLegacyDec())
	}

	tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
	if err != nil {
		return nil, err
	}
	priceImpact := tradingAssetPrice.Sub(executionPrice).Quo(tradingAssetPrice)

	returnAmountAtClosingPrice := math.ZeroInt()
	if req.ClosingPrice.IsPositive() {
		if mtp.Position == types.Position_LONG {
			borrowInterestPaymentInCustodyAtClosingPrice := unpaidInterestLiability.ToLegacyDec().Quo(req.ClosingPrice).TruncateInt()
			custodyAfterInterests := mtp.Custody.Sub(borrowInterestPaymentInCustodyAtClosingPrice)
			repayAmountAtClosingPrice := payingLiabilities.ToLegacyDec().Quo(req.ClosingPrice).TruncateInt()
			custodyAfterRepayAtClosingPrice := custodyAfterInterests.Sub(repayAmountAtClosingPrice)
			if custodyAfterRepayAtClosingPrice.IsPositive() {
				returnAmountAtClosingPrice = custodyAfterRepayAtClosingPrice
			}
		}
		if mtp.Position == types.Position_SHORT {
			// For short, liability is in trading asset, eg ATOM, custody is base currency eg USDC
			borrowInterestPaymentInCustodyAtClosingPrice := unpaidInterestLiability.ToLegacyDec().Mul(req.ClosingPrice).TruncateInt()
			custodyAfterInterests := mtp.Custody.Sub(borrowInterestPaymentInCustodyAtClosingPrice)
			repayAmountAtClosingPrice := payingLiabilities.ToLegacyDec().Mul(req.ClosingPrice).TruncateInt()
			custodyAfterRepayAtClosingPrice := custodyAfterInterests.Sub(repayAmountAtClosingPrice)
			if custodyAfterRepayAtClosingPrice.IsPositive() {
				returnAmountAtClosingPrice = custodyAfterRepayAtClosingPrice
			}
		}
	}

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
		PriceImpact:                   priceImpact,
		LiquidationPrice:              liquidationPrice,
		MaxCloseAmount:                maxCloseAmount,
		BorrowInterestUnpaidLiability: sdk.NewCoin(mtp.LiabilitiesAsset, unpaidInterestLiability),
		ReturningAmount:               sdk.NewCoin(mtp.CustodyAsset, returnAmount),
		PayingLiabilities:             sdk.NewCoin(mtp.LiabilitiesAsset, payingLiabilities),
		ReturnAmountAtClosingPrice:    sdk.NewCoin(mtp.CustodyAsset, returnAmountAtClosingPrice),
	}, nil
}
