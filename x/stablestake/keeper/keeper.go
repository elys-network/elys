package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/core/store"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmentkeeper "github.com/elys-network/elys/v6/x/commitment/keeper"
	"github.com/elys-network/elys/v6/x/stablestake/types"
)

type Keeper struct {
	cdc                codec.BinaryCodec
	storeService       store.KVStoreService
	authority          string
	bk                 types.BankKeeper
	commitmentKeeper   *commitmentkeeper.Keeper
	assetProfileKeeper types.AssetProfileKeeper
	oracleKeeper       types.OracleKeeper
	ammKeeper          types.AmmKeeper
	leverageLpKeeper   types.LeverageLpKeeper
	hooks              types.StableStakeHooks
}

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
	bk types.BankKeeper,
	commitmentKeeper *commitmentkeeper.Keeper,
	assetProfileKeeper types.AssetProfileKeeper,
	oracleKeeper types.OracleKeeper,
	ammKeeper types.AmmKeeper,
) *Keeper {

	// ensure that authority is a valid AccAddress
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	return &Keeper{
		cdc:                cdc,
		storeService:       storeService,
		authority:          authority,
		bk:                 bk,
		commitmentKeeper:   commitmentKeeper,
		assetProfileKeeper: assetProfileKeeper,
		oracleKeeper:       oracleKeeper,
		ammKeeper:          ammKeeper,
	}
}

func (k *Keeper) SetLeverageLpKeeper(v types.LeverageLpKeeper) {
	k.leverageLpKeeper = v
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetHooks set the epoch hooks
func (k *Keeper) SetHooks(eh types.StableStakeHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set stablestake hooks twice")
	}

	k.hooks = eh

	return k
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
