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
	"github.com/cosmos/cosmos-sdk/types/address"
	pkeeper "github.com/elys-network/elys/x/parameter/keeper"
	ptypes "github.com/elys-network/elys/x/parameter/types"

	"github.com/elys-network/elys/x/perpetual/keeper"

	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/perpetual/types/mocks"
	"github.com/stretchr/testify/require"
)

func PerpetualKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, *mocks.AssetProfileKeeper) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	govAddress := sdk.AccAddress(address.Module("gov"))

	assetProfileKeeper := mocks.NewAssetProfileKeeper(t)

	storeKeyP := sdk.NewKVStoreKey(ptypes.StoreKey)
	memStoreKeyP := storetypes.NewMemoryStoreKey(ptypes.MemStoreKey)
	stateStore.MountStoreWithDB(storeKeyP, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKeyP, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	parameterKeeper := pkeeper.NewKeeper(cdc,
		storeKeyP,
		memStoreKeyP,
		govAddress.String(),
	)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		govAddress.String(),
		nil,
		nil,
		nil,
		assetProfileKeeper,
		nil,
		parameterKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	params := types.DefaultParams()
	k.SetParams(ctx, &params)

	paramsP := ptypes.DefaultParams()
	paramsP.TotalBlocksPerYear = 86400 * 365
	parameterKeeper.SetParams(ctx, paramsP)

	return k, ctx, assetProfileKeeper
}
