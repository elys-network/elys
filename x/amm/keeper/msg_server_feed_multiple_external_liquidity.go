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
		price, found := oracleKeeper.GetAssetPrice(ctx, asset.Denom)
		if !found {
			return sdk.ZeroDec(), fmt.Errorf("asset price not set: %s", asset.Denom)
		} else {
			v := price.Price.Mul(asset.Amount)
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

		fmt.Println("tvl", tvl.String())
		fmt.Println("elValue", elValue.String())

		elRatio := elValue.Quo(tvl)
		fmt.Println("elRatio1", elRatio.String(), el.Depth)
		elRatio = elRatio.Quo(LiquidityRatioFromPriceDepth(el.Depth))
		fmt.Println("elRatio2", elRatio.String())
		if elRatio.LT(sdk.OneDec()) {
			elRatio = sdk.OneDec()
		}

		pool.PoolParams.ExternalLiquidityRatio = elRatio
		k.SetPool(ctx, pool)
	}

	return &types.MsgFeedMultipleExternalLiquidityResponse{}, nil
}
