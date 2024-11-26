package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) TVL(ctx sdk.Context, oracleKeeper types.OracleKeeper, baseCurrency string) math.LegacyDec {
	params := k.GetParams(ctx)
	totalDeposit := params.TotalValue
	price, denomDecimal := oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	return price.MulInt(totalDeposit).Quo(oraclekeeper.Pow10(denomDecimal))
}

func (k Keeper) ShareDenomPrice(ctx sdk.Context, oracleKeeper types.OracleKeeper, baseCurrency string) math.LegacyDec {
	params := k.GetParams(ctx)
	redemptionRate := params.RedemptionRate
	price, denomDecimal := oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	return price.Mul(redemptionRate).Quo(oraclekeeper.Pow10(denomDecimal))
}
