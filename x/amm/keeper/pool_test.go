package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Pool {
	items := make([]types.Pool, n)
	for i := range items {
		items[i].PoolId = uint64(i)
		items[i].TotalWeight = sdk.NewInt(100)
		items[i].PoolParams = types.PoolParams{
			SwapFee:                     sdk.ZeroDec(),
			ExitFee:                     sdk.ZeroDec(),
			UseOracle:                   false,
			WeightBreakingFeeMultiplier: sdk.ZeroDec(),
			ExternalLiquidityRatio:      sdk.ZeroDec(),
			LpFeePortion:                sdk.ZeroDec(),
			StakingFeePortion:           sdk.ZeroDec(),
			WeightRecoveryFeePortion:    sdk.ZeroDec(),
			ThresholdWeightDifference:   sdk.ZeroDec(),
			FeeDenom:                    ptypes.BaseCurrency,
		}

		keeper.SetPool(ctx, items[i])
	}
	return items
}

func TestPoolGet(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	items := createNPool(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetPool(ctx,
			item.PoolId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestPoolRemove(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	items := createNPool(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePool(ctx,
			item.PoolId,
		)
		_, found := keeper.GetPool(ctx,
			item.PoolId,
		)
		require.False(t, found)
	}
}

func TestPoolGetAll(t *testing.T) {
	keeper, ctx := keepertest.AmmKeeper(t)
	items := createNPool(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPool(ctx)),
	)
}
