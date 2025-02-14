package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) TVL(ctx sdk.Context, oracleKeeper types.OracleKeeper, baseCurrency string) math.LegacyDec {
	params := k.GetParams(ctx)
	totalDeposit := params.TotalValue
	price := oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	return price.MulInt(totalDeposit)
}
