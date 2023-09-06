package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
)

func AssetsValue(ctx sdk.Context, oracleKeeper types.OracleKeeper, elCoins sdk.DecCoins) (sdk.Dec, error) {
	totalValue := sdk.ZeroDec()
	for _, asset := range elCoins {
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Denom)
		if tokenPrice.IsZero() {
			return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", asset.Denom)
		} else {
			v := tokenPrice.Mul(asset.Amount)
			totalValue = totalValue.Add(v)
		}
	}
	return totalValue, nil
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

		elValue, err := AssetsValue(ctx, k.oracleKeeper, el.Amounts)
		if err != nil {
			return nil, err
		}

		elRatio := elValue.Quo(tvl)
		elRatio = elRatio.Quo(LiquidityRatioFromPriceDepth(el.Depth))
		if elRatio.LT(sdk.OneDec()) {
			elRatio = sdk.OneDec()
		}

		pool.PoolParams.ExternalLiquidityRatio = elRatio
		k.SetPool(ctx, pool)
	}

	return &types.MsgFeedMultipleExternalLiquidityResponse{}, nil
}
