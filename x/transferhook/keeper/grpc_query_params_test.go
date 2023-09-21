package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/transferhook/types"
)

func (s *KeeperTestSuite) TestParamsQuery() {
	// Test with amm enabled and amm enabled
	s.app.TransferhookKeeper.SetParams(s.ctx, types.Params{
		AmmActive: true,
	})
	queryResponse, err := s.app.TransferhookKeeper.Params(sdk.WrapSDKContext(s.ctx), &types.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().True(queryResponse.Params.AmmActive)

	// Test with amm disabled
	s.app.TransferhookKeeper.SetParams(s.ctx, types.Params{
		AmmActive: false,
	})
	queryResponse, err = s.app.TransferhookKeeper.Params(sdk.WrapSDKContext(s.ctx), &types.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().False(queryResponse.Params.AmmActive)
}
