package keeper

import (
	"cosmossdk.io/core/store"
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tier/types"
)

type (
	Keeper struct {
		cdc                codec.BinaryCodec
		storeService       store.KVStoreService
		bankKeeper         types.BankKeeper
		oracleKeeper       types.OracleKeeper
		assetProfileKeeper types.AssetProfileKeeper
		amm                types.AmmKeeper
		estaking           types.EstakingKeeper
		masterchef         types.MasterchefKeeper
		commitement        types.CommitmentKeeper
		perpetual          types.PerpetualKeeper
		stakingKeeper      types.StakingKeeper
		leveragelp         types.LeverageLpKeeper
		stablestakeKeeper  types.StablestakeKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	bankKeeper types.BankKeeper,
	oracleKeeper types.OracleKeeper,
	assetProfileKeeper types.AssetProfileKeeper,
	amm types.AmmKeeper,
	estaking types.EstakingKeeper,
	masterchef types.MasterchefKeeper,
	commitement types.CommitmentKeeper,
	stakingKeeper types.StakingKeeper,
	perpetual types.PerpetualKeeper,
	leveragelp types.LeverageLpKeeper,
	stablestakeKeeper types.StablestakeKeeper,
) *Keeper {

	return &Keeper{
		cdc:                cdc,
		storeService:       storeService,
		bankKeeper:         bankKeeper,
		oracleKeeper:       oracleKeeper,
		assetProfileKeeper: assetProfileKeeper,
		amm:                amm,
		estaking:           estaking,
		masterchef:         masterchef,
		commitement:        commitement,
		stakingKeeper:      stakingKeeper,
		perpetual:          perpetual,
		leveragelp:         leveragelp,
		stablestakeKeeper:  stablestakeKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
