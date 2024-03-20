package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"

	epochkeeper "github.com/elys-network/elys/x/epochs/keeper"
	"github.com/elys-network/elys/x/launchpad/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		authority  string

		bankKeeper         types.BankKeeper
		oracleKeeper       types.OracleKeeper
		epochKeeper        *epochkeeper.Keeper
		commitmentKeeper   *commitmentkeeper.Keeper
		assetProfileKeeper types.AssetProfileKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	authority string,
	bankKeeper types.BankKeeper,
	oracleKeeper types.OracleKeeper,
	epochKeeper *epochkeeper.Keeper,
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
		authority:          authority,
		bankKeeper:         bankKeeper,
		oracleKeeper:       oracleKeeper,
		epochKeeper:        epochKeeper,
		commitmentKeeper:   commitmentKeeper,
		assetProfileKeeper: assetProfileKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
