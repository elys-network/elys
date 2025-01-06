package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
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
			Aasset.Token.Amount.ToLegacyDec(), Aasset.Weight.ToLegacyDec(),
			Basset.Token.Amount.ToLegacyDec(), Basset.Weight.ToLegacyDec(),
		), nil
	}

	priceA, decimalA := oracleKeeper.GetRawAssetPriceFromDenom(ctx, tokenA)
	if priceA.IsZero() {
		return sdkmath.LegacyZeroDec(), fmt.Errorf("token price not set: %s", tokenA)
	}
	priceB, decimalB := oracleKeeper.GetRawAssetPriceFromDenom(ctx, tokenB)
	if priceB.IsZero() {
		return sdkmath.LegacyZeroDec(), fmt.Errorf("token price not set: %s", tokenB)
	}

	return priceA.Mul(Pow10(decimalB)).Quo(Pow10(decimalA).Mul(priceB)), nil
}
