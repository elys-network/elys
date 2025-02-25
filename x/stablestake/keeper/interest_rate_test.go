package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/stretchr/testify/require"

	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/stablestake/keeper"
	"github.com/elys-network/elys/x/stablestake/types"
)

func CreateNInterest(keeper *keeper.Keeper, ctx sdk.Context, n int) ([]types.InterestBlock, int64) {
	items := make([]types.InterestBlock, n)
	ctx = ctx.WithBlockHeight(1000)
	curBlock := ctx.BlockHeight()
	for i := range items {
		items[i].InterestRate = sdkmath.LegacyNewDec(int64(i))
		items[i].BlockTime = int64(i * 10)

		curBlock++
		keeper.SetInterestForPool(ctx, 1, uint64(curBlock), items[i])
	}
	return items, curBlock
}

func TestInterestGet(t *testing.T) {
	keeper, ctx := keepertest.StablestakeKeeper(t)
	_, lastBlock := CreateNInterest(keeper, ctx, 10)
	ctx = ctx.WithBlockHeight(lastBlock)

	// 1st case
	res := keeper.GetInterestForPool(ctx, uint64(ctx.BlockHeight()-2), uint64(ctx.BlockTime().Unix()-1), sdkmath.LegacyNewDec(86400*365), 1)
	require.Equal(t, res.Int64(), int64(8))

	// 2nd case
	res = keeper.GetInterestForPool(ctx, uint64(ctx.BlockHeight()-20), uint64(ctx.BlockTime().Unix()-1), sdkmath.LegacyNewDec(86400*365), 1)
	require.Equal(t, res.Int64(), int64(2))

	// 3rd case
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1000)
	res = keeper.GetInterestForPool(ctx, uint64(ctx.BlockHeight()-20), uint64(ctx.BlockTime().Unix()-1), sdkmath.LegacyNewDec(86400*365), 1)
	require.Equal(t, res.Int64(), int64(0))

	all := keeper.GetAllInterestForPool(ctx, 1)
	require.Equal(t, len(all), 10)

	keeper.DeleteInterestForPool(ctx, ctx.BlockHeight()-1, 1)
	all = keeper.GetAllInterestForPool(ctx, 1)
	require.Equal(t, len(all), 10)

	all = keeper.GetAllInterest(ctx)
	require.Equal(t, len(all), 10)
}

func (suite *KeeperTestSuite) TestInterestRateComputationForPool() {
	for _, tc := range []struct {
		desc    string
		pool    types.Pool
		expPass bool
		want    sdkmath.LegacyDec
	}{
		{
			desc: "interest calculation with zero total value",
			pool: types.Pool{
				TotalValue:           sdkmath.NewInt(0),
				InterestRate:         sdkmath.LegacyNewDec(5),
				InterestRateMax:      sdkmath.LegacyNewDec(10),
				InterestRateMin:      sdkmath.LegacyNewDec(1),
				InterestRateIncrease: sdkmath.LegacyNewDec(1),
				InterestRateDecrease: sdkmath.LegacyNewDec(1),
				HealthGainFactor:     sdkmath.LegacyNewDec(1),
			},
			expPass: true,
			want:    sdkmath.LegacyNewDec(5),
		},
		{
			desc: "interest calculation with zero max leverage",
			pool: types.Pool{
				TotalValue:           sdkmath.NewInt(0),
				InterestRate:         sdkmath.LegacyNewDec(5),
				InterestRateMax:      sdkmath.LegacyNewDec(10),
				InterestRateMin:      sdkmath.LegacyNewDec(1),
				InterestRateIncrease: sdkmath.LegacyNewDec(1),
				InterestRateDecrease: sdkmath.LegacyNewDec(1),
				HealthGainFactor:     sdkmath.LegacyNewDec(1),
				MaxLeverageRatio:     sdkmath.LegacyZeroDec(),
			},
			expPass: true,
			want:    sdkmath.LegacyNewDec(5),
		},
		{
			desc: "interest calculation with zero max leverage",
			pool: types.Pool{
				TotalValue:           sdkmath.NewInt(10000),
				InterestRate:         sdkmath.LegacyNewDec(12),
				InterestRateMax:      sdkmath.LegacyNewDec(17),
				InterestRateMin:      sdkmath.LegacyNewDec(12),
				InterestRateIncrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				InterestRateDecrease: sdkmath.LegacyMustNewDecFromStr("0.01"),
				HealthGainFactor:     sdkmath.LegacyNewDec(1),
				MaxLeverageRatio:     sdkmath.LegacyMustNewDecFromStr("0.8"),
				Id:                   1,
				DepositDenom:         ptypes.BaseCurrency,
			},
			expPass: true,
			want:    sdkmath.LegacyMustNewDecFromStr("12.01"),
		},
	} {
		suite.Run(tc.desc, func() {
			suite.SetupTest()

			got := suite.app.StablestakeKeeper.InterestRateComputationForPool(suite.ctx, tc.pool)
			suite.Require().Equal(tc.want, got)
		})
	}
}
