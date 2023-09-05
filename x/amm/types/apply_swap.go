package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// ApplySwap.
func (p *Pool) applySwap(ctx sdk.Context, tokensIn sdk.Coins, tokensOut sdk.Coins, swapFeeIn, swapFeeOut sdk.Dec, accPoolKeeper AccountedPoolKeeper) error {
	// Fixed gas consumption per swap to prevent spam
	ctx.GasMeter().ConsumeGas(BalancerGasFeeForSwap, "balancer swap computation")
	// Also ensures that len(tokensIn) = 1 = len(tokensOut)
	inPoolAsset, outPoolAsset, err := p.parsePoolAssetsCoins(tokensIn, tokensOut)
	if err != nil {
		return err
	}
	inTokensAfterFee := sdk.NewDecFromInt(tokensIn[0].Amount).Mul(sdk.OneDec().Sub(swapFeeIn)).TruncateInt()
	outTokensAfterFee := sdk.NewDecFromInt(tokensOut[0].Amount).Mul(sdk.OneDec().Sub(swapFeeOut)).TruncateInt()
	inPoolAsset.Token.Amount = inPoolAsset.Token.Amount.Add(inTokensAfterFee)
	outPoolAsset.Token.Amount = outPoolAsset.Token.Amount.Sub(outTokensAfterFee)

	return p.UpdatePoolAssetBalances(sdk.NewCoins(
		inPoolAsset.Token,
		outPoolAsset.Token,
	))
}
