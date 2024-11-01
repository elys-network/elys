package keeper_test

import (
	sdkmath "cosmossdk.io/math"
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
		Liabilities:           sdkmath.ZeroInt(),
		Custody:               sdkmath.ZeroInt(),
		TakeProfitLiabilities: sdkmath.ZeroInt(),
		TakeProfitCustody:     sdkmath.ZeroInt(),
		AssetDenom:            "",
	}
	testCases := []struct {
		name                 string
		expectErrMsg         string
		asset                string
		prerequisiteFunction func()
		postValidateFunction func(totalLiabilities sdkmath.Int)
	}{
		{
			"success: liabilities is 0",
			"",
			ptypes.ATOM,
			func() {
			},
			func(totalLiabilities sdkmath.Int) {
				suite.Require().True(totalLiabilities.IsZero())
			},
		},
		{
			"success: asset is uusdc",
			"",
			ptypes.BaseCurrency,
			func() {
				poolAsset.Liabilities = sdkmath.OneInt()
			},
			func(totalLiabilities sdkmath.Int) {
				suite.Require().True(totalLiabilities.Equal(sdkmath.OneInt()))
			},
		},
		{
			"amm pool not found",
			"pool does not exist",
			ptypes.ATOM,
			func() {
				poolAsset.Liabilities = sdkmath.OneInt()
			},
			func(totalLiabilities sdkmath.Int) {
			},
		},
		{
			"amm pool does not have enough funds",
			"amount too low",
			ptypes.ATOM,
			func() {
				amount := sdkmath.OneInt().MulRaw(1000_000)
				ammPool = suite.CreateNewAmmPool(addr[0], true, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ptypes.ATOM, amount, amount)
				poolAsset.Liabilities = amount.MulRaw(100)
			},
			func(totalLiabilities sdkmath.Int) {
			},
		},
		{
			"success swap",
			"",
			ptypes.ATOM,
			func() {
				poolAsset.Liabilities = sdkmath.OneInt().MulRaw(100)
			},
			func(totalLiabilities sdkmath.Int) {
				uusdcAmount, _, err := suite.app.AmmKeeper.CalcInAmtGivenOut(suite.ctx, uint64(1), suite.app.OracleKeeper, &ammPool, sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, poolAsset.Liabilities)), ptypes.BaseCurrency, sdkmath.LegacyZeroDec())
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
