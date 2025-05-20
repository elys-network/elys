package keeper

import (
	"testing"

	"cosmossdk.io/store/metrics"
	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cosmosdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v4/x/amm/keeper"
	"github.com/elys-network/elys/v4/x/amm/types"
	"github.com/elys-network/elys/v4/x/amm/types/mocks"
	"github.com/stretchr/testify/require"
)

func AmmKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, *mocks.AccountedPoolKeeper, *mocks.OracleKeeper) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	transientStoreKey := storetypes.NewTransientStoreKey(types.TStoreKey)

	db := cosmosdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	accountedPoolKeeper := mocks.NewAccountedPoolKeeper(t)
	oracleKeeper := mocks.NewOracleKeeper(t)

	k := keeper.NewKeeper(
		cdc,
		storeService,
		transientStoreKey,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		nil,
		nil,
		nil,
		oracleKeeper,
		nil,
		nil,
		accountedPoolKeeper,
		nil,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx, accountedPoolKeeper, oracleKeeper
}
