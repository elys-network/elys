package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *PerpetualKeeperTestSuite) TestCheckLowPoolHealth() {
	suite.ResetSuite()
	params := types.DefaultParams()
	err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
	suite.Require().NoError(err)
	addr := suite.AddAccounts(10, nil)
	amount := sdkmath.NewInt(1000)
	poolCreator := addr[0]
	suite.SetupCoinPrices()
	ammPool := suite.CreateNewAmmPool(poolCreator, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{
			"Pool not found",
			types.ErrPoolDoesNotExist.Error(),
			func() {
			},
		},
		// "Pool health is nil" case is not possible because Getter function always give 0 value of health
		{
			"Pool health is low LONG",
			"pool (1) base asset liabilities ratio (0.950000000000000000) too high for the operation",
			func() {
				pool := types.NewPool(ammPool, sdkmath.LegacyMustNewDecFromStr("10.5"))
				pool.BaseAssetLiabilitiesRatio = sdkmath.LegacyMustNewDecFromStr("0.95")
				suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
			},
		},
		{
			"Pool health is low SHORT",
			"pool (1) quote asset liabilities ratio (0.950000000000000000) too high for the operation",
			func() {
				pool := types.NewPool(ammPool, sdkmath.LegacyMustNewDecFromStr("10.5"))
				pool.QuoteAssetLiabilitiesRatio = sdkmath.LegacyMustNewDecFromStr("0.95")
				suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			err = suite.app.PerpetualKeeper.CheckLowPoolHealthAndMinimumCustody(suite.ctx, 1)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
