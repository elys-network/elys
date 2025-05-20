package keeper

import (
	"fmt"

	"cosmossdk.io/core/store"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmentkeeper "github.com/elys-network/elys/v4/x/commitment/keeper"
	"github.com/elys-network/elys/v4/x/stablestake/types"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeService       store.KVStoreService
		authority          string
		bk                 types.BankKeeper
		commitmentKeeper   *commitmentkeeper.Keeper
		assetProfileKeeper types.AssetProfileKeeper
		oracleKeeper       types.OracleKeeper
		ammKeeper          types.AmmKeeper
		hooks              types.StableStakeHooks
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
	bk types.BankKeeper,
	commitmentKeeper *commitmentkeeper.Keeper,
	assetProfileKeeper types.AssetProfileKeeper,
	oracleKeeper types.OracleKeeper,
	ammKeeper types.AmmKeeper,
) *Keeper {

	// ensure that authority is a valid AccAddress
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	return &Keeper{
		cdc:                cdc,
		storeService:       storeService,
		authority:          authority,
		bk:                 bk,
		commitmentKeeper:   commitmentKeeper,
		assetProfileKeeper: assetProfileKeeper,
		oracleKeeper:       oracleKeeper,
		ammKeeper:          ammKeeper,
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
