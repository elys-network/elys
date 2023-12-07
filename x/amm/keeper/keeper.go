package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/elys-network/elys/x/amm/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	pkeeper "github.com/elys-network/elys/x/parameter/keeper"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		storeKey          storetypes.StoreKey
		transientStoreKey storetypes.StoreKey
		paramstore        paramtypes.Subspace
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
	storeKey,
	transientStoreKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	authority string,

	parameterKeeper *pkeeper.Keeper,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	oracleKeeper types.OracleKeeper,
	commitmentKeeper *commitmentkeeper.Keeper,
	assetProfileKeeper types.AssetProfileKeeper,
	accountedPoolKeeper types.AccountedPoolKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		transientStoreKey: transientStoreKey,
		paramstore:        ps,
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
