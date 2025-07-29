package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v7/x/epochs/types"
)

// Keeper of this module maintains collections of epochs and hooks.
type Keeper struct {
	cdc          codec.Codec
	storeService store.KVStoreService
	hooks        types.EpochHooks
}

// NewKeeper returns a new instance of epochs Keeper
func NewKeeper(cdc codec.Codec, storeService store.KVStoreService) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeService: storeService,
	}
}

// SetHooks set the epoch hooks
func (k *Keeper) SetHooks(eh types.EpochHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set epochs hooks twice")
	}

	k.hooks = eh

	return k
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
