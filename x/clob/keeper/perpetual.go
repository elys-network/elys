package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetPerpetual(ctx sdk.Context, marketId, id uint64) (types.Perpetual, error) {
	key := types.GetPerpetualKey(marketId, id)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		return types.Perpetual{}, types.ErrPerpetualNotFound
	}

	var v types.Perpetual
	k.cdc.MustUnmarshal(b, &v)
	return v, nil
}

func (k Keeper) GetAllPerpetuals(ctx sdk.Context) []types.Perpetual {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.Perpetual

	for ; iterator.Valid(); iterator.Next() {
		var val types.Perpetual
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) SetPerpetual(ctx sdk.Context, p types.Perpetual) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualKey(p.MarketId, p.Id)
	b := k.cdc.MustMarshal(&p)
	store.Set(key, b)
}

func (k Keeper) DeletePerpetual(ctx sdk.Context, p types.Perpetual) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualKey(p.MarketId, p.Id)
	store.Delete(key)
}

func (k Keeper) GetPerpetualValue(ctx sdk.Context, perpetual types.Perpetual) math.LegacyDec {
	twapPrice := k.GetCurrentTwapPrice(ctx, perpetual.MarketId)
	if twapPrice.IsZero() {
		panic("twap price is zero while calculating perpetual value")
	}
	return twapPrice.Mul(perpetual.Quantity)
}

func (k Keeper) GetMaintenanceMargin(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) math.LegacyDec {
	currentValue := k.GetPerpetualValue(ctx, perpetual)
	return market.MaintenanceMarginRatio.Mul(currentValue)
}

// GetCurrentLeverage currentValue / balanceValue
func (k Keeper) GetCurrentLeverage(ctx sdk.Context, perpetual types.Perpetual) (math.LegacyDec, error) {
	currentValue := k.GetPerpetualValue(ctx, perpetual)
	subaccount, err := k.GetSubAccount(ctx, perpetual.GetOwnerAccAddress(), perpetual.MarketId)
	if err != nil {
		return math.LegacyDec{}, err
	}
	balanceValue, err := k.GetAvailableBalanceValue(ctx, subaccount)
	if err != nil {
		return math.LegacyDec{}, err
	}
	if balanceValue.IsZero() || balanceValue.IsNil() {
		return math.LegacyMaxSortableDec, nil
	}
	return currentValue.Quo(balanceValue), nil
}

// GetLiquidationPrice
// Long: Liquidation Price = Entry Price × (1 - 1/Current Leverage) / (1 - Maintenance Margin Rate)
// Short:Liquidation Price = Entry Price × (1 + 1/Current Leverage) / (1 + Maintenance Margin Rate)
func (k Keeper) GetLiquidationPrice(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) (math.LegacyDec, error) {
	leverage, err := k.GetCurrentLeverage(ctx, perpetual)
	if err != nil {
		return math.LegacyDec{}, err
	}
	numSub := math.LegacyOneDec().Quo(leverage)
	num := math.LegacyOneDec().Sub(numSub)
	den := math.LegacyOneDec().Sub(market.MaintenanceMarginRatio)
	if perpetual.IsShort() {
		numAdd := math.LegacyOneDec().Quo(leverage)
		num = math.LegacyOneDec().Add(numAdd)
		den = math.LegacyOneDec().Add(market.MaintenanceMarginRatio)
	}
	resultMult := num.Quo(den)
	return perpetual.EntryPrice.Mul(resultMult), nil
}
