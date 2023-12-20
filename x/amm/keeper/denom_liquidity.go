package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// SetDenomLiquidity sets a specific denomLiquidity in the store from its index
func (k Keeper) SetDenomLiquidity(ctx sdk.Context, denomLiquidity types.DenomLiquidity) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomLiquidityKeyPrefix))
	b := k.cdc.MustMarshal(&denomLiquidity)
	store.Set(types.DenomLiquidityKey(denomLiquidity.Denom), b)
}

// GetDenomLiquidity returns a denomLiquidity from its index
func (k Keeper) GetDenomLiquidity(ctx sdk.Context, denom string) (val types.DenomLiquidity, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomLiquidityKeyPrefix))
	b := store.Get(types.DenomLiquidityKey(denom))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDenomLiquidity removes a denomLiquidity from the store
func (k Keeper) RemoveDenomLiquidity(ctx sdk.Context, denom string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomLiquidityKeyPrefix))
	store.Delete(types.DenomLiquidityKey(denom))
}

// GetAllDenomLiquidity returns all denomLiquidity
func (k Keeper) GetAllDenomLiquidity(ctx sdk.Context) (list []types.DenomLiquidity) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomLiquidityKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.DenomLiquidity
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

// IncreaseDenomLiquidity increases the liquidity of a denom by a certain amount
func (k Keeper) IncreaseDenomLiquidity(ctx sdk.Context, denom string, amount math.Int) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomLiquidityKeyPrefix))
	b := store.Get(types.DenomLiquidityKey(denom))
	if b == nil {
		return fmt.Errorf("denom not found")
	}
	var denomLiquidity types.DenomLiquidity
	k.cdc.MustUnmarshal(b, &denomLiquidity)
	denomLiquidity.Liquidity = denomLiquidity.Liquidity.Add(amount)
	newB := k.cdc.MustMarshal(&denomLiquidity)
	store.Set(types.DenomLiquidityKey(denom), newB)
	return nil
}

// DecreaseDenomLiquidity decreases the liquidity of a denom by a certain amount
func (k Keeper) DecreaseDenomLiquidity(ctx sdk.Context, denom string, amount math.Int) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DenomLiquidityKeyPrefix))
	b := store.Get(types.DenomLiquidityKey(denom))
	if b == nil {
		return fmt.Errorf("denom not found")
	}
	var denomLiquidity types.DenomLiquidity
	k.cdc.MustUnmarshal(b, &denomLiquidity)
	if denomLiquidity.Liquidity.LT(amount) {
		return fmt.Errorf("not enough liquidity")
	}
	denomLiquidity.Liquidity = denomLiquidity.Liquidity.Sub(amount)
	newB := k.cdc.MustMarshal(&denomLiquidity)
	store.Set(types.DenomLiquidityKey(denom), newB)
	return nil
}

func (k Keeper) RecordTotalLiquidityIncrease(ctx sdk.Context, coins sdk.Coins) {
	for _, coin := range coins {
		_, found := k.GetDenomLiquidity(ctx, coin.Denom)
		if !found {
			k.SetDenomLiquidity(ctx, types.DenomLiquidity{Denom: coin.Denom, Liquidity: sdk.ZeroInt()})
		}

		k.IncreaseDenomLiquidity(ctx, coin.Denom, coin.Amount)
	}
}

func (k Keeper) RecordTotalLiquidityDecrease(ctx sdk.Context, coins sdk.Coins) {
	for _, coin := range coins {
		_, found := k.GetDenomLiquidity(ctx, coin.Denom)
		if !found {
			k.SetDenomLiquidity(ctx, types.DenomLiquidity{Denom: coin.Denom, Liquidity: sdk.ZeroInt()})
		}
		err := k.DecreaseDenomLiquidity(ctx, coin.Denom, coin.Amount)
		if err != nil {
			k.Logger(ctx).Error(err.Error())
			panic(err)
		}
	}
}
