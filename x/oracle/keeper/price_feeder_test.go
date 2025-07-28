package keeper_test

/*
import (
	"strconv"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/oracle/keeper"
	"github.com/elys-network/elys/v7/x/oracle/types"
)

func createNPriceFeeder(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.PriceFeeder {
	// Clear existing price feeders
	// Current it does have one price feeder from genesis.json
	priceFeeders := keeper.GetAllPriceFeeder(ctx)
	for _, p := range priceFeeders {
		keeper.RemovePriceFeeder(ctx, sdk.MustAccAddressFromBech32(p.Feeder))
	}

	// Add new n price feeders
	items := make([]types.PriceFeeder, n)
	for i := range items {
		items[i].Feeder = authtypes.NewModuleAddress(strconv.Itoa(i)).String()

		keeper.SetPriceFeeder(ctx, items[i])
	}
	return items
}

func (suite *KeeperTestSuite) TestPriceFeederGet() {
	keeper, ctx := suite.app.OracleKeeper, suite.ctx
	items := createNPriceFeeder(&keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetPriceFeeder(ctx, sdk.MustAccAddressFromBech32(item.Feeder))
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
		keeper.RemovePriceFeeder(ctx, sdk.MustAccAddressFromBech32(item.Feeder))
		_, found := keeper.GetPriceFeeder(ctx, sdk.MustAccAddressFromBech32(item.Feeder))
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
*/
