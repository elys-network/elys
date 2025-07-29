package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	leveragelpmodulekeeper "github.com/elys-network/elys/v7/x/leveragelp/keeper"
	leveragelpmoduletypes "github.com/elys-network/elys/v7/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *PerpetualKeeperTestSuite) TestBeginBlocker() {
	suite.SetupCoinPrices()

	addr := suite.AddAccounts(1, nil)
	amount := math.NewInt(1000)

	poolCreator := addr[0]
	ammPool := suite.CreateNewAmmPool(poolCreator, true, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ptypes.ATOM, amount.MulRaw(10), amount.MulRaw(10))
	enablePoolMsg := leveragelpmoduletypes.MsgAddPool{
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Pool: leveragelpmoduletypes.AddPool{
			AmmPoolId:   ammPool.PoolId,
			LeverageMax: math.LegacyMustNewDecFromStr("10"),
		},
	}
	_, err := leveragelpmodulekeeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper).AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	suite.app.PerpetualKeeper.BeginBlocker(suite.ctx)
}

func (suite *PerpetualKeeperTestSuite) TestComputeFundingRate() {
	ammPool := ammtypes.Pool{
		PoolId: 1,
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token: sdk.Coin{
					Denom:  "testAsset",
					Amount: math.NewInt(100),
				},
			},
		},
	}
	pool := types.NewPool(ammPool, math.LegacyMustNewDecFromStr("5.5"))
	pool.PoolAssetsLong = []types.PoolAsset{
		{
			Liabilities: math.NewInt(0),
			Custody:     math.NewInt(0),
			Collateral:  math.NewInt(0),
			AssetDenom:  "testAsset",
		},
	}
	pool.PoolAssetsShort = []types.PoolAsset{
		{
			Liabilities: math.NewInt(0),
			Custody:     math.NewInt(0),
			Collateral:  math.NewInt(0),
			AssetDenom:  "testAsset",
		},
	}

	params := suite.app.PerpetualKeeper.GetParams(suite.ctx)
	params.FixedFundingRate = math.LegacyMustNewDecFromStr("0.500000000000000000")
	suite.app.PerpetualKeeper.SetParams(suite.ctx, &params)

	testCases := []struct {
		name            string
		expectLongRate  math.LegacyDec
		expectShortRate math.LegacyDec
		Liabilities     math.Int
		Custody         math.Int
		Collateral      math.Int
	}{
		{
			"totalLongOpenInterest is zero",
			math.LegacyNewDec(0),
			math.LegacyNewDec(0),
			math.NewInt(0),
			math.NewInt(0),
			math.NewInt(0),
		},
		{
			"totalShortOpenInterest is zero",
			math.LegacyNewDec(0),
			math.LegacyNewDec(0),
			math.NewInt(0),
			math.NewInt(0),
			math.NewInt(0),
		},
		{
			"totalShortOpenInterest is less than totalLongOpenInterest",
			math.LegacyMustNewDecFromStr("0.166666666666666666"),
			math.LegacyNewDec(0),
			math.NewInt(500),
			math.NewInt(2000),
			math.NewInt(1000),
		},
		{
			"totalLongOpenInterest is less than totalShortOpenInterest",
			math.LegacyNewDec(0),
			math.LegacyMustNewDecFromStr("0.1"),
			math.NewInt(1500),
			math.NewInt(2000),
			math.NewInt(1000),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			pool.PoolAssetsLong[0].Custody = tc.Custody
			pool.PoolAssetsLong[0].Collateral = tc.Collateral
			pool.PoolAssetsShort[0].Liabilities = tc.Liabilities
			longRate, shortRate := suite.app.PerpetualKeeper.ComputeFundingRate(suite.ctx, pool)
			suite.Require().Equal(tc.expectLongRate, longRate)
			suite.Require().Equal(tc.expectShortRate, shortRate)
		})
	}
}
