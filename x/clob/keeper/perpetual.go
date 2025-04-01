package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
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

func (k Keeper) GetPerpetualValue(ctx sdk.Context, perpetual types.Perpetual) (math.Dec, error) {
	midPrice, err := k.GetMidPrice(ctx, perpetual.MarketId)
	if err != nil {
		return math.Dec{}, err
	}
	return midPrice.Mul(utils.IntToDec(perpetual.Quantity))
}

func (k Keeper) GetMaintenanceMargin(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) (math.Dec, error) {
	currentValue, err := k.GetPerpetualValue(ctx, perpetual)
	if err != nil {
		return math.Dec{}, err
	}
	return market.MaintenanceMarginRatio.Mul(currentValue)
}

// GetCurrentLeverage currentValue / balanceValue
func (k Keeper) GetCurrentLeverage(ctx sdk.Context, perpetual types.Perpetual) (math.Dec, error) {
	currentValue, err := k.GetPerpetualValue(ctx, perpetual)
	if err != nil {
		return math.Dec{}, err
	}
	subaccount, err := k.GetSubAccount(ctx, perpetual.GetOwnerAccAddress(), perpetual.MarketId)
	if err != nil {
		return math.Dec{}, err
	}
	balanceValue, err := k.GetAvailableBalanceValue(ctx, subaccount)
	if err != nil {
		return math.Dec{}, err
	}
	return currentValue.Quo(balanceValue)
}

// GetLiquidationPrice
// Long: Liquidation Price = Entry Price × (1 - 1/Leverage) / (1 - Maintenance Margin Rate)
// Short:Liquidation Price = Entry Price × (1 + 1/Leverage) / (1 + Maintenance Margin Rate)
func (k Keeper) GetLiquidationPrice(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket) (math.Dec, error) {
	leverage, err := k.GetCurrentLeverage(ctx, perpetual)
	if err != nil {
		return math.Dec{}, err
	}
	num_sub, err := utils.OneDec.Quo(leverage)
	if err != nil {
		return math.Dec{}, err
	}
	num, err := utils.OneDec.Sub(num_sub)
	if err != nil {
		return math.Dec{}, err
	}
	den, err := utils.OneDec.Sub(market.MaintenanceMarginRatio)
	if err != nil {
		return math.Dec{}, err
	}
	if perpetual.IsShort() {
		num_add, err := utils.OneDec.Quo(leverage)
		if err != nil {
			return math.Dec{}, err
		}
		num, err = utils.OneDec.Add(num_add)
		if err != nil {
			return math.Dec{}, err
		}
		den, err = utils.OneDec.Add(market.MaintenanceMarginRatio)
		if err != nil {
			return math.Dec{}, err
		}
	}
	result_mult, err := num.Quo(den)
	if err != nil {
		return math.Dec{}, err
	}
	return perpetual.EntryPrice.Mul(result_mult)
}
