package types

import (
	"errors"
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
)

type AssetWeight struct {
	Asset  string
	Weight elystypes.Dec34
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
			Weight: elystypes.NewDec34FromInt(asset.Weight).Quo(elystypes.NewDec34FromInt(totalWeight)),
		})
	}
	return poolWeights
}

func GetOraclePoolNormalizedWeights(ctx sdk.Context, poolId uint64, oracleKeeper OracleKeeper, poolAssets []PoolAsset) ([]AssetWeight, error) {
	oraclePoolWeights := []AssetWeight{}
	totalWeight := elystypes.ZeroDec34()
	for _, asset := range poolAssets {
		tokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Token.Denom)
		if tokenPrice.IsZero() {
			return oraclePoolWeights, fmt.Errorf("price for token not set: %s", asset.Token.Denom)
		}
		amount := asset.Token.Amount
		weight := tokenPrice.MulInt(amount)
		oraclePoolWeights = append(oraclePoolWeights, AssetWeight{
			Asset:  asset.Token.Denom,
			Weight: weight,
		})
		totalWeight = totalWeight.Add(weight)
	}

	if totalWeight.IsZero() {
		totalWeight = elystypes.OneDec34()
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
			return poolAssets, errors.New("negative pool amount after swap")
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

func (p Pool) WeightDistanceFromTarget(ctx sdk.Context, oracleKeeper OracleKeeper, poolAssets []PoolAsset) elystypes.Dec34 {
	oracleWeights, err := GetOraclePoolNormalizedWeights(ctx, p.PoolId, oracleKeeper, poolAssets)
	if err != nil {
		return elystypes.ZeroDec34()
	}
	targetWeights := NormalizedWeights(poolAssets)

	distanceSum := elystypes.ZeroDec34()
	for i := range poolAssets {
		distance := targetWeights[i].Weight.Sub(oracleWeights[i].Weight).Abs()
		distanceSum = distanceSum.Add(distance)
	}
	// Ensure len(p.PoolAssets) is not zero to avoid division by zero
	if len(p.PoolAssets) == 0 {
		return elystypes.ZeroDec34()
	}
	return distanceSum.QuoInt64(int64(len(p.PoolAssets)))
}

func GetDenomOracleAssetWeight(ctx sdk.Context, poolId uint64, oracleKeeper OracleKeeper, poolAssets []PoolAsset, denom string) elystypes.Dec34 {
	oracleWeights, err := GetOraclePoolNormalizedWeights(ctx, poolId, oracleKeeper, poolAssets)
	if err != nil {
		return elystypes.ZeroDec34()
	}
	for _, weight := range oracleWeights {
		if weight.Asset == denom {
			return weight.Weight
		}
	}
	return elystypes.ZeroDec34()
}

func GetDenomNormalizedWeight(poolAssets []PoolAsset, denom string) elystypes.Dec34 {
	targetWeights := NormalizedWeights(poolAssets)
	for _, weight := range targetWeights {
		if weight.Asset == denom {
			return weight.Weight
		}
	}
	return elystypes.ZeroDec34()
}

func (p Pool) CalcGivenInSlippage(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	accPoolKeeper AccountedPoolKeeper,
) (elystypes.Dec34, error) {
	balancerOutCoin, _, err := p.CalcOutAmtGivenIn(ctx, oracleKeeper, snapshot, tokensIn, tokenOutDenom, sdkmath.LegacyZeroDec(), accPoolKeeper)
	if err != nil {
		return elystypes.ZeroDec34(), err
	}

	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return elystypes.ZeroDec34(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenIn.Denom)
	if inTokenPrice.IsZero() {
		return elystypes.ZeroDec34(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
	if outTokenPrice.IsZero() {
		return elystypes.ZeroDec34(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	oracleOutAmount := elystypes.NewDec34FromInt(tokenIn.Amount).
		Mul(inTokenPrice).
		Quo(outTokenPrice)
	balancerOut := elystypes.NewDec34FromInt(balancerOutCoin.Amount)
	slippageAmount := oracleOutAmount.Sub(balancerOut)
	if slippageAmount.IsNegative() {
		return elystypes.ZeroDec34(), nil
	}
	return slippageAmount, nil
}

// SwapOutAmtGivenIn is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
// weightBreakingFeePerpetualFactor should be 1 if perpetual is not the one calling this function
// Pool, and it's bank balances are updated in keeper.UpdatePoolForSwap
func (p *Pool) SwapOutAmtGivenIn(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee sdkmath.LegacyDec,
	accPoolKeeper AccountedPoolKeeper,
	weightBreakingFeePerpetualFactor sdkmath.LegacyDec,
	params Params,
	takerFees sdkmath.LegacyDec,
) (tokenOut sdk.Coin, slippage, slippageAmount elystypes.Dec34, weightBalanceBonus elystypes.Dec34, oracleOutAmount elystypes.Dec34, swapFeeFinal sdkmath.LegacyDec, err error) {

	// Fixed gas consumption per swap to prevent spam
	ctx.GasMeter().ConsumeGas(BalancerGasFeeForSwap, "balancer swap computation")

	// early return with balancer swap if normal amm pool
	if !p.PoolParams.UseOracle {
		if len(tokensIn) != 1 {
			return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), errors.New("expected tokensIn to be of length one")
		}
		balancerOutCoin, slippage, err := p.CalcOutAmtGivenIn(ctx, oracleKeeper, snapshot, tokensIn, tokenOutDenom, swapFee, accPoolKeeper)
		if err != nil {
			return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), err
		}
		// Pools and it's bank balances are updated in keeper.UpdatePoolForSwap
		return balancerOutCoin, slippage, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), swapFee, nil
	}

	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenIn.Denom)
	if inTokenPrice.IsZero() {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOutDenom)
	if outTokenPrice.IsZero() {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	accountedAssets := p.GetAccountedBalance(ctx, accPoolKeeper, p.PoolAssets)

	// out amount is calculated in this formula
	// balancer slippage amount = Max(oracleOutAmount-balancerOutAmount, 0)
	// resizedAmount = tokenIn / externalLiquidityRatio
	// actualSlippageAmount = balancer slippage(resizedAmount)
	oracleOutAmount = elystypes.NewDec34FromInt(tokenIn.Amount).
		Mul(inTokenPrice).
		Quo(outTokenPrice)

	externalLiquidityRatio, err := p.GetAssetExternalLiquidityRatio(tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), err
	}
	externalLiquidityRatioDec34 := elystypes.NewDec34FromLegacyDec(externalLiquidityRatio)
	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatioDec34.IsZero() {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), ErrAmountTooLow
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
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), err
	}
	outAmountAfterSlippage := oracleOutAmount.Sub(slippageAmount.Mul(externalLiquidityRatioDec34))
	slippageAmount = slippageAmount.Mul(externalLiquidityRatioDec34)
	slippage = slippageAmount.Quo(oracleOutAmount)

	if slippage.LT(elystypes.NewDec34FromLegacyDec(params.MinSlippage)) {
		slippage = elystypes.NewDec34FromLegacyDec(params.MinSlippage)
		slippageAmount = oracleOutAmount.Mul(elystypes.NewDec34FromLegacyDec(params.MinSlippage))
	}

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
		sdk.Coins{sdk.NewCoin(tokenOutDenom, outAmountAfterSlippage.ToInt())}, accountedAssets,
	)
	if err != nil {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), err
	}

	weightBalanceBonus, weightBreakingFee, isSwapFee := p.CalculateWeightFees(ctx, oracleKeeper, accountedAssets, newAssetPools, tokenIn.Denom, params, weightBreakingFeePerpetualFactor)
	if !isSwapFee {
		swapFee = sdkmath.LegacyZeroDec()
	}

	if swapFee.GTE(sdkmath.LegacyOneDec()) {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), ErrTooMuchSwapFee
	}
	tokenAmountOutInt := outAmountAfterSlippage.
		Mul(elystypes.OneDec34().Sub(weightBreakingFee)).
		Mul(elystypes.OneDec34().Sub(elystypes.NewDec34FromLegacyDec(swapFee.Add(takerFees)))).ToInt() // We ignore the decimal component, as we round down the token amount out.
	oracleOutCoin := sdk.NewCoin(tokenOutDenom, tokenAmountOutInt)
	return oracleOutCoin, slippage, slippageAmount, weightBalanceBonus, oracleOutAmount, swapFee, nil
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
