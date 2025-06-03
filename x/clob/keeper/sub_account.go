package keeper

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetSubAccount(ctx sdk.Context, owner sdk.AccAddress, subAccountId uint64) (types.SubAccount, error) {
	key := types.GetSubAccountKey(owner, subAccountId)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		return types.SubAccount{}, errors.Wrapf(types.ErrSubAccountNotFound, "owner: %s, subAccountId: %d", owner.String(), subAccountId)
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
	key := types.GetSubAccountKey(s.GetOwnerAccAddress(), s.Id)
	b := k.cdc.MustMarshal(&s)
	store.Set(key, b)
}

func (k Keeper) SendFromSubAccount(ctx sdk.Context, subAccount types.SubAccount, to sdk.AccAddress, coins sdk.Coins) error {
	return k.bankKeeper.SendCoins(ctx, subAccount.GetTradingAccountAddress(), to, coins)
}

func (k Keeper) AddToSubAccount(ctx sdk.Context, from sdk.AccAddress, subAccount types.SubAccount, coins sdk.Coins) error {
	return k.bankKeeper.SendCoins(ctx, from, subAccount.GetTradingAccountAddress(), coins)
}

func (k Keeper) SendFromSubAccountToSubAccount(ctx sdk.Context, from types.SubAccount, to types.SubAccount, coins sdk.Coins) error {
	return k.bankKeeper.SendCoins(ctx, from.GetTradingAccountAddress(), to.GetTradingAccountAddress(), coins)
}

func (k Keeper) GetSubAccountBalance(ctx sdk.Context, subAccount types.SubAccount) sdk.Coins {
	return k.bankKeeper.GetAllBalances(ctx, subAccount.GetTradingAccountAddress())
}

func (k Keeper) GetSubAccountBalanceOf(ctx sdk.Context, subAccount types.SubAccount, denom string) sdk.Coin {
	return k.bankKeeper.GetBalance(ctx, subAccount.GetTradingAccountAddress(), denom)
}

//func (k Keeper) WithdrawableBalance(ctx sdk.Context, subAccount types.SubAccount) error {
//	if subAccount.IsIsolated() {
//		// No need to check for current positions as the margin amount has already been transferred to the market account
//		// Need to check for the open orders and the trading fees and margin amount for that
//		k.GetPerpetualOrder()
//	}
//}
