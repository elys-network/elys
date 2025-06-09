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

	k.DecrementTotalOpenPosition(ctx, p.MarketId)
}

func (k Keeper) GetMaintenanceMargin(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) (math.Int, error) {
	currentValue, err := k.GetPerpetualOracleValue(ctx, market, perpetual)
	if err != nil {
		return math.Int{}, err
	}
	quoteDenomPrice, err := k.GetDenomPrice(ctx, market.QuoteDenom)
	if err != nil {
		return math.Int{}, err
	}
	return market.MaintenanceMarginRatio.Mul(currentValue).Quo(quoteDenomPrice).RoundInt(), nil
}

func (k Keeper) GetPerpetualOracleValue(ctx sdk.Context, market types.PerpetualMarket, perpetual types.Perpetual) (math.LegacyDec, error) {
	currentPrice, err := k.GetAssetPriceFromDenom(ctx, market.BaseDenom)
	if err != nil {
		return math.LegacyDec{}, err
	}
	return currentPrice.Mul(perpetual.Quantity).Abs(), nil
}

// GetEquityValue = InitialMarginValue + UPnL
func (k Keeper) GetEquityValue(ctx sdk.Context, perpetual types.Perpetual, subAccount types.SubAccount, market types.PerpetualMarket) (math.LegacyDec, error) {
	if subAccount.IsIsolated() {
		// InitialMarginPosted + UnrealizedPNL
		baseAssetCurrentPrice, err := k.GetAssetPriceFromDenom(ctx, market.BaseDenom)
		if err != nil {
			return math.LegacyZeroDec(), err
		}
		unrealizedPnLValue, err := perpetual.CalculateUnrealizedPnLValue(baseAssetCurrentPrice)
		if err != nil {
			return math.LegacyDec{}, err
		}
		quoteAssetDenomPrice, err := k.GetDenomPrice(ctx, market.QuoteDenom)
		if err != nil {
			return math.LegacyDec{}, err
		}
		initialMarginValue := perpetual.MarginAmount.ToLegacyDec().Mul(quoteAssetDenomPrice)
		return initialMarginValue.Add(unrealizedPnLValue), nil
	} else {
		// TODO TotalAccountValue
		panic("implement me")
	}
}

// GetEffectiveLeverage PositionValue / EquityValue
func (k Keeper) GetEffectiveLeverage(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) (math.LegacyDec, error) {
	subAccount, err := k.GetSubAccount(ctx, perpetual.GetOwnerAccAddress(), perpetual.SubAccountId)
	if err != nil {
		return math.LegacyDec{}, err
	}
	currentValue, err := k.GetPerpetualOracleValue(ctx, market, perpetual)
	if err != nil {
		return math.LegacyDec{}, err
	}
	equityValue, err := k.GetEquityValue(ctx, perpetual, subAccount, market)
	if err != nil {
		return math.LegacyDec{}, err
	}
	return currentValue.Quo(equityValue), nil
}
