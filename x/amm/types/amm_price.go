package types

import (
	sdkmath "cosmossdk.io/math"
	"fmt"
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
) (rate sdkmath.LegacyDec, err error) {
	// balancer pricing if normal amm pool
	if !p.PoolParams.UseOracle {
		Aasset, Basset, err := p.parsePoolAssetsByDenoms(tokenA, tokenB)
		if err != nil {
			return sdkmath.LegacyZeroDec(), err
		}
		return CalculateTokenARate(
			sdkmath.LegacyNewDecFromInt(Aasset.Token.Amount), sdkmath.LegacyNewDecFromInt(Aasset.Weight),
			sdkmath.LegacyNewDecFromInt(Basset.Token.Amount), sdkmath.LegacyNewDecFromInt(Basset.Weight),
		), nil
	}

	priceA := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenA)
	if priceA.IsZero() {
		return sdkmath.LegacyZeroDec(), fmt.Errorf("token price not set: %s", tokenA)
	}
	priceB := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenB)
	if priceB.IsZero() {
		return sdkmath.LegacyZeroDec(), fmt.Errorf("token price not set: %s", tokenB)
	}

	return priceA.Quo(priceB), nil
}
