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
	"github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/elys-network/elys/x/tradeshield/types/mocks"
	"github.com/stretchr/testify/require"
)

func TradeshieldKeeper(t testing.TB) (*keeper.Keeper, sdk.Context, *mocks.AmmKeeper, *mocks.TierKeeper, *mocks.PerpetualKeeper) {
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

	ammKeeper := mocks.NewAmmKeeper(t)
	tierKeeper := mocks.NewTierKeeper(t)
	perpetualKeeper := mocks.NewPerpetualKeeper(t)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
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
