package keeper

import (
	"fmt"

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
