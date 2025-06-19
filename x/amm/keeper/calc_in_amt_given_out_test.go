package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *AmmKeeperTestSuite) TestCalcInAmtGivenOut() {
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
				_, _, err := suite.app.AmmKeeper.CalcInAmtGivenOut(suite.ctx, 0, suite.app.OracleKeeper, ammtypes.SnapshotPool{}, sdk.Coins{}, "", osmomath.ZeroBigDec())
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
				pool := suite.CreateNewAmmPool(addr, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
				snapshot := suite.app.AmmKeeper.GetPoolWithAccountedBalance(suite.ctx, pool.PoolId)
				_, _, err := suite.app.AmmKeeper.CalcInAmtGivenOut(suite.ctx, pool.PoolId, suite.app.OracleKeeper, snapshot, sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000))), ptypes.ATOM, osmomath.ZeroBigDec())
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
