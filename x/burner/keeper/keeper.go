package keeper

import (
	"cosmossdk.io/core/store"
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/burner/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService

		bankKeeper types.BankKeeper
		authority  string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,

	bankKeeper types.BankKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:          cdc,
		storeService: storeService,

		bankKeeper: bankKeeper,
		authority:  authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
