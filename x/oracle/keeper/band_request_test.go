package keeper_test

import (
	"github.com/elys-network/elys/x/oracle/types"
)

func (suite *KeeperTestSuite) TestBandRequestSetGet() {
	calldata := []types.BandPriceCallData{
		{
			Symbols:    []string{"BTC", "ETH"},
			Multiplier: 18,
		},
		{
			Symbols:    []string{"BTC"},
			Multiplier: 18,
		},
	}

	suite.app.OracleKeeper.SetBandRequest(suite.ctx, types.OracleRequestID(1), calldata[0])
	suite.app.OracleKeeper.SetBandRequest(suite.ctx, types.OracleRequestID(2), calldata[1])

	req, err := suite.app.OracleKeeper.GetBandRequest(suite.ctx, types.OracleRequestID(1))
	suite.Require().NoError(err)
	suite.Require().Equal(req, calldata[0])

	req, err = suite.app.OracleKeeper.GetBandRequest(suite.ctx, types.OracleRequestID(3))
	suite.Require().Error(err)
}
