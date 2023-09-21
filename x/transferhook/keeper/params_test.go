package keeper_test

import (
	"github.com/elys-network/elys/x/transferhook/types"
)

func (s *KeeperTestSuite) TestGetParams() {
	params := types.DefaultParams()
	params.AmmActive = false

	s.app.TransferhookKeeper.SetParams(s.ctx, params)

	s.Require().Equal(params, s.app.TransferhookKeeper.GetParams(s.ctx))
}
