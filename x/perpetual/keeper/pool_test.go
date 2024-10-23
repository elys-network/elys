package keeper_test

import (
	"cosmossdk.io/math"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

func createNPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Pool {
	items := make([]types.Pool, n)
	for i := range items {
		poolId := uint64(i)
		poolAssets := []ammtypes.PoolAsset{
			{
				Token:  sdk.NewCoin(ptypes.BaseCurrency, math.OneInt().MulRaw(1000_000)),
				Weight: sdk.NewInt(10),
			},
			{
				Token:  sdk.NewCoin(ptypes.ATOM, math.OneInt().MulRaw(1000_000)),
				Weight: sdk.NewInt(10),
			},
		}
		ammPool := ammtypes.Pool{
			PoolId:            poolId,
			Address:           ammtypes.NewPoolAddress(poolId).String(),
			RebalanceTreasury: ammtypes.NewPoolRebalanceTreasury(poolId).String(),
			PoolParams: ammtypes.PoolParams{
				UseOracle:                   true,
				ExternalLiquidityRatio:      sdk.NewDec(2),
				WeightBreakingFeeMultiplier: math.LegacyZeroDec(),
				WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
				WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
				ThresholdWeightDifference:   math.LegacyZeroDec(),
				SwapFee:                     math.LegacyZeroDec(),
				ExitFee:                     math.LegacyZeroDec(),
				FeeDenom:                    ptypes.BaseCurrency,
			},
			TotalShares: sdk.NewCoin("pool/1", sdk.NewInt(100)),
			PoolAssets:  poolAssets,
			TotalWeight: sdk.ZeroInt(),
		}
		items[i] = types.NewPool(ammPool)

		keeper.SetPool(ctx, items[i])
	}
	return items
}

func createNPoolResponse(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.PoolResponse {
	items := make([]types.PoolResponse, n)
	for i := range items {
		items[i] = types.PoolResponse{
			AmmPoolId:                            uint64(i),
			Health:                               sdk.NewDec(1),
			BorrowInterestRate:                   sdk.MustNewDecFromStr("0.000000000000000001"),
			PoolAssetsLong:                       []types.PoolAsset{},
			PoolAssetsShort:                      []types.PoolAsset{},
			LastHeightBorrowInterestRateComputed: 0,
			FundingRate:                          sdk.ZeroDec(),
			NetOpenInterest:                      sdk.ZeroInt(),
		}
		ammPool, _ := ammtypes.NewBalancerPool(uint64(i), ammtypes.PoolParams{}, []ammtypes.PoolAsset{}, ctx.BlockTime())
		keeper.SetPool(ctx, types.NewPool(ammPool))
	}
	return items
}

func TestPoolGet(t *testing.T) {
	keeper, ctx := keepertest.PerpetualKeeper(t)
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
	keeper, ctx := keepertest.PerpetualKeeper(t)
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
	keeper, ctx := keepertest.PerpetualKeeper(t)
	items := createNPool(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPools(ctx)),
	)
}
