package keeper

import (
	"cosmossdk.io/core/store"
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/tokenomics/types"
)

type (
	Keeper struct {
		cdc              codec.BinaryCodec
		storeService     store.KVStoreService
		commitmentKeeper *commitmentkeeper.Keeper
		// the address capable of executing a Msg* messages. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	commitmentKeeper *commitmentkeeper.Keeper,
	authority string,
) *Keeper {

	return &Keeper{
		cdc:              cdc,
		storeService:     storeService,
		commitmentKeeper: commitmentKeeper,
		authority:        authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
