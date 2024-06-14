package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/elys-network/elys/x/tier/types"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeKey           storetypes.StoreKey
		memKey             storetypes.StoreKey
		paramstore         paramtypes.Subspace
		bankKeeper         types.BankKeeper
		oracleKeeper       types.OracleKeeper
		assetProfileKeeper types.AssetProfileKeeper
		amm                types.AmmKeeper
		estaking           types.EstakingKeeper
		masterchef         types.MasterchefKeeper
		commitement        types.CommitmentKeeper
		perpetual          types.PerpetualKeeper
		stakingKeeper      types.StakingKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	bankKeeper types.BankKeeper,
	oracleKeeper types.OracleKeeper,
	assetProfileKeeper types.AssetProfileKeeper,
	amm types.AmmKeeper,
	estaking types.EstakingKeeper,
	masterchef types.MasterchefKeeper,
	commitement types.CommitmentKeeper,
	stakingKeeper types.StakingKeeper,
	perpetual types.PerpetualKeeper,
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
		bankKeeper:         bankKeeper,
		oracleKeeper:       oracleKeeper,
		assetProfileKeeper: assetProfileKeeper,
		amm:                amm,
		estaking:           estaking,
		masterchef:         masterchef,
		commitement:        commitement,
		stakingKeeper:      stakingKeeper,
		perpetual:          perpetual,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
