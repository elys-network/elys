package types

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (p Pool) CalcExitValueWithSlippage(ctx sdk.Context, oracleKeeper OracleKeeper, accPoolKeeper AccountedPoolKeeper,
	snapshot SnapshotPool, exitingShares sdkmath.Int, tokenOutDenom string,
	weightMultiplier osmomath.BigDec, applyFee bool, params Params) (osmomath.BigDec, osmomath.BigDec, sdk.Coins, sdk.Coin, error) {
	tvl, err := p.TVL(ctx, oracleKeeper, accPoolKeeper)
	if err != nil {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, sdk.Coin{}, err
	}

	// As this is 2 token pool, tokenOut will be
	tokenInDenom := ""
	for _, asset := range p.PoolAssets {
		if asset.Token.Denom == tokenOutDenom {
			continue
		}
		tokenInDenom = asset.Token.Denom
	}
	// Not possible, but we might require this when we have pools with assets more than 2
	if tokenInDenom == "" {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, sdk.Coin{}, fmt.Errorf("token in denom not found")
	}

	totalShares := p.GetTotalShares()
	refundedShares := osmomath.BigDecFromSDKInt(exitingShares)

	// Ensure totalShares is not zero to avoid division by zero
	if totalShares.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, sdk.Coin{}, ErrAmountTooLow
	}

	exitValue := tvl.Mul(refundedShares).Quo(osmomath.BigDecFromSDKInt(totalShares.Amount))

	inTokenPrice := oracleKeeper.GetDenomPrice(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, sdk.Coin{}, fmt.Errorf("token price not set: %s", tokenInDenom)
	}

	// tokenIn amount will be
	tokenInAmount := exitValue.Quo(inTokenPrice)
	weightedAmount := tokenInAmount.Mul(weightMultiplier)

	weightedInCoin := sdk.NewCoin(tokenInDenom, weightedAmount.Dec().TruncateInt())

	if !applyFee {
		return exitValue, osmomath.ZeroBigDec(), sdk.Coins{}, weightedInCoin, nil
	}

	outTokenPrice := oracleKeeper.GetDenomPrice(ctx, tokenOutDenom)
	if outTokenPrice.IsZero() {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, sdk.Coin{}, fmt.Errorf("token price not set: %s", tokenOutDenom)
	}

	externalLiquidityRatio, err := p.GetAssetExternalLiquidityRatio(tokenOutDenom)
	if err != nil {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, sdk.Coin{}, err
	}
	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatio.LT(osmomath.OneBigDec()) {
		externalLiquidityRatio = osmomath.OneBigDec()
	}

	resizedAmount := osmomath.BigDecFromSDKInt(weightedAmount.Dec().TruncateInt()).
		Quo(externalLiquidityRatio).Dec().RoundInt()
	slippageAmount, _, err := p.CalcGivenInSlippage(
		ctx,
		oracleKeeper,
		snapshot,
		sdk.Coins{sdk.NewCoin(tokenInDenom, resizedAmount)},
		tokenOutDenom,
	)
	if err != nil {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, sdk.Coin{}, err
	}
	slippageAmount = slippageAmount.Mul(externalLiquidityRatio)

	slippageValue := slippageAmount.Mul(outTokenPrice)
	slippage := slippageValue.Quo(exitValue)

	minSlippage := params.GetBigDecMinSlippage().Mul(weightMultiplier)
	if slippage.LT(minSlippage) {
		slippage = minSlippage
		slippageValue = exitValue.Mul(minSlippage)
	}

	exitValueWithSlippage := exitValue.Sub(slippageValue)

	if exitingShares.GTE(totalShares.Amount) {
		return osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, sdk.Coin{}, errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
	}

	slippageCoins := sdk.Coins{sdk.NewCoin(tokenOutDenom, slippageAmount.Dec().TruncateInt())}

	return exitValueWithSlippage, slippage, slippageCoins, weightedInCoin, nil
}

