package keeper_test

/*
import (
	"strconv"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/testutil/nullify"
	"github.com/elys-network/elys/v6/x/oracle/keeper"
	"github.com/elys-network/elys/v6/x/oracle/types"
)

func createNPrice(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Price {
	items := make([]types.Price, n)
	for i := range items {
		items[i].Provider = strconv.Itoa(i)
		items[i].Asset = "asset" + strconv.Itoa(i)
		items[i].Price = sdkmath.LegacyNewDec(1)
		items[i].Price = sdkmath.LegacyNewDec(1)
		items[i].Source = "elys"
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
			Price:     sdkmath.LegacyNewDec(1),
			Source:    "elys",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(2),
			Source:    "band",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(3),
			Source:    "band",
			Timestamp: 110000,
		},
	}
	for _, price := range prices {
		suite.app.OracleKeeper.SetPrice(suite.ctx, price)
	}
	price, found := suite.app.OracleKeeper.GetLatestPriceFromAssetAndSource(suite.ctx, "BTC", "elys")
	suite.Require().True(found)
	suite.Require().Equal(price, prices[0])
	price, found = suite.app.OracleKeeper.GetLatestPriceFromAssetAndSource(suite.ctx, "BTC", "band")
	suite.Require().True(found)
	suite.Require().Equal(price, prices[2])
}

func (suite *KeeperTestSuite) TestGetLatestPriceFromAnySource() {
	prices := []types.Price{
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(1),
			Source:    "elys",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(2),
			Source:    "band",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(3),
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

// TODO: add test GetAssetPrice
// TODO: add test GetDenomPrice
*/
