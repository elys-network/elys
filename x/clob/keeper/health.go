package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/clob/types"
)

// GetHealth Values from 0 to infinity, at 0, liquidation should happen as currentPrice == liquidationPrice
// if it's < 0, it means liquidation wasn't done on time
func (k Keeper) GetHealth(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) (math.LegacyDec, error) {
	liquidationPrice, err := k.GetLiquidationPrice(ctx, perpetual, market)
	if err != nil {
		return math.LegacyDec{}, err
	}
	currentPrice, err := k.GetAssetPriceFromDenom(ctx, market.BaseDenom)
	if err != nil {
		return math.LegacyDec{}, err
	}
	return currentPrice.Quo(liquidationPrice).Sub(math.LegacyOneDec()), nil
}
