package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestCheckLowPoolHealth() {
	suite.ResetSuite()
	params := types.DefaultParams()
	params.PoolOpenThreshold = sdk.OneDec()
	err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)
	suite.Require().NoError(err)
	addr := suite.AddAccounts(10, nil)
	amount := sdk.NewInt(1000)
	poolCreator := addr[0]
	ammPool := suite.SetAndGetAmmPool(poolCreator, 1, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
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
			"Pool health is low",
			"pool (1) health too low to open new positions",
			func() {
				pool := types.NewPool(ammPool)
				pool.Health = sdk.MustNewDecFromStr("0.5")
				suite.app.PerpetualKeeper.SetPool(suite.ctx, pool)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			err = suite.app.PerpetualKeeper.CheckLowPoolHealth(suite.ctx, 1)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
