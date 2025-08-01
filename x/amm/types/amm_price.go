package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/utils"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (p *Pool) GetTokenARate(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	tokenA string,
	tokenB string,
) (rate osmomath.BigDec, err error) {
	// balancer pricing if normal amm pool
	if !p.PoolParams.UseOracle {
		Aasset, Basset, err := p.ParsePoolAssetsByDenoms(tokenA, tokenB)
		if err != nil {
			return osmomath.ZeroBigDec(), errorsmod.Wrapf(err, "failed to parse pool assets for tokens %s and %s", tokenA, tokenB)
		}
		return CalculateTokenARate(
			osmomath.BigDecFromSDKInt(Aasset.Token.Amount), osmomath.BigDecFromSDKInt(Aasset.Weight),
			osmomath.BigDecFromSDKInt(Basset.Token.Amount), osmomath.BigDecFromSDKInt(Basset.Weight),
		), nil
	}

	priceA := oracleKeeper.GetDenomPrice(ctx, tokenA)
	if priceA.IsZero() {
		return osmomath.ZeroBigDec(), errorsmod.Wrapf(ErrInvalidMathApprox, "token price not set or zero for token: %s", tokenA)
	}
	priceB := oracleKeeper.GetDenomPrice(ctx, tokenB)
	if priceB.IsZero() {
		return osmomath.ZeroBigDec(), errorsmod.Wrapf(ErrInvalidMathApprox, "token price not set or zero for token: %s", tokenB)
	}

	return priceA.Quo(priceB), nil
}

func (p *Pool) GetTokenARateNormalized(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	tokenA string,
	tokenB string,
) (rate osmomath.BigDec, err error) {
	// Get the base rate without normalization
	baseRate, err := p.GetTokenARate(ctx, oracleKeeper, tokenA, tokenB)
	if err != nil {
		return osmomath.ZeroBigDec(), err
	}

	// Get token decimals from oracle keeper
	infoA, found := oracleKeeper.GetAssetInfo(ctx, tokenA)
	if !found {
		return osmomath.ZeroBigDec(), fmt.Errorf("asset info not found for token: %s", tokenA)
	}
	infoB, found := oracleKeeper.GetAssetInfo(ctx, tokenB)
	if !found {
		return osmomath.ZeroBigDec(), fmt.Errorf("asset info not found for token: %s", tokenB)
	}

	// Calculate decimal adjustment factor
	decimalDiff := int(infoB.Decimal) - int(infoA.Decimal)
	if decimalDiff > 0 {
		// If tokenB has more decimals, divide by 10^diff
		return baseRate.QuoInt64(utils.Pow10Int64((uint64(decimalDiff)))), nil
	} else if decimalDiff < 0 {
		// If tokenA has more decimals, multiply by 10^|diff|
		return baseRate.MulInt64(utils.Pow10Int64(uint64(-decimalDiff))), nil
	}
	// If decimals are equal, return base rate as is
	return baseRate, nil
}
