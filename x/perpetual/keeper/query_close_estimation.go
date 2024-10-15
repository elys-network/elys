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

func (k Keeper) CloseEstimation(goCtx context.Context, req *types.QueryCloseEstimationRequest) (res *types.QueryCloseEstimationResponse, err error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return k.HandleCloseEstimation(ctx, req)
}

func (k Keeper) HandleCloseEstimation(ctx sdk.Context, req *types.QueryCloseEstimationRequest) (res *types.QueryCloseEstimationResponse, err error) {
	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}
	mtp, err := k.CloseEstimationChecker.GetMTP(ctx, address, req.PositionId)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	// Retrieve Pool
	pool, found := k.CloseEstimationChecker.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return &types.QueryCloseEstimationResponse{}, errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
	}

	// Retrieve AmmPool
	ammPool, err := k.CloseEstimationChecker.GetAmmPool(ctx, mtp.AmmPoolId, mtp.CustodyAsset)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	// get base currency entry
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return &types.QueryCloseEstimationResponse{}, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// init repay amount
	var repayAmount math.Int

	// if position is long, repay in collateral asset
	if mtp.Position == types.Position_LONG {
		custodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, mtp.Custody)
		repayAmount, err = k.CloseEstimationChecker.EstimateSwap(ctx, custodyAmtTokenIn, mtp.CollateralAsset, ammPool)
		if err != nil {
			return &types.QueryCloseEstimationResponse{}, err
		}
	} else if mtp.Position == types.Position_SHORT {
		// if position is short, repay in trading asset
		custodyAmtTokenIn := sdk.NewCoin(mtp.CustodyAsset, mtp.Custody)
		repayAmount, err = k.CloseEstimationChecker.EstimateSwap(ctx, custodyAmtTokenIn, mtp.TradingAsset, ammPool)
		if err != nil {
			return &types.QueryCloseEstimationResponse{}, err
		}
	} else {
		return &types.QueryCloseEstimationResponse{}, types.ErrInvalidPosition
	}

	returnAmount, err := k.CalcReturnAmount(ctx, mtp, pool, ammPool, repayAmount, mtp.Custody, baseCurrency)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	// get swap fee param
	swapFee := k.GetSwapFee(ctx)

	// if collateral amount is not in base currency then convert it
	collateralAmountInBaseCurrency := mtp.Collateral
	if mtp.CollateralAsset != baseCurrency {
		var err error
		collateralAmountInBaseCurrency, err = k.CloseEstimationChecker.EstimateSwapGivenOut(ctx, sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral), baseCurrency, ammPool)
		if err != nil {
			return nil, err
		}
	}

	// calculate liquidation price
	// liquidation_price = open_price_value - collateral_amount / custody_amount
	liquidationPrice := mtp.OpenPrice.Sub(
		sdk.NewDecFromBigInt(collateralAmountInBaseCurrency.BigInt()).Quo(sdk.NewDecFromBigInt(mtp.Custody.BigInt())),
	)

	positionSizeInTradingAsset := mtp.Custody
	if mtp.Position == types.Position_SHORT {
		positionSizeInTradingAsset = mtp.Liabilities
	}

	return &types.QueryCloseEstimationResponse{
		Position:     mtp.Position,
		PositionSize: sdk.NewCoin(mtp.TradingAsset, positionSizeInTradingAsset),
		Custody:      sdk.NewCoin(mtp.CustodyAsset, mtp.Custody),
		Liabilities:  sdk.NewCoin(mtp.LiabilitiesAsset, mtp.Liabilities),
		// TODO: price impact calculation
		PriceImpact:      sdk.ZeroDec(),
		SwapFee:          swapFee,
		ReturnAmount:     sdk.NewCoin(mtp.CollateralAsset, returnAmount),
		LiquidationPrice: liquidationPrice,
	}, nil
}
