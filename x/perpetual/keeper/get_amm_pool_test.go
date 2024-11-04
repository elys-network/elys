package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestGetAmmPool() {
	ammPool := ammtypes.Pool{
		PoolId:            1,
		Address:           ammtypes.NewPoolAddress(1).String(),
		RebalanceTreasury: "",
		PoolParams: ammtypes.PoolParams{
			UseOracle:                   false,
			WeightBreakingFeeMultiplier: sdk.ZeroDec(),
			WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
			WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
			ThresholdWeightDifference:   sdk.ZeroDec(),
			WeightBreakingFeePortion:    sdk.NewDecWithPrec(50, 2), // 50%
			SwapFee:                     sdk.ZeroDec(),
			ExitFee:                     sdk.ZeroDec(),
			FeeDenom:                    ptypes.BaseCurrency,
		},
		TotalShares: sdk.NewCoin("pool/1", sdk.NewInt(100)),
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:                  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
				Weight:                 sdk.NewInt(10),
				ExternalLiquidityRatio: sdk.NewDec(2),
			},
			{
				Token:                  sdk.NewCoin("borrowAsset", sdk.NewInt(10000)),
				Weight:                 sdk.NewInt(10),
				ExternalLiquidityRatio: sdk.NewDec(2),
			},
		},
		TotalWeight: sdk.ZeroInt(),
	}
	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{
			"pool not found",
			types.ErrPoolDoesNotExist.Error(),
			func() {
			},
		},
		{
			"pool found",
			"",
			func() {
				suite.app.AmmKeeper.SetPool(suite.ctx, ammPool)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			pool, err := suite.app.PerpetualKeeper.GetAmmPool(suite.ctx, 1)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(ammPool, pool)
			}
		})
	}
}
