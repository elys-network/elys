package keeper

import (
	"cosmossdk.io/core/store"
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/transferhook/types"
)

type (
	Keeper struct {
		Cdc          codec.BinaryCodec
		storeService store.KVStoreService
		ammKeeper    ammkeeper.Keeper
	}
)

func NewKeeper(
	Cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	ammKeeper ammkeeper.Keeper,
) *Keeper {

	return &Keeper{
		Cdc:          Cdc,
		storeService: storeService,
		ammKeeper:    ammKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
