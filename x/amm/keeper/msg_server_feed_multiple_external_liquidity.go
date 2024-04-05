package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func AssetsValue(ctx sdk.Context, oracleKeeper types.OracleKeeper, amountDepthInfo []types.AssetAmountDepth) (sdk.Dec, sdk.Dec, error) {
	totalValue := sdk.ZeroDec()
	totalDepth := sdk.ZeroDec()
	if len(amountDepthInfo) == 0 {
		return sdk.ZeroDec(), sdk.ZeroDec(), nil
	}
	for _, asset := range amountDepthInfo {
		price, found := oracleKeeper.GetAssetPrice(ctx, asset.Asset)
		if !found {
			return sdk.ZeroDec(), sdk.ZeroDec(), fmt.Errorf("asset price not set: %s", asset.Asset)
		} else {
			v := price.Price.Mul(asset.Amount)
			totalValue = totalValue.Add(v)
		}
		totalDepth = totalDepth.Add(asset.Depth)
	}
	avgDepth := totalDepth.Quo(sdk.NewDec(int64(len(amountDepthInfo))))
	return totalValue, avgDepth, nil
}

func LiquidityRatioFromPriceDepth(depth sdk.Dec) sdk.Dec {
	if depth == sdk.OneDec() {
		return sdk.OneDec()
	}
	sqrt, err := sdk.OneDec().Sub(depth).ApproxSqrt()
	if err != nil {
		panic(err)
	}
	return sdk.OneDec().Sub(sqrt)
}

func (k msgServer) FeedMultipleExternalLiquidity(goCtx context.Context, msg *types.MsgFeedMultipleExternalLiquidity) (*types.MsgFeedMultipleExternalLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	feeder, found := k.oracleKeeper.GetPriceFeeder(ctx, msg.Sender)
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

		tvl, err := pool.TVL(ctx, k.oracleKeeper)
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
		if elRatio.LT(sdk.OneDec()) {
			elRatio = sdk.OneDec()
		}

		pool.PoolParams.ExternalLiquidityRatio = elRatio
		err = k.SetPool(ctx, pool)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgFeedMultipleExternalLiquidityResponse{}, nil
}
