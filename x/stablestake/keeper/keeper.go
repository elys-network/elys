package keeper

import (
	"cosmossdk.io/core/store"
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeService       store.KVStoreService
		authority          string
		bk                 types.BankKeeper
		commitmentKeeper   *commitmentkeeper.Keeper
		assetProfileKeeper types.AssetProfileKeeper
		hooks              types.StableStakeHooks
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
	bk types.BankKeeper,
	commitmentKeeper *commitmentkeeper.Keeper,
	assetProfileKeeper types.AssetProfileKeeper,
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
	}
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
