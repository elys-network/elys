package types

import (
	sdkmath "cosmossdk.io/math"
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AssetWeight struct {
	Asset  string
	Weight sdkmath.LegacyDec
}

func NormalizedWeights(poolAssets []PoolAsset) (poolWeights []AssetWeight) {
	poolWeights = []AssetWeight{}
	totalWeight := sdkmath.ZeroInt()
	for _, asset := range poolAssets {
		totalWeight = totalWeight.Add(asset.Weight)
	}
	if totalWeight.IsZero() {
		totalWeight = sdkmath.OneInt()
	}
	for _, asset := range poolAssets {
		poolWeights = append(poolWeights, AssetWeight{
			Asset:  asset.Token.Denom,
			Weight: sdkmath.LegacyNewDecFromInt(asset.Weight).Quo(sdkmath.LegacyNewDecFromInt(totalWeight)),
		})
	}
	return poolWeights
}

func OraclePoolNormalizedWeights(ctx sdk.Context, oracleKeeper OracleKeeper, poolAssets []PoolAsset) ([]AssetWeight, error) {
	oraclePoolWeights := []AssetWeight{}
	totalWeight := sdkmath.LegacyZeroDec()
	for _, asset := range poolAssets {
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Token.Denom)
		if tokenPrice.IsZero() {
			return oraclePoolWeights, fmt.Errorf("price for token not set: %s", asset.Token.Denom)
		}

		weight := tokenPrice.Mul(sdkmath.LegacyNewDecFromInt(asset.Token.Amount))
		oraclePoolWeights = append(oraclePoolWeights, AssetWeight{
			Asset:  asset.Token.Denom,
			Weight: weight,
		})
		totalWeight = totalWeight.Add(weight)
	}

	if totalWeight.IsZero() {
		totalWeight = sdkmath.LegacyOneDec()
	}
	for i, asset := range oraclePoolWeights {
		oraclePoolWeights[i].Weight = asset.Weight.Quo(totalWeight)
	}
	return oraclePoolWeights, nil
}

func (p Pool) NewPoolAssetsAfterSwap(inCoins sdk.Coins, outCoins sdk.Coins) (poolAssets []PoolAsset, err error) {
	for _, asset := range p.PoolAssets {
		denom := asset.Token.Denom
		amountAfterSwap := asset.Token.Amount.Add(inCoins.AmountOf(denom)).Sub(outCoins.AmountOf(denom))
		if amountAfterSwap.IsNegative() {
			return poolAssets, fmt.Errorf("negative pool amount after swap")
		}
		poolAssets = append(poolAssets, PoolAsset{
			Token:  sdk.NewCoin(denom, amountAfterSwap),
			Weight: asset.Weight,
		})
	}
	return
}

func (p Pool) StackedRatioFromSnapshot(ctx sdk.Context, oracleKeeper OracleKeeper, snapshot *Pool) sdkmath.LegacyDec {
	stackedRatio := sdkmath.LegacyZeroDec()
	for index, asset := range snapshot.PoolAssets {
		assetDiff := sdkmath.LegacyNewDecFromInt(p.PoolAssets[index].Token.Amount.Sub(asset.Token.Amount).Abs())
		// Ensure asset.Token is not zero to avoid division by zero
		if asset.Token.IsZero() {
			asset.Token.Amount = sdkmath.OneInt()
		}
		assetStacked := assetDiff.Quo(sdkmath.LegacyNewDecFromInt(asset.Token.Amount))
		stackedRatio = stackedRatio.Add(assetStacked)
	}

	return stackedRatio
}

func (p Pool) WeightDistanceFromTarget(ctx sdk.Context, oracleKeeper OracleKeeper, poolAssets []PoolAsset) sdkmath.LegacyDec {
	oracleWeights, err := OraclePoolNormalizedWeights(ctx, oracleKeeper, poolAssets)
	if err != nil {
		return sdkmath.LegacyZeroDec()
	}
	targetWeights := NormalizedWeights(p.PoolAssets)

	distanceSum := sdkmath.LegacyZeroDec()
	for i := range poolAssets {
		distance := targetWeights[i].Weight.Sub(oracleWeights[i].Weight).Abs()
		distanceSum = distanceSum.Add(distance)
	}
	// Ensure len(p.PoolAssets) is not zero to avoid division by zero
	if len(p.PoolAssets) == 0 {
		return sdkmath.LegacyZeroDec()
	}
	return distanceSum.Quo(sdkmath.LegacyNewDec(int64(len(p.PoolAssets))))
}

