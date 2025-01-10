package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
)

// SwapOutAmtGivenIn is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (p *Pool) GetTokenARate(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokenA string,
	tokenB string,
	accPoolKeeper AccountedPoolKeeper,
) (rate elystypes.Dec34, err error) {
	// balancer pricing if normal amm pool
	if !p.PoolParams.UseOracle {
		Aasset, Basset, err := p.parsePoolAssetsByDenoms(tokenA, tokenB)
		if err != nil {
			return elystypes.ZeroDec34(), err
		}
		return elystypes.NewDec34FromLegacyDec(CalculateTokenARate(
			Aasset.Token.Amount.ToLegacyDec(), Aasset.Weight.ToLegacyDec(),
			Basset.Token.Amount.ToLegacyDec(), Basset.Weight.ToLegacyDec(),
		)), nil
	}

	priceA, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenA)
	if priceA.IsZero() {
		return elystypes.ZeroDec34(), fmt.Errorf("token price not set: %s", tokenA)
	}
	priceB, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenB)
	if priceB.IsZero() {
		return elystypes.ZeroDec34(), fmt.Errorf("token price not set: %s", tokenB)
	}

	return priceA.Quo(priceB), nil
}
