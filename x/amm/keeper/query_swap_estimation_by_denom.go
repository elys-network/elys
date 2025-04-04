package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
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

	inRoute, outRoute, amount, spotPrice, swapFee, discount, availableLiquidity, slippage, weightBonus, priceImpact, err := k.CalcSwapEstimationByDenom(ctx, req.Amount, req.DenomIn, req.DenomOut, baseCurrency, req.Address, osmomath.ZeroBigDec(), entry.Decimals)
	if err != nil {
		return nil, err
	}
	recoveryReward := sdkmath.ZeroInt()
	if weightBonus.IsPositive() {
		recoveryReward = osmomath.BigDecFromSDKInt(amount.Amount).Mul(weightBonus).Dec().TruncateInt()
	}
	// Add weight balance amount here, not added in execution as out amount will be changed and that will impact the transfers
	amount.Amount = amount.Amount.Add(recoveryReward)

	return &types.QuerySwapEstimationByDenomResponse{
		InRoute:            inRoute,
		OutRoute:           outRoute,
		Amount:             amount,
		SpotPrice:          spotPrice.Dec(),
		SwapFee:            swapFee.Dec(),
		Discount:           discount.Dec(),
		AvailableLiquidity: availableLiquidity,
		Slippage:           slippage.Dec(),
		WeightBalanceRatio: weightBonus.Dec(),
		PriceImpact:        priceImpact.Dec(),
		// sdk.NewCoin() will panic in case of negative weightBonus
		WeightBalanceRewardAmount: sdk.Coin{Denom: amount.Denom, Amount: recoveryReward},
	}, nil
}
