package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/elys-network/elys/x/amm/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		hooks      types.AmmHooks

		bankKeeper       types.BankKeeper
		accountKeeper    types.AccountKeeper
		oracleKeeper     types.OracleKeeper
		commitmentKeeper *commitmentkeeper.Keeper
		apKeeper         types.AssetProfileKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	oracleKeeper types.OracleKeeper,
	commitmentKeeper *commitmentkeeper.Keeper,
	apKeeper types.AssetProfileKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,

		bankKeeper:       bankKeeper,
		accountKeeper:    accountKeeper,
		oracleKeeper:     oracleKeeper,
		commitmentKeeper: commitmentKeeper,
		apKeeper:         apKeeper,
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
