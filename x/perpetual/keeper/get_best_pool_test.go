package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestGetBestPool() {
	suite.SetupCoinPrices()
	addr := suite.AddAccounts(1, nil)
	amount := sdk.NewInt(10_000_000)
	expectAmmPool := ammtypes.Pool{}
	testCases := []struct {
		name                 string
		tradingAsset         string
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{
			"pool with assets not found",
			"unknownAsset",
			types.ErrPoolDoesNotExist.Error(),
			func() {
			},
		},
		{
			"pool exists not oracle pool",
			ptypes.ATOM,
			types.ErrPoolDoesNotExist.Error(),
			func() {
				_ = suite.SetAndGetAmmPool(addr[0], 1, false, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount, amount)
			},
		},
		{
			"only 1 pool exists",
			ptypes.ATOM,
			"",
			func() {
				expectAmmPool = suite.SetAndGetAmmPool(addr[0], 2, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount, amount)
			},
		},
		{
			"multiple pool exists",
			ptypes.ATOM,
			"",
			func() {
				expectAmmPool = suite.SetAndGetAmmPool(addr[0], 3, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			pool, err := suite.app.PerpetualKeeper.GetBestPool(suite.ctx, ptypes.BaseCurrency, tc.tradingAsset)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(expectAmmPool.PoolId, pool)
			}
		})
	}
}
