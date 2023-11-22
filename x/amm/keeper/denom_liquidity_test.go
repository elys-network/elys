package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNDenomLiquidity(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.DenomLiquidity {
	items := make([]types.DenomLiquidity, n)
	for i := range items {
		items[i].Denom = strconv.Itoa(i)
		items[i].Liquidity = sdk.ZeroInt()

		keeper.SetDenomLiquidity(ctx, items[i])
	}
	return items
}

func TestDenomLiquidityGet(t *testing.T) {
	keeper, ctx, _, _ := keepertest.AmmKeeper(t)
	items := createNDenomLiquidity(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetDenomLiquidity(ctx,
			item.Denom,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestDenomLiquidityRemove(t *testing.T) {
	keeper, ctx, _, _ := keepertest.AmmKeeper(t)
	items := createNDenomLiquidity(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveDenomLiquidity(ctx,
			item.Denom,
		)
		_, found := keeper.GetDenomLiquidity(ctx,
			item.Denom,
		)
		require.False(t, found)
	}
}

func TestDenomLiquidityGetAll(t *testing.T) {
	keeper, ctx, _, _ := keepertest.AmmKeeper(t)
	items := createNDenomLiquidity(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllDenomLiquidity(ctx)),
	)
}
