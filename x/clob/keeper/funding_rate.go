package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/clob/types"
)

func (k Keeper) GetFundingRate(ctx sdk.Context, marketId uint64) types.FundingRate {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetFundingRateKey(marketId)
	b := store.Get(key)
	if b == nil {
		return types.FundingRate{
			MarketId: marketId,
			Block:    uint64(ctx.BlockHeight()),
			Rate:     math.LegacyZeroDec(),
		}
	}

	var v types.FundingRate
	k.cdc.MustUnmarshal(b, &v)
	return v
}

func (k Keeper) SetFundingRate(ctx sdk.Context, p types.FundingRate) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetFundingRateKey(p.MarketId)
	b := k.cdc.MustMarshal(&p)
	store.Set(key, b)
}

func (k Keeper) GetAllFundingRate(ctx sdk.Context) []types.FundingRate {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.FundingRatePrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.FundingRate

	for ; iterator.Valid(); iterator.Next() {
		var val types.FundingRate
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

// UpdateFundingRate
// premium = TWAP(markPrice) - TWAP(indexPrice)
// fundingRate = clamp(premium / indexPrice, -cap, +cap)
func (k Keeper) UpdateFundingRate(ctx sdk.Context, market types.PerpetualMarket) error {
	twapMarkPrice, err := k.GetCurrentTwapPrice(ctx, market.Id)
	if err != nil {
		return err
	}

	// Additional nil safety check
	if twapMarkPrice.IsNil() {
		return fmt.Errorf("twap mark price is nil for market %d", market.Id)
	}

	indexPrice, err := k.GetAssetPriceFromDenom(ctx, market.BaseDenom)
	if err != nil {
		return err
	}

	// Additional nil/zero safety check for division
	if indexPrice.IsNil() || indexPrice.IsZero() {
		return fmt.Errorf("index price is nil or zero for market %d, cannot calculate funding rate", market.Id)
	}

	premium := twapMarkPrice.Sub(indexPrice)
	fundingRateCal := premium.Quo(indexPrice)

	lastFundingRate := k.GetFundingRate(ctx, market.Id)
	change := fundingRateCal.Sub(lastFundingRate.Rate)

	if !change.IsZero() {
		if change.IsPositive() && change.GT(market.MaxAbsFundingRateChange) {
			change = market.MaxAbsFundingRateChange
		}
		if change.IsNegative() && change.Abs().GT(market.MaxAbsFundingRateChange) {
			change = market.MaxAbsFundingRateChange.Neg()
		}
	}
	lastFundingRate.Rate = lastFundingRate.Rate.Add(change)
	if !lastFundingRate.Rate.IsZero() && lastFundingRate.Rate.Abs().GT(market.MaxAbsFundingRate) {
		positive := lastFundingRate.Rate.IsPositive()
		if positive {
			lastFundingRate.Rate = market.MaxAbsFundingRate
		} else {
			lastFundingRate.Rate = market.MaxAbsFundingRate.Neg()
		}
	}
	previousRate := lastFundingRate.Rate.Sub(change)
	lastFundingRate.Block = uint64(ctx.BlockHeight())
	k.SetFundingRate(ctx, lastFundingRate)

	// Emit funding rate update event
	k.EmitFundingRateUpdateEvent(ctx, market.Id, previousRate, lastFundingRate.Rate)

	return nil
}
