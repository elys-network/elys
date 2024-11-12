package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestCalcMinCollateral() {

	leverage := math.LegacyMustNewDecFromStr("2.5")
	price := math.LegacyMustNewDecFromStr("4.75")
	suite.app.PerpetualKeeper.SetParams(suite.ctx, &types.DefaultGenesis().Params)
	wantCollateral := math.NewInt(int64(1000033))

	collateral, err := suite.app.PerpetualKeeper.CalcMinCollateral(suite.ctx, leverage, price, 6)
	suite.Require().Nil(err)
	suite.Equal(wantCollateral, collateral)

	wrongLeverage := math.LegacyMustNewDecFromStr("0.89")
	_, err = suite.app.PerpetualKeeper.CalcMinCollateral(suite.ctx, wrongLeverage, price, 6)
	suite.Require().ErrorIs(err, types.ErrInvalidLeverage)
}
