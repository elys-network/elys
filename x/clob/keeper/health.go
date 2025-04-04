package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetHealth(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) (math.LegacyDec, error) {
	liquidationPrice, err := k.GetLiquidationPrice(ctx, perpetual, market)
	if err != nil {
		return math.LegacyDec{}, err
	}
	currentPrice := k.GetCurrentTwapPrice(ctx, market.Id)
	return currentPrice.Quo(liquidationPrice), nil
}
