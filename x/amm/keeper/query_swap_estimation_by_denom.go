package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

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

	// retrieve base currency denom
	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	// retrieve denom in decimals
	entry, found := k.assetProfileKeeper.GetEntryByDenom(ctx, req.DenomIn)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", req.DenomIn)
	}

	inRoute, outRoute, amount, spotPrice, swapFee, discount, availableLiquidity, slippage, weightBonus, priceImpact, err := k.CalcSwapEstimationByDenom(ctx, req.Amount, req.DenomIn, req.DenomOut, baseCurrency, req.Address, sdkmath.LegacyZeroDec(), entry.Decimals)
	if err != nil {
		return nil, err
	}
	// Add weight balance amount here, not added in execution as out amount will be changed and that will impact the transfers
	amount.Amount = amount.Amount.Add(((amount.Amount.ToLegacyDec()).Mul(weightBonus)).TruncateInt())

	return &types.QuerySwapEstimationByDenomResponse{
		InRoute:                   inRoute,
		OutRoute:                  outRoute,
		Amount:                    amount,
		SpotPrice:                 spotPrice,
		SwapFee:                   swapFee,
		Discount:                  discount,
		AvailableLiquidity:        availableLiquidity,
		Slippage:                  slippage,
		WeightBalanceRatio:        weightBonus,
		PriceImpact:               priceImpact,
		WeightBalanceRewardAmount: sdk.NewCoin(amount.Denom, ((amount.Amount.ToLegacyDec()).Mul(weightBonus)).TruncateInt()),
	}, nil
}
