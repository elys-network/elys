package keeper_test

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/masterchef/types"
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
	keeper, ctx := testkeeper.MasterchefKeeper(suite.T())
	wctx := sdk.WrapSDKContext(ctx)
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
	keeper.SetFeeInfo(ctx, newInfo, "2024-05-01")
	keeper.SetFeeInfo(ctx, newInfo, "2024-05-02")

	response, err := keeper.ListFeeInfo(wctx, &types.QueryListFeeInfoRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryListFeeInfoResponse{FeeInfo: []types.FeeInfo{newInfo, newInfo}}, response)
}
