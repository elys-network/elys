package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SwapOutAmtGivenIn is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (p *Pool) GetTokenARate(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokenA string,
	tokenB string,
	accPoolKeeper AccountedPoolKeeper,
) (rate sdk.Dec, err error) {
	// balancer pricing if normal amm pool
	if !p.PoolParams.UseOracle {
		Aasset, Basset, err := p.parsePoolAssetsByDenoms(tokenA, tokenB)
		if err != nil {
			return sdk.ZeroDec(), err
		}
		return CalculateTokenARate(
			sdk.NewDecFromInt(Aasset.Token.Amount), sdk.NewDecFromInt(Aasset.Weight),
			sdk.NewDecFromInt(Basset.Token.Amount), sdk.NewDecFromInt(Basset.Weight),
		), nil
	}

	priceA := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenA)
	if priceA.IsZero() {
		return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", tokenA)
	}
	priceB := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenB)
	if priceB.IsZero() {
		return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", tokenB)
	}

	return priceA.Quo(priceB), nil
}
