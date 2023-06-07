package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// SwapInAmtGivenOut is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (p *Pool) SwapInAmtGivenOut(
	ctx sdk.Context, tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec) (
	tokenIn sdk.Coin, err error,
) {
	tokenInCoin, err := p.CalcInAmtGivenOut(tokensOut, tokenInDenom, swapFee)
	if err != nil {
		return sdk.Coin{}, err
	}

	err = p.applySwap(ctx, sdk.Coins{tokenInCoin}, tokensOut)
	if err != nil {
		return sdk.Coin{}, err
	}
	return tokenInCoin, nil
}
