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

	_ = pool

	// Retrieve AmmPool
	ammPool, err := k.CloseEstimationChecker.GetAmmPool(ctx, mtp.AmmPoolId, mtp.CustodyAsset)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	_ = ammPool

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

	return &types.QueryCloseEstimationResponse{
		Position:     mtp.Position,
		PositionSize: sdk.NewCoin(mtp.CustodyAsset, mtp.Custody),
		Liabilities:  sdk.NewCoin(mtp.LiabilitiesAsset, mtp.Liabilities),
		// TODO: price impact calculation
		PriceImpact:  sdk.ZeroDec(),
		SwapFee:      swapFee,
		ReturnAmount: sdk.NewCoin(mtp.CollateralAsset, returnAmount),
	}, nil
}
