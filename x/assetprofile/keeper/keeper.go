package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	"github.com/elys-network/elys/v7/x/assetprofile/types"
)

type (
	Keeper struct {
		cdc            codec.BinaryCodec
		storeService   store.KVStoreService
		transferKeeper *ibctransferkeeper.Keeper
		// the address capable of executing a Msg* messages. Typically, this
		// should be the x/gov module account.
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	transferKeeper *ibctransferkeeper.Keeper,
	authority string,
) *Keeper {

	return &Keeper{
		cdc:            cdc,
		storeService:   storeService,
		transferKeeper: transferKeeper,
		authority:      authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
