package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) TVL(ctx sdk.Context, oracleKeeper types.OracleKeeper, baseCurrency string) elystypes.Dec34 {
	params := k.GetParams(ctx)
	totalDeposit := params.TotalValue
	price, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	return price.MulInt(totalDeposit)
}

func (k Keeper) ShareDenomPrice(ctx sdk.Context, oracleKeeper types.OracleKeeper, baseCurrency string) elystypes.Dec34 {
	params := k.GetParams(ctx)
	redemptionRate := params.RedemptionRate
	price, _ := oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	return price.MulLegacyDec(redemptionRate)
}
