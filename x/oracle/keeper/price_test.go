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
		items[i].Timestamp = uint64(i)

		keeper.SetPrice(ctx, items[i])
	}
	return items
}

func (suite *KeeperTestSuite) TestPriceGet() {
	k, ctx := suite.app.OracleKeeper, suite.ctx
	items := createNPrice(&k, ctx, 10)
	for _, item := range items {
		rst, found := k.GetPrice(ctx, item.Asset, item.Source, item.Timestamp)
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
		k.RemovePrice(ctx, item.Asset, item.Source, item.Timestamp)
		_, found := k.GetPrice(ctx, item.Asset, item.Source, item.Timestamp)
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

func (suite *KeeperTestSuite) TestGetLatestPriceFromAssetAndSource() {
	prices := []types.Price{
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(1),
			Source:    "binance",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(2),
			Source:    "band",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(3),
			Source:    "band",
			Timestamp: 110000,
		},
	}
	for _, price := range prices {
		suite.app.OracleKeeper.SetPrice(suite.ctx, price)
	}
	price, found := suite.app.OracleKeeper.GetLatestPriceFromAssetAndSource(suite.ctx, "BTC", "binance")
	suite.Require().True(found)
	suite.Require().Equal(price, prices[0])
	price, found = suite.app.OracleKeeper.GetLatestPriceFromAssetAndSource(suite.ctx, "BTC", "osmosis")
	suite.Require().False(found)
	price, found = suite.app.OracleKeeper.GetLatestPriceFromAssetAndSource(suite.ctx, "BTC", "band")
	suite.Require().True(found)
	suite.Require().Equal(price, prices[2])
}

func (suite *KeeperTestSuite) TestGetLatestPriceFromAnySource() {
	prices := []types.Price{
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(1),
			Source:    "binance",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(2),
			Source:    "band",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(3),
			Source:    "band",
			Timestamp: 110000,
		},
	}
	for _, price := range prices {
		suite.app.OracleKeeper.SetPrice(suite.ctx, price)
	}
	price, found := suite.app.OracleKeeper.GetLatestPriceFromAnySource(suite.ctx, "BTC")
	suite.Require().True(found)
	suite.Require().Equal(price, prices[0])
}
