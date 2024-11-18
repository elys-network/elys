package keeper_test

import (
	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *MasterchefKeeperTestSuite) TestCalculatePoolAprs() {

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(suite.app, suite.ctx, 1, sdkmath.NewInt(100010))

	// Mint 100000USDC + 10 ELYS (pool creation fee)
	coins := sdk.NewCoins(sdk.NewInt64Coin(ptypes.Elys, 10000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100000))
	suite.MintMultipleTokenToAddress(addr[0], coins)

	// Create pool
	var poolAssets []ammtypes.PoolAsset
	// Elys
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: sdkmath.NewInt(50),
		Token:  sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000)),
	})

	// USDC
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: sdkmath.NewInt(50),
		Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100)),
	})

	poolParams := &ammtypes.PoolParams{
		SwapFee:                     sdkmath.LegacyZeroDec(),
		ExitFee:                     sdkmath.LegacyZeroDec(),
		UseOracle:                   false,
		WeightBreakingFeeMultiplier: sdkmath.LegacyZeroDec(),
		WeightBreakingFeeExponent:   sdkmath.LegacyNewDecWithPrec(25, 1), // 2.5
		WeightRecoveryFeePortion:    sdkmath.LegacyNewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   sdkmath.LegacyZeroDec(),
		FeeDenom:                    "",
	}
	// Create a Elys+USDC pool
	ammPool := suite.CreateNewAmmPool(addr[0], poolAssets, poolParams)
	suite.Require().Equal(ammPool.PoolId, uint64(1))

	poolInfo, found := suite.app.MasterchefKeeper.GetPoolInfo(suite.ctx, ammPool.PoolId)
	suite.Require().True(found)

	poolInfo.DexApr = sdkmath.LegacyNewDecWithPrec(1, 2)  // 1%
	poolInfo.EdenApr = sdkmath.LegacyNewDecWithPrec(2, 2) // 2%
	suite.app.MasterchefKeeper.SetPoolInfo(suite.ctx, poolInfo)

	testCases := []struct {
		name          string
		ids           []uint64
		aprsExpectlen int
		expectValue   string
	}{
		{
			name:          "Empty poolIds",
			ids:           []uint64{},
			aprsExpectlen: 2, // setting it 2 because PoolId = math.MaxInt16 gets initiated in EndBlock
			expectValue:   "0.030000000000000000",
		},
		{
			name:          "Available pool's ids",
			ids:           []uint64{1},
			aprsExpectlen: 1,
			expectValue:   "0.030000000000000000",
		},
		{
			name:          "Pool not found, zero APRs",
			ids:           []uint64{4},
			aprsExpectlen: 1,
			expectValue:   "0.000000000000000000",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			aprs := suite.app.MasterchefKeeper.CalculatePoolAprs(suite.ctx, tc.ids)
			suite.Require().Equal(len(aprs), tc.aprsExpectlen)
			suite.Require().Equal(aprs[0].TotalApr.String(), tc.expectValue)
		})
	}
}