func OracleAssetWeight(ctx sdk.Context, oracleKeeper OracleKeeper, poolAssets []PoolAsset, denom string) sdkmath.LegacyDec {
	oracleWeights, err := OraclePoolNormalizedWeights(ctx, oracleKeeper, poolAssets)
	if err != nil {
		return sdkmath.LegacyZeroDec()
	}
	for _, weight := range oracleWeights {
		if weight.Asset == denom {
			return weight.Weight
		}
	}
	return sdkmath.LegacyZeroDec()
}

func NormalizedWeight(ctx sdk.Context, poolAssets []PoolAsset, denom string) sdkmath.LegacyDec {
	targetWeights := NormalizedWeights(poolAssets)
	for _, weight := range targetWeights {
		if weight.Asset == denom {
			return weight.Weight
		}
	}
	return sdkmath.LegacyZeroDec()
}

func (p Pool) CalcGivenInSlippage(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	accPoolKeeper AccountedPoolKeeper,
) (sdkmath.LegacyDec, error) {
	balancerOutCoin, _, err := p.CalcOutAmtGivenIn(ctx, oracleKeeper, snapshot, tokensIn, tokenOutDenom, sdkmath.LegacyZeroDec(), accPoolKeeper)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenIn.Denom)
	if inTokenPrice.IsZero() {
		return sdkmath.LegacyZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
	if outTokenPrice.IsZero() {
		return sdkmath.LegacyZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	oracleOutAmount := sdkmath.LegacyNewDecFromInt(tokenIn.Amount).Mul(inTokenPrice).Quo(outTokenPrice)
	balancerOut := sdkmath.LegacyNewDecFromInt(balancerOutCoin.Amount)
	slippageAmount := oracleOutAmount.Sub(balancerOut)
	if slippageAmount.IsNegative() {
		return sdkmath.LegacyZeroDec(), nil
	}
	return slippageAmount, nil
}

// SwapOutAmtGivenIn is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (p *Pool) SwapOutAmtGivenIn(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee sdkmath.LegacyDec,
	accPoolKeeper AccountedPoolKeeper,
) (tokenOut sdk.Coin, slippage, slippageAmount sdkmath.LegacyDec, weightBalanceBonus sdkmath.LegacyDec, err error) {
	balancerOutCoin, slippage, err := p.CalcOutAmtGivenIn(ctx, oracleKeeper, snapshot, tokensIn, tokenOutDenom, swapFee, accPoolKeeper)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	// early return with balancer swap if normal amm pool
	if !p.PoolParams.UseOracle {
		err = p.applySwap(ctx, tokensIn, sdk.Coins{balancerOutCoin}, sdkmath.LegacyZeroDec(), swapFee, accPoolKeeper)
		if err != nil {
			return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}
		return balancerOutCoin, slippage, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), nil
	}

	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenIn.Denom)
	if inTokenPrice.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
	if outTokenPrice.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, p.PoolAssets)

	// out amount is calculated in this formula
	// balancer slippage amount = Max(oracleOutAmount-balancerOutAmount, 0)
	// resizedAmount = tokenIn / externalLiquidityRatio
	// actualSlippageAmount = balancer slippage(resizedAmount)
	oracleOutAmount := sdkmath.LegacyNewDecFromInt(tokenIn.Amount).Mul(inTokenPrice).Quo(outTokenPrice)

	// Ensure p.PoolParams.ExternalLiquidityRatio is not zero to avoid division by zero
	if p.PoolParams.ExternalLiquidityRatio.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ErrAmountTooLow
	}

	resizedAmount := sdkmath.LegacyNewDecFromInt(tokenIn.Amount).Quo(p.PoolParams.ExternalLiquidityRatio).RoundInt()
	slippageAmount, err = p.CalcGivenInSlippage(
		ctx,
		oracleKeeper,
		snapshot,
		sdk.Coins{sdk.NewCoin(tokenIn.Denom, resizedAmount)},
		tokenOutDenom,
		accPoolKeeper,
	)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	outAmountAfterSlippage := oracleOutAmount.Sub(slippageAmount)
	slippage = slippageAmount.Quo(oracleOutAmount)

	// oracleOutAmount = 100 ATOM
	// BalancerOutAmount = 95 ATOM
	// balancerSlippageAmount = 5
	// slippageAmount = 5 * (1 - 99%) = 0.05 ATOM
	// Final amount = 99.95 ATOM
	// Osmosis liq=$100 million
	// Elys liq = $1 million
	// reduction = 99% // (100 - 1)/(100)

	// we know swap in amount - 1000 USDC
	// price impact for Osmosis pool - 1000/(50000000 + 1000) = roughly 0.002%
	// balancer price impact - balancerSlippageAmount / oracleOutAmount = 5%
	// 0.002% / 5% = 0.0004 != 0.01 (slippage reduction factor) (right?)

	// Elys normal amm pool = Osmosis normal amm pool (80/20 pool,
	// we can create same virtual pool on Elys and calculate slippage)

	// actual out amount = oracle out amount - slippage(Osmosis)

	// Oracle price
	// 1% depth
	// $1mil
	// Price impact for $1000
	// 0.001% - price impact
	// Out amount = (oracleOutAmount*(1-0.001%))
	// First $100, 0.0001%
	// For second $100, 0.0002%
	// Triangle in pricing
	// in amount = 100 ATOM
	// linear model USDC/USDT stable pool, BTC/USDC
	// Assume: it's linear model
	// out amount = ? USDC
	// Formula to calculate out amount
	// We won't use Elys pool data here
	// Reduction 98% - 99.9%
	// Slippage reduction is dynamic based on trade size
	// approximate value = slippage reduction
	// Dream's solution:
	// Dynamic slippage reduction
	// $1000 trade: 95%
	// $10000 trade: 80%

	// calculate weight distance difference to calculate bonus/cut on the operation
	newAssetPools, err := p.NewPoolAssetsAfterSwap(
		tokensIn,
		sdk.Coins{sdk.NewCoin(tokenOutDenom, outAmountAfterSlippage.TruncateInt())},
	)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// cut is valid when distance higher than original distance
	weightBreakingFee := sdkmath.LegacyZeroDec()
	if distanceDiff.IsPositive() {
		// old weight breaking fee implementation
		// weightBreakingFee = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff)

		// target weight
		targetWeightIn := NormalizedWeight(ctx, p.PoolAssets, tokenIn.Denom)
		targetWeightOut := NormalizedWeight(ctx, p.PoolAssets, tokenOutDenom)

		// weight breaking fee as in Plasma pool
		weightIn := OracleAssetWeight(ctx, oracleKeeper, newAssetPools, tokenIn.Denom)
		weightOut := OracleAssetWeight(ctx, oracleKeeper, newAssetPools, tokenOutDenom)
		weightBreakingFee = GetWeightBreakingFee(weightIn, weightOut, targetWeightIn, targetWeightOut, p.PoolParams)
	}

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus = weightBreakingFee.Neg()
	if initialWeightDistance.GT(p.PoolParams.ThresholdWeightDifference) && distanceDiff.IsNegative() {
		weightBalanceBonus = p.PoolParams.WeightBreakingFeeMultiplier.Mul(distanceDiff).Abs()
	}

	if swapFee.GTE(sdkmath.LegacyOneDec()) {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ErrTooMuchSwapFee
	}

	tokenAmountOutInt := outAmountAfterSlippage.
		Mul(sdkmath.LegacyOneDec().Sub(weightBreakingFee)).
		Mul(sdkmath.LegacyOneDec().Sub(swapFee)).TruncateInt()
	oracleOutCoin := sdk.NewCoin(tokenOutDenom, tokenAmountOutInt)
	err = p.applySwap(ctx, tokensIn, sdk.Coins{oracleOutCoin}, sdkmath.LegacyZeroDec(), swapFee, accPoolKeeper)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	return oracleOutCoin, slippage, slippageAmount, weightBalanceBonus, nil
}
