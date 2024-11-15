package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
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

func GetOraclePoolNormalizedWeights(ctx sdk.Context, poolId uint64, oracleKeeper OracleKeeper, poolAssets []PoolAsset) ([]AssetWeight, error) {
	oraclePoolWeights := []AssetWeight{}
	totalWeight := sdkmath.LegacyZeroDec()
	for _, asset := range poolAssets {
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Token.Denom)
		if tokenPrice.IsZero() {
			return oraclePoolWeights, fmt.Errorf("price for token not set: %s", asset.Token.Denom)
		}
		amount := asset.Token.Amount
		weight := amount.ToLegacyDec().Mul(tokenPrice)
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

func (p Pool) NewPoolAssetsAfterSwap(ctx sdk.Context, inCoins sdk.Coins, outCoins sdk.Coins, poolAssets []PoolAsset) ([]PoolAsset, error) {
	updatedAssets := []PoolAsset{}
	for _, asset := range poolAssets {
		denom := asset.Token.Denom
		beforeAmount := asset.Token.Amount
		amountAfterSwap := beforeAmount.Add(inCoins.AmountOf(denom)).Sub(outCoins.AmountOf(denom))
		if amountAfterSwap.IsNegative() {
			return poolAssets, fmt.Errorf("negative pool amount after swap")
		}
		updatedAssets = append(updatedAssets, PoolAsset{
			Token:  sdk.NewCoin(denom, amountAfterSwap),
			Weight: asset.Weight,
		})
	}
	return updatedAssets, nil
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
	oracleWeights, err := GetOraclePoolNormalizedWeights(ctx, p.PoolId, oracleKeeper, poolAssets)
	if err != nil {
		return sdkmath.LegacyZeroDec()
	}
	targetWeights := NormalizedWeights(poolAssets)

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

func GetDenomOracleAssetWeight(ctx sdk.Context, poolId uint64, oracleKeeper OracleKeeper, poolAssets []PoolAsset, denom string) sdkmath.LegacyDec {
	oracleWeights, err := GetOraclePoolNormalizedWeights(ctx, poolId, oracleKeeper, poolAssets)
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

func GetDenomNormalizedWeight(poolAssets []PoolAsset, denom string) sdkmath.LegacyDec {
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
// weightBreakingFeePerpetualFactor should be 1 if perpetual is not the one calling this function
func (p *Pool) SwapOutAmtGivenIn(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee sdkmath.LegacyDec,
	accPoolKeeper AccountedPoolKeeper,
	weightBreakingFeePerpetualFactor math.LegacyDec,
	params Params,
) (tokenOut sdk.Coin, slippage, slippageAmount sdkmath.LegacyDec, weightBalanceBonus sdkmath.LegacyDec, err error) {
	balancerOutCoin, slippage, err := p.CalcOutAmtGivenIn(ctx, oracleKeeper, snapshot, tokensIn, tokenOutDenom, swapFee, accPoolKeeper)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	// early return with balancer swap if normal amm pool
	if !p.PoolParams.UseOracle {
		err = p.applySwap(ctx, tokensIn, sdk.Coins{balancerOutCoin}, sdkmath.LegacyZeroDec(), swapFee)
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

	accountedAssets := p.GetAccountedBalance(ctx, accPoolKeeper, p.PoolAssets)
	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, accountedAssets)

	// out amount is calculated in this formula
	// balancer slippage amount = Max(oracleOutAmount-balancerOutAmount, 0)
	// resizedAmount = tokenIn / externalLiquidityRatio
	// actualSlippageAmount = balancer slippage(resizedAmount)
	oracleOutAmount := sdkmath.LegacyNewDecFromInt(tokenIn.Amount).Mul(inTokenPrice).Quo(outTokenPrice)

	externalLiquidityRatio, err := p.GetAssetExternalLiquidityRatio(tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatio.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ErrAmountTooLow
	}

	resizedAmount := sdkmath.LegacyNewDecFromInt(tokenIn.Amount).Quo(externalLiquidityRatio).RoundInt()
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
	newAssetPools, err := p.NewPoolAssetsAfterSwap(ctx,
		tokensIn,
		sdk.Coins{sdk.NewCoin(tokenOutDenom, outAmountAfterSlippage.TruncateInt())}, accountedAssets,
	)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// target weight
	// Asset weight remains same in new pool assets as in original pool assets
	targetWeightIn := GetDenomNormalizedWeight(newAssetPools, tokenIn.Denom)
	targetWeightOut := GetDenomNormalizedWeight(newAssetPools, tokenOutDenom)

	// weight breaking fee as in Plasma pool
	weightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, newAssetPools, tokenIn.Denom)
	weightOut := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, newAssetPools, tokenOutDenom)
	weightBreakingFee := GetWeightBreakingFee(weightIn, weightOut, targetWeightIn, targetWeightOut, distanceDiff, params)

	// weightBreakingFeePerpetualFactor is 1 if not send by perpetual
	weightBreakingFee = weightBreakingFee.Mul(weightBreakingFeePerpetualFactor)

	// weight recovery reward = weight breaking fee * weight recovery fee portion
	weightRecoveryReward := weightBreakingFee.Mul(params.WeightRecoveryFeePortion)

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus = weightBreakingFee.Neg()

	// If swap is improving weight, set weight breaking fee to zero
	if distanceDiff.IsNegative() {
		weightBreakingFee = sdkmath.LegacyZeroDec()
		weightBalanceBonus = sdkmath.LegacyZeroDec()

		// set weight breaking fee to zero if bonus is applied
		if initialWeightDistance.GT(params.ThresholdWeightDifference) {
			weightBalanceBonus = weightRecoveryReward
		}
	}

	if swapFee.GTE(sdkmath.LegacyOneDec()) {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ErrTooMuchSwapFee
	}

	tokenAmountOutInt := outAmountAfterSlippage.
		Mul(sdkmath.LegacyOneDec().Sub(weightBreakingFee)).
		Mul(sdkmath.LegacyOneDec().Sub(swapFee)).TruncateInt()
	oracleOutCoin := sdk.NewCoin(tokenOutDenom, tokenAmountOutInt)
	err = p.applySwap(ctx, tokensIn, sdk.Coins{oracleOutCoin}, sdkmath.LegacyZeroDec(), swapFee)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	return oracleOutCoin, slippage, slippageAmount, weightBalanceBonus, nil
}

// TODO: Ideally we should have a single DS for accounted pool to avoid confusion
// Or refactor/improve amm pool to use accounted pool
// Task has been added
func (p *Pool) GetAccountedBalance(ctx sdk.Context, accountedPoolKeeper AccountedPoolKeeper, poolAssets []PoolAsset) (updatedAssets []PoolAsset) {
	for _, asset := range poolAssets {
		if p.PoolParams.UseOracle {
			accountedPoolAmt := accountedPoolKeeper.GetAccountedBalance(ctx, p.PoolId, asset.Token.Denom)
			if accountedPoolAmt.IsPositive() {
				asset.Token.Amount = accountedPoolAmt
			}
		}
		updatedAssets = append(updatedAssets, asset)
	}
	return updatedAssets
}