// CalcExitPool returns how many tokens should come out, when exiting k LP shares against a "standard" CFMM
func (p Pool) CalcExitPool(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot SnapshotPool,
	accountedPoolKeeper AccountedPoolKeeper,
	exitingShares sdkmath.Int,
	tokenOutDenom string,
	params Params,
	takerFees osmomath.BigDec,
	applyFee bool,
) (exitCoins sdk.Coins, weightBalanceBonus osmomath.BigDec, slippage osmomath.BigDec, swapFee osmomath.BigDec, takerFeesFinal osmomath.BigDec, slippageCoins sdk.Coins, swapInfos []SwapInfo, err error) {
	totalShares := p.GetTotalShares()
	if exitingShares.GTE(totalShares.Amount) {
		return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, nil, errorsmod.Wrapf(ErrLimitMaxAmount, ErrMsgFormatSharesLargerThanMax, exitingShares, totalShares)
	}

	// refundedShares = exitingShares * (1 - exit fee)
	// with 0 exit fee optimization
	refundedShares := osmomath.BigDecFromSDKInt(exitingShares)

	shareOutRatio := refundedShares.Quo(osmomath.BigDecFromSDKInt(totalShares.Amount))
	// exitedCoins = shareOutRatio * pool liquidity
	exitedCoins := sdk.Coins{}

	if p.PoolParams.UseOracle && tokenOutDenom != "" {

		tokenPrice := oracleKeeper.GetDenomPrice(ctx, tokenOutDenom)

		initialWeightOut := GetDenomOracleAssetWeight(ctx, oracleKeeper, snapshot.PoolAssets, tokenOutDenom)
		initialWeightIn := osmomath.OneBigDec().Sub(initialWeightOut)

		exitValueWithSlippage, slippage, slippageCoins, swappedWeightedTokenIn, err := p.CalcExitValueWithSlippage(ctx, oracleKeeper, accountedPoolKeeper, snapshot, exitingShares, tokenOutDenom, initialWeightIn, applyFee, params)
		if err != nil {
			return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, nil, err
		}

		tokenOut, _, _, err := p.CalcOutAmtGivenIn(ctx, oracleKeeper, snapshot, sdk.Coins{swappedWeightedTokenIn}, tokenOutDenom, osmomath.ZeroBigDec())
		if err != nil {
			return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, nil, err
		}

		swapInfos = append(swapInfos, NewSwapInfo(swappedWeightedTokenIn, tokenOut))

		// Ensure tokenPrice is not zero to avoid division by zero
		if tokenPrice.IsZero() {
			return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, nil, ErrAmountTooLow
		}

		oracleOutAmount := exitValueWithSlippage.Quo(tokenPrice)

		tokenOutAmount := oracleOutAmount.Dec().RoundInt()
		weightBalanceBonus = osmomath.ZeroBigDec()
		takerFeesFinal = osmomath.ZeroBigDec()
		isSwapFee := true
		swapFee = osmomath.ZeroBigDec()

		if applyFee {
			newAssetPools, err := p.NewPoolAssetsAfterSwap(
				sdk.Coins{},
				sdk.Coins{sdk.NewCoin(tokenOutDenom, oracleOutAmount.Dec().RoundInt())}, snapshot.PoolAssets,
			)
			if err != nil {
				return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, nil, err
			}
			var tokenInDenom string
			for _, asset := range newAssetPools {
				if asset.Token.Amount.IsNegative() {
					return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, nil, errors.New("out amount exceeds liquidity balance")
				}

				// As we have two asset pool so other asset will be tokenIn
				if asset.Token.Denom != tokenOutDenom {
					tokenInDenom = asset.Token.Denom
				}
			}

			var weightBreakingFee osmomath.BigDec
			weightBalanceBonus, weightBreakingFee, isSwapFee = p.CalculateWeightFees(ctx, oracleKeeper, snapshot.PoolAssets, newAssetPools, tokenInDenom, params, osmomath.OneBigDec())
			// apply percentage to fees, consider improvement or reduction of other token
			// Other denom weight ratio to reduce the weight breaking fees
			weightBreakingFee = weightBreakingFee.Mul(initialWeightIn)
			weightBalanceBonus = weightBalanceBonus.Mul(initialWeightIn)

			if isSwapFee {
				swapFee = p.GetPoolParams().GetBigDecSwapFee().Mul(initialWeightIn)
			}

			takerFeesFinal = takerFees.Mul(initialWeightIn)

			tokenOutAmount = (oracleOutAmount.
				Mul(osmomath.OneBigDec().Sub(weightBreakingFee)).
				Mul(osmomath.OneBigDec().Sub(swapFee.Add(takerFeesFinal)))).Dec().RoundInt()
		}

		return sdk.Coins{sdk.NewCoin(tokenOutDenom, tokenOutAmount)}, weightBalanceBonus, slippage, swapFee, takerFeesFinal, slippageCoins, swapInfos, nil
	}

	// Real balances
	poolLiquidity := p.GetTotalPoolLiquidity()

	for _, accountedAsset := range snapshot.PoolAssets {
		// round down here, due to not wanting to over-exit
		exitAmt := shareOutRatio.Mul(osmomath.BigDecFromSDKInt(accountedAsset.Token.Amount)).Dec().TruncateInt()
		if exitAmt.LTE(sdkmath.ZeroInt()) {
			continue
		}
		for _, pooledAsset := range poolLiquidity {
			if pooledAsset.Denom == accountedAsset.Token.Denom && exitAmt.GTE(pooledAsset.Amount) {
				return sdk.Coins{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, nil, errors.New("too many shares out")
			}
		}
		exitedCoins = exitedCoins.Add(sdk.NewCoin(accountedAsset.Token.Denom, exitAmt))
	}

	return exitedCoins, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), sdk.Coins{}, swapInfos, nil
}
