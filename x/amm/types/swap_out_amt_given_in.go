package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AssetWeight struct {
	Asset  string
	Weight sdk.Dec
}

func NormalizedWeights(poolAssets []*PoolAsset) []AssetWeight {
	totalWeight := sdk.ZeroInt()
	for _, asset := range poolAssets {
		totalWeight = totalWeight.Add(asset.Weight)
	}
	if totalWeight.IsZero() {
		totalWeight = sdk.OneInt()
	}
	poolWeights := []AssetWeight{}
	for _, asset := range poolAssets {
		poolWeights = append(poolWeights, AssetWeight{
			Asset:  asset.Token.Denom,
			Weight: sdk.NewDecFromInt(asset.Weight).Quo(sdk.NewDecFromInt(totalWeight)),
		})
	}
	return poolWeights
}

func OraclePoolNormalizedWeights(ctx sdk.Context, oracleKeeper OracleKeeper, poolAssets []*PoolAsset) ([]AssetWeight, error) {
	oraclePoolWeights := []AssetWeight{}
	totalWeight := sdk.ZeroDec()
	for _, asset := range poolAssets {
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Token.Denom)
		if tokenPrice.IsZero() {
			return oraclePoolWeights, fmt.Errorf("price for token not set: %s", asset.Token.Denom)
		}

		weight := tokenPrice.Mul(sdk.NewDecFromInt(asset.Token.Amount))
		oraclePoolWeights = append(oraclePoolWeights, AssetWeight{
			Asset:  asset.Token.Denom,
			Weight: weight,
		})
		totalWeight = totalWeight.Add(weight)
	}

	if totalWeight.IsZero() {
		totalWeight = sdk.OneDec()
	}
	for i, asset := range oraclePoolWeights {
		oraclePoolWeights[i].Weight = asset.Weight.Quo(totalWeight)
	}
	return oraclePoolWeights, nil
}

func (p Pool) NewPoolAssetsAfterSwap(inCoins sdk.Coins, outCoins sdk.Coins) (poolAssets []*PoolAsset, err error) {
	for _, asset := range p.PoolAssets {
		denom := asset.Token.Denom
		amountAfterSwap := asset.Token.Amount.Add(inCoins.AmountOf(denom)).Sub(outCoins.AmountOf(denom))
		if amountAfterSwap.IsNegative() {
			return poolAssets, fmt.Errorf("negative pool amount after swap")
		}
		poolAssets = append(poolAssets, &PoolAsset{
			Token:  sdk.NewCoin(denom, amountAfterSwap),
			Weight: asset.Weight,
		})
	}
	return
}

func (p Pool) WeightDistanceFromTarget(ctx sdk.Context, oracleKeeper OracleKeeper, poolAssets []*PoolAsset) sdk.Dec {
	oracleWeights, err := OraclePoolNormalizedWeights(ctx, oracleKeeper, poolAssets)
	if err != nil {
		return sdk.ZeroDec()
	}
	targetWeights := NormalizedWeights(p.PoolAssets)

	distanceSum := sdk.ZeroDec()
	for i := range poolAssets {
		distance := targetWeights[i].Weight.Sub(oracleWeights[i].Weight).Abs()
		distanceSum = distanceSum.Add(distance)
	}
	return distanceSum.Quo(sdk.NewDec(int64(len(p.PoolAssets))))
}

// SwapOutAmtGivenIn is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (p *Pool) SwapOutAmtGivenIn(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee sdk.Dec,
) (tokenOut sdk.Coin, weightBalanceBonus sdk.Dec, err error) {
	balancerOutCoin, err := p.CalcOutAmtGivenIn(tokensIn, tokenOutDenom, swapFee)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}

	// early return with balancer swap if normal amm pool
	if !p.PoolParams.UseOracle {
		err = p.applySwap(ctx, tokensIn, sdk.Coins{balancerOutCoin}, sdk.ZeroDec(), swapFee)
		if err != nil {
			return sdk.Coin{}, sdk.ZeroDec(), err
		}
		return balancerOutCoin, sdk.ZeroDec(), nil
	}

	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenIn.Denom)
	if inTokenPrice.IsZero() {
		return sdk.Coin{}, sdk.ZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
	if outTokenPrice.IsZero() {
		return sdk.Coin{}, sdk.ZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, p.PoolAssets)

	// out amount is calculated in this formula
	// slippage = (oracleOutAmount-balancerOutAmount)*slippageReduction
	// outAmountAfterSlippage = oracleOutAmount - slippage
	// TODO: consider when slippage is minus
	oracleOutAmount := sdk.NewDecFromInt(tokenIn.Amount).Mul(inTokenPrice).Quo(outTokenPrice)
	balancerOutWithoutFee := sdk.NewDecFromInt(balancerOutCoin.Amount).Quo(sdk.OneDec().Sub(swapFee))
	balancerSlippage := oracleOutAmount.Sub(balancerOutWithoutFee)
	slippage := balancerSlippage.Mul(sdk.OneDec().Sub(p.PoolParams.SlippageReduction))
	outAmountAfterSlippage := oracleOutAmount.Sub(slippage)

	// calculate weight distance difference to calculate bonus/cut on the operation
	newAssetPools, err := p.NewPoolAssetsAfterSwap(
		tokensIn,
		sdk.Coins{sdk.NewCoin(tokenOutDenom, outAmountAfterSlippage.TruncateInt())},
	)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// cut is valid when distance higher than original distance
	weightBreakingFee := sdk.ZeroDec()
	if distanceDiff.IsPositive() {
		weightBreakingFee = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff)
	}

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus = sdk.ZeroDec()
	if initialWeightDistance.GT(p.PoolParams.ThresholdWeightDifference) && distanceDiff.IsNegative() {
		weightBalanceBonus = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff).Abs()
		// TODO: we might skip swap fee in case it's a balance recovery operation
	}
	tokenAmountOutInt := outAmountAfterSlippage.
		Mul(sdk.OneDec().Sub(weightBreakingFee)).
		Mul(sdk.OneDec().Sub(swapFee)).TruncateInt()
	oracleOutCoin := sdk.NewCoin(tokenOutDenom, tokenAmountOutInt)
	err = p.applySwap(ctx, tokensIn, sdk.Coins{oracleOutCoin}, sdk.ZeroDec(), swapFee)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}
	return oracleOutCoin, weightBalanceBonus, nil
}
