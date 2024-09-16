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

func createNPendingPerpetualOrder(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.PerpetualOrder {
	items := make([]types.PerpetualOrder, n)
	for i := range items {
		items[i].OrderId = keeper.AppendPendingPerpetualOrder(ctx, items[i])
	}
	return items
}

func TestPendingPerpetualOrderGet(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetPendingPerpetualOrder(ctx, item.OrderId)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestPendingPerpetualOrderRemove(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePendingPerpetualOrder(ctx, item.OrderId)
		_, found := keeper.GetPendingPerpetualOrder(ctx, item.OrderId)
		require.False(t, found)
	}
}

func TestPendingPerpetualOrderGetAll(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPendingPerpetualOrder(ctx)),
	)
}

func TestPendingPerpetualOrderCount(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)
	items := createNPendingPerpetualOrder(keeper, ctx, 10)
	count := uint64(len(items))
	require.Equal(t, count, keeper.GetPendingSpotOrderCount(ctx)-1)
}

func TestSortedPerpetualOrder(t *testing.T) {
	keeper, ctx := keepertest.TradeshieldKeeper(t)

	// Set to main storage
	keeper.AppendPendingPerpetualOrder(ctx, types.PerpetualOrder{
		OwnerAddress:       "address",
		OrderId:            0,
		PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
		TriggerPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdk.NewDec(1),
		},
	})

	order, _ := keeper.GetPendingPerpetualOrder(ctx, 1)

	err := keeper.InsertPerptualSortedOrder(ctx, order)
	require.NoError(t, err)

	res, _ := keeper.GetAllSortedPerpetualOrder(ctx)

	assert.Equal(t, res, [][]uint64{{1}})

	// Insert two more elements
	// Set to main storage
	keeper.AppendPendingPerpetualOrder(ctx, types.PerpetualOrder{
		OwnerAddress:       "address1",
		OrderId:            0,
		PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
		TriggerPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdk.NewDec(20),
		},
	})

	keeper.AppendPendingPerpetualOrder(ctx, types.PerpetualOrder{
		OwnerAddress:       "address2",
		OrderId:            0,
		PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
		TriggerPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdk.NewDec(5),
		},
	})

	keeper.AppendPendingPerpetualOrder(ctx, types.PerpetualOrder{
		OwnerAddress:       "address3",
		OrderId:            0,
		PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
		TriggerPrice: &types.OrderPrice{
			BaseDenom:  "base",
			QuoteDenom: "quote",
			Rate:       sdk.NewDec(25),
		},
	})

	order2, _ := keeper.GetPendingPerpetualOrder(ctx, 2)
	order3, _ := keeper.GetPendingPerpetualOrder(ctx, 3)
	order4, _ := keeper.GetPendingPerpetualOrder(ctx, 4)

	err = keeper.InsertPerptualSortedOrder(ctx, order2)
	require.NoError(t, err)
	err = keeper.InsertPerptualSortedOrder(ctx, order3)
	require.NoError(t, err)
	err = keeper.InsertPerptualSortedOrder(ctx, order4)
	require.NoError(t, err)

	res, _ = keeper.GetAllSortedPerpetualOrder(ctx)

	// Should store in sorted order
	assert.Equal(t, res, [][]uint64{{1, 3, 2, 4}})

	// Test binary search, search with rate 5
	index, err := keeper.PerpetualBinarySearch(ctx, sdk.NewDec(5), []uint64{1, 3, 2, 4})
	require.NoError(t, err)

	// second element
	assert.Equal(t, index, 1)

	// Test remove sorted order
	keeper.RemovePerpetualSortedOrder(ctx, 2)
	res, _ = keeper.GetAllSortedPerpetualOrder(ctx)

	// Should store in sorted order
	assert.Equal(t, res, [][]uint64{{1, 3, 4}})
}
