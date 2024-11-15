package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestCalcOutAmtGivenIn() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"pool not found",
			func() {
				suite.ResetSuite()
			},
			func() {
				_, _, err := suite.app.AmmKeeper.CalcOutAmtGivenIn(suite.ctx, 0, suite.app.OracleKeeper, nil, sdk.Coins{}, "", math.LegacyZeroDec())
				suite.Require().Error(err)
			},
		},
		{
			"pool found",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				addr := suite.AddAccounts(1, nil)[0]

				amount := math.NewInt(100000000000)
				pool := suite.CreateNewAmmPool(addr, true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))

				_, _, err := suite.app.AmmKeeper.CalcOutAmtGivenIn(suite.ctx, pool.PoolId, suite.app.OracleKeeper, &pool, sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000))), ptypes.ATOM, math.LegacyZeroDec())
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
