package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestQuerySwapEstimationExactAmountOut() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"nil request",
			func() {
				suite.ResetSuite()
			},
			func() {
				_, err := suite.app.AmmKeeper.SwapEstimationExactAmountOut(suite.ctx, nil)
				suite.Require().Error(err)
			},
		},
		{
			"no routes",
			func() {
				suite.ResetSuite()
			},
			func() {
				_, err := suite.app.AmmKeeper.SwapEstimationExactAmountOut(suite.ctx, &types.QuerySwapEstimationExactAmountOutRequest{})
				suite.Require().Error(err)
			},
		},
		{
			"valid request",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				addr := suite.AddAccounts(1, nil)[0]

				amount := math.NewInt(100000000000)
				pool := suite.CreateNewAmmPool(addr, true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))

				_, err := suite.app.AmmKeeper.SwapEstimationExactAmountOut(suite.ctx, &types.QuerySwapEstimationExactAmountOutRequest{
					Routes: []*types.SwapAmountOutRoute{
						{
							PoolId:       pool.PoolId,
							TokenInDenom: ptypes.ATOM,
						},
					},
					TokenOut: sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000)),
					Discount: math.LegacyZeroDec(),
				})
				suite.Require().NoError(err)
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
