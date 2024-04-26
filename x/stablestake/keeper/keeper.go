package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeKey           storetypes.StoreKey
		memKey             storetypes.StoreKey
		paramstore         paramtypes.Subspace
		bk                 types.BankKeeper
		commitmentKeeper   *commitmentkeeper.Keeper
		assetProfileKeeper types.AssetProfileKeeper
		hooks              types.StableStakeHooks
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bk types.BankKeeper,
	commitmentKeeper *commitmentkeeper.Keeper,
	assetProfileKeeper types.AssetProfileKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		memKey:             memKey,
		paramstore:         ps,
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
