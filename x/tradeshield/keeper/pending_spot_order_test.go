package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createNPendingSpotOrder(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.SpotOrder {
	items := make([]types.SpotOrder, n)
	for i := range items {
		items[i].OrderId = keeper.AppendPendingSpotOrder(ctx, items[i])
	}
	return items
}

func TestPendingSpotOrderGet(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetPendingSpotOrder(ctx, item.OrderId)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestPendingSpotOrderRemove(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePendingSpotOrder(ctx, item.OrderId)
		_, found := keeper.GetPendingSpotOrder(ctx, item.OrderId)
		require.False(t, found)
	}
}

func TestPendingSpotOrderGetAll(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPendingSpotOrder(ctx)),
	)
}

func TestPendingSpotOrderCount(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingSpotOrder(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetPendingSpotOrderCount(ctx)-1)
}

func TestSortedSpotOrder(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)

	// Set to main storage
	keeper.AppendPendingSpotOrder(ctx, types.SpotOrder{
		OwnerAddress: "address",
		OrderId:      0,
		OrderType:    types.SpotOrderType_LIMITBUY,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdk.NewDec(1),
		},
	})

	order, _ := keeper.GetPendingSpotOrder(ctx, 1)

	err := keeper.InsertSpotSortedOrder(ctx, order)
	require.NoError(t, err)

	res, _ := keeper.GetAllSortedSpotOrder(ctx)

	assert.Equal(t, res, [][]uint64{{1}})

	// Insert two more elements
	// Set to main storage
	keeper.AppendPendingSpotOrder(ctx, types.SpotOrder{
		OwnerAddress: "address1",
		OrderId:      0,
		OrderType:    types.SpotOrderType_LIMITBUY,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdk.NewDec(20),
		},
	})

	keeper.AppendPendingSpotOrder(ctx, types.SpotOrder{
		OwnerAddress: "address2",
		OrderId:      0,
		OrderType:    types.SpotOrderType_LIMITBUY,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdk.NewDec(5),
		},
	})

	order2, _ := keeper.GetPendingSpotOrder(ctx, 2)
	order3, _ := keeper.GetPendingSpotOrder(ctx, 3)

	err = keeper.InsertSpotSortedOrder(ctx, order2)
	require.NoError(t, err)
	err = keeper.InsertSpotSortedOrder(ctx, order3)
	require.NoError(t, err)

	res, _ = keeper.GetAllSortedSpotOrder(ctx)

	// Should store in sorted order
	assert.Equal(t, res, [][]uint64{{1, 3, 2}})

	// Test binary search, search with rate 5
	index, err := keeper.SpotBinarySearch(ctx, sdk.NewDec(5), []uint64{1, 3, 2})
	require.NoError(t, err)

	// second element
	assert.Equal(t, index, 1)

	// Test remove sorted order
	keeper.RemoveSpotSortedOrder(ctx, 2)
	res, _ = keeper.GetAllSortedSpotOrder(ctx)

	// Should store in sorted order
	assert.Equal(t, res, [][]uint64{{1, 3}})
}

// TestExecuteStopLossOrder
func TestExecuteStopLossOrder(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)

	// Set to main storage
	keeper.AppendPendingSpotOrder(ctx, types.SpotOrder{
		OwnerAddress: "address",
		OrderId:      0,
		OrderType:    types.SpotOrderType_STOPLOSS,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdk.NewDec(1),
		},
	})

	order, _ := keeper.GetPendingSpotOrder(ctx, 1)

	err := keeper.InsertSpotSortedOrder(ctx, order)
	require.NoError(t, err)

	err = keeper.ExecuteStopLossOrder(ctx, order)
	require.NoError(t, err)

	// Should remove from sorted order
	res, _ := keeper.GetAllSortedSpotOrder(ctx)
	assert.Equal(t, res, [][]uint64{})

	// Should remove from main storage
	_, found := keeper.GetPendingSpotOrder(ctx, 1)
	assert.False(t, found)
}

// TestExecuteLimitSellOrder
func TestExecuteLimitSellOrder(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)

	// Set to main storage
	keeper.AppendPendingSpotOrder(ctx, types.SpotOrder{
		OwnerAddress: "address",
		OrderId:      0,
		OrderType:    types.SpotOrderType_LIMITSELL,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdk.NewDec(1),
		},
	})

	order, _ := keeper.GetPendingSpotOrder(ctx, 1)

	err := keeper.InsertSpotSortedOrder(ctx, order)
	require.NoError(t, err)

	err = keeper.ExecuteLimitSellOrder(ctx, order)
	require.NoError(t, err)

	// Should remove from sorted order
	res, _ := keeper.GetAllSortedSpotOrder(ctx)
	assert.Equal(t, res, [][]uint64{})

	// Should remove from main storage
	_, found := keeper.GetPendingSpotOrder(ctx, 1)
	assert.False(t, found)
}

// TestExecuteLimitBuyOrder
func TestExecuteLimitBuyOrder(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)

	// Set to main storage
	keeper.AppendPendingSpotOrder(ctx, types.SpotOrder{
		OwnerAddress: "address",
		OrderId:      0,
		OrderType:    types.SpotOrderType_LIMITBUY,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdk.NewDec(1),
		},
	})

	order, _ := keeper.GetPendingSpotOrder(ctx, 1)

	err := keeper.InsertSpotSortedOrder(ctx, order)
	require.NoError(t, err)

	err = keeper.ExecuteLimitBuyOrder(ctx, order)
	require.NoError(t, err)

	// Should remove from sorted order
	res, _ := keeper.GetAllSortedSpotOrder(ctx)
	assert.Equal(t, res, [][]uint64{})

	// Should remove from main storage
	_, found := keeper.GetPendingSpotOrder(ctx, 1)
	assert.False(t, found)
}

// TestExecuteMarketBuyOrder
func TestExecuteMarketBuyOrder(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)

	// Set to main storage
	keeper.AppendPendingSpotOrder(ctx, types.SpotOrder{
		OwnerAddress: "address",
		OrderId:      0,
		OrderType:    types.SpotOrderType_MARKETBUY,
		OrderPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdk.NewDec(1),
		},
	})

	order, _ := keeper.GetPendingSpotOrder(ctx, 1)

	err := keeper.InsertSpotSortedOrder(ctx, order)
	require.NoError(t, err)

	err = keeper.ExecuteMarketBuyOrder(ctx, order)
	require.NoError(t, err)

	// Should remove from sorted order
	res, _ := keeper.GetAllSortedSpotOrder(ctx)
	assert.Equal(t, res, [][]uint64{})

	// Should remove from main storage
	_, found := keeper.GetPendingSpotOrder(ctx, 1)
	assert.False(t, found)
}
