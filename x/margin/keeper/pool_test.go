package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/margin/keeper"
	"github.com/elys-network/elys/x/margin/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Pool {
	items := make([]types.Pool, n)
	for i := range items {
		items[i].AmmPoolId = (uint64)(i)
		items[i].Health = sdk.ZeroDec()
		items[i].Enabled = true
		items[i].Closed = false
		items[i].ExternalLiabilities = sdk.ZeroInt()
		items[i].ExternalCustody = sdk.ZeroInt()
		items[i].NativeLiabilities = sdk.ZeroInt()
		items[i].NativeCustody = sdk.ZeroInt()
		items[i].InterestRate = sdk.ZeroDec()
		items[i].NativeAssetBalance = sdk.ZeroInt()
		items[i].ExternalAssetBalance = sdk.ZeroInt()
		items[i].UnsettledExternalLiabilities = sdk.ZeroInt()
		items[i].BlockInterestNative = sdk.ZeroInt()

		keeper.SetPool(ctx, items[i])
	}
	return items
}

func TestPoolGet(t *testing.T) {
	keeper, ctx := keepertest.MarginKeeper(t)
	items := createNPool(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetPool(ctx,
			item.AmmPoolId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestPoolRemove(t *testing.T) {
	keeper, ctx := keepertest.MarginKeeper(t)
	items := createNPool(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePool(ctx,
			item.AmmPoolId,
		)
		_, found := keeper.GetPool(ctx,
			item.AmmPoolId,
		)
		require.False(t, found)
	}
}

func TestPoolGetAll(t *testing.T) {
	keeper, ctx := keepertest.MarginKeeper(t)
	items := createNPool(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPools(ctx)),
	)
}
