package keeper_test

import (
	"github.com/osmosis-labs/osmosis/osmomath"
	"strconv"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/testutil/nullify"
	"github.com/elys-network/elys/v7/x/oracle/keeper"
	"github.com/elys-network/elys/v7/x/oracle/types"
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

func (suite *KeeperTestSuite) TestGetAllAssetPrice() {
	priceData := []types.Price{
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(23),
			Source:    "elys",
			Timestamp: 20,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(12),
			Source:    "band",
			Timestamp: 10,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(76),
			Source:    "band",
			Timestamp: 40,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(55),
			Source:    "band",
			Timestamp: 30,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(89),
			Source:    "band",
			Timestamp: 25,
		},
	}
	for _, price := range priceData {
		suite.app.OracleKeeper.SetPrice(suite.ctx, price)
	}
	allPrices := suite.app.OracleKeeper.GetAllAssetPrice(suite.ctx, "BTC")
	suite.Require().Equal(len(priceData), len(allPrices))
}

func (suite *KeeperTestSuite) TestGetAssetPriceAndGetDenomPrice() {
	type Data struct {
		Asset   string
		Denom   string
		Price   sdkmath.LegacyDec
		Decimal uint64
	}
	prices := []Data{
		{
			Asset:   "BTC",
			Denom:   "sat",
			Price:   sdkmath.LegacyMustNewDecFromStr("110538.3425"),
			Decimal: 8,
		},
		{
			Asset:   "ATOM",
			Denom:   "uatom",
			Price:   sdkmath.LegacyMustNewDecFromStr("45.262"),
			Decimal: 6,
		},
		{
			Asset:   "ETH",
			Denom:   "wei",
			Price:   sdkmath.LegacyMustNewDecFromStr("50244.3179"),
			Decimal: 18,
		},
		{
			Asset:   "FET",
			Denom:   "afet",
			Price:   sdkmath.LegacyMustNewDecFromStr("1.2342"),
			Decimal: 18,
		},
		{
			Asset:   "MEME",
			Denom:   "ameme",
			Price:   sdkmath.LegacyMustNewDecFromStr("0.345232"),
			Decimal: 18,
		},
		{
			Asset:   "WORST_MEME",
			Denom:   "a_worst_meme",
			Price:   sdkmath.LegacyMustNewDecFromStr("0.003445232"),
			Decimal: 23,
		},
	}
	for _, price := range prices {
		suite.app.OracleKeeper.SetAssetInfo(suite.ctx, types.AssetInfo{
			Denom:      price.Denom,
			Display:    price.Asset,
			BandTicker: price.Asset,
			ElysTicker: price.Asset,
			Decimal:    price.Decimal,
		})
		suite.app.OracleKeeper.SetPrice(suite.ctx, types.Price{
			Asset:       price.Asset,
			Price:       price.Price,
			Source:      "elys",
			Provider:    "elys",
			Timestamp:   100000,
			BlockHeight: 100000,
		})
	}
	price, found := suite.app.OracleKeeper.GetAssetPrice(suite.ctx, prices[0].Asset)
	suite.Require().True(found)
	suite.Require().Equal(price, prices[0].Price)
	denomPrice := suite.app.OracleKeeper.GetDenomPrice(suite.ctx, prices[0].Denom)
	suite.Require().NotEqual(osmomath.ZeroBigDec(), denomPrice)
	suite.Require().NotEqual(price.String(), denomPrice.String())
	suite.Require().Equal("0.001105383425000000000000000000000000", denomPrice.String())

	price, found = suite.app.OracleKeeper.GetAssetPrice(suite.ctx, prices[1].Asset)
	suite.Require().True(found)
	suite.Require().Equal(price, prices[1].Price)
	denomPrice = suite.app.OracleKeeper.GetDenomPrice(suite.ctx, prices[1].Denom)
	suite.Require().NotEqual(osmomath.ZeroBigDec(), denomPrice)
	suite.Require().NotEqual(price.String(), denomPrice.String())
	suite.Require().Equal("0.000045262000000000000000000000000000", denomPrice.String())

	price, found = suite.app.OracleKeeper.GetAssetPrice(suite.ctx, prices[2].Asset)
	suite.Require().True(found)
	suite.Require().Equal(price, prices[2].Price)
	denomPrice = suite.app.OracleKeeper.GetDenomPrice(suite.ctx, prices[2].Denom)
	suite.Require().NotEqual(osmomath.ZeroBigDec(), denomPrice)
	suite.Require().NotEqual(price.String(), denomPrice.String())
	suite.Require().Equal("0.000000000000050244317900000000000000", denomPrice.String())

	price, found = suite.app.OracleKeeper.GetAssetPrice(suite.ctx, prices[3].Asset)
	suite.Require().True(found)
	suite.Require().Equal(price, prices[3].Price)
	denomPrice = suite.app.OracleKeeper.GetDenomPrice(suite.ctx, prices[3].Denom)
	suite.Require().NotEqual(osmomath.ZeroBigDec(), denomPrice)
	suite.Require().NotEqual(price.String(), denomPrice.String())
	suite.Require().Equal("0.000000000000000001234200000000000000", denomPrice.String())

	price, found = suite.app.OracleKeeper.GetAssetPrice(suite.ctx, prices[4].Asset)
	suite.Require().True(found)
	suite.Require().Equal(price, prices[4].Price)
	denomPrice = suite.app.OracleKeeper.GetDenomPrice(suite.ctx, prices[4].Denom)
	suite.Require().NotEqual(osmomath.ZeroBigDec(), denomPrice)
	suite.Require().NotEqual(price.String(), denomPrice.String())
	suite.Require().Equal("0.000000000000000000345232000000000000", denomPrice.String())

	price, found = suite.app.OracleKeeper.GetAssetPrice(suite.ctx, prices[5].Asset)
	suite.Require().True(found)
	suite.Require().Equal(price, prices[5].Price)
	denomPrice = suite.app.OracleKeeper.GetDenomPrice(suite.ctx, prices[5].Denom)
	suite.Require().NotEqual(osmomath.ZeroBigDec(), denomPrice)
	suite.Require().NotEqual(price.String(), denomPrice.String())
	suite.Require().Equal("0.000000000000000000000000034452320000", denomPrice.String())
}
