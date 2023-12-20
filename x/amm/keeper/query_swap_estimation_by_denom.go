package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SwapEstimationByDenom(goCtx context.Context, req *types.QuerySwapEstimationByDenomRequest) (*types.QuerySwapEstimationByDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	_ = baseCurrency

	inRoute, outRoute, amount, spotPrice, swapFee, discount, availableLiquidity, weightBonus, priceImpact, err := k.CalcSwapEstimationByDenom(ctx, req.Amount, req.DenomIn, req.DenomOut, baseCurrency, req.Discount, sdk.ZeroDec())
	if err != nil {
		return nil, err
	}

	return &types.QuerySwapEstimationByDenomResponse{
		InRoute:            inRoute,
		OutRoute:           outRoute,
		Amount:             amount,
		SpotPrice:          spotPrice,
		SwapFee:            swapFee,
		Discount:           discount,
		AvailableLiquidity: availableLiquidity,
		WeightBalanceRatio: weightBonus,
		PriceImpact:        priceImpact,
	}, nil
}
