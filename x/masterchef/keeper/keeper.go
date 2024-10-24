package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	estakingkeeper "github.com/elys-network/elys/x/estaking/keeper"

	"github.com/elys-network/elys/x/masterchef/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace

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
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
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
		storeKey:            storeKey,
		memKey:              memKey,
		paramstore:          ps,
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

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
