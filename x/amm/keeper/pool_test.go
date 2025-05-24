package keeper_test

import (
	"strconv"
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v5/app"
	keepertest "github.com/elys-network/elys/v5/testutil/keeper"
	"github.com/elys-network/elys/v5/testutil/nullify"
	"github.com/elys-network/elys/v5/x/amm/keeper"
	"github.com/elys-network/elys/v5/x/amm/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/stretchr/testify/require"
)

func createNPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Pool {
	items := make([]types.Pool, n)
	for i := range items {
		items[i].PoolId = uint64(i)
		items[i].Address = types.NewPoolAddress(uint64(i)).String()
		items[i].TotalWeight = sdkmath.NewInt(100)
		items[i].TotalShares = sdk.NewCoin(types.GetPoolShareDenom(uint64(i)), types.OneShare)
		items[i].PoolParams = types.PoolParams{
			SwapFee:   sdkmath.LegacyZeroDec(),
			UseOracle: false,
			FeeDenom:  ptypes.BaseCurrency,
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
	app := simapp.InitElysTestApp(initChain, t)
	ctx := app.BaseApp.NewContext(initChain)
	keeper := app.AmmKeeper
	items := createNPool(keeper, ctx, 10)

	// Add assets to some pools for testing
	for i, item := range items {
		item.PoolAssets = append(item.PoolAssets, types.PoolAsset{
			Token: sdk.Coin{Denom: "denom" + strconv.Itoa(i), Amount: sdkmath.NewInt(1000)},
		})

		if i%2 == 0 { // Add "usdc" to every other pool
			item.PoolAssets = append(item.PoolAssets, types.PoolAsset{
				Token: sdk.Coin{Denom: "usdc", Amount: sdkmath.NewInt(1000)},
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

func (suite *AmmKeeperTestSuite) TestPool() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"GetAllLegacyPool returns empty list",
			func() {
				suite.ResetSuite()
			},
			func() {
				list := suite.app.AmmKeeper.GetAllLegacyPool(suite.ctx)
				suite.Require().Len(list, 0)
			},
		},
		{
			"GetAllLegacyPool returns list with one item",
			func() {
				suite.ResetSuite()

				suite.app.AmmKeeper.SetPool(suite.ctx, types.Pool{
					PoolId:  1,
					Address: types.NewPoolAddress(1).String(),
				})
			},
			func() {
				list := suite.app.AmmKeeper.GetAllLegacyPool(suite.ctx)
				suite.Require().Len(list, 1)
			},
		},
		{
			"IterateLiquidityPools",
			func() {
				suite.ResetSuite()

				suite.app.AmmKeeper.SetPool(suite.ctx, types.Pool{
					PoolId:  1,
					Address: types.NewPoolAddress(1).String(),
				})
			},
			func() {
				suite.app.AmmKeeper.IterateLiquidityPools(suite.ctx, func(pool types.Pool) (stop bool) {
					suite.Require().Equal(uint64(1), pool.PoolId)
					return true
				})
			},
		},
		{
			"AddToPoolBalance with non-existent pool asset",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				addr := suite.AddAccounts(1, nil)[0]

				amount := sdkmath.NewInt(100000000000)
				pool := suite.CreateNewAmmPool(addr, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))

				err := suite.app.AmmKeeper.AddToPoolBalanceAndUpdateLiquidity(suite.ctx, &pool, sdkmath.ZeroInt(), sdk.NewCoins(sdk.NewCoin("non-existant-denom", amount)))
				suite.Require().Error(err)
			},
		},
		{
			"AddToPoolBalance with existent pool",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				addr := suite.AddAccounts(1, nil)[0]

				amount := sdkmath.NewInt(100000000000)
				pool := suite.CreateNewAmmPool(addr, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))

				err := suite.app.AmmKeeper.AddToPoolBalanceAndUpdateLiquidity(suite.ctx, &pool, sdkmath.ZeroInt(), sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, amount)))
				suite.Require().NoError(err)
			},
		},
		{
			"RemoveFromPoolBalance with non-existent pool asset",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				addr := suite.AddAccounts(1, nil)[0]

				amount := sdkmath.NewInt(100000000000)
				pool := suite.CreateNewAmmPool(addr, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))

				err := suite.app.AmmKeeper.RemoveFromPoolBalanceAndUpdateLiquidity(suite.ctx, &pool, sdkmath.ZeroInt(), sdk.NewCoins(sdk.NewCoin("non-existant-denom", amount)))
				suite.Require().Error(err)
			},
		},
		{
			"RemoveFromPoolBalance with existent pool",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				addr := suite.AddAccounts(1, nil)[0]

				amount := sdkmath.NewInt(100000000000)
				pool := suite.CreateNewAmmPool(addr, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))

				err := suite.app.AmmKeeper.RemoveFromPoolBalanceAndUpdateLiquidity(suite.ctx, &pool, sdkmath.ZeroInt(), sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, amount)))
				suite.Require().NoError(err)
			},
		},
		{
			"PoolExists",
			func() {
				suite.ResetSuite()
			},
			func() {
				suite.Require().False(suite.app.AmmKeeper.PoolExists(suite.ctx, 1))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			tc.postValidateFunction()
		})
	}
}
