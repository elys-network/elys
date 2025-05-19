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

func (k Keeper) GetMaintenanceMargin(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) math.LegacyDec {
	currentValue := k.GetPerpetualAbsValue(ctx, perpetual)
	return market.MaintenanceMarginRatio.Mul(currentValue)
}

func (k Keeper) GetPerpetualAbsValue(ctx sdk.Context, perpetual types.Perpetual) math.LegacyDec {
	twapPrice := k.GetCurrentTwapPrice(ctx, perpetual.MarketId)
	if twapPrice.IsZero() {
		panic("twap price is zero while calculating perpetual value")
	}
	return twapPrice.Mul(perpetual.Quantity).Abs()
}

// GetEquityValue = InitialMarginValue + UPnL
func (k Keeper) GetEquityValue(ctx sdk.Context, perpetual types.Perpetual, subAccount types.SubAccount, market types.PerpetualMarket) (math.LegacyDec, error) {
	if subAccount.IsIsolated() {
		// InitialMarginPosted + UnrealizedPNL
		markPrice := k.GetCurrentTwapPrice(ctx, perpetual.MarketId)
		unrealizedPnLValue, err := perpetual.CalculateUnrealizedPnLValue(markPrice)
		if err != nil {
			return math.LegacyDec{}, err
		}
		price, err := k.GetDenomPrice(ctx, market.QuoteDenom)
		if err != nil {
			return math.LegacyDec{}, err
		}
		initialMarginValue := perpetual.MarginAmount.ToLegacyDec().Mul(price)
		return initialMarginValue.Add(unrealizedPnLValue), nil
	} else {
		// TODO TotalAccountValue
		panic("implement me")
	}
}

// GetEffectiveLeverage PositionValue / EquityValue
func (k Keeper) GetEffectiveLeverage(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) (math.LegacyDec, error) {
	subAccount, err := k.GetSubAccount(ctx, perpetual.GetOwnerAccAddress(), perpetual.MarketId)
	if err != nil {
		return math.LegacyDec{}, err
	}
	currentValue := k.GetPerpetualAbsValue(ctx, perpetual)
	equityValue, err := k.GetEquityValue(ctx, perpetual, subAccount, market)
	if err != nil {
		return math.LegacyDec{}, err
	}
	return currentValue.Quo(equityValue), nil
}
