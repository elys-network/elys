package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) TVL(ctx sdk.Context, oracleKeeper OracleKeeper) (sdk.Dec, error) {
	// OracleAssetsTVL * TotalWeight / OracleAssetsWeight
	// E.g. JUNO / USDT / USDC (30:30:30)
	// TVL = USDC_USDT_liquidity * 90 / 60

	oracleAssetsTVL := sdk.ZeroDec()
	totalWeight := sdk.ZeroInt()
	oracleAssetsWeight := sdk.ZeroInt()
	for _, asset := range p.PoolAssets {
		tokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, asset.Token.Denom)
		fmt.Println("getprice", asset.Token.Denom, tokenPrice.String())
		totalWeight = totalWeight.Add(asset.Weight)
		if tokenPrice.IsZero() {
			if p.PoolParams.UseOracle {
				return sdk.ZeroDec(), fmt.Errorf("token price not set: %s", asset.Token.Denom)
			}
		} else {
			v := tokenPrice.Mul(sdk.NewDecFromInt(asset.Token.Amount))
			oracleAssetsTVL = oracleAssetsTVL.Add(v)
			oracleAssetsWeight = oracleAssetsWeight.Add(asset.Weight)
		}
	}

	fmt.Println("oracleAssetsWeight", oracleAssetsWeight.String(), totalWeight.String())
	if oracleAssetsWeight.IsZero() {
		return sdk.ZeroDec(), nil
	}

	return oracleAssetsTVL.Mul(sdk.NewDecFromInt(totalWeight)).Quo(sdk.NewDecFromInt(oracleAssetsWeight)), nil
}
