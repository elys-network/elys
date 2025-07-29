package keeper_test

import (
	"cosmossdk.io/math"

	"github.com/elys-network/elys/v7/x/masterchef/types"
)

func (suite *MasterchefKeeperTestSuite) TestShowFeeInfoQuery() {
	newInfo := types.FeeInfo{
		GasLp:        math.NewInt(300),
		GasStakers:   math.NewInt(150),
		GasProtocol:  math.NewInt(75),
		DexLp:        math.NewInt(400),
		DexStakers:   math.NewInt(200),
		DexProtocol:  math.NewInt(100),
		PerpLp:       math.NewInt(500),
		PerpStakers:  math.NewInt(250),
		PerpProtocol: math.NewInt(125),
		EdenLp:       math.NewInt(2000),
	}
	suite.app.MasterchefKeeper.SetFeeInfo(suite.ctx, newInfo, "2024-05-01")

	response, err := suite.app.MasterchefKeeper.ShowFeeInfo(suite.ctx, &types.QueryShowFeeInfoRequest{Date: "2024-05-01"})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryShowFeeInfoResponse{FeeInfo: newInfo}, response)
}

func (suite *MasterchefKeeperTestSuite) TestListFeeInfoQuery() {
	newInfo := types.FeeInfo{
		GasLp:        math.NewInt(300),
		GasStakers:   math.NewInt(150),
		GasProtocol:  math.NewInt(75),
		DexLp:        math.NewInt(400),
		DexStakers:   math.NewInt(200),
		DexProtocol:  math.NewInt(100),
		PerpLp:       math.NewInt(500),
		PerpStakers:  math.NewInt(250),
		PerpProtocol: math.NewInt(125),
		EdenLp:       math.NewInt(2000),
	}
	currData := suite.ctx.BlockTime().Format("2006-01-02")
	suite.app.MasterchefKeeper.SetFeeInfo(suite.ctx, newInfo, currData)
	currData = suite.ctx.BlockTime().AddDate(0, 0, 1).Format("2006-01-02")
	suite.app.MasterchefKeeper.SetFeeInfo(suite.ctx, newInfo, currData)

	response, err := suite.app.MasterchefKeeper.ListFeeInfo(suite.ctx, &types.QueryListFeeInfoRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryListFeeInfoResponse{FeeInfo: []types.FeeInfo{newInfo, newInfo}}, response)
}
