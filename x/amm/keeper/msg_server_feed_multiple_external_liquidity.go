package keeper

import (
	"context"
	sdkmath "cosmossdk.io/math"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func AssetsValue(ctx sdk.Context, oracleKeeper types.OracleKeeper, amountDepthInfo []types.AssetAmountDepth) (sdkmath.LegacyDec, sdkmath.LegacyDec, error) {
	totalValue := sdkmath.LegacyZeroDec()
	totalDepth := sdkmath.LegacyZeroDec()
	if len(amountDepthInfo) == 0 {
		return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), nil
	}
	for _, asset := range amountDepthInfo {
		price, found := oracleKeeper.GetAssetPrice(ctx, asset.Asset)
		if !found {
			return sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), fmt.Errorf("asset price not set: %s", asset.Asset)
		} else {
			v := price.Price.Mul(asset.Amount)
			totalValue = totalValue.Add(v)
		}
		totalDepth = totalDepth.Add(asset.Depth)
	}
	avgDepth := totalDepth.Quo(sdkmath.LegacyNewDec(int64(len(amountDepthInfo))))
	return totalValue, avgDepth, nil
}

func LiquidityRatioFromPriceDepth(depth sdkmath.LegacyDec) sdkmath.LegacyDec {
	if depth == sdkmath.LegacyOneDec() {
		return sdkmath.LegacyOneDec()
	}
	sqrt, err := sdkmath.LegacyOneDec().Sub(depth).ApproxSqrt()
	if err != nil {
		panic(err)
	}
	return sdkmath.LegacyOneDec().Sub(sqrt)
}

func (k msgServer) FeedMultipleExternalLiquidity(goCtx context.Context, msg *types.MsgFeedMultipleExternalLiquidity) (*types.MsgFeedMultipleExternalLiquidityResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	feeder, found := k.oracleKeeper.GetPriceFeeder(ctx, sdk.MustAccAddressFromBech32(msg.Sender))
	if !found {
		return nil, oracletypes.ErrNotAPriceFeeder
	}

	if !feeder.IsActive {
		return nil, oracletypes.ErrPriceFeederNotActive
	}

	for _, el := range msg.Liquidity {
		pool, found := k.GetPool(ctx, el.PoolId)
		if !found {
			return nil, types.ErrInvalidPoolId
		}

		tvl, err := pool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
		if err != nil {
			return nil, err
		}

		elValue, elDepth, err := AssetsValue(ctx, k.oracleKeeper, el.AmountDepthInfo)
		if err != nil {
			return nil, err
		}

		// Ensure tvl is not zero to avoid division by zero
		if tvl.IsZero() {
			return nil, types.ErrAmountTooLow
		}

		elRatio := elValue.Quo(tvl)

		// calculate liquidity ratio
		liquidityRatio := LiquidityRatioFromPriceDepth(elDepth)

		// Ensure tvl is not zero to avoid division by zero
		if liquidityRatio.IsZero() {
			return nil, types.ErrAmountTooLow
		}

		elRatio = elRatio.Quo(liquidityRatio)
		if elRatio.LT(sdkmath.LegacyOneDec()) {
			elRatio = sdkmath.LegacyOneDec()
		}

		pool.PoolParams.ExternalLiquidityRatio = elRatio
		k.SetPool(ctx, pool)
	}

	return &types.MsgFeedMultipleExternalLiquidityResponse{}, nil
}
