package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestCalcMinCollateral() {

	leverage := sdk.MustNewDecFromStr("2.5")
	price := sdk.MustNewDecFromStr("4.75")
	suite.app.PerpetualKeeper.SetParams(suite.ctx, &types.DefaultGenesis().Params)
	wantCollateral := sdk.NewInt(int64(1000033))

	collateral, err := suite.app.PerpetualKeeper.CalcMinCollateral(suite.ctx, leverage, price, 6)
	suite.Require().Nil(err)
	suite.Equal(wantCollateral, collateral)

	wrongLeverage := sdk.MustNewDecFromStr("0.89")
	_, err = suite.app.PerpetualKeeper.CalcMinCollateral(suite.ctx, wrongLeverage, price, 6)
	suite.Require().ErrorIs(err, types.ErrInvalidLeverage)
}
