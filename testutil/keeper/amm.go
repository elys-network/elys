package keeper

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/amm/types/mocks"
	"github.com/stretchr/testify/require"
)

func AmmKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, *mocks.AccountedPoolKeeper, *mocks.OracleKeeper) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	transientStoreKey := storetypes.NewTransientStoreKey(types.TStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	accountedPoolKeeper := mocks.NewAccountedPoolKeeper(t)
	oracleKeeper := mocks.NewOracleKeeper(t)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		transientStoreKey,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		nil,
		nil,
		nil,
		oracleKeeper,
		nil,
		nil,
		accountedPoolKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx, accountedPoolKeeper, oracleKeeper
}
