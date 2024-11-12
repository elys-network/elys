package keeper

import (
	"cosmossdk.io/core/store"
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/amm/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	pkeeper "github.com/elys-network/elys/x/parameter/keeper"
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
