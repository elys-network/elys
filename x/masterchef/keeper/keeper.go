package keeper

import (
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	estakingkeeper "github.com/elys-network/elys/x/estaking/keeper"

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
		estakingKeeper      *estakingkeeper.Keeper

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
	estakingKeeper *estakingkeeper.Keeper,
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
		estakingKeeper:      estakingKeeper,
		authority:           authority,
	}
}
