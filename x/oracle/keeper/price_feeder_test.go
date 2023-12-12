package keeper_test

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/oracle/types"
)

func createNPriceFeeder(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.PriceFeeder {
	// Clear existing price feeders
	// Current it does have one price feeder from genesis.json
	priceFeeders := keeper.GetAllPriceFeeder(ctx)
	for _, p := range priceFeeders {
		keeper.RemovePriceFeeder(ctx, p.GetFeeder())
	}

	// Add new n price feeders
	items := make([]types.PriceFeeder, n)
	for i := range items {
		items[i].Feeder = strconv.Itoa(i)

		keeper.SetPriceFeeder(ctx, items[i])
	}
	return items
}

func (suite *KeeperTestSuite) TestPriceFeederGet() {
	keeper, ctx := suite.app.OracleKeeper, suite.ctx
	items := createNPriceFeeder(&keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetPriceFeeder(ctx, item.Feeder)
		suite.Require().True(found)
		suite.Require().Equal(
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}

func (suite *KeeperTestSuite) TestPriceFeederRemove() {
	keeper, ctx := suite.app.OracleKeeper, suite.ctx
	items := createNPriceFeeder(&keeper, ctx, 10)
	for _, item := range items {
		keeper.RemovePriceFeeder(ctx, item.Feeder)
		_, found := keeper.GetPriceFeeder(ctx, item.Feeder)
		suite.Require().False(found)
	}
}

func (suite *KeeperTestSuite) TestPriceFeederGetAll() {
	keeper, ctx := suite.app.OracleKeeper, suite.ctx
	items := createNPriceFeeder(&keeper, ctx, 10)
	suite.Require().ElementsMatch(
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllPriceFeeder(ctx)),
	)
}
