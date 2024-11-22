package keeper_test

import (
	"github.com/elys-network/elys/x/tradeshield/types"
)

<<<<<<< HEAD
func (suite *TradeshieldKeeperTestSuite) TestGetParams() {
=======
func TestGetParams(t *testing.T) {
	k, ctx, _, _, _, _ := testkeeper.TradeshieldKeeper(t)
>>>>>>> main
	params := types.DefaultParams()
	suite.app.TradeshieldKeeper.SetParams(suite.ctx, &params)

	suite.Require().EqualValues(params, suite.app.TradeshieldKeeper.GetParams(suite.ctx))
}
