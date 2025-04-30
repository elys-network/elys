package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// SwapOutAmtGivenIn is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (p *Pool) GetTokenARate(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokenA string,
	tokenB string,
	accPoolKeeper AccountedPoolKeeper,
) (rate osmomath.BigDec, err error) {
	// balancer pricing if normal amm pool
	if !p.PoolParams.UseOracle {
		Aasset, Basset, err := p.parsePoolAssetsByDenoms(tokenA, tokenB)
		if err != nil {
			return osmomath.ZeroBigDec(), err
		}
		return CalculateTokenARate(
			osmomath.BigDecFromSDKInt(Aasset.Token.Amount), osmomath.BigDecFromSDKInt(Aasset.Weight),
			osmomath.BigDecFromSDKInt(Basset.Token.Amount), osmomath.BigDecFromSDKInt(Basset.Weight),
		), nil
	}

	priceA := oracleKeeper.GetDenomPrice(ctx, tokenA)
	if priceA.IsZero() {
		return osmomath.ZeroBigDec(), fmt.Errorf("token price not set: %s", tokenA)
	}
	priceB := oracleKeeper.GetDenomPrice(ctx, tokenB)
	if priceB.IsZero() {
		return osmomath.ZeroBigDec(), fmt.Errorf("token price not set: %s", tokenB)
	}

	return priceA.Quo(priceB), nil
}
