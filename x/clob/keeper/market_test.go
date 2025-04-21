package keeper_test

import (
	"github.com/elys-network/elys/x/clob/types"
)

func (suite *KeeperTestSuite) TestGetPerpetualMarket() {

	market := types.PerpetualMarket{
		Id: 1,
	}
	suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)
	// Define test cases
	testCases := []struct {
		name           string
		id             uint64
		expectedErrMsg string
	}{
		{
			"market not found",
			2,
			types.ErrPerpetualMarketNotFound.Error(),
		},
		{
			"market found",
			1,
			"",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {

			res, err := suite.app.ClobKeeper.GetPerpetualMarket(suite.ctx, tc.id)

			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.id, res.Id)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGetAllPerpetualMarkets() {

	market1 := types.PerpetualMarket{
		Id: 1,
	}
	market2 := types.PerpetualMarket{
		Id: 2,
	}
	testCases := []struct {
		name  string
		ids   []uint64
		setup func()
	}{
		{
			"empty markets",
			[]uint64{},
			func() {

			},
		},
		{
			"2 markets",
			[]uint64{1, 2},
			func() {
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market1)
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market2)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setup()
			res := suite.app.ClobKeeper.GetAllPerpetualMarket(suite.ctx)
			suite.Equal(len(res), len(tc.ids))
		})
	}
}

func (suite *KeeperTestSuite) TestCountAllPerpetualMarkets() {

	market1 := types.PerpetualMarket{
		Id: 1,
	}
	market2 := types.PerpetualMarket{
		Id: 2,
	}
	testCases := []struct {
		name   string
		result uint64
		setup  func()
	}{
		{
			"empty markets",
			0,
			func() {

			},
		},
		{
			"2 markets",
			2,
			func() {
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market1)
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market2)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setup()
			res := suite.app.ClobKeeper.CountAllPerpetualMarket(suite.ctx)
			suite.Equal(res, tc.result)
		})
	}
}

func (suite *KeeperTestSuite) TestCheckPerpetualMarketAlreadyExists() {

	baseDenom := "uatom"
	quoteDenom := QuoteDenom

	market1 := types.PerpetualMarket{
		Id:         1,
		BaseDenom:  baseDenom,
		QuoteDenom: quoteDenom,
	}
	testCases := []struct {
		name   string
		result bool
		setup  func()
	}{
		{
			"does not exist",
			false,
			func() {
			},
		},
		{
			"exists",
			true,
			func() {
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market1)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setup()
			res := suite.app.ClobKeeper.CheckPerpetualMarketAlreadyExists(suite.ctx, baseDenom, quoteDenom)
			suite.Equal(res, tc.result)
		})
	}
}
