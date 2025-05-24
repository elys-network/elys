package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v5/x/amm/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestGetFundingPaymentRates() {
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
		name                   string
		expectLongRate         math.LegacyDec
		expectShortRate        math.LegacyDec
		expectLongFundingRate  math.LegacyDec
		expectShortFundingRate math.LegacyDec
		Liabilities            math.Int
		Custody                math.Int
		Collateral             math.Int
	}{
		{
			"fundingRateLong is zero, totalLongOpenInterest is non zero",
			math.LegacyNewDec(0),
			math.LegacyMustNewDecFromStr("0.1"),
			math.LegacyMustNewDecFromStr("-0.15"),
			math.LegacyMustNewDecFromStr("0.1"),
			math.NewInt(1500),
			math.NewInt(2000),
			math.NewInt(1000),
		},
		{
			"fundingRateLong is zero, totalLongOpenInterest is zero",
			math.LegacyNewDec(0),
			math.LegacyNewDec(0),
			math.LegacyNewDec(0),
			math.LegacyNewDec(0),
			math.NewInt(0),
			math.NewInt(0),
			math.NewInt(0),
		},
		{
			"fundingRateLong is non zero, totalShortOpenInterest is non zero",
			math.LegacyMustNewDecFromStr("0.166666666666666666"),
			math.LegacyNewDec(0),
			math.LegacyMustNewDecFromStr("0.166666666666666666"),
			math.LegacyMustNewDecFromStr("-0.333333333333333332"),
			math.NewInt(500),
			math.NewInt(2000),
			math.NewInt(1000),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			pool.PoolAssetsLong[0].Custody = tc.Custody
			pool.PoolAssetsLong[0].Collateral = tc.Collateral
			pool.PoolAssetsShort[0].Liabilities = tc.Liabilities
			longRate, shortRate := suite.app.PerpetualKeeper.GetFundingPaymentRates(suite.ctx, pool)
			suite.Require().Equal(tc.expectLongFundingRate, longRate)
			suite.Require().Equal(tc.expectShortFundingRate, shortRate)
		})
	}
}
