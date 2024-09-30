package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/elys-network/elys/x/amm/keeper"
)

func (suite *KeeperTestSuite) TestLiquidityRatioFromPriceDepth() {
	depth := sdkmath.LegacyNewDecWithPrec(1, 2) // 1%
	suite.Require().Equal("0.005012562893380045", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(2, 2) // 2%
	suite.Require().Equal("0.010050506338833466", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(5, 2) // 5%
	suite.Require().Equal("0.025320565519103609", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(10, 2) // 10%
	suite.Require().Equal("0.051316701949486200", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(30, 2) // 30%
	suite.Require().Equal("0.163339973465924452", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(50, 2) // 50%
	suite.Require().Equal("0.292893218813452475", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(70, 2) // 70%
	suite.Require().Equal("0.452277442494833886", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(90, 2) // 90%
	suite.Require().Equal("0.683772233983162067", keeper.LiquidityRatioFromPriceDepth(depth).String())
	depth = sdkmath.LegacyNewDecWithPrec(100, 2) // 100%
	suite.Require().Equal("1.000000000000000000", keeper.LiquidityRatioFromPriceDepth(depth).String())
}
