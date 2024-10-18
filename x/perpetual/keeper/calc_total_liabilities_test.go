package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite PerpetualKeeperTestSuite) TestCalcTotalLiabilities() {
	suite.SetupCoinPrices()
	addr := suite.AddAccounts(1, nil)
	var ammPool ammtypes.Pool
	poolId := uint64(1)
	poolAsset := types.PoolAsset{
		Liabilities:           sdk.ZeroInt(),
		Custody:               sdk.ZeroInt(),
		TakeProfitLiabilities: sdk.ZeroInt(),
		TakeProfitCustody:     sdk.ZeroInt(),
		AssetDenom:            "",
	}
	testCases := []struct {
		name                 string
		expectErrMsg         string
		asset                string
		prerequisiteFunction func()
		postValidateFunction func(totalLiabilities sdk.Int)
	}{
		{
			"success: liabilities is 0",
			"",
			ptypes.ATOM,
			func() {
			},
			func(totalLiabilities sdk.Int) {
				suite.Require().True(totalLiabilities.IsZero())
			},
		},
		{
			"success: asset is uusdc",
			"",
			ptypes.BaseCurrency,
			func() {
				poolAsset.Liabilities = sdk.OneInt()
			},
			func(totalLiabilities sdk.Int) {
				suite.Require().True(totalLiabilities.Equal(sdk.OneInt()))
			},
		},
		{
			"amm pool not found",
			"pool does not exist",
			ptypes.ATOM,
			func() {
				poolAsset.Liabilities = sdk.OneInt()
			},
			func(totalLiabilities sdk.Int) {
			},
		},
		{
			"amm pool does not have enough funds",
			"amount too low",
			ptypes.ATOM,
			func() {
				amount := sdk.OneInt().MulRaw(1000_000)
				ammPool = suite.SetAndGetAmmPool(addr[0], poolId, true, sdk.ZeroDec(), sdk.ZeroDec(), ptypes.ATOM, amount, amount)
				poolAsset.Liabilities = amount
			},
			func(totalLiabilities sdk.Int) {
			},
		},
		{
			"success swap",
			"",
			ptypes.ATOM,
			func() {
				poolAsset.Liabilities = sdk.OneInt().MulRaw(100)
			},
			func(totalLiabilities sdk.Int) {
				uusdcAmount, _, err := suite.app.AmmKeeper.CalcInAmtGivenOut(suite.ctx, uint64(1), suite.app.OracleKeeper, &ammPool, sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, poolAsset.Liabilities)), ptypes.BaseCurrency, sdk.ZeroDec())
				suite.Require().NoError(err)
				suite.Require().True(totalLiabilities.Equal(uusdcAmount.Amount))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			poolAsset.AssetDenom = tc.asset
			totalLiabilities, err := suite.app.PerpetualKeeper.CalcTotalLiabilities(suite.ctx, []types.PoolAsset{poolAsset}, poolId, "uusdc")
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
			tc.postValidateFunction(totalLiabilities)
		})
	}
}
