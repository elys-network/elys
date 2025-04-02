package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetFundingRate(ctx sdk.Context, marketId uint64) types.FundingRate {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetFundingRateKey(marketId)
	b := store.Get(key)
	if b == nil {
		return types.FundingRate{
			MarketId: marketId,
			Block:    uint64(ctx.BlockHeight()),
			Rate:     utils.ZeroDec,
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

//premium = TWAP(markPrice) - TWAP(indexPrice)

// fundingRate = clamp(premium / indexPrice, -cap, +cap)
func (k Keeper) UpdateFundingRate(ctx sdk.Context, market types.PerpetualMarket) error {
	twapMarkPrice, err := k.GetCurrentTwapPrice(ctx, market.Id)
	if err != nil {
		return err
	}
	assetInfo, found := k.oracleKeeper.GetAssetInfo(ctx, market.BaseDenom)
	if !found {
		return fmt.Errorf("asset info (%s) not found", market.BaseDenom)
	}
	oraclePrice, found := k.oracleKeeper.GetAssetPrice(ctx, assetInfo.Display)
	if !found {
		return fmt.Errorf("asset price (%s) not found", assetInfo.Display)
	}
	indexPrice, err := math.DecFromLegacyDec(oraclePrice.Price)
	if err != nil {
		return err
	}
	premium, err := twapMarkPrice.Sub(indexPrice)
	if err != nil {
		return err
	}
	fundingRateCal, err := premium.Quo(indexPrice)
	if err != nil {
		return err
	}
	lastFundingRate := k.GetFundingRate(ctx, market.Id)
	change, err := fundingRateCal.Sub(lastFundingRate.Rate)
	if err != nil {
		return err
	}
	if !change.IsZero() {
		if change.IsPositive() && change.Cmp(market.MaxFundingRateChange) > 0 {
			change = market.MaxFundingRateChange
		}
		if change.IsNegative() && utils.Abs(change).Cmp(market.MaxFundingRateChange) > 0 {
			change = utils.Neg(market.MaxFundingRateChange)
		}
	}
	lastFundingRate.Rate, err = lastFundingRate.Rate.Add(change)
	if err != nil {
		return err
	}
	lastFundingRate.Block = uint64(ctx.BlockHeight())
	k.SetFundingRate(ctx, lastFundingRate)
	return nil
}
