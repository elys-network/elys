package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// SwapOutAmtGivenIn is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (p *Pool) SwapOutAmtGivenIn(
	ctx sdk.Context,
	tokensIn sdk.Coins,
	tokenOutDenom string,
) (tokenOut sdk.Coin, err error) {
	tokenOutCoin, err := p.CalcOutAmtGivenIn(tokensIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, err
	}

	err = p.applySwap(ctx, tokensIn, sdk.Coins{tokenOutCoin})
	if err != nil {
		return sdk.Coin{}, err
	}
	return tokenOutCoin, nil
}
