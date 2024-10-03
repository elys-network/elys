package keeper

import (
	"cosmossdk.io/core/store"
	"fmt"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

type (
	Keeper struct {
		cdc                 codec.BinaryCodec
		storeService        store.KVStoreService
		parameterKeeper     types.ParameterKeeper
		commitmentKeeper    types.CommitmentKeeper
		amm                 types.AmmKeeper
		oracleKeeper        types.OracleKeeper
		assetProfileKeeper  types.AssetProfileKeeper
		accountedPoolKeeper types.AccountedPoolKeeper
		stableKeeper        types.StableStakeKeeper
		tokenomicsKeeper    types.TokenomicsKeeper
		authKeeper          types.AccountKeeper
		bankKeeper          types.BankKeeper
		perpetualKeeper     types.PeperpetualKeeper

		authority string // gov module addresss
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	parameterKeeper types.ParameterKeeper,
	ck types.CommitmentKeeper,
	amm types.AmmKeeper,
	ok types.OracleKeeper,
	ap types.AssetProfileKeeper,
	accountedPoolKeeper types.AccountedPoolKeeper,
	stableKeeper types.StableStakeKeeper,
	tokenomicsKeeper types.TokenomicsKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	perpetualKeeper types.PeperpetualKeeper,
	authority string,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		storeService:        storeService,
		parameterKeeper:     parameterKeeper,
		commitmentKeeper:    ck,
		amm:                 amm,
		oracleKeeper:        ok,
		assetProfileKeeper:  ap,
		accountedPoolKeeper: accountedPoolKeeper,
		stableKeeper:        stableKeeper,
		tokenomicsKeeper:    tokenomicsKeeper,
		authKeeper:          ak,
		bankKeeper:          bk,
		perpetualKeeper:     perpetualKeeper,
		authority:           authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
