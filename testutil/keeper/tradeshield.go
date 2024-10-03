package keeper

import (
	"cosmossdk.io/store/metrics"
	"github.com/cosmos/cosmos-sdk/runtime"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cosmosdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/elys-network/elys/x/tradeshield/types/mocks"
	"github.com/stretchr/testify/require"
)

func TradeshieldKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, *mocks.AmmKeeper, *mocks.TierKeeper, *mocks.PerpetualKeeper) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)

	db := cosmosdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	govAddress := sdk.AccAddress(address.Module("gov"))

	ammKeeper := mocks.NewAmmKeeper(t)
	tierKeeper := mocks.NewTierKeeper(t)
	perpetualKeeper := mocks.NewPerpetualKeeper(t)

	k := keeper.NewKeeper(
		cdc,
		storeService,
		govAddress.String(),
		ammKeeper,
		tierKeeper,
		perpetualKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	params := types.DefaultParams()
	k.SetParams(ctx, &params)

	return k, ctx, ammKeeper, tierKeeper, perpetualKeeper
}
