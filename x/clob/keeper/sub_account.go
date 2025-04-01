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

func (k Keeper) GetSubAccount(ctx sdk.Context, owner sdk.AccAddress, marketId uint64) (types.SubAccount, error) {
	key := types.GetSubAccountKey(owner, marketId)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		return types.SubAccount{}, types.ErrSubAccountNotFound
	}

	var val types.SubAccount
	k.cdc.MustUnmarshal(b, &val)
	return val, nil
}

func (k Keeper) GetAllOwnerSubAccount(ctx sdk.Context, addr sdk.AccAddress) []types.SubAccount {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.GetAddressSubAccountPrefixKey(addr))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.SubAccount

	for ; iterator.Valid(); iterator.Next() {
		var val types.SubAccount
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) GetAllSubAccount(ctx sdk.Context) []types.SubAccount {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.SubAccountPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.SubAccount

	for ; iterator.Valid(); iterator.Next() {
		var val types.SubAccount
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) SetSubAccount(ctx sdk.Context, s types.SubAccount) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetSubAccountKey(s.GetOwnerAccAddress(), s.MarketId)
	b := k.cdc.MustMarshal(&s)
	store.Set(key, b)
}

func (k Keeper) SendFromSubAccount(ctx sdk.Context, subAccount types.SubAccount, to sdk.AccAddress, coins sdk.Coins) error {
	subAccount.AvailableBalance = subAccount.AvailableBalance.Sub(coins...)
	k.SetSubAccount(ctx, subAccount)
	return k.bankKeeper.SendCoins(ctx, subAccount.GetTradingAccountAddress(), to, coins)
}

func (k Keeper) AddToSubAccount(ctx sdk.Context, from sdk.AccAddress, subAccount types.SubAccount, coins sdk.Coins) error {
	subAccount.AvailableBalance = subAccount.AvailableBalance.Add(coins...)
	k.SetSubAccount(ctx, subAccount)
	return k.bankKeeper.SendCoins(ctx, from, subAccount.GetTradingAccountAddress(), coins)
}

func (k Keeper) GetAvailableBalanceValue(ctx sdk.Context, subAccount types.SubAccount) (totalValue math.Dec, err error) {
	for _, coin := range subAccount.AvailableBalance {
		price, err := k.GetDenomPrice(ctx, coin.Denom)
		if err != nil {
			return totalValue, err
		}
		v, addErr := price.Mul(utils.IntToDec(coin.Amount))
		if addErr != nil {
			return totalValue, addErr
		}
		totalValue, err = totalValue.Add(v)
		if err != nil {
			return totalValue, addErr
		}
	}
	return totalValue, nil
}
