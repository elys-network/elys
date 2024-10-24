package keeper_test

import (
	"strconv"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func createNPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Pool {
	items := make([]types.Pool, n)
	for i := range items {
		items[i].PoolId = uint64(i)
		items[i].TotalWeight = sdk.NewInt(100)
		items[i].TotalShares = sdk.NewCoin(types.GetPoolShareDenom(uint64(i)), types.OneShare)
		items[i].PoolParams = types.PoolParams{
			SwapFee:                     sdk.ZeroDec(),
			ExitFee:                     sdk.ZeroDec(),
			UseOracle:                   false,
			WeightBreakingFeeMultiplier: sdk.ZeroDec(),
			WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
			ExternalLiquidityRatio:      sdk.ZeroDec(),
			WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
			ThresholdWeightDifference:   sdk.ZeroDec(),
			WeightBreakingFeePortion:    sdk.NewDecWithPrec(50, 2), // 50%
			FeeDenom:                    ptypes.BaseCurrency,
		}

		keeper.SetPool(ctx, items[i])
	}
	return items
}

func TestPoolGet(t *testing.T) {
	keeper, ctx, _, _ := keepertest.AmmKeeper(t)
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
	keeper, ctx, _, _ := keepertest.AmmKeeper(t)
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
	keeper, ctx, _, _ := keepertest.AmmKeeper(t)
	items := createNPool(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPool(ctx)),
	)
}

func TestGetBestPoolWithDenoms(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})
	keeper := app.AmmKeeper
	items := createNPool(&keeper, ctx, 10)

	// Add assets to some pools for testing
	for i, item := range items {
		item.PoolAssets = append(item.PoolAssets, types.PoolAsset{
			Token: sdk.Coin{Denom: "denom" + strconv.Itoa(i), Amount: sdk.NewInt(1000)},
		})

		if i%2 == 0 { // Add "usdc" to every other pool
			item.PoolAssets = append(item.PoolAssets, types.PoolAsset{
				Token: sdk.Coin{Denom: "usdc", Amount: sdk.NewInt(1000)},
			})
		}

		keeper.SetPool(ctx, item)
	}

	// Test case where pool is found
	pool, found := keeper.GetBestPoolWithDenoms(ctx, []string{"denom2", "usdc"}, false)
	require.True(t, found)
	require.Equal(t, uint64(2), pool.PoolId)

	// Test case where pool is not found
	_, found = keeper.GetBestPoolWithDenoms(ctx, []string{"nonexistent", "usdc"}, false)
	require.False(t, found)
}
