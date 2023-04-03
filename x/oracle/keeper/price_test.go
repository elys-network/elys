package keeper_test

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/testutil/nullify"
	"github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/oracle/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNPrice(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Price {
	items := make([]types.Price, n)
	for i := range items {
		items[i].Provider = strconv.Itoa(i)
		items[i].Asset = "asset" + strconv.Itoa(i)
		items[i].Price = sdk.NewDec(1)
		items[i].Price = sdk.NewDec(1)
		items[i].Source = "binance"

		keeper.SetPrice(ctx, items[i])
	}
	return items
}

func (suite *KeeperTestSuite) TestPriceGet() {
	k, ctx := suite.app.OracleKeeper, suite.ctx
	items := createNPrice(&k, ctx, 10)
	for _, item := range items {
		rst, found := k.GetPrice(ctx, item.Asset)
		suite.Require().True(found)
		suite.Require().Equal(
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func (suite *KeeperTestSuite) TestPriceRemove() {
	k, ctx := suite.app.OracleKeeper, suite.ctx
	items := createNPrice(&k, ctx, 10)
	for _, item := range items {
		k.RemovePrice(ctx, item.Asset)
		_, found := k.GetPrice(ctx, item.Asset)
		suite.Require().False(found)
	}
}

func (suite *KeeperTestSuite) TestPriceGetAll() {
	k, ctx := suite.app.OracleKeeper, suite.ctx
	items := createNPrice(&k, ctx, 10)
	suite.Require().ElementsMatch(
		nullify.Fill(items),
		nullify.Fill(k.GetAllPrice(ctx)),
	)
}
