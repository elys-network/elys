package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/core/store"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v6/x/amm/types"
	commitmentkeeper "github.com/elys-network/elys/v6/x/commitment/keeper"
	pkeeper "github.com/elys-network/elys/v6/x/parameter/keeper"
	tierkeeper "github.com/elys-network/elys/v6/x/tier/keeper"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeService      store.KVStoreService
		transientStoreKey storetypes.StoreKey
		authority         string
		hooks             types.AmmHooks

		parameterKeeper     *pkeeper.Keeper
		bankKeeper          types.BankKeeper
		accountKeeper       types.AccountKeeper
		oracleKeeper        types.OracleKeeper
		commitmentKeeper    *commitmentkeeper.Keeper
		assetProfileKeeper  types.AssetProfileKeeper
		accountedPoolKeeper types.AccountedPoolKeeper
		tierKeeper          *tierkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	transientStoreKey storetypes.StoreKey,
	authority string,

	parameterKeeper *pkeeper.Keeper,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	oracleKeeper types.OracleKeeper,
	commitmentKeeper *commitmentkeeper.Keeper,
	assetProfileKeeper types.AssetProfileKeeper,
	accountedPoolKeeper types.AccountedPoolKeeper,
	tierKeeper *tierkeeper.Keeper,
) *Keeper {

	return &Keeper{
		cdc:               cdc,
		storeService:      storeService,
		transientStoreKey: transientStoreKey,
		authority:         authority,

		parameterKeeper:     parameterKeeper,
		bankKeeper:          bankKeeper,
		accountKeeper:       accountKeeper,
		oracleKeeper:        oracleKeeper,
		commitmentKeeper:    commitmentKeeper,
		assetProfileKeeper:  assetProfileKeeper,
		accountedPoolKeeper: accountedPoolKeeper,
		tierKeeper:          tierKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Set the amm hooks.
func (k *Keeper) SetHooks(gh types.AmmHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set amm hooks twice")
	}

	k.hooks = gh

	return k
}

func (k *Keeper) SetTierKeeper(tk *tierkeeper.Keeper) {
	k.tierKeeper = tk
}

func (k *Keeper) GetTierKeeper() *tierkeeper.Keeper {
	return k.tierKeeper
}

func (k *Keeper) GetCommitmentKeeper() *commitmentkeeper.Keeper {
	return k.commitmentKeeper
}

func (k *Keeper) GetLastProccessed(ctx sdk.Context, id uint64) sdk.AccAddress {
	kvStore := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetLastProcessedKey(id)
	bz := kvStore.Get(key)
	if bz != nil {
		return sdk.AccAddress(bz)
	} else {
		return nil
	}
}

func (k *Keeper) SetLastProccessed(ctx sdk.Context, id uint64, addr sdk.AccAddress) {
	kvStore := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetLastProcessedKey(id)
	kvStore.Set(key, addr)
}

func (k *Keeper) DeleteLastProccessed(ctx sdk.Context, id uint64) {
	kvStore := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetLastProcessedKey(id)
	kvStore.Delete(key)
}

func (k Keeper) SetAddressInMigrationQueue(ctx sdk.Context, addr sdk.AccAddress) {
	kvStore := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := append(types.MigrationQueuePrefix, addr.Bytes()...)
	kvStore.Set(key, addr)
}

// RemoveAddressFromMigrationQueue deletes the account from the queue once processed.
func (k Keeper) RemoveAddressFromMigrationQueue(ctx sdk.Context, addr sdk.AccAddress) {
	kvStore := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := append(types.MigrationQueuePrefix, addr.Bytes()...)
	kvStore.Delete(key)
}

// GetMigrationQueueIterator returns an iterator over all queued addresses.
func (k Keeper) GetMigrationQueueIterator(ctx sdk.Context) storetypes.Iterator {
	kvStore := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return storetypes.KVStorePrefixIterator(kvStore, types.MigrationQueuePrefix)
}

func (k Keeper) SetMigrationReceipt(ctx sdk.Context, receipt types.BalanceMigrationReceipt) {
	kvStore := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	addr, _ := sdk.AccAddressFromBech32(receipt.Address)

	key := append(types.MigrationHistoryPrefix, addr.Bytes()...)
	value := k.cdc.MustMarshal(&receipt)

	kvStore.Set(key, value)
}

func (k Keeper) GetMigrationReceipt(ctx sdk.Context, addr sdk.AccAddress) types.BalanceMigrationReceipt {
	kvStore := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	key := append(types.MigrationHistoryPrefix, addr.Bytes()...)
	bz := kvStore.Get(key)

	var receipt types.BalanceMigrationReceipt
	if bz != nil {
		k.cdc.MustUnmarshal(bz, &receipt)
	}
	return receipt
}
